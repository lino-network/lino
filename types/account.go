package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/go-crypto"
)

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
	Coins    sdk.Coins   `json:"coins"`
	Username AccountKey  `json:"Username"`
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence       int64  `json:"sequence"`
	LastActivity   Height `json:"last_activity"`
	ActivityBurden uint64 `json:"activity_burden"`
	LastABBlock    Height `json:"last_activity_burden_block"`
}

// Follower records all follower belong to one user
type Follower struct {
	Follower []AccountKey `json:"follower"`
}

// Following records all follower belong to one user
type Following struct {
	Following []AccountKey `json:"following"`
}

// AccountManager stores and retrieves accounts from stores
// retrieved from the context.
type AccountManager interface {
	// Account getter/setter
	CreateAccount(ctx sdk.Context, accKey AccountKey, pubkey crypto.PubKey, accBank *AccountBank) (*AccountInfo, sdk.Error)
	AccountExist(ctx sdk.Context, accKey AccountKey) bool
	GetInfo(ctx sdk.Context, accKey AccountKey) (*AccountInfo, sdk.Error)
	SetInfo(ctx sdk.Context, accKey AccountKey, accInfo *AccountInfo) sdk.Error

	GetBankFromAccountKey(ctx sdk.Context, accKey AccountKey) (*AccountBank, sdk.Error)
	GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*AccountBank, sdk.Error)
	SetBank(ctx sdk.Context, address sdk.Address, accBank *AccountBank) sdk.Error

	GetMeta(ctx sdk.Context, accKey AccountKey) (*AccountMeta, sdk.Error)
	SetMeta(ctx sdk.Context, accKey AccountKey, accMeta *AccountMeta) sdk.Error

	GetFollower(ctx sdk.Context, accKey AccountKey) (*Follower, sdk.Error)
	SetFollower(ctx sdk.Context, accKey AccountKey, follower *Follower) sdk.Error

	GetFollowing(ctx sdk.Context, accKey AccountKey) (*Following, sdk.Error)
	SetFollowing(ctx sdk.Context, accKey AccountKey, following *Following) sdk.Error
}

func RegisterWireLinoAccount(cdc *wire.Codec) {
	// Register crypto.[PubKey] types.
	wire.RegisterCrypto(cdc)
}
