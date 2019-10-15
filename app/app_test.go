//nolint:deadcode,unused
package app

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	globalModel "github.com/lino-network/lino/x/global/model"
	posttypes "github.com/lino-network/lino/x/post/types"
)

var (
	user1 = "validator0"
	priv1 = secp256k1.GenPrivKey()
	addr1 = priv1.PubKey().Address()
	priv2 = secp256k1.GenPrivKey()
	addr2 = priv2.PubKey().Address()

	genesisTotalCoin    = types.NewCoinFromInt64(2100000000 * types.Decimals)
	coinPerValidator    = types.NewCoinFromInt64(100000000 * types.Decimals)
	growthRate          = types.NewDecFromRat(98, 1000)
	validatorAllocation = types.NewDecFromRat(5, 100)
)

func loggerAndDB() (logger log.Logger, db dbm.DB) {
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "lino/app")
	db = dbm.NewMemDB()
	return
}

func newLinoBlockchain(t *testing.T, numOfValidators int) *LinoBlockchain {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db, nil)

	genesisState := GenesisState{
		InitCoinPrice: types.NewMiniDollar(1200),
		Accounts:      []GenesisAccount{},
	}

	// Generate 21 validators
	genesisAcc := GenesisAccount{
		Name:           user1,
		Coin:           coinPerValidator,
		ResetKey:       priv1.PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
		AppKey:         secp256k1.GenPrivKey().PubKey(),
		IsValidator:    true,
		ValPubKey:      priv2.PubKey(),
	}
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	for i := 1; i < numOfValidators; i++ {
		genesisAcc := GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Coin:           coinPerValidator,
			ResetKey:       secp256k1.GenPrivKey().PubKey(),
			TransactionKey: secp256k1.GenPrivKey().PubKey(),
			AppKey:         secp256k1.GenPrivKey().PubKey(),
			IsValidator:    true,
			ValPubKey:      secp256k1.GenPrivKey().PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}
	genesisState.InitGlobalMeta = globalModel.InitParamList{
		MaxTPS:                       sdk.NewDec(1000),
		ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
		ConsumptionFrictionRate:      types.NewDecFromRat(5, 100),
	}

	result, err := wire.MarshalJSONIndent(lb.cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{AppStateBytes: json.RawMessage(result)})
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{Height: 1, ChainID: "Lino", Time: time.Unix(0, 0)}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	return lb
}

func TestGenesisAcc(t *testing.T) {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db, nil)

	accs := []struct {
		genesisAccountName string
		coin               types.Coin
		resetKey           crypto.PubKey
		transactionKey     crypto.PubKey
		appKey             crypto.PubKey
		isValidator        bool
		valPubKey          crypto.PubKey
	}{
		{"lino", types.NewCoinFromInt64(9000000000 * types.Decimals), secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			true, secp256k1.GenPrivKey().PubKey()},
		{"genesis", types.NewCoinFromInt64(500000000 * types.Decimals), secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			true, secp256k1.GenPrivKey().PubKey()},
		{"nonvalidator", types.NewCoinFromInt64(500000000 * types.Decimals), secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			false, secp256k1.GenPrivKey().PubKey()},
		{"developer", types.NewCoinFromInt64(500000000 * types.Decimals), secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			false, secp256k1.GenPrivKey().PubKey()},
	}
	genesisState := GenesisState{
		InitCoinPrice: types.NewMiniDollar(1200),
		Accounts:      []GenesisAccount{},
	}
	for _, acc := range accs {
		genesisAcc := GenesisAccount{
			Name:           acc.genesisAccountName,
			Coin:           acc.coin,
			ResetKey:       acc.resetKey,
			TransactionKey: acc.transactionKey,
			AppKey:         acc.appKey,
			IsValidator:    acc.isValidator,
			ValPubKey:      acc.valPubKey,
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}
	// TODO(yumin): add developer genesis test back.
	// genesisAppDeveloper := GenesisAppDeveloper{
	// 	Name:        "developer",
	// 	Deposit:     types.NewCoinFromInt64(1000000 * types.Decimals),
	// 	Website:     "https://lino.network/",
	// 	Description: "",
	// 	AppMetaData: "",
	// }
	// genesisState.Developers = append(genesisState.Developers, genesisAppDeveloper)
	result, err := wire.MarshalJSONIndent(lb.cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{AppStateBytes: json.RawMessage(result)})
	lb.Commit()

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	for _, acc := range accs {
		expectBalance := acc.coin
		if acc.isValidator {
			param := lb.paramHolder.GetValidatorParam(ctx)
			expectBalance = expectBalance.Minus(param.ValidatorMinDeposit)
		}
		if acc.genesisAccountName == "developer" {
			param, _ := lb.paramHolder.GetDeveloperParam(ctx)
			expectBalance = expectBalance.Minus(param.DeveloperMinDeposit)
			// TODO(yumin): add developer genesis test back.
			continue
		}
		saving, err :=
			lb.accountManager.GetSavingFromUsername(ctx, types.AccountKey(acc.genesisAccountName))
		assert.Nil(t, err)
		assert.Equal(
			t, expectBalance, saving,
			"account %s saving is %s, expect is %s, struct: %+v", acc.genesisAccountName, saving.String(), expectBalance.String(), acc)
	}
}

func TestGenesisFromConfig(t *testing.T) {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db, nil)
	genesisState := GenesisState{
		InitCoinPrice: types.NewMiniDollar(1200),
		Accounts:      []GenesisAccount{},
	}
	genesisState.GenesisParam = GenesisParam{
		true,
		param.GlobalAllocationParam{
			GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
			ContentCreatorAllocation: types.NewDecFromRat(85, 100),
			DeveloperAllocation:      types.NewDecFromRat(10, 100),
			ValidatorAllocation:      types.NewDecFromRat(5, 100),
		},
		param.VoteParam{
			MinStakeIn:                     types.NewCoinFromInt64(1000 * types.Decimals),
			VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
			VoterCoinReturnTimes:           int64(7),
			DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
			DelegatorCoinReturnTimes:       int64(7),
		},
		param.ProposalParam{
			ContentCensorshipDecideSec:  int64(24 * 7 * 3600),
			ContentCensorshipPassRatio:  types.NewDecFromRat(50, 100),
			ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
			ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

			ChangeParamDecideSec:  int64(24 * 7 * 3600),
			ChangeParamPassRatio:  types.NewDecFromRat(70, 100),
			ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
			ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

			ProtocolUpgradeDecideSec:  int64(24 * 7 * 3600),
			ProtocolUpgradePassRatio:  types.NewDecFromRat(80, 100),
			ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
			ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
		},
		param.DeveloperParam{
			DeveloperMinDeposit:            types.NewCoinFromInt64(1000000 * types.Decimals),
			DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
			DeveloperCoinReturnTimes:       int64(7),
		},
		param.ValidatorParam{
			ValidatorMinDeposit:            types.NewCoinFromInt64(200000 * types.Decimals),
			ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
			ValidatorCoinReturnTimes:       int64(7),
			PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
			PenaltyByzantine:               types.NewCoinFromInt64(1000 * types.Decimals),
			AbsentCommitLimitation:         int64(600), // 30min
			OncallSize:                     int64(22),
			StandbySize:                    int64(7),
			ValidatorRevokePendingSec:      int64(7 * 24 * 3600),
			OncallInflationWeight:          int64(2),
			StandbyInflationWeight:         int64(1),
			MaxVotedValidators:             int64(3),
			SlashLimitation:                int64(5),
		},
		param.CoinDayParam{
			SecondsToRecoverCoinDay: int64(7 * 24 * 3600),
		},
		param.BandwidthParam{
			SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
			CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
			VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
			GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
			GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
			AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
			AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
			ExpectedMaxMPS:              types.NewDecFromRat(1000, 1),
			MsgFeeFactorA:               types.NewDecFromRat(6, 1),
			MsgFeeFactorB:               types.NewDecFromRat(10, 1),
			MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
			AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
			AppVacancyFactor:            types.NewDecFromRat(69, 100),
			AppPunishmentFactor:         types.NewDecFromRat(14, 5),
		},
		param.AccountParam{
			MinimumBalance:               types.NewCoinFromInt64(1 * types.Decimals),
			RegisterFee:                  types.NewCoinFromInt64(0),
			FirstDepositFullCoinDayLimit: types.NewCoinFromInt64(0),
			MaxNumFrozenMoney:            10,
		},
		param.PostParam{
			ReportOrUpvoteIntervalSec: 24 * 3600,
			PostIntervalSec:           600,
			MaxReportReputation:       types.NewCoinFromInt64(100 * types.Decimals),
		},
		param.ReputationParam{
			BestContentIndexN: 200,
			UserMaxN:          50,
		},
		param.PriceParam{
			TestnetMode:     true,
			UpdateEverySec:  int64(time.Hour.Seconds()),
			FeedEverySec:    int64((10 * time.Minute).Seconds()),
			HistoryMaxLen:   71,
			PenaltyMissFeed: types.NewCoinFromInt64(10000 * types.Decimals),
		},
	}
	genesisState.InitGlobalMeta = globalModel.InitParamList{
		MaxTPS:                       sdk.NewDec(1000),
		ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
		ConsumptionFrictionRate:      types.NewDecFromRat(5, 100),
	}
	result, err := wire.MarshalJSONIndent(lb.cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{AppStateBytes: json.RawMessage(result)})
	lb.Commit()
	assert.True(t, genesisState.GenesisParam.InitFromConfig)
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	accParam, err := lb.paramHolder.GetAccountParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.AccountParam, *accParam)
	postParam, err := lb.paramHolder.GetPostParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.PostParam, *postParam)
	bandwidthParam, err := lb.paramHolder.GetBandwidthParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.BandwidthParam, *bandwidthParam)
	coinDayParam, err := lb.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.CoinDayParam, *coinDayParam)
	validatorParam := lb.paramHolder.GetValidatorParam(ctx)
	assert.Equal(t, genesisState.GenesisParam.ValidatorParam, *validatorParam)
	voteParam, err := lb.paramHolder.GetVoteParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.VoteParam, *voteParam)
	proposalParam, err := lb.paramHolder.GetProposalParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.ProposalParam, *proposalParam)
	globalParam, err := lb.paramHolder.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.GlobalAllocationParam, *globalParam)
}

func TestDistributeInflationToValidators(t *testing.T) {
	lb := newLinoBlockchain(t, 21)

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	remainValidatorPool := types.DecToCoin(
		genesisTotalCoin.ToDec().Mul(
			growthRate.Mul(validatorAllocation)))
	param := lb.paramHolder.GetValidatorParam(ctx)

	expectBaseBalance := coinPerValidator.Minus(param.ValidatorMinDeposit)
	expectBalanceList := make([]types.Coin, 21)
	for i := 0; i < len(expectBalanceList); i++ {
		expectBalanceList[i] = expectBaseBalance
	}
	err := lb.globalManager.DistributeHourlyInflation(ctx)
	if err != nil {
		panic(err)
	}
	if err := lb.valManager.DistributeInflationToValidator(ctx); err != nil {
		panic(err)
	}
	// simulate app
	// hourly inflation
	inflationForValidator :=
		types.DecToCoin(remainValidatorPool.ToDec().Mul(
			types.NewDecFromRat(1, types.HoursPerYear)))
	// expectBalance for all validators
	for i := 0; i < 21; i++ {
		inflation := types.DecToCoin(
			inflationForValidator.ToDec().Quo(sdk.NewDec(int64(21 - i))))
		expectBalanceList[i] = expectBalanceList[i].Plus(inflation)
		inflationForValidator = inflationForValidator.Minus(inflation)
		saving, err :=
			lb.accountManager.GetSavingFromUsername(
				ctx, types.AccountKey("validator"+strconv.Itoa(i)))
		assert.Nil(t, err)
		assert.Equal(t, expectBalanceList[i], saving)
	}
}

func TestFireByzantineValidators(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			Height:  lb.LastBlockHeight() + 1,
			ChainID: "Lino", Time: time.Unix(time.Now().Unix()+200, 0)},
		ByzantineValidators: []abci.Evidence{
			{
				Validator: abci.Validator{
					Address: priv2.PubKey().Address(),
					Power:   1000,
				},
			},
		},
	})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Now()})
	lst := lb.valManager.GetValidatorList(ctx)
	assert.Equal(t, 20, len(lst.Oncall)+len(lst.Standby))
}

func TestDistributeInflationToValidator(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	cases := map[string]struct {
		beforeDistributionInflationPool types.Coin
		pastMinutes                     int64
	}{
		"first distribution": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     types.MinutesPerMonth,
		},
		"last distribution of first year": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     types.MinutesPerMonth * 12,
		},
		"first distribution of second year": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     types.MinutesPerMonth * 13,
		},
		"last distribution of second year": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     types.MinutesPerMonth * 24,
		},
	}
	for testName, cs := range cases {
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		err := lb.globalManager.SetPastMinutes(ctx, cs.pastMinutes)
		if err != nil {
			t.Errorf("%s: failed to set past minutes, got err %v", testName, err)
		}
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)

		err = globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			ValidatorInflationPool: cs.beforeDistributionInflationPool,
		})
		if err != nil {
			t.Errorf("%s: failed to set inflation pool, got err %v", testName, err)
		}

		if err := lb.valManager.DistributeInflationToValidator(ctx); err != nil {
			t.Errorf("%s: failed to set inflation pool, got err %v", testName, err)
		}
		inflationPool, err := globalStore.GetInflationPool(ctx)
		if err != nil {
			t.Errorf("%s: failed to get inflation pool, got err %v", testName, err)
		}
		if !inflationPool.ValidatorInflationPool.IsZero() {
			t.Errorf(
				"%s: diff validator inflation pool, got %v, want %v",
				testName, inflationPool.ValidatorInflationPool,
				types.NewCoinFromInt64(0))
			return
		}
	}
}

func TestHourlyEvent(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	gs := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	globalMeta, err := gs.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	globalAllocation, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)

	expectConsumptionPool := types.NewCoinFromInt64(0)
	for i := 1; i < types.MinutesPerMonth/10; i++ {
		ctx = lb.BaseApp.NewContext(true, abci.Header{Time: time.Unix(int64(i*60), 0)})
		lb.increaseMinute(ctx)

		ctx = lb.BaseApp.NewContext(true, abci.Header{Time: time.Unix(int64(i*60), 0)})
		pastMinutes, err := lb.globalManager.GetPastMinutes(ctx)
		assert.Nil(t, err)
		assert.Equal(t, pastMinutes, int64(i))
		if i%60 == 0 {
			hourlyInflation :=
				types.DecToCoin(
					globalMeta.TotalLinoCoin.ToDec().
						Mul(globalAllocation.GlobalGrowthRate).Mul(types.NewDecFromRat(1, types.HoursPerYear)))
			consumptionMeta, err := gs.GetConsumptionMeta(ctx)
			assert.Nil(t, err)
			expectConsumptionPool =
				expectConsumptionPool.Plus(
					types.DecToCoin(hourlyInflation.ToDec().Mul(globalAllocation.ContentCreatorAllocation)))
			assert.Equal(t, expectConsumptionPool, consumptionMeta.ConsumptionRewardPool)

			inflationPool, _ := gs.GetInflationPool(ctx)
			assert.Equal(t, types.NewCoinFromInt64(0), inflationPool.ValidatorInflationPool)
		}
	}
}

func TestIncreaseMinute(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	validatorParam := lb.paramHolder.GetValidatorParam(ctx)
	minVotingDeposit, _ := validatorParam.ValidatorMinDeposit.ToInt64()
	initStake := minVotingDeposit * 21
	gs := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
	expectLinoStakeStat := globalModel.LinoStakeStat{
		TotalConsumptionFriction: types.NewCoinFromInt64(0),
		TotalLinoStake:           types.NewCoinFromInt64(initStake),
		UnclaimedFriction:        types.NewCoinFromInt64(0),
		UnclaimedLinoStake:       types.NewCoinFromInt64(initStake),
	}
	for i := 1; i < types.MinutesPerMonth/10; i++ {
		// simulate add lino stake and friction at previous block
		ctx := lb.BaseApp.NewContext(true, abci.Header{Time: time.Unix(int64((i-1)*60), 0)})
		err := lb.globalManager.AddLinoStakeToStat(ctx, types.NewCoinFromInt64(1))
		if err != nil {
			panic(err)
		}
		err = lb.globalManager.AddFrictionAndRegisterContentRewardEvent(
			ctx, posttypes.RewardEvent{}, types.NewCoinFromInt64(2), types.NewMiniDollar(1))
		if err != nil {
			panic(err)
		}
		expectLinoStakeStat.TotalConsumptionFriction =
			expectLinoStakeStat.TotalConsumptionFriction.Plus(types.NewCoinFromInt64(2))
		expectLinoStakeStat.UnclaimedFriction =
			expectLinoStakeStat.UnclaimedFriction.Plus(types.NewCoinFromInt64(2))
		expectLinoStakeStat.TotalLinoStake =
			expectLinoStakeStat.TotalLinoStake.Plus(types.NewCoinFromInt64(1))
		expectLinoStakeStat.UnclaimedLinoStake =
			expectLinoStakeStat.UnclaimedLinoStake.Plus(types.NewCoinFromInt64(1))

		// increase minutes after previous block finished
		ctx = lb.BaseApp.NewContext(true, abci.Header{Time: time.Unix(int64(i*60), 0)})
		lb.increaseMinute(ctx)

		ctx = lb.BaseApp.NewContext(true, abci.Header{Time: time.Unix(int64(i*60), 0)})
		pastMinutes, err := lb.globalManager.GetPastMinutes(ctx)
		assert.Nil(t, err)
		assert.Equal(t, pastMinutes, int64(i))
		if i%(60*24) == 0 {
			linoStakeStat, err := gs.GetLinoStakeStat(ctx, int64(i/(60*24)))
			assert.Nil(t, err)
			assert.Equal(t, linoStakeStat.TotalConsumptionFriction, types.NewCoinFromInt64(0))
			assert.Equal(t, linoStakeStat.UnclaimedFriction, types.NewCoinFromInt64(0))
			assert.Equal(t, linoStakeStat.TotalLinoStake, expectLinoStakeStat.TotalLinoStake)
			assert.Equal(t, linoStakeStat.UnclaimedLinoStake, expectLinoStakeStat.UnclaimedLinoStake)
			linoStakeStat, err = gs.GetLinoStakeStat(ctx, int64(i/(60*24)-1))
			assert.Nil(t, err)
			assert.Equal(t, linoStakeStat.TotalConsumptionFriction, expectLinoStakeStat.TotalConsumptionFriction)
			assert.Equal(t, linoStakeStat.UnclaimedFriction, expectLinoStakeStat.UnclaimedFriction)
			assert.Equal(t, linoStakeStat.TotalLinoStake, expectLinoStakeStat.TotalLinoStake)
			assert.Equal(t, linoStakeStat.UnclaimedLinoStake, expectLinoStakeStat.UnclaimedLinoStake)
			expectLinoStakeStat.TotalConsumptionFriction = types.NewCoinFromInt64(0)
			expectLinoStakeStat.UnclaimedFriction = types.NewCoinFromInt64(0)

		}
	}
}

func TestGlobalTime(t *testing.T) {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db, nil)

	genesisState := GenesisState{
		InitCoinPrice: types.NewMiniDollar(1200),
		Accounts: []GenesisAccount{
			{
				Name:           user1,
				Coin:           coinPerValidator,
				ResetKey:       priv1.PubKey(),
				TransactionKey: secp256k1.GenPrivKey().PubKey(),
				AppKey:         secp256k1.GenPrivKey().PubKey(),
				IsValidator:    true,
				ValPubKey:      priv2.PubKey(),
			},
		},
	}

	result, err := wire.MarshalJSONIndent(lb.cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{AppStateBytes: json.RawMessage(result)})
	lb.Commit()

	baseTime := time.Now().Unix()

	cases := []struct {
		testName            string
		baseTime            int64
		expectStartTime     int64
		expectPastMintues   int64
		expectLastBlockTime int64
	}{
		{
			testName:            "init start time",
			baseTime:            baseTime,
			expectStartTime:     baseTime,
			expectPastMintues:   0,
			expectLastBlockTime: baseTime,
		},
		{
			testName:            "past minutes",
			baseTime:            baseTime + 61,
			expectStartTime:     baseTime,
			expectPastMintues:   1,
			expectLastBlockTime: baseTime + 61,
		},
		{
			testName:            "past two minutes",
			baseTime:            baseTime + 121,
			expectStartTime:     baseTime,
			expectPastMintues:   2,
			expectLastBlockTime: baseTime + 121,
		},
		{
			testName:            "past an hour minutes",
			baseTime:            baseTime + 3601,
			expectStartTime:     baseTime,
			expectPastMintues:   60,
			expectLastBlockTime: baseTime + 3601,
		},
	}
	for _, cs := range cases {
		lb := NewLinoBlockchain(logger, db, nil)
		// XXX(yumin): db is shared among all cases, so height is 2..3
		lb.BeginBlock(abci.RequestBeginBlock{
			Header: abci.Header{Height: lb.LastBlockHeight() + 1, ChainID: "Lino", Time: time.Unix(cs.baseTime, 0)}})
		lb.EndBlock(abci.RequestEndBlock{})
		lb.Commit()
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		startTime, err := lb.globalManager.GetChainStartTime(ctx)
		if err != nil {
			t.Errorf("%s: failed to get chain start time, got err %v", cs.testName, err)
		}
		pastMinutes, err := lb.globalManager.GetPastMinutes(ctx)
		if err != nil {
			t.Errorf("%s:failed to get past minutes, got err %v", cs.testName, err)
		}
		lastBlockTime, err := lb.globalManager.GetLastBlockTime(ctx)
		if err != nil {
			t.Errorf("%s:failed to get last block time, got err %v", cs.testName, err)
		}
		assert.Equal(t, cs.expectStartTime, startTime)
		assert.Equal(t, cs.expectPastMintues, pastMinutes)
		assert.Equal(t, cs.expectLastBlockTime, lastBlockTime)
	}
}
