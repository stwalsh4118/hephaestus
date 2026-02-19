package templates

import (
	"errors"
	"strconv"
	"sync"
)

// Default port range for host port allocation.
const (
	DefaultMinPort = 10000
	DefaultMaxPort = 19999
)

// ErrPortsExhausted is returned when the port range has no available ports.
var ErrPortsExhausted = errors.New("port range exhausted")

// PortAllocator assigns unique host ports from a configurable range.
// It is safe for concurrent use.
type PortAllocator struct {
	mu       sync.Mutex
	minPort  int
	maxPort  int
	nextPort int
	used     map[int]bool
}

// NewPortAllocator creates a PortAllocator that allocates ports in [minPort, maxPort].
func NewPortAllocator(minPort, maxPort int) *PortAllocator {
	return &PortAllocator{
		minPort:  minPort,
		maxPort:  maxPort,
		nextPort: minPort,
		used:     make(map[int]bool),
	}
}

// allocateLocked returns the next available port. Caller must hold a.mu.
func (a *PortAllocator) allocateLocked() (int, error) {
	rangeSize := a.maxPort - a.minPort + 1
	for tried := 0; tried < rangeSize; tried++ {
		port := a.nextPort
		a.nextPort++
		if a.nextPort > a.maxPort {
			a.nextPort = a.minPort
		}
		if !a.used[port] {
			a.used[port] = true
			return port, nil
		}
	}
	return 0, ErrPortsExhausted
}

// Allocate returns the next available port as a string.
func (a *PortAllocator) Allocate() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	port, err := a.allocateLocked()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(port), nil
}

// AllocateN atomically allocates n ports and returns them as a string slice.
// If the range cannot satisfy all n ports, no ports are consumed (rollback).
func (a *PortAllocator) AllocateN(n int) ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	allocated := make([]int, 0, n)
	for i := 0; i < n; i++ {
		port, err := a.allocateLocked()
		if err != nil {
			// Rollback already-allocated ports.
			for _, p := range allocated {
				delete(a.used, p)
			}
			return nil, err
		}
		allocated = append(allocated, port)
	}

	ports := make([]string, n)
	for i, p := range allocated {
		ports[i] = strconv.Itoa(p)
	}
	return ports, nil
}

// Reset clears all allocations, allowing ports to be reused.
func (a *PortAllocator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.nextPort = a.minPort
	a.used = make(map[int]bool)
}
