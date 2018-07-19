package model

import (
	"github.com/lino-network/lino/types"

	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountInfo stores general Lino Account information
type AccountInfo struct {
	Username        types.AccountKey `json:"username"`
	CreatedAt       int64            `json:"created_at"`
	MasterKey       crypto.PubKey    `json:"master_key"`
	TransactionKey  crypto.PubKey    `json:"transaction_key"`
	MicropaymentKey crypto.PubKey    `json:"micropayment_key"`
	PostKey         crypto.PubKey    `json:"post_key"`
}

// AccountBank uses Address as the key instead of Username
type AccountBank struct {
	Saving          types.Coin    `json:"saving"`
	Stake           types.Coin    `json:"stake"`
	FrozenMoneyList []FrozenMoney `json:"frozen_money_list"`
	NumOfTx         int64         `json:"number_of_transaction"`
	NumOfReward     int64         `json:"number_of_reward"`
}

type FrozenMoney struct {
	Amount   types.Coin `json:"amount"`
	StartAt  int64      `json:"start_at"`
	Times    int64      `json:"times"`
	Interval int64      `json:"interval"`
}

// PendingStakeQueue stores a list of pending stake and total number of coin waiting in list
type PendingStakeQueue struct {
	LastUpdatedAt    int64          `json:"last_updated_at"`
	StakeCoinInQueue sdk.Rat        `json:"stake_coin_in_queue"`
	TotalCoin        types.Coin     `json:"total_coin"`
	PendingStakeList []PendingStake `json:"pending_stake_list"`
}

// pending stake in the list
type PendingStake struct {
	StartTime int64      `json:"start_time"`
	EndTime   int64      `json:"end_time"`
	Coin      types.Coin `json:"coin"`
}

type GrantPubKey struct {
	Username   types.AccountKey `json:"username"`
	Permission types.Permission `json:"permission"`
	LeftTimes  int64            `json:"left_times"`
	CreatedAt  int64            `json:"created_at"`
	ExpiresAt  int64            `json:"expires_at"`
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence             int64      `json:"sequence"`
	LastActivityAt       int64      `json:"last_activity_at"`
	TransactionCapacity  types.Coin `json:"transaction_capacity"`
	JSONMeta             string     `json:"json_meta"`
	LastReportOrUpvoteAt int64      `json:"last_report_or_upvote_at"`
}

// AccountInfraConsumption records infra utility consumption
// type AccountInfraConsumption struct {
// 	Storage   int64 `json:"storage"`
// 	Bandwidth int64 `json:"bandwidth"`
// }

// FollowerMeta record all meta info about this relation
type FollowerMeta struct {
	CreatedAt    int64            `json:"created_at"`
	FollowerName types.AccountKey `json:"follower_name"`
}

// FollowingMeta record all meta info about this relation
type FollowingMeta struct {
	CreatedAt     int64            `json:"created_at"`
	FollowingName types.AccountKey `json:"following_name"`
}

// Reward get from the inflation pool, only 1% of total income
type Reward struct {
	TotalIncome     types.Coin `json:"total_income"`
	OriginalIncome  types.Coin `json:"original_income"`
	FrictionIncome  types.Coin `json:"friction_income"`
	InflationIncome types.Coin `json:"inflation_income"`
	UnclaimReward   types.Coin `json:"unclaim_reward"`
}

type RewardDetail struct {
	OriginalDonation types.Coin       `json:"original_donation"`
	FrictionDonation types.Coin       `json:"friction_donation"`
	ActualReward     types.Coin       `json:"actual_reward"`
	Consumer         types.AccountKey `json:"consumer"`
	PostAuthor       types.AccountKey `json:"post_author"`
	PostID           string           `json:"post_id`
}

type RewardHistory struct {
	Details []RewardDetail `json:"details"`
}

// Relationship between accounts
type Relationship struct {
	DonationTimes int64 `json:"donation_times"`
}

// BalanceHistory records all transactions belong to the user
// Currently one balance history bundle can store at most 1000 transactions
// If number of transaction exceeds the limitation, a new bundle will be
// generated in KVStore
// Total number of history bundle is defined in metadata
type BalanceHistory struct {
	Details []Detail `json:"details"`
}

type Detail struct {
	DetailType types.TransferDetailType `json:"detail_type"`
	From       types.AccountKey         `json:"from"`
	To         types.AccountKey         `json:"to"`
	Amount     types.Coin               `json:"amount"`
	CreatedAt  int64                    `json:"created_at"`
	Memo       string                   `json:"memo"`
}
