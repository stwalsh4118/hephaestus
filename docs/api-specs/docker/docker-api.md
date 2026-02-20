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
    Cmd         []string          `json:"cmd,omitempty"`         // override image CMD
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

---

## Service-to-Container Mapping

Package: `backend/internal/docker/templates`

### Interfaces

```go
type ContainerTemplate interface {
    Build(node model.DiagramNode, hostPort string, hostPorts ...string) (docker.ContainerConfig, error)
}
```

### Types

```go
type TemplateRegistry map[string]ContainerTemplate

type PortAllocator struct { /* thread-safe port allocator */ }

type Translator struct { /* not safe for concurrent use */ }
```

### Image Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `ImageAPIService` | `"stoplight/prism:latest"` | API Service (Prism mock server) |
| `ImagePostgreSQL` | `"postgres:16"` | PostgreSQL database |
| `ImageRedis` | `"redis:7"` | Redis cache |
| `ImageNginx` | `"nginx:latest"` | Nginx web server |
| `ImageRabbitMQ` | `"rabbitmq:3-management"` | RabbitMQ message broker |

### Port Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `PortAPIService` | `"4010"` | Prism default port |
| `PortPostgreSQL` | `"5432"` | PostgreSQL default port |
| `PortRedis` | `"6379"` | Redis default port |
| `PortNginx` | `"80"` | Nginx HTTP port |
| `PortRabbitMQAMQP` | `"5672"` | RabbitMQ AMQP port |
| `PortRabbitMQManagement` | `"15672"` | RabbitMQ management UI port |
| `DefaultMinPort` | `10000` | Host port allocation range start |
| `DefaultMaxPort` | `19999` | Host port allocation range end |

### Priority Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `PriorityInfrastructure` | `0` | postgresql, redis, rabbitmq (start first) |
| `PriorityApplication` | `1` | nginx, api-service |

### Constructors

```go
func NewRegistry() TemplateRegistry
func NewPortAllocator(minPort, maxPort int) *PortAllocator
func NewTranslator() *Translator
func DefaultPostgresEnv() map[string]string
```

### Template Implementations

| Service Type | Struct | File |
|-------------|--------|------|
| `api-service` | `APIServiceTemplate` | `api_service.go` |
| `postgresql` | `PostgreSQLTemplate` | `postgresql.go` |
| `redis` | `RedisTemplate` | `redis.go` |
| `nginx` | `NginxTemplate` | `nginx.go` |
| `rabbitmq` | `RabbitMQTemplate` | `rabbitmq.go` |

### PortAllocator Methods

```go
func (a *PortAllocator) Allocate() (string, error)
func (a *PortAllocator) AllocateN(n int) ([]string, error)  // atomic with rollback
func (a *PortAllocator) Reset()
```

### Dependency Resolver

```go
func ResolveDependencies(nodes []model.DiagramNode, edges []model.DiagramEdge) ([]string, error)
```

### Translator Method

```go
func (t *Translator) Translate(diagram model.Diagram) ([]docker.ContainerConfig, error)
```

---

## OpenAPI Spec Generator

Package: `backend/internal/openapi`

### Functions

```go
func GenerateSpec(endpoints []model.Endpoint, title string) ([]byte, error)
```

Converts a slice of endpoint definitions into a valid OpenAPI 3.0.0 JSON document.

- Groups endpoints by path (multiple methods on the same path share a path item)
- Validates HTTP methods (GET, POST, PUT, DELETE, PATCH)
- Parses `responseSchema`: valid JSON object → use directly, empty → `{"type":"object"}`, invalid JSON → wrap as string example
- Returns indented JSON bytes

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `openAPIVersion` | `"3.0.0"` | OpenAPI specification version |
| `contentTypeJSON` | `"application/json"` | Response content type |
