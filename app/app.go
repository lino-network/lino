package app

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/auth"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post"
	"github.com/lino-network/lino/x/proposal"

	acc "github.com/lino-network/lino/x/account"
	developer "github.com/lino-network/lino/x/developer"
	infra "github.com/lino-network/lino/x/infra"
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
	CapKeyParamStore     *sdk.KVStoreKey
	CapKeyProposalStore  *sdk.KVStoreKey

	// Manager for different KVStore
	accountManager   acc.AccountManager
	postManager      post.PostManager
	valManager       val.ValidatorManager
	globalManager    global.GlobalManager
	voteManager      vote.VoteManager
	infraManager     infra.InfraManager
	developerManager developer.DeveloperManager
	proposalManager  proposal.ProposalManager

	// global param
	paramHolder param.ParamHolder
	// time related
	chainStartTime int64
	lastBlockTime  int64
	pastMinutes    int64
}

func NewLinoBlockchain(
	logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *LinoBlockchain {
	// create your application object
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, cdc, logger, db, baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	var lb = &LinoBlockchain{
		BaseApp:              bApp,
		cdc:                  cdc,
		CapKeyMainStore:      sdk.NewKVStoreKey(types.MainKVStoreKey),
		CapKeyAccountStore:   sdk.NewKVStoreKey(types.AccountKVStoreKey),
		CapKeyPostStore:      sdk.NewKVStoreKey(types.PostKVStoreKey),
		CapKeyValStore:       sdk.NewKVStoreKey(types.ValidatorKVStoreKey),
		CapKeyVoteStore:      sdk.NewKVStoreKey(types.VoteKVStoreKey),
		CapKeyInfraStore:     sdk.NewKVStoreKey(types.InfraKVStoreKey),
		CapKeyDeveloperStore: sdk.NewKVStoreKey(types.DeveloperKVStoreKey),
		CapKeyGlobalStore:    sdk.NewKVStoreKey(types.GlobalKVStoreKey),
		CapKeyParamStore:     sdk.NewKVStoreKey(types.ParamKVStoreKey),
		CapKeyProposalStore:  sdk.NewKVStoreKey(types.ProposalKVStoreKey),
	}
	lb.paramHolder = param.NewParamHolder(lb.CapKeyParamStore)
	lb.accountManager = acc.NewAccountManager(lb.CapKeyAccountStore, lb.paramHolder)
	lb.postManager = post.NewPostManager(lb.CapKeyPostStore, lb.paramHolder)
	lb.valManager = val.NewValidatorManager(lb.CapKeyValStore, lb.paramHolder)
	lb.globalManager = global.NewGlobalManager(lb.CapKeyGlobalStore, lb.paramHolder)
	RegisterEvent(lb.globalManager.WireCodec())

	lb.voteManager = vote.NewVoteManager(lb.CapKeyVoteStore, lb.paramHolder)
	lb.infraManager = infra.NewInfraManager(lb.CapKeyInfraStore, lb.paramHolder)
	lb.developerManager = developer.NewDeveloperManager(lb.CapKeyDeveloperStore, lb.paramHolder)
	lb.proposalManager = proposal.NewProposalManager(lb.CapKeyProposalStore, lb.paramHolder)

	lb.Router().
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(
			lb.postManager, lb.accountManager, lb.globalManager, lb.developerManager)).
		AddRoute(types.VoteRouterName, vote.NewHandler(lb.voteManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.DeveloperRouterName, developer.NewHandler(
			lb.developerManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.ProposalRouterName, proposal.NewHandler(
			lb.accountManager, lb.proposalManager, lb.postManager, lb.globalManager, lb.voteManager)).
		AddRoute(types.InfraRouterName, infra.NewHandler(lb.infraManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(
			lb.accountManager, lb.valManager, lb.voteManager, lb.globalManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager, lb.globalManager))
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoresIAVL(
		lb.CapKeyMainStore, lb.CapKeyAccountStore, lb.CapKeyPostStore, lb.CapKeyValStore,
		lb.CapKeyVoteStore, lb.CapKeyInfraStore, lb.CapKeyDeveloperStore, lb.CapKeyGlobalStore,
		lb.CapKeyParamStore, lb.CapKeyProposalStore)
	if err := lb.LoadLatestVersion(lb.CapKeyMainStore); err != nil {
		cmn.Exit(err.Error())
	}
	return lb
}

// custom logic for transaction decoding
func (lb *LinoBlockchain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = cauth.StdTx{}

	// StdTx.Msg is an interface.
	err := lb.cdc.UnmarshalJSON(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("")
	}
	return tx, nil
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	acc.RegisterWire(cdc)
	post.RegisterWire(cdc)
	developer.RegisterWire(cdc)
	infra.RegisterWire(cdc)
	vote.RegisterWire(cdc)
	val.RegisterWire(cdc)
	proposal.RegisterWire(cdc)

	RegisterEvent(cdc)
	return cdc
}

func RegisterEvent(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(post.RewardEvent{}, "lino/eventReward", nil)
	cdc.RegisterConcrete(acc.ReturnCoinEvent{}, "lino/eventReturn", nil)
	cdc.RegisterConcrete(param.ChangeParamEvent{}, "lino/eventCpe", nil)
	cdc.RegisterConcrete(proposal.DecideProposalEvent{}, "lino/eventDpe", nil)
}

// custom logic for basecoin initialization
func (lb *LinoBlockchain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	genesisState := new(GenesisState)
	if err := lb.cdc.UnmarshalJSON(stateJSON, genesisState); err != nil {
		panic(err)
	}

	if err := lb.paramHolder.InitParam(ctx); err != nil {
		panic(err)
	}

	totalCoin := types.NewCoinFromInt64(0)

	for _, gacc := range genesisState.Accounts {
		coin, err := types.LinoToCoin(gacc.Lino)
		if err != nil {
			panic(err)
		}
		totalCoin = totalCoin.Plus(coin)
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
	if err := lb.proposalManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.valManager.InitGenesis(ctx); err != nil {
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
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga GenesisAccount) sdk.Error {
	// send coins using address (even no account bank associated with this addr)
	coin, err := types.LinoToCoin(ga.Lino)
	if err != nil {
		panic(err)
	}
	if lb.accountManager.DoesAccountExist(ctx, types.AccountKey(ga.Name)) {
		panic(errors.New("genesis account already exist"))
	}
	if err := lb.accountManager.CreateAccount(
		ctx, types.AccountKey(ga.Name), types.AccountKey(ga.Name),
		ga.ResetKey, ga.TransactionKey, ga.MicropaymentKey, ga.PostKey, coin); err != nil {
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
			valParam.ValidatorMinCommitingDeposit.Plus(valParam.ValidatorMinVotingDeposit),
			"", "", types.ValidatorDeposit); err != nil {
			panic(err)
		}

		if err := lb.voteManager.AddVoter(
			ctx, types.AccountKey(ga.Name), valParam.ValidatorMinVotingDeposit); err != nil {
			panic(err)
		}
		if err := lb.valManager.RegisterValidator(
			ctx, types.AccountKey(ga.Name), ga.ValPubKey,
			valParam.ValidatorMinCommitingDeposit, ""); err != nil {
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
	coin, err := types.LinoToCoin(types.LNO(developer.Deposit))
	if err != nil {
		return err
	}

	if err := lb.accountManager.MinusSavingCoin(
		ctx, types.AccountKey(developer.Name), coin,
		"", "", types.DeveloperDeposit); err != nil {
		return err
	}

	if err := lb.developerManager.RegisterDeveloper(
		ctx, types.AccountKey(developer.Name), coin, developer.Website,
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
	if lb.chainStartTime == 0 {
		lb.chainStartTime = ctx.BlockHeader().Time
		lb.lastBlockTime = ctx.BlockHeader().Time
	}

	for (ctx.BlockHeader().Time-lb.chainStartTime)/60 > lb.pastMinutes {
		lb.increaseMinute(ctx)
	}

	tags := global.BeginBlocker(ctx, req, lb.globalManager, lb.lastBlockTime)
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
	if lb.pastMinutes%types.MinutesPerYear == 0 {
		lb.executeAnnuallyEvent(ctx)
	}
}

// execute hourly event, distribute inflation to validators and
// add hourly inflation to content creator reward pool
func (lb *LinoBlockchain) executeHourlyEvent(ctx sdk.Context) {
	lb.distributeInflationToConsumptionRewardPool(ctx)
	lb.distributeInflationToValidator(ctx)
}

// execute monthly event, distribute inflation to infra and application
func (lb *LinoBlockchain) executeMonthlyEvent(ctx sdk.Context) {
	lb.distributeInflationToInfraProvider(ctx)
	lb.distributeInflationToDeveloper(ctx)
}

func (lb *LinoBlockchain) executeAnnuallyEvent(ctx sdk.Context) {
	if err := lb.globalManager.RecalculateAnnuallyInflation(ctx); err != nil {
		panic(err)
	}
}

// distribute inflation to validators
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToConsumptionRewardPool(ctx sdk.Context) {
	pastHoursMinusOneThisYear := lb.getPastHoursMinusOneThisYear()
	if err := lb.globalManager.AddHourlyInflationToRewardPool(
		ctx, pastHoursMinusOneThisYear); err != nil {
		panic(err)
	}
}

// distribute inflation to validators
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToValidator(ctx sdk.Context) {
	pastHoursMinusOneThisYear := lb.getPastHoursMinusOneThisYear()
	lst, err := lb.valManager.GetValidatorList(ctx)
	if err != nil {
		panic(err)
	}
	coin, err := lb.globalManager.GetValidatorHourlyInflation(ctx, pastHoursMinusOneThisYear)
	if err != nil {
		panic(err)
	}
	// give inflation to each validator evenly
	for i, validator := range lst.OncallValidators {
		ratPerValidator := coin.ToRat().Quo(sdk.NewRat(int64(len(lst.OncallValidators) - i)))
		coinPerValidator := types.RatToCoin(ratPerValidator)
		lb.accountManager.AddSavingCoin(
			ctx, validator, coinPerValidator, "", "", types.ValidatorInflation)
		coin = coin.Minus(coinPerValidator)
	}
}

// distribute inflation to infra provider monthly
// TODO: encaptulate module event inside module
func (lb *LinoBlockchain) distributeInflationToInfraProvider(ctx sdk.Context) {
	pastMonthMinusOneThisYear := lb.getPastMonthMinusOneThisYear()
	inflation, err := lb.globalManager.GetInfraMonthlyInflation(ctx, pastMonthMinusOneThisYear)
	if err != nil {
		panic(err)
	}

	lst, err := lb.infraManager.GetInfraProviderList(ctx)
	if err != nil {
		panic(err)
	}
	for _, provider := range lst.AllInfraProviders {
		percentage, err := lb.infraManager.GetUsageWeight(ctx, provider)
		if err != nil {
			panic(err)
		}
		myShareRat := inflation.ToRat().Mul(percentage)
		myShareCoin := types.RatToCoin(myShareRat)
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
	pastMonthMinusOneThisYear := lb.getPastMonthMinusOneThisYear()
	inflation, err := lb.globalManager.GetDeveloperMonthlyInflation(ctx, pastMonthMinusOneThisYear)
	if err != nil {
		panic(err)
	}

	lst, err := lb.developerManager.GetDeveloperList(ctx)
	if err != nil {
		panic(err)
	}

	for _, developer := range lst.AllDevelopers {
		percentage, err := lb.developerManager.GetConsumptionWeight(ctx, developer)
		if err != nil {
			panic(err)
		}
		myShareRat := inflation.ToRat().Mul(percentage)
		myShareCoin := types.RatToCoin(myShareRat)
		lb.accountManager.AddSavingCoin(
			ctx, developer, myShareCoin, "", "", types.DeveloperInflation)
	}

	if err := lb.developerManager.ClearConsumption(ctx); err != nil {
		panic(err)
	}
}

func (lb *LinoBlockchain) syncInfoWithVoteManager(ctx sdk.Context) {
	// tell voting committe the newest validators
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

func (lb *LinoBlockchain) getPastHoursMinusOneThisYear() int64 {
	return (lb.pastMinutes/60 - 1) % types.HoursPerYear
}

func (lb *LinoBlockchain) getPastMonthMinusOneThisYear() int64 {
	return (lb.pastMinutes/types.MinutesPerMonth - 1) % 12
}

// Custom logic for state export
func (lb *LinoBlockchain) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	//ctx := lb.NewContext(true, abci.Header{})

	// // iterate to get the accounts
	// accounts := []*GenesisAccount{}
	// appendAccount := func(acc auth.Account) (stop bool) {
	// 	account := &types.GenesisAccount{
	// 		Address: acc.GetAddress(),
	// 		Coins:   acc.GetCoins(),
	// 	}
	// 	accounts = append(accounts, account)
	// 	return false
	// }
	// app.accountMapper.IterateAccounts(ctx, appendAccount)

	// genState := types.GenesisState{
	// 	Accounts:    accounts,
	// 	POWGenesis:  pow.WriteGenesis(ctx, app.powKeeper),
	// 	CoolGenesis: cool.WriteGenesis(ctx, app.coolKeeper),
	// }
	// appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	// if err != nil {
	// 	return nil, nil, err
	// }
	return appState, validators, nil
}
