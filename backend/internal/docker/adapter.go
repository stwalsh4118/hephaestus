package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// sdkClientAdapter wraps the Docker SDK *client.Client to conform to the
// dockerAPIClient interface. This adapter omits the platform parameter from
// ContainerCreate (always nil) so that the orchestrator and tests use a
// simpler interface.
type sdkClientAdapter struct {
	cli *client.Client
}

func (a *sdkClientAdapter) NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error) {
	return a.cli.NetworkCreate(ctx, name, options)
}

func (a *sdkClientAdapter) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
	return a.cli.NetworkList(ctx, options)
}

func (a *sdkClientAdapter) NetworkRemove(ctx context.Context, networkID string) error {
	return a.cli.NetworkRemove(ctx, networkID)
}

func (a *sdkClientAdapter) ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error) {
	return a.cli.ImagePull(ctx, refStr, options)
}

func (a *sdkClientAdapter) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.CreateResponse, error) {
	return a.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
}

func (a *sdkClientAdapter) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	return a.cli.ContainerStart(ctx, containerID, options)
}

func (a *sdkClientAdapter) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	return a.cli.ContainerStop(ctx, containerID, options)
}

func (a *sdkClientAdapter) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	return a.cli.ContainerRemove(ctx, containerID, options)
}

func (a *sdkClientAdapter) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
	return a.cli.ContainerList(ctx, options)
}

func (a *sdkClientAdapter) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error) {
	return a.cli.ContainerInspect(ctx, containerID)
}
