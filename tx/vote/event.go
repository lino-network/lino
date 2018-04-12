package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	//val "github.com/lino-network/lino/tx/validator"
	types "github.com/lino-network/lino/types"
)

type ReturnCoinEvent struct {
	Username types.AccountKey `json:"username"`
	Amount   types.Coin       `json:"amount"`
}

type DecideProposalEvent struct{}

func (rce ReturnCoinEvent) Execute(ctx sdk.Context, voteManager VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	if !am.IsAccountExist(ctx, rce.Username) {
		return acc.ErrUsernameNotFound()
	}

	if err := am.AddCoin(ctx, rce.Username, rce.Amount); err != nil {
		return err
	}

	return nil
}

func (dpe DecideProposalEvent) Execute(ctx sdk.Context, voteManager VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	// update the ongoing and past proposal list
	curID, updateErr := dpe.updateProposalList(ctx, voteManager)
	if updateErr != nil {
		return updateErr
	}

	// calculate voting result and set absent validators
	pass, calErr := dpe.calculateVotingResult(ctx, curID, voteManager)
	if calErr != nil {
		return calErr
	}

	// majority disagree this proposal
	if !pass {
		return nil
	}

	// change parameter
	if err := dpe.changeParameter(ctx, curID, voteManager, gm); err != nil {
		return err
	}
	return nil
}

func (dpe DecideProposalEvent) updateProposalList(ctx sdk.Context, voteManager VoteManager) (types.ProposalKey, sdk.Error) {
	lst, getErr := voteManager.storage.GetProposalList(ctx)
	if getErr != nil {
		return types.ProposalKey(""), getErr
	}

	curID := lst.OngoingProposal[0]
	lst.OngoingProposal = lst.OngoingProposal[1:]
	lst.PastProposal = append(lst.PastProposal, curID)

	if setErr := voteManager.storage.SetProposalList(ctx, lst); setErr != nil {
		return curID, setErr
	}
	return curID, nil
}

func (dpe DecideProposalEvent) calculateVotingResult(ctx sdk.Context, curID types.ProposalKey, vm VoteManager) (bool, sdk.Error) {
	// get all votes to calculate the voting result
	votes, getErr := vm.storage.GetAllVotes(ctx, curID)
	if getErr != nil {
		return false, getErr
	}

	validators := make([]types.AccountKey, len(vm.OncallValidators))
	copy(validators, vm.OncallValidators)

	// get the proposal we are going to decide
	proposal, err := vm.storage.GetProposal(ctx, curID)
	if err != nil {
		return false, err
	}

	for _, vote := range votes {
		voterPower, err := vm.GetVotingPower(ctx, vote.Voter)
		if err != nil {
			continue
		}
		if vote.Result == true {
			proposal.AgreeVote = proposal.AgreeVote.Plus(voterPower)
		} else {
			proposal.DisagreeVote = proposal.DisagreeVote.Plus(voterPower)
		}

		// remove from list if the validator voted
		for idx, validator := range validators {
			if validator == vote.Voter {
				validators = append(validators[:idx], validators[idx+1:]...)
				break
			}
		}
		vm.storage.DeleteVote(ctx, curID, vote.Voter)
	}

	if err := vm.storage.SetProposal(ctx, curID, proposal); err != nil {
		return false, err
	}

	// put all validators who didn't vote into penalty list
	for _, validator := range validators {
		vm.PenaltyValidators = append(vm.PenaltyValidators, validator)
	}
	return true, nil
}

func (dpe DecideProposalEvent) changeParameter(ctx sdk.Context, curID types.ProposalKey, voteManager VoteManager, gm global.GlobalManager) sdk.Error {
	proposal, getErr := voteManager.storage.GetProposal(ctx, curID)
	if getErr != nil {
		return getErr
	}
	des := proposal.ChangeParameterDescription
	if err := gm.ChangeInfraInternalInflation(ctx, des.StorageAllocation, des.CDNAllocation); err != nil {
		return err
	}

	if err := gm.ChangeGlobalInflation(ctx, des.InfraAllocation, des.ContentCreatorAllocation,
		des.DeveloperAllocation, des.ValidatorAllocation); err != nil {
		return err
	}
	return nil
}
