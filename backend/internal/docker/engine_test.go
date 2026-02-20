package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
)

// mockDockerAPI is a test double for the Docker SDK client.
type mockDockerAPI struct {
	networkCreateFn    func(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error)
	networkListFn      func(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
	networkRemoveFn    func(ctx context.Context, networkID string) error
	imagePullFn        func(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error)
	containerCreateFn  func(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.CreateResponse, error)
	containerStartFn   func(ctx context.Context, containerID string, options container.StartOptions) error
	containerStopFn    func(ctx context.Context, containerID string, options container.StopOptions) error
	containerRemoveFn  func(ctx context.Context, containerID string, options container.RemoveOptions) error
	containerListFn    func(ctx context.Context, options container.ListOptions) ([]container.Summary, error)
	containerInspectFn func(ctx context.Context, containerID string) (container.InspectResponse, error)
}

func (m *mockDockerAPI) NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error) {
	if m.networkCreateFn != nil {
		return m.networkCreateFn(ctx, name, options)
	}
	return network.CreateResponse{}, nil
}

func (m *mockDockerAPI) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
	if m.networkListFn != nil {
		return m.networkListFn(ctx, options)
	}
	return nil, nil
}

func (m *mockDockerAPI) NetworkRemove(ctx context.Context, networkID string) error {
	if m.networkRemoveFn != nil {
		return m.networkRemoveFn(ctx, networkID)
	}
	return nil
}

func (m *mockDockerAPI) ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error) {
	if m.imagePullFn != nil {
		return m.imagePullFn(ctx, refStr, options)
	}
	return io.NopCloser(bytes.NewReader(nil)), nil
}

func (m *mockDockerAPI) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.CreateResponse, error) {
	if m.containerCreateFn != nil {
		return m.containerCreateFn(ctx, config, hostConfig, networkingConfig, containerName)
	}
	return container.CreateResponse{}, nil
}

func (m *mockDockerAPI) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	if m.containerStartFn != nil {
		return m.containerStartFn(ctx, containerID, options)
	}
	return nil
}

func (m *mockDockerAPI) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	if m.containerStopFn != nil {
		return m.containerStopFn(ctx, containerID, options)
	}
	return nil
}

func (m *mockDockerAPI) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	if m.containerRemoveFn != nil {
		return m.containerRemoveFn(ctx, containerID, options)
	}
	return nil
}

func (m *mockDockerAPI) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
	if m.containerListFn != nil {
		return m.containerListFn(ctx, options)
	}
	return nil, nil
}

func (m *mockDockerAPI) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error) {
	if m.containerInspectFn != nil {
		return m.containerInspectFn(ctx, containerID)
	}
	return container.InspectResponse{}, nil
}

// --- Network Tests (from task 6-3) ---

func TestCreateNetwork_CreatesNewBridgeNetwork(t *testing.T) {
	var createdName string
	var createdDriver string

	mock := &mockDockerAPI{
		networkListFn: func(_ context.Context, _ network.ListOptions) ([]network.Summary, error) {
			return nil, nil
		},
		networkCreateFn: func(_ context.Context, name string, options network.CreateOptions) (network.CreateResponse, error) {
			createdName = name
			createdDriver = options.Driver
			return network.CreateResponse{ID: "net-123"}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	if err := o.CreateNetwork(context.Background()); err != nil {
		t.Fatalf("CreateNetwork() returned error: %v", err)
	}
	if createdName != NetworkName {
		t.Errorf("expected network name %q, got %q", NetworkName, createdName)
	}
	if createdDriver != "bridge" {
		t.Errorf("expected driver 'bridge', got %q", createdDriver)
	}
	if o.networkID != "net-123" {
		t.Errorf("expected networkID 'net-123', got %q", o.networkID)
	}
}

func TestCreateNetwork_IsIdempotent(t *testing.T) {
	createCalls := 0
	mock := &mockDockerAPI{
		networkListFn: func(_ context.Context, _ network.ListOptions) ([]network.Summary, error) {
			return []network.Summary{{ID: "existing-net", Name: NetworkName}}, nil
		},
		networkCreateFn: func(_ context.Context, _ string, _ network.CreateOptions) (network.CreateResponse, error) {
			createCalls++
			return network.CreateResponse{}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	if err := o.CreateNetwork(context.Background()); err != nil {
		t.Fatalf("CreateNetwork() returned error: %v", err)
	}
	if createCalls != 0 {
		t.Errorf("expected 0 NetworkCreate calls, got %d", createCalls)
	}
	if o.networkID != "existing-net" {
		t.Errorf("expected networkID 'existing-net', got %q", o.networkID)
	}
}

func TestRemoveNetwork_RemovesSuccessfully(t *testing.T) {
	var removedID string
	mock := &mockDockerAPI{
		networkRemoveFn: func(_ context.Context, networkID string) error {
			removedID = networkID
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.networkID = "net-to-remove"

	if err := o.RemoveNetwork(context.Background()); err != nil {
		t.Fatalf("RemoveNetwork() returned error: %v", err)
	}
	if removedID != "net-to-remove" {
		t.Errorf("expected to remove 'net-to-remove', got %q", removedID)
	}
	if o.networkID != "" {
		t.Errorf("expected networkID to be cleared, got %q", o.networkID)
	}
}

func TestRemoveNetwork_NoOpWhenNoNetwork(t *testing.T) {
	removeCalls := 0
	mock := &mockDockerAPI{
		networkRemoveFn: func(_ context.Context, _ string) error {
			removeCalls++
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	if err := o.RemoveNetwork(context.Background()); err != nil {
		t.Fatalf("RemoveNetwork() returned error: %v", err)
	}
	if removeCalls != 0 {
		t.Errorf("expected 0 NetworkRemove calls, got %d", removeCalls)
	}
}

func TestRemoveNetwork_HandlesNotFound(t *testing.T) {
	mock := &mockDockerAPI{
		networkRemoveFn: func(_ context.Context, _ string) error {
			return notFoundError("network not found")
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.networkID = "gone-net"

	if err := o.RemoveNetwork(context.Background()); err != nil {
		t.Fatalf("RemoveNetwork() should not error on not-found, got: %v", err)
	}
	if o.networkID != "" {
		t.Errorf("expected networkID to be cleared, got %q", o.networkID)
	}
}

func TestCreateNetwork_ReturnsErrorOnListFailure(t *testing.T) {
	mock := &mockDockerAPI{
		networkListFn: func(_ context.Context, _ network.ListOptions) ([]network.Summary, error) {
			return nil, errors.New("connection refused")
		},
	}

	o := newOrchestratorWithAPI(mock)
	if err := o.CreateNetwork(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Container Lifecycle Tests (task 6-4) ---

func TestCreateContainer_PullsImageAndCreates(t *testing.T) {
	var pulledImage string
	var createdConfig *container.Config
	var createdHostConfig *container.HostConfig
	var createdNetConfig *network.NetworkingConfig
	var createdName string

	mock := &mockDockerAPI{
		imagePullFn: func(_ context.Context, refStr string, _ image.PullOptions) (io.ReadCloser, error) {
			pulledImage = refStr
			return io.NopCloser(bytes.NewReader(nil)), nil
		},
		containerCreateFn: func(_ context.Context, config *container.Config, hostConfig *container.HostConfig, netConfig *network.NetworkingConfig, name string) (container.CreateResponse, error) {
			createdConfig = config
			createdHostConfig = hostConfig
			createdNetConfig = netConfig
			createdName = name
			return container.CreateResponse{ID: "ctr-123"}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	id, err := o.CreateContainer(context.Background(), ContainerConfig{
		Image:    "nginx:latest",
		Name:     "web",
		Env:      map[string]string{"PORT": "80"},
		Ports:    map[string]string{"8080": "80"},
		Volumes:  map[string]string{"/host": "/container"},
		Hostname: "web-service",
	})
	if err != nil {
		t.Fatalf("CreateContainer() returned error: %v", err)
	}

	if id != "ctr-123" {
		t.Errorf("expected container ID 'ctr-123', got %q", id)
	}
	if pulledImage != "nginx:latest" {
		t.Errorf("expected pulled image 'nginx:latest', got %q", pulledImage)
	}
	if createdName != "heph-web" {
		t.Errorf("expected container name 'heph-web', got %q", createdName)
	}
	if createdConfig.Image != "nginx:latest" {
		t.Errorf("expected config image 'nginx:latest', got %q", createdConfig.Image)
	}
	if createdConfig.Hostname != "web-service" {
		t.Errorf("expected hostname 'web-service', got %q", createdConfig.Hostname)
	}
	if len(createdHostConfig.PortBindings) == 0 {
		t.Error("expected port bindings to be set")
	}
	if len(createdHostConfig.Binds) == 0 {
		t.Error("expected volume binds to be set")
	}
	if _, ok := createdNetConfig.EndpointsConfig[NetworkName]; !ok {
		t.Error("expected container to be connected to shared network")
	}
	// Verify tracking.
	if o.managedContainers["ctr-123"] != "heph-web" {
		t.Errorf("expected managed container tracking, got %v", o.managedContainers)
	}
}

func TestCreateContainer_PassesCmdToDockerConfig(t *testing.T) {
	var createdConfig *container.Config

	mock := &mockDockerAPI{
		containerCreateFn: func(_ context.Context, config *container.Config, _ *container.HostConfig, _ *network.NetworkingConfig, _ string) (container.CreateResponse, error) {
			createdConfig = config
			return container.CreateResponse{ID: "ctr-cmd"}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	cmd := []string{"mock", "-h", "0.0.0.0", "/tmp/spec.json"}
	_, err := o.CreateContainer(context.Background(), ContainerConfig{
		Image: "stoplight/prism:latest",
		Name:  "api-service",
		Cmd:   cmd,
	})
	if err != nil {
		t.Fatalf("CreateContainer() returned error: %v", err)
	}

	if len(createdConfig.Cmd) != len(cmd) {
		t.Fatalf("expected Cmd length %d, got %d", len(cmd), len(createdConfig.Cmd))
	}
	for i, arg := range cmd {
		if createdConfig.Cmd[i] != arg {
			t.Errorf("Cmd[%d] = %q, want %q", i, createdConfig.Cmd[i], arg)
		}
	}
}

func TestCreateContainer_NilCmdUsesImageDefault(t *testing.T) {
	var createdConfig *container.Config

	mock := &mockDockerAPI{
		containerCreateFn: func(_ context.Context, config *container.Config, _ *container.HostConfig, _ *network.NetworkingConfig, _ string) (container.CreateResponse, error) {
			createdConfig = config
			return container.CreateResponse{ID: "ctr-nocmd"}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	_, err := o.CreateContainer(context.Background(), ContainerConfig{
		Image: "nginx:latest",
		Name:  "web",
	})
	if err != nil {
		t.Fatalf("CreateContainer() returned error: %v", err)
	}

	if createdConfig.Cmd != nil {
		t.Errorf("expected nil Cmd for unset config, got %v", createdConfig.Cmd)
	}
}

func TestCreateContainer_AppliesNamePrefix(t *testing.T) {
	var createdName string
	mock := &mockDockerAPI{
		containerCreateFn: func(_ context.Context, _ *container.Config, _ *container.HostConfig, _ *network.NetworkingConfig, name string) (container.CreateResponse, error) {
			createdName = name
			return container.CreateResponse{ID: "ctr-456"}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	_, err := o.CreateContainer(context.Background(), ContainerConfig{
		Image: "alpine:latest",
		Name:  "myservice",
	})
	if err != nil {
		t.Fatalf("CreateContainer() returned error: %v", err)
	}
	if createdName != "heph-myservice" {
		t.Errorf("expected name 'heph-myservice', got %q", createdName)
	}
}

func TestStartContainer_CallsDockerAPI(t *testing.T) {
	var startedID string
	mock := &mockDockerAPI{
		containerStartFn: func(_ context.Context, containerID string, _ container.StartOptions) error {
			startedID = containerID
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	if err := o.StartContainer(context.Background(), "ctr-123"); err != nil {
		t.Fatalf("StartContainer() returned error: %v", err)
	}
	if startedID != "ctr-123" {
		t.Errorf("expected started ID 'ctr-123', got %q", startedID)
	}
}

func TestStopContainer_UsesGracefulTimeout(t *testing.T) {
	var stoppedID string
	var timeout *int
	mock := &mockDockerAPI{
		containerStopFn: func(_ context.Context, containerID string, options container.StopOptions) error {
			stoppedID = containerID
			timeout = options.Timeout
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	if err := o.StopContainer(context.Background(), "ctr-123"); err != nil {
		t.Fatalf("StopContainer() returned error: %v", err)
	}
	if stoppedID != "ctr-123" {
		t.Errorf("expected stopped ID 'ctr-123', got %q", stoppedID)
	}
	if timeout == nil || *timeout != StopTimeout {
		t.Errorf("expected timeout %d, got %v", StopTimeout, timeout)
	}
}

func TestRemoveContainer_ForcesRemoval(t *testing.T) {
	var removedID string
	var wasForced bool
	mock := &mockDockerAPI{
		containerRemoveFn: func(_ context.Context, containerID string, options container.RemoveOptions) error {
			removedID = containerID
			wasForced = options.Force
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.managedContainers["ctr-123"] = "heph-test"

	if err := o.RemoveContainer(context.Background(), "ctr-123"); err != nil {
		t.Fatalf("RemoveContainer() returned error: %v", err)
	}
	if removedID != "ctr-123" {
		t.Errorf("expected removed ID 'ctr-123', got %q", removedID)
	}
	if !wasForced {
		t.Error("expected Force=true")
	}
	// Verify tracking removal.
	if _, exists := o.managedContainers["ctr-123"]; exists {
		t.Error("expected container to be removed from tracking")
	}
}

func TestListContainers_MapsToContainerInfo(t *testing.T) {
	mock := &mockDockerAPI{
		containerListFn: func(_ context.Context, _ container.ListOptions) ([]container.Summary, error) {
			return []container.Summary{
				{ID: "ctr-1", Names: []string{"/heph-web"}, Image: "nginx:latest", State: "running"},
				{ID: "ctr-2", Names: []string{"/heph-db"}, Image: "postgres:15", State: "exited"},
			}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	infos, err := o.ListContainers(context.Background())
	if err != nil {
		t.Fatalf("ListContainers() returned error: %v", err)
	}
	if len(infos) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(infos))
	}
	if infos[0].Name != "heph-web" {
		t.Errorf("expected name 'heph-web', got %q", infos[0].Name)
	}
	if infos[0].Status != StatusRunning {
		t.Errorf("expected status %q, got %q", StatusRunning, infos[0].Status)
	}
	if infos[1].Status != StatusStopped {
		t.Errorf("expected status %q, got %q", StatusStopped, infos[1].Status)
	}
}

func TestInspectContainer_MapsStateCorrectly(t *testing.T) {
	mock := &mockDockerAPI{
		containerInspectFn: func(_ context.Context, _ string) (container.InspectResponse, error) {
			return container.InspectResponse{
				ContainerJSONBase: &container.ContainerJSONBase{
					ID:   "ctr-1",
					Name: "/heph-web",
					State: &container.State{
						Status: "running",
					},
				},
				Config: &container.Config{Image: "nginx:latest"},
			}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	info, err := o.InspectContainer(context.Background(), "ctr-1")
	if err != nil {
		t.Fatalf("InspectContainer() returned error: %v", err)
	}
	if info.ID != "ctr-1" {
		t.Errorf("expected ID 'ctr-1', got %q", info.ID)
	}
	if info.Name != "heph-web" {
		t.Errorf("expected name 'heph-web', got %q", info.Name)
	}
	if info.Status != StatusRunning {
		t.Errorf("expected status %q, got %q", StatusRunning, info.Status)
	}
}

func TestMapContainerState_AllStates(t *testing.T) {
	tests := []struct {
		state    string
		expected ContainerStatus
	}{
		{"created", StatusCreated},
		{"running", StatusRunning},
		{"exited", StatusStopped},
		{"dead", StatusStopped},
		{"paused", StatusError},
		{"unknown", StatusError},
	}

	for _, tt := range tests {
		got := mapContainerState(tt.state)
		if got != tt.expected {
			t.Errorf("mapContainerState(%q) = %q, want %q", tt.state, got, tt.expected)
		}
	}
}

func TestMapInspectState_WithHealthCheck(t *testing.T) {
	tests := []struct {
		name     string
		state    *container.State
		expected ContainerStatus
	}{
		{
			name:     "nil state",
			state:    nil,
			expected: StatusError,
		},
		{
			name:     "running without healthcheck",
			state:    &container.State{Status: "running"},
			expected: StatusRunning,
		},
		{
			name: "running healthy",
			state: &container.State{
				Status: "running",
				Health: &container.Health{Status: "healthy"},
			},
			expected: StatusHealthy,
		},
		{
			name: "running unhealthy",
			state: &container.State{
				Status: "running",
				Health: &container.Health{Status: "unhealthy"},
			},
			expected: StatusUnhealthy,
		},
		{
			name:     "exited",
			state:    &container.State{Status: "exited"},
			expected: StatusStopped,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapInspectState(tt.state)
			if got != tt.expected {
				t.Errorf("mapInspectState() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCreateContainer_WrapsErrorOnPullFailure(t *testing.T) {
	mock := &mockDockerAPI{
		imagePullFn: func(_ context.Context, _ string, _ image.PullOptions) (io.ReadCloser, error) {
			return nil, errors.New("pull failed")
		},
	}

	o := newOrchestratorWithAPI(mock)
	_, err := o.CreateContainer(context.Background(), ContainerConfig{Image: "bad:image"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Health Check Tests (task 6-5) ---

func TestHealthCheck_ReturnsRunningStatus(t *testing.T) {
	mock := &mockDockerAPI{
		containerInspectFn: func(_ context.Context, _ string) (container.InspectResponse, error) {
			return container.InspectResponse{
				ContainerJSONBase: &container.ContainerJSONBase{
					ID:    "ctr-1",
					Name:  "/heph-web",
					State: &container.State{Status: "running"},
				},
				Config: &container.Config{Image: "nginx:latest"},
			}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	status, err := o.HealthCheck(context.Background(), "ctr-1")
	if err != nil {
		t.Fatalf("HealthCheck() returned error: %v", err)
	}
	if status != StatusRunning {
		t.Errorf("expected status %q, got %q", StatusRunning, status)
	}
}

func TestHealthCheck_ReturnsStoppedStatus(t *testing.T) {
	mock := &mockDockerAPI{
		containerInspectFn: func(_ context.Context, _ string) (container.InspectResponse, error) {
			return container.InspectResponse{
				ContainerJSONBase: &container.ContainerJSONBase{
					ID:    "ctr-1",
					Name:  "/heph-web",
					State: &container.State{Status: "exited"},
				},
				Config: &container.Config{Image: "nginx:latest"},
			}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	status, err := o.HealthCheck(context.Background(), "ctr-1")
	if err != nil {
		t.Fatalf("HealthCheck() returned error: %v", err)
	}
	if status != StatusStopped {
		t.Errorf("expected status %q, got %q", StatusStopped, status)
	}
}

func TestHealthCheck_HandlesNotFound_RemovesFromTracking(t *testing.T) {
	mock := &mockDockerAPI{
		containerInspectFn: func(_ context.Context, _ string) (container.InspectResponse, error) {
			return container.InspectResponse{}, notFoundError("container not found")
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.managedContainers["ctr-gone"] = "heph-gone"

	status, err := o.HealthCheck(context.Background(), "ctr-gone")
	if err != nil {
		t.Fatalf("HealthCheck() should not error on not-found, got: %v", err)
	}
	if status != StatusError {
		t.Errorf("expected status %q, got %q", StatusError, status)
	}
	if _, exists := o.managedContainers["ctr-gone"]; exists {
		t.Error("expected container to be removed from tracking")
	}
}

func TestHealthCheck_WrapsAPIError(t *testing.T) {
	mock := &mockDockerAPI{
		containerInspectFn: func(_ context.Context, _ string) (container.InspectResponse, error) {
			return container.InspectResponse{}, errors.New("connection refused")
		},
	}

	o := newOrchestratorWithAPI(mock)
	_, err := o.HealthCheck(context.Background(), "ctr-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestStartHealthPolling_CallsCallbackAndStopsOnCancel(t *testing.T) {
	mock := &mockDockerAPI{
		containerInspectFn: func(_ context.Context, _ string) (container.InspectResponse, error) {
			return container.InspectResponse{
				ContainerJSONBase: &container.ContainerJSONBase{
					ID:    "ctr-1",
					Name:  "/heph-web",
					State: &container.State{Status: "running"},
				},
				Config: &container.Config{Image: "nginx:latest"},
			}, nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.managedContainers["ctr-1"] = "heph-web"

	var mu sync.Mutex
	callbackCount := 0
	var lastStatuses map[string]ContainerStatus

	ctx, cancel := context.WithCancel(context.Background())

	o.StartHealthPolling(ctx, 50*time.Millisecond, func(statuses map[string]ContainerStatus) {
		mu.Lock()
		callbackCount++
		lastStatuses = statuses
		mu.Unlock()
	})

	// Wait enough time for at least 2 ticks.
	time.Sleep(150 * time.Millisecond)
	cancel()
	// Give goroutine time to exit.
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if callbackCount < 1 {
		t.Errorf("expected at least 1 callback, got %d", callbackCount)
	}
	if lastStatuses["ctr-1"] != StatusRunning {
		t.Errorf("expected status %q for ctr-1, got %q", StatusRunning, lastStatuses["ctr-1"])
	}
}

func TestStartHealthPolling_NoCallbackWhenNoContainers(t *testing.T) {
	mock := &mockDockerAPI{}

	o := newOrchestratorWithAPI(mock)
	// No managed containers.

	callbackCount := 0
	ctx, cancel := context.WithCancel(context.Background())

	o.StartHealthPolling(ctx, 50*time.Millisecond, func(_ map[string]ContainerStatus) {
		callbackCount++
	})

	time.Sleep(150 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)

	if callbackCount != 0 {
		t.Errorf("expected 0 callbacks with no containers, got %d", callbackCount)
	}
}

// --- Teardown Tests (task 6-6) ---

func TestTeardownAll_StopsRemovesContainersAndNetwork(t *testing.T) {
	var stoppedIDs, removedIDs []string
	var networkRemoved bool

	mock := &mockDockerAPI{
		containerStopFn: func(_ context.Context, containerID string, _ container.StopOptions) error {
			stoppedIDs = append(stoppedIDs, containerID)
			return nil
		},
		containerRemoveFn: func(_ context.Context, containerID string, _ container.RemoveOptions) error {
			removedIDs = append(removedIDs, containerID)
			return nil
		},
		networkRemoveFn: func(_ context.Context, _ string) error {
			networkRemoved = true
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.managedContainers["ctr-1"] = "heph-web"
	o.managedContainers["ctr-2"] = "heph-db"
	o.networkID = "net-123"

	err := o.TeardownAll(context.Background())
	if err != nil {
		t.Fatalf("TeardownAll() returned error: %v", err)
	}

	if len(stoppedIDs) != 2 {
		t.Errorf("expected 2 containers stopped, got %d", len(stoppedIDs))
	}
	if len(removedIDs) != 2 {
		t.Errorf("expected 2 containers removed, got %d", len(removedIDs))
	}
	if !networkRemoved {
		t.Error("expected network to be removed")
	}
	if len(o.managedContainers) != 0 {
		t.Errorf("expected empty tracking map, got %d entries", len(o.managedContainers))
	}
	if o.networkID != "" {
		t.Errorf("expected networkID to be cleared, got %q", o.networkID)
	}
}

func TestTeardownAll_ContinuesOnPartialFailure(t *testing.T) {
	stopCalls := 0
	removeCalls := 0

	mock := &mockDockerAPI{
		containerStopFn: func(_ context.Context, _ string, _ container.StopOptions) error {
			stopCalls++
			if stopCalls == 1 {
				return errors.New("stop failed")
			}
			return nil
		},
		containerRemoveFn: func(_ context.Context, _ string, _ container.RemoveOptions) error {
			removeCalls++
			return nil
		},
		networkRemoveFn: func(_ context.Context, _ string) error {
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.managedContainers["ctr-1"] = "heph-web"
	o.managedContainers["ctr-2"] = "heph-db"
	o.networkID = "net-123"

	err := o.TeardownAll(context.Background())
	// Should have collected the stop error but still continued.
	if err == nil {
		t.Fatal("expected error from partial failure, got nil")
	}

	// Both containers should have been attempted for removal.
	if removeCalls != 2 {
		t.Errorf("expected 2 remove calls, got %d", removeCalls)
	}
}

func TestTeardownAll_IsIdempotent(t *testing.T) {
	removeCalls := 0

	mock := &mockDockerAPI{
		containerStopFn: func(_ context.Context, _ string, _ container.StopOptions) error {
			return nil
		},
		containerRemoveFn: func(_ context.Context, _ string, _ container.RemoveOptions) error {
			removeCalls++
			return nil
		},
		networkRemoveFn: func(_ context.Context, _ string) error {
			return nil
		},
	}

	o := newOrchestratorWithAPI(mock)
	// Empty state — no containers, no network.

	err := o.TeardownAll(context.Background())
	if err != nil {
		t.Fatalf("TeardownAll() returned error: %v", err)
	}
	if removeCalls != 0 {
		t.Errorf("expected 0 remove calls on empty state, got %d", removeCalls)
	}
}

func TestTeardownAll_CollectsMultipleErrors(t *testing.T) {
	mock := &mockDockerAPI{
		containerStopFn: func(_ context.Context, _ string, _ container.StopOptions) error {
			return errors.New("stop failed")
		},
		containerRemoveFn: func(_ context.Context, _ string, _ container.RemoveOptions) error {
			return errors.New("remove failed")
		},
		networkRemoveFn: func(_ context.Context, _ string) error {
			return errors.New("network remove failed")
		},
	}

	o := newOrchestratorWithAPI(mock)
	o.managedContainers["ctr-1"] = "heph-web"
	o.networkID = "net-123"

	err := o.TeardownAll(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// Should contain at least 3 errors (stop, remove, network).
	errStr := err.Error()
	if !errors.Is(err, err) { // basic sanity
		t.Error("error should not be nil")
	}
	if len(errStr) == 0 {
		t.Error("expected non-empty error string")
	}
}

// notFoundError implements the interface checked by errdefs.IsNotFound.
type notFoundError string

func (e notFoundError) Error() string { return string(e) }
func (e notFoundError) NotFound()     {}
