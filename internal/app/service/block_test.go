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
	"github.com/buni/tx-parser/internal/pkg/transactionmanager"
)

type BlockServiceTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	mockBlock     *contract_mock.MockBlockRepository
	mockTxManager *transactionmanager.NoopTxm
	blockSvc      *service.BlockService
	ctx           context.Context
}

func (s *BlockServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockBlock = contract_mock.NewMockBlockRepository(s.ctrl)
	s.mockTxManager = transactionmanager.NewNoopTxm()
	s.blockSvc = service.NewBlockService(s.mockBlock, s.mockTxManager)
	s.ctx = context.Background()
}

func (s *BlockServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *BlockServiceTestSuite) TestGetCurrentBlockSuccess() {
	expectedHeight := "123456"
	req := &dto.GetCurrentBlockRequest{
		TokenType: entity.TokenTypeETH,
	}

	s.mockBlock.EXPECT().GetHeight(s.ctx, req.TokenType).Return(expectedHeight, nil)

	height, err := s.blockSvc.GetCurrentBlock(s.ctx, req)
	s.NoError(err)
	s.Equal(expectedHeight, height)
}

func (s *BlockServiceTestSuite) TestGetCurrentBlockError() {
	req := &dto.GetCurrentBlockRequest{
		TokenType: entity.TokenTypeETH,
	}

	s.mockBlock.EXPECT().GetHeight(s.ctx, req.TokenType).Return("", context.DeadlineExceeded)

	height, err := s.blockSvc.GetCurrentBlock(s.ctx, req)
	s.ErrorIs(err, context.DeadlineExceeded)
	s.Empty(height)
}

func TestBlockServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BlockServiceTestSuite))
}
