package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var (
	infraProviderSubstore     = []byte{0x00}
	infraProviderListSubstore = []byte{0x01}
)

type InfraProviderStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewInfraProviderStorage(key sdk.StoreKey) InfraProviderStorage {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	return InfraProviderStorage{
		key: key,
		cdc: cdc,
	}
}

func (is InfraProviderStorage) InitGenesis(ctx sdk.Context) error {
	if err := is.SetInfraProviderList(ctx, &InfraProviderList{}); err != nil {
		return err
	}
	return nil
}

func (is InfraProviderStorage) DoesInfraProviderExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(is.key)
	return store.Has(GetInfraProviderKey(accKey))
}

func (is InfraProviderStorage) GetInfraProvider(
	ctx sdk.Context, accKey types.AccountKey) (*InfraProvider, sdk.Error) {
	store := ctx.KVStore(is.key)
	providerByte := store.Get(GetInfraProviderKey(accKey))
	if providerByte == nil {
		return nil, ErrGetInfraProvider()
	}
	provider := new(InfraProvider)
	if err := is.cdc.UnmarshalJSON(providerByte, provider); err != nil {
		return nil, ErrInfraProviderUnmarshalError(err)
	}
	return provider, nil
}

func (is InfraProviderStorage) SetInfraProvider(
	ctx sdk.Context, accKey types.AccountKey, InfraProvider *InfraProvider) sdk.Error {
	store := ctx.KVStore(is.key)
	InfraProviderByte, err := is.cdc.MarshalJSON(*InfraProvider)
	if err != nil {
		return ErrInfraProviderMarshalError(err)
	}
	store.Set(GetInfraProviderKey(accKey), InfraProviderByte)
	return nil
}

func (is InfraProviderStorage) GetInfraProviderList(ctx sdk.Context) (*InfraProviderList, sdk.Error) {
	store := ctx.KVStore(is.key)
	listByte := store.Get(GetInfraProviderListKey())
	if listByte == nil {
		return nil, ErrGetInfraProviderList()
	}
	lst := new(InfraProviderList)
	if err := is.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrInfraProviderUnmarshalError(err)
	}
	return lst, nil
}

func (is InfraProviderStorage) SetInfraProviderList(ctx sdk.Context, lst *InfraProviderList) sdk.Error {
	store := ctx.KVStore(is.key)
	listByte, err := is.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrSetInfraProviderList()
	}
	store.Set(GetInfraProviderListKey(), listByte)
	return nil
}

func GetInfraProviderKey(accKey types.AccountKey) []byte {
	return append(infraProviderSubstore, accKey...)
}

func GetInfraProviderListKey() []byte {
	return infraProviderListSubstore
}
