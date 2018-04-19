package types

const (
	// Total decimals in Lino Blockchain
	Decimals = 100000

	// AccountKVStoreKey presents store which keeps account related value
	AccountKVStoreKey = "account"
	// PostKVStoreKey presents store which keeps post related value
	PostKVStoreKey = "post"
	// ValidatorKVStoreKey presents store which keeps validator related value
	ValidatorKVStoreKey = "validator"
	// EventKVStoreKey presents store which keeps event related value
	GlobalKVStoreKey = "global"
	// VoteKVStoreKey presents store which keeps vote related value
	VoteKVStoreKey = "vote"
	// InfraKVStoreKey presents store which keeps infra related value
	InfraKVStoreKey = "infra"
	// DeveloperKVStoreKey presents store which keeps developer related value
	DeveloperKVStoreKey = "developer"

	// RegisterRouterName is used for routing in app
	RegisterRouterName = "register"
	// AccountRouterName is used for routing in app
	AccountRouterName = "account"
	// PostRouterName is used for routing in app
	PostRouterName = "post"
	// ValidatorRouterName is used for routing in app
	ValidatorRouterName = "validator"
	// VoteRounterName is used for routing in app
	VoteRouterName = "vote"
	// InfraRounterName is used for routing in app
	InfraRouterName = "infra"
	// DeveloperRounterName is used for routing in app
	DeveloperRouterName = "developer"

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

	// KeySeparator used to separate different key component
	KeySeparator = "/"

	// as defined by a julian year of 365.25 days
	HoursPerYear = 8766
)
