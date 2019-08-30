package global

//go:generate mockery -name GlobalKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

type GlobalKeeper interface {
	AddFrictionAndRegisterContentRewardEvent(
		ctx sdk.Context, event types.Event, friction types.Coin, evaluate types.MiniDollar) sdk.Error
	GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Dec, sdk.Error)
	GetRewardAndPopFromWindow(ctx sdk.Context, evaluate types.MiniDollar) (types.Coin, sdk.Error)
	AddToValidatorInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error
	GetLastBlockTime(ctx sdk.Context) (int64, sdk.Error)
}

var _ GlobalKeeper = &GlobalManager{}
