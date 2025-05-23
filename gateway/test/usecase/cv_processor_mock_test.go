// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/proto/v1/cv_processor_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=pkg/proto/v1/cv_processor_grpc.pb.go -destination=gateway/test/usecase/cv_processor_mock_test.go -package=usecase_test
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

// MockCVProcessorServiceClient is a mock of CVProcessorServiceClient interface.
type MockCVProcessorServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCVProcessorServiceClientMockRecorder
	isgomock struct{}
}

// MockCVProcessorServiceClientMockRecorder is the mock recorder for MockCVProcessorServiceClient.
type MockCVProcessorServiceClientMockRecorder struct {
	mock *MockCVProcessorServiceClient
}

// NewMockCVProcessorServiceClient creates a new mock instance.
func NewMockCVProcessorServiceClient(ctrl *gomock.Controller) *MockCVProcessorServiceClient {
	mock := &MockCVProcessorServiceClient{ctrl: ctrl}
	mock.recorder = &MockCVProcessorServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCVProcessorServiceClient) EXPECT() *MockCVProcessorServiceClientMockRecorder {
	return m.recorder
}

// FetchSummarizedPdfHistory mocks base method.
func (m *MockCVProcessorServiceClient) FetchSummarizedPdfHistory(ctx context.Context, in *v1.FetchSummarizedPdfHistoryRequest, opts ...grpc.CallOption) (*v1.BatchPDFProcessResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchSummarizedPdfHistory", varargs...)
	ret0, _ := ret[0].(*v1.BatchPDFProcessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchSummarizedPdfHistory indicates an expected call of FetchSummarizedPdfHistory.
func (mr *MockCVProcessorServiceClientMockRecorder) FetchSummarizedPdfHistory(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchSummarizedPdfHistory", reflect.TypeOf((*MockCVProcessorServiceClient)(nil).FetchSummarizedPdfHistory), varargs...)
}

// ProcessBatchPDF mocks base method.
func (m *MockCVProcessorServiceClient) ProcessBatchPDF(ctx context.Context, in *v1.BatchPDFProcessRequest, opts ...grpc.CallOption) (*v1.BatchPDFProcessResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ProcessBatchPDF", varargs...)
	ret0, _ := ret[0].(*v1.BatchPDFProcessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessBatchPDF indicates an expected call of ProcessBatchPDF.
func (mr *MockCVProcessorServiceClientMockRecorder) ProcessBatchPDF(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessBatchPDF", reflect.TypeOf((*MockCVProcessorServiceClient)(nil).ProcessBatchPDF), varargs...)
}

// MockCVProcessorServiceServer is a mock of CVProcessorServiceServer interface.
type MockCVProcessorServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockCVProcessorServiceServerMockRecorder
	isgomock struct{}
}

// MockCVProcessorServiceServerMockRecorder is the mock recorder for MockCVProcessorServiceServer.
type MockCVProcessorServiceServerMockRecorder struct {
	mock *MockCVProcessorServiceServer
}

// NewMockCVProcessorServiceServer creates a new mock instance.
func NewMockCVProcessorServiceServer(ctrl *gomock.Controller) *MockCVProcessorServiceServer {
	mock := &MockCVProcessorServiceServer{ctrl: ctrl}
	mock.recorder = &MockCVProcessorServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCVProcessorServiceServer) EXPECT() *MockCVProcessorServiceServerMockRecorder {
	return m.recorder
}

// FetchSummarizedPdfHistory mocks base method.
func (m *MockCVProcessorServiceServer) FetchSummarizedPdfHistory(arg0 context.Context, arg1 *v1.FetchSummarizedPdfHistoryRequest) (*v1.BatchPDFProcessResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchSummarizedPdfHistory", arg0, arg1)
	ret0, _ := ret[0].(*v1.BatchPDFProcessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchSummarizedPdfHistory indicates an expected call of FetchSummarizedPdfHistory.
func (mr *MockCVProcessorServiceServerMockRecorder) FetchSummarizedPdfHistory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchSummarizedPdfHistory", reflect.TypeOf((*MockCVProcessorServiceServer)(nil).FetchSummarizedPdfHistory), arg0, arg1)
}

// ProcessBatchPDF mocks base method.
func (m *MockCVProcessorServiceServer) ProcessBatchPDF(arg0 context.Context, arg1 *v1.BatchPDFProcessRequest) (*v1.BatchPDFProcessResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessBatchPDF", arg0, arg1)
	ret0, _ := ret[0].(*v1.BatchPDFProcessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessBatchPDF indicates an expected call of ProcessBatchPDF.
func (mr *MockCVProcessorServiceServerMockRecorder) ProcessBatchPDF(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessBatchPDF", reflect.TypeOf((*MockCVProcessorServiceServer)(nil).ProcessBatchPDF), arg0, arg1)
}

// mustEmbedUnimplementedCVProcessorServiceServer mocks base method.
func (m *MockCVProcessorServiceServer) mustEmbedUnimplementedCVProcessorServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCVProcessorServiceServer")
}

// mustEmbedUnimplementedCVProcessorServiceServer indicates an expected call of mustEmbedUnimplementedCVProcessorServiceServer.
func (mr *MockCVProcessorServiceServerMockRecorder) mustEmbedUnimplementedCVProcessorServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCVProcessorServiceServer", reflect.TypeOf((*MockCVProcessorServiceServer)(nil).mustEmbedUnimplementedCVProcessorServiceServer))
}

// MockUnsafeCVProcessorServiceServer is a mock of UnsafeCVProcessorServiceServer interface.
type MockUnsafeCVProcessorServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeCVProcessorServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeCVProcessorServiceServerMockRecorder is the mock recorder for MockUnsafeCVProcessorServiceServer.
type MockUnsafeCVProcessorServiceServerMockRecorder struct {
	mock *MockUnsafeCVProcessorServiceServer
}

// NewMockUnsafeCVProcessorServiceServer creates a new mock instance.
func NewMockUnsafeCVProcessorServiceServer(ctrl *gomock.Controller) *MockUnsafeCVProcessorServiceServer {
	mock := &MockUnsafeCVProcessorServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeCVProcessorServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeCVProcessorServiceServer) EXPECT() *MockUnsafeCVProcessorServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedCVProcessorServiceServer mocks base method.
func (m *MockUnsafeCVProcessorServiceServer) mustEmbedUnimplementedCVProcessorServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCVProcessorServiceServer")
}

// mustEmbedUnimplementedCVProcessorServiceServer indicates an expected call of mustEmbedUnimplementedCVProcessorServiceServer.
func (mr *MockUnsafeCVProcessorServiceServerMockRecorder) mustEmbedUnimplementedCVProcessorServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCVProcessorServiceServer", reflect.TypeOf((*MockUnsafeCVProcessorServiceServer)(nil).mustEmbedUnimplementedCVProcessorServiceServer))
}
