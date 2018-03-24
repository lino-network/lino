package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	abci "github.com/tendermint/abci/types"
)

// Validator Account
type ValidatorAccount struct {
	abci.Validator
	ValidatorName acc.AccountKey `json:"validator_name"`
	Votes         []Vote         `json:"votes"`
	Deposit       sdk.Coins      `json:"deposit"`
}

// Validator list
type ValidatorList struct {
	Validators      []acc.AccountKey `json:"validators"`
	ValidatorPool   []acc.AccountKey `json:"validatorPool"`
	LowestPower     sdk.Coins        `json:"lowest_power"`
	LowestValidator acc.AccountKey   `json:"lowest_validator"`
}

// User's vote
type Vote struct {
	voter         acc.AccountKey `json:"voter"`
	power         sdk.Coins      `json:"power"`
	validatorName acc.AccountKey `json:"validator_name"`
}
