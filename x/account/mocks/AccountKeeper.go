// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/tendermint/crypto"

	linotypes "github.com/lino-network/lino/types"

	mock "github.com/stretchr/testify/mock"

	model "github.com/lino-network/lino/x/account/model"

	types "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper is an autogenerated mock type for the AccountKeeper type
type AccountKeeper struct {
	mock.Mock
}

// AddCoinToAddress provides a mock function with given fields: ctx, addr, coin
func (_m *AccountKeeper) AddCoinToAddress(ctx types.Context, addr types.AccAddress, coin linotypes.Coin) types.Error {
	ret := _m.Called(ctx, addr, coin)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, addr, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// AddCoinToUsername provides a mock function with given fields: ctx, username, coin
func (_m *AccountKeeper) AddCoinToUsername(ctx types.Context, username linotypes.AccountKey, coin linotypes.Coin) types.Error {
	ret := _m.Called(ctx, username, coin)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, username, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// AddFrozenMoney provides a mock function with given fields: ctx, username, amount, start, interval, times
func (_m *AccountKeeper) AddFrozenMoney(ctx types.Context, username linotypes.AccountKey, amount linotypes.Coin, start int64, interval int64, times int64) types.Error {
	ret := _m.Called(ctx, username, amount, start, interval, times)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.Coin, int64, int64, int64) types.Error); ok {
		r0 = rf(ctx, username, amount, start, interval, times)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// AuthorizePermission provides a mock function with given fields: ctx, me, grantTo, validityPeriod, grantLevel, amount
func (_m *AccountKeeper) AuthorizePermission(ctx types.Context, me linotypes.AccountKey, grantTo linotypes.AccountKey, validityPeriod int64, grantLevel linotypes.Permission, amount linotypes.Coin) types.Error {
	ret := _m.Called(ctx, me, grantTo, validityPeriod, grantLevel, amount)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, int64, linotypes.Permission, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, me, grantTo, validityPeriod, grantLevel, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// CheckSigningPubKeyOwner provides a mock function with given fields: ctx, me, signKey, permission, amount
func (_m *AccountKeeper) CheckSigningPubKeyOwner(ctx types.Context, me linotypes.AccountKey, signKey crypto.PubKey, permission linotypes.Permission, amount linotypes.Coin) (linotypes.AccountKey, types.Error) {
	ret := _m.Called(ctx, me, signKey, permission, amount)

	var r0 linotypes.AccountKey
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, crypto.PubKey, linotypes.Permission, linotypes.Coin) linotypes.AccountKey); ok {
		r0 = rf(ctx, me, signKey, permission, amount)
	} else {
		r0 = ret.Get(0).(linotypes.AccountKey)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey, crypto.PubKey, linotypes.Permission, linotypes.Coin) types.Error); ok {
		r1 = rf(ctx, me, signKey, permission, amount)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// CheckSigningPubKeyOwnerByAddress provides a mock function with given fields: ctx, addr, signkey
func (_m *AccountKeeper) CheckSigningPubKeyOwnerByAddress(ctx types.Context, addr types.AccAddress, signkey crypto.PubKey) types.Error {
	ret := _m.Called(ctx, addr, signkey)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, crypto.PubKey) types.Error); ok {
		r0 = rf(ctx, addr, signkey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// CreateAccount provides a mock function with given fields: ctx, username, signingKey, transactionKey
func (_m *AccountKeeper) CreateAccount(ctx types.Context, username linotypes.AccountKey, signingKey crypto.PubKey, transactionKey crypto.PubKey) types.Error {
	ret := _m.Called(ctx, username, signingKey, transactionKey)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, crypto.PubKey, crypto.PubKey) types.Error); ok {
		r0 = rf(ctx, username, signingKey, transactionKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// DoesAccountExist provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) DoesAccountExist(ctx types.Context, username linotypes.AccountKey) bool {
	ret := _m.Called(ctx, username)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) bool); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetAddress provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetAddress(ctx types.Context, username linotypes.AccountKey) (types.AccAddress, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 types.AccAddress
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) types.AccAddress); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.AccAddress)
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

// GetAllGrantPubKeys provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetAllGrantPubKeys(ctx types.Context, username linotypes.AccountKey) ([]*model.GrantPermission, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 []*model.GrantPermission
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) []*model.GrantPermission); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GrantPermission)
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

// GetBank provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetBank(ctx types.Context, username linotypes.AccountKey) (*model.AccountBank, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 *model.AccountBank
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) *model.AccountBank); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccountBank)
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

// GetBankByAddress provides a mock function with given fields: ctx, addr
func (_m *AccountKeeper) GetBankByAddress(ctx types.Context, addr types.AccAddress) (*model.AccountBank, types.Error) {
	ret := _m.Called(ctx, addr)

	var r0 *model.AccountBank
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) *model.AccountBank); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccountBank)
		}
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, types.AccAddress) types.Error); ok {
		r1 = rf(ctx, addr)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetFrozenMoneyList provides a mock function with given fields: ctx, addr
func (_m *AccountKeeper) GetFrozenMoneyList(ctx types.Context, addr types.Address) ([]model.FrozenMoney, types.Error) {
	ret := _m.Called(ctx, addr)

	var r0 []model.FrozenMoney
	if rf, ok := ret.Get(0).(func(types.Context, types.Address) []model.FrozenMoney); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FrozenMoney)
		}
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, types.Address) types.Error); ok {
		r1 = rf(ctx, addr)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetGrantPubKeys provides a mock function with given fields: ctx, username, grantTo
func (_m *AccountKeeper) GetGrantPubKeys(ctx types.Context, username linotypes.AccountKey, grantTo linotypes.AccountKey) ([]*model.GrantPermission, types.Error) {
	ret := _m.Called(ctx, username, grantTo)

	var r0 []*model.GrantPermission
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey) []*model.GrantPermission); ok {
		r0 = rf(ctx, username, grantTo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GrantPermission)
		}
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, username, grantTo)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetInfo provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetInfo(ctx types.Context, username linotypes.AccountKey) (*model.AccountInfo, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 *model.AccountInfo
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) *model.AccountInfo); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccountInfo)
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

// GetMeta provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetMeta(ctx types.Context, username linotypes.AccountKey) (*model.AccountMeta, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 *model.AccountMeta
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) *model.AccountMeta); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccountMeta)
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

// GetSavingFromUsername provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetSavingFromUsername(ctx types.Context, username linotypes.AccountKey) (linotypes.Coin, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 linotypes.Coin
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) linotypes.Coin); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(linotypes.Coin)
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

// GetSequence provides a mock function with given fields: ctx, address
func (_m *AccountKeeper) GetSequence(ctx types.Context, address types.Address) (uint64, types.Error) {
	ret := _m.Called(ctx, address)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(types.Context, types.Address) uint64); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, types.Address) types.Error); ok {
		r1 = rf(ctx, address)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetSigningKey provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetSigningKey(ctx types.Context, username linotypes.AccountKey) (crypto.PubKey, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 crypto.PubKey
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) crypto.PubKey); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crypto.PubKey)
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

// GetTransactionKey provides a mock function with given fields: ctx, username
func (_m *AccountKeeper) GetTransactionKey(ctx types.Context, username linotypes.AccountKey) (crypto.PubKey, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 crypto.PubKey
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) crypto.PubKey); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crypto.PubKey)
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

// ImportFromFile provides a mock function with given fields: ctx, cdc, filepath
func (_m *AccountKeeper) ImportFromFile(ctx types.Context, cdc *amino.Codec, filepath string) error {
	ret := _m.Called(ctx, cdc, filepath)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *amino.Codec, string) error); ok {
		r0 = rf(ctx, cdc, filepath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IncreaseSequenceByOne provides a mock function with given fields: ctx, address
func (_m *AccountKeeper) IncreaseSequenceByOne(ctx types.Context, address types.Address) types.Error {
	ret := _m.Called(ctx, address)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, types.Address) types.Error); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// MinusCoinFromAddress provides a mock function with given fields: ctx, addr, coin
func (_m *AccountKeeper) MinusCoinFromAddress(ctx types.Context, addr types.AccAddress, coin linotypes.Coin) types.Error {
	ret := _m.Called(ctx, addr, coin)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, addr, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// MinusCoinFromUsername provides a mock function with given fields: ctx, username, coin
func (_m *AccountKeeper) MinusCoinFromUsername(ctx types.Context, username linotypes.AccountKey, coin linotypes.Coin) types.Error {
	ret := _m.Called(ctx, username, coin)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, username, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// MoveCoin provides a mock function with given fields: ctx, sender, receiver, coin
func (_m *AccountKeeper) MoveCoin(ctx types.Context, sender linotypes.AccountKey, receiver linotypes.AccountKey, coin linotypes.Coin) types.Error {
	ret := _m.Called(ctx, sender, receiver, coin)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, sender, receiver, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// RegisterAccount provides a mock function with given fields: ctx, referrer, registerFee, username, signingKey, transactionKey
func (_m *AccountKeeper) RegisterAccount(ctx types.Context, referrer linotypes.AccountKey, registerFee linotypes.Coin, username linotypes.AccountKey, signingKey crypto.PubKey, transactionKey crypto.PubKey) types.Error {
	ret := _m.Called(ctx, referrer, registerFee, username, signingKey, transactionKey)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.Coin, linotypes.AccountKey, crypto.PubKey, crypto.PubKey) types.Error); ok {
		r0 = rf(ctx, referrer, registerFee, username, signingKey, transactionKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// RevokePermission provides a mock function with given fields: ctx, me, grantTo, permission
func (_m *AccountKeeper) RevokePermission(ctx types.Context, me linotypes.AccountKey, grantTo linotypes.AccountKey, permission linotypes.Permission) types.Error {
	ret := _m.Called(ctx, me, grantTo, permission)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, linotypes.Permission) types.Error); ok {
		r0 = rf(ctx, me, grantTo, permission)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// UpdateJSONMeta provides a mock function with given fields: ctx, username, JSONMeta
func (_m *AccountKeeper) UpdateJSONMeta(ctx types.Context, username linotypes.AccountKey, JSONMeta string) types.Error {
	ret := _m.Called(ctx, username, JSONMeta)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, string) types.Error); ok {
		r0 = rf(ctx, username, JSONMeta)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}
