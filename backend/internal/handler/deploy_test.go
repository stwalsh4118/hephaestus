package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/deploy"
	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// mockDeployOrchestrator implements docker.Orchestrator for deploy handler tests.
type mockDeployOrchestrator struct {
	createContainerFn func(ctx context.Context, cfg docker.ContainerConfig) (string, error)
	startContainerFn  func(ctx context.Context, id string) error
	teardownAllFn     func(ctx context.Context) error
	healthCheckFn     func(ctx context.Context, id string) (docker.ContainerStatus, error)
}

func newMockDeployOrch() *mockDeployOrchestrator {
	return &mockDeployOrchestrator{
		createContainerFn: func(_ context.Context, cfg docker.ContainerConfig) (string, error) {
			return "ctr-" + cfg.Name, nil
		},
		startContainerFn:  func(_ context.Context, _ string) error { return nil },
		teardownAllFn:     func(_ context.Context) error { return nil },
		healthCheckFn: func(_ context.Context, _ string) (docker.ContainerStatus, error) {
			return docker.StatusRunning, nil
		},
	}
}

func (m *mockDeployOrchestrator) CreateContainer(ctx context.Context, cfg docker.ContainerConfig) (string, error) {
	return m.createContainerFn(ctx, cfg)
}
func (m *mockDeployOrchestrator) StartContainer(ctx context.Context, id string) error {
	return m.startContainerFn(ctx, id)
}
func (m *mockDeployOrchestrator) StopContainer(_ context.Context, _ string) error   { return nil }
func (m *mockDeployOrchestrator) RemoveContainer(_ context.Context, _ string) error { return nil }
func (m *mockDeployOrchestrator) ListContainers(_ context.Context) ([]docker.ContainerInfo, error) {
	return nil, nil
}
func (m *mockDeployOrchestrator) InspectContainer(_ context.Context, _ string) (*docker.ContainerInfo, error) {
	return nil, nil
}
func (m *mockDeployOrchestrator) CreateNetwork(_ context.Context) error  { return nil }
func (m *mockDeployOrchestrator) RemoveNetwork(_ context.Context) error  { return nil }
func (m *mockDeployOrchestrator) HealthCheck(ctx context.Context, id string) (docker.ContainerStatus, error) {
	return m.healthCheckFn(ctx, id)
}
func (m *mockDeployOrchestrator) TeardownAll(ctx context.Context) error {
	return m.teardownAllFn(ctx)
}

func deployTestDiagramJSON() string {
	d := model.Diagram{
		ID:   "deploy-test-1",
		Name: "test",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1", Position: &model.Position{X: 0, Y: 0}},
		},
		Edges: []model.DiagramEdge{},
	}
	b, _ := json.Marshal(d)
	return string(b)
}

func newTestDeployHandler() (*DeployHandler, *mockDeployOrchestrator) {
	orch := newMockDeployOrch()
	mgr := deploy.NewDeploymentManager(orch)
	return NewDeployHandler(mgr), orch
}

func TestDeployHandler_Deploy_Success(t *testing.T) {
	h, _ := newTestDeployHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(deployTestDiagramJSON()))
	rec := httptest.NewRecorder()

	h.Deploy(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected 202, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp deploy.DeployResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != "deploying" {
		t.Errorf("expected status deploying, got %s", resp.Status)
	}
}

func TestDeployHandler_Deploy_InvalidJSON(t *testing.T) {
	h, _ := newTestDeployHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader("not json"))
	rec := httptest.NewRecorder()

	h.Deploy(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestDeployHandler_Deploy_InvalidDiagram(t *testing.T) {
	h, _ := newTestDeployHandler()

	// Valid JSON but missing required fields (no id, no name, no nodes, no edges)
	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(`{"name":""}`))
	rec := httptest.NewRecorder()

	h.Deploy(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestDeployHandler_Deploy_AlreadyDeploying(t *testing.T) {
	h, _ := newTestDeployHandler()

	// First deploy
	req1 := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(deployTestDiagramJSON()))
	rec1 := httptest.NewRecorder()
	h.Deploy(rec1, req1)

	// Second deploy → conflict
	req2 := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(deployTestDiagramJSON()))
	rec2 := httptest.NewRecorder()
	h.Deploy(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d: %s", rec2.Code, rec2.Body.String())
	}
}

func TestDeployHandler_Teardown_Success(t *testing.T) {
	h, _ := newTestDeployHandler()

	// Deploy first
	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(deployTestDiagramJSON()))
	rec := httptest.NewRecorder()
	h.Deploy(rec, req)

	// Teardown
	req2 := httptest.NewRequest(http.MethodDelete, "/api/deploy", nil)
	rec2 := httptest.NewRecorder()
	h.Teardown(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec2.Code, rec2.Body.String())
	}

	var resp deploy.DeployResponse
	if err := json.NewDecoder(rec2.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != "idle" {
		t.Errorf("expected status idle, got %s", resp.Status)
	}
}

func TestDeployHandler_Teardown_NotDeployed(t *testing.T) {
	h, _ := newTestDeployHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/deploy", nil)
	rec := httptest.NewRecorder()
	h.Teardown(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestDeployHandler_Status_Idle(t *testing.T) {
	h, _ := newTestDeployHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/deploy/status", nil)
	rec := httptest.NewRecorder()
	h.Status(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp deploy.StatusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.DeployStatus != deploy.StatusIdle {
		t.Errorf("expected idle, got %s", resp.DeployStatus)
	}
}

func TestDeployHandler_Update_AddedNode(t *testing.T) {
	h, _ := newTestDeployHandler()

	// Deploy first
	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(deployTestDiagramJSON()))
	rec := httptest.NewRecorder()
	h.Deploy(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("deploy: expected 202, got %d: %s", rec.Code, rec.Body.String())
	}

	// Update with an extra node
	updatedDiagram := model.Diagram{
		ID:   "deploy-test-1",
		Name: "test",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1", Position: &model.Position{X: 0, Y: 0}},
			{ID: "n2", Type: model.ServiceTypePostgreSQL, Name: "pg-1", Position: &model.Position{X: 100, Y: 0}},
		},
		Edges: []model.DiagramEdge{},
	}
	b, _ := json.Marshal(updatedDiagram)

	req2 := httptest.NewRequest(http.MethodPut, "/api/deploy", strings.NewReader(string(b)))
	rec2 := httptest.NewRecorder()
	h.Update(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec2.Code, rec2.Body.String())
	}

	var resp deploy.StatusResponse
	if err := json.NewDecoder(rec2.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.DeployStatus != deploy.StatusDeployed {
		t.Errorf("expected deployed, got %s", resp.DeployStatus)
	}
	if len(resp.NodeStatuses) != 2 {
		t.Errorf("expected 2 node statuses, got %d", len(resp.NodeStatuses))
	}
}

func TestDeployHandler_Update_RemovedNode(t *testing.T) {
	h, _ := newTestDeployHandler()

	// Deploy with two nodes first
	twoNodeDiagram := model.Diagram{
		ID:   "deploy-test-2",
		Name: "test",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1", Position: &model.Position{X: 0, Y: 0}},
			{ID: "n2", Type: model.ServiceTypePostgreSQL, Name: "pg-1", Position: &model.Position{X: 100, Y: 0}},
		},
		Edges: []model.DiagramEdge{},
	}
	b, _ := json.Marshal(twoNodeDiagram)
	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(string(b)))
	rec := httptest.NewRecorder()
	h.Deploy(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("deploy: expected 202, got %d: %s", rec.Code, rec.Body.String())
	}

	// Update with one node removed
	oneNodeDiagram := model.Diagram{
		ID:   "deploy-test-2",
		Name: "test",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1", Position: &model.Position{X: 0, Y: 0}},
		},
		Edges: []model.DiagramEdge{},
	}
	b2, _ := json.Marshal(oneNodeDiagram)
	req2 := httptest.NewRequest(http.MethodPut, "/api/deploy", strings.NewReader(string(b2)))
	rec2 := httptest.NewRecorder()
	h.Update(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec2.Code, rec2.Body.String())
	}

	var resp deploy.StatusResponse
	if err := json.NewDecoder(rec2.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.NodeStatuses) != 1 {
		t.Errorf("expected 1 node status, got %d", len(resp.NodeStatuses))
	}
}

func TestDeployHandler_Update_NotDeployed(t *testing.T) {
	h, _ := newTestDeployHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/deploy", strings.NewReader(deployTestDiagramJSON()))
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestDeployHandler_Status_Deployed(t *testing.T) {
	h, _ := newTestDeployHandler()

	// Deploy first
	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(deployTestDiagramJSON()))
	rec := httptest.NewRecorder()
	h.Deploy(rec, req)

	// Get status
	req2 := httptest.NewRequest(http.MethodGet, "/api/deploy/status", nil)
	rec2 := httptest.NewRecorder()
	h.Status(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec2.Code)
	}

	var resp deploy.StatusResponse
	if err := json.NewDecoder(rec2.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.DeployStatus != deploy.StatusDeployed {
		t.Errorf("expected deployed, got %s", resp.DeployStatus)
	}
	if len(resp.NodeStatuses) != 1 {
		t.Errorf("expected 1 node status, got %d", len(resp.NodeStatuses))
	}
}
