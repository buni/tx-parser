package testutils

import (
	"context"

	"github.com/buni/tx-parser/internal/pkg/pubsub"
)

// NoopTransactionManager implements database.TransactionManager.
type NoopPublisher struct{}

func NewNoopPublisher() NoopPublisher {
	return NoopPublisher{}
}

func (n NoopPublisher) Publish(_ context.Context, _ *pubsub.Message, _ ...pubsub.PublishOption) error {
	return nil
}
