// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	param "github.com/lino-network/lino/param"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ParamKeeper is an autogenerated mock type for the ParamKeeper type
type ParamKeeper struct {
	mock.Mock
}

// GetAccountParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetAccountParam(ctx types.Context) *param.AccountParam {
	ret := _m.Called(ctx)

	var r0 *param.AccountParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.AccountParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.AccountParam)
		}
	}

	return r0
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

// GetGlobalAllocationParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetGlobalAllocationParam(ctx types.Context) *param.GlobalAllocationParam {
	ret := _m.Called(ctx)

	var r0 *param.GlobalAllocationParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.GlobalAllocationParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.GlobalAllocationParam)
		}
	}

	return r0
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

// GetPriceParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetPriceParam(ctx types.Context) *param.PriceParam {
	ret := _m.Called(ctx)

	var r0 *param.PriceParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.PriceParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.PriceParam)
		}
	}

	return r0
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
func (_m *ParamKeeper) GetReputationParam(ctx types.Context) *param.ReputationParam {
	ret := _m.Called(ctx)

	var r0 *param.ReputationParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.ReputationParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.ReputationParam)
		}
	}

	return r0
}

// GetValidatorParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetValidatorParam(ctx types.Context) *param.ValidatorParam {
	ret := _m.Called(ctx)

	var r0 *param.ValidatorParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.ValidatorParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.ValidatorParam)
		}
	}

	return r0
}

// GetVoteParam provides a mock function with given fields: ctx
func (_m *ParamKeeper) GetVoteParam(ctx types.Context) *param.VoteParam {
	ret := _m.Called(ctx)

	var r0 *param.VoteParam
	if rf, ok := ret.Get(0).(func(types.Context) *param.VoteParam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*param.VoteParam)
		}
	}

	return r0
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
