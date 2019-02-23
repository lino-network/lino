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
	NumOfTx         int64         `json:"number_of_transaction"`
	NumOfReward     int64         `json:"number_of_reward"`
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
		TotalCoinDay:    p.TotalCoinDay.FloatString(),
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

// GrantPubKey - user grant permission to a public key with a certain permission
type GrantPubKey struct {
	Username   types.AccountKey `json:"username"`
	Permission types.Permission `json:"permission"`
	CreatedAt  int64            `json:"created_at"`
	ExpiresAt  int64            `json:"expires_at"`
	Amount     types.Coin       `json:"amount"`
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

// FollowerMeta - record all meta info about this relation
type FollowerMeta struct {
	CreatedAt    int64            `json:"created_at"`
	FollowerName types.AccountKey `json:"follower_name"`
}

// FollowingMeta - record all meta info about this relation
type FollowingMeta struct {
	CreatedAt     int64            `json:"created_at"`
	FollowingName types.AccountKey `json:"following_name"`
}

// Reward - get from the inflation pool
type Reward struct {
	TotalIncome     types.Coin `json:"total_income"`
	OriginalIncome  types.Coin `json:"original_income"`
	FrictionIncome  types.Coin `json:"friction_income"`
	InflationIncome types.Coin `json:"inflation_income"`
	UnclaimReward   types.Coin `json:"unclaim_reward"`
}

// RewardDetail - reward detail
type RewardDetail struct {
	OriginalDonation types.Coin       `json:"original_donation"`
	FrictionDonation types.Coin       `json:"friction_donation"`
	ActualReward     types.Coin       `json:"actual_reward"`
	Consumer         types.AccountKey `json:"consumer"`
	PostAuthor       types.AccountKey `json:"post_author"`
	PostID           string           `json:"post_id"`
}

// RewardHistory - reward history
type RewardHistory struct {
	Details []RewardDetail `json:"details"`
}

// Relationship - relation between two users
type Relationship struct {
	DonationTimes int64 `json:"donation_times"`
}

// BalanceHistory - records all transactions belong to the user
// Currently one balance history bundle can store at most 1000 transactions
// If number of transaction exceeds the limitation, a new bundle will be
// generated in KVStore
// Total number of history bundle is defined in metadata
type BalanceHistory struct {
	Details []Detail `json:"details"`
}

// Detail - detail of each income and outcome
type Detail struct {
	DetailType types.TransferDetailType `json:"detail_type"`
	From       types.AccountKey         `json:"from"`
	To         types.AccountKey         `json:"to"`
	Amount     types.Coin               `json:"amount"`
	Balance    types.Coin               `json:"balance"`
	CreatedAt  int64                    `json:"created_at"`
	Memo       string                   `json:"memo"`
}
