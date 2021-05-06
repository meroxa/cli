// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/meroxa/root/resources/list.go

// Package mock_cmd is a generated GoMock package.
package mock_cmd

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meroxa "github.com/meroxa/meroxa-go"
)

// MockListResourcesClient is a mock of listResourcesClient interface.
type MockListResourcesClient struct {
	ctrl     *gomock.Controller
	recorder *MockListResourcesClientMockRecorder
}

// MockListResourcesClientMockRecorder is the mock recorder for MockListResourcesClient.
type MockListResourcesClientMockRecorder struct {
	mock *MockListResourcesClient
}

// NewMockListResourcesClient creates a new mock instance.
func NewMockListResourcesClient(ctrl *gomock.Controller) *MockListResourcesClient {
	mock := &MockListResourcesClient{ctrl: ctrl}
	mock.recorder = &MockListResourcesClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockListResourcesClient) EXPECT() *MockListResourcesClientMockRecorder {
	return m.recorder
}

// ListResourceTypes mocks base method.
func (m *MockListResourcesClient) ListResourceTypes(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResourceTypes", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResourceTypes indicates an expected call of ListResourceTypes.
func (mr *MockListResourcesClientMockRecorder) ListResourceTypes(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResourceTypes", reflect.TypeOf((*MockListResourcesClient)(nil).ListResourceTypes), ctx)
}

// ListResources mocks base method.
func (m *MockListResourcesClient) ListResources(ctx context.Context) ([]*meroxa.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResources", ctx)
	ret0, _ := ret[0].([]*meroxa.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResources indicates an expected call of ListResources.
func (mr *MockListResourcesClientMockRecorder) ListResources(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResources", reflect.TypeOf((*MockListResourcesClient)(nil).ListResources), ctx)
}