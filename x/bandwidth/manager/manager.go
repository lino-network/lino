package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/bandwidth/model"
	global "github.com/lino-network/lino/x/global"
)

// BandwidthManager - bandwidth manager
type BandwidthManager struct {
	storage     model.BandwidthStorage
	paramHolder param.ParamKeeper
	// in-memory storage
	blockStatsCache model.BlockStatsCache
	// deps
	gm global.GlobalKeeper
}

func NewBandwidthManager(key sdk.StoreKey, holder param.ParamKeeper, gm global.GlobalKeeper) *BandwidthManager {
	return &BandwidthManager{
		storage:         model.NewBandwidthStorage(key),
		paramHolder:     holder,
		gm:              gm,
		blockStatsCache: model.BlockStatsCache{},
	}
}

// InitGenesis - initialize KV Store
func (bm *BandwidthManager) InitGenesis(ctx sdk.Context) error {
	if err := bm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (bm *BandwidthManager) IsUserMsgFeeEnough(ctx sdk.Context, fee auth.StdFee) bool {
	curFeeCoin := types.DecToCoin(bm.blockStatsCache.CurMsgFee.Mul(sdk.NewDec(100000)))
	providedFee := types.NewCoinFromInt64(fee.Amount.AmountOf("lino").Int64())
	return providedFee.IsGT(curFeeCoin)
}

func (bm *BandwidthManager) AddMsgSignedByApp(ctx sdk.Context, num uint32) sdk.Error {
	bm.blockStatsCache.TotalMsgSignedByApp += num
	return nil
}

func (bm *BandwidthManager) AddMsgSignedByUser(ctx sdk.Context, num uint32) sdk.Error {
	bm.blockStatsCache.TotalMsgSignedByUser += num
	return nil
}

func (bm *BandwidthManager) ClearBlockStatsCache(ctx sdk.Context) sdk.Error {
	bm.blockStatsCache.TotalMsgSignedByApp = 0
	bm.blockStatsCache.TotalMsgSignedByUser = 0
	bm.blockStatsCache.CurMsgFee = sdk.NewDec(0)
	return nil
}

// calcuate the new EMA at the end of each block
func (bm *BandwidthManager) UpdateMaxMPSAndEMA(ctx sdk.Context) sdk.Error {
	lastBlockTime, err := bm.gm.GetLastBlockTime(ctx)
	if err != nil {
		return err
	}

	if lastBlockTime == ctx.BlockHeader().Time.Unix() {
		return nil
	}

	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	// EMA_general = EMA_general_prev * (1 - k_general) + generalMPS * k_general
	generalMPS := types.NewDecFromRat(int64(bm.blockStatsCache.TotalMsgSignedByUser), ctx.BlockHeader().Time.Unix()-lastBlockTime)
	bandwidthInfo.GeneralMsgEMA = bandwidthInfo.GeneralMsgEMA.Mul(sdk.NewDec(1).Sub(params.GeneralMsgEMAFactor)).Add(generalMPS.Mul(params.GeneralMsgEMAFactor))

	// EMA_app = EMA_app_prev * (1 - k_app) + appMPS * k_app
	appMPS := types.NewDecFromRat(int64(bm.blockStatsCache.TotalMsgSignedByApp), ctx.BlockHeader().Time.Unix()-lastBlockTime)
	bandwidthInfo.AppMsgEMA = bandwidthInfo.AppMsgEMA.Mul(sdk.NewDec(1).Sub(params.AppMsgEMAFactor)).Add(appMPS.Mul(params.AppMsgEMAFactor))

	// MaxMPS = max( (totalMsgSignedByUser + totalMsgSignedByApp)/(curBlockTime - lastBlockTime), MaxMPS)
	totalMPS := types.NewDecFromRat(int64(bm.blockStatsCache.TotalMsgSignedByUser)+int64(bm.blockStatsCache.TotalMsgSignedByApp), ctx.BlockHeader().Time.Unix()-lastBlockTime)
	if totalMPS.GT(bandwidthInfo.MaxMPS) {
		bandwidthInfo.MaxMPS = totalMPS
	}

	if err := bm.storage.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}

	// store the current block stats into lastBlockInfo
	info, err := bm.storage.GetLastBlockInfo(ctx)
	if err != nil {
		return err
	}

	info.TotalMsgSignedByApp = bm.blockStatsCache.TotalMsgSignedByApp
	info.TotalMsgSignedByUser = bm.blockStatsCache.TotalMsgSignedByUser

	if err := bm.storage.SetLastBlockInfo(ctx, info); err != nil {
		return err
	}

	return nil
}

// calcuate the current msg fee based on last block info at the begining of each block
func (bm *BandwidthManager) CalculateCurMsgFee(ctx sdk.Context) sdk.Error {
	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return err
	}

	curMaxMPS := sdk.NewDec(0)
	if params.ExpectedMaxMPS.GT(bandwidthInfo.MaxMPS) {
		curMaxMPS = params.ExpectedMaxMPS
	} else {
		curMaxMPS = bandwidthInfo.MaxMPS
	}

	generalMsgQuota := params.GeneralMsgQuotaRatio.Mul(curMaxMPS)

	expResult := bm.approximateExp(bandwidthInfo.GeneralMsgEMA.Sub(generalMsgQuota).Quo(generalMsgQuota).Mul(params.MsgFeeFactorA))
	bm.blockStatsCache.CurMsgFee = expResult.Mul(params.MsgFeeFactorB)
	return nil
}

func (bm *BandwidthManager) approximateExp(x sdk.Dec) sdk.Dec {
	prev := x
	x = sdk.NewDec(1).Add(x.Abs().Quo(sdk.NewDec(1024)))
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)
	x = x.Mul(x)

	if prev.LT(sdk.NewDec(0)) {
		return sdk.NewDec(1).Quo(x)
	}
	return x
}
