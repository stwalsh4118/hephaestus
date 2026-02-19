package templates

import (
	"encoding/json"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// TestE2E_FullDiagramTranslation verifies AC1–AC6 with a realistic diagram
// containing all 5 service types and dependency edges.
func TestE2E_FullDiagramTranslation(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "e2e-full",
		Name: "Full Stack Architecture",
		Nodes: []model.DiagramNode{
			{
				ID:   "pg-1",
				Type: model.ServiceTypePostgreSQL,
				Name: "Primary DB",
				Config: json.RawMessage(`{"type":"postgresql","engine":"PostgreSQL","version":"16"}`),
			},
			{
				ID:   "redis-1",
				Type: model.ServiceTypeRedis,
				Name: "Session Cache",
				Config: json.RawMessage(`{"type":"redis","maxMemory":"512mb","evictionPolicy":"allkeys-lru"}`),
			},
			{
				ID:   "rmq-1",
				Type: model.ServiceTypeRabbitMQ,
				Name: "Event Bus",
				Config: json.RawMessage(`{"type":"rabbitmq","vhost":"/events"}`),
			},
			{
				ID:   "api-1",
				Type: model.ServiceTypeAPIService,
				Name: "User API",
				Config: json.RawMessage(`{"type":"api-service","endpoints":[{"method":"GET","path":"/users","responseSchema":"{}"}],"port":4010}`),
			},
			{
				ID:   "nginx-1",
				Type: model.ServiceTypeNginx,
				Name: "API Gateway",
				Config: json.RawMessage(`{"type":"nginx","upstreamServers":["user-api"]}`),
			},
		},
		Edges: []model.DiagramEdge{
			{ID: "e1", Source: "api-1", Target: "pg-1"},    // API depends on PostgreSQL
			{ID: "e2", Source: "api-1", Target: "redis-1"},  // API depends on Redis
			{ID: "e3", Source: "api-1", Target: "rmq-1"},    // API depends on RabbitMQ
			{ID: "e4", Source: "nginx-1", Target: "api-1"},  // Nginx depends on API
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	// AC1: Each of the 5 service types has a container template that produces valid Docker configurations.
	t.Run("AC1_AllFiveTypesProduceConfigs", func(t *testing.T) {
		if len(configs) != 5 {
			t.Fatalf("expected 5 configs, got %d", len(configs))
		}
	})

	// AC2: Templates use correct official Docker images for each service type.
	t.Run("AC2_CorrectImages", func(t *testing.T) {
		imageByName := make(map[string]string)
		for _, c := range configs {
			imageByName[c.Name] = c.Image
		}

		expected := map[string]string{
			"primary-db":    ImagePostgreSQL,
			"session-cache": ImageRedis,
			"event-bus":     ImageRabbitMQ,
			"user-api":      ImageAPIService,
			"api-gateway":   ImageNginx,
		}
		for name, img := range expected {
			if imageByName[name] != img {
				t.Errorf("node %q: expected image %q, got %q", name, img, imageByName[name])
			}
		}
	})

	// AC3: Service-specific configuration (env vars, volumes, ports) is correctly applied.
	t.Run("AC3_ServiceSpecificConfig", func(t *testing.T) {
		cfgByName := make(map[string]docker.ContainerConfig)
		for _, c := range configs {
			cfgByName[c.Name] = c
		}

		// PostgreSQL env vars.
		pg := cfgByName["primary-db"]
		if pg.Env["POSTGRES_USER"] != "hephaestus" {
			t.Errorf("PostgreSQL: expected POSTGRES_USER=hephaestus, got %q", pg.Env["POSTGRES_USER"])
		}
		if pg.Env["POSTGRES_PASSWORD"] != "hephaestus" {
			t.Errorf("PostgreSQL: expected POSTGRES_PASSWORD=hephaestus, got %q", pg.Env["POSTGRES_PASSWORD"])
		}
		if pg.Env["POSTGRES_DB"] != "hephaestus" {
			t.Errorf("PostgreSQL: expected POSTGRES_DB=hephaestus, got %q", pg.Env["POSTGRES_DB"])
		}

		// Redis config from node.
		redis := cfgByName["session-cache"]
		if redis.Env["REDIS_MAXMEMORY"] != "512mb" {
			t.Errorf("Redis: expected REDIS_MAXMEMORY=512mb, got %q", redis.Env["REDIS_MAXMEMORY"])
		}
		if redis.Env["REDIS_EVICTION_POLICY"] != "allkeys-lru" {
			t.Errorf("Redis: expected REDIS_EVICTION_POLICY=allkeys-lru, got %q", redis.Env["REDIS_EVICTION_POLICY"])
		}

		// RabbitMQ vhost.
		rmq := cfgByName["event-bus"]
		if rmq.Env["RABBITMQ_DEFAULT_VHOST"] != "/events" {
			t.Errorf("RabbitMQ: expected RABBITMQ_DEFAULT_VHOST=/events, got %q", rmq.Env["RABBITMQ_DEFAULT_VHOST"])
		}

		// Nginx upstreams.
		nginx := cfgByName["api-gateway"]
		if nginx.Env["NGINX_UPSTREAMS"] != "user-api" {
			t.Errorf("Nginx: expected NGINX_UPSTREAMS=user-api, got %q", nginx.Env["NGINX_UPSTREAMS"])
		}

		// All configs have correct hostname.
		for _, c := range configs {
			if c.Hostname == "" {
				t.Errorf("config %q has empty hostname", c.Name)
			}
			if c.Hostname != c.Name {
				t.Errorf("config %q: hostname %q should match name", c.Name, c.Hostname)
			}
		}
	})

	// AC4: Dependency ordering ensures databases and caches start before services that depend on them.
	t.Run("AC4_DependencyOrdering", func(t *testing.T) {
		idx := make(map[string]int)
		for i, c := range configs {
			idx[c.Name] = i
		}

		// Infrastructure before application.
		infra := []string{"primary-db", "session-cache", "event-bus"}
		app := []string{"user-api", "api-gateway"}

		for _, i := range infra {
			for _, a := range app {
				if idx[i] > idx[a] {
					t.Errorf("%s (infrastructure) should start before %s (application), got positions %d and %d",
						i, a, idx[i], idx[a])
				}
			}
		}

		// Specific edge dependencies.
		if idx["primary-db"] > idx["user-api"] {
			t.Error("primary-db should start before user-api")
		}
		if idx["session-cache"] > idx["user-api"] {
			t.Error("session-cache should start before user-api")
		}
		if idx["event-bus"] > idx["user-api"] {
			t.Error("event-bus should start before user-api")
		}
		if idx["user-api"] > idx["api-gateway"] {
			t.Error("user-api should start before api-gateway")
		}
	})

	// AC5: Port allocation assigns unique host ports without conflicts.
	t.Run("AC5_UniquePortAllocation", func(t *testing.T) {
		allHostPorts := make(map[string]string) // host port → service name
		for _, c := range configs {
			for hp := range c.Ports {
				if existing, ok := allHostPorts[hp]; ok {
					t.Errorf("port conflict: host port %s used by both %q and %q", hp, existing, c.Name)
				}
				allHostPorts[hp] = c.Name
			}
		}

		// RabbitMQ should have 2 ports.
		rmqPorts := 0
		for _, c := range configs {
			if c.Image == ImageRabbitMQ {
				rmqPorts = len(c.Ports)
			}
		}
		if rmqPorts != 2 {
			t.Errorf("RabbitMQ should have 2 port mappings, got %d", rmqPorts)
		}

		// Total unique ports: 4 single-port services + 1 dual-port = 6.
		if len(allHostPorts) != 6 {
			t.Errorf("expected 6 total port mappings, got %d", len(allHostPorts))
		}
	})

	// AC6: A diagram with all 5 service types can be translated into container configs
	// that the orchestration engine accepts.
	t.Run("AC6_OrchestratorCompatibility", func(t *testing.T) {
		for _, c := range configs {
			if c.Image == "" {
				t.Errorf("config %q: Image must not be empty", c.Name)
			}
			if c.Name == "" {
				t.Error("config has empty Name")
			}
			if len(c.Ports) == 0 {
				t.Errorf("config %q: must have at least one port mapping", c.Name)
			}
			if c.NetworkName != docker.NetworkName {
				t.Errorf("config %q: expected NetworkName %q, got %q", c.Name, docker.NetworkName, c.NetworkName)
			}
			// Verify port mappings have valid container-side ports.
			for hp, cp := range c.Ports {
				if hp == "" {
					t.Errorf("config %q: empty host port", c.Name)
				}
				if cp == "" {
					t.Errorf("config %q: empty container port for host port %s", c.Name, hp)
				}
			}
		}
	})
}

// TestE2E_ServiceSpecificConfigAccuracy verifies AC3 with populated config JSON.
func TestE2E_ServiceSpecificConfigAccuracy(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "e2e-config",
		Name: "Config Test",
		Nodes: []model.DiagramNode{
			{
				ID:     "pg",
				Type:   model.ServiceTypePostgreSQL,
				Name:   "DB",
				Config: json.RawMessage(`{"type":"postgresql","engine":"PostgreSQL","version":"16"}`),
			},
			{
				ID:     "redis",
				Type:   model.ServiceTypeRedis,
				Name:   "Cache",
				Config: json.RawMessage(`{"type":"redis","maxMemory":"1gb","evictionPolicy":"volatile-lru"}`),
			},
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	for _, c := range configs {
		switch c.Image {
		case ImagePostgreSQL:
			if c.Env["POSTGRES_USER"] == "" || c.Env["POSTGRES_PASSWORD"] == "" || c.Env["POSTGRES_DB"] == "" {
				t.Error("PostgreSQL missing required env vars")
			}
		case ImageRedis:
			if c.Env["REDIS_MAXMEMORY"] != "1gb" {
				t.Errorf("Redis maxMemory: expected 1gb, got %q", c.Env["REDIS_MAXMEMORY"])
			}
			if c.Env["REDIS_EVICTION_POLICY"] != "volatile-lru" {
				t.Errorf("Redis evictionPolicy: expected volatile-lru, got %q", c.Env["REDIS_EVICTION_POLICY"])
			}
		}
	}
}

// TestE2E_DependencyOrderingCorrectness verifies AC4 with a clear dependency chain.
func TestE2E_DependencyOrderingCorrectness(t *testing.T) {
	tr := NewTranslator()

	// Chain: nginx → api → pg (pg starts first, then api, then nginx).
	diagram := model.Diagram{
		ID:   "e2e-order",
		Name: "Order Test",
		Nodes: []model.DiagramNode{
			{ID: "nginx", Type: model.ServiceTypeNginx, Name: "LB"},
			{ID: "api", Type: model.ServiceTypeAPIService, Name: "App"},
			{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "Store"},
		},
		Edges: []model.DiagramEdge{
			{ID: "e1", Source: "api", Target: "pg"},
			{ID: "e2", Source: "nginx", Target: "api"},
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	idx := make(map[string]int)
	for i, c := range configs {
		idx[c.Name] = i
	}

	if idx["store"] > idx["app"] {
		t.Error("store (PostgreSQL) should start before app (API)")
	}
	if idx["app"] > idx["lb"] {
		t.Error("app (API) should start before lb (Nginx)")
	}
}

// TestE2E_PortAllocationUniqueness verifies AC5 exhaustively.
func TestE2E_PortAllocationUniqueness(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "e2e-ports",
		Name: "Port Test",
		Nodes: []model.DiagramNode{
			{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "db1"},
			{ID: "redis", Type: model.ServiceTypeRedis, Name: "cache1"},
			{ID: "rmq", Type: model.ServiceTypeRabbitMQ, Name: "mq1"},
			{ID: "api", Type: model.ServiceTypeAPIService, Name: "svc1"},
			{ID: "nginx", Type: model.ServiceTypeNginx, Name: "lb1"},
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	allPorts := make(map[string]bool)
	for _, c := range configs {
		for hp := range c.Ports {
			if allPorts[hp] {
				t.Fatalf("duplicate host port: %s", hp)
			}
			allPorts[hp] = true
		}
	}

	// 4 single-port + 1 dual-port (RabbitMQ) = 6 total.
	if len(allPorts) != 6 {
		t.Errorf("expected 6 unique host ports, got %d", len(allPorts))
	}
}

// TestE2E_OrchestratorContractCompliance verifies AC6 — each config has all
// required fields for docker.Orchestrator.CreateContainer.
func TestE2E_OrchestratorContractCompliance(t *testing.T) {
	tr := NewTranslator()

	diagram := model.Diagram{
		ID:   "e2e-contract",
		Name: "Contract Test",
		Nodes: []model.DiagramNode{
			{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "db"},
			{ID: "redis", Type: model.ServiceTypeRedis, Name: "cache"},
			{ID: "rmq", Type: model.ServiceTypeRabbitMQ, Name: "mq"},
			{ID: "api", Type: model.ServiceTypeAPIService, Name: "api"},
			{ID: "nginx", Type: model.ServiceTypeNginx, Name: "nginx"},
		},
	}

	configs, err := tr.Translate(diagram)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	for _, c := range configs {
		// Required for CreateContainer.
		if c.Image == "" {
			t.Errorf("%q: Image is required", c.Name)
		}
		if c.Name == "" {
			t.Error("Name is required")
		}
		if c.NetworkName == "" {
			t.Errorf("%q: NetworkName is required", c.Name)
		}
		if len(c.Ports) == 0 {
			t.Errorf("%q: at least one port mapping is required", c.Name)
		}

		// Verify container-side ports are well-known ports for each service.
		for _, cp := range c.Ports {
			validPorts := map[string]bool{
				PortAPIService:         true,
				PortPostgreSQL:         true,
				PortRedis:              true,
				PortNginx:              true,
				PortRabbitMQAMQP:       true,
				PortRabbitMQManagement: true,
			}
			if !validPorts[cp] {
				t.Errorf("%q: unexpected container port %q", c.Name, cp)
			}
		}
	}
}
