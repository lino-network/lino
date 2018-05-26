package post

import (
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	dev "github.com/lino-network/lino/tx/developer"
)

func init() {
	cdc := wire.NewCodec()

	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(RewardEvent{}, "event/reward", nil)
}

type RewardEvent struct {
	PostAuthor types.AccountKey `json:"post_author"`
	PostID     string           `json:"post_id"`
	Consumer   types.AccountKey `json:"consumer"`
	Evaluate   types.Coin       `json:"evaluate"`
	Original   types.Coin       `json:"original"`
	Friction   types.Coin       `json:"friction"`
	FromApp    types.AccountKey `json:"from_app"`
}

func (event RewardEvent) Execute(
	ctx sdk.Context, pm PostManager, am acc.AccountManager,
	gm global.GlobalManager, dm dev.DeveloperManager) sdk.Error {

	permLink := types.GetPermLink(event.PostAuthor, event.PostID)
	paneltyScore, err := pm.GetPenaltyScore(ctx, permLink)
	if err != nil {
		return err
	}
	reward, err := gm.GetRewardAndPopFromWindow(ctx, event.Evaluate, paneltyScore)
	if err != nil {
		return err
	}
	if dm.IsDeveloperExist(ctx, event.FromApp) {
		dm.ReportConsumption(ctx, event.FromApp, reward)
	}
	if !am.IsAccountExist(ctx, event.PostAuthor) {
		return acc.ErrUsernameNotFound()
	}
	if !pm.IsPostExist(ctx, permLink) {
		return ErrDonatePostNotFound(permLink)
	}
	if err := pm.AddDonation(ctx, permLink, event.Consumer, reward, types.Inflation); err != nil {
		return err
	}
	if err := am.AddIncomeAndReward(
		ctx, event.PostAuthor, event.Original, event.Friction, reward); err != nil {
		return err
	}
	return nil
}
