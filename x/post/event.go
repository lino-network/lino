package post

import (
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"
	rep "github.com/lino-network/lino/x/reputation"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	dev "github.com/lino-network/lino/x/developer"
	vote "github.com/lino-network/lino/x/vote"
)

func init() {
	cdc := wire.NewCodec()

	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(RewardEvent{}, "event/reward", nil)
}

// RewardEvent - when donation occurred, a reward event will be register
// at 7 days later. After 7 days reward event will be executed and send
// inflation to author.
type RewardEvent struct {
	PostAuthor types.AccountKey `json:"post_author"`
	PostID     string           `json:"post_id"`
	Consumer   types.AccountKey `json:"consumer"`
	Evaluate   types.Coin       `json:"evaluate"`
	Original   types.Coin       `json:"original"`
	Friction   types.Coin       `json:"friction"`
	FromApp    types.AccountKey `json:"from_app"`
}

// Execute - execute reward event after 7 days
func (event RewardEvent) Execute(
	ctx sdk.Context, pm PostManager, am acc.AccountManager,
	gm global.GlobalManager, dm dev.DeveloperManager,
	vm vote.VoteManager, rm rep.ReputationManager) sdk.Error {

	permlink := types.GetPermlink(event.PostAuthor, event.PostID)
	paneltyScore, err := pm.GetPenaltyScore(ctx, permlink)
	if err != nil {
		return err
	}
	// check if post is deleted
	if isDeleted, err := pm.IsDeleted(ctx, permlink); isDeleted || err != nil {
		paneltyScore = sdk.OneRat()
	}
	reward, err := gm.GetRewardAndPopFromWindow(ctx, event.Evaluate, paneltyScore)
	if err != nil {
		return err
	}
	// if developer exist, add to developer consumption
	if dm.DoesDeveloperExist(ctx, event.FromApp) {
		dm.ReportConsumption(ctx, event.FromApp, reward)
	}
	if !am.DoesAccountExist(ctx, event.PostAuthor) {
		return ErrAccountNotFound(event.PostAuthor)
	}
	if !pm.DoesPostExist(ctx, permlink) {
		return ErrPostNotFound(permlink)
	}
	// add donation information to post
	if err := pm.AddDonation(ctx, permlink, event.Consumer, reward, types.Inflation); err != nil {
		return err
	}
	// add half reward to user
	addToReward := types.RatToCoin(reward.ToRat().Mul(sdk.NewRat(1, 2)))
	addToStake := reward.Minus(addToReward)
	if err := am.AddIncomeAndReward(
		ctx, event.PostAuthor, event.Original, event.Friction, addToReward, event.Consumer, event.PostAuthor, event.PostID); err != nil {
		return err
	}
	if !addToStake.IsZero() {
		if err := vote.AddStake(ctx, event.PostAuthor, addToStake, vm, gm, am, rm); err != nil {
			return err
		}
	}
	return nil
}
