package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/types"
)

var (
	validatorSubstore     = []byte{0x00}
	validatorListSubstore = []byte{0x01}
)

type ValidatorStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewValidatorStorage(key sdk.StoreKey) ValidatorStorage {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	vs := ValidatorStorage{
		key: key,
		cdc: cdc,
	}
	return vs
}

func (vs ValidatorStorage) InitGenesis(ctx sdk.Context) error {
	lst := &ValidatorList{
		LowestPower: types.NewCoinFromInt64(0),
	}

	if err := vs.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vs ValidatorStorage) DoesValidatorExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetValidatorKey(accKey))
}

func (vs ValidatorStorage) GetValidator(ctx sdk.Context, accKey types.AccountKey) (*Validator, sdk.Error) {
	store := ctx.KVStore(vs.key)
	validatorByte := store.Get(GetValidatorKey(accKey))
	if validatorByte == nil {
		return nil, ErrValidatorNotFound()
	}
	validator := new(Validator)
	if err := vs.cdc.UnmarshalJSON(validatorByte, validator); err != nil {
		return nil, ErrFailedToUnmarshalValidator(err)
	}
	return validator, nil
}

func (vs ValidatorStorage) SetValidator(ctx sdk.Context, accKey types.AccountKey, validator *Validator) sdk.Error {
	store := ctx.KVStore(vs.key)
	validatorByte, err := vs.cdc.MarshalJSON(*validator)
	if err != nil {
		return ErrFailedToMarshalValidator(err)
	}
	store.Set(GetValidatorKey(accKey), validatorByte)
	return nil
}

func (vs ValidatorStorage) DeleteValidator(ctx sdk.Context, username types.AccountKey) sdk.Error {
	store := ctx.KVStore(vs.key)
	store.Delete(GetValidatorKey(username))
	return nil
}

func (vs ValidatorStorage) GetValidatorList(ctx sdk.Context) (*ValidatorList, sdk.Error) {
	store := ctx.KVStore(vs.key)
	listByte := store.Get(GetValidatorListKey())
	if listByte == nil {
		return nil, ErrValidatorListNotFound()
	}
	lst := new(ValidatorList)
	if err := vs.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrFailedToUnmarshalValidatorList(err)
	}
	return lst, nil
}

func (vs ValidatorStorage) SetValidatorList(ctx sdk.Context, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vs.key)
	listByte, err := vs.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrFailedToMarshalValidatorList(err)
	}
	store.Set(GetValidatorListKey(), listByte)
	return nil
}

func GetValidatorKey(accKey types.AccountKey) []byte {
	return append(validatorSubstore, accKey...)
}

func GetValidatorListKey() []byte {
	return validatorListSubstore
}
