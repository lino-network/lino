package manager

import (
	"sort"
	"time"

	// codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/price/model"
	"github.com/lino-network/lino/x/price/types"
	// "github.com/lino-network/lino/x/validator"
)

const (
	OneHour = 1 * time.Hour
	// Make it a param.
	WeightPeriod = 71 // price is median from last 71 hours.
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
	val  FakeValidator
	vote FakeVote
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
		UpdateAt: ctx.BlockTime(),
	}
	wm.store.SetCurrentPrice(ctx, priceTime)
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
func (wm WeightedMedianPriceManager) UpdateHourlyPrice(ctx sdk.Context) sdk.Error {
	defer wm.store.SetLastValidators(ctx, wm.val.GetValidators())
	wvals := wm.getWeightedValidators(ctx)
	if len(wvals) == 0 {
		return types.ErrNoValidator()
	}
	blocktime := ctx.BlockTime()
	wvals = wm.filterAndSlash(ctx, wvals)
	if len(wvals) == 0 {
		// no valid price this hour, use the same price from last hour.
		curr, err := wm.store.GetCurrentPrice(ctx)
		if err != nil {
			// as long as genesis was inited correctly, curr price should never
			// return error, so panic when err.
			panic(err)
		}
		wm.updateNewPrice(ctx, model.TimePrice{
			Price:    curr.Price,
			UpdateAt: blocktime,
		})
	} else {
		median := wm.calcWeightedMedian(wvals)
		wm.updateNewPrice(ctx, model.TimePrice{
			Price:    median,
			UpdateAt: blocktime,
		})
	}

	return nil
}

// updateNewPrice update history and current price.
// 0. remove the oldest price history entry, if history is full.
// 1. append the price to the price history
// 2. save price history.
// 3. sort price by (price, time)
// 4. set the median as the current price
func (wm WeightedMedianPriceManager) updateNewPrice(ctx sdk.Context, timePrice model.TimePrice) {
	history := wm.store.GetPriceHistory(ctx)
	if len(history)+1 > WeightPeriod {
		history = history[len(history)+1-WeightPeriod:]
	}
	history = append(history, timePrice)
	wm.store.SetPriceHistory(ctx, history)

	// update current price
	sort.SliceStable(history, func(i int, j int) bool {
		left := history[i]
		right := history[j]
		if left.Price.Equal(right.Price) {
			return left.UpdateAt.Unix() < right.UpdateAt.Unix()
		}
		return left.Price.LT(right.Price)
	})
	// when the length is an even number, use lower, indead of (mid + next) / 2.
	mid := len(history) / 2
	current := history[mid]
	wm.store.SetCurrentPrice(ctx, current)
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
		if err != nil || blocktime.Sub(fedPrice.FedTime) > OneHour {
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
