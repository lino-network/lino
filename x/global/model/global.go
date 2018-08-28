package model

import (
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GlobalMeta - global statistic information
type GlobalMeta struct {
	TotalLinoCoin                 types.Coin `json:"total_lino_coin"`
	LastYearTotalLinoCoin         types.Coin `json:"last_year_total_lino_coin"`
	LastYearCumulativeConsumption types.Coin `json:"last_year_cumulative_consumption"`
	CumulativeConsumption         types.Coin `json:"cumulative_consumption"`
}

// GlobalTime - global time
type GlobalTime struct {
	ChainStartTime int64 `json:"chain_start_time"`
	LastBlockTime  int64 `json:"last_block_time"`
	PastMinutes    int64 `json:"past_minutes"`
}

// TPS - transaction per section
type TPS struct {
	CurrentTPS sdk.Rat `json:"current_tps"`
	MaxTPS     sdk.Rat `json:"max_tps"`
}

// InflationPool, determined by GlobalAllocation
// InfraInflationPool inflation pool for infra
// TotalContentCreatorInflationPool total inflation pool for content creator this year
// DistributedContentCreatorInflationPool inflation alrady distributed
// DeveloperInflationPool inflation pool for developer
// ValidatorInflationPool inflation pool for validator
type InflationPool struct {
	InfraInflationPool     types.Coin `json:"infra_inflation_pool"`
	DeveloperInflationPool types.Coin `json:"developer_inflation_pool"`
	ValidatorInflationPool types.Coin `json:"validator_inflation_pool"`
}

// ConsumptionMeta
// ConsumptionFrictionRate: percentage the user consumption deducted and added to the TotalLinoInflationPool
// ConsumptionWindow records all content related consumption within the freezing period
// ConsumptionFreezingPeriodHr is the time content createor can get remain consumption after friction
type ConsumptionMeta struct {
	ConsumptionFrictionRate      sdk.Rat    `json:"consumption_friction_rate"`
	ConsumptionWindow            types.Coin `json:"consumption_window"`
	ConsumptionRewardPool        types.Coin `json:"consumption_reward_pool"`
	ConsumptionFreezingPeriodSec int64      `json:"consumption_freezing_period_second"`
}

// InitParamList - genesis parameters
type InitParamList struct {
	MaxTPS                       sdk.Rat `json:"max_tps"`
	ConsumptionFrictionRate      sdk.Rat `json:"consumption_friction_rate"`
	ConsumptionFreezingPeriodSec int64   `json:"consumption_freezing_period_second"`
}
