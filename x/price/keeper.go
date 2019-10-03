package price

//go:generate mockery -name PriceKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
	linotypes "github.com/lino-network/lino/types"
)

// PriceKeeper - conversion between Coin/MiniDollar at current consensus price.
type PriceKeeper interface {
	// set initial price of LINO
	InitGenesis(ctx sdk.Context, initPrice linotypes.MiniDollar) sdk.Error

	// feed price.
	FeedPrice(ctx sdk.Context, validator linotypes.AccountKey, price linotypes.MiniDollar) sdk.Error

	// UpdatePrice is the hourly event.
	UpdatePrice(ctx sdk.Context) sdk.Error

	// CoinToMiniDollar - convert minidollar to coin
	// since internally, every coin have a price of minidollar, so any amount of coin
	// can all be converted into minidollar.
	CoinToMiniDollar(ctx sdk.Context, coin types.Coin) (bought types.MiniDollar, err sdk.Error)

	// MiniDollarToCoin - return the maximum coins that @p dollar can buy and
	// the amount of dollar used. The returned value is a pair of (new token, used previous token).
	// As there is a minimum price of coin, for dollars that are less than price of one coin
	// they are not used.
	MiniDollarToCoin(ctx sdk.Context, dollar types.MiniDollar) (bought types.Coin, used types.MiniDollar, err sdk.Error)
}
