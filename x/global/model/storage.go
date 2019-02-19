package model

import (
	"strconv"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	timeEventListSubStore   = []byte{0x00} // SubStore for time event list
	globalMetaSubStore      = []byte{0x01} // SubStore for global meta
	inflationPoolSubStore   = []byte{0x02} // SubStore for allocation
	consumptionMetaSubStore = []byte{0x03} // SubStore for consumption meta
	tpsSubStore             = []byte{0x04} // SubStore for tps
	timeSubStore            = []byte{0x05} // SubStore for time
	linoStakeStatSubStore   = []byte{0x06} // SubStore for lino power statistic
)

// GlobalStorage - global storage
type GlobalStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewGlobalStorage - new global storage
func NewGlobalStorage(key sdk.StoreKey) GlobalStorage {
	cdc := wire.NewCodec()
	cdc.RegisterInterface((*param.Parameter)(nil), nil)
	cdc.RegisterConcrete(param.EvaluateOfContentValueParam{}, "param/contentValue", nil)
	cdc.RegisterConcrete(param.GlobalAllocationParam{}, "param/allocation", nil)
	cdc.RegisterConcrete(param.InfraInternalAllocationParam{}, "param/infaAllocation", nil)
	cdc.RegisterConcrete(param.VoteParam{}, "param/vote", nil)
	cdc.RegisterConcrete(param.ProposalParam{}, "param/proposal", nil)
	cdc.RegisterConcrete(param.DeveloperParam{}, "param/developer", nil)
	cdc.RegisterConcrete(param.ValidatorParam{}, "param/validator", nil)
	cdc.RegisterConcrete(param.CoinDayParam{}, "param/coinDay", nil)
	cdc.RegisterConcrete(param.BandwidthParam{}, "param/bandwidth", nil)
	cdc.RegisterConcrete(param.AccountParam{}, "param/account", nil)
	cdc.RegisterConcrete(param.PostParam{}, "param/post", nil)

	wire.RegisterCrypto(cdc)
	return GlobalStorage{
		key: key,
		cdc: cdc,
	}
}

// WireCodec - access to global storage codec
func (gs GlobalStorage) WireCodec() *wire.Codec {
	return gs.cdc
}

// InitGlobalStateWithConfig - initialization based on genesis config file
func (gs GlobalStorage) InitGlobalStateWithConfig(
	ctx sdk.Context, totalLino types.Coin, param InitParamList) sdk.Error {
	globalMeta := &GlobalMeta{
		TotalLinoCoin:                 totalLino,
		LastYearTotalLinoCoin:         totalLino,
		LastYearCumulativeConsumption: types.NewCoinFromInt64(0),
		CumulativeConsumption:         types.NewCoinFromInt64(0),
	}
	if err := gs.SetGlobalMeta(ctx, globalMeta); err != nil {
		return err
	}

	inflationPool := &InflationPool{}
	if err := gs.SetInflationPool(ctx, inflationPool); err != nil {
		return err
	}

	globalTime := &GlobalTime{}
	if err := gs.SetGlobalTime(ctx, globalTime); err != nil {
		return err
	}

	consumptionMeta := &ConsumptionMeta{
		ConsumptionFrictionRate:      param.ConsumptionFrictionRate,
		ConsumptionWindow:            types.NewCoinFromInt64(0),
		ConsumptionRewardPool:        types.NewCoinFromInt64(0),
		ConsumptionFreezingPeriodSec: param.ConsumptionFreezingPeriodSec,
	}
	if err := gs.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return err
	}
	tps := &TPS{
		CurrentTPS: sdk.ZeroRat(),
		MaxTPS:     param.MaxTPS,
	}
	if err := gs.SetTPS(ctx, tps); err != nil {
		return err
	}
	linoStakeStat := &LinoStakeStat{
		TotalConsumptionFriction: types.NewCoinFromInt64(0),
		TotalLinoStake:           types.NewCoinFromInt64(0),
		UnclaimedFriction:        types.NewCoinFromInt64(0),
		UnclaimedLinoStake:       types.NewCoinFromInt64(0),
	}
	if err := gs.SetLinoStakeStat(ctx, 0, linoStakeStat); err != nil {
		return err
	}
	return nil
}

// InitGlobalState - initialization based on code
func (gs GlobalStorage) InitGlobalState(
	ctx sdk.Context, totalLino types.Coin) sdk.Error {
	initParamList := InitParamList{
		MaxTPS: sdk.NewRat(1000),
		ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
		ConsumptionFrictionRate:      sdk.NewRat(5, 100),
	}
	return gs.InitGlobalStateWithConfig(ctx, totalLino, initParamList)
}

// GetTimeEventList - get time event list at given unix time
func (gs GlobalStorage) GetTimeEventList(ctx sdk.Context, unixTime int64) (*types.TimeEventList, sdk.Error) {
	store := ctx.KVStore(gs.key)
	listByte := store.Get(GetTimeEventListKey(unixTime))
	// event doesn't exist
	if listByte == nil {
		return nil, nil
	}
	lst := new(types.TimeEventList)
	if err := gs.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrFailedToUnmarshalTimeEventList(err)
	}
	return lst, nil
}

// SetTimeEventList - set time event list at given unix time
func (gs GlobalStorage) SetTimeEventList(ctx sdk.Context, unixTime int64, lst *types.TimeEventList) sdk.Error {
	store := ctx.KVStore(gs.key)
	listByte, err := gs.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrFailedToMarshalTimeEventList(err)
	}
	store.Set(GetTimeEventListKey(unixTime), listByte)
	return nil
}

// RemoveTimeEventList - remove time event list at given unix time
func (gs GlobalStorage) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	store := ctx.KVStore(gs.key)
	store.Delete(GetTimeEventListKey(unixTime))
	return nil
}

// SetLinoStakeStat - set lino power statistic at given day
func (gs GlobalStorage) SetLinoStakeStat(ctx sdk.Context, day int64, lps *LinoStakeStat) sdk.Error {
	store := ctx.KVStore(gs.key)
	lpsByte, err := gs.cdc.MarshalJSON(*lps)
	if err != nil {
		return ErrFailedToMarshalTimeEventList(err)
	}
	store.Set(GetLinoStakeStatKey(day), lpsByte)
	return nil
}

// GetLinoStakeStat - get lino power statistic at given day
func (gs GlobalStorage) GetLinoStakeStat(ctx sdk.Context, day int64) (*LinoStakeStat, sdk.Error) {
	store := ctx.KVStore(gs.key)
	linoStakeStatBytes := store.Get(GetLinoStakeStatKey(day))
	if linoStakeStatBytes == nil {
		return nil, ErrLinoStakeStatisticNotFound()
	}
	linoStakeStat := new(LinoStakeStat)
	if err := gs.cdc.UnmarshalJSON(linoStakeStatBytes, linoStakeStat); err != nil {
		return nil, ErrFailedToUnmarshalLinoStakeStatistic(err)
	}
	return linoStakeStat, nil
}

// GetGlobalMeta - get global meta from KVStore
func (gs GlobalStorage) GetGlobalMeta(ctx sdk.Context) (*GlobalMeta, sdk.Error) {
	store := ctx.KVStore(gs.key)
	globalMetaBytes := store.Get(GetGlobalMetaKey())
	if globalMetaBytes == nil {
		return nil, ErrGlobalMetaNotFound()
	}
	globalMeta := new(GlobalMeta)
	if err := gs.cdc.UnmarshalJSON(globalMetaBytes, globalMeta); err != nil {
		return nil, ErrFailedToUnmarshalGlobalMeta(err)
	}
	return globalMeta, nil
}

// SetGlobalMeta - set global meta to KVStore
func (gs GlobalStorage) SetGlobalMeta(ctx sdk.Context, globalMeta *GlobalMeta) sdk.Error {
	store := ctx.KVStore(gs.key)
	globalMetaBytes, err := gs.cdc.MarshalJSON(*globalMeta)
	if err != nil {
		return ErrFailedToMarshalGlobalMeta(err)
	}
	store.Set(GetGlobalMetaKey(), globalMetaBytes)
	return nil
}

// GetInflationPool - get inflation pool from KVStore
func (gs GlobalStorage) GetInflationPool(ctx sdk.Context) (*InflationPool, sdk.Error) {
	store := ctx.KVStore(gs.key)
	inflationPoolBytes := store.Get(GetInflationPoolKey())
	if inflationPoolBytes == nil {
		return nil, ErrInflationPoolNotFound()
	}
	inflationPool := new(InflationPool)
	if err := gs.cdc.UnmarshalJSON(inflationPoolBytes, inflationPool); err != nil {
		return nil, ErrFailedToUnmarshalInflationPool(err)
	}
	return inflationPool, nil
}

// SetInflationPool - set inflation pool to KVStore
func (gs GlobalStorage) SetInflationPool(ctx sdk.Context, inflationPool *InflationPool) sdk.Error {
	store := ctx.KVStore(gs.key)
	inflationPoolBytes, err := gs.cdc.MarshalJSON(*inflationPool)
	if err != nil {
		return ErrFailedToMarshalInflationPool(err)
	}
	store.Set(GetInflationPoolKey(), inflationPoolBytes)
	return nil
}

// GetConsumptionMeta - get consumption meta from KVStore
func (gs GlobalStorage) GetConsumptionMeta(ctx sdk.Context) (*ConsumptionMeta, sdk.Error) {
	store := ctx.KVStore(gs.key)
	consumptionMetaBytes := store.Get(GetConsumptionMetaKey())
	if consumptionMetaBytes == nil {
		return nil, ErrGlobalConsumptionMetaNotFound()
	}
	consumptionMeta := new(ConsumptionMeta)
	if err := gs.cdc.UnmarshalJSON(consumptionMetaBytes, consumptionMeta); err != nil {
		return nil, ErrFailedToUnmarshalConsumptionMeta(err)
	}
	return consumptionMeta, nil
}

// SetConsumptionMeta - set consumption meta to KVStore
func (gs GlobalStorage) SetConsumptionMeta(ctx sdk.Context, consumptionMeta *ConsumptionMeta) sdk.Error {
	store := ctx.KVStore(gs.key)
	consumptionMetaBytes, err := gs.cdc.MarshalJSON(*consumptionMeta)
	if err != nil {
		return ErrFailedToMarshalConsumptionMeta(err)
	}
	store.Set(GetConsumptionMetaKey(), consumptionMetaBytes)
	return nil
}

// GetTPS - get tps from KVStore
func (gs GlobalStorage) GetTPS(ctx sdk.Context) (*TPS, sdk.Error) {
	store := ctx.KVStore(gs.key)
	tpsBytes := store.Get(GetTPSKey())
	if tpsBytes == nil {
		return nil, ErrGlobalTPSNotFound()
	}
	tps := new(TPS)
	if err := gs.cdc.UnmarshalJSON(tpsBytes, tps); err != nil {
		return nil, ErrFailedToUnmarshalTPS(err)
	}
	return tps, nil
}

// SetTPS - set tps to KVStore
func (gs GlobalStorage) SetTPS(ctx sdk.Context, tps *TPS) sdk.Error {
	store := ctx.KVStore(gs.key)
	tpsBytes, err := gs.cdc.MarshalJSON(*tps)
	if err != nil {
		return ErrFailedToMarshalTPS(err)
	}
	store.Set(GetTPSKey(), tpsBytes)
	return nil
}

// GetGlobalTime - get global time from KVStore
func (gs GlobalStorage) GetGlobalTime(ctx sdk.Context) (*GlobalTime, sdk.Error) {
	store := ctx.KVStore(gs.key)
	timeBytes := store.Get(GetTimeKey())
	if timeBytes == nil {
		return nil, ErrGlobalTimeNotFound()
	}
	globalTime := new(GlobalTime)
	if err := gs.cdc.UnmarshalJSON(timeBytes, globalTime); err != nil {
		return nil, ErrFailedToUnmarshalTime(err)
	}
	return globalTime, nil
}

// SetGlobalTime - set global time to KVStore
func (gs GlobalStorage) SetGlobalTime(ctx sdk.Context, globalTime *GlobalTime) sdk.Error {
	store := ctx.KVStore(gs.key)
	timeBytes, err := gs.cdc.MarshalJSON(*globalTime)
	if err != nil {
		return ErrFailedToMarshalTime(err)
	}
	store.Set(GetTimeKey(), timeBytes)
	return nil
}

// GetLinoStakeStatKey - get lino power statistic at day from KVStore
func GetLinoStakeStatKey(day int64) []byte {
	return append(linoStakeStatSubStore, strconv.FormatInt(day, 10)...)
}

// GetTimeEventListKey - get time event list from KVStore
func GetTimeEventListKey(unixTime int64) []byte {
	return append(timeEventListSubStore, strconv.FormatInt(unixTime, 10)...)
}

// GetGlobalMetaKey - "global meta substore"
func GetGlobalMetaKey() []byte {
	return globalMetaSubStore
}

// GetInflationPoolKey - "inflation pool substore"
func GetInflationPoolKey() []byte {
	return inflationPoolSubStore
}

// GetConsumptionMetaKey - "consumption meta substore"
func GetConsumptionMetaKey() []byte {
	return consumptionMetaSubStore
}

// GetTPSKey - "tps substore"
func GetTPSKey() []byte {
	return tpsSubStore
}

// GetTimeKey - "time substore"
func GetTimeKey() []byte {
	return timeSubStore
}
