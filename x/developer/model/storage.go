package model

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	developerSubstore     = []byte{0x00}
	developerListSubstore = []byte{0x01}
)

// DeveloperStorage - developer storage
type DeveloperStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

// DeveloperStorage - new developer storage
func NewDeveloperStorage(key sdk.StoreKey) DeveloperStorage {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	return DeveloperStorage{
		key: key,
		cdc: cdc,
	}
}

// InitGenesis - initialize developer storage
func (ds DeveloperStorage) InitGenesis(ctx sdk.Context) error {
	if err := ds.SetDeveloperList(ctx, &DeveloperList{}); err != nil {
		return err
	}
	return nil
}

// DoesDeveloperExist - check if developer in KVStore or not
func (ds DeveloperStorage) DoesDeveloperExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(ds.key)
	return store.Has(GetDeveloperKey(accKey))
}

// GetDeveloper - get developer from KVStore
func (ds DeveloperStorage) GetDeveloper(
	ctx sdk.Context, accKey types.AccountKey) (*Developer, sdk.Error) {
	store := ctx.KVStore(ds.key)
	providerByte := store.Get(GetDeveloperKey(accKey))
	if providerByte == nil {
		return nil, ErrDeveloperNotFound()
	}
	provider := new(Developer)
	if err := ds.cdc.UnmarshalJSON(providerByte, provider); err != nil {
		return nil, ErrFailedToUnmarshalDeveloper(err)
	}
	return provider, nil
}

// SetDeveloper - set developer to KVStore
func (ds DeveloperStorage) SetDeveloper(
	ctx sdk.Context, accKey types.AccountKey, developer *Developer) sdk.Error {
	store := ctx.KVStore(ds.key)
	developerByte, err := ds.cdc.MarshalJSON(*developer)
	if err != nil {
		return ErrFailedToMarshalDeveloper(err)
	}
	store.Set(GetDeveloperKey(accKey), developerByte)
	return nil
}

// DeleteDeveloper - delete developer from KVStore
func (ds DeveloperStorage) DeleteDeveloper(ctx sdk.Context, username types.AccountKey) sdk.Error {
	store := ctx.KVStore(ds.key)
	store.Delete(GetDeveloperKey(username))
	return nil
}

// GetDeveloperList - get developer list from KVStore
func (ds DeveloperStorage) GetDeveloperList(ctx sdk.Context) (*DeveloperList, sdk.Error) {
	store := ctx.KVStore(ds.key)
	listByte := store.Get(GetDeveloperListKey())
	if listByte == nil {
		return nil, ErrDeveloperListNotFound()
	}
	lst := new(DeveloperList)
	if err := ds.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrFailedToUnmarshalDeveloperList(err)
	}
	return lst, nil
}

// SetDeveloperList - set developer list to KVStore
func (ds DeveloperStorage) SetDeveloperList(ctx sdk.Context, lst *DeveloperList) sdk.Error {
	store := ctx.KVStore(ds.key)
	listByte, err := ds.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrFailedToMarshalDeveloperList(err)
	}
	store.Set(GetDeveloperListKey(), listByte)
	return nil
}

// Export developer storage state
func (ds DeveloperStorage) Export(ctx sdk.Context) *DeveloperTables {
	tables := &DeveloperTables{}
	store := ctx.KVStore(ds.key)
	// export table.Developers
	func() {
		itr := sdk.KVStorePrefixIterator(store, developerSubstore)
		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			k := itr.Key()
			username := types.AccountKey(k[1:])
			dev, err := ds.GetDeveloper(ctx, username)
			if err != nil {
				panic("failed to read developer: " + err.Error())
			}
			row := DeveloperRow{
				Username:  username,
				Developer: *dev,
			}
			tables.Developers = append(tables.Developers, row)
		}
	}()
	// export table.DeveloperList
	list, err := ds.GetDeveloperList(ctx)
	if err != nil {
		panic("failed to get developer list: " + err.Error())
	}
	tables.DeveloperList = DeveloperListTable{
		List: *list,
	}
	return tables
}

// GetDeveloperKey - "developer substore" + "developer"
func GetDeveloperKey(accKey types.AccountKey) []byte {
	return append(developerSubstore, accKey...)
}

// GetDeveloperListKey - "developerlist substore"
func GetDeveloperListKey() []byte {
	return developerListSubstore
}
