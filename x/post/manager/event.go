package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	// "github.com/lino-network/lino/x/post/types"
)

// RewardEvent - when donation occurred, a reward event will be register
// at 7 days later. After 7 days reward event will be executed and send
// inflation to author.
type RewardEvent struct {
	PostAuthor linotypes.AccountKey `json:"post_author"`
	PostID     string               `json:"post_id"`
	Consumer   linotypes.AccountKey `json:"consumer"`
	Evaluate   linotypes.MiniDollar `json:"evaluate"`
	FromApp    linotypes.AccountKey `json:"from_app"`
}

// Execute - execute reward event after 7 days
func (event RewardEvent) Execute(ctx sdk.Context, pm PostManager) sdk.Error {
	// check if post is deleted, Note that if post is deleted, it's ok to just
	// skip this event. It does not return an error because errors will panic in events.
	permlink := linotypes.GetPermlink(event.PostAuthor, event.PostID)
	if !pm.DoesPostExist(ctx, permlink) {
		return nil
	}

	// pop out rewards
	reward, err := pm.gm.GetRewardAndPopFromWindow(ctx, event.Evaluate)
	if err != nil {
		return err
	}
	// if developer exist, add to developer consumption
	if pm.dev.DoesDeveloperExist(ctx, event.FromApp) {
		// ignore report consumption err.
		_ = pm.dev.ReportConsumption(ctx, event.FromApp, event.Evaluate)
	}

	// previsously rewards were added to account's reward, now it's added directly to balance.
	err = pm.am.AddCoinToUsername(ctx, event.PostAuthor, reward)
	return err
}
