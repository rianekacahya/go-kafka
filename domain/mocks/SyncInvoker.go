// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SyncInvoker is an autogenerated mock type for the SyncInvoker type
type SyncInvoker struct {
	mock.Mock
}

// EventListener provides a mock function with given fields: _a0, _a1
func (_m *SyncInvoker) EventListener(_a0 context.Context, _a1 string) <-chan error {
	ret := _m.Called(_a0, _a1)

	var r0 <-chan error
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan error); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan error)
		}
	}

	return r0
}

// EventWhisper provides a mock function with given fields: _a0, _a1
func (_m *SyncInvoker) EventWhisper(_a0 context.Context, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
