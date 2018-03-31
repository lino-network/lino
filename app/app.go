package app

import (
	"encoding/json"

	abci "github.com/tendermint/abci/types"
	oldwire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/auth"
	"github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/tx/register"
	"github.com/lino-network/lino/tx/validator"
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
	valManager     validator.ValidatorManager
	globalManager  global.GlobalManager
}

func NewLinoBlockchain(logger log.Logger, db dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:            bam.NewBaseApp(appName, logger, db),
		cdc:                MakeCodec(),
		capKeyAccountStore: sdk.NewKVStoreKey(types.AccountKVStoreKey),
		capKeyPostStore:    sdk.NewKVStoreKey(types.PostKVStoreKey),
		capKeyValStore:     sdk.NewKVStoreKey(types.ValidatorKVStoreKey),
		capKeyGlobalStore:  sdk.NewKVStoreKey(types.GlobalKVStoreKey),
		capKeyIBCStore:     sdk.NewKVStoreKey("ibc"),
	}
	lb.accountManager = acc.NewLinoAccountManager(lb.capKeyAccountStore)
	lb.postManager = post.NewPostMananger(lb.capKeyPostStore)
	lb.valManager = validator.NewValidatorMananger(lb.capKeyValStore)
	lb.globalManager = global.NewGlobalManager(lb.capKeyGlobalStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(lb.accountManager), nil).
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager), nil).
		AddRoute(types.PostRouterName, post.NewHandler(lb.postManager, lb.accountManager, lb.globalManager), lb.globalManager.InitGenesis).
		AddRoute(types.ValidatorRouterName, validator.NewHandler(lb.valManager, lb.accountManager), lb.valManager.InitGenesis)

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetEndBlocker(lb.endBlocker)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532
	lb.MountStoresIAVL(lb.capKeyAccountStore)
	lb.MountStoresIAVL(lb.capKeyPostStore)
	lb.MountStoresIAVL(lb.capKeyValStore)
	lb.MountStoresIAVL(lb.capKeyGlobalStore)
	lb.MountStoresIAVL(lb.capKeyIBCStore)
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager))
	if err := lb.LoadLatestVersion(lb.capKeyAccountStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyPostStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyValStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyGlobalStore); err != nil {
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
	genesisState := new(genesis.GenesisState)
	//err := oldwire.UnmarshalJSON(stateJSON, genesisState)
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
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga *genesis.GenesisAccount) sdk.Error {
	// send coins using address (even no account bank associated with this addr)
	bank, err := lb.accountManager.GetBankFromAddress(ctx, ga.PubKey.Address())
	if err == nil {
		// account bank exists
		panic(sdk.ErrGenesisParse("genesis bank already exist"))
	} else {
		coin, err := types.LinoToCoin(ga.Lino)
		if err != nil {
			panic(err)
		}
		// account bank not found, create a new one for this address
		bank = &acc.AccountBank{
			Address: ga.PubKey.Address(),
			Balance: coin,
		}
		if setErr := lb.accountManager.SetBankFromAddress(ctx, bank.Address, bank); setErr != nil {
			panic(sdk.ErrGenesisParse("set genesis bank failed"))
		}
		account := acc.NewProxyAccount(acc.AccountKey(ga.Name), &lb.accountManager)
		if account.IsAccountExist(ctx) {
			panic(sdk.ErrGenesisParse("genesis account already exist"))
		}
		if err := account.CreateAccount(ctx, acc.AccountKey(ga.Name), ga.PubKey, bank); err != nil {
			panic(err)
		}

		deposit := types.Coin{1000 * types.Decimals}
		// withdraw money from validator's bank
		if err := account.MinusCoin(ctx, deposit); err != nil {
			panic(err)
		}
		val := &validator.Validator{
			ABCIValidator: abci.Validator{PubKey: ga.ValPubKey.Bytes(), Power: deposit.Amount},
			Username:      account.GetUsername(ctx),
			Deposit:       deposit,
		}
		if setErr := lb.valManager.SetValidator(ctx, account.GetUsername(ctx), val); setErr != nil {
			panic(setErr)
		}

		if addErr := lb.valManager.AddToCandidatePool(ctx, account.GetUsername(ctx)); addErr != nil {
			panic(addErr)
		}
		if joinErr := lb.valManager.TryBecomeOncallValidator(ctx, account.GetUsername(ctx)); joinErr != nil {
			panic(joinErr)
		}
		if err := account.Apply(ctx); err != nil {
			panic(err)
		}
	}
	return nil
}

func (lb *LinoBlockchain) beginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	absentValidators := req.GetAbsentValidators()
	if absentValidators != nil {
		// TODO Err handling
		lb.valManager.UpdateAbsentValidator(ctx, absentValidators)
	}
	// TODO Err handling
	lb.valManager.FireIncompetentValidator(ctx, req.GetByzantineValidators())
	return abci.ResponseBeginBlock{}
}

func (lb *LinoBlockchain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	valList, err := lb.valManager.GetOncallValList(ctx)
	if err != nil {
		panic(err)
	}
	ABCIValList := make([]abci.Validator, len(valList))
	for i, validator := range valList {
		ABCIValList[i] = validator.ABCIValidator
	}
	return abci.ResponseEndBlock{ValidatorUpdates: ABCIValList}
}
