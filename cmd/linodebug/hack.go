package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/auth"
	"github.com/lino-network/lino/x/developer"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/infra"
	"github.com/lino-network/lino/x/post"
	"github.com/lino-network/lino/x/proposal"
	"github.com/lino-network/lino/x/vote"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/wire"
	acc "github.com/lino-network/lino/x/account"
	rep "github.com/lino-network/lino/x/reputation"
	val "github.com/lino-network/lino/x/validator"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	cmn "github.com/tendermint/tmlibs/common"
)

var hackCmd = &cobra.Command{
	Use:   "hack",
	Short: "Boilerplate to Hack on an existing state by scripting some Go...",
	RunE:  runHackCmd,
}

func runHackCmd(cmd *cobra.Command, args []string) error {

	if len(args) != 2 {
		return fmt.Errorf("Expected 1 arg")
	}

	// ".lino"
	dataDir := args[0]
	dataDir = path.Join(dataDir, "data")

	fmt.Println("data dir:", dataDir, "height:", args[1])
	// load the app
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	db, err := dbm.NewGoLevelDB("lino", dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app := NewLinoBlockchain(logger, db, os.Stdout)

	// print some info
	id := app.LastCommitID()
	lastBlockHeight := app.LastBlockHeight()
	fmt.Println("ID", id)
	fmt.Println("LastBlockHeight", lastBlockHeight)

	checkHeight, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		panic(err)
	}
	// load the given version of the state
	err = app.LoadVersion(checkHeight, app.CapKeyMainStore)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := app.NewContext(true, abci.Header{})

	// check for the powerkey and the validator from the store
	fmt.Println("last commit ID:", app.LastCommitID(), ", last block height:", app.LastBlockHeight())
	keyList := map[string]*sdk.KVStoreKey{
		"reputation": app.CapKeyReputationStore}
	// resultArray := [][32]byte{}
	for name, key := range keyList {
		store := ctx.KVStore(key)
		result := iterateStore(store)
		fmt.Println(name, result)
	}
	return nil
}

type bigInt = *big.Int
type Rep = bigInt
type RoundId = int64
type Uid = string
type Time = int64
type Pid = string
type Dp = bigInt // donation power

// used in topN.
type PostDpPair struct {
	Pid   Pid
	SumDp Dp
}

type userMeta struct {
	CustomerScore     Rep
	FreeScore         Rep
	LastSettled       RoundId
	LastDonationRound RoundId
}
type postMeta struct {
	SumRep Rep
}

type roundMeta struct {
	Result  []Pid
	SumDp   Dp
	StartAt Time
	TopN    []PostDpPair
}

func iterateStore(store sdk.KVStore) [32]byte {
	storeResult := ""
	iter := sdk.KVStorePrefixIterator(store, []byte{0x00})
	for {
		if !iter.Valid() {
			break
		}
		rst := &userMeta{}
		val := iter.Value()
		dec := gob.NewDecoder(bytes.NewBuffer(val))
		dec.Decode(rst)
		fmt.Println(rst, iter.Key(), hex.EncodeToString(iter.Value()))
		storeResult += string(val)
		// fmt.Println(string(val))
		iter.Next()
	}
	iter = sdk.KVStorePrefixIterator(store, []byte{0x01})
	for {
		if !iter.Valid() {
			break
		}
		rst := &postMeta{}
		val := iter.Value()
		dec := gob.NewDecoder(bytes.NewBuffer(val))
		storeResult += string(val)
		dec.Decode(rst)
		fmt.Println(rst, iter.Key(), hex.EncodeToString(iter.Value()))
		// fmt.Println(string(val))
		iter.Next()
	}
	iter = sdk.KVStorePrefixIterator(store, []byte{0x03})
	for {
		if !iter.Valid() {
			break
		}
		rst := &roundMeta{}
		val := iter.Value()
		dec := gob.NewDecoder(bytes.NewBuffer(val))
		storeResult += string(val)
		dec.Decode(rst)
		fmt.Println(rst, iter.Key(), hex.EncodeToString(iter.Value()))
		// fmt.Println(string(val))
		iter.Next()
	}
	return sha256.Sum256([]byte(storeResult))
}

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
}

// NewLinoBlockchain - create a Lino Blockchain instance
func NewLinoBlockchain(
	logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *LinoBlockchain {
	// create your application object
	cdc := app.MakeCodec()
	bApp := bam.NewBaseApp("LinoBlockchain", logger, db, app.DefaultTxDecoder(cdc), baseAppOptions...)
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
	lb.paramHolder = param.NewParamHolder(lb.CapKeyParamStore)
	lb.accountManager = acc.NewAccountManager(lb.CapKeyAccountStore, lb.paramHolder)
	lb.postManager = post.NewPostManager(lb.CapKeyPostStore, lb.paramHolder)
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
	genesisState := new(app.GenesisState)
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
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga app.GenesisAccount) sdk.Error {
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
	ctx sdk.Context, developer app.GenesisAppDeveloper) sdk.Error {
	if !lb.accountManager.DoesAccountExist(ctx, types.AccountKey(developer.Name)) {
		return app.ErrGenesisFailed("genesis developer account doesn't exist")
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
	ctx sdk.Context, infra app.GenesisInfraProvider) sdk.Error {
	if !lb.accountManager.DoesAccountExist(ctx, types.AccountKey(infra.Name)) {
		return app.ErrGenesisFailed("genesis infra account doesn't exist")
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
		ratPerValidator := coin.ToRat().Quo(sdk.NewRat(int64(len(lst.OncallValidators) - i))).Round(types.PrecisionFactor)
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
