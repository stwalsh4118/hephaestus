package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

// DockerAPI defines the contract for Docker daemon connectivity.
type DockerAPI interface {
	// Ping verifies connectivity to the Docker daemon.
	Ping(ctx context.Context) error

	// Close releases resources held by the client.
	Close() error
}

// Compile-time assertion that Client implements DockerAPI.
var _ DockerAPI = (*Client)(nil)

// Client wraps the Docker SDK client with connection management.
type Client struct {
	cli *client.Client
}

// NewClient creates a Client using the default Docker socket/environment.
// It configures API version negotiation so the client works with any daemon version.
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}
	return &Client{cli: cli}, nil
}

// Ping verifies connectivity to the Docker daemon.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ping docker daemon: %w", err)
	}
	return nil
}

// Close releases resources held by the underlying Docker SDK client.
func (c *Client) Close() error {
	if err := c.cli.Close(); err != nil {
		return fmt.Errorf("close docker client: %w", err)
	}
	return nil
}

// apiClient returns the underlying Docker SDK client for use within the package.
func (c *Client) apiClient() *client.Client {
	return c.cli
}
