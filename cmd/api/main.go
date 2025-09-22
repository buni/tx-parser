package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	httpin_integration "github.com/ggicci/httpin/integration" //nolint
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

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
	"github.com/buni/tx-parser/internal/pkg/testing/testutils"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
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
	ethClient := ethclient.NewClient(srv.Config.RPCEndpoint)
	noopTxm := transactionmanager.NewNoopTxm()
	noopPublisher := testutils.NewNoopPublisher()
	blockSvc := service.NewBlockService(blockRepo, noopTxm)
	blockHandler := handler.NewBlockHandler(blockSvc)

	transactionSvc := service.NewTransactionService(transactionRepo, noopPublisher, noopTxm)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

	subscriptionSvc := service.NewSubscriberService(subscriptionRepo, noopTxm)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionSvc)

	ethParser := service.NewEthParser(subscriptionRepo, blockRepo, transactionSvc, ethClient, noopTxm)
	ehtBlockParserTask := task.NewEthBlockParser(ethParser)

	err = seedLocalEnv(ctx, srv, blockRepo, subscriptionRepo)
	if err != nil {
		return fmt.Errorf("seed local env: %w", err)
	}

	srv.Router.Route("/v1/{tokenType}", func(r chi.Router) {
		blockHandler.RegisterRoutes(r)
		transactionHandler.RegisterRoutes(r)
		subscriptionHandler.RegisterRoutes(r)
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

func seedLocalEnv(ctx context.Context, srv *server.Server, blockRepo contract.BlockRepository, subscriptionRepo contract.SubscriptionRepository) error {
	if srv.Config.Environment != configuration.EnvLocal {
		return nil
	}

	if srv.Config.InitialHeight != "" {
		err := blockRepo.SetHeight(ctx, entity.TokenTypeETH, srv.Config.InitialHeight)
		if err != nil {
			return fmt.Errorf("set initial height: %w", err)
		}
	}

	for _, addr := range srv.Config.SeedAddresses {
		sub, err := entity.NewSubscription(
			entity.TokenTypeETH,
			uuid.Must(uuid.NewV7()).String(),
			addr,
		)
		if err != nil {
			return fmt.Errorf("new subscription: %w", err)
		}

		err = subscriptionRepo.Create(ctx, sub)
		if err != nil {
			return fmt.Errorf("create subscription: %w", err)
		}
	}

	newAddresses, err := generateEthAddress(500000)
	if err != nil {
		return fmt.Errorf("generate eth addresses: %w", err)
	}

	for addr := range newAddresses {
		sub, err := entity.NewSubscription(
			entity.TokenTypeETH,
			uuid.Must(uuid.NewV7()).String(),
			addr,
		)
		if err != nil {
			return fmt.Errorf("new subscription: %w", err)
		}

		err = subscriptionRepo.Create(ctx, sub)
		if err != nil {
			return fmt.Errorf("create subscription: %w", err)
		}
	}

	return nil
}

func generateEthAddress(n int) (map[string]struct{}, error) {
	if n <= 0 {
		n = 1
	}

	addresses := make(map[string]struct{}, n)

	for i := 0; i < n; i++ {
		bytes := make([]byte, 20)

		_, err := rand.Read(bytes)
		if err != nil {
			return nil, fmt.Errorf("rand read: %w", err)
		}

		addresses[fmt.Sprintf("0x%040x", bytes)] = struct{}{}
	}

	return addresses, nil
}
