package model

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var (
	heightEventListSubStore              = []byte{0x00} // SubStore for height event list
	timeEventListSubStore                = []byte{0x01} // SubStore for time event list
	statisticsSubStore                   = []byte{0x02} // SubStore for statistics
	globalMetaSubStore                   = []byte{0x03} // SubStore for global meta
	allocationParamSubStore              = []byte{0x04} // SubStore for allocation
	inflationPoolSubStore                = []byte{0x05} // SubStore for allocation
	infraInternalAllocationParamSubStore = []byte{0x06} // SubStore for infrat internal allocation
	consumptionMetaSubStore              = []byte{0x07} // SubStore for consumption meta
	tpsSubStore                          = []byte{0x08} // SubStore for tps
	evaluateOfContentValueParamSubStore  = []byte{0x09} // Substore for evaluate of content value
	developerParamSubStore               = []byte{0x10} // Substore for developer param
	voteParamSubStore                    = []byte{0x11} // Substore for vote param
	proposalParamSubStore                = []byte{0x12} // Substore for proposal param
	validatorParamSubStore               = []byte{0x13} // Substore for validator param
)

type GlobalStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewGlobalStorage(key sdk.StoreKey) GlobalStorage {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	return GlobalStorage{
		key: key,
		cdc: cdc,
	}
}

func (gs GlobalStorage) WireCodec() *wire.Codec {
	return gs.cdc
}

func (gs GlobalStorage) InitGlobalState(ctx sdk.Context, totalLino types.Coin) error {
	globalMeta := &GlobalMeta{
		TotalLinoCoin:                 totalLino,
		LastYearCumulativeConsumption: types.NewCoin(0),
		CumulativeConsumption:         types.NewCoin(0),
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

	globalAllocationParam := &GlobalAllocationParam{
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(50, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(10, 100),
	}
	if err := gs.SetGlobalAllocationParam(ctx, globalAllocationParam); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	infraInternalAllocationParam := &InfraInternalAllocationParam{
		StorageAllocation: sdk.NewRat(50, 100),
		CDNAllocation:     sdk.NewRat(50, 100),
	}
	if err := gs.SetInfraInternalAllocationParam(ctx, infraInternalAllocationParam); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	infraInflationCoin := totalLino.ToRat().Mul(globalMeta.GrowthRate).Mul(globalAllocationParam.InfraAllocation)

	contentCreatorCoin := totalLino.ToRat().Mul(globalMeta.GrowthRate).Mul(globalAllocationParam.ContentCreatorAllocation)

	developerCoin := totalLino.ToRat().Mul(globalMeta.GrowthRate).Mul(globalAllocationParam.DeveloperAllocation)

	validatorCoin := totalLino.ToRat().Mul(globalMeta.GrowthRate).Mul(globalAllocationParam.ValidatorAllocation)

	inflationPool := &InflationPool{
		InfraInflationPool:          types.RatToCoin(infraInflationCoin),
		ContentCreatorInflationPool: types.RatToCoin(contentCreatorCoin),
		DeveloperInflationPool:      types.RatToCoin(developerCoin),
		ValidatorInflationPool:      types.RatToCoin(validatorCoin),
	}
	if err := gs.SetInflationPool(ctx, inflationPool); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	consumptionMeta := &ConsumptionMeta{
		ConsumptionFrictionRate:     sdk.NewRat(5, 100),
		ReportStakeWindow:           sdk.ZeroRat,
		DislikeStakeWindow:          sdk.ZeroRat,
		ConsumptionWindow:           types.NewCoin(0),
		ConsumptionRewardPool:       types.NewCoin(0),
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
	evaluateOfContentValueParam := &EvaluateOfContentValueParam{
		ConsumptionTimeAdjustBase:      3153600,
		ConsumptionTimeAdjustOffset:    5,
		NumOfConsumptionOnAuthorOffset: 7,
		TotalAmountOfConsumptionBase:   1000 * types.Decimals,
		TotalAmountOfConsumptionOffset: 5,
		AmountOfConsumptionExponent:    sdk.NewRat(8, 10),
	}
	if err := gs.SetEvaluateOfContentValueParam(ctx, evaluateOfContentValueParam); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	developerParam := &DeveloperParam{
		DeveloperMinDeposit:           types.NewCoin(100000 * types.Decimals),
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(7),
	}
	if err := gs.SetDeveloperParam(ctx, developerParam); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	validatorParam := &ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}
	if err := gs.SetValidatorParam(ctx, validatorParam); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	voteParam := &VoteParam{
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}
	if err := gs.SetVoteParam(ctx, voteParam); err != nil {
		return ErrGlobalStorageGenesisFailed().TraceCause(err, "")
	}

	// TODO need to init other proposal params
	proposalParam := &ProposalParam{
		NextProposalID: int64(0),
	}
	if err := gs.SetProposalParam(ctx, proposalParam); err != nil {
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
		return nil, ErrGlobalAllocationParamNotFound()
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

func (gs GlobalStorage) GetEvaluateOfContentValueParam(
	ctx sdk.Context) (*EvaluateOfContentValueParam, sdk.Error) {
	store := ctx.KVStore(gs.key)
	paraBytes := store.Get(GetEvaluateOfContentValueParamKey())
	if paraBytes == nil {
		return nil, ErrEvluateOfContentValueParam()
	}
	para := new(EvaluateOfContentValueParam)
	if err := gs.cdc.UnmarshalJSON(paraBytes, para); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return para, nil
}

func (gs GlobalStorage) SetEvaluateOfContentValueParam(
	ctx sdk.Context, para *EvaluateOfContentValueParam) sdk.Error {
	store := ctx.KVStore(gs.key)
	paraBytes, err := gs.cdc.MarshalJSON(*para)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetEvaluateOfContentValueParamKey(), paraBytes)
	return nil
}

func (gs GlobalStorage) GetGlobalAllocationParam(
	ctx sdk.Context) (*GlobalAllocationParam, sdk.Error) {
	store := ctx.KVStore(gs.key)
	allocationBytes := store.Get(GetAllocationParamKey())
	if allocationBytes == nil {
		return nil, ErrGlobalAllocationParamNotFound()
	}
	allocation := new(GlobalAllocationParam)
	if err := gs.cdc.UnmarshalJSON(allocationBytes, allocation); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return allocation, nil
}

func (gs GlobalStorage) SetGlobalAllocationParam(
	ctx sdk.Context, allocation *GlobalAllocationParam) sdk.Error {
	store := ctx.KVStore(gs.key)
	allocationBytes, err := gs.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetAllocationParamKey(), allocationBytes)
	return nil
}

func (gs GlobalStorage) GetInfraInternalAllocationParam(
	ctx sdk.Context) (*InfraInternalAllocationParam, sdk.Error) {
	store := ctx.KVStore(gs.key)
	allocationBytes := store.Get(GetInfraInternalAllocationParamKey())
	if allocationBytes == nil {
		return nil, ErrInfraAllocationParamNotFound()
	}
	allocation := new(InfraInternalAllocationParam)
	if err := gs.cdc.UnmarshalJSON(allocationBytes, allocation); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return allocation, nil
}

func (gs GlobalStorage) SetInfraInternalAllocationParam(
	ctx sdk.Context, allocation *InfraInternalAllocationParam) sdk.Error {
	store := ctx.KVStore(gs.key)
	allocationBytes, err := gs.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetInfraInternalAllocationParamKey(), allocationBytes)
	return nil
}

func (gs GlobalStorage) GetDeveloperParam(ctx sdk.Context) (*DeveloperParam, sdk.Error) {
	store := ctx.KVStore(gs.key)
	paramBytes := store.Get(GetDeveloperParamKey())
	if paramBytes == nil {
		return nil, ErrDeveloperParamNotFound()
	}
	param := new(DeveloperParam)
	if err := gs.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (gs GlobalStorage) SetDeveloperParam(ctx sdk.Context, param *DeveloperParam) sdk.Error {
	store := ctx.KVStore(gs.key)
	paramBytes, err := gs.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetDeveloperParamKey(), paramBytes)
	return nil
}

func (gs GlobalStorage) GetVoteParam(ctx sdk.Context) (*VoteParam, sdk.Error) {
	store := ctx.KVStore(gs.key)
	paramBytes := store.Get(GetVoteParamKey())
	if paramBytes == nil {
		return nil, ErrVoteParamNotFound()
	}
	param := new(VoteParam)
	if err := gs.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (gs GlobalStorage) SetVoteParam(ctx sdk.Context, param *VoteParam) sdk.Error {
	store := ctx.KVStore(gs.key)
	paramBytes, err := gs.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetVoteParamKey(), paramBytes)
	return nil
}

func (gs GlobalStorage) GetProposalParam(ctx sdk.Context) (*ProposalParam, sdk.Error) {
	store := ctx.KVStore(gs.key)
	paramBytes := store.Get(GetProposalParamKey())
	if paramBytes == nil {
		return nil, ErrProposalParamNotFound()
	}
	param := new(ProposalParam)
	if err := gs.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (gs GlobalStorage) SetProposalParam(ctx sdk.Context, param *ProposalParam) sdk.Error {
	store := ctx.KVStore(gs.key)
	paramBytes, err := gs.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetProposalParamKey(), paramBytes)
	return nil
}

func (gs GlobalStorage) GetValidatorParam(ctx sdk.Context) (*ValidatorParam, sdk.Error) {
	store := ctx.KVStore(gs.key)
	paramBytes := store.Get(GetValidatorParamKey())
	if paramBytes == nil {
		return nil, ErrValidatorParamNotFound()
	}
	param := new(ValidatorParam)
	if err := gs.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (gs GlobalStorage) SetValidatorParam(ctx sdk.Context, param *ValidatorParam) sdk.Error {
	store := ctx.KVStore(gs.key)
	paramBytes, err := gs.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetValidatorParamKey(), paramBytes)
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

func GetEvaluateOfContentValueParamKey() []byte {
	return evaluateOfContentValueParamSubStore
}

func GetAllocationParamKey() []byte {
	return allocationParamSubStore
}

func GetInfraInternalAllocationParamKey() []byte {
	return infraInternalAllocationParamSubStore
}

func GetDeveloperParamKey() []byte {
	return developerParamSubStore
}

func GetVoteParamKey() []byte {
	return voteParamSubStore
}

func GetValidatorParamKey() []byte {
	return validatorParamSubStore
}

func GetProposalParamKey() []byte {
	return proposalParamSubStore
}
