package transactionmanager

import "context"

// TransactionManager is defines the contract for implementing a unit of work pattern.
type TransactionManager interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
