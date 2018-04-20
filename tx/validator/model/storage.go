package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var ValidatorSubstore = []byte("Validator/")
var ValidatorListSubstore = []byte("ValidatorList/ValidatorListKey")

type ValidatorStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewValidatorStorage(key sdk.StoreKey) ValidatorStorage {
	cdc := wire.NewCodec()
	vs := ValidatorStorage{
		key: key,
		cdc: cdc,
	}
	return vs
}

func (vs ValidatorStorage) InitGenesis(ctx sdk.Context) error {
	lst := &ValidatorList{
		LowestPower: types.Coin{0},
	}

	if err := vs.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vs ValidatorStorage) GetValidator(ctx sdk.Context, accKey types.AccountKey) (*Validator, sdk.Error) {
	store := ctx.KVStore(vs.key)
	validatorByte := store.Get(GetValidatorKey(accKey))
	if validatorByte == nil {
		return nil, ErrGetValidator()
	}
	validator := new(Validator)
	if err := vs.cdc.UnmarshalJSON(validatorByte, validator); err != nil {
		return nil, ErrValidatorUnmarshalError(err)
	}
	return validator, nil
}

func (vs ValidatorStorage) SetValidator(ctx sdk.Context, accKey types.AccountKey, validator *Validator) sdk.Error {
	store := ctx.KVStore(vs.key)
	validatorByte, err := vs.cdc.MarshalJSON(*validator)
	if err != nil {
		return ErrValidatorMarshalError(err)
	}
	store.Set(GetValidatorKey(accKey), validatorByte)
	return nil
}

func (vs ValidatorStorage) GetValidatorList(ctx sdk.Context) (*ValidatorList, sdk.Error) {
	store := ctx.KVStore(vs.key)
	listByte := store.Get(GetValidatorListKey())
	if listByte == nil {
		return nil, ErrGetValidatorList()
	}
	lst := new(ValidatorList)
	if err := vs.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrValidatorUnmarshalError(err)
	}
	return lst, nil
}

func (vs ValidatorStorage) SetValidatorList(ctx sdk.Context, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vs.key)
	listByte, err := vs.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrSetValidatorList()
	}
	store.Set(GetValidatorListKey(), listByte)
	return nil
}

func GetValidatorKey(accKey types.AccountKey) []byte {
	return append(ValidatorSubstore, accKey...)
}

func GetValidatorListKey() []byte {
	return ValidatorListSubstore
}
