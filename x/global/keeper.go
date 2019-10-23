package global

//go:generate mockery -name GlobalKeeper

import (
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global/model"
)

type GlobalKeeper interface {
	AddFrictionAndRegisterContentRewardEvent(
		ctx sdk.Context, event types.Event, friction types.Coin, evaluate types.MiniDollar) sdk.Error
	GetConsumptionFrictionRate(ctx sdk.Context) (sdk.Dec, sdk.Error)
	GetRewardAndPopFromWindow(ctx sdk.Context, evaluate types.MiniDollar) (types.Coin, sdk.Error)
	AddToValidatorInflationPool(ctx sdk.Context, coin types.Coin) sdk.Error
	GetLastBlockTime(ctx sdk.Context) (int64, sdk.Error)
	GetInterestSince(ctx sdk.Context, unixTime int64, linoStake types.Coin) (types.Coin, sdk.Error)
	RegisterCoinReturnEvent(
		ctx sdk.Context, events []types.Event, times int64, intervalSec int64) sdk.Error
	RegisterEventAtTime(ctx sdk.Context, unixTime int64, event types.Event) sdk.Error
	// pop out developer monthly inflation from pool.
	PopDeveloperMonthlyInflation(ctx sdk.Context) (types.Coin, sdk.Error)
	AddLinoStakeToStat(ctx sdk.Context, linoStake types.Coin) sdk.Error
	MinusLinoStakeFromStat(ctx sdk.Context, linoStake types.Coin) sdk.Error
	GetValidatorHourlyInflation(ctx sdk.Context) (types.Coin, sdk.Error)

	// import export
	ValidatorInflationPool(ctx sdk.Context) types.Coin
	DevInflationPool(ctx sdk.Context) types.Coin
	CCInflationPool(ctx sdk.Context) types.Coin
	StakeinPool(ctx sdk.Context) types.Coin
	StakeStats(ctx sdk.Context) (rst []model.LinoStakeStat, days []int64)
	ConsumptionWindow(ctx sdk.Context) types.MiniDollar
	FrictionPool(ctx sdk.Context) types.Coin
	ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
	ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
}

var _ GlobalKeeper = &GlobalManager{}
