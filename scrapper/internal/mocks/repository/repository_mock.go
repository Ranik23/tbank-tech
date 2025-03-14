// Code generated by MockGen. DO NOT EDIT.
// Source: /home/anton/tbank-tech/scrapper/internal/repository/repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"
	models "tbank/scrapper/internal/models"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v5"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *MockRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockRepositoryMockRecorder) BeginTx(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockRepository)(nil).BeginTx), ctx)
}

// CommitTx mocks base method.
func (m *MockRepository) CommitTx(ctx context.Context, tx pgx.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitTx", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitTx indicates an expected call of CommitTx.
func (mr *MockRepositoryMockRecorder) CommitTx(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitTx", reflect.TypeOf((*MockRepository)(nil).CommitTx), ctx, tx)
}

// CreateLink mocks base method.
func (m *MockRepository) CreateLink(ctx context.Context, link string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLink", ctx, link)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLink indicates an expected call of CreateLink.
func (mr *MockRepositoryMockRecorder) CreateLink(ctx, link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLink", reflect.TypeOf((*MockRepository)(nil).CreateLink), ctx, link)
}

// CreateLinkUser mocks base method.
func (m *MockRepository) CreateLinkUser(ctx context.Context, linkID, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLinkUser", ctx, linkID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLinkUser indicates an expected call of CreateLinkUser.
func (mr *MockRepositoryMockRecorder) CreateLinkUser(ctx, linkID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLinkUser", reflect.TypeOf((*MockRepository)(nil).CreateLinkUser), ctx, linkID, userID)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, userID uint, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, userID, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, userID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, userID, name)
}

// DeleteLink mocks base method.
func (m *MockRepository) DeleteLink(ctx context.Context, linkID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLink", ctx, linkID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLink indicates an expected call of DeleteLink.
func (mr *MockRepositoryMockRecorder) DeleteLink(ctx, linkID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLink", reflect.TypeOf((*MockRepository)(nil).DeleteLink), ctx, linkID)
}

// DeleteLinkUser mocks base method.
func (m *MockRepository) DeleteLinkUser(ctx context.Context, linkID, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLinkUser", ctx, linkID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLinkUser indicates an expected call of DeleteLinkUser.
func (mr *MockRepositoryMockRecorder) DeleteLinkUser(ctx, linkID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLinkUser", reflect.TypeOf((*MockRepository)(nil).DeleteLinkUser), ctx, linkID, userID)
}

// DeleteUser mocks base method.
func (m *MockRepository) DeleteUser(ctx context.Context, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockRepositoryMockRecorder) DeleteUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockRepository)(nil).DeleteUser), ctx, userID)
}

// GetLinkByID mocks base method.
func (m *MockRepository) GetLinkByID(ctx context.Context, id uint) (*models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkByID", ctx, id)
	ret0, _ := ret[0].(*models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkByID indicates an expected call of GetLinkByID.
func (mr *MockRepositoryMockRecorder) GetLinkByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkByID", reflect.TypeOf((*MockRepository)(nil).GetLinkByID), ctx, id)
}

// GetLinkByURL mocks base method.
func (m *MockRepository) GetLinkByURL(ctx context.Context, url string) (*models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkByURL", ctx, url)
	ret0, _ := ret[0].(*models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkByURL indicates an expected call of GetLinkByURL.
func (mr *MockRepositoryMockRecorder) GetLinkByURL(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkByURL", reflect.TypeOf((*MockRepository)(nil).GetLinkByURL), ctx, url)
}

// GetURLS mocks base method.
func (m *MockRepository) GetURLS(ctx context.Context, userID uint) ([]models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLS", ctx, userID)
	ret0, _ := ret[0].([]models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLS indicates an expected call of GetURLS.
func (mr *MockRepositoryMockRecorder) GetURLS(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLS", reflect.TypeOf((*MockRepository)(nil).GetURLS), ctx, userID)
}

// RollbackTx mocks base method.
func (m *MockRepository) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackTx", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RollbackTx indicates an expected call of RollbackTx.
func (mr *MockRepositoryMockRecorder) RollbackTx(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTx", reflect.TypeOf((*MockRepository)(nil).RollbackTx), ctx, tx)
}
