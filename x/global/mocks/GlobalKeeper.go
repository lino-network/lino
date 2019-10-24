// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	amino "github.com/tendermint/go-amino"

	linotypes "github.com/lino-network/lino/types"

	mock "github.com/stretchr/testify/mock"

	model "github.com/lino-network/lino/x/global/model"

	types "github.com/cosmos/cosmos-sdk/types"
)

// GlobalKeeper is an autogenerated mock type for the GlobalKeeper type
type GlobalKeeper struct {
	mock.Mock
}

// ExecuteEvents provides a mock function with given fields: ctx, exec
func (_m *GlobalKeeper) ExecuteEvents(ctx types.Context, exec func(types.Context, linotypes.Event) types.Error) {
	_m.Called(ctx, exec)
}

// ExportToFile provides a mock function with given fields: ctx, cdc, filepath
func (_m *GlobalKeeper) ExportToFile(ctx types.Context, cdc *amino.Codec, filepath string) error {
	ret := _m.Called(ctx, cdc, filepath)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *amino.Codec, string) error); ok {
		r0 = rf(ctx, cdc, filepath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBCEventErrors provides a mock function with given fields: ctx
func (_m *GlobalKeeper) GetBCEventErrors(ctx types.Context) []linotypes.BCEventErr {
	ret := _m.Called(ctx)

	var r0 []linotypes.BCEventErr
	if rf, ok := ret.Get(0).(func(types.Context) []linotypes.BCEventErr); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]linotypes.BCEventErr)
		}
	}

	return r0
}

// GetEventErrors provides a mock function with given fields: ctx
func (_m *GlobalKeeper) GetEventErrors(ctx types.Context) []linotypes.EventError {
	ret := _m.Called(ctx)

	var r0 []linotypes.EventError
	if rf, ok := ret.Get(0).(func(types.Context) []linotypes.EventError); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]linotypes.EventError)
		}
	}

	return r0
}

// GetGlobalTime provides a mock function with given fields: ctx
func (_m *GlobalKeeper) GetGlobalTime(ctx types.Context) model.GlobalTime {
	ret := _m.Called(ctx)

	var r0 model.GlobalTime
	if rf, ok := ret.Get(0).(func(types.Context) model.GlobalTime); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(model.GlobalTime)
	}

	return r0
}

// GetLastBlockTime provides a mock function with given fields: ctx
func (_m *GlobalKeeper) GetLastBlockTime(ctx types.Context) int64 {
	ret := _m.Called(ctx)

	var r0 int64
	if rf, ok := ret.Get(0).(func(types.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// GetPastDay provides a mock function with given fields: ctx, unixTime
func (_m *GlobalKeeper) GetPastDay(ctx types.Context, unixTime int64) int64 {
	ret := _m.Called(ctx, unixTime)

	var r0 int64
	if rf, ok := ret.Get(0).(func(types.Context, int64) int64); ok {
		r0 = rf(ctx, unixTime)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// ImportFromFile provides a mock function with given fields: ctx, cdc, filepath
func (_m *GlobalKeeper) ImportFromFile(ctx types.Context, cdc *amino.Codec, filepath string) error {
	ret := _m.Called(ctx, cdc, filepath)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *amino.Codec, string) error); ok {
		r0 = rf(ctx, cdc, filepath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitGenesis provides a mock function with given fields: ctx
func (_m *GlobalKeeper) InitGenesis(ctx types.Context) {
	_m.Called(ctx)
}

// OnBeginBlock provides a mock function with given fields: ctx
func (_m *GlobalKeeper) OnBeginBlock(ctx types.Context) {
	_m.Called(ctx)
}

// OnEndBlock provides a mock function with given fields: ctx
func (_m *GlobalKeeper) OnEndBlock(ctx types.Context) {
	_m.Called(ctx)
}

// RegisterEventAtTime provides a mock function with given fields: ctx, unixTime, event
func (_m *GlobalKeeper) RegisterEventAtTime(ctx types.Context, unixTime int64, event linotypes.Event) types.Error {
	ret := _m.Called(ctx, unixTime, event)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, int64, linotypes.Event) types.Error); ok {
		r0 = rf(ctx, unixTime, event)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}
