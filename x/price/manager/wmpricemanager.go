package manager

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/price/model"
	"github.com/lino-network/lino/x/price/types"
	"github.com/lino-network/lino/x/validator"
)

type WeightedMedianPriceManager struct {
	store model.PriceStorage

	// deps
	param param.ParamKeeper
	val   validator.ValidatorKeeper
}

func NewWeightedMedianPriceManager(key sdk.StoreKey, val validator.ValidatorKeeper, param param.ParamKeeper) WeightedMedianPriceManager {
	return WeightedMedianPriceManager{
		store: model.NewPriceStorage(key),
		param: param,
		val:   val,
	}
}

type weightedValidator struct {
	validator linotypes.AccountKey
	weight    linotypes.Coin
	price     linotypes.MiniDollar
}

// set current price.
func (wm WeightedMedianPriceManager) InitGenesis(ctx sdk.Context, initPrice linotypes.MiniDollar) sdk.Error {
	if !initPrice.IsPositive() {
		return types.ErrInvalidPriceFeed(initPrice)
	}
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
	defer wm.updateLastValidatorSet(ctx)
	wvals := wm.getWeightedValidators(ctx)
	if len(wvals) == 0 {
		return types.ErrNoValidator()
	}
	blocktime := ctx.BlockTime().Unix()
	wvals, err := wm.filterAndSlash(ctx, wvals)
	if err != nil {
		return err
	}
	var price linotypes.MiniDollar
	if len(wvals) == 0 {
		// no valid price this hour, use the same price from last hour.
		// this is irrelevant to testnet mode, CANNOT use CurrPrice.
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

func (wm WeightedMedianPriceManager) updateLastValidatorSet(ctx sdk.Context) {
	vals := wm.val.GetCommittingValidators(ctx)
	wm.store.SetLastValidators(ctx, vals)
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
	blocktime := ctx.BlockTime().Unix()
	last, err := wm.store.GetFedPrice(ctx, validator)
	feedEverySec := wm.param.GetPriceParam(ctx).FeedEverySec
	// have fed price before(err is nil) and too frequent.
	if err == nil && blocktime-last.UpdateAt < feedEverySec {
		return types.ErrPriceFeedRateLimited()
	}

	wm.store.SetFedPrice(ctx, &model.FedPrice{
		Validator: validator,
		Price:     price,
		UpdateAt:  blocktime,
	})
	return nil
}

func (wm WeightedMedianPriceManager) CoinToMiniDollar(ctx sdk.Context, coin linotypes.Coin) (linotypes.MiniDollar, sdk.Error) {
	price, err := wm.CurrPrice(ctx)
	if err != nil {
		return linotypes.NewMiniDollar(0), err
	}
	return coinToMiniDollar(coin, price), nil
}

func (wm WeightedMedianPriceManager) MiniDollarToCoin(ctx sdk.Context, dollar linotypes.MiniDollar) (linotypes.Coin, linotypes.MiniDollar, sdk.Error) {
	price, err := wm.CurrPrice(ctx)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), linotypes.NewMiniDollar(0), err
	}
	bought, used := miniDollarToCoin(dollar, price)
	return bought, used, nil
}

func (wm WeightedMedianPriceManager) CurrPrice(ctx sdk.Context) (linotypes.MiniDollar, sdk.Error) {
	if wm.param.GetPriceParam(ctx).TestnetMode {
		return linotypes.TestnetPrice, nil
	}
	curr, err := wm.store.GetCurrentPrice(ctx)
	if err != nil {
		return linotypes.NewMiniDollar(0), err
	}
	return curr.Price, nil
}

func (wm WeightedMedianPriceManager) isValidator(ctx sdk.Context, user linotypes.AccountKey) bool {
	vals := wm.val.GetCommittingValidators(ctx)
	return linotypes.FindAccountInList(user, vals) != -1
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
	// XXX(yumin): history MUST BE set before it get sorted.
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
	// when the length is an even number, use higher, e.g. 4 / 2 = 2, which is [0, 1, 2, 3].
	mid := len(history) / 2
	current := history[mid]
	wm.store.SetCurrentPrice(ctx, &current)
}

// getWeightedValidators return weighted validators, sorted by (weight, namestr), increasingly.
// price fields are empty value.
func (wm WeightedMedianPriceManager) getWeightedValidators(ctx sdk.Context) []weightedValidator {
	wvals := make([]weightedValidator, 0)
	vals := wm.val.GetCommittingValidatorVoteStatus(ctx)
	for _, val := range vals {
		wvals = append(wvals, weightedValidator{
			validator: val.ValidatorName,
			weight:    val.ReceivedVotes,
			price:     linotypes.NewMiniDollar(0),
		})
	}
	return wvals
}

// filterAndSlash slash validators that missed price feeding.
// premise: fedPrice needs to be validated upon validators send update message.
func (wm WeightedMedianPriceManager) filterAndSlash(ctx sdk.Context, wvals []weightedValidator) (rst []weightedValidator, err sdk.Error) {
	lastValidatorSet := wm.lastRoundValidatorSet(ctx)
	blocktime := ctx.BlockTime().Unix()
	for i := range wvals {
		valname := wvals[i].validator
		fedPrice, err := wm.store.GetFedPrice(ctx, valname)
		updateEverySec := wm.param.GetPriceParam(ctx).UpdateEverySec
		if err != nil || blocktime-fedPrice.UpdateAt > updateEverySec {
			// unless the validator is not in the last set, slash.
			if lastValidatorSet[valname] {
				if !wm.param.GetPriceParam(ctx).TestnetMode {
					err := wm.val.PunishCommittingValidator(
						ctx, valname,
						wm.param.GetPriceParam(ctx).PenaltyMissFeed,
						linotypes.PunishNoPriceFed)
					if err != nil {
						return nil, err
					}
				}
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
