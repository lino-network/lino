package model

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
)

func NewPriceDumper(store PriceStorage) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc)
	dumper.RegisterType(&FedPrice{}, "lino/price/fedprice", fedPriceSubStore)
	dumper.RegisterType(&[]TimePrice{}, "lino/price/history", priceHistorySubStore)
	dumper.RegisterType(&TimePrice{}, "lino/price/current", currentPriceSubStore)
	dumper.RegisterType(&[]linotypes.AccountKey{}, "lino/price/lastvals", lastValidatorsSubStore)
	return dumper
}
