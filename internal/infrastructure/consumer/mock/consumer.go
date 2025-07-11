// Code generated by MockGen. DO NOT EDIT.
// Source: order-service/internal/infrastructure/consumer (interfaces: Consumer)
//
// Generated by this command:
//
//	mockgen -destination=./mock/consumer.go -package=mock order-service/internal/infrastructure/consumer Consumer
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockConsumer is a mock of Consumer interface.
type MockConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerMockRecorder
	isgomock struct{}
}

// MockConsumerMockRecorder is the mock recorder for MockConsumer.
type MockConsumerMockRecorder struct {
	mock *MockConsumer
}

// NewMockConsumer creates a new mock instance.
func NewMockConsumer(ctrl *gomock.Controller) *MockConsumer {
	mock := &MockConsumer{ctrl: ctrl}
	mock.recorder = &MockConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumer) EXPECT() *MockConsumerMockRecorder {
	return m.recorder
}

// Consume mocks base method.
func (m *MockConsumer) Consume(ctx context.Context) ([]byte, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", ctx)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Consume indicates an expected call of Consume.
func (mr *MockConsumerMockRecorder) Consume(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockConsumer)(nil).Consume), ctx)
}
