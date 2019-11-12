package account

//go:generate mockery -name AccountKeeper

import (
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/manager"
	"github.com/lino-network/lino/x/account/model"
)

type AccountKeeper interface {
	InitGenesis(ctx sdk.Context, total types.Coin, pools []model.Pool)
	// core bank APIs.
	MoveCoin(ctx sdk.Context, sender, receiver types.AccOrAddr, coin types.Coin) sdk.Error
	MoveFromPool(
		ctx sdk.Context, poolName types.PoolName, dest types.AccOrAddr, amount types.Coin) sdk.Error
	MoveToPool(
		ctx sdk.Context, poolName types.PoolName, from types.AccOrAddr, amount types.Coin) sdk.Error
	MoveBetweenPools(ctx sdk.Context, from, to types.PoolName, amount types.Coin) sdk.Error
	Mint(ctx sdk.Context) sdk.Error

	DoesAccountExist(ctx sdk.Context, username types.AccountKey) bool
	GenesisAccount(ctx sdk.Context, username types.AccountKey,
		signingKey, transactionKey crypto.PubKey) sdk.Error
	RegisterAccount(
		ctx sdk.Context, referrer types.AccOrAddr, registerFee types.Coin,
		username types.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error
	UpdateJSONMeta(ctx sdk.Context, username types.AccountKey, JSONMeta string) sdk.Error
	GetPool(ctx sdk.Context, poolName types.PoolName) (types.Coin, sdk.Error)
	GetTransactionKey(ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error)
	GetSigningKey(ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error)
	GetSavingFromUsername(ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error)
	GetSequence(ctx sdk.Context, address sdk.Address) (uint64, sdk.Error)
	GetAddress(ctx sdk.Context, username types.AccountKey) (sdk.AccAddress, sdk.Error)
	GetSupply(ctx sdk.Context) model.Supply
	IncreaseSequenceByOne(ctx sdk.Context, address sdk.Address) sdk.Error
	AddPending(
		ctx sdk.Context, username types.AccountKey, amount types.Coin) sdk.Error
	CheckSigningPubKeyOwner(
		ctx sdk.Context, me types.AccountKey, signKey crypto.PubKey) (types.AccountKey, sdk.Error)
	CheckSigningPubKeyOwnerByAddress(
		ctx sdk.Context, addr sdk.AccAddress, signkey crypto.PubKey, isPaid bool) sdk.Error
	RecoverAccount(
		ctx sdk.Context, username types.AccountKey, newTransactionPubKey, newSigningKey crypto.PubKey) sdk.Error

	// getter
	GetInfo(ctx sdk.Context, username types.AccountKey) (*model.AccountInfo, sdk.Error)
	GetBank(ctx sdk.Context, username types.AccountKey) (*model.AccountBank, sdk.Error)
	GetBankByAddress(ctx sdk.Context, addr sdk.AccAddress) (*model.AccountBank, sdk.Error)
	GetMeta(ctx sdk.Context, username types.AccountKey) (*model.AccountMeta, sdk.Error)

	// import export
	ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
	ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
}

var _ AccountKeeper = manager.AccountManager{}
