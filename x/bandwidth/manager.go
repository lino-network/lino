package bandwidth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/bandwidth/model"
)

// BandwidthManager - bandwidth manager
type BandwidthManager struct {
	storage     model.BandwidthStorage
	paramHolder param.ParamHolder
}

func NewBandwidthManager(key sdk.StoreKey, holder param.ParamHolder) BandwidthManager {
	return BandwidthManager{
		storage:     model.NewBandwidthStorage(key),
		paramHolder: holder,
	}
}

// InitGenesis - initialize KV Store
func (bm BandwidthManager) InitGenesis(ctx sdk.Context) error {
	if err := bm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) IsUserMsgFeeEnough(ctx sdk.Context, fee auth.StdFee) bool {
	blockInfo, err := bm.storage.GetCurBlockInfo(ctx)
	if err != nil {
		return false
	}

	curFeeCoin := types.DecToCoin(blockInfo.CurMsgFee.Mul(sdk.NewDec(100000)))
	providedFee := types.NewCoinFromInt64(fee.Amount.AmountOf("LNO").Int64())
	return providedFee.IsGT(curFeeCoin)
}

func (bm BandwidthManager) AddMsgSignedByApp(ctx sdk.Context, num uint32) sdk.Error {
	blockInfo, err := bm.storage.GetCurBlockInfo(ctx)
	if err != nil {
		return err
	}

	blockInfo.TotalMsgSignedByApp += num
	if err := bm.storage.SetCurBlockInfo(ctx, blockInfo); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) AddMsgSignedByUser(ctx sdk.Context, num uint32) sdk.Error {
	blockInfo, err := bm.storage.GetCurBlockInfo(ctx)
	if err != nil {
		return err
	}

	blockInfo.TotalMsgSignedByUser += num
	if err := bm.storage.SetCurBlockInfo(ctx, blockInfo); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) ClearCurBlockInfo(ctx sdk.Context) sdk.Error {
	blockInfo, err := bm.storage.GetCurBlockInfo(ctx)
	if err != nil {
		return err
	}

	blockInfo.TotalMsgSignedByUser = 0
	blockInfo.TotalMsgSignedByApp = 0

	if err := bm.storage.SetCurBlockInfo(ctx, blockInfo); err != nil {
		return err
	}
	return nil
}

// calcuate the new EMA at the end of each block
func (bm BandwidthManager) UpdateMaxMPSAndEMA(ctx sdk.Context, lastBlockTime int64) sdk.Error {
	if lastBlockTime == ctx.BlockHeader().Time.Unix() {
		return nil
	}

	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return err
	}

	blockInfo, err := bm.storage.GetCurBlockInfo(ctx)
	if err != nil {
		return err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	// EMA_general = EMA_general_prev * (1 - k_general) + generalMPS * k_general
	generalMPS := types.NewDecFromRat(int64(blockInfo.TotalMsgSignedByUser), ctx.BlockHeader().Time.Unix()-lastBlockTime)
	bandwidthInfo.GeneralMsgEMA = bandwidthInfo.GeneralMsgEMA.Mul(sdk.NewDec(1).Sub(params.GeneralMsgEMAFactor)).Add(generalMPS.Mul(params.GeneralMsgEMAFactor))

	// EMA_app = EMA_app_prev * (1 - k_app) + appMPS * k_app
	appMPS := types.NewDecFromRat(int64(blockInfo.TotalMsgSignedByApp), ctx.BlockHeader().Time.Unix()-lastBlockTime)
	bandwidthInfo.AppMsgEMA = bandwidthInfo.AppMsgEMA.Mul(sdk.NewDec(1).Sub(params.AppMsgEMAFactor)).Add(appMPS.Mul(params.AppMsgEMAFactor))

	// MaxMPS = max( (totalMsgSignedByUser + totalMsgSignedByApp)/(curBlockTime - lastBlockTime), MaxMPS)
	totalMPS := types.NewDecFromRat(int64(blockInfo.TotalMsgSignedByUser)+int64(blockInfo.TotalMsgSignedByApp), ctx.BlockHeader().Time.Unix()-lastBlockTime)
	if totalMPS.GT(bandwidthInfo.MaxMPS) {
		bandwidthInfo.MaxMPS = totalMPS
	}

	if err := bm.storage.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}
	return nil
}

// calcuate the current msg fee at the begining of each block
func (bm BandwidthManager) CalculateCurMsgFee(ctx sdk.Context) sdk.Error {
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
	blockInfo, err := bm.storage.GetCurBlockInfo(ctx)
	if err != nil {
		return err
	}

	expResult := bm.approximateExp(bandwidthInfo.GeneralMsgEMA.Sub(generalMsgQuota).Quo(generalMsgQuota).Mul(params.MsgFeeFactorA))
	blockInfo.CurMsgFee = expResult.Mul(params.MsgFeeFactorB)

	if err := bm.storage.SetCurBlockInfo(ctx, blockInfo); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) approximateExp(x sdk.Dec) sdk.Dec {
	fmt.Println(x)
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
