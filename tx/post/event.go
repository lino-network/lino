package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
)

type RewardEvent struct {
	PostAuthor acc.AccountKey `json:"post_author"`
	PostID     string         `json:"post_id"`
	Consumer   acc.AccountKey `json:"consumer"`
	Amount     types.Coin     `json:"amount"`
}

func (event RewardEvent) Execute(ctx sdk.Context, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	globalProxy := global.NewGlobalProxy(&gm)
	authorAccount := acc.NewProxyAccount(event.PostAuthor, &am)
	if !authorAccount.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound()
	}
	post := NewPostProxy(event.PostAuthor, event.PostID, &pm)
	if !post.IsPostExist(ctx) {
		return ErrDonatePostDoesntExist()
	}
	reward, err := globalProxy.GetRewardAndPopFromWindow(ctx, event.Amount)
	if err != nil {
		return err
	}
	if err := post.AddDonation(ctx, event.Consumer, reward); err != nil {
		return err
	}
	if err := authorAccount.AddCoin(ctx, reward); err != nil {
		return err
	}
	if err := authorAccount.Apply(ctx); err != nil {
		return err
	}
	if err := post.Apply(ctx); err != nil {
		return err
	}
	return nil
}
