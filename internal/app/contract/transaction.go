package contract

import (
	"context"

	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
)

//go:generate go tool mockgen -source=transaction.go -destination=mock/transaction_mocks.go -package contract_mock

type TransactionRepository interface {
	Create(ctx context.Context, tx entity.Transaction) error
	List(ctx context.Context, tokenType entity.TokenType, address string) ([]entity.Transaction, error)
}

type TransactionService interface {
	ListAddressTransactions(ctx context.Context, req *dto.ListAddressTransactionsRequest) (txs []entity.Transaction, err error)
	Create(ctx context.Context, tx entity.Transaction) error
}
