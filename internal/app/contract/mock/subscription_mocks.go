// Code generated by MockGen. DO NOT EDIT.
// Source: subscription.go
//
// Generated by this command:
//
//	mockgen -source=subscription.go -destination=mock/subscription_mocks.go -package contract_mock
//

// Package contract_mock is a generated GoMock package.
package contract_mock

import (
	context "context"
	reflect "reflect"

	dto "github.com/buni/tx-parser/internal/app/dto"
	entity "github.com/buni/tx-parser/internal/app/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockSubscriptionRepositorty is a mock of SubscriptionRepositorty interface.
type MockSubscriptionRepositorty struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionRepositortyMockRecorder
	isgomock struct{}
}

// MockSubscriptionRepositortyMockRecorder is the mock recorder for MockSubscriptionRepositorty.
type MockSubscriptionRepositortyMockRecorder struct {
	mock *MockSubscriptionRepositorty
}

// NewMockSubscriptionRepositorty creates a new mock instance.
func NewMockSubscriptionRepositorty(ctrl *gomock.Controller) *MockSubscriptionRepositorty {
	mock := &MockSubscriptionRepositorty{ctrl: ctrl}
	mock.recorder = &MockSubscriptionRepositortyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionRepositorty) EXPECT() *MockSubscriptionRepositortyMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSubscriptionRepositorty) Create(ctx context.Context, tokenType entity.TokenType, address string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, tokenType, address)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSubscriptionRepositortyMockRecorder) Create(ctx, tokenType, address any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSubscriptionRepositorty)(nil).Create), ctx, tokenType, address)
}

// List mocks base method.
func (m *MockSubscriptionRepositorty) List(ctx context.Context, tokenType entity.TokenType) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, tokenType)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockSubscriptionRepositortyMockRecorder) List(ctx, tokenType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSubscriptionRepositorty)(nil).List), ctx, tokenType)
}

// MockSubscriptionService is a mock of SubscriptionService interface.
type MockSubscriptionService struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionServiceMockRecorder
	isgomock struct{}
}

// MockSubscriptionServiceMockRecorder is the mock recorder for MockSubscriptionService.
type MockSubscriptionServiceMockRecorder struct {
	mock *MockSubscriptionService
}

// NewMockSubscriptionService creates a new mock instance.
func NewMockSubscriptionService(ctrl *gomock.Controller) *MockSubscriptionService {
	mock := &MockSubscriptionService{ctrl: ctrl}
	mock.recorder = &MockSubscriptionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionService) EXPECT() *MockSubscriptionServiceMockRecorder {
	return m.recorder
}

// Subscribe mocks base method.
func (m *MockSubscriptionService) Subscribe(ctx context.Context, req *dto.SubscribeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockSubscriptionServiceMockRecorder) Subscribe(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockSubscriptionService)(nil).Subscribe), ctx, req)
}
