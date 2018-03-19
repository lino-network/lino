package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ABCI Response Codes
	// Base SDK reserves 0 ~ 99.
	// Coin errors reserve 100 ~ 199.
	// Lino authentication errors reserve 200 ~ 299.
	// Lino register handler errors reserve 300 ~ 309.
	CodeInvalidUsername   sdk.CodeType = 301
	CodeAccRegisterFailed sdk.CodeType = 302

	// Lino account handler errors reserve 310 ~ 399
	CodeAccountManagerFail sdk.CodeType = 310

	// Lino post handler errors reserve 400 ~ 499
	// CodePostMarshalError indicates error occurs during marshal
	CodePostMarshalError sdk.CodeType = 400
	// CodePostUnmarshalError indicates error occurs during unmarshal
	CodePostUnmarshalError sdk.CodeType = 401
	// CodePostNotFound indicates the post is not in store.
	CodePostNotFound sdk.CodeType = 402
	// CodePostCreateError occurs when create msg fails some precondition
	CodePostCreateError sdk.CodeType = 403

	// RegisterRouterName is used for routing in app
	RegisterRouterName = "register"

	// AccountRouterName is used for routing in app
	AccountRouterName = "account"

	// UsernameReCheck is used to check user registration
	UsernameReCheck = "^[a-zA-Z0-9]([a-zA-Z0-9_-]){2,20}$"

	// MinimumUsernameLength minimum username length
	MinimumUsernameLength = 3

	// MaximumUsernameLength maximum username length
	MaximumUsernameLength = 20

	// DefaultAcitivityBurden for user when account is registered
	DefaultActivityBurden = 100

	// MsgType is uesd to register App codec
	msgTypeRegister = 0x1
)
