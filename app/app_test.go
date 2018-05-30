package app

import (
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	globalModel "github.com/lino-network/lino/tx/global/model"
)

var (
	user1 = "validator0"
	priv1 = crypto.GenPrivKeyEd25519()
	addr1 = priv1.PubKey().Address()
	priv2 = crypto.GenPrivKeyEd25519()
	addr2 = priv2.PubKey().Address()

	genesisTotalLino    types.LNO  = "10000000000"
	genesisTotalCoin    types.Coin = types.NewCoinFromInt64(10000000000 * types.Decimals)
	LNOPerValidator     types.LNO  = "100000000"
	growthRate          sdk.Rat    = sdk.NewRat(98, 1000)
	validatorAllocation sdk.Rat    = sdk.NewRat(10, 100)
)

func loggerAndDB() (log.Logger, dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return logger, db
}

func newLinoBlockchain(t *testing.T, numOfValidators int) *LinoBlockchain {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db)

	genesisState := genesis.GenesisState{
		Accounts:  []genesis.GenesisAccount{},
		TotalLino: genesisTotalLino,
	}

	// Generate 21 validators
	genesisAcc := genesis.GenesisAccount{
		Name:           user1,
		Lino:           LNOPerValidator,
		MasterKey:      priv1.PubKey(),
		TransactionKey: crypto.GenPrivKeyEd25519().PubKey(),
		PostKey:        crypto.GenPrivKeyEd25519().PubKey(),
		IsValidator:    true,
		ValPubKey:      priv2.PubKey(),
	}
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	for i := 1; i < numOfValidators; i++ {
		genesisAcc := genesis.GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Lino:           LNOPerValidator,
			MasterKey:      crypto.GenPrivKeyEd25519().PubKey(),
			TransactionKey: crypto.GenPrivKeyEd25519().PubKey(),
			PostKey:        crypto.GenPrivKeyEd25519().PubKey(),
			IsValidator:    true,
			ValPubKey:      crypto.GenPrivKeyEd25519().PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}

	result, err := genesis.GetGenesisJson(genesisState)
	assert.Nil(t, err)

	vals := []abci.Validator{}
	lb.InitChain(abci.RequestInitChain{vals, json.RawMessage(result)})
	lb.Commit()

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{ChainID: "Lino"}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	return lb
}

func TestGenesisAcc(t *testing.T) {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db)

	accs := []struct {
		genesisAccountName string
		numOfLino          types.LNO
		masterKey          crypto.PubKey
		transactionKey     crypto.PubKey
		postKey            crypto.PubKey
		isValidator        bool
		valPubKey          crypto.PubKey
	}{
		{"Lino", "9000000000", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			true, crypto.GenPrivKeyEd25519().PubKey()},
		{"Genesis", "500000000", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			true, crypto.GenPrivKeyEd25519().PubKey()},
		{"NonValidator", "500000000", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			false, crypto.GenPrivKeyEd25519().PubKey()},
	}
	genesisState := genesis.GenesisState{
		Accounts:  []genesis.GenesisAccount{},
		TotalLino: genesisTotalLino,
	}
	for _, acc := range accs {
		genesisAcc := genesis.GenesisAccount{
			Name:           acc.genesisAccountName,
			Lino:           acc.numOfLino,
			MasterKey:      acc.masterKey,
			TransactionKey: acc.transactionKey,
			PostKey:        acc.postKey,
			IsValidator:    acc.isValidator,
			ValPubKey:      acc.valPubKey,
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}

	result, err := genesis.GetGenesisJson(genesisState)
	assert.Nil(t, err)

	vals := []abci.Validator{}
	lb.InitChain(abci.RequestInitChain{vals, json.RawMessage(result)})
	lb.Commit()

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	param, _ := lb.paramHolder.GetValidatorParam(ctx)
	for _, acc := range accs {
		expectBalance, err := types.LinoToCoin(acc.numOfLino)
		assert.Nil(t, err)
		if acc.isValidator {
			expectBalance = expectBalance.Minus(
				param.ValidatorMinCommitingDeposit.Plus(param.ValidatorMinVotingDeposit))
		}
		saving, err :=
			lb.accountManager.GetSavingFromBank(ctx, types.AccountKey(acc.genesisAccountName))
		assert.Nil(t, err)
		assert.Equal(t, expectBalance, saving)
	}

	// reload app and ensure the account is still there
	lb = NewLinoBlockchain(logger, db)
	ctx = lb.BaseApp.NewContext(true, abci.Header{})
	for _, acc := range accs {
		expectBalance, err := types.LinoToCoin(acc.numOfLino)
		assert.Nil(t, err)
		if acc.isValidator {
			expectBalance = expectBalance.Minus(
				param.ValidatorMinCommitingDeposit.Plus(param.ValidatorMinVotingDeposit))
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
	remainValidatorPool, _ := types.RatToCoin(
		new(big.Rat).Mul(genesisTotalCoin.ToRat(),
			new(big.Rat).Mul(growthRate.GetRat(), validatorAllocation.GetRat())))
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
				inflationForValidator, _ :=
					types.RatToCoin(new(big.Rat).Mul(remainValidatorPool.ToRat(),
						big.NewRat(1, types.HoursPerYear-lb.pastMinutes/60+1)))
				remainValidatorPool = remainValidatorPool.Minus(inflationForValidator)
				// expectBalance for all validators
				ctx := lb.BaseApp.NewContext(true, abci.Header{})
				for i := 0; i < 21; i++ {
					inflation, _ := types.RatToCoin(new(big.Rat).
						Quo(inflationForValidator.ToRat(), big.NewRat(int64(21-i), 1)))
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
		ByzantineValidators: []abci.Evidence{abci.Evidence{PubKey: priv2.PubKey().Bytes()}}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	lst, err := lb.valManager.GetValidatorList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 20, len(lst.OncallValidators))
}

func TestDistributeInflationToConsumptionRewardPool(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	cases := map[string]struct {
		beforeDistributionRewardPool    types.Coin
		beforeDistributionInflationPool types.Coin
		pastMinutes                     int64
	}{
		"first distribution": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			beforeDistributionRewardPool:    types.NewCoinFromInt64(0),
			pastMinutes:                     60,
		},
		"end of year distribution": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			beforeDistributionRewardPool:    types.NewCoinFromInt64(0),
			pastMinutes:                     60 * types.HoursPerYear,
		},
		"next year first hour": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			beforeDistributionRewardPool:    types.NewCoinFromInt64(0),
			pastMinutes:                     60*types.HoursPerYear + 60,
		},
		"end of next year": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			beforeDistributionRewardPool:    types.NewCoinFromInt64(0),
			pastMinutes:                     60 * types.HoursPerYear * 2,
		},
	}
	for testName, cs := range cases {
		lb.pastMinutes = cs.pastMinutes
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err := globalStore.SetConsumptionMeta(
			ctx, &globalModel.ConsumptionMeta{
				ConsumptionRewardPool: cs.beforeDistributionRewardPool,
			},
		)
		assert.Nil(t, err)
		err = globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			ContentCreatorInflationPool: cs.beforeDistributionInflationPool,
		})
		assert.Nil(t, err)
		lb.distributeInflationToConsumptionRewardPool(ctx)
		inflationPool, err := globalStore.GetInflationPool(ctx)
		assert.Nil(t, err)
		consumption, err := globalStore.GetConsumptionMeta(ctx)
		assert.Nil(t, err)
		expectInflation, _ := types.RatToCoin(new(big.Rat).Quo(
			cs.beforeDistributionInflationPool.ToRat(),
			new(big.Rat).SetInt64(types.HoursPerYear-lb.getPastHoursMinusOneThisYear())))
		if !cs.beforeDistributionRewardPool.Plus(expectInflation).
			IsEqual(consumption.ConsumptionRewardPool) {
			t.Errorf(
				"%s: expect inflation pool %v, got %v",
				testName, cs.beforeDistributionRewardPool.Plus(expectInflation),
				inflationPool.ContentCreatorInflationPool)
			return
		}
		if !cs.beforeDistributionInflationPool.Minus(expectInflation).
			IsEqual(inflationPool.ContentCreatorInflationPool) {
			t.Errorf(
				"%s: expect inflation pool %v, got %v",
				testName, cs.beforeDistributionRewardPool.Plus(expectInflation),
				inflationPool.ContentCreatorInflationPool)
			return
		}
	}
}

func TestDistributeInflationToValidator(t *testing.T) {
	lb := newLinoBlockchain(t, 21)
	cases := map[string]struct {
		beforeDistributionInflationPool types.Coin
		pastMinutes                     int64
	}{
		"first distribution": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     60,
		},
		"end of year distribution": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     60 * types.HoursPerYear,
		},
		"next year first hour": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     60*types.HoursPerYear + 60,
		},
		"end of next year": {
			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
			pastMinutes:                     60 * types.HoursPerYear * 2,
		},
	}
	for testName, cs := range cases {
		lb.pastMinutes = cs.pastMinutes
		ctx := lb.BaseApp.NewContext(true, abci.Header{})
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err := globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			ValidatorInflationPool: cs.beforeDistributionInflationPool,
		})
		assert.Nil(t, err)
		lb.distributeInflationToValidator(ctx)
		inflationPool, err := globalStore.GetInflationPool(ctx)
		assert.Nil(t, err)
		expectInflation, _ := types.RatToCoin(new(big.Rat).Quo(
			cs.beforeDistributionInflationPool.ToRat(),
			new(big.Rat).SetInt64(types.HoursPerYear-lb.getPastHoursMinusOneThisYear())))
		if !cs.beforeDistributionInflationPool.Minus(expectInflation).
			IsEqual(inflationPool.ValidatorInflationPool) {
			t.Errorf(
				"%s: expect inflation pool %v, got %v",
				testName, cs.beforeDistributionInflationPool.Minus(expectInflation),
				inflationPool.ValidatorInflationPool)
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
		assert.Nil(t, err)
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err = globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			InfraInflationPool: cs.beforeDistributionInflationPool,
		})
		assert.Nil(t, err)
		lb.distributeInflationToInfraProvider(ctx)
		inflationPool, err := globalStore.GetInflationPool(ctx)
		assert.Nil(t, err)
		expectInflation, _ := types.RatToCoin(new(big.Rat).Quo(
			cs.beforeDistributionInflationPool.ToRat(),
			new(big.Rat).SetInt64(12-lb.getPastMonthMinusOneThisYear())))
		if !cs.beforeDistributionInflationPool.Minus(expectInflation).
			IsEqual(inflationPool.InfraInflationPool) {
			t.Errorf(
				"%s: expect inflation pool %v, got %v",
				testName, cs.beforeDistributionInflationPool.Minus(expectInflation),
				inflationPool.InfraInflationPool)
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
		err := lb.developerManager.RegisterDeveloper(ctx, "Lino", types.NewCoinFromInt64(1000000*types.Decimals))
		assert.Nil(t, err)
		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
		err = globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
			DeveloperInflationPool: cs.beforeDistributionInflationPool,
		})
		assert.Nil(t, err)
		lb.distributeInflationToDeveloper(ctx)
		inflationPool, err := globalStore.GetInflationPool(ctx)
		assert.Nil(t, err)
		expectInflation, _ := types.RatToCoin(new(big.Rat).Quo(
			cs.beforeDistributionInflationPool.ToRat(),
			new(big.Rat).SetInt64(12-lb.getPastMonthMinusOneThisYear())))
		if !cs.beforeDistributionInflationPool.Minus(expectInflation).
			IsEqual(inflationPool.DeveloperInflationPool) {
			t.Errorf(
				"%s: expect inflation pool %v, got %v",
				testName, cs.beforeDistributionInflationPool.Minus(expectInflation),
				inflationPool.DeveloperInflationPool)
			return
		}
	}
}
