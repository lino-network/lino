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
