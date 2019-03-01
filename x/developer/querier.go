package developer

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/model"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "developer"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryDeveloper     = "dev"
	QueryDeveloperList = "devList"
)

// creates a querier for developer REST endpoints
func NewQuerier(dm DeveloperManager) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryDeveloper:
			return queryDeveloper(ctx, cdc, path[1:], req, dm)
		case QueryDeveloperList:
			return queryDeveloperList(ctx, cdc, path[1:], req, dm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown developer query endpoint")
		}
	}
}

func queryDeveloper(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, dm DeveloperManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	developer, err := dm.storage.GetDeveloper(ctx, types.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(developer)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryDeveloperList(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, dm DeveloperManager) ([]byte, sdk.Error) {
	developerList, err := dm.storage.GetDeveloperList(ctx)
	if err != nil {
		return nil, err
	}
	developers := make(map[string]*model.Developer)
	for _, username := range developerList.AllDevelopers {
		developer, err := dm.storage.GetDeveloper(ctx, username)
		if err != nil {
			return nil, err
		}
		developers[string(username)] = developer
	}
	res, marshalErr := cdc.MarshalJSON(developers)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
