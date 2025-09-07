package service_test

import (
	"context"
	"testing"

	contract_mock "github.com/buni/tx-parser/internal/app/contract/mock"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/app/service"
	"github.com/buni/tx-parser/internal/pkg/testing/testutils"
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type SubscriberServiceTestSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	mockSubscription *contract_mock.MockSubscriptionRepository
	mockTxManager    *transactionmanager.NoopTxm
	subscriber       *service.SubscriberService
	ctx              context.Context
}

func (s *SubscriberServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockSubscription = contract_mock.NewMockSubscriptionRepository(s.ctrl)
	s.mockTxManager = transactionmanager.NewNoopTxm()
	s.subscriber = service.NewSubscriberService(s.mockSubscription, s.mockTxManager)
	s.ctx = context.Background()
}

func (s *SubscriberServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *SubscriberServiceTestSuite) TestSubscribeSuccess() {
	req := &dto.SubscribeRequest{
		Address:   "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70",
		UserID:    "user1",
		TokenType: entity.TokenTypeETH,
	}

	sub := entity.Subscription{
		TokenType: req.TokenType,
		Address:   req.Address,
		UserID:    req.UserID,
	}

	s.mockSubscription.EXPECT().Create(s.ctx, testutils.NewMatcher(sub, cmpopts.IgnoreFields(entity.Subscription{}, "ID", "CreatedAt"))).Return(nil)

	err := s.subscriber.Subscribe(s.ctx, req)
	s.NoError(err)
}

func (s *SubscriberServiceTestSuite) TestSubscribeError() {
	req := &dto.SubscribeRequest{
		Address:   "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70",
		UserID:    "user1",
		TokenType: entity.TokenTypeETH,
	}

	sub := entity.Subscription{
		TokenType: req.TokenType,
		Address:   req.Address,
		UserID:    req.UserID,
	}

	s.mockSubscription.EXPECT().Create(s.ctx, testutils.NewMatcher(sub, cmpopts.IgnoreFields(entity.Subscription{}, "ID", "CreatedAt"))).Return(context.DeadlineExceeded)

	err := s.subscriber.Subscribe(s.ctx, req)
	s.ErrorIs(err, context.DeadlineExceeded)
}

func TestSubscriberServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriberServiceTestSuite))
}
