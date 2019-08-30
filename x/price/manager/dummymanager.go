package manager

import (
	"github.com/lino-network/lino/types"
)

type DummyPriceManager struct {
}

func (d DummyPriceManager) CoinToMiniDollar(coin types.Coin) types.MiniDollar {
	return types.NewMiniDollarFromBig(coin.Amount.BigInt())
}

func (d DummyPriceManager) MiniDollarToCoin(dollar types.MiniDollar) types.Coin {
	return types.NewCoinFromBigInt(dollar.BigInt())
}
