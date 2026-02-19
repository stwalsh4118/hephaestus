package docker

import "context"

// Orchestrator defines the contract for Docker container lifecycle management.
type Orchestrator interface {
	// CreateContainer creates a new container from the given config and returns its ID.
	CreateContainer(ctx context.Context, config ContainerConfig) (string, error)

	// StartContainer starts a previously created container.
	StartContainer(ctx context.Context, containerID string) error

	// StopContainer gracefully stops a running container.
	StopContainer(ctx context.Context, containerID string) error

	// RemoveContainer removes a container, force-removing if still running.
	RemoveContainer(ctx context.Context, containerID string) error

	// ListContainers returns info for all containers managed by this orchestrator.
	ListContainers(ctx context.Context) ([]ContainerInfo, error)

	// InspectContainer returns detailed info for a single container.
	InspectContainer(ctx context.Context, containerID string) (*ContainerInfo, error)

	// CreateNetwork creates the shared Docker bridge network.
	CreateNetwork(ctx context.Context) error

	// RemoveNetwork removes the shared Docker bridge network.
	RemoveNetwork(ctx context.Context) error

	// HealthCheck inspects a container and returns its current status.
	HealthCheck(ctx context.Context, containerID string) (ContainerStatus, error)

	// TeardownAll stops and removes all managed containers and the shared network.
	TeardownAll(ctx context.Context) error
}
