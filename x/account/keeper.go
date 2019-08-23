package account

//go:generate mockery -name AccountKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

type AccountKeeper interface {
	DoesAccountExist(ctx sdk.Context, username types.AccountKey) bool
	CreateAccount(
		ctx sdk.Context, referrer types.AccountKey, username types.AccountKey,
		resetKey, transactionKey, appKey crypto.PubKey, registerDeposit types.Coin) sdk.Error
	GetCoinDay(
		ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error)
	AddSavingCoin(
		ctx sdk.Context, username types.AccountKey, coin types.Coin, from types.AccountKey, memo string,
		detailType types.TransferDetailType) (err sdk.Error)
	AddSavingCoinWithFullCoinDay(
		ctx sdk.Context, username types.AccountKey, coin types.Coin, from types.AccountKey, memo string,
		detailType types.TransferDetailType) (err sdk.Error)
	MinusSavingCoin(
		ctx sdk.Context, username types.AccountKey, coin types.Coin, to types.AccountKey,
		memo string, detailType types.TransferDetailType) (err sdk.Error)
	MinusSavingCoinWithFullCoinDay(
		ctx sdk.Context, username types.AccountKey, coin types.Coin, to types.AccountKey,
		memo string, detailType types.TransferDetailType) (types.Coin, sdk.Error)
}

var _ AccountKeeper = AccountManager{}
