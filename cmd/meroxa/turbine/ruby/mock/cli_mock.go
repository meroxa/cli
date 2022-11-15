// Code generated by MockGen. DO NOT EDIT.
// Source: cli.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockturbineServer is a mock of turbineServer interface.
type MockturbineServer struct {
	ctrl     *gomock.Controller
	recorder *MockturbineServerMockRecorder
}

// MockturbineServerMockRecorder is the mock recorder for MockturbineServer.
type MockturbineServerMockRecorder struct {
	mock *MockturbineServer
}

// NewMockturbineServer creates a new mock instance.
func NewMockturbineServer(ctrl *gomock.Controller) *MockturbineServer {
	mock := &MockturbineServer{ctrl: ctrl}
	mock.recorder = &MockturbineServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockturbineServer) EXPECT() *MockturbineServerMockRecorder {
	return m.recorder
}

// GracefulStop mocks base method.
func (m *MockturbineServer) GracefulStop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GracefulStop")
}

// GracefulStop indicates an expected call of GracefulStop.
func (mr *MockturbineServerMockRecorder) GracefulStop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GracefulStop", reflect.TypeOf((*MockturbineServer)(nil).GracefulStop))
}

// Run mocks base method.
func (m *MockturbineServer) Run(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run.
func (mr *MockturbineServerMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockturbineServer)(nil).Run), arg0)
}
