package proposal

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"
)

func NewHandler(am acc.AccountManager, pm ProposalManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ChangeParamMsg:
			return handleChangeParamMsg(ctx, am, pm, gm, msg)
		case ContentCensorshipMsg:
			return handleContentCensorshipMsg(ctx, am, pm, gm, msg)
		case ProtocolUpgradeMsg:
			return handleProtocolUpgradeMsg(ctx, am, pm, gm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized proposal Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleContentCensorshipMsg(
	ctx sdk.Context, am acc.AccountManager, pm ProposalManager, gm global.GlobalManager, msg ContentCensorshipMsg) sdk.Result {
	return sdk.Result{}
}

func handleProtocolUpgradeMsg(
	ctx sdk.Context, am acc.AccountManager, pm ProposalManager, gm global.GlobalManager, msg ProtocolUpgradeMsg) sdk.Result {
	return sdk.Result{}
}
func handleChangeParamMsg(
	ctx sdk.Context, am acc.AccountManager, pm ProposalManager, gm global.GlobalManager, msg ChangeParamMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.GetCreator()) {
		return ErrUsernameNotFound().Result()
	}

	if _, err := pm.AddProposal(ctx, msg.GetCreator(), msg.GetDescription(), gm); err != nil {
		return err.Result()
	}
	//  set a time event to decide the proposal in 7 days
	if err := pm.CreateDecideProposalEvent(ctx, gm); err != nil {
		return err.Result()
	}

	// minus coin from account and return when deciding the proposal
	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err = am.MinusCoin(ctx, msg.GetCreator(), param.TypeBProposalMinDeposit); err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.GetCreator(), gm, am, int64(1),
		param.TypeBProposalDecideHr, param.TypeBProposalMinDeposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm global.GlobalManager, am acc.AccountManager,
	times int64, interval int64, coin types.Coin) sdk.Error {
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

	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time, interval, times); err != nil {
		return err
	}

	if err := gm.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
