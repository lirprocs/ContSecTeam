package service

import (
	"ContSecTeam/internal/model"
	"ContSecTeam/pkg/worker"
	"context"
	"errors"
	"testing"
	"time"
)

func TestService_EnqueueStoresQueued(t *testing.T) {
	s := NewService(1)
	s.Start(context.Background(), 1)
	t.Cleanup(func() { s.Stop() })

	task := &model.Task{ID: "t1", MaxRetries: 0}
	if err := s.Enqueue(task); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}

	if v, ok := s.Store.Load("t1"); !ok {
		t.Fatalf("task not found in Store")
	} else {
		status := v.(string)
		if status != model.StatusQueued && status != model.StatusRunning && status != model.StatusDone {
			t.Fatalf("unexpected status: %v", status)
		}
	}
}

func TestService_EnqueueQueueOverflow(t *testing.T) {
	s := NewService(1)
	s.Start(context.Background(), 1)
	t.Cleanup(func() { s.Stop() })

	if err := s.Enqueue(&model.Task{ID: "t1", MaxRetries: 0}); err != nil {
		t.Fatalf("first Enqueue() error = %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := s.Enqueue(&model.Task{ID: "t2", MaxRetries: 0}); err != nil {
		t.Fatalf("second Enqueue() error = %v", err)
	}

	err := s.Enqueue(&model.Task{ID: "t3", MaxRetries: 0})
	if !errors.Is(err, worker.ErrQueueOverflow) {
		t.Fatalf("expected ErrQueueOverflow, got %v", err)
	}
}

func TestService_StopPreventsSubmit(t *testing.T) {
	s := NewService(1)
	s.Start(context.Background(), 0)
	s.Stop()

	err := s.Enqueue(&model.Task{ID: "t1"})
	if !errors.Is(err, worker.ErrPoolStopped) {
		t.Fatalf("expected ErrPoolStopped, got %v", err)
	}
}

func TestService_ProcessEventuallyUpdatesStatus(t *testing.T) {
	s := NewService(4)
	s.Start(context.Background(), 1)
	t.Cleanup(func() { s.Stop() })

	task := &model.Task{ID: "tX", MaxRetries: 1}
	if err := s.Enqueue(task); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if v, ok := s.Store.Load("tX"); ok {
			status := v.(string)
			if status == model.StatusDone || status == model.StatusFailed || status == model.StatusRunning {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("status did not progress in time")
}
