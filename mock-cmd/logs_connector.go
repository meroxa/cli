// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/meroxa/root/connectors/logs_connector.go

// Package mock_connectors is a generated GoMock package.
package mock_cmd

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLogsConnectorClient is a mock of logsConnectorClient interface.
type MockLogsConnectorClient struct {
	ctrl     *gomock.Controller
	recorder *MockLogsConnectorClientMockRecorder
}

// MockLogsConnectorClientMockRecorder is the mock recorder for MockLogsConnectorClient.
type MockLogsConnectorClientMockRecorder struct {
	mock *MockLogsConnectorClient
}

// NewMockLogsConnectorClient creates a new mock instance.
func NewMockLogsConnectorClient(ctrl *gomock.Controller) *MockLogsConnectorClient {
	mock := &MockLogsConnectorClient{ctrl: ctrl}
	mock.recorder = &MockLogsConnectorClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogsConnectorClient) EXPECT() *MockLogsConnectorClientMockRecorder {
	return m.recorder
}

// GetConnectorLogs mocks base method.
func (m *MockLogsConnectorClient) GetConnectorLogs(ctx context.Context, connectorName string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectorLogs", ctx, connectorName)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectorLogs indicates an expected call of GetConnectorLogs.
func (mr *MockLogsConnectorClientMockRecorder) GetConnectorLogs(ctx, connectorName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectorLogs", reflect.TypeOf((*MockLogsConnectorClient)(nil).GetConnectorLogs), ctx, connectorName)
}
