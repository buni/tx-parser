package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	contract_mock "github.com/buni/tx-parser/internal/app/contract/mock"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/handler"
	handlerwrap "github.com/buni/tx-parser/internal/pkg/handler"
	"github.com/buni/tx-parser/internal/pkg/render"
	"github.com/buni/tx-parser/internal/pkg/testing/testutils"
	"github.com/go-chi/chi/v5"
	"github.com/kinbiko/jsonassert"
	"github.com/test-go/testify/suite"
	"go.uber.org/mock/gomock"
)

type BlockHandlerTestSuite struct {
	suite.Suite
	svcMock *contract_mock.MockBlockService
	handler *handler.BlockHandler
	ctx     context.Context
	ctrl    *gomock.Controller
}

func (s *BlockHandlerTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.svcMock = contract_mock.NewMockBlockService(s.ctrl)
	s.handler = handler.NewBlockHandler(s.svcMock)
}

func (s *BlockHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *BlockHandlerTestSuite) statusCompare(gotCode, expectedCode int, gotBody string, expectedBody any) {
	s.Equal(expectedCode, gotCode)
	if expectedBody != nil && expectedBody != "" {
		jsonassert.New(s.T()).Assertf(gotBody, testutils.ToJSON(s.T(), expectedBody))
	}
}

func (s *BlockHandlerTestSuite) buildContext() context.Context {
	chiContext := chi.NewRouteContext()
	return context.WithValue(s.ctx, chi.RouteCtxKey, chiContext)
}

func (s *BlockHandlerTestSuite) TestGetCurrentBlockSuccess() {
	req := &dto.GetCurrentBlockRequest{}
	expectedBody := dto.GetCurrentBlockResponse{
		Height: "123",
	}

	s.ctx = s.buildContext()

	s.svcMock.EXPECT().GetCurrentBlock(s.ctx).Return("123", nil)

	recorder := httptest.NewRecorder()

	handlerwrap.WrapDefaultBasic(s.handler.GetCurrentBlock).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", testutils.ToJSONReader(s.T(), req)).WithContext(s.ctx))
	s.statusCompare(recorder.Code, http.StatusOK, recorder.Body.String(), expectedBody)
}

func (s *BlockHandlerTestSuite) TestGetCurrentBlockFailure() {
	req := &dto.GetCurrentBlockRequest{}

	s.ctx = s.buildContext()

	s.svcMock.EXPECT().GetCurrentBlock(s.ctx).Return("", context.DeadlineExceeded)

	recorder := httptest.NewRecorder()

	handlerwrap.WrapDefaultBasic(s.handler.GetCurrentBlock).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", testutils.ToJSONReader(s.T(), req)).WithContext(s.ctx))
	s.statusCompare(recorder.Code, http.StatusInternalServerError, recorder.Body.String(), render.ErrorResponse{
		Error: &render.Error{
			Status:  render.InternalServerError,
			Message: "internal server error",
		},
	})
}

func TestBlockHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BlockHandlerTestSuite))
}
