package app

import (
	"encoding/json"

	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/cosmos/cosmos-sdk/examples/basecoin/types"
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
	accountMapper sdk.AccountMapper
}

func NewLinoBlockchain(logger log.Logger, db dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:         bam.NewBaseApp(appName, logger, db),
		cdc:             MakeCodec(),
		capKeyMainStore: sdk.NewKVStoreKey("main"),
		capKeyIBCStore:  sdk.NewKVStoreKey("ibc"),
	}

	// define the accountMapper
	// lb.accountMapper = auth.NewAccountMapperSealed(
	// 	lb.capKeyMainStore, // target store
	// 	&types.AppAccount{}, // prototype
	// )

	// TODO(Lino): add handler
	// lb.Router().AddRoute("bank", bank.NewHandler(coinKeeper))

	// initialize BaseApp
	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532
	lb.MountStoresIAVL(lb.capKeyMainStore) // , app.capKeyIBCStore)
	// TODO(Lino): add antehandler here
	// app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper))
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
	// const msgTypeSend = 0x1
	// const msgTypeIssue = 0x2
	// const msgTypeQuiz = 0x3
	// const msgTypeSetTrend = 0x4
	// var _ = oldwire.RegisterInterface(
	// 	struct{ sdk.Msg }{},
	// )

	// const accTypeApp = 0x1
	// var _ = oldwire.RegisterInterface(
	// 	struct{ sdk.Account }{},
	// 	oldwire.ConcreteType{&types.AppAccount{}, accTypeApp},
	// )
	cdc := wire.NewCodec()

	// cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	// bank.RegisterWire(cdc)   // Register bank.[SendMsg,IssueMsg] types.
	// crypto.RegisterWire(cdc) // Register crypto.[PubKey,PrivKey,Signature] types.
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
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := json.Unmarshal(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	for _, gacc := range genesisState.Accounts {
		_, err := gacc.ToAppAccount()
		if err != nil {
			panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		// lb.accountMapper.SetAccount(ctx, acc)
	}
	return abci.ResponseInitChain{}
}
