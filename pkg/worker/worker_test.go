package worker

import (
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(2, 1)
	defer pool.Stop()

	executed := false
	err := pool.Submit(func() {
		executed = true
	})

	if err != nil {
		t.Errorf("Submit() failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	if !executed {
		t.Error("Task was not executed")
	}
}

func TestWorkerPoolMultipleTasks(t *testing.T) {
	pool := NewWorkerPool(10, 2)
	defer pool.Stop()

	var mu sync.Mutex
	executed := 0

	for i := 0; i < 5; i++ {
		err := pool.Submit(func() {
			mu.Lock()
			executed++
			mu.Unlock()
		})
		if err != nil {
			t.Errorf("Submit() failed: %v", err)
		}
	}

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	if executed != 5 {
		t.Errorf("Expected 5 tasks executed, got %d", executed)
	}
	mu.Unlock()
}

func TestWorkerPoolQueueOverflow(t *testing.T) {
	pool := NewWorkerPool(1, 1)
	defer pool.Stop()

	blocked := make(chan struct{})
	err := pool.Submit(func() {
		<-blocked
	})
	if err != nil {
		t.Errorf("First Submit() failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	err = pool.Submit(func() {})
	if err != nil {
		t.Errorf("Second Submit() failed: %v", err)
	}

	err = pool.Submit(func() {})
	if err != ErrQueueOverflow {
		t.Errorf("Expected ErrQueueOverflow, got %v", err)
	}

	close(blocked)
}

func TestWorkerPoolStop(t *testing.T) {
	pool := NewWorkerPool(2, 1)

	err := pool.Submit(func() {
		time.Sleep(100 * time.Millisecond)
	})
	if err != nil {
		t.Errorf("Submit() failed: %v", err)
	}

	err = pool.Stop()
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}

	err = pool.Submit(func() {})
	if err != ErrPoolStopped {
		t.Errorf("Expected ErrPoolStopped, got %v", err)
	}
}

func TestWorkerPoolStopMultiple(t *testing.T) {
	pool := NewWorkerPool(2, 1)

	err1 := pool.Stop()
	err2 := pool.Stop()

	if err1 != nil {
		t.Errorf("First Stop() failed: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Second Stop() failed: %v", err2)
	}
}
