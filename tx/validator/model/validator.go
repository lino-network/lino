package model

import (
	types "github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
)

// Validator is basic structure records all validator information
type Validator struct {
	ABCIValidator abci.Validator
	Username      types.AccountKey `json:"username"`
	Deposit       types.Coin       `json:"deposit"`
	AbsentVote    int              `json:"absent_vote"`
	IsByzantine   bool             `json:"is_byzantine"`
}

// Validator list
type ValidatorList struct {
	OncallValidators []types.AccountKey `json:"oncall_validators"`
	AllValidators    []types.AccountKey `json:"all_validators"`
	LowestPower      types.Coin         `json:"lowest_power"`
	LowestValidator  types.AccountKey   `json:"lowest_validator"`
}
