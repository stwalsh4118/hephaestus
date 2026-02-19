package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
	"github.com/stwalsh4118/hephaestus/backend/internal/storage"
)

// DiagramHandler provides HTTP handlers for diagram CRUD operations.
type DiagramHandler struct {
	store storage.DiagramStore
}

// NewDiagramHandler creates a DiagramHandler backed by the given store.
func NewDiagramHandler(store storage.DiagramStore) *DiagramHandler {
	return &DiagramHandler{store: store}
}

// RegisterRoutes registers diagram CRUD routes on the given mux.
func (h *DiagramHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/diagrams", h.Create)
	mux.HandleFunc("GET /api/diagrams/{id}", h.Get)
	mux.HandleFunc("PUT /api/diagrams/{id}", h.Update)
}

// Create handles POST /api/diagrams.
func (h *DiagramHandler) Create(w http.ResponseWriter, r *http.Request) {
	var d model.Diagram
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := model.ValidateDiagram(&d); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.store.Create(&d)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create diagram")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": created.ID})
}

// Get handles GET /api/diagrams/{id}.
func (h *DiagramHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	d, err := h.store.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "diagram not found")
			return
		}
		if errors.Is(err, storage.ErrInvalidID) {
			writeError(w, http.StatusBadRequest, "invalid diagram ID")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to retrieve diagram")
		return
	}

	writeJSON(w, http.StatusOK, d)
}

// Update handles PUT /api/diagrams/{id}.
func (h *DiagramHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var d model.Diagram
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := model.ValidateDiagram(&d); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.store.Update(id, &d)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "diagram not found")
			return
		}
		if errors.Is(err, storage.ErrInvalidID) {
			writeError(w, http.StatusBadRequest, "invalid diagram ID")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update diagram")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}
