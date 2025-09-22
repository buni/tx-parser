package repository

import (
	"context"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/syncx"
)

var _ contract.SubscriptionRepository = (*InMemorySubscriptionRepository)(nil)

type InMemorySubscriptionRepository struct {
	m *syncx.SyncMap[entity.TokenType, *syncx.SyncMap[string, entity.Subscription]]
}

func NewInMemorySubscriptionRepository() *InMemorySubscriptionRepository {
	return &InMemorySubscriptionRepository{
		m: &syncx.SyncMap[entity.TokenType, *syncx.SyncMap[string, entity.Subscription]]{},
	}
}

func (r *InMemorySubscriptionRepository) Create(_ context.Context, sub entity.Subscription) error {
	tokenMap, _ := r.m.LoadOrStore(sub.TokenType, &syncx.SyncMap[string, entity.Subscription]{})
	tokenMap.Store(sub.Address, sub)
	return nil
}

func (r *InMemorySubscriptionRepository) ListByAddresses(_ context.Context, tokenType entity.TokenType, addresses []string) (result []entity.Subscription, err error) {
	if len(addresses) == 0 {
		return nil, nil
	}

	tokenMap, ok := r.m.Load(tokenType)
	if !ok {
		return nil, nil
	}

	for _, address := range addresses {
		if sub, ok := tokenMap.Load(address); ok {
			result = append(result, sub)
		}
	}

	return result, nil
}
