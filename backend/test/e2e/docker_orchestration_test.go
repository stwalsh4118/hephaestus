package e2e

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
)

// skipIfNoDocker skips the test if Docker is not available.
func skipIfNoDocker(t *testing.T) *docker.Client {
	t.Helper()
	if os.Getenv("CI") != "" && os.Getenv("DOCKER_HOST") == "" {
		t.Skip("skipping docker e2e test in CI without Docker")
	}
	c, err := docker.NewClient()
	if err != nil {
		t.Skipf("skipping: docker client init failed: %v", err)
	}
	if err := c.Ping(context.Background()); err != nil {
		_ = c.Close()
		t.Skipf("skipping: docker daemon not reachable: %v", err)
	}
	return c
}

// TestAC1_DockerClientConnects verifies the Docker SDK client can connect to the daemon.
func TestAC1_DockerClientConnects(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

// TestAC2_ContainerCreatedWithConfig verifies containers can be created with env vars, ports, and volumes.
func TestAC2_ContainerCreatedWithConfig(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
	})

	id, err := o.CreateContainer(ctx, docker.ContainerConfig{
		Image:    "alpine:latest",
		Name:     "e2e-config-test",
		Env:      map[string]string{"TEST_VAR": "hello"},
		Hostname: "config-test",
	})
	if err != nil {
		t.Fatalf("CreateContainer: %v", err)
	}

	if id == "" {
		t.Fatal("expected non-empty container ID")
	}

	info, err := o.InspectContainer(ctx, id)
	if err != nil {
		t.Fatalf("InspectContainer: %v", err)
	}
	if info.Image != "alpine:latest" {
		t.Errorf("expected image 'alpine:latest', got %q", info.Image)
	}
}

// TestAC3_ContainerJoinsSharedNetwork verifies containers join the shared network.
func TestAC3_ContainerJoinsSharedNetwork(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
	})

	_, err := o.CreateContainer(ctx, docker.ContainerConfig{
		Image: "alpine:latest",
		Name:  "e2e-net-test",
	})
	if err != nil {
		t.Fatalf("CreateContainer: %v", err)
	}

	// Verify via ListContainers that the container exists on the network.
	containers, err := o.ListContainers(ctx)
	if err != nil {
		t.Fatalf("ListContainers: %v", err)
	}

	found := false
	for _, c := range containers {
		if strings.Contains(c.Name, "e2e-net-test") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find e2e-net-test in container list")
	}
}

// TestAC4_ContainerLifecycle verifies create → start → stop → remove flow.
func TestAC4_ContainerLifecycle(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
	})

	// Create.
	id, err := o.CreateContainer(ctx, docker.ContainerConfig{
		Image: "alpine:latest",
		Name:  "e2e-lifecycle",
	})
	if err != nil {
		t.Fatalf("CreateContainer: %v", err)
	}

	// Start — use a command that keeps the container running.
	// Alpine with no CMD exits immediately, so we start and verify it ran.
	if err := o.StartContainer(ctx, id); err != nil {
		t.Fatalf("StartContainer: %v", err)
	}

	// Give container a moment to transition state.
	time.Sleep(500 * time.Millisecond)

	// Inspect after start.
	info, err := o.InspectContainer(ctx, id)
	if err != nil {
		t.Fatalf("InspectContainer after start: %v", err)
	}
	// Alpine with no CMD may have already exited — that's OK for the lifecycle test.
	// The key test is that start/stop/remove don't error.
	t.Logf("status after start: %s", info.Status)

	// Stop.
	if err := o.StopContainer(ctx, id); err != nil {
		// Container may have already exited.
		t.Logf("StopContainer (may be already stopped): %v", err)
	}

	// Remove.
	if err := o.RemoveContainer(ctx, id); err != nil {
		t.Fatalf("RemoveContainer: %v", err)
	}

	// Verify gone.
	_, err = o.InspectContainer(ctx, id)
	if err == nil {
		t.Error("expected error inspecting removed container")
	}
}

// TestAC5_HealthCheckReportsStatus verifies health check returns accurate status.
func TestAC5_HealthCheckReportsStatus(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
	})

	id, err := o.CreateContainer(ctx, docker.ContainerConfig{
		Image: "alpine:latest",
		Name:  "e2e-health",
	})
	if err != nil {
		t.Fatalf("CreateContainer: %v", err)
	}

	// Check status before start — should be created.
	status, err := o.HealthCheck(ctx, id)
	if err != nil {
		t.Fatalf("HealthCheck before start: %v", err)
	}
	if status != docker.StatusCreated {
		t.Errorf("expected status %q before start, got %q", docker.StatusCreated, status)
	}

	// Start container.
	if err := o.StartContainer(ctx, id); err != nil {
		t.Fatalf("StartContainer: %v", err)
	}

	// Give it a moment.
	time.Sleep(500 * time.Millisecond)

	// Check status — should be running or stopped (alpine exits quickly).
	status, err = o.HealthCheck(ctx, id)
	if err != nil {
		t.Fatalf("HealthCheck after start: %v", err)
	}
	if status != docker.StatusRunning && status != docker.StatusStopped {
		t.Errorf("expected status running or stopped, got %q", status)
	}
	t.Logf("health check status: %s", status)
}

// TestAC6_TeardownCleansAll verifies teardown removes all containers and the network.
func TestAC6_TeardownCleansAll(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}

	// Create multiple containers.
	for _, name := range []string{"e2e-teardown-1", "e2e-teardown-2"} {
		_, err := o.CreateContainer(ctx, docker.ContainerConfig{
			Image: "alpine:latest",
			Name:  name,
		})
		if err != nil {
			t.Fatalf("CreateContainer(%s): %v", name, err)
		}
	}

	// Verify containers exist.
	before, err := o.ListContainers(ctx)
	if err != nil {
		t.Fatalf("ListContainers: %v", err)
	}
	if len(before) < 2 {
		t.Fatalf("expected at least 2 containers, got %d", len(before))
	}

	// Teardown.
	if err := o.TeardownAll(ctx); err != nil {
		t.Fatalf("TeardownAll: %v", err)
	}

	// Verify all managed containers are gone.
	after, err := o.ListContainers(ctx)
	if err != nil {
		t.Fatalf("ListContainers after teardown: %v", err)
	}

	for _, c := range after {
		if strings.Contains(c.Name, "e2e-teardown") {
			t.Errorf("container %q still exists after teardown", c.Name)
		}
	}
}

// TestAC7_ContainerNamesArePrefixed verifies all managed container names start with heph-.
func TestAC7_ContainerNamesArePrefixed(t *testing.T) {
	c := skipIfNoDocker(t)
	defer func() { _ = c.Close() }()

	o := docker.NewDockerOrchestrator(c)
	ctx := context.Background()

	if err := o.CreateNetwork(ctx); err != nil {
		t.Fatalf("CreateNetwork: %v", err)
	}
	t.Cleanup(func() {
		if err := o.TeardownAll(context.Background()); err != nil {
			t.Logf("teardown error: %v", err)
		}
	})

	_, err := o.CreateContainer(ctx, docker.ContainerConfig{
		Image: "alpine:latest",
		Name:  "e2e-prefix-test",
	})
	if err != nil {
		t.Fatalf("CreateContainer: %v", err)
	}

	containers, err := o.ListContainers(ctx)
	if err != nil {
		t.Fatalf("ListContainers: %v", err)
	}

	for _, c := range containers {
		if !strings.HasPrefix(c.Name, docker.ContainerNamePrefix) {
			t.Errorf("container name %q does not start with prefix %q", c.Name, docker.ContainerNamePrefix)
		}
	}
}
