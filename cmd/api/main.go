package main

import (
	"context"
	"fmt"
	"time"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/app/handler"
	"github.com/buni/tx-parser/internal/app/repository"
	"github.com/buni/tx-parser/internal/app/service"
	"github.com/buni/tx-parser/internal/app/task"
	"github.com/buni/tx-parser/internal/pkg/configuration"
	"github.com/buni/tx-parser/internal/pkg/ethclient"
	"github.com/buni/tx-parser/internal/pkg/scheduler"
	"github.com/buni/tx-parser/internal/pkg/server"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
	httpin_integration "github.com/ggicci/httpin/integration" //nolint
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		zap.L().Fatal("failed to start service", zap.Error(err))
	}
}

func run() error {
	httpin_integration.UseGochiURLParam("path", chi.URLParam)

	ctx := context.Background()

	srv, err := server.NewServer(ctx)
	if err != nil {
		return fmt.Errorf("new server: %w", err)
	}

	taskScheduler := scheduler.NewPollingScheduler(time.Second*1, srv.Logger)

	blockRepo := repository.NewInMemoryBlockRepository()
	transactionRepo := repository.NewInMemoryTransactionRepository()
	subscriptionRepo := repository.NewInMemorySubscriptionRepository()
	ethClient := ethclient.NewClient(srv.Config.Ethereum.RPCEndpoint)
	noopTxm := transactionmanager.NewNoopTxm()

	blockSvc := service.NewBlockService(blockRepo, noopTxm)
	blockHandler := handler.NewBlockHandler(blockSvc)

	transactionSvc := service.NewTransactionService(transactionRepo, noopTxm)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

	subscriptionSvc := service.NewSubscriberService(subscriptionRepo, noopTxm)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionSvc)

	ethParser := service.NewEthParser(subscriptionRepo, blockRepo, transactionRepo, ethClient, noopTxm)
	ehtBlockParserTask := task.NewEthBlockParser(ethParser)

	err = seedLocalEnv(ctx, srv, blockRepo, subscriptionRepo)
	if err != nil {
		return fmt.Errorf("seed local env: %w", err)
	}

	srv.Router.Route("/v1", func(r chi.Router) {
		r.Route("/eth", func(r chi.Router) {
			blockHandler.RegisterRoutes(r)
			transactionHandler.RegisterRoutes(r)
			subscriptionHandler.RegisterRoutes(r)
		})
	})

	err = taskScheduler.Start(ctx, ehtBlockParserTask)
	if err != nil {
		return fmt.Errorf("start task scheduler: %w", err)
	}

	if err := srv.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	srv.Wait()

	return nil
}

func seedLocalEnv(ctx context.Context, srv *server.Server, blockRepo contract.BlockRepository, subscriptionRepo contract.SubscriptionRepositorty) error {
	if srv.Config.Service.Environment != configuration.EnvLocal {
		return nil
	}

	if srv.Config.Ethereum.InitialHeight != "" {
		err := blockRepo.SetHeight(ctx, entity.TokenTypeETH, srv.Config.Ethereum.InitialHeight)
		if err != nil {
			return fmt.Errorf("set initial height: %w", err)
		}
	}

	for _, address := range srv.Config.Ethereum.SeedAddresses {
		err := subscriptionRepo.Create(ctx, entity.TokenTypeETH, address)
		if err != nil {
			return fmt.Errorf("create subscription: %w", err)
		}
	}

	return nil
}
