package model

import (
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
)

var (
	heightEventListSubStore = []byte{0x00} // SubStore for height event list
	timeEventListSubStore   = []byte{0x01} // SubStore for time event list
	statisticsSubStore      = []byte{0x02} // SubStore for statistics
	globalMetaSubStore      = []byte{0x03} // SubStore for global meta
	inflationPoolSubStore   = []byte{0x04} // SubStore for allocation
	consumptionMetaSubStore = []byte{0x05} // SubStore for consumption meta
	tpsSubStore             = []byte{0x06} // SubStore for tps
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

	wire.RegisterCrypto(cdc)
	return GlobalStorage{
		key: key,
		cdc: cdc,
	}
}

func (gs GlobalStorage) WireCodec() *wire.Codec {
	return gs.cdc
}

func (gs GlobalStorage) InitGlobalState(
	ctx sdk.Context, totalLino types.Coin, param *param.GlobalAllocationParam) sdk.Error {
	globalMeta := &GlobalMeta{
		TotalLinoCoin:                 totalLino,
		LastYearCumulativeConsumption: types.NewCoinFromInt64(0),
		CumulativeConsumption:         types.NewCoinFromInt64(0),
		GrowthRate:                    sdk.NewRat(98, 1000),
		Ceiling:                       sdk.NewRat(98, 1000),
		Floor:                         sdk.NewRat(30, 1000),
	}

	if err := gs.SetGlobalMeta(ctx, globalMeta); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}
	if err := gs.SetGlobalStatistics(ctx, &GlobalStatistics{}); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	infraInflationCoin, err := types.RatToCoin(new(big.Rat).Mul(
		totalLino.ToRat(),
		(new(big.Rat).Mul(
			globalMeta.GrowthRate.GetRat(),
			param.InfraAllocation.GetRat()))))
	if err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}
	contentCreatorCoin, err := types.RatToCoin(new(big.Rat).Mul(
		totalLino.ToRat(),
		(new(big.Rat).Mul(
			globalMeta.GrowthRate.GetRat(),
			param.ContentCreatorAllocation.GetRat()))))
	if err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}
	developerCoin, err := types.RatToCoin(new(big.Rat).Mul(
		totalLino.ToRat(),
		(new(big.Rat).Mul(
			globalMeta.GrowthRate.GetRat(),
			param.DeveloperAllocation.GetRat()))))
	if err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}
	validatorCoin, err := types.RatToCoin(new(big.Rat).Mul(
		totalLino.ToRat(),
		(new(big.Rat).Mul(
			globalMeta.GrowthRate.GetRat(),
			param.ValidatorAllocation.GetRat()))))

	if err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	inflationPool := &InflationPool{
		InfraInflationPool:          infraInflationCoin,
		ContentCreatorInflationPool: contentCreatorCoin,
		DeveloperInflationPool:      developerCoin,
		ValidatorInflationPool:      validatorCoin,
	}
	if err := gs.SetInflationPool(ctx, inflationPool); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	consumptionMeta := &ConsumptionMeta{
		ConsumptionFrictionRate:     sdk.NewRat(5, 100),
		ReportStakeWindow:           sdk.ZeroRat,
		DislikeStakeWindow:          sdk.ZeroRat,
		ConsumptionWindow:           types.NewCoinFromInt64(0),
		ConsumptionRewardPool:       types.NewCoinFromInt64(0),
		ConsumptionFreezingPeriodHr: 24 * 7,
	}
	if err := gs.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}
	tps := &TPS{
		CurrentTPS: sdk.ZeroRat,
		MaxTPS:     sdk.NewRat(1000),
	}
	if err := gs.SetTPS(ctx, tps); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}
	return nil
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
		return nil, ErrEventUnmarshalError(err)
	}
	return lst, nil
}

func (gs GlobalStorage) SetTimeEventList(ctx sdk.Context, unixTime int64, lst *types.TimeEventList) sdk.Error {
	store := ctx.KVStore(gs.key)
	listByte, err := gs.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetTimeEventListKey(unixTime), listByte)
	return nil
}

func (gs GlobalStorage) RemoveTimeEventList(ctx sdk.Context, unixTime int64) sdk.Error {
	store := ctx.KVStore(gs.key)
	store.Delete(GetTimeEventListKey(unixTime))
	return nil
}

func (gs GlobalStorage) GetGlobalStatistics(ctx sdk.Context) (*GlobalStatistics, sdk.Error) {
	store := ctx.KVStore(gs.key)
	statisticsBytes := store.Get(GetGlobalStatisticsKey())
	if statisticsBytes == nil {
		return nil, ErrGlobalStatisticsNotFound()
	}
	statistics := new(GlobalStatistics)
	if err := gs.cdc.UnmarshalJSON(statisticsBytes, statistics); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return statistics, nil
}

func (gs GlobalStorage) SetGlobalStatistics(ctx sdk.Context, statistics *GlobalStatistics) sdk.Error {
	store := ctx.KVStore(gs.key)
	statisticsBytes, err := gs.cdc.MarshalJSON(*statistics)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetGlobalStatisticsKey(), statisticsBytes)
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
		return nil, ErrEventUnmarshalError(err)
	}
	return globalMeta, nil
}

func (gs GlobalStorage) SetGlobalMeta(ctx sdk.Context, globalMeta *GlobalMeta) sdk.Error {
	store := ctx.KVStore(gs.key)
	globalMetaBytes, err := gs.cdc.MarshalJSON(*globalMeta)
	if err != nil {
		return ErrEventMarshalError(err)
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
		return nil, ErrEventUnmarshalError(err)
	}
	return inflationPool, nil
}

func (gs GlobalStorage) SetInflationPool(ctx sdk.Context, inflationPool *InflationPool) sdk.Error {
	store := ctx.KVStore(gs.key)
	inflationPoolBytes, err := gs.cdc.MarshalJSON(*inflationPool)
	if err != nil {
		return ErrEventMarshalError(err)
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
		return nil, ErrEventUnmarshalError(err)
	}
	return consumptionMeta, nil
}

func (gs GlobalStorage) SetConsumptionMeta(ctx sdk.Context, consumptionMeta *ConsumptionMeta) sdk.Error {
	store := ctx.KVStore(gs.key)
	consumptionMetaBytes, err := gs.cdc.MarshalJSON(*consumptionMeta)
	if err != nil {
		return ErrEventMarshalError(err)
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
		return nil, ErrEventUnmarshalError(err)
	}
	return tps, nil
}

func (gs GlobalStorage) SetTPS(ctx sdk.Context, tps *TPS) sdk.Error {
	store := ctx.KVStore(gs.key)
	tpsBytes, err := gs.cdc.MarshalJSON(*tps)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetTPSKey(), tpsBytes)
	return nil
}

func GetHeightEventListKey(height int64) []byte {
	return append(heightEventListSubStore, strconv.FormatInt(height, 10)...)
}

func GetTimeEventListKey(unixTime int64) []byte {
	return append(timeEventListSubStore, strconv.FormatInt(unixTime, 10)...)
}

func GetGlobalStatisticsKey() []byte {
	return statisticsSubStore
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
