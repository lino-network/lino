package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	val "github.com/lino-network/lino/tx/validator"
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

func (dpe DecideProposalEvent) Execute(ctx sdk.Context, voteManager VoteManager, valManager val.ValidatorManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	// update the ongoing and past proposal list
	curID, updateErr := dpe.updateProposalList(ctx, voteManager)
	if updateErr != nil {
		return updateErr
	}

	// calculate voting result and set absent validators
	pass, calErr := dpe.calculateVotingResult(ctx, curID, voteManager, valManager)
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

func (dpe DecideProposalEvent) calculateVotingResult(ctx sdk.Context, curID types.ProposalKey, voteManager VoteManager, valManager val.ValidatorManager) (bool, sdk.Error) {
	// get all votes to calculate the voting result
	votes, getErr := voteManager.storage.GetAllVotes(ctx, curID)
	if getErr != nil {
		return false, getErr
	}

	// get the proposal we are going to decide
	proposal, err := voteManager.storage.GetProposal(ctx, curID)
	if err != nil {
		return false, err
	}

	oncallValidators, getListErr := valManager.GetOncallValList(ctx)
	if getListErr != nil {
		return false, getListErr
	}

	for _, vote := range votes {
		voterPower, err := voteManager.GetVotingPower(ctx, vote.Voter)
		if err != nil {
			continue
		}
		if vote.Result == true {
			proposal.AgreeVote = proposal.AgreeVote.Plus(voterPower)
		} else {
			proposal.DisagreeVote = proposal.DisagreeVote.Plus(voterPower)
		}

		// remove from list if the validator voted
		for idx, validator := range oncallValidators {
			if validator.Username == vote.Voter {
				oncallValidators = append(oncallValidators[:idx], oncallValidators[idx+1:]...)
				break
			}
		}
		voteManager.storage.DeleteVote(ctx, curID, vote.Voter)
	}

	if err := voteManager.storage.SetProposal(ctx, curID, proposal); err != nil {
		return false, err
	}

	// add absent vote for all validator didn't vote
	for _, validator := range oncallValidators {
		if err := valManager.MarkAbsentVote(ctx, validator.Username); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (dpe DecideProposalEvent) changeParameter(ctx sdk.Context, curID types.ProposalKey, voteManager VoteManager, gm global.GlobalManager) sdk.Error {
	_, getErr := voteManager.storage.GetProposal(ctx, curID)
	if getErr != nil {
		return getErr
	}
	return nil
}
