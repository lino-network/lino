package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
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

// AccountManager stores and retrieves accounts from stores
// retrieved from the context.
type AccountManager interface {
	// Account getter/setter
	GetInfo(ctx sdk.Context, accKey AccountKey) (*AccountInfo, sdk.Error)
	SetInfo(ctx sdk.Context, accKey AccountKey, accInfo *AccountInfo) sdk.Error

	GetBankFromAccountKey(ctx sdk.Context, accKey AccountKey) (*AccountBank, sdk.Error)
	GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*AccountBank, sdk.Error)
	SetBank(ctx sdk.Context, address sdk.Address, accBank *AccountBank) sdk.Error

	GetMeta(ctx sdk.Context, accKey AccountKey) (*AccountMeta, sdk.Error)
	SetMeta(ctx sdk.Context, accKey AccountKey, accMeta *AccountMeta) sdk.Error

	GetFollowers(ctx sdk.Context, accKey AccountKey) (*Followers, sdk.Error)
	SetFollowers(ctx sdk.Context, accKey AccountKey, followers *Followers) sdk.Error

	GetFollowings(ctx sdk.Context, accKey AccountKey) (*Followings, sdk.Error)
	SetFollowings(ctx sdk.Context, accKey AccountKey, followings *Followings) sdk.Error
}

func RegisterWireLinoAccount(cdc *wire.Codec) {
	// Register crypto.[PubKey] types.
	wire.RegisterCrypto(cdc)
}