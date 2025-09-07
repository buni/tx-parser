package service

import (
	"context"
	"fmt"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/dto"
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

func (p *BlockService) GetCurrentBlock(ctx context.Context, req *dto.GetCurrentBlockRequest) (height string, err error) {
	err = p.txm.Run(ctx, func(ctx context.Context) error {
		height, err = p.blockRepo.GetHeight(ctx, req.TokenType)
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
