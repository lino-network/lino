package post

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/types"
)

const (
	QueryPostInfo = "info"
)

// creates a querier for post REST endpoints
func NewQuerier(pm PostKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryPostInfo:
			return queryPostInfo(ctx, cdc, path[1:], req, pm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown post query endpoint")
		}
	}
}

func queryPostInfo(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm PostKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	postInfo, err := pm.GetPost(ctx, linotypes.Permlink(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(postInfo)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}
