package proposal

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "proposal"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryNextProposal    = "next"
	QueryOngoingProposal = "ongoing"
	QueryExpiredProposal = "expired"
)

// creates a querier for proposal REST endpoints
func NewQuerier(pm ProposalManager) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryOngoingProposal:
			return queryOngoingProposal(ctx, cdc, path[1:], req, pm)
		case QueryExpiredProposal:
			return queryExpiredProposal(ctx, cdc, path[1:], req, pm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown proposal query endpoint")
		}
	}
}

func queryOngoingProposal(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm ProposalManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	proposal, err := pm.storage.GetOngoingProposal(ctx, types.ProposalKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(proposal)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryExpiredProposal(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, pm ProposalManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	proposal, err := pm.storage.GetExpiredProposal(ctx, types.ProposalKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(proposal)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
