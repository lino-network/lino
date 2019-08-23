package reputation

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "reputation"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryReputation = "rep"
)

// creates a querier for vote REST endpoints
func NewQuerier(rm ReputationKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryReputation:
			return queryReputation(ctx, cdc, path[1:], req, rm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown reputation query endpoint")
		}
	}
}

func queryReputation(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, rm ReputationKeeper) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	reputation, err := rm.GetReputation(ctx, types.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(reputation)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
