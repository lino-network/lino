package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
)

var (
	ValidatorSubstore     = []byte{0x00}
	ValidatorListSubstore = []byte{0x01}
)

type ValidatorStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewValidatorStorage(key sdk.StoreKey) ValidatorStorage {
	cdc := wire.New()
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
	if err := vs.cdc.UnmarshalBinaryLengthPrefixed(validatorByte, validator); err != nil {
		return nil, ErrFailedToUnmarshalValidator(err)
	}
	return validator, nil
}

func (vs ValidatorStorage) SetValidator(ctx sdk.Context, accKey types.AccountKey, validator *Validator) sdk.Error {
	store := ctx.KVStore(vs.key)
	validatorByte, err := vs.cdc.MarshalBinaryLengthPrefixed(*validator)
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
	if err := vs.cdc.UnmarshalBinaryLengthPrefixed(listByte, lst); err != nil {
		return nil, ErrFailedToUnmarshalValidatorList(err)
	}
	return lst, nil
}

func (vs ValidatorStorage) SetValidatorList(ctx sdk.Context, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vs.key)
	listByte, err := vs.cdc.MarshalBinaryLengthPrefixed(*lst)
	if err != nil {
		return ErrFailedToMarshalValidatorList(err)
	}
	store.Set(GetValidatorListKey(), listByte)
	return nil
}

// Export state of validators.
// StoreMap - map of all substores
func (vs ValidatorStorage) StoreMap(ctx sdk.Context) utils.StoreMap {
	store := ctx.KVStore(vs.key)
	substores := []utils.SubStore{
		{
			Store:      store,
			Prefix:     ValidatorSubstore,
			ValCreator: func() interface{} { return new(Validator) },
			Decoder:    vs.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
		{
			Store:      store,
			Prefix:     ValidatorListSubstore,
			ValCreator: func() interface{} { return new(ValidatorList) },
			Decoder:    vs.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
	}
	return utils.NewStoreMap(substores)
}

func GetValidatorKey(accKey types.AccountKey) []byte {
	return append(ValidatorSubstore, accKey...)
}

func GetValidatorListKey() []byte {
	return ValidatorListSubstore
}
