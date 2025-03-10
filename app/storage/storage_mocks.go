// Code generated by MockGen. DO NOT EDIT.
// Source: app/storage/storage.go
//
// Generated by this command:
//
//	mockgen -source=app/storage/storage.go -destination=app/storage/storage_mocks.go -package=storage Engine WAL
//

// Package storage is a generated GoMock package.
package storage

import (
	wal "concurrency/app/storage/wal"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEngine is a mock of Engine interface.
type MockEngine struct {
	ctrl     *gomock.Controller
	recorder *MockEngineMockRecorder
}

// MockEngineMockRecorder is the mock recorder for MockEngine.
type MockEngineMockRecorder struct {
	mock *MockEngine
}

// NewMockEngine creates a new mock instance.
func NewMockEngine(ctrl *gomock.Controller) *MockEngine {
	mock := &MockEngine{ctrl: ctrl}
	mock.recorder = &MockEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEngine) EXPECT() *MockEngineMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockEngine) Delete(ctx context.Context, key string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", ctx, key)
}

// Delete indicates an expected call of Delete.
func (mr *MockEngineMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockEngine)(nil).Delete), ctx, key)
}

// Get mocks base method.
func (m *MockEngine) Get(ctx context.Context, key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockEngineMockRecorder) Get(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockEngine)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockEngine) Set(ctx context.Context, key, value string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", ctx, key, value)
}

// Set indicates an expected call of Set.
func (mr *MockEngineMockRecorder) Set(ctx, key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockEngine)(nil).Set), ctx, key, value)
}

// MockWAL is a mock of WAL interface.
type MockWAL struct {
	ctrl     *gomock.Controller
	recorder *MockWALMockRecorder
}

// MockWALMockRecorder is the mock recorder for MockWAL.
type MockWALMockRecorder struct {
	mock *MockWAL
}

// NewMockWAL creates a new mock instance.
func NewMockWAL(ctrl *gomock.Controller) *MockWAL {
	mock := &MockWAL{ctrl: ctrl}
	mock.recorder = &MockWALMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWAL) EXPECT() *MockWALMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockWAL) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockWALMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockWAL)(nil).Delete), ctx, key)
}

// Recover mocks base method.
func (m *MockWAL) Recover() ([]wal.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recover")
	ret0, _ := ret[0].([]wal.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recover indicates an expected call of Recover.
func (mr *MockWALMockRecorder) Recover() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recover", reflect.TypeOf((*MockWAL)(nil).Recover))
}

// Set mocks base method.
func (m *MockWAL) Set(ctx context.Context, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockWALMockRecorder) Set(ctx, key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockWAL)(nil).Set), ctx, key, value)
}
