// Code generated by MockGen. DO NOT EDIT.
// Source: basic_client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	http "net/http"
	url "net/url"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockBasicClient is a mock of BasicClient interface.
type MockBasicClient struct {
	ctrl     *gomock.Controller
	recorder *MockBasicClientMockRecorder
}

// MockBasicClientMockRecorder is the mock recorder for MockBasicClient.
type MockBasicClientMockRecorder struct {
	mock *MockBasicClient
}

// NewMockBasicClient creates a new mock instance.
func NewMockBasicClient(ctrl *gomock.Controller) *MockBasicClient {
	mock := &MockBasicClient{ctrl: ctrl}
	mock.recorder = &MockBasicClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBasicClient) EXPECT() *MockBasicClientMockRecorder {
	return m.recorder
}

// AddHeader mocks base method.
func (m *MockBasicClient) AddHeader(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddHeader", arg0, arg1)
}

// AddHeader indicates an expected call of AddHeader.
func (mr *MockBasicClientMockRecorder) AddHeader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHeader", reflect.TypeOf((*MockBasicClient)(nil).AddHeader), arg0, arg1)
}

// CollectionRequest mocks base method.
func (m *MockBasicClient) CollectionRequest(arg0 context.Context, arg1, arg2, arg3 string, arg4 interface{}, arg5 url.Values) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectionRequest", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CollectionRequest indicates an expected call of CollectionRequest.
func (mr *MockBasicClientMockRecorder) CollectionRequest(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectionRequest", reflect.TypeOf((*MockBasicClient)(nil).CollectionRequest), arg0, arg1, arg2, arg3, arg4, arg5)
}

// CollectionRequestMultipart mocks base method.
func (m *MockBasicClient) CollectionRequestMultipart(arg0 context.Context, arg1, arg2, arg3 string, arg4 interface{}, arg5 url.Values, arg6 map[string]string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectionRequestMultipart", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CollectionRequestMultipart indicates an expected call of CollectionRequestMultipart.
func (mr *MockBasicClientMockRecorder) CollectionRequestMultipart(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectionRequestMultipart", reflect.TypeOf((*MockBasicClient)(nil).CollectionRequestMultipart), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// ResetBaseURL mocks base method.
func (m *MockBasicClient) ResetBaseURL() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetBaseURL")
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetBaseURL indicates an expected call of ResetBaseURL.
func (mr *MockBasicClientMockRecorder) ResetBaseURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetBaseURL", reflect.TypeOf((*MockBasicClient)(nil).ResetBaseURL))
}

// SetTimeout mocks base method.
func (m *MockBasicClient) SetTimeout(arg0 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTimeout", arg0)
}

// SetTimeout indicates an expected call of SetTimeout.
func (mr *MockBasicClientMockRecorder) SetTimeout(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTimeout", reflect.TypeOf((*MockBasicClient)(nil).SetTimeout), arg0)
}

// URLRequest mocks base method.
func (m *MockBasicClient) URLRequest(arg0 context.Context, arg1, arg2 string, arg3 interface{}, arg4 url.Values, arg5 http.Header) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URLRequest", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// URLRequest indicates an expected call of URLRequest.
func (mr *MockBasicClientMockRecorder) URLRequest(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URLRequest", reflect.TypeOf((*MockBasicClient)(nil).URLRequest), arg0, arg1, arg2, arg3, arg4, arg5)
}
