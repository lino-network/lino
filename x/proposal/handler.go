package proposal

import (
	"fmt"
	"reflect"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post"
	"github.com/lino-network/lino/x/vote"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
)

// NewHandler - Handle all "proposal" type messages.
func NewHandler(
	am acc.AccountManager, proposalManager ProposalManager,
	postManager post.PostKeeper, gm *global.GlobalManager, vm vote.VoteManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ChangeParamMsg:
			return handleChangeParamMsg(ctx, am, proposalManager, gm, msg)
		case ContentCensorshipMsg:
			return handleContentCensorshipMsg(ctx, am, proposalManager, postManager, gm, msg)
		case ProtocolUpgradeMsg:
			return handleProtocolUpgradeMsg(ctx, am, proposalManager, gm, msg)
		case VoteProposalMsg:
			return handleVoteProposalMsg(ctx, proposalManager, vm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized proposal Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleChangeParamMsg(
	ctx sdk.Context, am acc.AccountManager, pm ProposalManager, gm *global.GlobalManager,
	msg ChangeParamMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.GetCreator()) {
		return ErrAccountNotFound().Result()
	}

	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err.Result()
	}

	proposal := pm.CreateChangeParamProposal(ctx, msg.GetParameter(), msg.GetReason())
	proposalID, err := pm.AddProposal(ctx, msg.GetCreator(), proposal, param.ChangeParamDecideSec)
	if err != nil {
		return err.Result()
	}
	//  set a time event to decide the proposal
	event := pm.CreateDecideProposalEvent(ctx, types.ChangeParam, proposalID)

	if err := gm.RegisterProposalDecideEvent(ctx, param.ChangeParamDecideSec, event); err != nil {
		return err.Result()
	}

	// minus coin from account and return when deciding the proposal
	if err = am.MinusSavingCoin(
		ctx, msg.GetCreator(), param.ChangeParamMinDeposit, "",
		string(proposalID), types.ProposalDeposit); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.GetCreator(), gm, am, int64(1),
		param.ChangeParamDecideSec, param.ChangeParamMinDeposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleProtocolUpgradeMsg(
	ctx sdk.Context, am acc.AccountManager, pm ProposalManager, gm *global.GlobalManager,
	msg ProtocolUpgradeMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.GetCreator()) {
		return ErrAccountNotFound().Result()
	}

	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err.Result()
	}

	proposal := pm.CreateProtocolUpgradeProposal(ctx, msg.GetLink(), msg.GetReason())
	proposalID, err := pm.AddProposal(ctx, msg.GetCreator(), proposal, param.ProtocolUpgradeDecideSec)
	if err != nil {
		return err.Result()
	}
	//  set a time event to decide the proposal
	event := pm.CreateDecideProposalEvent(ctx, types.ProtocolUpgrade, proposalID)

	if err := gm.RegisterProposalDecideEvent(ctx, param.ProtocolUpgradeDecideSec, event); err != nil {
		return err.Result()
	}

	// minus coin from account and return when deciding the proposal
	if err = am.MinusSavingCoin(
		ctx, msg.GetCreator(), param.ProtocolUpgradeMinDeposit,
		"", string(proposalID), types.ProposalDeposit); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.GetCreator(), gm, am, int64(1),
		param.ProtocolUpgradeDecideSec, param.ProtocolUpgradeMinDeposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleContentCensorshipMsg(
	ctx sdk.Context, am acc.AccountManager, proposalManager ProposalManager,
	postManager post.PostKeeper, gm *global.GlobalManager, msg ContentCensorshipMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.GetCreator()) {
		return ErrAccountNotFound().Result()
	}

	if !postManager.DoesPostExist(ctx, msg.GetPermlink()) {
		return ErrPostNotFound().Result()
	}

	param, err := proposalManager.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err.Result()
	}

	proposal :=
		proposalManager.CreateContentCensorshipProposal(
			ctx, msg.GetPermlink(), msg.GetReason())
	proposalID, err :=
		proposalManager.AddProposal(
			ctx, msg.GetCreator(), proposal, param.ContentCensorshipDecideSec)
	if err != nil {
		return err.Result()
	}
	//  set a time event to decide the proposal
	event := proposalManager.CreateDecideProposalEvent(ctx, types.ContentCensorship, proposalID)
	// minus coin from account and return when deciding the proposal
	if err = am.MinusSavingCoin(
		ctx, msg.GetCreator(), param.ContentCensorshipMinDeposit,
		"", string(proposalID), types.ProposalDeposit); err != nil {
		return err.Result()
	}

	if err := gm.RegisterProposalDecideEvent(ctx, param.ContentCensorshipDecideSec, event); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.GetCreator(), gm, am, int64(1),
		param.ContentCensorshipDecideSec, param.ContentCensorshipMinDeposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleVoteProposalMsg(ctx sdk.Context, proposalManager ProposalManager, vm vote.VoteManager, msg VoteProposalMsg) sdk.Result {
	if !vm.DoesVoterExist(ctx, msg.Voter) {
		return ErrVoterNotFound().Result()
	}

	if !proposalManager.IsOngoingProposal(ctx, msg.ProposalID) {
		return ErrNotOngoingProposal().Result()
	}

	if err := vm.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Result); err != nil {
		return err.Result()
	}

	v, err := vm.GetVote(ctx, msg.ProposalID, msg.Voter)
	if err != nil {
		return err.Result()
	}

	err = proposalManager.UpdateProposalVotingStatus(ctx, msg.ProposalID, msg.Voter, v.Result, v.VotingPower)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm *global.GlobalManager, am acc.AccountManager,
	times int64, interval int64, coin types.Coin) sdk.Error {
	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time.Unix(), interval, times); err != nil {
		return err
	}

	events, err := acc.CreateCoinReturnEvents(ctx, name, times, interval, coin, types.ProposalReturnCoin)
	if err != nil {
		return err
	}

	if err := gm.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
