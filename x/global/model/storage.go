package model

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/wire"
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
	timeSubStore            = []byte{0x05} // SubStore for tps
)

type GlobalStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

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

func (gs GlobalStorage) WireCodec() *wire.Codec {
	return gs.cdc
}

func (gs GlobalStorage) InitGlobalStateWithConfig(
	ctx sdk.Context, totalLino types.Coin, param InitParamList) sdk.Error {
	globalMeta := &GlobalMeta{
		TotalLinoCoin:                 totalLino,
		LastYearCumulativeConsumption: types.NewCoinFromInt64(0),
		CumulativeConsumption:         types.NewCoinFromInt64(0),
		GrowthRate:                    param.GrowthRate,
		Ceiling:                       param.Ceiling,
		Floor:                         param.Floor,
		AnnualInflation:               types.RatToCoin(totalLino.ToRat().Mul(param.GrowthRate)),
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
	return nil
}

func (gs GlobalStorage) InitGlobalState(
	ctx sdk.Context, totalLino types.Coin) sdk.Error {
	initParamList := InitParamList{
		GrowthRate: sdk.NewRat(98, 1000),
		Ceiling:    sdk.NewRat(98, 1000),
		Floor:      sdk.NewRat(3, 100),
		MaxTPS:     sdk.NewRat(1000),
		ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
		ConsumptionFrictionRate:      sdk.NewRat(5, 100),
	}
	return gs.InitGlobalStateWithConfig(ctx, totalLino, initParamList)
}

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

func (gs GlobalStorage) SetTimeEventList(ctx sdk.Context, unixTime int64, lst *types.TimeEventList) sdk.Error {
	store := ctx.KVStore(gs.key)
	listByte, err := gs.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrFailedToMarshalTimeEventList(err)
	}
	store.Set(GetTimeEventListKey(unixTime), listByte)
	return nil
}

func (gs GlobalStorage) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	store := ctx.KVStore(gs.key)
	store.Delete(GetTimeEventListKey(unixTime))
	return nil
}

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

func (gs GlobalStorage) SetGlobalMeta(ctx sdk.Context, globalMeta *GlobalMeta) sdk.Error {
	store := ctx.KVStore(gs.key)
	globalMetaBytes, err := gs.cdc.MarshalJSON(*globalMeta)
	if err != nil {
		return ErrFailedToMarshalGlobalMeta(err)
	}
	store.Set(GetGlobalMetaKey(), globalMetaBytes)
	return nil
}

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

func (gs GlobalStorage) SetInflationPool(ctx sdk.Context, inflationPool *InflationPool) sdk.Error {
	store := ctx.KVStore(gs.key)
	inflationPoolBytes, err := gs.cdc.MarshalJSON(*inflationPool)
	if err != nil {
		return ErrFailedToMarshalInflationPool(err)
	}
	store.Set(GetInflationPoolKey(), inflationPoolBytes)
	return nil
}

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

func (gs GlobalStorage) SetConsumptionMeta(ctx sdk.Context, consumptionMeta *ConsumptionMeta) sdk.Error {
	store := ctx.KVStore(gs.key)
	consumptionMetaBytes, err := gs.cdc.MarshalJSON(*consumptionMeta)
	if err != nil {
		return ErrFailedToMarshalConsumptionMeta(err)
	}
	store.Set(GetConsumptionMetaKey(), consumptionMetaBytes)
	return nil
}

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

func (gs GlobalStorage) SetTPS(ctx sdk.Context, tps *TPS) sdk.Error {
	store := ctx.KVStore(gs.key)
	tpsBytes, err := gs.cdc.MarshalJSON(*tps)
	if err != nil {
		return ErrFailedToMarshalTPS(err)
	}
	store.Set(GetTPSKey(), tpsBytes)
	return nil
}

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

func (gs GlobalStorage) SetGlobalTime(ctx sdk.Context, globalTime *GlobalTime) sdk.Error {
	store := ctx.KVStore(gs.key)
	timeBytes, err := gs.cdc.MarshalJSON(*globalTime)
	if err != nil {
		return ErrFailedToMarshalTime(err)
	}
	store.Set(GetTimeKey(), timeBytes)
	return nil
}

func GetTimeEventListKey(unixTime int64) []byte {
	return append(timeEventListSubStore, strconv.FormatInt(unixTime, 10)...)
}

func GetGlobalMetaKey() []byte {
	return globalMetaSubStore
}

func GetInflationPoolKey() []byte {
	return inflationPoolSubStore
}

func GetConsumptionMetaKey() []byte {
	return consumptionMetaSubStore
}

func GetTPSKey() []byte {
	return tpsSubStore
}

func GetTimeKey() []byte {
	return timeSubStore
}
