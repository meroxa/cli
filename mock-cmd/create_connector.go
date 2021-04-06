// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/create_connector.go

// Package mock_cmd is a generated GoMock package.
package mock_cmd

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meroxa "github.com/meroxa/meroxa-go"
)

// MockCreateConnectorClient is a mock of CreateConnectorClient interface.
type MockCreateConnectorClient struct {
	ctrl     *gomock.Controller
	recorder *MockCreateConnectorClientMockRecorder
}

// MockCreateConnectorClientMockRecorder is the mock recorder for MockCreateConnectorClient.
type MockCreateConnectorClientMockRecorder struct {
	mock *MockCreateConnectorClient
}

// NewMockCreateConnectorClient creates a new mock instance.
func NewMockCreateConnectorClient(ctrl *gomock.Controller) *MockCreateConnectorClient {
	mock := &MockCreateConnectorClient{ctrl: ctrl}
	mock.recorder = &MockCreateConnectorClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateConnectorClient) EXPECT() *MockCreateConnectorClientMockRecorder {
	return m.recorder
}

// CreateConnector mocks base method.
func (m *MockCreateConnectorClient) CreateConnector(ctx context.Context, input meroxa.CreateConnectorInput) (*meroxa.Connector, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateConnector", ctx, input)
	ret0, _ := ret[0].(*meroxa.Connector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateConnector indicates an expected call of CreateConnector.
func (mr *MockCreateConnectorClientMockRecorder) CreateConnector(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateConnector", reflect.TypeOf((*MockCreateConnectorClient)(nil).CreateConnector), ctx, input)
}

// GetResourceByName mocks base method.
func (m *MockCreateConnectorClient) GetResourceByName(ctx context.Context, name string) (*meroxa.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceByName", ctx, name)
	ret0, _ := ret[0].(*meroxa.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourceByName indicates an expected call of GetResourceByName.
func (mr *MockCreateConnectorClientMockRecorder) GetResourceByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceByName", reflect.TypeOf((*MockCreateConnectorClient)(nil).GetResourceByName), ctx, name)
}
