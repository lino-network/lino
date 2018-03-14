package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/go-crypto"
)

// Account is general Lino Account information
type Account struct {
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
