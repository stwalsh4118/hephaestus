package templates

import (
	"errors"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

func TestResolveDependencies_LinearChain(t *testing.T) {
	// DB → API → Nginx (DB must start first)
	nodes := []model.DiagramNode{
		{ID: "nginx", Type: model.ServiceTypeNginx, Name: "nginx"},
		{ID: "api", Type: model.ServiceTypeAPIService, Name: "api"},
		{ID: "db", Type: model.ServiceTypePostgreSQL, Name: "db"},
	}
	edges := []model.DiagramEdge{
		{ID: "e1", Source: "api", Target: "db"},    // api depends on db
		{ID: "e2", Source: "nginx", Target: "api"},  // nginx depends on api
	}

	order, err := ResolveDependencies(nodes, edges)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// db must come before api, api before nginx.
	idx := indexMap(order)
	if idx["db"] > idx["api"] {
		t.Errorf("db should start before api: %v", order)
	}
	if idx["api"] > idx["nginx"] {
		t.Errorf("api should start before nginx: %v", order)
	}
}

func TestResolveDependencies_DiamondGraph(t *testing.T) {
	// DB ← API1 ← Nginx
	// DB ← API2 ← Nginx
	nodes := []model.DiagramNode{
		{ID: "db", Type: model.ServiceTypePostgreSQL, Name: "db"},
		{ID: "api1", Type: model.ServiceTypeAPIService, Name: "api1"},
		{ID: "api2", Type: model.ServiceTypeAPIService, Name: "api2"},
		{ID: "nginx", Type: model.ServiceTypeNginx, Name: "nginx"},
	}
	edges := []model.DiagramEdge{
		{ID: "e1", Source: "api1", Target: "db"},
		{ID: "e2", Source: "api2", Target: "db"},
		{ID: "e3", Source: "nginx", Target: "api1"},
		{ID: "e4", Source: "nginx", Target: "api2"},
	}

	order, err := ResolveDependencies(nodes, edges)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(order))
	}

	idx := indexMap(order)
	if idx["db"] > idx["api1"] || idx["db"] > idx["api2"] {
		t.Errorf("db should start before both APIs: %v", order)
	}
	if idx["api1"] > idx["nginx"] || idx["api2"] > idx["nginx"] {
		t.Errorf("both APIs should start before nginx: %v", order)
	}
}

func TestResolveDependencies_NoEdges_PriorityOrdering(t *testing.T) {
	nodes := []model.DiagramNode{
		{ID: "nginx", Type: model.ServiceTypeNginx, Name: "nginx"},
		{ID: "redis", Type: model.ServiceTypeRedis, Name: "redis"},
		{ID: "api", Type: model.ServiceTypeAPIService, Name: "api"},
		{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "pg"},
		{ID: "rmq", Type: model.ServiceTypeRabbitMQ, Name: "rmq"},
	}

	order, err := ResolveDependencies(nodes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 5 {
		t.Fatalf("expected 5 nodes, got %d", len(order))
	}

	// Infrastructure (pg, redis, rmq) should come before application (api, nginx).
	idx := indexMap(order)
	for _, infra := range []string{"pg", "redis", "rmq"} {
		for _, app := range []string{"api", "nginx"} {
			if idx[infra] > idx[app] {
				t.Errorf("%s (infrastructure) should start before %s (application): %v", infra, app, order)
			}
		}
	}
}

func TestResolveDependencies_SingleNode(t *testing.T) {
	nodes := []model.DiagramNode{
		{ID: "db", Type: model.ServiceTypePostgreSQL, Name: "db"},
	}

	order, err := ResolveDependencies(nodes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 1 || order[0] != "db" {
		t.Errorf("expected [db], got %v", order)
	}
}

func TestResolveDependencies_CycleDetection(t *testing.T) {
	nodes := []model.DiagramNode{
		{ID: "a", Type: model.ServiceTypeAPIService, Name: "a"},
		{ID: "b", Type: model.ServiceTypeAPIService, Name: "b"},
	}
	edges := []model.DiagramEdge{
		{ID: "e1", Source: "a", Target: "b"},
		{ID: "e2", Source: "b", Target: "a"},
	}

	_, err := ResolveDependencies(nodes, edges)
	if err == nil {
		t.Fatal("expected cycle error")
	}
	if !errors.Is(err, ErrCyclicDependency) {
		t.Errorf("expected ErrCyclicDependency, got %v", err)
	}
}

func TestResolveDependencies_EmptyInput(t *testing.T) {
	order, err := ResolveDependencies(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order != nil {
		t.Errorf("expected nil for empty input, got %v", order)
	}
}

func TestResolveDependencies_MixedEdgesAndPriority(t *testing.T) {
	// redis has no dependencies, api depends on redis, pg has no dependencies.
	// Expected: pg and redis first (infra, no deps), then api.
	nodes := []model.DiagramNode{
		{ID: "api", Type: model.ServiceTypeAPIService, Name: "api"},
		{ID: "redis", Type: model.ServiceTypeRedis, Name: "redis"},
		{ID: "pg", Type: model.ServiceTypePostgreSQL, Name: "pg"},
	}
	edges := []model.DiagramEdge{
		{ID: "e1", Source: "api", Target: "redis"},
	}

	order, err := ResolveDependencies(nodes, edges)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	idx := indexMap(order)
	// pg (infra, no deps) should be first or tied with redis.
	// redis (infra, no deps) should come before api.
	if idx["redis"] > idx["api"] {
		t.Errorf("redis should start before api: %v", order)
	}
	if idx["pg"] > idx["api"] {
		t.Errorf("pg should start before api: %v", order)
	}
}

// indexMap creates a map from node ID to its position in the order slice.
func indexMap(order []string) map[string]int {
	m := make(map[string]int, len(order))
	for i, id := range order {
		m[id] = i
	}
	return m
}
