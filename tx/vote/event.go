package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	types "github.com/lino-network/lino/types"
)

type DecideProposalEvent struct{}

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
	lst, err := voteManager.storage.GetProposalList(ctx)
	if err != nil {
		return types.ProposalKey(""), err
	}

	curID := lst.OngoingProposal[0]
	lst.OngoingProposal = lst.OngoingProposal[1:]
	lst.PastProposal = append(lst.PastProposal, curID)

	if err := voteManager.storage.SetProposalList(ctx, lst); err != nil {
		return curID, err
	}
	return curID, nil
}

func (dpe DecideProposalEvent) calculateVotingResult(ctx sdk.Context, curID types.ProposalKey, vm VoteManager) (bool, sdk.Error) {
	// get all votes to calculate the voting result
	votes, err := vm.storage.GetAllVotes(ctx, curID)
	if err != nil {
		return false, err
	}
	referenceList, err := vm.storage.GetValidatorReferenceList(ctx)
	if err != nil {
		return false, err
	}
	validators := make([]types.AccountKey, len(referenceList.OncallValidators))
	copy(validators, referenceList.OncallValidators)

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
		referenceList.PenaltyValidators = append(referenceList.PenaltyValidators, validator)
	}
	if err := vm.storage.SetValidatorReferenceList(ctx, referenceList); err != nil {
		return false, err
	}
	return proposal.AgreeVote.IsGT(proposal.DisagreeVote), nil
}

func (dpe DecideProposalEvent) changeParameter(ctx sdk.Context, curID types.ProposalKey, voteManager VoteManager, gm global.GlobalManager) sdk.Error {
	proposal, err := voteManager.storage.GetProposal(ctx, curID)
	if err != nil {
		return err
	}
	des := proposal.ChangeParameterDescription
	if err := gm.ChangeInfraInternalInflationParam(ctx, des.StorageAllocation, des.CDNAllocation); err != nil {
		return err
	}

	if err := gm.ChangeGlobalInflationParam(ctx, des.InfraAllocation, des.ContentCreatorAllocation,
		des.DeveloperAllocation, des.ValidatorAllocation); err != nil {
		return err
	}
	return nil
}
