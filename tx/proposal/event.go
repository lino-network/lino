package proposal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	val "github.com/lino-network/lino/tx/validator"
	"github.com/lino-network/lino/tx/vote"
	types "github.com/lino-network/lino/types"
)

type DecideProposalEvent struct{}

func (dpe DecideProposalEvent) Execute(
	ctx sdk.Context, voteManager vote.VoteManager, valManager val.ValidatorManager,
	am acc.AccountManager, pm ProposalManager, gm global.GlobalManager) sdk.Error {
	// get the proposal ID we are going to decide
	curID, err := pm.GetCurrentProposal(ctx)
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
	proposalRes, err := pm.UpdateProposalStatus(ctx, votingRes)
	if err != nil {
		return err
	}

	// majority disagree this proposal
	if proposalRes == types.ProposalNotPass {
		return nil
	}

	// create parameter change event
	if err := pm.CreateParamChangeEvent(ctx, curID, gm); err != nil {
		return err
	}
	return nil
}
