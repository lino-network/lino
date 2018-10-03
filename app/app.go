package app

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/recorder"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/auth"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post"
	"github.com/lino-network/lino/x/proposal"

	acc "github.com/lino-network/lino/x/account"
	accModel "github.com/lino-network/lino/x/account/model"
	developer "github.com/lino-network/lino/x/developer"
	infra "github.com/lino-network/lino/x/infra"
	rep "github.com/lino-network/lino/x/reputation"
	val "github.com/lino-network/lino/x/validator"
	vote "github.com/lino-network/lino/x/vote"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cauth "github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	appName = "LinoBlockchain"
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
	CapKeyMainStore       *sdk.KVStoreKey
	CapKeyAccountStore    *sdk.KVStoreKey
	CapKeyPostStore       *sdk.KVStoreKey
	CapKeyValStore        *sdk.KVStoreKey
	CapKeyVoteStore       *sdk.KVStoreKey
	CapKeyInfraStore      *sdk.KVStoreKey
	CapKeyDeveloperStore  *sdk.KVStoreKey
	CapKeyIBCStore        *sdk.KVStoreKey
	CapKeyGlobalStore     *sdk.KVStoreKey
	CapKeyParamStore      *sdk.KVStoreKey
	CapKeyProposalStore   *sdk.KVStoreKey
	CapKeyReputationStore *sdk.KVStoreKey

	// manager for different KVStore
	accountManager    acc.AccountManager
	postManager       post.PostManager
	valManager        val.ValidatorManager
	globalManager     global.GlobalManager
	voteManager       vote.VoteManager
	infraManager      infra.InfraManager
	developerManager  developer.DeveloperManager
	proposalManager   proposal.ProposalManager
	reputationManager rep.ReputationManager

	// global param
	paramHolder param.ParamHolder

	// recorder
	recorder recorder.Recorder
}

// NewLinoBlockchain - create a Lino Blockchain instance
func NewLinoBlockchain(
	logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *LinoBlockchain {
	// create your application object
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	var lb = &LinoBlockchain{
		BaseApp:               bApp,
		cdc:                   cdc,
		CapKeyMainStore:       sdk.NewKVStoreKey(types.MainKVStoreKey),
		CapKeyAccountStore:    sdk.NewKVStoreKey(types.AccountKVStoreKey),
		CapKeyPostStore:       sdk.NewKVStoreKey(types.PostKVStoreKey),
		CapKeyValStore:        sdk.NewKVStoreKey(types.ValidatorKVStoreKey),
		CapKeyVoteStore:       sdk.NewKVStoreKey(types.VoteKVStoreKey),
		CapKeyInfraStore:      sdk.NewKVStoreKey(types.InfraKVStoreKey),
		CapKeyDeveloperStore:  sdk.NewKVStoreKey(types.DeveloperKVStoreKey),
		CapKeyGlobalStore:     sdk.NewKVStoreKey(types.GlobalKVStoreKey),
		CapKeyParamStore:      sdk.NewKVStoreKey(types.ParamKVStoreKey),
		CapKeyProposalStore:   sdk.NewKVStoreKey(types.ProposalKVStoreKey),
		CapKeyReputationStore: sdk.NewKVStoreKey(types.ReputationKVStoreKey),
	}
	lb.recorder = recorder.NewRecorder()
	lb.paramHolder = param.NewParamHolder(lb.CapKeyParamStore)
	lb.accountManager = acc.NewAccountManager(lb.CapKeyAccountStore, lb.paramHolder)
	lb.postManager = post.NewPostManager(lb.CapKeyPostStore, lb.paramHolder, lb.recorder)
	lb.valManager = val.NewValidatorManager(lb.CapKeyValStore, lb.paramHolder)
	lb.globalManager = global.NewGlobalManager(lb.CapKeyGlobalStore, lb.paramHolder)
	registerEvent(lb.globalManager.WireCodec())

	lb.reputationManager = rep.NewReputationManager(lb.CapKeyReputationStore, lb.paramHolder)
	lb.voteManager = vote.NewVoteManager(lb.CapKeyVoteStore, lb.paramHolder)
	lb.infraManager = infra.NewInfraManager(lb.CapKeyInfraStore, lb.paramHolder)
	lb.developerManager = developer.NewDeveloperManager(lb.CapKeyDeveloperStore, lb.paramHolder)
	lb.proposalManager = proposal.NewProposalManager(lb.CapKeyProposalStore, lb.paramHolder)

	lb.Router().
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager, lb.globalManager)).
		AddRoute(types.PostRouterName, post.NewHandler(
			lb.postManager, lb.accountManager, lb.globalManager, lb.developerManager, lb.reputationManager)).
		AddRoute(types.VoteRouterName, vote.NewHandler(
			lb.voteManager, lb.accountManager, lb.globalManager, lb.reputationManager)).
		AddRoute(types.DeveloperRouterName, developer.NewHandler(
			lb.developerManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.ProposalRouterName, proposal.NewHandler(
			lb.accountManager, lb.proposalManager, lb.postManager, lb.globalManager, lb.voteManager)).
		AddRoute(types.InfraRouterName, infra.NewHandler(lb.infraManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(
			lb.accountManager, lb.valManager, lb.voteManager, lb.globalManager))

	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager, lb.globalManager))
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoresIAVL(
		lb.CapKeyMainStore, lb.CapKeyAccountStore, lb.CapKeyPostStore, lb.CapKeyValStore,
		lb.CapKeyVoteStore, lb.CapKeyInfraStore, lb.CapKeyDeveloperStore, lb.CapKeyGlobalStore,
		lb.CapKeyParamStore, lb.CapKeyProposalStore, lb.CapKeyReputationStore)
	if err := lb.LoadLatestVersion(lb.CapKeyMainStore); err != nil {
		cmn.Exit(err.Error())
	}

	lb.Seal()

	return lb
}

// DefaultTxDecoder - default tx decoder, decode tx before authenticate handler
func DefaultTxDecoder(cdc *wire.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (tx sdk.Tx, err sdk.Error) {
		defer func() {
			if r := recover(); r != nil {
				err = sdk.ErrTxDecode("tx decode panic")
			}
		}()
		tx = cauth.StdTx{}

		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}

		// StdTx.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		unmarshalErr := cdc.UnmarshalJSON(txBytes, &tx)
		if unmarshalErr != nil {
			return nil, sdk.ErrTxDecode("")
		}
		return tx, nil
	}
}

// MackCodec - codec for application, used by command line tool and authenticate handler
func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()
	cdc.RegisterConcrete(cauth.StdTx{}, "auth/StdTx", nil)
	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)

	acc.RegisterWire(cdc)
	post.RegisterWire(cdc)
	developer.RegisterWire(cdc)
	infra.RegisterWire(cdc)
	vote.RegisterWire(cdc)
	val.RegisterWire(cdc)
	proposal.RegisterWire(cdc)

	cdc.Seal()

	return cdc
}

func registerEvent(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(post.RewardEvent{}, "lino/eventReward", nil)
	cdc.RegisterConcrete(acc.ReturnCoinEvent{}, "lino/eventReturn", nil)
	cdc.RegisterConcrete(param.ChangeParamEvent{}, "lino/eventCpe", nil)
	cdc.RegisterConcrete(proposal.DecideProposalEvent{}, "lino/eventDpe", nil)
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
			genesisState.GenesisParam.EvaluateOfContentValueParam,
			genesisState.GenesisParam.DeveloperParam,
			genesisState.GenesisParam.ValidatorParam,
			genesisState.GenesisParam.VoteParam,
			genesisState.GenesisParam.ProposalParam,
			genesisState.GenesisParam.CoinDayParam,
			genesisState.GenesisParam.BandwidthParam,
			genesisState.GenesisParam.AccountParam,
			genesisState.GenesisParam.ReputationParam); err != nil {
			panic(err)
		}
	} else {
		if err := lb.paramHolder.InitParam(ctx); err != nil {
			panic(err)
		}
	}

	totalCoin := types.NewCoinFromInt64(0)

	// calculate total lino coin
	for _, gacc := range genesisState.Accounts {
		totalCoin = totalCoin.Plus(gacc.Coin)
	}
	if err := lb.globalManager.InitGlobalManagerWithConfig(
		ctx, totalCoin, genesisState.InitGlobalMeta); err != nil {
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
	if err := lb.proposalManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.valManager.InitGenesis(ctx); err != nil {
		panic(err)
	}

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
	return abci.ResponseInitChain{}
}

// convert GenesisAccount to AppAccount
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga GenesisAccount) sdk.Error {
	if lb.accountManager.DoesAccountExist(ctx, types.AccountKey(ga.Name)) {
		panic(errors.New("genesis account already exist"))
	}
	if err := lb.accountManager.CreateAccount(
		ctx, types.AccountKey(ga.Name), types.AccountKey(ga.Name),
		ga.ResetKey, ga.TransactionKey, ga.AppKey, ga.Coin); err != nil {
		panic(err)
	}

	valParam, err := lb.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		panic(err)
	}

	if ga.IsValidator {
		// withdraw money from validator's bank
		if err := lb.accountManager.MinusSavingCoin(
			ctx, types.AccountKey(ga.Name),
			valParam.ValidatorMinCommittingDeposit.Plus(valParam.ValidatorMinVotingDeposit),
			"", "", types.ValidatorDeposit); err != nil {
			panic(err)
		}
		if err := vote.AddStake(
			ctx, types.AccountKey(ga.Name), valParam.ValidatorMinVotingDeposit,
			lb.voteManager, lb.globalManager, lb.accountManager,
			lb.reputationManager); err != nil {
			panic(err)
		}
		if err := lb.voteManager.AddVoter(
			ctx, types.AccountKey(ga.Name), valParam.ValidatorMinVotingDeposit); err != nil {
			panic(err)
		}
		if err := lb.valManager.RegisterValidator(
			ctx, types.AccountKey(ga.Name), ga.ValPubKey,
			valParam.ValidatorMinCommittingDeposit, ""); err != nil {
			panic(err)
		}
		if err := lb.valManager.TryBecomeOncallValidator(ctx, types.AccountKey(ga.Name)); err != nil {
			panic(err)
		}
	}
	return nil
}

// convert GenesisDeveloper to AppDeveloper
func (lb *LinoBlockchain) toAppDeveloper(
	ctx sdk.Context, developer GenesisAppDeveloper) sdk.Error {
	if !lb.accountManager.DoesAccountExist(ctx, types.AccountKey(developer.Name)) {
		return ErrGenesisFailed("genesis developer account doesn't exist")
	}

	if err := lb.accountManager.MinusSavingCoin(
		ctx, types.AccountKey(developer.Name), developer.Deposit,
		"", "", types.DeveloperDeposit); err != nil {
		return err
	}

	if err := lb.developerManager.RegisterDeveloper(
		ctx, types.AccountKey(developer.Name), developer.Deposit, developer.Website,
		developer.Description, developer.AppMetaData); err != nil {
		return err
	}
	return nil
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
		lb.globalManager.SetChainStartTime(ctx, ctx.BlockHeader().Time.Unix())
		lb.globalManager.SetLastBlockTime(ctx, ctx.BlockHeader().Time.Unix())
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

	tags := global.BeginBlocker(ctx, req, lb.globalManager)
	actualPenalty := val.BeginBlocker(ctx, req, lb.valManager)

	// add coins back to inflation pool
	if err := lb.globalManager.AddToValidatorInflationPool(ctx, actualPenalty); err != nil {
		panic(err)
	}

	lb.syncInfoWithVoteManager(ctx)
	lb.executeTimeEvents(ctx)
	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
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
			lb.globalManager.RemoveTimeEventList(ctx, i)
		}
	}
	if err := lb.globalManager.SetLastBlockTime(ctx, currentTime); err != nil {
		panic(err)
	}
}

// execute events in list based on their type
func (lb *LinoBlockchain) executeEvents(ctx sdk.Context, eventList []types.Event) sdk.Error {
	for _, event := range eventList {
		switch e := event.(type) {
		case post.RewardEvent:
			if err := e.Execute(
				ctx, lb.postManager, lb.accountManager, lb.globalManager,
				lb.developerManager, lb.voteManager, lb.reputationManager); err != nil {
				panic(err)
			}
		case acc.ReturnCoinEvent:
			if err := e.Execute(ctx, lb.accountManager); err != nil {
				panic(err)
			}
		case proposal.DecideProposalEvent:
			if err := e.Execute(
				ctx, lb.voteManager, lb.valManager, lb.accountManager, lb.proposalManager,
				lb.postManager, lb.globalManager); err != nil {
				panic(err)
			}
		case param.ChangeParamEvent:
			if err := e.Execute(ctx, lb.paramHolder); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

// udpate validator set and renew reputation round
func (lb *LinoBlockchain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	ABCIValList, err := lb.valManager.GetUpdateValidatorList(ctx)
	if err != nil {
		panic(err)
	}
	rep.EndBlocker(ctx, req, lb.reputationManager)

	return abci.ResponseEndBlock{ValidatorUpdates: ABCIValList}
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
	if pastMinutes%types.MinutesPerYear == 0 {
		lb.executeAnnuallyEvent(ctx)
	}
}

// execute hourly event, distribute inflation to validators and
// add hourly inflation to content creator reward pool
func (lb *LinoBlockchain) executeHourlyEvent(ctx sdk.Context) {
	lb.globalManager.DistributeHourlyInflation(ctx)
	lb.distributeInflationToValidator(ctx)
}

// execute daily event, record consumption friction and lino power
func (lb *LinoBlockchain) executeDailyEvent(ctx sdk.Context) {
	lb.globalManager.RecordConsumptionAndLinoStake(ctx)
}

// execute monthly event, distribute inflation to infra and application
func (lb *LinoBlockchain) executeMonthlyEvent(ctx sdk.Context) {
	lb.distributeInflationToInfraProvider(ctx)
	lb.distributeInflationToDeveloper(ctx)
}

func (lb *LinoBlockchain) executeAnnuallyEvent(ctx sdk.Context) {
	if err := lb.globalManager.SetTotalLinoAndRecalculateGrowthRate(ctx); err != nil {
		panic(err)
	}
}

// distribute inflation to validators
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToValidator(ctx sdk.Context) {
	lst, err := lb.valManager.GetValidatorList(ctx)
	if err != nil {
		panic(err)
	}
	coin, err := lb.globalManager.GetValidatorHourlyInflation(ctx)
	if err != nil {
		panic(err)
	}
	// give inflation to each validator evenly
	for i, validator := range lst.OncallValidators {
		var ratPerValidator sdk.Rat
		if ctx.BlockHeader().Height > types.LinoBlockchainFirstUpdateHeight {
			ratPerValidator = coin.ToRat().Quo(sdk.NewRat(int64(len(lst.OncallValidators) - i)))
		} else {
			ratPerValidator = coin.ToRat().Quo(sdk.NewRat(int64(len(lst.OncallValidators) - i))).Round(types.PrecisionFactor)
		}
		coinPerValidator := types.RatToCoin(ratPerValidator)
		lb.accountManager.AddSavingCoin(
			ctx, validator, coinPerValidator, "", "", types.ValidatorInflation)
		coin = coin.Minus(coinPerValidator)
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
			lb.accountManager.AddSavingCoin(
				ctx, provider, inflation.Minus(totalDistributedInflation), "", "", types.InfraInflation)
			break
		}
		percentage, err := lb.infraManager.GetUsageWeight(ctx, provider)
		if err != nil {
			panic(err)
		}
		myShareRat := inflation.ToRat().Mul(percentage)
		myShareCoin := types.RatToCoin(myShareRat)
		totalDistributedInflation = totalDistributedInflation.Plus(myShareCoin)
		lb.accountManager.AddSavingCoin(
			ctx, provider, myShareCoin, "", "", types.InfraInflation)
	}
	if err := lb.infraManager.ClearUsage(ctx); err != nil {
		panic(err)
	}
}

// distribute inflation to developer monthly
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToDeveloper(ctx sdk.Context) {
	inflation, err := lb.globalManager.GetDeveloperMonthlyInflation(ctx)
	if err != nil {
		panic(err)
	}

	lst, err := lb.developerManager.GetDeveloperList(ctx)
	if err != nil {
		panic(err)
	}

	totalDistributedInflation := types.NewCoinFromInt64(0)
	for idx, developer := range lst.AllDevelopers {
		if idx == (len(lst.AllDevelopers) - 1) {
			lb.accountManager.AddSavingCoin(
				ctx, developer, inflation.Minus(totalDistributedInflation), "", "", types.DeveloperInflation)
			break
		}
		percentage, err := lb.developerManager.GetConsumptionWeight(ctx, developer)
		if err != nil {
			panic(err)
		}
		myShareRat := inflation.ToRat().Mul(percentage)
		myShareCoin := types.RatToCoin(myShareRat)
		totalDistributedInflation = totalDistributedInflation.Plus(myShareCoin)
		lb.accountManager.AddSavingCoin(
			ctx, developer, myShareCoin, "", "", types.DeveloperInflation)
	}

	if err := lb.developerManager.ClearConsumption(ctx); err != nil {
		panic(err)
	}
}

func (lb *LinoBlockchain) syncInfoWithVoteManager(ctx sdk.Context) {
	// tell voting committee the newest validators
	validatorList, err := lb.valManager.GetValidatorList(ctx)
	if err != nil {
		panic(err)
	}

	referenceList, err := lb.voteManager.GetValidatorReferenceList(ctx)
	if err != nil {
		panic(err)
	}
	referenceList.AllValidators = validatorList.AllValidators
	if err := lb.voteManager.SetValidatorReferenceList(ctx, referenceList); err != nil {
		panic(err)
	}
}

// Custom logic for state export
func (lb *LinoBlockchain) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})

	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(accInfo accModel.AccountInfo, accBank accModel.AccountBank) (stop bool) {
		saving := accBank.Saving
		deposit, err := lb.valManager.GetValidatorDeposit(ctx, accInfo.Username)
		if err != nil {
			saving = saving.Plus(deposit)
		}
		account := GenesisAccount{
			Name:           string(accInfo.Username),
			ResetKey:       accInfo.ResetKey,
			TransactionKey: accInfo.TransactionKey,
			AppKey:         accInfo.AppKey,
			IsValidator:    false,
			Coin:           saving,
		}
		accounts = append(accounts, account)
		return false
	}
	lb.accountManager.IterateAccounts(ctx, appendAccount)

	genesisState := GenesisState{
		Accounts:   accounts,
		Developers: []GenesisAppDeveloper{},
		Infra:      []GenesisInfraProvider{},
	}
	appState, err = wire.MarshalJSONIndent(lb.cdc, genesisState)
	if err != nil {
		return nil, nil, err
	}
	return appState, validators, nil
}
