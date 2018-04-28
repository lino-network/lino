package app

import (
	"encoding/json"

	"github.com/lino-network/lino/genesis"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/auth"
	developer "github.com/lino-network/lino/tx/developer"
	"github.com/lino-network/lino/tx/global"
	infra "github.com/lino-network/lino/tx/infra"
	"github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/tx/register"
	val "github.com/lino-network/lino/tx/validator"
	vote "github.com/lino-network/lino/tx/vote"
	"github.com/lino-network/lino/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	abci "github.com/tendermint/abci/types"
	oldwire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

const (
	appName = "LinoBlockchain"

	msgTypeRegister          = 0x1
	msgTypeFollow            = 0x2
	msgTypeUnfollow          = 0x3
	msgTypeTransfer          = 0x4
	msgTypePost              = 0x5
	msgTypeLike              = 0x6
	msgTypeDonate            = 0x7
	msgTypeValidatorDeposit  = 0x8
	msgTypeValidatorWithdraw = 0x9
	msgTypeValidatorRevoke   = 0x10
	msgTypeClaim             = 0x11
	msgTypeVoterDeposit      = 0x12
	msgTypeVoterRevoke       = 0x13
	msgTypeVoterWithdraw     = 0x14
	msgTypeDelegate          = 0x15
	msgTypeDelegatorWithdraw = 0x16
	msgTypeRevokeDelegation  = 0x17
	msgTypeVote              = 0x18
	msgTypeCreateProposal    = 0x19
	msgTypeDeveloperRegister = 0x20
	msgTypeDeveloperRevoke   = 0x21
	msgTypeProviderReport    = 0x22

	eventTypeReward     = 0x1
	eventTypeReturnCoin = 0x2
)

// Extended ABCI application
type LinoBlockchain struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the KVStore
	CapKeyMainStore      *sdk.KVStoreKey
	CapKeyAccountStore   *sdk.KVStoreKey
	CapKeyPostStore      *sdk.KVStoreKey
	CapKeyValStore       *sdk.KVStoreKey
	CapKeyVoteStore      *sdk.KVStoreKey
	CapKeyInfraStore     *sdk.KVStoreKey
	CapKeyDeveloperStore *sdk.KVStoreKey
	CapKeyIBCStore       *sdk.KVStoreKey
	CapKeyGlobalStore    *sdk.KVStoreKey

	// Manager for different KVStore
	accountManager   acc.AccountManager
	postManager      post.PostManager
	valManager       val.ValidatorManager
	globalManager    global.GlobalManager
	voteManager      vote.VoteManager
	infraManager     infra.InfraManager
	developerManager developer.DeveloperManager

	// time related
	chainStartTime int64
	lastBlockTime  int64
	pastMinutes    int64
}

func NewLinoBlockchain(logger log.Logger, dbs map[string]dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:              bam.NewBaseApp(appName, logger, dbs["main"]),
		cdc:                  MakeCodec(),
		CapKeyMainStore:      sdk.NewKVStoreKey(types.MainKVStoreKey),
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
		AddRoute(types.RegisterRouterName, register.NewHandler(lb.accountManager)).
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(lb.postManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.VoteRouterName, vote.NewHandler(lb.voteManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.DeveloperRouterName, developer.NewHandler(lb.developerManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.InfraRouterName, infra.NewHandler(lb.infraManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(lb.accountManager, lb.valManager, lb.voteManager, lb.globalManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoreWithDB(lb.CapKeyMainStore, sdk.StoreTypeIAVL, dbs["main"])
	lb.MountStoreWithDB(lb.CapKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	lb.MountStoreWithDB(lb.CapKeyPostStore, sdk.StoreTypeIAVL, dbs["post"])
	lb.MountStoreWithDB(lb.CapKeyValStore, sdk.StoreTypeIAVL, dbs["val"])
	lb.MountStoreWithDB(lb.CapKeyVoteStore, sdk.StoreTypeIAVL, dbs["vote"])
	lb.MountStoreWithDB(lb.CapKeyInfraStore, sdk.StoreTypeIAVL, dbs["infra"])
	lb.MountStoreWithDB(lb.CapKeyDeveloperStore, sdk.StoreTypeIAVL, dbs["developer"])
	lb.MountStoreWithDB(lb.CapKeyGlobalStore, sdk.StoreTypeIAVL, dbs["global"])
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager, lb.globalManager))
	if err := lb.LoadLatestVersion(lb.CapKeyMainStore); err != nil {
		cmn.Exit(err.Error())
	}
	return lb
}

// custom tx codec
// TODO: use new go-wire
func MakeCodec() *wire.Codec {

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
		oldwire.ConcreteType{vote.VoterDepositMsg{}, msgTypeVoterDeposit},
		oldwire.ConcreteType{vote.VoterRevokeMsg{}, msgTypeVoterRevoke},
		oldwire.ConcreteType{vote.VoterWithdrawMsg{}, msgTypeVoterWithdraw},
		oldwire.ConcreteType{vote.DelegateMsg{}, msgTypeDelegate},
		oldwire.ConcreteType{vote.DelegatorWithdrawMsg{}, msgTypeDelegatorWithdraw},
		oldwire.ConcreteType{vote.RevokeDelegationMsg{}, msgTypeRevokeDelegation},
		oldwire.ConcreteType{vote.VoteMsg{}, msgTypeVote},
		oldwire.ConcreteType{vote.CreateProposalMsg{}, msgTypeCreateProposal},
		oldwire.ConcreteType{developer.DeveloperRegisterMsg{}, msgTypeDeveloperRegister},
		oldwire.ConcreteType{developer.DeveloperRevokeMsg{}, msgTypeDeveloperRevoke},
		oldwire.ConcreteType{infra.ProviderReportMsg{}, msgTypeProviderReport},
	)

	var _ = oldwire.RegisterInterface(
		struct{ types.Event }{},
		oldwire.ConcreteType{post.RewardEvent{}, eventTypeReward},
		oldwire.ConcreteType{acc.ReturnCoinEvent{}, eventTypeReturnCoin},
	)

	cdc := wire.NewCodec()

	return cdc
}

// custom logic for transaction decoding
func (lb *LinoBlockchain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx{}

	// StdTx.Msg is an interface.
	err := lb.cdc.UnmarshalJSON(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}

// custom logic for basecoin initialization
func (lb *LinoBlockchain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	genesisState := new(genesis.GenesisState)
	if err := json.Unmarshal(stateJSON, genesisState); err != nil {
		panic(err)
	}

	if err := lb.valManager.InitGenesis(ctx); err != nil {
		panic(err)
	}

	totalCoin, err := types.LinoToCoin(genesisState.TotalLino)
	if err != nil {
		panic(err)
	}
	if err := lb.globalManager.InitGlobalManager(ctx, totalCoin); err != nil {
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
			panic(err)
		}
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
	coin, err := types.LinoToCoin(ga.Lino)
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
	coin, err := types.LinoToCoin(types.LNO(developer.Deposit))
	if err != nil {
		return err
	}

	if err := lb.accountManager.MinusCoin(ctx, types.AccountKey(developer.Name), coin); err != nil {
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

// init process for a block, execute time events and fire incompetent validators
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

	validatorList, getErr := lb.valManager.GetValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	validatorList.PreBlockValidators = validatorList.OncallValidators
	if err := lb.valManager.SetValidatorList(ctx, validatorList); err != nil {
		panic(err)
	}

	absentValidators := req.GetAbsentValidators()
	if absentValidators != nil {
		if err := lb.valManager.UpdateAbsentValidator(ctx, absentValidators); err != nil {
			panic(err)
		}
	}

	if err := lb.valManager.FireIncompetentValidator(ctx, req.GetByzantineValidators(), lb.globalManager); err != nil {
		panic(err)
	}
	lb.syncValidatorWithVoteManager(ctx)
	lb.executeTimeEvents(ctx)
	lb.punishValidatorsDidntVote(ctx)
	return abci.ResponseBeginBlock{}
}

// execute events between last block time and current block time
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

// execute events in list based on their type
func (lb *LinoBlockchain) executeEvents(ctx sdk.Context, eventList []types.Event) sdk.Error {
	for _, event := range eventList {
		switch e := event.(type) {
		case post.RewardEvent:
			if err := e.Execute(
				ctx, lb.postManager, lb.accountManager, lb.globalManager, lb.developerManager); err != nil {
				continue
			}
		case acc.ReturnCoinEvent:
			if err := e.Execute(ctx, lb.accountManager); err != nil {
				continue
			}
		}
	}
	return nil
}

// udpate validator set
func (lb *LinoBlockchain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	ABCIValList, err := lb.valManager.GetUpdateValidatorList(ctx)
	if err != nil {
		panic(err)
	}

	return abci.ResponseEndBlock{ValidatorUpdates: ABCIValList}
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

// execute hourly event, distribute inflation to validators and
// add hourly inflation to content creator reward pool
func (lb *LinoBlockchain) executeHourlyEvent(ctx sdk.Context) {
	if err := lb.globalManager.AddHourlyInflationToRewardPool(
		ctx, (lb.pastMinutes/60)%types.HoursPerYear); err != nil {
		panic(err)
	}
	lb.distributeInflationToValidator(ctx)
}

// execute monthly event, distribute inflation to infra and application
func (lb *LinoBlockchain) executeMonthlyEvent(ctx sdk.Context) {
	lb.distributeInflationToInfraProvider(ctx)
	lb.distributeInflationToDeveloper(ctx)
}

// distribute inflation to validators
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToValidator(ctx sdk.Context) {
	lst, getErr := lb.valManager.GetValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	pastHoursThisYear := (lb.pastMinutes / 60) % types.HoursPerYear
	coin, err := lb.globalManager.GetValidatorHourlyInflation(ctx, pastHoursThisYear)
	if err != nil {
		panic(err)
	}
	// give inflation to each validator evenly
	ratPerValidator := coin.ToRat().Quo(sdk.NewRat(int64(len(lst.OncallValidators))))
	for _, validator := range lst.OncallValidators {
		lb.accountManager.AddCoin(ctx, validator, types.RatToCoin(ratPerValidator))
	}
}

// distribute inflation to infra provider monthly
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToInfraProvider(ctx sdk.Context) {
	pastMonthMinusOneThisYear := (lb.pastMinutes/types.MinutesPerMonth - 1) % 12
	inflation, err := lb.globalManager.GetInfraMonthlyInflation(ctx, pastMonthMinusOneThisYear)
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

// distribute inflation to developer monthly
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToDeveloper(ctx sdk.Context) {
	pastMonthMinusOneThisYear := (lb.pastMinutes/types.MinutesPerMonth - 1) % 12
	inflation, err := lb.globalManager.GetDeveloperMonthlyInflation(ctx, pastMonthMinusOneThisYear)
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

func (lb *LinoBlockchain) syncValidatorWithVoteManager(ctx sdk.Context) {
	// tell voting committe the newest validators
	validatorList, getErr := lb.valManager.GetValidatorList(ctx)
	if getErr != nil {
		panic(getErr)
	}

	referenceList, getErr := lb.voteManager.GetValidatorReferenceList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	referenceList.OncallValidators = validatorList.OncallValidators
	referenceList.AllValidators = validatorList.AllValidators
	if err := lb.voteManager.SetValidatorReferenceList(ctx, referenceList); err != nil {
		panic(err)
	}
}

// validators are required to vote
func (lb *LinoBlockchain) punishValidatorsDidntVote(ctx sdk.Context) {
	lst, getErr := lb.voteManager.GetValidatorReferenceList(ctx)
	if getErr != nil {
		panic(getErr)
	}
	// punish these validators who didn't vote
	for _, validator := range lst.PenaltyValidators {
		if err := lb.valManager.PunishOncallValidator(ctx, validator, types.PenaltyMissVote, lb.globalManager, false); err != nil {
			panic(err)
		}
	}
	lst.PenaltyValidators = lst.PenaltyValidators[:0]
	if err := lb.voteManager.SetValidatorReferenceList(ctx, lst); err != nil {
		panic(err)
	}
}
