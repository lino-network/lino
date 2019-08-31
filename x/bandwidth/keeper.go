package bandwidth

//go:generate mockery -name BandwidthKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lino-network/lino/x/bandwidth/manager"
)

type BandwidthKeeper interface {
	IsUserMsgFeeEnough(ctx sdk.Context, fee auth.StdFee) bool
	AddMsgSignedByApp(ctx sdk.Context, num int64) sdk.Error
	AddMsgSignedByUser(ctx sdk.Context, num int64) sdk.Error
	ClearBlockInfo(ctx sdk.Context) sdk.Error
	UpdateMaxMPSAndEMA(ctx sdk.Context) sdk.Error
	CalculateCurMsgFee(ctx sdk.Context) sdk.Error
	InitGenesis(ctx sdk.Context) error
	DecayMaxMPS(ctx sdk.Context) sdk.Error
}

var _ BandwidthKeeper = manager.BandwidthManager{}
