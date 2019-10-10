package model

import (
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/types"
)

// VoterIR - pk: username
type VoterIR struct {
	Username          linotypes.AccountKey `json:"username"`
	LinoStake         linotypes.Coin       `json:"lino_stake"`
	LastPowerChangeAt int64                `json:"last_power_change_at"`
	Interest          linotypes.Coin       `json:"interest"`
	Duty              types.VoterDuty      `json:"duty"`
	FrozenAmount      linotypes.Coin       `json:"frozen_amount"`
}

// LinoStakeStatIR - records the information needed by
// lino power deposit, update and store daily.
type LinoStakeStatIR struct {
	TotalConsumptionFriction linotypes.Coin `json:"total_consumption_friction"`
	UnclaimedFriction        linotypes.Coin `json:"unclaimed_friction"`
	TotalLinoStake           linotypes.Coin `json:"total_lino_power"`
	UnclaimedLinoStake       linotypes.Coin `json:"unclaimed_lino_power"`
}

// StakeStatDayIR - stake stats of a day, pk: day
type StakeStatDayIR struct {
	Day       int64           `json:"day"`
	StakeStat LinoStakeStatIR `json:"stake_stat"`
}

// VoterTablesIR - state of voter
type VoterTablesIR struct {
	Version    int              `json:"version"`
	Voters     []VoterIR        `json:"voters"`
	StakeStats []StakeStatDayIR `json:"stake_stats"`
}
