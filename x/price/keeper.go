package price

//go:generate mockery -name PriceKeeper

import (
	"github.com/lino-network/lino/types"
)

// PriceKeeper - conversion between Coin/MiniDollar at current consensus price.
type PriceKeeper interface {
	// CoinToMiniDollar - convert minidollar to coin
	// since internally, every coin have a price of minidollar, so any amount of coin
	// can all be converted into minidollar.
	CoinToMiniDollar(coin types.Coin) (bought types.MiniDollar)

	// MiniDollarToCoin - return the maximum coins that @p dollar can buy and
	// the amount of dollar used. The returned value is a pair of (new token, used previous token).
	// As there is a minimum price of coin, for dollars that are less than price of one coin
	// they are not used.
	MiniDollarToCoin(dollar types.MiniDollar) (bought types.Coin, used types.MiniDollar)
}
