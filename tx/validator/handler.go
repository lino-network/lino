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
		case ValidatorRegisterMsg:
			return handleRegisterMsg(ctx, vm, am, msg)
		case VoteMsg:
			return handleVoteMsg(ctx, vm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleRegisterMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg ValidatorRegisterMsg) sdk.Result {
	proxyAcc := acc.NewProxyAccount(msg.Username, &am)
	// This name has been registered
	if vm.IsValidatorExist(ctx, msg.Username) {
		return ErrValidatorHandlerFail("validator exists").Result()
	}

	// Must have an normal acount before becoming a validator
	if !proxyAcc.IsAccountExist(ctx) {
		return ErrValidatorHandlerFail("user account not found").Result()
	}

	// withdraw money from validator's bank
	if err := proxyAcc.MinusCoins(ctx, msg.Deposit); err != nil {
		return err.Result()
	}
	ownerKey, getErr := proxyAcc.GetOwnerKey(ctx)
	if getErr != nil {
		return getErr.Result()
	}
	account := &Validator{
		ABCIValidator: abci.Validator{PubKey: ownerKey.Bytes(), Power: msg.Deposit.AmountOf(types.Denom)},
		Username:      msg.Username,
		Deposit:       msg.Deposit,
	}

	if setErr := vm.SetValidator(ctx, msg.Username, account); setErr != nil {
		return setErr.Result()
	}

	// add to pool and try to add to validator list
	if joinErr := vm.TryJoinValidatorList(ctx, msg.Username, true); joinErr != nil {
		return joinErr.Result()
	}
	if err := proxyAcc.Apply(ctx); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle VoteMsg
func handleVoteMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg VoteMsg) sdk.Result {
	// Validator not found
	if !vm.IsValidatorExist(ctx, msg.ValidatorName) {
		return ErrValidatorHandlerFail("validator not found").Result()
	}

	// withdraw money from voter's bank
	proxyAcc := acc.NewProxyAccount(msg.Voter, &am)
	if err := proxyAcc.MinusCoins(ctx, msg.Power); err != nil {
		return err.Result()
	}
	validator, getErr := vm.GetValidator(ctx, msg.ValidatorName)
	if getErr != nil {
		return getErr.Result()
	}
	vote := Vote{
		Voter: msg.Voter,
		Power: msg.Power,
	}

	validator.Votes = append(validator.Votes, vote)
	validator.ABCIValidator.Power += msg.Power.AmountOf(types.Denom)

	if setErr := vm.SetValidator(ctx, msg.ValidatorName, validator); setErr != nil {
		return setErr.Result()
	}
	if joinErr := vm.TryJoinValidatorList(ctx, msg.ValidatorName, false); joinErr != nil {
		return joinErr.Result()
	}

	if err := proxyAcc.Apply(ctx); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
