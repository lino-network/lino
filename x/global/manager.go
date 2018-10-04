package global

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/recorder"
	"github.com/lino-network/lino/recorder/inflation"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global/model"
)

// GlobalManager - encapsulates all basic struct
type GlobalManager struct {
	storage     model.GlobalStorage
	paramHolder param.ParamHolder
	recorder    recorder.Recorder
}

// NewGlobalManager - return the global manager
func NewGlobalManager(key sdk.StoreKey, holder param.ParamHolder, recorder recorder.Recorder) GlobalManager {
	return GlobalManager{
		storage:     model.NewGlobalStorage(key),
		paramHolder: holder,
		recorder:    recorder,
	}
}

// WireCodec - access to global manager codec
func (gm GlobalManager) WireCodec() *wire.Codec {
	return gm.storage.WireCodec()
}

// InitGlobalManager - initialize global manager based on code
func (gm GlobalManager) InitGlobalManager(ctx sdk.Context, totalLino types.Coin) sdk.Error {
	return gm.storage.InitGlobalState(ctx, totalLino)
}

// InitGlobalManagerWithConfig - initialize global manager based on genesis file
func (gm GlobalManager) InitGlobalManagerWithConfig(
	ctx sdk.Context, totalLino types.Coin, param model.InitParamList) sdk.Error {
	return gm.storage.InitGlobalStateWithConfig(ctx, totalLino, param)
}

func (gm GlobalManager) registerEventAtTime(ctx sdk.Context, unixTime int64, event types.Event) sdk.Error {
	if unixTime < ctx.BlockHeader().Time.Unix() {
		return ErrRegisterExpiredEvent(unixTime)
	}
	eventList, err := gm.storage.GetTimeEventList(ctx, unixTime)
	if err != nil {
		return err
	}
	if eventList == nil {
		eventList = &types.TimeEventList{Events: []types.Event{}}
	}
	eventList.Events = append(eventList.Events, event)
	if err := gm.storage.SetTimeEventList(ctx, unixTime, eventList); err != nil {
		return err
	}
	return nil
}

// GetTimeEventListAtTime - get time event list at given time
func (gm GlobalManager) GetTimeEventListAtTime(ctx sdk.Context, unixTime int64) *types.TimeEventList {
	eventList, _ := gm.storage.GetTimeEventList(ctx, unixTime)
	return eventList
}

// GetLastBlockTime - get last block time from KVStore
func (gm GlobalManager) GetLastBlockTime(ctx sdk.Context) (int64, sdk.Error) {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return 0, err
	}
	return globalTime.LastBlockTime, nil
}

// SetLastBlockTime - set last block time to KVStore
func (gm GlobalManager) SetLastBlockTime(ctx sdk.Context, unixTime int64) sdk.Error {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return err
	}
	globalTime.LastBlockTime = unixTime
	return gm.storage.SetGlobalTime(ctx, globalTime)
}

// GetPastDay - get start time from KVStore to calculate past day
func (gm GlobalManager) GetPastDay(ctx sdk.Context, unixTime int64) (int64, sdk.Error) {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return 0, err
	}
	pastDay := (unixTime - globalTime.ChainStartTime) / (3600 * 24)
	if pastDay < 0 {
		return 0, nil
	}
	return pastDay, nil
}

// GetChainStartTime - get chain start time from KVStore
func (gm GlobalManager) GetChainStartTime(ctx sdk.Context) (int64, sdk.Error) {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return 0, err
	}
	return globalTime.ChainStartTime, nil
}

// SetChainStartTime - set chain start time to KVStore
func (gm GlobalManager) SetChainStartTime(ctx sdk.Context, unixTime int64) sdk.Error {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return err
	}
	globalTime.ChainStartTime = unixTime
	return gm.storage.SetGlobalTime(ctx, globalTime)
}

// GetPastMinutes - get past minutes from KVStore
func (gm GlobalManager) GetPastMinutes(ctx sdk.Context) (int64, sdk.Error) {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return 0, err
	}
	return globalTime.PastMinutes, nil
}

// SetPastMinutes - set past minutes to KVStore
func (gm GlobalManager) SetPastMinutes(ctx sdk.Context, minutes int64) sdk.Error {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return err
	}
	globalTime.PastMinutes = minutes
	return gm.storage.SetGlobalTime(ctx, globalTime)
}

// RemoveTimeEventList - remove time event list from KVstore at given time
func (gm GlobalManager) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	return gm.storage.RemoveTimeEventList(ctx, unixTime)
}

// GetConsumptionFrictionRate - get consumption friction rate
func (gm GlobalManager) GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return sdk.Dec{}, err
	}
	return consumptionMeta.ConsumptionFrictionRate, nil
}

// GetConsumption - get this year consumption
func (gm GlobalManager) GetConsumption(ctx sdk.Context) (types.Coin, sdk.Error) {
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return globalMeta.CumulativeConsumption, nil
}

// AddFrictionAndRegisterContentRewardEvent - register reward calculation event at 7 days later
func (gm GlobalManager) AddFrictionAndRegisterContentRewardEvent(
	ctx sdk.Context, event types.Event, friction types.Coin, evaluate types.Coin) sdk.Error {
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	pastDay, err := gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
	if err != nil {
		return err
	}
	linoStakeStat, err := gm.storage.GetLinoStakeStat(ctx, pastDay)
	if err != nil {
		return err
	}
	// consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Plus(friction)
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Plus(evaluate)
	linoStakeStat.TotalConsumptionFriction = linoStakeStat.TotalConsumptionFriction.Plus(friction)
	linoStakeStat.UnclaimedFriction = linoStakeStat.UnclaimedFriction.Plus(friction)
	if err := gm.registerEventAtTime(
		ctx, ctx.BlockHeader().Time.Unix()+consumptionMeta.ConsumptionFreezingPeriodSec, event); err != nil {
		return err
	}
	if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}
	if err := gm.storage.SetLinoStakeStat(ctx, pastDay, linoStakeStat); err != nil {
		return err
	}
	return nil
}

// AddLinoStakeToStat - add lino power to total lino power at current day
func (gm GlobalManager) AddLinoStakeToStat(ctx sdk.Context, linoStake types.Coin) sdk.Error {
	pastDay, err := gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
	if err != nil {
		return err
	}
	linoStakeStat, err := gm.storage.GetLinoStakeStat(ctx, pastDay)
	if err != nil {
		return err
	}
	linoStakeStat.TotalLinoStake = linoStakeStat.TotalLinoStake.Plus(linoStake)
	linoStakeStat.UnclaimedLinoStake = linoStakeStat.UnclaimedLinoStake.Plus(linoStake)
	if err := gm.storage.SetLinoStakeStat(ctx, pastDay, linoStakeStat); err != nil {
		return err
	}
	return nil
}

// MinusLinoStakeFromStat - minus lino power from total lino power at current day
func (gm GlobalManager) MinusLinoStakeFromStat(ctx sdk.Context, linoStake types.Coin) sdk.Error {
	pastDay, err := gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
	if err != nil {
		return err
	}
	linoStakeStat, err := gm.storage.GetLinoStakeStat(ctx, pastDay)
	if err != nil {
		return err
	}
	linoStakeStat.TotalLinoStake = linoStakeStat.TotalLinoStake.Minus(linoStake)
	linoStakeStat.UnclaimedLinoStake = linoStakeStat.UnclaimedLinoStake.Minus(linoStake)
	if err := gm.storage.SetLinoStakeStat(ctx, pastDay, linoStakeStat); err != nil {
		return err
	}
	return nil
}

// GetInterestSince - get interest from unix time till now (exclusive)
func (gm GlobalManager) GetInterestSince(ctx sdk.Context, unixTime int64, linoStake types.Coin) (types.Coin, sdk.Error) {
	startDay, err := gm.GetPastDay(ctx, unixTime)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	endDay, err := gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	totalInterest := types.NewCoinFromInt64(0)
	for day := startDay; day < endDay; day++ {
		linoStakeStat, err := gm.storage.GetLinoStakeStat(ctx, day)
		if err != nil {
			return types.NewCoinFromInt64(0), err
		}
		if linoStakeStat.UnclaimedLinoStake.IsZero() {
			continue
		}
		interest :=
			types.DecToCoin(linoStakeStat.UnclaimedFriction.ToDec().Mul(
				linoStake.ToDec().Quo(linoStakeStat.UnclaimedLinoStake.ToDec())))
		totalInterest = totalInterest.Plus(interest)
		linoStakeStat.UnclaimedFriction = linoStakeStat.UnclaimedFriction.Minus(interest)
		linoStakeStat.UnclaimedLinoStake = linoStakeStat.UnclaimedLinoStake.Minus(linoStake)
		if err := gm.storage.SetLinoStakeStat(ctx, day, linoStakeStat); err != nil {
			return types.NewCoinFromInt64(0), err
		}
	}
	return totalInterest, nil
}

// RecordConsumptionAndLinoStake - records consumption and lino power to LinoStakeStat and renew to new slot
func (gm GlobalManager) RecordConsumptionAndLinoStake(ctx sdk.Context) sdk.Error {
	pastMinutes, err := gm.GetPastMinutes(ctx)
	if err != nil {
		return err
	}
	lastLinoStakeStat, err := gm.storage.GetLinoStakeStat(ctx, (pastMinutes/types.MinutesPerDay)-1)
	if err != nil {
		return err
	}
	// If lino stake exist last day, the consumption will keep for lino stake holder that day
	if !lastLinoStakeStat.TotalLinoStake.IsZero() {
		lastLinoStakeStat.TotalConsumptionFriction = types.NewCoinFromInt64(0)
		lastLinoStakeStat.UnclaimedFriction = types.NewCoinFromInt64(0)
	}
	if err := gm.storage.SetLinoStakeStat(ctx, pastMinutes/types.MinutesPerDay, lastLinoStakeStat); err != nil {
		return err
	}
	return nil
}

// RegisterCoinReturnEvent - register coin return event with time interval
func (gm GlobalManager) RegisterCoinReturnEvent(
	ctx sdk.Context, events []types.Event, times int64, intervalSec int64) sdk.Error {
	for i := int64(0); i < times; i++ {
		if err := gm.registerEventAtTime(
			ctx, ctx.BlockHeader().Time.Unix()+(intervalSec*(i+1)), events[i]); err != nil {
			return err
		}
	}
	return nil
}

// RegisterProposalDecideEvent - register proposal decide event
func (gm GlobalManager) RegisterProposalDecideEvent(
	ctx sdk.Context, decideSec int64, event types.Event) sdk.Error {
	if err := gm.registerEventAtTime(
		ctx, ctx.BlockHeader().Time.Unix()+decideSec, event); err != nil {
		return err
	}
	return nil
}

// RegisterParamChangeEvent - register parameter change event
func (gm GlobalManager) RegisterParamChangeEvent(ctx sdk.Context, event types.Event) sdk.Error {
	// param will be changed in one day
	proposalParam, err := gm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err
	}
	if ctx.BlockHeader().Height > types.LinoBlockchainFirstUpdateHeight {
		if err := gm.registerEventAtTime(
			ctx, ctx.BlockHeader().Time.Unix()+3600, event); err != nil {
			return err
		}
	} else {
		if err := gm.registerEventAtTime(
			ctx, ctx.BlockHeader().Time.Unix()+proposalParam.ChangeParamExecutionSec, event); err != nil {
			return err
		}
	}
	return nil
}

// DistributeHourlyInflation - distribute inflation hourly
func (gm GlobalManager) DistributeHourlyInflation(ctx sdk.Context) sdk.Error {
	// param will be changed in one day
	globalAllocation, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	if err != nil {
		return err
	}
	if ctx.BlockHeader().Height > types.LinoBlockchainFirstUpdateHeight {
		globalAllocation.DeveloperAllocation = types.NewDecFromRat(35, 100)
		globalAllocation.ContentCreatorAllocation = types.NewDecFromRat(40, 1000)
	}
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return err
	}
	// get hourly inflation
	thisHourInflation :=
		types.DecToCoin(
			globalMeta.LastYearTotalLinoCoin.ToDec().
				Mul(globalAllocation.GlobalGrowthRate).
				Mul(types.NewDecFromRat(1, types.HoursPerYear)))
	if err := gm.storage.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}

	// distribute content creator inflation to consumption meta
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	contentCreatorInflation :=
		types.DecToCoin(thisHourInflation.ToDec().Mul(globalAllocation.ContentCreatorAllocation))
	validatorInflation :=
		types.DecToCoin(thisHourInflation.ToDec().Mul(globalAllocation.ValidatorAllocation))
	infraInflation :=
		types.DecToCoin(thisHourInflation.ToDec().Mul(globalAllocation.InfraAllocation))
	developerInflation :=
		thisHourInflation.Minus(contentCreatorInflation).Minus(validatorInflation).Minus(infraInflation)
	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Plus(contentCreatorInflation)
	if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}

	// distribute inflation to validator inflation pool
	pool.InfraInflationPool = pool.InfraInflationPool.Plus(infraInflation)
	pool.ValidatorInflationPool = pool.ValidatorInflationPool.Plus(validatorInflation)
	pool.DeveloperInflationPool = pool.DeveloperInflationPool.Plus(developerInflation)
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return err
	}

	// record
	infraInflationPoolInt, _ := pool.InfraInflationPool.ToInt64()
	devInflationPoolInt, _ := pool.DeveloperInflationPool.ToInt64()
	creatorInflationPoolInt, _ := consumptionMeta.ConsumptionRewardPool.ToInt64()
	validatorInflationPoolInt, _ := pool.ValidatorInflationPool.ToInt64()

	infraInflationInt, _ := infraInflation.ToInt64()
	devInflationInt, _ := developerInflation.ToInt64()
	creatorInflationInt, _ := contentCreatorInflation.ToInt64()
	validatorInflationInt, _ := validatorInflation.ToInt64()

	inflation := &inflation.Inflation{
		InfraPool:          infraInflationPoolInt,
		DevPool:            devInflationPoolInt,
		CreatorPool:        creatorInflationPoolInt,
		ValidatorPool:      validatorInflationPoolInt,
		InfraInflation:     infraInflationInt,
		DevInflation:       devInflationInt,
		CreatorInflation:   creatorInflationInt,
		ValidatorInflation: validatorInflationInt,
		Timestamp:          ctx.BlockHeader().Time.Unix(),
	}
	gm.recorder.InflationRepo.Add(inflation)

	return nil
}

// SetTotalLinoAndRecalculateGrowthRate - recalculate annually inflation based on consumption growth rate
func (gm GlobalManager) SetTotalLinoAndRecalculateGrowthRate(ctx sdk.Context) sdk.Error {
	var growthRate sdk.Dec
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	globalAllocationParam, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	if err != nil {
		return err
	}
	// if last year cumulative consumption is zero, we use the same growth rate as the last year
	if globalMeta.LastYearCumulativeConsumption.IsZero() {
		growthRate = globalAllocationParam.GlobalGrowthRate
	} else {
		// growthRate = (consumption this year - consumption last year) / consumption last year
		lastYearConsumptionRat := globalMeta.LastYearCumulativeConsumption.ToDec()
		thisYearConsumptionRat := globalMeta.CumulativeConsumption.ToDec()
		consumptionIncrement := thisYearConsumptionRat.Sub(lastYearConsumptionRat)
		if consumptionIncrement.GTE(lastYearConsumptionRat) {
			growthRate = sdk.OneDec()
		} else {
			growthRate = consumptionIncrement.Quo(lastYearConsumptionRat)
		}
	}
	globalMeta.LastYearCumulativeConsumption = globalMeta.CumulativeConsumption
	globalMeta.CumulativeConsumption = types.NewCoinFromInt64(0)
	globalMeta.LastYearTotalLinoCoin = globalMeta.TotalLinoCoin
	if err := gm.storage.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}
	return gm.paramHolder.UpdateGlobalGrowthRate(ctx, growthRate)
}

// GetRewardAndPopFromWindow - after 7 days, one consumption needs to claim its reward from consumption reward pool
func (gm GlobalManager) GetRewardAndPopFromWindow(
	ctx sdk.Context, evaluate types.Coin, penaltyScore sdk.Dec) (types.Coin, sdk.Error) {
	if evaluate.IsZero() {
		return types.NewCoinFromInt64(0), nil
	}

	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	// consumptionRatio = (this consumption * penalty score) / (total consumption in 7 days window)
	consumptionRatio :=
		evaluate.ToDec().Mul(sdk.OneDec().Sub(penaltyScore)).Quo(
			consumptionMeta.ConsumptionWindow.ToDec())
	// reward = (consumption reward pool) * (consumptionRatio)
	reward := types.DecToCoin(
		consumptionMeta.ConsumptionRewardPool.ToDec().Mul(consumptionRatio))
	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Minus(reward)
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Minus(evaluate)
	if err := gm.addTotalLinoCoin(ctx, reward); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return reward, nil
}

// AddConsumption - add consumption to global meta, which is used to compute GDP
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

// AddToDeveloperInflationPool - add coin to developer inflation pool
func (gm GlobalManager) AddToDeveloperInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error {
	inflationPool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return err
	}
	inflationPool.DeveloperInflationPool = inflationPool.DeveloperInflationPool.Plus(coin)

	if err := gm.storage.SetInflationPool(ctx, inflationPool); err != nil {
		return err
	}
	return nil
}

// AddToValidatorInflationPool - add validator inflation to pool
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

// GetValidatorHourlyInflation - get validator hourly inflation
func (gm GlobalManager) GetValidatorHourlyInflation(ctx sdk.Context) (types.Coin, sdk.Error) {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	resCoin := pool.ValidatorInflationPool
	pool.ValidatorInflationPool = types.NewCoinFromInt64(0)
	if err := gm.addTotalLinoCoin(ctx, resCoin); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return resCoin, nil
}

// GetInfraMonthlyInflation - get infra monthly inflation
func (gm GlobalManager) GetInfraMonthlyInflation(ctx sdk.Context) (types.Coin, sdk.Error) {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	resCoin := pool.InfraInflationPool
	pool.InfraInflationPool = types.NewCoinFromInt64(0)
	if err := gm.addTotalLinoCoin(ctx, resCoin); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return resCoin, nil
}

// GetDeveloperMonthlyInflation - get developer monthly inflation
func (gm GlobalManager) GetDeveloperMonthlyInflation(ctx sdk.Context) (types.Coin, sdk.Error) {
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	resCoin := pool.DeveloperInflationPool
	pool.DeveloperInflationPool = types.NewCoinFromInt64(0)
	if err := gm.addTotalLinoCoin(ctx, resCoin); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := gm.storage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoinFromInt64(0), err
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

// UpdateTPS - update current tps based on current block information
func (gm GlobalManager) UpdateTPS(ctx sdk.Context) sdk.Error {
	tps, err := gm.storage.GetTPS(ctx)
	if err != nil {
		return err
	}
	lastBlockTime, err := gm.GetLastBlockTime(ctx)
	if err != nil {
		return err
	}

	if ctx.BlockHeader().Time.Unix() == lastBlockTime {
		tps.CurrentTPS = sdk.ZeroDec()
	} else {
		tps.CurrentTPS = types.NewDecFromRat(int64(ctx.BlockHeader().NumTxs), ctx.BlockHeader().Time.Unix()-lastBlockTime)
	}
	if tps.CurrentTPS.GT(tps.MaxTPS) {
		tps.MaxTPS = tps.CurrentTPS
	}

	if err := gm.storage.SetTPS(ctx, tps); err != nil {
		return err
	}
	return nil
}

// GetTPSCapacityRatio - get transaction per second ratio
func (gm GlobalManager) GetTPSCapacityRatio(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	tps, err := gm.storage.GetTPS(ctx)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	return tps.CurrentTPS.Quo(tps.MaxTPS), nil
}

// EvaluateConsumption - evaluate consumption
// XXX(yumin): deprecated, now just return @p coin, need @refactor.
func (gm GlobalManager) EvaluateConsumption(
	coin types.Coin,
	// ctx sdk.Context, coin types.Coin,
	// numOfConsumptionOnAuthor int64, created int64,
	// totalReward types.Coin
) (types.Coin, sdk.Error) {
	return coin, nil
}

// get and set params
// TODO add more change methods
// func (gm GlobalManager) ChangeGlobalInflationParam(ctx sdk.Context, InfraAllocation sdk.Dec,
// 	ContentCreatorAllocation sdk.Dec, DeveloperAllocation sdk.Dec, ValidatorAllocation sdk.Dec) sdk.Error {
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
// 	ctx sdk.Context, StorageAllocation sdk.Dec, CDNAllocation sdk.Dec) sdk.Error {
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
