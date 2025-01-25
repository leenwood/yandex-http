package mocks

import (
	gomock "github.com/golang/mock/gomock"
	url "leenwood/yandex-http/internal/domain/url"
	reflect "reflect"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// FindById mocks base method
func (m *MockRepositoryInterface) FindById(id string) (*url.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", id)
	ret0, _ := ret[0].(*url.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById
func (mr *MockRepositoryInterfaceMockRecorder) FindById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockRepositoryInterface)(nil).FindById), id)
}

// FindByUrl mocks base method
func (m *MockRepositoryInterface) FindByUrl(originalUrl string) (*url.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUrl", originalUrl)
	ret0, _ := ret[0].(*url.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUrl indicates an expected call of FindByUrl
func (mr *MockRepositoryInterfaceMockRecorder) FindByUrl(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUrl", reflect.TypeOf((*MockRepositoryInterface)(nil).FindByUrl), url)
}

// Save mocks base method
func (m *MockRepositoryInterface) Save(originalUrl string, shortUuid string) (*url.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", originalUrl, shortUuid)
	ret0, _ := ret[0].(*url.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save
func (mr *MockRepositoryInterfaceMockRecorder) Save(originalUrl, shortUuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepositoryInterface)(nil).Save), originalUrl, shortUuid)
}

func (m *MockRepositoryInterface) FindAll(page, limit int) ([]*url.Url, error) {
	ret := m.ctrl.Call(m, "FindAll", page, limit)
	ret0, _ := ret[0].([]*url.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockRepositoryInterfaceMockRecorder) FindAll(page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockRepositoryInterface)(nil).FindAll), page, limit)
}

func (m *MockRepositoryInterface) Update(originalUrl *url.Url) (*url.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", originalUrl)
	ret0, _ := ret[0].(*url.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockRepositoryInterfaceMockRecorder) Update(originalUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepositoryInterface)(nil).Update), originalUrl)
}
