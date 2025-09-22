package task

import (
	"context"
	"net/http"
	"testing"

	contract_mock "github.com/buni/tx-parser/internal/app/contract/mock"
	"github.com/buni/tx-parser/internal/pkg/ethclient"
	"github.com/buni/tx-parser/internal/pkg/scheduler"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type EthBlockParserTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	svcMock     *contract_mock.MockParserService
	blockParser scheduler.Task
	ctx         context.Context
}

func (s *EthBlockParserTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.svcMock = contract_mock.NewMockParserService(s.ctrl)
	s.blockParser = NewEthBlockParser(s.svcMock)
	s.ctx = context.Background()
}

func (s *EthBlockParserTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *EthBlockParserTestSuite) TestHandleSuccess() {
	s.svcMock.EXPECT().ParseNextBlock(s.ctx).Return(nil)
	clientError := &ethclient.ErrorResponse{StatusCode: http.StatusNotFound}
	s.svcMock.EXPECT().ParseNextBlock(s.ctx).Return(clientError)

	err := s.blockParser.Handle(s.ctx)
	s.NoError(err)
}

func (s *EthBlockParserTestSuite) TestHandleClientError() {
	clientError := &ethclient.ErrorResponse{StatusCode: http.StatusNotFound}
	s.svcMock.EXPECT().ParseNextBlock(s.ctx).Return(clientError)

	err := s.blockParser.Handle(s.ctx)
	s.NoError(err)
}

func (s *EthBlockParserTestSuite) TestHandleOtherError() {
	s.svcMock.EXPECT().ParseNextBlock(s.ctx).Return(context.DeadlineExceeded)

	err := s.blockParser.Handle(s.ctx)
	s.ErrorIs(err, context.DeadlineExceeded)
}

func (s *EthBlockParserTestSuite) TestHandleCancelContext() {
	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	err := s.blockParser.Handle(ctx)
	s.NoError(err)
}

func TestEthBlockParserTestSuite(t *testing.T) {
	suite.Run(t, new(EthBlockParserTestSuite))
}
