package param

import (
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Parameter - parameter in Lino Blockchain
type Parameter interface{}

// EvaluateOfContentValueParam - parameters used to evaluate content value
type EvaluateOfContentValueParam struct {
	ConsumptionTimeAdjustBase      int64 `json:"consumption_time_adjust_base"`
	ConsumptionTimeAdjustOffset    int64 `json:"consumption_time_adjust_offset"`
	NumOfConsumptionOnAuthorOffset int64 `json:"num_of_consumption_on_author_offset"`
	TotalAmountOfConsumptionBase   int64 `json:"total_amount_of_consumption_base"`
	TotalAmountOfConsumptionOffset int64 `json:"total_amount_of_consumption_offset"`
}

// GlobalAllocationParam - global allocation parameters
// InfraAllocation - percentage for all infra related allocation
// ContentCreatorAllocation - percentage for all content creator related allocation
// DeveloperAllocation - percentage of inflation for developers
// ValidatorAllocation - percentage of inflation for validators
type GlobalAllocationParam struct {
	GlobalGrowthRate         sdk.Dec `json:"global_growth_rate"`
	InfraAllocation          sdk.Dec `json:"infra_allocation"`
	ContentCreatorAllocation sdk.Dec `json:"content_creator_allocation"`
	DeveloperAllocation      sdk.Dec `json:"developer_allocation"`
	ValidatorAllocation      sdk.Dec `json:"validator_allocation"`
}

// InfraInternalAllocationParam - infra internal allocation parameters
// StorageAllocation - percentage for storage provider (not in use now)
// CDNAllocation - percentage for CDN provider (not in use now)
type InfraInternalAllocationParam struct {
	StorageAllocation sdk.Dec `json:"storage_allocation"`
	CDNAllocation     sdk.Dec `json:"CDN_allocation"`
}

// VoteParam - vote paramters
// MinStakeIn - minimum stake for stake in msg
// VoterCoinReturnIntervalSec - when withdraw or revoke, the deposit return to voter by return event
// VoterCoinReturnTimes - when withdraw or revoke, the deposit return to voter by return event
// DelegatorCoinReturnIntervalSec - when withdraw or revoke, the deposit return to delegator by return event
// DelegatorCoinReturnTimes - when withdraw or revoke, the deposit return to delegator by return event
type VoteParam struct {
	MinStakeIn                     types.Coin `json:"min_stake_in"`
	VoterCoinReturnIntervalSec     int64      `json:"voter_coin_return_interval_second"`
	VoterCoinReturnTimes           int64      `json:"voter_coin_return_times"`
	DelegatorCoinReturnIntervalSec int64      `json:"delegator_coin_return_interval_second"`
	DelegatorCoinReturnTimes       int64      `json:"delegator_coin_return_times"`
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
// ValidatorMinWithdraw - minimum withdraw requirement
// ValidatorMinVotingDeposit - minimum voting deposit requirement for user wanna be validator
// ValidatorMinCommittingDeposit - minimum committing (validator) deposit requirement for user wanna be validator
// ValidatorCoinReturnIntervalSec - when withdraw or revoke, coin return to validator by coin return event
// ValidatorCoinReturnTimes - when withdraw or revoke, coin return to validator by coin return event
// PenaltyMissVote - when missing vote for content censorship or protocol upgrade proposal,
// minus PenaltyMissCommit amount of Coin from validator deposit
// PenaltyMissCommit - when missing block till AbsentCommitLimitation, minus PenaltyMissCommit amount of Coin from validator deposit
// PenaltyByzantine - when validator acts as byzantine (double sign, for example),
// minus PenaltyByzantine amount of Coin from validator deposit
// ValidatorListSize - size of oncall validator
// AbsentCommitLimitation - absent block limitation till penalty
type ValidatorParam struct {
	ValidatorMinWithdraw           types.Coin `json:"validator_min_withdraw"`
	ValidatorMinVotingDeposit      types.Coin `json:"validator_min_voting_deposit"`
	ValidatorMinCommittingDeposit  types.Coin `json:"validator_min_committing_deposit"`
	ValidatorCoinReturnIntervalSec int64      `json:"validator_coin_return_second"`
	ValidatorCoinReturnTimes       int64      `json:"validator_coin_return_times"`
	PenaltyMissVote                types.Coin `json:"penalty_miss_vote"`
	PenaltyMissCommit              types.Coin `json:"penalty_miss_commit"`
	PenaltyByzantine               types.Coin `json:"penalty_byzantine"`
	ValidatorListSize              int64      `json:"validator_list_size"`
	AbsentCommitLimitation         int64      `json:"absent_commit_limitation"`
}

// CoinDayParam - coin day parameters
// SecondsToRecoverCoinDay - seconds for each incoming balance coin day fully charged
type CoinDayParam struct {
	SecondsToRecoverCoinDay int64 `json:"seconds_to_recover_coin_day"`
}

// BandwidthParam - bandwidth parameters
// SecondsToRecoverBandwidth - seconds for user tps capacity fully charged
// CapacityUsagePerTransaction - capacity usage per transaction, dynamic changed based on traffic
type BandwidthParam struct {
	SecondsToRecoverBandwidth   int64      `json:"seconds_to_recover_bandwidth"`
	CapacityUsagePerTransaction types.Coin `json:"capacity_usage_per_transaction"`
	VirtualCoin                 types.Coin `json:"virtual_coin"`
}

// AccountParam - account parameters
// MinimumBalance - minimum balance each account need to maintain
// RegisterFee - register fee need to pay to developer inflation pool for each account registration
// FirstDepositFullCoinDayLimit - when register account, some of coin day of register fee to newly open account will be fully charged
// MaxNumFrozenMoney - the upper limit for each person's ongoing frozen money
type AccountParam struct {
	MinimumBalance               types.Coin `json:"minimum_balance"`
	RegisterFee                  types.Coin `json:"register_fee"`
	FirstDepositFullCoinDayLimit types.Coin `json:"first_deposit_full_coin_day_limit"`
	MaxNumFrozenMoney            int64      `json:"max_num_frozen_money"`
}

// PostParam - post parameters
// ReportOrUpvoteIntervalSec - report interval second
// PostIntervalSec - post interval second
type PostParam struct {
	ReportOrUpvoteIntervalSec int64      `json:"report_or_upvote_interval_second"`
	PostIntervalSec           int64      `json:"post_interval_sec"`
	MaxReportReputation       types.Coin `json:"max_report_reputation"`
}

// BestContentIndexN - hard cap of how many content can be indexed every round.
type ReputationParam struct {
	BestContentIndexN int `json:"best_content_index_n"`
}
