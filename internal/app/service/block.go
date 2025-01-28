package service

import (
	"context"
	"fmt"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

type BlockService struct {
	blockRepo contract.BlockRepository
	txm       transactionmanager.TransactionManager
}

func NewBlockService(blockRepo contract.BlockRepository, txm transactionmanager.TransactionManager) *BlockService {
	return &BlockService{
		blockRepo: blockRepo,
		txm:       txm,
	}
}

func (p *BlockService) GetCurrentBlock(ctx context.Context) (height string, err error) {
	err = p.txm.Run(ctx, func(ctx context.Context) error {
		height, err = p.blockRepo.GetHeight(ctx, entity.TokenTypeETH)
		if err != nil {
			return fmt.Errorf("get current block: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("transaction manager: %w", err)
	}

	return height, nil
}
