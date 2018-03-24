package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	abci "github.com/tendermint/abci/types"
)

// Validator Account
type ValidatorAccount struct {
	abci.Validator
	validatorName acc.AccountKey `json:"validator_name"`
	votes         []Vote         `json:"votes"`
	deposit       sdk.Coins      `json:"deposit"`
}

// Validator list
type ValidatorList struct {
	validatorListKey acc.AccountKey   `json:"validator_list_key"`
	validators       []acc.AccountKey `json:"validators"`
	validatorPool    []acc.AccountKey `json:"validatorPool"`
	lowestPower      sdk.Coins        `json:"lowest_power"`
	lowestValidator  acc.AccountKey   `json:"lowest_validator"`
}

// User's vote
type Vote struct {
	voter         acc.AccountKey `json:"voter"`
	power         sdk.Coins      `json:"power"`
	validatorName acc.AccountKey `json:"validator_name"`
}
