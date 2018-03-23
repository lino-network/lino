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
	// This name has been registered
	if vm.IsValidatorExist(ctx, msg.ValidatorName) {
		return ErrValidatorManagerFail("account exist").Result()
	}

	// Must have an acount before becoming a validator
	if !am.AccountExist(ctx, msg.ValidatorName) {
		return ErrValidatorManagerFail("account exist").Result()
	}

	// withdraw money from validator's bank
	proxyAcc := acc.NewProxyAccount(msg.ValidatorName, &am)
	if err := proxyAcc.MinusCoins(ctx, msg.Deposit); err != nil {
		return ErrValidatorManagerFail("Withdraw money from validator's bank failed").Result()
	}
	// TODO: publick key?
	account := &ValidatorAccount{
		validatorName: msg.ValidatorName,
		//totalWeight:   msg.Deposit,
		deposit: msg.Deposit,
	}

	vm.SetValidatorAccount(ctx, msg.ValidatorName, account)

	// add to validator list
	// TODO: key?
	lstPtr, _ := vm.GetValidatorList(ctx, "validatoryKey")
	lstPtr.validators = append(lstPtr.validators, *account)
	vm.SetValidatorList(ctx, "validatoryKey", lstPtr)

	proxyAcc.Apply(ctx)
	return sdk.Result{}
}

// Handle VoteMsg
func handleVoteMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg VoteMsg) sdk.Result {
	// Validator not found
	if !vm.IsValidatorExist(ctx, msg.ValidatorName) {
		return ErrValidatorManagerFail("validator not found").Result()
	}

	validator, _ := vm.GetValidatorAccount(ctx, msg.ValidatorName)
	vote := Vote{
		voter:         msg.Voter,
		weight:        msg.Weight,
		validatorName: msg.ValidatorName,
	}

	validator.votes = append(validator.votes, vote)
	//validator.totalWeight += msg.Weight
	vm.SetValidatorAccount(ctx, msg.ValidatorName, validator)
	return sdk.Result{}
}
