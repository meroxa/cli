// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/meroxa/root/environments/list.go

// Package mock_environments is a generated GoMock package.
package mock_cmd

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meroxa "github.com/meroxa/meroxa-go"
)

// MocklistEnvironmentsClient is a mock of listEnvironmentsClient interface.
type MocklistEnvironmentsClient struct {
	ctrl     *gomock.Controller
	recorder *MocklistEnvironmentsClientMockRecorder
}

// MocklistEnvironmentsClientMockRecorder is the mock recorder for MocklistEnvironmentsClient.
type MocklistEnvironmentsClientMockRecorder struct {
	mock *MocklistEnvironmentsClient
}

// NewMocklistEnvironmentsClient creates a new mock instance.
func NewMocklistEnvironmentsClient(ctrl *gomock.Controller) *MocklistEnvironmentsClient {
	mock := &MocklistEnvironmentsClient{ctrl: ctrl}
	mock.recorder = &MocklistEnvironmentsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocklistEnvironmentsClient) EXPECT() *MocklistEnvironmentsClientMockRecorder {
	return m.recorder
}

// ListEnvironments mocks base method.
func (m *MocklistEnvironmentsClient) ListEnvironments(ctx context.Context) ([]*meroxa.Environment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEnvironments", ctx)
	ret0, _ := ret[0].([]*meroxa.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEnvironments indicates an expected call of ListEnvironments.
func (mr *MocklistEnvironmentsClientMockRecorder) ListEnvironments(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEnvironments", reflect.TypeOf((*MocklistEnvironmentsClient)(nil).ListEnvironments), ctx)
}