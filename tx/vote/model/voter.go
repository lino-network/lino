package model

import (
	types "github.com/lino-network/lino/types"
)

type Voter struct {
	Username       types.AccountKey `json:"username"`
	Deposit        types.Coin       `json:"deposit"`
	DelegatedPower types.Coin       `json:"delegated_power"`
}

type Vote struct {
	Voter  types.AccountKey `json:"voter"`
	Result bool             `json:"result"`
}

type Delegation struct {
	Delegator types.AccountKey `json:"delegator"`
	Amount    types.Coin       `json:"amount"`
}

type ValidatorReferenceList struct {
	PenaltyValidators []types.AccountKey `json:"penalty_validators"`
	AllValidators     []types.AccountKey `json:"all_validators"`
	OncallValidators  []types.AccountKey `json:"oncall_validators"`
}
