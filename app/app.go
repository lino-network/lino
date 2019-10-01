package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	acctypes "github.com/lino-network/lino/x/account/types"
	"github.com/lino-network/lino/x/auth"
	bandwidth "github.com/lino-network/lino/x/bandwidth"
	bandwidthmn "github.com/lino-network/lino/x/bandwidth/manager"
	bandwidthtypes "github.com/lino-network/lino/x/bandwidth/types"
	dev "github.com/lino-network/lino/x/developer"
	devmn "github.com/lino-network/lino/x/developer/manager"
	devtypes "github.com/lino-network/lino/x/developer/types"
	"github.com/lino-network/lino/x/global"
	post "github.com/lino-network/lino/x/post"
	postmn "github.com/lino-network/lino/x/post/manager"
	posttypes "github.com/lino-network/lino/x/post/types"
	price "github.com/lino-network/lino/x/price"
	pricemn "github.com/lino-network/lino/x/price/manager"
	votemn "github.com/lino-network/lino/x/vote/manager"
	votetypes "github.com/lino-network/lino/x/vote/types"

	infra "github.com/lino-network/lino/x/infra"
	inframodel "github.com/lino-network/lino/x/infra/model"
	"github.com/lino-network/lino/x/proposal"

	rep "github.com/lino-network/lino/x/reputation"
	val "github.com/lino-network/lino/x/validator"
	valmn "github.com/lino-network/lino/x/validator/manager"
	valmodel "github.com/lino-network/lino/x/validator/model"
	valtypes "github.com/lino-network/lino/x/validator/types"
	vote "github.com/lino-network/lino/x/vote"
	votemodel "github.com/lino-network/lino/x/vote/model"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cauth "github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	appName = "LinoBlockchain"

	// state files
	prevStateFolder     = "prevstates/"
	currStateFolder     = "currstates/"
	accountStateFile    = "account"
	developerStateFile  = "developer"
	postStateFile       = "post"
	globalStateFile     = "global"
	infraStateFile      = "infra"
	validatorStateFile  = "validator"
	reputationStateFile = "reputation"
	voterStateFile      = "voter"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.linocli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.lino")
)

// LinoBlockchain - Extended ABCI application
type LinoBlockchain struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the KVStore
	CapKeyMainStore         *sdk.KVStoreKey
	CapKeyAccountStore      *sdk.KVStoreKey
	CapKeyPostStore         *sdk.KVStoreKey
	CapKeyValStore          *sdk.KVStoreKey
	CapKeyVoteStore         *sdk.KVStoreKey
	CapKeyInfraStore        *sdk.KVStoreKey
	CapKeyDeveloperStore    *sdk.KVStoreKey
	CapKeyIBCStore          *sdk.KVStoreKey
	CapKeyGlobalStore       *sdk.KVStoreKey
	CapKeyParamStore        *sdk.KVStoreKey
	CapKeyProposalStore     *sdk.KVStoreKey
	CapKeyReputationV2Store *sdk.KVStoreKey
	CapKeyBandwidthStore    *sdk.KVStoreKey
	CapKeyPriceStore        *sdk.KVStoreKey

	// manager for different KVStore
	accountManager    acc.AccountKeeper
	postManager       post.PostKeeper
	valManager        val.ValidatorKeeper
	globalManager     global.GlobalManager
	voteManager       vote.VoteKeeper
	infraManager      infra.InfraManager
	developerManager  dev.DeveloperKeeper
	proposalManager   proposal.ProposalManager
	reputationManager rep.ReputationKeeper
	bandwidthManager  bandwidth.BandwidthKeeper
	priceManager      price.PriceKeeper

	// global param
	paramHolder param.ParamHolder

	// auth
	auth sdk.AnteHandler
}

// NewLinoBlockchain - create a Lino Blockchain instance
func NewLinoBlockchain(
	logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *LinoBlockchain {
	// create your application object
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, types.TxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	var lb = &LinoBlockchain{
		BaseApp:                 bApp,
		cdc:                     cdc,
		CapKeyMainStore:         sdk.NewKVStoreKey(types.MainKVStoreKey),
		CapKeyAccountStore:      sdk.NewKVStoreKey(types.AccountKVStoreKey),
		CapKeyPostStore:         sdk.NewKVStoreKey(types.PostKVStoreKey),
		CapKeyValStore:          sdk.NewKVStoreKey(types.ValidatorKVStoreKey),
		CapKeyVoteStore:         sdk.NewKVStoreKey(types.VoteKVStoreKey),
		CapKeyInfraStore:        sdk.NewKVStoreKey(types.InfraKVStoreKey),
		CapKeyDeveloperStore:    sdk.NewKVStoreKey(types.DeveloperKVStoreKey),
		CapKeyGlobalStore:       sdk.NewKVStoreKey(types.GlobalKVStoreKey),
		CapKeyParamStore:        sdk.NewKVStoreKey(types.ParamKVStoreKey),
		CapKeyProposalStore:     sdk.NewKVStoreKey(types.ProposalKVStoreKey),
		CapKeyReputationV2Store: sdk.NewKVStoreKey(types.ReputationV2KVStoreKey),
		CapKeyBandwidthStore:    sdk.NewKVStoreKey(types.BandwidthKVStoreKey),
		CapKeyPriceStore:        sdk.NewKVStoreKey(types.PriceKVStoreKey),
	}
	// layer-1: basics
	lb.paramHolder = param.NewParamHolder(lb.CapKeyParamStore)
	lb.globalManager = global.NewGlobalManager(lb.CapKeyGlobalStore, lb.paramHolder)
	registerEvent(lb.globalManager.WireCodec())
	lb.accountManager = accmn.NewAccountManager(lb.CapKeyAccountStore, lb.paramHolder, &lb.globalManager)
	lb.reputationManager = rep.NewReputationManager(lb.CapKeyReputationV2Store, lb.paramHolder)
	lb.infraManager = infra.NewInfraManager(lb.CapKeyInfraStore, lb.paramHolder)
	lb.proposalManager = proposal.NewProposalManager(lb.CapKeyProposalStore, lb.paramHolder)

	// layer-2: middlewares
	//// vote <--> validator
	voteManager := votemn.NewVoteManager(lb.CapKeyVoteStore, lb.paramHolder, lb.accountManager, &lb.globalManager)
	lb.valManager = valmn.NewValidatorManager(lb.CapKeyValStore, lb.paramHolder, &voteManager, &lb.globalManager, lb.accountManager)
	lb.voteManager = *voteManager.SetHooks(votemn.NewMultiStakingHooks(lb.valManager.Hooks()))
	//// price -> vote, validator
	lb.priceManager = pricemn.NewWeightedMedianPriceManager(lb.CapKeyPriceStore, lb.valManager, lb.voteManager, &lb.globalManager, lb.paramHolder)

	// layer-3: applications
	lb.developerManager = devmn.NewDeveloperManager(
		lb.CapKeyDeveloperStore, lb.paramHolder,
		&voteManager, lb.accountManager, lb.priceManager, &lb.globalManager)
	lb.postManager = postmn.NewPostManager(
		lb.CapKeyPostStore, lb.accountManager,
		&lb.globalManager, lb.developerManager, lb.reputationManager, lb.priceManager)
	lb.bandwidthManager = bandwidthmn.NewBandwidthManager(
		lb.CapKeyBandwidthStore, lb.paramHolder,
		&lb.globalManager, &voteManager, lb.developerManager, lb.accountManager)
	lb.auth = auth.NewAnteHandler(lb.accountManager, lb.bandwidthManager)

	lb.Router().
		AddRoute(acctypes.RouterKey, acc.NewHandler(lb.accountManager)).
		AddRoute(posttypes.RouterKey, post.NewHandler(lb.postManager)).
		AddRoute(votetypes.RouterKey, vote.NewHandler(lb.voteManager)).
		AddRoute(devtypes.RouterKey, dev.NewHandler(lb.developerManager)).
		AddRoute(proposal.RouterKey, proposal.NewHandler(
			lb.accountManager, lb.proposalManager, lb.postManager, &lb.globalManager, lb.voteManager)).
		AddRoute(infra.RouterKey, infra.NewHandler(lb.infraManager)).
		AddRoute(val.RouterKey, val.NewHandler(lb.valManager))

	lb.QueryRouter().
		AddRoute(acctypes.QuerierRoute, acc.NewQuerier(lb.accountManager)).
		AddRoute(posttypes.QuerierRoute, post.NewQuerier(lb.postManager)).
		AddRoute(votetypes.QuerierRoute, vote.NewQuerier(lb.voteManager)).
		AddRoute(devtypes.QuerierRoute, dev.NewQuerier(lb.developerManager)).
		AddRoute(proposal.QuerierRoute, proposal.NewQuerier(lb.proposalManager)).
		AddRoute(infra.QuerierRoute, infra.NewQuerier(lb.infraManager)).
		AddRoute(val.QuerierRoute, val.NewQuerier(lb.valManager)).
		AddRoute(global.QuerierRoute, global.NewQuerier(lb.globalManager)).
		AddRoute(param.QuerierRoute, param.NewQuerier(lb.paramHolder)).
		AddRoute(bandwidthtypes.QuerierRoute, bandwidth.NewQuerier(lb.bandwidthManager)).
		AddRoute(rep.QuerierRoute, rep.NewQuerier(lb.reputationManager))

	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	lb.SetAnteHandler(lb.auth)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStores(
		lb.CapKeyMainStore, lb.CapKeyAccountStore, lb.CapKeyPostStore, lb.CapKeyValStore,
		lb.CapKeyVoteStore, lb.CapKeyInfraStore, lb.CapKeyDeveloperStore, lb.CapKeyGlobalStore,
		lb.CapKeyParamStore, lb.CapKeyProposalStore, lb.CapKeyReputationV2Store, lb.CapKeyBandwidthStore, lb.CapKeyPriceStore)
	if err := lb.LoadLatestVersion(lb.CapKeyMainStore); err != nil {
		panic(err)
	}

	lb.Seal()

	return lb
}

// MackCodec - codec for application, used by command line tool and authenticate handler
func MakeCodec() *wire.Codec {
	cdc := wire.New()
	cdc.RegisterConcrete(cauth.StdTx{}, "auth/StdTx", nil)
	wire.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	types.RegisterWire(cdc) // types.Msg and types.AddrMsg

	acctypes.RegisterWire(cdc)
	posttypes.RegisterCodec(cdc)
	devtypes.RegisterWire(cdc)
	infra.RegisterWire(cdc)
	votetypes.RegisterWire(cdc)
	valtypes.RegisterCodec(cdc)
	proposal.RegisterWire(cdc)
	registerEvent(cdc)

	cdc.Seal()

	return cdc
}

func registerEvent(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.Event)(nil), nil)
	// TODO(yumin): remove this on upgrade3.
	cdc.RegisterConcrete(posttypes.RewardEventV1{}, "lino/eventReward", nil)
	cdc.RegisterConcrete(posttypes.RewardEvent{}, "lino/eventRewardV2", nil)
	cdc.RegisterConcrete(accmn.ReturnCoinEvent{}, "lino/eventReturn", nil)
	cdc.RegisterConcrete(param.ChangeParamEvent{}, "lino/eventCpe", nil)
	cdc.RegisterConcrete(proposal.DecideProposalEvent{}, "lino/eventDpe", nil)
	cdc.RegisterConcrete(votetypes.UnassignDutyEvent{}, "lino/eventUde", nil)
}

// custom logic for lino blockchain initialization
func (lb *LinoBlockchain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	// set init time to zero
	blockHeader := ctx.BlockHeader()
	blockHeader.Time = time.Unix(0, 0)
	ctx = ctx.WithBlockHeader(blockHeader)

	stateJSON := req.AppStateBytes
	genesisState := new(GenesisState)
	if err := lb.cdc.UnmarshalJSON(stateJSON, genesisState); err != nil {
		panic(err)
	}

	// init parameter holder
	if genesisState.GenesisParam.InitFromConfig {
		if err := lb.paramHolder.InitParamFromConfig(
			ctx,
			genesisState.GenesisParam.GlobalAllocationParam,
			genesisState.GenesisParam.InfraInternalAllocationParam,
			genesisState.GenesisParam.PostParam,
			genesisState.GenesisParam.DeveloperParam,
			genesisState.GenesisParam.ValidatorParam,
			genesisState.GenesisParam.VoteParam,
			genesisState.GenesisParam.ProposalParam,
			genesisState.GenesisParam.CoinDayParam,
			genesisState.GenesisParam.BandwidthParam,
			genesisState.GenesisParam.AccountParam,
			genesisState.GenesisParam.ReputationParam,
			genesisState.GenesisParam.PriceParam,
		); err != nil {
			panic(err)
		}
	} else {
		if err := lb.paramHolder.InitParam(ctx); err != nil {
			panic(err)
		}
	}

	// calculate total lino coin
	totalCoin := types.NewCoinFromInt64(0)
	for _, gacc := range genesisState.Accounts {
		totalCoin = totalCoin.Plus(gacc.Coin)
	}
	totalCoin = totalCoin.Plus(genesisState.ReservePool)
	// global state will then be override if during importing.
	if err := lb.globalManager.InitGlobalManagerWithConfig(
		ctx, totalCoin, genesisState.InitGlobalMeta); err != nil {
		panic(err)
	}

	// set up init state, like empty lists in state.
	if err := lb.developerManager.InitGenesis(ctx, genesisState.ReservePool); err != nil {
		panic(err)
	}
	if err := lb.priceManager.InitGenesis(ctx, genesisState.InitCoinPrice); err != nil {
		panic(err)
	}
	if err := lb.infraManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.proposalManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.bandwidthManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	lb.valManager.InitGenesis(ctx)

	// import from prev state, do not read from genesis.
	if genesisState.LoadPrevStates {
		lb.ImportFromFiles(ctx)
	} else {
		// init genesis accounts
		for _, gacc := range genesisState.Accounts {
			if err := lb.toAppAccount(ctx, gacc); err != nil {
				panic(err)
			}
		}

		// init genesis developers
		for _, developer := range genesisState.Developers {
			if err := lb.toAppDeveloper(ctx, developer); err != nil {
				panic(err)
			}
		}

		// init genesis infra
		for _, infra := range genesisState.Infra {
			if err := lb.toAppInfra(ctx, infra); err != nil {
				panic(err)
			}
		}
	}

	// generate respoinse init message.
	validators, err := lb.valManager.GetInitValidators(ctx)
	if err != nil {
		panic(err)
	}

	return abci.ResponseInitChain{
		ConsensusParams: req.ConsensusParams,
		Validators:      validators,
	}
}

// convert GenesisAccount to AppAccount
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga GenesisAccount) sdk.Error {
	if lb.accountManager.DoesAccountExist(ctx, types.AccountKey(ga.Name)) {
		panic(errors.New("genesis account already exist"))
	}
	if err := lb.accountManager.CreateAccount(
		ctx, types.AccountKey(ga.Name), ga.TransactionKey, ga.ResetKey); err != nil {
		panic(err)
	}
	if err := lb.accountManager.AddCoinToUsername(ctx, types.AccountKey(ga.Name), ga.Coin); err != nil {
		panic(err)
	}

	valParam := lb.paramHolder.GetValidatorParam(ctx)
	if ga.IsValidator {
		if err := lb.voteManager.StakeIn(
			ctx, types.AccountKey(ga.Name), valParam.ValidatorMinDeposit); err != nil {
			panic(err)
		}
		if err := lb.valManager.RegisterValidator(ctx, types.AccountKey(ga.Name), ga.ValPubKey, ""); err != nil {
			panic(err)
		}
	}
	return nil
}

// convert GenesisDeveloper to AppDeveloper
func (lb *LinoBlockchain) toAppDeveloper(
	ctx sdk.Context, developer GenesisAppDeveloper) sdk.Error {
	// TODO(yumin): this is broke. App must first stake then it apply for app.
	// this should be implemented after vote module is ready.
	panic("Unimplemetend genesis to app developer")
	// if !lb.accountManager.DoesAccountExist(ctx, types.AccountKey(developer.Name)) {
	// 	return ErrGenesisFailed("genesis developer account doesn't exist")
	// }

	// if err := lb.accountManager.MinusCoinFromUsername(
	// 	ctx, types.AccountKey(developer.Name), developer.Deposit); err != nil {
	// 	return err
	// }

	// if err := lb.developerManager.RegisterDeveloper(
	// 	ctx, types.AccountKey(developer.Name), developer.Website,
	// 	developer.Description, developer.AppMetaData); err != nil {
	// 	return err
	// }
	// return nil
}

// convert GenesisInfra to AppInfra
func (lb *LinoBlockchain) toAppInfra(
	ctx sdk.Context, infra GenesisInfraProvider) sdk.Error {
	if !lb.accountManager.DoesAccountExist(ctx, types.AccountKey(infra.Name)) {
		return ErrGenesisFailed("genesis infra account doesn't exist")
	}
	if err := lb.infraManager.RegisterInfraProvider(ctx, types.AccountKey(infra.Name)); err != nil {
		return err
	}
	return nil
}

// init process for a block, execute time events and fire incompetent validators
func (lb *LinoBlockchain) beginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	chainStartTime, err := lb.globalManager.GetChainStartTime(ctx)
	if err != nil {
		panic(err)
	}
	if chainStartTime == 0 {
		err := lb.globalManager.SetChainStartTime(ctx, ctx.BlockHeader().Time.Unix())
		if err != nil {
			panic(err)
		}
		err = lb.globalManager.SetLastBlockTime(ctx, ctx.BlockHeader().Time.Unix())
		if err != nil {
			panic(err)
		}
		chainStartTime = ctx.BlockHeader().Time.Unix()
	}

	pastMinutes, err := lb.globalManager.GetPastMinutes(ctx)
	if err != nil {
		panic(err)
	}
	for (ctx.BlockHeader().Time.Unix()-chainStartTime)/60 > pastMinutes {
		lb.increaseMinute(ctx)
		pastMinutes, err = lb.globalManager.GetPastMinutes(ctx)
		if err != nil {
			panic(err)
		}
	}

	global.BeginBlocker(ctx, req, &lb.globalManager)
	bandwidth.BeginBlocker(ctx, req, lb.bandwidthManager)
	val.BeginBlocker(ctx, req, lb.valManager)

	// lb.syncInfoWithVoteManager(ctx)
	lb.executeTimeEvents(ctx)
	return abci.ResponseBeginBlock{}
}

// execute events between last block time and current block time
func (lb *LinoBlockchain) executeTimeEvents(ctx sdk.Context) {
	currentTime := ctx.BlockHeader().Time.Unix()
	lastBlockTime, err := lb.globalManager.GetLastBlockTime(ctx)
	if err != nil {
		panic(err)
	}
	for i := lastBlockTime; i < currentTime; i++ {
		if timeEvents := lb.globalManager.GetTimeEventListAtTime(ctx, i); timeEvents != nil {
			lb.executeEvents(ctx, timeEvents.Events)
			err := lb.globalManager.RemoveTimeEventList(ctx, i)
			if err != nil {
				panic(err)
			}
		}
	}
}

// execute events in list based on their type
func (lb *LinoBlockchain) executeEvents(ctx sdk.Context, eventList []types.Event) {
	for _, event := range eventList {
		switch e := event.(type) {
		case posttypes.RewardEvent:
			if err := lb.postManager.ExecRewardEvent(ctx, e); err != nil {
				panic(err)
			}
		case accmn.ReturnCoinEvent:
			if err := e.Execute(ctx, lb.accountManager.(accmn.AccountManager)); err != nil {
				panic(err)
			}
		case proposal.DecideProposalEvent:
			if err := e.Execute(
				ctx, lb.voteManager, lb.valManager, lb.accountManager, lb.proposalManager,
				lb.postManager, &lb.globalManager); err != nil {
				panic(err)
			}
		case param.ChangeParamEvent:
			if err := e.Execute(ctx, lb.paramHolder); err != nil {
				panic(err)
			}
		case votetypes.UnassignDutyEvent:
			if err := lb.voteManager.ExecUnassignDutyEvent(ctx, e); err != nil {
				panic(err)
			}
		default:
			ctx.Logger().Error(fmt.Sprintf("skipping event: %+v", e))
		}
	}
}

// udpate validator set and renew reputation round
func (lb *LinoBlockchain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	rep.EndBlocker(ctx, req, lb.reputationManager)
	global.EndBlocker(ctx, req, &lb.globalManager)
	bandwidth.EndBlocker(ctx, req, lb.bandwidthManager)

	// update last block time
	if err := lb.globalManager.SetLastBlockTime(ctx, ctx.BlockHeader().Time.Unix()); err != nil {
		panic(err)
	}
	// update validator set.
	validatorUpdates, err := lb.valManager.GetValidatorUpdates(ctx)
	if err != nil {
		panic(err)
	}

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
	}
}

func (lb *LinoBlockchain) increaseMinute(ctx sdk.Context) {
	pastMinutes, err := lb.globalManager.GetPastMinutes(ctx)
	if err != nil {
		panic(err)
	}
	pastMinutes++
	if err := lb.globalManager.SetPastMinutes(ctx, pastMinutes); err != nil {
		panic(err)
	}
	if pastMinutes%60 == 0 {
		lb.executeHourlyEvent(ctx)
	}
	if pastMinutes%types.MinutesPerDay == 0 {
		lb.executeDailyEvent(ctx)
	}
	if pastMinutes%types.MinutesPerMonth == 0 {
		lb.executeMonthlyEvent(ctx)
	}
	// 	if pastMinutes%types.MinutesPerYear == 0 {
	// 		lb.executeAnnuallyEvent(ctx)
	// 	}
}

// execute hourly event, distribute inflation to validators and
// add hourly inflation to content creator reward pool
func (lb *LinoBlockchain) executeHourlyEvent(ctx sdk.Context) {
	if err := lb.globalManager.DistributeHourlyInflation(ctx); err != nil {
		panic(err)
	}
	if err := lb.valManager.DistributeInflationToValidator(ctx); err != nil {
		panic(err)
	}
	if err := lb.bandwidthManager.ReCalculateAppBandwidthInfo(ctx); err != nil {
		panic(err)
	}
	if err := lb.priceManager.UpdatePrice(ctx); err != nil {
		panic(err)
	}
}

// execute daily event, record consumption friction and lino power
func (lb *LinoBlockchain) executeDailyEvent(ctx sdk.Context) {
	err := lb.globalManager.RecordConsumptionAndLinoStake(ctx)
	if err != nil {
		panic(err)
	}
	err = lb.bandwidthManager.DecayMaxMPS(ctx)
	if err != nil {
		panic(err)
	}
}

// execute monthly event, distribute inflation to infra and application
func (lb *LinoBlockchain) executeMonthlyEvent(ctx sdk.Context) {
	lb.distributeInflationToInfraProvider(ctx)
	// distributeInflationToDeveloper
	err := lb.developerManager.DistributeDevInflation(ctx)
	if err != nil {
		panic(err)
	}

}

func (lb *LinoBlockchain) executeAnnuallyEvent(ctx sdk.Context) { //nolint:unused
	if err := lb.globalManager.SetTotalLinoAndRecalculateGrowthRate(ctx); err != nil {
		panic(err)
	}
}

// distribute inflation to infra provider monthly
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToInfraProvider(ctx sdk.Context) {
	inflation, err := lb.globalManager.GetInfraMonthlyInflation(ctx)
	if err != nil {
		panic(err)
	}

	lst, err := lb.infraManager.GetInfraProviderList(ctx)
	if err != nil {
		panic(err)
	}
	totalDistributedInflation := types.NewCoinFromInt64(0)
	for idx, provider := range lst.AllInfraProviders {
		if idx == (len(lst.AllInfraProviders) - 1) {
			err := lb.accountManager.AddCoinToUsername(ctx, provider, inflation.Minus(totalDistributedInflation))
			if err != nil {
				panic(fmt.Errorf("%s: %s", err, provider))
			}
			break
		}
		percentage, err := lb.infraManager.GetUsageWeight(ctx, provider)
		if err != nil {
			panic(err)
		}
		myShareRat := inflation.ToDec().Mul(percentage)
		myShareCoin := types.DecToCoin(myShareRat)
		totalDistributedInflation = totalDistributedInflation.Plus(myShareCoin)
		err = lb.accountManager.AddCoinToUsername(ctx, provider, myShareCoin)
		if err != nil {
			panic(err)
		}
	}
	if err := lb.infraManager.ClearUsage(ctx); err != nil {
		panic(err)
	}
}

// func (lb *LinoBlockchain) syncInfoWithVoteManager(ctx sdk.Context) {
// 	// tell voting committee the newest validators
// 	validatorList, err := lb.valManager.GetValidatorList(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	referenceList, err := lb.voteManager.GetValidatorReferenceList(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	referenceList.AllValidators = validatorList.AllValidators
// 	if err := lb.voteManager.SetValidatorReferenceList(ctx, referenceList); err != nil {
// 		panic(err)
// 	}
// }

// Custom logic for state export
func (lb *LinoBlockchain) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := lb.NewContext(true, abci.Header{})

	exportPath := DefaultNodeHome + "/" + currStateFolder
	err = os.MkdirAll(exportPath, os.ModePerm)
	if err != nil {
		panic("failed to create export dir due to: " + err.Error())
	}

	// exportToFile := func(filename string, exporter func(sdk.Context) interface{}) {
	// 	f, err := os.Create(exportPath + filename)
	// 	if err != nil {
	// 		panic("failed to create account")
	// 	}
	// 	defer f.Close()
	// 	jsonbytes, err := lb.cdc.MarshalJSON(exporter(ctx))
	// 	if err != nil {
	// 		panic("failed to marshal json for " + filename + " due to " + err.Error())
	// 	}
	// 	_, err = f.Write(jsonbytes)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("export for %s done: %d bytes\n", filename, len(jsonbytes))
	// 	err = f.Sync()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// TODO(yumin): accountStateFile
	// exportToFile(accountStateFile, func(ctx sdk.Context) interface{} {
	// 	return lb.accountManager.Export(ctx).ToIR()
	// })
	// TODO(yumin): dev export is not implemented yet.
	// exportToFile(developerStateFile, func(ctx sdk.Context) interface{} {
	// 	return lb.developerManager.Export(ctx).ToIR()
	// })
	// TODO(yumin): post export is not implemented yet.
	// exportToFile(postStateFile, func(ctx sdk.Context) interface{} {
	// 	return lb.postManager.Export(ctx).ToIR()
	// })
	// TODO(yumin): global export is not implemented yet.
	// exportToFile(globalStateFile, func(ctx sdk.Context) interface{} {
	// 	return lb.globalManager.Export(ctx).ToIR()
	// })
	// exportToFile(infraStateFile, func(ctx sdk.Context) interface{} {
	// 	return lb.infraManager.Export(ctx).ToIR()
	// })
	// exportToFile(validatorStateFile, func(ctx sdk.Context) interface{} {
	// 	return lb.valManager.Export(ctx).ToIR()
	// })
	// exportToFile(voterStateFile, func(ctx sdk.Context) interface{} {
	// 	return lb.voteManager.Export(ctx).ToIR()
	// })
	err = lb.reputationManager.ExportToFile(ctx, exportPath+"reputation")
	if err != nil {
		panic(err)
	}

	genesisState := GenesisState{}

	appState, err = wire.MarshalJSONIndent(lb.cdc, genesisState)
	if err != nil {
		return nil, nil, err
	}
	return appState, validators, nil
}

// ImportFromFiles Custom logic for state export
func (lb *LinoBlockchain) ImportFromFiles(ctx sdk.Context) {
	check := func(err error) {
		if err != nil {
			panic("failed to unmarshal " + err.Error())
		}
	}

	prevStateDir := DefaultNodeHome + "/" + prevStateFolder
	// import account
	err := lb.accountManager.ImportFromFile(
		ctx, lb.cdc, prevStateDir+accountStateFile)
	if err != nil {
		panic(err)
	}
	// import post.
	err = lb.postManager.ImportFromFile(
		ctx, lb.cdc, prevStateDir+postStateFile)
	if err != nil {
		panic(err)
	}
	// import dev
	err = lb.developerManager.ImportFromFile(
		ctx, lb.cdc, prevStateDir+developerStateFile)
	if err != nil {
		panic(err)
	}
	// import global
	err = lb.globalManager.ImportFromFile(
		ctx, lb.cdc, prevStateDir+globalStateFile)
	if err != nil {
		panic(err)
	}
	// import reputation
	err = lb.reputationManager.ImportFromFile(
		ctx, DefaultNodeHome+"/"+prevStateFolder+reputationStateFile)
	if err != nil {
		panic(err)
	}

	importFromFile := func(filename string, tables interface{}) {
		// XXX(yumin): does not support customized node home import.
		f, err := os.Open(DefaultNodeHome + "/" + prevStateFolder + filename)
		if err != nil {
			panic("failed to open " + err.Error())
		}
		defer f.Close()
		bytes, err := ioutil.ReadAll(f)
		check(err)
		// XXX(yumin): ugly, trying found a better way.
		switch t := tables.(type) {
		case *inframodel.InfraTablesIR:
			err = lb.cdc.UnmarshalJSON(bytes, t)
			check(err)
			fmt.Printf("%s state parsed: %T\n", filename, t)
			lb.infraManager.Import(ctx, t)
		case *valmodel.ValidatorTablesIR:
			err = lb.cdc.UnmarshalJSON(bytes, t)
			check(err)
			fmt.Printf("%s state parsed: %T\n", filename, t)
			// lb.valManager.Import(ctx, t)
		// case *votemodel.VoterTablesIR:
		// 	err = lb.cdc.UnmarshalJSON(bytes, t)
		// 	check(err)
		// 	fmt.Printf("%s state parsed: %T\n", filename, t)
		// 	lb.voteManager.Import(ctx, t)
		default:
			panic(fmt.Sprintf("Unknown import type: %T", t))
		}
		fmt.Printf("%s loaded, total %d bytes\n", filename, len(bytes))
	}

	importFromFile(infraStateFile, &inframodel.InfraTablesIR{})
	importFromFile(validatorStateFile, &valmodel.ValidatorTablesIR{})
	importFromFile(voterStateFile, &votemodel.VoterTablesIR{})
}
