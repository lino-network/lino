package model

import (
	crypto "github.com/tendermint/tendermint/crypto"

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
	Reward              Reward                `json:"reward"`
	PendingCoinDayQueue PendingCoinDayQueueIR `json:"pending_coin_day_queue"`
}

// AccountRowIRV1 account related information when migrate, pk: Username
type AccountRowIRV1 struct {
	Username            types.AccountKey      `json:"username"`
	Info                AccountInfoV1         `json:"info"`
	Bank                AccountBankV1         `json:"bank"`
	Meta                AccountMetaV1         `json:"meta"`
	Reward              Reward                `json:"reward"`
	PendingCoinDayQueue PendingCoinDayQueueIR `json:"pending_coin_day_queue"`
}

// GrantPermissionIR - user grant permission to a user with a certain permission
// XXX(yumin): note that there is a field name change during upgrade-1.
type GrantPermissionIR struct {
	Username   types.AccountKey `json:"username"`
	Permission types.Permission `json:"permission"`
	CreatedAt  int64            `json:"created_at"`
	ExpiresAt  int64            `json:"expires_at"`
	Amount     types.Coin       `json:"amount"`
}

// ToState - convert IR back to state.
func (g GrantPermissionIR) ToState() *GrantPermission {
	return &GrantPermission{
		GrantTo:    g.Username,
		Permission: g.Permission,
		CreatedAt:  g.CreatedAt,
		ExpiresAt:  g.ExpiresAt,
		Amount:     g.Amount,
	}
}

// GrantPubKeyRowIR also in account, pk: (Username, pubKey)
type GrantPubKeyRowIR struct {
	Username    types.AccountKey  `json:"username"`
	PubKey      crypto.PubKey     `json:"pub_key"`
	GrantPubKey GrantPermissionIR `json:"grant_pub_key"`
}

// AccountTablesIR -
type AccountTablesIR struct {
	Accounts            []AccountRowIR     `json:"accounts"`
	AccountGrantPubKeys []GrantPubKeyRowIR `json:"account_grant_pub_keys"`
}

// AccountInfoV1 - user information
type AccountInfoV1 struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      int64            `json:"created_at"`
	ResetKey       crypto.PubKey    `json:"reset_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	AppKey         crypto.PubKey    `json:"app_key"`
}

// AccountBankV1 - user balance
type AccountBankV1 struct {
	Saving          types.Coin    `json:"saving"`
	CoinDay         types.Coin    `json:"coin_day"`
	FrozenMoneyList []FrozenMoney `json:"frozen_money_list"`
}

// AccountMetaV1 - stores tiny and frequently updated fields.
type AccountMetaV1 struct {
	Sequence             uint64     `json:"sequence"`
	LastActivityAt       int64      `json:"last_activity_at"`
	TransactionCapacity  types.Coin `json:"transaction_capacity"`
	JSONMeta             string     `json:"json_meta"`
	LastReportOrUpvoteAt int64      `json:"last_report_or_upvote_at"`
	LastPostAt           int64      `json:"last_post_at"`
}
