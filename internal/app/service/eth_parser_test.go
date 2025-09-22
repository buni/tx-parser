package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	contract_mock "github.com/buni/tx-parser/internal/app/contract/mock"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/app/repository"
	"github.com/buni/tx-parser/internal/app/service"
	"github.com/buni/tx-parser/internal/pkg/ethclient"
	ethclient_mock "github.com/buni/tx-parser/internal/pkg/ethclient/mock"
	"github.com/buni/tx-parser/internal/pkg/testing/testutils"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

func BenchmarkParseBlock(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	ctx := context.Background()
	blockRepo := repository.NewInMemoryBlockRepository()
	transactionRepo := repository.NewInMemoryTransactionRepository()
	transactionSvc := service.NewTransactionService(transactionRepo, testutils.NewNoopPublisher(), transactionmanager.NewNoopTxm())
	subscriptionRepo := repository.NewInMemorySubscriptionRepository()
	ethClient := ethclient.NewClient("https://ethereum-rpc.publicnode.com/")
	noopTxm := transactionmanager.NewNoopTxm()
	sub, err := entity.NewSubscription(entity.TokenTypeETH, "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70", "user1")
	if err != nil {
		zap.L().Fatal("failed to create subscription", zap.Error(err))
	}

	err = subscriptionRepo.Create(ctx, sub)
	if err != nil {
		zap.L().Fatal("failed to create subscription", zap.Error(err))
	}

	err = blockRepo.SetHeight(ctx, entity.TokenTypeETH, "")
	if err != nil {
		zap.L().Fatal("failed to set initial height", zap.Error(err))
	}

	ethParserSvc := service.NewEthParser(subscriptionRepo, blockRepo, transactionSvc, ethClient, noopTxm)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err = ethParserSvc.ParseBlock(ctx, "0x14b66a0")
		if err != nil {
			b.Fatal(err)
		}
	}
}

type EthereumTxParserTestSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	mockSubscription *contract_mock.MockSubscriptionRepository
	mockBlock        *contract_mock.MockBlockRepository
	mockTransaction  *contract_mock.MockTransactionService
	mockEthClient    *ethclient_mock.MockClient
	mockTxManager    *transactionmanager.NoopTxm
	parser           *service.EthereumTxParser
	ctx              context.Context
}

func (s *EthereumTxParserTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockSubscription = contract_mock.NewMockSubscriptionRepository(s.ctrl)
	s.mockBlock = contract_mock.NewMockBlockRepository(s.ctrl)
	s.mockTransaction = contract_mock.NewMockTransactionService(s.ctrl)
	s.mockEthClient = ethclient_mock.NewMockClient(s.ctrl)
	s.mockTxManager = &transactionmanager.NoopTxm{}
	s.parser = service.NewEthParser(s.mockSubscription, s.mockBlock, s.mockTransaction, s.mockEthClient, s.mockTxManager)
	s.ctx = context.Background()
}

func (s *EthereumTxParserTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *EthereumTxParserTestSuite) TestParseNextBlockSuccess() {
	s.mockBlock.EXPECT().GetHeight(s.ctx, entity.TokenTypeETH).Return("", entity.ErrBlockHightNotSet)
	s.mockEthClient.EXPECT().GetCurrentBlock(s.ctx, nil).Return(&ethclient.GetCurrentBlockResponse{Response: ethclient.Response[string]{
		Result: testutils.ToPtr("123456"),
	}}, nil)
	s.mockSubscription.EXPECT().ListByAddresses(s.ctx, entity.TokenTypeETH, gomock.InAnyOrder([]string{"to1", "from1"})).Return([]entity.Subscription{{
		Address: "from1",
	}}, nil)
	s.mockEthClient.EXPECT().GetBlockByNumber(s.ctx, &ethclient.GetBlockByNumberRequest{
		Number: "0x1e240",
	}).Return(&ethclient.GetBlockByNumberResponse{
		Response: ethclient.Response[ethclient.Block]{
			Result: &ethclient.Block{
				Transactions: []ethclient.Transaction{
					{
						Hash:     "hash1",
						From:     "from1",
						To:       "to1",
						GasPrice: ethclient.NewHexBigInt(big.NewInt(1)),
						Value:    ethclient.NewHexBigInt(big.NewInt(100)),
					},
				},
			},
		},
	}, nil)
	s.mockTransaction.EXPECT().Create(s.ctx, gomock.Any()).Return(nil)

	s.mockBlock.EXPECT().SetHeight(s.ctx, entity.TokenTypeETH, "123456").Return(nil)

	err := s.parser.ParseNextBlock(s.ctx)

	s.NoError(err)
}

func (s *EthereumTxParserTestSuite) TestParseNextBlockExistingBlockSuccess() {
	s.mockBlock.EXPECT().GetHeight(s.ctx, entity.TokenTypeETH).Return("123455", nil)
	s.mockEthClient.EXPECT().GetCurrentBlock(s.ctx, nil).Return(&ethclient.GetCurrentBlockResponse{Response: ethclient.Response[string]{
		Result: testutils.ToPtr("123456"),
	}}, nil)
	s.mockSubscription.EXPECT().ListByAddresses(s.ctx, entity.TokenTypeETH, gomock.InAnyOrder([]string{"to1", "from1"})).Return([]entity.Subscription{{
		Address: "from1",
	}}, nil)
	s.mockEthClient.EXPECT().GetBlockByNumber(s.ctx, &ethclient.GetBlockByNumberRequest{
		Number: "0x1e240",
	}).Return(&ethclient.GetBlockByNumberResponse{
		Response: ethclient.Response[ethclient.Block]{
			Result: &ethclient.Block{
				Transactions: []ethclient.Transaction{
					{
						Hash:     "hash1",
						From:     "from1",
						To:       "to1",
						GasPrice: ethclient.NewHexBigInt(big.NewInt(1)),
						Value:    ethclient.NewHexBigInt(big.NewInt(100)),
					},
				},
			},
		},
	}, nil)
	s.mockTransaction.EXPECT().Create(s.ctx, gomock.Any()).Return(nil)
	s.mockBlock.EXPECT().SetHeight(s.ctx, entity.TokenTypeETH, "123456").Return(nil)

	err := s.parser.ParseNextBlock(s.ctx)
	s.NoError(err)
}

func (s *EthereumTxParserTestSuite) TestParseNextBlockNoSubscriptionsSuccess() {
	s.mockBlock.EXPECT().GetHeight(s.ctx, entity.TokenTypeETH).Return("", context.DeadlineExceeded)

	err := s.parser.ParseNextBlock(s.ctx)
	s.ErrorIs(err, context.DeadlineExceeded)
}

func (s *EthereumTxParserTestSuite) TestParseNextBlockError() {
	s.mockBlock.EXPECT().GetHeight(s.ctx, entity.TokenTypeETH).Return("", entity.ErrBlockHightNotSet)
	s.mockEthClient.EXPECT().GetCurrentBlock(s.ctx, nil).Return(nil, errors.New("some error"))

	err := s.parser.ParseNextBlock(s.ctx)
	s.Error(err)
}

func TestEthereumTxParserTestSuite(t *testing.T) {
	suite.Run(t, new(EthereumTxParserTestSuite))
}
