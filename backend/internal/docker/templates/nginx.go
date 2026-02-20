package templates

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// NginxTemplate builds a ContainerConfig for Nginx nodes.
type NginxTemplate struct{}

// Build creates a docker.ContainerConfig for an Nginx service node.
func (t *NginxTemplate) Build(node model.DiagramNode, hostPort string, _ ...string) (docker.ContainerConfig, error) {
	hostname := SanitizeName(node.Name)

	env := map[string]string{}

	if len(node.Config) > 0 {
		var cfg model.NginxConfig
		if err := json.Unmarshal(node.Config, &cfg); err != nil {
			return docker.ContainerConfig{}, fmt.Errorf("parse nginx config for node %q: %w", node.ID, err)
		}
		// Store upstream servers as comma-separated env var for runtime config.
		if len(cfg.UpstreamServers) > 0 {
			env["NGINX_UPSTREAMS"] = strings.Join(cfg.UpstreamServers, ",")
		}
	}

	return docker.ContainerConfig{
		Image:       ImageNginx,
		Name:        hostname,
		Env:         env,
		Ports:       map[string]string{hostPort: PortNginx},
		Hostname:    hostname,
		NetworkName: docker.NetworkName,
	}, nil
}
