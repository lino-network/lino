package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	acc "github.com/lino-network/lino/x/account"
	dev "github.com/lino-network/lino/x/developer"
	global "github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post/model"
	"github.com/lino-network/lino/x/post/types"
	price "github.com/lino-network/lino/x/price"
	rep "github.com/lino-network/lino/x/reputation"
)

type PostManager struct {
	postStorage model.PostStorage

	// deps
	am    acc.AccountKeeper
	gm    global.GlobalKeeper
	dev   dev.DeveloperKeeper
	rep   rep.ReputationKeeper
	price price.PriceKeeper
}

// NewPostManager - create a new post manager
func NewPostManager(key sdk.StoreKey, am acc.AccountKeeper, gm global.GlobalKeeper, dev dev.DeveloperKeeper, rep rep.ReputationKeeper, price price.PriceKeeper) PostManager {
	return PostManager{
		postStorage: model.NewPostStorage(key),
		am:          am,
		gm:          gm,
		dev:         dev,
		rep:         rep,
		price:       price,
	}
}

// DoesPostExist - check if post exist
func (pm PostManager) DoesPostExist(ctx sdk.Context, permlink linotypes.Permlink) bool {
	return pm.postStorage.HasPost(ctx, permlink)
}

// GetPost - return post.
func (pm PostManager) GetPost(ctx sdk.Context, permlink linotypes.Permlink) (model.Post, sdk.Error) {
	post, err := pm.postStorage.GetPost(ctx, permlink)
	if err != nil {
		return model.Post{}, err
	}
	return *post, nil
}

// CreatePost validate and handles CreatePostMsg
// stateful validation;
// 1. both author and post id exists.
// 2. if createdBy is not author, then it must be an app.
// 3. post does not exists.
func (pm PostManager) CreatePost(ctx sdk.Context, author linotypes.AccountKey, postID string, createdBy linotypes.AccountKey, content string, title string) sdk.Error {
	if !pm.am.DoesAccountExist(ctx, author) {
		return types.ErrAccountNotFound(author)
	}
	if !pm.am.DoesAccountExist(ctx, createdBy) {
		return types.ErrAccountNotFound(createdBy)
	}
	// if created by app, then app must exist.
	if author != createdBy && !pm.dev.DoesDeveloperExist(ctx, createdBy) {
		return types.ErrDeveloperNotFound(createdBy)
	}
	permlink := linotypes.GetPermlink(author, postID)
	if pm.DoesPostExist(ctx, permlink) {
		return types.ErrPostAlreadyExist(permlink)
	}

	createdAt := ctx.BlockHeader().Time.Unix()
	postInfo := &model.Post{
		PostID:    postID,
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedBy: createdBy,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	pm.postStorage.SetPost(ctx, postInfo)
	return nil
}

// UpdatePost - update post title, content and links.
// stateful validation:
// 1. author exist.
// 2. post exist.
func (pm PostManager) UpdatePost(ctx sdk.Context, author linotypes.AccountKey, postID, title, content string) sdk.Error {
	if !pm.am.DoesAccountExist(ctx, author) {
		return types.ErrAccountNotFound(author)
	}
	permlink := linotypes.GetPermlink(author, postID)
	postInfo, err := pm.postStorage.GetPost(ctx, permlink)
	if err != nil {
		// post not exists
		return err
	}
	postInfo.Title = title
	postInfo.Content = content
	postInfo.UpdatedAt = ctx.BlockHeader().Time.Unix()
	pm.postStorage.SetPost(ctx, postInfo)
	return nil
}

// DeletePost - delete post by author or content censorship
// stateful validation:
// 1. author exists.
// 2. permlink exists.
func (pm PostManager) DeletePost(ctx sdk.Context, permlink linotypes.Permlink) sdk.Error {
	if !pm.DoesPostExist(ctx, permlink) {
		return types.ErrPostNotFound(permlink)
	}
	pm.postStorage.DeletePost(ctx, permlink)
	return nil
}

// LinoDonate handles donation using lino.
// stateful validation:
// 1. post exits
// 2. from/to account exists.
// 3. no self donation.
// 4. if app is not empty, then developer must exist.
// 5. amount positive > 0.
func (pm PostManager) LinoDonate(ctx sdk.Context, from linotypes.AccountKey, amount linotypes.Coin, author linotypes.AccountKey, postID string, app linotypes.AccountKey) sdk.Error {
	if err := pm.validateLinoDonation(ctx, from, amount, author, postID, app); err != nil {
		return err
	}
	// donation.
	permlink := linotypes.GetPermlink(author, postID)
	rate, err := pm.gm.GetConsumptionFrictionRate(ctx)
	if err != nil {
		return err
	}
	frictionCoin := linotypes.DecToCoin(amount.ToDec().Mul(rate))
	// dp is the evaluated result.
	dp, err := pm.rep.DonateAt(ctx, from, permlink, pm.price.CoinToMiniDollar(amount))
	if err != nil {
		return err
	}
	rewardEvent := RewardEvent{
		PostAuthor: author,
		PostID:     postID,
		Consumer:   from,
		Evaluate:   dp,
		FromApp:    app,
	}
	if err := pm.gm.AddFrictionAndRegisterContentRewardEvent(
		ctx, rewardEvent, frictionCoin, dp); err != nil {
		return err
	}
	// memo is deprecated.
	_, err = pm.am.MinusSavingCoinWithFullCoinDay(
		ctx, from, amount, author, "", linotypes.DonationOut)
	if err != nil {
		return err
	}
	directDeposit := amount.Minus(frictionCoin)
	if err := pm.am.AddSavingCoin(
		ctx, author, directDeposit, from, "", linotypes.DonationIn); err != nil {
		return err
	}
	return nil
}

// IDADonate - handle IDA donation.
func (pm PostManager) IDADonate(ctx sdk.Context, from linotypes.AccountKey, n linotypes.IDA, author linotypes.AccountKey, postID string, app linotypes.AccountKey) sdk.Error {
	if err := pm.validateIDADonate(ctx, from, n, author, postID, app); err != nil {
		return err
	}
	permlink := linotypes.GetPermlink(author, postID)
	idaPrice, err := pm.dev.GetIDAPrice(app)
	if err != nil {
		return err
	}
	dollarAmount := linotypes.IDAToMiniDollar(n, idaPrice) // unit conversion

	rate, err := pm.gm.GetConsumptionFrictionRate(ctx)
	if err != nil {
		return err
	}
	tax := linotypes.NewMiniDollarFromInt(dollarAmount.ToDec().Mul(rate).TruncateInt())
	dollarTransfer := linotypes.NewMiniDollarFromInt(dollarAmount.Sub(tax.Int))

	// dp is the evaluated result.
	dp, err := pm.rep.DonateAt(ctx, from, permlink, dollarAmount)
	if err != nil {
		return err
	}
	rewardEvent := RewardEvent{
		PostAuthor: author,
		PostID:     postID,
		Consumer:   from,
		Evaluate:   dp,
		FromApp:    app,
	}
	if err := pm.gm.AddFrictionAndRegisterContentRewardEvent(
		ctx, rewardEvent, pm.price.MiniDollarToCoin(tax), dp); err != nil {
		return err
	}
	if err := pm.dev.MoveIDA(app, from, author, dollarTransfer); err != nil {
		return err
	}
	return nil
}

// donation stateful basic validation:
// 1. post exits
// 2. from/to account exists.
// 3. no self donation.
func (pm PostManager) validateDonationBasic(ctx sdk.Context, from linotypes.AccountKey, author linotypes.AccountKey, postID string) sdk.Error {
	if from == author {
		return types.ErrCannotDonateToSelf(from)
	}
	if !pm.am.DoesAccountExist(ctx, from) {
		return types.ErrAccountNotFound(from)
	}
	if !pm.am.DoesAccountExist(ctx, author) {
		return types.ErrAccountNotFound(author)
	}
	permlink := linotypes.GetPermlink(author, postID)
	if !pm.DoesPostExist(ctx, permlink) {
		return types.ErrPostNotFound(permlink)
	}
	return nil
}

// lino donation stateful.
// 1. basic validation
// 2. lino amount > 0.
// 3. if app is not empty, then developer must exist.
func (pm PostManager) validateLinoDonation(ctx sdk.Context, from linotypes.AccountKey, amount linotypes.Coin, author linotypes.AccountKey, postID string, app linotypes.AccountKey) sdk.Error {
	err := pm.validateDonationBasic(ctx, from, author, postID)
	if err != nil {
		return err
	}
	if app != "" && !pm.dev.DoesDeveloperExist(ctx, app) {
		return types.ErrDeveloperNotFound(app)
	}
	if !amount.IsPositive() {
		return types.ErrInvalidDonationAmount(amount)
	}
	return nil
}

// IDA donation stateful.
// 1. basic validation
// 2. lino amount > 0.
// 3. app cannot be empty and the developer must exist.
func (pm PostManager) validateIDADonate(ctx sdk.Context, from linotypes.AccountKey, n linotypes.IDA, author linotypes.AccountKey, postID string, app linotypes.AccountKey) sdk.Error {
	err := pm.validateDonationBasic(ctx, from, author, postID)
	if err != nil {
		return err
	}
	if app == "" || !pm.dev.DoesDeveloperExist(ctx, app) {
		return types.ErrDeveloperNotFound(app)
	}
	if !n.IsPositive() {
		return types.ErrNonPositiveIDAAmount(n)
	}
	return nil
}

// Export - adaptor to storage Export.
func (pm PostManager) ExportToFile(ctx sdk.Context, filepath string) error {
	panic("post export unimplemented")
}

// Import - adaptor to storage Export.
func (pm PostManager) ImportFromFile(ctx sdk.Context, filepath string) error {
	rst, err := utils.Load(filepath, func() interface{} { return &model.PostTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.PostTablesIR)
	ctx.Logger().Info("%s state parsed\n", filepath)

	// upgrade2 has simplied the post structure to just one post.
	for _, v := range table.Posts {
		pm.postStorage.SetPost(ctx, &model.Post{
			PostID:    v.Info.PostID,
			Title:     v.Info.Title,
			Content:   v.Info.Content,
			Author:    v.Info.Author,
			CreatedBy: v.Info.Author,
			CreatedAt: v.Meta.CreatedAt,
			UpdatedAt: v.Meta.LastUpdatedAt,
		})
	}
	ctx.Logger().Info("%s state imported\n", filepath)
	return nil
}
