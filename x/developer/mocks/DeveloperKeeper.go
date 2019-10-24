// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	amino "github.com/tendermint/go-amino"

	linotypes "github.com/lino-network/lino/types"

	mock "github.com/stretchr/testify/mock"

	model "github.com/lino-network/lino/x/developer/model"

	types "github.com/cosmos/cosmos-sdk/types"
)

// DeveloperKeeper is an autogenerated mock type for the DeveloperKeeper type
type DeveloperKeeper struct {
	mock.Mock
}

// AppTransferIDA provides a mock function with given fields: ctx, appname, signer, amount, from, to
func (_m *DeveloperKeeper) AppTransferIDA(ctx types.Context, appname linotypes.AccountKey, signer linotypes.AccountKey, amount types.Int, from linotypes.AccountKey, to linotypes.AccountKey) types.Error {
	ret := _m.Called(ctx, appname, signer, amount, from, to)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, types.Int, linotypes.AccountKey, linotypes.AccountKey) types.Error); ok {
		r0 = rf(ctx, appname, signer, amount, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// BurnIDA provides a mock function with given fields: ctx, app, user, amount
func (_m *DeveloperKeeper) BurnIDA(ctx types.Context, app linotypes.AccountKey, user linotypes.AccountKey, amount linotypes.MiniDollar) (linotypes.Coin, types.Error) {
	ret := _m.Called(ctx, app, user, amount)

	var r0 linotypes.Coin
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, linotypes.MiniDollar) linotypes.Coin); ok {
		r0 = rf(ctx, app, user, amount)
	} else {
		r0 = ret.Get(0).(linotypes.Coin)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, linotypes.MiniDollar) types.Error); ok {
		r1 = rf(ctx, app, user, amount)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// DoesDeveloperExist provides a mock function with given fields: ctx, username
func (_m *DeveloperKeeper) DoesDeveloperExist(ctx types.Context, username linotypes.AccountKey) bool {
	ret := _m.Called(ctx, username)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) bool); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ExportToFile provides a mock function with given fields: ctx, cdc, filepath
func (_m *DeveloperKeeper) ExportToFile(ctx types.Context, cdc *amino.Codec, filepath string) error {
	ret := _m.Called(ctx, cdc, filepath)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *amino.Codec, string) error); ok {
		r0 = rf(ctx, cdc, filepath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAffiliated provides a mock function with given fields: ctx, app
func (_m *DeveloperKeeper) GetAffiliated(ctx types.Context, app linotypes.AccountKey) []linotypes.AccountKey {
	ret := _m.Called(ctx, app)

	var r0 []linotypes.AccountKey
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) []linotypes.AccountKey); ok {
		r0 = rf(ctx, app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]linotypes.AccountKey)
		}
	}

	return r0
}

// GetAffiliatingApp provides a mock function with given fields: ctx, username
func (_m *DeveloperKeeper) GetAffiliatingApp(ctx types.Context, username linotypes.AccountKey) (linotypes.AccountKey, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 linotypes.AccountKey
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) linotypes.AccountKey); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(linotypes.AccountKey)
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

// GetDeveloper provides a mock function with given fields: ctx, username
func (_m *DeveloperKeeper) GetDeveloper(ctx types.Context, username linotypes.AccountKey) (model.Developer, types.Error) {
	ret := _m.Called(ctx, username)

	var r0 model.Developer
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) model.Developer); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(model.Developer)
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

// GetIDA provides a mock function with given fields: ctx, app
func (_m *DeveloperKeeper) GetIDA(ctx types.Context, app linotypes.AccountKey) (model.AppIDA, types.Error) {
	ret := _m.Called(ctx, app)

	var r0 model.AppIDA
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) model.AppIDA); ok {
		r0 = rf(ctx, app)
	} else {
		r0 = ret.Get(0).(model.AppIDA)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, app)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetIDABank provides a mock function with given fields: ctx, app, user
func (_m *DeveloperKeeper) GetIDABank(ctx types.Context, app linotypes.AccountKey, user linotypes.AccountKey) (model.IDABank, types.Error) {
	ret := _m.Called(ctx, app, user)

	var r0 model.IDABank
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey) model.IDABank); ok {
		r0 = rf(ctx, app, user)
	} else {
		r0 = ret.Get(0).(model.IDABank)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, app, user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetIDAStats provides a mock function with given fields: ctx, app
func (_m *DeveloperKeeper) GetIDAStats(ctx types.Context, app linotypes.AccountKey) (model.AppIDAStats, types.Error) {
	ret := _m.Called(ctx, app)

	var r0 model.AppIDAStats
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) model.AppIDAStats); ok {
		r0 = rf(ctx, app)
	} else {
		r0 = ret.Get(0).(model.AppIDAStats)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, app)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetLiveDevelopers provides a mock function with given fields: ctx
func (_m *DeveloperKeeper) GetLiveDevelopers(ctx types.Context) []model.Developer {
	ret := _m.Called(ctx)

	var r0 []model.Developer
	if rf, ok := ret.Get(0).(func(types.Context) []model.Developer); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Developer)
		}
	}

	return r0
}

// GetMiniIDAPrice provides a mock function with given fields: ctx, app
func (_m *DeveloperKeeper) GetMiniIDAPrice(ctx types.Context, app linotypes.AccountKey) (linotypes.MiniDollar, types.Error) {
	ret := _m.Called(ctx, app)

	var r0 linotypes.MiniDollar
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey) linotypes.MiniDollar); ok {
		r0 = rf(ctx, app)
	} else {
		r0 = ret.Get(0).(linotypes.MiniDollar)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.AccountKey) types.Error); ok {
		r1 = rf(ctx, app)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// GetReservePool provides a mock function with given fields: ctx
func (_m *DeveloperKeeper) GetReservePool(ctx types.Context) model.ReservePool {
	ret := _m.Called(ctx)

	var r0 model.ReservePool
	if rf, ok := ret.Get(0).(func(types.Context) model.ReservePool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(model.ReservePool)
	}

	return r0
}

// GrantPermission provides a mock function with given fields: ctx, app, user, duration, level, amount
func (_m *DeveloperKeeper) GrantPermission(ctx types.Context, app linotypes.AccountKey, user linotypes.AccountKey, duration int64, level linotypes.Permission, amount string) types.Error {
	ret := _m.Called(ctx, app, user, duration, level, amount)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, int64, linotypes.Permission, string) types.Error); ok {
		r0 = rf(ctx, app, user, duration, level, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// ImportFromFile provides a mock function with given fields: ctx, cdc, filepath
func (_m *DeveloperKeeper) ImportFromFile(ctx types.Context, cdc *amino.Codec, filepath string) error {
	ret := _m.Called(ctx, cdc, filepath)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, *amino.Codec, string) error); ok {
		r0 = rf(ctx, cdc, filepath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitGenesis provides a mock function with given fields: ctx, reservePoolAmount
func (_m *DeveloperKeeper) InitGenesis(ctx types.Context, reservePoolAmount linotypes.Coin) types.Error {
	ret := _m.Called(ctx, reservePoolAmount)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, reservePoolAmount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// IssueIDA provides a mock function with given fields: ctx, appname, idaName, idaPrice
func (_m *DeveloperKeeper) IssueIDA(ctx types.Context, appname linotypes.AccountKey, idaName string, idaPrice int64) types.Error {
	ret := _m.Called(ctx, appname, idaName, idaPrice)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, string, int64) types.Error); ok {
		r0 = rf(ctx, appname, idaName, idaPrice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// MintIDA provides a mock function with given fields: ctx, appname, amount
func (_m *DeveloperKeeper) MintIDA(ctx types.Context, appname linotypes.AccountKey, amount linotypes.Coin) types.Error {
	ret := _m.Called(ctx, appname, amount)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.Coin) types.Error); ok {
		r0 = rf(ctx, appname, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// MonthlyDistributeDevInflation provides a mock function with given fields: ctx
func (_m *DeveloperKeeper) MonthlyDistributeDevInflation(ctx types.Context) types.Error {
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

// MoveIDA provides a mock function with given fields: ctx, app, from, to, amount
func (_m *DeveloperKeeper) MoveIDA(ctx types.Context, app linotypes.AccountKey, from linotypes.AccountKey, to linotypes.AccountKey, amount linotypes.MiniDollar) types.Error {
	ret := _m.Called(ctx, app, from, to, amount)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, linotypes.AccountKey, linotypes.MiniDollar) types.Error); ok {
		r0 = rf(ctx, app, from, to, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// RegisterDeveloper provides a mock function with given fields: ctx, username, website, description, appMetaData
func (_m *DeveloperKeeper) RegisterDeveloper(ctx types.Context, username linotypes.AccountKey, website string, description string, appMetaData string) types.Error {
	ret := _m.Called(ctx, username, website, description, appMetaData)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, string, string, string) types.Error); ok {
		r0 = rf(ctx, username, website, description, appMetaData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// ReportConsumption provides a mock function with given fields: ctx, username, consumption
func (_m *DeveloperKeeper) ReportConsumption(ctx types.Context, username linotypes.AccountKey, consumption linotypes.MiniDollar) types.Error {
	ret := _m.Called(ctx, username, consumption)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.MiniDollar) types.Error); ok {
		r0 = rf(ctx, username, consumption)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// RevokePermission provides a mock function with given fields: ctx, user, app, perm
func (_m *DeveloperKeeper) RevokePermission(ctx types.Context, user linotypes.AccountKey, app linotypes.AccountKey, perm linotypes.Permission) types.Error {
	ret := _m.Called(ctx, user, app, perm)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, linotypes.Permission) types.Error); ok {
		r0 = rf(ctx, user, app, perm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// UpdateAffiliated provides a mock function with given fields: ctx, appname, username, activate
func (_m *DeveloperKeeper) UpdateAffiliated(ctx types.Context, appname linotypes.AccountKey, username linotypes.AccountKey, activate bool) types.Error {
	ret := _m.Called(ctx, appname, username, activate)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, bool) types.Error); ok {
		r0 = rf(ctx, appname, username, activate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// UpdateDeveloper provides a mock function with given fields: ctx, username, website, description, appMetadata
func (_m *DeveloperKeeper) UpdateDeveloper(ctx types.Context, username linotypes.AccountKey, website string, description string, appMetadata string) types.Error {
	ret := _m.Called(ctx, username, website, description, appMetadata)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, string, string, string) types.Error); ok {
		r0 = rf(ctx, username, website, description, appMetadata)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// UpdateIDAAuth provides a mock function with given fields: ctx, app, username, active
func (_m *DeveloperKeeper) UpdateIDAAuth(ctx types.Context, app linotypes.AccountKey, username linotypes.AccountKey, active bool) types.Error {
	ret := _m.Called(ctx, app, username, active)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.AccountKey, bool) types.Error); ok {
		r0 = rf(ctx, app, username, active)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}
