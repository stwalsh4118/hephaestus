package templates

import (
	"strconv"
	"sync"
	"testing"
)

func TestPortAllocator_SequentialAllocation(t *testing.T) {
	a := NewPortAllocator(10000, 10004)

	seen := make(map[string]bool)
	for i := 0; i < 5; i++ {
		p, err := a.Allocate()
		if err != nil {
			t.Fatalf("allocation %d: unexpected error: %v", i, err)
		}
		if seen[p] {
			t.Fatalf("allocation %d: duplicate port %s", i, p)
		}
		seen[p] = true

		n, _ := strconv.Atoi(p)
		if n < 10000 || n > 10004 {
			t.Fatalf("allocation %d: port %d out of range [10000, 10004]", i, n)
		}
	}

	if len(seen) != 5 {
		t.Errorf("expected 5 unique ports, got %d", len(seen))
	}
}

func TestPortAllocator_AllocateN(t *testing.T) {
	a := NewPortAllocator(10000, 10009)

	ports, err := a.AllocateN(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(ports))
	}

	seen := make(map[string]bool)
	for _, p := range ports {
		if seen[p] {
			t.Errorf("duplicate port %s", p)
		}
		seen[p] = true
	}
}

func TestPortAllocator_Reset(t *testing.T) {
	a := NewPortAllocator(10000, 10001)

	_, _ = a.Allocate()
	_, _ = a.Allocate()

	// Range exhausted.
	_, err := a.Allocate()
	if err != ErrPortsExhausted {
		t.Fatalf("expected ErrPortsExhausted, got %v", err)
	}

	a.Reset()

	// Should be able to allocate again.
	p, err := a.Allocate()
	if err != nil {
		t.Fatalf("unexpected error after reset: %v", err)
	}
	if p != "10000" {
		t.Errorf("expected port 10000 after reset, got %s", p)
	}
}

func TestPortAllocator_Exhaustion(t *testing.T) {
	a := NewPortAllocator(10000, 10002)

	for i := 0; i < 3; i++ {
		_, err := a.Allocate()
		if err != nil {
			t.Fatalf("allocation %d: unexpected error: %v", i, err)
		}
	}

	_, err := a.Allocate()
	if err != ErrPortsExhausted {
		t.Errorf("expected ErrPortsExhausted, got %v", err)
	}
}

func TestPortAllocator_AllocateN_Exhaustion(t *testing.T) {
	a := NewPortAllocator(10000, 10001)

	_, err := a.AllocateN(3)
	if err != ErrPortsExhausted {
		t.Errorf("expected ErrPortsExhausted, got %v", err)
	}
}

func TestPortAllocator_AllocateN_RollbackOnFailure(t *testing.T) {
	// Range has exactly 2 ports; requesting 3 should fail and roll back.
	a := NewPortAllocator(10000, 10001)

	_, err := a.AllocateN(3)
	if err != ErrPortsExhausted {
		t.Fatalf("expected ErrPortsExhausted, got %v", err)
	}

	// After rollback, all 2 ports should be available again.
	ports, err := a.AllocateN(2)
	if err != nil {
		t.Fatalf("unexpected error after rollback: %v", err)
	}
	if len(ports) != 2 {
		t.Errorf("expected 2 ports after rollback, got %d", len(ports))
	}
}

func TestPortAllocator_ConcurrentSafety(t *testing.T) {
	a := NewPortAllocator(10000, 10999)

	var wg sync.WaitGroup
	results := make(chan string, 100)
	errCh := make(chan error, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p, err := a.Allocate()
			if err != nil {
				errCh <- err
				return
			}
			results <- p
		}()
	}

	wg.Wait()
	close(results)
	close(errCh)

	for err := range errCh {
		t.Fatalf("concurrent allocation error: %v", err)
	}

	seen := make(map[string]bool)
	for p := range results {
		if seen[p] {
			t.Fatalf("concurrent allocation produced duplicate port: %s", p)
		}
		seen[p] = true
	}

	if len(seen) != 100 {
		t.Errorf("expected 100 unique ports, got %d", len(seen))
	}
}
