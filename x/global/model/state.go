package model

import (
	"github.com/lino-network/lino/types"
)

// GlobalTimeEventTimeRow - events, pk: UnixTime
type GlobalTimeEventTimeRow struct {
	UnixTime      int64               `json:"unix_time"`
	TimeEventList types.TimeEventList `json:"time_event_list"`
}

// GlobalStakeStatDayRow - stake stats of a day, pk: day
type GlobalStakeStatDayRow struct {
	Day       int64         `json:"day"`
	StakeStat LinoStakeStat `json:"stake_stat"`
}

// GlobalMisc - a bunch of global variables with no pk, pk: none
type GlobalMisc struct {
	Meta            GlobalMeta      `json:"meta"`
	InflationPool   InflationPool   `json:"inflation_pool"`
	ConsumptionMeta ConsumptionMeta `json:"consumption_meta"`
	TPS             TPS             `json:"tps"`
	Time            GlobalTime      `json:"time"`
}

// ToIR -
func (g GlobalMisc) ToIR() GlobalMiscIR {
	return GlobalMiscIR{
		Meta:            g.Meta,
		InflationPool:   g.InflationPool,
		ConsumptionMeta: g.ConsumptionMeta.ToIR(),
		TPS:             g.TPS.ToIR(),
		Time:            g.Time,
	}
}

// GlobalTables state of global.
type GlobalTables struct {
	GlobalTimeEventLists []GlobalTimeEventTimeRow `json:"global_time_event_lists"`
	GlobalStakeStats     []GlobalStakeStatDayRow  `json:"global_stake_stats"`
	GlobalMisc           GlobalMisc               `json:"global_misc"`
}
