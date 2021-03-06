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
	InitGenesis(ctx sdk.Context) error
	DecayMaxMPS(ctx sdk.Context) sdk.Error
	ReCalculateAppBandwidthInfo(ctx sdk.Context) sdk.Error
	CheckBandwidth(ctx sdk.Context, addr sdk.AccAddress, fee auth.StdFee) sdk.Error
	EndBlocker(ctx sdk.Context) sdk.Error
	BeginBlocker(ctx sdk.Context) sdk.Error

	// getter
	GetBandwidthInfo(ctx sdk.Context) (*model.BandwidthInfo, sdk.Error)
	GetBlockInfo(ctx sdk.Context) (*model.BlockInfo, sdk.Error)
	GetAppBandwidthInfo(ctx sdk.Context, accKey linotypes.AccountKey) (*model.AppBandwidthInfo, sdk.Error)
}

var _ BandwidthKeeper = manager.BandwidthManager{}
