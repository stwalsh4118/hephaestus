package templates

import (
	"encoding/json"
	"fmt"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// RedisTemplate builds a ContainerConfig for Redis nodes.
type RedisTemplate struct{}

// Build creates a docker.ContainerConfig for a Redis service node.
func (t *RedisTemplate) Build(node model.DiagramNode, hostPort string, _ ...string) (docker.ContainerConfig, error) {
	hostname := SanitizeName(node.Name)

	env := map[string]string{}

	if len(node.Config) > 0 {
		var cfg model.RedisConfig
		if err := json.Unmarshal(node.Config, &cfg); err != nil {
			return docker.ContainerConfig{}, fmt.Errorf("parse redis config for node %q: %w", node.ID, err)
		}
		if cfg.MaxMemory != "" {
			env["REDIS_MAXMEMORY"] = cfg.MaxMemory
		}
		if cfg.EvictionPolicy != "" {
			env["REDIS_EVICTION_POLICY"] = cfg.EvictionPolicy
		}
	}

	return docker.ContainerConfig{
		Image:       ImageRedis,
		Name:        hostname,
		Env:         env,
		Ports:       map[string]string{hostPort: PortRedis},
		Hostname:    hostname,
		NetworkName: docker.NetworkName,
	}, nil
}
