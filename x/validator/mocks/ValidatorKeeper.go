// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"

	linotypes "github.com/lino-network/lino/types"

	manager "github.com/lino-network/lino/x/validator/manager"

	mock "github.com/stretchr/testify/mock"

	model "github.com/lino-network/lino/x/validator/model"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorKeeper is an autogenerated mock type for the ValidatorKeeper type
type ValidatorKeeper struct {
	mock.Mock
}

// DistributeInflationToValidator provides a mock function with given fields: ctx
func (_m *ValidatorKeeper) DistributeInflationToValidator(ctx types.Context) types.Error {
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

// GetCommittingValidatorVoteStatus provides a mock function with given fields: ctx
func (_m *ValidatorKeeper) GetCommittingValidatorVoteStatus(ctx types.Context) []model.ReceivedVotesStatus {
	ret := _m.Called(ctx)

	var r0 []model.ReceivedVotesStatus
	if rf, ok := ret.Get(0).(func(types.Context) []model.ReceivedVotesStatus); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ReceivedVotesStatus)
		}
	}

	return r0
}

// GetCommittingValidators provides a mock function with given fields: ctx
func (_m *ValidatorKeeper) GetCommittingValidators(ctx types.Context) []linotypes.AccountKey {
	ret := _m.Called(ctx)

	var r0 []linotypes.AccountKey
	if rf, ok := ret.Get(0).(func(types.Context) []linotypes.AccountKey); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]linotypes.AccountKey)
		}
	}

	return r0
}

// GetElectionVoteList provides a mock function with given fields: ctx, accKey
func (_m *ValidatorKeeper) GetElectionVoteList(ctx types.Context, accKey linotypes.AccountKey) *model.ElectionVoteList {
	ret := _m.Called(ctx, accKey)

	var r0 *model.ElectionVoteList
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) *model.ElectionVoteList); ok {
		r0 = rf(ctx, accKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ElectionVoteList)
		}
	}

	return r0
}

// GetInitValidators provides a mock function with given fields: ctx
func (_m *ValidatorKeeper) GetInitValidators(ctx types.Context) ([]abcitypes.ValidatorUpdate, types.Error) {
	ret := _m.Called(ctx)

	var r0 []abcitypes.ValidatorUpdate
	if rf, ok := ret.Get(0).(func(types.Context) []abcitypes.ValidatorUpdate); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]abcitypes.ValidatorUpdate)
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

// GetValidator provides a mock function with given fields: ctx, username
func (_m *ValidatorKeeper) GetValidator(ctx types.Context, username linotypes.AccountKey) (*model.Validator, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 *model.Validator
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) *model.Validator); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Validator)
		}
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, username)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetValidatorList provides a mock function with given fields: ctx
func (_m *ValidatorKeeper) GetValidatorList(ctx types.Context) *model.ValidatorList {
	ret := _m.Called(ctx)

	var r0 *model.ValidatorList
	if rf, ok := ret.Get(0).(func(types.Context) *model.ValidatorList); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ValidatorList)
		}
	}

	return r0
}

// GetValidatorUpdates provides a mock function with given fields: ctx
func (_m *ValidatorKeeper) GetValidatorUpdates(ctx types.Context) ([]abcitypes.ValidatorUpdate, types.Error) {
	ret := _m.Called(ctx)

	var r0 []abcitypes.ValidatorUpdate
	if rf, ok := ret.Get(0).(func(types.Context) []abcitypes.ValidatorUpdate); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]abcitypes.ValidatorUpdate)
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

// Hooks provides a mock function with given fields:
func (_m *ValidatorKeeper) Hooks() manager.Hooks {
	ret := _m.Called()

	var r0 manager.Hooks
	if rf, ok := ret.Get(0).(func() manager.Hooks); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(manager.Hooks)
	}

	return r0
}

// InitGenesis provides a mock function with given fields: ctx
func (_m *ValidatorKeeper) InitGenesis(ctx types.Context) {
	_m.Called(ctx)
}

// OnBeginBlock provides a mock function with given fields: ctx, req
func (_m *ValidatorKeeper) OnBeginBlock(ctx types.Context, req abcitypes.RequestBeginBlock) {
	_m.Called(ctx, req)
}

// RegisterValidator provides a mock function with given fields: ctx, username, valPubKey, link
func (_m *ValidatorKeeper) RegisterValidator(ctx types.Context, username linotypes.AccountKey, valPubKey crypto.PubKey, link string) types.Error {
	ret := _m.Called(ctx, username, valPubKey, link)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, crypto.PubKey, string) types.Error); ok {
		r0 = rf(ctx, username, valPubKey, link)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// RevokeValidator provides a mock function with given fields: ctx, username
func (_m *ValidatorKeeper) RevokeValidator(ctx types.Context, username linotypes.AccountKey) types.Error {
	ret := _m.Called(ctx, username)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// VoteValidator provides a mock function with given fields: ctx, username, votedValidators
func (_m *ValidatorKeeper) VoteValidator(ctx types.Context, username linotypes.AccountKey, votedValidators []linotypes.AccountKey) types.Error {
	ret := _m.Called(ctx, username, votedValidators)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, []linotypes.AccountKey) types.Error); ok {
		r0 = rf(ctx, username, votedValidators)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}
