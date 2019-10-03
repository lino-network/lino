package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	account "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/bandwidth/model"
	"github.com/lino-network/lino/x/bandwidth/types"
	developer "github.com/lino-network/lino/x/developer"
	global "github.com/lino-network/lino/x/global"
	vote "github.com/lino-network/lino/x/vote"
)

var BandwidthManagerTestMode bool = false

// BandwidthManager - bandwidth manager
type BandwidthManager struct {
	storage     model.BandwidthStorage
	paramHolder param.ParamKeeper
	// deps
	gm global.GlobalKeeper
	vm vote.VoteKeeper
	dm developer.DeveloperKeeper
	am account.AccountKeeper
}

func NewBandwidthManager(key sdk.StoreKey, holder param.ParamKeeper, gm global.GlobalKeeper, vm vote.VoteKeeper, dm developer.DeveloperKeeper, am account.AccountKeeper) *BandwidthManager {
	return &BandwidthManager{
		storage:     model.NewBandwidthStorage(key),
		paramHolder: holder,
		gm:          gm,
		vm:          vm,
		dm:          dm,
		am:          am,
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
		CurU:                 sdk.NewDec(1),
	}

	if err := bm.storage.SetBlockInfo(ctx, blockInfo); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) PrecheckAndConsumeBandwidthCredit(ctx sdk.Context, accKey linotypes.AccountKey) sdk.Error {
	appBandwidthInfo, err := bm.storage.GetAppBandwidthInfo(ctx, accKey)
	if err != nil {
		return err
	}

	blockInfo, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}
	if appBandwidthInfo.CurBandwidthCredit.LT(blockInfo.CurU) {
		return types.ErrAppBandwidthNotEnough()
	}

	// consume u at first
	appBandwidthInfo.CurBandwidthCredit = appBandwidthInfo.CurBandwidthCredit.Sub(blockInfo.CurU)
	if err := bm.storage.SetAppBandwidthInfo(ctx, accKey, appBandwidthInfo); err != nil {
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

func (bm BandwidthManager) AddMsgSignedByApp(ctx sdk.Context, accKey linotypes.AccountKey, num int64) sdk.Error {
	info, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}
	info.TotalMsgSignedByApp += num

	if err := bm.storage.SetBlockInfo(ctx, info); err != nil {
		return err
	}

	appBandwidthInfo, err := bm.storage.GetAppBandwidthInfo(ctx, accKey)
	if err != nil {
		return err
	}
	appBandwidthInfo.MessagesInCurBlock += 1
	if err := bm.storage.SetAppBandwidthInfo(ctx, accKey, appBandwidthInfo); err != nil {
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
	generalMPS := linotypes.NewDecFromRat(blockInfo.TotalMsgSignedByUser, pastTime)
	bandwidthInfo.GeneralMsgEMA = bm.calculateEMA(bandwidthInfo.GeneralMsgEMA, params.GeneralMsgEMAFactor, generalMPS)

	// EMA_app = EMA_app_prev * (1 - k_app) + appMPS * k_app
	appMPS := linotypes.NewDecFromRat(blockInfo.TotalMsgSignedByApp, pastTime)
	bandwidthInfo.AppMsgEMA = bm.calculateEMA(bandwidthInfo.AppMsgEMA, params.AppMsgEMAFactor, appMPS)

	// MaxMPS = max( (totalMsgSignedByUser + totalMsgSignedByApp)/(curBlockTime - lastBlockTime), MaxMPS)
	totalMPS := linotypes.NewDecFromRat(blockInfo.TotalMsgSignedByUser+blockInfo.TotalMsgSignedByApp, pastTime)
	if totalMPS.GT(bandwidthInfo.MaxMPS) {
		bandwidthInfo.MaxMPS = totalMPS
	}

	if err := bm.storage.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}

	return nil
}

// calcuate the current msg fee based on last block info at the beginning of each block
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

// calcuate the current vacancy coeef u based on last block info at the beginning of each block
func (bm BandwidthManager) CalculateCurU(ctx sdk.Context) sdk.Error {
	bandwidthInfo, err := bm.storage.GetBandwidthInfo(ctx)
	if err != nil {
		return err
	}

	params, err := bm.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	appMsgQuota, err := bm.getAppMsgQuota(ctx)
	if err != nil {
		return err
	}

	blockInfo, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}

	delta := bandwidthInfo.AppMsgEMA.Sub(appMsgQuota)
	blockInfo.CurU = bm.approximateExp(delta.Quo(appMsgQuota).Mul(params.AppVacancyFactor))
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

	if err := bm.storage.SetAppBandwidthInfo(ctx, accKey, info); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) GetPunishmentCoeff(ctx sdk.Context, accKey linotypes.AccountKey) (sdk.Dec, sdk.Error) {
	lastBlockTime, err := bm.gm.GetLastBlockTime(ctx)
	if err != nil {
		return sdk.NewDec(1), err
	}

	pastTime := ctx.BlockHeader().Time.Unix() - lastBlockTime
	if pastTime <= 0 {
		return sdk.NewDec(1), nil
	}
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

func (bm BandwidthManager) ConsumeBandwidthCredit(ctx sdk.Context, u sdk.Dec, p sdk.Dec, accKey linotypes.AccountKey) sdk.Error {
	info, err := bm.storage.GetAppBandwidthInfo(ctx, accKey)
	if err != nil {
		return err
	}
	numMsgs := sdk.NewDec(info.MessagesInCurBlock)
	// add back pre-check consumed bandwidth credit
	info.CurBandwidthCredit = info.CurBandwidthCredit.Add(numMsgs.Mul(u))
	// consume credit, currently allow the credit be negative
	costPerMsg := bm.GetBandwidthCostPerMsg(ctx, u, p)
	info.CurBandwidthCredit = info.CurBandwidthCredit.Sub(numMsgs.Mul(costPerMsg))
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
			// refill the bandwidth credit with old expectedMPS (refill rate)
			if err := bm.RefillAppBandwidthCredit(ctx, app.Username); err != nil {
				return err
			}
			// update it's max and expected MPS
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

func (bm BandwidthManager) CheckBandwidth(ctx sdk.Context, addr sdk.AccAddress, fee auth.StdFee) sdk.Error {
	bank, err := bm.am.GetBankByAddress(ctx, addr)
	if err != nil {
		return err
	}
	if bank.Username != "" {
		appName, err := bm.dm.GetAffiliatingApp(ctx, bank.Username)
		if err == nil {
			// refill bandwidth for apps with messages in current block
			if err := bm.RefillAppBandwidthCredit(ctx, appName); err != nil {
				return err
			}

			// app bandwidth model
			if err := bm.PrecheckAndConsumeBandwidthCredit(ctx, appName); err != nil {
				return err
			}

			// add app message stats
			if err := bm.AddMsgSignedByApp(ctx, appName, 1); err != nil {
				return err
			}
			return nil
		}
	}

	// msg fee for general message
	if !bm.IsUserMsgFeeEnough(ctx, fee) {
		return types.ErrUserMsgFeeNotEnough()
	}

	// minus message fee
	info, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}

	//  minus msg fee
	if !BandwidthManagerTestMode {
		if err := bm.am.MinusCoinFromAddress(ctx, addr, info.CurMsgFee); err != nil {
			return err
		}

		if err := bm.gm.AddToValidatorInflationPool(ctx, info.CurMsgFee); err != nil {
			return err
		}
	}
	// add general message stats
	if err := bm.AddMsgSignedByUser(ctx, 1); err != nil {
		return err
	}

	return nil
}

func (bm BandwidthManager) BeginBlocker(ctx sdk.Context) sdk.Error {
	// calculate the new general msg fee for the current block
	if err := bm.CalculateCurMsgFee(ctx); err != nil {
		return err
	}

	// calculate the new vacancy coeff
	if err := bm.CalculateCurU(ctx); err != nil {
		return err
	}

	// clear stats for block info
	if err := bm.ClearBlockInfo(ctx); err != nil {
		return err
	}
	return nil
}

func (bm BandwidthManager) EndBlocker(ctx sdk.Context) sdk.Error {
	// update maxMPS and EMA for different msgs and store cur block info
	if err := bm.UpdateMaxMPSAndEMA(ctx); err != nil {
		return err
	}

	blockInfo, err := bm.storage.GetBlockInfo(ctx)
	if err != nil {
		return err
	}

	// get all app bandwidth info
	allInfo, err := bm.GetAllAppInfo(ctx)
	if err != nil {
		return err
	}

	for _, info := range allInfo {
		if info.MessagesInCurBlock == 0 {
			continue
		}
		// calculate cost and consume bandwidth credit
		p, err := bm.GetPunishmentCoeff(ctx, info.Username)
		if err != nil {
			return err
		}
		if err := bm.ConsumeBandwidthCredit(ctx, blockInfo.CurU, p, info.Username); err != nil {
			return err
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

	var curMaxMPS sdk.Dec
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

	var curMaxMPS sdk.Dec
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

// getter
func (bm BandwidthManager) GetBandwidthInfo(ctx sdk.Context) (*model.BandwidthInfo, sdk.Error) {
	return bm.storage.GetBandwidthInfo(ctx)
}

func (bm BandwidthManager) GetBlockInfo(ctx sdk.Context) (*model.BlockInfo, sdk.Error) {
	return bm.storage.GetBlockInfo(ctx)
}

func (bm BandwidthManager) GetAppBandwidthInfo(ctx sdk.Context, accKey linotypes.AccountKey) (*model.AppBandwidthInfo, sdk.Error) {
	return bm.storage.GetAppBandwidthInfo(ctx, accKey)
}
