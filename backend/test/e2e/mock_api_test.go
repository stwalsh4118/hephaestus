package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/docker/templates"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
	"github.com/stwalsh4118/hephaestus/backend/internal/openapi"
)

// waitForPrism polls the given URL until it responds with HTTP 200 or the timeout expires.
func waitForPrism(t *testing.T, url string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("Prism did not become ready at %s within %v", url, timeout)
}

// prismReadyTimeout is how long to wait for a Prism container to start serving.
const prismReadyTimeout = 30 * time.Second

// TestPBI8_AC1_OpenAPISpecGeneration verifies that endpoint definitions are
// translated into valid OpenAPI 3.0 specs.
func TestPBI8_AC1_OpenAPISpecGeneration(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/users", ResponseSchema: `{"type":"array","items":{"type":"object","properties":{"id":{"type":"integer"},"name":{"type":"string"}}}}`},
		{Method: "POST", Path: "/users", ResponseSchema: `{"type":"object","properties":{"id":{"type":"integer"},"name":{"type":"string"}}}`},
		{Method: "GET", Path: "/users/{id}", ResponseSchema: `{"type":"object","properties":{"id":{"type":"integer"},"name":{"type":"string"}}}`},
	}

	specBytes, err := openapi.GenerateSpec(endpoints, "User Service")
	if err != nil {
		t.Fatalf("GenerateSpec failed: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(specBytes, &doc); err != nil {
		t.Fatalf("spec is not valid JSON: %v", err)
	}

	// Verify OpenAPI version.
	if doc["openapi"] != "3.0.0" {
		t.Errorf("expected openapi 3.0.0, got %v", doc["openapi"])
	}

	// Verify paths match input.
	paths, ok := doc["paths"].(map[string]any)
	if !ok {
		t.Fatal("expected paths object in spec")
	}
	expectedPaths := []string{"/users", "/users/{id}"}
	for _, p := range expectedPaths {
		if _, ok := paths[p]; !ok {
			t.Errorf("expected path %q in spec", p)
		}
	}

	// Verify /users has both GET and POST.
	usersPath, ok := paths["/users"].(map[string]any)
	if !ok {
		t.Fatal("expected /users to be an object")
	}
	if _, ok := usersPath["get"]; !ok {
		t.Error("expected GET operation on /users")
	}
	if _, ok := usersPath["post"]; !ok {
		t.Error("expected POST operation on /users")
	}

	// Verify response schemas contain correct types.
	getOp, ok := usersPath["get"].(map[string]any)
	if !ok {
		t.Fatal("expected GET operation to be an object")
	}
	responses, ok := getOp["responses"].(map[string]any)
	if !ok {
		t.Fatal("expected responses to be an object")
	}
	resp200, ok := responses["200"].(map[string]any)
	if !ok {
		t.Fatal("expected 200 response to be an object")
	}
	content, ok := resp200["content"].(map[string]any)
	if !ok {
		t.Fatal("expected content to be an object")
	}
	jsonContent, ok := content["application/json"].(map[string]any)
	if !ok {
		t.Fatal("expected application/json content to be an object")
	}
	schema, ok := jsonContent["schema"].(map[string]any)
	if !ok {
		t.Fatal("expected schema to be an object")
	}
	if schema["type"] != "array" {
		t.Errorf("expected GET /users response schema type 'array', got %v", schema["type"])
	}
}

// TestPBI8_AC2_SpecValidation verifies that generated specs have all required
// OpenAPI fields and are structurally valid.
func TestPBI8_AC2_SpecValidation(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/health", ResponseSchema: `{"type":"object","properties":{"status":{"type":"string"}}}`},
	}

	specBytes, err := openapi.GenerateSpec(endpoints, "Validation Test")
	if err != nil {
		t.Fatalf("GenerateSpec failed: %v", err)
	}

	// Must be valid JSON.
	if !json.Valid(specBytes) {
		t.Fatal("generated spec is not valid JSON")
	}

	var doc map[string]any
	if err := json.Unmarshal(specBytes, &doc); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Required OpenAPI fields.
	for _, field := range []string{"openapi", "info", "paths"} {
		if _, ok := doc[field]; !ok {
			t.Errorf("missing required OpenAPI field %q", field)
		}
	}

	// Info must have title and version.
	info, ok := doc["info"].(map[string]any)
	if !ok {
		t.Fatal("info is not an object")
	}
	if info["title"] != "Validation Test" {
		t.Errorf("expected title 'Validation Test', got %v", info["title"])
	}
	if _, ok := info["version"]; !ok {
		t.Error("info missing version field")
	}

	// Path items contain correct operations and response definitions.
	paths, _ := doc["paths"].(map[string]any)
	healthPath, ok := paths["/health"].(map[string]any)
	if !ok {
		t.Fatal("expected /health path")
	}
	getOp, ok := healthPath["get"].(map[string]any)
	if !ok {
		t.Fatal("expected GET operation on /health")
	}
	responses, ok := getOp["responses"].(map[string]any)
	if !ok {
		t.Fatal("expected responses in GET /health")
	}
	if _, ok := responses["200"]; !ok {
		t.Error("expected 200 response in GET /health")
	}
}

// TestPBI8_AC3_PrismContainerServesEndpoints verifies that a Prism container
// starts with a mounted spec and serves the defined endpoints.
func TestPBI8_AC3_PrismContainerServesEndpoints(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
		// Clean up spec files.
		specDir := filepath.Join(os.TempDir(), "heph-specs")
		_ = os.RemoveAll(specDir)
	})

	// Build config via template.
	tmpl := &templates.APIServiceTemplate{}
	node := model.DiagramNode{
		ID:   "prism-e2e",
		Type: model.ServiceTypeAPIService,
		Name: "E2E API",
		Config: json.RawMessage(`{
			"type": "api-service",
			"endpoints": [
				{"method": "GET", "path": "/users", "responseSchema": "{\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"integer\"},\"name\":{\"type\":\"string\"}}}}"},
				{"method": "POST", "path": "/users", "responseSchema": "{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"integer\"},\"name\":{\"type\":\"string\"}}}"},
				{"method": "GET", "path": "/health", "responseSchema": "{\"type\":\"object\",\"properties\":{\"status\":{\"type\":\"string\"}}}"}
			],
			"port": 4010
		}`),
	}

	hostPort := "14010"
	cfg, err := tmpl.Build(node, hostPort)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	id, err := o.CreateContainer(ctx, cfg)
	if err != nil {
		t.Fatalf("CreateContainer: %v", err)
	}

	if err := o.StartContainer(ctx, id); err != nil {
		t.Fatalf("StartContainer: %v", err)
	}

	baseURL := fmt.Sprintf("http://localhost:%s", hostPort)
	waitForPrism(t, baseURL+"/health", prismReadyTimeout)

	// Test GET /users → 200.
	t.Run("GET_users_returns_200", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/users")
		if err != nil {
			t.Fatalf("GET /users failed: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		if !json.Valid(body) {
			t.Error("response body is not valid JSON")
		}
	})

	// Test POST /users → 200.
	t.Run("POST_users_returns_200", func(t *testing.T) {
		resp, err := http.Post(baseURL+"/users", "application/json", nil)
		if err != nil {
			t.Fatalf("POST /users failed: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		if !json.Valid(body) {
			t.Error("response body is not valid JSON")
		}
	})

	// Test GET /health → 200.
	t.Run("GET_health_returns_200", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("GET /health failed: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})
}

// TestPBI8_AC4_ResponseSchemaMatching verifies that Prism returns responses
// matching the defined schemas.
func TestPBI8_AC4_ResponseSchemaMatching(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
		specDir := filepath.Join(os.TempDir(), "heph-specs")
		_ = os.RemoveAll(specDir)
	})

	tmpl := &templates.APIServiceTemplate{}
	node := model.DiagramNode{
		ID:   "prism-schema",
		Type: model.ServiceTypeAPIService,
		Name: "Schema API",
		Config: json.RawMessage(`{
			"type": "api-service",
			"endpoints": [
				{"method": "GET", "path": "/item", "responseSchema": "{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"integer\"},\"name\":{\"type\":\"string\"}},\"required\":[\"id\",\"name\"]}"}
			],
			"port": 4010
		}`),
	}

	hostPort := "14011"
	cfg, err := tmpl.Build(node, hostPort)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	id, err := o.CreateContainer(ctx, cfg)
	if err != nil {
		t.Fatalf("CreateContainer: %v", err)
	}
	if err := o.StartContainer(ctx, id); err != nil {
		t.Fatalf("StartContainer: %v", err)
	}

	baseURL := fmt.Sprintf("http://localhost:%s", hostPort)
	waitForPrism(t, baseURL+"/item", prismReadyTimeout)

	resp, err := http.Get(baseURL + "/item")
	if err != nil {
		t.Fatalf("GET /item failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	if !json.Valid(body) {
		t.Fatalf("response is not valid JSON: %s", string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Prism should return an object with id (number) and name (string).
	if _, ok := result["id"].(float64); !ok {
		t.Errorf("expected 'id' to be a number, got %T", result["id"])
	}
	if _, ok := result["name"].(string); !ok {
		t.Errorf("expected 'name' to be a string, got %T", result["name"])
	}

	t.Logf("Prism response: %s", string(body))
}

// TestPBI8_AC5_CrossServiceDNSResolution verifies that containers can reach
// each other via Docker DNS within the shared network.
func TestPBI8_AC5_CrossServiceDNSResolution(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
		specDir := filepath.Join(os.TempDir(), "heph-specs")
		_ = os.RemoveAll(specDir)
	})

	// Start the API service (Prism).
	tmpl := &templates.APIServiceTemplate{}
	apiNode := model.DiagramNode{
		ID:   "prism-dns",
		Type: model.ServiceTypeAPIService,
		Name: "user-service",
		Config: json.RawMessage(`{
			"type": "api-service",
			"endpoints": [
				{"method": "GET", "path": "/users", "responseSchema": "{\"type\":\"array\",\"items\":{\"type\":\"object\"}}"}
			],
			"port": 4010
		}`),
	}

	apiCfg, err := tmpl.Build(apiNode, "14012")
	if err != nil {
		t.Fatalf("Build API: %v", err)
	}

	apiID, err := o.CreateContainer(ctx, apiCfg)
	if err != nil {
		t.Fatalf("CreateContainer API: %v", err)
	}
	if err := o.StartContainer(ctx, apiID); err != nil {
		t.Fatalf("StartContainer API: %v", err)
	}

	// Wait for Prism to be ready via host port.
	waitForPrism(t, "http://localhost:14012/users", prismReadyTimeout)

	// Start a helper container (curlimages/curl) that will curl the API service via Docker DNS.
	helperID, err := o.CreateContainer(ctx, docker.ContainerConfig{
		Image:    "curlimages/curl:latest",
		Name:     "dns-helper",
		Cmd:      []string{"curl", "-sf", "--max-time", "10", "http://user-service:4010/users"},
		Hostname: "dns-helper",
	})
	if err != nil {
		t.Fatalf("CreateContainer helper: %v", err)
	}
	if err := o.StartContainer(ctx, helperID); err != nil {
		t.Fatalf("StartContainer helper: %v", err)
	}

	// Wait for the helper container to finish executing.
	deadline := time.Now().Add(30 * time.Second)
	var helperInfo *docker.ContainerInfo
	for time.Now().Before(deadline) {
		helperInfo, err = o.InspectContainer(ctx, helperID)
		if err != nil {
			t.Fatalf("InspectContainer helper: %v", err)
		}
		if helperInfo.Status == docker.StatusStopped {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if helperInfo.Status != docker.StatusStopped {
		t.Fatalf("helper container did not finish, status: %s", helperInfo.Status)
	}

	// The helper container ran curl against user-service:4010/users via Docker DNS.
	// If it exited cleanly (status = stopped, not error), curl succeeded, meaning
	// DNS resolution worked and Prism responded.
	t.Log("Cross-service DNS resolution verified: helper container reached user-service:4010 via Docker DNS")
}
