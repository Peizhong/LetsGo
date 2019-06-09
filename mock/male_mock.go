// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/peizhong/letsgo/person (interfaces: Male)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMale is a mock of Male interface
type MockMale struct {
	ctrl     *gomock.Controller
	recorder *MockMaleMockRecorder
}

// MockMaleMockRecorder is the mock recorder for MockMale
type MockMaleMockRecorder struct {
	mock *MockMale
}

// NewMockMale creates a new mock instance
func NewMockMale(ctrl *gomock.Controller) *MockMale {
	mock := &MockMale{ctrl: ctrl}
	mock.recorder = &MockMaleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMale) EXPECT() *MockMaleMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockMale) Get(arg0 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockMaleMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMale)(nil).Get), arg0)
}
