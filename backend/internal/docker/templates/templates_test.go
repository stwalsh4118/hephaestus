package templates

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

func TestPostgreSQLTemplate_Build(t *testing.T) {
	tmpl := &PostgreSQLTemplate{}
	node := model.DiagramNode{
		ID:   "pg-1",
		Type: model.ServiceTypePostgreSQL,
		Name: "My Database",
		Config: json.RawMessage(`{"type":"postgresql","engine":"PostgreSQL","version":"16"}`),
	}

	cfg, err := tmpl.Build(node, "15432")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Image != ImagePostgreSQL {
		t.Errorf("expected image %q, got %q", ImagePostgreSQL, cfg.Image)
	}
	if cfg.Hostname != "my-database" {
		t.Errorf("expected hostname %q, got %q", "my-database", cfg.Hostname)
	}
	if cfg.Ports["15432"] != PortPostgreSQL {
		t.Errorf("expected port mapping 15432→%s, got %q", PortPostgreSQL, cfg.Ports["15432"])
	}
	if cfg.Env["POSTGRES_USER"] != "hephaestus" {
		t.Errorf("expected POSTGRES_USER=hephaestus, got %q", cfg.Env["POSTGRES_USER"])
	}
	if cfg.Env["POSTGRES_PASSWORD"] != "hephaestus" {
		t.Errorf("expected POSTGRES_PASSWORD=hephaestus, got %q", cfg.Env["POSTGRES_PASSWORD"])
	}
	if cfg.Env["POSTGRES_DB"] != "hephaestus" {
		t.Errorf("expected POSTGRES_DB=hephaestus, got %q", cfg.Env["POSTGRES_DB"])
	}
	if cfg.NetworkName != docker.NetworkName {
		t.Errorf("expected network %q, got %q", docker.NetworkName, cfg.NetworkName)
	}
}

func TestPostgreSQLTemplate_Build_EmptyConfig(t *testing.T) {
	tmpl := &PostgreSQLTemplate{}
	node := model.DiagramNode{
		ID:   "pg-2",
		Type: model.ServiceTypePostgreSQL,
		Name: "postgres",
	}

	cfg, err := tmpl.Build(node, "15432")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Image != ImagePostgreSQL {
		t.Errorf("expected image %q, got %q", ImagePostgreSQL, cfg.Image)
	}
}

func TestRedisTemplate_Build(t *testing.T) {
	tmpl := &RedisTemplate{}
	node := model.DiagramNode{
		ID:   "redis-1",
		Type: model.ServiceTypeRedis,
		Name: "Redis Cache",
		Config: json.RawMessage(`{"type":"redis","maxMemory":"256mb","evictionPolicy":"allkeys-lru"}`),
	}

	cfg, err := tmpl.Build(node, "16379")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Image != ImageRedis {
		t.Errorf("expected image %q, got %q", ImageRedis, cfg.Image)
	}
	if cfg.Hostname != "redis-cache" {
		t.Errorf("expected hostname %q, got %q", "redis-cache", cfg.Hostname)
	}
	if cfg.Ports["16379"] != PortRedis {
		t.Errorf("expected port mapping 16379→%s, got %q", PortRedis, cfg.Ports["16379"])
	}
	if cfg.Env["REDIS_MAXMEMORY"] != "256mb" {
		t.Errorf("expected REDIS_MAXMEMORY=256mb, got %q", cfg.Env["REDIS_MAXMEMORY"])
	}
	if cfg.Env["REDIS_EVICTION_POLICY"] != "allkeys-lru" {
		t.Errorf("expected REDIS_EVICTION_POLICY=allkeys-lru, got %q", cfg.Env["REDIS_EVICTION_POLICY"])
	}
}

func TestRedisTemplate_Build_EmptyConfig(t *testing.T) {
	tmpl := &RedisTemplate{}
	node := model.DiagramNode{
		ID:   "redis-2",
		Type: model.ServiceTypeRedis,
		Name: "redis",
	}

	cfg, err := tmpl.Build(node, "16379")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Env) != 0 {
		t.Errorf("expected empty env for unconfigured redis, got %v", cfg.Env)
	}
}

func TestNginxTemplate_Build(t *testing.T) {
	tmpl := &NginxTemplate{}
	node := model.DiagramNode{
		ID:   "nginx-1",
		Type: model.ServiceTypeNginx,
		Name: "Load Balancer",
		Config: json.RawMessage(`{"type":"nginx","upstreamServers":["api-1","api-2"]}`),
	}

	cfg, err := tmpl.Build(node, "18080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Image != ImageNginx {
		t.Errorf("expected image %q, got %q", ImageNginx, cfg.Image)
	}
	if cfg.Hostname != "load-balancer" {
		t.Errorf("expected hostname %q, got %q", "load-balancer", cfg.Hostname)
	}
	if cfg.Ports["18080"] != PortNginx {
		t.Errorf("expected port mapping 18080→%s, got %q", PortNginx, cfg.Ports["18080"])
	}
	if cfg.Env["NGINX_UPSTREAMS"] != "api-1,api-2" {
		t.Errorf("expected NGINX_UPSTREAMS=api-1,api-2, got %q", cfg.Env["NGINX_UPSTREAMS"])
	}
}

func TestNginxTemplate_Build_EmptyConfig(t *testing.T) {
	tmpl := &NginxTemplate{}
	node := model.DiagramNode{
		ID:   "nginx-2",
		Type: model.ServiceTypeNginx,
		Name: "nginx",
	}

	cfg, err := tmpl.Build(node, "18080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Env["NGINX_UPSTREAMS"]; ok {
		t.Error("expected no NGINX_UPSTREAMS for unconfigured nginx")
	}
}

func TestRabbitMQTemplate_Build(t *testing.T) {
	tmpl := &RabbitMQTemplate{}
	node := model.DiagramNode{
		ID:   "rmq-1",
		Type: model.ServiceTypeRabbitMQ,
		Name: "Message Broker",
		Config: json.RawMessage(`{"type":"rabbitmq","vhost":"/events"}`),
	}

	cfg, err := tmpl.Build(node, "15672", "25672")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Image != ImageRabbitMQ {
		t.Errorf("expected image %q, got %q", ImageRabbitMQ, cfg.Image)
	}
	if cfg.Hostname != "message-broker" {
		t.Errorf("expected hostname %q, got %q", "message-broker", cfg.Hostname)
	}
	if cfg.Ports["15672"] != PortRabbitMQAMQP {
		t.Errorf("expected port mapping 15672→%s, got %q", PortRabbitMQAMQP, cfg.Ports["15672"])
	}
	if cfg.Ports["25672"] != PortRabbitMQManagement {
		t.Errorf("expected port mapping 25672→%s, got %q", PortRabbitMQManagement, cfg.Ports["25672"])
	}
	if cfg.Env["RABBITMQ_DEFAULT_VHOST"] != "/events" {
		t.Errorf("expected RABBITMQ_DEFAULT_VHOST=/events, got %q", cfg.Env["RABBITMQ_DEFAULT_VHOST"])
	}
}

func TestRabbitMQTemplate_Build_DefaultVhost(t *testing.T) {
	tmpl := &RabbitMQTemplate{}
	node := model.DiagramNode{
		ID:   "rmq-2",
		Type: model.ServiceTypeRabbitMQ,
		Name: "rabbitmq",
	}

	cfg, err := tmpl.Build(node, "15672")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Env["RABBITMQ_DEFAULT_VHOST"] != "/" {
		t.Errorf("expected default vhost /, got %q", cfg.Env["RABBITMQ_DEFAULT_VHOST"])
	}
}

func TestAPIServiceTemplate_Build_NoConfig(t *testing.T) {
	tmpl := &APIServiceTemplate{}
	node := model.DiagramNode{
		ID:   "api-1",
		Type: model.ServiceTypeAPIService,
		Name: "User API",
	}

	cfg, err := tmpl.Build(node, "14010")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Image != ImageAPIService {
		t.Errorf("expected image %q, got %q", ImageAPIService, cfg.Image)
	}
	if cfg.Hostname != "user-api" {
		t.Errorf("expected hostname %q, got %q", "user-api", cfg.Hostname)
	}
	if cfg.Ports["14010"] != PortAPIService {
		t.Errorf("expected port mapping 14010→%s, got %q", PortAPIService, cfg.Ports["14010"])
	}
	if cfg.NetworkName != docker.NetworkName {
		t.Errorf("expected network %q, got %q", docker.NetworkName, cfg.NetworkName)
	}
	// Verify Prism Cmd is set.
	expectedCmd := []string{"mock", "-h", "0.0.0.0", "/tmp/spec.json"}
	if len(cfg.Cmd) != len(expectedCmd) {
		t.Fatalf("expected Cmd length %d, got %d", len(expectedCmd), len(cfg.Cmd))
	}
	for i, arg := range expectedCmd {
		if cfg.Cmd[i] != arg {
			t.Errorf("Cmd[%d] = %q, want %q", i, cfg.Cmd[i], arg)
		}
	}
	// Verify a spec file volume is mounted.
	if len(cfg.Volumes) == 0 {
		t.Fatal("expected at least one volume mount for spec file")
	}
	for _, containerPath := range cfg.Volumes {
		if containerPath != "/tmp/spec.json" {
			t.Errorf("expected container path /tmp/spec.json, got %q", containerPath)
		}
	}

	// Cleanup spec file.
	for hostPath := range cfg.Volumes {
		_ = os.Remove(hostPath)
	}
}

func TestAPIServiceTemplate_Build_WithEndpoints(t *testing.T) {
	tmpl := &APIServiceTemplate{}
	node := model.DiagramNode{
		ID:   "api-2",
		Type: model.ServiceTypeAPIService,
		Name: "Order Service",
		Config: json.RawMessage(`{
			"type": "api-service",
			"endpoints": [
				{"method": "GET", "path": "/orders", "responseSchema": "{\"type\":\"array\"}"},
				{"method": "POST", "path": "/orders", "responseSchema": "{\"type\":\"object\"}"}
			],
			"port": 4010
		}`),
	}

	cfg, err := tmpl.Build(node, "14010")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Cmd.
	if len(cfg.Cmd) != 4 || cfg.Cmd[0] != "mock" {
		t.Errorf("expected Prism mock command, got %v", cfg.Cmd)
	}

	// Verify volume mount.
	if len(cfg.Volumes) != 1 {
		t.Fatalf("expected 1 volume mount, got %d", len(cfg.Volumes))
	}
	var hostSpecPath string
	for hp, cp := range cfg.Volumes {
		hostSpecPath = hp
		if cp != "/tmp/spec.json" {
			t.Errorf("expected container path /tmp/spec.json, got %q", cp)
		}
	}

	// Verify spec file was written and is valid JSON.
	specData, err := os.ReadFile(hostSpecPath)
	if err != nil {
		t.Fatalf("failed to read spec file %q: %v", hostSpecPath, err)
	}
	if !json.Valid(specData) {
		t.Error("spec file is not valid JSON")
	}

	// Verify spec contains our endpoints.
	var spec map[string]any
	if err := json.Unmarshal(specData, &spec); err != nil {
		t.Fatalf("failed to parse spec: %v", err)
	}
	if spec["openapi"] != "3.0.0" {
		t.Errorf("expected openapi 3.0.0, got %v", spec["openapi"])
	}
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatal("expected paths object in spec")
	}
	if _, ok := paths["/orders"]; !ok {
		t.Error("expected /orders path in spec")
	}

	// Cleanup.
	_ = os.Remove(hostSpecPath)
}

func TestAPIServiceTemplate_Build_EmptyEndpoints(t *testing.T) {
	tmpl := &APIServiceTemplate{}
	node := model.DiagramNode{
		ID:   "api-3",
		Type: model.ServiceTypeAPIService,
		Name: "Empty API",
		Config: json.RawMessage(`{
			"type": "api-service",
			"endpoints": [],
			"port": 4010
		}`),
	}

	cfg, err := tmpl.Build(node, "14010")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still produce a valid config with Cmd and Volumes.
	if len(cfg.Cmd) == 0 {
		t.Error("expected Prism Cmd even with empty endpoints")
	}
	if len(cfg.Volumes) == 0 {
		t.Error("expected volume mount even with empty endpoints")
	}

	// Verify spec is valid JSON with empty paths.
	for hostPath := range cfg.Volumes {
		specData, err := os.ReadFile(hostPath)
		if err != nil {
			t.Fatalf("failed to read spec file: %v", err)
		}
		var spec map[string]any
		if err := json.Unmarshal(specData, &spec); err != nil {
			t.Fatalf("failed to parse spec: %v", err)
		}
		paths, ok := spec["paths"].(map[string]any)
		if !ok {
			t.Fatal("expected paths object")
		}
		if len(paths) != 0 {
			t.Errorf("expected 0 paths for empty endpoints, got %d", len(paths))
		}
		_ = os.Remove(hostPath)
	}
}

func TestAPIServiceTemplate_Build_InvalidConfig(t *testing.T) {
	tmpl := &APIServiceTemplate{}
	node := model.DiagramNode{
		ID:     "api-bad",
		Type:   model.ServiceTypeAPIService,
		Name:   "bad-api",
		Config: json.RawMessage(`{invalid json`),
	}

	_, err := tmpl.Build(node, "14010")
	if err == nil {
		t.Fatal("expected error for invalid config JSON")
	}
}

func TestNewRegistry_ContainsAll5ServiceTypes(t *testing.T) {
	reg := NewRegistry()

	expected := []string{
		model.ServiceTypeAPIService,
		model.ServiceTypePostgreSQL,
		model.ServiceTypeRedis,
		model.ServiceTypeNginx,
		model.ServiceTypeRabbitMQ,
	}

	if len(reg) != len(expected) {
		t.Fatalf("expected %d templates, got %d", len(expected), len(reg))
	}

	for _, svcType := range expected {
		if _, ok := reg[svcType]; !ok {
			t.Errorf("registry missing template for %q", svcType)
		}
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My Database", "my-database"},
		{"Redis Cache", "redis-cache"},
		{"nginx", "nginx"},
		{"API Service!", "api-service"},
		{"Test  Name@#$", "test--name"},
		{"UPPERCASE", "uppercase"},
	}

	for _, tc := range tests {
		got := sanitizeName(tc.input)
		if got != tc.expected {
			t.Errorf("sanitizeName(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestPostgreSQLTemplate_Build_InvalidConfig(t *testing.T) {
	tmpl := &PostgreSQLTemplate{}
	node := model.DiagramNode{
		ID:     "pg-bad",
		Type:   model.ServiceTypePostgreSQL,
		Name:   "postgres",
		Config: json.RawMessage(`{invalid json`),
	}

	_, err := tmpl.Build(node, "15432")
	if err == nil {
		t.Fatal("expected error for invalid config JSON")
	}
}

func TestRedisTemplate_Build_InvalidConfig(t *testing.T) {
	tmpl := &RedisTemplate{}
	node := model.DiagramNode{
		ID:     "redis-bad",
		Type:   model.ServiceTypeRedis,
		Name:   "redis",
		Config: json.RawMessage(`{invalid json`),
	}

	_, err := tmpl.Build(node, "16379")
	if err == nil {
		t.Fatal("expected error for invalid config JSON")
	}
}

func TestNginxTemplate_Build_InvalidConfig(t *testing.T) {
	tmpl := &NginxTemplate{}
	node := model.DiagramNode{
		ID:     "nginx-bad",
		Type:   model.ServiceTypeNginx,
		Name:   "nginx",
		Config: json.RawMessage(`{invalid json`),
	}

	_, err := tmpl.Build(node, "18080")
	if err == nil {
		t.Fatal("expected error for invalid config JSON")
	}
}

func TestRabbitMQTemplate_Build_InvalidConfig(t *testing.T) {
	tmpl := &RabbitMQTemplate{}
	node := model.DiagramNode{
		ID:     "rmq-bad",
		Type:   model.ServiceTypeRabbitMQ,
		Name:   "rabbitmq",
		Config: json.RawMessage(`{invalid json`),
	}

	_, err := tmpl.Build(node, "15672")
	if err == nil {
		t.Fatal("expected error for invalid config JSON")
	}
}
