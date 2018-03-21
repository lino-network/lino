package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/go-crypto"
)

type Memo uint64

// AccountInfo stores general Lino Account information
type AccountInfo struct {
	Username AccountKey    `json:"key"`
	Created  Height        `json:"created"`
	PostKey  crypto.PubKey `json:"post_key"`
	OwnerKey crypto.PubKey `json:"owner_key"`
	Address  sdk.Address   `json:"address"`
}

// AccountBank uses Address as the key instead of Username
type AccountBank struct {
	Address  sdk.Address `json:"address"`
	Balance  sdk.Coins   `json:"coins"`
	Username AccountKey  `json:"Username"`
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence       int64  `json:"sequence"`
	LastActivity   Height `json:"last_activity"`
	ActivityBurden int64  `json:"activity_burden"`
}

// Follower records all follower belong to one user
type Follower struct {
	Follower []AccountKey `json:"follower"`
}

// Following records all follower belong to one user
type Following struct {
	Following []AccountKey `json:"following"`
}

func RegisterWireLinoAccount(cdc *wire.Codec) {
	// Register crypto.[PubKey] types.
	wire.RegisterCrypto(cdc)
}
