package repository

import (
	"context"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/syncx"
)

var _ contract.SubscriptionRepositorty = (*InMemorySubscriptionRepository)(nil)

type InMemorySubscriptionRepository struct {
	m *syncx.SyncMap[entity.TokenType, []string]
}

func NewInMemorySubscriptionRepository() *InMemorySubscriptionRepository {
	return &InMemorySubscriptionRepository{
		m: &syncx.SyncMap[entity.TokenType, []string]{},
	}
}

func (r *InMemorySubscriptionRepository) Create(_ context.Context, tokenType entity.TokenType, address string) error {
	addresses, ok := r.m.Load(tokenType)
	if !ok {
		addresses = make([]string, 0, 1) // technically this is not needed, but I prefer to be explicit
	}

	addresses = append(addresses, address)
	r.m.Store(tokenType, addresses)
	return nil
}

func (r *InMemorySubscriptionRepository) List(_ context.Context, tokenType entity.TokenType) ([]string, error) {
	addresses, ok := r.m.Load(tokenType)
	if !ok {
		return nil, entity.ErrNotFound
	}

	return addresses, nil
}
