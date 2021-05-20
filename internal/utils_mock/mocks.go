// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kbence/conndetect/internal/utils (interfaces: Printer,Time)

// Package utils_mock is a generated GoMock package.
package utils_mock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockPrinter is a mock of Printer interface.
type MockPrinter struct {
	ctrl     *gomock.Controller
	recorder *MockPrinterMockRecorder
}

// MockPrinterMockRecorder is the mock recorder for MockPrinter.
type MockPrinterMockRecorder struct {
	mock *MockPrinter
}

// NewMockPrinter creates a new mock instance.
func NewMockPrinter(ctrl *gomock.Controller) *MockPrinter {
	mock := &MockPrinter{ctrl: ctrl}
	mock.recorder = &MockPrinterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrinter) EXPECT() *MockPrinterMockRecorder {
	return m.recorder
}

// Printf mocks base method.
func (m *MockPrinter) Printf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Printf", varargs...)
}

// Printf indicates an expected call of Printf.
func (mr *MockPrinterMockRecorder) Printf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Printf", reflect.TypeOf((*MockPrinter)(nil).Printf), varargs...)
}

// MockTime is a mock of Time interface.
type MockTime struct {
	ctrl     *gomock.Controller
	recorder *MockTimeMockRecorder
}

// MockTimeMockRecorder is the mock recorder for MockTime.
type MockTimeMockRecorder struct {
	mock *MockTime
}

// NewMockTime creates a new mock instance.
func NewMockTime(ctrl *gomock.Controller) *MockTime {
	mock := &MockTime{ctrl: ctrl}
	mock.recorder = &MockTimeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTime) EXPECT() *MockTimeMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockTime) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockTimeMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockTime)(nil).Now))
}