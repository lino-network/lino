package proposal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/tx/post"
	val "github.com/lino-network/lino/tx/validator"
	"github.com/lino-network/lino/tx/vote"
	types "github.com/lino-network/lino/types"
)

type DecideProposalEvent struct {
	ProposalType types.ProposalType `json:"proposal_type"`
	ProposalID   types.ProposalKey  `json:"proposal_id"`
}

func (dpe DecideProposalEvent) Execute(
	ctx sdk.Context, voteManager vote.VoteManager, valManager val.ValidatorManager,
	am acc.AccountManager, proposalManager ProposalManager, postManager post.PostManager, gm global.GlobalManager) sdk.Error {
	// get the proposal ID we are going to decide
	curID, err := proposalManager.GetCurrentProposal(ctx)
	if err != nil {
		return err
	}

	// get all oncall validators (make sure they voted on certain type of proposal)
	lst, err := valManager.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	// calculate voting result
	votingRes, err := voteManager.CalculateVotingResult(ctx, curID, lst.OncallValidators)
	if err != nil {
		return err
	}

	// punish validators who didn't vote
	actualPenalty, err := valManager.PunishValidatorsDidntVote(ctx, votingRes.PenaltyList)
	if err != nil {
		return err
	}

	// add coins back to inflation pool
	if err := gm.AddToValidatorInflationPool(ctx, actualPenalty); err != nil {
		return err
	}

	// update the ongoing and past proposal list
	proposalRes, err := proposalManager.UpdateProposalStatus(
		ctx, votingRes, dpe.ProposalType, dpe.ProposalID)
	if err != nil {
		return err
	}

	// majority disagree this proposal
	if proposalRes == types.ProposalNotPass {
		return nil
	}

	// execute proposal
	switch dpe.ProposalType {
	case types.ChangeParam:
		if err := dpe.ExecuteChangeParam(ctx, curID, proposalManager, gm); err != nil {
			return err
		}
	case types.ContentCensorship:
		if err := dpe.ExecuteContentCensorship(ctx, curID, proposalManager, postManager); err != nil {
			return err
		}
	case types.ProtocolUpgrade:
		if err := dpe.ExecuteProtocolUpgrade(ctx, curID, proposalManager); err != nil {
			return err
		}
	}
	return nil
}

func (dpe DecideProposalEvent) ExecuteChangeParam(
	ctx sdk.Context, curID types.ProposalKey, proposalManager ProposalManager, gm global.GlobalManager) sdk.Error {
	event, err := proposalManager.CreateParamChangeEvent(ctx, curID)
	if err != nil {
		return err
	}
	if err := gm.RegisterParamChangeEvent(ctx, event); err != nil {
		return err
	}
	return nil
}

func (dpe DecideProposalEvent) ExecuteContentCensorship(
	ctx sdk.Context, curID types.ProposalKey, proposalManager ProposalManager, postManager post.PostManager) sdk.Error {
	_, err := proposalManager.GetPermLink(ctx, curID)
	if err != nil {
		return err
	}

	// TODO add content censorship logic
	return nil
}

func (dpe DecideProposalEvent) ExecuteProtocolUpgrade(
	ctx sdk.Context, curID types.ProposalKey, proposalManager ProposalManager) sdk.Error {
	return nil
}
