package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

type SubscriberService struct {
	subscriptionRepo contract.SubscriptionRepositorty
	txm              transactionmanager.TransactionManager
}

func NewSubscriberService(subscriptionRepo contract.SubscriptionRepositorty, txm transactionmanager.TransactionManager) *SubscriberService {
	return &SubscriberService{
		subscriptionRepo: subscriptionRepo,
		txm:              txm,
	}
}

func (p *SubscriberService) Subscribe(ctx context.Context, req *dto.SubscribeRequest) (err error) {
	err = p.txm.Run(ctx, func(ctx context.Context) error {
		err = p.subscriptionRepo.Create(ctx, entity.TokenTypeETH, strings.ToLower(req.Address))
		if err != nil {
			return fmt.Errorf("create subscription: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction manager: %w", err)
	}

	return nil
}
