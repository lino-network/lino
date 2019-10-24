package utils

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

// query store
type StoreQuerier = func(args ...string) (interface{}, sdk.Error)

// return query result
type QueryResolver = func(ctx sdk.Context, cdc *wire.Codec, path []string) ([]byte, sdk.Error)

func NewQueryResolver(numArgs int, resolver StoreQuerier) QueryResolver {
	return func(ctx sdk.Context, cdc *wire.Codec, path []string) ([]byte, sdk.Error) {
		if len(path) < 1 {
			return nil, linotypes.ErrQueryFailed("invalid query")
		}
		substore := path[0]
		path = path[1:]
		if err := linotypes.CheckPathContentAndMinLength(path, numArgs); err != nil {
			return nil, err
		}
		rst, err := resolver(path...)
		if err != nil {
			return nil, err
		}
		res, marshalErr := cdc.MarshalJSON(rst)
		if marshalErr != nil {
			return nil, linotypes.ErrQueryFailed(fmt.Sprintf("query %s failed", substore))
		}
		return res, nil
	}
}
