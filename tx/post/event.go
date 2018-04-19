package post

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

type RewardEvent struct {
	PostAuthor types.AccountKey `json:"post_author"`
	PostID     string           `json:"post_id"`
	Consumer   types.AccountKey `json:"consumer"`
	Evaluate   types.Coin       `json:"evaluate"`
	Original   types.Coin       `json:"original"`
	Friction   types.Coin       `json:"friction"`
}

func (event RewardEvent) Execute(
	ctx sdk.Context, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {

	postKey := types.GetPostKey(event.PostAuthor, event.PostID)
	paneltyScore, err := pm.GetPenaltyScore(ctx, postKey)
	if err != nil {
		return err
	}
	reward, err := gm.GetRewardAndPopFromWindow(ctx, event.Evaluate, paneltyScore)
	if err != nil {
		return err
	}
	if !am.IsAccountExist(ctx, event.PostAuthor) {
		return acc.ErrUsernameNotFound()
	}
	if !pm.IsPostExist(ctx, postKey) {
		return ErrDonatePostDoesntExist(postKey)
	}
	if err := pm.AddDonation(ctx, postKey, event.Consumer, reward); err != nil {
		return err
	}
	fmt.Println("add reward:", reward, event.PostAuthor)
	if err := am.AddIncomeAndReward(
		ctx, event.PostAuthor, event.Original, event.Friction, reward); err != nil {
		return err
	}
	return nil
}
