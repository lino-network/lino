package price

import (
	"strings"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/price/types"
)

// creates a querier for price REST endpoints
func NewQuerier(pm PriceKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryPriceCurrent:
			return queryPriceCurrent(ctx, cdc, path[1:], req, pm)
		case types.QueryPriceHistory:
			return queryPriceHistory(ctx, cdc, path[1:], req, pm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint:" + strings.Join(path, "/"))
		}
	}
}

func queryPriceCurrent(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm PriceKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 0); err != nil {
		return nil, err
	}
	price, err := pm.CurrPrice(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(price)
	if marshalErr != nil {
		return nil, linotypes.ErrQueryFailed(marshalErr.Error())
	}
	return res, nil
}

func queryPriceHistory(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm PriceKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 0); err != nil {
		return nil, err
	}
	price := pm.HistoryPrice(ctx)
	res, marshalErr := cdc.MarshalJSON(price)
	if marshalErr != nil {
		return nil, linotypes.ErrQueryFailed(marshalErr.Error())
	}
	return res, nil
}
