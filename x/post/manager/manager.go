package manager

import (
	"fmt"

	codec "github.com/cosmos/cosmos-sdk/codec"
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

const (
	exportVersion = 1
	importVersion = 1
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
// 1. permlink kv exists.
// 2. post is not marked as deleted.
func (pm PostManager) DoesPostExist(ctx sdk.Context, permlink linotypes.Permlink) bool {
	if !pm.postStorage.HasPost(ctx, permlink) {
		return false
	}
	post, _ := pm.postStorage.GetPost(ctx, permlink)
	return !post.IsDeleted
}

// GetPost - return post.
// return err if post is deleted.
func (pm PostManager) GetPost(ctx sdk.Context, permlink linotypes.Permlink) (model.Post, sdk.Error) {
	post, err := pm.postStorage.GetPost(ctx, permlink)
	if err != nil {
		return model.Post{}, err
	}
	if post.IsDeleted {
		return model.Post{}, types.ErrPostDeleted(permlink)
	}
	return *post, nil
}

// CreatePost validate and handles CreatePostMsg
// stateful validation;
// 1. both author and post id exists.
// 2. if createdBy is not author, then it must be an app.
// 3. post's permlink does not exists.
func (pm PostManager) CreatePost(ctx sdk.Context, author linotypes.AccountKey, postID string, createdBy linotypes.AccountKey, content string, title string) sdk.Error {
	if !pm.am.DoesAccountExist(ctx, author) {
		return types.ErrAccountNotFound(author)
	}
	if !pm.am.DoesAccountExist(ctx, createdBy) {
		return types.ErrAccountNotFound(createdBy)
	}
	permlink := linotypes.GetPermlink(author, postID)
	if pm.postStorage.HasPost(ctx, permlink) {
		return types.ErrPostAlreadyExist(permlink)
	}
	if author != createdBy {
		// if created by app, then createdBy must either be the app or an affiliated account of app.
		dev := createdBy
		var err error
		createdBy, err = pm.dev.GetAffiliatingApp(ctx, createdBy)
		if err != nil {
			return types.ErrDeveloperNotFound(dev)
		}
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
// 1. permlink exists.
// 2. post not deleted.
// Delete does not delete the post in kv store, as that will make `permlink` not permanent.
// It is marked as deleted, then on deleted posts,
// 1. manager.DoesPostExist will return false.
// 2. manager.GetPost will return ErrPermlinkDeleted.
// 3. manager.CreatePost will return ErrPostAlreadyExist.
func (pm PostManager) DeletePost(ctx sdk.Context, permlink linotypes.Permlink) sdk.Error {
	post, err := pm.postStorage.GetPost(ctx, permlink)
	if err != nil {
		return err
	}
	if post.IsDeleted {
		return types.ErrPostDeleted(permlink)
	}
	post.IsDeleted = true
	post.Title = ""
	post.Content = ""
	pm.postStorage.SetPost(ctx, post)
	return nil
}

// LinoDonate handles donation using lino.
// stateful validation:
// 1. post exits
// 2. from/to account exists.
// 3. no self donation.
// 4. if app is not empty, then developer must exist.
// 5. amount positive > 0.
// 6. 9.9% of amount > 0 coin.
func (pm PostManager) LinoDonate(ctx sdk.Context, from linotypes.AccountKey, amount linotypes.Coin, author linotypes.AccountKey, postID string, app linotypes.AccountKey) sdk.Error {
	if err := pm.validateLinoDonation(ctx, from, amount, author, postID, app); err != nil {
		return err
	}
	// donation.
	rate := sdk.MustNewDecFromStr(linotypes.ConsumptionFrictionRate)
	frictionCoin := linotypes.DecToCoin(amount.ToDec().Mul(rate))
	if frictionCoin.IsZero() {
		return types.ErrDonateAmountTooLittle()
	}
	// friction goes to the friction pool for voters.
	err := pm.am.MoveToPool(ctx,
		linotypes.VoteFrictionPool, linotypes.NewAccOrAddrFromAcc(from), frictionCoin)
	if err != nil {
		return err
	}

	// rest goes to the author.
	err = pm.am.MoveCoin(ctx, linotypes.NewAccOrAddrFromAcc(from),
		linotypes.NewAccOrAddrFromAcc(author), amount.Minus(frictionCoin))
	if err != nil {
		return err
	}

	mdamount, err := pm.price.CoinToMiniDollar(ctx, amount)
	if err != nil {
		return err
	}
	return pm.afterDonation(ctx, author, postID, from, mdamount, frictionCoin, app)
}

// IDADonate - handle IDA donation.
func (pm PostManager) IDADonate(ctx sdk.Context, from linotypes.AccountKey, n linotypes.MiniIDA, author linotypes.AccountKey, postID string, app, signer linotypes.AccountKey) sdk.Error {
	if err := pm.validateIDADonate(ctx, from, n, author, postID, app); err != nil {
		return err
	}
	signerApp, err := pm.dev.GetAffiliatingApp(ctx, signer)
	if err != nil || signerApp != app {
		return types.ErrInvalidSigner()
	}
	idaPrice, err := pm.dev.GetMiniIDAPrice(ctx, app)
	if err != nil {
		return err
	}

	// amount = tax + dollarTransfer
	// tax: burned to lino
	// dollarTransfer: moved from sender to receipient.
	rate := sdk.MustNewDecFromStr(linotypes.ConsumptionFrictionRate)
	dollarAmount := linotypes.MiniIDAToMiniDollar(n, idaPrice) // unit conversion
	tax := linotypes.NewMiniDollarFromInt(dollarAmount.ToDec().Mul(rate).TruncateInt())

	// burn and check taxable coins.
	// tax will be subtracted from @p from's IDA account, and converted to coins and
	// saved in the account.
	taxcoins, err := pm.dev.BurnIDA(ctx, app, from, tax)
	if err != nil {
		return err
	}
	if !taxcoins.IsPositive() {
		return types.ErrDonateAmountTooLittle()
	}

	// friction goes to the friction pool for voters.
	err = pm.am.MoveToPool(ctx,
		linotypes.VoteFrictionPool, linotypes.NewAccOrAddrFromAcc(from), taxcoins)
	if err != nil {
		return err
	}

	// rest goes to the author
	if err := pm.dev.MoveIDA(ctx, app, from, author, dollarAmount.Minus(tax)); err != nil {
		return err
	}

	return pm.afterDonation(ctx, author, postID, from, dollarAmount, taxcoins, app)
}

func (pm PostManager) afterDonation(ctx sdk.Context, author linotypes.AccountKey, postID string, from linotypes.AccountKey, damount linotypes.MiniDollar, friction linotypes.Coin, app linotypes.AccountKey) sdk.Error {
	// impact is the evaluated consumption.
	impact, err := pm.rep.DonateAt(ctx, from, linotypes.GetPermlink(author, postID), damount)
	if err != nil {
		return err
	}

	// update consumptionm window
	consumptionWindow := pm.postStorage.GetConsumptionWindow(ctx)
	pm.postStorage.SetConsumptionWindow(ctx, consumptionWindow.Plus(impact))

	// record friction stats.
	err = pm.gm.RecordFriction(ctx, friction)
	if err != nil {
		return err
	}

	// add content bonus return event.
	rewardEvent := types.RewardEvent{
		PostAuthor: author,
		PostID:     postID,
		Consumer:   from,
		Evaluate:   impact,
		FromApp:    app,
	}
	eventTime := ctx.BlockHeader().Time.Unix() + linotypes.ConsumptionFreezingPeriodSec
	if err := pm.gm.RegisterEventAtTime(ctx, eventTime, rewardEvent); err != nil {
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
func (pm PostManager) validateIDADonate(ctx sdk.Context, from linotypes.AccountKey, n linotypes.MiniIDA, author linotypes.AccountKey, postID string, app linotypes.AccountKey) sdk.Error {
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

// ExecRewardEvent - execute reward events.
func (pm PostManager) ExecRewardEvent(ctx sdk.Context, event types.RewardEvent) sdk.Error {
	// check if post is deleted, Note that if post is deleted, it's ok to just
	// skip this event. It does not return an error because errors will panic in events.
	permlink := linotypes.GetPermlink(event.PostAuthor, event.PostID)
	if !pm.DoesPostExist(ctx, permlink) {
		return nil
	}

	// if developer exist, add to developer consumption
	if pm.dev.DoesDeveloperExist(ctx, event.FromApp) {
		// ignore report consumption err.
		_ = pm.dev.ReportConsumption(ctx, event.FromApp, event.Evaluate)
	}

	return pm.allocContentBonus(ctx, event.Evaluate, event.PostAuthor)
}

func (pm PostManager) allocContentBonus(ctx sdk.Context, impact linotypes.MiniDollar, author linotypes.AccountKey) sdk.Error {
	if impact.IsZero() {
		return nil
	}

	// get consumption window and update the window, this mutation is the reason why
	// this function is called "pop***"
	consumptionWindow := pm.postStorage.GetConsumptionWindow(ctx)
	pm.postStorage.SetConsumptionWindow(ctx, consumptionWindow.Minus(impact))

	// XXX(yumin): the ratio is zero when the window is zero, because the consumption window
	// as the sum of past donation impacts, shall be large than zero as the @p impact is nonzero.
	// consumptionRatio = (this consumption * penalty score) / (total consumption in 7 days window)
	consumptionRatio := sdk.ZeroDec()
	if !consumptionWindow.ToDec().IsZero() {
		consumptionRatio = impact.ToDec().Quo(consumptionWindow.ToDec())
	}
	rewardPool, err := pm.am.GetPool(ctx, linotypes.InflationConsumptionPool)
	if err != nil {
		return err
	}
	// reward = (consumption reward pool) * (consumptionRatio)
	reward := types.DecToCoin(rewardPool.ToDec().Mul(consumptionRatio))
	return pm.am.MoveFromPool(ctx, linotypes.InflationConsumptionPool, author, reward)
}

// Export - to file.
func (pm PostManager) ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	state := &model.PostTablesIR{
		Version: exportVersion,
	}
	storeList := pm.postStorage.StoreMap(ctx)

	// export posts
	posts := make([]model.PostIR, 0)
	postSubStore := storeList[string(model.PostSubStore)]
	postSubStore.Iterate(func(key []byte, val interface{}) bool {
		post := val.(*model.Post)
		posts = append(posts, model.PostIR(*post))
		return false
	})
	state.Posts = posts

	return utils.Save(filepath, cdc, state)
}

// Import - from file
func (pm PostManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.PostTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.PostTablesIR)

	if table.Version != importVersion {
		return fmt.Errorf("unsupported import version: %d", table.Version)
	}

	for _, v := range table.Posts {
		pm.postStorage.SetPost(ctx, &model.Post{
			PostID:    v.PostID,
			Title:     v.Title,
			Content:   v.Content,
			Author:    v.Author,
			CreatedBy: v.Author,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			IsDeleted: v.IsDeleted,
		})
	}
	return nil
}
