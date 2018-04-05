package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
)

type ReturnCoinEvent struct {
	Username acc.AccountKey `json:"username"`
	Amount   types.Coin     `json:"amount"`
}

type DecideProposalEvent struct {
	//ProposalID ProposalKey `json:"proposal_id"`
}

func (event ReturnCoinEvent) Execute(ctx sdk.Context, vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	account := acc.NewProxyAccount(event.Username, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound()
	}

	if err := account.AddCoin(ctx, event.Amount); err != nil {
		return err
	}
	if err := account.Apply(ctx); err != nil {
		return err
	}

	return nil
}

func (event DecideProposalEvent) Execute(ctx sdk.Context, vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	lst, getErr := vm.GetProposalList(ctx)
	if getErr != nil {
		return getErr
	}

	curID := lst.OngoingProposal[0]
	lst.OngoingProposal = lst.OngoingProposal[1:]
	lst.PastProposal = append(lst.PastProposal, curID)

	if setErr := vm.SetProposalList(ctx, lst); setErr != nil {
		return setErr
	}

	votes, getErr := vm.GetAllVotes(ctx, curID)
	if getErr != nil {
		return getErr
	}

	proposal, err := vm.GetProposal(ctx, curID)
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
	}

	if err := vm.SetProposal(ctx, curID, proposal); err != nil {
		return err
	}
	// majority disagree this proposal
	if proposal.DisagreeVote.IsGTE(proposal.AgreeVote) {
		return nil
	}

	// change parameter

	// delete all votes
	return nil
}
