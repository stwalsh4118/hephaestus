package templates

import (
	"encoding/json"
	"fmt"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// RabbitMQTemplate builds a ContainerConfig for RabbitMQ nodes.
type RabbitMQTemplate struct{}

// Build creates a docker.ContainerConfig for a RabbitMQ service node.
// RabbitMQ requires two host ports: the first for AMQP (5672), the second
// for the management UI (15672). The management port is passed via hostPorts[0].
func (t *RabbitMQTemplate) Build(node model.DiagramNode, hostPort string, hostPorts ...string) (docker.ContainerConfig, error) {
	hostname := SanitizeName(node.Name)

	env := map[string]string{}
	vhost := "/"

	if len(node.Config) > 0 {
		var cfg model.RabbitMQConfig
		if err := json.Unmarshal(node.Config, &cfg); err != nil {
			return docker.ContainerConfig{}, fmt.Errorf("parse rabbitmq config for node %q: %w", node.ID, err)
		}
		if cfg.Vhost != "" {
			vhost = cfg.Vhost
		}
	}
	env["RABBITMQ_DEFAULT_VHOST"] = vhost

	ports := map[string]string{
		hostPort: PortRabbitMQAMQP,
	}
	if len(hostPorts) > 0 {
		ports[hostPorts[0]] = PortRabbitMQManagement
	}

	return docker.ContainerConfig{
		Image:       ImageRabbitMQ,
		Name:        hostname,
		Env:         env,
		Ports:       ports,
		Hostname:    hostname,
		NetworkName: docker.NetworkName,
	}, nil
}
