package app

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/abci/types"
	oldwire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/lino-network/lino/genesis"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/auth"
	"github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/tx/register"
	"github.com/lino-network/lino/types"
)

const (
	appName = "LinoBlockchain"
)

// Extended ABCI application
type LinoBlockchain struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	capKeyAccountStore *sdk.KVStoreKey
	capKeyPostStore    *sdk.KVStoreKey
	capKeyValStore     *sdk.KVStoreKey
	capKeyIBCStore     *sdk.KVStoreKey
	capKeyGlobalStore  *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountManager acc.AccountManager
	postManager    post.PostManager
}

func NewLinoBlockchain(logger log.Logger, dbs map[string]dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:            bam.NewBaseApp(appName, logger, dbs["acc"]),
		cdc:                MakeCodec(),
		capKeyAccountStore: sdk.NewKVStoreKey(types.AccountKVStoreKey),
		capKeyPostStore:    sdk.NewKVStoreKey(types.PostKVStoreKey),
	}
	lb.accountManager = acc.NewLinoAccountManager(lb.capKeyAccountStore)
	lb.postManager = post.NewPostMananger(lb.capKeyPostStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(lb.accountManager), nil).
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager), nil).
		AddRoute(types.PostRouterName, post.NewHandler(lb.postManager, lb.accountManager), nil)

	lb.SetTxDecoder(lb.txDecoder)

	lb.MountStoreWithDB(lb.capKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	lb.MountStoreWithDB(lb.capKeyPostStore, sdk.StoreTypeIAVL, dbs["post"])
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager))
	if err := lb.LoadLatestVersion(lb.capKeyAccountStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyPostStore); err != nil {
		cmn.Exit(err.Error())
	}
	return lb
}

// custom tx codec
// TODO: use new go-wire
func MakeCodec() *wire.Codec {
	const msgTypeRegister = 0x1
	const msgTypeFollow = 0x2
	const msgTypeUnfollow = 0x3
	const msgTypeTransfer = 0x4
	const msgTypePost = 0x5
	const msgTypeLike = 0x6
	const msgTypeDonate = 0x7
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{register.RegisterMsg{}, msgTypeRegister},
		oldwire.ConcreteType{acc.FollowMsg{}, msgTypeFollow},
		oldwire.ConcreteType{acc.UnfollowMsg{}, msgTypeUnfollow},
		oldwire.ConcreteType{acc.TransferMsg{}, msgTypeTransfer},
		oldwire.ConcreteType{post.CreatePostMsg{}, msgTypePost},
		oldwire.ConcreteType{post.LikeMsg{}, msgTypeLike},
		oldwire.ConcreteType{post.DonateMsg{}, msgTypeDonate},
	)

	cdc := wire.NewCodec()

	return cdc
}

// custom logic for transaction decoding
func (lb *LinoBlockchain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx{}

	err := lb.cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}

// custom logic for lino blockchain initialization
func (lb *LinoBlockchain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	genesisState := new(genesis.GenesisState)
	if err := json.Unmarshal(stateJSON, genesisState); err != nil {
		panic(err)
	}

	for _, gacc := range genesisState.Accounts {
		if err := lb.toAppAccount(ctx, gacc); err != nil {
			panic(err)
		}
	}
	return abci.ResponseInitChain{}
}

// convert GenesisAccount to AppAccount
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga *genesis.GenesisAccount) sdk.Error {
	bank, err := lb.accountManager.GetBankFromAddress(ctx, ga.PubKey.Address())
	if err == nil {
		// account bank exists
		panic(sdk.ErrGenesisParse("genesis bank already exist"))
	} else {
		fmt.Println(ga)
		coin, err := types.LinoToCoin(types.TestLNO(sdk.NewRat(ga.Lino)))
		if err != nil {
			panic(err)
		}
		bank = &acc.AccountBank{
			Address: ga.PubKey.Address(),
			Balance: coin,
		}
		if setErr := lb.accountManager.SetBankFromAddress(ctx, bank.Address, bank); setErr != nil {
			panic(sdk.ErrGenesisParse("set genesis bank failed"))
		}
		account := acc.NewAccountProxy(acc.AccountKey(ga.Name), &lb.accountManager)
		if account.IsAccountExist(ctx) {
			panic(sdk.ErrGenesisParse("genesis account already exist"))
		}
		if err := account.CreateAccount(ctx, acc.AccountKey(ga.Name), ga.PubKey, bank); err != nil {
			panic(err)
		}
		if err := account.Apply(ctx); err != nil {
			panic(err)
		}
	}
	return nil
}
