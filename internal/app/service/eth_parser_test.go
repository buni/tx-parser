package service_test

import (
	"context"
	"errors"
	"testing"

	contract_mock "github.com/buni/tx-parser/internal/app/contract/mock"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/app/repository"
	"github.com/buni/tx-parser/internal/app/service"
	"github.com/buni/tx-parser/internal/pkg/ethclient"
	ethclient_mock "github.com/buni/tx-parser/internal/pkg/ethclient/mock"
	"github.com/buni/tx-parser/internal/pkg/testing/testutils"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func BenchmarkParseBlock(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	ctx := context.Background()
	blockRepo := repository.NewInMemoryBlockRepository()
	transactionRepo := repository.NewInMemoryTransactionRepository()
	subscriptionRepo := repository.NewInMemorySubscriptionRepository()
	ethClient := ethclient.NewClient("https://ethereum-rpc.publicnode.com/")
	noopTxm := transactionmanager.NewNoopTxm()

	err := subscriptionRepo.Create(ctx, entity.TokenTypeETH, "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70")
	if err != nil {
		zap.L().Fatal("failed to create subscription", zap.Error(err))
	}

	err = blockRepo.SetHeight(ctx, entity.TokenTypeETH, "")
	if err != nil {
		zap.L().Fatal("failed to set initial height", zap.Error(err))
	}

	ethParserSvc := service.NewEthParser(subscriptionRepo, blockRepo, transactionRepo, ethClient, noopTxm)
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
	mockSubscription *contract_mock.MockSubscriptionRepositorty
	mockBlock        *contract_mock.MockBlockRepository
	mockTransaction  *contract_mock.MockTransactionRepository
	mockEthClient    *ethclient_mock.MockClient
	mockTxManager    *transactionmanager.NoopTxm
	parser           *service.EthereumTxParser
	ctx              context.Context
}

func (s *EthereumTxParserTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockSubscription = contract_mock.NewMockSubscriptionRepositorty(s.ctrl)
	s.mockBlock = contract_mock.NewMockBlockRepository(s.ctrl)
	s.mockTransaction = contract_mock.NewMockTransactionRepository(s.ctrl)
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
	s.mockSubscription.EXPECT().List(s.ctx, entity.TokenTypeETH).Return([]string{"0x2527d2ed1dd0e7de193cf121f1630caefc23ac70"}, nil)
	s.mockEthClient.EXPECT().GetBlockByNumber(s.ctx, &ethclient.GetBlockByNumberRequest{
		Number: "0x1e240",
	}).Return(&ethclient.GetBlockByNumberResponse{
		Response: ethclient.Response[ethclient.Block]{
			Result: &ethclient.Block{
				Transactions: []ethclient.Transaction{
					{
						Hash:  "hash1",
						From:  "from1",
						To:    "to1",
						Value: "100",
					},
				},
			},
		},
	}, nil)
	s.mockTransaction.EXPECT().BatchCreate(s.ctx, gomock.Any()).Return(nil)
	s.mockBlock.EXPECT().SetHeight(s.ctx, entity.TokenTypeETH, "123456").Return(nil)

	err := s.parser.ParseNextBlock(s.ctx)
	s.NoError(err)
}

func (s *EthereumTxParserTestSuite) TestParseNextBlockExistingBlockSuccess() {
	s.mockBlock.EXPECT().GetHeight(s.ctx, entity.TokenTypeETH).Return("123455", nil)
	s.mockEthClient.EXPECT().GetCurrentBlock(s.ctx, nil).Return(&ethclient.GetCurrentBlockResponse{Response: ethclient.Response[string]{
		Result: testutils.ToPtr("123456"),
	}}, nil)
	s.mockSubscription.EXPECT().List(s.ctx, entity.TokenTypeETH).Return([]string{"0x2527d2ed1dd0e7de193cf121f1630caefc23ac70"}, nil)
	s.mockEthClient.EXPECT().GetBlockByNumber(s.ctx, &ethclient.GetBlockByNumberRequest{
		Number: "0x1e240",
	}).Return(&ethclient.GetBlockByNumberResponse{
		Response: ethclient.Response[ethclient.Block]{
			Result: &ethclient.Block{
				Transactions: []ethclient.Transaction{
					{
						Hash:  "hash1",
						From:  "from1",
						To:    "to1",
						Value: "100",
					},
				},
			},
		},
	}, nil)
	s.mockTransaction.EXPECT().BatchCreate(s.ctx, gomock.Any()).Return(nil)
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
