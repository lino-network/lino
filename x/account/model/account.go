package model

import (
	"github.com/lino-network/lino/types"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountInfo - user information
type AccountInfo struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      int64            `json:"created_at"`
	ResetKey       crypto.PubKey    `json:"reset_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	AppKey         crypto.PubKey    `json:"app_key"`
}

// AccountBank - user balance
type AccountBank struct {
	Saving          types.Coin    `json:"saving"`
	CoinDay         types.Coin    `json:"coin_day"`
	FrozenMoneyList []FrozenMoney `json:"frozen_money_list"`
}

// FrozenMoney - frozen money
type FrozenMoney struct {
	Amount   types.Coin `json:"amount"`
	StartAt  int64      `json:"start_at"`
	Times    int64      `json:"times"`
	Interval int64      `json:"interval"`
}

// PendingCoinDayQueue - stores a list of pending coin day and total number of coin waiting in list
type PendingCoinDayQueue struct {
	LastUpdatedAt   int64            `json:"last_updated_at"`
	TotalCoinDay    sdk.Dec          `json:"total_coin_day"`
	TotalCoin       types.Coin       `json:"total_coin"`
	PendingCoinDays []PendingCoinDay `json:"pending_coin_days"`
}

// ToIR coin.
func (p PendingCoinDayQueue) ToIR() PendingCoinDayQueueIR {
	return PendingCoinDayQueueIR{
		LastUpdatedAt:   p.LastUpdatedAt,
		TotalCoinDay:    p.TotalCoinDay.String(),
		TotalCoin:       p.TotalCoin,
		PendingCoinDays: p.PendingCoinDays,
	}
}

// PendingCoinDay - pending coin day in the list
type PendingCoinDay struct {
	StartTime int64      `json:"start_time"`
	EndTime   int64      `json:"end_time"`
	Coin      types.Coin `json:"coin"`
}

// GrantPermission - user grant permission to a user with a certain permission
type GrantPermission struct {
	GrantTo    types.AccountKey `json:"grant_to"`
	Permission types.Permission `json:"permission"`
	CreatedAt  int64            `json:"created_at"`
	ExpiresAt  int64            `json:"expires_at"`
	Amount     types.Coin       `json:"amount"`
}

// ToIR - name change, username -> GrantTo
func (g GrantPermission) ToIR() GrantPermissionIR {
	return GrantPermissionIR{
		Username:   g.GrantTo,
		Permission: g.Permission,
		CreatedAt:  g.CreatedAt,
		ExpiresAt:  g.ExpiresAt,
		Amount:     g.Amount,
	}
}

// AccountMeta - stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence             uint64     `json:"sequence"`
	LastActivityAt       int64      `json:"last_activity_at"`
	TransactionCapacity  types.Coin `json:"transaction_capacity"`
	JSONMeta             string     `json:"json_meta"`
	LastReportOrUpvoteAt int64      `json:"last_report_or_upvote_at"`
	LastPostAt           int64      `json:"last_post_at"`
}

// AccountInfraConsumption records infra utility consumption
// type AccountInfraConsumption struct {
// 	Storage   int64 `json:"storage"`
// 	Bandwidth int64 `json:"bandwidth"`
// }

// Reward - get from the inflation pool
type Reward struct {
	TotalIncome     types.Coin `json:"total_income"`
	OriginalIncome  types.Coin `json:"original_income"`
	FrictionIncome  types.Coin `json:"friction_income"`
	InflationIncome types.Coin `json:"inflation_income"`
	UnclaimReward   types.Coin `json:"unclaim_reward"`
}
