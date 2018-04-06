package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
)

type ReturnCoinEvent struct {
	Username types.AccountKey `json:"username"`
	Amount   types.Coin       `json:"amount"`
}

type DecideProposalEvent struct{}

func (event ReturnCoinEvent) Execute(ctx sdk.Context, vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	if !am.IsAccountExist(ctx, event.Username) {
		return acc.ErrUsernameNotFound()
	}

	if err := am.AddCoin(ctx, event.Username, event.Amount); err != nil {
		return err
	}

	return nil
}

func (event DecideProposalEvent) Execute(ctx sdk.Context, vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	// update the proposal list
	lst, getErr := vm.storage.GetProposalList(ctx)
	if getErr != nil {
		return getErr
	}

	curID := lst.OngoingProposal[0]
	lst.OngoingProposal = lst.OngoingProposal[1:]
	lst.PastProposal = append(lst.PastProposal, curID)

	if setErr := vm.storage.SetProposalList(ctx, lst); setErr != nil {
		return setErr
	}

	// get all votes to calculate the voting result
	votes, getErr := vm.storage.GetAllVotes(ctx, curID)
	if getErr != nil {
		return getErr
	}

	proposal, err := vm.storage.GetProposal(ctx, curID)
	if err != nil {
		return err
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
		// delete this vote
		vm.storage.DeleteVote(ctx, curID, vote.Voter)
	}

	if err := vm.storage.SetProposal(ctx, curID, proposal); err != nil {
		return err
	}
	// majority disagree this proposal
	if proposal.DisagreeVote.IsGTE(proposal.AgreeVote) {
		return nil
	}

	// change parameter
	return nil
}
