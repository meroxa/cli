// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/meroxa/root/environments/create.go

// Package mock_environments is a generated GoMock package.
package mock_cmd

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meroxa "github.com/meroxa/meroxa-go"
)

// MockCreateEnvironmentClient is a mock of createEnvironmentClient interface.
type MockCreateEnvironmentClient struct {
	ctrl     *gomock.Controller
	recorder *MockCreateEnvironmentClientMockRecorder
}

// MockCreateEnvironmentClientMockRecorder is the mock recorder for MockCreateEnvironmentClient.
type MockCreateEnvironmentClientMockRecorder struct {
	mock *MockCreateEnvironmentClient
}

// NewMockCreateEnvironmentClient creates a new mock instance.
func NewMockCreateEnvironmentClient(ctrl *gomock.Controller) *MockCreateEnvironmentClient {
	mock := &MockCreateEnvironmentClient{ctrl: ctrl}
	mock.recorder = &MockCreateEnvironmentClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateEnvironmentClient) EXPECT() *MockCreateEnvironmentClientMockRecorder {
	return m.recorder
}

// CreateEnvironment mocks base method.
func (m *MockCreateEnvironmentClient) CreateEnvironment(ctx context.Context, body *meroxa.CreateEnvironmentInput) (*meroxa.Environment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEnvironment", ctx, body)
	ret0, _ := ret[0].(*meroxa.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEnvironment indicates an expected call of CreateEnvironment.
func (mr *MockCreateEnvironmentClientMockRecorder) CreateEnvironment(ctx, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEnvironment", reflect.TypeOf((*MockCreateEnvironmentClient)(nil).CreateEnvironment), ctx, body)
}
