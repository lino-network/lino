package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/types"
)

var (
	heightEventListSubStore         = []byte{0x00} // SubStore for height event list
	timeEventListSubStore           = []byte{0x01} // SubStore for time event list
	statisticsSubStore              = []byte{0x02} // SubStore for statistics
	globalMetaSubStore              = []byte{0x03} // SubStore for global meta
	allocationSubStore              = []byte{0x04} // SubStore for allocation
	inflationPoolSubStore           = []byte{0x05} // SubStore for allocation
	infraInternalAllocationSubStore = []byte{0x06} // SubStore for infrat internal allocation
	consumptionMetaSubStore         = []byte{0x07} // SubStore for consumption meta
)

const eventTypePostReward = 0x1
const eventTypeDonateReward = 0x2

type GlobalManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewGlobalManager(key sdk.StoreKey) GlobalManager {
	cdc := wire.NewCodec()
	gm := GlobalManager{
		key: key,
		cdc: cdc,
	}
	return gm
}

func (gm GlobalManager) InitGlobalState(ctx sdk.Context, state genesis.GlobalState) error {
	globalMeta := &GlobalMeta{
		TotalLino:  sdk.NewRat(state.TotalLino),
		GrowthRate: state.GrowthRate,
	}

	if err := gm.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}
	if err := gm.SetGlobalStatistics(ctx, &GlobalStatistics{}); err != nil {
		return err
	}
	if !state.InfraAllocation.
		Add(state.ContentCreatorAllocation).
		Add(state.DeveloperAllocation).
		Add(state.ValidatorAllocation).
		Equal(sdk.NewRat(1)) {
		return ErrInflationGenesisError()
	}

	globalAllocation := &GlobalAllocation{
		InfraAllocation:          state.InfraAllocation,
		ContentCreatorAllocation: state.ContentCreatorAllocation,
		DeveloperAllocation:      state.DeveloperAllocation,
		ValidatorAllocation:      state.ValidatorAllocation,
	}
	if err := gm.SetGlobalAllocation(ctx, globalAllocation); err != nil {
		return err
	}
	inflaInflationCoin, err := types.LinoToCoin(types.LNO(sdk.Rat(globalMeta.TotalLino).Mul(globalMeta.GrowthRate).Mul(globalAllocation.InfraAllocation)))
	if err != nil {
		return err
	}

	contentCreatorCoin, err := types.LinoToCoin(types.LNO(sdk.Rat(globalMeta.TotalLino).Mul(globalMeta.GrowthRate).Mul(globalAllocation.ContentCreatorAllocation)))
	if err != nil {
		return err
	}
	developerCoin, err := types.LinoToCoin(types.LNO(sdk.Rat(globalMeta.TotalLino).Mul(globalMeta.GrowthRate).Mul(globalAllocation.DeveloperAllocation)))
	if err != nil {
		return err
	}
	validatorCoin, err := types.LinoToCoin(types.LNO(sdk.Rat(globalMeta.TotalLino).Mul(globalMeta.GrowthRate).Mul(globalAllocation.ValidatorAllocation)))
	if err != nil {
		return err
	}
	inflationPool := &InflationPool{
		InfraInflationPool:          inflaInflationCoin,
		ContentCreatorInflationPool: contentCreatorCoin,
		DeveloperInflationPool:      developerCoin,
		ValidatorInflationPool:      validatorCoin,
	}
	if err := gm.SetInflationPool(ctx, inflationPool); err != nil {
		return err
	}

	consumptionMeta := &ConsumptionMeta{
		ConsumptionFrictionRate: state.ConsumptionFrictionRate,
		FreezingPeriodHr:        state.FreezingPeriodHr,
	}
	if err := gm.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}
	return nil
}

func (gm GlobalManager) GetHeightEventList(ctx sdk.Context, key EventListKey) (*HeightEventList, sdk.Error) {
	store := ctx.KVStore(gm.key)
	listByte := store.Get(GetHeightEventListKey(key))
	if listByte == nil {
		return nil, ErrEventNotFound(GetHeightEventListKey(key))
	}
	lst := new(HeightEventList)
	if err := gm.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return lst, nil
}

func (gm GlobalManager) SetHeightEventList(ctx sdk.Context, key EventListKey, lst *HeightEventList) sdk.Error {
	store := ctx.KVStore(gm.key)
	listByte, err := gm.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetHeightEventListKey(key), listByte)
	return nil
}

func (gm GlobalManager) GetTimeEventList(ctx sdk.Context, key EventListKey) (*TimeEventList, sdk.Error) {
	store := ctx.KVStore(gm.key)
	listByte := store.Get(GetTimeEventListKey(key))
	if listByte == nil {
		return nil, ErrEventNotFound(GetTimeEventListKey(key))
	}
	lst := new(TimeEventList)
	if err := gm.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return lst, nil
}

func (gm GlobalManager) SetTimeEventList(ctx sdk.Context, key EventListKey, lst *TimeEventList) sdk.Error {
	store := ctx.KVStore(gm.key)
	listByte, err := gm.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetTimeEventListKey(key), listByte)
	return nil
}

func (gm GlobalManager) RemoveHeightEventList(ctx sdk.Context, key EventListKey) sdk.Error {
	store := ctx.KVStore(gm.key)
	store.Delete(GetHeightEventListKey(key))
	return nil
}

func (gm GlobalManager) RemoveTimeEventList(ctx sdk.Context, key EventListKey) sdk.Error {
	store := ctx.KVStore(gm.key)
	store.Delete(GetTimeEventListKey(key))
	return nil
}

func (gm GlobalManager) GetGlobalStatistics(ctx sdk.Context) (*GlobalStatistics, sdk.Error) {
	store := ctx.KVStore(gm.key)
	statisticsBytes := store.Get(GetGlobalStatisticsKey())
	if statisticsBytes == nil {
		return nil, ErrGlobalStatisticsNotFound()
	}
	statistics := new(GlobalStatistics)
	if err := gm.cdc.UnmarshalJSON(statisticsBytes, statistics); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return statistics, nil
}

func (gm GlobalManager) SetGlobalStatistics(ctx sdk.Context, statistics *GlobalStatistics) sdk.Error {
	store := ctx.KVStore(gm.key)
	statisticsBytes, err := gm.cdc.MarshalJSON(*statistics)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetGlobalStatisticsKey(), statisticsBytes)
	return nil
}

func (gm GlobalManager) GetGlobalMeta(ctx sdk.Context) (*GlobalMeta, sdk.Error) {
	store := ctx.KVStore(gm.key)
	globalMetaBytes := store.Get(GetGlobalMetaKey())
	if globalMetaBytes == nil {
		return nil, ErrGlobalMetaNotFound()
	}
	globalMeta := new(GlobalMeta)
	if err := gm.cdc.UnmarshalJSON(globalMetaBytes, globalMeta); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return globalMeta, nil
}

func (gm GlobalManager) SetGlobalMeta(ctx sdk.Context, globalMeta *GlobalMeta) sdk.Error {
	store := ctx.KVStore(gm.key)
	globalMetaBytes, err := gm.cdc.MarshalJSON(*globalMeta)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetGlobalMetaKey(), globalMetaBytes)
	return nil
}

func (gm GlobalManager) GetGlobalAllocation(ctx sdk.Context) (*GlobalAllocation, sdk.Error) {
	store := ctx.KVStore(gm.key)
	allocationBytes := store.Get(GetAllocationKey())
	if allocationBytes == nil {
		return nil, ErrGlobalAllocationNotFound()
	}
	allocation := new(GlobalAllocation)
	if err := gm.cdc.UnmarshalJSON(allocationBytes, allocation); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return allocation, nil
}

func (gm GlobalManager) SetGlobalAllocation(ctx sdk.Context, allocation *GlobalAllocation) sdk.Error {
	store := ctx.KVStore(gm.key)
	allocationBytes, err := gm.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetAllocationKey(), allocationBytes)
	return nil
}

func (gm GlobalManager) GetInflationPool(ctx sdk.Context) (*InflationPool, sdk.Error) {
	store := ctx.KVStore(gm.key)
	inflationPoolBytes := store.Get(GetInflationPoolKey())
	if inflationPoolBytes == nil {
		return nil, ErrGlobalAllocationNotFound()
	}
	inflationPool := new(InflationPool)
	if err := gm.cdc.UnmarshalJSON(inflationPoolBytes, inflationPool); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return inflationPool, nil
}

func (gm GlobalManager) SetInflationPool(ctx sdk.Context, inflationPool *InflationPool) sdk.Error {
	store := ctx.KVStore(gm.key)
	inflationPoolBytes, err := gm.cdc.MarshalJSON(*inflationPool)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetInflationPoolKey(), inflationPoolBytes)
	return nil
}

func (gm GlobalManager) GetConsumptionMeta(ctx sdk.Context) (*ConsumptionMeta, sdk.Error) {
	store := ctx.KVStore(gm.key)
	consumptionMetaBytes := store.Get(GetConsumptionMetaKey())
	if consumptionMetaBytes == nil {
		return nil, ErrGlobalConsumptionMetaNotFound()
	}
	consumptionMeta := new(ConsumptionMeta)
	if err := gm.cdc.UnmarshalJSON(consumptionMetaBytes, consumptionMeta); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return consumptionMeta, nil
}

func (gm GlobalManager) SetConsumptionMeta(ctx sdk.Context, consumptionMeta *ConsumptionMeta) sdk.Error {
	store := ctx.KVStore(gm.key)
	consumptionMetaBytes, err := gm.cdc.MarshalJSON(*consumptionMeta)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetConsumptionMetaKey(), consumptionMetaBytes)
	return nil
}

func GetHeightEventListKey(eventListKey EventListKey) []byte {
	return append(heightEventListSubStore, eventListKey...)
}

func GetTimeEventListKey(eventListKey EventListKey) []byte {
	return append(timeEventListSubStore, eventListKey...)
}

func GetGlobalStatisticsKey() []byte {
	return statisticsSubStore
}

func GetGlobalMetaKey() []byte {
	return globalMetaSubStore
}

func GetAllocationKey() []byte {
	return allocationSubStore
}

func GetInflationPoolKey() []byte {
	return inflationPoolSubStore
}

func GetInfraInternalAllocationKey() []byte {
	return infraInternalAllocationSubStore
}

func GetConsumptionMetaKey() []byte {
	return consumptionMetaSubStore
}
