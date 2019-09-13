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
	if err := is.cdc.UnmarshalBinaryLengthPrefixed(providerByte, provider); err != nil {
		return nil, ErrFailedToUnmarshalInfraProvider(err)
	}
	return provider, nil
}

// SetInfraProvider - set infra provider to KVStore
func (is InfraProviderStorage) SetInfraProvider(
	ctx sdk.Context, accKey types.AccountKey, infraProvider *InfraProvider) sdk.Error {
	store := ctx.KVStore(is.key)
	InfraProviderByte, err := is.cdc.MarshalBinaryLengthPrefixed(*infraProvider)
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
	if err := is.cdc.UnmarshalBinaryLengthPrefixed(listByte, lst); err != nil {
		return nil, ErrFailedToUnmarshalInfraProviderList(err)
	}
	return lst, nil
}

// SetInfraProviderList - set infra provider list to KVStore
func (is InfraProviderStorage) SetInfraProviderList(ctx sdk.Context, lst *InfraProviderList) sdk.Error {
	store := ctx.KVStore(is.key)
	listByte, err := is.cdc.MarshalBinaryLengthPrefixed(*lst)
	if err != nil {
		return ErrFailedToMarshalInfraProviderList(err)
	}
	store.Set(GetInfraProviderListKey(), listByte)
	return nil
}

// Export - infra state
func (is InfraProviderStorage) Export(ctx sdk.Context) *InfraTables {
	tables := &InfraTables{}
	store := ctx.KVStore(is.key)
	// export table.providers
	func() {
		itr := sdk.KVStorePrefixIterator(store, infraProviderSubstore)
		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			k := itr.Key()
			username := types.AccountKey(k[1:])
			provider, err := is.GetInfraProvider(ctx, username)
			if err != nil {
				panic("failed to read developer: " + err.Error())
			}
			row := InfraProviderRow{
				App:      username,
				Provider: *provider,
			}
			tables.InfraProviders = append(tables.InfraProviders, row)
		}
	}()
	// export table.DeveloperList
	list, err := is.GetInfraProviderList(ctx)
	if err != nil {
		panic("failed to get developer list: " + err.Error())
	}
	tables.InfraProviderList = InfraProviderListRow{
		List: *list,
	}
	return tables
}

// Import from tablesIR.
func (is InfraProviderStorage) Import(ctx sdk.Context, tb *InfraTablesIR) {
	check := func(e error) {
		if e != nil {
			panic("[is] Failed to import: " + e.Error())
		}
	}
	// import table.Providers
	for _, v := range tb.InfraProviders {
		err := is.SetInfraProvider(ctx, v.App, &v.Provider)
		check(err)
	}
	// import ProviderList
	err := is.SetInfraProviderList(ctx, &tb.InfraProviderList.List)
	check(err)
}

// GetInfraProviderKey - get infra provider key in infra provider substore
func GetInfraProviderKey(accKey types.AccountKey) []byte {
	return append(infraProviderSubstore, accKey...)
}

// GetInfraProviderListKey - get infra provider list key in infra provider list substore
func GetInfraProviderListKey() []byte {
	return infraProviderListSubstore
}
