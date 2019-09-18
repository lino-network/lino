package proposal

import (
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post"
	"github.com/lino-network/lino/x/vote"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	val "github.com/lino-network/lino/x/validator"
)

// DecideProposalEvent - a 7 days event to determine the result and status of ongoing proposal
type DecideProposalEvent struct {
	ProposalType types.ProposalType `json:"proposal_type"`
	ProposalID   types.ProposalKey  `json:"proposal_id"`
}

// Execute - execute proposal event, check vote and update status
func (dpe DecideProposalEvent) Execute(
	ctx sdk.Context, vk vote.VoteKeeper, valManager val.ValidatorManager,
	am acc.AccountKeeper, proposalManager ProposalManager, postManager post.PostKeeper,
	gm *global.GlobalManager) sdk.Error {
	// check it is ongoing proposal
	if !proposalManager.IsOngoingProposal(ctx, dpe.ProposalID) {
		return ErrOngoingProposalNotFound()
	}

	// update the ongoing and past proposal list
	proposalRes, err := proposalManager.UpdateProposalPassStatus(
		ctx, dpe.ProposalType, dpe.ProposalID)
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
		if err := dpe.ExecuteChangeParam(ctx, dpe.ProposalID, proposalManager, gm); err != nil {
			return err
		}
	case types.ContentCensorship:
		if err := dpe.ExecuteContentCensorship(ctx, dpe.ProposalID, proposalManager, postManager); err != nil {
			return err
		}
	case types.ProtocolUpgrade:
		if err := dpe.ExecuteProtocolUpgrade(ctx, dpe.ProposalID, proposalManager); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteChangeParam - reigster parameter change event
func (dpe DecideProposalEvent) ExecuteChangeParam(
	ctx sdk.Context, curID types.ProposalKey, proposalManager ProposalManager,
	gm *global.GlobalManager) sdk.Error {
	event, err := proposalManager.CreateParamChangeEvent(ctx, curID)
	if err != nil {
		return err
	}
	if err := gm.RegisterEventAtTime(ctx,
		ctx.BlockHeader().Time.Unix()+types.ParamChangeTimeout, event); err != nil {
		return err
	}
	return nil
}

// ExecuteContentCensorship - delete target post
func (dpe DecideProposalEvent) ExecuteContentCensorship(
	ctx sdk.Context, curID types.ProposalKey, proposalManager ProposalManager,
	postManager post.PostKeeper) sdk.Error {
	permlink, err := proposalManager.GetPermlink(ctx, curID)
	if err != nil {
		return err
	}

	// TODO add content censorship logic
	if exist := postManager.DoesPostExist(ctx, permlink); !exist {
		return ErrCensorshipPostNotFound()
	}
	if err := postManager.DeletePost(ctx, permlink); err != nil {
		return err
	}
	return nil
}

// ExecuteProtocolUpgrade - since execute protocol upgrade engage code change, the process need to be done manually
func (dpe DecideProposalEvent) ExecuteProtocolUpgrade(
	ctx sdk.Context, curID types.ProposalKey, proposalManager ProposalManager) sdk.Error {
	return nil
}
