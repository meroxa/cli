// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/add_resource.go

// Package mock_cmd is a generated GoMock package.
package mock_cmd

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meroxa "github.com/meroxa/meroxa-go"
)

// MockAddResourceClient is a mock of AddResourceClient interface.
type MockAddResourceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAddResourceClientMockRecorder
}

// MockAddResourceClientMockRecorder is the mock recorder for MockAddResourceClient.
type MockAddResourceClientMockRecorder struct {
	mock *MockAddResourceClient
}

// NewMockAddResourceClient creates a new mock instance.
func NewMockAddResourceClient(ctrl *gomock.Controller) *MockAddResourceClient {
	mock := &MockAddResourceClient{ctrl: ctrl}
	mock.recorder = &MockAddResourceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAddResourceClient) EXPECT() *MockAddResourceClientMockRecorder {
	return m.recorder
}

// CreateResource mocks base method.
func (m *MockAddResourceClient) CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateResource", ctx, resource)
	ret0, _ := ret[0].(*meroxa.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateResource indicates an expected call of CreateResource.
func (mr *MockAddResourceClientMockRecorder) CreateResource(ctx, resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateResource", reflect.TypeOf((*MockAddResourceClient)(nil).CreateResource), ctx, resource)
}
