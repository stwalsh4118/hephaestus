package e2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stwalsh4118/hephaestus/backend/internal/deploy"
	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/handler"
	"github.com/stwalsh4118/hephaestus/backend/internal/middleware"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// mockOrchestrator is a test orchestrator that tracks container state in memory.
type mockOrchestrator struct {
	mu         sync.Mutex
	containers map[string]docker.ContainerStatus // containerID → status
	counter    int
}

func newMockOrchestrator() *mockOrchestrator {
	return &mockOrchestrator{containers: make(map[string]docker.ContainerStatus)}
}

func (m *mockOrchestrator) CreateContainer(_ context.Context, cfg docker.ContainerConfig) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter++
	id := "ctr-" + cfg.Name
	m.containers[id] = docker.StatusCreated
	return id, nil
}

func (m *mockOrchestrator) StartContainer(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.containers[id] = docker.StatusRunning
	return nil
}

func (m *mockOrchestrator) StopContainer(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.containers[id] = docker.StatusStopped
	return nil
}

func (m *mockOrchestrator) RemoveContainer(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.containers, id)
	return nil
}

func (m *mockOrchestrator) ListContainers(_ context.Context) ([]docker.ContainerInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	infos := make([]docker.ContainerInfo, 0, len(m.containers))
	for id, status := range m.containers {
		infos = append(infos, docker.ContainerInfo{ID: id, Status: status})
	}
	return infos, nil
}

func (m *mockOrchestrator) InspectContainer(_ context.Context, id string) (*docker.ContainerInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	status, ok := m.containers[id]
	if !ok {
		return nil, nil
	}
	return &docker.ContainerInfo{ID: id, Status: status}, nil
}

func (m *mockOrchestrator) CreateNetwork(_ context.Context) error  { return nil }
func (m *mockOrchestrator) RemoveNetwork(_ context.Context) error  { return nil }
func (m *mockOrchestrator) TeardownAll(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.containers = make(map[string]docker.ContainerStatus)
	return nil
}

func (m *mockOrchestrator) HealthCheck(_ context.Context, id string) (docker.ContainerStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	status, ok := m.containers[id]
	if !ok {
		return docker.StatusError, nil
	}
	return status, nil
}

func (m *mockOrchestrator) containerCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.containers)
}

// deployTestEnv holds all components needed for deploy e2e tests.
type deployTestEnv struct {
	server    *httptest.Server
	orch      *mockOrchestrator
	deployMgr *deploy.DeploymentManager
	statusHub *handler.StatusHub
}

// buildDeployServer creates a test server with deploy endpoints and a mock orchestrator.
func buildDeployServer(t *testing.T) *deployTestEnv {
	t.Helper()

	orch := newMockOrchestrator()
	deployMgr := deploy.NewDeploymentManager(orch)
	statusHub := handler.NewStatusHub()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	deployHandler := handler.NewDeployHandler(deployMgr)
	deployHandler.RegisterRoutes(mux)

	wsHandler := handler.NewWebSocketHandler(statusHub)
	wsHandler.RegisterRoutes(mux)

	corsHandler := middleware.CORS()(mux)

	return &deployTestEnv{
		server:    httptest.NewServer(corsHandler),
		orch:      orch,
		deployMgr: deployMgr,
		statusHub: statusHub,
	}
}

func deployDiagram() model.Diagram {
	return model.Diagram{
		ID:   "deploy-e2e-1",
		Name: "Deploy E2E",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1", Position: &model.Position{X: 0, Y: 0}},
			{ID: "n2", Type: model.ServiceTypePostgreSQL, Name: "pg-1", Position: &model.Position{X: 100, Y: 0}},
		},
		Edges: []model.DiagramEdge{},
	}
}

func marshalDiagram(d model.Diagram) string {
	b, _ := json.Marshal(d)
	return string(b)
}

// --- PBI 9 AC1: Deploy button sends diagram to POST /api/deploy ---

func TestPBI9_AC1_DeployDiagram(t *testing.T) {
	env := buildDeployServer(t)
	server := env.server
	defer server.Close()

	d := deployDiagram()
	resp, err := http.Post(server.URL+"/api/deploy", "application/json", strings.NewReader(marshalDiagram(d)))
	if err != nil {
		t.Fatalf("POST /api/deploy: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 202, got %d: %s", resp.StatusCode, body)
	}

	var result deploy.DeployResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.Status != "deploying" {
		t.Errorf("expected status deploying, got %s", result.Status)
	}
}

// --- PBI 9 AC2: Backend creates containers and reports status ---

func TestPBI9_AC2_ContainersCreatedAndStatus(t *testing.T) {
	env := buildDeployServer(t)
	server, orch := env.server, env.orch
	defer server.Close()

	d := deployDiagram()
	resp, err := http.Post(server.URL+"/api/deploy", "application/json", strings.NewReader(marshalDiagram(d)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	// Verify containers were created
	if count := orch.containerCount(); count != 2 {
		t.Errorf("expected 2 containers, got %d", count)
	}

	// Get status endpoint
	statusResp, err := http.Get(server.URL + "/api/deploy/status")
	if err != nil {
		t.Fatalf("GET /api/deploy/status: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, statusResp.Body); _ = statusResp.Body.Close() }()

	if statusResp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", statusResp.StatusCode)
	}

	var status deploy.StatusResponse
	if err := json.NewDecoder(statusResp.Body).Decode(&status); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if status.DeployStatus != deploy.StatusDeployed {
		t.Errorf("expected deployed, got %s", status.DeployStatus)
	}
	if len(status.NodeStatuses) != 2 {
		t.Errorf("expected 2 node statuses, got %d", len(status.NodeStatuses))
	}
}

// --- PBI 9 AC3: WebSocket status stream connects and broadcasts ---

func TestPBI9_AC3_WebSocketStatusStream(t *testing.T) {
	env := buildDeployServer(t)
	server := env.server
	defer server.Close()

	// 1. Connect WebSocket
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/status"
	dialer := websocket.Dialer{}

	conn, wsResp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	if wsResp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("status: got %d, want 101", wsResp.StatusCode)
	}

	// 2. Deploy a diagram (creates containers)
	d := deployDiagram()
	resp, err := http.Post(server.URL+"/api/deploy", "application/json", strings.NewReader(marshalDiagram(d)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	// 3. Simulate health polling by building and broadcasting a status message
	containerStatuses := make(map[string]docker.ContainerStatus)
	env.orch.mu.Lock()
	for id, status := range env.orch.containers {
		containerStatuses[id] = status
	}
	env.orch.mu.Unlock()
	statusMsg := env.deployMgr.BuildStatusMessage(containerStatuses)
	env.statusHub.Broadcast(statusMsg)

	// 4. Read the broadcast message from WebSocket
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, rawMsg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read ws message: %v", err)
	}

	var received deploy.StatusMessage
	if err := json.Unmarshal(rawMsg, &received); err != nil {
		t.Fatalf("unmarshal ws message: %v", err)
	}
	if received.Type != deploy.StatusMessageType {
		t.Errorf("expected type %q, got %q", deploy.StatusMessageType, received.Type)
	}
	if received.DeployStatus != deploy.StatusDeployed {
		t.Errorf("expected deployed, got %s", received.DeployStatus)
	}
	if len(received.NodeStatuses) != 2 {
		t.Errorf("expected 2 node statuses, got %d", len(received.NodeStatuses))
	}
}

// --- PBI 9 AC4: Adding a node deploys only the new container ---

func TestPBI9_AC4_IncrementalAddNode(t *testing.T) {
	env := buildDeployServer(t)
	server, orch := env.server, env.orch
	defer server.Close()

	// Initial deploy with 2 nodes
	d := deployDiagram()
	resp, err := http.Post(server.URL+"/api/deploy", "application/json", strings.NewReader(marshalDiagram(d)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if count := orch.containerCount(); count != 2 {
		t.Fatalf("expected 2 containers after deploy, got %d", count)
	}

	// Add a third node via PUT
	d.Nodes = append(d.Nodes, model.DiagramNode{
		ID: "n3", Type: model.ServiceTypeNginx, Name: "nginx-1",
		Position: &model.Position{X: 200, Y: 0},
	})

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/deploy", strings.NewReader(marshalDiagram(d)))
	req.Header.Set("Content-Type", "application/json")
	putResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, putResp.Body); _ = putResp.Body.Close() }()

	if putResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(putResp.Body)
		t.Fatalf("expected 200, got %d: %s", putResp.StatusCode, body)
	}

	// Now should have 3 containers
	if count := orch.containerCount(); count != 3 {
		t.Errorf("expected 3 containers after add, got %d", count)
	}

	var status deploy.StatusResponse
	if err := json.NewDecoder(putResp.Body).Decode(&status); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(status.NodeStatuses) != 3 {
		t.Errorf("expected 3 node statuses, got %d", len(status.NodeStatuses))
	}
}

// --- PBI 9 AC5: Removing a node removes only that container ---

func TestPBI9_AC5_IncrementalRemoveNode(t *testing.T) {
	env := buildDeployServer(t)
	server, orch := env.server, env.orch
	defer server.Close()

	// Deploy with 2 nodes
	d := deployDiagram()
	resp, err := http.Post(server.URL+"/api/deploy", "application/json", strings.NewReader(marshalDiagram(d)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if count := orch.containerCount(); count != 2 {
		t.Fatalf("expected 2 containers after deploy, got %d", count)
	}

	// Remove second node via PUT
	d.Nodes = d.Nodes[:1]
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/deploy", strings.NewReader(marshalDiagram(d)))
	req.Header.Set("Content-Type", "application/json")
	putResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, putResp.Body); _ = putResp.Body.Close() }()

	if putResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(putResp.Body)
		t.Fatalf("expected 200, got %d: %s", putResp.StatusCode, body)
	}

	// Now should have 1 container
	if count := orch.containerCount(); count != 1 {
		t.Errorf("expected 1 container after remove, got %d", count)
	}
}

// --- PBI 9 AC6: Teardown removes all containers and resets status ---

func TestPBI9_AC6_Teardown(t *testing.T) {
	env := buildDeployServer(t)
	server, orch := env.server, env.orch
	defer server.Close()

	// Deploy
	d := deployDiagram()
	resp, err := http.Post(server.URL+"/api/deploy", "application/json", strings.NewReader(marshalDiagram(d)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if count := orch.containerCount(); count != 2 {
		t.Fatalf("expected 2 containers, got %d", count)
	}

	// Teardown
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/deploy", nil)
	delResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, delResp.Body); _ = delResp.Body.Close() }()

	if delResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(delResp.Body)
		t.Fatalf("expected 200, got %d: %s", delResp.StatusCode, body)
	}

	// Verify no containers
	if count := orch.containerCount(); count != 0 {
		t.Errorf("expected 0 containers after teardown, got %d", count)
	}

	// Verify status is idle
	statusResp, err := http.Get(server.URL + "/api/deploy/status")
	if err != nil {
		t.Fatalf("GET status: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, statusResp.Body); _ = statusResp.Body.Close() }()

	var status deploy.StatusResponse
	_ = json.NewDecoder(statusResp.Body).Decode(&status)
	if status.DeployStatus != deploy.StatusIdle {
		t.Errorf("expected idle, got %s", status.DeployStatus)
	}
}

// --- PBI 9 AC7: Error surfaced in API response ---

func TestPBI9_AC7_ErrorSurfaced(t *testing.T) {
	env := buildDeployServer(t)
	server := env.server
	defer server.Close()

	// Try to teardown when nothing is deployed
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/deploy", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", resp.StatusCode)
	}

	var body map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] == "" {
		t.Error("expected error message in response")
	}
}

// --- Full Deploy Lifecycle ---

func TestPBI9_FullDeployLifecycle(t *testing.T) {
	env := buildDeployServer(t)
	server, orch := env.server, env.orch
	defer server.Close()

	// 1. Status should be idle
	statusResp, err := http.Get(server.URL + "/api/deploy/status")
	if err != nil {
		t.Fatalf("GET status: %v", err)
	}
	var status deploy.StatusResponse
	_ = json.NewDecoder(statusResp.Body).Decode(&status)
	_, _ = io.Copy(io.Discard, statusResp.Body)
	_ = statusResp.Body.Close()
	if status.DeployStatus != deploy.StatusIdle {
		t.Fatalf("initial status: expected idle, got %s", status.DeployStatus)
	}

	// 2. Deploy with 2 nodes
	d := deployDiagram()
	resp, err := http.Post(server.URL+"/api/deploy", "application/json", strings.NewReader(marshalDiagram(d)))
	if err != nil {
		t.Fatalf("POST deploy: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("deploy status: %d", resp.StatusCode)
	}

	// 3. Status should be deployed with 2 nodes
	statusResp2, err := http.Get(server.URL + "/api/deploy/status")
	if err != nil {
		t.Fatalf("GET status after deploy: %v", err)
	}
	var status2 deploy.StatusResponse
	_ = json.NewDecoder(statusResp2.Body).Decode(&status2)
	_, _ = io.Copy(io.Discard, statusResp2.Body)
	_ = statusResp2.Body.Close()
	if status2.DeployStatus != deploy.StatusDeployed {
		t.Fatalf("after deploy: expected deployed, got %s", status2.DeployStatus)
	}
	if len(status2.NodeStatuses) != 2 {
		t.Fatalf("after deploy: expected 2 nodes, got %d", len(status2.NodeStatuses))
	}

	// 4. Add a node via PUT
	d.Nodes = append(d.Nodes, model.DiagramNode{
		ID: "n3", Type: model.ServiceTypeNginx, Name: "nginx-1",
		Position: &model.Position{X: 200, Y: 0},
	})
	putReq, _ := http.NewRequest(http.MethodPut, server.URL+"/api/deploy", strings.NewReader(marshalDiagram(d)))
	putReq.Header.Set("Content-Type", "application/json")
	putResp, err := http.DefaultClient.Do(putReq)
	if err != nil {
		t.Fatalf("PUT add: %v", err)
	}
	_, _ = io.Copy(io.Discard, putResp.Body)
	_ = putResp.Body.Close()
	if putResp.StatusCode != http.StatusOK {
		t.Fatalf("PUT add: got %d", putResp.StatusCode)
	}
	if count := orch.containerCount(); count != 3 {
		t.Fatalf("after add: expected 3 containers, got %d", count)
	}

	// 5. Remove a node via PUT
	d.Nodes = d.Nodes[1:] // Remove first node (n1)
	putReq2, _ := http.NewRequest(http.MethodPut, server.URL+"/api/deploy", strings.NewReader(marshalDiagram(d)))
	putReq2.Header.Set("Content-Type", "application/json")
	putResp2, err := http.DefaultClient.Do(putReq2)
	if err != nil {
		t.Fatalf("PUT remove: %v", err)
	}
	_, _ = io.Copy(io.Discard, putResp2.Body)
	_ = putResp2.Body.Close()
	if putResp2.StatusCode != http.StatusOK {
		t.Fatalf("PUT remove: got %d", putResp2.StatusCode)
	}
	if count := orch.containerCount(); count != 2 {
		t.Fatalf("after remove: expected 2 containers, got %d", count)
	}

	// 6. Teardown
	delReq, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/deploy", nil)
	delResp, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatalf("DELETE teardown: %v", err)
	}
	_, _ = io.Copy(io.Discard, delResp.Body)
	_ = delResp.Body.Close()
	if delResp.StatusCode != http.StatusOK {
		t.Fatalf("teardown: got %d", delResp.StatusCode)
	}
	if count := orch.containerCount(); count != 0 {
		t.Fatalf("after teardown: expected 0 containers, got %d", count)
	}

	// 7. Status should be idle again
	statusResp3, err := http.Get(server.URL + "/api/deploy/status")
	if err != nil {
		t.Fatalf("GET status after teardown: %v", err)
	}
	var status3 deploy.StatusResponse
	_ = json.NewDecoder(statusResp3.Body).Decode(&status3)
	_, _ = io.Copy(io.Discard, statusResp3.Body)
	_ = statusResp3.Body.Close()
	if status3.DeployStatus != deploy.StatusIdle {
		t.Fatalf("final status: expected idle, got %s", status3.DeployStatus)
	}

	// 8. Connect WebSocket to verify status streaming is available
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/status"
	dialer := websocket.Dialer{}
	conn, wsResp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("ws dial: %v", err)
	}
	defer func() { _ = conn.Close() }()
	if wsResp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("ws status: got %d", wsResp.StatusCode)
	}
	// Set a short read deadline to verify we can read (no data expected since idle)
	_ = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, readErr := conn.ReadMessage()
	// Expected: deadline exceeded (no messages while idle)
	if readErr == nil {
		t.Error("expected read timeout, got message")
	}
}
