package account

//go:generate mockery -name AccountKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
)

type AccountKeeper interface {
	DoesAccountExist(ctx sdk.Context, username types.AccountKey) bool
	CreateAccount(
		ctx sdk.Context, username types.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error
	MoveCoinFromUsernameToUsername(
		ctx sdk.Context, sender, receiver types.AccountKey, coin types.Coin) sdk.Error
	AddCoinToUsername(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error
	AddCoinToAddress(ctx sdk.Context, addr sdk.Address, coin types.Coin) sdk.Error
	MinusCoinFromUsername(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error
	MinusCoinFromAddress(ctx sdk.Context, address sdk.Address, coin types.Coin) sdk.Error
	UpdateJSONMeta(ctx sdk.Context, username types.AccountKey, JSONMeta string) sdk.Error
	GetTransactionKey(ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error)
	GetSigningKey(ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error)
	GetSavingFromAddress(ctx sdk.Context, address sdk.Address) (types.Coin, sdk.Error)
	GetSequence(ctx sdk.Context, address sdk.Address) (uint64, sdk.Error)
	GetAddress(ctx sdk.Context, username types.AccountKey) (sdk.Address, sdk.Error)
	GetFrozenMoneyList(ctx sdk.Context, addr sdk.Address) ([]model.FrozenMoney, sdk.Error)
	IncreaseSequenceByOne(ctx sdk.Context, address sdk.Address) sdk.Error
}

var _ AccountKeeper = AccountManager{}
