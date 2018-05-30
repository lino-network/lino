package model

import (
	"github.com/lino-network/lino/types"

	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountInfo stores general Lino Account information
type AccountInfo struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      int64            `json:"created_at"`
	MasterKey      crypto.PubKey    `json:"master_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	PostKey        crypto.PubKey    `json:"post_key"`
	Address        sdk.Address      `json:"address"`
}

// AccountBank uses Address as the key instead of Username
type AccountBank struct {
	Address         sdk.Address      `json:"address"`
	Saving          types.Coin       `json:"saving"`
	Username        types.AccountKey `json:"username"`
	Stake           types.Coin       `json:"stake"`
	FrozenMoneyList []FrozenMoney    `json:"frozen_money_list"`
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

// GrantKeyList stores a list of key authenticated by the use
type GrantKeyList struct {
	GrantPubKeyList []GrantPubKey `json:"grant_public_key_list"`
}

type GrantPubKey struct {
	Username  types.AccountKey `json:"username"`
	PubKey    crypto.PubKey    `json:"public_key"`
	ExpiresAt int64            `json:"expires_at"`
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence            int64      `json:"sequence"`
	LastActivityAt      int64      `json:"last_activity_at"`
	TransactionCapacity types.Coin `json:"transaction_capacity"`
}

// AccountInfraConsumption records infra utility consumption
type AccountInfraConsumption struct {
	Storage   int64 `json:"storage"`
	Bandwidth int64 `json:"bandwidth"`
}

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
	OriginalIncome types.Coin `json:"original_income"`
	FrictionIncome types.Coin `json:"friction_income"`
	ActualReward   types.Coin `json:"actual_reward"`
	UnclaimReward  types.Coin `json:"unclaim_reward"`
}

// Relationship between accounts
type Relationship struct {
	DonationTimes int64 `json:"donation_times"`
}

// BalanceHistory records all transactions in a certain time period
type BalanceHistory struct {
	Details []Detail `json:"details"`
}

// Detail is information about each transaction related to balance
type Detail struct {
	DetailType types.BalanceHistoryDetailType `json:"detail"`
	Amount     types.Coin                     `json:"amount"`
	CreatedAt  int64                          `json:"created_at"`
}
