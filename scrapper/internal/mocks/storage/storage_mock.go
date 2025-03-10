// Code generated by MockGen. DO NOT EDIT.
// Source: /home/anton/tbank-tech/scrapper/internal/storage/storage.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"
	models "tbank/scrapper/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// CreateLink mocks base method.
func (m *MockStorage) CreateLink(ctx context.Context, link string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLink", ctx, link)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLink indicates an expected call of CreateLink.
func (mr *MockStorageMockRecorder) CreateLink(ctx, link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLink", reflect.TypeOf((*MockStorage)(nil).CreateLink), ctx, link)
}

// CreateLinkUser mocks base method.
func (m *MockStorage) CreateLinkUser(ctx context.Context, linkID, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLinkUser", ctx, linkID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLinkUser indicates an expected call of CreateLinkUser.
func (mr *MockStorageMockRecorder) CreateLinkUser(ctx, linkID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLinkUser", reflect.TypeOf((*MockStorage)(nil).CreateLinkUser), ctx, linkID, userID)
}

// CreateUser mocks base method.
func (m *MockStorage) CreateUser(ctx context.Context, userID uint, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, userID, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStorageMockRecorder) CreateUser(ctx, userID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStorage)(nil).CreateUser), ctx, userID, name)
}

// DeleteLink mocks base method.
func (m *MockStorage) DeleteLink(ctx context.Context, linkID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLink", ctx, linkID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLink indicates an expected call of DeleteLink.
func (mr *MockStorageMockRecorder) DeleteLink(ctx, linkID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLink", reflect.TypeOf((*MockStorage)(nil).DeleteLink), ctx, linkID)
}

// DeleteLinkUser mocks base method.
func (m *MockStorage) DeleteLinkUser(ctx context.Context, linkID, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLinkUser", ctx, linkID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLinkUser indicates an expected call of DeleteLinkUser.
func (mr *MockStorageMockRecorder) DeleteLinkUser(ctx, linkID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLinkUser", reflect.TypeOf((*MockStorage)(nil).DeleteLinkUser), ctx, linkID, userID)
}

// DeleteUser mocks base method.
func (m *MockStorage) DeleteUser(ctx context.Context, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStorageMockRecorder) DeleteUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStorage)(nil).DeleteUser), ctx, userID)
}

// GetLinkByID mocks base method.
func (m *MockStorage) GetLinkByID(ctx context.Context, id uint) (*models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkByID", ctx, id)
	ret0, _ := ret[0].(*models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkByID indicates an expected call of GetLinkByID.
func (mr *MockStorageMockRecorder) GetLinkByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkByID", reflect.TypeOf((*MockStorage)(nil).GetLinkByID), ctx, id)
}

// GetLinkByURL mocks base method.
func (m *MockStorage) GetLinkByURL(ctx context.Context, url string) (*models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkByURL", ctx, url)
	ret0, _ := ret[0].(*models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkByURL indicates an expected call of GetLinkByURL.
func (mr *MockStorageMockRecorder) GetLinkByURL(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkByURL", reflect.TypeOf((*MockStorage)(nil).GetLinkByURL), ctx, url)
}

// GetURLS mocks base method.
func (m *MockStorage) GetURLS(ctx context.Context, userID uint) ([]models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLS", ctx, userID)
	ret0, _ := ret[0].([]models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLS indicates an expected call of GetURLS.
func (mr *MockStorageMockRecorder) GetURLS(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLS", reflect.TypeOf((*MockStorage)(nil).GetURLS), ctx, userID)
}
