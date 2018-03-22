package validator

import (
	acc "github.com/lino-network/lino/tx/account"
)

// ValidatorKey key format in KVStore
type ValidatorKey string

// Validator Account
type ValidatorAccount struct {
	validatorID ValidatorKey `json:"key"`
	votes       []Vote       `json:"votes"`
	totalWeight int64        `json:"total_weight"`
	deposit     int64        `json:"deposit"`
}

// Validator candidate list
type ValidatorList struct {
	validatorListKey string             `json:"validator_list_key"`
	validators       []ValidatorAccount `json:"validators"`
}

// User's vote
type Vote struct {
	username acc.AccountKey `json:"username"`
	weight   int64          `json:"weight"`
}
