package model

import (
	"strconv"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
)

var (
	TimeEventListSubStore = []byte{0x00} // SubStore for time event list
	TimeSubStore          = []byte{0x01} // SubStore for time
	EventErrorSubStore    = []byte{0x02} // SubStore failed events and errors
	BCErrorSubStore       = []byte{0x03} // SubStore failed blockchain event and errors
)

// GlobalStorage - global storage
type GlobalStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewGlobalStorage - new global storage
func NewGlobalStorage(key sdk.StoreKey, cdc *wire.Codec) GlobalStorage {
	return GlobalStorage{
		key: key,
		cdc: cdc,
	}
}

func (gs GlobalStorage) CanEncode(event types.Event) bool {
	_, err := gs.cdc.MarshalBinaryLengthPrefixed(event)
	return err == nil
}

// GetTimeEventList - get time event list at given unix time
func (gs GlobalStorage) GetTimeEventList(ctx sdk.Context, unixTime int64) *types.TimeEventList {
	store := ctx.KVStore(gs.key)
	listByte := store.Get(GetTimeEventListKey(unixTime))
	// event doesn't exist
	if listByte == nil {
		return &types.TimeEventList{
			Events: nil,
		}
	}
	lst := new(types.TimeEventList)
	gs.cdc.MustUnmarshalBinaryLengthPrefixed(listByte, lst)
	return lst
}

// SetTimeEventList - set time event list at given unix time
func (gs GlobalStorage) SetTimeEventList(ctx sdk.Context, unixTime int64, lst *types.TimeEventList) {
	store := ctx.KVStore(gs.key)
	listByte := gs.cdc.MustMarshalBinaryLengthPrefixed(*lst)
	store.Set(GetTimeEventListKey(unixTime), listByte)
}

// RemoveTimeEventList - remove time event list at given unix time
func (gs GlobalStorage) RemoveTimeEventList(ctx sdk.Context, unixTime int64) {
	store := ctx.KVStore(gs.key)
	store.Delete(GetTimeEventListKey(unixTime))
}

// GetGlobalTime - get global time from KVStore
func (gs GlobalStorage) GetGlobalTime(ctx sdk.Context) *GlobalTime {
	store := ctx.KVStore(gs.key)
	timeBytes := store.Get(GetTimeKey())
	if timeBytes == nil {
		panic("Global Time is not Initialized at genesis")
	}
	globalTime := new(GlobalTime)
	gs.cdc.MustUnmarshalBinaryLengthPrefixed(timeBytes, globalTime)
	return globalTime
}

// SetGlobalTime - set global time to KVStore
func (gs GlobalStorage) SetGlobalTime(ctx sdk.Context, globalTime *GlobalTime) {
	store := ctx.KVStore(gs.key)
	timeBytes := gs.cdc.MustMarshalBinaryLengthPrefixed(*globalTime)
	store.Set(GetTimeKey(), timeBytes)
}

// GetEventErrors - get global time from KVStore
func (gs GlobalStorage) GetEventErrors(ctx sdk.Context) []EventError {
	store := ctx.KVStore(gs.key)
	bz := store.Get(GetEventErrorKey())
	if bz == nil {
		return nil
	}
	errors := make([]EventError, 0)
	gs.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &errors)
	return errors
}

// SetGlobalTime - set global time to KVStore
func (gs GlobalStorage) SetEventErrors(ctx sdk.Context, errs []EventError) {
	store := ctx.KVStore(gs.key)
	bz := gs.cdc.MustMarshalBinaryLengthPrefixed(errs)
	store.Set(GetEventErrorKey(), bz)
}

// GetEventErrors - get global time from KVStore
func (gs GlobalStorage) GetBCErrors(ctx sdk.Context) []types.BCEventErr {
	store := ctx.KVStore(gs.key)
	bz := store.Get(GetBCErrorKey())
	if bz == nil {
		return nil
	}
	errors := make([]types.BCEventErr, 0)
	gs.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &errors)
	return errors
}

// SetGlobalTime - set global time to KVStore
func (gs GlobalStorage) SetBCErrors(ctx sdk.Context, errs []types.BCEventErr) {
	store := ctx.KVStore(gs.key)
	bz := gs.cdc.MustMarshalBinaryLengthPrefixed(errs)
	store.Set(GetBCErrorKey(), bz)
}

func (gs GlobalStorage) PartialStoreMap(ctx sdk.Context) utils.StoreMap {
	store := ctx.KVStore(gs.key)
	stores := []utils.SubStore{
		{
			Store:      store,
			Prefix:     TimeEventListSubStore,
			ValCreator: func() interface{} { return new(types.TimeEventList) },
			Decoder:    gs.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
	}
	return utils.NewStoreMap(stores)
}

// GetTimeEventListKey - get time event list from KVStore
func GetTimeEventListKey(unixTime int64) []byte {
	return append(TimeEventListSubStore, strconv.FormatInt(unixTime, 10)...)
}

// GetTimeKey - "time substore"
func GetTimeKey() []byte {
	return TimeSubStore
}

// GetEventErrorKey - "event error substore"
func GetEventErrorKey() []byte {
	return EventErrorSubStore
}

// GetEventErrorKey - "bc event error substore"
func GetBCErrorKey() []byte {
	return BCErrorSubStore
}
