package validator

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
)

func NewHandler(vm ValidatorManager, am acc.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ValidatorRegisterMsg:
			return handleRegisterMsg(ctx, vm, am, msg)
		case VoteMsg:
			return handleVoteMsg(ctx, vm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleRegisterMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg ValidatorRegisterMsg) sdk.Result {
	return sdk.Result{}
}

// Handle VoteMsg
func handleVoteMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg VoteMsg) sdk.Result {
	return sdk.Result{}
}
