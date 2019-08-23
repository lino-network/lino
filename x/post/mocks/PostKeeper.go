// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import linotypes "github.com/lino-network/lino/types"
import mock "github.com/stretchr/testify/mock"
import model "github.com/lino-network/lino/x/post/model"

import types "github.com/cosmos/cosmos-sdk/types"

// PostKeeper is an autogenerated mock type for the PostKeeper type
type PostKeeper struct {
	mock.Mock
}

// CreatePost provides a mock function with given fields: ctx, author, postID, createdBy, content, title
func (_m *PostKeeper) CreatePost(ctx types.Context, author linotypes.AccountKey, postID string, createdBy linotypes.AccountKey, content string, title string) types.Error {
	ret := _m.Called(ctx, author, postID, createdBy, content, title)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, string, linotypes.AccountKey, string, string) types.Error); ok {
		r0 = rf(ctx, author, postID, createdBy, content, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// DeletePost provides a mock function with given fields: ctx, permlink
func (_m *PostKeeper) DeletePost(ctx types.Context, permlink linotypes.Permlink) types.Error {
	ret := _m.Called(ctx, permlink)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.Permlink) types.Error); ok {
		r0 = rf(ctx, permlink)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// DoesPostExist provides a mock function with given fields: ctx, permlink
func (_m *PostKeeper) DoesPostExist(ctx types.Context, permlink linotypes.Permlink) bool {
	ret := _m.Called(ctx, permlink)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.Permlink) bool); ok {
		r0 = rf(ctx, permlink)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ExportToFile provides a mock function with given fields: ctx, filepath
func (_m *PostKeeper) ExportToFile(ctx types.Context, filepath string) error {
	ret := _m.Called(ctx, filepath)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string) error); ok {
		r0 = rf(ctx, filepath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPost provides a mock function with given fields: ctx, permlink
func (_m *PostKeeper) GetPost(ctx types.Context, permlink linotypes.Permlink) (model.Post, types.Error) {
	ret := _m.Called(ctx, permlink)

	var r0 model.Post
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.Permlink) model.Post); ok {
		r0 = rf(ctx, permlink)
	} else {
		r0 = ret.Get(0).(model.Post)
	}

	var r1 types.Error
	if rf, ok := ret.Get(1).(func(types.Context, linotypes.Permlink) types.Error); ok {
		r1 = rf(ctx, permlink)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(types.Error)
		}
	}

	return r0, r1
}

// IDADonate provides a mock function with given fields: ctx, from, n, author, postID, app
func (_m *PostKeeper) IDADonate(ctx types.Context, from linotypes.AccountKey, n types.Int, author linotypes.AccountKey, postID string, app linotypes.AccountKey) types.Error {
	ret := _m.Called(ctx, from, n, author, postID, app)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, types.Int, linotypes.AccountKey, string, linotypes.AccountKey) types.Error); ok {
		r0 = rf(ctx, from, n, author, postID, app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// ImportFromFile provides a mock function with given fields: ctx, filepath
func (_m *PostKeeper) ImportFromFile(ctx types.Context, filepath string) error {
	ret := _m.Called(ctx, filepath)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string) error); ok {
		r0 = rf(ctx, filepath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LinoDonate provides a mock function with given fields: ctx, from, amount, author, postID, app
func (_m *PostKeeper) LinoDonate(ctx types.Context, from linotypes.AccountKey, amount linotypes.Coin, author linotypes.AccountKey, postID string, app linotypes.AccountKey) types.Error {
	ret := _m.Called(ctx, from, amount, author, postID, app)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, linotypes.Coin, linotypes.AccountKey, string, linotypes.AccountKey) types.Error); ok {
		r0 = rf(ctx, from, amount, author, postID, app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}

// UpdatePost provides a mock function with given fields: ctx, author, postID, title, content
func (_m *PostKeeper) UpdatePost(ctx types.Context, author linotypes.AccountKey, postID string, title string, content string) types.Error {
	ret := _m.Called(ctx, author, postID, title, content)

	var r0 types.Error
	if rf, ok := ret.Get(0).(func(types.Context, linotypes.AccountKey, string, string, string) types.Error); ok {
		r0 = rf(ctx, author, postID, title, content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Error)
		}
	}

	return r0
}