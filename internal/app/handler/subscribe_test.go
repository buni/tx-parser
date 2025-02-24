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

type SubscribeHandlerTestSuite struct {
	suite.Suite
	svcMock *contract_mock.MockSubscriptionService
	handler *handler.SubscriptionHandler
	ctx     context.Context
	ctrl    *gomock.Controller
}

func (s *SubscribeHandlerTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.svcMock = contract_mock.NewMockSubscriptionService(s.ctrl)
	s.handler = handler.NewSubscriptionHandler(s.svcMock)
}

func (s *SubscribeHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *SubscribeHandlerTestSuite) statusCompare(gotCode, expectedCode int, gotBody string, expectedBody any) {
	s.Equal(expectedCode, gotCode)
	if expectedBody != nil && expectedBody != "" {
		jsonassert.New(s.T()).Assertf(gotBody, testutils.ToJSON(s.T(), expectedBody))
	}
}

func (s *SubscribeHandlerTestSuite) buildContext() context.Context {
	chiContext := chi.NewRouteContext()
	return context.WithValue(s.ctx, chi.RouteCtxKey, chiContext)
}

func (s *SubscribeHandlerTestSuite) TestSubscribeSuccess() {
	req := &dto.SubscribeRequest{
		Address: "123",
	}
	expectedBody := dto.SubscribeResponse{}

	s.ctx = s.buildContext()

	s.svcMock.EXPECT().Subscribe(s.ctx, req).Return(nil)

	recorder := httptest.NewRecorder()

	handlerwrap.WrapDefault(s.handler.Subscribe).ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/", testutils.ToJSONReader(s.T(), req)).WithContext(s.ctx))
	s.statusCompare(recorder.Code, http.StatusCreated, recorder.Body.String(), expectedBody)
}

func (s *SubscribeHandlerTestSuite) TestSubscribeFailure() {
	req := &dto.SubscribeRequest{
		Address: "123",
	}

	s.ctx = s.buildContext()

	s.svcMock.EXPECT().Subscribe(s.ctx, req).Return(context.DeadlineExceeded)

	recorder := httptest.NewRecorder()

	handlerwrap.WrapDefault(s.handler.Subscribe).ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/", testutils.ToJSONReader(s.T(), req)).WithContext(s.ctx))
	s.statusCompare(recorder.Code, http.StatusInternalServerError, recorder.Body.String(), render.ErrorResponse{
		Error: &render.Error{
			Status:  render.InternalServerError,
			Message: "internal server error",
		},
	})
}

func TestSubscribeHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SubscribeHandlerTestSuite))
}
