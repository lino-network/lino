package vote

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

func NewHandler(vm VoteManager, am acc.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case VoterDepositMsg:
			return handleDepositMsg(ctx, vm, am, msg)
		case VoterWithdrawMsg:
			return handleWithdrawMsg(ctx, vm, am, msg)
		case VoterRevokeMsg:
			return handleRevokeMsg(ctx, vm, msg)
		case DelegateMsg:
			return handleDelegateMsg(ctx, vm, msg)
		case RevokeDelegationMsg:
			return handleRevokeDelegationMsg(ctx, vm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle DepositMsg
func handleDepositMsg(ctx sdk.Context, vm VoteManager, am acc.AccountManager, msg VoterDepositMsg) sdk.Result {
	proxyAcc := acc.NewProxyAccount(msg.Username, &am)
	// Must have an normal acount
	if !proxyAcc.IsAccountExist(ctx) {
		return ErrUsernameNotFound().Result()
	}

	coin, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	// withdraw money from voter's bank
	err = proxyAcc.MinusCoin(ctx, coin)
	if err != nil {
		return err.Result()
	}
	if err := proxyAcc.Apply(ctx); err != nil {
		return err.Result()
	}

	// Register the user if this name has not been registered
	if !vm.IsVoterExist(ctx, msg.Username) {
		if err := vm.RegisterVoter(ctx, msg.Username, coin); err != nil {
			return err.Result()
		}
	} else {
		// Deposit coins
		if err := vm.Deposit(ctx, msg.Username, coin); err != nil {
			return err.Result()
		}
	}
	return sdk.Result{}
}

// Handle Withdraw Msg
func handleWithdrawMsg(ctx sdk.Context, vm VoteManager, am acc.AccountManager, msg VoterWithdrawMsg) sdk.Result {

	return sdk.Result{}
}

// Handle RevokeMsg
func handleRevokeMsg(ctx sdk.Context, vm VoteManager, msg VoterRevokeMsg) sdk.Result {
	// if err := vm.RemoveValidatorFromAllLists(ctx, msg.Username); err != nil {
	// 	return err.Result()
	// }
	return sdk.Result{}
}

// Handle DelegateMsg
func handleDelegateMsg(ctx sdk.Context, vm VoteManager, msg DelegateMsg) sdk.Result {
	// if err := vm.RemoveValidatorFromAllLists(ctx, msg.Username); err != nil {
	// 	return err.Result()
	// }
	return sdk.Result{}
}

// Handle RevokeDelegationMsg
func handleRevokeDelegationMsg(ctx sdk.Context, vm VoteManager, msg RevokeDelegationMsg) sdk.Result {
	// if err := vm.RemoveValidatorFromAllLists(ctx, msg.Username); err != nil {
	// 	return err.Result()
	// }
	return sdk.Result{}
}
