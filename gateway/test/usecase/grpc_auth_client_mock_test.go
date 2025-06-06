// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/proto/v1/auth_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=pkg/proto/v1/auth_grpc.pb.go -destination=gateway/test/usecase/grpc_auth_client_mock_test.go -package=usecase_test
//

// Package usecase_test is a generated GoMock package.
package usecase_test

import (
	context "context"
	reflect "reflect"

	v1 "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAuthenticationServiceClient is a mock of AuthenticationServiceClient interface.
type MockAuthenticationServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticationServiceClientMockRecorder
	isgomock struct{}
}

// MockAuthenticationServiceClientMockRecorder is the mock recorder for MockAuthenticationServiceClient.
type MockAuthenticationServiceClientMockRecorder struct {
	mock *MockAuthenticationServiceClient
}

// NewMockAuthenticationServiceClient creates a new mock instance.
func NewMockAuthenticationServiceClient(ctrl *gomock.Controller) *MockAuthenticationServiceClient {
	mock := &MockAuthenticationServiceClient{ctrl: ctrl}
	mock.recorder = &MockAuthenticationServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticationServiceClient) EXPECT() *MockAuthenticationServiceClientMockRecorder {
	return m.recorder
}

// UserLogin mocks base method.
func (m *MockAuthenticationServiceClient) UserLogin(ctx context.Context, in *v1.LoginRequest, opts ...grpc.CallOption) (*v1.LoginResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UserLogin", varargs...)
	ret0, _ := ret[0].(*v1.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserLogin indicates an expected call of UserLogin.
func (mr *MockAuthenticationServiceClientMockRecorder) UserLogin(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserLogin", reflect.TypeOf((*MockAuthenticationServiceClient)(nil).UserLogin), varargs...)
}

// UserRegister mocks base method.
func (m *MockAuthenticationServiceClient) UserRegister(ctx context.Context, in *v1.RegisterRequest, opts ...grpc.CallOption) (*v1.RegisterResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UserRegister", varargs...)
	ret0, _ := ret[0].(*v1.RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserRegister indicates an expected call of UserRegister.
func (mr *MockAuthenticationServiceClientMockRecorder) UserRegister(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserRegister", reflect.TypeOf((*MockAuthenticationServiceClient)(nil).UserRegister), varargs...)
}

// MockAuthenticationServiceServer is a mock of AuthenticationServiceServer interface.
type MockAuthenticationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticationServiceServerMockRecorder
	isgomock struct{}
}

// MockAuthenticationServiceServerMockRecorder is the mock recorder for MockAuthenticationServiceServer.
type MockAuthenticationServiceServerMockRecorder struct {
	mock *MockAuthenticationServiceServer
}

// NewMockAuthenticationServiceServer creates a new mock instance.
func NewMockAuthenticationServiceServer(ctrl *gomock.Controller) *MockAuthenticationServiceServer {
	mock := &MockAuthenticationServiceServer{ctrl: ctrl}
	mock.recorder = &MockAuthenticationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticationServiceServer) EXPECT() *MockAuthenticationServiceServerMockRecorder {
	return m.recorder
}

// UserLogin mocks base method.
func (m *MockAuthenticationServiceServer) UserLogin(arg0 context.Context, arg1 *v1.LoginRequest) (*v1.LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserLogin", arg0, arg1)
	ret0, _ := ret[0].(*v1.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserLogin indicates an expected call of UserLogin.
func (mr *MockAuthenticationServiceServerMockRecorder) UserLogin(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserLogin", reflect.TypeOf((*MockAuthenticationServiceServer)(nil).UserLogin), arg0, arg1)
}

// UserRegister mocks base method.
func (m *MockAuthenticationServiceServer) UserRegister(arg0 context.Context, arg1 *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserRegister", arg0, arg1)
	ret0, _ := ret[0].(*v1.RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserRegister indicates an expected call of UserRegister.
func (mr *MockAuthenticationServiceServerMockRecorder) UserRegister(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserRegister", reflect.TypeOf((*MockAuthenticationServiceServer)(nil).UserRegister), arg0, arg1)
}

// mustEmbedUnimplementedAuthenticationServiceServer mocks base method.
func (m *MockAuthenticationServiceServer) mustEmbedUnimplementedAuthenticationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthenticationServiceServer")
}

// mustEmbedUnimplementedAuthenticationServiceServer indicates an expected call of mustEmbedUnimplementedAuthenticationServiceServer.
func (mr *MockAuthenticationServiceServerMockRecorder) mustEmbedUnimplementedAuthenticationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthenticationServiceServer", reflect.TypeOf((*MockAuthenticationServiceServer)(nil).mustEmbedUnimplementedAuthenticationServiceServer))
}

// MockUnsafeAuthenticationServiceServer is a mock of UnsafeAuthenticationServiceServer interface.
type MockUnsafeAuthenticationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeAuthenticationServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeAuthenticationServiceServerMockRecorder is the mock recorder for MockUnsafeAuthenticationServiceServer.
type MockUnsafeAuthenticationServiceServerMockRecorder struct {
	mock *MockUnsafeAuthenticationServiceServer
}

// NewMockUnsafeAuthenticationServiceServer creates a new mock instance.
func NewMockUnsafeAuthenticationServiceServer(ctrl *gomock.Controller) *MockUnsafeAuthenticationServiceServer {
	mock := &MockUnsafeAuthenticationServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeAuthenticationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeAuthenticationServiceServer) EXPECT() *MockUnsafeAuthenticationServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedAuthenticationServiceServer mocks base method.
func (m *MockUnsafeAuthenticationServiceServer) mustEmbedUnimplementedAuthenticationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthenticationServiceServer")
}

// mustEmbedUnimplementedAuthenticationServiceServer indicates an expected call of mustEmbedUnimplementedAuthenticationServiceServer.
func (mr *MockUnsafeAuthenticationServiceServerMockRecorder) mustEmbedUnimplementedAuthenticationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthenticationServiceServer", reflect.TypeOf((*MockUnsafeAuthenticationServiceServer)(nil).mustEmbedUnimplementedAuthenticationServiceServer))
}
