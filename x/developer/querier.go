package developer

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/model"
	"github.com/lino-network/lino/x/developer/types"
)

// creates a querier for developer REST endpoints
func NewQuerier(dm DeveloperKeeper) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryDeveloper:
			return queryDeveloper(ctx, cdc, path[1:], req, dm)
		case types.QueryDeveloperList:
			return queryDeveloperList(ctx, cdc, path[1:], req, dm)
		case types.QueryIDA:
			return queryIDA(ctx, cdc, path[1:], req, dm)
		case types.QueryIDABalance:
			return queryIDABalance(ctx, cdc, path[1:], req, dm)
		case types.QueryAffiliated:
			return queryAffiliated(ctx, cdc, path[1:], req, dm)
		default:
			return nil, sdk.ErrUnknownRequest(
				fmt.Sprintf("unknown developer query endpoint: %s", path[0]),
			)
		}
	}
}

func queryDeveloper(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, dm DeveloperKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	developer, err := dm.GetDeveloper(ctx, linotypes.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(developer)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryDeveloperList(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, dm DeveloperKeeper) ([]byte, sdk.Error) {
	developers := dm.GetLiveDevelopers(ctx)
	devmap := make(map[string]model.Developer)
	for _, dev := range developers {
		devmap[string(dev.Username)] = dev
	}
	res, marshalErr := cdc.MarshalJSON(devmap)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryIDABalance(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, dm DeveloperKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 2); err != nil {
		return nil, err
	}
	dev := linotypes.AccountKey(path[0])
	user := linotypes.AccountKey(path[1])
	price, err := dm.GetMiniIDAPrice(ctx, dev)
	if err != nil {
		return nil, err
	}
	bank, err := dm.GetIDABank(ctx, dev, user)
	if err != nil {
		return nil, err
	}

	idaAmount := bank.Balance.Quo(price.Int).ToDec().Quo(sdk.NewDec(linotypes.Decimals)).String()
	res, marshalErr := cdc.MarshalJSON(types.QueryResultIDABalance{
		Amount:   idaAmount,
		Unauthed: bank.Unauthed,
	})
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryIDA(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, dm DeveloperKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	app := linotypes.AccountKey(path[0])
	ida, err := dm.GetIDA(ctx, app)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(ida)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}

func queryAffiliated(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, dm DeveloperKeeper) ([]byte, sdk.Error) {
	if err := linotypes.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	app := linotypes.AccountKey(path[0])
	accounts := dm.GetAffiliated(ctx, app)
	res, marshalErr := cdc.MarshalJSON(accounts)
	if marshalErr != nil {
		return nil, types.ErrQueryFailed()
	}
	return res, nil
}
