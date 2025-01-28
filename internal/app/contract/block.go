package contract

import (
	"context"

	"github.com/buni/tx-parser/internal/app/entity"
)

//go:generate mockgen -source=block.go -destination=mock/block_mocks.go -package contract_mock

type BlockRepository interface {
	SetHeight(ctx context.Context, tokenType entity.TokenType, height string) error
	GetHeight(ctx context.Context, tokenType entity.TokenType) (string, error)
}

type BlockService interface {
	GetCurrentBlock(ctx context.Context) (height string, err error)
}
