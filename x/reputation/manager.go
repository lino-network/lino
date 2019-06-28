package reputation

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	model "github.com/lino-network/lino/x/reputation/internal"
	repv2 "github.com/lino-network/lino/x/reputation/repv2"
)

// ReputationManager - adaptor for reputation math model and cosmos application.
type ReputationManager struct {
	v1key       sdk.StoreKey
	v2key       sdk.StoreKey
	paramHolder param.ParamHolder
}

// NewReputationManager - require holder for BestContentIndexN
func NewReputationManager(v1key sdk.StoreKey, v2key sdk.StoreKey, holder param.ParamHolder) ReputationManager {
	return ReputationManager{
		v1key:       v1key,
		v2key:       v2key,
		paramHolder: holder,
	}
}

// construct a handler.
func (rep ReputationManager) getHandlerV2(ctx sdk.Context) repv2.Reputation {
	store := ctx.KVStore(rep.v2key)
	repStore := repv2.NewReputationStore(store, repv2.DefaultInitialReputation)
	handler := repv2.NewReputation(repStore, 200, 50, 25*3600, 10, 10)
	return handler
}

// construct a handler.
func (rep ReputationManager) getHandler(ctx sdk.Context) (model.Reputation, sdk.Error) {
	store := ctx.KVStore(rep.v1key)
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

func (rep ReputationManager) migrate(handler model.Reputation, repv2 repv2.Reputation, uid model.Uid) {
	if repv2.RequireMigrate(uid) {
		prev := handler.GetReputation(uid)
		// only when user's reputation is larger than initial, migrate it.
		if prev.Cmp(big.NewInt(model.InitialCustomerScore)) > 0 {
			repv2.MigrateFromV1(uid, handler.GetReputation(uid))
		}
	}
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

	// Update6, start to use new reputation algorithm.
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		repv2 := rep.getHandlerV2(ctx)
		rep.migrate(handler, repv2, uid)
		dp := repv2.DonateAt(uid, pid, coinDay.Amount.BigInt())
		return types.NewCoinFromBigInt(dp), nil
	}

	dp := handler.DonateAt(uid, pid, coinDay.Amount.BigInt())
	return types.NewCoinFromBigInt(dp), nil
}

// ReportAt - @p username report @p post.
func (rep ReputationManager) ReportAt(ctx sdk.Context,
	username types.AccountKey, post types.Permlink) (types.Coin, sdk.Error) {
	// Update6, report is deprecated.
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		return types.NewCoinFromInt64(0), nil
	}
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
	_ = rep.incFreeScore(ctx, username, incAmount)
}

// OnStakeOut - on @p username stakeout @p amount
func (rep ReputationManager) OnStakeOut(ctx sdk.Context,
	username types.AccountKey, amount types.Coin) {
	incAmount := rep.calcFreeScore(amount)
	incAmount.Neg(incAmount)
	_ = rep.incFreeScore(ctx, username, incAmount)
}

func (rep ReputationManager) incFreeScore(ctx sdk.Context,
	username types.AccountKey, score *big.Int) sdk.Error {
	// After Update6, no-op on StakeIn/StakeOut.
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		return nil
	}
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
	// Update6
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		repv2 := rep.getHandlerV2(ctx)
		repv2.Update(ctx.BlockHeader().Time.Unix())
		return nil
	}

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

	// Update6
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		repv2 := rep.getHandlerV2(ctx)
		if err != nil {
			return types.NewCoinFromInt64(0), err
		}
		rep.migrate(handler, repv2, uid)
		return types.NewCoinFromBigInt(repv2.GetReputation(uid)), nil
	}
	return types.NewCoinFromBigInt(handler.GetReputation(uid)), nil
}

// GetSumRep of @p post
func (rep ReputationManager) GetSumRep(ctx sdk.Context, post types.Permlink) (types.Coin, sdk.Error) {
	// Update6, sumrep deprecated.
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		return types.NewCoinFromInt64(0), nil
	}
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

	// Update6
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		repv2 := rep.getHandlerV2(ctx)
		if err != nil {
			return 0, err
		}
		_, ts := repv2.GetCurrentRound()
		return ts, nil
	}

	_, ts := handler.GetCurrentRound()
	return ts, nil
}

// ExportToFile state of reputation system.
func (rep ReputationManager) ExportToFile(ctx sdk.Context, file string) error {
	param, err := rep.paramHolder.GetReputationParam(ctx)
	if err != nil {
		return err
	}
	v1RepStore := model.NewReputationStore(ctx.KVStore(rep.v1key), param.BestContentIndexN)
	v1handler := model.NewReputation(v1RepStore)
	// Update6, if a user does not donate after update6, his reputation is reset to 0.
	if ctx.BlockHeight() >= types.BlockchainUpgrade1Update6Height {
		repv2 := rep.getHandlerV2(ctx)
		if err != nil {
			return err
		}
		v1RepStore.IterateUsers(func(uid string) bool {
			rep.migrate(v1handler, repv2, uid)
			return false
		})
		repv2.ExportToFile(file)
		return nil
	}

	v1handler.ExportToFile(file)
	return nil
}

// ImportFromFile state of reputation system.
// after update6's code is merged, V2 is the only version that will exist.
func (rep ReputationManager) ImportFromFile(ctx sdk.Context, file string) error {
	handler := rep.getHandlerV2(ctx)
	handler.ImportFromFile(file)
	return nil
}
