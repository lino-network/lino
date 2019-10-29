package price

import (
	"strings"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/price/types"
)

// creates a querier for price REST endpoints
func NewQuerier(pm PriceKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryPriceCurrent:
			return utils.NewQueryResolver(0, func(args ...string) (interface{}, sdk.Error) {
				return pm.CurrPrice(ctx)
			})(ctx, cdc, path)
		case types.QueryPriceHistory:
			return utils.NewQueryResolver(0, func(args ...string) (interface{}, sdk.Error) {
				return pm.HistoryPrice(ctx), nil
			})(ctx, cdc, path)
		case types.QueryLastFeed:
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				return pm.LastFeed(ctx, linotypes.AccountKey(args[0]))
			})(ctx, cdc, path)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint:" + strings.Join(path, "/"))
		}
	}
}
