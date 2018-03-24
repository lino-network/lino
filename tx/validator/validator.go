package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
)

// Validator Account
type ValidatorAccount struct {
	validatorName acc.AccountKey `json:"validator_name"`
	votes         []Vote         `json:"votes"`
	totalWeight   int64          `json:"total_weight"`
	deposit       sdk.Coins      `json:"deposit"`
}

// Validator candidate list
type ValidatorList struct {
	validatorListKey acc.AccountKey   `json:"validator_list_key"`
	validators       []acc.AccountKey `json:"validators"`
	minWeight        int64            `json:"min_weight"`
}

// User's vote
type Vote struct {
	voter         acc.AccountKey `json:"voter"`
	weight        int64          `json:"weight"`
	validatorName acc.AccountKey `json:"validator_name"`
}
