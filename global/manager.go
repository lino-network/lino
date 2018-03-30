package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

var (
	heightEventListSubStore         = []byte{0x00} // SubStore for height event list
	timeEventListSubStore           = []byte{0x01} // SubStore for time event list
	statisticsSubStore              = []byte{0x02} // SubStore for statistics
	globalMetaSubStore              = []byte{0x03} // SubStore for global meta
	allocationSubStore              = []byte{0x04} // SubStore for allocation
	infraInternalAllocationSubStore = []byte{0x05} // SubStore for infrat internal allocation
	consumptionMetaSubStore         = []byte{0x06} // SubStore for consumption meta
)

const eventTypePostReward = 0x1
const eventTypeDonateReward = 0x2

var _ = oldwire.RegisterInterface(
	struct{ Event }{},
	oldwire.ConcreteType{PostRewardEvent{}, eventTypePostReward},
	oldwire.ConcreteType{DonateRewardEvent{}, eventTypeDonateReward},
)

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

func (gm GlobalManager) removeHeightEventList(ctx sdk.Context, key EventListKey) sdk.Error {
	store := ctx.KVStore(gm.key)
	store.Delete(GetHeightEventListKey(key))
	return nil
}

func (gm GlobalManager) ExecuteHeightEvents(ctx sdk.Context, key EventListKey) sdk.Error {
	lst, err := gm.GetHeightEventList(ctx, key)
	if err != nil {
		return err
	}

	for _, event := range lst.Events {
		switch event := event.(type) {
		case PostRewardEvent:
			if err := event.execute(); err != nil {
				return err
			}
		case DonateRewardEvent:
			if err := event.execute(); err != nil {
				return err
			}
		default:
			return ErrWrongEventType()
		}
	}

	if err := gm.removeHeightEventList(ctx, key); err != nil {
		return err
	}
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

func GetInfraInternalAllocationKey() []byte {
	return infraInternalAllocationSubStore
}

func GetConsumptionMetaKey() []byte {
	return consumptionMetaSubStore
}
