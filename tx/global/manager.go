package global

import (
	"math"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/global/model"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GlobalManager encapsulates all basic struct
type GlobalManager struct {
	storage     model.GlobalStorage `json:"global_manager"`
	paramHolder param.ParamHolder   `json:"param_holder"`
}

// NewGlobalManager return the global proxy pointer
func NewGlobalManager(key sdk.StoreKey, holder param.ParamHolder) GlobalManager {
	return GlobalManager{
		storage:     model.NewGlobalStorage(key),
		paramHolder: holder,
	}
}

func (gm GlobalManager) WireCodec() *wire.Codec {
	return gm.storage.WireCodec()
}

func (gm GlobalManager) InitGlobalManager(ctx sdk.Context, totalLino types.Coin) sdk.Error {
	allocationParam, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	if err != nil {
		return err
	}
	return gm.storage.InitGlobalState(ctx, totalLino, allocationParam)
}

func (gm GlobalManager) registerEventAtTime(ctx sdk.Context, unixTime int64, event types.Event) sdk.Error {
	if unixTime < ctx.BlockHeader().Time {
		return ErrGlobalManagerRegisterExpiredEvent(unixTime)
	}
	eventList, _ := gm.storage.GetTimeEventList(ctx, unixTime)
	if eventList == nil {
		eventList = &types.TimeEventList{Events: []types.Event{}}
	}
	eventList.Events = append(eventList.Events, event)
	if err := gm.storage.SetTimeEventList(ctx, unixTime, eventList); err != nil {
		return ErrGlobalManagerRegisterEventAtTime(unixTime).TraceCause(err, "")
	}
	return nil
}

func (gm GlobalManager) GetTimeEventListAtTime(ctx sdk.Context, unixTime int64) *types.TimeEventList {
	eventList, _ := gm.storage.GetTimeEventList(ctx, unixTime)
	return eventList
}

func (gm GlobalManager) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	return gm.storage.RemoveTimeEventList(ctx, unixTime)
}

func (gm GlobalManager) GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return sdk.Rat{}, err
	}
	return consumptionMeta.ConsumptionFrictionRate, nil
}

// register reward calculation event at 7 days later
func (gm GlobalManager) AddFrictionAndRegisterContentRewardEvent(
	ctx sdk.Context, event types.Event, friction types.Coin, evaluate types.Coin) sdk.Error {
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Plus(friction)
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Plus(evaluate)

	if err := gm.registerEventAtTime(
		ctx, ctx.BlockHeader().Time+
			(consumptionMeta.ConsumptionFreezingPeriodHr*3600), event); err != nil {
		return err
	}
	if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}
	return nil
}

// register coin return event with a time interval
func (gm GlobalManager) RegisterCoinReturnEvent(
	ctx sdk.Context, events []types.Event, times int64, interval int64) sdk.Error {
	for i := int64(0); i < times; i++ {
		if err := gm.registerEventAtTime(
			ctx, ctx.BlockHeader().Time+(interval*3600*(i+1)), events[i]); err != nil {
			return err
		}
	}
	return nil
}

func (gm GlobalManager) RegisterProposalDecideEvent(
	ctx sdk.Context, decideHr int64, event types.Event) sdk.Error {
	if err := gm.registerEventAtTime(
		ctx, ctx.BlockHeader().Time+(decideHr*3600), event); err != nil {
		return err
	}
	return nil
}

func (gm GlobalManager) RegisterParamChangeEvent(ctx sdk.Context, event types.Event) sdk.Error {
	// param will be changed in one day
	if err := gm.registerEventAtTime(
		ctx, ctx.BlockHeader().Time+(24*3600), event); err != nil {
		return err
	}
	return nil
}

// put hourly inflation to reward pool
func (gm GlobalManager) AddHourlyInflationToRewardPool(ctx sdk.Context, pastHoursThisYear int64) sdk.Error {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return err
	}
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	resRat := pool.ContentCreatorInflationPool.ToRat().
		Mul(sdk.NewRat(1, types.HoursPerYear-pastHoursThisYear+1))
	resCoin := types.RatToCoin(resRat)
	pool.ContentCreatorInflationPool = pool.ContentCreatorInflationPool.Minus(resCoin)
	if err := gm.addTotalLinoCoin(ctx, resCoin); err != nil {
		return err
	}
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return err
	}

	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Plus(resCoin)

	if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}
	return nil
}

// recalculate annually inflation based on consumption growth rate
func (gm GlobalManager) RecalculateAnnuallyInflation(ctx sdk.Context) sdk.Error {
	growthRate, err := gm.getGrowthRate(ctx)
	if err != nil {
		return err
	}
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	allocation, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	if err != nil {
		return err
	}

	infraInflationCoin :=
		globalMeta.TotalLinoCoin.ToRat().Mul(growthRate).Mul(allocation.InfraAllocation)
	contentCreatorCoin :=
		globalMeta.TotalLinoCoin.ToRat().Mul(growthRate).Mul(allocation.ContentCreatorAllocation)
	developerCoin := globalMeta.TotalLinoCoin.ToRat().Mul(growthRate).Mul(allocation.DeveloperAllocation)
	validatorCoin := globalMeta.TotalLinoCoin.ToRat().Mul(growthRate).Mul(allocation.ValidatorAllocation)

	inflationPool := &model.InflationPool{
		InfraInflationPool:          types.RatToCoin(infraInflationCoin),
		ContentCreatorInflationPool: types.RatToCoin(contentCreatorCoin),
		DeveloperInflationPool:      types.RatToCoin(developerCoin),
		ValidatorInflationPool:      types.RatToCoin(validatorCoin),
	}
	if err := gm.storage.SetInflationPool(ctx, inflationPool); err != nil {
		return err
	}
	return nil
}

// get growth rate based on consumption growth rate
func (gm GlobalManager) getGrowthRate(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	var growthRate sdk.Rat
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return sdk.ZeroRat, err
	}
	// if last year cumulative consumption is zero, we use the same growth rate as the last year
	if globalMeta.LastYearCumulativeConsumption.IsZero() {
		growthRate = globalMeta.GrowthRate
	} else {
		// growthRate = (consumption this year - consumption last year) / consumption last year
		lastYearConsumptionRat := globalMeta.LastYearCumulativeConsumption.ToRat()
		thisYearConsumptionRat := globalMeta.CumulativeConsumption.ToRat()
		growthRate =
			(thisYearConsumptionRat.Sub(lastYearConsumptionRat)).Quo(lastYearConsumptionRat).Round(types.PrecisionFactor)
		if growthRate.GT(globalMeta.Ceiling) {
			growthRate = globalMeta.Ceiling
		} else if growthRate.LT(globalMeta.Floor) {
			growthRate = globalMeta.Floor
		}
	}
	globalMeta.LastYearCumulativeConsumption = globalMeta.CumulativeConsumption
	globalMeta.CumulativeConsumption = types.NewCoin(0)
	globalMeta.GrowthRate = growthRate.Round(types.PrecisionFactor)
	if err := gm.storage.SetGlobalMeta(ctx, globalMeta); err != nil {
		return sdk.ZeroRat, err
	}
	return growthRate, nil
}

// after 7 days, one consumption needs to claim its reward from consumption reward pool
func (gm GlobalManager) GetRewardAndPopFromWindow(
	ctx sdk.Context, evaluate types.Coin, penaltyScore sdk.Rat) (types.Coin, sdk.Error) {
	if evaluate.IsZero() {
		return types.NewCoin(0), nil
	}

	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return types.NewCoin(0), ErrGetRewardAndPopFromWindow().TraceCause(err, "")
	}

	// consumptionRatio = this consumption * penalty score) / (total consumption in 7 days window)
	consumptionRatio := evaluate.ToRat().Mul(sdk.OneRat.Sub(penaltyScore)).
		Quo(consumptionMeta.ConsumptionWindow.ToRat()).Round(types.PrecisionFactor)
	// reward = (consumption reward pool) * (consumptionRatio)
	reward := types.RatToCoin(
		consumptionMeta.ConsumptionRewardPool.ToRat().Mul(consumptionRatio))

	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Minus(reward)
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Minus(evaluate)

	if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return types.NewCoin(0), ErrGetRewardAndPopFromWindow().TraceCause(err, "")
	}
	return reward, nil
}

// add consumption to global meta, which is used to compute GDP
func (gm GlobalManager) AddConsumption(ctx sdk.Context, coin types.Coin) sdk.Error {
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	globalMeta.CumulativeConsumption = globalMeta.CumulativeConsumption.Plus(coin)

	if err := gm.storage.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}
	return nil
}

// add inflation to pool
func (gm GlobalManager) AddToValidatorInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return err
	}
	pool.ValidatorInflationPool = pool.ValidatorInflationPool.Plus(coin)
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return err
	}
	return nil
}

// get validator hourly inflation
func (gm GlobalManager) GetValidatorHourlyInflation(
	ctx sdk.Context, pastHoursThisYear int64) (types.Coin, sdk.Error) {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return types.NewCoin(0), err
	}

	resRat := pool.ValidatorInflationPool.ToRat().Mul(sdk.NewRat(1, types.HoursPerYear-pastHoursThisYear+1))
	resCoin := types.RatToCoin(resRat)
	pool.ValidatorInflationPool = pool.ValidatorInflationPool.Minus(resCoin)
	if err := gm.addTotalLinoCoin(ctx, resCoin); err != nil {
		return types.NewCoin(0), err
	}
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

// get infra monthly inflation
func (gm GlobalManager) GetInfraMonthlyInflation(
	ctx sdk.Context, pastMonthMinusOneThisYear int64) (types.Coin, sdk.Error) {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return types.NewCoin(0), err
	}

	resRat := pool.InfraInflationPool.ToRat().Mul(sdk.NewRat(1, 12-pastMonthMinusOneThisYear))
	resCoin := types.RatToCoin(resRat)
	pool.InfraInflationPool = pool.InfraInflationPool.Minus(resCoin)
	if err := gm.addTotalLinoCoin(ctx, resCoin); err != nil {
		return types.NewCoin(0), err
	}
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

// get developer monthly inflation
func (gm GlobalManager) GetDeveloperMonthlyInflation(
	ctx sdk.Context, pastMonthMinusOneThisYear int64) (types.Coin, sdk.Error) {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return types.NewCoin(0), err
	}

	resRat := pool.DeveloperInflationPool.ToRat().Mul(sdk.NewRat(1, 12-pastMonthMinusOneThisYear))
	resCoin := types.RatToCoin(resRat)
	pool.DeveloperInflationPool = pool.DeveloperInflationPool.Minus(resCoin)
	if err := gm.addTotalLinoCoin(ctx, resCoin); err != nil {
		return types.NewCoin(0), err
	}
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

func (gm GlobalManager) addTotalLinoCoin(ctx sdk.Context, newCoin types.Coin) sdk.Error {
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	globalMeta.TotalLinoCoin = globalMeta.TotalLinoCoin.Plus(newCoin)

	if err := gm.storage.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}
	return nil
}

// update current tps based on current block information
func (gm GlobalManager) UpdateTPS(ctx sdk.Context, lastBlockTime int64) sdk.Error {
	tps, err := gm.storage.GetTPS(ctx)
	if err != nil {
		return err
	}
	if ctx.BlockHeader().Time == lastBlockTime {
		tps.CurrentTPS = sdk.ZeroRat
	} else {
		tps.CurrentTPS = sdk.NewRat(int64(ctx.BlockHeader().NumTxs), ctx.BlockHeader().Time-lastBlockTime)
	}
	if tps.CurrentTPS.GT(tps.MaxTPS) {
		tps.MaxTPS = tps.CurrentTPS
	}

	if err := gm.storage.SetTPS(ctx, tps); err != nil {
		return err
	}
	return nil
}

func (gm GlobalManager) GetTPSCapacityRatio(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	tps, err := gm.storage.GetTPS(ctx)
	if err != nil {
		return sdk.ZeroRat, err
	}
	return tps.CurrentTPS.Quo(tps.MaxTPS), nil
}

func (gm GlobalManager) EvaluateConsumption(
	ctx sdk.Context, coin types.Coin, numOfConsumptionOnAuthor int64, created int64,
	totalReward types.Coin) (types.Coin, sdk.Error) {
	paras, err := gm.paramHolder.GetEvaluateOfContentValueParam(ctx)
	if err != nil {
		return types.NewCoin(0), err
	}
	// evaluate result coin^0.8 * total consumption adjustment *
	// post time adjustment * consumption times adjustment
	expPara, _ := paras.AmountOfConsumptionExponent.GetRat().Float64()
	return types.NewCoin(
		int64(math.Pow(float64(coin.ToInt64()), expPara) *
			PostTotalConsumptionAdjustment(totalReward, paras) *
			PostTimeAdjustment(ctx.BlockHeader().Time-created, paras) *
			PostConsumptionTimesAdjustment(numOfConsumptionOnAuthor, paras))), nil
}

// get and set params
// TODO add more change methods
// func (gm GlobalManager) ChangeGlobalInflationParam(ctx sdk.Context, InfraAllocation sdk.Rat,
// 	ContentCreatorAllocation sdk.Rat, DeveloperAllocation sdk.Rat, ValidatorAllocation sdk.Rat) sdk.Error {
// 	allocation, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	allocation.ContentCreatorAllocation = ContentCreatorAllocation
// 	allocation.DeveloperAllocation = DeveloperAllocation
// 	allocation.InfraAllocation = InfraAllocation
// 	allocation.ValidatorAllocation = ValidatorAllocation
//
// 	if err := gm.paramHolder.SetGlobalAllocationParam(ctx, allocation); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (gm GlobalManager) ChangeInfraInternalInflationParam(
// 	ctx sdk.Context, StorageAllocation sdk.Rat, CDNAllocation sdk.Rat) sdk.Error {
// 	allocation, err := gm.storage.GetInfraInternalAllocationParam(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	allocation.CDNAllocation = CDNAllocation
// 	allocation.StorageAllocation = StorageAllocation
// 	if err := gm.storage.SetInfraInternalAllocationParam(ctx, allocation); err != nil {
// 		return err
// 	}
// 	return nil
// }

// total consumption adjustment = 1/(1+e^(c/base - offset)) + 1
func PostTotalConsumptionAdjustment(
	totalReward types.Coin, paras *param.EvaluateOfContentValueParam) float64 {
	return (1.0 / (1.0 + math.Exp(
		(float64(totalReward.ToInt64())/float64(paras.TotalAmountOfConsumptionBase) -
			float64(paras.TotalAmountOfConsumptionOffset))))) + 1.0
}

// post time adjustment = 1/(1+e^(t/base - offset))
func PostTimeAdjustment(
	elapseTime int64, paras *param.EvaluateOfContentValueParam) float64 {
	return (1.0 / (1.0 + math.Exp(
		(float64(elapseTime)/float64(paras.ConsumptionTimeAdjustBase) -
			float64(paras.ConsumptionTimeAdjustOffset)))))
}

// consumption times adjustment = 1/(1+e^(n-offset)) + 1
func PostConsumptionTimesAdjustment(
	numOfConsumptionOnAuthor int64, paras *param.EvaluateOfContentValueParam) float64 {
	return (1.0/(1.0+math.Exp(
		(float64(numOfConsumptionOnAuthor)-float64(paras.NumOfConsumptionOnAuthorOffset)))) + 1.0) + 1.0
}
