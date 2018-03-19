package app

import (
	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/auth"
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
	capKeyMainStore *sdk.KVStoreKey
	capKeyIBCStore  *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountManager types.AccountManager
}

func NewLinoBlockchain(logger log.Logger, db dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:         bam.NewBaseApp(appName, logger, db),
		cdc:             MakeCodec(),
		capKeyMainStore: sdk.NewKVStoreKey("main"),
		capKeyIBCStore:  sdk.NewKVStoreKey("ibc"),
	}
	lb.accountManager = account.NewLinoAccountManager(lb.capKeyMainStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(lb.accountManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532
	lb.MountStoresIAVL(lb.capKeyMainStore)
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager))
	err := lb.LoadLatestVersion(lb.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return lb
}

// custom tx codec
// TODO: use new go-wire
func MakeCodec() *wire.Codec {

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
		return nil, sdk.ErrTxParse("").TraceCause(err, "")
	}
	return tx, nil
}

// custom logic for basecoin initialization
func (lb *LinoBlockchain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	// stateJSON := req.AppStateBytes

	// genesisState := new(sdk.GenesisState)
	// err := json.Unmarshal(stateJSON, genesisState)
	// if err != nil {
	// 	panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
	// 	// return sdk.ErrGenesisParse("").TraceCause(err, "")
	// }

	// for _, gacc := range genesisState.Accounts {
	// 	_, err := gacc.ToAppAccount()
	// 	if err != nil {
	// 		panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
	// 		//	return sdk.ErrGenesisParse("").TraceCause(err, "")
	// 	}
	// 	// lb.accountMapper.SetAccount(ctx, acc)
	// }
	return abci.ResponseInitChain{}
}
