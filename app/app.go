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
	val "github.com/lino-network/lino/tx/validator"
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
	accountManager *acc.AccountManager
	postManager    *post.PostManager
	valManager     *val.ValidatorManager
	globalManager  *global.GlobalManager

	lastBlockTime int64
	// for recurring time based event
	pastMinutes int64
}

func NewLinoBlockchain(logger log.Logger, dbs map[string]dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:            bam.NewBaseApp(appName, logger, dbs["acc"]),
		cdc:                MakeCodec(),
		capKeyAccountStore: sdk.NewKVStoreKey(types.AccountKVStoreKey),
		capKeyPostStore:    sdk.NewKVStoreKey(types.PostKVStoreKey),
		capKeyValStore:     sdk.NewKVStoreKey(types.ValidatorKVStoreKey),
		capKeyGlobalStore:  sdk.NewKVStoreKey(types.GlobalKVStoreKey),
		capKeyIBCStore:     sdk.NewKVStoreKey("ibc"),
	}
	lb.accountManager = acc.NewAccountManager(lb.capKeyAccountStore)
	lb.postManager = post.NewPostManager(lb.capKeyPostStore)
	lb.valManager = val.NewValidatorManager(lb.capKeyValStore)
	lb.globalManager = global.NewGlobalManager(lb.capKeyGlobalStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(*lb.accountManager)).
		AddRoute(types.AccountRouterName, acc.NewHandler(*lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(*lb.postManager, *lb.accountManager, *lb.globalManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(*lb.valManager, *lb.accountManager, *lb.globalManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoreWithDB(lb.capKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	lb.MountStoreWithDB(lb.capKeyPostStore, sdk.StoreTypeIAVL, dbs["post"])
	lb.MountStoreWithDB(lb.capKeyValStore, sdk.StoreTypeIAVL, dbs["val"])
	lb.MountStoreWithDB(lb.capKeyGlobalStore, sdk.StoreTypeIAVL, dbs["global"])
	lb.SetAnteHandler(auth.NewAnteHandler(*lb.accountManager))
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
	const msgTypeValidatorDeposit = 0x8
	const msgTypeValidatorWithdraw = 0x9
	const msgTypeValidatorRevoke = 0x10

	const msgTypeClaim = 0x11
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{register.RegisterMsg{}, msgTypeRegister},
		oldwire.ConcreteType{acc.FollowMsg{}, msgTypeFollow},
		oldwire.ConcreteType{acc.UnfollowMsg{}, msgTypeUnfollow},
		oldwire.ConcreteType{acc.TransferMsg{}, msgTypeTransfer},
		oldwire.ConcreteType{post.CreatePostMsg{}, msgTypePost},
		oldwire.ConcreteType{post.LikeMsg{}, msgTypeLike},
		oldwire.ConcreteType{post.DonateMsg{}, msgTypeDonate},
		oldwire.ConcreteType{val.ValidatorDepositMsg{}, msgTypeValidatorDeposit},
		oldwire.ConcreteType{val.ValidatorWithdrawMsg{}, msgTypeValidatorWithdraw},
		oldwire.ConcreteType{val.ValidatorRevokeMsg{}, msgTypeValidatorRevoke},
		oldwire.ConcreteType{acc.ClaimMsg{}, msgTypeClaim},
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
	if err := json.Unmarshal(stateJSON, genesisState); err != nil {
		panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
	}

	if err := lb.valManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.globalManager.InitGlobalManager(ctx, genesisState.GlobalState); err != nil {
		panic(err)
	}
	for _, gacc := range genesisState.Accounts {
		if err := lb.toAppAccount(ctx, gacc); err != nil {
			panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		// lb.accountMapper.SetAccount(ctx, acc)
	}
	if lb.lastBlockTime == 0 {
		lb.lastBlockTime = ctx.BlockHeader().Time
	}

	if lb.pastMinutes == 0 {
		lb.pastMinutes = ctx.BlockHeader().Time / 60
	}

	return abci.ResponseInitChain{}
}

// convert GenesisAccount to AppAccount
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga genesis.GenesisAccount) sdk.Error {
	// send coins using address (even no account bank associated with this addr)
	coin, err := types.LinoToCoin(types.LNO(sdk.NewRat(ga.Lino)))
	if err != nil {
		panic(err)
	}
	if setErr := lb.accountManager.AddCoinToAddress(ctx, ga.PubKey.Address(), coin); setErr != nil {
		panic(sdk.ErrGenesisParse("set genesis bank failed"))
	}
	if lb.accountManager.IsAccountExist(ctx, types.AccountKey(ga.Name)) {
		panic(sdk.ErrGenesisParse("genesis account already exist"))
	}
	if err := lb.accountManager.CreateAccount(ctx, types.AccountKey(ga.Name), ga.PubKey, types.NewCoin(0)); err != nil {
		panic(err)
	}

	deposit := types.NewCoin(1000 * types.Decimals)
	// withdraw money from validator's bank
	if err := lb.accountManager.MinusCoin(ctx, types.AccountKey(ga.Name), deposit); err != nil {
		panic(err)
	}

	if addErr := lb.valManager.RegisterValidator(ctx, types.AccountKey(ga.Name), ga.ValPubKey.Bytes(), deposit); addErr != nil {
		panic(addErr)
	}
	if joinErr := lb.valManager.TryBecomeOncallValidator(ctx, types.AccountKey(ga.Name)); joinErr != nil {
		panic(joinErr)
	}
	return nil
}

func (lb *LinoBlockchain) beginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	if ctx.BlockHeader().Time/60 > lb.pastMinutes {
		lb.increaseMinute(ctx)
	}
	if err := lb.valManager.SetPreRoundValidators(ctx); err != nil {
		panic(err)
	}
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
	if heightEvents := lb.globalManager.GetHeightEventListAtHeight(ctx, ctx.BlockHeight()); heightEvents != nil {
		lb.executeEvents(ctx, heightEvents.Events)
		lb.globalManager.RemoveHeightEventList(ctx, ctx.BlockHeight())
	}
	currentTime := ctx.BlockHeader().Time
	for i := lb.lastBlockTime; i < currentTime; i += 1 {
		if timeEvents := lb.globalManager.GetTimeEventListAtTime(ctx, i); timeEvents != nil {
			lb.executeEvents(ctx, timeEvents.Events)
			lb.globalManager.RemoveTimeEventList(ctx, i)
		}
	}
	lb.lastBlockTime = ctx.BlockHeader().Time
	ABCIValList, err := lb.valManager.GetUpdateValidatorList(ctx)
	if err != nil {
		panic(err)
	}

	return abci.ResponseEndBlock{ValidatorUpdates: ABCIValList}
}

func (lb *LinoBlockchain) executeEvents(ctx sdk.Context, eventList []types.Event) sdk.Error {
	for _, event := range eventList {
		switch e := event.(type) {
		case post.RewardEvent:
			if err := e.Execute(ctx, *lb.postManager, *lb.accountManager, *lb.globalManager); err != nil {
				continue
			}
		}
	}
	return nil
}

func (lb *LinoBlockchain) increaseMinute(ctx sdk.Context) {
	lb.pastMinutes += 1
	if lb.pastMinutes%60 == 0 {
		lb.executeHourlyEvent(ctx)
	}
}

func (lb *LinoBlockchain) executeHourlyEvent(ctx sdk.Context) {
	lb.distributeInflationToValidator(ctx)
}

func (lb *LinoBlockchain) distributeInflationToValidator(ctx sdk.Context) {
	validators, getErr := lb.valManager.GetOncallValList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	coin, err := lb.globalManager.GetValidatorHourlyInflation(ctx, lb.pastMinutes/60)
	if err != nil {
		panic(err)
	}
	// give inflation to each validator evenly
	ratPerValidator := coin.ToRat().Quo(sdk.NewRat(int64(len(validators))))
	for _, validator := range validators {
		lb.accountManager.AddCoin(ctx, validator.Username, types.RatToCoin(ratPerValidator))
	}
}
