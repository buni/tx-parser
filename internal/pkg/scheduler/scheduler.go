package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Task interface {
	Name() string
	Handle(ctx context.Context) error
	Interval() time.Duration
}

type PollingScheduler struct {
	logger       *slog.Logger
	tickInterval time.Duration
	errsChan     chan error
	wg           sync.WaitGroup
}

func NewPollingScheduler(tickInterval time.Duration, logger *slog.Logger) *PollingScheduler {
	return &PollingScheduler{
		logger:       logger,
		tickInterval: tickInterval,
		errsChan:     make(chan error),
		wg:           sync.WaitGroup{},
	}
}

func (s *PollingScheduler) Start(ctx context.Context, tasks ...Task) error {
	for _, task := range tasks {
		s.wg.Add(1)
		go s.processTask(ctx, task)
	}

	return nil
}

func (s *PollingScheduler) processTask(ctx context.Context, task Task) {
	logger := s.logger.With(slog.String("task_name", task.Name()))
	defer s.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(task.Interval()):
			logger.InfoContext(ctx, "starting task execution")

			err := task.Handle(ctx)
			if err != nil {
				logger.ErrorContext(ctx, "task execution failed", slog.String("error", err.Error()))
				continue
			}
		}
		logger.InfoContext(ctx, "task execution completed")
	}
}

func (s *PollingScheduler) Wait() {
	s.wg.Wait()
}
