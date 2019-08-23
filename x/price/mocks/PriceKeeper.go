// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

import types "github.com/lino-network/lino/types"

// PriceKeeper is an autogenerated mock type for the PriceKeeper type
type PriceKeeper struct {
	mock.Mock
}

// CoinToMiniDollar provides a mock function with given fields: coin
func (_m *PriceKeeper) CoinToMiniDollar(coin types.Coin) types.MiniDollar {
	ret := _m.Called(coin)

	var r0 types.MiniDollar
	if rf, ok := ret.Get(0).(func(types.Coin) types.MiniDollar); ok {
		r0 = rf(coin)
	} else {
		r0 = ret.Get(0).(types.MiniDollar)
	}

	return r0
}

// MiniDollarToCoin provides a mock function with given fields: dollar
func (_m *PriceKeeper) MiniDollarToCoin(dollar types.MiniDollar) types.Coin {
	ret := _m.Called(dollar)

	var r0 types.Coin
	if rf, ok := ret.Get(0).(func(types.MiniDollar) types.Coin); ok {
		r0 = rf(dollar)
	} else {
		r0 = ret.Get(0).(types.Coin)
	}

	return r0
}