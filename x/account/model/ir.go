package model

import (
	crypto "github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

// AccountRowIR account related information when migrate, pk: Username
type AccountRowIR struct {
	Username types.AccountKey `json:"username"`
	Info     AccountInfoV1    `json:"info"`
	Bank     AccountBankV1    `json:"bank"`
	Meta     AccountMetaV1    `json:"meta"`
	Reward   RewardV1         `json:"reward"`
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

// RewardV1 - get from the inflation pool
type RewardV1 struct {
	TotalIncome     types.Coin `json:"total_income"`
	OriginalIncome  types.Coin `json:"original_income"`
	FrictionIncome  types.Coin `json:"friction_income"`
	InflationIncome types.Coin `json:"inflation_income"`
	UnclaimReward   types.Coin `json:"unclaim_reward"`
}
