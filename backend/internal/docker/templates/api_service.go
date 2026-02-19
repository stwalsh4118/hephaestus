package templates

import (
	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// APIServiceTemplate builds a ContainerConfig for API Service (Prism) nodes.
type APIServiceTemplate struct{}

// Build creates a docker.ContainerConfig for an API service node.
// OpenAPI spec mounting is deferred to PBI 8.
func (t *APIServiceTemplate) Build(node model.DiagramNode, hostPort string, _ ...string) (docker.ContainerConfig, error) {
	hostname := sanitizeName(node.Name)

	return docker.ContainerConfig{
		Image:       ImageAPIService,
		Name:        hostname,
		Env:         map[string]string{},
		Ports:       map[string]string{hostPort: PortAPIService},
		Hostname:    hostname,
		NetworkName: docker.NetworkName,
	}, nil
}
