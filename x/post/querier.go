package post

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "post"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryPostInfo           = "info"
	QueryPostMeta           = "meta"
	QueryPostReportOrUpvote = "reportOrUpvote"
	QueryPostComment        = "comment"
	QueryPostView           = "view"
)

// creates a querier for post REST endpoints
func NewQuerier(pm PostManager) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryPostInfo:
			return queryPostInfo(ctx, cdc, path[1:], req, pm)
		case QueryPostMeta:
			return queryPostMeta(ctx, cdc, path[1:], req, pm)
		case QueryPostReportOrUpvote:
			return queryReportOrUpvote(ctx, cdc, path[1:], req, pm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown post query endpoint")
		}
	}
}

func queryPostInfo(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm PostManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	postInfo, err := pm.postStorage.GetPostInfo(ctx, types.Permlink(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(postInfo)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryPostMeta(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm PostManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	postMeta, err := pm.postStorage.GetPostMeta(ctx, types.Permlink(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(postMeta)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryReportOrUpvote(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm PostManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 2); err != nil {
		return nil, err
	}
	reportOrUpvote, err := pm.postStorage.GetPostReportOrUpvote(ctx, types.Permlink(path[0]), types.AccountKey(path[1]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(reportOrUpvote)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
