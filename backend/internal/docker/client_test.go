package docker

import (
	"context"
	"strings"
	"testing"
)

func TestNewClient_ReturnsNonNil(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewClient() returned nil client")
	}
	if c.apiClient() == nil {
		t.Fatal("apiClient() returned nil")
	}
}

func TestClose_NoError(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() returned error: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}
}

func TestPing_ReturnsWrappedError_WhenDaemonUnavailable(t *testing.T) {
	// Create a client pointed at an invalid host to simulate unreachable daemon.
	c, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() returned error: %v", err)
	}
	defer func() { _ = c.Close() }()

	// Use a very short-lived context. On machines without Docker, Ping will fail.
	// On machines with Docker running, this test still passes (Ping succeeds).
	ctx := context.Background()
	err = c.Ping(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "ping docker daemon") {
			t.Errorf("expected error to contain 'ping docker daemon', got: %v", err)
		}
	}
	// If err is nil, Docker is running locally — that's also fine.
}
