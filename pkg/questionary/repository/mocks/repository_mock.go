// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/ismaeljpv/qa-api/pkg/questionary/domain"
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

// AddAnswer mocks base method.
func (m *MockRepository) AddAnswer(ctx context.Context, answer domain.Answer) (domain.QuestionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAnswer", ctx, answer)
	ret0, _ := ret[0].(domain.QuestionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAnswer indicates an expected call of AddAnswer.
func (mr *MockRepositoryMockRecorder) AddAnswer(ctx, answer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAnswer", reflect.TypeOf((*MockRepository)(nil).AddAnswer), ctx, answer)
}

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, question domain.Question) (domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, question)
	ret0, _ := ret[0].(domain.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, question interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, question)
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockRepository) FindAll(ctx context.Context) ([]domain.QuestionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]domain.QuestionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockRepository)(nil).FindAll), ctx)
}

// FindByID mocks base method.
func (m *MockRepository) FindByID(ctx context.Context, id string) (domain.QuestionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(domain.QuestionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockRepository)(nil).FindByID), ctx, id)
}

// FindByUser mocks base method.
func (m *MockRepository) FindByUser(ctx context.Context, userId string) ([]domain.QuestionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUser", ctx, userId)
	ret0, _ := ret[0].([]domain.QuestionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUser indicates an expected call of FindByUser.
func (mr *MockRepositoryMockRecorder) FindByUser(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUser", reflect.TypeOf((*MockRepository)(nil).FindByUser), ctx, userId)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, questionInfo domain.QuestionInfo) (domain.QuestionInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, questionInfo)
	ret0, _ := ret[0].(domain.QuestionInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, questionInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, questionInfo)
}
