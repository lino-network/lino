package model

import (
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/types"
)

// Voter - a voter in blockchain is account with voter deposit, who can vote for a proposal
type Voter struct {
	Username          linotypes.AccountKey `json:"username"`
	LinoStake         linotypes.Coin       `json:"lino_stake"`
	LastPowerChangeAt int64                `json:"last_power_change_at"`
	Interest          linotypes.Coin       `json:"interest"`
	Duty              types.VoterDuty      `json:"duty"`
	FrozenAmount      linotypes.Coin       `json:"frozen_amount"`
}

// LinoStakeStat - records the information needed by
// lino power deposit, update and store daily.
type LinoStakeStat struct {
	TotalConsumptionFriction linotypes.Coin `json:"total_consumption_friction"`
	UnclaimedFriction        linotypes.Coin `json:"unclaimed_friction"`
	TotalLinoStake           linotypes.Coin `json:"total_lino_power"`
	UnclaimedLinoStake       linotypes.Coin `json:"unclaimed_lino_power"`
}
