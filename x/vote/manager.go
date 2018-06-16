package vote

import (
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteManager struct {
	storage     model.VoteStorage `json:"vote_storage"`
	paramHolder param.ParamHolder `json:"param_holder"`
}

func NewVoteManager(key sdk.StoreKey, holder param.ParamHolder) VoteManager {
	return VoteManager{
		storage:     model.NewVoteStorage(key),
		paramHolder: holder,
	}
}

func (vm VoteManager) InitGenesis(ctx sdk.Context) error {
	if err := vm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) IsVoterExist(ctx sdk.Context, accKey types.AccountKey) bool {
	voterByte, _ := vm.storage.GetVoter(ctx, accKey)
	return voterByte != nil
}

func (vm VoteManager) IsVoteExist(ctx sdk.Context, proposalID types.ProposalKey, accKey types.AccountKey) bool {
	voteByte, _ := vm.storage.GetVote(ctx, proposalID, accKey)
	return voteByte != nil
}

func (vm VoteManager) IsInValidatorList(ctx sdk.Context, username types.AccountKey) bool {
	lst, err := vm.storage.GetReferenceList(ctx)
	if err != nil {
		return false
	}
	for _, validator := range lst.AllValidators {
		if validator == username {
			return true
		}
	}
	return false
}

func (vm VoteManager) IsDelegationExist(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) bool {
	delegationByte, _ := vm.storage.GetDelegation(ctx, voter, delegator)
	return delegationByte != nil
}

func (vm VoteManager) IsLegalVoterWithdraw(
	ctx sdk.Context, username types.AccountKey, coin types.Coin) bool {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return false
	}
	// reject if this is a validator
	if vm.IsInValidatorList(ctx, username) {
		return false
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return false
	}
	// reject if withdraw is less than minimum voter withdraw
	if !coin.IsGTE(param.VoterMinWithdraw) {
		return false
	}
	//reject if the remaining coins are less than voter minimum deposit
	remaining := voter.Deposit.Minus(coin)
	if !remaining.IsGTE(param.VoterMinDeposit) {
		return false
	}
	return true
}

func (vm VoteManager) IsLegalDelegatorWithdraw(
	ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, coin types.Coin) bool {
	delegation, err := vm.storage.GetDelegation(ctx, voterName, delegatorName)
	if err != nil {
		return false
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return false
	}

	// reject if withdraw is less than minimum delegator withdraw
	if !coin.IsGTE(param.DelegatorMinWithdraw) {
		return false
	}
	//reject if the remaining delegation are less than zero
	res := delegation.Amount.Minus(coin)
	return res.IsNotNegative()
}

func (vm VoteManager) CanBecomeValidator(ctx sdk.Context, username types.AccountKey) bool {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return false
	}
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return false
	}
	// check minimum voting deposit for validator
	return voter.Deposit.IsGTE(param.ValidatorMinVotingDeposit)
}

// only support change parameter proposal now
func (vm VoteManager) AddVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey, res bool) sdk.Error {
	// check if the vote exist
	if vm.IsVoteExist(ctx, proposalID, voter) {
		return ErrVoteExist()
	}

	votingPower, err := vm.GetVotingPower(ctx, voter)
	if err != nil {
		return err
	}

	vote := model.Vote{
		Voter:       voter,
		Result:      res,
		VotingPower: votingPower,
	}

	if err := vm.storage.SetVote(ctx, proposalID, voter, &vote); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) GetVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey) (*model.Vote, sdk.Error) {
	return vm.storage.GetVote(ctx, proposalID, voter)
}

func (vm VoteManager) AddDelegation(ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, coin types.Coin) sdk.Error {
	var delegation *model.Delegation
	var err sdk.Error

	if !vm.IsDelegationExist(ctx, voterName, delegatorName) {
		delegation = &model.Delegation{
			Delegator: delegatorName,
		}
		if err := vm.AddDelegatee(ctx, delegatorName, voterName); err != nil {
			return err
		}
	} else {
		delegation, err = vm.storage.GetDelegation(ctx, voterName, delegatorName)
		if err != nil {
			return err
		}
	}

	voter, err := vm.storage.GetVoter(ctx, voterName)
	if err != nil {
		return err
	}

	voter.DelegatedPower = voter.DelegatedPower.Plus(coin)
	delegation.Amount = delegation.Amount.Plus(coin)

	if err := vm.storage.SetDelegation(ctx, voterName, delegatorName, delegation); err != nil {
		return err
	}
	if err := vm.storage.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) AddDelegatee(
	ctx sdk.Context, me types.AccountKey, delegatee types.AccountKey) sdk.Error {
	delegateeList, err := vm.storage.GetDelegateeList(ctx, me)
	if err != nil {
		return err
	}
	if delegateeList == nil {
		delegateeList = &model.DelegateeList{[]types.AccountKey{}}
	}
	delegateeList.DelegateeList = append(delegateeList.DelegateeList, delegatee)
	return vm.storage.SetDelegateeList(ctx, me, delegateeList)
}

func (vm VoteManager) AddVoter(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	voter := &model.Voter{
		Username: username,
		Deposit:  coin,
	}

	param, err := vm.paramHolder.GetVoteParam(ctx)
	if err != nil {
		return err
	}

	// check minimum requirements for registering as a voter
	if !coin.IsGTE(param.VoterMinDeposit) {
		return ErrRegisterFeeNotEnough()
	}

	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) Deposit(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	voter.Deposit = voter.Deposit.Plus(coin)
	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	return nil
}

// this method won't check if it is a legal withdraw, caller should check by itself
func (vm VoteManager) VoterWithdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return ErrNoCoinToWithdraw()
	}
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	voter.Deposit = voter.Deposit.Minus(coin)

	if voter.Deposit.IsZero() {
		if err := vm.storage.DeleteVoter(ctx, username); err != nil {
			return err
		}
	} else {
		if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
			return err
		}
	}

	return nil
}

func (vm VoteManager) VoterWithdrawAll(ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := vm.VoterWithdraw(ctx, username, voter.Deposit); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return voter.Deposit, nil
}

func (vm VoteManager) DelegatorWithdraw(
	ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return ErrNoCoinToWithdraw()
	}
	// change voter's delegated power
	voter, err := vm.storage.GetVoter(ctx, voterName)
	if err != nil {
		return err
	}
	voter.DelegatedPower = voter.DelegatedPower.Minus(coin)
	if err := vm.storage.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}

	// change this delegation's amount
	delegation, err := vm.storage.GetDelegation(ctx, voterName, delegatorName)
	if err != nil {
		return err
	}
	delegation.Amount = delegation.Amount.Minus(coin)

	if delegation.Amount.IsZero() {
		if err := vm.storage.DeleteDelegation(ctx, voterName, delegatorName); err != nil {
			return err
		}
		if err := vm.RemoveDelegatee(ctx, delegatorName, voterName); err != nil {
			return err
		}
	} else {
		vm.storage.SetDelegation(ctx, voterName, delegatorName, delegation)
	}

	return nil
}

func (vm VoteManager) RemoveDelegatee(
	ctx sdk.Context, me types.AccountKey, delegateeToRemove types.AccountKey) sdk.Error {
	delegateeList, err := vm.storage.GetDelegateeList(ctx, me)
	if err != nil {
		return err
	}
	idx := 0
	for idx < len(delegateeList.DelegateeList) {
		if delegateeList.DelegateeList[idx] == delegateeToRemove {
			delegateeList.DelegateeList =
				append(delegateeList.DelegateeList[:idx], delegateeList.DelegateeList[idx+1:]...)
			continue
		}
		idx += 1
	}
	return vm.storage.SetDelegateeList(ctx, me, delegateeList)
}

func (vm VoteManager) DelegatorWithdrawAll(ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey) (types.Coin, sdk.Error) {
	delegation, err := vm.storage.GetDelegation(ctx, voterName, delegatorName)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := vm.DelegatorWithdraw(ctx, voterName, delegatorName, delegation.Amount); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return delegation.Amount, nil
}

func (vm VoteManager) GetVotingPower(ctx sdk.Context, voterName types.AccountKey) (types.Coin, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, voterName)
	if err != nil {
		return types.Coin{}, err
	}
	res := voter.Deposit.Plus(voter.DelegatedPower)
	return res, nil
}

func (vm VoteManager) GetPenaltyList(
	ctx sdk.Context, proposalID types.ProposalKey, proposalType types.ProposalType,
	oncallValidators []types.AccountKey) (types.PenaltyList, sdk.Error) {
	penaltyList := types.PenaltyList{
		PenaltyList: []types.AccountKey{},
	}

	// get all votes to calculate the voting result
	votes, err := vm.storage.GetAllVotes(ctx, proposalID)
	if err != nil {
		return penaltyList, err
	}

	for _, vote := range votes {
		// remove from list if the validator voted
		for idx, validator := range oncallValidators {
			if validator == vote.Voter {
				oncallValidators = append(oncallValidators[:idx], oncallValidators[idx+1:]...)
				break
			}
		}
		// TODO: Decide if we wanna delete vote.
		// vm.storage.DeleteVote(ctx, proposalID, vote.Voter)
	}

	// put all validators who didn't vote on these two types proposal into penalty list
	if proposalType == types.ChangeParam || proposalType == types.ProtocolUpgrade {
		penaltyList.PenaltyList = oncallValidators
	}
	return penaltyList, nil
}

func (vm VoteManager) GetVoterDeposit(ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, accKey)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return voter.Deposit, nil
}

func (vm VoteManager) GetAllDelegators(ctx sdk.Context, voterName types.AccountKey) ([]types.AccountKey, sdk.Error) {
	return vm.storage.GetAllDelegators(ctx, voterName)
}

func (vm VoteManager) GetValidatorReferenceList(ctx sdk.Context) (*model.ReferenceList, sdk.Error) {
	return vm.storage.GetReferenceList(ctx)
}

func (vm VoteManager) SetValidatorReferenceList(ctx sdk.Context, lst *model.ReferenceList) sdk.Error {
	return vm.storage.SetReferenceList(ctx, lst)
}
