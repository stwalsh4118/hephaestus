package deploy

import (
	"context"
	"errors"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// mockOrchestrator implements docker.Orchestrator for testing.
type mockOrchestrator struct {
	createNetworkCalls int
	createContainerFn  func(ctx context.Context, cfg docker.ContainerConfig) (string, error)
	startContainerFn   func(ctx context.Context, id string) error
	stopContainerFn    func(ctx context.Context, id string) error
	removeContainerFn  func(ctx context.Context, id string) error
	healthCheckFn      func(ctx context.Context, id string) (docker.ContainerStatus, error)
	teardownAllFn      func(ctx context.Context) error

	createdConfigs  []docker.ContainerConfig
	startedIDs      []string
	stoppedIDs      []string
	removedIDs      []string
}

func newMockOrchestrator() *mockOrchestrator {
	m := &mockOrchestrator{}
	containerCount := 0
	m.createContainerFn = func(_ context.Context, cfg docker.ContainerConfig) (string, error) {
		containerCount++
		return "container-" + cfg.Name, nil
	}
	m.startContainerFn = func(_ context.Context, _ string) error { return nil }
	m.stopContainerFn = func(_ context.Context, _ string) error { return nil }
	m.removeContainerFn = func(_ context.Context, _ string) error { return nil }
	m.healthCheckFn = func(_ context.Context, _ string) (docker.ContainerStatus, error) {
		return docker.StatusRunning, nil
	}
	m.teardownAllFn = func(_ context.Context) error { return nil }
	return m
}

func (m *mockOrchestrator) CreateContainer(ctx context.Context, cfg docker.ContainerConfig) (string, error) {
	m.createdConfigs = append(m.createdConfigs, cfg)
	return m.createContainerFn(ctx, cfg)
}
func (m *mockOrchestrator) StartContainer(ctx context.Context, id string) error {
	m.startedIDs = append(m.startedIDs, id)
	return m.startContainerFn(ctx, id)
}
func (m *mockOrchestrator) StopContainer(ctx context.Context, id string) error {
	m.stoppedIDs = append(m.stoppedIDs, id)
	return m.stopContainerFn(ctx, id)
}
func (m *mockOrchestrator) RemoveContainer(ctx context.Context, id string) error {
	m.removedIDs = append(m.removedIDs, id)
	return m.removeContainerFn(ctx, id)
}
func (m *mockOrchestrator) ListContainers(_ context.Context) ([]docker.ContainerInfo, error) {
	return nil, nil
}
func (m *mockOrchestrator) InspectContainer(_ context.Context, _ string) (*docker.ContainerInfo, error) {
	return nil, nil
}
func (m *mockOrchestrator) CreateNetwork(_ context.Context) error {
	m.createNetworkCalls++
	return nil
}
func (m *mockOrchestrator) RemoveNetwork(_ context.Context) error { return nil }
func (m *mockOrchestrator) HealthCheck(ctx context.Context, id string) (docker.ContainerStatus, error) {
	return m.healthCheckFn(ctx, id)
}
func (m *mockOrchestrator) TeardownAll(ctx context.Context) error {
	return m.teardownAllFn(ctx)
}

func testDiagram() model.Diagram {
	return model.Diagram{
		ID:   "test-diag",
		Name: "test",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1"},
			{ID: "n2", Type: model.ServiceTypeNginx, Name: "nginx-1"},
		},
		Edges: []model.DiagramEdge{
			{ID: "e1", Source: "n2", Target: "n1", Label: "connects"},
		},
	}
}

func TestDeploy_Success(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	if err := mgr.Deploy(ctx, testDiagram()); err != nil {
		t.Fatalf("Deploy failed: %v", err)
	}

	if mgr.Status() != StatusDeployed {
		t.Errorf("expected status deployed, got %s", mgr.Status())
	}
	if orch.createNetworkCalls != 1 {
		t.Errorf("expected 1 CreateNetwork call, got %d", orch.createNetworkCalls)
	}
	if len(orch.createdConfigs) != 2 {
		t.Errorf("expected 2 containers created, got %d", len(orch.createdConfigs))
	}
	if len(orch.startedIDs) != 2 {
		t.Errorf("expected 2 containers started, got %d", len(orch.startedIDs))
	}
}

func TestDeploy_AlreadyDeploying(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	_ = mgr.Deploy(ctx, testDiagram())

	err := mgr.Deploy(ctx, testDiagram())
	if !errors.Is(err, ErrAlreadyDeploying) {
		t.Errorf("expected ErrAlreadyDeploying, got %v", err)
	}
}

func TestDeploy_StateTransition(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	if mgr.Status() != StatusIdle {
		t.Errorf("initial status should be idle, got %s", mgr.Status())
	}

	_ = mgr.Deploy(ctx, testDiagram())

	if mgr.Status() != StatusDeployed {
		t.Errorf("post-deploy status should be deployed, got %s", mgr.Status())
	}
}

func TestTeardown_Success(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	_ = mgr.Deploy(ctx, testDiagram())

	if err := mgr.Teardown(ctx); err != nil {
		t.Fatalf("Teardown failed: %v", err)
	}

	if mgr.Status() != StatusIdle {
		t.Errorf("expected idle after teardown, got %s", mgr.Status())
	}
}

func TestTeardown_NotDeployed(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)

	err := mgr.Teardown(context.Background())
	if !errors.Is(err, ErrNotDeployed) {
		t.Errorf("expected ErrNotDeployed, got %v", err)
	}
}

func TestGetStatus_Idle(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)

	status, nodes, err := mgr.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status != StatusIdle {
		t.Errorf("expected idle, got %s", status)
	}
	if len(nodes) != 0 {
		t.Errorf("expected 0 node statuses, got %d", len(nodes))
	}
}

func TestGetStatus_Deployed(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	_ = mgr.Deploy(ctx, testDiagram())

	status, nodes, err := mgr.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status != StatusDeployed {
		t.Errorf("expected deployed, got %s", status)
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 node statuses, got %d", len(nodes))
	}
	for _, ns := range nodes {
		if ns.Status != docker.StatusRunning {
			t.Errorf("expected running status for node %s, got %s", ns.NodeID, ns.Status)
		}
	}
}

func TestApplyDiff_AddedNodes(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	_ = mgr.Deploy(ctx, testDiagram())

	orch.createdConfigs = nil
	orch.startedIDs = nil

	added := []model.DiagramNode{
		{ID: "n3", Type: model.ServiceTypePostgreSQL, Name: "pg-1"},
	}

	if err := mgr.ApplyDiff(ctx, added, nil, nil); err != nil {
		t.Fatalf("ApplyDiff failed: %v", err)
	}

	if len(orch.createdConfigs) != 1 {
		t.Errorf("expected 1 container created, got %d", len(orch.createdConfigs))
	}
	if len(orch.startedIDs) != 1 {
		t.Errorf("expected 1 container started, got %d", len(orch.startedIDs))
	}
}

func TestApplyDiff_RemovedNodes(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	_ = mgr.Deploy(ctx, testDiagram())

	removed := []model.DiagramNode{
		{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1"},
	}

	if err := mgr.ApplyDiff(ctx, nil, removed, nil); err != nil {
		t.Fatalf("ApplyDiff failed: %v", err)
	}

	if len(orch.stoppedIDs) != 1 {
		t.Errorf("expected 1 container stopped, got %d", len(orch.stoppedIDs))
	}
	if len(orch.removedIDs) != 1 {
		t.Errorf("expected 1 container removed, got %d", len(orch.removedIDs))
	}
}

func TestApplyDiff_NotDeployed(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)

	err := mgr.ApplyDiff(context.Background(), nil, nil, nil)
	if !errors.Is(err, ErrNotDeployed) {
		t.Errorf("expected ErrNotDeployed, got %v", err)
	}
}

func TestNodeIDForContainer(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	_ = mgr.Deploy(ctx, testDiagram())

	// We know the mock generates container IDs as "container-<name>".
	nodeID, ok := mgr.NodeIDForContainer("container-redis-1")
	if !ok {
		t.Fatal("expected to find node for container-redis-1")
	}
	if nodeID != "n1" {
		t.Errorf("expected node n1, got %s", nodeID)
	}

	_, ok = mgr.NodeIDForContainer("nonexistent")
	if ok {
		t.Error("expected not to find node for nonexistent container")
	}
}

