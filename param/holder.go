package param

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	allocationParamSubStore              = []byte{0x00} // SubStore for allocation
	infraInternalAllocationParamSubStore = []byte{0x01} // SubStore for infra internal allocation
	evaluateOfContentValueParamSubStore  = []byte{0x02} // Substore for evaluate of content value
	developerParamSubStore               = []byte{0x03} // Substore for developer param
	voteParamSubStore                    = []byte{0x04} // Substore for vote param
	proposalParamSubStore                = []byte{0x05} // Substore for proposal param
	validatorParamSubStore               = []byte{0x06} // Substore for validator param
	coinDayParamSubStore                 = []byte{0x07} // Substore for coin day param
	bandwidthParamSubStore               = []byte{0x08} // Substore for bandwidth param
	accountParamSubstore                 = []byte{0x09} // Substore for account param
	postParamSubStore                    = []byte{0x10} // Substore for evaluate of content value
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

func (ph ParamHolder) InitParam(ctx sdk.Context) error {
	globalAllocationParam := &GlobalAllocationParam{
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(65, 100),
		DeveloperAllocation:      sdk.NewRat(10, 100),
		ValidatorAllocation:      sdk.NewRat(5, 100),
	}
	if err := ph.setGlobalAllocationParam(ctx, globalAllocationParam); err != nil {
		return err
	}

	infraInternalAllocationParam := &InfraInternalAllocationParam{
		StorageAllocation: sdk.NewRat(50, 100),
		CDNAllocation:     sdk.NewRat(50, 100),
	}
	if err := ph.setInfraInternalAllocationParam(ctx, infraInternalAllocationParam); err != nil {
		return err
	}

	postParam := &PostParam{
		MicropaymentLimitation: types.NewCoinFromInt64(10 * types.Decimals),
		ReportOrUpvoteInterval: 24 * 3600,
	}
	if err := ph.setPostParam(ctx, postParam); err != nil {
		return err
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
		return err
	}

	developerParam := &DeveloperParam{
		DeveloperMinDeposit:           types.NewCoinFromInt64(1000000 * types.Decimals),
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(7),
	}
	if err := ph.setDeveloperParam(ctx, developerParam); err != nil {
		return err
	}

	validatorParam := &ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoinFromInt64(300000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoinFromInt64(20000 * types.Decimals),
		PenaltyMissCommit:             types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoinFromInt64(1000000 * types.Decimals),
		ValidatorListSize:             int64(21),
		AbsentCommitLimitation:        int64(100),
	}
	if err := ph.setValidatorParam(ctx, validatorParam); err != nil {
		return err
	}

	voteParam := &VoteParam{
		VoterMinDeposit:               types.NewCoinFromInt64(2000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoinFromInt64(2 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoinFromInt64(2 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}
	if err := ph.setVoteParam(ctx, voteParam); err != nil {
		return err
	}

	proposalParam := &ProposalParam{
		ContentCensorshipDecideHr:   int64(24 * 7),
		ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamDecideHr:   int64(24 * 7),
		ChangeParamPassRatio:  sdk.NewRat(70, 100),
		ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideHr:   int64(24 * 7),
		ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}
	if err := ph.setProposalParam(ctx, proposalParam); err != nil {
		return err
	}

	coinDayParam := &CoinDayParam{
		DaysToRecoverCoinDayStake:    int64(7),
		SecondsToRecoverCoinDayStake: int64(7 * 24 * 3600),
	}
	if err := ph.setCoinDayParam(ctx, coinDayParam); err != nil {
		return err
	}

	bandwidthParam := &BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
	}
	if err := ph.setBandwidthParam(ctx, bandwidthParam); err != nil {
		return err
	}

	accountParam := &AccountParam{
		MinimumBalance:                types.NewCoinFromInt64(1 * types.Decimals),
		RegisterFee:                   types.NewCoinFromInt64(1 * types.Decimals),
		BalanceHistoryBundleSize:      100,
		MaximumMicropaymentGrantTimes: 20,
		RewardHistoryBundleSize:       100,
	}
	if err := ph.setAccountParam(ctx, accountParam); err != nil {
		return err
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
		return nil, ErrFailedToUnmarshalEvaluateOfContentValueParam(err)
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
		return nil, ErrFailedToUnmarshalGlobalAllocationParam(err)
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
		return nil, ErrFailedToUnmarshalInfraInternalAllocationParam(err)
	}
	return allocation, nil
}

func (ph ParamHolder) GetPostParam(ctx sdk.Context) (*PostParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetPostParamKey())
	if paramBytes == nil {
		return nil, ErrPostParamNotFound()
	}
	param := new(PostParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrFailedToUnmarshalPostParam(err)
	}
	return param, nil
}

func (ph ParamHolder) GetDeveloperParam(ctx sdk.Context) (*DeveloperParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetDeveloperParamKey())
	if paramBytes == nil {
		return nil, ErrDeveloperParamNotFound()
	}
	param := new(DeveloperParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrFailedToUnmarshalDeveloperParam(err)
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
		return nil, ErrFailedToUnmarshalVoteParam(err)
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
		return nil, ErrFailedToUnmarshalProposalParam(err)
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
		return nil, ErrFailedToUnmarshalValidatorParam(err)
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
		return nil, ErrFailedToUnmarshalCoinDayParam(err)
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
		return nil, ErrFailedToUnmarshalBandwidthParam(err)
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
		return nil, ErrFailedToUnmarshalAccountParam(err)
	}
	return param, nil
}

func (ph ParamHolder) setValidatorParam(ctx sdk.Context, param *ValidatorParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalValidatorParam(err)
	}
	store.Set(GetValidatorParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setEvaluateOfContentValueParam(
	ctx sdk.Context, para *EvaluateOfContentValueParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paraBytes, err := ph.cdc.MarshalJSON(*para)
	if err != nil {
		return ErrFailedToMarshalEvaluateOfContentValueParam(err)
	}
	store.Set(GetEvaluateOfContentValueParamKey(), paraBytes)
	return nil
}

func (ph ParamHolder) setGlobalAllocationParam(
	ctx sdk.Context, allocation *GlobalAllocationParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	allocationBytes, err := ph.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrFailedToMarshalGlobalAllocationParam(err)
	}
	store.Set(GetAllocationParamKey(), allocationBytes)
	return nil
}

func (ph ParamHolder) setInfraInternalAllocationParam(
	ctx sdk.Context, allocation *InfraInternalAllocationParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	allocationBytes, err := ph.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrFailedToMarshalInfraInternalAllocationParam(err)
	}
	store.Set(GetInfraInternalAllocationParamKey(), allocationBytes)
	return nil
}

func (ph ParamHolder) setPostParam(
	ctx sdk.Context, para *PostParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paraBytes, err := ph.cdc.MarshalJSON(*para)
	if err != nil {
		return ErrFailedToMarshalPostParam(err)
	}
	store.Set(GetPostParamKey(), paraBytes)
	return nil
}

func (ph ParamHolder) setDeveloperParam(ctx sdk.Context, param *DeveloperParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalDeveloperParam(err)
	}
	store.Set(GetDeveloperParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setVoteParam(ctx sdk.Context, param *VoteParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalVoteParam(err)
	}
	store.Set(GetVoteParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setProposalParam(ctx sdk.Context, param *ProposalParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalProposalParam(err)
	}
	store.Set(GetProposalParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setCoinDayParam(ctx sdk.Context, param *CoinDayParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalCoinDayParam(err)
	}
	store.Set(GetCoinDayParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setBandwidthParam(ctx sdk.Context, param *BandwidthParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	bandwidthBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalBandwidthParam(err)
	}
	store.Set(GetBandwidthParamKey(), bandwidthBytes)
	return nil
}

func (ph ParamHolder) setAccountParam(ctx sdk.Context, param *AccountParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	accountBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalAccountParam(err)
	}
	store.Set(GetAccountParamKey(), accountBytes)
	return nil
}

func GetPostParamKey() []byte {
	return postParamSubStore
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
