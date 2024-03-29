// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jianghushinian/blog-go-example/test/testable/reader (interfaces: IReader)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIReader is a mock of IReader interface.
type MockIReader struct {
	ctrl     *gomock.Controller
	recorder *MockIReaderMockRecorder
}

// MockIReaderMockRecorder is the mock recorder for MockIReader.
type MockIReaderMockRecorder struct {
	mock *MockIReader
}

// NewMockIReader creates a new mock instance.
func NewMockIReader(ctrl *gomock.Controller) *MockIReader {
	mock := &MockIReader{ctrl: ctrl}
	mock.recorder = &MockIReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIReader) EXPECT() *MockIReaderMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockIReader) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockIReaderMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIReader)(nil).Read), arg0)
}
