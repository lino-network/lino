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
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/auth"
	developer "github.com/lino-network/lino/tx/developer"
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
	CapKeyAccountStore   *sdk.KVStoreKey
	CapKeyPostStore      *sdk.KVStoreKey
	CapKeyValStore       *sdk.KVStoreKey
	CapKeyVoteStore      *sdk.KVStoreKey
	CapKeyInfraStore     *sdk.KVStoreKey
	CapKeyDeveloperStore *sdk.KVStoreKey
	CapKeyIBCStore       *sdk.KVStoreKey
	CapKeyGlobalStore    *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountManager   *acc.AccountManager
	postManager      *post.PostManager
	valManager       *val.ValidatorManager
	globalManager    *global.GlobalManager
	voteManager      *vote.VoteManager
	infraManager     *infra.InfraManager
	developerManager *developer.DeveloperManager

	chainStartTime int64
	lastBlockTime  int64
	// for recurring time based event
	pastMinutes int64
}

func NewLinoBlockchain(logger log.Logger, dbs map[string]dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:              bam.NewBaseApp(appName, logger, dbs["acc"]),
		cdc:                  MakeCodec(),
		CapKeyAccountStore:   sdk.NewKVStoreKey(types.AccountKVStoreKey),
		CapKeyPostStore:      sdk.NewKVStoreKey(types.PostKVStoreKey),
		CapKeyValStore:       sdk.NewKVStoreKey(types.ValidatorKVStoreKey),
		CapKeyVoteStore:      sdk.NewKVStoreKey(types.VoteKVStoreKey),
		CapKeyInfraStore:     sdk.NewKVStoreKey(types.InfraKVStoreKey),
		CapKeyDeveloperStore: sdk.NewKVStoreKey(types.DeveloperKVStoreKey),
		CapKeyGlobalStore:    sdk.NewKVStoreKey(types.GlobalKVStoreKey),
	}
	lb.accountManager = acc.NewAccountManager(lb.CapKeyAccountStore)
	lb.postManager = post.NewPostManager(lb.CapKeyPostStore)
	lb.valManager = val.NewValidatorManager(lb.CapKeyValStore)
	lb.globalManager = global.NewGlobalManager(lb.CapKeyGlobalStore)
	lb.voteManager = vote.NewVoteManager(lb.CapKeyVoteStore)
	lb.infraManager = infra.NewInfraManager(lb.CapKeyInfraStore)
	lb.developerManager = developer.NewDeveloperManager(lb.CapKeyDeveloperStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(*lb.accountManager)).
		AddRoute(types.AccountRouterName, acc.NewHandler(*lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(*lb.postManager, *lb.accountManager, *lb.globalManager)).
		AddRoute(types.VoteRouterName, vote.NewHandler(*lb.voteManager, *lb.accountManager, *lb.globalManager)).
		AddRoute(types.DeveloperRouterName, developer.NewHandler(*lb.developerManager, *lb.accountManager, *lb.globalManager)).
		AddRoute(types.InfraRouterName, infra.NewHandler(*lb.infraManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(*lb.accountManager, *lb.valManager, *lb.voteManager, *lb.globalManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoreWithDB(lb.CapKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	lb.MountStoreWithDB(lb.CapKeyPostStore, sdk.StoreTypeIAVL, dbs["post"])
	lb.MountStoreWithDB(lb.CapKeyValStore, sdk.StoreTypeIAVL, dbs["val"])
	lb.MountStoreWithDB(lb.CapKeyVoteStore, sdk.StoreTypeIAVL, dbs["vote"])
	lb.MountStoreWithDB(lb.CapKeyInfraStore, sdk.StoreTypeIAVL, dbs["infra"])
	lb.MountStoreWithDB(lb.CapKeyDeveloperStore, sdk.StoreTypeIAVL, dbs["developer"])
	lb.MountStoreWithDB(lb.CapKeyGlobalStore, sdk.StoreTypeIAVL, dbs["global"])
	lb.SetAnteHandler(auth.NewAnteHandler(*lb.accountManager, *lb.globalManager))
	if err := lb.LoadLatestVersion(lb.CapKeyAccountStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.CapKeyPostStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.CapKeyValStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.CapKeyGlobalStore); err != nil {
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

	const eventTypeReward = 0x1
	const eventTypeReturnCoin = 0x2

	var _ = oldwire.RegisterInterface(
		struct{ types.Event }{},
		oldwire.ConcreteType{post.RewardEvent{}, eventTypeReward},
		oldwire.ConcreteType{acc.ReturnCoinEvent{}, eventTypeReturnCoin},
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
	if err := lb.globalManager.InitGlobalManager(
		ctx, types.NewCoin(genesisState.TotalLino*types.Decimals)); err != nil {
		panic(err)
	}
	if err := lb.developerManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.infraManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.voteManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	for _, gacc := range genesisState.Accounts {
		if err := lb.toAppAccount(ctx, gacc); err != nil {
			panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		// lb.accountMapper.SetAccount(ctx, acc)
	}

	for _, developer := range genesisState.Developers {
		if err := lb.toAppDeveloper(ctx, developer); err != nil {
			panic(err)
		}
	}

	for _, infra := range genesisState.Infra {
		if err := lb.toAppInfra(ctx, infra); err != nil {
			panic(err)
		}
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
	if ga.IsValidator {
		commitingDeposit := types.ValidatorMinCommitingDeposit
		votingDeposit := types.ValidatorMinVotingDeposit
		// withdraw money from validator's bank
		if err := lb.accountManager.MinusCoin(
			ctx, types.AccountKey(ga.Name), commitingDeposit.Plus(votingDeposit)); err != nil {
			panic(err)
		}

		if addErr := lb.voteManager.AddVoter(ctx, types.AccountKey(ga.Name), votingDeposit); addErr != nil {
			panic(addErr)
		}
		if registerErr := lb.valManager.RegisterValidator(
			ctx, types.AccountKey(ga.Name), ga.ValPubKey.Bytes(), commitingDeposit); registerErr != nil {
			panic(registerErr)
		}
		if joinErr := lb.valManager.TryBecomeOncallValidator(ctx, types.AccountKey(ga.Name)); joinErr != nil {
			panic(joinErr)
		}
	}
	return nil
}

// convert GenesisDeveloper to AppDeveloper
func (lb *LinoBlockchain) toAppDeveloper(
	ctx sdk.Context, developer genesis.GenesisAppDeveloper) sdk.Error {
	if !lb.accountManager.IsAccountExist(ctx, types.AccountKey(developer.Name)) {
		return sdk.ErrGenesisParse("genesis developer account doesn't exist")
	}
	coin, err := types.LinoToCoin(types.LNO(sdk.NewRat(developer.Deposit)))
	if err != nil {
		return err
	}
	if err := lb.developerManager.RegisterDeveloper(ctx, types.AccountKey(developer.Name), coin); err != nil {
		return err
	}
	return nil
}

// convert GenesisInfra to AppInfra
func (lb *LinoBlockchain) toAppInfra(
	ctx sdk.Context, infra genesis.GenesisInfraProvider) sdk.Error {
	if !lb.accountManager.IsAccountExist(ctx, types.AccountKey(infra.Name)) {
		return sdk.ErrGenesisParse("genesis infra account doesn't exist")
	}
	if err := lb.infraManager.RegisterInfraProvider(ctx, types.AccountKey(infra.Name)); err != nil {
		return err
	}
	return nil
}

func (lb *LinoBlockchain) beginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	if lb.chainStartTime == 0 {
		lb.chainStartTime = ctx.BlockHeader().Time
		lb.lastBlockTime = ctx.BlockHeader().Time
	}
	if err := lb.globalManager.UpdateTPS(ctx, lb.lastBlockTime); err != nil {
		panic(err)
	}
	for (ctx.BlockHeader().Time-lb.chainStartTime)/60 > lb.pastMinutes {
		lb.increaseMinute(ctx)
	}

	preBlockValidators, getErr := lb.valManager.GetOncallValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	ctx = val.WithPreBlockValidators(ctx, preBlockValidators)
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
	ctx = lb.syncValidatorWithVoteManager(ctx)
	lb.executeTimeEvents(ctx)
	lb.punishValidatorsDidntVote(ctx)

	ABCIValList, err := lb.valManager.GetUpdateValidatorList(ctx)
	if err != nil {
		panic(err)
	}

	return abci.ResponseEndBlock{ValidatorUpdates: ABCIValList}
}

func (lb *LinoBlockchain) executeTimeEvents(ctx sdk.Context) {
	currentTime := ctx.BlockHeader().Time
	for i := lb.lastBlockTime; i < currentTime; i += 1 {
		if timeEvents := lb.globalManager.GetTimeEventListAtTime(ctx, i); timeEvents != nil {
			fmt.Println("execute time event:", i)
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
			if err := e.Execute(ctx, *lb.postManager, *lb.accountManager, *lb.globalManager, *lb.developerManager); err != nil {
				continue
			}
		case acc.ReturnCoinEvent:
			if err := e.Execute(ctx, *lb.accountManager); err != nil {
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
	if lb.pastMinutes%types.MinutesPerMonth == 0 {
		lb.executeMonthlyEvent(ctx)
	}
}

func (lb *LinoBlockchain) executeHourlyEvent(ctx sdk.Context) {
	if err := lb.globalManager.AddHourlyInflationToRewardPool(
		ctx, (lb.pastMinutes/60)%types.HoursPerYear); err != nil {
		panic(err)
	}
	lb.distributeInflationToValidator(ctx)
}

func (lb *LinoBlockchain) executeMonthlyEvent(ctx sdk.Context) {
	lb.distributeInflationToInfraProvider(ctx)
	lb.distributeInflationToDeveloper(ctx)
}

func (lb *LinoBlockchain) distributeInflationToValidator(ctx sdk.Context) {
	validators, getErr := lb.valManager.GetOncallValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	coin, err := lb.globalManager.GetValidatorHourlyInflation(ctx, (lb.pastMinutes/60)%types.HoursPerYear)
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
	inflation, err := lb.globalManager.GetInfraMonthlyInflation(ctx, (lb.pastMinutes/types.MinutesPerMonth-1)%12)
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
		myShare := inflation.ToRat().Mul(percentage)
		lb.accountManager.AddCoin(ctx, provider, types.RatToCoin(myShare))
	}

	if err := lb.infraManager.ClearUsage(ctx); err != nil {
		panic(err)
	}
}

func (lb *LinoBlockchain) distributeInflationToDeveloper(ctx sdk.Context) {
	inflation, err := lb.globalManager.GetDeveloperMonthlyInflation(ctx, (lb.pastMinutes/types.MinutesPerMonth-1)%12)
	if err != nil {
		panic(err)
	}

	lst, getErr := lb.developerManager.GetDeveloperList(ctx)
	if getErr != nil {
		panic(getErr)
	}

	for _, developer := range lst.AllDevelopers {
		percentage, getErr := lb.developerManager.GetConsumptionWeight(ctx, developer)
		if getErr != nil {
			panic(getErr)
		}
		myShare := inflation.ToRat().Mul(percentage)
		lb.accountManager.AddCoin(ctx, developer, types.RatToCoin(myShare))
	}

	if err := lb.developerManager.ClearConsumption(ctx); err != nil {
		panic(err)
	}
}

func (lb *LinoBlockchain) syncValidatorWithVoteManager(ctx sdk.Context) sdk.Context {
	// tell voting committe the newest validators
	oncallValidators, getErr := lb.valManager.GetOncallValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	ctx = vote.WithOncallValidators(ctx, oncallValidators)

	allValidators, getErr := lb.valManager.GetAllValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	ctx = vote.WithAllValidators(ctx, allValidators)
	return ctx
}

func (lb *LinoBlockchain) punishValidatorsDidntVote(ctx sdk.Context) {
	lst, getErr := lb.voteManager.GetValidatorPenaltyList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	// punish these validators who didn't vote
	for _, validator := range lst.Validators {
		if err := lb.valManager.PunishOncallValidator(ctx, validator, types.PenaltyMissVote, *lb.globalManager, false); err != nil {
			panic(err)
		}
	}
	lst.Validators = lst.Validators[:0]
	if err := lb.voteManager.SetValidatorPenaltyList(ctx, lst); err != nil {
		panic(err)
	}
}
