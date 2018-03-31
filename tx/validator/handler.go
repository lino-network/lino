package validator

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
)

func NewHandler(vm ValidatorManager, am acc.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ValidatorDepositMsg:
			return handleDepositMsg(ctx, vm, am, msg)
		case ValidatorWithdrawMsg:
			return handleWithdrawMsg(ctx, vm, am, msg)
		case ValidatorRevokeMsg:
			return handleRevokeMsg(ctx, vm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle DepositMsg
func handleDepositMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg ValidatorDepositMsg) sdk.Result {
	proxyAcc := acc.NewProxyAccount(msg.Username, &am)
	// Must have an normal acount
	if !proxyAcc.IsAccountExist(ctx) {
		return ErrUsernameNotFound().Result()
	}
	// withdraw money from validator's bank
	err := proxyAcc.MinusCoins(ctx, msg.Deposit)
	if err != nil {
		return err.Result()
	}

	var validator *Validator
	// This name has not been registered
	if !vm.IsValidatorExist(ctx, msg.Username) {
		validator = &Validator{
			ABCIValidator: abci.Validator{PubKey: msg.ValPubKey.Bytes(), Power: msg.Deposit.AmountOf(types.Denom)},
			Username:      msg.Username,
			Deposit:       msg.Deposit,
			IsByzantine:   false,
		}
		if setErr := vm.SetValidator(ctx, msg.Username, validator); setErr != nil {
			return setErr.Result()
		}
		vm.AddToCandidatePool(ctx, msg.Username)
	} else {
		validator, err = vm.GetValidator(ctx, msg.Username)
		if err != nil {
			return err.Result()
		}
		validator.Deposit = validator.Deposit.Plus(msg.Deposit)
		validator.ABCIValidator.Power = validator.Deposit.AmountOf("lino")
	}

	if setErr := vm.SetValidator(ctx, msg.Username, validator); setErr != nil {
		return setErr.Result()
	}
	// add to pool and try to become oncall validator
	if joinErr := vm.TryBecomeOncallValidator(ctx, msg.Username); joinErr != nil {
		return joinErr.Result()
	}
	if err := proxyAcc.Apply(ctx); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle Withdraw Msg
func handleWithdrawMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg ValidatorWithdrawMsg) sdk.Result {
	validator, getErr := vm.GetValidator(ctx, msg.Username)
	if getErr != nil {
		return getErr.Result()
	}
	// check the deposit is available now
	if ctx.BlockHeight() < int64(validator.WithdrawAvailableAt) {
		return ErrDepositNotAvailable().Result()
	}
	if !validator.Deposit.IsPositive() {
		return ErrNoDeposit().Result()
	}
	// add money to validator's bank
	proxyAcc := acc.NewProxyAccount(msg.Username, &am)
	if err := proxyAcc.AddCoins(ctx, validator.Deposit); err != nil {
		return err.Result()
	}
	if err := proxyAcc.Apply(ctx); err != nil {
		return err.Result()
	}

	validator.Deposit = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}}
	if err := vm.SetValidator(ctx, msg.Username, validator); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle RevokeMsg
func handleRevokeMsg(ctx sdk.Context, vm ValidatorManager, msg ValidatorRevokeMsg) sdk.Result {
	if err := vm.RemoveValidatorFromAllLists(ctx, msg.Username); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
