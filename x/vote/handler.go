package vote

import (
	"fmt"
	"reflect"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
)

// NewHandler - Handle all "vote" type messages.
func NewHandler(vm VoteManager, am acc.AccountKeeper, gm *global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case StakeInMsg:
			return handleStakeInMsg(ctx, vm, gm, am, msg)
		case StakeOutMsg:
			return handleStakeOutMsg(ctx, vm, gm, am, msg)
		case DelegateMsg:
			return handleDelegateMsg(ctx, vm, gm, am, msg)
		case DelegatorWithdrawMsg:
			return handleDelegatorWithdrawMsg(ctx, vm, gm, am, msg)
		case ClaimInterestMsg:
			return handleClaimInterestMsg(ctx, vm, gm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized vote msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleStakeInMsg(
	ctx sdk.Context, vm VoteManager, gm *global.GlobalManager,
	am acc.AccountKeeper, msg StakeInMsg) sdk.Result {
	// Must have an normal acount
	// if !am.DoesAccountExist(ctx, msg.Username) {
	// 	return ErrAccountNotFound().Result()
	// }

	coin, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}

	if param.MinStakeIn.IsGT(coin) {
		return ErrInsufficientDeposit().Result()
	}

	// withdraw money from voter's bank
	if err := am.MinusCoinFromUsername(ctx, msg.Username, coin); err != nil {
		return err.Result()
	}

	if err := AddStake(ctx, msg.Username, coin, vm, gm, am); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleStakeOutMsg(
	ctx sdk.Context, vm VoteManager, gm *global.GlobalManager,
	am acc.AccountKeeper, msg StakeOutMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	if !vm.IsLegalVoterWithdraw(ctx, msg.Username, coin) {
		return ErrIllegalWithdraw().Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := MinusStake(ctx, msg.Username, coin, vm, gm, am); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Username, gm, am, param.VoterCoinReturnTimes,
		param.VoterCoinReturnIntervalSec, coin, types.VoteReturnCoin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDelegateMsg(
	ctx sdk.Context, vm VoteManager, gm *global.GlobalManager, am acc.AccountKeeper, msg DelegateMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}

	if param.MinStakeIn.IsGT(coin) {
		return ErrInsufficientDeposit().Result()
	}

	// withdraw money from delegator's bank
	if err := am.MinusCoinFromUsername(ctx, msg.Delegator, coin); err != nil {
		return err.Result()
	}

	if err := AddStake(ctx, msg.Delegator, coin, vm, gm, am); err != nil {
		return err.Result()
	}

	// add delegation relation
	if addErr := vm.AddDelegation(ctx, msg.Voter, msg.Delegator, coin); addErr != nil {
		return addErr.Result()
	}
	return sdk.Result{}
}

func handleDelegatorWithdrawMsg(
	ctx sdk.Context, vm VoteManager, gm *global.GlobalManager,
	am acc.AccountKeeper, msg DelegatorWithdrawMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if !vm.IsLegalDelegatorWithdraw(ctx, msg.Voter, msg.Delegator, coin) {
		return ErrIllegalWithdraw().Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}
	if err := MinusStake(ctx, msg.Delegator, coin, vm, gm, am); err != nil {
		return err.Result()
	}
	if err := vm.DelegatorWithdraw(ctx, msg.Voter, msg.Delegator, coin); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Delegator, gm, am, param.DelegatorCoinReturnTimes,
		param.DelegatorCoinReturnIntervalSec, coin, types.DelegationReturnCoin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleClaimInterestMsg(ctx sdk.Context, vm VoteManager, gm *global.GlobalManager, am acc.AccountKeeper, msg ClaimInterestMsg) sdk.Result {
	if err := calculateAndAddInterest(ctx, vm, gm, am, msg.Username); err != nil {
		return err.Result()
	}
	// claim interest
	interest, err := vm.ClaimInterest(ctx, msg.Username)
	if err != nil {
		return err.Result()
	}
	if err := am.AddCoinToUsername(ctx, msg.Username, interest); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func AddStake(
	ctx sdk.Context, username types.AccountKey, stake types.Coin, vm VoteManager,
	gm *global.GlobalManager, am acc.AccountKeeper) sdk.Error {
	// Register the user if this name has not been registered
	if !vm.DoesVoterExist(ctx, username) {
		if err := vm.AddVoter(ctx, username, types.NewCoinFromInt64(0)); err != nil {
			return err
		}
	}
	if err := calculateAndAddInterest(ctx, vm, gm, am, username); err != nil {
		return err
	}

	// add linoStake to voter account
	if err := vm.AddLinoStake(ctx, username, stake); err != nil {
		return err
	}

	// add linoStake to global stat
	if err := gm.AddLinoStakeToStat(ctx, stake); err != nil {
		return err
	}
	return nil
}

func MinusStake(
	ctx sdk.Context, username types.AccountKey, stake types.Coin, vm VoteManager,
	gm *global.GlobalManager, am acc.AccountKeeper) sdk.Error {
	if err := calculateAndAddInterest(ctx, vm, gm, am, username); err != nil {
		return err
	}

	// minus linoStake to voter account
	if err := vm.MinusLinoStake(ctx, username, stake); err != nil {
		return err
	}
	// add linoStake to global stat
	if err := gm.MinusLinoStakeFromStat(ctx, stake); err != nil {
		return err
	}
	return nil
}

func calculateAndAddInterest(ctx sdk.Context, vm VoteManager, gm *global.GlobalManager,
	am acc.AccountKeeper, name types.AccountKey) sdk.Error {
	userLinoStake, err := vm.GetLinoStake(ctx, name)
	if err != nil {
		return err
	}

	LSLastChangedAt, err := vm.GetLinoStakeLastChangedAt(ctx, name)
	if err != nil {
		return err
	}

	interest, err := gm.GetInterestSince(ctx, LSLastChangedAt, userLinoStake)
	if err != nil {
		return err
	}

	if err := vm.AddInterest(ctx, name, interest); err != nil {
		return err
	}

	if err := vm.SetLinoStakeLastChangedAt(ctx, name, ctx.BlockHeader().Time.Unix()); err != nil {
		return err
	}

	return nil
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm *global.GlobalManager, am acc.AccountKeeper,
	times int64, interval int64, coin types.Coin, returnType types.TransferDetailType) sdk.Error {

	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time.Unix(), interval, times); err != nil {
		return err
	}

	events, err := accmn.CreateCoinReturnEvents(ctx, name, times, interval, coin, returnType)
	if err != nil {
		return err
	}

	if err := gm.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
