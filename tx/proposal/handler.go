package proposal

import (
	"fmt"
	"reflect"

	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
)

func NewHandler(
	am acc.AccountManager, proposalManager ProposalManager,
	postManager post.PostManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ChangeParamMsg:
			return handleChangeParamMsg(ctx, am, proposalManager, gm, msg)
		case ContentCensorshipMsg:
			return handleContentCensorshipMsg(ctx, am, proposalManager, postManager, gm, msg)
		case ProtocolUpgradeMsg:
			return handleProtocolUpgradeMsg(ctx, am, proposalManager, gm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized proposal Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleChangeParamMsg(
	ctx sdk.Context, am acc.AccountManager, pm ProposalManager, gm global.GlobalManager,
	msg ChangeParamMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.GetCreator()) {
		return ErrUsernameNotFound().Result()
	}

	proposal := pm.CreateChangeParamProposal(ctx, msg.GetParameter())
	proposalID, err := pm.AddProposal(ctx, msg.GetCreator(), proposal)
	if err != nil {
		return err.Result()
	}
	//  set a time event to decide the proposal
	event, err := pm.CreateDecideProposalEvent(ctx, types.ChangeParam, proposalID)
	if err != nil {
		return err.Result()
	}

	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := gm.RegisterProposalDecideEvent(ctx, param.ChangeParamDecideHr, event); err != nil {
		return err.Result()
	}

	// minus coin from account and return when deciding the proposal
	if err = am.MinusSavingCoin(ctx, msg.GetCreator(), param.ChangeParamMinDeposit); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.GetCreator(), gm, am, int64(1),
		param.ChangeParamDecideHr, param.ChangeParamMinDeposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleProtocolUpgradeMsg(
	ctx sdk.Context, am acc.AccountManager, pm ProposalManager, gm global.GlobalManager,
	msg ProtocolUpgradeMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.GetCreator()) {
		return ErrUsernameNotFound().Result()
	}

	proposal := pm.CreateProtocolUpgradeProposal(ctx, msg.GetLink())
	proposalID, err := pm.AddProposal(ctx, msg.GetCreator(), proposal)
	if err != nil {
		return err.Result()
	}
	//  set a time event to decide the proposal
	event, err := pm.CreateDecideProposalEvent(ctx, types.ProtocolUpgrade, proposalID)
	if err != nil {
		return err.Result()
	}

	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := gm.RegisterProposalDecideEvent(ctx, param.ProtocolUpgradeDecideHr, event); err != nil {
		return err.Result()
	}

	// minus coin from account and return when deciding the proposal
	if err = am.MinusSavingCoin(ctx, msg.GetCreator(), param.ProtocolUpgradeMinDeposit); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.GetCreator(), gm, am, int64(1),
		param.ProtocolUpgradeDecideHr, param.ProtocolUpgradeMinDeposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleContentCensorshipMsg(
	ctx sdk.Context, am acc.AccountManager, proposalManager ProposalManager,
	postManager post.PostManager, gm global.GlobalManager, msg ContentCensorshipMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.GetCreator()) {
		return ErrUsernameNotFound().Result()
	}

	if !postManager.IsPostExist(ctx, msg.GetPermLink()) {
		return ErrPostNotFound().Result()
	}

	if isDeleted, err := postManager.IsDeleted(ctx, msg.GetPermLink()); isDeleted || err != nil {
		return ErrCensorshipPostIsDeleted(msg.GetPermLink()).Result()
	}

	proposal := proposalManager.CreateContentCensorshipProposal(ctx, msg.GetPermLink())
	proposalID, err := proposalManager.AddProposal(ctx, msg.GetCreator(), proposal)
	if err != nil {
		return err.Result()
	}
	//  set a time event to decide the proposal
	event, err := proposalManager.CreateDecideProposalEvent(ctx, types.ContentCensorship, proposalID)
	if err != nil {
		return err.Result()
	}

	param, err := proposalManager.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err.Result()
	}

	// minus coin from account and return when deciding the proposal
	if err = am.MinusSavingCoin(ctx, msg.GetCreator(), param.ContentCensorshipMinDeposit); err != nil {
		return err.Result()
	}

	if err := gm.RegisterProposalDecideEvent(ctx, param.ContentCensorshipDecideHr, event); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.GetCreator(), gm, am, int64(1),
		param.ContentCensorshipDecideHr, param.ContentCensorshipMinDeposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm global.GlobalManager, am acc.AccountManager,
	times int64, interval int64, coin types.Coin) sdk.Error {
	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time, interval, times); err != nil {
		return err
	}

	events := []types.Event{}
	for i := int64(0); i < times; i++ {
		pieceRat := coin.ToRat().Quo(sdk.NewRat(times - i))
		piece := types.RatToCoin(pieceRat)
		coin = coin.Minus(piece)

		event := acc.ReturnCoinEvent{
			Username: name,
			Amount:   piece,
		}
		events = append(events, event)
	}

	if err := gm.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
