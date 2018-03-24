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
	capKeyIBCStore     *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountManager acc.AccountManager
	postManager    post.PostManager
}

func NewLinoBlockchain(logger log.Logger, db dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:            bam.NewBaseApp(appName, logger, db),
		cdc:                MakeCodec(),
		capKeyAccountStore: sdk.NewKVStoreKey("account"),
		capKeyPostStore:    sdk.NewKVStoreKey("post"),
		capKeyIBCStore:     sdk.NewKVStoreKey("ibc"),
	}
	lb.accountManager = acc.NewLinoAccountManager(lb.capKeyAccountStore)
	lb.postManager = post.NewPostMananger(lb.capKeyPostStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(lb.accountManager)).
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(lb.postManager, lb.accountManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532
	lb.MountStoresIAVL(lb.capKeyAccountStore)
	lb.MountStoresIAVL(lb.capKeyPostStore)
	lb.MountStoresIAVL(lb.capKeyIBCStore)
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager))
	if err := lb.LoadLatestVersion(lb.capKeyAccountStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyPostStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyIBCStore); err != nil {
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

	// TODO(Lino): Register msg type and model.
	cdc := wire.NewCodec()

	return cdc
}

// custom logic for transaction decoding
func (lb *LinoBlockchain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx{}

	// StdTx.Msg is an interface. The concrete types
	// are registered by MakeTxCodec in bank.RegisterWire.
	err := lb.cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}

// custom logic for basecoin initialization
func (lb *LinoBlockchain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(acc.GenesisState)
	err := json.Unmarshal(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
	}

	for _, gacc := range genesisState.Accounts {
		if err := lb.toAppAccount(ctx, gacc); err != nil {
			panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		// lb.accountMapper.SetAccount(ctx, acc)
	}
	return abci.ResponseInitChain{}
}

// convert GenesisAccount to AppAccount
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga *acc.GenesisAccount) sdk.Error {
	// send coins using address (even no account bank associated with this addr)
	fmt.Println("===============genesis account: ", ga.Address)
	bank, err := lb.accountManager.GetBankFromAddress(ctx, ga.Address)
	if err == nil {
		// account bank exists
		panic(sdk.ErrGenesisParse("genesis account already exist"))
	} else {
		// account bank not found, create a new one for this address
		bank = &acc.AccountBank{
			Address: ga.Address,
			Balance: ga.Coins,
		}
		if setErr := lb.accountManager.SetBankFromAddress(ctx, ga.Address, bank); setErr != nil {
			panic(sdk.ErrGenesisParse("set genesis account failed"))
		}
	}
	return nil
}
