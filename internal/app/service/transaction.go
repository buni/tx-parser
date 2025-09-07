package service

import (
	"context"
	"fmt"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/pubsub"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

var _ contract.TransactionService = (*TransactionService)(nil)

type TransactionService struct {
	transactionRepo contract.TransactionRepository
	publisher       pubsub.Publisher
	txm             transactionmanager.TransactionManager
}

func NewTransactionService(
	transactionRepo contract.TransactionRepository,
	publisher pubsub.Publisher,
	txm transactionmanager.TransactionManager,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		publisher:       publisher,
		txm:             txm,
	}
}

func (s *TransactionService) ListAddressTransactions(ctx context.Context, req *dto.ListAddressTransactionsRequest) (result []entity.Transaction, err error) {
	err = s.txm.Run(ctx, func(ctx context.Context) error {
		result, err = s.transactionRepo.List(ctx, req.TokenType, req.Address)
		if err != nil {
			return fmt.Errorf("list transactions: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction manager: %w", err)
	}

	return result, nil
}

func (s *TransactionService) Create(ctx context.Context, tx entity.Transaction) (err error) {
	err = s.txm.Run(ctx, func(ctx context.Context) error {
		if err = s.transactionRepo.Create(ctx, tx); err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		event, err := entity.NewTransactionCreatedEvent(tx)
		if err != nil {
			return fmt.Errorf("create event: %w", err)
		}

		if err = s.publisher.Publish(ctx, event); err != nil {
			return fmt.Errorf("publish: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction manager: %w", err)
	}

	return nil
}
