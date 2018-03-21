package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-crypto"
)

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
