package types

const (
	// Total decimals in Lino Blockchain
	Decimals = 100000

	// KVStoreKey presents store which used by app
	MainKVStoreKey      = "main"
	AccountKVStoreKey   = "account"
	PostKVStoreKey      = "post"
	ValidatorKVStoreKey = "validator"
	GlobalKVStoreKey    = "global"
	VoteKVStoreKey      = "vote"
	InfraKVStoreKey     = "infra"
	DeveloperKVStoreKey = "developer"
	ParamKVStoreKey     = "param"
	ProposalKVStoreKey  = "proposal"

	// RouterName for msg routing in app
	AccountRouterName   = "account"
	PostRouterName      = "post"
	ValidatorRouterName = "validator"
	VoteRouterName      = "vote"
	InfraRouterName     = "infra"
	DeveloperRouterName = "developer"
	ProposalRouterName  = "proposal"

	// Different permission level for msg
	UnknownPermission           = Permission(0)
	PostPermission              = Permission(1)
	MicropaymentPermission      = Permission(2)
	TransactionPermission       = Permission(3)
	MasterPermission            = Permission(4)
	GrantPostPermission         = Permission(5)
	GrantMicropaymentPermission = Permission(6)

	// Different proposal result
	ProposalNotPass = ProposalResult(0)
	ProposalPass    = ProposalResult(1)
	ProposalRevoked = ProposalResult(2)

	// Different proposal types
	ChangeParam       = ProposalType(0)
	ContentCensorship = ProposalType(1)
	ProtocolUpgrade   = ProposalType(2)

	// Different donation types
	DirectDeposit = DonationType(0)
	Inflation     = DonationType(1)

	// Different possible incomes
	TransferIn           = TransferDetailType(0)
	DonationIn           = TransferDetailType(1)
	ClaimReward          = TransferDetailType(2)
	ValidatorInflation   = TransferDetailType(3)
	DeveloperInflation   = TransferDetailType(4)
	InfraInflation       = TransferDetailType(5)
	VoteReturnCoin       = TransferDetailType(6)
	DelegationReturnCoin = TransferDetailType(7)
	ValidatorReturnCoin  = TransferDetailType(8)
	DeveloperReturnCoin  = TransferDetailType(9)
	InfraReturnCoin      = TransferDetailType(10)
	ProposalReturnCoin   = TransferDetailType(11)
	GenesisCoin          = TransferDetailType(12)

	// Different possible outcomes
	TransferOut      = TransferDetailType(13)
	DonationOut      = TransferDetailType(14)
	Delegate         = TransferDetailType(15)
	VoterDeposit     = TransferDetailType(16)
	ValidatorDeposit = TransferDetailType(17)
	DeveloperDeposit = TransferDetailType(18)
	InfraDeposit     = TransferDetailType(19)
	ProposalDeposit  = TransferDetailType(20)

	// UsernameReCheck is used to check user registration
	UsernameReCheck = "^[a-z0-9]([a-z0-9_-]){2,20}$"

	// MinimumUsernameLength minimum username length
	MinimumUsernameLength = 3

	// MaximumUsernameLength maximum username length
	MaximumUsernameLength = 20

	// MaximumMemoLength denotes the maximum length of memo
	MaximumMemoLength = 100

	// MaximumJSONMeta denotes the maximum length of account JSON meta
	MaximumJSONMetaLength = 500

	// DefaultAcitivityBurden for user when account is registered
	DefaultActivityBurden = 100

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

	// as defined by a julian year of 365.25 days
	MinutesPerYear = HoursPerYear * 60

	// as defined by a julian year of 365.25 days
	MinutesPerMonth = MinutesPerYear / 12

	// all decimals will around to allow at most 3 decimals
	PrecisionFactor = 1000

	// precision used in sdk NewRatFromDecimal
	NewRatFromDecimalPrecision = 5

	// Maximum length of sdk.Rat can pass into blockchain
	MaximumSdkRatLength = 10

	// Maximum length of Links identifier
	MaximumLinkIdentifier = 20

	// Maximum length of Links URL
	MaximumLinkURL = 50

	// Maximum length of post ID
	MaximumLengthOfPostID = 50

	// Maximum number of links per post
	MaximumNumOfLinks = 10

	// Init account with full stake memo
	InitAccountWithFullStakeMemo = "init register fee"

	// Init account deposit fee memo
	InitAccountRegisterDepositMemo = "init deposit"

	// Permlink separator
	PermlinkSeparator = "#"
)
