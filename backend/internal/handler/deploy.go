package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/stwalsh4118/hephaestus/backend/internal/deploy"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// DeployHandler provides HTTP handlers for deploy operations.
type DeployHandler struct {
	manager *deploy.DeploymentManager
}

// NewDeployHandler creates a DeployHandler with the given DeploymentManager.
func NewDeployHandler(manager *deploy.DeploymentManager) *DeployHandler {
	return &DeployHandler{manager: manager}
}

// RegisterRoutes registers deploy routes on the given mux.
func (h *DeployHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/deploy", h.Deploy)
	mux.HandleFunc("PUT /api/deploy", h.Update)
	mux.HandleFunc("DELETE /api/deploy", h.Teardown)
	mux.HandleFunc("GET /api/deploy/status", h.Status)
}

// Deploy handles POST /api/deploy.
func (h *DeployHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	var diagram model.Diagram
	if err := json.NewDecoder(r.Body).Decode(&diagram); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := model.ValidateDiagram(&diagram); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.manager.Deploy(r.Context(), diagram); err != nil {
		if errors.Is(err, deploy.ErrAlreadyDeploying) {
			writeError(w, http.StatusConflict, "deployment already in progress")
			return
		}
		writeError(w, http.StatusInternalServerError, "deploy failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, deploy.DeployResponse{Status: "deploying"})
}

// Teardown handles DELETE /api/deploy.
func (h *DeployHandler) Teardown(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.Teardown(r.Context()); err != nil {
		if errors.Is(err, deploy.ErrNotDeployed) {
			writeError(w, http.StatusConflict, "no active deployment")
			return
		}
		writeError(w, http.StatusInternalServerError, "teardown failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, deploy.DeployResponse{Status: "idle"})
}

// Update handles PUT /api/deploy — incremental topology update.
func (h *DeployHandler) Update(w http.ResponseWriter, r *http.Request) {
	var diagram model.Diagram
	if err := json.NewDecoder(r.Body).Decode(&diagram); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := model.ValidateDiagram(&diagram); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lastDiagram := h.manager.LastDiagram()
	if lastDiagram == nil {
		writeError(w, http.StatusConflict, "no active deployment")
		return
	}

	diff := deploy.ComputeDiff(lastDiagram.Nodes, diagram.Nodes)

	if err := h.manager.ApplyDiff(r.Context(), diff.Added, diff.Removed, diagram.Edges); err != nil {
		if errors.Is(err, deploy.ErrNotDeployed) {
			writeError(w, http.StatusConflict, "no active deployment")
			return
		}
		writeError(w, http.StatusInternalServerError, "update deploy failed: "+err.Error())
		return
	}

	h.manager.UpdateLastDiagram(diagram)

	deployStatus, nodeStatuses, err := h.manager.GetStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get status: "+err.Error())
		return
	}

	if nodeStatuses == nil {
		nodeStatuses = []deploy.NodeStatus{}
	}

	writeJSON(w, http.StatusOK, deploy.StatusResponse{
		DeployStatus: deployStatus,
		NodeStatuses: nodeStatuses,
	})
}

// Status handles GET /api/deploy/status.
func (h *DeployHandler) Status(w http.ResponseWriter, r *http.Request) {
	deployStatus, nodeStatuses, err := h.manager.GetStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get status: "+err.Error())
		return
	}

	if nodeStatuses == nil {
		nodeStatuses = []deploy.NodeStatus{}
	}

	writeJSON(w, http.StatusOK, deploy.StatusResponse{
		DeployStatus: deployStatus,
		NodeStatuses: nodeStatuses,
	})
}
