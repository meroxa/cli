// Code generated by MockGen. DO NOT EDIT.
// Source: vendor/github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core/turbine_grpc.pb.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	core "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// MockTurbineServiceClient is a mock of TurbineServiceClient interface.
type MockTurbineServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockTurbineServiceClientMockRecorder
}

// MockTurbineServiceClientMockRecorder is the mock recorder for MockTurbineServiceClient.
type MockTurbineServiceClientMockRecorder struct {
	mock *MockTurbineServiceClient
}

// NewMockTurbineServiceClient creates a new mock instance.
func NewMockTurbineServiceClient(ctrl *gomock.Controller) *MockTurbineServiceClient {
	mock := &MockTurbineServiceClient{ctrl: ctrl}
	mock.recorder = &MockTurbineServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTurbineServiceClient) EXPECT() *MockTurbineServiceClientMockRecorder {
	return m.recorder
}

// AddProcessToCollection mocks base method.
func (m *MockTurbineServiceClient) AddProcessToCollection(ctx context.Context, in *core.ProcessCollectionRequest, opts ...grpc.CallOption) (*core.Collection, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddProcessToCollection", varargs...)
	ret0, _ := ret[0].(*core.Collection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProcessToCollection indicates an expected call of AddProcessToCollection.
func (mr *MockTurbineServiceClientMockRecorder) AddProcessToCollection(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProcessToCollection", reflect.TypeOf((*MockTurbineServiceClient)(nil).AddProcessToCollection), varargs...)
}

// GetResource mocks base method.
func (m *MockTurbineServiceClient) GetResource(ctx context.Context, in *core.GetResourceRequest, opts ...grpc.CallOption) (*core.Resource, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetResource", varargs...)
	ret0, _ := ret[0].(*core.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResource indicates an expected call of GetResource.
func (mr *MockTurbineServiceClientMockRecorder) GetResource(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockTurbineServiceClient)(nil).GetResource), varargs...)
}

// GetSpec mocks base method.
func (m *MockTurbineServiceClient) GetSpec(ctx context.Context, in *core.GetSpecRequest, opts ...grpc.CallOption) (*core.GetSpecResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSpec", varargs...)
	ret0, _ := ret[0].(*core.GetSpecResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpec indicates an expected call of GetSpec.
func (mr *MockTurbineServiceClientMockRecorder) GetSpec(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpec", reflect.TypeOf((*MockTurbineServiceClient)(nil).GetSpec), varargs...)
}

// HasFunctions mocks base method.
func (m *MockTurbineServiceClient) HasFunctions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HasFunctions", varargs...)
	ret0, _ := ret[0].(*wrapperspb.BoolValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasFunctions indicates an expected call of HasFunctions.
func (mr *MockTurbineServiceClientMockRecorder) HasFunctions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasFunctions", reflect.TypeOf((*MockTurbineServiceClient)(nil).HasFunctions), varargs...)
}

// Init mocks base method.
func (m *MockTurbineServiceClient) Init(ctx context.Context, in *core.InitRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Init", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Init indicates an expected call of Init.
func (mr *MockTurbineServiceClientMockRecorder) Init(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockTurbineServiceClient)(nil).Init), varargs...)
}

// ListResources mocks base method.
func (m *MockTurbineServiceClient) ListResources(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*core.ListResourcesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListResources", varargs...)
	ret0, _ := ret[0].(*core.ListResourcesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResources indicates an expected call of ListResources.
func (mr *MockTurbineServiceClientMockRecorder) ListResources(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResources", reflect.TypeOf((*MockTurbineServiceClient)(nil).ListResources), varargs...)
}

// ReadCollection mocks base method.
func (m *MockTurbineServiceClient) ReadCollection(ctx context.Context, in *core.ReadCollectionRequest, opts ...grpc.CallOption) (*core.Collection, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadCollection", varargs...)
	ret0, _ := ret[0].(*core.Collection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCollection indicates an expected call of ReadCollection.
func (mr *MockTurbineServiceClientMockRecorder) ReadCollection(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCollection", reflect.TypeOf((*MockTurbineServiceClient)(nil).ReadCollection), varargs...)
}

// RegisterSecret mocks base method.
func (m *MockTurbineServiceClient) RegisterSecret(ctx context.Context, in *core.Secret, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RegisterSecret", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterSecret indicates an expected call of RegisterSecret.
func (mr *MockTurbineServiceClientMockRecorder) RegisterSecret(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterSecret", reflect.TypeOf((*MockTurbineServiceClient)(nil).RegisterSecret), varargs...)
}

// WriteCollectionToResource mocks base method.
func (m *MockTurbineServiceClient) WriteCollectionToResource(ctx context.Context, in *core.WriteCollectionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteCollectionToResource", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteCollectionToResource indicates an expected call of WriteCollectionToResource.
func (mr *MockTurbineServiceClientMockRecorder) WriteCollectionToResource(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteCollectionToResource", reflect.TypeOf((*MockTurbineServiceClient)(nil).WriteCollectionToResource), varargs...)
}

// MockTurbineServiceServer is a mock of TurbineServiceServer interface.
type MockTurbineServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockTurbineServiceServerMockRecorder
}

// MockTurbineServiceServerMockRecorder is the mock recorder for MockTurbineServiceServer.
type MockTurbineServiceServerMockRecorder struct {
	mock *MockTurbineServiceServer
}

// NewMockTurbineServiceServer creates a new mock instance.
func NewMockTurbineServiceServer(ctrl *gomock.Controller) *MockTurbineServiceServer {
	mock := &MockTurbineServiceServer{ctrl: ctrl}
	mock.recorder = &MockTurbineServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTurbineServiceServer) EXPECT() *MockTurbineServiceServerMockRecorder {
	return m.recorder
}

// AddProcessToCollection mocks base method.
func (m *MockTurbineServiceServer) AddProcessToCollection(arg0 context.Context, arg1 *core.ProcessCollectionRequest) (*core.Collection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProcessToCollection", arg0, arg1)
	ret0, _ := ret[0].(*core.Collection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProcessToCollection indicates an expected call of AddProcessToCollection.
func (mr *MockTurbineServiceServerMockRecorder) AddProcessToCollection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProcessToCollection", reflect.TypeOf((*MockTurbineServiceServer)(nil).AddProcessToCollection), arg0, arg1)
}

// GetResource mocks base method.
func (m *MockTurbineServiceServer) GetResource(arg0 context.Context, arg1 *core.GetResourceRequest) (*core.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResource", arg0, arg1)
	ret0, _ := ret[0].(*core.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResource indicates an expected call of GetResource.
func (mr *MockTurbineServiceServerMockRecorder) GetResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockTurbineServiceServer)(nil).GetResource), arg0, arg1)
}

// GetSpec mocks base method.
func (m *MockTurbineServiceServer) GetSpec(arg0 context.Context, arg1 *core.GetSpecRequest) (*core.GetSpecResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpec", arg0, arg1)
	ret0, _ := ret[0].(*core.GetSpecResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpec indicates an expected call of GetSpec.
func (mr *MockTurbineServiceServerMockRecorder) GetSpec(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpec", reflect.TypeOf((*MockTurbineServiceServer)(nil).GetSpec), arg0, arg1)
}

// HasFunctions mocks base method.
func (m *MockTurbineServiceServer) HasFunctions(arg0 context.Context, arg1 *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasFunctions", arg0, arg1)
	ret0, _ := ret[0].(*wrapperspb.BoolValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasFunctions indicates an expected call of HasFunctions.
func (mr *MockTurbineServiceServerMockRecorder) HasFunctions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasFunctions", reflect.TypeOf((*MockTurbineServiceServer)(nil).HasFunctions), arg0, arg1)
}

// Init mocks base method.
func (m *MockTurbineServiceServer) Init(arg0 context.Context, arg1 *core.InitRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Init indicates an expected call of Init.
func (mr *MockTurbineServiceServerMockRecorder) Init(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockTurbineServiceServer)(nil).Init), arg0, arg1)
}

// ListResources mocks base method.
func (m *MockTurbineServiceServer) ListResources(arg0 context.Context, arg1 *emptypb.Empty) (*core.ListResourcesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResources", arg0, arg1)
	ret0, _ := ret[0].(*core.ListResourcesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResources indicates an expected call of ListResources.
func (mr *MockTurbineServiceServerMockRecorder) ListResources(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResources", reflect.TypeOf((*MockTurbineServiceServer)(nil).ListResources), arg0, arg1)
}

// ReadCollection mocks base method.
func (m *MockTurbineServiceServer) ReadCollection(arg0 context.Context, arg1 *core.ReadCollectionRequest) (*core.Collection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCollection", arg0, arg1)
	ret0, _ := ret[0].(*core.Collection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCollection indicates an expected call of ReadCollection.
func (mr *MockTurbineServiceServerMockRecorder) ReadCollection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCollection", reflect.TypeOf((*MockTurbineServiceServer)(nil).ReadCollection), arg0, arg1)
}

// RegisterSecret mocks base method.
func (m *MockTurbineServiceServer) RegisterSecret(arg0 context.Context, arg1 *core.Secret) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterSecret", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterSecret indicates an expected call of RegisterSecret.
func (mr *MockTurbineServiceServerMockRecorder) RegisterSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterSecret", reflect.TypeOf((*MockTurbineServiceServer)(nil).RegisterSecret), arg0, arg1)
}

// WriteCollectionToResource mocks base method.
func (m *MockTurbineServiceServer) WriteCollectionToResource(arg0 context.Context, arg1 *core.WriteCollectionRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteCollectionToResource", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteCollectionToResource indicates an expected call of WriteCollectionToResource.
func (mr *MockTurbineServiceServerMockRecorder) WriteCollectionToResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteCollectionToResource", reflect.TypeOf((*MockTurbineServiceServer)(nil).WriteCollectionToResource), arg0, arg1)
}

// MockUnsafeTurbineServiceServer is a mock of UnsafeTurbineServiceServer interface.
type MockUnsafeTurbineServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeTurbineServiceServerMockRecorder
}

// MockUnsafeTurbineServiceServerMockRecorder is the mock recorder for MockUnsafeTurbineServiceServer.
type MockUnsafeTurbineServiceServerMockRecorder struct {
	mock *MockUnsafeTurbineServiceServer
}

// NewMockUnsafeTurbineServiceServer creates a new mock instance.
func NewMockUnsafeTurbineServiceServer(ctrl *gomock.Controller) *MockUnsafeTurbineServiceServer {
	mock := &MockUnsafeTurbineServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeTurbineServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeTurbineServiceServer) EXPECT() *MockUnsafeTurbineServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedTurbineServiceServer mocks base method.
func (m *MockUnsafeTurbineServiceServer) mustEmbedUnimplementedTurbineServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTurbineServiceServer")
}

// mustEmbedUnimplementedTurbineServiceServer indicates an expected call of mustEmbedUnimplementedTurbineServiceServer.
func (mr *MockUnsafeTurbineServiceServerMockRecorder) mustEmbedUnimplementedTurbineServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTurbineServiceServer", reflect.TypeOf((*MockUnsafeTurbineServiceServer)(nil).mustEmbedUnimplementedTurbineServiceServer))
}