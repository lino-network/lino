package vote

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

func NewHandler(vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case VoterDepositMsg:
			return handleVoterDepositMsg(ctx, vm, am, msg)
		case VoterWithdrawMsg:
			return handleVoterWithdrawMsg(ctx, vm, gm, msg)
		case VoterRevokeMsg:
			return handleVoterRevokeMsg(ctx, vm, gm, msg)
		case DelegateMsg:
			return handleDelegateMsg(ctx, vm, am, msg)
		case DelegatorWithdrawMsg:
			return handleDelegatorWithdrawMsg(ctx, vm, gm, msg)
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
func handleVoterDepositMsg(ctx sdk.Context, vm VoteManager, am acc.AccountManager, msg VoterDepositMsg) sdk.Result {
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
		if err := vm.AddVoter(ctx, msg.Username, coin); err != nil {
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
func handleVoterWithdrawMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg VoterWithdrawMsg) sdk.Result {
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

	if err := returnCoinTo(ctx, msg.Username, gm, types.CoinReturnTimes, types.CoinReturnIntervalHr, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle VoterRevokeMsg
func handleVoterRevokeMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg VoterRevokeMsg) sdk.Result {
	fmt.Println(vm.AllValidators)
	// reject if this is a validator
	if vm.IsInValidatorList(ctx, msg.Username) {
		return ErrValidatorCannotRevoke().Result()
	}

	delegators, getErr := vm.GetAllDelegators(ctx, msg.Username)
	if getErr != nil {
		return getErr.Result()
	}

	for _, delegator := range delegators {
		coin, withdrawErr := vm.DelegatorWithdrawAll(ctx, msg.Username, delegator)
		if withdrawErr != nil {
			return withdrawErr.Result()
		}
		if err := returnCoinTo(ctx, delegator, gm, types.CoinReturnTimes, types.CoinReturnIntervalHr, coin); err != nil {
			return err.Result()
		}
	}

	coin, withdrawErr := vm.VoterWithdrawAll(ctx, msg.Username)
	if withdrawErr != nil {
		return withdrawErr.Result()
	}

	if err := returnCoinTo(ctx, msg.Username, gm, types.CoinReturnTimes, types.CoinReturnIntervalHr, coin); err != nil {
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

// Handle DelegatorWithdrawMsg
func handleDelegatorWithdrawMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg DelegatorWithdrawMsg) sdk.Result {
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

	if err := returnCoinTo(ctx, msg.Delegator, gm, types.CoinReturnTimes, types.CoinReturnIntervalHr, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle RevokeDelegationMsg
func handleRevokeDelegationMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg RevokeDelegationMsg) sdk.Result {
	coin, withdrawErr := vm.DelegatorWithdrawAll(ctx, msg.Voter, msg.Delegator)
	if withdrawErr != nil {
		return withdrawErr.Result()
	}

	if err := returnCoinTo(ctx, msg.Delegator, gm, types.CoinReturnTimes, types.CoinReturnIntervalHr, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle VoteMsg
func handleVoteMsg(ctx sdk.Context, vm VoteManager, gm global.GlobalManager, msg VoteMsg) sdk.Result {
	if !vm.IsVoterExist(ctx, msg.Voter) {
		return ErrGetVoter().Result()
	}

	if !vm.IsProposalExist(ctx, msg.ProposalID) {
		return ErrGetProposal().Result()
	}

	if err := vm.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Result); err != nil {
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
	if err := am.MinusCoin(ctx, msg.Creator, types.ProposalRegisterFee); err != nil {
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

func returnCoinTo(ctx sdk.Context, name types.AccountKey, gm global.GlobalManager, times int64, interval int64, coin types.Coin) sdk.Error {
	pieceRat := coin.ToRat().Quo(sdk.NewRat(times))
	piece := types.RatToCoin(pieceRat)
	event := acc.ReturnCoinEvent{
		Username: name,
		Amount:   piece,
	}

	if err := gm.RegisterCoinReturnEvent(ctx, event, times, interval); err != nil {
		return err
	}
	return nil
}
