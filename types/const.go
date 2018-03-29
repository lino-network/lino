package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Denom = "lino"
	// ABCI Response Codes
	// Base SDK reserves 0 ~ 99.
	// Coin errors reserve 100 ~ 199.
	// Lino authentication errors reserve 200 ~ 299.
	// Lino register handler errors reserve 300 ~ 309.
	CodeInvalidUsername   sdk.CodeType = 301
	CodeAccRegisterFailed sdk.CodeType = 302
	CodeUsernameNotFound  sdk.CodeType = 303

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
	// CodePostLikeError occurs when like msg fails
	CodePostLikeError sdk.CodeType = 404
	// CodePostDonateError occurs when donate msg fails
	CodePostDonateError sdk.CodeType = 405

	// validator errors reserve 500 ~ 599
	CodeValidatorHandlerFailed sdk.CodeType = 500
	CodeValidatorManagerFailed sdk.CodeType = 501

	// Event errors reserve 600 ~ 699
	CodeEventExecuteError sdk.CodeType = 600

	// AccountKVStoreKey presents store which keeps account related value
	AccountKVStoreKey = "account"
	// PostKVStoreKey presents store which keeps post related value
	PostKVStoreKey = "post"
	// ValidatorKVStoreKey presents store which keeps validator related value
	ValidatorKVStoreKey = "validator"

	// RegisterRouterName is used for routing in app
	RegisterRouterName = "register"

	// AccountRouterName is used for routing in app
	AccountRouterName = "account"

	// PostRouterName is used for routing in app
	PostRouterName = "post"

	// ValidatorRouterName is used for routing in app
	ValidatorRouterName = "validator"

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

	// MinimumUsernameLength minimum username length
	MaxPostTitleLength = 50

	// MaximumUsernameLength maximum username length
	MaxPostContentLength = 1000

	// MaxLikeWeight indicates the 100.00% maximum like weight.
	MaxLikeWeight = 10000

	// MinLikeWeight indicates the -100.00% maximum like weight.
	MinLikeWeight = -10000

	KeySeparator = "/"
)
