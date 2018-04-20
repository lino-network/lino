package global

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/global/model"
	"github.com/lino-network/lino/types"
)

// GlobalManager encapsulates all basic struct
type GlobalManager struct {
	globalStorage model.GlobalStorage `json:"global_manager"`
}

// NewGlobalManager return the global proxy pointer
func NewGlobalManager(key sdk.StoreKey) GlobalManager {
	return GlobalManager{
		globalStorage: model.NewGlobalStorage(key),
	}
}

func (gm GlobalManager) InitGlobalManager(ctx sdk.Context, totalLino types.Coin) error {
	return gm.globalStorage.InitGlobalState(ctx, totalLino)
}

func (gm GlobalManager) registerEventAtTime(ctx sdk.Context, unixTime int64, event types.Event) sdk.Error {
	if unixTime < ctx.BlockHeader().Time {
		return ErrGlobalManagerRegisterExpiredEvent(unixTime)
	}
	eventList, _ := gm.globalStorage.GetTimeEventList(ctx, unixTime)
	if eventList == nil {
		eventList = &types.TimeEventList{Events: []types.Event{}}
	}
	eventList.Events = append(eventList.Events, event)
	if err := gm.globalStorage.SetTimeEventList(ctx, unixTime, eventList); err != nil {
		return ErrGlobalManagerRegisterEventAtTime(unixTime).TraceCause(err, "")
	}
	return nil
}

func (gm GlobalManager) GetTimeEventListAtTime(ctx sdk.Context, unixTime int64) *types.TimeEventList {
	eventList, _ := gm.globalStorage.GetTimeEventList(ctx, unixTime)
	return eventList
}

func (gm GlobalManager) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	return gm.globalStorage.RemoveTimeEventList(ctx, unixTime)
}

func (gm GlobalManager) GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
	if err != nil {
		return sdk.Rat{}, err
	}
	return consumptionMeta.ConsumptionFrictionRate, nil
}

// register reward calculation event at 7 days later
func (gm GlobalManager) AddFrictionAndRegisterContentRewardEvent(
	ctx sdk.Context, event types.Event, friction types.Coin, evaluate types.Coin) sdk.Error {
	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
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
	if err := gm.globalStorage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
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

func (gm GlobalManager) RegisterProposalDecideEvent(ctx sdk.Context, event types.Event) sdk.Error {
	if err := gm.registerEventAtTime(
		ctx, ctx.BlockHeader().Time+(types.ProposalDecideHr*3600), event); err != nil {
		return err
	}
	return nil
}

// put hourly inflation to reward pool
func (gm GlobalManager) AddHourlyInflationToRewardPool(ctx sdk.Context, pastHoursThisYear int64) sdk.Error {
	pool, getErr := gm.globalStorage.GetInflationPool(ctx)
	if getErr != nil {
		return getErr
	}
	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	resRat := pool.ContentCreatorInflationPool.ToRat().
		Mul(sdk.NewRat(1, types.HoursPerYear-pastHoursThisYear+1))
	resCoin := types.RatToCoin(resRat)
	pool.ContentCreatorInflationPool = pool.ContentCreatorInflationPool.Minus(resCoin)

	if err := gm.globalStorage.SetInflationPool(ctx, pool); err != nil {
		return err
	}

	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Plus(resCoin)

	if err := gm.globalStorage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}
	return nil
}

// after 7 days, one consumption needs to claim its reward from consumption reward pool
func (gm GlobalManager) GetRewardAndPopFromWindow(
	ctx sdk.Context, coin types.Coin, penaltyScore sdk.Rat) (types.Coin, sdk.Error) {
	if coin.IsZero() {
		return types.NewCoin(0), nil
	}

	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
	if err != nil {
		return types.NewCoin(0), ErrGetRewardAndPopFromWindow().TraceCause(err, "")
	}

	// reward = (consumption reward pool) * ((this consumption * penalty score) / (total consumption in 7 days window))
	reward := types.RatToCoin(consumptionMeta.ConsumptionRewardPool.ToRat().
		Mul(coin.ToRat().Mul(sdk.OneRat.Sub(penaltyScore)).
			Quo(consumptionMeta.ConsumptionWindow.ToRat())))

	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Minus(reward)
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Minus(coin)

	if err := gm.globalStorage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return types.NewCoin(0), ErrGetRewardAndPopFromWindow().TraceCause(err, "")
	}
	return reward, nil
}

// add consumption to global meta, which is used to compute GDP
func (gm GlobalManager) AddConsumption(ctx sdk.Context, coin types.Coin) sdk.Error {
	globalMeta, err := gm.globalStorage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	globalMeta.CumulativeConsumption = globalMeta.CumulativeConsumption.Plus(coin)

	if err := gm.globalStorage.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}
	return nil
}

func (gm GlobalManager) AddToValidatorInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error {
	pool, getErr := gm.globalStorage.GetInflationPool(ctx)
	if getErr != nil {
		return getErr
	}
	pool.ValidatorInflationPool = pool.ValidatorInflationPool.Plus(coin)
	if err := gm.globalStorage.SetInflationPool(ctx, pool); err != nil {
		return err
	}
	return nil
}

func (gm GlobalManager) GetValidatorHourlyInflation(
	ctx sdk.Context, pastHoursThisYear int64) (types.Coin, sdk.Error) {
	pool, getErr := gm.globalStorage.GetInflationPool(ctx)
	if getErr != nil {
		return types.NewCoin(0), getErr
	}

	resRat := pool.ValidatorInflationPool.ToRat().Mul(sdk.NewRat(1, types.HoursPerYear-pastHoursThisYear+1))
	resCoin := types.RatToCoin(resRat)
	pool.ValidatorInflationPool = pool.ValidatorInflationPool.Minus(resCoin)

	if err := gm.globalStorage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

func (gm GlobalManager) GetInfraMonthlyInflation(
	ctx sdk.Context, pastMonthMinusOneThisYear int64) (types.Coin, sdk.Error) {
	pool, getErr := gm.globalStorage.GetInflationPool(ctx)
	if getErr != nil {
		return types.NewCoin(0), getErr
	}

	resRat := pool.InfraInflationPool.ToRat().Mul(sdk.NewRat(1, 12-pastMonthMinusOneThisYear))
	resCoin := types.RatToCoin(resRat)
	pool.InfraInflationPool = pool.InfraInflationPool.Minus(resCoin)

	if err := gm.globalStorage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

func (gm GlobalManager) GetDeveloperMonthlyInflation(
	ctx sdk.Context, pastMonthMinusOneThisYear int64) (types.Coin, sdk.Error) {
	pool, getErr := gm.globalStorage.GetInflationPool(ctx)
	if getErr != nil {
		return types.NewCoin(0), getErr
	}

	resRat := pool.DeveloperInflationPool.ToRat().Mul(sdk.NewRat(1, 12-pastMonthMinusOneThisYear))
	resCoin := types.RatToCoin(resRat)
	pool.DeveloperInflationPool = pool.DeveloperInflationPool.Minus(resCoin)

	if err := gm.globalStorage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

func (gm GlobalManager) ChangeInfraInternalInflation(
	ctx sdk.Context, StorageAllocation sdk.Rat, CDNAllocation sdk.Rat) sdk.Error {
	allocation, getErr := gm.globalStorage.GetInfraInternalAllocation(ctx)
	if getErr != nil {
		return getErr
	}
	allocation.CDNAllocation = CDNAllocation
	allocation.StorageAllocation = StorageAllocation
	if err := gm.globalStorage.SetInfraInternalAllocation(ctx, allocation); err != nil {
		return err
	}
	return nil
}

// update current tps based on current block information
func (gm GlobalManager) UpdateTPS(ctx sdk.Context, lastBlockTime int64) sdk.Error {
	tps, err := gm.globalStorage.GetTPS(ctx)
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

	if err := gm.globalStorage.SetTPS(ctx, tps); err != nil {
		return err
	}
	return nil
}

func (gm GlobalManager) ChangeGlobalInflation(ctx sdk.Context, InfraAllocation sdk.Rat,
	ContentCreatorAllocation sdk.Rat, DeveloperAllocation sdk.Rat, ValidatorAllocation sdk.Rat) sdk.Error {
	allocation, getErr := gm.globalStorage.GetGlobalAllocation(ctx)
	if getErr != nil {
		return getErr
	}
	allocation.ContentCreatorAllocation = ContentCreatorAllocation
	allocation.DeveloperAllocation = DeveloperAllocation
	allocation.InfraAllocation = InfraAllocation
	allocation.ValidatorAllocation = ValidatorAllocation

	if err := gm.globalStorage.SetGlobalAllocation(ctx, allocation); err != nil {
		return err
	}
	return nil
}

func (gm GlobalManager) GetTPSCapacityRatio(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	tps, err := gm.globalStorage.GetTPS(ctx)
	if err != nil {
		return sdk.ZeroRat, err
	}
	return tps.CurrentTPS.Quo(tps.MaxTPS), nil
}

func (gm GlobalManager) EvaluateConsumption(
	ctx sdk.Context, coin types.Coin, numOfConsumptionOnAuthor int64, created int64,
	totalReward types.Coin) (types.Coin, sdk.Error) {
	paras, err := gm.globalStorage.GetEvaluateOfContentValuePara(ctx)
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

// total consumption adjustment = 1/(1+e^(c/base - offset)) + 1
func PostTotalConsumptionAdjustment(
	totalReward types.Coin, paras *model.EvaluateOfContentValuePara) float64 {
	return (1.0 / (1.0 + math.Exp(
		(float64(totalReward.ToInt64())/float64(paras.TotalAmountOfConsumptionBase) -
			float64(paras.TotalAmountOfConsumptionOffset))))) + 1.0
}

// post time adjustment = 1/(1+e^(t/base - offset))
func PostTimeAdjustment(
	elapseTime int64, paras *model.EvaluateOfContentValuePara) float64 {
	return (1.0 / (1.0 + math.Exp(
		(float64(elapseTime)/float64(paras.ConsumptionTimeAdjustBase) -
			float64(paras.ConsumptionTimeAdjustOffset)))))
}

// consumption times adjustment = 1/(1+e^(n-offset))
func PostConsumptionTimesAdjustment(
	numOfConsumptionOnAuthor int64, paras *model.EvaluateOfContentValuePara) float64 {
	return (1.0/(1.0+math.Exp(
		(float64(numOfConsumptionOnAuthor)-float64(paras.NumOfConsumptionOnAuthorOffset)))) + 1.0)
}
