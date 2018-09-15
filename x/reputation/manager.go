package reputation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"

	model "github.com/lino-network/lino/x/reputation/internal"
)

type ReputationManager struct {
	storeKey    sdk.StoreKey
	paramHolder param.ParamHolder
}

func NewReputationManager(key sdk.StoreKey, holder param.ParamHolder) ReputationManager {
	return ReputationManager{
		storeKey:    key,
		paramHolder: holder,
	}
}

func (rep ReputationManager) getHandler(ctx sdk.Context) model.Reputation {
	store := ctx.KVStore(rep.storeKey)
	// TODO(yumin): use parameter.
	repStore := model.NewReputationStoreDefaultN(store)
	handler := model.NewReputation(repStore)
	return handler
}

// It's caller's responsibility that parameters are all correct, although we do have some checks.
func (rep ReputationManager) DonateAt(ctx sdk.Context,
	username types.AccountKey, post types.Permlink, stake types.Coin) (types.Coin, sdk.Error) {
	handler := rep.getHandler(ctx)
	uid := string(username)
	pid := string(post)
	dp := handler.DonateAt(uid, pid, stake.Amount.BigInt())
	return types.NewCoinFromBigInt(dp), nil
}

func (rep ReputationManager) ReportAt(ctx sdk.Context,
	username types.AccountKey, post types.Permlink) types.Coin {
	handler := rep.getHandler(ctx)
	uid := string(username)
	pid := string(post)
	sumRep := handler.ReportAt(uid, pid)
	return types.NewCoinFromBigInt(sumRep)
}

func (rep ReputationManager) IncFreeScore(ctx sdk.Context,
	username types.AccountKey, score types.Coin) {
	handler := rep.getHandler(ctx)
	uid := string(username)
	handler.IncFreeScore(uid, score.Amount.BigInt())
}

func (rep ReputationManager) Update(ctx sdk.Context) {
	handler := rep.getHandler(ctx)
	handler.Update(ctx.BlockHeader().Time.Unix())
}

func (rep ReputationManager) GetReputation(ctx sdk.Context, username types.AccountKey) types.Coin {
	handler := rep.getHandler(ctx)
	return types.NewCoinFromBigInt(handler.GetReputation(string(username)))
}

func (rep ReputationManager) GetSumRep(ctx sdk.Context, post types.Permlink) types.Coin {
	handler := rep.getHandler(ctx)
	return types.NewCoinFromBigInt(handler.GetSumRep(string(post)))
}

func (rep ReputationManager) GetCurrentRound(ctx sdk.Context) int64 {
	handler := rep.getHandler(ctx)
	_, ts := handler.GetCurrentRound()
	return ts
}
