package templates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
	"github.com/stwalsh4118/hephaestus/backend/internal/openapi"
)

// Prism container configuration constants.
const (
	// specDir is the host-side directory where generated OpenAPI specs are written.
	specDir = "heph-specs"
	// containerSpecPath is the path inside the Prism container where the spec is mounted.
	containerSpecPath = "/tmp/spec.json"
)

// newPrismCmd returns the command passed to the Prism container to serve a mounted spec.
func newPrismCmd() []string {
	return []string{"mock", "-h", "0.0.0.0", containerSpecPath}
}

// APIServiceTemplate builds a ContainerConfig for API Service (Prism) nodes.
type APIServiceTemplate struct{}

// Build creates a docker.ContainerConfig for an API service node.
// It parses endpoint config, generates an OpenAPI spec, writes it to disk,
// and mounts it into the Prism container.
func (t *APIServiceTemplate) Build(node model.DiagramNode, hostPort string, _ ...string) (docker.ContainerConfig, error) {
	hostname := SanitizeName(node.Name)

	endpoints, err := parseEndpoints(node)
	if err != nil {
		return docker.ContainerConfig{}, err
	}

	specBytes, err := openapi.GenerateSpec(endpoints, node.Name)
	if err != nil {
		return docker.ContainerConfig{}, fmt.Errorf("generate openapi spec for node %q: %w", node.ID, err)
	}

	hostSpecPath, err := writeSpecFile(hostname, specBytes)
	if err != nil {
		return docker.ContainerConfig{}, fmt.Errorf("write spec file for node %q: %w", node.ID, err)
	}

	return docker.ContainerConfig{
		Image:       ImageAPIService,
		Name:        hostname,
		Cmd:         newPrismCmd(),
		Env:         map[string]string{},
		Ports:       map[string]string{hostPort: PortAPIService},
		Volumes:     map[string]string{hostSpecPath: containerSpecPath},
		Hostname:    hostname,
		NetworkName: docker.NetworkName,
	}, nil
}

// parseEndpoints extracts the endpoint definitions from a node's config JSON.
// Returns nil endpoints if the config is empty. Returns an error if the config
// is present but cannot be parsed.
func parseEndpoints(node model.DiagramNode) ([]model.Endpoint, error) {
	if len(node.Config) == 0 {
		return nil, nil
	}

	var cfg model.ApiServiceConfig
	if err := json.Unmarshal(node.Config, &cfg); err != nil {
		return nil, fmt.Errorf("parse api-service config for node %q: %w", node.ID, err)
	}

	return cfg.Endpoints, nil
}

// writeSpecFile writes the OpenAPI spec bytes to a temp file on the host.
// Returns the absolute path to the written file.
func writeSpecFile(containerName string, specBytes []byte) (string, error) {
	dir := filepath.Join(os.TempDir(), specDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create spec directory %q: %w", dir, err)
	}

	specPath := filepath.Join(dir, containerName+".json")
	if err := os.WriteFile(specPath, specBytes, 0o644); err != nil {
		return "", fmt.Errorf("write spec file %q: %w", specPath, err)
	}

	return specPath, nil
}
