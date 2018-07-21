package param

import (
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Parameter interface{}

// parameters can be changed by proposal
type EvaluateOfContentValueParam struct {
	ConsumptionTimeAdjustBase      int64   `json:"consumption_time_adjust_base"`
	ConsumptionTimeAdjustOffset    int64   `json:"consumption_time_adjust_offset"`
	NumOfConsumptionOnAuthorOffset int64   `json:"num_of_consumption_on_author_offset"`
	TotalAmountOfConsumptionBase   int64   `json:"total_amount_of_consumption_base"`
	TotalAmountOfConsumptionOffset int64   `json:"total_amount_of_consumption_offset"`
	AmountOfConsumptionExponent    sdk.Rat `json:"amount_of_consumption_exponent"`
}

// GlobalAllocation
// TotalLinoInflationPool: total Lino inflation for all community roles
// InfraAllocation percentage for all infra related allocation
// ContentCreatorAllocation percentage for all content creator related allocation
// DeveloperAllocation percentage of inflation for developers
// ValidatorAllocation percentage of inflation for validators
type GlobalAllocationParam struct {
	InfraAllocation          sdk.Rat `json:"infra_allocation"`
	ContentCreatorAllocation sdk.Rat `json:"content_creator_allocation"`
	DeveloperAllocation      sdk.Rat `json:"developer_allocation"`
	ValidatorAllocation      sdk.Rat `json:"validator_allocation"`
}

// TODO: year period

type InfraInternalAllocationParam struct {
	StorageAllocation sdk.Rat `json:"storage_allocation"`
	CDNAllocation     sdk.Rat `json:"CDN_allocation"`
}

type VoteParam struct {
	VoterMinDeposit               types.Coin `json:"voter_min_deposit"`
	VoterMinWithdraw              types.Coin `json:"voter_min_withdraw"`
	DelegatorMinWithdraw          types.Coin `json:"delegator_min_withdraw"`
	VoterCoinReturnIntervalHr     int64      `json:"voter_coin_return_interval"`
	VoterCoinReturnTimes          int64      `json:"voter_coin_return_times"`
	DelegatorCoinReturnIntervalHr int64      `json:"delegator_coin_return_interval"`
	DelegatorCoinReturnTimes      int64      `json:"delegator_coin_return_times"`
}

type ProposalParam struct {
	ContentCensorshipDecideHr   int64      `json:"content_censorship_decide_hr"`
	ContentCensorshipMinDeposit types.Coin `json:"content_censorship_min_deposit"`
	ContentCensorshipPassRatio  sdk.Rat    `json:"content_censorship_pass_ratio"`
	ContentCensorshipPassVotes  types.Coin `json:"content_censorship_pass_votes"`
	ChangeParamDecideHr         int64      `json:"change_param_decide_hr"`
	ChangeParamMinDeposit       types.Coin `json:"change_param_min_deposit"`
	ChangeParamPassRatio        sdk.Rat    `json:"change_param_pass_ratio"`
	ChangeParamPassVotes        types.Coin `json:"change_param_pass_votes"`
	ProtocolUpgradeDecideHr     int64      `json:"protocol_upgrade_decide_hr"`
	ProtocolUpgradeMinDeposit   types.Coin `json:"protocol_upgrade_min_deposit"`
	ProtocolUpgradePassRatio    sdk.Rat    `json:"protocol_upgrade_pass_ratio"`
	ProtocolUpgradePassVotes    types.Coin `json:"protocol_upgrade_pass_votes"`
}

type DeveloperParam struct {
	DeveloperMinDeposit           types.Coin `json:"developer_min_deposit"`
	DeveloperCoinReturnIntervalHr int64      `json:"developer_coin_return_interval"`
	DeveloperCoinReturnTimes      int64      `json:"developer_coin_return_times"`
}

// TODO: number of validators 20 -> 41
type ValidatorParam struct {
	ValidatorMinWithdraw          types.Coin `json:"validator_min_withdraw"`
	ValidatorMinVotingDeposit     types.Coin `json:"validator_min_voting_deposit"`
	ValidatorMinCommitingDeposit  types.Coin `json:"validator_min_commiting_deposit"`
	ValidatorCoinReturnIntervalHr int64      `json:"validator_coin_return_interval"`
	ValidatorCoinReturnTimes      int64      `json:"validator_coin_return_times"`
	PenaltyMissVote               types.Coin `json:"penalty_miss_vote"`
	PenaltyMissCommit             types.Coin `json:"penalty_miss_commit"`
	PenaltyByzantine              types.Coin `json:"penalty_byzantine"`
	ValidatorListSize             int64      `json:"validator_list_size"`
	AbsentCommitLimitation        int64      `json:"absent_commit_limitation"`
}

type CoinDayParam struct {
	DaysToRecoverCoinDayStake    int64 `json:"days_to_recover_coin_day_stake"`
	SecondsToRecoverCoinDayStake int64 `json:"seconds_to_recover_coin_day_stake"`
}

type BandwidthParam struct {
	SecondsToRecoverBandwidth   int64      `json:"seconds_to_recover_bandwidth"`
	CapacityUsagePerTransaction types.Coin `json:"capacity_usage_per_transaction"`
}

// AccountParam includes params related to account
type AccountParam struct {
	MinimumBalance           types.Coin `json:"minimum_balance"`
	RegisterFee              types.Coin `json:"register_fee"`
	BalanceHistoryBundleSize int64      `json:"balance_history_bundle_size"`
	RewardHistoryBundleSize  int64      `json:"reward_history_bundle_size"`
}

type PostParam struct {
	ReportOrUpvoteInterval int64 `json:"report_or_upvote_interval"`
}
