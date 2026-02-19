package templates

import (
	"encoding/json"
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

func TestAPIServiceTemplate_Build(t *testing.T) {
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
