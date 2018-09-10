package vote

import (
	"fmt"
	"reflect"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
)

// NewHandler - Handle all "vote" type messages.
func NewHandler(vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case VoterDepositMsg:
			return handleVoterDepositMsg(ctx, vm, am, msg)
		case VoterWithdrawMsg:
			return handleVoterWithdrawMsg(ctx, vm, gm, am, msg)
		case VoterRevokeMsg:
			return handleVoterRevokeMsg(ctx, vm, gm, am, msg)
		case DelegateMsg:
			return handleDelegateMsg(ctx, vm, am, msg)
		case DelegatorWithdrawMsg:
			return handleDelegatorWithdrawMsg(ctx, vm, gm, am, msg)
		case RevokeDelegationMsg:
			return handleRevokeDelegationMsg(ctx, vm, gm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized vote msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleVoterDepositMsg(
	ctx sdk.Context, vm VoteManager, am acc.AccountManager, msg VoterDepositMsg) sdk.Result {
	// Must have an normal acount
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound().Result()
	}

	coin, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	// withdraw money from voter's bank
	if err := am.MinusSavingCoin(ctx, msg.Username, coin, "", "", types.VoterDeposit); err != nil {
		return err.Result()
	}

	// Register the user if this name has not been registered
	if !vm.DoesVoterExist(ctx, msg.Username) {
		if err := vm.AddVoter(ctx, msg.Username, coin); err != nil {
			return err.Result()
		}
		return sdk.Result{}
	}

	// Deposit coins
	if err := vm.AddLinoPower(ctx, msg.Username, coin); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleVoterWithdrawMsg(
	ctx sdk.Context, vm VoteManager, gm global.GlobalManager, am acc.AccountManager, msg VoterWithdrawMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	if !vm.IsLegalVoterWithdraw(ctx, msg.Username, coin) {
		return ErrIllegalWithdraw().Result()
	}

	if err := vm.VoterWithdraw(ctx, msg.Username, coin); err != nil {
		return err.Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}
	if err := returnCoinTo(
		ctx, msg.Username, gm, am, param.VoterCoinReturnTimes,
		param.VoterCoinReturnIntervalSec, coin, types.VoteReturnCoin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleVoterRevokeMsg(
	ctx sdk.Context, vm VoteManager, gm global.GlobalManager,
	am acc.AccountManager, msg VoterRevokeMsg) sdk.Result {
	// reject if this is a validator
	if vm.IsInValidatorList(ctx, msg.Username) {
		return ErrValidatorCannotRevoke().Result()
	}

	delegators, err := vm.GetAllDelegators(ctx, msg.Username)
	if err != nil {
		return err.Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}
	// return coins to all delegators
	for _, delegator := range delegators {
		coin, withdrawErr := vm.DelegatorWithdrawAll(ctx, msg.Username, delegator)
		if withdrawErr != nil {
			return withdrawErr.Result()
		}
		if err := returnCoinTo(
			ctx, delegator, gm, am, param.DelegatorCoinReturnTimes,
			param.DelegatorCoinReturnIntervalSec, coin, types.DelegationReturnCoin); err != nil {
			return err.Result()
		}
	}

	// return coins to voter
	coin, withdrawErr := vm.VoterWithdrawAll(ctx, msg.Username)
	if withdrawErr != nil {
		return withdrawErr.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Username, gm, am, param.VoterCoinReturnTimes,
		param.VoterCoinReturnIntervalSec, coin, types.VoteReturnCoin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDelegateMsg(ctx sdk.Context, vm VoteManager, am acc.AccountManager, msg DelegateMsg) sdk.Result {
	// Must have an normal acount
	if !am.DoesAccountExist(ctx, msg.Voter) {
		return ErrAccountNotFound().Result()
	}

	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	// withdraw money from delegator's bank
	if err := am.MinusSavingCoin(
		ctx, msg.Delegator, coin, msg.Voter, "", types.Delegate); err != nil {
		return err.Result()
	}
	// add delegation relation
	if addErr := vm.AddDelegation(ctx, msg.Voter, msg.Delegator, coin); addErr != nil {
		return addErr.Result()
	}
	return sdk.Result{}
}

func handleDelegatorWithdrawMsg(
	ctx sdk.Context, vm VoteManager, gm global.GlobalManager,
	am acc.AccountManager, msg DelegatorWithdrawMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if !vm.IsLegalDelegatorWithdraw(ctx, msg.Voter, msg.Delegator, coin) {
		return ErrIllegalWithdraw().Result()
	}

	if err := vm.DelegatorWithdraw(ctx, msg.Voter, msg.Delegator, coin); err != nil {
		return err.Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Delegator, gm, am, param.DelegatorCoinReturnTimes,
		param.DelegatorCoinReturnIntervalSec, coin, types.DelegationReturnCoin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleRevokeDelegationMsg(
	ctx sdk.Context, vm VoteManager, gm global.GlobalManager,
	am acc.AccountManager, msg RevokeDelegationMsg) sdk.Result {
	coin, withdrawErr := vm.DelegatorWithdrawAll(ctx, msg.Voter, msg.Delegator)
	if withdrawErr != nil {
		return withdrawErr.Result()
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Delegator, gm, am, param.DelegatorCoinReturnTimes,
		param.DelegatorCoinReturnIntervalSec, coin, types.DelegationReturnCoin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm global.GlobalManager, am acc.AccountManager,
	times int64, interval int64, coin types.Coin, returnType types.TransferDetailType) sdk.Error {

	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time.Unix(), interval, times); err != nil {
		return err
	}

	events, err := acc.CreateCoinReturnEvents(name, times, interval, coin, returnType)
	if err != nil {
		return err
	}

	if err := gm.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
