package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

// GlobalMetaIR - global statistic information
type GlobalMetaIR struct {
	TotalLinoCoin         types.Coin `json:"total_lino_coin"`
	LastYearTotalLinoCoin types.Coin `json:"last_year_total_lino_coin"`
}

// InflationPoolIR -
type InflationPoolIR struct {
	InfraInflationPool     types.Coin `json:"infra_inflation_pool"`
	DeveloperInflationPool types.Coin `json:"developer_inflation_pool"`
	ValidatorInflationPool types.Coin `json:"validator_inflation_pool"`
}

// ConsumptionMetaIR - ConsumptionFrictionRate rat -> float string
type ConsumptionMetaIR struct {
	ConsumptionFrictionRate      sdk.Dec          `json:"consumption_friction_rate"`
	ConsumptionWindow            types.MiniDollar `json:"consumption_window"`
	ConsumptionRewardPool        types.Coin       `json:"consumption_reward_pool"`
	ConsumptionFreezingPeriodSec int64            `json:"consumption_freezing_period_second"`
}

// TPSIR - all from rat to float string
type TPSIR struct {
	CurrentTPS sdk.Dec `json:"current_tps"`
	MaxTPS     sdk.Dec `json:"max_tps"`
}

// GlobalTimeIR - global time
type GlobalTimeIR struct {
	ChainStartTime int64 `json:"chain_start_time"`
	LastBlockTime  int64 `json:"last_block_time"`
	PastMinutes    int64 `json:"past_minutes"`
}

// LinoStakeStatIR - records the information needed by
// lino power deposit, update and store daily.
type LinoStakeStatIR struct {
	TotalConsumptionFriction types.Coin `json:"total_consumption_friction"`
	UnclaimedFriction        types.Coin `json:"unclaimed_friction"`
	TotalLinoStake           types.Coin `json:"total_lino_power"`
	UnclaimedLinoStake       types.Coin `json:"unclaimed_lino_power"`
}

// GlobalStakeStatDayIR - stake stats of a day, pk: day
type GlobalStakeStatDayIR struct {
	Day       int64           `json:"day"`
	StakeStat LinoStakeStatIR `json:"stake_stat"`
}

// GlobalTimeEventsIR - events, pk: UnixTime
type GlobalTimeEventsIR struct {
	UnixTime      int64               `json:"unix_time"`
	TimeEventList types.TimeEventList `json:"time_event_list"`
}

// GlobalTablesIR - state
type GlobalTablesIR struct {
	Version              int                    `json:"version"`
	GlobalTimeEventLists []GlobalTimeEventsIR   `json:"global_time_event_lists"`
	GlobalStakeStats     []GlobalStakeStatDayIR `json:"global_stake_stats"`
	Meta                 GlobalMetaIR           `json:"meta"`
	InflationPool        InflationPoolIR        `json:"inflation_pool"`
	ConsumptionMeta      ConsumptionMetaIR      `json:"consumption_meta"`
	TPS                  TPSIR                  `json:"tps"`
	Time                 GlobalTimeIR           `json:"time"`
}
