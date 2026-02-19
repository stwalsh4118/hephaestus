package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stwalsh4118/hephaestus/backend/internal/handler"
	"github.com/stwalsh4118/hephaestus/backend/internal/middleware"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
	"github.com/stwalsh4118/hephaestus/backend/internal/storage"
)

// buildServer creates a fully-wired test server identical to the real server.
func buildServer(t *testing.T, storageDir string) *httptest.Server {
	t.Helper()

	store, err := storage.NewFileStore(storageDir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	diagramHandler := handler.NewDiagramHandler(store)
	diagramHandler.RegisterRoutes(mux)

	wsHandler := handler.NewWebSocketHandler()
	wsHandler.RegisterRoutes(mux)

	corsHandler := middleware.CORS()(mux)

	return httptest.NewServer(corsHandler)
}

func validDiagramPayload() model.Diagram {
	return model.Diagram{
		ID:   "ignored",
		Name: "E2E Test Diagram",
		Nodes: []model.DiagramNode{
			{
				ID:          "n1",
				Type:        model.ServiceTypeAPIService,
				Name:        "My API",
				Description: "Test API service",
				Position:    &model.Position{X: 100, Y: 200},
				Config:      json.RawMessage(`{"type":"api-service","endpoints":[],"port":8080}`),
			},
			{
				ID:          "n2",
				Type:        model.ServiceTypePostgreSQL,
				Name:        "Main DB",
				Description: "Primary database",
				Position:    &model.Position{X: 300, Y: 200},
				Config:      json.RawMessage(`{"type":"postgresql","engine":"PostgreSQL","version":"16"}`),
			},
		},
		Edges: []model.DiagramEdge{
			{
				ID:     "e1",
				Source: "n1",
				Target: "n2",
				Label:  "connects to",
			},
		},
	}
}

// --- AC1: Server starts and listens on a configurable port ---

func TestAC1_ServerStartsAndListens(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("health status: got %q, want %q", body["status"], "ok")
	}
}

// --- AC2: POST /api/diagrams creates a new diagram and returns its ID ---

func TestAC2_CreateDiagram(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	d := validDiagramPayload()
	data, _ := json.Marshal(d)

	resp, err := http.Post(server.URL+"/api/diagrams", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status: got %d, want %d; body: %s", resp.StatusCode, http.StatusCreated, body)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result["id"] == "" {
		t.Error("expected non-empty id in create response")
	}
}

// --- AC3: GET /api/diagrams/:id retrieves a saved diagram ---

func TestAC3_GetDiagram(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	// Create first
	d := validDiagramPayload()
	data, _ := json.Marshal(d)
	createResp, err := http.Post(server.URL+"/api/diagrams", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	var created map[string]string
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	_, _ = io.Copy(io.Discard, createResp.Body)
	_ = createResp.Body.Close()

	// Get
	resp, err := http.Get(server.URL + "/api/diagrams/" + created["id"])
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var diagram model.Diagram
	if err := json.NewDecoder(resp.Body).Decode(&diagram); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if diagram.ID != created["id"] {
		t.Errorf("ID: got %q, want %q", diagram.ID, created["id"])
	}
	if diagram.Name != "E2E Test Diagram" {
		t.Errorf("Name: got %q, want %q", diagram.Name, "E2E Test Diagram")
	}
	if len(diagram.Nodes) != 2 {
		t.Errorf("Nodes: got %d, want 2", len(diagram.Nodes))
	}
	if len(diagram.Edges) != 1 {
		t.Errorf("Edges: got %d, want 1", len(diagram.Edges))
	}
}

// --- AC4: PUT /api/diagrams/:id updates an existing diagram ---

func TestAC4_UpdateDiagram(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	// Create
	d := validDiagramPayload()
	data, _ := json.Marshal(d)
	createResp, err := http.Post(server.URL+"/api/diagrams", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	var created map[string]string
	_ = json.NewDecoder(createResp.Body).Decode(&created)
	_, _ = io.Copy(io.Discard, createResp.Body)
	_ = createResp.Body.Close()

	// Update
	d.Name = "Updated E2E Diagram"
	d.Nodes = d.Nodes[:1] // Remove second node
	updateData, _ := json.Marshal(d)

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/diagrams/"+created["id"], bytes.NewReader(updateData))
	req.Header.Set("Content-Type", "application/json")
	updateResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, updateResp.Body); _ = updateResp.Body.Close() }()

	if updateResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(updateResp.Body)
		t.Fatalf("status: got %d, want %d; body: %s", updateResp.StatusCode, http.StatusOK, body)
	}

	// Verify update via GET
	getResp, err := http.Get(server.URL + "/api/diagrams/" + created["id"])
	if err != nil {
		t.Fatalf("GET after update: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, getResp.Body); _ = getResp.Body.Close() }()

	var diagram model.Diagram
	_ = json.NewDecoder(getResp.Body).Decode(&diagram)
	if diagram.Name != "Updated E2E Diagram" {
		t.Errorf("Name: got %q, want %q", diagram.Name, "Updated E2E Diagram")
	}
	if len(diagram.Nodes) != 1 {
		t.Errorf("Nodes: got %d, want 1", len(diagram.Nodes))
	}
}

// --- AC5: Diagrams persist to storage and survive server restarts ---

func TestAC5_PersistenceAcrossRestarts(t *testing.T) {
	dir := t.TempDir()

	// Start server 1, create diagram
	server1 := buildServer(t, dir)
	d := validDiagramPayload()
	data, _ := json.Marshal(d)
	createResp, err := http.Post(server1.URL+"/api/diagrams", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	var created map[string]string
	_ = json.NewDecoder(createResp.Body).Decode(&created)
	_, _ = io.Copy(io.Discard, createResp.Body)
	_ = createResp.Body.Close()

	// Shut down server 1
	server1.Close()

	// Start server 2 with same storage directory
	server2 := buildServer(t, dir)
	defer server2.Close()

	// Retrieve diagram from server 2
	resp, err := http.Get(server2.URL + "/api/diagrams/" + created["id"])
	if err != nil {
		t.Fatalf("GET from server2: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var diagram model.Diagram
	_ = json.NewDecoder(resp.Body).Decode(&diagram)
	if diagram.ID != created["id"] {
		t.Errorf("ID: got %q, want %q", diagram.ID, created["id"])
	}
	if diagram.Name != "E2E Test Diagram" {
		t.Errorf("Name: got %q, want %q", diagram.Name, "E2E Test Diagram")
	}
}

// --- AC6: WebSocket at /ws/status accepts connections and responds to ping ---

func TestAC6_WebSocketPingPong(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/status"

	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusSwitchingProtocols)
	}

	// Send ping, verify connection stays alive
	if err := conn.WriteMessage(websocket.PingMessage, []byte("e2e")); err != nil {
		t.Fatalf("write ping: %v", err)
	}

	// Connection should still be writable after ping
	if err := conn.WriteMessage(websocket.TextMessage, []byte("alive")); err != nil {
		t.Fatalf("write after ping: %v", err)
	}

	// Graceful close
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
		t.Fatalf("write close: %v", err)
	}
}

// --- AC7: CORS configured to allow frontend origin ---

func TestAC7_CORSHeaders(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CORS_ORIGIN", "http://localhost:3000")
	server := buildServer(t, dir)
	defer server.Close()

	// Preflight OPTIONS request
	req, _ := http.NewRequest(http.MethodOptions, server.URL+"/api/diagrams", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("OPTIONS: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("preflight status: got %d, want %d", resp.StatusCode, http.StatusNoContent)
	}

	origin := resp.Header.Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:3000" {
		t.Errorf("Allow-Origin: got %q, want %q", origin, "http://localhost:3000")
	}

	methods := resp.Header.Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("missing Access-Control-Allow-Methods")
	}

	headers := resp.Header.Get("Access-Control-Allow-Headers")
	if headers == "" {
		t.Error("missing Access-Control-Allow-Headers")
	}

	// Regular request should also have CORS headers
	getResp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, getResp.Body); _ = getResp.Body.Close() }()

	corsOrigin := getResp.Header.Get("Access-Control-Allow-Origin")
	if corsOrigin != "http://localhost:3000" {
		t.Errorf("CORS origin on GET: got %q, want %q", corsOrigin, "http://localhost:3000")
	}
}

// --- Error handling ---

func TestErrorHandling_GetNotFound(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/diagrams/nonexistent-id")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	var body map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] == "" {
		t.Error("expected error message in response")
	}
}

func TestErrorHandling_PostInvalidBody(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/diagrams", "application/json", bytes.NewReader([]byte("not json")))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

// --- Full CRUD lifecycle ---

func TestFullCRUDLifecycle(t *testing.T) {
	dir := t.TempDir()
	server := buildServer(t, dir)
	defer server.Close()

	// 1. Create
	d := validDiagramPayload()
	data, _ := json.Marshal(d)
	createResp, err := http.Post(server.URL+"/api/diagrams", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create status: %d", createResp.StatusCode)
	}
	var created map[string]string
	_ = json.NewDecoder(createResp.Body).Decode(&created)
	_, _ = io.Copy(io.Discard, createResp.Body)
	_ = createResp.Body.Close()
	id := created["id"]

	// 2. Get
	getResp, err := http.Get(server.URL + "/api/diagrams/" + id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get status: %d", getResp.StatusCode)
	}
	var fetched model.Diagram
	_ = json.NewDecoder(getResp.Body).Decode(&fetched)
	_, _ = io.Copy(io.Discard, getResp.Body)
	_ = getResp.Body.Close()

	if fetched.Name != "E2E Test Diagram" {
		t.Errorf("get name: %q", fetched.Name)
	}

	// 3. Update
	fetched.Name = "Lifecycle Updated"
	updateData, _ := json.Marshal(fetched)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/diagrams/"+id, bytes.NewReader(updateData))
	req.Header.Set("Content-Type", "application/json")
	updateResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updateResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(updateResp.Body)
		t.Fatalf("update status: %d; body: %s", updateResp.StatusCode, body)
	}
	_, _ = io.Copy(io.Discard, updateResp.Body)
	_ = updateResp.Body.Close()

	// 4. Get (confirm update)
	getResp2, err := http.Get(server.URL + "/api/diagrams/" + id)
	if err != nil {
		t.Fatalf("get2: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, getResp2.Body); _ = getResp2.Body.Close() }()
	var updated model.Diagram
	_ = json.NewDecoder(getResp2.Body).Decode(&updated)

	if updated.Name != "Lifecycle Updated" {
		t.Errorf("updated name: got %q, want %q", updated.Name, "Lifecycle Updated")
	}
}

