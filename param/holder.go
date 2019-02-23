package param

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
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
	postParamSubStore                    = []byte{0x0a} // Substore for evaluate of content value
	reputationParamSubStore              = []byte{0x0b} // Substore for reputation parameters

	// AnnualInflationCeiling - annual inflation upper bound
	AnnualInflationCeiling = types.NewDecFromRat(98, 1000)
	// AnnualInflationFloor - annual inflation lower bound
	AnnualInflationFloor = types.NewDecFromRat(3, 100)
)

// ParamHolder - parameter KVStore
type ParamHolder struct {
	// The (unexposed) key used to access the store from the Context
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewParamHolder - create a new parameter KVStore
func NewParamHolder(key sdk.StoreKey) ParamHolder {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return ParamHolder{
		key: key,
		cdc: cdc,
	}
}

// InitParam - init all parameters based on code
func (ph ParamHolder) InitParam(ctx sdk.Context) error {
	globalAllocationParam := &GlobalAllocationParam{
		GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
		InfraAllocation:          types.NewDecFromRat(20, 100),
		ContentCreatorAllocation: types.NewDecFromRat(65, 100),
		DeveloperAllocation:      types.NewDecFromRat(10, 100),
		ValidatorAllocation:      types.NewDecFromRat(5, 100),
	}
	if err := ph.setGlobalAllocationParam(ctx, globalAllocationParam); err != nil {
		return err
	}

	infraInternalAllocationParam := &InfraInternalAllocationParam{
		StorageAllocation: types.NewDecFromRat(50, 100),
		CDNAllocation:     types.NewDecFromRat(50, 100),
	}
	if err := ph.setInfraInternalAllocationParam(ctx, infraInternalAllocationParam); err != nil {
		return err
	}

	postParam := &PostParam{
		ReportOrUpvoteIntervalSec: 24 * 3600,
		PostIntervalSec:           600,
		MaxReportReputation:       types.NewCoinFromInt64(100 * types.Decimals),
	}
	if err := ph.setPostParam(ctx, postParam); err != nil {
		return err
	}

	developerParam := &DeveloperParam{
		DeveloperMinDeposit:            types.NewCoinFromInt64(1000000 * types.Decimals),
		DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DeveloperCoinReturnTimes:       int64(7),
	}
	if err := ph.setDeveloperParam(ctx, developerParam); err != nil {
		return err
	}

	validatorParam := &ValidatorParam{
		ValidatorMinWithdraw:           types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:      types.NewCoinFromInt64(300000 * types.Decimals),
		ValidatorMinCommittingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissVote:                types.NewCoinFromInt64(20000 * types.Decimals),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000000 * types.Decimals),
		ValidatorListSize:              int64(21),
		AbsentCommitLimitation:         int64(600), // 30min
	}
	if err := ph.setValidatorParam(ctx, validatorParam); err != nil {
		return err
	}

	voteParam := &VoteParam{
		MinStakeIn:                     types.NewCoinFromInt64(1000 * types.Decimals),
		VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
		VoterCoinReturnTimes:           int64(7),
		DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DelegatorCoinReturnTimes:       int64(7),
	}
	if err := ph.setVoteParam(ctx, voteParam); err != nil {
		return err
	}

	proposalParam := &ProposalParam{
		ContentCensorshipDecideSec:  int64(7 * 24 * 3600),
		ContentCensorshipPassRatio:  types.NewDecFromRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamDecideSec:    int64(7 * 24 * 3600),
		ChangeParamPassRatio:    types.NewDecFromRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(7 * 24 * 3600),
		ProtocolUpgradePassRatio:  types.NewDecFromRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}
	if err := ph.setProposalParam(ctx, proposalParam); err != nil {
		return err
	}

	coinDayParam := &CoinDayParam{
		SecondsToRecoverCoinDay: int64(7 * 24 * 3600),
	}
	if err := ph.setCoinDayParam(ctx, coinDayParam); err != nil {
		return err
	}

	bandwidthParam := &BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
	}
	if err := ph.setBandwidthParam(ctx, bandwidthParam); err != nil {
		return err
	}

	accountParam := &AccountParam{
		MinimumBalance:               types.NewCoinFromInt64(0),
		RegisterFee:                  types.NewCoinFromInt64(1 * types.Decimals),
		FirstDepositFullCoinDayLimit: types.NewCoinFromInt64(1 * types.Decimals),
		MaxNumFrozenMoney:            10,
	}
	if err := ph.setAccountParam(ctx, accountParam); err != nil {
		return err
	}

	reputationParam := &ReputationParam{
		BestContentIndexN: 10,
	}
	if err := ph.setReputationParam(ctx, reputationParam); err != nil {
		return err
	}

	return nil
}

// InitParamFromConfig - init all parameters based on pass in args
func (ph ParamHolder) InitParamFromConfig(
	ctx sdk.Context,
	globalParam GlobalAllocationParam,
	infraInternalParam InfraInternalAllocationParam,
	postParam PostParam,
	developerParam DeveloperParam,
	validatorParam ValidatorParam,
	voteParam VoteParam,
	proposalParam ProposalParam,
	coinDayParam CoinDayParam,
	bandwidthParam BandwidthParam,
	accParam AccountParam,
	repParam ReputationParam) error {
	if err := ph.setGlobalAllocationParam(ctx, &globalParam); err != nil {
		return err
	}

	if err := ph.setInfraInternalAllocationParam(ctx, &infraInternalParam); err != nil {
		return err
	}

	if err := ph.setPostParam(ctx, &postParam); err != nil {
		return err
	}

	if err := ph.setDeveloperParam(ctx, &developerParam); err != nil {
		return err
	}

	if err := ph.setValidatorParam(ctx, &validatorParam); err != nil {
		return err
	}
	if err := ph.setVoteParam(ctx, &voteParam); err != nil {
		return err
	}
	if err := ph.setProposalParam(ctx, &proposalParam); err != nil {
		return err
	}
	if err := ph.setCoinDayParam(ctx, &coinDayParam); err != nil {
		return err
	}

	if err := ph.setBandwidthParam(ctx, &bandwidthParam); err != nil {
		return err
	}

	if err := ph.setAccountParam(ctx, &accParam); err != nil {
		return err
	}

	if err := ph.setReputationParam(ctx, &repParam); err != nil {
		return err
	}

	return nil
}

// GetGlobalAllocationParam - get global allocation param
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

// GetInfraInternalAllocationParam - get infra internal allocation allocation param
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

// GetPostParam - get post param
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

// GetDeveloperParam - get developer param
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

// GetVoteParam - get vote param
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

// GetProposalParam - get proposal param
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

// GetValidatorParam - get validator param
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

// GetCoinDayParam - get coin day param
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

// GetBandwidthParam - get bandwidth param
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

// GetAccountParam - get account param
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

// GetReputationParam - get reputation param
func (ph ParamHolder) GetReputationParam(ctx sdk.Context) (*ReputationParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetReputationParamKey())
	if paramBytes == nil {
		return nil, ErrReputationParamNotFound()
	}
	param := new(ReputationParam)
	if err := ph.cdc.UnmarshalJSON(paramBytes, param); err != nil {
		return nil, ErrFailedToUnmarshalReputationParam(err)
	}
	return param, nil
}

// UpdateGlobalGrowthRate - update global growth rate
func (ph ParamHolder) UpdateGlobalGrowthRate(ctx sdk.Context, growthRate sdk.Dec) sdk.Error {
	store := ctx.KVStore(ph.key)
	allocationBytes := store.Get(GetAllocationParamKey())
	if allocationBytes == nil {
		return ErrGlobalAllocationParamNotFound()
	}
	allocation := new(GlobalAllocationParam)
	if err := ph.cdc.UnmarshalJSON(allocationBytes, allocation); err != nil {
		return ErrFailedToUnmarshalGlobalAllocationParam(err)
	}

	if growthRate.GT(AnnualInflationCeiling) {
		growthRate = AnnualInflationCeiling
	} else if growthRate.LT(AnnualInflationFloor) {
		growthRate = AnnualInflationFloor
	}
	allocation.GlobalGrowthRate = growthRate
	allocationBytes, err := ph.cdc.MarshalJSON(*allocation)
	if err != nil {
		return ErrFailedToMarshalGlobalAllocationParam(err)
	}
	store.Set(GetAllocationParamKey(), allocationBytes)
	return nil
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

func (ph ParamHolder) setReputationParam(ctx sdk.Context, param *ReputationParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	reputationBytes, err := ph.cdc.MarshalJSON(*param)
	if err != nil {
		return ErrFailedToMarshalReputationParam(err)
	}
	store.Set(GetReputationParamKey(), reputationBytes)
	return nil
}

// GetPostParamKey - "post param substore"
func GetPostParamKey() []byte {
	return postParamSubStore
}

// GetEvaluateOfContentValueParamKey - "evaluate of content value param substore"
func GetEvaluateOfContentValueParamKey() []byte {
	return evaluateOfContentValueParamSubStore
}

// GetAllocationParamKey - "allocation param substore"
func GetAllocationParamKey() []byte {
	return allocationParamSubStore
}

// GetInfraInternalAllocationParamKey - "infra internal allocation param substore"
func GetInfraInternalAllocationParamKey() []byte {
	return infraInternalAllocationParamSubStore
}

// GetDeveloperParamKey - "developer param substore"
func GetDeveloperParamKey() []byte {
	return developerParamSubStore
}

// GetVoteParamKey - "vote param substore"
func GetVoteParamKey() []byte {
	return voteParamSubStore
}

// GetValidatorParamKey - "validator param substore"
func GetValidatorParamKey() []byte {
	return validatorParamSubStore
}

// GetProposalParamKey - "proposal param substore"
func GetProposalParamKey() []byte {
	return proposalParamSubStore
}

// GetCoinDayParamKey - "coin day param substore"
func GetCoinDayParamKey() []byte {
	return coinDayParamSubStore
}

// GetBandwidthParamKey - "bandwidth param substore"
func GetBandwidthParamKey() []byte {
	return bandwidthParamSubStore
}

// GetAccountParamKey - "account param substore"
func GetAccountParamKey() []byte {
	return accountParamSubstore
}

// GetAccountParamKey - "account param substore"
func GetReputationParamKey() []byte {
	return reputationParamSubStore
}
