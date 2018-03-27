package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	abci "github.com/tendermint/abci/types"
)

// Validator is basic structure records all validator information
type Validator struct {
	ABCIValidator abci.Validator
	Username      acc.AccountKey `json:"username"`
	Votes         []Vote         `json:"votes"`
	Deposit       sdk.Coins      `json:"deposit"`
	AbsentVote    int            `json:"absent_vote"`
}

// Validator list
type ValidatorList struct {
	OncallValidators []acc.AccountKey `json:"oncall_validators"`
	AllValidators    []acc.AccountKey `json:"all_validators"`
	LowestPower      sdk.Coins        `json:"lowest_power"`
	LowestValidator  acc.AccountKey   `json:"lowest_validator"`
}

// User's vote
type Vote struct {
	Voter acc.AccountKey `json:"voter"`
	Power sdk.Coins      `json:"power"`
}
