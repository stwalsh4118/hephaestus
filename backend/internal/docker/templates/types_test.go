package templates

import (
	"strconv"
	"testing"
)

func TestImageConstants_AreNonEmpty(t *testing.T) {
	images := map[string]string{
		"ImageAPIService": ImageAPIService,
		"ImagePostgreSQL": ImagePostgreSQL,
		"ImageRedis":      ImageRedis,
		"ImageNginx":      ImageNginx,
		"ImageRabbitMQ":   ImageRabbitMQ,
	}
	for name, val := range images {
		if val == "" {
			t.Errorf("%s must not be empty", name)
		}
	}
}

func TestPortConstants_AreValidPorts(t *testing.T) {
	ports := map[string]string{
		"PortAPIService":         PortAPIService,
		"PortPostgreSQL":         PortPostgreSQL,
		"PortRedis":              PortRedis,
		"PortNginx":              PortNginx,
		"PortRabbitMQAMQP":       PortRabbitMQAMQP,
		"PortRabbitMQManagement": PortRabbitMQManagement,
	}
	for name, val := range ports {
		n, err := strconv.Atoi(val)
		if err != nil {
			t.Errorf("%s: %q is not a valid integer: %v", name, val, err)
			continue
		}
		if n < 1 || n > 65535 {
			t.Errorf("%s: port %d is out of valid range 1-65535", name, n)
		}
	}
}

func TestNewRegistry_ReturnsPopulatedRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if len(reg) != 5 {
		t.Errorf("expected 5 entries in registry, got %d", len(reg))
	}
}

func TestDefaultPostgresEnv_HasRequiredKeys(t *testing.T) {
	env := DefaultPostgresEnv()
	required := []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"}
	for _, key := range required {
		val, ok := env[key]
		if !ok {
			t.Errorf("DefaultPostgresEnv missing required key %q", key)
		} else if val == "" {
			t.Errorf("DefaultPostgresEnv[%q] must not be empty", key)
		}
	}
}

func TestDefaultPostgresEnv_ReturnsFreshCopy(t *testing.T) {
	env1 := DefaultPostgresEnv()
	env2 := DefaultPostgresEnv()
	env1["POSTGRES_USER"] = "mutated"
	if env2["POSTGRES_USER"] == "mutated" {
		t.Error("DefaultPostgresEnv should return independent copies")
	}
}
