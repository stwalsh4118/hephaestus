# Docker Orchestration API

Package: `backend/internal/docker`

## Interfaces

### `DockerAPI`
```go
type DockerAPI interface {
    Ping(ctx context.Context) error
    Close() error
}
```

### `Orchestrator`
```go
type Orchestrator interface {
    CreateContainer(ctx context.Context, config ContainerConfig) (string, error)
    StartContainer(ctx context.Context, containerID string) error
    StopContainer(ctx context.Context, containerID string) error
    RemoveContainer(ctx context.Context, containerID string) error
    ListContainers(ctx context.Context) ([]ContainerInfo, error)
    InspectContainer(ctx context.Context, containerID string) (*ContainerInfo, error)
    CreateNetwork(ctx context.Context) error
    RemoveNetwork(ctx context.Context) error
    HealthCheck(ctx context.Context, containerID string) (ContainerStatus, error)
    TeardownAll(ctx context.Context) error
}
```

## Types

```go
type ContainerConfig struct {
    Image       string            `json:"image"`
    Name        string            `json:"name"`
    Env         map[string]string `json:"env,omitempty"`
    Ports       map[string]string `json:"ports,omitempty"`       // host → container
    Volumes     map[string]string `json:"volumes,omitempty"`     // host → container
    Hostname    string            `json:"hostname,omitempty"`
    NetworkName string            `json:"networkName,omitempty"`
}

type ContainerInfo struct {
    ID     string            `json:"id"`
    Name   string            `json:"name"`
    Image  string            `json:"image"`
    Status ContainerStatus   `json:"status"`
    Ports  map[string]string `json:"ports,omitempty"`
}

type ContainerStatus string // "created" | "running" | "stopped" | "error" | "healthy" | "unhealthy"

type HealthStatusCallback func(statuses map[string]ContainerStatus)
```

## Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `ContainerNamePrefix` | `"heph-"` | Prefix for all managed container names |
| `NetworkName` | `"heph-network"` | Shared Docker bridge network name |
| `StopTimeout` | `10` | Graceful stop timeout in seconds |
| `DefaultHealthCheckInterval` | `5s` | Default polling interval for health checks |

## Constructors

```go
func NewClient() (*Client, error)
func NewDockerOrchestrator(c *Client) *DockerOrchestrator
```

## Additional Methods (on DockerOrchestrator)

```go
func (o *DockerOrchestrator) StartHealthPolling(ctx context.Context, interval time.Duration, callback HealthStatusCallback)
```
