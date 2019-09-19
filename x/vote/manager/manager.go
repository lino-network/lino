package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/lino-network/lino/x/vote/types"
)

// VoteManager - vote manager
type VoteManager struct {
	am          acc.AccountKeeper
	storage     model.VoteStorage
	paramHolder param.ParamHolder
	gm          global.GlobalKeeper
}

// NewVoteManager - new vote manager
func NewVoteManager(key sdk.StoreKey, holder param.ParamHolder, am acc.AccountKeeper, gm global.GlobalKeeper) VoteManager {
	return VoteManager{
		am:          am,
		storage:     model.NewVoteStorage(key),
		paramHolder: holder,
		gm:          gm,
	}
}

// InitGenesis - initialize KV Store
func (vm VoteManager) InitGenesis(ctx sdk.Context) error {
	if err := vm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

// DoesVoterExist - check if voter exist or not
func (vm VoteManager) DoesVoterExist(ctx sdk.Context, username linotypes.AccountKey) bool {
	return vm.storage.DoesVoterExist(ctx, username)
}

func (vm VoteManager) StakeIn(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err
	}

	if param.MinStakeIn.IsGT(amount) {
		return types.ErrInsufficientDeposit()
	}

	// withdraw money from voter's bank
	if err := vm.am.MinusCoinFromUsername(ctx, username, amount); err != nil {
		return err
	}

	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		if err.Code() != linotypes.CodeVoterNotFound {
			return err
		}
		voter = &model.Voter{
			Username:          username,
			LinoStake:         linotypes.NewCoinFromInt64(0),
			LastPowerChangeAt: ctx.BlockHeader().Time.Unix(),
			Duty:              types.DutyVoter,
			Interest:          linotypes.NewCoinFromInt64(0),
		}
	}

	interest, err := vm.gm.GetInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}
	voter.Interest = voter.Interest.Plus(interest)
	voter.LinoStake = voter.LinoStake.Plus(amount)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()

	return vm.storage.SetVoter(ctx, username, voter)
}

func (vm VoteManager) StakeOut(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}

	if !voter.LinoStake.IsGTE(amount) {
		return types.ErrInsufficientStake()
	}

	interest, err := vm.gm.GetInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}
	voter.Interest = voter.Interest.Plus(interest)
	voter.LinoStake = voter.LinoStake.Minus(amount)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()

	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err
	}

	if err := vm.am.AddFrozenMoney(
		ctx, username, amount, ctx.BlockHeader().Time.Unix(),
		param.VoterCoinReturnIntervalSec, param.VoterCoinReturnTimes); err != nil {
		return err
	}

	events, err := accmn.CreateCoinReturnEvents(
		ctx, username, param.VoterCoinReturnTimes, param.VoterCoinReturnIntervalSec, amount, linotypes.VoteReturnCoin)
	if err != nil {
		return err
	}

	if err := vm.gm.RegisterCoinReturnEvent(
		ctx, events, param.VoterCoinReturnTimes, param.VoterCoinReturnIntervalSec); err != nil {
		return err
	}

	return nil
}

// ClaimInterest - add lino power interst to user balance
func (vm VoteManager) ClaimInterest(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}

	interest, err := vm.gm.GetInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}

	if err := vm.am.AddCoinToUsername(ctx, username, voter.Interest.Plus(interest)); err != nil {
		return err
	}
	voter.Interest = linotypes.NewCoinFromInt64(0)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()
	return vm.storage.SetVoter(ctx, username, voter)
}

func (vm VoteManager) GetVoterDuty(ctx sdk.Context, username linotypes.AccountKey) (types.VoterDuty, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return types.DutyVoter, err
	}
	return voter.Duty, nil
}

// AssignDuty froze some amount of stake and assign a duty to user.
func (vm VoteManager) AssignDuty(
	ctx sdk.Context, username linotypes.AccountKey, duty types.VoterDuty, frozenAmount linotypes.Coin) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	if voter.Duty != types.DutyVoter {
		return types.ErrNotAVoterOrHasDuty()
	}

	if !voter.LinoStake.IsGTE(frozenAmount) {
		return types.ErrInsufficientStake()
	}

	interest, err := vm.gm.GetInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}

	voter.Interest = voter.Interest.Plus(interest)
	voter.Duty = duty
	voter.FrozenAmount = frozenAmount
	// voter.LinoStake = voter.LinoStake.Minus(frozenAmount)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()

	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	return nil
}

// Export storage state.
func (vm VoteManager) Export(ctx sdk.Context) *model.VoterTables {
	return vm.storage.Export(ctx)
}

// Import storage state.
func (vm VoteManager) Import(ctx sdk.Context, voter *model.VoterTablesIR) {
	vm.storage.Import(ctx, voter)
}
