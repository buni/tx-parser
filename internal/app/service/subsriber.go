package service

import (
	"context"
	"fmt"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

type SubscriberService struct {
	subscriptionRepo contract.SubscriptionRepository
	txm              transactionmanager.TransactionManager
}

func NewSubscriberService(subscriptionRepo contract.SubscriptionRepository, txm transactionmanager.TransactionManager) *SubscriberService {
	return &SubscriberService{
		subscriptionRepo: subscriptionRepo,
		txm:              txm,
	}
}

func (p *SubscriberService) Subscribe(ctx context.Context, req *dto.SubscribeRequest) (err error) {
	err = p.txm.Run(ctx, func(ctx context.Context) error {
		sub, err := entity.NewSubscription(
			req.TokenType,
			req.UserID,
			req.Address,
		)
		if err != nil {
			return fmt.Errorf("new subscription: %w", err)
		}

		if err = p.subscriptionRepo.Create(ctx, sub); err != nil {
			return fmt.Errorf("create subscription: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction manager: %w", err)
	}

	return nil
}
