package contract

import (
	"context"

	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
)

//go:generate go tool mockgen -source=subscription.go -destination=mock/subscription_mocks.go -package contract_mock

type SubscriptionRepository interface {
	Create(ctx context.Context, sub entity.Subscription) error
	ListByAddresses(ctx context.Context, tokenType entity.TokenType, addresses []string) ([]entity.Subscription, error)
}

type SubscriptionService interface {
	Subscribe(ctx context.Context, req *dto.SubscribeRequest) error
}
