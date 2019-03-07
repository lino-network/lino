package reputation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"

	"math/big"

	model "github.com/lino-network/lino/x/reputation/internal"
)

// ReputationManager - adaptor for reputation math model and cosmos application.
type ReputationManager struct {
	storeKey    sdk.StoreKey
	paramHolder param.ParamHolder
}

// NewReputationManager - require holder for BestContentIndexN
func NewReputationManager(key sdk.StoreKey, holder param.ParamHolder) ReputationManager {
	return ReputationManager{
		storeKey:    key,
		paramHolder: holder,
	}
}

// construct a handler.
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
		return ErrAccountNotFound("")
	}
	return nil
}

func (rep ReputationManager) checkPost(pid model.Pid) sdk.Error {
	if len(pid) == 0 {
		return ErrPostNotFound("")
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

// DonateAt - It's caller's responsibility that parameters are all correct,
// although we do have some checks.
func (rep ReputationManager) DonateAt(ctx sdk.Context,
	username types.AccountKey, post types.Permlink, coinDay types.Coin) (types.Coin, sdk.Error) {
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

	dp := handler.DonateAt(uid, pid, coinDay.Amount.BigInt())
	return types.NewCoinFromBigInt(dp), nil
}

// ReportAt - @p username report @p post.
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

func (rep ReputationManager) calcFreeScore(amount types.Coin) *big.Int {
	score := amount.Amount.BigInt()
	score.Mul(score, big.NewInt(15))
	score.Div(score, big.NewInt(10000)) // 0.15% freescore if you lock down.
	return score
}

// OnStakeIn - on @p username stakein @p amount.
func (rep ReputationManager) OnStakeIn(ctx sdk.Context,
	username types.AccountKey, amount types.Coin) {
	incAmount := rep.calcFreeScore(amount)
	rep.incFreeScore(ctx, username, incAmount)
}

// OnStakeOut - on @p username stakeout @p amount
func (rep ReputationManager) OnStakeOut(ctx sdk.Context,
	username types.AccountKey, amount types.Coin) {
	incAmount := rep.calcFreeScore(amount)
	incAmount.Neg(incAmount)
	rep.incFreeScore(ctx, username, incAmount)
}

func (rep ReputationManager) incFreeScore(ctx sdk.Context,
	username types.AccountKey, score *big.Int) sdk.Error {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return err
	}
	uid := string(username)
	err = rep.checkUsername(uid)
	if err != nil {
		return err
	}

	handler.IncFreeScore(uid, score)
	return nil
}

// Update - on blocker end, update reputation time related information.
func (rep ReputationManager) Update(ctx sdk.Context) sdk.Error {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return err
	}

	handler.Update(ctx.BlockHeader().Time.Unix())
	return nil
}

// GetRepution - return reputation of @p username, costomnerScore + freeScore.
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

// GetSumRep of @p post
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

// GetCurrentRound of now
func (rep ReputationManager) GetCurrentRound(ctx sdk.Context) (int64, sdk.Error) {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return 0, err
	}

	_, ts := handler.GetCurrentRound()
	return ts, nil
}

// ExportToFile state of reputation system.
func (rep ReputationManager) ExportToFile(ctx sdk.Context, file string) error {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return err
	}
	handler.ExportToFile(file)
	return nil
}

// ImportFromFile state of reputation system.
func (rep ReputationManager) ImportFromFile(ctx sdk.Context, file string) error {
	handler, err := rep.getHandler(ctx)
	if err != nil {
		return err
	}
	handler.ImportFromFile(file)
	return nil
}
