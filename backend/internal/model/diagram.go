package model

import "encoding/json"

// Service type constants matching the frontend ServiceType union.
const (
	ServiceTypeAPIService = "api-service"
	ServiceTypePostgreSQL = "postgresql"
	ServiceTypeRedis      = "redis"
	ServiceTypeNginx      = "nginx"
	ServiceTypeRabbitMQ   = "rabbitmq"
)

// ValidServiceTypes is the set of allowed service type values.
var ValidServiceTypes = map[string]bool{
	ServiceTypeAPIService: true,
	ServiceTypePostgreSQL: true,
	ServiceTypeRedis:      true,
	ServiceTypeNginx:      true,
	ServiceTypeRabbitMQ:   true,
}

// Position represents x/y coordinates on the canvas.
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// DiagramNode represents a service node in the diagram.
type DiagramNode struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Position    *Position       `json:"position"`
	Config      json.RawMessage `json:"config,omitempty"`
}

// DiagramEdge represents a connection between two nodes.
type DiagramEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
}

// Diagram is the top-level structure matching the frontend DiagramJson schema.
type Diagram struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Nodes []DiagramNode `json:"nodes"`
	Edges []DiagramEdge `json:"edges"`
}

// Endpoint represents an API service endpoint definition.
type Endpoint struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	ResponseSchema string `json:"responseSchema"`
}

// ApiServiceConfig is the configuration for api-service nodes.
type ApiServiceConfig struct {
	Type      string     `json:"type"`
	Endpoints []Endpoint `json:"endpoints"`
	Port      int        `json:"port"`
}

// PostgresqlConfig is the configuration for postgresql nodes.
type PostgresqlConfig struct {
	Type    string `json:"type"`
	Engine  string `json:"engine"`
	Version string `json:"version"`
}

// RedisConfig is the configuration for redis nodes.
type RedisConfig struct {
	Type           string `json:"type"`
	MaxMemory      string `json:"maxMemory"`
	EvictionPolicy string `json:"evictionPolicy"`
}

// NginxConfig is the configuration for nginx nodes.
type NginxConfig struct {
	Type            string   `json:"type"`
	UpstreamServers []string `json:"upstreamServers"`
}

// RabbitMQConfig is the configuration for rabbitmq nodes.
type RabbitMQConfig struct {
	Type  string `json:"type"`
	Vhost string `json:"vhost"`
}
