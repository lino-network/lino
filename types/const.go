package types

const (
	// Decimals - Total decimals in Lino Blockchain
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
	UnknownPermission          = Permission(0)
	AppPermission              = Permission(1)
	TransactionPermission      = Permission(2)
	ResetPermission            = Permission(3)
	GrantAppPermission         = Permission(4)
	PreAuthorizationPermission = Permission(5)

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

	// punishment type
	UnknownPunish      = PunishType(0)
	PunishByzantine    = PunishType(1)
	PunishAbsentCommit = PunishType(2)
	PunishDidntVote    = PunishType(3)

	// UsernameReCheck - UsernameReCheck is used to check user registration
	UsernameReCheck        = "^[a-z]([a-z0-9-\\.]){1,19}[a-z0-9]$"
	IllegalUsernameReCheck = "^[a-z0-9\\.-]*([-\\.]){2,}[a-z0-9\\.-]*$"

	// MinimumUsernameLength - minimum username length
	MinimumUsernameLength = 3

	// MaximumUsernameLength - maximum username length
	MaximumUsernameLength = 20

	// MaximumMemoLength - maximum length of memo
	MaximumMemoLength = 100

	// MaximumJSONMetaLength - maximum length of account JSON meta
	MaximumJSONMetaLength = 500

	// MaxPostTitleLength - maximum length of post title
	MaxPostTitleLength = 50

	// MaxPostContentLength - maximum length of post content
	MaxPostContentLength = 1000

	// KeySeparator - separate different key component
	KeySeparator = "/"

	// HoursPerYear - as defined by a julian year of 365.25 days
	HoursPerYear = 8766

	// MinutesPerYear - as defined by a julian year of 365.25 days
	MinutesPerYear = HoursPerYear * 60

	// MinutesPerMonth - as defined by a julian year of 365.25 days
	MinutesPerMonth = MinutesPerYear / 12

	// PrecisionFactor - all decimals will around to allow at most 7 decimals
	PrecisionFactor = 10000000

	// NewRatFromDecimalPrecision - precision used in sdk NewRatFromDecimal
	NewRatFromDecimalPrecision = 5

	// MaximumSdkRatLength - maximum length of sdk.Rat can pass into blockchain
	MaximumSdkRatLength = 10

	// MaximumLinkIdentifier - maximum length of Links identifier
	MaximumLinkIdentifier = 20

	// MaximumLinkURL - maximum length of Links URL
	MaximumLinkURL = 100

	// MaximumLengthOfPostID - maximum length of post ID
	MaximumLengthOfPostID = 50

	// MaximumNumOfLinks - maximum number of links per post
	MaximumNumOfLinks = 10

	// MaximumLengthOfDeveloperWebsite - maximum length of developer website
	MaximumLengthOfDeveloperWebsite = 100

	// MaximumLengthOfDeveloperDesctiption - maximum length of developer description
	MaximumLengthOfDeveloperDesctiption = 1000

	// MaximumLengthOfAppMetadata - maximum length of developer App meta data
	MaximumLengthOfAppMetadata = 1000

	// MaximumLengthOfProposalReason - maximum length of proposal reason
	MaximumLengthOfProposalReason = 1000

	// InitAccountWithFullStakeMemo - init account with full stake memo
	InitAccountWithFullStakeMemo = "open account deposit"

	// InitAccountRegisterDepositMemo - init account deposit fee memo
	InitAccountRegisterDepositMemo = "init deposit"

	// PermlinkSeparator - permlink separator
	PermlinkSeparator = "#"

	// BalanceHistoryBundleSize - bundle size for balance history
	BalanceHistoryBundleSize = 100

	// RewardHistoryBundleSize - bundle size for reward history
	RewardHistoryBundleSize = 100

	// CoinDayRecordIntervalSec - coin day record in the same interval bucket will be merged
	CoinDayRecordIntervalSec = 1200
)
