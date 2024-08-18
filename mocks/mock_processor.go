// Code generated by MockGen. DO NOT EDIT.
// Source: zt-event-logger/pkg/events (interfaces: Processor)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=mock_processor.go -package=mocks zt-event-logger/pkg/events Processor
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	events "zt-event-logger/pkg/events"

	ztchooks "github.com/zerotier/ztchooks"
	gomock "go.uber.org/mock/gomock"
)

// MockProcessor is a mock of Processor interface.
type MockProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockProcessorMockRecorder
}

// MockProcessorMockRecorder is the mock recorder for MockProcessor.
type MockProcessorMockRecorder struct {
	mock *MockProcessor
}

// NewMockProcessor creates a new mock instance.
func NewMockProcessor(ctrl *gomock.Controller) *MockProcessor {
	mock := &MockProcessor{ctrl: ctrl}
	mock.recorder = &MockProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessor) EXPECT() *MockProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockProcessor) Process(arg0 []byte, arg1 ...events.SignatureOpt) (*ztchooks.HookBase, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Process", varargs...)
	ret0, _ := ret[0].(*ztchooks.HookBase)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockProcessorMockRecorder) Process(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockProcessor)(nil).Process), varargs...)
}
