package docker

import "time"

// ContainerNamePrefix is prepended to all managed container names to avoid
// conflicts with the user's other Docker containers.
const ContainerNamePrefix = "heph-"

// DefaultHealthCheckInterval is the default polling interval for container health checks.
const DefaultHealthCheckInterval = 5 * time.Second

// HealthStatusCallback is called by StartHealthPolling with the current status of all
// managed containers, keyed by container ID.
type HealthStatusCallback func(statuses map[string]ContainerStatus)

// ContainerStatus represents the current state of a managed container.
type ContainerStatus string

const (
	StatusCreated   ContainerStatus = "created"
	StatusRunning   ContainerStatus = "running"
	StatusStopped   ContainerStatus = "stopped"
	StatusError     ContainerStatus = "error"
	StatusHealthy   ContainerStatus = "healthy"
	StatusUnhealthy ContainerStatus = "unhealthy"
)

// ContainerConfig holds the configuration needed to create a container.
type ContainerConfig struct {
	Image       string            `json:"image"`
	Name        string            `json:"name"`
	Cmd         []string          `json:"cmd,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Ports       map[string]string `json:"ports,omitempty"`   // host port → container port
	Volumes     map[string]string `json:"volumes,omitempty"` // host path → container path
	Hostname    string            `json:"hostname,omitempty"`
	NetworkName string            `json:"networkName,omitempty"`
}

// ContainerInfo represents the current state of a running or stopped container.
type ContainerInfo struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Image  string            `json:"image"`
	Status ContainerStatus   `json:"status"`
	Ports  map[string]string `json:"ports,omitempty"` // host port → container port
}
