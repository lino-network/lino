package validator

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/global"
	vote "github.com/lino-network/lino/x/vote"
)

func NewHandler(am acc.AccountManager, valManager ValidatorManager, voteManager vote.VoteManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ValidatorDepositMsg:
			return handleDepositMsg(ctx, valManager, voteManager, am, msg)
		case ValidatorWithdrawMsg:
			return handleWithdrawMsg(ctx, valManager, gm, am, msg)
		case ValidatorRevokeMsg:
			return handleRevokeMsg(ctx, valManager, gm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleDepositMsg(
	ctx sdk.Context, valManager ValidatorManager, voteManager vote.VoteManager,
	am acc.AccountManager, msg ValidatorDepositMsg) sdk.Result {
	// Must have a normal acount
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound().Result()
	}

	coin, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	// withdraw money from validator's bank
	if err = am.MinusSavingCoin(ctx, msg.Username, coin, "", "", types.ValidatorDeposit); err != nil {
		return err.Result()
	}

	// Register the user if this name has not been registered
	if !valManager.DoesValidatorExist(ctx, msg.Username) {
		// check validator minimum voting deposit requirement
		if !voteManager.CanBecomeValidator(ctx, msg.Username) {
			return ErrInsufficientDeposit().Result()
		}
		if err := valManager.RegisterValidator(
			ctx, msg.Username, msg.ValPubKey, coin, msg.Link); err != nil {
			return err.Result()
		}
	} else {
		// Deposit coins
		if err := valManager.Deposit(ctx, msg.Username, coin, msg.Link); err != nil {
			return err.Result()
		}
	}

	// Deposit must be balanced
	votingDeposit, err := voteManager.GetVoterDeposit(ctx, msg.Username)
	if err != nil {
		return err.Result()
	}

	if !valManager.IsBalancedAccount(ctx, msg.Username, votingDeposit) {
		return ErrUnbalancedAccount().Result()
	}

	// Try to become oncall validator
	if err := valManager.TryBecomeOncallValidator(ctx, msg.Username); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle Withdraw Msg
func handleWithdrawMsg(
	ctx sdk.Context, vm ValidatorManager, gm global.GlobalManager, am acc.AccountManager,
	msg ValidatorWithdrawMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	if !vm.IsLegalWithdraw(ctx, msg.Username, coin) {
		return ErrIllegalWithdraw().Result()
	}

	if err := vm.ValidatorWithdraw(ctx, msg.Username, coin); err != nil {
		return err.Result()
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Username, gm, am, param.ValidatorCoinReturnTimes,
		param.ValidatorCoinReturnIntervalHr, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleRevokeMsg(
	ctx sdk.Context, vm ValidatorManager, gm global.GlobalManager, am acc.AccountManager,
	msg ValidatorRevokeMsg) sdk.Result {
	coin, withdrawErr := vm.ValidatorWithdrawAll(ctx, msg.Username)
	if withdrawErr != nil {
		return withdrawErr.Result()
	}

	if err := vm.RemoveValidatorFromAllLists(ctx, msg.Username); err != nil {
		return err.Result()
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Username, gm, am, param.ValidatorCoinReturnTimes, param.ValidatorCoinReturnIntervalHr, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm global.GlobalManager, am acc.AccountManager,
	times int64, interval int64, coin types.Coin) sdk.Error {

	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time.Unix(), interval, times); err != nil {
		return err
	}

	events, err := acc.CreateCoinReturnEvents(name, times, interval, coin, types.ValidatorReturnCoin)
	if err != nil {
		return err
	}

	if err := gm.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
