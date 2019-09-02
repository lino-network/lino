// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import auth "github.com/cosmos/cosmos-sdk/x/auth"

import linotypes "github.com/lino-network/lino/types"
import mock "github.com/stretchr/testify/mock"
import model "github.com/lino-network/lino/x/bandwidth/model"
import types "github.com/cosmos/cosmos-sdk/types"

// BandwidthKeeper is an autogenerated mock type for the BandwidthKeeper type
type BandwidthKeeper struct {
	mock.Mock
}

// AddMsgSignedByApp provides a mock function with given fields: ctx, num
func (_m *BandwidthKeeper) AddMsgSignedByApp(ctx types.Context, num int64) types.Error {
	ret := _m.Called(ctx, num)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, int64) types.Error); ok {
		r0 = rf(ctx, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// AddMsgSignedByUser provides a mock function with given fields: ctx, num
func (_m *BandwidthKeeper) AddMsgSignedByUser(ctx types.Context, num int64) types.Error {
	ret := _m.Called(ctx, num)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, int64) types.Error); ok {
		r0 = rf(ctx, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// CalculateCurMsgFee provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) CalculateCurMsgFee(ctx types.Context) types.Error {
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

// ClearBlockInfo provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) ClearBlockInfo(ctx types.Context) types.Error {
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

// ConsumeBandwidthCredit provides a mock function with given fields: ctx, costPerMsg, accKey
func (_m *BandwidthKeeper) ConsumeBandwidthCredit(ctx types.Context, costPerMsg types.Dec, accKey linotypes.AccountKey) types.Error {
	ret := _m.Called(ctx, costPerMsg, accKey)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, types.Dec, linotypes.AccountKey) types.Error); ok {
		r0 = rf(ctx, costPerMsg, accKey)
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

// GetAllAppInfo provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) GetAllAppInfo(ctx types.Context) ([]*model.AppBandwidthInfo, types.Error) {
	ret := _m.Called(ctx)

	var r0 []*model.AppBandwidthInfo
	if rf, ok := ret.Get(0).(func(types.Context) []*model.AppBandwidthInfo); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AppBandwidthInfo)
		}
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context) types.Error); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetBandwidthCostPerMsg provides a mock function with given fields: ctx, u, p
func (_m *BandwidthKeeper) GetBandwidthCostPerMsg(ctx types.Context, u types.Dec, p types.Dec) types.Dec {
	ret := _m.Called(ctx, u, p)

	var r0 types.Dec
	if rf, ok := ret.Get(0).(func(types.Context, types.Dec, types.Dec) types.Dec); ok {
		r0 = rf(ctx, u, p)
	} else {
		r0 = ret.Get(0).(types.Dec)
	}

	return r0
}

// GetPunishmentCoeff provides a mock function with given fields: ctx, accKey
func (_m *BandwidthKeeper) GetPunishmentCoeff(ctx types.Context, accKey linotypes.AccountKey) (types.Dec, types.Error) {
	ret := _m.Called(ctx, accKey)

	var r0 types.Dec
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) types.Dec); ok {
		r0 = rf(ctx, accKey)
	} else {
		r0 = ret.Get(0).(types.Dec)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, accKey)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetVacancyCoeff provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) GetVacancyCoeff(ctx types.Context) (types.Dec, types.Error) {
	ret := _m.Called(ctx)

	var r0 types.Dec
	if rf, ok := ret.Get(0).(func(types.Context) types.Dec); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(types.Dec)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context) types.Error); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
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

// IsUserMsgFeeEnough provides a mock function with given fields: ctx, fee
func (_m *BandwidthKeeper) IsUserMsgFeeEnough(ctx types.Context, fee auth.StdFee) bool {
	ret := _m.Called(ctx, fee)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context, auth.StdFee) bool); ok {
		r0 = rf(ctx, fee)
	} else {
		r0 = ret.Get(0).(bool)
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

// RefillAppBandwidthCredit provides a mock function with given fields: ctx, accKey
func (_m *BandwidthKeeper) RefillAppBandwidthCredit(ctx types.Context, accKey linotypes.AccountKey) types.Error {
	ret := _m.Called(ctx, accKey)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r0 = rf(ctx, accKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// UpdateMaxMPSAndEMA provides a mock function with given fields: ctx
func (_m *BandwidthKeeper) UpdateMaxMPSAndEMA(ctx types.Context) types.Error {
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
