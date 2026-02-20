package templates

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// PostgreSQLTemplate builds a ContainerConfig for PostgreSQL nodes.
type PostgreSQLTemplate struct{}

// Build creates a docker.ContainerConfig for a PostgreSQL service node.
func (t *PostgreSQLTemplate) Build(node model.DiagramNode, hostPort string, _ ...string) (docker.ContainerConfig, error) {
	hostname := SanitizeName(node.Name)

	env := DefaultPostgresEnv()

	// Validate config JSON if present. No overridable fields yet —
	// PostgresqlConfig only contains Type/Engine/Version which don't
	// affect the container configuration for MVP.
	if len(node.Config) > 0 {
		var cfg model.PostgresqlConfig
		if err := json.Unmarshal(node.Config, &cfg); err != nil {
			return docker.ContainerConfig{}, fmt.Errorf("parse postgresql config for node %q: %w", node.ID, err)
		}
		_ = cfg // validated but no overridable fields yet
	}

	return docker.ContainerConfig{
		Image:       ImagePostgreSQL,
		Name:        hostname,
		Env:         env,
		Ports:       map[string]string{hostPort: PortPostgreSQL},
		Hostname:    hostname,
		NetworkName: docker.NetworkName,
	}, nil
}

// SanitizeName converts a node name into a valid container name:
// lowercase, spaces replaced with hyphens, non-alphanumeric chars removed.
func SanitizeName(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
