package repository

import (
	"context"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/syncx"
)

var _ contract.BlockRepository = (*InMemoryBlockRepository)(nil)

type InMemoryBlockRepository struct {
	m *syncx.SyncMap[entity.TokenType, string]
}

func NewInMemoryBlockRepository() *InMemoryBlockRepository {
	return &InMemoryBlockRepository{
		m: &syncx.SyncMap[entity.TokenType, string]{},
	}
}

func (r *InMemoryBlockRepository) SetHeight(_ context.Context, tokenType entity.TokenType, height string) error {
	r.m.Store(tokenType, height)
	return nil
}

func (r *InMemoryBlockRepository) GetHeight(_ context.Context, tokenType entity.TokenType) (string, error) {
	height, ok := r.m.Load(tokenType)
	if height == "" || !ok {
		return "", entity.ErrBlockHightNotSet
	}

	return height, nil
}
