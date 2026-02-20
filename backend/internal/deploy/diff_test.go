package deploy

import (
	"sort"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

func nodeIDs(nodes []model.DiagramNode) []string {
	ids := make([]string, len(nodes))
	for i, n := range nodes {
		ids[i] = n.ID
	}
	sort.Strings(ids)
	return ids
}

func TestComputeDiff_EmptyToCurrent(t *testing.T) {
	newNodes := []model.DiagramNode{
		{ID: "a", Type: "redis", Name: "redis-1"},
		{ID: "b", Type: "nginx", Name: "nginx-1"},
	}

	result := ComputeDiff(nil, newNodes)

	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
	if len(result.Unchanged) != 0 {
		t.Errorf("expected 0 unchanged, got %d", len(result.Unchanged))
	}

	addedIDs := nodeIDs(result.Added)
	if len(addedIDs) != 2 || addedIDs[0] != "a" || addedIDs[1] != "b" {
		t.Errorf("expected added [a, b], got %v", addedIDs)
	}
}

func TestComputeDiff_CurrentToEmpty(t *testing.T) {
	current := []model.DiagramNode{
		{ID: "a", Type: "redis", Name: "redis-1"},
		{ID: "b", Type: "nginx", Name: "nginx-1"},
	}

	result := ComputeDiff(current, nil)

	if len(result.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(result.Added))
	}
	if len(result.Unchanged) != 0 {
		t.Errorf("expected 0 unchanged, got %d", len(result.Unchanged))
	}

	removedIDs := nodeIDs(result.Removed)
	if len(removedIDs) != 2 || removedIDs[0] != "a" || removedIDs[1] != "b" {
		t.Errorf("expected removed [a, b], got %v", removedIDs)
	}
}

func TestComputeDiff_MixedChanges(t *testing.T) {
	current := []model.DiagramNode{
		{ID: "a", Type: "redis", Name: "redis-1"},
		{ID: "b", Type: "nginx", Name: "nginx-1"},
		{ID: "c", Type: "postgresql", Name: "pg-1"},
	}
	newNodes := []model.DiagramNode{
		{ID: "b", Type: "nginx", Name: "nginx-1"},
		{ID: "c", Type: "postgresql", Name: "pg-1"},
		{ID: "d", Type: "rabbitmq", Name: "rmq-1"},
	}

	result := ComputeDiff(current, newNodes)

	addedIDs := nodeIDs(result.Added)
	removedIDs := nodeIDs(result.Removed)
	unchangedIDs := nodeIDs(result.Unchanged)

	if len(addedIDs) != 1 || addedIDs[0] != "d" {
		t.Errorf("expected added [d], got %v", addedIDs)
	}
	if len(removedIDs) != 1 || removedIDs[0] != "a" {
		t.Errorf("expected removed [a], got %v", removedIDs)
	}
	if len(unchangedIDs) != 2 || unchangedIDs[0] != "b" || unchangedIDs[1] != "c" {
		t.Errorf("expected unchanged [b, c], got %v", unchangedIDs)
	}
}

func TestComputeDiff_SameIDDifferentConfig(t *testing.T) {
	current := []model.DiagramNode{
		{ID: "a", Type: "redis", Name: "redis-1"},
	}
	newNodes := []model.DiagramNode{
		{ID: "a", Type: "redis", Name: "redis-modified"},
	}

	result := ComputeDiff(current, newNodes)

	if len(result.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(result.Added))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
	if len(result.Unchanged) != 1 || result.Unchanged[0].ID != "a" {
		t.Errorf("expected unchanged [a], got %v", nodeIDs(result.Unchanged))
	}
}

func TestComputeDiff_BothEmpty(t *testing.T) {
	result := ComputeDiff(nil, nil)

	if len(result.Added) != 0 || len(result.Removed) != 0 || len(result.Unchanged) != 0 {
		t.Error("expected all empty for nil inputs")
	}
}
