package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

type TestnetPriceManager struct {
}

func (tm TestnetPriceManager) CoinToMiniDollar(ctx sdk.Context, coin linotypes.Coin) (bought linotypes.MiniDollar, err sdk.Error) {
	return coinToMiniDollar(coin, linotypes.TestnetPrice), nil
}

// convert minidollar to coin
func (tm TestnetPriceManager) MiniDollarToCoin(ctx sdk.Context, dollar linotypes.MiniDollar) (bought linotypes.Coin, used linotypes.MiniDollar, err sdk.Error) {
	bought, used = miniDollarToCoin(dollar, linotypes.TestnetPrice)
	return bought, used, nil
}

func (tm TestnetPriceManager) InitGenesis(ctx sdk.Context, initPrice linotypes.MiniDollar) sdk.Error {
	return nil
}

func (tm TestnetPriceManager) FeedPrice(ctx sdk.Context, validator linotypes.AccountKey, price linotypes.MiniDollar) sdk.Error {
	return nil
}

func (tm TestnetPriceManager) UpdatePrice(ctx sdk.Context) sdk.Error {
	return nil
}
