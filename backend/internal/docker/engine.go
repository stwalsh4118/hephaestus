package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// NetworkName is the name of the shared Docker bridge network used by all
// managed containers.
const NetworkName = "heph-network"

// StopTimeout is the graceful stop timeout in seconds for containers.
const StopTimeout = 10

// ErrNotImplemented is returned by stub methods that are not yet implemented.
var ErrNotImplemented = errors.New("not implemented")

// Compile-time assertion that DockerOrchestrator implements Orchestrator.
var _ Orchestrator = (*DockerOrchestrator)(nil)

// dockerAPIClient defines the subset of the Docker SDK client used by the
// orchestrator. This enables unit testing with mocks.
type dockerAPIClient interface {
	// Network operations
	NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error)
	NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
	NetworkRemove(ctx context.Context, networkID string) error

	// Container operations
	ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error)
	ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error)
}

// DockerOrchestrator manages Docker containers and networks via the Docker SDK.
type DockerOrchestrator struct {
	api               dockerAPIClient
	mu                sync.Mutex
	networkID         string
	managedContainers map[string]string // container ID → name
}

// NewDockerOrchestrator creates an orchestrator using the provided Docker client.
func NewDockerOrchestrator(c *Client) *DockerOrchestrator {
	return &DockerOrchestrator{
		api:               &sdkClientAdapter{c.cli},
		managedContainers: make(map[string]string),
	}
}

// newOrchestratorWithAPI creates an orchestrator with a custom API client (for testing).
func newOrchestratorWithAPI(api dockerAPIClient) *DockerOrchestrator {
	return &DockerOrchestrator{
		api:               api,
		managedContainers: make(map[string]string),
	}
}

// CreateNetwork creates the shared Docker bridge network. If the network
// already exists, it reuses the existing one (idempotent).
func (o *DockerOrchestrator) CreateNetwork(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	existing, err := o.api.NetworkList(ctx, network.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", NetworkName)),
	})
	if err != nil {
		return fmt.Errorf("list networks: %w", err)
	}

	for _, n := range existing {
		if n.Name == NetworkName {
			o.networkID = n.ID
			return nil
		}
	}

	resp, err := o.api.NetworkCreate(ctx, NetworkName, network.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return fmt.Errorf("create network %q: %w", NetworkName, err)
	}

	o.networkID = resp.ID
	return nil
}

// RemoveNetwork removes the shared Docker bridge network. Returns nil if the
// network does not exist.
func (o *DockerOrchestrator) RemoveNetwork(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.networkID == "" {
		return nil
	}

	if err := o.api.NetworkRemove(ctx, o.networkID); err != nil {
		if errdefs.IsNotFound(err) {
			o.networkID = ""
			return nil
		}
		return fmt.Errorf("remove network %q: %w", NetworkName, err)
	}

	o.networkID = ""
	return nil
}

// CreateContainer pulls the image (if needed), creates a container with the
// given configuration, and connects it to the shared network.
func (o *DockerOrchestrator) CreateContainer(ctx context.Context, cfg ContainerConfig) (string, error) {
	// Pull image — drain the reader to complete the pull.
	reader, err := o.api.ImagePull(ctx, cfg.Image, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("pull image %q: %w", cfg.Image, err)
	}
	if _, err := io.Copy(io.Discard, reader); err != nil {
		_ = reader.Close()
		return "", fmt.Errorf("read image pull response for %q: %w", cfg.Image, err)
	}
	_ = reader.Close()

	prefixedName := ContainerNamePrefix + cfg.Name

	// Build environment slice.
	env := make([]string, 0, len(cfg.Env))
	for k, v := range cfg.Env {
		env = append(env, k+"="+v)
	}

	// Build port bindings and exposed ports.
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}
	for hostPort, containerPort := range cfg.Ports {
		cp := nat.Port(containerPort + "/tcp")
		exposedPorts[cp] = struct{}{}
		portBindings[cp] = []nat.PortBinding{{HostPort: hostPort}}
	}

	// Build volume binds.
	binds := make([]string, 0, len(cfg.Volumes))
	for hostPath, containerPath := range cfg.Volumes {
		binds = append(binds, hostPath+":"+containerPath)
	}

	// Determine hostname.
	hostname := cfg.Hostname
	if hostname == "" {
		hostname = cfg.Name
	}

	// Build networking config to attach to the shared network.
	networkCfg := &network.NetworkingConfig{}
	netName := cfg.NetworkName
	if netName == "" {
		netName = NetworkName
	}
	networkCfg.EndpointsConfig = map[string]*network.EndpointSettings{
		netName: {},
	}

	resp, err := o.api.ContainerCreate(ctx,
		&container.Config{
			Image:        cfg.Image,
			Env:          env,
			ExposedPorts: exposedPorts,
			Hostname:     hostname,
		},
		&container.HostConfig{
			PortBindings: portBindings,
			Binds:        binds,
		},
		networkCfg,
		prefixedName,
	)
	if err != nil {
		return "", fmt.Errorf("create container %q: %w", prefixedName, err)
	}

	o.mu.Lock()
	o.managedContainers[resp.ID] = prefixedName
	o.mu.Unlock()

	return resp.ID, nil
}

// StartContainer starts a previously created container.
func (o *DockerOrchestrator) StartContainer(ctx context.Context, containerID string) error {
	if err := o.api.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("start container %q: %w", containerID, err)
	}
	return nil
}

// StopContainer gracefully stops a running container with a timeout.
func (o *DockerOrchestrator) StopContainer(ctx context.Context, containerID string) error {
	timeout := StopTimeout
	if err := o.api.ContainerStop(ctx, containerID, container.StopOptions{
		Timeout: &timeout,
	}); err != nil {
		return fmt.Errorf("stop container %q: %w", containerID, err)
	}
	return nil
}

// RemoveContainer force-removes a container.
func (o *DockerOrchestrator) RemoveContainer(ctx context.Context, containerID string) error {
	if err := o.api.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: true,
	}); err != nil {
		return fmt.Errorf("remove container %q: %w", containerID, err)
	}

	o.mu.Lock()
	delete(o.managedContainers, containerID)
	o.mu.Unlock()

	return nil
}

// ListContainers returns info for all containers with the managed prefix.
func (o *DockerOrchestrator) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := o.api.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", ContainerNamePrefix)),
	})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	infos := make([]ContainerInfo, 0, len(containers))
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		infos = append(infos, ContainerInfo{
			ID:     c.ID,
			Name:   name,
			Image:  c.Image,
			Status: mapContainerState(c.State),
		})
	}
	return infos, nil
}

// InspectContainer returns detailed info for a single container.
func (o *DockerOrchestrator) InspectContainer(ctx context.Context, containerID string) (*ContainerInfo, error) {
	resp, err := o.api.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("inspect container %q: %w", containerID, err)
	}

	info := &ContainerInfo{
		ID:     resp.ID,
		Name:   strings.TrimPrefix(resp.Name, "/"),
		Image:  resp.Config.Image,
		Status: mapInspectState(resp.State),
	}
	return info, nil
}

// mapContainerState maps Docker's short state string to ContainerStatus.
func mapContainerState(state string) ContainerStatus {
	switch state {
	case "created":
		return StatusCreated
	case "running":
		return StatusRunning
	case "exited", "dead":
		return StatusStopped
	default:
		return StatusError
	}
}

// mapInspectState maps the detailed container state from inspect to ContainerStatus.
func mapInspectState(state *container.State) ContainerStatus {
	if state == nil {
		return StatusError
	}

	switch state.Status {
	case "created":
		return StatusCreated
	case "running":
		if state.Health != nil {
			switch state.Health.Status {
			case "healthy":
				return StatusHealthy
			case "unhealthy":
				return StatusUnhealthy
			}
		}
		return StatusRunning
	case "exited", "dead":
		return StatusStopped
	default:
		return StatusError
	}
}

// HealthCheck inspects a single container and returns its current status.
// If the container has disappeared, it reports StatusError and removes the
// container from internal tracking.
func (o *DockerOrchestrator) HealthCheck(ctx context.Context, containerID string) (ContainerStatus, error) {
	resp, err := o.api.ContainerInspect(ctx, containerID)
	if err != nil {
		if errdefs.IsNotFound(err) {
			o.mu.Lock()
			delete(o.managedContainers, containerID)
			o.mu.Unlock()
			return StatusError, nil
		}
		return "", fmt.Errorf("health check container %q: %w", containerID, err)
	}
	return mapInspectState(resp.State), nil
}

// StartHealthPolling runs a background goroutine that polls all managed
// containers at the given interval and calls the callback with their statuses.
// It stops when the context is cancelled.
func (o *DockerOrchestrator) StartHealthPolling(ctx context.Context, interval time.Duration, callback HealthStatusCallback) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				o.mu.Lock()
				ids := make([]string, 0, len(o.managedContainers))
				for id := range o.managedContainers {
					ids = append(ids, id)
				}
				o.mu.Unlock()

				statuses := make(map[string]ContainerStatus, len(ids))
				for _, id := range ids {
					status, _ := o.HealthCheck(ctx, id)
					statuses[id] = status
				}

				if len(statuses) > 0 {
					callback(statuses)
				}
			}
		}
	}()
}

// TeardownAll stops and removes all managed containers, then removes the
// shared network. It continues even if individual operations fail, collecting
// all errors. It is idempotent — safe to call multiple times.
func (o *DockerOrchestrator) TeardownAll(ctx context.Context) error {
	o.mu.Lock()
	ids := make([]string, 0, len(o.managedContainers))
	for id := range o.managedContainers {
		ids = append(ids, id)
	}
	o.mu.Unlock()

	var errs []error

	// Stop and remove each managed container.
	for _, id := range ids {
		timeout := StopTimeout
		if err := o.api.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
			if !errdefs.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("stop container %q: %w", id, err))
			}
		}

		if err := o.api.ContainerRemove(ctx, id, container.RemoveOptions{Force: true}); err != nil {
			if !errdefs.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("remove container %q: %w", id, err))
			}
		}
	}

	// Clear tracking and capture network ID under lock.
	o.mu.Lock()
	for k := range o.managedContainers {
		delete(o.managedContainers, k)
	}
	netID := o.networkID
	o.networkID = ""
	o.mu.Unlock()

	// Remove the shared network.
	if netID != "" {
		if err := o.api.NetworkRemove(ctx, netID); err != nil {
			if !errdefs.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("remove network: %w", err))
			}
		}
	}

	return errors.Join(errs...)
}
