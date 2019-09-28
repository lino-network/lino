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
	paramHolder param.ParamKeeper
	gm          global.GlobalKeeper
	hooks       StakingHooks
}

// NewVoteManager - new vote manager
func NewVoteManager(
	key sdk.StoreKey, holder param.ParamKeeper, am acc.AccountKeeper, gm global.GlobalKeeper) VoteManager {
	return VoteManager{
		am:          am,
		storage:     model.NewVoteStorage(key),
		paramHolder: holder,
		gm:          gm,
	}
}

// Set the validator hooks
func (vm *VoteManager) SetHooks(sh StakingHooks) *VoteManager {
	if vm.hooks != nil {
		panic("cannot set vote hooks twice")
	}
	vm.hooks = sh
	return vm
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
	return vm.AddStake(ctx, username, amount)
}

func (vm VoteManager) AddStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
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

	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	// add linoStake to global stat
	if err := vm.gm.AddLinoStakeToStat(ctx, amount); err != nil {
		return err
	}
	return vm.AfterAddingStake(ctx, username)
}

func (vm VoteManager) StakeOut(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	if err := vm.MinusStake(ctx, username, amount); err != nil {
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

func (vm VoteManager) MinusStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}

	// make sure stake is sufficient excludes frozen amount
	if !voter.LinoStake.Minus(voter.FrozenAmount).IsGTE(amount) {
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
	// add linoStake to global stat
	if err := vm.gm.MinusLinoStakeFromStat(ctx, amount); err != nil {
		return err
	}
	return vm.AfterSubtractingStake(ctx, username)
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

func (vm VoteManager) GetLinoStake(ctx sdk.Context, username linotypes.AccountKey) (linotypes.Coin, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}
	return voter.LinoStake, nil
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

	if voter.FrozenAmount.IsPositive() {
		return types.ErrFrozenAmountIsNotEmpty()
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

// UnassignDuty register unassign duty event with time after waitingPeriodSec seconds.
func (vm VoteManager) UnassignDuty(ctx sdk.Context, username linotypes.AccountKey, waitingPeriodSec int64) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	if voter.Duty == types.DutyVoter {
		return types.ErrNoDuty()
	}
	return vm.gm.RegisterEventAtTime(
		ctx, ctx.BlockHeader().Time.Unix()+waitingPeriodSec, types.UnassignDutyEvent{Username: username})
}

// SlashStake - slash as much as it can, regardless of frozen money
func (vm VoteManager) SlashStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) (slashedAmount linotypes.Coin, err sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}

	interest, err := vm.gm.GetInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}

	voter.Interest = voter.Interest.Plus(interest)
	if !voter.LinoStake.IsGTE(amount) {
		slashedAmount = voter.LinoStake
		voter.LinoStake = linotypes.NewCoinFromInt64(0)
	} else {
		slashedAmount = amount
		voter.LinoStake = voter.LinoStake.Minus(amount)
	}
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()

	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}
	if err := vm.AfterSlashing(ctx, username); err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}
	return slashedAmount, nil
}

// ExecUnassignDutyEvent - execute unassign duty events.
func (vm VoteManager) ExecUnassignDutyEvent(ctx sdk.Context, event types.UnassignDutyEvent) sdk.Error {
	// Check if it is voter or not
	voter, err := vm.storage.GetVoter(ctx, event.Username)
	if err != nil {
		return err
	}

	// set frozen amount to zero and duty to voter
	voter.FrozenAmount = linotypes.NewCoinFromInt64(0)
	voter.Duty = types.DutyVoter
	return vm.storage.SetVoter(ctx, event.Username, voter)
}

func (vm VoteManager) GetVoter(ctx sdk.Context, username linotypes.AccountKey) (*model.Voter, sdk.Error) {
	return vm.storage.GetVoter(ctx, username)
}

// Export storage state.
func (vm VoteManager) Export(ctx sdk.Context) *model.VoterTables {
	return vm.storage.Export(ctx)
}

// Import storage state.
func (vm VoteManager) Import(ctx sdk.Context, voter *model.VoterTablesIR) {
	vm.storage.Import(ctx, voter)
}
