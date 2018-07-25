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
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/param"
	globalModel "github.com/lino-network/lino/x/global/model"
)

var (
	user1 = "validator0"
	priv1 = crypto.GenPrivKeySecp256k1()
	addr1 = priv1.PubKey().Address()
	priv2 = crypto.GenPrivKeySecp256k1()
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
		TransactionKey: crypto.GenPrivKeySecp256k1().PubKey(),
		AppKey:         crypto.GenPrivKeySecp256k1().PubKey(),
		IsValidator:    true,
		ValPubKey:      priv2.PubKey(),
	}
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	for i := 1; i < numOfValidators; i++ {
		genesisAcc := GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Lino:           LNOPerValidator,
			ResetKey:       crypto.GenPrivKeySecp256k1().PubKey(),
			TransactionKey: crypto.GenPrivKeySecp256k1().PubKey(),
			AppKey:         crypto.GenPrivKeySecp256k1().PubKey(),
			IsValidator:    true,
			ValPubKey:      crypto.GenPrivKeySecp256k1().PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
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
		{"lino", "9000000000", crypto.GenPrivKeySecp256k1().PubKey(),
			crypto.GenPrivKeySecp256k1().PubKey(), crypto.GenPrivKeySecp256k1().PubKey(),
			true, crypto.GenPrivKeySecp256k1().PubKey()},
		{"genesis", "500000000", crypto.GenPrivKeySecp256k1().PubKey(),
			crypto.GenPrivKeySecp256k1().PubKey(), crypto.GenPrivKeySecp256k1().PubKey(),
			true, crypto.GenPrivKeySecp256k1().PubKey()},
		{"nonvalidator", "500000000", crypto.GenPrivKeySecp256k1().PubKey(),
			crypto.GenPrivKeySecp256k1().PubKey(), crypto.GenPrivKeySecp256k1().PubKey(),
			false, crypto.GenPrivKeySecp256k1().PubKey()},
		{"developer", "500000000", crypto.GenPrivKeySecp256k1().PubKey(),
			crypto.GenPrivKeySecp256k1().PubKey(), crypto.GenPrivKeySecp256k1().PubKey(),
			false, crypto.GenPrivKeySecp256k1().PubKey()},
		{"infra", "500000000", crypto.GenPrivKeySecp256k1().PubKey(),
			crypto.GenPrivKeySecp256k1().PubKey(), crypto.GenPrivKeySecp256k1().PubKey(),
			false, crypto.GenPrivKeySecp256k1().PubKey()},
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

func TestDistributeInflationToDeveloper(t *testing.T) {
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
		err := lb.developerManager.RegisterDeveloper(ctx, "Lino", types.NewCoinFromInt64(1000000*types.Decimals), "", "", "")
		if err != nil {
			t.Errorf("%s: failed to register developer, got err %v", testName, err)
		}

		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err = globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
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
			assert.Equal(t, expectConsumptionPool, consumptionMeta.ConsumptionRewardPool)
		}
	}
}
