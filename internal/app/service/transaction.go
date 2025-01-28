package service

import (
	"context"
	"fmt"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

var _ contract.TransactionService = (*TransactionService)(nil)

type TransactionService struct {
	transactionRepo contract.TransactionRepository
	txm             transactionmanager.TransactionManager
}

func NewTransactionService(
	transactionRepo contract.TransactionRepository,
	txm transactionmanager.TransactionManager,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		txm:             txm,
	}
}

func (p *TransactionService) ListAddressTransactions(ctx context.Context, req *dto.ListAddressTransactionsRequest) (txs []entity.Transaction, err error) {
	err = p.txm.Run(ctx, func(ctx context.Context) error {
		// TODO: check if address exists?
		txs, err = p.transactionRepo.List(ctx, entity.TokenTypeETH, req.Address)
		if err != nil {
			return fmt.Errorf("list transactions: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction manager: %w", err)
	}

	return txs, nil
}
