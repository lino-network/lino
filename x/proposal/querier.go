package proposal

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for proposal REST endpoints
func NewQuerier(pm ProposalKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryOngoingProposal:
			return queryOngoingProposal(ctx, cdc, path[1:], req, pm)
		case types.QueryExpiredProposal:
			return queryExpiredProposal(ctx, cdc, path[1:], req, pm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown proposal query endpoint")
		}
	}
}

func queryOngoingProposal(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm ProposalKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	proposal, err := pm.GetOngoingProposal(ctx, linotypes.ProposalKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(proposal)
	if marshalErr != nil {
		return nil, linotypes.ErrQueryFailed(marshalErr.Error())
	}
	return res, nil
}

func queryExpiredProposal(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm ProposalKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	proposal, err := pm.GetExpiredProposal(ctx, linotypes.ProposalKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(proposal)
	if marshalErr != nil {
		return nil, linotypes.ErrQueryFailed(marshalErr.Error())
	}
	return res, nil
}
