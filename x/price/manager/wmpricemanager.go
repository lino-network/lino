package manager

import (
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/price/model"
	"github.com/lino-network/lino/x/price/types"
	// "github.com/lino-network/lino/x/validator"
)

type FakeValidator interface {
	GetValidators() []linotypes.AccountKey
	Slash(u linotypes.AccountKey)
}

type FakeVote interface {
	GetVote(u linotypes.AccountKey) linotypes.Coin
}

type WeightedMedianPriceManager struct {
	store model.PriceStorage

	// deps
	param param.ParamHolder
	val   FakeValidator
	vote  FakeVote
}

type weightedValidator struct {
	validator linotypes.AccountKey
	weight    linotypes.Coin
	price     linotypes.MiniDollar
}

// set current price.
func (wm WeightedMedianPriceManager) InitGenesis(ctx sdk.Context, initPrice linotypes.MiniDollar) sdk.Error {
	priceTime := model.TimePrice{
		Price:    initPrice,
		UpdateAt: ctx.BlockTime().Unix(),
	}
	wm.store.SetCurrentPrice(ctx, &priceTime)
	wm.store.SetPriceHistory(ctx, []model.TimePrice{priceTime})
	return nil
}

// UpdateHourlyPrice - update hourly weighted price.
// premise: FedPrice is positive.
// 1. Get Current Validator List, with weight.
// 2. set prices of validators.
// 3. remove invalid.
// 4. get weighted median if at least one validator.
// 5. otherwise, use the previsous price.
func (wm WeightedMedianPriceManager) UpdatePrice(ctx sdk.Context) sdk.Error {
	defer wm.store.SetLastValidators(ctx, wm.val.GetValidators())
	wvals := wm.getWeightedValidators(ctx)
	if len(wvals) == 0 {
		return types.ErrNoValidator()
	}
	blocktime := ctx.BlockTime().Unix()
	wvals = wm.filterAndSlash(ctx, wvals)
	var price linotypes.MiniDollar
	if len(wvals) == 0 {
		// no valid price this hour, use the same price from last hour.
		curr, err := wm.store.GetCurrentPrice(ctx)
		if err != nil {
			// as long as genesis was inited correctly, curr price should never
			// return error, so panic when err.
			panic(err)
		}
		price = curr.Price
	} else {
		price = wm.calcWeightedMedian(wvals)
	}
	wm.updateNewPrice(ctx, model.TimePrice{
		Price:    price,
		UpdateAt: blocktime,
	})
	return nil
}

// FeedPrice - validator update price.
// validation:
// 1. price is positive.
// 2. feeder is a validator.
// 3. can only update after FeedEvery.
func (wm WeightedMedianPriceManager) FeedPrice(ctx sdk.Context, validator linotypes.AccountKey, price linotypes.MiniDollar) sdk.Error {
	if !price.IsPositive() {
		return types.ErrInvalidPriceFeed(price)
	}
	if !wm.isValidator(ctx, validator) {
		return types.ErrNotAValidator(validator)
	}
	blocktime := ctx.BlockTime()
	last, err := wm.store.GetFedPrice(ctx, validator)
	feedEvery := wm.param.GetPriceParam(ctx).FeedEvery
	if err == nil && blocktime.Sub(time.Unix(last.UpdateAt, 0)) < feedEvery {
		return types.ErrPriceFeedRateLimited()
	}

	wm.store.SetFedPrice(ctx, &model.FedPrice{
		Validator: validator,
		Price:     price,
	})
	return nil
}

func (wm WeightedMedianPriceManager) CoinToMiniDollar(ctx sdk.Context, coin linotypes.Coin) (linotypes.MiniDollar, sdk.Error) {
	curr, err := wm.store.GetCurrentPrice(ctx)
	if err != nil {
		return linotypes.NewMiniDollar(0), err
	}
	return coinToMiniDollar(coin, curr.Price), nil
}

func (wm WeightedMedianPriceManager) MiniDollarToCoin(ctx sdk.Context, dollar linotypes.MiniDollar) (linotypes.Coin, linotypes.MiniDollar, sdk.Error) {
	curr, err := wm.store.GetCurrentPrice(ctx)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), linotypes.NewMiniDollar(0), err
	}
	bought, used := miniDollarToCoin(dollar, curr.Price)
	return bought, used, nil
}

func (wm WeightedMedianPriceManager) isValidator(ctx sdk.Context, user linotypes.AccountKey) bool {
	validators := wm.val.GetValidators()
	for _, val := range validators {
		if val == user {
			return true
		}
	}
	return false
}

// updateNewPrice update history and current price.
// 0. remove the oldest price history entry, if history is full.
// 1. append the price to the price history
// 2. save price history.
// 3. sort price by (price, time)
// 4. set the median as the current price
func (wm WeightedMedianPriceManager) updateNewPrice(ctx sdk.Context, timePrice model.TimePrice) {
	history := wm.store.GetPriceHistory(ctx)
	historyMaxLen := wm.param.GetPriceParam(ctx).HistoryMaxLen
	if len(history)+1 > historyMaxLen {
		history = history[len(history)+1-historyMaxLen:]
	}
	history = append(history, timePrice)
	wm.store.SetPriceHistory(ctx, history)

	// update current price
	sort.SliceStable(history, func(i int, j int) bool {
		left := history[i]
		right := history[j]
		if left.Price.Equal(right.Price) {
			return left.UpdateAt < right.UpdateAt
		}
		return left.Price.LT(right.Price)
	})
	// when the length is an even number, use lower, indead of (mid + next) / 2.
	mid := len(history) / 2
	current := history[mid]
	wm.store.SetCurrentPrice(ctx, &current)
}

// getWeightedValidators return weighted validators, sorted by (weight, namestr), increasingly.
// price fields are empty value.
func (wm WeightedMedianPriceManager) getWeightedValidators(ctx sdk.Context) []weightedValidator {
	wvals := make([]weightedValidator, 0)
	vals := wm.val.GetValidators()
	for _, val := range vals {
		wvals = append(wvals, weightedValidator{
			validator: val,
			weight:    wm.vote.GetVote(val),
			price:     linotypes.NewMiniDollar(0),
		})
	}
	return wvals
}

// filterAndSlash slash validators that missed price feeding.
// premise: fedPrice needs to be validated upon validators send update message.
func (wm WeightedMedianPriceManager) filterAndSlash(ctx sdk.Context, wvals []weightedValidator) (rst []weightedValidator) {
	lastValidatorSet := wm.lastRoundValidatorSet(ctx)
	blocktime := ctx.BlockTime()
	for i := range wvals {
		valname := wvals[i].validator
		fedPrice, err := wm.store.GetFedPrice(ctx, valname)
		updateEvery := wm.param.GetPriceParam(ctx).UpdateEvery
		if err != nil || blocktime.Sub(time.Unix(fedPrice.UpdateAt, 0)) > updateEvery {
			// unless the validator is not in the last set, slash.
			if lastValidatorSet[valname] {
				wm.val.Slash(valname)
			}
		} else {
			wvals[i].price = fedPrice.Price
			rst = append(rst, wvals[i])
		}
	}
	return
}

// calcWeightedMedian - return weighted median. pre: len(vals) > 0
func (wm WeightedMedianPriceManager) calcWeightedMedian(wvals []weightedValidator) linotypes.MiniDollar {
	// sort
	sort.Slice(wvals, func(i, j int) bool {
		left := wvals[i]
		right := wvals[j]
		if left.weight.IsEqual(right.weight) {
			return left.validator < right.validator
		}
		return !left.weight.IsGTE(right.weight)
	})

	totalPower := sdk.NewInt(0)
	for _, v := range wvals {
		totalPower = totalPower.Add(v.weight.Amount)
	}
	median := totalPower.QuoRaw(2)
	for _, val := range wvals {
		if median.LT(val.weight.Amount) {
			return val.price
		}
		median = median.Sub(val.weight.Amount)
	}
	// impossible to hit this path.
	return wvals[0].price
}

func (wm WeightedMedianPriceManager) lastRoundValidatorSet(ctx sdk.Context) map[linotypes.AccountKey]bool {
	rst := make(map[linotypes.AccountKey]bool)
	for _, val := range wm.store.GetLastValidators(ctx) {
		rst[val] = true
	}
	return rst
}
