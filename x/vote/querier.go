package vote

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "vote"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryVoter         = "voter"
	QueryVote          = "vote"
	QueryReferenceList = "refList"
)

// creates a querier for vote REST endpoints
func NewQuerier(vm VoteManager) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryVoter:
			return queryVoter(ctx, cdc, path[1:], req, vm)
		case QueryVote:
			return queryVote(ctx, cdc, path[1:], req, vm)
		case QueryReferenceList:
			return queryReferenceList(ctx, cdc, path[1:], req, vm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown vote query endpoint")
		}
	}
}

func queryVoter(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, vm VoteManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	voter, err := vm.storage.GetVoter(ctx, types.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(voter)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryVote(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, vm VoteManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 2); err != nil {
		return nil, err
	}
	vote, err := vm.storage.GetVote(ctx, types.ProposalKey(path[0]), types.AccountKey(path[1]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(vote)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryReferenceList(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, vm VoteManager) ([]byte, sdk.Error) {
	referenceList, err := vm.storage.GetReferenceList(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(referenceList)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
