package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/bandwidth/model"
	"github.com/lino-network/lino/x/bandwidth/types"
	global "github.com/lino-network/lino/x/global"
)

// BandwidthManager - bandwidth manager
type BandwidthManager struct {
	storage     model.BandwidthStorage
	paramHolder param.ParamKeeper
	// deps
	gm global.GlobalKeeper
}

func NewBandwidthManager(key sdk.StoreKey, holder param.ParamKeeper, gm global.GlobalKeeper) *BandwidthManager {
	return &BandwidthManager{
		storage:     model.NewBandwidthStorage(key),
		paramHolder: holder,
		gm:          gm,
	}
}

// InitGenesis - initialize KV Store
func (bm BandwidthManager) InitGenesis(ctx sdk.Context) error {
	bandwidthInfo := &model.BandwidthInfo{
		GeneralMsgEMA: sdk.NewDec(0),
		AppMsgEMA:     sdk.NewDec(0),
		MaxMPS:        sdk.NewDec(0),
	}

	if err := bm.storage.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}

	blockInfo := &model.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            linotypes.NewCoinFromInt64(int64(0)),
	}

	if err := bm.storage.SetBlockInfo(ctx, blockInfo); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) IsUserMsgFeeEnough(ctx sdk.Context, fee auth.StdFee) bool {
	if !fee.Amount.IsValid() {
		return false
	}
	providedFee := linotypes.NewCoinFromInt64(fee.Amount.AmountOf(linotypes.LinoCoinDenom).Int64())
	info, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return false
	}
	return providedFee.IsGTE(info.CurMsgFee)
}

func (bm BandwidthManager) AddMsgSignedByApp(ctx sdk.Context, num int64) sdk.Error {
	info, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}
	info.TotalMsgSignedByApp += num

	if err := bm.storage.SetBlockInfo(ctx, info); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) AddMsgSignedByUser(ctx sdk.Context, num int64) sdk.Error {
	info, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}
	info.TotalMsgSignedByUser += num

	if err := bm.storage.SetBlockInfo(ctx, info); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) ClearBlockInfo(ctx sdk.Context) sdk.Error {
	info, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}
	info.TotalMsgSignedByUser = 0
	info.TotalMsgSignedByApp = 0

	if err := bm.storage.SetBlockInfo(ctx, info); err != nil {
		return err
	}
	return nil
}

// calcuate the new EMA at the end of each block
func (bm BandwidthManager) UpdateMaxMPSAndEMA(ctx sdk.Context) sdk.Error {
	lastBlockTime, err := bm.gm.GetLastBlockTime(ctx)
	if err != nil {
		return err
	}

	if lastBlockTime >= ctx.BlockHeader().Time.Unix() {
		return nil
	}

	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return err
	}

	blockInfo, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	pastTime := ctx.BlockHeader().Time.Unix() - lastBlockTime

	// EMA_general = EMA_general_prev * (1 - k_general) + generalMPS * k_general
	generalMPS := linotypes.NewDecFromRat(int64(blockInfo.TotalMsgSignedByUser), pastTime)
	bandwidthInfo.GeneralMsgEMA = bm.calculateEMA(bandwidthInfo.GeneralMsgEMA, params.GeneralMsgEMAFactor, generalMPS)

	// EMA_app = EMA_app_prev * (1 - k_app) + appMPS * k_app
	appMPS := linotypes.NewDecFromRat(int64(blockInfo.TotalMsgSignedByApp), pastTime)
	bandwidthInfo.AppMsgEMA = bm.calculateEMA(bandwidthInfo.AppMsgEMA, params.AppMsgEMAFactor, appMPS)

	// MaxMPS = max( (totalMsgSignedByUser + totalMsgSignedByApp)/(curBlockTime - lastBlockTime), MaxMPS)
	totalMPS := linotypes.NewDecFromRat(int64(blockInfo.TotalMsgSignedByUser)+int64(blockInfo.TotalMsgSignedByApp), pastTime)
	if totalMPS.GT(bandwidthInfo.MaxMPS) {
		bandwidthInfo.MaxMPS = totalMPS
	}

	if err := bm.storage.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}

	return nil
}

// calcuate the current msg fee based on last block info at the begining of each block
func (bm BandwidthManager) CalculateCurMsgFee(ctx sdk.Context) sdk.Error {
	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return err
	}

	blockInfo, err := bm.storage.GetBlockInfo(ctx)
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
	if !generalMsgQuota.IsPositive() {
		return types.ErrInvalidMsgQuota()
	}

	expResult := bm.approximateExp(bandwidthInfo.GeneralMsgEMA.Sub(generalMsgQuota).Quo(generalMsgQuota).Mul(params.MsgFeeFactorA))
	msgFeeLino := expResult.Mul(params.MsgFeeFactorB)
	blockInfo.CurMsgFee = linotypes.NewCoinFromInt64(msgFeeLino.Mul(sdk.NewDec(linotypes.Decimals)).RoundInt64())

	if err := bm.storage.SetBlockInfo(ctx, blockInfo); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) DecayMaxMPS(ctx sdk.Context) sdk.Error {
	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	bandwidthInfo.MaxMPS = bandwidthInfo.MaxMPS.Mul(params.MaxMPSDecayRate)
	if err := bm.storage.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) calculateEMA(prevEMA sdk.Dec, k sdk.Dec, curMPS sdk.Dec) sdk.Dec {
	pre := prevEMA.Mul(sdk.NewDec(1).Sub(k))
	cur := curMPS.Mul(k)
	return pre.Add(cur)
}

func (bm BandwidthManager) approximateExp(x sdk.Dec) sdk.Dec {
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
