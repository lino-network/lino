package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/param"
	simAcc "github.com/lino-network/lino/sim/account"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	globalModel "github.com/lino-network/lino/x/global/model"
	"github.com/lino-network/lino/x/post"
	"github.com/syndtr/goleveldb/leveldb"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

// construct some global keys and addrs.
var (
	DefaultNumOfVal  int64 = 21
	genesisTotalCoin       = types.NewCoinFromInt64(10000000000 * types.Decimals)
	coinPerValidator       = types.NewCoinFromInt64(100000000 * types.Decimals)
	GenesisUser            = "genesis"

	signingValidatorList = []abci.SigningValidator{}

	SimIntervalSec         int64 = 120
	Pool                   *globalModel.InflationPool
	GlobalMeta             *globalModel.GlobalMeta
	ConsumptionMeta        *globalModel.ConsumptionMeta
	GlobalAllocationParam  *param.GlobalAllocationParam
	NewRegistrationArrival float64                 = 0.01
	EventMap               map[int64][]types.Event = map[int64][]types.Event{}
)

func NewSimLinoBlockchain() *app.LinoBlockchain {
	logger, db := loggerAndDB()
	lb := app.NewLinoBlockchain(logger, db, nil)
	genesisState := app.GenesisState{
		Accounts: []app.GenesisAccount{},
		InitGlobalMeta: globalModel.InitParamList{
			MaxTPS: sdk.NewRat(1000),
			ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
			ConsumptionFrictionRate:      sdk.NewRat(5, 100),
		},
	}

	// Generate 21 validators
	for i := 0; i < int(DefaultNumOfVal); i++ {
		sa := simAcc.NewSimAcc("validator"+strconv.Itoa(i), coinPerValidator)
		sa.IsValidator = true
		sa.ValidatorPrivKey = ed25519.GenPrivKey()
		sa.IsVoter = true
		genesisAcc := app.GenesisAccount{
			Name:           string(sa.Username),
			Coin:           coinPerValidator,
			ResetKey:       sa.ResetPrivKey.PubKey(),
			TransactionKey: sa.TxPrivKey.PubKey(),
			AppKey:         sa.AppPrivKey.PubKey(),
			IsValidator:    true,
			ValPubKey:      sa.ValidatorPrivKey.PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
		simAcc.AccountList = append(simAcc.AccountList, sa)
		simAcc.VoterList = append(simAcc.VoterList, sa)
		signingValidatorList = append(signingValidatorList, abci.SigningValidator{Validator: abci.Validator{Address: sa.ValidatorPrivKey.PubKey().Address()}, SignedLastBlock: true})
	}
	genesisAccCoin := genesisTotalCoin.Minus(types.RatToCoin(coinPerValidator.ToRat().Mul(sdk.NewRat(21))))
	simAcc.GenesisAccount = *simAcc.NewSimAcc(GenesisUser, genesisAccCoin)
	genesisAcc := app.GenesisAccount{
		Name:           GenesisUser,
		Coin:           genesisAccCoin,
		ResetKey:       simAcc.GenesisAccount.ResetPrivKey.PubKey(),
		TransactionKey: simAcc.GenesisAccount.TxPrivKey.PubKey(),
		AppKey:         simAcc.GenesisAccount.AppPrivKey.PubKey(),
		IsValidator:    false,
	}
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	simAcc.GenesisAccount.IsDeveloper = true
	genesisAppDeveloper := app.GenesisAppDeveloper{
		Name:    GenesisUser,
		Deposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}
	simAcc.GenesisAccount.IsInfraProvider = true
	genesisInfraProvider := app.GenesisInfraProvider{
		Name: GenesisUser,
	}
	genesisState.Developers = []app.GenesisAppDeveloper{genesisAppDeveloper}
	genesisState.Infra = []app.GenesisInfraProvider{genesisInfraProvider}
	simAcc.AccountList = append(simAcc.AccountList, &simAcc.GenesisAccount)

	appState, _ := wire.MarshalJSONIndent(app.MakeCodec(), genesisState)
	lb.InitChain(abci.RequestInitChain{AppStateBytes: appState})
	lb.Commit()

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{ChainID: "Lino", Time: time.Now()}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	return lb
}

type MuteWriter struct {
}

func (mw MuteWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func loggerAndDB() (log.Logger, dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return logger, db
}

func main() {
	db, err := leveldb.OpenFile("tmp/db", nil)
	if err != nil {
		fmt.Println("Create db failed")
		return
	}
	defer db.Close()
	lb := NewSimLinoBlockchain()
	if lb == nil {
		fmt.Println("Create app failed")
		return
	}
	pastHour := 0
	pastMonth := 0
	startTime := time.Now().Unix()
	lastRegistration := startTime
	for i := startTime; i < startTime+types.HoursPerYear*3600*2; i += SimIntervalSec {
		if (i-startTime) != 0 && (i-startTime)%3600 == 0 {
			pastHour++
			executeHourlyEvent(i, startTime, pastHour)
		}

		if (i-startTime) != 0 && (i-startTime)%(types.MinutesPerMonth*60) == 0 {
			pastMonth++
			executeMonthlyEvent(i, startTime, pastMonth)
		}
		if (i-startTime)%(types.HoursPerYear*3600) == 0 {
			ctx := lb.BaseApp.NewContext(true, abci.Header{})
			globalStorage := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
			ph := param.NewParamHolder(lb.CapKeyParamStore)
			var growthRate sdk.Rat
			GlobalMeta, _ = globalStorage.GetGlobalMeta(ctx)
			ConsumptionMeta, _ = globalStorage.GetConsumptionMeta(ctx)
			GlobalAllocationParam, _ = ph.GetGlobalAllocationParam(ctx)
			if (i - startTime) != 0 {
				if GlobalMeta.LastYearCumulativeConsumption.IsZero() {
					growthRate = GlobalAllocationParam.GlobalGrowthRate
				} else {
					// growthRate = (consumption this year - consumption last year) / consumption last year
					lastYearConsumptionRat := GlobalMeta.LastYearCumulativeConsumption.ToRat()
					thisYearConsumptionRat := GlobalMeta.CumulativeConsumption.ToRat()
					consumptionIncrement := thisYearConsumptionRat.Sub(lastYearConsumptionRat)

					growthRate = consumptionIncrement.Quo(lastYearConsumptionRat).Round(types.PrecisionFactor)
				}
				GlobalMeta.LastYearCumulativeConsumption = GlobalMeta.CumulativeConsumption
				GlobalMeta.CumulativeConsumption = types.NewCoinFromInt64(0)
				GlobalMeta.LastYearTotalLinoCoin = GlobalMeta.TotalLinoCoin
			}
			if growthRate.LT(param.AnnualInflationFloor) {
				growthRate = param.AnnualInflationFloor
			}
			if growthRate.GT(param.AnnualInflationCeiling) {
				growthRate = param.AnnualInflationCeiling
			}
			GlobalAllocationParam.GlobalGrowthRate = growthRate
			Pool = &globalModel.InflationPool{types.NewCoinFromInt64(0), types.NewCoinFromInt64(0), types.NewCoinFromInt64(0)}
		}
		if i != startTime {
			for k := i - SimIntervalSec; k < i; k++ {
				eventList, ok := EventMap[k]
				if !ok {
					continue
				}
				for _, event := range eventList {
					//fmt.Println("loop event")
					switch e := event.(type) {
					case post.RewardEvent:
						for _, acc := range simAcc.AccountList {
							if acc.Username == e.PostAuthor {
								rewardPoolRat := ConsumptionMeta.ConsumptionRewardPool.ToRat()
								actualRewardRat := rewardPoolRat.Mul(e.Evaluate.ToRat().Quo(ConsumptionMeta.ConsumptionWindow.ToRat()).Round(types.PrecisionFactor))
								actualReward := types.RatToCoin(actualRewardRat)
								acc.ActualReward = acc.ActualReward.Plus(actualReward)
								fmt.Println("add actual reward:", acc.Username, actualReward, acc.ActualReward,
									", reward pool:", ConsumptionMeta.ConsumptionRewardPool, e.Evaluate, ConsumptionMeta.ConsumptionWindow,
									"the event time:", k)
								ConsumptionMeta.ConsumptionRewardPool = ConsumptionMeta.ConsumptionRewardPool.Minus(actualReward)
								ConsumptionMeta.ConsumptionWindow = ConsumptionMeta.ConsumptionWindow.Minus(e.Evaluate)
								break
							}
						}
					}
				}
			}
		}
		lb.BeginBlock(abci.RequestBeginBlock{
			LastCommitInfo: abci.LastCommitInfo{Validators: signingValidatorList},
			Header: abci.Header{
				ChainID: "Lino", Time: time.Unix(i, 0), Height: (i - startTime) / 6,
			},
		})

		if len(simAcc.AccountList) < 50 && rand.Float64() < 1-math.Exp(-NewRegistrationArrival*float64((i-lastRegistration)/SimIntervalSec)) {
			lastRegistration = i
			genesisDeposit := rand.Intn(10000)
			if result := simAcc.CreateAccount(
				"account"+strconv.Itoa(len(simAcc.AccountList)), lb, strconv.Itoa(genesisDeposit)); result == true {
				Pool.DeveloperInflationPool = Pool.DeveloperInflationPool.Plus(types.NewCoinFromInt64(1 * types.Decimals))
			}

		}
		for _, acc := range simAcc.AccountList {
			acc.Action(lb)
			event := acc.Donation(lb, i)
			if event == nil {
				continue
			}
			ConsumptionMeta.ConsumptionWindow = ConsumptionMeta.ConsumptionWindow.Plus(event.Evaluate)
			ConsumptionMeta.ConsumptionRewardPool = ConsumptionMeta.ConsumptionRewardPool.Plus(event.Friction)
			eventList, ok := EventMap[i+24*7*3600]
			if !ok {
				EventMap[i+24*7*3600] = []types.Event{*event}
			} else {
				EventMap[i+24*7*3600] = append(eventList, *event)
			}
		}

		lb.EndBlock(abci.RequestEndBlock{})
		lb.Commit()

		if (i-startTime)%36000 == 0 {
			fmt.Println("check status:", (i-startTime)/3600)
			fmt.Println("total Account:", len(simAcc.AccountList))
			fmt.Println("statistic param:", simAcc.StatisticParam)
			ctx := lb.BaseApp.NewContext(true, abci.Header{Time: time.Unix(i, 0), Height: (i - startTime) / 6})
			ph := param.NewParamHolder(lb.CapKeyParamStore)
			accManager := acc.NewAccountManager(lb.CapKeyAccountStore, ph)
			globalStorage := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
			totalStake := types.NewCoinFromInt64(0)
			totalBalance := types.NewCoinFromInt64(0)
			totalExpect := types.NewCoinFromInt64(0)
			for _, acc := range simAcc.AccountList {
				if !acc.CheckSelfBalance(lb) {
					panic(fmt.Sprintf("check balance failed: %v", acc.Username))
				}
				stake, _ := accManager.GetStake(ctx, acc.Username)
				//fmt.Println("check stake in simP:", stake, acc.Username)
				if !stake.IsNotNegative() {
					panic("Negative stake")
				}
				totalStake = totalStake.Plus(stake)
				saving, _ := accManager.GetSavingFromBank(ctx, acc.Username)
				totalBalance = totalBalance.Plus(saving)
				totalExpect = totalExpect.Plus(acc.ExpectCoin)
			}
			fmt.Println("whole net stake:", totalStake, ", whole net balance:", totalBalance)
			fmt.Println("whole net expect:", totalExpect, ", consumption pool:", ConsumptionMeta.ConsumptionRewardPool)
			globalMeta, _ := globalStorage.GetGlobalMeta(ctx)
			fmt.Println("in sim total lino:", totalExpect.Plus(ConsumptionMeta.ConsumptionRewardPool), ", total lino coin:", globalMeta.TotalLinoCoin)
		}
	}
}

func executeHourlyEvent(currentTime, startTime int64, pastHour int) {
	fmt.Println("simulation past hour:", (currentTime-startTime)/3600)
	thisHourInflation :=
		types.RatToCoin(
			GlobalMeta.LastYearTotalLinoCoin.ToRat().
				Mul(GlobalAllocationParam.GlobalGrowthRate).
				Mul(sdk.NewRat(1, types.HoursPerYear)))
	contentCreatorInflation :=
		types.RatToCoin(thisHourInflation.ToRat().Mul(GlobalAllocationParam.ContentCreatorAllocation))
	validatorInflation :=
		types.RatToCoin(thisHourInflation.ToRat().Mul(GlobalAllocationParam.ValidatorAllocation))
	infraInflation :=
		types.RatToCoin(thisHourInflation.ToRat().Mul(GlobalAllocationParam.InfraAllocation))
	developerInflation :=
		thisHourInflation.Minus(contentCreatorInflation).Minus(validatorInflation).Minus(infraInflation)
	var numOfValidators int
	for _, acc := range simAcc.AccountList {
		if acc.IsValidator {
			numOfValidators++
		}
	}
	for i, acc := range simAcc.AccountList {
		if acc.IsValidator {
			expectIncome := types.RatToCoin(validatorInflation.ToRat().Quo(sdk.NewRat(int64(numOfValidators-i), 1)))
			acc.ExpectCoin = acc.ExpectCoin.Plus(expectIncome)
			validatorInflation = validatorInflation.Minus(expectIncome)
			//fmt.Println("in sim, add ", expectIncome, " to validator:", acc.Username)
		}
	}
	Pool.DeveloperInflationPool = Pool.DeveloperInflationPool.Plus(developerInflation)
	Pool.InfraInflationPool = Pool.InfraInflationPool.Plus(infraInflation)
	ConsumptionMeta.ConsumptionRewardPool = ConsumptionMeta.ConsumptionRewardPool.Plus(contentCreatorInflation)
	//fmt.Println("consumption window in sim", ConsumptionMeta.ConsumptionWindow)
}

func executeMonthlyEvent(currentTime, startTime int64, pastMonth int) {
	fmt.Println("simulate pass month:", (currentTime-startTime)/(43830*60))
	simAcc.GenesisAccount.ExpectCoin = simAcc.GenesisAccount.ExpectCoin.Plus(Pool.DeveloperInflationPool)
	Pool.DeveloperInflationPool = types.NewCoinFromInt64(0)
	simAcc.GenesisAccount.ExpectCoin = simAcc.GenesisAccount.ExpectCoin.Plus(Pool.InfraInflationPool)
	Pool.InfraInflationPool = types.NewCoinFromInt64(0)
}
