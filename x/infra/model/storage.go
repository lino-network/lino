package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

var (
	infraProviderSubstore     = []byte{0x00}
	infraProviderListSubstore = []byte{0x01}
)

// InfraProviderStorage - infra provider storage
type InfraProviderStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewInfraProviderStorage - create a new infra provider storage
func NewInfraProviderStorage(key sdk.StoreKey) InfraProviderStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return InfraProviderStorage{
		key: key,
		cdc: cdc,
	}
}

// InitGenesis - initialize infra provider manager
func (is InfraProviderStorage) InitGenesis(ctx sdk.Context) error {
	if err := is.SetInfraProviderList(ctx, &InfraProviderList{}); err != nil {
		return err
	}
	return nil
}

// DoesInfraProviderExist - check infra provider exists in KVStore or not
func (is InfraProviderStorage) DoesInfraProviderExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(is.key)
	return store.Has(GetInfraProviderKey(accKey))
}

// GetInfraProvider - get infra provider from KVStore
func (is InfraProviderStorage) GetInfraProvider(
	ctx sdk.Context, accKey types.AccountKey) (*InfraProvider, sdk.Error) {
	store := ctx.KVStore(is.key)
	providerByte := store.Get(GetInfraProviderKey(accKey))
	if providerByte == nil {
		return nil, ErrInfraProviderNotFound()
	}
	provider := new(InfraProvider)
	if err := is.cdc.UnmarshalJSON(providerByte, provider); err != nil {
		return nil, ErrFailedToUnmarshalInfraProvider(err)
	}
	return provider, nil
}

// SetInfraProvider - set infra provider to KVStore
func (is InfraProviderStorage) SetInfraProvider(
	ctx sdk.Context, accKey types.AccountKey, InfraProvider *InfraProvider) sdk.Error {
	store := ctx.KVStore(is.key)
	InfraProviderByte, err := is.cdc.MarshalJSON(*InfraProvider)
	if err != nil {
		return ErrFailedToMarshalInfraProvider(err)
	}
	store.Set(GetInfraProviderKey(accKey), InfraProviderByte)
	return nil
}

// GetInfraProviderList - get infra provider list from KVStore
func (is InfraProviderStorage) GetInfraProviderList(ctx sdk.Context) (*InfraProviderList, sdk.Error) {
	store := ctx.KVStore(is.key)
	listByte := store.Get(GetInfraProviderListKey())
	if listByte == nil {
		return nil, ErrInfraProviderListNotFound()
	}
	lst := new(InfraProviderList)
	if err := is.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrFailedToUnmarshalInfraProviderList(err)
	}
	return lst, nil
}

// SetInfraProviderList - set infra provider list to KVStore
func (is InfraProviderStorage) SetInfraProviderList(ctx sdk.Context, lst *InfraProviderList) sdk.Error {
	store := ctx.KVStore(is.key)
	listByte, err := is.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrFailedToMarshalInfraProviderList(err)
	}
	store.Set(GetInfraProviderListKey(), listByte)
	return nil
}

// GetInfraProviderKey - get infra provider key in infra provider substore
func GetInfraProviderKey(accKey types.AccountKey) []byte {
	return append(infraProviderSubstore, accKey...)
}

// GetInfraProviderListKey - get infra provider list key in infra provider list substore
func GetInfraProviderListKey() []byte {
	return infraProviderListSubstore
}
