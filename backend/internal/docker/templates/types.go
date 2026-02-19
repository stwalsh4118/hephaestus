package templates

import (
	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// Docker image constants for each service type.
const (
	ImageAPIService = "stoplight/prism:latest"
	ImagePostgreSQL = "postgres:16"
	ImageRedis      = "redis:7"
	ImageNginx      = "nginx:latest"
	ImageRabbitMQ   = "rabbitmq:3-management"
)

// Default container-side port constants for each service type.
const (
	PortAPIService        = "4010"
	PortPostgreSQL        = "5432"
	PortRedis             = "6379"
	PortNginx             = "80"
	PortRabbitMQAMQP      = "5672"
	PortRabbitMQManagement = "15672"
)

// DefaultPostgresEnv returns a fresh copy of the default PostgreSQL environment variables.
func DefaultPostgresEnv() map[string]string {
	return map[string]string{
		"POSTGRES_USER":     "hephaestus",
		"POSTGRES_PASSWORD": "hephaestus",
		"POSTGRES_DB":       "hephaestus",
	}
}

// ContainerTemplate builds a docker.ContainerConfig from a diagram node.
// The hostPort parameter is the allocated host port; multi-port services
// receive additional ports via the hostPorts variadic parameter.
type ContainerTemplate interface {
	Build(node model.DiagramNode, hostPort string, hostPorts ...string) (docker.ContainerConfig, error)
}

// TemplateRegistry maps service type strings to their ContainerTemplate.
type TemplateRegistry map[string]ContainerTemplate

// NewRegistry returns a TemplateRegistry populated with all 5 service templates.
func NewRegistry() TemplateRegistry {
	return TemplateRegistry{
		model.ServiceTypeAPIService: &APIServiceTemplate{},
		model.ServiceTypePostgreSQL: &PostgreSQLTemplate{},
		model.ServiceTypeRedis:      &RedisTemplate{},
		model.ServiceTypeNginx:      &NginxTemplate{},
		model.ServiceTypeRabbitMQ:   &RabbitMQTemplate{},
	}
}
