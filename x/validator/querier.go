package validator

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "validator"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryValidator     = "validator"
	QueryValidatorList = "valList"
)

// creates a querier for validator REST endpoints
func NewQuerier(vm ValidatorManager) sdk.Querier {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryValidator:
			return queryValidator(ctx, cdc, path[1:], req, vm)
		case QueryValidatorList:
			return queryValidatorList(ctx, cdc, path[1:], req, vm)
		default:
			return nil, sdk.ErrUnknownRequest("unknown validator query endpoint")
		}
	}
}

func queryValidator(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, vm ValidatorManager) ([]byte, sdk.Error) {
	if err := types.CheckPathContentAndMinLength(path, 1); err != nil {
		return nil, err
	}
	validator, err := vm.storage.GetValidator(ctx, types.AccountKey(path[0]))
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(validator)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}

func queryValidatorList(ctx sdk.Context, cdc *wire.Codec, path []string, req abci.RequestQuery, vm ValidatorManager) ([]byte, sdk.Error) {
	validatorList, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return nil, err
	}
	res, marshalErr := cdc.MarshalJSON(validatorList)
	if marshalErr != nil {
		return nil, ErrQueryFailed()
	}
	return res, nil
}
