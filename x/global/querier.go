package global

import (
	// "strconv"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	// linotypes "github.com/lino-network/lino/types"
	// "github.com/lino-network/lino/x/global/types"
)

// creates a querier for global REST endpoints
func NewQuerier(gm GlobalKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		// case types.QueryTimeEventList:
		// 	return queryTimeEventList(ctx, cdc, path[1:], req, gm)
		// case types.QueryGlobalTime:
		// 	return queryGlobalTime(ctx, cdc, path[1:], req, gm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown global query endpoint")
		}
	}
}

// func queryTimeEventList(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, gm GlobalKeeper) ([]byte, sdk.Error) {
// 	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
// 		return nil, err
// 	}
// 	unixTime, convertErr := strconv.ParseInt(path[0], 10, 64)
// 	if convertErr != nil {
// 		return nil, types.ErrQueryFailed()
// 	}
// 	eventList, err := gm.storage.GetTimeEventList(ctx, unixTime)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, marshalErr := cdc.MarshalJSON(eventList)
// 	if marshalErr != nil {
// 		return nil, ErrQueryFailed()
// 	}
// 	return res, nil
// }

// func queryGlobalTime(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, gm GlobalManager) ([]byte, sdk.Error) {
// 	globalTime, err := gm.storage.GetGlobalTime(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, marshalErr := cdc.MarshalJSON(globalTime)
// 	if marshalErr != nil {
// 		return nil, ErrQueryFailed()
// 	}
// 	return res, nil
// }
