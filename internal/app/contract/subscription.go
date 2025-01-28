package contract

import (
	"context"

	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
)

//go:generate mockgen -source=subscription.go -destination=mock/subscription_mocks.go -package contract_mock

type SubscriptionRepositorty interface {
	Create(ctx context.Context, tokenType entity.TokenType, address string) error
	List(ctx context.Context, tokenType entity.TokenType) ([]string, error)
}

type SubscriptionService interface {
	Subscribe(ctx context.Context, req *dto.SubscribeRequest) error
}
