package infra

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "infra"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryInfraProvider = "infra"
	QueryInfraList     = "infraList"
)

// creates a querier for infra REST endpoints
func NewQuerier(im InfraManager) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryInfraProvider:
			return queryInfraProvider(ctx, cdc, path[1:], req, im)
		case QueryInfraList:
			return queryInfraList(ctx, cdc, path[1:], req, im)
		default:
			return nil, sdk.ErrUnknownRequest("unknown infra query endpoint")
		}
	}
}

func queryInfraProvider(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, im InfraManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	postInfo, err := im.storage.GetInfraProvider(ctx, types.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(postInfo)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryInfraList(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, im InfraManager) ([]byte, sdk.Error) {
	postMeta, err := im.storage.GetInfraProviderList(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(postMeta)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
