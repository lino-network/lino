package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/validator/types"
)

var (
	ValidatorSubstore        = []byte{0x00}
	ValidatorListSubstore    = []byte{0x01}
	ElectionVoteListSubstore = []byte{0x02}
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

// func (vs ValidatorStorage) DoesValidatorExist(ctx sdk.Context, accKey linotypes.AccountKey) bool {
// 	store := ctx.KVStore(vs.key)
// 	return store.Has(GetValidatorKey(accKey))
// }

func (vs ValidatorStorage) GetValidator(ctx sdk.Context, accKey linotypes.AccountKey) (*Validator, sdk.Error) {
	store := ctx.KVStore(vs.key)
	validatorByte := store.Get(GetValidatorKey(accKey))
	if validatorByte == nil {
		return nil, types.ErrValidatorNotFound(accKey)
	}
	validator := new(Validator)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(validatorByte, validator)
	return validator, nil
}

func (vs ValidatorStorage) SetValidator(ctx sdk.Context, accKey linotypes.AccountKey, validator *Validator) {
	store := ctx.KVStore(vs.key)
	validatorByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*validator)
	store.Set(GetValidatorKey(accKey), validatorByte)
}

func (vs ValidatorStorage) GetValidatorList(ctx sdk.Context) *ValidatorList {
	store := ctx.KVStore(vs.key)
	listByte := store.Get(GetValidatorListKey())
	if listByte == nil {
		panic("Validator List should be initialized during genesis")
	}
	lst := new(ValidatorList)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(listByte, lst)
	return lst
}

func (vs ValidatorStorage) SetValidatorList(ctx sdk.Context, lst *ValidatorList) {
	store := ctx.KVStore(vs.key)
	listByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*lst)
	store.Set(GetValidatorListKey(), listByte)
}

func (vs ValidatorStorage) GetElectionVoteList(ctx sdk.Context, accKey linotypes.AccountKey) *ElectionVoteList {
	store := ctx.KVStore(vs.key)
	lstByte := store.Get(GetElectionVoteListKey(accKey))
	if lstByte == nil {
		// valid empty value.
		return &ElectionVoteList{}
	}
	lst := new(ElectionVoteList)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(lstByte, lst)
	return lst
}

func (vs ValidatorStorage) SetElectionVoteList(ctx sdk.Context, accKey linotypes.AccountKey, lst *ElectionVoteList) {
	store := ctx.KVStore(vs.key)
	lstByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*lst)
	store.Set(GetElectionVoteListKey(accKey), lstByte)
}

func (vs ValidatorStorage) StoreMap(ctx sdk.Context) utils.StoreMap {
	store := ctx.KVStore(vs.key)
	stores := []utils.SubStore{
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
		{
			Store:      store,
			Prefix:     ElectionVoteListSubstore,
			ValCreator: func() interface{} { return new(ElectionVoteList) },
			Decoder:    vs.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
	}
	return utils.NewStoreMap(stores)
}

func GetValidatorKey(accKey linotypes.AccountKey) []byte {
	return append(ValidatorSubstore, accKey...)
}

func GetElectionVoteListKey(accKey linotypes.AccountKey) []byte {
	return append(ElectionVoteListSubstore, accKey...)
}

func GetValidatorListKey() []byte {
	return ValidatorListSubstore
}
