package validator

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/x/validator/types"
)

// NewHandler - Handle all "validator" type messages.
func NewHandler(vm ValidatorKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.ValidatorRegisterMsg:
			return handleValidatorRegisterMsg(ctx, vm, msg)
		case types.ValidatorRevokeMsg:
			return handleValidatorRevokeMsg(ctx, vm, msg)
		case types.VoteValidatorMsg:
			return handleVoteValidatorMsg(ctx, vm, msg)
		case types.ValidatorUpdateMsg:
			return handleValidatorUpdateMsg(ctx, vm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleValidatorRegisterMsg(ctx sdk.Context, vm ValidatorKeeper, msg types.ValidatorRegisterMsg) sdk.Result {
	if err := vm.RegisterValidator(ctx, msg.Username, msg.ValPubKey, msg.Link); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleValidatorRevokeMsg(
	ctx sdk.Context, vm ValidatorKeeper, msg types.ValidatorRevokeMsg) sdk.Result {
	if err := vm.RevokeValidator(ctx, msg.Username); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleVoteValidatorMsg(
	ctx sdk.Context, vm ValidatorKeeper, msg types.VoteValidatorMsg) sdk.Result {
	if err := vm.VoteValidator(ctx, msg.Username, msg.VotedValidators); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleValidatorUpdateMsg(ctx sdk.Context, vm ValidatorKeeper, msg types.ValidatorUpdateMsg) sdk.Result {
	if err := vm.UpdateValidator(ctx, msg.Username, msg.Link); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
