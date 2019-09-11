// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	linotypes "github.com/lino-network/lino/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// BandwidthKeeper is an autogenerated mock type for the BandwidthKeeper type
type BandwidthKeeper struct {
	mock.Mock
}

// BeginBlocker provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) BeginBlocker(ctx types.Context) types.Error {
	ret := _m.Called(ctx)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context) types.Error); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// CheckBandwidth provides a mock function with given fields: ctx, accKey, fee
func (_m *BandwidthKeeper) CheckBandwidth(ctx types.Context, accKey linotypes.AccountKey, fee authtypes.StdFee) types.Error {
	ret := _m.Called(ctx, accKey, fee)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, authtypes.StdFee) types.Error); ok {
		r0 = rf(ctx, accKey, fee)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// DecayMaxMPS provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) DecayMaxMPS(ctx types.Context) types.Error {
	ret := _m.Called(ctx)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context) types.Error); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// EndBlocker provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) EndBlocker(ctx types.Context) types.Error {
	ret := _m.Called(ctx)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context) types.Error); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// InitGenesis provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) InitGenesis(ctx types.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReCalculateAppBandwidthInfo provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) ReCalculateAppBandwidthInfo(ctx types.Context) types.Error {
	ret := _m.Called(ctx)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context) types.Error); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}
