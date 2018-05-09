package param

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	allocationParamSubStore              = []byte{0x00} // SubStore for allocation
	infraInternalAllocationParamSubStore = []byte{0x01} // SubStore for infrat internal allocation
	evaluateOfContentValueParamSubStore  = []byte{0x02} // Substore for evaluate of content value
	developerParamSubStore               = []byte{0x03} // Substore for developer param
	voteParamSubStore                    = []byte{0x04} // Substore for vote param
	proposalParamSubStore                = []byte{0x05} // Substore for proposal param
	validatorParamSubStore               = []byte{0x06} // Substore for validator param
	coinDayParamSubStore                 = []byte{0x07} // Substore for coin day param
	bandwidthParamSubStore               = []byte{0x08} // Substore for bandwidth param
	accountParamSubstore                 = []byte{0x09} // Substore for account param
)

type ParamHolder struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewParamHolder(key sdk.StoreKey) ParamHolder {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	return ParamHolder{
		key: key,
		cdc: cdc,
	}
}

func (ph ParamHolder) WireCodec() *wire.Codec {
	return ph.cdc
}

func (ph ParamHolder) InjectParam(ctx sdk.Context, minimum types.Coin) error {
	err := ph.InitParam(ctx)
	if err != nil {
		return err
	}

	accountParam := &AccountParam{
		MinimumBalance: minimum,
	}
	if err := ph.setAccountParam(ctx, accountParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}
	return nil
}

func (ph ParamHolder) InitParam(ctx sdk.Context) error {
	globalAllocationParam := &GlobalAllocationParam{
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(50, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(10, 100),
	}
	if err := ph.setGlobalAllocationParam(ctx, globalAllocationParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	infraInternalAllocationParam := &InfraInternalAllocationParam{
		StorageAllocation: sdk.NewRat(50, 100),
		CDNAllocation:     sdk.NewRat(50, 100),
	}
	if err := ph.setInfraInternalAllocationParam(ctx, infraInternalAllocationParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	evaluateOfContentValueParam := &EvaluateOfContentValueParam{
		ConsumptionTimeAdjustBase:      3153600,
		ConsumptionTimeAdjustOffset:    5,
		NumOfConsumptionOnAuthorOffset: 7,
		TotalAmountOfConsumptionBase:   1000 * types.Decimals,
		TotalAmountOfConsumptionOffset: 5,
		AmountOfConsumptionExponent:    sdk.NewRat(8, 10),
	}
	if err := ph.setEvaluateOfContentValueParam(ctx, evaluateOfContentValueParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	developerParam := &DeveloperParam{
		DeveloperMinDeposit:           types.NewCoin(100000 * types.Decimals),
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(7),
	}
	if err := ph.setDeveloperParam(ctx, developerParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
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
	if err := ph.setValidatorParam(ctx, validatorParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
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
	if err := ph.setVoteParam(ctx, voteParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	// TODO need to init other proposal params
	proposalParam := &ProposalParam{
		TypeAProposalDecideHr: int64(24 * 7),
		NextProposalID:        int64(0),
	}
	if err := ph.setProposalParam(ctx, proposalParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	coinDayParam := &CoinDayParam{
		DaysToRecoverCoinDayStake:    int64(7),
		SecondsToRecoverCoinDayStake: int64(7 * 24 * 3600),
	}
	if err := ph.setCoinDayParam(ctx, coinDayParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	bandwidthParam := &BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoin(1 * types.Decimals),
	}
	if err := ph.setBandwidthParam(ctx, bandwidthParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	accountParam := &AccountParam{
		MinimumBalance: types.NewCoin(1 * types.Decimals),
	}
	if err := ph.setAccountParam(ctx, accountParam); err != nil {
		return ErrParamHolderGenesisFailed().TraceCause(err, "")
	}

	return nil
}

func (ph ParamHolder) GetEvaluateOfContentValueParam(
	ctx sdk.Context) (*EvaluateOfContentValueParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paraBytes := store.Get(GetEvaluateOfContentValueParamKey())
	if paraBytes == nil {
		return nil, ErrEvaluateOfContentValueParamNotFound()
	}
	para := new(EvaluateOfContentValueParam)
	if err := ph.cdc.UnmarshalJSON(paraBytes, para); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return para, nil
}

func (ph ParamHolder) GetGlobalAllocationParam(
	ctx sdk.Context) (*GlobalAllocationParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	allocationBytes := store.Get(GetAllocationParamKey())
	if allocationBytes == nil {
		return nil, ErrGlobalAllocationParamNotFound()
	}
	allocation := new(GlobalAllocationParam)
	if err := ph.cdc.UnmarshalJSON(allocationBytes, allocation); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return allocation, nil
}

func (ph ParamHolder) GetInfraInternalAllocationParam(
	ctx sdk.Context) (*InfraInternalAllocationParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	allocationBytes := store.Get(GetInfraInternalAllocationParamKey())
	if allocationBytes == nil {
		return nil, ErrInfraAllocationParamNotFound()
	}
	allocation := new(InfraInternalAllocationParam)
	if err := ph.cdc.UnmarshalJSON(allocationBytes, allocation); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return allocation, nil
}

func (ph ParamHolder) GetDeveloperParam(ctx sdk.Context) (*DeveloperParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetDeveloperParamKey())
	if paramBytes == nil {
		return nil, ErrDeveloperParamNotFound()
	}
	param := new(DeveloperParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (ph ParamHolder) GetVoteParam(ctx sdk.Context) (*VoteParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetVoteParamKey())
	if paramBytes == nil {
		return nil, ErrVoteParamNotFound()
	}
	param := new(VoteParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (ph ParamHolder) GetProposalParam(ctx sdk.Context) (*ProposalParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetProposalParamKey())
	if paramBytes == nil {
		return nil, ErrProposalParamNotFound()
	}
	param := new(ProposalParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (ph ParamHolder) GetValidatorParam(ctx sdk.Context) (*ValidatorParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetValidatorParamKey())
	if paramBytes == nil {
		return nil, ErrValidatorParamNotFound()
	}
	param := new(ValidatorParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (ph ParamHolder) GetCoinDayParam(ctx sdk.Context) (*CoinDayParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetCoinDayParamKey())
	if paramBytes == nil {
		return nil, ErrCoinDayParamNotFound()
	}
	param := new(CoinDayParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (ph ParamHolder) GetBandwidthParam(ctx sdk.Context) (*BandwidthParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetBandwidthParamKey())
	if paramBytes == nil {
		return nil, ErrBandwidthParamNotFound()
	}
	param := new(BandwidthParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return param, nil
}

func (ph ParamHolder) GetAccountParam(ctx sdk.Context) (*AccountParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetAccountParamKey())
	if paramBytes == nil {
		return nil, ErrAccountParamNotFound()
	}
	param := new(AccountParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrEventMarshalError(err)
	}
	return param, nil
}

func (ph ParamHolder) GetNextProposalID(ctx sdk.Context) (types.ProposalKey, sdk.Error) {
	param, err := ph.GetProposalParam(ctx)
	if err != nil {
		return types.ProposalKey(""), err
	}
	param.NextProposalID += 1
	if err := ph.setProposalParam(ctx, param); err != nil {
		return types.ProposalKey(""), err
	}
	return types.ProposalKey(strconv.FormatInt(param.NextProposalID, 10)), nil
}

func (ph ParamHolder) setValidatorParam(ctx sdk.Context, param *ValidatorParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetValidatorParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setEvaluateOfContentValueParam(
	ctx sdk.Context, para *EvaluateOfContentValueParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paraBytes, err := ph.cdc.MarshalJSON(*para)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetEvaluateOfContentValueParamKey(), paraBytes)
	return nil
}

func (ph ParamHolder) setGlobalAllocationParam(
	ctx sdk.Context, allocation *GlobalAllocationParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	allocationBytes, err := ph.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetAllocationParamKey(), allocationBytes)
	return nil
}

func (ph ParamHolder) setInfraInternalAllocationParam(
	ctx sdk.Context, allocation *InfraInternalAllocationParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	allocationBytes, err := ph.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetInfraInternalAllocationParamKey(), allocationBytes)
	return nil
}

func (ph ParamHolder) setDeveloperParam(ctx sdk.Context, param *DeveloperParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetDeveloperParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setVoteParam(ctx sdk.Context, param *VoteParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetVoteParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setProposalParam(ctx sdk.Context, param *ProposalParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetProposalParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setCoinDayParam(ctx sdk.Context, param *CoinDayParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetCoinDayParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setBandwidthParam(ctx sdk.Context, param *BandwidthParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	bandwidthBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetBandwidthParamKey(), bandwidthBytes)
	return nil
}

func (ph ParamHolder) setAccountParam(ctx sdk.Context, param *AccountParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	accountBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(GetAccountParamKey(), accountBytes)
	return nil
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

func GetCoinDayParamKey() []byte {
	return coinDayParamSubStore
}

func GetBandwidthParamKey() []byte {
	return bandwidthParamSubStore
}

func GetAccountParamKey() []byte {
	return accountParamSubstore
}
