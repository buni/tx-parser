package repository

import (
	"context"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/syncx"
)

var _ contract.TransactionRepository = (*InMemoryTransactionRepository)(nil)

type InMemoryTransactionRepository struct {
	m *syncx.SyncMap[entity.TokenType, *syncx.SyncMap[string, []entity.Transaction]]
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		m: &syncx.SyncMap[entity.TokenType, *syncx.SyncMap[string, []entity.Transaction]]{},
	}
}

func (r *InMemoryTransactionRepository) Create(_ context.Context, tx entity.Transaction) error {
	txMap, ok := r.m.Load(tx.TokenType)
	if !ok {
		txMap = &syncx.SyncMap[string, []entity.Transaction]{}
	}

	currentTxs, ok := txMap.Load(tx.Address)
	if !ok {
		currentTxs = make([]entity.Transaction, 0, 1)
	}

	currentTxs = append(currentTxs, tx)

	txMap.Store(tx.Address, currentTxs)
	r.m.Store(tx.TokenType, txMap)

	return nil
}

func (r *InMemoryTransactionRepository) List(_ context.Context, tokenType entity.TokenType, address string) ([]entity.Transaction, error) {
	txMap, ok := r.m.Load(tokenType)
	if !ok {
		return []entity.Transaction{}, nil
	}

	txs, ok := txMap.Load(address)
	if !ok {
		return []entity.Transaction{}, nil
	}

	return txs, nil
}
