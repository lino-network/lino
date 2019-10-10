package manager

import (
	"fmt"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/lino-network/lino/x/vote/types"
)

const (
	exportVersion = 1
	importVersion = 1
)

// VoteManager - vote manager
type VoteManager struct {
	storage model.VoteStorage

	// deps
	paramHolder param.ParamKeeper
	am          acc.AccountKeeper
	gm          global.GlobalKeeper

	// state
	hooks StakingHooks
}

// NewVoteManager - new vote manager
func NewVoteManager(key sdk.StoreKey, holder param.ParamKeeper, am acc.AccountKeeper, gm global.GlobalKeeper) VoteManager {
	return VoteManager{
		am:          am,
		storage:     model.NewVoteStorage(key),
		paramHolder: holder,
		gm:          gm,
	}
}

func (vm VoteManager) InitGenesis(ctx sdk.Context) {
	linoStakeStat := &model.LinoStakeStat{
		TotalConsumptionFriction: linotypes.NewCoinFromInt64(0),
		TotalLinoStake:           linotypes.NewCoinFromInt64(0),
		UnclaimedFriction:        linotypes.NewCoinFromInt64(0),
		UnclaimedLinoStake:       linotypes.NewCoinFromInt64(0),
	}
	vm.storage.SetLinoStakeStat(ctx, 0, linoStakeStat)
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
	if vm.paramHolder.GetVoteParam(ctx).MinStakeIn.IsGT(amount) {
		return types.ErrInsufficientDeposit()
	}

	err := vm.am.MoveToPool(
		ctx, linotypes.VoteStakeInPool, linotypes.NewAccOrAddrFromAcc(username), amount)
	if err != nil {
		return err
	}

	return vm.AddStake(ctx, username, amount)
}

func (vm VoteManager) StakeInFor(ctx sdk.Context, sender linotypes.AccountKey,
	receiver linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	if vm.paramHolder.GetVoteParam(ctx).MinStakeIn.IsGT(amount) {
		return types.ErrInsufficientDeposit()
	}

	// withdraw money from sender's bank and add stake to receiver
	err := vm.am.MoveToPool(
		ctx, linotypes.VoteStakeInPool, linotypes.NewAccOrAddrFromAcc(sender), amount)
	if err != nil {
		return err
	}

	return vm.AddStake(ctx, receiver, amount)
}

func (vm VoteManager) AddStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	if !vm.am.DoesAccountExist(ctx, username) {
		return types.ErrAccountNotFound()
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

	interest, err := vm.popInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}
	voter.Interest = voter.Interest.Plus(interest)
	voter.LinoStake = voter.LinoStake.Plus(amount)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()

	vm.storage.SetVoter(ctx, username, voter)
	// add linoStake to global stat
	if err := vm.updateLinoStakeStat(ctx, amount, true); err != nil {
		return err
	}
	return vm.AfterAddingStake(ctx, username)
}

func (vm VoteManager) StakeOut(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	// move stake to stake return pool.
	err := vm.am.MoveBetweenPools(
		ctx, linotypes.VoteStakeInPool, linotypes.VoteStakeReturnPool, amount)
	if err != nil {
		return err
	}

	// minus stake stats
	if err := vm.MinusStake(ctx, username, amount); err != nil {
		return err
	}

	// create coin return events to return coins from stake return pool.
	//// add frozen money for records.
	param := vm.paramHolder.GetVoteParam(ctx)
	if err := vm.am.AddFrozenMoney(
		ctx, username, amount, ctx.BlockHeader().Time.Unix(),
		param.VoterCoinReturnIntervalSec, param.VoterCoinReturnTimes); err != nil {
		return err
	}

	//// create and register the events.
	events := accmn.CreateCoinReturnEvents(
		username, ctx.BlockTime().Unix(),
		param.VoterCoinReturnIntervalSec, param.VoterCoinReturnTimes,
		amount, linotypes.VoteReturnCoin, linotypes.VoteStakeReturnPool)
	for _, event := range events {
		err := vm.gm.RegisterEventAtTime(ctx, event.At, event)
		if err != nil {
			return err
		}
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

	interest, err := vm.popInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}
	voter.Interest = voter.Interest.Plus(interest)
	voter.LinoStake = voter.LinoStake.Minus(amount)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()
	vm.storage.SetVoter(ctx, username, voter)

	// minus linoStake from global stat
	if err := vm.updateLinoStakeStat(ctx, amount, false); err != nil {
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

	interest, err := vm.popInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}

	if err := vm.am.MoveFromPool(ctx,
		linotypes.VoteFrictionPool, linotypes.NewAccOrAddrFromAcc(username),
		voter.Interest.Plus(interest)); err != nil {
		return err
	}

	voter.Interest = linotypes.NewCoinFromInt64(0)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()
	vm.storage.SetVoter(ctx, username, voter)
	return nil
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

	interest, err := vm.popInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
	if err != nil {
		return err
	}

	voter.Interest = voter.Interest.Plus(interest)
	voter.Duty = duty
	voter.FrozenAmount = frozenAmount
	// voter.LinoStake = voter.LinoStake.Minus(frozenAmount)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()

	vm.storage.SetVoter(ctx, username, voter)
	return nil
}

// UnassignDuty register unassign duty event with time after waitingPeriodSec seconds.
func (vm VoteManager) UnassignDuty(ctx sdk.Context, username linotypes.AccountKey, waitingPeriodSec int64) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	if voter.Duty == types.DutyVoter || voter.Duty == types.DutyPending {
		return types.ErrNoDuty()
	}
	if err := vm.gm.RegisterEventAtTime(
		ctx, ctx.BlockHeader().Time.Unix()+waitingPeriodSec, types.UnassignDutyEvent{Username: username}); err != nil {
		return err
	}
	voter.Duty = types.DutyPending
	vm.storage.SetVoter(ctx, username, voter)
	return nil
}

// SlashStake - slash as much as it can, regardless of frozen money
func (vm VoteManager) SlashStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) (slashedAmount linotypes.Coin, err sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}

	interest, err := vm.popInterestSince(ctx, voter.LastPowerChangeAt, voter.LinoStake)
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

	vm.storage.SetVoter(ctx, username, voter)
	// minus linoStake from global stat
	if err := vm.updateLinoStakeStat(ctx, slashedAmount, false); err != nil {
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
	vm.storage.SetVoter(ctx, event.Username, voter)
	return nil
}

// popInterestSince - pop interest from unix time till now (exclusive)
func (vm VoteManager) popInterestSince(ctx sdk.Context, unixTime int64, linoStake linotypes.Coin) (linotypes.Coin, sdk.Error) {
	startDay := vm.gm.GetPastDay(ctx, unixTime)
	endDay := vm.gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
	totalInterest := linotypes.NewCoinFromInt64(0)
	for day := startDay; day < endDay; day++ {
		linoStakeStat, err := vm.storage.GetLinoStakeStat(ctx, day)
		if err != nil {
			return linotypes.NewCoinFromInt64(0), err
		}
		if linoStakeStat.UnclaimedLinoStake.IsZero() || !linoStakeStat.UnclaimedLinoStake.IsGTE(linoStake) {
			continue
		}
		interest :=
			linotypes.DecToCoin(linoStakeStat.UnclaimedFriction.ToDec().Mul(
				linoStake.ToDec().Quo(linoStakeStat.UnclaimedLinoStake.ToDec())))
		totalInterest = totalInterest.Plus(interest)
		linoStakeStat.UnclaimedFriction = linoStakeStat.UnclaimedFriction.Minus(interest)
		linoStakeStat.UnclaimedLinoStake = linoStakeStat.UnclaimedLinoStake.Minus(linoStake)
		vm.storage.SetLinoStakeStat(ctx, day, linoStakeStat)
	}
	return totalInterest, nil
}

// updateLinoStakeStat - add/sub lino power to total lino power at current day
func (vm VoteManager) updateLinoStakeStat(ctx sdk.Context, linoStake linotypes.Coin, isAdd bool) sdk.Error {
	pastDay := vm.gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
	linoStakeStat, err := vm.storage.GetLinoStakeStat(ctx, pastDay)
	if err != nil {
		return err
	}
	if isAdd {
		linoStakeStat.TotalLinoStake = linoStakeStat.TotalLinoStake.Plus(linoStake)
		linoStakeStat.UnclaimedLinoStake = linoStakeStat.UnclaimedLinoStake.Plus(linoStake)
	} else {
		linoStakeStat.TotalLinoStake = linoStakeStat.TotalLinoStake.Minus(linoStake)
		linoStakeStat.UnclaimedLinoStake = linoStakeStat.UnclaimedLinoStake.Minus(linoStake)
	}
	vm.storage.SetLinoStakeStat(ctx, pastDay, linoStakeStat)
	return nil
}

func (vm VoteManager) RecordFriction(ctx sdk.Context, friction linotypes.Coin) sdk.Error {
	pastDay := vm.gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
	linoStakeStat, err := vm.storage.GetLinoStakeStat(ctx, pastDay)
	if err != nil {
		return err
	}
	linoStakeStat.TotalConsumptionFriction = linoStakeStat.TotalConsumptionFriction.Plus(friction)
	linoStakeStat.UnclaimedFriction = linoStakeStat.UnclaimedFriction.Plus(friction)
	vm.storage.SetLinoStakeStat(ctx, pastDay, linoStakeStat)
	return nil
}

// AdvanceLinoStakeStats - save consumption and lino power to LinoStakeStat of a new day.
// It need to be executed daily.
func (vm VoteManager) AdvanceLinoStakeStats(ctx sdk.Context) sdk.Error {
	nDay := vm.gm.GetPastDay(ctx, ctx.BlockTime().Unix())
	if nDay < 1 {
		return nil
	}
	lastLinoStakeStat, err := vm.storage.GetLinoStakeStat(ctx, nDay-1)
	if err != nil {
		return err
	}
	// If lino stake exist last day, the consumption will keep for lino stake holder that day
	if !lastLinoStakeStat.TotalLinoStake.IsZero() {
		lastLinoStakeStat.TotalConsumptionFriction = linotypes.NewCoinFromInt64(0)
		lastLinoStakeStat.UnclaimedFriction = linotypes.NewCoinFromInt64(0)
	}
	vm.storage.SetLinoStakeStat(ctx, nDay, lastLinoStakeStat)
	return nil
}

func (vm VoteManager) GetVoter(ctx sdk.Context, username linotypes.AccountKey) (*model.Voter, sdk.Error) {
	return vm.storage.GetVoter(ctx, username)
}

// Export storage state.
func (vm VoteManager) ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	state := &model.VoterTablesIR{
		Version: exportVersion,
	}
	storeMap := vm.storage.StoreMap(ctx)
	storeMap[string(model.VoterSubstore)].Iterate(func(key []byte, val interface{}) bool {
		voter := val.(*model.Voter)
		state.Voters = append(state.Voters, model.VoterIR(*voter))
		return false
	})
	return utils.Save(filepath, cdc, state)
}

// Import storage state.
func (vm VoteManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.VoterTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.VoterTablesIR)

	if table.Version != importVersion {
		return fmt.Errorf("unsupported import version: %d", table.Version)
	}

	for _, voterir := range table.Voters {
		voter := model.Voter(voterir)
		vm.storage.SetVoter(ctx, voter.Username, &voter)
	}
	return nil
}
