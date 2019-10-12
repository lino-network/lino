package proposal

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/x/proposal/types"
)

// NewHandler - Handle all "proposal" type messages.
func NewHandler(pm ProposalKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.ChangeParamMsg:
			return handleChangeParamMsg(ctx, msg, pm)
		case types.ContentCensorshipMsg:
			return handleContentCensorshipMsg(ctx, msg, pm)
		case types.ProtocolUpgradeMsg:
			return handleProtocolUpgradeMsg(ctx, msg, pm)
		// case types.VoteProposalMsg:
		// 	return handleVoteProposalMsg(ctx, msg, pm)
		default:
			errMsg := fmt.Sprintf("Unrecognized proposal Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleChangeParamMsg(ctx sdk.Context, msg types.ChangeParamMsg, pm ProposalKeeper) sdk.Result {
	if err := pm.ChangeParam(ctx, msg.GetCreator(), msg.GetReason(), msg.GetParameter()); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleProtocolUpgradeMsg(ctx sdk.Context, msg types.ProtocolUpgradeMsg, pm ProposalKeeper) sdk.Result {
	if err := pm.ProtocolUpgrade(ctx, msg.GetCreator(), msg.GetReason(), msg.GetLink()); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleContentCensorshipMsg(ctx sdk.Context, msg types.ContentCensorshipMsg, pm ProposalKeeper) sdk.Result {
	if err := pm.ContentCensorship(ctx, msg.GetCreator(), msg.GetReason(), msg.GetPermlink()); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// func handleVoteProposalMsg(ctx sdk.Context, proposalManager ProposalManager, vm vote.VoteKeeper, msg VoteProposalMsg) sdk.Result {
// 	if !vm.DoesVoterExist(ctx, msg.Voter) {
// 		return ErrVoterNotFound().Result()
// 	}

// 	if !proposalManager.IsOngoingProposal(ctx, msg.ProposalID) {
// 		return ErrNotOngoingProposal().Result()
// 	}

// 	// if err := vm.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Result); err != nil {
// 	// 	return err.Result()
// 	// }

// 	// v, err := vm.GetVote(ctx, msg.ProposalID, msg.Voter)
// 	// if err != nil {
// 	// 	return err.Result()
// 	// }

// 	// err = proposalManager.UpdateProposalVotingStatus(ctx, msg.ProposalID, msg.Voter, v.Result, v.VotingPower)
// 	// if err != nil {
// 	// 	return err.Result()
// 	// }

// 	return sdk.Result{}
// }
