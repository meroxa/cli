// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	log "github.com/meroxa/cli/log"
)

// MockCLI is a mock of CLI interface.
type MockCLI struct {
	ctrl     *gomock.Controller
	recorder *MockCLIMockRecorder
}

// MockCLIMockRecorder is the mock recorder for MockCLI.
type MockCLIMockRecorder struct {
	mock *MockCLI
}

// NewMockCLI creates a new mock instance.
func NewMockCLI(ctrl *gomock.Controller) *MockCLI {
	mock := &MockCLI{ctrl: ctrl}
	mock.recorder = &MockCLIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCLI) EXPECT() *MockCLIMockRecorder {
	return m.recorder
}

// CheckUncommittedChanges mocks base method.
func (m *MockCLI) CheckUncommittedChanges(arg0 context.Context, arg1 log.Logger, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUncommittedChanges", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckUncommittedChanges indicates an expected call of CheckUncommittedChanges.
func (mr *MockCLIMockRecorder) CheckUncommittedChanges(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUncommittedChanges", reflect.TypeOf((*MockCLI)(nil).CheckUncommittedChanges), arg0, arg1, arg2)
}

// CleanupDockerfile mocks base method.
func (m *MockCLI) CleanupDockerfile(arg0 log.Logger, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CleanupDockerfile", arg0, arg1)
}

// CleanupDockerfile indicates an expected call of CleanupDockerfile.
func (mr *MockCLIMockRecorder) CleanupDockerfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupDockerfile", reflect.TypeOf((*MockCLI)(nil).CleanupDockerfile), arg0, arg1)
}

// CreateDockerfile mocks base method.
func (m *MockCLI) CreateDockerfile(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDockerfile", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDockerfile indicates an expected call of CreateDockerfile.
func (mr *MockCLIMockRecorder) CreateDockerfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDockerfile", reflect.TypeOf((*MockCLI)(nil).CreateDockerfile), arg0, arg1)
}

// GetDeploymentSpec mocks base method.
func (m *MockCLI) GetDeploymentSpec(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeploymentSpec", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeploymentSpec indicates an expected call of GetDeploymentSpec.
func (mr *MockCLIMockRecorder) GetDeploymentSpec(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeploymentSpec", reflect.TypeOf((*MockCLI)(nil).GetDeploymentSpec), arg0, arg1)
}

// GetGitSha mocks base method.
func (m *MockCLI) GetGitSha(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitSha", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitSha indicates an expected call of GetGitSha.
func (mr *MockCLIMockRecorder) GetGitSha(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitSha", reflect.TypeOf((*MockCLI)(nil).GetGitSha), arg0, arg1)
}

// GetVersion mocks base method.
func (m *MockCLI) GetVersion(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersion indicates an expected call of GetVersion.
func (mr *MockCLIMockRecorder) GetVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockCLI)(nil).GetVersion), arg0)
}

// GitInit mocks base method.
func (m *MockCLI) GitInit(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GitInit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GitInit indicates an expected call of GitInit.
func (mr *MockCLIMockRecorder) GitInit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GitInit", reflect.TypeOf((*MockCLI)(nil).GitInit), arg0, arg1)
}

// Init mocks base method.
func (m *MockCLI) Init(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockCLIMockRecorder) Init(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockCLI)(nil).Init), arg0, arg1)
}

// Run mocks base method.
func (m *MockCLI) Run(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockCLIMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCLI)(nil).Run), arg0)
}

// StartGrpcServer mocks base method.
func (m *MockCLI) StartGrpcServer(arg0 context.Context, arg1 string) (func(), error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartGrpcServer", arg0, arg1)
	ret0, _ := ret[0].(func())
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartGrpcServer indicates an expected call of StartGrpcServer.
func (mr *MockCLIMockRecorder) StartGrpcServer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartGrpcServer", reflect.TypeOf((*MockCLI)(nil).StartGrpcServer), arg0, arg1)
}
