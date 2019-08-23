package price

//go:generate mockery -name PriceKeeper

import (
	"github.com/lino-network/lino/types"
)

type PriceKeeper interface {
	// convert coin to MiniDollar at current consensus price.
	CoinToMiniDollar(coin types.Coin) types.MiniDollar
	// convert minidollar to coin
	MiniDollarToCoin(dollar types.MiniDollar) types.Coin
}
