package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

type Memo uint64

// AccountKey key format in KVStore
type AccountKey string

// AccountInfo stores general Lino Account information
type AccountInfo struct {
	Username AccountKey    `json:"username"`
	Created  types.Height  `json:"created"`
	PostKey  crypto.PubKey `json:"post_key"`
	OwnerKey crypto.PubKey `json:"owner_key"`
	Address  sdk.Address   `json:"address"`
}

// AccountBank uses Address as the key instead of Username
type AccountBank struct {
	Address  sdk.Address `json:"address"`
	Balance  types.Coin  `json:"balance"`
	Username AccountKey  `json:"username"`
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence       int64        `json:"sequence"`
	LastActivity   types.Height `json:"last_activity"`
	ActivityBurden int64        `json:"activity_burden"`
}

// AccountInfraConsumption records infra utility consumption
type AccountInfraConsumption struct {
	Storage   int64 `json:"storage"`
	Bandwidth int64 `json:"bandwidth"`
}

// record all meta info about this relation
type FollowerMeta struct {
	CreatedAt    types.Height `json:"created_at"`
	FollowerName AccountKey   `json:"follower_name"`
}

// record all meta info about this relation
type FollowingMeta struct {
	CreatedAt    types.Height `json:"created_at"`
	FolloweeName AccountKey   `json:"followee_name"`
}
