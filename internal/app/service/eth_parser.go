package service

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math/big"
	"slices"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/pkg/ethclient"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

var _ contract.ParserService = (*EthereumTxParser)(nil)

type EthereumTxParser struct {
	subscriptionRepo contract.SubscriptionRepository
	blockRepo        contract.BlockRepository
	transactionSvc   contract.TransactionService
	ethClient        ethclient.Client
	txm              transactionmanager.TransactionManager
}

func NewEthParser(subscriptionRepo contract.SubscriptionRepository, blockRepo contract.BlockRepository, transactionSvc contract.TransactionService, ethClient ethclient.Client, txm transactionmanager.TransactionManager) *EthereumTxParser {
	return &EthereumTxParser{
		subscriptionRepo: subscriptionRepo,
		blockRepo:        blockRepo,
		transactionSvc:   transactionSvc,
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

			if err = p.ParseBlock(ctx, height); err != nil {
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

		if err = p.ParseBlock(ctx, height); err != nil {
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
	tokenType := entity.TokenTypeETH

	bigIntHeight, ok := new(big.Int).SetString(height, 0)
	if !ok {
		return fmt.Errorf("parse height: %w", err)
	}

	err = p.txm.Run(ctx, func(ctx context.Context) error {
		uniqAddrs := make(map[string]struct{}, 300)
		blockResp, err := p.ethClient.GetBlockByNumber(ctx, &ethclient.GetBlockByNumberRequest{
			Number: fmt.Sprintf("%#x", bigIntHeight),
		})
		if err != nil {
			return fmt.Errorf("get block by number: %w", err)
		}

		for k := range blockResp.Response.Result.Transactions {
			uniqAddrs[blockResp.Response.Result.Transactions[k].To] = struct{}{}
			uniqAddrs[blockResp.Response.Result.Transactions[k].From] = struct{}{}
		}

		addrs := slices.AppendSeq(make([]string, 0, len(uniqAddrs)), maps.Keys(uniqAddrs))

		subscriptions, err := p.subscriptionRepo.ListByAddresses(ctx, entity.TokenTypeETH, addrs)
		if err != nil {
			return fmt.Errorf("list subscriptions: %w", err)
		}

		if len(subscriptions) == 0 {
			if err := p.blockRepo.SetHeight(ctx, tokenType, height); err != nil {
				return fmt.Errorf("set height: %w", err)
			}
			return nil
		}

		subscriptionMap := make(map[string]entity.Subscription, len(subscriptions))
		for _, v := range subscriptions {
			subscriptionMap[v.Address] = v
		}

		// I've omited the block.Result.Withdraws, as I wasn't sure if those count as "transactions" in the context of this task.
		blockTxs := blockResp.Response.Result.Transactions
		for k := range blockTxs {
			subTo, ok := subscriptionMap[blockTxs[k].To]
			if ok {
				if err := p.parseAndStoreTransaction(ctx, blockTxs[k], blockTxs[k].To, subTo, entity.TransactionTypeDebit); err != nil {
					return fmt.Errorf("parse and store transaction sender: %w", err)
				}
			}

			subFrom, ok := subscriptionMap[blockTxs[k].From]
			if ok {
				if err := p.parseAndStoreTransaction(ctx, blockTxs[k], blockTxs[k].From, subFrom, entity.TransactionTypeCredit); err != nil {
					return fmt.Errorf("parse and store transaction sender: %w", err)
				}
			}
		}

		if err := p.blockRepo.SetHeight(ctx, tokenType, height); err != nil {
			return fmt.Errorf("set height: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction manager: %w", err)
	}

	return nil
}

func (p *EthereumTxParser) parseAndStoreTransaction(ctx context.Context, blockTx ethclient.Transaction, addr string, sub entity.Subscription, txType entity.TransactionType) (err error) {
	if addr == "" {
		return nil
	}

	if blockTx.Value == nil || blockTx.GasPrice == nil {
		return nil
	}

	tx, err := entity.NewTransaction(entity.TokenTypeETH, txType, sub.UserID, addr, blockTx.To, blockTx.From, blockTx.Hash, blockTx.Value.String(), blockTx.BlockNumber)
	if err != nil {
		return fmt.Errorf("new transaction: %w", err)
	}

	if err := p.transactionSvc.Create(ctx, tx); err != nil {
		return fmt.Errorf("create transactions: %w", err)
	}
	return nil
}
