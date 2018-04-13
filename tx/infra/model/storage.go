package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var InfraProviderSubstore = []byte("InfraProvider/")
var InfraProviderListSubstore = []byte("InfraProvider/InfraProviderListKey")

type InfraProviderStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewInfraProviderStorage(key sdk.StoreKey) *InfraProviderStorage {
	cdc := wire.NewCodec()
	storage := InfraProviderStorage{
		key: key,
		cdc: cdc,
	}
	return &storage
}

func (vs InfraProviderStorage) GetInfraProvider(ctx sdk.Context, accKey types.AccountKey) (*InfraProvider, sdk.Error) {
	store := ctx.KVStore(vs.key)
	providerByte := store.Get(GetInfraProviderKey(accKey))
	if providerByte == nil {
		return nil, ErrGetInfraProvider()
	}
	provider := new(InfraProvider)
	if err := vs.cdc.UnmarshalJSON(providerByte, provider); err != nil {
		return nil, ErrInfraProviderUnmarshalError(err)
	}
	return provider, nil
}

func (vs InfraProviderStorage) SetInfraProvider(ctx sdk.Context, accKey types.AccountKey, InfraProvider *InfraProvider) sdk.Error {
	store := ctx.KVStore(vs.key)
	InfraProviderByte, err := vs.cdc.MarshalJSON(*InfraProvider)
	if err != nil {
		return ErrInfraProviderMarshalError(err)
	}
	store.Set(GetInfraProviderKey(accKey), InfraProviderByte)
	return nil
}

func (vs InfraProviderStorage) GetInfraProviderList(ctx sdk.Context) (*InfraProviderList, sdk.Error) {
	store := ctx.KVStore(vs.key)
	listByte := store.Get(GetInfraProviderListKey())
	if listByte == nil {
		return nil, ErrGetInfraProviderList()
	}
	lst := new(InfraProviderList)
	if err := vs.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrInfraProviderUnmarshalError(err)
	}
	return lst, nil
}

func (vs InfraProviderStorage) SetInfraProviderList(ctx sdk.Context, lst *InfraProviderList) sdk.Error {
	store := ctx.KVStore(vs.key)
	listByte, err := vs.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrSetInfraProviderList()
	}
	store.Set(GetInfraProviderListKey(), listByte)
	return nil
}

func GetInfraProviderKey(accKey types.AccountKey) []byte {
	return append(InfraProviderSubstore, accKey...)
}

func GetInfraProviderListKey() []byte {
	return InfraProviderListSubstore
}
