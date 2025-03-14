// Code generated by MockGen. DO NOT EDIT.
// Source: /home/anton/tbank-tech/scrapper/internal/hub/hub.go

// Package hub is a generated GoMock package.
package hub

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHub is a mock of Hub interface.
type MockHub struct {
	ctrl     *gomock.Controller
	recorder *MockHubMockRecorder
}

// MockHubMockRecorder is the mock recorder for MockHub.
type MockHubMockRecorder struct {
	mock *MockHub
}

// NewMockHub creates a new mock instance.
func NewMockHub(ctrl *gomock.Controller) *MockHub {
	mock := &MockHub{ctrl: ctrl}
	mock.recorder = &MockHubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHub) EXPECT() *MockHubMockRecorder {
	return m.recorder
}

// AddLink mocks base method.
func (m *MockHub) AddLink(link string, userID uint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddLink", link, userID)
}

// AddLink indicates an expected call of AddLink.
func (mr *MockHubMockRecorder) AddLink(link, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLink", reflect.TypeOf((*MockHub)(nil).AddLink), link, userID)
}

// RemoveLink mocks base method.
func (m *MockHub) RemoveLink(link string, userID uint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveLink", link, userID)
}

// RemoveLink indicates an expected call of RemoveLink.
func (mr *MockHubMockRecorder) RemoveLink(link, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLink", reflect.TypeOf((*MockHub)(nil).RemoveLink), link, userID)
}

// Run mocks base method.
func (m *MockHub) Run() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run")
}

// Run indicates an expected call of Run.
func (mr *MockHubMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockHub)(nil).Run))
}

// Stop mocks base method.
func (m *MockHub) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockHubMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockHub)(nil).Stop))
}
