package manager

import (
	"github.com/lino-network/lino/types"
)

type TestnetPriceManager struct {
}

func (tm TestnetPriceManager) CoinToMiniDollar(coin types.Coin) (bought types.MiniDollar) {
	return types.TestnetPrice.Multiply(types.NewMiniDollarFromInt(coin.Amount))
}

// convert minidollar to coin
func (tm TestnetPriceManager) MiniDollarToCoin(dollar types.MiniDollar) (bought types.Coin, used types.MiniDollar) {
	c := dollar.Quo(types.TestnetPrice.Int)
	bought = types.NewCoin(c)
	used = types.NewMiniDollarFromInt(c.Mul(types.TestnetPrice.Int))
	return
}
