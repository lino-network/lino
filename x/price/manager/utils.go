package manager

import (
	"github.com/lino-network/lino/types"
)

// price: 1 coin = ? minidollar
func coinToMiniDollar(coin types.Coin, price types.MiniDollar) (bought types.MiniDollar) {
	return price.Multiply(types.NewMiniDollarFromInt(coin.Amount))
}

// convert minidollar to coin
// price: 1 coin = ? minidollar.
func miniDollarToCoin(dollar types.MiniDollar, price types.MiniDollar) (bought types.Coin, used types.MiniDollar) {
	c := dollar.Quo(price.Int)
	bought = types.NewCoin(c)
	used = types.NewMiniDollarFromInt(c.Mul(price.Int))
	return
}
