package scheduler_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/buni/tx-parser/internal/pkg/scheduler"
	"github.com/stretchr/testify/suite"
)

type PollingSchedulerTestSuite struct {
	suite.Suite
	logger    *slog.Logger
	scheduler *scheduler.PollingScheduler
}

func (s *PollingSchedulerTestSuite) SetupTest() {
	s.logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	s.scheduler = scheduler.NewPollingScheduler(time.Millisecond*10, s.logger)
}

func (s *PollingSchedulerTestSuite) TestStartAndProcessTaskSuccess() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task := &mockTask{
		name:     "mockTask",
		interval: time.Millisecond * 10,
		handle: func(ctx context.Context) error {
			return nil
		},
	}

	err := s.scheduler.Start(ctx, task)
	s.NoError(err)

	// Allow some time for the task to be processed
	time.Sleep(time.Millisecond * 50) // TODO: replace with chan blocking
}

func (s *PollingSchedulerTestSuite) TestStartAndProcessTaskError() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task := &mockTask{
		name:     "mockTask",
		interval: time.Millisecond * 10,
		handle: func(ctx context.Context) error {
			return errors.New("task error")
		},
	}

	err := s.scheduler.Start(ctx, task)
	s.NoError(err)

	// Allow some time for the task to be processed
	time.Sleep(time.Millisecond * 50)
}

func (s *PollingSchedulerTestSuite) TestWait() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task := &mockTask{
		name:     "mockTask",
		interval: time.Millisecond * 10,
		handle: func(ctx context.Context) error {
			return nil
		},
	}

	err := s.scheduler.Start(ctx, task)
	s.NoError(err)

	// Allow some time for the task to be processed
	time.Sleep(time.Millisecond * 50)

	cancel()
	s.scheduler.Wait()
}

func TestPollingSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(PollingSchedulerTestSuite))
}

type mockTask struct {
	name     string
	interval time.Duration
	handle   func(ctx context.Context) error
}

func (m *mockTask) Name() string {
	return m.name
}

func (m *mockTask) Handle(ctx context.Context) error {
	return m.handle(ctx)
}

func (m *mockTask) Interval() time.Duration {
	return m.interval
}
