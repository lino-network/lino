package model

import (
	crypto "github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

// AccountRow account related information when migrate, pk: Username
type AccountRow struct {
	Username            types.AccountKey    `json:"username"`
	Info                AccountInfo         `json:"info"`
	Bank                AccountBank         `json:"bank"`
	Meta                AccountMeta         `json:"meta"`
	Reward              Reward              `json:"reward"`
	PendingCoinDayQueue PendingCoinDayQueue `json:"pending_coin_day_queue"`
}

// ToIR -
func (a AccountRow) ToIR() AccountRowIR {
	return AccountRowIR{
		Username:            a.Username,
		Info:                a.Info,
		Bank:                a.Bank,
		Meta:                a.Meta,
		Reward:              a.Reward,
		PendingCoinDayQueue: a.PendingCoinDayQueue.ToIR(),
	}
}

// GrantPubKeyRow also in account, pk: (Username, pubKey)
type GrantPubKeyRow struct {
	Username    types.AccountKey `json:"username"`
	PubKey      crypto.PubKey    `json:"pub_key"`
	GrantPubKey GrantPermission  `json:"grant_pub_key"`
}

// AccountTables is the state of account storage, organized as a table.
type AccountTables struct {
	Accounts            []AccountRow     `json:"accounts"`
	AccountGrantPubKeys []GrantPubKeyRow `json:"account_grant_pub_keys"`
}

// ToIR -
func (a AccountTables) ToIR() *AccountTablesIR {
	tables := &AccountTablesIR{}
	for _, v := range a.Accounts {
		tables.Accounts = append(tables.Accounts, v.ToIR())
	}
	tables.AccountGrantPubKeys = a.AccountGrantPubKeys
	return tables
}
