package docker

import "testing"

func TestContainerConfig_Instantiation(t *testing.T) {
	cfg := ContainerConfig{
		Image:       "nginx:latest",
		Name:        "test-service",
		Env:         map[string]string{"PORT": "8080"},
		Ports:       map[string]string{"8080": "80"},
		Volumes:     map[string]string{"/data": "/var/data"},
		Hostname:    "test-service",
		NetworkName: "heph-network",
	}

	if cfg.Image != "nginx:latest" {
		t.Errorf("expected image 'nginx:latest', got %q", cfg.Image)
	}
	if cfg.Name != "test-service" {
		t.Errorf("expected name 'test-service', got %q", cfg.Name)
	}
}

func TestContainerInfo_Instantiation(t *testing.T) {
	info := ContainerInfo{
		ID:     "abc123",
		Name:   "heph-test",
		Image:  "alpine:latest",
		Status: StatusRunning,
		Ports:  map[string]string{"3000": "3000"},
	}

	if info.Status != StatusRunning {
		t.Errorf("expected status %q, got %q", StatusRunning, info.Status)
	}
}

func TestContainerStatus_ConstantsAreDistinct(t *testing.T) {
	statuses := []ContainerStatus{
		StatusCreated,
		StatusRunning,
		StatusStopped,
		StatusError,
		StatusHealthy,
		StatusUnhealthy,
	}

	seen := make(map[ContainerStatus]bool, len(statuses))
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("duplicate ContainerStatus value: %q", s)
		}
		seen[s] = true
	}
}

func TestContainerNamePrefix_IsNonEmpty(t *testing.T) {
	if ContainerNamePrefix == "" {
		t.Fatal("ContainerNamePrefix must not be empty")
	}
}
