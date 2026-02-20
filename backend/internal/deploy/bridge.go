package deploy

import "github.com/stwalsh4118/hephaestus/backend/internal/docker"

// BuildStatusMessage converts a map of container ID → ContainerStatus into a
// StatusMessage with node IDs, using the DeploymentManager's internal mapping.
func (m *DeploymentManager) BuildStatusMessage(containerStatuses map[string]docker.ContainerStatus) StatusMessage {
	m.mu.Lock()
	status := m.status
	nc := make(map[string]string, len(m.containerNodes))
	for k, v := range m.containerNodes {
		nc[k] = v
	}
	m.mu.Unlock()

	nodeStatuses := make([]NodeStatus, 0, len(nc))
	for containerID, nodeID := range nc {
		containerStatus, ok := containerStatuses[containerID]
		if !ok {
			containerStatus = docker.StatusError
		}
		nodeStatuses = append(nodeStatuses, NodeStatus{
			NodeID:      nodeID,
			ContainerID: containerID,
			Status:      containerStatus,
		})
	}

	return StatusMessage{
		Type:         StatusMessageType,
		DeployStatus: status,
		NodeStatuses: nodeStatuses,
	}
}
