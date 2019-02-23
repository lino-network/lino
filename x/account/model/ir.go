package model

import (
	"github.com/lino-network/lino/types"
)

// PendingCoinDayQueueIR - TotalCoinDay: rat -> string
type PendingCoinDayQueueIR struct {
	LastUpdatedAt   int64            `json:"last_updated_at"`
	TotalCoinDay    string           `json:"total_coin_day"`
	TotalCoin       types.Coin       `json:"total_coin"`
	PendingCoinDays []PendingCoinDay `json:"pending_coin_days"`
}

// AccountRowIR account related information when migrate, pk: Username
type AccountRowIR struct {
	Username            types.AccountKey      `json:"username"`
	Info                AccountInfo           `json:"info"`
	Bank                AccountBank           `json:"bank"`
	Meta                AccountMeta           `json:"meta"`
	PendingCoinDayQueue PendingCoinDayQueueIR `json:"pending_coin_day_queue"`
}

// AccountTablesIR -
type AccountTablesIR struct {
	Accounts            []AccountRowIR   `json:"accounts"`
	AccountGrantPubKeys []GrantPubKeyRow `json:"account_grant_pub_keys"`
}
