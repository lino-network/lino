package vote

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for vote REST endpoints
func NewQuerier(vk VoteKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryVoter:
			return queryVoter(ctx, cdc, path[1:], req, vk)
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
