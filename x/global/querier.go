package global

import (
	"strconv"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "global"

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryTimeEventList   = "timeEventList"
	QueryGlobalMeta      = "globalMeta"
	QueryInflationPool   = "inflationPool"
	QueryConsumptionMeta = "consumptionMeta"
	QueryTPS             = "tps"
	QueryLinoStakeStat   = "linoStakeStat"
)

// creates a querier for global REST endpoints
func NewQuerier(gm GlobalManager) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryTimeEventList:
			return queryTimeEventList(ctx, cdc, path[1:], req, gm)
		case QueryGlobalMeta:
			return queryGlobalMeta(ctx, cdc, path[1:], req, gm)
		case QueryConsumptionMeta:
			return queryConsumptionMeta(ctx, cdc, path[1:], req, gm)
		case QueryTPS:
			return queryTPS(ctx, cdc, path[1:], req, gm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown global query endpoint")
		}
	}
}

func queryTimeEventList(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, gm GlobalManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	unixTime, convertErr := strconv.ParseInt(path[0], 10, 64)
	if convertErr == nil {
		return nil, ErrQueryFailed()
	}
	eventList, err := gm.storage.GetTimeEventList(ctx, unixTime)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(eventList)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryGlobalMeta(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, gm GlobalManager) ([]byte, sdk.Error) {
	postMeta, err := gm.storage.GetGlobalMeta(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(postMeta)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryInflationPool(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, gm GlobalManager) ([]byte, sdk.Error) {
	inflationPool, err := gm.storage.GetInflationPool(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(inflationPool)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryConsumptionMeta(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, gm GlobalManager) ([]byte, sdk.Error) {
	consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(consumptionMeta)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryTPS(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, gm GlobalManager) ([]byte, sdk.Error) {
	tps, err := gm.storage.GetTPS(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(tps)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
