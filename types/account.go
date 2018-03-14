package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/go-crypto"
)

// AccountInfo stores general Lino Account information
type AccountInfo struct {
	Key      AccountKey    `json:"key"`
	Created  Height        `json:"created"`
	PostKey  crypto.PubKey `json:"post_key"`
	OwnerKey crypto.PubKey `json:"owner_key"`
	Address  sdk.Address   `json:"address"`
}

// AccountBank embeds base account, handle the balance, which implements sdk.Account
type AccountBank struct {
	auth.BaseAccount
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	LastActivity   Height `json:"last_activity"`
	ActivityBurden uint64 `json:"activity_burden"`
	LastABBlock    Height `json:"last_activity_burden_block"`
}

// Followers records all followers belong to one user
type Followers struct {
	Followers []AccountKey `json:"followers"`
}

// Followings records all followers belong to one user
type Followings struct {
	Followings []AccountKey `json:"followings"`
}

// Account interface is general interfacc for account. Embed sdk Account.
type Account interface {
	// embed sdk.Account
	sdk.Account

	GetUsername() AccountKey
	SetUsername(AccountKey) error // errors if already set.

	GetPostKey() crypto.PubKey
	SetPostKey(crypto.PubKey) error

	GetOwnerKey() crypto.PubKey
	SetOwnerKey(crypto.PubKey) error

	GetCreated() Height
	SetCreated(Height) error // errors if already set.

	GetLastActivity() Height
	SetLastActivity(Height) error

	GetActivityBurden() uint64
	SetActivityBurden(uint64) error // set AB Block too.

	GetLastABBlock() Height

	GetFollowers() Followers
	SetFollowers(Followers) error

	GetFollowings() Followings
	SetFollowings(Followings) error
}

// AccountManager stores and retrieves accounts from stores
// retrieved from the context.
type AccountManager interface {
	AccountExist(ctx sdk.Context, accKey AccountKey) bool
	// Account getter/setter
	GetAccount(ctx sdk.Context, accKey AccountKey) Account
	SetAccount(ctx sdk.Context, acc Account)
}
