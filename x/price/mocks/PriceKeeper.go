// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	linotypes "github.com/lino-network/lino/types"
	mock "github.com/stretchr/testify/mock"

	model "github.com/lino-network/lino/x/price/model"

	types "github.com/cosmos/cosmos-sdk/types"
)

// PriceKeeper is an autogenerated mock type for the PriceKeeper type
type PriceKeeper struct {
	mock.Mock
}

// CoinToMiniDollar provides a mock function with given fields: ctx, coin
func (_m *PriceKeeper) CoinToMiniDollar(ctx types.Context, coin linotypes.Coin) (linotypes.MiniDollar, types.Error) {
	ret := _m.Called(ctx, coin)

	var r0 linotypes.MiniDollar
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.Coin) linotypes.MiniDollar); ok {
		r0 = rf(ctx, coin)
	} else {
		r0 = ret.Get(0).(linotypes.MiniDollar)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.Coin) types.Error); ok {
		r1 = rf(ctx, coin)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// CurrPrice provides a mock function with given fields: ctx
func (_m *PriceKeeper) CurrPrice(ctx types.Context) (linotypes.MiniDollar, types.Error) {
	ret := _m.Called(ctx)

	var r0 linotypes.MiniDollar
	if rf, ok := ret.Get(0).(func(types.Context) linotypes.MiniDollar); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(linotypes.MiniDollar)
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

// FeedPrice provides a mock function with given fields: ctx, validator, _a2
func (_m *PriceKeeper) FeedPrice(ctx types.Context, validator linotypes.AccountKey, _a2 linotypes.MiniDollar) types.Error {
	ret := _m.Called(ctx, validator, _a2)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.MiniDollar) types.Error); ok {
		r0 = rf(ctx, validator, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// HistoryPrice provides a mock function with given fields: ctx
func (_m *PriceKeeper) HistoryPrice(ctx types.Context) []model.FeedHistory {
	ret := _m.Called(ctx)

	var r0 []model.FeedHistory
	if rf, ok := ret.Get(0).(func(types.Context) []model.FeedHistory); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FeedHistory)
		}
	}

	return r0
}

// InitGenesis provides a mock function with given fields: ctx, initPrice
func (_m *PriceKeeper) InitGenesis(ctx types.Context, initPrice linotypes.MiniDollar) types.Error {
	ret := _m.Called(ctx, initPrice)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.MiniDollar) types.Error); ok {
		r0 = rf(ctx, initPrice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// LastFeed provides a mock function with given fields: ctx, validator
func (_m *PriceKeeper) LastFeed(ctx types.Context, validator linotypes.AccountKey) (*model.FedPrice, types.Error) {
	ret := _m.Called(ctx, validator)

	var r0 *model.FedPrice
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) *model.FedPrice); ok {
		r0 = rf(ctx, validator)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FedPrice)
		}
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, validator)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// MiniDollarToCoin provides a mock function with given fields: ctx, dollar
func (_m *PriceKeeper) MiniDollarToCoin(ctx types.Context, dollar linotypes.MiniDollar) (linotypes.Coin, linotypes.MiniDollar, types.Error) {
	ret := _m.Called(ctx, dollar)

	var r0 linotypes.Coin
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.MiniDollar) linotypes.Coin); ok {
		r0 = rf(ctx, dollar)
	} else {
		r0 = ret.Get(0).(linotypes.Coin)
	}

	var r1 linotypes.MiniDollar
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.MiniDollar) linotypes.MiniDollar); ok {
		r1 = rf(ctx, dollar)
	} else {
		r1 = ret.Get(1).(linotypes.MiniDollar)
	}

	var r2 types.Error
	if rf, ok := ret.Get(2).(func(types.Context, linotypes.MiniDollar) types.Error); ok {
		r2 = rf(ctx, dollar)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(types.Error)
		}
	}

	return r0, r1, r2
}

// UpdatePrice provides a mock function with given fields: ctx
func (_m *PriceKeeper) UpdatePrice(ctx types.Context) types.Error {
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
