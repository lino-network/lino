package account

//go:generate mockery -name AccountKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/manager"
	"github.com/lino-network/lino/x/account/model"
)

type AccountKeeper interface {
	DoesAccountExist(ctx sdk.Context, username types.AccountKey) bool
	RegisterAccount(
		ctx sdk.Context, referrerAddr sdk.AccAddress, registerFee types.Coin,
		username types.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error
	CreateAccount(
		ctx sdk.Context, username types.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error
	MoveCoinFromUsernameToUsername(
		ctx sdk.Context, sender, receiver types.AccountKey, coin types.Coin) sdk.Error
	AddCoinToUsername(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error
	MinusCoinFromUsername(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error
	UpdateJSONMeta(ctx sdk.Context, username types.AccountKey, JSONMeta string) sdk.Error
	GetTransactionKey(ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error)
	GetSigningKey(ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error)
	GetSavingFromUsername(ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error)
	GetSequence(ctx sdk.Context, address sdk.Address) (uint64, sdk.Error)
	GetAddress(ctx sdk.Context, username types.AccountKey) (sdk.AccAddress, sdk.Error)
	GetFrozenMoneyList(ctx sdk.Context, addr sdk.Address) ([]model.FrozenMoney, sdk.Error)
	IncreaseSequenceByOne(ctx sdk.Context, address sdk.Address) sdk.Error
	AddFrozenMoney(
		ctx sdk.Context, username types.AccountKey, amount types.Coin, start, interval, times int64) sdk.Error
	CheckSigningPubKeyOwner(
		ctx sdk.Context, me types.AccountKey, signKey crypto.PubKey,
		permission types.Permission, amount types.Coin) (types.AccountKey, sdk.Error)
	AuthorizePermission(
		ctx sdk.Context, me types.AccountKey, grantTo types.AccountKey,
		validityPeriod int64, grantLevel types.Permission, amount types.Coin) sdk.Error
	RevokePermission(
		ctx sdk.Context, me, grantTo types.AccountKey, permission types.Permission) sdk.Error

	// getter
	GetInfo(ctx sdk.Context, username types.AccountKey) (*model.AccountInfo, sdk.Error)
	GetBank(ctx sdk.Context, username types.AccountKey) (*model.AccountBank, sdk.Error)
	GetMeta(ctx sdk.Context, username types.AccountKey) (*model.AccountMeta, sdk.Error)
	GetReward(ctx sdk.Context, username types.AccountKey) (*model.Reward, sdk.Error)
	GetGrantPubKeys(ctx sdk.Context, username, grantTo types.AccountKey) ([]*model.GrantPermission, sdk.Error)
	GetAllGrantPubKeys(ctx sdk.Context, username types.AccountKey) ([]*model.GrantPermission, sdk.Error)
}

var _ AccountKeeper = manager.AccountManager{}
