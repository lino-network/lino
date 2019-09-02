package bandwidth

//go:generate mockery -name BandwidthKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/bandwidth/manager"
	"github.com/lino-network/lino/x/bandwidth/model"
)

type BandwidthKeeper interface {
	IsUserMsgFeeEnough(ctx sdk.Context, fee auth.StdFee) bool
	AddMsgSignedByApp(ctx sdk.Context, accKey linotypes.AccountKey, num int64) sdk.Error
	AddMsgSignedByUser(ctx sdk.Context, num int64) sdk.Error
	ClearBlockInfo(ctx sdk.Context) sdk.Error
	UpdateMaxMPSAndEMA(ctx sdk.Context) sdk.Error
	CalculateCurMsgFee(ctx sdk.Context) sdk.Error
	InitGenesis(ctx sdk.Context) error
	DecayMaxMPS(ctx sdk.Context) sdk.Error
	RefillAppBandwidthCredit(ctx sdk.Context, accKey linotypes.AccountKey) sdk.Error
	GetVacancyCoeff(ctx sdk.Context) (sdk.Dec, sdk.Error)
	GetPunishmentCoeff(ctx sdk.Context, accKey linotypes.AccountKey) (sdk.Dec, sdk.Error)
	GetBandwidthCostPerMsg(ctx sdk.Context, u sdk.Dec, p sdk.Dec) sdk.Dec
	ConsumeBandwidthCredit(ctx sdk.Context, costPerMsg sdk.Dec, accKey linotypes.AccountKey) sdk.Error
	ReCalculateAppBandwidthInfo(ctx sdk.Context) sdk.Error
	CheckBandwidth(ctx sdk.Context, accKey linotypes.AccountKey, fee auth.StdFee) sdk.Error
	// getter
	GetAllAppInfo(ctx sdk.Context) ([]*model.AppBandwidthInfo, sdk.Error)
}

var _ BandwidthKeeper = manager.BandwidthManager{}
