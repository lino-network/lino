package bandwidth

//go:generate mockery -name BandwidthKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type BandwidthKeeper interface {
	IsUserMsgFeeEnough(ctx sdk.Context, fee auth.StdFee) bool
	AddMsgSignedByApp(ctx sdk.Context, num uint32) sdk.Error
	AddMsgSignedByUser(ctx sdk.Context, num uint32) sdk.Error
	ClearCurBlockInfo(ctx sdk.Context) sdk.Error
	UpdateMaxMPSAndEMA(ctx sdk.Context, lastBlockTime int64) sdk.Error
	CalculateCurMsgFee(ctx sdk.Context) sdk.Error
	InitGenesis(ctx sdk.Context) error
}

var _ BandwidthKeeper = BandwidthManager{}
