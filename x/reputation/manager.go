package reputation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	model "github.com/lino-network/lino/x/reputation/internal"
	repv2 "github.com/lino-network/lino/x/reputation/repv2"
)

// ReputationManager - adaptor for reputation math model and cosmos application.
type ReputationManager struct {
	storeKey    sdk.StoreKey
	paramHolder param.ParamHolder
}

// NewReputationManager - require holder for BestContentIndexN
func NewReputationManager(storeKey sdk.StoreKey, holder param.ParamHolder) ReputationKeeper {
	return ReputationManager{
		storeKey:    storeKey,
		paramHolder: holder,
	}
}

// construct a handler.
func (rep ReputationManager) getHandlerV2(ctx sdk.Context) repv2.Reputation {
	store := ctx.KVStore(rep.storeKey)
	repStore := repv2.NewReputationStore(store, repv2.DefaultInitialReputation)
	param := rep.paramHolder.GetReputationParam(ctx)
	handler := repv2.NewReputation(
		repStore, param.BestContentIndexN, param.UserMaxN,
		repv2.DefaultRoundDurationSeconds,
		repv2.DefaultSampleWindowSize,
		repv2.DefaultDecayFactor)
	return handler
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
	username types.AccountKey, post types.Permlink, amount types.MiniDollar) (types.MiniDollar, sdk.Error) {
	uid := string(username)
	pid := string(post)
	err := rep.basicCheck(uid, pid)
	if err != nil {
		return types.NewMiniDollar(0), err
	}

	// Update6, start to use new reputation algorithm.
	handler := rep.getHandlerV2(ctx)
	dp := handler.DonateAt(repv2.Uid(uid), repv2.Pid(pid), repv2.NewIntFromBig(amount.Int.BigInt()))
	return types.NewMiniDollarFromBig(dp.Int), nil
}

// Update - on blocker end, update reputation time related information.
func (rep ReputationManager) Update(ctx sdk.Context) sdk.Error {
	handler := rep.getHandlerV2(ctx)
	handler.Update(repv2.Time(ctx.BlockHeader().Time.Unix()))
	return nil
}

// GetRepution - return reputation of @p username, costomnerScore + freeScore.
func (rep ReputationManager) GetReputation(ctx sdk.Context, username types.AccountKey) (types.MiniDollar, sdk.Error) {
	uid := string(username)
	err := rep.checkUsername(uid)
	if err != nil {
		return types.NewMiniDollar(0), err
	}

	handler := rep.getHandlerV2(ctx)
	return types.NewMiniDollarFromBig(handler.GetReputation(repv2.Uid(uid)).Int), nil
}

// GetCurrentRound of now
func (rep ReputationManager) GetCurrentRound(ctx sdk.Context) (int64, sdk.Error) {
	repv2 := rep.getHandlerV2(ctx)
	_, ts := repv2.GetCurrentRound()
	return int64(ts), nil
}

// ExportToFile state of reputation system.
func (rep ReputationManager) ExportToFile(ctx sdk.Context, file string) error {
	repv2 := rep.getHandlerV2(ctx)
	return repv2.ExportToFile(file)
}

// ImportFromFile state of reputation system.
// after update6's code is merged, V2 is the only version that will exist.
func (rep ReputationManager) ImportFromFile(ctx sdk.Context, file string) error {
	handler := rep.getHandlerV2(ctx)
	return handler.ImportFromFile(file)
}
