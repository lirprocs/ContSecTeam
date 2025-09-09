package service

import (
	"ContSecTeam/internal/model"
	"ContSecTeam/pkg"
	"ContSecTeam/pkg/worker"
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Service struct {
	p         worker.Pool
	queueSize int
	Store     sync.Map
	ctx       context.Context
}

func NewService(queueSize int) *Service {
	return &Service{queueSize: queueSize}
}

func (s *Service) Start(ctx context.Context, workers int) {
	s.ctx = ctx
	s.p = worker.NewWorkerPool(s.queueSize, workers)
}

func (s *Service) Stop() {
	if s.p != nil {
		_ = s.p.Stop()
	}
}

func (s *Service) Enqueue(task *model.Task) error {
	if task == nil {
		return errors.New("nil task")
	}
	task.Status = model.StatusQueued
	s.Store.Store(task.ID, task.Status)

	return s.p.Submit(func() {
		s.processTask(s.ctx, task)
	})
}

func (s *Service) processTask(ctx context.Context, task *model.Task) {
	for {
		task.Attempts++
		task.Status = model.StatusRunning
		s.Store.Store(task.ID, task.Status)

		pkg.SleepRandom()

		if pkg.ShouldFail() {
			if task.Attempts <= task.MaxRetries {
				select {
				case <-ctx.Done():
					return
				case <-time.After(pkg.Backoff(task.Attempts)):
				}
				continue
			}
			task.Status = model.StatusFailed
			s.Store.Store(task.ID, task.Status)
			log.Printf("Task %s failed", task.ID)
			return
		}

		task.Status = model.StatusDone
		s.Store.Store(task.ID, task.Status)
		log.Printf("Task %s done", task.ID)
		return
	}
}
