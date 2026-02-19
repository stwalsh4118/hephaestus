# Docker SDK for Go — Research Guide

**Date**: 2026-02-19
**Source**: https://pkg.go.dev/github.com/docker/docker/client
**Version**: v28.5.2+incompatible

## Client Initialization

```go
import "github.com/docker/docker/client"

// Create client from environment with API version negotiation
cli, err := client.NewClientWithOpts(
    client.FromEnv,
    client.WithAPIVersionNegotiation(),
)
if err != nil {
    return fmt.Errorf("create docker client: %w", err)
}
defer cli.Close()
```

**Key options**:
- `FromEnv` — reads DOCKER_HOST, DOCKER_API_VERSION, DOCKER_CERT_PATH, DOCKER_TLS_VERIFY
- `WithAPIVersionNegotiation()` — auto-negotiates API version with daemon
- `WithHost(host)` — override Docker host
- `WithTimeout(d)` — set request timeout

## Ping (Connectivity Check)

```go
ping, err := cli.Ping(ctx)
// ping.APIVersion contains the daemon API version
```

## Close

```go
err := cli.Close()
```

## Container Operations

### Create

```go
import (
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/network"
    "github.com/docker/go-connections/nat"
)

resp, err := cli.ContainerCreate(ctx,
    &container.Config{
        Image:    "nginx:latest",
        Env:      []string{"KEY=value"},
        Hostname: "my-service",
    },
    &container.HostConfig{
        PortBindings: nat.PortMap{
            "80/tcp": []nat.PortBinding{{HostPort: "8080"}},
        },
        Binds: []string{"/host/path:/container/path"},
    },
    &network.NetworkingConfig{
        EndpointsConfig: map[string]*network.EndpointSettings{
            "my-network": {},
        },
    },
    nil, // platform
    "container-name",
)
// resp.ID is the container ID
```

### Start

```go
err := cli.ContainerStart(ctx, containerID, container.StartOptions{})
```

### Stop

```go
timeout := 10 // seconds
err := cli.ContainerStop(ctx, containerID, container.StopOptions{
    Timeout: &timeout,
})
```

### Remove

```go
err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
    Force: true,
})
```

### List

```go
import "github.com/docker/docker/api/types/filters"

containers, err := cli.ContainerList(ctx, container.ListOptions{
    All: true,
    Filters: filters.NewArgs(filters.Arg("name", "heph-")),
})
for _, c := range containers {
    fmt.Printf("%s %s %s\n", c.ID[:12], c.Image, c.State)
}
```

### Inspect

```go
info, err := cli.ContainerInspect(ctx, containerID)
// info.State.Status — "created", "running", "paused", "restarting", "removing", "exited", "dead"
// info.State.Health — non-nil if container has healthcheck configured
// info.State.Health.Status — "starting", "healthy", "unhealthy"
// info.NetworkSettings.Networks — map of attached networks
```

## Network Operations

### Create Network

```go
import "github.com/docker/docker/api/types/network"

resp, err := cli.NetworkCreate(ctx, "heph-network", network.CreateOptions{
    Driver: "bridge",
})
// resp.ID is the network ID
```

### List Networks

```go
networks, err := cli.NetworkList(ctx, network.ListOptions{
    Filters: filters.NewArgs(filters.Arg("name", "heph-network")),
})
```

### Remove Network

```go
err := cli.NetworkRemove(ctx, networkID)
```

## Image Pull

```go
import "github.com/docker/docker/api/types/image"

reader, err := cli.ImagePull(ctx, "alpine:latest", image.PullOptions{})
if err != nil {
    return err
}
defer reader.Close()
// Must drain the reader for the pull to complete
io.Copy(io.Discard, reader)
```

## Key Types

| Import Path | Types |
|-------------|-------|
| `github.com/docker/docker/client` | `Client`, `Opt`, `FromEnv`, `WithAPIVersionNegotiation` |
| `github.com/docker/docker/api/types/container` | `Config`, `HostConfig`, `StartOptions`, `StopOptions`, `RemoveOptions`, `ListOptions`, `Summary`, `InspectResponse` |
| `github.com/docker/docker/api/types/network` | `NetworkingConfig`, `EndpointSettings`, `CreateOptions`, `CreateResponse`, `ListOptions`, `Summary` |
| `github.com/docker/docker/api/types/image` | `PullOptions` |
| `github.com/docker/docker/api/types/filters` | `NewArgs`, `Arg` |
| `github.com/docker/go-connections/nat` | `PortMap`, `PortBinding`, `Port` |
