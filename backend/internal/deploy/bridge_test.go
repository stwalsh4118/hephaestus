package deploy

import (
	"context"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

func TestBuildStatusMessage(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	diagram := model.Diagram{
		ID:   "test",
		Name: "test",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1"},
			{ID: "n2", Type: model.ServiceTypeNginx, Name: "nginx-1"},
		},
		Edges: []model.DiagramEdge{
			{ID: "e1", Source: "n2", Target: "n1", Label: "connects"},
		},
	}

	_ = mgr.Deploy(ctx, diagram)

	containerStatuses := map[string]docker.ContainerStatus{
		"container-redis-1": docker.StatusRunning,
		"container-nginx-1": docker.StatusHealthy,
	}

	msg := mgr.BuildStatusMessage(containerStatuses)

	if msg.Type != StatusMessageType {
		t.Errorf("expected type %s, got %s", StatusMessageType, msg.Type)
	}
	if msg.DeployStatus != StatusDeployed {
		t.Errorf("expected deployed, got %s", msg.DeployStatus)
	}
	if len(msg.NodeStatuses) != 2 {
		t.Fatalf("expected 2 node statuses, got %d", len(msg.NodeStatuses))
	}

	statusByNode := make(map[string]docker.ContainerStatus)
	for _, ns := range msg.NodeStatuses {
		statusByNode[ns.NodeID] = ns.Status
	}

	if statusByNode["n1"] != docker.StatusRunning {
		t.Errorf("expected n1 running, got %s", statusByNode["n1"])
	}
	if statusByNode["n2"] != docker.StatusHealthy {
		t.Errorf("expected n2 healthy, got %s", statusByNode["n2"])
	}
}

func TestBuildStatusMessage_MissingContainer(t *testing.T) {
	orch := newMockOrchestrator()
	mgr := NewDeploymentManager(orch)
	ctx := context.Background()

	diagram := model.Diagram{
		ID:   "test",
		Name: "test",
		Nodes: []model.DiagramNode{
			{ID: "n1", Type: model.ServiceTypeRedis, Name: "redis-1"},
		},
		Edges: []model.DiagramEdge{},
	}

	_ = mgr.Deploy(ctx, diagram)

	// Empty statuses map — container not found.
	msg := mgr.BuildStatusMessage(map[string]docker.ContainerStatus{})

	if len(msg.NodeStatuses) != 1 {
		t.Fatalf("expected 1 node status, got %d", len(msg.NodeStatuses))
	}
	if msg.NodeStatuses[0].Status != docker.StatusError {
		t.Errorf("expected error status for missing container, got %s", msg.NodeStatuses[0].Status)
	}
}
