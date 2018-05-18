package app

import (
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/auth"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/tx/proposal"
	"github.com/lino-network/lino/tx/register"
	"github.com/lino-network/lino/types"

	acc "github.com/lino-network/lino/tx/account"
	developer "github.com/lino-network/lino/tx/developer"
	infra "github.com/lino-network/lino/tx/infra"
	val "github.com/lino-network/lino/tx/validator"
	vote "github.com/lino-network/lino/tx/vote"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

const (
	appName = "LinoBlockchain"
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

func NewLinoBlockchain(logger log.Logger, db dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:              bam.NewBaseApp(appName, logger, db),
		cdc:                  MakeCodec(),
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
	lb.voteManager = vote.NewVoteManager(lb.CapKeyVoteStore, lb.paramHolder)
	lb.infraManager = infra.NewInfraManager(lb.CapKeyInfraStore, lb.paramHolder)
	lb.developerManager = developer.NewDeveloperManager(lb.CapKeyDeveloperStore, lb.paramHolder)
	lb.proposalManager = proposal.NewProposalManager(lb.CapKeyProposalStore, lb.paramHolder)

	RegisterEvent(lb.globalManager.WireCodec())

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(lb.accountManager)).
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(lb.postManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.VoteRouterName, vote.NewHandler(lb.voteManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.DeveloperRouterName, developer.NewHandler(
			lb.developerManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.ProposalRouterName, proposal.NewHandler(
			lb.accountManager, lb.proposalManager, lb.postManager, lb.globalManager)).
		AddRoute(types.InfraRouterName, infra.NewHandler(lb.infraManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(
			lb.accountManager, lb.valManager, lb.voteManager, lb.globalManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetBeginBlocker(lb.beginBlocker)
	lb.SetEndBlocker(lb.endBlocker)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoresIAVL(
		lb.CapKeyMainStore, lb.CapKeyAccountStore, lb.CapKeyPostStore, lb.CapKeyValStore,
		lb.CapKeyVoteStore, lb.CapKeyInfraStore, lb.CapKeyDeveloperStore, lb.CapKeyGlobalStore,
		lb.CapKeyParamStore, lb.CapKeyProposalStore)
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager, lb.globalManager))
	if err := lb.LoadLatestVersion(lb.CapKeyMainStore); err != nil {
		cmn.Exit(err.Error())
	}
	return lb
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(register.RegisterMsg{}, "register", nil)
	cdc.RegisterConcrete(acc.FollowMsg{}, "follow", nil)
	cdc.RegisterConcrete(acc.UnfollowMsg{}, "unfollow", nil)
	cdc.RegisterConcrete(acc.TransferMsg{}, "transfer", nil)
	cdc.RegisterConcrete(acc.ClaimMsg{}, "claim", nil)
	cdc.RegisterConcrete(post.CreatePostMsg{}, "post", nil)
	cdc.RegisterConcrete(post.UpdatePostMsg{}, "update/post", nil)
	cdc.RegisterConcrete(post.DeletePostMsg{}, "delete/post", nil)
	cdc.RegisterConcrete(post.LikeMsg{}, "like", nil)
	cdc.RegisterConcrete(post.DonateMsg{}, "donate", nil)
	cdc.RegisterConcrete(post.ReportOrUpvoteMsg{}, "reportOrUpvote", nil)
	cdc.RegisterConcrete(val.ValidatorDepositMsg{}, "val/deposit", nil)
	cdc.RegisterConcrete(val.ValidatorWithdrawMsg{}, "val/withdraw", nil)
	cdc.RegisterConcrete(val.ValidatorRevokeMsg{}, "val/revoke", nil)
	cdc.RegisterConcrete(vote.VoterDepositMsg{}, "vote/deposit", nil)
	cdc.RegisterConcrete(vote.VoterRevokeMsg{}, "vote/revoke", nil)
	cdc.RegisterConcrete(vote.VoterWithdrawMsg{}, "vote/withdraw", nil)
	cdc.RegisterConcrete(vote.DelegateMsg{}, "delegate", nil)
	cdc.RegisterConcrete(vote.DelegatorWithdrawMsg{}, "delegate/withdraw", nil)
	cdc.RegisterConcrete(vote.RevokeDelegationMsg{}, "delegate/revoke", nil)
	cdc.RegisterConcrete(vote.VoteMsg{}, "vote", nil)
	cdc.RegisterConcrete(developer.DeveloperRegisterMsg{}, "developer/register", nil)
	cdc.RegisterConcrete(developer.DeveloperRevokeMsg{}, "developer/revoke", nil)
	cdc.RegisterConcrete(infra.ProviderReportMsg{}, "provider/report", nil)
	cdc.RegisterConcrete(developer.GrantDeveloperMsg{}, "grant/developer", nil)

	wire.RegisterCrypto(cdc)
	return cdc
}

func RegisterEvent(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(post.RewardEvent{}, "event/reward", nil)
	cdc.RegisterConcrete(acc.ReturnCoinEvent{}, "event/return", nil)
	cdc.RegisterConcrete(param.ChangeParamEvent{}, "event/cpe", nil)
	cdc.RegisterConcrete(proposal.DecideProposalEvent{}, "event/dpe", nil)
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
	if err := lb.cdc.UnmarshalJSON(stateJSON, genesisState); err != nil {
		panic(err)
	}

	if err := lb.paramHolder.InitParam(ctx); err != nil {
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
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga genesis.GenesisAccount) sdk.Error {
	// send coins using address (even no account bank associated with this addr)
	coin, err := types.LinoToCoin(ga.Lino)
	if err != nil {
		panic(err)
	}
	if err := lb.accountManager.AddSavingCoinToAddress(ctx, ga.MasterKey.Address(), coin); err != nil {
		panic(sdk.ErrGenesisParse("set genesis bank failed"))
	}
	if lb.accountManager.IsAccountExist(ctx, types.AccountKey(ga.Name)) {
		panic(sdk.ErrGenesisParse("genesis account already exist"))
	}
	if err := lb.accountManager.CreateAccount(
		ctx, types.AccountKey(ga.Name),
		ga.MasterKey, ga.TransactionKey, ga.PostKey); err != nil {
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
			valParam.ValidatorMinCommitingDeposit.Plus(valParam.ValidatorMinVotingDeposit)); err != nil {
			panic(err)
		}

		if err := lb.voteManager.AddVoter(
			ctx, types.AccountKey(ga.Name), valParam.ValidatorMinVotingDeposit); err != nil {
			panic(err)
		}
		if err := lb.valManager.RegisterValidator(
			ctx, types.AccountKey(ga.Name), ga.ValPubKey.Bytes(), valParam.ValidatorMinCommitingDeposit, ""); err != nil {
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
	ctx sdk.Context, developer genesis.GenesisAppDeveloper) sdk.Error {
	if !lb.accountManager.IsAccountExist(ctx, types.AccountKey(developer.Name)) {
		return sdk.ErrGenesisParse("genesis developer account doesn't exist")
	}
	coin, err := types.LinoToCoin(types.LNO(developer.Deposit))
	if err != nil {
		return err
	}

	if err := lb.accountManager.MinusSavingCoin(ctx, types.AccountKey(developer.Name), coin); err != nil {
		return err
	}

	if err := lb.developerManager.RegisterDeveloper(
		ctx, types.AccountKey(developer.Name), coin); err != nil {
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

	validatorList, err := lb.valManager.GetValidatorList(ctx)
	if err != nil {
		panic(err)
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

	actualPenalty, err := lb.valManager.FireIncompetentValidator(ctx, req.GetByzantineValidators())
	if err != nil {
		panic(err)
	}

	// add coins back to inflation pool
	if err := lb.globalManager.AddToValidatorInflationPool(ctx, actualPenalty); err != nil {
		panic(err)
	}

	lb.syncInfoWithVoteManager(ctx)
	lb.executeTimeEvents(ctx)
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

func (lb *LinoBlockchain) executeAnnuallyEvent(ctx sdk.Context) {
	if err := lb.globalManager.RecalculateAnnuallyInflation(ctx); err != nil {
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
	pastHoursThisYear := (lb.pastMinutes / 60) % types.HoursPerYear
	coin, err := lb.globalManager.GetValidatorHourlyInflation(ctx, pastHoursThisYear)
	if err != nil {
		panic(err)
	}
	// give inflation to each validator evenly
	ratPerValidator := coin.ToRat().Quo(sdk.NewRat(int64(len(lst.OncallValidators))))
	for _, validator := range lst.OncallValidators {
		lb.accountManager.AddSavingCoin(ctx, validator, types.RatToCoin(ratPerValidator))
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

	lst, err := lb.infraManager.GetInfraProviderList(ctx)
	if err != nil {
		panic(err)
	}

	for _, provider := range lst.AllInfraProviders {
		percentage, err := lb.infraManager.GetUsageWeight(ctx, provider)
		if err != nil {
			panic(err)
		}
		myShare := inflation.ToRat().Mul(percentage)
		lb.accountManager.AddSavingCoin(ctx, provider, types.RatToCoin(myShare))
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

	lst, err := lb.developerManager.GetDeveloperList(ctx)
	if err != nil {
		panic(err)
	}

	for _, developer := range lst.AllDevelopers {
		percentage, err := lb.developerManager.GetConsumptionWeight(ctx, developer)
		if err != nil {
			panic(err)
		}
		myShare := inflation.ToRat().Mul(percentage)
		lb.accountManager.AddSavingCoin(ctx, developer, types.RatToCoin(myShare))
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

	proposalList, err := lb.proposalManager.GetProposalList(ctx)
	if err != nil {
		panic(err)
	}

	referenceList, err := lb.voteManager.GetValidatorReferenceList(ctx)
	if err != nil {
		panic(err)
	}
	referenceList.AllValidators = validatorList.AllValidators
	referenceList.OngoingProposal = proposalList.OngoingProposal
	if err := lb.voteManager.SetValidatorReferenceList(ctx, referenceList); err != nil {
		panic(err)
	}
}
