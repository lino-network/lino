package param

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "param"

	// RouterKey is the message route for param
	RouterKey = ModuleName

	// QuerierRoute is the querier route for param
	QuerierRoute = ModuleName

	QueryAllocationParam = "allocation"
	QueryDeveloperParam  = "developer"
	QueryVoteParam       = "vote"
	QueryProposalParam   = "proposal"
	QueryValidatorParam  = "validator"
	QueryBandwidthParam  = "bandwidth"
	QueryAccountParam    = "account"
	QueryPostParam       = "post"
	QueryReputationParam = "reputation"
	QueryPriceParam      = "price"
)

// creates a querier for account REST endpoints
func NewQuerier(ph ParamHolder) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryAllocationParam:
			return queryAllocationParam(ctx, cdc, path[1:], req, ph)
		case QueryDeveloperParam:
			return queryDeveloperParam(ctx, cdc, path[1:], req, ph)
		case QueryVoteParam:
			return queryVoteParam(ctx, cdc, path[1:], req, ph)
		case QueryProposalParam:
			return queryProposalParam(ctx, cdc, path[1:], req, ph)
		case QueryValidatorParam:
			return queryValidatorParam(ctx, cdc, path[1:], req, ph)
		case QueryBandwidthParam:
			return queryBandwidthParam(ctx, cdc, path[1:], req, ph)
		case QueryAccountParam:
			return queryAccountParam(ctx, cdc, path[1:], req, ph)
		case QueryPostParam:
			return queryPostParam(ctx, cdc, path[1:], req, ph)
		case QueryReputationParam:
			return queryReputationParam(ctx, cdc, path[1:], req, ph)
		case QueryPriceParam:
			return queryPriceParam(ctx, cdc, path[1:], req, ph)
		default:
			return nil, sdk.ErrUnknownRequest("unknown param query endpoint")
		}
	}
}

func queryAllocationParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	globalAllocationParam, err := ph.GetGlobalAllocationParam(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(globalAllocationParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryDeveloperParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	devParam, err := ph.GetDeveloperParam(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(devParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryVoteParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	voteParam, err := ph.GetVoteParam(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(voteParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryProposalParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	proposalParam, err := ph.GetProposalParam(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(proposalParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryValidatorParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	valParam := ph.GetValidatorParam(ctx)
	res, marshalErr := cdc.MarshalJSON(valParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryBandwidthParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	bandwidthParam, err := ph.GetBandwidthParam(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(bandwidthParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryAccountParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	accParam, err := ph.GetAccountParam(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(accParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryPostParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	postParam, err := ph.GetPostParam(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(postParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryReputationParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	repParam := ph.GetReputationParam(ctx)
	res, marshalErr := cdc.MarshalJSON(repParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryPriceParam(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, ph ParamHolder) ([]byte, sdk.Error) {
	priceParam := ph.GetPriceParam(ctx)
	res, marshalErr := cdc.MarshalJSON(priceParam)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
