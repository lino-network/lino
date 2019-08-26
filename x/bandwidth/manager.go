package bandwidth

import (
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

func (bm BandwidthManager) IsAppBandwidthEnough(ctx sdk.Context, username types.AccountKey) bool {
	return false
}
func (bm BandwidthManager) IsUserMsgFeeEnough(ctx sdk.Context, username types.AccountKey, fee auth.StdFee) bool {
	return false
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
func (bm BandwidthManager) UpdateEMA(ctx sdk.Context) sdk.Error {
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

	// EMA_general = EMA_general_prev * (1 - k_general) + numOfGeneralMsgs * k_general
	totalMsgSignedByUser := sdk.NewDec(int64(blockInfo.TotalMsgSignedByUser))
	bandwidthInfo.GeneralMsgEMA = bandwidthInfo.GeneralMsgEMA.Mul(sdk.NewDec(1).Sub(params.GeneralMsgEMAFactor)).Add(totalMsgSignedByUser.Mul(params.GeneralMsgEMAFactor))

	// EMA_app = EMA_app_prev * (1 - k_app) + numOfAppMsgs * k_app
	totalMsgSignedByApp := sdk.NewDec(int64(blockInfo.TotalMsgSignedByApp))
	bandwidthInfo.AppMsgEMA = bandwidthInfo.AppMsgEMA.Mul(sdk.NewDec(1).Sub(params.AppMsgEMAFactor)).Add(totalMsgSignedByApp.Mul(params.AppMsgEMAFactor))

	if err := bm.storage.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}
	return nil
}
