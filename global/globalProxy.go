package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// GlobalProxy encapsulates all basic struct
type GlobalProxy struct {
	writeGlobalMeta              bool                     `json:"write_global_meta"`
	writeStatistics              bool                     `json:"write_statistics"`
	writeGlobalAllocation        bool                     `json:"write_global_allocation"`
	writeInfraInternalAllocation bool                     `json:"write_infra_internal_allocation"`
	writeConsumptionMeta         bool                     `json:"write_consumption_meta"`
	globalManager                *GlobalManager           `json:"global_manager"`
	globalMeta                   *GlobalMeta              `json:"global_meta"`
	globalStatistics             *GlobalStatistics        `json:"statistics"`
	globalAllocation             *GlobalAllocation        `json:"global_allocation"`
	infraInternalAllocation      *InfraInternalAllocation `json:"infra_internal_allocation"`
	inflationPool                *InflationPool           `json:"inflation_pool"`
	consumptionMeta              *ConsumptionMeta         `json:"consumption_meta"`
}

// NewGlobalProxy return the global proxy pointer
func NewGlobalProxy(gm *GlobalManager) *GlobalProxy {
	return &GlobalProxy{
		globalManager: gm,
	}
}

func (gp *GlobalProxy) RegisterEventAtHeight(ctx sdk.Context, height types.Height, event Event) sdk.Error {
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

func (gp *GlobalProxy) GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	if err := gp.checkConsumptionMeta(ctx); err != nil {
		return sdk.Rat{}, err
	}
	return gp.consumptionMeta.ConsumptionFrictionRate, nil
}

func (gp *GlobalProxy) checkGlobalMeta(ctx sdk.Context) (err sdk.Error) {
	if gp.globalMeta == nil {
		gp.globalMeta, err = gp.globalManager.GetGlobalMeta(ctx)
	}
	return err
}

func (gp *GlobalProxy) checkGlobalStatistics(ctx sdk.Context) (err sdk.Error) {
	if gp.globalStatistics == nil {
		gp.globalStatistics, err = gp.globalManager.GetGlobalStatistics(ctx)
	}
	return err
}

func (gp *GlobalProxy) checkGlobalAllocation(ctx sdk.Context) (err sdk.Error) {
	if gp.globalAllocation == nil {
		gp.globalAllocation, err = gp.globalManager.GetGlobalAllocation(ctx)
	}
	return err
}

func (gp *GlobalProxy) checkConsumptionMeta(ctx sdk.Context) (err sdk.Error) {
	if gp.consumptionMeta == nil {
		gp.consumptionMeta, err = gp.globalManager.GetConsumptionMeta(ctx)
	}
	return err
}
