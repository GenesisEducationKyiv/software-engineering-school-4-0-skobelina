// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/mails/service.go

// Package mails is a generated GoMock package.
package mails

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMailServiceInterface is a mock of MailServiceInterface interface.
type MockMailServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMailServiceInterfaceMockRecorder
}

// MockMailServiceInterfaceMockRecorder is the mock recorder for MockMailServiceInterface.
type MockMailServiceInterfaceMockRecorder struct {
	mock *MockMailServiceInterface
}

// NewMockMailServiceInterface creates a new mock instance.
func NewMockMailServiceInterface(ctrl *gomock.Controller) *MockMailServiceInterface {
	mock := &MockMailServiceInterface{ctrl: ctrl}
	mock.recorder = &MockMailServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailServiceInterface) EXPECT() *MockMailServiceInterfaceMockRecorder {
	return m.recorder
}

// SendBatch mocks base method.
func (m *MockMailServiceInterface) SendBatch(messages ...*Message) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range messages {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendBatch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendBatch indicates an expected call of SendBatch.
func (mr *MockMailServiceInterfaceMockRecorder) SendBatch(messages ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendBatch", reflect.TypeOf((*MockMailServiceInterface)(nil).SendBatch), messages...)
}

// SendEmail mocks base method.
func (m *MockMailServiceInterface) SendEmail(recipients []string, subject string, temp Template) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", recipients, subject, temp)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockMailServiceInterfaceMockRecorder) SendEmail(recipients, subject, temp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockMailServiceInterface)(nil).SendEmail), recipients, subject, temp)
}
