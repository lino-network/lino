package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/global/model"
	"github.com/lino-network/lino/types"
)

// GlobalManager encapsulates all basic struct
type GlobalManager struct {
	globalStorage *model.GlobalStorage `json:"global_manager"`
}

// NewGlobalManager return the global proxy pointer
func NewGlobalManager(key sdk.StoreKey) *GlobalManager {
	return &GlobalManager{
		globalStorage: model.NewGlobalStorage(key),
	}
}

func (gm *GlobalManager) InitGlobalManager(ctx sdk.Context, state genesis.GlobalState) error {
	return gm.globalStorage.InitGlobalState(ctx, state)
}

func (gm *GlobalManager) registerEventAtHeight(ctx sdk.Context, height int64, event types.Event) sdk.Error {
	eventList, _ := gm.globalStorage.GetHeightEventList(ctx, height)
	if eventList == nil {
		eventList = &types.HeightEventList{Events: []types.Event{}}
	}
	eventList.Events = append(eventList.Events, event)
	if err := gm.globalStorage.SetHeightEventList(ctx, height, eventList); err != nil {
		return ErrGlobalManagerRegisterEventAtHeight(height).TraceCause(err, "")
	}
	return nil
}

func (gm *GlobalManager) registerEventAtTime(ctx sdk.Context, unixTime int64, event types.Event) sdk.Error {
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

func (gm *GlobalManager) GetHeightEventListAtHeight(ctx sdk.Context, height int64) *types.HeightEventList {
	eventList, _ := gm.globalStorage.GetHeightEventList(ctx, height)
	return eventList
}

func (gm *GlobalManager) RemoveHeightEventList(ctx sdk.Context, height int64) sdk.Error {
	return gm.globalStorage.RemoveHeightEventList(ctx, height)
}

func (gm *GlobalManager) GetTimeEventListAtTime(ctx sdk.Context, unixTime int64) *types.TimeEventList {
	eventList, _ := gm.globalStorage.GetTimeEventList(ctx, unixTime)
	return eventList
}

func (gm *GlobalManager) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	return gm.globalStorage.RemoveTimeEventList(ctx, unixTime)
}

func (gm *GlobalManager) GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
	if err != nil {
		return sdk.Rat{}, err
	}
	return consumptionMeta.ConsumptionFrictionRate, nil
}

// register reward calculation event at 7 days later
func (gm *GlobalManager) RegisterContentRewardEvent(ctx sdk.Context, event types.Event) sdk.Error {
	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	if err := gm.registerEventAtTime(ctx, ctx.BlockHeader().Time+(consumptionMeta.FreezingPeriodHr*3600), event); err != nil {
		return err
	}
	return nil
}

// register coin return event with a time interval
func (gm *GlobalManager) RegisterCoinReturnEvent(ctx sdk.Context, event types.Event) sdk.Error {
	for i := int64(1); i <= types.CoinReturnTimes; i++ {
		if err := gm.registerEventAtTime(ctx, ctx.BlockHeader().Time+(types.CoinReturnIntervalHr*3600*i), event); err != nil {
			return err
		}
	}
	return nil
}

func (gm *GlobalManager) RegisterProposalDecideEvent(ctx sdk.Context, event types.Event) sdk.Error {
	if err := gm.registerEventAtTime(ctx, ctx.BlockHeader().Time+(types.ProposalDecideHr*3600), event); err != nil {
		return err
	}
	return nil
}

// put a friction of user consumption to reward pool
func (gm *GlobalManager) AddConsumptionFrictionToRewardPool(ctx sdk.Context, coin types.Coin) sdk.Error {
	// skip micro micro payment (etc: 0.0001 LNO)
	if coin.IsZero() {
		return nil
	}

	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
	if err != nil {
		return ErrAddConsumptionFrictionToRewardPool().TraceCause(err, "")
	}

	// reward pool consists of a small friction of user consumption and hourly content creator reward
	// consumption window will be used to calculate the percentage of reward to claim for this consumption
	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Plus(coin)
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Plus(coin)

	if err := gm.globalStorage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return ErrAddConsumptionFrictionToRewardPool().TraceCause(err, "")
	}
	return nil
}

// after 7 days, one consumption needs to claim its reward from consumption reward pool
func (gm *GlobalManager) GetRewardAndPopFromWindow(ctx sdk.Context, coin types.Coin) (types.Coin, sdk.Error) {
	if coin.IsZero() {
		return types.NewCoin(0), nil
	}

	consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
	if err != nil {
		return types.NewCoin(0), ErrGetRewardAndPopFromWindow().TraceCause(err, "")
	}

	// reward = (consumption reward pool) * ((this consumption) / (total consumption in 7 days window))
	reward := types.RatToCoin(consumptionMeta.ConsumptionRewardPool.ToRat().
		Mul(coin.ToRat().Quo(consumptionMeta.ConsumptionWindow.ToRat())))

	consumptionMeta.ConsumptionRewardPool = consumptionMeta.ConsumptionRewardPool.Minus(reward)
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Minus(coin)

	if err := gm.globalStorage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return types.NewCoin(0), ErrGetRewardAndPopFromWindow().TraceCause(err, "")
	}
	return reward, nil
}

// add consumption to global meta, which is used to compute GDP
func (gm *GlobalManager) AddConsumption(ctx sdk.Context, coin types.Coin) sdk.Error {
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

func (gm *GlobalManager) AddToValidatorInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error {
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

func (gm *GlobalManager) GetValidatorHourlyInflation(ctx sdk.Context, pastHours int64) (types.Coin, sdk.Error) {
	pool, getErr := gm.globalStorage.GetInflationPool(ctx)
	if getErr != nil {
		return types.NewCoin(0), getErr
	}

	resRat := pool.ValidatorInflationPool.ToRat().Mul(sdk.NewRat(1, types.HoursPerYear-pastHours+1))
	resCoin := types.RatToCoin(resRat)
	pool.ValidatorInflationPool = pool.ValidatorInflationPool.Minus(resCoin)

	if err := gm.globalStorage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

func (gm *GlobalManager) GetInfraHourlyInflation(ctx sdk.Context, pastHours int64) (types.Coin, sdk.Error) {
	pool, getErr := gm.globalStorage.GetInflationPool(ctx)
	if getErr != nil {
		return types.NewCoin(0), getErr
	}

	resRat := pool.ValidatorInflationPool.ToRat().Mul(sdk.NewRat(1, types.HoursPerYear-pastHours+1))
	resCoin := types.RatToCoin(resRat)
	pool.InfraInflationPool = pool.InfraInflationPool.Minus(resCoin)

	if err := gm.globalStorage.SetInflationPool(ctx, pool); err != nil {
		return types.NewCoin(0), err
	}
	return resCoin, nil
}

func (gm *GlobalManager) ChangeInfraInternalInflation(ctx sdk.Context, StorageAllocation sdk.Rat, CDNAllocation sdk.Rat) sdk.Error {
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
func (gm *GlobalManager) UpdateTPS(ctx sdk.Context, lastBlockTime int64) sdk.Error {
	tps, err := gm.globalStorage.GetTPS(ctx)
	if err != nil {
		return err
	}
	tps.CurrentTPS = sdk.NewRat(Txs, ctx.BlockHeader().Time-lastBlockTime)
	if tps.CurrentTPS.GT(tps.MaxTPS) {
		tps.MaxTPS = tps.CurrentTPS
	}

	if err := gm.globalStorage.SetTPS(ctx, tps); err != nil {
		return err
	}
	return nil
}

func (gm *GlobalManager) ChangeGlobalInflation(ctx sdk.Context, InfraAllocation sdk.Rat,
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

func (gm *GlobalManager) GetTPSCapacityRatio(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	tps, err := gm.globalStorage.GetTPS(ctx)
	if err != nil {
		return sdk.ZeroRat, err
	}
	return tps.CurrentTPS.Quo(tps.MaxTPS)
}
