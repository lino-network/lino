package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/bandwidth/model"
	"github.com/lino-network/lino/x/bandwidth/types"
	developer "github.com/lino-network/lino/x/developer"
	global "github.com/lino-network/lino/x/global"
	vote "github.com/lino-network/lino/x/vote"
)

// BandwidthManager - bandwidth manager
type BandwidthManager struct {
	storage     model.BandwidthStorage
	paramHolder param.ParamKeeper
	// deps
	gm global.GlobalKeeper
	vm vote.VoteKeeper
	dm developer.DeveloperKeeper
}

func NewBandwidthManager(key sdk.StoreKey, holder param.ParamKeeper, gm global.GlobalKeeper, vm vote.VoteKeeper, dm developer.DeveloperKeeper) *BandwidthManager {
	return &BandwidthManager{
		storage:     model.NewBandwidthStorage(key),
		paramHolder: holder,
		gm:          gm,
		vm:          vm,
		dm:          dm,
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

	generalMsgQuota, err := bm.getGeneralMsgQuota(ctx)
	if err != nil {
		return err
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

func (bm BandwidthManager) RefillAppBandwidthCredit(ctx sdk.Context, accKey linotypes.AccountKey) sdk.Error {
	info, err := bm.storage.GetAppBandwidthInfo(ctx, accKey)
	if err != nil {
		return err
	}

	curTime := ctx.BlockHeader().Time.Unix()
	if info.LastRefilledAt >= curTime {
		return nil
	}

	if info.CurBandwidthCredit.GTE(info.MaxBandwidthCredit) {
		return nil
	}

	pastSeconds := curTime - info.LastRefilledAt
	// assume refill rate is equal to expectedMPS
	newCredit := info.ExpectedMPS.Mul(sdk.NewDec(pastSeconds)).Add(info.CurBandwidthCredit)
	if newCredit.GTE(info.MaxBandwidthCredit) {
		info.CurBandwidthCredit = info.MaxBandwidthCredit
	} else {
		info.CurBandwidthCredit = newCredit
	}
	info.LastRefilledAt = curTime

	if bm.storage.SetAppBandwidthInfo(ctx, accKey, info); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) GetVacancyCoeff(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}

	appMsgQuota, err := bm.getAppMsgQuota(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}
	delta := bandwidthInfo.AppMsgEMA.Sub(appMsgQuota)
	return bm.approximateExp(delta.Quo(appMsgQuota).Mul(params.AppVacancyFactor)), nil
}

func (bm BandwidthManager) GetPunishmentCoeff(ctx sdk.Context, accKey linotypes.AccountKey) (sdk.Dec, sdk.Error) {
	lastBlockTime, err := bm.gm.GetLastBlockTime(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}
	if lastBlockTime >= ctx.BlockHeader().Time.Unix() {
		return sdk.NewDec(1), nil
	}
	pastTime := ctx.BlockHeader().Time.Unix() - lastBlockTime

	appInfo, err := bm.storage.GetAppBandwidthInfo(ctx, accKey)
	if err != nil {
		return sdk.NewDec(1), err
	}

	if !appInfo.ExpectedMPS.IsPositive() {
		return sdk.NewDec(1), types.ErrInvalidExpectedMPS()
	}

	curMPS := linotypes.NewDecFromRat(appInfo.MessagesInCurBlock, pastTime)
	delta := curMPS.Sub(appInfo.ExpectedMPS)
	if delta.IsNegative() {
		delta = sdk.NewDec(0)
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}
	return bm.approximateExp(delta.Quo(appInfo.ExpectedMPS).Mul(params.AppPunishmentFactor)), nil
}

func (bm BandwidthManager) ConsumeBandwidthCredit(ctx sdk.Context, costPerMsg sdk.Dec, accKey linotypes.AccountKey) sdk.Error {
	info, err := bm.storage.GetAppBandwidthInfo(ctx, accKey)
	if err != nil {
		return err
	}

	// currently allow the credit be negative
	info.CurBandwidthCredit = info.CurBandwidthCredit.Sub(sdk.NewDec(info.MessagesInCurBlock).Mul(costPerMsg))
	info.MessagesInCurBlock = 0
	if err := bm.storage.SetAppBandwidthInfo(ctx, accKey, info); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) ReCalculateAppBandwidthInfo(ctx sdk.Context) sdk.Error {
	totalAppStakeCoin := linotypes.NewCoinFromInt64(0)

	// calculate all app total stake
	for _, app := range bm.dm.GetLiveDevelopers(ctx) {
		appStakeCoin, err := bm.vm.GetLinoStake(ctx, app.Username)
		if err != nil {
			return err
		}
		totalAppStakeCoin = totalAppStakeCoin.Plus(appStakeCoin)
	}

	if !totalAppStakeCoin.IsPositive() {
		return nil
	}

	// calculate all app MPS quota
	appMsgQuota, err := bm.getAppMsgQuota(ctx)
	if err != nil {
		return err
	}

	// calculate each app's share and expected MPS
	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	for _, app := range bm.dm.GetLiveDevelopers(ctx) {
		appStakeCoin, err := bm.vm.GetLinoStake(ctx, app.Username)
		if err != nil {
			return err
		}
		appStakePct := appStakeCoin.ToDec().Quo(totalAppStakeCoin.ToDec())
		expectedMPS := appMsgQuota.Mul(appStakePct)
		maxBandwidthCredit := expectedMPS.Mul(params.AppBandwidthPoolSize)

		// first time app, refill it's current pool to the full
		if !bm.storage.DoesAppBandwidthInfoExist(ctx, app.Username) {
			newAppInfo := model.AppBandwidthInfo{
				Username:           app.Username,
				MaxBandwidthCredit: maxBandwidthCredit,
				CurBandwidthCredit: maxBandwidthCredit,
				MessagesInCurBlock: 0,
				ExpectedMPS:        expectedMPS,
				LastRefilledAt:     ctx.BlockHeader().Time.Unix(),
			}

			if err := bm.storage.SetAppBandwidthInfo(ctx, app.Username, &newAppInfo); err != nil {
				return err
			}

		} else {
			// only update it's max and expected MPS, refill will happen when the app actually send messages.
			info, err := bm.storage.GetAppBandwidthInfo(ctx, app.Username)
			if err != nil {
				return err
			}
			info.ExpectedMPS = expectedMPS
			info.MaxBandwidthCredit = maxBandwidthCredit

			if err := bm.storage.SetAppBandwidthInfo(ctx, app.Username, info); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bm BandwidthManager) GetBandwidthCostPerMsg(ctx sdk.Context, u sdk.Dec, p sdk.Dec) sdk.Dec {
	return u.Mul(p)
}

func (bm BandwidthManager) GetAllAppInfo(ctx sdk.Context) ([]*model.AppBandwidthInfo, sdk.Error) {
	return bm.storage.GetAllAppBandwidthInfo(ctx)
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

func (bm BandwidthManager) getAppMsgQuota(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}

	curMaxMPS := sdk.NewDec(1)
	if params.ExpectedMaxMPS.GT(bandwidthInfo.MaxMPS) {
		curMaxMPS = params.ExpectedMaxMPS
	} else {
		curMaxMPS = bandwidthInfo.MaxMPS
	}

	appMsgQuota := params.AppMsgQuotaRatio.Mul(curMaxMPS)
	if !appMsgQuota.IsPositive() {
		return sdk.NewDec(1), types.ErrInvalidMsgQuota()
	}
	return appMsgQuota, nil
}

func (bm BandwidthManager) getGeneralMsgQuota(ctx sdk.Context) (sdk.Dec, sdk.Error) {
	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}

	curMaxMPS := sdk.NewDec(1)
	if params.ExpectedMaxMPS.GT(bandwidthInfo.MaxMPS) {
		curMaxMPS = params.ExpectedMaxMPS
	} else {
		curMaxMPS = bandwidthInfo.MaxMPS
	}

	generalMsgQuota := params.GeneralMsgQuotaRatio.Mul(curMaxMPS)
	if !generalMsgQuota.IsPositive() {
		return sdk.NewDec(1), types.ErrInvalidMsgQuota()
	}
	return generalMsgQuota, nil
}
