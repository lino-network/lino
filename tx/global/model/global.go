package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type GlobalMeta struct {
	TotalLinoCoin                 types.Coin `json:"total_lino_coin"`
	LastYearCumulativeConsumption types.Coin `json:"last_year_cumulative_consumption"`
	CumulativeConsumption         types.Coin `json:"cumulative_consumption"`
	GrowthRate                    sdk.Rat    `json:"growth_rate"`
	Ceiling                       sdk.Rat    `json:"ceiling"`
	Floor                         sdk.Rat    `json:"floor"`
}

type TPS struct {
	CurrentTPS sdk.Rat `json:"current_tps"`
	MaxTPS     sdk.Rat `json:"max_tps"`
}

// GlobalStatistics
type GlobalStatistics struct {
	numOfAccount  int64 `json:"number_of_account"`
	numOfPost     int64 `json:"number_of_post"`
	numOfComment  int64 `json:"number_of_comment"`
	numOfTransfer int64 `json:"number_of_transfer"`
	numOfLike     int64 `json:"number_of_like"`
	numOfDonation int64 `json:"number_of_donation"`
}

// InflationPool, determined by GlobalAllocation
// InfraInflationPool inflation pool for infra
// TotalContentCreatorInflationPool total inflation pool for content creator this year
// DistributedContentCreatorInflationPool inflation alrady distributed
// DeveloperInflationPool inflation pool for developer
// ValidatorInflationPool inflation pool for validator
type InflationPool struct {
	InfraInflationPool          types.Coin `json:"infra_inflation_pool"`
	ContentCreatorInflationPool types.Coin `json:"content_creator_inflation_pool"`
	DeveloperInflationPool      types.Coin `json:"developer_inflation_pool"`
	ValidatorInflationPool      types.Coin `json:"validator_inflation_pool"`
}

// ConsumptionMeta
// ConsumptionFrictionRate: percentage the user consumption deducted and added to the TotalLinoInflationPool
// ReportStakeWindow used to evaluate the panelty of the post within the freezing period
// DislikeStakeWindow used to evaluate the panelty of the post within the freezing period
// ConsumptionWindow records all content related consumption within the freezing period
// ConsumptionFreezingPeriodHr is the time content createor can get remain consumption after friction
type ConsumptionMeta struct {
	ConsumptionFrictionRate     sdk.Rat    `json:"consumption_friction_rate"`
	ReportStakeWindow           sdk.Rat    `json:"report_stake_window"`
	DislikeStakeWindow          sdk.Rat    `json:"dislike_stake_window"`
	ConsumptionWindow           types.Coin `json:"consumption_window"`
	ConsumptionRewardPool       types.Coin `json:"consumption_reward_pool"`
	ConsumptionFreezingPeriodHr int64      `json:"consumption_freezing_period"`
}

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

type InfraInternalAllocationParam struct {
	StorageAllocation sdk.Rat `json:"storage_allocation"`
	CDNAllocation     sdk.Rat `json:"CDN_allocation"`
}

// var VoterMinDeposit = NewCoin(1000 * Decimals)
// var VoterMinWithdraw = NewCoin(1 * Decimals)
//var ProposalDecideHr = int64(7 * 24)
//var ProposalRegisterFee = NewCoin(2000 * Decimals)
//var CoinReturnIntervalHr = int64(7 * 24)
//var CoinReturnTimes = int64(7)
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
	NextProposalID          int64      `json:"next_proposal_id"`
	TypeAProposalDecideHr   int64      `json:"type_a_proposal_decide_hr"`
	TypeAProposalMinDeposit types.Coin `json:"type_a_proposal_min_deposit"`
	TypeAProposalPassRatio  sdk.Rat    `json:"type_a_proposal_pass_ratio"`
	TypeBProposalDecideHr   int64      `json:"type_b_proposal_decide_hr"`
	TypeBProposalMinDeposit types.Coin `json:"type_b_proposal_min_deposit"`
	TypeBProposalPassRatio  sdk.Rat    `json:"type_b_proposal_pass_ratio"`
	TypeCProposalDecideHr   int64      `json:"type_c_proposal_decide_hr"`
	TypeCProposalMinDeposit types.Coin `json:"type_c_proposal_min_deposit"`
	TypeCProposalPassRatio  sdk.Rat    `json:"type_c_proposal_pass_ratio"`
}

type DeveloperParam struct {
	DeveloperMinDeposit           types.Coin `json:"developer_min_deposit"`
	DeveloperCoinReturnIntervalHr int64      `json:"developer_coin_return_interval"`
	DeveloperCoinReturnTimes      int64      `json:"developer_coin_return_times"`
}

type ValidatorParam struct {
	ValidatorMinWithdraw          types.Coin `json:"validator_min_withdraw"`
	ValidatorMinVotingDeposit     types.Coin `json:"validator_min_voting_deposit"`
	ValidatorMinCommitingDeposit  types.Coin `json:"validator_min_commiting_deposit"`
	ValidatorCoinReturnIntervalHr int64      `json:"validator_coin_return_interval"`
	ValidatorCoinReturnTimes      int64      `json:"validator_coin_return_times"`
	PenaltyMissVote               types.Coin `json:"penalty_miss_vote"`
	PenaltyMissCommit             types.Coin `json:"penalty_miss_commit"`
	PenaltyByzantine              types.Coin `json:"penalty_byzantine"`
}
