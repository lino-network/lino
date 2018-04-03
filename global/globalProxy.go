package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// GlobalProxy encapsulates all basic struct
type GlobalProxy struct {
	globalManager *GlobalManager `json:"global_manager"`
}

// NewGlobalProxy return the global proxy pointer
func NewGlobalProxy(gm *GlobalManager) *GlobalProxy {
	return &GlobalProxy{
		globalManager: gm,
	}
}

func (gp *GlobalProxy) RegisterEventAtHeight(ctx sdk.Context, height int64, event Event) sdk.Error {
	eventList, _ := gp.globalManager.GetHeightEventList(ctx, HeightToEventListKey(height))
	if eventList == nil {
		eventList = &HeightEventList{Events: []Event{}}
	}
	eventList.Events = append(eventList.Events, event)
	if err := gp.globalManager.SetHeightEventList(ctx, HeightToEventListKey(height), eventList); err != nil {
		return err
	}
	return nil
}

func (gp *GlobalProxy) RegisterEventAtTime(ctx sdk.Context, unixTime int64, event Event) sdk.Error {
	eventList, _ := gp.globalManager.GetTimeEventList(ctx, UnixTimeToEventListKey(unixTime))
	if eventList == nil {
		eventList = &TimeEventList{Events: []Event{}}
	}
	eventList.Events = append(eventList.Events, event)
	if err := gp.globalManager.SetTimeEventList(ctx, UnixTimeToEventListKey(unixTime), eventList); err != nil {
		return err
	}
	return nil
}

func (gp *GlobalProxy) GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	consumptionMeta, err := gp.globalManager.GetConsumptionMeta(ctx)
	if err != nil {
		return sdk.Rat{}, err
	}
	return consumptionMeta.ConsumptionFrictionRate, nil
}

func (gp *GlobalProxy) RegisterRedistributionEvent(ctx sdk.Context, event Event) sdk.Error {
	consumptionMeta, err := gp.globalManager.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	if err := gp.RegisterEventAtTime(ctx, ctx.BlockHeader().Time+(consumptionMeta.FreezingPeriodHr*3600), event); err != nil {
		return err
	}
	return nil
}

func (gp *GlobalProxy) AddRedistributeCoin(ctx sdk.Context, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return nil
	}
	inflationPool, err := gp.globalManager.GetInflationPool(ctx)
	if err != nil {
		return err
	}
	inflationPool.ContentCreatorInflationPool = inflationPool.ContentCreatorInflationPool.Plus(coin)

	if err := gp.globalManager.SetInflationPool(ctx, inflationPool); err != nil {
		return err
	}

	consumptionMeta, err := gp.globalManager.GetConsumptionMeta(ctx)
	if err != nil {
		return err
	}
	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Plus(coin)

	if err := gp.globalManager.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}
	return nil
}

func (gp *GlobalProxy) GetRewardAndPopFromWindow(ctx sdk.Context, coin types.Coin) (types.Coin, sdk.Error) {
	if coin.IsZero() {
		return types.Coin{}, nil
	}
	inflationPool, err := gp.globalManager.GetInflationPool(ctx)
	if err != nil {
		return types.Coin{}, err
	}

	consumptionMeta, err := gp.globalManager.GetConsumptionMeta(ctx)
	if err != nil {
		return types.Coin{}, err
	}

	reward := types.Coin{sdk.NewRat(coin.Amount, consumptionMeta.ConsumptionWindow.Amount).
		Mul(sdk.NewRat(inflationPool.ContentCreatorInflationPool.Amount)).Evaluate()}

	inflationPool.ContentCreatorInflationPool = inflationPool.ContentCreatorInflationPool.Minus(reward)
	if err := gp.globalManager.SetInflationPool(ctx, inflationPool); err != nil {
		return types.Coin{}, err
	}

	consumptionMeta.ConsumptionWindow = consumptionMeta.ConsumptionWindow.Minus(coin)

	if err := gp.globalManager.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return types.Coin{}, err
	}
	return reward, nil
}

func (gp *GlobalProxy) AddConsumption(ctx sdk.Context, coin types.Coin) sdk.Error {
	globalMeta, err := gp.globalManager.GetGlobalMeta(ctx)
	if err != nil {
		return err
	}
	globalMeta.CumulativeConsumption = globalMeta.CumulativeConsumption.Plus(coin)

	if err := gp.globalManager.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}
	return nil
}
