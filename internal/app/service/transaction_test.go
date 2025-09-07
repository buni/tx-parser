package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	contract_mock "github.com/buni/tx-parser/internal/app/contract/mock"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/app/service"
	pubsub_mock "github.com/buni/tx-parser/internal/pkg/pubsub/mock"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockPublisher   *pubsub_mock.MockPublisher
	mockTransaction *contract_mock.MockTransactionRepository
	mockTxManager   *transactionmanager.NoopTxm
	transactionSvc  *service.TransactionService
	ctx             context.Context
}

func (s *TransactionServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockTransaction = contract_mock.NewMockTransactionRepository(s.ctrl)
	s.mockTxManager = transactionmanager.NewNoopTxm()
	s.mockPublisher = pubsub_mock.NewMockPublisher(s.ctrl)
	s.transactionSvc = service.NewTransactionService(s.mockTransaction, s.mockPublisher, s.mockTxManager)
	s.ctx = context.Background()
}

func (s *TransactionServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *TransactionServiceTestSuite) TestListAddressTransactionsSuccess() {
	req := &dto.ListAddressTransactionsRequest{Address: "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70", TokenType: entity.TokenTypeETH}
	want := []entity.Transaction{
		{
			ID:          "tx1",
			TokenType:   entity.TokenTypeETH,
			Hash:        "hash1",
			FromAddress: "from1",
			ToAddress:   "to1",
			Value:       "100",
			Address:     "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70",
		},
	}

	s.mockTransaction.EXPECT().List(s.ctx, entity.TokenTypeETH, req.Address).Return(want, nil)

	got, err := s.transactionSvc.ListAddressTransactions(s.ctx, req)
	s.NoError(err)
	s.Equal(want, got)
}

func (s *TransactionServiceTestSuite) TestListAddressTransactionsError() {
	req := &dto.ListAddressTransactionsRequest{
		TokenType: entity.TokenTypeETH,
		Address:   "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70",
	}

	s.mockTransaction.EXPECT().List(s.ctx, entity.TokenTypeETH, req.Address).Return(nil, context.DeadlineExceeded)

	got, err := s.transactionSvc.ListAddressTransactions(s.ctx, req)
	s.ErrorIs(err, context.DeadlineExceeded)
	s.Nil(got)
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
