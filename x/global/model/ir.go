package model

import (
	"github.com/lino-network/lino/types"
)

// ConsumptionMetaIR - ConsumptionFrictionRate rat -> float string
type ConsumptionMetaIR struct {
	ConsumptionFrictionRate string `json:"consumption_friction_rate"`
	// Although ConsumptionWindow's type is types.Coin, its actually type is specified by
	// IsConsumptionWindowDollarUnit, from upgrade2.
	ConsumptionWindow            types.Coin `json:"consumption_window"`
	ConsumptionRewardPool        types.Coin `json:"consumption_reward_pool"`
	ConsumptionFreezingPeriodSec int64      `json:"consumption_freezing_period_second"`
	// since upgrade-2
	IsConsumptionWindowDollarUnit bool `json:"is_consumption_window_dollar_unit"`
}

// TPSIR - all from rat to float string
type TPSIR struct {
	CurrentTPS string `json:"current_tps"`
	MaxTPS     string `json:"max_tps"`
}

// GlobalMiscIR - ConsumptionMeta changed.
type GlobalMiscIR struct {
	Meta            GlobalMeta        `json:"meta"`
	InflationPool   InflationPool     `json:"inflation_pool"`
	ConsumptionMeta ConsumptionMetaIR `json:"consumption_meta"`
	TPS             TPSIR             `json:"tps"`
	Time            GlobalTime        `json:"time"`
}

// GlobalTablesIR - GlobalMisc changed.
type GlobalTablesIR struct {
	GlobalTimeEventLists []GlobalTimeEventTimeRow `json:"global_time_event_lists"`
	GlobalStakeStats     []GlobalStakeStatDayRow  `json:"global_stake_stats"`
	GlobalMisc           GlobalMiscIR             `json:"global_misc"`
}
