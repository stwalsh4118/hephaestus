package storage

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

func sampleDiagram() *model.Diagram {
	return &model.Diagram{
		Name: "Test Diagram",
		Nodes: []model.DiagramNode{
			{
				ID:       "n1",
				Type:     model.ServiceTypeAPIService,
				Name:     "API",
				Position: &model.Position{X: 10, Y: 20},
				Config:   json.RawMessage(`{"type":"api-service","endpoints":[],"port":8080}`),
			},
		},
		Edges: []model.DiagramEdge{
			{
				ID:     "e1",
				Source: "n1",
				Target: "n2",
				Label:  "calls",
			},
		},
	}
}

func newTestStore(t *testing.T) *FileStore {
	t.Helper()
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}
	return store
}

func TestCreate_ReturnsGeneratedUUID(t *testing.T) {
	store := newTestStore(t)
	d := sampleDiagram()

	created, err := store.Create(d)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == "" {
		t.Error("expected non-empty ID")
	}
	if created.Name != "Test Diagram" {
		t.Errorf("Name: got %q, want %q", created.Name, "Test Diagram")
	}
}

func TestGet_ReturnsCreatedDiagram(t *testing.T) {
	store := newTestStore(t)
	created, err := store.Create(sampleDiagram())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := store.Get(created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID: got %q, want %q", got.ID, created.ID)
	}
	if got.Name != created.Name {
		t.Errorf("Name: got %q, want %q", got.Name, created.Name)
	}
	if len(got.Nodes) != 1 {
		t.Fatalf("Nodes: got %d, want 1", len(got.Nodes))
	}
	if got.Nodes[0].ID != "n1" {
		t.Errorf("Node.ID: got %q, want %q", got.Nodes[0].ID, "n1")
	}
}

func TestUpdate_PersistsChanges(t *testing.T) {
	store := newTestStore(t)
	created, err := store.Create(sampleDiagram())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	updated := sampleDiagram()
	updated.Name = "Updated Diagram"

	result, err := store.Update(created.ID, updated)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if result.Name != "Updated Diagram" {
		t.Errorf("Name: got %q, want %q", result.Name, "Updated Diagram")
	}

	got, err := store.Get(created.ID)
	if err != nil {
		t.Fatalf("Get after update: %v", err)
	}
	if got.Name != "Updated Diagram" {
		t.Errorf("Name after re-read: got %q, want %q", got.Name, "Updated Diagram")
	}
}

func TestGet_NotFound(t *testing.T) {
	store := newTestStore(t)
	_, err := store.Get("nonexistent-id")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	store := newTestStore(t)
	_, err := store.Update("nonexistent-id", sampleDiagram())
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestPersistence_AcrossInstances(t *testing.T) {
	dir := t.TempDir()

	store1, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore (1): %v", err)
	}
	created, err := store1.Create(sampleDiagram())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Simulate server restart: create a new store pointing to the same directory
	store2, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore (2): %v", err)
	}
	got, err := store2.Get(created.ID)
	if err != nil {
		t.Fatalf("Get from new instance: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID: got %q, want %q", got.ID, created.ID)
	}
	if got.Name != created.Name {
		t.Errorf("Name: got %q, want %q", got.Name, created.Name)
	}
}

func TestConcurrentCreate(t *testing.T) {
	store := newTestStore(t)
	const n = 20

	var wg sync.WaitGroup
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Create(sampleDiagram())
			if err != nil {
				errs <- err
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent Create error: %v", err)
	}
}

func TestNewFileStore_NestedDir(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir + "/sub/nested")
	if err != nil {
		t.Fatalf("NewFileStore with nested dir: %v", err)
	}
	if store.dir != dir+"/sub/nested" {
		t.Errorf("dir: got %q, want %q", store.dir, dir+"/sub/nested")
	}
}

func TestCreate_DoesNotMutateInput(t *testing.T) {
	store := newTestStore(t)
	d := sampleDiagram()
	originalID := d.ID

	_, err := store.Create(d)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if d.ID != originalID {
		t.Errorf("input diagram ID was mutated: got %q, want %q", d.ID, originalID)
	}
}

func TestGet_PathTraversal(t *testing.T) {
	store := newTestStore(t)
	maliciousIDs := []string{
		"../etc/passwd",
		"../../etc/shadow",
		"foo/bar",
		"",
	}
	for _, id := range maliciousIDs {
		_, err := store.Get(id)
		if !errors.Is(err, ErrInvalidID) {
			t.Errorf("Get(%q): expected ErrInvalidID, got: %v", id, err)
		}
	}
}

func TestUpdate_PathTraversal(t *testing.T) {
	store := newTestStore(t)
	_, err := store.Update("../etc/passwd", sampleDiagram())
	if !errors.Is(err, ErrInvalidID) {
		t.Errorf("expected ErrInvalidID, got: %v", err)
	}
}
