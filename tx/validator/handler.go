package validator

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
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

	// Must have an normal acount before becoming a validator
	if !am.AccountExist(ctx, msg.ValidatorName) {
		return ErrValidatorManagerFail("normal account not found").Result()
	}

	// withdraw money from validator's bank
	proxyAcc := acc.NewProxyAccount(msg.ValidatorName, &am)
	if err := proxyAcc.MinusCoins(ctx, msg.Deposit); err != nil {
		return ErrValidatorManagerFail("Withdraw money from validator's bank failed").Result()
	}

	account := &ValidatorAccount{
		Validator:     abci.Validator{PubKey: msg.PubKey.Bytes(), Power: msg.Deposit.AmountOf("lino")},
		ValidatorName: msg.ValidatorName,
		Deposit:       msg.Deposit,
	}

	vm.SetValidatorAccount(ctx, msg.ValidatorName, account)

	// add to pool and try to add to validator list
	vm.TryJoinValidatorList(ctx, msg.ValidatorName, true)

	proxyAcc.Apply(ctx)
	return sdk.Result{}
}

// Handle VoteMsg
func handleVoteMsg(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, msg VoteMsg) sdk.Result {
	// Validator not found
	if !vm.IsValidatorExist(ctx, msg.ValidatorName) {
		return ErrValidatorManagerFail("validator not found").Result()
	}

	// withdraw money from voter's bank
	proxyAcc := acc.NewProxyAccount(msg.Voter, &am)
	if err := proxyAcc.MinusCoins(ctx, msg.Power); err != nil {
		return ErrValidatorManagerFail("Withdraw money from voter's bank failed").Result()
	}
	validator, _ := vm.GetValidatorAccount(ctx, msg.ValidatorName)
	vote := Vote{
		voter:         msg.Voter,
		power:         msg.Power,
		validatorName: msg.ValidatorName,
	}

	validator.Votes = append(validator.Votes, vote)
	validator.Power += msg.Power.AmountOf("lino")
	vm.SetValidatorAccount(ctx, msg.ValidatorName, validator)
	vm.TryJoinValidatorList(ctx, msg.ValidatorName, false)

	proxyAcc.Apply(ctx)
	return sdk.Result{}
}
