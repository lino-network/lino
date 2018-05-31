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

	// Msg key
	IsRegister      = "is_register"
	PermissionLevel = "permission_level"

	// Different permission level for msg
	PostPermission        = Permission(0)
	TransactionPermission = Permission(1)
	MasterPermission      = Permission(2)

	// Different proposal result
	ProposalPass    = ProposalResult(0)
	ProposalNotPass = ProposalResult(1)
	ProposalRevoked = ProposalResult(2)

	// Different proposal types
	ChangeParam       = ProposalType(0)
	ContentCensorship = ProposalType(1)
	ProtocolUpgrade   = ProposalType(2)

	// Different donation types
	DirectDeposit = DonationType(0)
	Inflation     = DonationType(1)

	// Different possible incomes
	TransferIn           = BalanceHistoryDetailType(0)
	DonationIn           = BalanceHistoryDetailType(1)
	ClaimReward          = BalanceHistoryDetailType(2)
	ValidatorInflation   = BalanceHistoryDetailType(3)
	DeveloperInflation   = BalanceHistoryDetailType(4)
	InfraInflation       = BalanceHistoryDetailType(5)
	VoteReturnCoin       = BalanceHistoryDetailType(6)
	DelegationReturnCoin = BalanceHistoryDetailType(7)
	ValidatorReturnCoin  = BalanceHistoryDetailType(8)
	DeveloperReturnCoin  = BalanceHistoryDetailType(9)
	InfraReturnCoin      = BalanceHistoryDetailType(10)
	ProposalReturnCoin   = BalanceHistoryDetailType(11)
	GenesisCoin          = BalanceHistoryDetailType(12)

	// Different possible outcomes
	TransferOut      = BalanceHistoryDetailType(13)
	DonationOut      = BalanceHistoryDetailType(14)
	Delegate         = BalanceHistoryDetailType(15)
	VoterDeposit     = BalanceHistoryDetailType(16)
	ValidatorDeposit = BalanceHistoryDetailType(17)
	DeveloperDeposit = BalanceHistoryDetailType(18)
	InfraDeposit     = BalanceHistoryDetailType(19)
	ProposalDeposit  = BalanceHistoryDetailType(20)

	// UsernameReCheck is used to check user registration
	UsernameReCheck = "^[a-zA-Z0-9]([a-zA-Z0-9_-]){2,20}$"

	// MinimumUsernameLength minimum username length
	MinimumUsernameLength = 3

	// MaximumUsernameLength maximum username length
	MaximumUsernameLength = 20

	// MaximumMemoLength denotes the maximum length of memo
	MaximumMemoLength = 100

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

	PrecisionFactor = 1000
)
