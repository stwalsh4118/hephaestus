package deploy

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/docker/templates"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// Sentinel errors for deploy operations.
var (
	ErrAlreadyDeploying = errors.New("deployment already in progress")
	ErrNotDeployed      = errors.New("no active deployment")
)

// DeploymentManager coordinates deploying diagrams to Docker containers,
// tracking the mapping between diagram nodes and container IDs.
type DeploymentManager struct {
	orchestrator docker.Orchestrator

	mu             sync.Mutex
	status         DeployStatus
	nodeContainers map[string]string // node ID → container ID
	containerNodes map[string]string // container ID → node ID
	lastDiagram    *model.Diagram
}

// NewDeploymentManager creates a DeploymentManager with the given dependencies.
func NewDeploymentManager(orchestrator docker.Orchestrator) *DeploymentManager {
	return &DeploymentManager{
		orchestrator:   orchestrator,
		status:         StatusIdle,
		nodeContainers: make(map[string]string),
		containerNodes: make(map[string]string),
	}
}

// Deploy translates the diagram into containers and starts them.
// Returns ErrAlreadyDeploying if a deployment is already in progress.
func (m *DeploymentManager) Deploy(ctx context.Context, diagram model.Diagram) error {
	m.mu.Lock()
	if m.status == StatusDeploying || m.status == StatusDeployed {
		m.mu.Unlock()
		return ErrAlreadyDeploying
	}
	m.status = StatusDeploying
	m.nodeContainers = make(map[string]string)
	m.containerNodes = make(map[string]string)
	m.mu.Unlock()

	// Create a fresh translator per call to avoid data races on the internal
	// PortAllocator (Translator is not safe for concurrent use).
	translator := templates.NewTranslator()
	configs, err := translator.Translate(diagram)
	if err != nil {
		m.mu.Lock()
		m.status = StatusError
		m.mu.Unlock()
		return fmt.Errorf("translate diagram: %w", err)
	}

	if err := m.orchestrator.CreateNetwork(ctx); err != nil {
		m.mu.Lock()
		m.status = StatusError
		m.mu.Unlock()
		return fmt.Errorf("create network: %w", err)
	}

	// Build a map from sanitized container name to node ID for correlation.
	// Templates use SanitizeName(node.Name) as the container config name.
	nodeByName := make(map[string]string, len(diagram.Nodes))
	for _, node := range diagram.Nodes {
		nodeByName[templates.SanitizeName(node.Name)] = node.ID
	}

	for _, cfg := range configs {
		containerID, err := m.orchestrator.CreateContainer(ctx, cfg)
		if err != nil {
			m.mu.Lock()
			m.status = StatusError
			m.mu.Unlock()
			return fmt.Errorf("create container %q: %w", cfg.Name, err)
		}

		if err := m.orchestrator.StartContainer(ctx, containerID); err != nil {
			m.mu.Lock()
			m.status = StatusError
			m.mu.Unlock()
			return fmt.Errorf("start container %q: %w", cfg.Name, err)
		}

		// Track the mapping immediately so partial state is visible on error.
		if nodeID, ok := nodeByName[cfg.Name]; ok {
			m.mu.Lock()
			m.nodeContainers[nodeID] = containerID
			m.containerNodes[containerID] = nodeID
			m.mu.Unlock()
		}
	}

	m.mu.Lock()
	m.status = StatusDeployed
	m.lastDiagram = &diagram
	m.mu.Unlock()

	return nil
}

// Teardown stops and removes all deployed containers, resets state to idle.
// Returns ErrNotDeployed if there is no active deployment.
func (m *DeploymentManager) Teardown(ctx context.Context) error {
	m.mu.Lock()
	if m.status != StatusDeployed && m.status != StatusError {
		m.mu.Unlock()
		return ErrNotDeployed
	}
	m.status = StatusTearingDown
	m.mu.Unlock()

	err := m.orchestrator.TeardownAll(ctx)

	m.mu.Lock()
	m.status = StatusIdle
	m.nodeContainers = make(map[string]string)
	m.containerNodes = make(map[string]string)
	m.lastDiagram = nil
	m.mu.Unlock()

	if err != nil {
		return fmt.Errorf("teardown: %w", err)
	}
	return nil
}

// GetStatus returns the current deployment status and per-node statuses.
func (m *DeploymentManager) GetStatus(ctx context.Context) (DeployStatus, []NodeStatus, error) {
	m.mu.Lock()
	status := m.status
	nc := make(map[string]string, len(m.nodeContainers))
	for k, v := range m.nodeContainers {
		nc[k] = v
	}
	m.mu.Unlock()

	if status == StatusIdle || len(nc) == 0 {
		return status, nil, nil
	}

	nodeStatuses := make([]NodeStatus, 0, len(nc))
	for nodeID, containerID := range nc {
		containerStatus, err := m.orchestrator.HealthCheck(ctx, containerID)
		if err != nil {
			containerStatus = docker.StatusError
		}
		nodeStatuses = append(nodeStatuses, NodeStatus{
			NodeID:      nodeID,
			ContainerID: containerID,
			Status:      containerStatus,
		})
	}

	return status, nodeStatuses, nil
}

// ApplyDiff deploys new containers for added nodes and removes containers for
// removed nodes, without touching unchanged containers.
func (m *DeploymentManager) ApplyDiff(ctx context.Context, added []model.DiagramNode, removed []model.DiagramNode, edges []model.DiagramEdge) error {
	m.mu.Lock()
	if m.status != StatusDeployed {
		m.mu.Unlock()
		return ErrNotDeployed
	}
	m.mu.Unlock()

	// Remove containers for removed nodes.
	for _, node := range removed {
		m.mu.Lock()
		containerID, exists := m.nodeContainers[node.ID]
		m.mu.Unlock()

		if !exists {
			continue
		}

		if err := m.orchestrator.StopContainer(ctx, containerID); err != nil {
			return fmt.Errorf("stop container for node %q: %w", node.ID, err)
		}
		if err := m.orchestrator.RemoveContainer(ctx, containerID); err != nil {
			return fmt.Errorf("remove container for node %q: %w", node.ID, err)
		}

		m.mu.Lock()
		delete(m.nodeContainers, node.ID)
		delete(m.containerNodes, containerID)
		m.mu.Unlock()
	}

	// Deploy containers for added nodes.
	if len(added) > 0 {
		addedDiagram := model.Diagram{
			Nodes: added,
			Edges: edges,
		}

		translator := templates.NewTranslator()
		configs, err := translator.Translate(addedDiagram)
		if err != nil {
			return fmt.Errorf("translate added nodes: %w", err)
		}

		nodeByName := make(map[string]string, len(added))
		for _, node := range added {
			nodeByName[templates.SanitizeName(node.Name)] = node.ID
		}

		for _, cfg := range configs {
			containerID, err := m.orchestrator.CreateContainer(ctx, cfg)
			if err != nil {
				return fmt.Errorf("create container for added node %q: %w", cfg.Name, err)
			}

			if err := m.orchestrator.StartContainer(ctx, containerID); err != nil {
				return fmt.Errorf("start container for added node %q: %w", cfg.Name, err)
			}

			if nodeID, ok := nodeByName[cfg.Name]; ok {
				m.mu.Lock()
				m.nodeContainers[nodeID] = containerID
				m.containerNodes[containerID] = nodeID
				m.mu.Unlock()
			}
		}
	}

	return nil
}

// NodeIDForContainer returns the diagram node ID for a given container ID.
func (m *DeploymentManager) NodeIDForContainer(containerID string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	nodeID, ok := m.containerNodes[containerID]
	return nodeID, ok
}

// Status returns the current deploy status.
func (m *DeploymentManager) Status() DeployStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

// LastDiagram returns the last deployed diagram, or nil if none.
func (m *DeploymentManager) LastDiagram() *model.Diagram {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastDiagram
}

// UpdateLastDiagram updates the stored last-deployed diagram snapshot.
func (m *DeploymentManager) UpdateLastDiagram(diagram model.Diagram) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastDiagram = &diagram
}
