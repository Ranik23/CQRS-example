// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/antonfedorov/Projects/clothes/catalog-service/internal/infrastructure/txmanager/tx_manager.go
//
// Generated by this command:
//
//	mockgen --source=/Users/antonfedorov/Projects/clothes/catalog-service/internal/infrastructure/txmanager/tx_manager.go --destination=/Users/antonfedorov/Projects/clothes/catalog-service/internal/infrastructure/txmanager/mock/tx_manager.go --package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	pgx "github.com/jackc/pgx/v5"
	gomock "go.uber.org/mock/gomock"
)

// MockTxManager is a mock of TxManager interface.
type MockTxManager struct {
	ctrl     *gomock.Controller
	recorder *MockTxManagerMockRecorder
	isgomock struct{}
}

// MockTxManagerMockRecorder is the mock recorder for MockTxManager.
type MockTxManagerMockRecorder struct {
	mock *MockTxManager
}

// NewMockTxManager creates a new mock instance.
func NewMockTxManager(ctrl *gomock.Controller) *MockTxManager {
	mock := &MockTxManager{ctrl: ctrl}
	mock.recorder = &MockTxManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxManager) EXPECT() *MockTxManagerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockTxManager) Run(ctx context.Context, fn func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockTxManagerMockRecorder) Run(ctx, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockTxManager)(nil).Run), ctx, fn)
}

// Tx mocks base method.
func (m *MockTxManager) Tx(ctx context.Context) pgx.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx", ctx)
	ret0, _ := ret[0].(pgx.Tx)
	return ret0
}

// Tx indicates an expected call of Tx.
func (mr *MockTxManagerMockRecorder) Tx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockTxManager)(nil).Tx), ctx)
}
