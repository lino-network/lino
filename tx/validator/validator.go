package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
)

// Validator is basic structure records all validator information
type Validator struct {
	ABCIValidator       abci.Validator
	Username            acc.AccountKey `json:"username"`
	Deposit             sdk.Coins      `json:"deposit"`
	AbsentVote          int            `json:"absent_vote"`
	WithdrawAvailableAt types.Height   `json:"withdraw_available_at"`
	IsByzantine         bool           `json:"is_byzantine"`
}

// Validator list
type ValidatorList struct {
	OncallValidators []acc.AccountKey `json:"oncall_validators"`
	AllValidators    []acc.AccountKey `json:"all_validators"`
	LowestPower      sdk.Coins        `json:"lowest_power"`
	LowestValidator  acc.AccountKey   `json:"lowest_validator"`
}
