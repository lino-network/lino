package vote

import (
	"strconv"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/vote/types"
)

// creates a querier for vote REST endpoints
func NewQuerier(vk VoteKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryVoter:
			return queryVoter(ctx, cdc, path[1:], req, vk)
		case types.QueryStakeStats:
			return utils.NewQueryResolver(1, func(args ...string) (interface{}, sdk.Error) {
				day, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return nil, linotypes.ErrInvalidQueryPath()
				}
				return vk.GetStakeStatsOfDay(ctx, day)
			})(ctx, cdc, path)
		default:
			return nil, sdk.ErrUnknownRequest("unknown vote query endpoint")
		}
	}
}

func queryVoter(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, vk VoteKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	voter, err := vk.GetVoter(ctx, linotypes.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(voter)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}
