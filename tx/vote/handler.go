package vote

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
)

func NewHandler(vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case VoterDepositMsg:
			return handleDepositMsg(ctx, vm, am, msg)
		case VoterWithdrawMsg:
			return handleWithdrawMsg(ctx, vm, gm, msg)
		case VoterRevokeMsg:
			return handleRevokeMsg(ctx, vm, gm, msg)
		case DelegateMsg:
			return handleDelegateMsg(ctx, vm, am, msg)
		case RevokeDelegationMsg:
			return handleRevokeDelegationMsg(ctx, vm, gm, msg)
		case VoteMsg:
			return handleVoteMsg(ctx, vm, gm, msg)
		case CreateProposalMsg:
			return handleCreateProposalMsg(ctx, vm, am, gm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle DepositMsg
func handleDepositMsg(ctx sdk.Context, vm VoteManager, am acc.AccountManager, msg VoterDepositMsg) sdk.Result {
	// Must have an normal acount
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrUsernameNotFound().Result()
	}

	coin, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	// withdraw money from voter's bank
	if err := am.MinusCoin(ctx, msg.Username, coin); err != nil {
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
func handleWithdrawMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg VoterWithdrawMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	if !vm.IsLegalWithdraw(ctx, msg.Username, coin) {
		return ErrIllegalWithdraw().Result()
	}

	if err := vm.Withdraw(ctx, msg.Username, coin, gm); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle RevokeMsg
func handleRevokeMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg VoterRevokeMsg) sdk.Result {
	// TODO also a Validator
	delegators, getErr := vm.storage.GetAllDelegators(ctx, msg.Username)
	if getErr != nil {
		return getErr.Result()
	}

	for _, delegator := range delegators {
		if err := vm.ReturnCoinToDelegator(ctx, msg.Username, delegator, gm); err != nil {
			return err.Result()
		}
	}

	if err := vm.WithdrawAll(ctx, msg.Username, gm); err != nil {
		return err.Result()
	}

	if err := vm.storage.DeleteVoter(ctx, msg.Username); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle DelegateMsg
func handleDelegateMsg(ctx sdk.Context, vm VoteManager, am acc.AccountManager, msg DelegateMsg) sdk.Result {
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	// withdraw money from delegator's bank
	if err := am.MinusCoin(ctx, msg.Delegator, coin); err != nil {
		return err.Result()
	}
	// add delegation relation
	if addErr := vm.AddDelegation(ctx, msg.Voter, msg.Delegator, coin); addErr != nil {
		return addErr.Result()
	}
	return sdk.Result{}
}

// Handle RevokeDelegationMsg
func handleRevokeDelegationMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg RevokeDelegationMsg) sdk.Result {
	if err := vm.ReturnCoinToDelegator(ctx, msg.Voter, msg.Delegator, gm); err != nil {
		return err.Result()
	}
	if err := vm.storage.DeleteDelegation(ctx, msg.Voter, msg.Delegator); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle VoteMsg
func handleVoteMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg VoteMsg) sdk.Result {
	if !vm.IsVoterExist(ctx, msg.Voter) {
		return ErrGetVoter().Result()
	}
	vote := model.Vote{
		Voter:  msg.Voter,
		Result: msg.Result,
	}

	if !vm.IsProposalExist(ctx, msg.ProposalID) {
		return ErrGetProposal().Result()
	}
	// will overwrite the old vote
	if err := vm.storage.SetVote(ctx, msg.ProposalID, msg.Voter, &vote); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle CreateProposalMsg
func handleCreateProposalMsg(ctx sdk.Context, vm VoteManager, am acc.AccountManager, gm global.GlobalManager, msg CreateProposalMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Creator) {
		return ErrUsernameNotFound().Result()
	}

	// withdraw money from creator's bank
	if err := am.MinusCoin(ctx, msg.Creator, model.ProposalRegisterFee); err != nil {
		return err.Result()
	}

	if addErr := vm.AddProposal(ctx, msg.Creator, &msg.ChangeParameterDescription); addErr != nil {
		return addErr.Result()
	}
	//  set a time event to decide the proposal in 7 days
	if err := vm.CreateDecideProposalEvent(ctx, gm); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
