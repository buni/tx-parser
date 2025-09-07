package transactionmanager

import "context"

type NoopTxm struct{}

func NewNoopTxm() *NoopTxm {
	return &NoopTxm{}
}

// Run executes the given function.
// It is a no-op implementation of the TransactionManager interface, and as such doesn't rollback any changes on failure.
func (n *NoopTxm) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
