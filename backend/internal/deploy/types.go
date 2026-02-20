package deploy

import (
	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// DeployStatus represents the overall state of the deployment.
type DeployStatus string

const (
	StatusIdle       DeployStatus = "idle"
	StatusDeploying  DeployStatus = "deploying"
	StatusDeployed   DeployStatus = "deployed"
	StatusTearingDown DeployStatus = "tearing_down"
	StatusError      DeployStatus = "error"
)

// NodeStatus represents the status of a single deployed node.
type NodeStatus struct {
	NodeID      string                `json:"nodeId"`
	ContainerID string                `json:"containerId"`
	Status      docker.ContainerStatus `json:"status"`
}

// StatusMessage is the WebSocket message format for status updates.
type StatusMessage struct {
	Type         string       `json:"type"`
	DeployStatus DeployStatus `json:"deployStatus"`
	NodeStatuses []NodeStatus `json:"nodeStatuses"`
}

// StatusResponse is the HTTP response for GET /api/deploy/status.
type StatusResponse struct {
	DeployStatus DeployStatus `json:"deployStatus"`
	NodeStatuses []NodeStatus `json:"nodeStatuses"`
}

// DeployResponse is the HTTP response for POST /api/deploy and DELETE /api/deploy.
type DeployResponse struct {
	Status string `json:"status"`
}

// DeployRequest wraps a model.Diagram for deployment.
type DeployRequest struct {
	Diagram model.Diagram `json:"diagram"`
}

// StatusMessageType is the type field value for WebSocket status messages.
const StatusMessageType = "status_update"
