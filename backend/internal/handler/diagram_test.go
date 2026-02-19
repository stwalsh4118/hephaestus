package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
	"github.com/stwalsh4118/hephaestus/backend/internal/storage"
)

func setupTest(t *testing.T) (*http.ServeMux, *storage.FileStore) {
	t.Helper()
	dir := t.TempDir()
	store, err := storage.NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}
	h := NewDiagramHandler(store)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux, store
}

func validDiagramJSON() []byte {
	d := model.Diagram{
		ID:   "placeholder",
		Name: "Test Diagram",
		Nodes: []model.DiagramNode{
			{
				ID:       "n1",
				Type:     model.ServiceTypeAPIService,
				Name:     "My API",
				Position: &model.Position{X: 100, Y: 200},
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
	data, _ := json.Marshal(d)
	return data
}

func createDiagram(t *testing.T, mux *http.ServeMux) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/diagrams", bytes.NewReader(validDiagramJSON()))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/diagrams: got %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}
	return resp["id"]
}

func TestCreate_Valid(t *testing.T) {
	mux, _ := setupTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/diagrams", bytes.NewReader(validDiagramJSON()))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusCreated)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != contentTypeJSON {
		t.Errorf("Content-Type: got %q, want %q", ct, contentTypeJSON)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["id"] == "" {
		t.Error("expected non-empty id in response")
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	mux, _ := setupTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/diagrams", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error == "" {
		t.Error("expected error message")
	}
}

func TestCreate_ValidationError(t *testing.T) {
	mux, _ := setupTest(t)

	invalid := model.Diagram{Name: "test"}
	data, _ := json.Marshal(invalid)

	req := httptest.NewRequest(http.MethodPost, "/api/diagrams", bytes.NewReader(data))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGet_Existing(t *testing.T) {
	mux, _ := setupTest(t)
	id := createDiagram(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/diagrams/"+id, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != contentTypeJSON {
		t.Errorf("Content-Type: got %q, want %q", ct, contentTypeJSON)
	}

	var d model.Diagram
	if err := json.Unmarshal(rec.Body.Bytes(), &d); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if d.ID != id {
		t.Errorf("ID: got %q, want %q", d.ID, id)
	}
	if d.Name != "Test Diagram" {
		t.Errorf("Name: got %q, want %q", d.Name, "Test Diagram")
	}
}

func TestGet_NotFound(t *testing.T) {
	mux, _ := setupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/diagrams/nonexistent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != contentTypeJSON {
		t.Errorf("Content-Type: got %q, want %q", ct, contentTypeJSON)
	}
}

func TestUpdate_Valid(t *testing.T) {
	mux, _ := setupTest(t)
	id := createDiagram(t, mux)

	updated := model.Diagram{
		ID:   id,
		Name: "Updated Diagram",
		Nodes: []model.DiagramNode{
			{
				ID:       "n1",
				Type:     model.ServiceTypeRedis,
				Name:     "Cache",
				Position: &model.Position{X: 50, Y: 50},
			},
		},
		Edges: []model.DiagramEdge{},
	}
	data, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/api/diagrams/"+id, bytes.NewReader(data))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var d model.Diagram
	if err := json.Unmarshal(rec.Body.Bytes(), &d); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if d.Name != "Updated Diagram" {
		t.Errorf("Name: got %q, want %q", d.Name, "Updated Diagram")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	mux, _ := setupTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/diagrams/nonexistent", bytes.NewReader(validDiagramJSON()))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestUpdate_InvalidBody(t *testing.T) {
	mux, _ := setupTest(t)
	id := createDiagram(t, mux)

	req := httptest.NewRequest(http.MethodPut, "/api/diagrams/"+id, bytes.NewReader([]byte("bad")))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAllResponses_HaveJSONContentType(t *testing.T) {
	mux, _ := setupTest(t)

	tests := []struct {
		name   string
		method string
		path   string
		body   []byte
	}{
		{"POST valid", http.MethodPost, "/api/diagrams", validDiagramJSON()},
		{"POST invalid", http.MethodPost, "/api/diagrams", []byte("bad")},
		{"GET not found", http.MethodGet, "/api/diagrams/missing", nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body *bytes.Reader
			if tc.body != nil {
				body = bytes.NewReader(tc.body)
			}
			var req *http.Request
			if body != nil {
				req = httptest.NewRequest(tc.method, tc.path, body)
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			ct := rec.Header().Get("Content-Type")
			if ct != contentTypeJSON {
				t.Errorf("Content-Type: got %q, want %q", ct, contentTypeJSON)
			}
		})
	}
}
