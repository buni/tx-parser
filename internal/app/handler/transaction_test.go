package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	contract_mock "github.com/buni/tx-parser/internal/app/contract/mock"
	"github.com/buni/tx-parser/internal/app/dto"
	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/buni/tx-parser/internal/app/handler"
	handlerwrap "github.com/buni/tx-parser/internal/pkg/handler"
	"github.com/buni/tx-parser/internal/pkg/render"
	"github.com/buni/tx-parser/internal/pkg/testing/testutils"
	"github.com/go-chi/chi/v5"
	"github.com/kinbiko/jsonassert"
	"github.com/test-go/testify/suite"
	"go.uber.org/mock/gomock"
)

type TransactionHandlerTestSuite struct {
	suite.Suite
	svcMock *contract_mock.MockTransactionService
	handler *handler.TransactionHandler
	ctx     context.Context
	ctrl    *gomock.Controller
}

func (s *TransactionHandlerTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.svcMock = contract_mock.NewMockTransactionService(s.ctrl)
	s.handler = handler.NewTransactionHandler(s.svcMock)
}

func (s *TransactionHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *TransactionHandlerTestSuite) statusCompare(gotCode, expectedCode int, gotBody string, expectedBody any) {
	s.Equal(expectedCode, gotCode)
	if expectedBody != nil && expectedBody != "" {
		jsonassert.New(s.T()).Assertf(gotBody, testutils.ToJSON(s.T(), expectedBody))
	}
}

func (s *TransactionHandlerTestSuite) buildContext(address string) context.Context {
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("address", address)

	return context.WithValue(s.ctx, chi.RouteCtxKey, chiContext)
}

func (s *TransactionHandlerTestSuite) TestListAddressTransactionsSuccess() {
	req := &dto.ListAddressTransactionsRequest{
		Address: "123",
	}
	expectedBody := dto.ListAddressTransactionsResponse{
		Transactions: []dto.Transaction{
			{
				ID:        "1",
				TokenType: entity.TokenTypeETH.String(),
				To:        "1",
				From:      "1",
				Address:   "1",
				Hash:      "1",
				Value:     "1",
			},
			{
				ID:        "1",
				TokenType: entity.TokenTypeETH.String(),
				To:        "1",
				From:      "1",
				Address:   "1",
				Hash:      "1",
				Value:     "1",
			},
		},
	}

	s.ctx = s.buildContext("123")

	s.svcMock.EXPECT().ListAddressTransactions(s.ctx, req).Return([]entity.Transaction{
		{
			ID:        "1",
			TokenType: entity.TokenTypeETH,
			To:        "1",
			From:      "1",
			Address:   "1",
			Hash:      "1",
			Value:     "1",
		},
		{
			ID:        "1",
			TokenType: entity.TokenTypeETH,
			To:        "1",
			From:      "1",
			Address:   "1",
			Hash:      "1",
			Value:     "1",
		},
	}, nil)

	recorder := httptest.NewRecorder()

	handlerwrap.WrapDefaultBasic(s.handler.ListAddressTransactions).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", testutils.ToJSONReader(s.T(), req)).WithContext(s.ctx))
	s.statusCompare(recorder.Code, http.StatusOK, recorder.Body.String(), expectedBody)
}

func (s *TransactionHandlerTestSuite) TestListAddressTransactionsFailure() {
	req := &dto.ListAddressTransactionsRequest{
		Address: "123",
	}

	s.ctx = s.buildContext("123")

	s.svcMock.EXPECT().ListAddressTransactions(s.ctx, req).Return(nil, context.DeadlineExceeded)

	recorder := httptest.NewRecorder()

	handlerwrap.WrapDefaultBasic(s.handler.ListAddressTransactions).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", testutils.ToJSONReader(s.T(), req)).WithContext(s.ctx))
	s.statusCompare(recorder.Code, http.StatusInternalServerError, recorder.Body.String(), render.ErrorResponse{
		Error: &render.Error{
			Status:  render.InternalServerError,
			Message: "internal server error",
		},
	})
}

func TestTransactionHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionHandlerTestSuite))
}
