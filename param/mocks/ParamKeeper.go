// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import param "github.com/lino-network/lino/param"
import types "github.com/cosmos/cosmos-sdk/types"

// ParamKeeper is an autogenerated mock type for the ParamKeeper type
type ParamKeeper struct {
	mock.Mock
}

// GetAccountParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetAccountParam(ctx types.Context) (*param.AccountParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.AccountParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.AccountParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.AccountParam)
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

// GetBandwidthParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetBandwidthParam(ctx types.Context) (*param.BandwidthParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.BandwidthParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.BandwidthParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.BandwidthParam)
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

// GetCoinDayParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetCoinDayParam(ctx types.Context) (*param.CoinDayParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.CoinDayParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.CoinDayParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.CoinDayParam)
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

// GetDeveloperParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetDeveloperParam(ctx types.Context) (*param.DeveloperParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.DeveloperParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.DeveloperParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.DeveloperParam)
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

// GetInfraInternalAllocationParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetInfraInternalAllocationParam(ctx types.Context) (*param.InfraInternalAllocationParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.InfraInternalAllocationParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.InfraInternalAllocationParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.InfraInternalAllocationParam)
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

// GetPostParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetPostParam(ctx types.Context) (*param.PostParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.PostParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.PostParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.PostParam)
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

// GetProposalParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetProposalParam(ctx types.Context) (*param.ProposalParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.ProposalParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.ProposalParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.ProposalParam)
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

// GetReputationParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetReputationParam(ctx types.Context) (*param.ReputationParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.ReputationParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.ReputationParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.ReputationParam)
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

// GetValidatorParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetValidatorParam(ctx types.Context) (*param.ValidatorParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.ValidatorParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.ValidatorParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.ValidatorParam)
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

// GetVoteParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetVoteParam(ctx types.Context) (*param.VoteParam, types.Error) {
	ret := _m.Called(ctx)

	var r0 *param.VoteParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.VoteParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.VoteParam)
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

// UpdateGlobalGrowthRate provides a mock function with given fields: ctx, growthRate
func (_m *ParamKeeper) UpdateGlobalGrowthRate(ctx types.Context, growthRate types.Dec) types.Error {
	ret := _m.Called(ctx, growthRate)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, types.Dec) types.Error); ok {
		r0 = rf(ctx, growthRate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}
