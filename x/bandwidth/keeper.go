package bandwidth

//go:generate mockery -name BandwidthKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lino-network/lino/types"
)

type BandwidthKeeper interface {
	IsAppBandwidthEnough(ctx sdk.Context, username types.AccountKey) bool
	IsUserMsgFeeEnough(ctx sdk.Context, username types.AccountKey, fee auth.StdFee) bool
	AddMsgSignedByApp(ctx sdk.Context, num uint32) sdk.Error
	AddMsgSignedByUser(ctx sdk.Context, num uint32) sdk.Error
	ClearCurBlockInfo(ctx sdk.Context) sdk.Error
	UpdateEMA(ctx sdk.Context) sdk.Error
}

var _ BandwidthKeeper = BandwidthManager{}
