package app

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/param"
	devModel "github.com/lino-network/lino/x/developer/model"
	globalModel "github.com/lino-network/lino/x/global/model"
	infraModel "github.com/lino-network/lino/x/infra/model"
)

var (
	user1 = "validator0"
	priv1 = secp256k1.GenPrivKey()
	addr1 = priv1.PubKey().Address()
	priv2 = secp256k1.GenPrivKey()
	addr2 = priv2.PubKey().Address()

	genesisTotalLino    types.LNO  = "10000000000"
	genesisTotalCoin    types.Coin = types.NewCoinFromInt64(10000000000 * types.Decimals)
	LNOPerValidator     types.LNO  = "100000000"
	growthRate          sdk.Rat    = sdk.NewRat(98, 1000)
	validatorAllocation sdk.Rat    = sdk.NewRat(10, 100)
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
		Accounts: []GenesisAccount{},
	}

	// Generate 21 validators
	genesisAcc := GenesisAccount{
		Name:           user1,
		Lino:           LNOPerValidator,
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
			Lino:           LNOPerValidator,
			ResetKey:       secp256k1.GenPrivKey().PubKey(),
			TransactionKey: secp256k1.GenPrivKey().PubKey(),
			AppKey:         secp256k1.GenPrivKey().PubKey(),
			IsValidator:    true,
			ValPubKey:      secp256k1.GenPrivKey().PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}
	genesisState.InitGlobalMeta = globalModel.InitParamList{
		GrowthRate: sdk.NewRat(98, 1000),
		Ceiling:    sdk.NewRat(98, 1000),
		Floor:      sdk.NewRat(3, 100),
		MaxTPS:     sdk.NewRat(1000),
		ConsumptionFreezingPeriodHr: 7 * 24,
		ConsumptionFrictionRate:     sdk.NewRat(5, 100),
	}

	result, err := wire.MarshalJSONIndent(lb.cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{AppStateBytes: json.RawMessage(result)})
	lb.Commit()

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{ChainID: "Lino"}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	return lb
}

func TestGenesisAcc(t *testing.T) {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db, nil)

	accs := []struct {
		genesisAccountName string
		numOfLino          types.LNO
		resetKey           crypto.PubKey
		transactionKey     crypto.PubKey
		appKey             crypto.PubKey
		isValidator        bool
		valPubKey          crypto.PubKey
	}{
		{"lino", "9000000000", secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			true, secp256k1.GenPrivKey().PubKey()},
		{"genesis", "500000000", secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			true, secp256k1.GenPrivKey().PubKey()},
		{"nonvalidator", "500000000", secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			false, secp256k1.GenPrivKey().PubKey()},
		{"developer", "500000000", secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			false, secp256k1.GenPrivKey().PubKey()},
		{"infra", "500000000", secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			false, secp256k1.GenPrivKey().PubKey()},
	}
	genesisState := GenesisState{
		Accounts: []GenesisAccount{},
	}
	for _, acc := range accs {
		genesisAcc := GenesisAccount{
			Name:           acc.genesisAccountName,
			Lino:           acc.numOfLino,
			ResetKey:       acc.resetKey,
			TransactionKey: acc.transactionKey,
			AppKey:         acc.appKey,
			IsValidator:    acc.isValidator,
			ValPubKey:      acc.valPubKey,
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}
	genesisAppDeveloper := GenesisAppDeveloper{
		Name:        "developer",
		Deposit:     "1000000",
		Website:     "https://lino.network/",
		Description: "",
		AppMetaData: "",
	}
	genesisInfraProvider := GenesisInfraProvider{
		Name: "infra",
	}
	genesisState.Developers = append(genesisState.Developers, genesisAppDeveloper)
	genesisState.Infra = append(genesisState.Infra, genesisInfraProvider)
	result, err := wire.MarshalJSONIndent(lb.cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{AppStateBytes: json.RawMessage(result)})
	lb.Commit()

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	for _, acc := range accs {
		expectBalance, err := types.LinoToCoin(acc.numOfLino)
		assert.Nil(t, err)
		if acc.isValidator {
			param, _ := lb.paramHolder.GetValidatorParam(ctx)
			expectBalance = expectBalance.Minus(
				param.ValidatorMinCommitingDeposit.Plus(param.ValidatorMinVotingDeposit))
		}
		if acc.genesisAccountName == "developer" {
			param, _ := lb.paramHolder.GetDeveloperParam(ctx)
			expectBalance = expectBalance.Minus(param.DeveloperMinDeposit)
		}
		saving, err :=
			lb.accountManager.GetSavingFromBank(ctx, types.AccountKey(acc.genesisAccountName))
		assert.Nil(t, err)
		assert.Equal(t, expectBalance, saving)
	}
}

func TestGenesisFromConfig(t *testing.T) {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db, nil)
	genesisState := GenesisState{
		Accounts: []GenesisAccount{},
	}
	genesisState.GenesisParam = GenesisParam{
		true,
		param.EvaluateOfContentValueParam{
			ConsumptionTimeAdjustBase:      3153600,
			ConsumptionTimeAdjustOffset:    5,
			NumOfConsumptionOnAuthorOffset: 7,
			TotalAmountOfConsumptionBase:   1000 * types.Decimals,
			TotalAmountOfConsumptionOffset: 5,
			AmountOfConsumptionExponent:    sdk.NewRat(8, 10),
		},
		param.GlobalAllocationParam{
			InfraAllocation:          sdk.NewRat(20, 100),
			ContentCreatorAllocation: sdk.NewRat(65, 100),
			DeveloperAllocation:      sdk.NewRat(10, 100),
			ValidatorAllocation:      sdk.NewRat(5, 100),
		},
		param.InfraInternalAllocationParam{
			StorageAllocation: sdk.NewRat(50, 100),
			CDNAllocation:     sdk.NewRat(50, 100),
		},
		param.VoteParam{
			VoterMinDeposit:               types.NewCoinFromInt64(2000 * types.Decimals),
			VoterMinWithdraw:              types.NewCoinFromInt64(2 * types.Decimals),
			DelegatorMinWithdraw:          types.NewCoinFromInt64(2 * types.Decimals),
			VoterCoinReturnIntervalHr:     int64(7 * 24),
			VoterCoinReturnTimes:          int64(7),
			DelegatorCoinReturnIntervalHr: int64(7 * 24),
			DelegatorCoinReturnTimes:      int64(7),
		},
		param.ProposalParam{
			ContentCensorshipDecideHr:   int64(24 * 7),
			ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
			ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
			ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

			ChangeParamDecideHr:   int64(24 * 7),
			ChangeParamPassRatio:  sdk.NewRat(70, 100),
			ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
			ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

			ProtocolUpgradeDecideHr:   int64(24 * 7),
			ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
			ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
			ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
		},
		param.DeveloperParam{
			DeveloperMinDeposit:           types.NewCoinFromInt64(1000000 * types.Decimals),
			DeveloperCoinReturnIntervalHr: int64(7 * 24),
			DeveloperCoinReturnTimes:      int64(7),
		},
		param.ValidatorParam{
			ValidatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
			ValidatorMinVotingDeposit:     types.NewCoinFromInt64(300000 * types.Decimals),
			ValidatorMinCommitingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
			ValidatorCoinReturnIntervalHr: int64(7 * 24),
			ValidatorCoinReturnTimes:      int64(7),
			PenaltyMissVote:               types.NewCoinFromInt64(20000 * types.Decimals),
			PenaltyMissCommit:             types.NewCoinFromInt64(200 * types.Decimals),
			PenaltyByzantine:              types.NewCoinFromInt64(1000000 * types.Decimals),
			ValidatorListSize:             int64(21),
			AbsentCommitLimitation:        int64(600), // 30min
		},
		param.CoinDayParam{
			DaysToRecoverCoinDayStake:    int64(7),
			SecondsToRecoverCoinDayStake: int64(7 * 24 * 3600),
		},
		param.BandwidthParam{
			SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
			CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		},
		param.AccountParam{
			MinimumBalance:             types.NewCoinFromInt64(1 * types.Decimals),
			RegisterFee:                types.NewCoinFromInt64(0),
			FirstDepositFullStakeLimit: types.NewCoinFromInt64(0),
		},
		param.PostParam{
			ReportOrUpvoteInterval: 24 * 3600,
		},
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
	validatorParam, err := lb.paramHolder.GetValidatorParam(ctx)
	assert.Nil(t, err)
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
	infraAllocationParam, err := lb.paramHolder.GetInfraInternalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, genesisState.GenesisParam.InfraInternalAllocationParam, *infraAllocationParam)
}

func TestDistributeInflationToValidators(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	baseTime := time.Now().Unix()
	remainValidatorPool := types.RatToCoin(
		genesisTotalCoin.ToRat().Mul(
			growthRate.Mul(validatorAllocation)))
	coinPerValidator, _ := types.LinoToCoin(LNOPerValidator)
	param, _ := lb.paramHolder.GetValidatorParam(ctx)

	expectBaseBalance := coinPerValidator.Minus(
		param.ValidatorMinCommitingDeposit.Plus(param.ValidatorMinVotingDeposit))
	expectBalanceList := make([]types.Coin, 21)
	for i := 0; i < len(expectBalanceList); i++ {
		expectBalanceList[i] = expectBaseBalance
	}

	testPastMinutes := int64(0)
	for i := baseTime; i < baseTime+3600*1; i += 50 {
		lb.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "Lino", Time: baseTime + i}})
		lb.EndBlock(abci.RequestEndBlock{})
		lb.Commit()
		// simulate app
		if (baseTime+int64(i)-lb.chainStartTime)/int64(60) > testPastMinutes {
			testPastMinutes += 1
			if testPastMinutes%60 == 0 {
				// hourly inflation
				inflationForValidator :=
					types.RatToCoin(remainValidatorPool.ToRat().Mul(
						sdk.NewRat(1, types.HoursPerYear-lb.pastMinutes/60+1)))
				remainValidatorPool = remainValidatorPool.Minus(inflationForValidator)
				// expectBalance for all validators
				ctx := lb.BaseApp.NewContext(true, abci.Header{})
				for i := 0; i < 21; i++ {
					inflation := types.RatToCoin(
						inflationForValidator.ToRat().Quo(sdk.NewRat(int64(21 - i))))
					expectBalanceList[i] = expectBalanceList[i].Plus(inflation)
					saving, err :=
						lb.accountManager.GetSavingFromBank(
							ctx, types.AccountKey("validator"+strconv.Itoa(i)))
					assert.Nil(t, err)
					assert.Equal(t, expectBalanceList[i], saving)
				}
			}
		}
	}
}

func TestFireByzantineValidators(t *testing.T) {
	lb := newLinoBlockchain(t, 21)

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: "Lino", Time: time.Now().Unix()},
		ByzantineValidators: []abci.Evidence{
			abci.Evidence{Validator: abci.Validator{PubKey: tmtypes.TM2PB.PubKey(priv2.PubKey())}}}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	lst, err := lb.valManager.GetValidatorList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 20, len(lst.OncallValidators))
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
		lb.pastMinutes = cs.pastMinutes
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)

		err := globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			ValidatorInflationPool: cs.beforeDistributionInflationPool,
		})
		if err != nil {
			t.Errorf("%s: failed to set inflation pool, got err %v", testName, err)
		}

		lb.distributeInflationToValidator(ctx)
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

func TestDistributeInflationToInfraProvider(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	cases := map[string]struct {
		beforeDistributionInflationPool types.Coin
		pastMinutes                     int64
		numberOfInfraProvider           int
		consumptionList                 []int64
	}{
		"first distribution": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     types.MinutesPerMonth,
			numberOfInfraProvider:           1,
			consumptionList:                 []int64{0},
		},
		"test distribution need to be rounded case": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     types.MinutesPerMonth,
			numberOfInfraProvider:           3,
			consumptionList:                 []int64{0, 0, 0},
		},
		"test distribution based on consumption": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     types.MinutesPerMonth,
			numberOfInfraProvider:           3,
			consumptionList:                 []int64{10, 0, 20},
		},
	}
	for testName, cs := range cases {
		lb := newLinoBlockchain(t, 21)
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		infraStorage := infraModel.NewInfraProviderStorage(lb.CapKeyInfraStore)
		totalWeight := int64(0)
		for i := 0; i < cs.numberOfInfraProvider; i++ {
			err := lb.accountManager.CreateAccount(
				ctx, "", types.AccountKey("infra"+strconv.Itoa(i)),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), types.NewCoinFromInt64(0))
			if err != nil {
				t.Errorf("%s: failed to register account, got err %v", testName, err)
			}
			err = lb.infraManager.RegisterInfraProvider(ctx, types.AccountKey("infra"+strconv.Itoa(i)))
			if err != nil {
				t.Errorf("%s: failed to register infra provider, got err %v", testName, err)
			}
			infra, _ := infraStorage.GetInfraProvider(ctx, types.AccountKey("infra"+strconv.Itoa(i)))
			infra.Usage = cs.consumptionList[i]
			infraStorage.SetInfraProvider(ctx, types.AccountKey("infra"+strconv.Itoa(i)), infra)
			totalWeight = totalWeight + cs.consumptionList[i]
		}
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err := globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			InfraInflationPool: cs.beforeDistributionInflationPool,
		})
		if err != nil {
			t.Errorf("%s: failed to set inflation pool, got err %v", testName, err)
		}

		lb.distributeInflationToInfraProvider(ctx)
		inflationPool, err := globalStore.GetInflationPool(ctx)
		if err != nil {
			t.Errorf("%s: failed to get inflation pool, got err %v", testName, err)
		}

		if !inflationPool.InfraInflationPool.IsZero() {
			t.Errorf(
				"%s: diff infra inflation pool, got %v, want %v",
				testName, inflationPool.InfraInflationPool,
				types.NewCoinFromInt64(0))
			return
		}

		actualInflation := types.NewCoinFromInt64(0)
		for i := 0; i < cs.numberOfInfraProvider; i++ {
			saving, err :=
				lb.accountManager.GetSavingFromBank(
					ctx, types.AccountKey("infra"+strconv.Itoa(i)))
			assert.Nil(t, err)
			inflation := types.NewCoinFromInt64(0)
			if totalWeight == 0 {
				inflation =
					types.RatToCoin(
						sdk.NewRat(1, int64(cs.numberOfInfraProvider)).Round(types.PrecisionFactor).
							Mul(cs.beforeDistributionInflationPool.ToRat()))
			} else {
				inflation =
					types.RatToCoin(
						sdk.NewRat(cs.consumptionList[i], totalWeight).Round(types.PrecisionFactor).
							Mul(cs.beforeDistributionInflationPool.ToRat()))
			}
			if i == (cs.numberOfInfraProvider - 1) {
				inflation = cs.beforeDistributionInflationPool.Minus(actualInflation)
			}
			actualInflation = actualInflation.Plus(inflation)
			if !saving.IsEqual(inflation) {
				t.Errorf(
					"%s: diff inflation for %v, got %v, want %v",
					testName, "dev"+strconv.Itoa(i), inflation,
					saving)
				return
			}
			infra, err := infraStorage.GetInfraProvider(ctx, types.AccountKey("infra"+strconv.Itoa(i)))
			assert.Nil(t, err)
			assert.Equal(t, infra.Usage, int64(0))
		}
	}
	for testName, cs := range cases {
		lb.pastMinutes = cs.pastMinutes
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		err := lb.infraManager.RegisterInfraProvider(ctx, "Lino")
		if err != nil {
			t.Errorf("%s: failed to register infra provider, got err %v", testName, err)
		}

		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err = globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			InfraInflationPool: cs.beforeDistributionInflationPool,
		})
		if err != nil {
			t.Errorf("%s: failed to set inflation pool, got err %v", testName, err)
		}

		lb.distributeInflationToInfraProvider(ctx)
		inflationPool, err := globalStore.GetInflationPool(ctx)
		if err != nil {
			t.Errorf("%s: failed to get inflation pool, got err %v", testName, err)
		}

		if !inflationPool.InfraInflationPool.IsZero() {
			t.Errorf(
				"%s: diff infra inflation pool, got %v, want %v",
				testName, inflationPool.InfraInflationPool,
				types.NewCoinFromInt64(0))
			return
		}
	}
}

func TestDistributeInflationToDevelopers(t *testing.T) {
	cases := map[string]struct {
		beforeDistributionInflationPool types.Coin
		pastMinutes                     int64
		numberOfDevelopers              int
		consumptionList                 []types.Coin
	}{
		"distribute to one developer with zero consumption": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			numberOfDevelopers:              1,
			pastMinutes:                     types.MinutesPerMonth,
			consumptionList:                 []types.Coin{types.NewCoinFromInt64(0)},
		},
		"distribute to five developers with zero consumption": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			numberOfDevelopers:              5,
			pastMinutes:                     types.MinutesPerMonth,
			consumptionList: []types.Coin{
				types.NewCoinFromInt64(0),
				types.NewCoinFromInt64(0),
				types.NewCoinFromInt64(0),
				types.NewCoinFromInt64(0),
				types.NewCoinFromInt64(0)},
		},
		"test inflation need to be rounded case": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(100 * types.Decimals),
			numberOfDevelopers:              3,
			pastMinutes:                     types.MinutesPerMonth,
			consumptionList: []types.Coin{
				types.NewCoinFromInt64(0),
				types.NewCoinFromInt64(0),
				types.NewCoinFromInt64(0),
			},
		},
		"test different consumption case": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(100 * types.Decimals),
			numberOfDevelopers:              3,
			pastMinutes:                     types.MinutesPerMonth,
			consumptionList: []types.Coin{
				types.NewCoinFromInt64(1000 * types.Decimals),
				types.NewCoinFromInt64(2000 * types.Decimals),
				types.NewCoinFromInt64(20),
			},
		},
	}
	for testName, cs := range cases {
		lb := newLinoBlockchain(t, 21)
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		devStorage := devModel.NewDeveloperStorage(lb.CapKeyDeveloperStore)
		totalConsumption := types.NewCoinFromInt64(0)
		for i := 0; i < cs.numberOfDevelopers; i++ {
			err := lb.accountManager.CreateAccount(
				ctx, "", types.AccountKey("dev"+strconv.Itoa(i)),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), types.NewCoinFromInt64(0))
			if err != nil {
				t.Errorf("%s: failed to register account, got err %v", testName, err)
			}
			err = lb.developerManager.RegisterDeveloper(
				ctx, types.AccountKey("dev"+strconv.Itoa(i)), types.NewCoinFromInt64(1000000*types.Decimals), "", "", "")
			if err != nil {
				t.Errorf("%s: failed to register developer, got err %v", testName, err)
			}
			developer, _ := devStorage.GetDeveloper(ctx, types.AccountKey("dev"+strconv.Itoa(i)))
			developer.AppConsumption = cs.consumptionList[i]
			devStorage.SetDeveloper(ctx, types.AccountKey("dev"+strconv.Itoa(i)), developer)
			totalConsumption = totalConsumption.Plus(cs.consumptionList[i])
		}
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err := globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			DeveloperInflationPool: cs.beforeDistributionInflationPool,
		})
		if err != nil {
			t.Errorf("%s: failed to set inflation pool, got err %v", testName, err)
		}

		lb.distributeInflationToDeveloper(ctx)
		inflationPool, err := globalStore.GetInflationPool(ctx)
		if err != nil {
			t.Errorf("%s: failed to get inflation pool, got err %v", testName, err)
		}

		if !inflationPool.DeveloperInflationPool.IsZero() {
			t.Errorf(
				"%s: diff developer inflation pool, got %v, want %v",
				testName, inflationPool.DeveloperInflationPool,
				types.NewCoinFromInt64(0))
			return
		}

		actualInflation := types.NewCoinFromInt64(0)
		for i := 0; i < cs.numberOfDevelopers; i++ {
			saving, err :=
				lb.accountManager.GetSavingFromBank(
					ctx, types.AccountKey("dev"+strconv.Itoa(i)))
			assert.Nil(t, err)
			inflation := types.NewCoinFromInt64(0)
			if totalConsumption.IsZero() {
				inflation =
					types.RatToCoin(
						sdk.NewRat(1, int64(len(cs.consumptionList))).Round(types.PrecisionFactor).
							Mul(cs.beforeDistributionInflationPool.ToRat()))
			} else {
				inflation =
					types.RatToCoin(
						cs.consumptionList[i].ToRat().
							Quo(totalConsumption.ToRat()).Round(types.PrecisionFactor).
							Mul(cs.beforeDistributionInflationPool.ToRat()))
			}
			if i == (cs.numberOfDevelopers - 1) {
				inflation = cs.beforeDistributionInflationPool.Minus(actualInflation)
			}
			actualInflation = actualInflation.Plus(inflation)
			if !saving.IsEqual(inflation) {
				t.Errorf(
					"%s: diff inflation for %v, got %v, want %v",
					testName, "dev"+strconv.Itoa(i), inflation,
					saving)
				return
			}
			developer, err := devStorage.GetDeveloper(ctx, types.AccountKey("dev"+strconv.Itoa(i)))
			assert.Nil(t, err)
			assert.True(t, developer.AppConsumption.IsZero())
		}
	}
}

func TestIncreaseMinute(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	gs := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	globalMeta, err := gs.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	globalAllocation, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)

	inflation := globalMeta.AnnualInflation
	expectConsumptionPool := types.NewCoinFromInt64(0)
	expectInfraPool := types.NewCoinFromInt64(0)
	for i := 0; i < types.MinutesPerMonth/10; i++ {
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		lb.increaseMinute(ctx)
		if i > 0 && i%60 == 0 {
			hourlyInflation :=
				types.RatToCoin(
					inflation.ToRat().Mul(sdk.NewRat(1, types.HoursPerYear-int64(i/60-1))))
			inflation = inflation.Minus(hourlyInflation)
			consumptionMeta, err := gs.GetConsumptionMeta(ctx)
			assert.Nil(t, err)
			expectConsumptionPool =
				expectConsumptionPool.Plus(
					types.RatToCoin(hourlyInflation.ToRat().Mul(globalAllocation.ContentCreatorAllocation)))
			expectInfraPool =
				expectInfraPool.Plus(
					types.RatToCoin(hourlyInflation.ToRat().Mul(globalAllocation.InfraAllocation)))
			assert.Equal(t, expectConsumptionPool, consumptionMeta.ConsumptionRewardPool)

			inflationPool, err := gs.GetInflationPool(ctx)
			assert.Equal(t, types.NewCoinFromInt64(0), inflationPool.ValidatorInflationPool)

			if i%types.MinutesPerMonth == 0 {
				expectInfraPool = types.NewCoinFromInt64(0)
			}

			assert.Equal(t, expectInfraPool, inflationPool.InfraInflationPool)
		}
	}
}
