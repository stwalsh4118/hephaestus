package templates

import (
	"encoding/json"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

func TestTranslator_AllFiveServiceTypes(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "d1",
		Name: "Full Stack",
		Nodes: []model.DiagramNode{
			{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "Database", Config: json.RawMessage(`{"type":"postgresql","engine":"PostgreSQL","version":"16"}`)},
			{ID: "redis", Type: model.ServiceTypeRedis, Name: "Cache", Config: json.RawMessage(`{"type":"redis","maxMemory":"256mb","evictionPolicy":"allkeys-lru"}`)},
			{ID: "rmq", Type: model.ServiceTypeRabbitMQ, Name: "Queue", Config: json.RawMessage(`{"type":"rabbitmq","vhost":"/"}`)},
			{ID: "api", Type: model.ServiceTypeAPIService, Name: "API"},
			{ID: "nginx", Type: model.ServiceTypeNginx, Name: "Gateway"},
		},
		Edges: []model.DiagramEdge{
			{ID: "e1", Source: "api", Target: "pg"},
			{ID: "e2", Source: "api", Target: "redis"},
			{ID: "e3", Source: "nginx", Target: "api"},
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(configs) != 5 {
		t.Fatalf("expected 5 configs, got %d", len(configs))
	}

	// Verify dependency ordering: infrastructure before application.
	cfgIndex := make(map[string]int)
	for i, c := range configs {
		cfgIndex[c.Name] = i
	}

	// pg, redis, rmq should come before api; api before nginx.
	for _, infra := range []string{"database", "cache", "queue"} {
		if cfgIndex[infra] > cfgIndex["api"] {
			t.Errorf("%s should start before api, got order: %v", infra, configNames(configs))
		}
	}
	if cfgIndex["api"] > cfgIndex["gateway"] {
		t.Errorf("api should start before gateway, got order: %v", configNames(configs))
	}

	// Verify all ports are unique.
	allPorts := make(map[string]bool)
	for _, cfg := range configs {
		for hostPort := range cfg.Ports {
			if allPorts[hostPort] {
				t.Errorf("duplicate host port: %s", hostPort)
			}
			allPorts[hostPort] = true
		}
	}

	// Verify correct images.
	imageByName := make(map[string]string)
	for _, c := range configs {
		imageByName[c.Name] = c.Image
	}
	expectedImages := map[string]string{
		"database": ImagePostgreSQL,
		"cache":    ImageRedis,
		"queue":    ImageRabbitMQ,
		"api":      ImageAPIService,
		"gateway":  ImageNginx,
	}
	for name, expected := range expectedImages {
		if imageByName[name] != expected {
			t.Errorf("expected image %q for %s, got %q", expected, name, imageByName[name])
		}
	}

	// Verify all configs have NetworkName set.
	for _, cfg := range configs {
		if cfg.NetworkName != docker.NetworkName {
			t.Errorf("config %q: expected network %q, got %q", cfg.Name, docker.NetworkName, cfg.NetworkName)
		}
	}
}

func TestTranslator_UnknownServiceType(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "d2",
		Name: "Bad",
		Nodes: []model.DiagramNode{
			{ID: "x", Type: "unknown-service", Name: "x"},
		},
	}

	_, err := tr.Translate(diagram)
	if err == nil {
		t.Fatal("expected error for unknown service type")
	}
}

func TestTranslator_EmptyDiagram(t *testing.T) {
	tr := NewTranslator()

	configs, err := tr.Translate(model.Diagram{ID: "d3", Name: "Empty"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if configs != nil {
		t.Errorf("expected nil for empty diagram, got %v", configs)
	}
}

func TestTranslator_SingleNode(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "d4",
		Name: "Single",
		Nodes: []model.DiagramNode{
			{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "db"},
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}
	if configs[0].Image != ImagePostgreSQL {
		t.Errorf("expected image %q, got %q", ImagePostgreSQL, configs[0].Image)
	}
}

func TestTranslator_PortAllocatorResetsPerCall(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "d5",
		Name: "Reset Test",
		Nodes: []model.DiagramNode{
			{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "db"},
		},
	}

	// Translate twice — both should succeed with the same starting port.
	configs1, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("first translate: %v", err)
	}
	configs2, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("second translate: %v", err)
	}

	// Port should be the same since allocator resets.
	port1 := firstHostPort(configs1[0])
	port2 := firstHostPort(configs2[0])
	if port1 != port2 {
		t.Errorf("expected same port across translates (allocator reset), got %s and %s", port1, port2)
	}
}

func TestTranslator_RabbitMQGetsTwoPorts(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "d6",
		Name: "RabbitMQ Ports",
		Nodes: []model.DiagramNode{
			{ID: "rmq", Type: model.ServiceTypeRabbitMQ, Name: "broker"},
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(configs[0].Ports) != 2 {
		t.Errorf("expected 2 port mappings for RabbitMQ, got %d", len(configs[0].Ports))
	}
}

func configNames(configs []docker.ContainerConfig) []string {
	names := make([]string, len(configs))
	for i, c := range configs {
		names[i] = c.Name
	}
	return names
}

func firstHostPort(cfg docker.ContainerConfig) string {
	for hp := range cfg.Ports {
		return hp
	}
	return ""
}
