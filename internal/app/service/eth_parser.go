package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/ethclient"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

var _ contract.ParserService = (*EthereumTxParser)(nil)

type EthereumTxParser struct {
	subscriptionRepo contract.SubscriptionRepositorty
	blockRepo        contract.BlockRepository
	transactionRepo  contract.TransactionRepository
	ethClient        ethclient.Client
	txm              transactionmanager.TransactionManager
}

func NewEthParser(subscriptionRepo contract.SubscriptionRepositorty, blockRepo contract.BlockRepository, transactionRepo contract.TransactionRepository, ethClient ethclient.Client, txm transactionmanager.TransactionManager) *EthereumTxParser {
	return &EthereumTxParser{
		subscriptionRepo: subscriptionRepo,
		blockRepo:        blockRepo,
		transactionRepo:  transactionRepo,
		ethClient:        ethClient,
		txm:              txm,
	}
}

func (p *EthereumTxParser) ParseNextBlock(ctx context.Context) (err error) {
	err = p.txm.Run(ctx, func(ctx context.Context) error {
		height, err := p.blockRepo.GetHeight(ctx, entity.TokenTypeETH)
		if err != nil && !errors.Is(err, entity.ErrBlockHightNotSet) {
			return fmt.Errorf("get current height: %w", err)
		}

		blockResp, err := p.ethClient.GetCurrentBlock(ctx, nil)
		if err != nil {
			return fmt.Errorf("get current block: %w", err)
		}

		if blockResp.Response.Result == nil {
			return fmt.Errorf("current block: %w", entity.ErrCurrentBlockNil)
		}

		if height == "" {
			height = *blockResp.Response.Result

			err = p.ParseBlock(ctx, height)
			if err != nil {
				return fmt.Errorf("initial process block: %w", err)
			}

			return nil
		}

		heightInt, ok := new(big.Int).SetString(height, 0)
		if !ok {
			return fmt.Errorf("parse height: %w", err)
		}

		heightInt = heightInt.Add(heightInt, big.NewInt(1)) // increment the height
		height = heightInt.String()

		err = p.ParseBlock(ctx, height)
		if err != nil {
			return fmt.Errorf("process block: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction manager: %w", err)
	}

	return nil
}

func (p *EthereumTxParser) ParseBlock(ctx context.Context, height string) (err error) {
	bigIntHeight, ok := new(big.Int).SetString(height, 0)
	if !ok {
		return fmt.Errorf("parse height: %w", err) // TODO: return sentitel error
	}

	err = p.txm.Run(ctx, func(ctx context.Context) error {
		var txs []entity.Transaction

		subscriptions, err := p.subscriptionRepo.List(ctx, entity.TokenTypeETH)
		if err != nil {
			return fmt.Errorf("list subscriptions: %w", err)
		}

		subscriptionMap := sliceToMap(subscriptions)
		if len(subscriptions) == 0 {
			return nil
		}

		blockResp, err := p.ethClient.GetBlockByNumber(ctx, &ethclient.GetBlockByNumberRequest{
			Number: fmt.Sprintf("%#x", bigIntHeight),
		})
		if err != nil {
			return fmt.Errorf("get block by number: %w", err)
		}

		// I've omited the block.Result.Withdraws, as I wasn't sure if those count as "transactions" in the context of this task.
		blockTxs := blockResp.Response.Result.Transactions
		for k := range blockTxs {
			addr := ""

			_, ok := subscriptionMap[blockTxs[k].To] // a small optimization by indexing the slice instead of relying on for _, v := range ... which copies each value, and the transaction struct is quite substantial in size
			if ok {
				addr = blockTxs[k].To
			}

			_, ok = subscriptionMap[blockTxs[k].From]
			if ok {
				addr = blockTxs[k].From
			}

			if addr == "" {
				continue
			}

			bigIntValue, ok := new(big.Int).SetString(blockTxs[k].Value, 0)
			if !ok {
				return fmt.Errorf("parse value: %w", err) // TODO: return sentitel error
			}

			txs = append(txs, entity.Transaction{
				ID:        blockTxs[k].Hash, // ideally we should have some sort of unique constraint for the transaction, with the composite of the hash and the token type and address
				TokenType: entity.TokenTypeETH,
				Hash:      blockTxs[k].Hash,
				From:      blockTxs[k].From,
				To:        blockTxs[k].To,
				Value:     bigIntValue.String(),
				Address:   addr,
			})
		}

		if err := p.transactionRepo.BatchCreate(ctx, txs); err != nil {
			return fmt.Errorf("create transactions: %w", err)
		}

		if err := p.blockRepo.SetHeight(ctx, entity.TokenTypeETH, height); err != nil {
			return fmt.Errorf("set height: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction manager: %w", err)
	}

	return nil
}

func sliceToMap[T comparable](s []T) map[T]struct{} {
	m := make(map[T]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}
