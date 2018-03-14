package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/go-crypto"
)

type AccountKey []byte

// Lino Account
// key: Username AccountKey
type Account struct {
	Created  uint64        `json:"created"`
	PostKey  crypto.PubKey `json:"post_key"`
	OwnerKey crypto.PubKey `json:"owner_key"`
	Address  sdk.Address   `json:"address"`
}

// AccountBank embeds base account, handle the balence, which implements sdk.Account
type AccountBank struct {
	auth.BaseAccount
}

// AccountMeta stores tiny and frequently updated fields.
// key: Username AccountKey
type AccountMeta struct {
	LastActivity   uint64 `json:"last_activity"`
	ActivityBurden uint64 `json:"activity_burden"`
	LastABBlock    uint64 `json:"last_activity_burden_block"`
}

// key: Username AccountKey
type Follower struct {
	Followers []AccountKey `json:"followers"`
}

// key: Username AccountKey
type Followings struct {
	Followings []AccountKey `json:"followings"`
}
