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

// VoterTablesIR - state of voter
type VoterTablesIR struct {
	Version int       `json:"version"`
	Voters  []VoterIR `json:"voters"`
}
