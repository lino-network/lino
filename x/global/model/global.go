package model

import (
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
