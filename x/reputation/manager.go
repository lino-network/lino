package reputation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"

	acc "github.com/lino-network/lino/x/account"
	post "github.com/lino-network/lino/x/post"
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

func (rep ReputationManager) getHandler(ctx sdk.Context) (model.Reputation, sdk.Error) {
	store := ctx.KVStore(rep.storeKey)
	param, err := rep.paramHolder.GetReputationParam(ctx)
	if err != nil {
		return nil, err
	}
	repStore := model.NewReputationStore(store, param.BestContentIndexN)
	handler := model.NewReputation(repStore)
	return handler, nil
}

func (rep ReputationManager) checkUsername(uid model.Uid) sdk.Error {
	if len(uid) == 0 {
		return acc.ErrAccountNotFound("")
	}
	return nil
}

func (rep ReputationManager) checkPost(pid model.Pid) sdk.Error {
	if len(pid) == 0 {
		return post.ErrPostNotFound("")
	}
	return nil
}

func (rep ReputationManager) basicCheck(uid model.Uid, pid model.Pid) sdk.Error {
	err := rep.checkUsername(uid)
	if err != nil {
		return err
	}
	err = rep.checkPost(pid)
	return err
}

// It's caller's responsibility that parameters are all correct, although we do have some checks.
func (rep ReputationManager) DonateAt(ctx sdk.Context,
	username types.AccountKey, post types.Permlink, stake types.Coin) (types.Coin, sdk.Error) {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	uid := string(username)
	pid := string(post)
	err = rep.basicCheck(uid, pid)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	dp := handler.DonateAt(uid, pid, stake.Amount.BigInt())
	return types.NewCoinFromBigInt(dp), nil
}

func (rep ReputationManager) ReportAt(ctx sdk.Context,
	username types.AccountKey, post types.Permlink) (types.Coin, sdk.Error) {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	uid := string(username)
	pid := string(post)
	err = rep.basicCheck(uid, pid)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	sumRep := handler.ReportAt(uid, pid)
	return types.NewCoinFromBigInt(sumRep), nil
}

func (rep ReputationManager) IncFreeScore(ctx sdk.Context,
	username types.AccountKey, score types.Coin) sdk.Error {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return err
	}

	uid := string(username)
	err = rep.checkUsername(uid)
	if err != nil {
		return err
	}

	handler.IncFreeScore(uid, score.Amount.BigInt())
	return nil
}

func (rep ReputationManager) Update(ctx sdk.Context) sdk.Error {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return err
	}

	handler.Update(ctx.BlockHeader().Time.Unix())
	return nil
}

func (rep ReputationManager) GetReputation(ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	uid := string(username)
	err = rep.checkUsername(uid)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	return types.NewCoinFromBigInt(handler.GetReputation(uid)), nil
}

func (rep ReputationManager) GetSumRep(ctx sdk.Context, post types.Permlink) (types.Coin, sdk.Error) {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	pid := string(post)
	err = rep.checkPost(pid)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	return types.NewCoinFromBigInt(handler.GetSumRep(pid)), nil
}

func (rep ReputationManager) GetCurrentRound(ctx sdk.Context) (int64, sdk.Error) {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return 0, err
	}

	_, ts := handler.GetCurrentRound()
	return ts, nil
}
