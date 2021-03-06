package param

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

// Parameter - parameter in Lino Blockchain
type Parameter interface{}

// GlobalAllocationParam - global allocation parameters
// ContentCreatorAllocation - percentage for all content creator related allocation
// DeveloperAllocation - percentage of inflation for developers
// ValidatorAllocation - percentage of inflation for validators
type GlobalAllocationParam struct {
	GlobalGrowthRate         sdk.Dec `json:"global_growth_rate"`
	ContentCreatorAllocation sdk.Dec `json:"content_creator_allocation"`
	DeveloperAllocation      sdk.Dec `json:"developer_allocation"`
	ValidatorAllocation      sdk.Dec `json:"validator_allocation"`
}

func (gp GlobalAllocationParam) IsValid() bool {
	sum := sdk.NewDec(0)
	sum = sum.Add(gp.ContentCreatorAllocation)
	sum = sum.Add(gp.DeveloperAllocation)
	sum = sum.Add(gp.ValidatorAllocation)
	return sum.Equal(sdk.NewDec(1))
}

// VoteParam - vote parameters
// MinStakeIn - minimum stake for stake in msg
// VoterCoinReturnIntervalSec - when withdraw or revoke, the deposit return to voter by return event
// VoterCoinReturnTimes - when withdraw or revoke, the deposit return to voter by return event
type VoteParam struct {
	MinStakeIn                 types.Coin `json:"min_stake_in"`
	VoterCoinReturnIntervalSec int64      `json:"voter_coin_return_interval_second"`
	VoterCoinReturnTimes       int64      `json:"voter_coin_return_times"`
}

// ProposalParam - proposal parameters
// ContentCensorshipDecideSec - seconds after content censorship proposal created till expired
// ContentCensorshipMinDeposit - minimum deposit to propose content censorship proposal
// ContentCensorshipPassRatio - upvote and downvote ratio for content censorship proposal
// ContentCensorshipPassVotes - minimum voting power required to pass content censorship proposal
// ChangeParamDecideSec - seconds after parameter change proposal created till expired
// ChangeParamExecutionSec - seconds after parameter change proposal pass till execution
// ChangeParamMinDeposit - minimum deposit to propose parameter change proposal
// ChangeParamPassRatio - upvote and downvote ratio for parameter change proposal
// ChangeParamPassVotes - minimum voting power required to pass parameter change proposal
// ProtocolUpgradeDecideSec - seconds after protocol upgrade proposal created till expired
// ProtocolUpgradeMinDeposit - minimum deposit to propose protocol upgrade proposal
// ProtocolUpgradePassRatio - upvote and downvote ratio for protocol upgrade proposal
// ProtocolUpgradePassVotes - minimum voting power required to pass protocol upgrade proposal
type ProposalParam struct {
	ContentCensorshipDecideSec  int64      `json:"content_censorship_decide_second"`
	ContentCensorshipMinDeposit types.Coin `json:"content_censorship_min_deposit"`
	ContentCensorshipPassRatio  sdk.Dec    `json:"content_censorship_pass_ratio"`
	ContentCensorshipPassVotes  types.Coin `json:"content_censorship_pass_votes"`
	ChangeParamDecideSec        int64      `json:"change_param_decide_second"`
	ChangeParamExecutionSec     int64      `json:"change_param_execution_second"`
	ChangeParamMinDeposit       types.Coin `json:"change_param_min_deposit"`
	ChangeParamPassRatio        sdk.Dec    `json:"change_param_pass_ratio"`
	ChangeParamPassVotes        types.Coin `json:"change_param_pass_votes"`
	ProtocolUpgradeDecideSec    int64      `json:"protocol_upgrade_decide_second"`
	ProtocolUpgradeMinDeposit   types.Coin `json:"protocol_upgrade_min_deposit"`
	ProtocolUpgradePassRatio    sdk.Dec    `json:"protocol_upgrade_pass_ratio"`
	ProtocolUpgradePassVotes    types.Coin `json:"protocol_upgrade_pass_votes"`
}

// DeveloperParam - developer parameters
// DeveloperMinDeposit - minimum deposit to become a developer
// DeveloperCoinReturnIntervalSec - when withdraw or revoke, coin return to developer by coin return event
// DeveloperCoinReturnTimes - when withdraw or revoke, coin return to developer by coin return event
type DeveloperParam struct {
	DeveloperMinDeposit            types.Coin `json:"developer_min_deposit"`
	DeveloperCoinReturnIntervalSec int64      `json:"developer_coin_return_interval_second"`
	DeveloperCoinReturnTimes       int64      `json:"developer_coin_return_times"`
}

// ValidatorParam - validator parameters
// ValidatorMinDeposit - minimum deposit requirement for user wanna be validator
// ValidatorCoinReturnIntervalSec - when withdraw or revoke, coin return to validator by coin return event
// ValidatorCoinReturnTimes - when withdraw or revoke, coin return to validator by coin return event
// minus PenaltyMissCommit amount of Coin from validator deposit
// PenaltyMissCommit - when missing block till AbsentCommitLimitation, minus PenaltyMissCommit amount of Coin from validator deposit
// PenaltyByzantine - when validator acts as byzantine (double sign, for example),
// minus PenaltyByzantine amount of Coin from validator deposit
// AbsentCommitLimitation - absent block limitation till penalty
// OncallSize - the size of oncall validators
// StandbySize - the size of standby validators
// ValidatorRevokePendingSec - how many seconds before unassign validator duty
// OncallInflationWeight - oncall validator's weight when distributing inflation
// StandbyInflationWeight - standby validator's weight when distributing inflation
// MaxVotedValidators - the number of max validators one voter can vote
// SlashLimitation - slash limitation till into jail
type ValidatorParam struct {
	ValidatorMinDeposit            types.Coin `json:"validator_min_deposit"`
	ValidatorCoinReturnIntervalSec int64      `json:"validator_coin_return_second"`
	ValidatorCoinReturnTimes       int64      `json:"validator_coin_return_times"`
	PenaltyMissCommit              types.Coin `json:"penalty_miss_commit"`
	PenaltyByzantine               types.Coin `json:"penalty_byzantine"`
	AbsentCommitLimitation         int64      `json:"absent_commit_limitation"`
	OncallSize                     int64      `json:"oncall_size"`
	StandbySize                    int64      `json:"standby_size"`
	ValidatorRevokePendingSec      int64      `json:"validator_revoke_pending_sec"`
	OncallInflationWeight          int64      `json:"oncall_inflation_weight"`
	StandbyInflationWeight         int64      `json:"standby_inflation_weight"`
	MaxVotedValidators             int64      `json:"max_voted_validators"`
	SlashLimitation                int64      `json:"slash_limitation"`
}

// BandwidthParam - bandwidth parameters
// SecondsToRecoverBandwidth - seconds for user tps capacity fully charged
// CapacityUsagePerTransaction - capacity usage per transaction, dynamic changed based on traffic
// GeneralMsgQuotaRatio - the ratio for reserved general messages per second
// GeneralMsgEMAFactor - the multiplier for weighting the general message EMA
// AppMsgQuotaRatio - the ratio for reserved app messages per second
// AppMsgEMAFactor - the multiplier for weighting the app message EMA
// ExpectedMaxMPS - the expected max messages per second
// MsgFeeFactorA - factor A for calculating msg fee
// MsgFeeFactorB - factor B for calculating msg fee
// MaxMPSDecayRate - decay rate for historical max message per seconds
// AppBandwidthPoolSize - the depth for app bandwidth pool
// AppVacancyFactor - app vacancy factor for calculating u
// AppPunishmentFactor - app punishment factor for calculating p

type BandwidthParam struct {
	SecondsToRecoverBandwidth   int64      `json:"seconds_to_recover_bandwidth"`
	CapacityUsagePerTransaction types.Coin `json:"capacity_usage_per_transaction"`
	VirtualCoin                 types.Coin `json:"virtual_coin"`
	GeneralMsgQuotaRatio        sdk.Dec    `json:"general_msg_quota_ratio"`
	GeneralMsgEMAFactor         sdk.Dec    `json:"general_msg_ema_factor"`
	AppMsgQuotaRatio            sdk.Dec    `json:"app_msg_quota_ratio"`
	AppMsgEMAFactor             sdk.Dec    `json:"app_msg_ema_factor"`
	ExpectedMaxMPS              sdk.Dec    `json:"expected_max_mps"`
	MsgFeeFactorA               sdk.Dec    `json:"msg_fee_factor_a"`
	MsgFeeFactorB               sdk.Dec    `json:"msg_fee_factor_b"`
	MaxMPSDecayRate             sdk.Dec    `json:"max_mps_decay_rate"`
	AppBandwidthPoolSize        sdk.Dec    `json:"app_bandwidth_pool_size"`
	AppVacancyFactor            sdk.Dec    `json:"app_vacancy_factor"`
	AppPunishmentFactor         sdk.Dec    `json:"app_punishment_factor"`
}

// AccountParam - account parameters
// MinimumBalance - minimum balance each account need to maintain
// RegisterFee - register fee need to pay to developer inflation pool for each account registration
type AccountParam struct {
	MinimumBalance types.Coin `json:"minimum_balance"`
	RegisterFee    types.Coin `json:"register_fee"`
}

// PostParam - empty, reserved.
type PostParam struct {
}

// ReputationParam: parameters of reputation
// BestContentIndexN - hard cap of how many content can be indexed every round.
// UserMaxN - maximum different donation target counted every round.
type ReputationParam struct {
	BestContentIndexN int `json:"best_content_index_n"`
	UserMaxN          int `json:"user_max_n"`
}

// PriceParam - parameters of price module.
type PriceParam struct {
	TestnetMode     bool       `json:"testnet_mode"`
	UpdateEverySec  int64      `json:"update_every"`
	FeedEverySec    int64      `json:"feed_every"`
	HistoryMaxLen   int        `json:"history_max_len"`
	PenaltyMissFeed types.Coin `json:"penalty_miss_feed"`
}
