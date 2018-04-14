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
	infra "github.com/lino-network/lino/tx/infra"
	"github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/tx/register"
	val "github.com/lino-network/lino/tx/validator"
	vote "github.com/lino-network/lino/tx/vote"
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
	capKeyVoteStore    *sdk.KVStoreKey
	capKeyInfraStore   *sdk.KVStoreKey
	capKeyIBCStore     *sdk.KVStoreKey
	capKeyGlobalStore  *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountManager *acc.AccountManager
	postManager    *post.PostManager
	valManager     *val.ValidatorManager
	globalManager  *global.GlobalManager
	voteManager    *vote.VoteManager
	infraManager   *infra.InfraManager

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
		capKeyVoteStore:    sdk.NewKVStoreKey(types.VoteKVStoreKey),
		capKeyInfraStore:   sdk.NewKVStoreKey(types.InfraKVStoreKey),
		capKeyGlobalStore:  sdk.NewKVStoreKey(types.GlobalKVStoreKey),
		capKeyIBCStore:     sdk.NewKVStoreKey("ibc"),
	}
	lb.accountManager = acc.NewAccountManager(lb.capKeyAccountStore)
	lb.postManager = post.NewPostManager(lb.capKeyPostStore)
	lb.valManager = val.NewValidatorManager(lb.capKeyValStore)
	lb.globalManager = global.NewGlobalManager(lb.capKeyGlobalStore)
	lb.voteManager = vote.NewVoteManager(lb.capKeyVoteStore)
	lb.infraManager = infra.NewInfraManager(lb.capKeyInfraStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(*lb.accountManager)).
		AddRoute(types.AccountRouterName, acc.NewHandler(*lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(*lb.postManager, *lb.accountManager, *lb.globalManager)).
		AddRoute(types.VoteRouterName, vote.NewHandler(*lb.voteManager, *lb.accountManager, *lb.globalManager)).
		AddRoute(types.InfraRouterName, infra.NewHandler(*lb.infraManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(*lb.accountManager, *lb.valManager, *lb.voteManager, *lb.globalManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoreWithDB(lb.capKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	lb.MountStoreWithDB(lb.capKeyPostStore, sdk.StoreTypeIAVL, dbs["post"])
	lb.MountStoreWithDB(lb.capKeyValStore, sdk.StoreTypeIAVL, dbs["val"])
	lb.MountStoreWithDB(lb.capKeyVoteStore, sdk.StoreTypeIAVL, dbs["vote"])
	lb.MountStoreWithDB(lb.capKeyInfraStore, sdk.StoreTypeIAVL, dbs["infra"])
	lb.MountStoreWithDB(lb.capKeyGlobalStore, sdk.StoreTypeIAVL, dbs["global"])
	lb.SetAnteHandler(auth.NewAnteHandler(*lb.accountManager, *lb.globalManager))
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

	if err := lb.voteManager.InitGenesis(ctx); err != nil {
		panic(err)
	}

	if err := lb.infraManager.InitGenesis(ctx); err != nil {
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

	commitingDeposit := types.ValidatorMinCommitingDeposit
	votingDeposit := types.ValidatorMinVotingDeposit
	// withdraw money from validator's bank
	if err := lb.accountManager.MinusCoin(ctx, types.AccountKey(ga.Name), commitingDeposit.Plus(votingDeposit)); err != nil {
		panic(err)
	}

	if addErr := lb.voteManager.AddVoter(ctx, types.AccountKey(ga.Name), votingDeposit); addErr != nil {
		panic(addErr)
	}
	if registerErr := lb.valManager.RegisterValidator(ctx, types.AccountKey(ga.Name), ga.ValPubKey.Bytes(), commitingDeposit); registerErr != nil {
		panic(registerErr)
	}
	if joinErr := lb.valManager.TryBecomeOncallValidator(ctx, types.AccountKey(ga.Name)); joinErr != nil {
		panic(joinErr)
	}
	return nil
}

func (lb *LinoBlockchain) beginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	if lb.lastBlockTime == 0 {
		lb.lastBlockTime = ctx.BlockHeader().Time
	}

	if lb.pastMinutes == 0 {
		lb.pastMinutes = ctx.BlockHeader().Time / 60
	}

	if ctx.BlockHeader().Time/60 > lb.pastMinutes {
		lb.increaseMinute(ctx)
	}
	if err := lb.valManager.SetPreBlockValidators(ctx); err != nil {
		panic(err)
	}
	absentValidators := req.GetAbsentValidators()
	if absentValidators != nil {
		if err := lb.valManager.UpdateAbsentValidator(ctx, absentValidators); err != nil {
			panic(err)
		}
	}

	if err := lb.valManager.FireIncompetentValidator(ctx, req.GetByzantineValidators(), *lb.globalManager); err != nil {
		panic(err)
	}
	return abci.ResponseBeginBlock{}
}

func (lb *LinoBlockchain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	lb.syncValidatorWithVoteManager(ctx)
	lb.executeHeightEvents(ctx)
	lb.executeTimeEvents(ctx)
	lb.punishValidatorsDidntVote(ctx)

	ABCIValList, err := lb.valManager.GetUpdateValidatorList(ctx)
	if err != nil {
		panic(err)
	}

	return abci.ResponseEndBlock{ValidatorUpdates: ABCIValList}
}

func (lb *LinoBlockchain) executeHeightEvents(ctx sdk.Context) {
	if heightEvents := lb.globalManager.GetHeightEventListAtHeight(ctx, ctx.BlockHeight()); heightEvents != nil {
		lb.executeEvents(ctx, heightEvents.Events)
		lb.globalManager.RemoveHeightEventList(ctx, ctx.BlockHeight())
	}
}

func (lb *LinoBlockchain) executeTimeEvents(ctx sdk.Context) {
	currentTime := ctx.BlockHeader().Time
	for i := lb.lastBlockTime; i < currentTime; i += 1 {
		if timeEvents := lb.globalManager.GetTimeEventListAtTime(ctx, i); timeEvents != nil {
			lb.executeEvents(ctx, timeEvents.Events)
			lb.globalManager.RemoveTimeEventList(ctx, i)
		}
	}
	lb.lastBlockTime = ctx.BlockHeader().Time
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
	lb.distributeInflationToInfraProvider(ctx)
}

func (lb *LinoBlockchain) distributeInflationToValidator(ctx sdk.Context) {
	validators, getErr := lb.valManager.GetOncallValidatorList(ctx)
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
		lb.accountManager.AddCoin(ctx, validator, types.RatToCoin(ratPerValidator))
	}
}

func (lb *LinoBlockchain) distributeInflationToInfraProvider(ctx sdk.Context) {
	coin, err := lb.globalManager.GetInfraHourlyInflation(ctx, lb.pastMinutes/60)
	if err != nil {
		panic(err)
	}

	lst, getErr := lb.infraManager.GetInfraProviderList(ctx)
	if getErr != nil {
		panic(getErr)
	}

	for _, provider := range lst.AllInfraProviders {
		percentage, getErr := lb.infraManager.GetUsageWeight(ctx, provider)
		if getErr != nil {
			panic(getErr)
		}
		myShare := coin.ToRat().Mul(percentage)
		lb.accountManager.AddCoin(ctx, provider, types.RatToCoin(myShare))
	}

	if err := lb.infraManager.ClearUsage(ctx); err != nil {
		panic(err)
	}
	
}

func (lb *LinoBlockchain) syncValidatorWithVoteManager(ctx sdk.Context) {
	// tell voting committe the newest validators
	oncallValidators, getErr := lb.valManager.GetOncallValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	lb.voteManager.OncallValidators = oncallValidators

	allValidators, getErr := lb.valManager.GetAllValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	lb.voteManager.AllValidators = allValidators
}

func (lb *LinoBlockchain) punishValidatorsDidntVote(ctx sdk.Context) {
	// punish these validators who didn't vote
	for _, validator := range lb.voteManager.PenaltyValidators {
		if err := lb.valManager.PunishOncallValidator(ctx, validator, types.PenaltyMissVote, *lb.globalManager, false); err != nil {
			panic(err)
		}
	}
	lb.voteManager.PenaltyValidators = lb.voteManager.PenaltyValidators[:0]
}
