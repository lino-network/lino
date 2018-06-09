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
	TransferIn           = TransferInDetail(0)
	DonationIn           = TransferInDetail(1)
	ClaimReward          = TransferInDetail(2)
	ValidatorInflation   = TransferInDetail(3)
	DeveloperInflation   = TransferInDetail(4)
	InfraInflation       = TransferInDetail(5)
	VoteReturnCoin       = TransferInDetail(6)
	DelegationReturnCoin = TransferInDetail(7)
	ValidatorReturnCoin  = TransferInDetail(8)
	DeveloperReturnCoin  = TransferInDetail(9)
	InfraReturnCoin      = TransferInDetail(10)
	ProposalReturnCoin   = TransferInDetail(11)
	GenesisCoin          = TransferInDetail(12)

	// Different possible outcomes
	TransferOut      = TransferOutDetail(0)
	DonationOut      = TransferOutDetail(1)
	Delegate         = TransferOutDetail(2)
	VoterDeposit     = TransferOutDetail(3)
	ValidatorDeposit = TransferOutDetail(4)
	DeveloperDeposit = TransferOutDetail(5)
	InfraDeposit     = TransferOutDetail(6)
	ProposalDeposit  = TransferOutDetail(7)

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

	PrecisionFactor = 1000

	// Detail from sources
	FromRewardPool         = InternalObject("reward pool")
	FromCoinReturnEvent    = InternalObject("coin return event")
	FromValidatorInflation = InternalObject("validator inflation")
	FromInfraInflation     = InternalObject("infra inflation")
	FromDeveloperInflation = InternalObject("developer inflation")

	// Detail to target
	ToDeveloperDeposit = InternalObject("developer deposit")
	ToVoterDeposit     = InternalObject("voter deposit")
	ToValidatorDeposit = InternalObject("validator deposit")
	ToProposalDeposit  = InternalObject("proposal deposit")
)
