package global

import (
	"fmt"

	codec "github.com/cosmos/cosmos-sdk/codec"
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/global/model"

	// only for import from last version.
	posttypes "github.com/lino-network/lino/x/post/types"
)

// GlobalManager - encapsulates all basic struct
type GlobalManager struct {
	storage                 model.GlobalStorage
	paramHolder             param.ParamHolder
	deliverTxEventCacheList []*model.EventCache
}

// NewGlobalManager - return the global manager
func NewGlobalManager(key sdk.StoreKey, holder param.ParamHolder) GlobalManager {
	return GlobalManager{
		storage:     model.NewGlobalStorage(key),
		paramHolder: holder,
	}
}

// WireCodec - access to global manager codec
func (gm *GlobalManager) WireCodec() *wire.Codec {
	return gm.storage.WireCodec()
}

// InitGlobalManager - initialize global manager based on code
func (gm *GlobalManager) InitGlobalManager(ctx sdk.Context, totalLino types.Coin) sdk.Error {
	return gm.storage.InitGlobalState(ctx, totalLino)
}

// InitGlobalManagerWithConfig - initialize global manager based on genesis file
func (gm *GlobalManager) InitGlobalManagerWithConfig(
	ctx sdk.Context, totalLino types.Coin, param model.InitParamList) sdk.Error {
	return gm.storage.InitGlobalStateWithConfig(ctx, totalLino, param)
}

func (gm *GlobalManager) RegisterEventAtTime(ctx sdk.Context, unixTime int64, event types.Event) sdk.Error {
	if unixTime < ctx.BlockHeader().Time.Unix() {
		return ErrRegisterExpiredEvent(unixTime)
	}
	if ctx.IsCheckTx() {
		return nil
	}
	for _, eventCache := range gm.deliverTxEventCacheList {
		if unixTime == eventCache.UnixTime {
			eventCache.EventList = append(eventCache.EventList, event)
			return nil
		}
	}
	gm.deliverTxEventCacheList = append(gm.deliverTxEventCacheList, &model.EventCache{
		UnixTime:  unixTime,
		EventList: []types.Event{event},
	})
	return nil
}

// GetTimeEventListAtTime - get time event list at given time
func (gm *GlobalManager) GetTimeEventListAtTime(ctx sdk.Context, unixTime int64) *types.TimeEventList {
	eventList, _ := gm.storage.GetTimeEventList(ctx, unixTime)
	return eventList
}

// GetLastBlockTime - get last block time from KVStore
func (gm *GlobalManager) GetLastBlockTime(ctx sdk.Context) (int64, sdk.Error) {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return 0, err
	}
	return globalTime.LastBlockTime, nil
}

// SetLastBlockTime - set last block time to KVStore
func (gm *GlobalManager) SetLastBlockTime(ctx sdk.Context, unixTime int64) sdk.Error {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return err
	}
	globalTime.LastBlockTime = unixTime
	return gm.storage.SetGlobalTime(ctx, globalTime)
}

// GetPastDay - get start time from KVStore to calculate past day
func (gm *GlobalManager) GetPastDay(ctx sdk.Context, unixTime int64) (int64, sdk.Error) {
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
func (gm *GlobalManager) GetChainStartTime(ctx sdk.Context) (int64, sdk.Error) {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return 0, err
	}
	return globalTime.ChainStartTime, nil
}

// SetChainStartTime - set chain start time to KVStore
func (gm *GlobalManager) SetChainStartTime(ctx sdk.Context, unixTime int64) sdk.Error {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return err
	}
	globalTime.ChainStartTime = unixTime
	return gm.storage.SetGlobalTime(ctx, globalTime)
}

// GetPastMinutes - get past minutes from KVStore
func (gm *GlobalManager) GetPastMinutes(ctx sdk.Context) (int64, sdk.Error) {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return 0, err
	}
	return globalTime.PastMinutes, nil
}

// SetPastMinutes - set past minutes to KVStore
func (gm *GlobalManager) SetPastMinutes(ctx sdk.Context, minutes int64) sdk.Error {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		return err
	}
	globalTime.PastMinutes = minutes
	return gm.storage.SetGlobalTime(ctx, globalTime)
}

// RemoveTimeEventList - remove time event list from KVstore at given time
func (gm *GlobalManager) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	return gm.storage.RemoveTimeEventList(ctx, unixTime)
}

// GetConsumptionFrictionRate - get consumption friction rate
func (gm *GlobalManager) GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return sdk.Dec{}, err
	}
	return consumptionMeta.ConsumptionFrictionRate, nil
}

// AddFrictionAndRegisterContentRewardEvent - register reward calculation event at 7 days later
func (gm *GlobalManager) AddFrictionAndRegisterContentRewardEvent(
	ctx sdk.Context, event types.Event, friction types.Coin, evaluate types.MiniDollar) sdk.Error {
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
	consumptionMeta.ConsumptionWindow = types.NewMiniDollarFromInt(consumptionMeta.ConsumptionWindow.Add(evaluate.Int))
	linoStakeStat.TotalConsumptionFriction = linoStakeStat.TotalConsumptionFriction.Plus(friction)
	linoStakeStat.UnclaimedFriction = linoStakeStat.UnclaimedFriction.Plus(friction)
	if err := gm.RegisterEventAtTime(
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
func (gm *GlobalManager) AddLinoStakeToStat(ctx sdk.Context, linoStake types.Coin) sdk.Error {
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
func (gm *GlobalManager) MinusLinoStakeFromStat(ctx sdk.Context, linoStake types.Coin) sdk.Error {
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
func (gm *GlobalManager) GetInterestSince(ctx sdk.Context, unixTime int64, linoStake types.Coin) (types.Coin, sdk.Error) {
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
		if linoStakeStat.UnclaimedLinoStake.IsZero() && !linoStakeStat.UnclaimedLinoStake.IsGTE(linoStake) {
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
func (gm *GlobalManager) RecordConsumptionAndLinoStake(ctx sdk.Context) sdk.Error {
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
func (gm *GlobalManager) RegisterCoinReturnEvent(
	ctx sdk.Context, events []types.Event, times int64, intervalSec int64) sdk.Error {
	for i := int64(0); i < times; i++ {
		if err := gm.RegisterEventAtTime(
			ctx, ctx.BlockHeader().Time.Unix()+(intervalSec*(i+1)), events[i]); err != nil {
			return err
		}
	}
	return nil
}

// DistributeHourlyInflation - distribute inflation hourly
func (gm *GlobalManager) DistributeHourlyInflation(ctx sdk.Context) sdk.Error {
	// param will be changed in one day
	globalAllocation, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	if err != nil {
		return err
	}
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}

	if ctx.BlockHeader().Height > types.Upgrade2Update3Height {
		globalMeta.LastYearTotalLinoCoin = types.NewCoinFromInt64(10000000000 * types.Decimals)
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return err
	}

	// get hourly inflation
	growthRate := globalAllocation.GlobalGrowthRate
	thisHourInflation :=
		types.DecToCoin(
			globalMeta.LastYearTotalLinoCoin.ToDec().
				Mul(growthRate).
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
	return nil
}

// SetTotalLinoAndRecalculateGrowthRate - recalculate annually inflation based on consumption growth rate
func (gm *GlobalManager) SetTotalLinoAndRecalculateGrowthRate(ctx sdk.Context) sdk.Error {
	var growthRate sdk.Dec
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	globalAllocationParam, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	if err != nil {
		return err
	}
	growthRate = globalAllocationParam.GlobalGrowthRate
	globalMeta.LastYearTotalLinoCoin = globalMeta.TotalLinoCoin
	if err := gm.storage.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}
	return gm.paramHolder.UpdateGlobalGrowthRate(ctx, growthRate)
}

// GetRewardAndPopFromWindow - after 7 days, one consumption needs to claim its reward from consumption reward pool
func (gm *GlobalManager) GetRewardAndPopFromWindow(ctx sdk.Context, evaluate types.MiniDollar) (types.Coin, sdk.Error) {
	if evaluate.IsZero() {
		return types.NewCoinFromInt64(0), nil
	}

	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	// consumptionRatio = (this consumption * penalty score) / (total consumption in 7 days window)
	consumptionRatio := sdk.ZeroDec()
	if !consumptionMeta.ConsumptionWindow.ToDec().IsZero() {
		consumptionRatio =
			evaluate.ToDec().Quo(consumptionMeta.ConsumptionWindow.ToDec())
	}
	// reward = (consumption reward pool) * (consumptionRatio)
	reward := types.DecToCoin(
		consumptionMeta.ConsumptionRewardPool.ToDec().Mul(consumptionRatio))
	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Minus(reward)
	consumptionMeta.ConsumptionWindow = types.NewMiniDollarFromInt(consumptionMeta.ConsumptionWindow.Sub(evaluate.Int))
	if err := gm.addTotalLinoCoin(ctx, reward); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return reward, nil
}

// AddToDeveloperInflationPool - add coin to developer inflation pool
func (gm *GlobalManager) AddToDeveloperInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error {
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
func (gm *GlobalManager) AddToValidatorInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error {
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
func (gm *GlobalManager) GetValidatorHourlyInflation(ctx sdk.Context) (types.Coin, sdk.Error) {
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
func (gm *GlobalManager) GetInfraMonthlyInflation(ctx sdk.Context) (types.Coin, sdk.Error) {
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

// PopDeveloperMonthlyInflation - pop out developer monthly inflation and
// reset current inflation pool to zero. Total coin amount will be updated.
func (gm *GlobalManager) PopDeveloperMonthlyInflation(ctx sdk.Context) (types.Coin, sdk.Error) {
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

func (gm *GlobalManager) addTotalLinoCoin(ctx sdk.Context, newCoin types.Coin) sdk.Error {
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

// GetTPSCapacityRatio - get transaction per second ratio
func (gm *GlobalManager) GetTPSCapacityRatio(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	tps, err := gm.storage.GetTPS(ctx)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	return tps.CurrentTPS.Quo(tps.MaxTPS), nil
}

// CommitEvent - append event to event list
// Commit event cache will only be committed at the endblocker
func (gm *GlobalManager) CommitEventCache(ctx sdk.Context) sdk.Error {
	for _, eventCache := range gm.deliverTxEventCacheList {
		eventList, err := gm.storage.GetTimeEventList(ctx, eventCache.UnixTime)
		if err != nil {
			return err
		}
		if eventList == nil {
			eventList = &types.TimeEventList{Events: []types.Event{}}
		}
		eventList.Events = append(eventList.Events, eventCache.EventList...)
		if err := gm.storage.SetTimeEventList(ctx, eventCache.UnixTime, eventList); err != nil {
			return err
		}
	}
	gm.deliverTxEventCacheList = []*model.EventCache{}
	return nil
}

// ClearEventCache - clear event cache
// clear event cache will only be committed at the beginblocker
func (gm *GlobalManager) ClearEventCache(ctx sdk.Context) sdk.Error {
	gm.deliverTxEventCacheList = []*model.EventCache{}
	return nil
}

func (gm *GlobalManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.GlobalTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.GlobalTablesIR)

	ctx.Logger().Info(fmt.Sprintf("%s state parsed", filepath))

	// import table.TimeEventLists
	for _, v := range table.GlobalTimeEventLists {
		newevents := make([]types.Event, 0)
		events := v.TimeEventList.Events
		for _, event := range events {
			switch e := event.(type) {
			case posttypes.RewardEventV1:
				newevent := posttypes.RewardEvent{
					PostAuthor: e.PostAuthor,
					PostID:     e.PostID,
					Consumer:   e.Consumer,
					Evaluate:   types.NewMiniDollarFromTestnetCoin(e.Evaluate),
					FromApp:    e.FromApp,
				}
				newevents = append(newevents, newevent)
				ctx.Logger().Info(fmt.Sprintf("event conversion: from %+v to %+v ", e, newevent))
			default:
				newevents = append(newevents, event)
			}
		}
		err := gm.storage.SetTimeEventList(ctx, v.UnixTime, &types.TimeEventList{
			Events: newevents,
		})
		if err != nil {
			return err
		}
	}
	// import table.GlobalStakeStats
	for _, v := range table.GlobalStakeStats {
		err := gm.storage.SetLinoStakeStat(ctx, v.Day, &v.StakeStat)
		if err != nil {
			return err
		}
	}
	// import table.Misc
	misc := table.GlobalMisc

	err = gm.storage.SetInflationPool(ctx, &misc.InflationPool)
	if err != nil {
		return err
	}

	err = gm.storage.SetGlobalTime(ctx, &misc.Time)
	if err != nil {
		return err
	}

	// type diff in IR
	var cwindow types.MiniDollar
	if misc.ConsumptionMeta.IsConsumptionWindowDollarUnit {
		cwindow = types.NewMiniDollarFromInt(misc.ConsumptionMeta.ConsumptionWindow.Amount)
	} else {
		cwindow = types.NewMiniDollarFromTestnetCoin(misc.ConsumptionMeta.ConsumptionWindow)
		ctx.Logger().Info(
			fmt.Sprintf("consumption window converting, from %+v to %+v ",
				misc.ConsumptionMeta.ConsumptionWindow,
				cwindow,
			))
	}
	err = gm.storage.SetConsumptionMeta(ctx, &model.ConsumptionMeta{
		ConsumptionFrictionRate: sdk.MustNewDecFromStr(
			misc.ConsumptionMeta.ConsumptionFrictionRate),
		ConsumptionWindow:            cwindow,
		ConsumptionRewardPool:        misc.ConsumptionMeta.ConsumptionRewardPool,
		ConsumptionFreezingPeriodSec: misc.ConsumptionMeta.ConsumptionFreezingPeriodSec,
	})
	if err != nil {
		return err
	}

	err = gm.storage.SetTPS(ctx, &model.TPS{
		CurrentTPS: sdk.MustNewDecFromStr(misc.TPS.CurrentTPS),
		MaxTPS:     sdk.MustNewDecFromStr(misc.TPS.MaxTPS),
	})
	if err != nil {
		return err
	}

	ctx.Logger().Info(fmt.Sprintf("%s state imported", filepath))
	return nil
}

// get and set params
// TODO add more change methods
// func (gm *GlobalManager) ChangeGlobalInflationParam(ctx sdk.Context, InfraAllocation sdk.Dec,
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
// func (gm *GlobalManager) ChangeInfraInternalInflationParam(
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
