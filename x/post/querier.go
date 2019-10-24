package post

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/post/types"
)

// creates a querier for post REST endpoints
func NewQuerier(pm PostKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryPostInfo:
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				return pm.GetPost(ctx, linotypes.Permlink(args[0]))
			})(ctx, cdc, path)
		case types.QueryConsumptionWindow:
			return utils.NewQueryResolver(0, func(args ...string) (interface{}, sdk.Error) {
				return pm.GetComsumptionWindow(ctx), nil
			})(ctx, cdc, path)
		default:
			return nil, sdk.ErrUnknownRequest("unknown post query endpoint")
		}
	}
}
