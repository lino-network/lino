package param

import (
	"fmt"
	"time"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

var (
	allocationParamSubStore             = []byte{0x00} // SubStore for allocation
	evaluateOfContentValueParamSubStore = []byte{0x02} // Substore for evaluate of content value
	developerParamSubStore              = []byte{0x03} // Substore for developer param
	voteParamSubStore                   = []byte{0x04} // Substore for vote param
	proposalParamSubStore               = []byte{0x05} // Substore for proposal param
	validatorParamSubStore              = []byte{0x06} // Substore for validator param
	bandwidthParamSubStore              = []byte{0x08} // Substore for bandwidth param
	accountParamSubstore                = []byte{0x09} // Substore for account param
	postParamSubStore                   = []byte{0x0a} // Substore for evaluate of content value
	reputationParamSubStore             = []byte{0x0b} // Substore for reputation parameters
	priceParamSubStore                  = []byte{0x0c} // Substore for price parameters

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
		ContentCreatorAllocation: types.NewDecFromRat(85, 100),
		DeveloperAllocation:      types.NewDecFromRat(10, 100),
		ValidatorAllocation:      types.NewDecFromRat(5, 100),
	}
	if err := ph.setGlobalAllocationParam(ctx, globalAllocationParam); err != nil {
		return err
	}

	postParam := &PostParam{}
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
		ValidatorMinDeposit:            types.NewCoinFromInt64(200000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000 * types.Decimals),
		AbsentCommitLimitation:         int64(600), // 30min
		OncallSize:                     int64(22),
		StandbySize:                    int64(7),
		ValidatorRevokePendingSec:      int64(7 * 24 * 3600),
		OncallInflationWeight:          int64(2),
		StandbyInflationWeight:         int64(1),
		MaxVotedValidators:             int64(3),
		SlashLimitation:                int64(5),
	}
	if err := ph.setValidatorParam(ctx, validatorParam); err != nil {
		return err
	}

	voteParam := &VoteParam{
		MinStakeIn:                 types.NewCoinFromInt64(1000 * types.Decimals),
		VoterCoinReturnIntervalSec: int64(7 * 24 * 3600),
		VoterCoinReturnTimes:       int64(7),
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

	bandwidthParam := &BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
		GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
		GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
		AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
		AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
		ExpectedMaxMPS:              types.NewDecFromRat(300, 1),
		MsgFeeFactorA:               types.NewDecFromRat(6, 1),
		MsgFeeFactorB:               types.NewDecFromRat(10, 1),
		MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
		AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
		AppVacancyFactor:            types.NewDecFromRat(69, 100),
		AppPunishmentFactor:         types.NewDecFromRat(14, 5),
	}
	if err := ph.setBandwidthParam(ctx, bandwidthParam); err != nil {
		return err
	}

	accountParam := &AccountParam{
		MinimumBalance: types.NewCoinFromInt64(0),
		RegisterFee:    types.NewCoinFromInt64(1 * types.Decimals),
	}
	if err := ph.setAccountParam(ctx, accountParam); err != nil {
		return err
	}

	reputationParam := &ReputationParam{
		BestContentIndexN: 200,
		UserMaxN:          50,
	}
	if err := ph.setReputationParam(ctx, reputationParam); err != nil {
		return err
	}

	priceParam := &PriceParam{
		TestnetMode:     true,
		UpdateEverySec:  int64(time.Hour.Seconds()),
		FeedEverySec:    int64((10 * time.Minute).Seconds()),
		HistoryMaxLen:   71,
		PenaltyMissFeed: types.NewCoinFromInt64(10000 * types.Decimals),
	}
	ph.setPriceParam(ctx, priceParam)

	return nil
}

// InitParamFromConfig - init all parameters based on pass in args
func (ph ParamHolder) InitParamFromConfig(
	ctx sdk.Context,
	globalParam GlobalAllocationParam,
	postParam PostParam,
	developerParam DeveloperParam,
	validatorParam ValidatorParam,
	voteParam VoteParam,
	proposalParam ProposalParam,
	bandwidthParam BandwidthParam,
	accParam AccountParam,
	repParam ReputationParam,
	priceParam PriceParam) error {
	if !globalParam.IsValid() {
		return fmt.Errorf("invalid global allocation param: %+v", globalParam)
	}
	if err := ph.setGlobalAllocationParam(ctx, &globalParam); err != nil {
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
	if err := ph.setBandwidthParam(ctx, &bandwidthParam); err != nil {
		return err
	}

	if err := ph.setAccountParam(ctx, &accParam); err != nil {
		return err
	}

	if err := ph.setReputationParam(ctx, &repParam); err != nil {
		return err
	}

	ph.setPriceParam(ctx, &priceParam)
	return nil
}

// GetGlobalAllocationParam - get global allocation param
func (ph ParamHolder) GetGlobalAllocationParam(ctx sdk.Context) *GlobalAllocationParam {
	store := ctx.KVStore(ph.key)
	allocationBytes := store.Get(GetAllocationParamKey())
	if allocationBytes == nil {
		panic("Global Allocation Param Not Initialized")
	}
	allocation := new(GlobalAllocationParam)
	ph.cdc.MustUnmarshalBinaryLengthPrefixed(allocationBytes, allocation)
	return allocation
}

// GetPostParam - get post param
func (ph ParamHolder) GetPostParam(ctx sdk.Context) (*PostParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetPostParamKey())
	if paramBytes == nil {
		return nil, ErrPostParamNotFound()
	}
	param := new(PostParam)
	if err := ph.cdc.UnmarshalBinaryLengthPrefixed(paramBytes, param); err != nil {
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
	if err := ph.cdc.UnmarshalBinaryLengthPrefixed(paramBytes, param); err != nil {
		return nil, ErrFailedToUnmarshalDeveloperParam(err)
	}
	return param, nil
}

// GetVoteParam - get vote param
func (ph ParamHolder) GetVoteParam(ctx sdk.Context) *VoteParam {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetVoteParamKey())
	if paramBytes == nil {
		panic("Vote Param Not Initialized")
	}
	param := new(VoteParam)
	ph.cdc.MustUnmarshalBinaryLengthPrefixed(paramBytes, param)
	return param
}

// GetProposalParam - get proposal param
func (ph ParamHolder) GetProposalParam(ctx sdk.Context) (*ProposalParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetProposalParamKey())
	if paramBytes == nil {
		return nil, ErrProposalParamNotFound()
	}
	param := new(ProposalParam)
	if err := ph.cdc.UnmarshalBinaryLengthPrefixed(paramBytes, param); err != nil {
		return nil, ErrFailedToUnmarshalProposalParam(err)
	}
	return param, nil
}

// GetValidatorParam - get validator param
func (ph ParamHolder) GetValidatorParam(ctx sdk.Context) *ValidatorParam {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetValidatorParamKey())
	if paramBytes == nil {
		panic("Validator Param Not FOund")
	}
	param := new(ValidatorParam)
	ph.cdc.MustUnmarshalBinaryLengthPrefixed(paramBytes, param)
	return param
}

// GetBandwidthParam - get bandwidth param
func (ph ParamHolder) GetBandwidthParam(ctx sdk.Context) (*BandwidthParam, sdk.Error) {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetBandwidthParamKey())
	if paramBytes == nil {
		return nil, ErrBandwidthParamNotFound()
	}
	param := new(BandwidthParam)
	if err := ph.cdc.UnmarshalBinaryLengthPrefixed(paramBytes, param); err != nil {
		return nil, ErrFailedToUnmarshalBandwidthParam(err)
	}
	return param, nil
}

// GetAccountParam - get account param
func (ph ParamHolder) GetAccountParam(ctx sdk.Context) *AccountParam {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetAccountParamKey())
	if paramBytes == nil {
		panic("Account Param Not Initialized")
	}
	param := new(AccountParam)
	ph.cdc.MustUnmarshalBinaryLengthPrefixed(paramBytes, param)
	return param
}

// GetReputationParam - get reputation param
func (ph ParamHolder) GetReputationParam(ctx sdk.Context) *ReputationParam {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetReputationParamKey())
	if paramBytes == nil {
		panic("Reputation Param Not Initialized")
	}
	param := new(ReputationParam)
	ph.cdc.MustUnmarshalBinaryLengthPrefixed(paramBytes, param)
	return param
}

// GetPriceParam - get price param
func (ph ParamHolder) GetPriceParam(ctx sdk.Context) *PriceParam {
	store := ctx.KVStore(ph.key)
	paramBytes := store.Get(GetPriceParamKey())
	if paramBytes == nil {
		panic("Price Param Not Initialized")
	}
	param := new(PriceParam)
	ph.cdc.MustUnmarshalBinaryLengthPrefixed(paramBytes, param)
	return param
}

// UpdateGlobalGrowthRate - update global growth rate
func (ph ParamHolder) UpdateGlobalGrowthRate(ctx sdk.Context, growthRate sdk.Dec) sdk.Error {
	store := ctx.KVStore(ph.key)
	allocationBytes := store.Get(GetAllocationParamKey())
	if allocationBytes == nil {
		return ErrGlobalAllocationParamNotFound()
	}
	allocation := new(GlobalAllocationParam)
	if err := ph.cdc.UnmarshalBinaryLengthPrefixed(allocationBytes, allocation); err != nil {
		return ErrFailedToUnmarshalGlobalAllocationParam(err)
	}

	if growthRate.GT(AnnualInflationCeiling) {
		growthRate = AnnualInflationCeiling
	} else if growthRate.LT(AnnualInflationFloor) {
		growthRate = AnnualInflationFloor
	}
	allocation.GlobalGrowthRate = growthRate
	allocationBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*allocation)
	if err != nil {
		return ErrFailedToMarshalGlobalAllocationParam(err)
	}
	store.Set(GetAllocationParamKey(), allocationBytes)
	return nil
}

func (ph ParamHolder) setValidatorParam(ctx sdk.Context, param *ValidatorParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*param)
	if err != nil {
		return ErrFailedToMarshalValidatorParam(err)
	}
	store.Set(GetValidatorParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setGlobalAllocationParam(
	ctx sdk.Context, allocation *GlobalAllocationParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	allocationBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*allocation)
	if err != nil {
		return ErrFailedToMarshalGlobalAllocationParam(err)
	}
	store.Set(GetAllocationParamKey(), allocationBytes)
	return nil
}

func (ph ParamHolder) setPostParam(
	ctx sdk.Context, para *PostParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paraBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*para)
	if err != nil {
		return ErrFailedToMarshalPostParam(err)
	}
	store.Set(GetPostParamKey(), paraBytes)
	return nil
}

func (ph ParamHolder) setDeveloperParam(ctx sdk.Context, param *DeveloperParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*param)
	if err != nil {
		return ErrFailedToMarshalDeveloperParam(err)
	}
	store.Set(GetDeveloperParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setVoteParam(ctx sdk.Context, param *VoteParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*param)
	if err != nil {
		return ErrFailedToMarshalVoteParam(err)
	}
	store.Set(GetVoteParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setProposalParam(ctx sdk.Context, param *ProposalParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	paramBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*param)
	if err != nil {
		return ErrFailedToMarshalProposalParam(err)
	}
	store.Set(GetProposalParamKey(), paramBytes)
	return nil
}

func (ph ParamHolder) setBandwidthParam(ctx sdk.Context, param *BandwidthParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	bandwidthBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*param)
	if err != nil {
		return ErrFailedToMarshalBandwidthParam(err)
	}
	store.Set(GetBandwidthParamKey(), bandwidthBytes)
	return nil
}

func (ph ParamHolder) setAccountParam(ctx sdk.Context, param *AccountParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	accountBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*param)
	if err != nil {
		return ErrFailedToMarshalAccountParam(err)
	}
	store.Set(GetAccountParamKey(), accountBytes)
	return nil
}

func (ph ParamHolder) setReputationParam(ctx sdk.Context, param *ReputationParam) sdk.Error {
	store := ctx.KVStore(ph.key)
	reputationBytes, err := ph.cdc.MarshalBinaryLengthPrefixed(*param)
	if err != nil {
		return ErrFailedToMarshalReputationParam(err)
	}
	store.Set(GetReputationParamKey(), reputationBytes)
	return nil
}

func (ph ParamHolder) setPriceParam(ctx sdk.Context, param *PriceParam) {
	store := ctx.KVStore(ph.key)
	bytes := ph.cdc.MustMarshalBinaryLengthPrefixed(*param)
	store.Set(GetPriceParamKey(), bytes)
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

// GetPriceParamKey
func GetPriceParamKey() []byte {
	return priceParamSubStore
}
