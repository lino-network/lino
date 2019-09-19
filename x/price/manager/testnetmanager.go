package manager

import (
	"github.com/lino-network/lino/types"
)

type TestnetPriceManager struct {
}

func (tm TestnetPriceManager) CoinToMiniDollar(coin types.Coin) (bought types.MiniDollar) {
	return coinToMiniDollar(coin, types.TestnetPrice)
}

// convert minidollar to coin
func (tm TestnetPriceManager) MiniDollarToCoin(dollar types.MiniDollar) (bought types.Coin, used types.MiniDollar) {
	bought, used = miniDollarToCoin(dollar, types.TestnetPrice)
	return
}
