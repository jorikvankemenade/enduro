// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/artefactual-labs/enduro/internal/amclient (interfaces: IngestService)

// Package fake is a generated GoMock package.
package fake

import (
	context "context"
	amclient "github.com/artefactual-labs/enduro/internal/amclient"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIngestService is a mock of IngestService interface
type MockIngestService struct {
	ctrl     *gomock.Controller
	recorder *MockIngestServiceMockRecorder
}

// MockIngestServiceMockRecorder is the mock recorder for MockIngestService
type MockIngestServiceMockRecorder struct {
	mock *MockIngestService
}

// NewMockIngestService creates a new mock instance
func NewMockIngestService(ctrl *gomock.Controller) *MockIngestService {
	mock := &MockIngestService{ctrl: ctrl}
	mock.recorder = &MockIngestServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIngestService) EXPECT() *MockIngestServiceMockRecorder {
	return m.recorder
}

// Hide mocks base method
func (m *MockIngestService) Hide(arg0 context.Context, arg1 string) (*amclient.IngestHideResponse, *amclient.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hide", arg0, arg1)
	ret0, _ := ret[0].(*amclient.IngestHideResponse)
	ret1, _ := ret[1].(*amclient.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Hide indicates an expected call of Hide
func (mr *MockIngestServiceMockRecorder) Hide(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hide", reflect.TypeOf((*MockIngestService)(nil).Hide), arg0, arg1)
}

// Status mocks base method
func (m *MockIngestService) Status(arg0 context.Context, arg1 string) (*amclient.IngestStatusResponse, *amclient.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", arg0, arg1)
	ret0, _ := ret[0].(*amclient.IngestStatusResponse)
	ret1, _ := ret[1].(*amclient.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Status indicates an expected call of Status
func (mr *MockIngestServiceMockRecorder) Status(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockIngestService)(nil).Status), arg0, arg1)
}