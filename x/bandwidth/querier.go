package bandwidth

import (
	linotypes "github.com/lino-network/lino/types"

	"github.com/lino-network/lino/x/bandwidth/types"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryBandwidthInfo    = "bandwidthinfo"
	QueryBlockInfo        = "blockinfo"
	QueryAppBandwidthInfo = "appinfo"
)

// creates a querier for account REST endpoints
func NewQuerier(bm BandwidthKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBandwidthInfo:
			return queryBandwidthInfo(ctx, cdc, req, bm)
		case QueryBlockInfo:
			return queryBlockInfo(ctx, cdc, req, bm)
		case QueryAppBandwidthInfo:
			return queryAppBandwidthInfo(ctx, cdc, path[1:], req, bm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown bandwidth query endpoint")
		}
	}
}

func queryBandwidthInfo(ctx sdk.Context, cdc *wire.Codec, req abci.RequestQuery, bm BandwidthKeeper) ([]byte, sdk.Error) {
	bandwidthInfo, err := bm.GetBandwidthInfo(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(bandwidthInfo)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryBlockInfo(ctx sdk.Context, cdc *wire.Codec, req abci.RequestQuery, bm BandwidthKeeper) ([]byte, sdk.Error) {
	blockInfo, err := bm.GetBlockInfo(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(blockInfo)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryAppBandwidthInfo(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, bm BandwidthKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	appInfo, err := bm.GetAppBandwidthInfo(ctx, linotypes.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(appInfo)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}
