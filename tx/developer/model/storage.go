package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var (
	DeveloperSubstore     = []byte("Developer/")
	DeveloperListSubstore = []byte("Developer/DeveloperListKey")
)

type DeveloperStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewDeveloperStorage(key sdk.StoreKey) DeveloperStorage {
	cdc := wire.NewCodec()
	return DeveloperStorage{
		key: key,
		cdc: cdc,
	}
}

func (ds DeveloperStorage) InitGenesis(ctx sdk.Context) error {
	if err := ds.SetDeveloperList(ctx, &DeveloperList{}); err != nil {
		return err
	}
	return nil
}

func (ds DeveloperStorage) GetDeveloper(
	ctx sdk.Context, accKey types.AccountKey) (*Developer, sdk.Error) {
	store := ctx.KVStore(ds.key)
	providerByte := store.Get(GetDeveloperKey(accKey))
	if providerByte == nil {
		return nil, ErrGetDeveloper()
	}
	provider := new(Developer)
	if err := ds.cdc.UnmarshalJSON(providerByte, provider); err != nil {
		return nil, ErrDeveloperUnmarshalError(err)
	}
	return provider, nil
}

func (ds DeveloperStorage) SetDeveloper(
	ctx sdk.Context, accKey types.AccountKey, Developer *Developer) sdk.Error {
	store := ctx.KVStore(ds.key)
	DeveloperByte, err := ds.cdc.MarshalJSON(*Developer)
	if err != nil {
		return ErrDeveloperMarshalError(err)
	}
	store.Set(GetDeveloperKey(accKey), DeveloperByte)
	return nil
}

func (ds DeveloperStorage) DeleteDeveloper(ctx sdk.Context, username types.AccountKey) sdk.Error {
	store := ctx.KVStore(ds.key)
	store.Delete(GetDeveloperKey(username))
	return nil
}

func (ds DeveloperStorage) GetDeveloperList(ctx sdk.Context) (*DeveloperList, sdk.Error) {
	store := ctx.KVStore(ds.key)
	listByte := store.Get(GetDeveloperListKey())
	if listByte == nil {
		return nil, ErrGetDeveloperList()
	}
	lst := new(DeveloperList)
	if err := ds.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrDeveloperUnmarshalError(err)
	}
	return lst, nil
}

func (ds DeveloperStorage) SetDeveloperList(ctx sdk.Context, lst *DeveloperList) sdk.Error {
	store := ctx.KVStore(ds.key)
	listByte, err := ds.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrSetDeveloperList()
	}
	store.Set(GetDeveloperListKey(), listByte)
	return nil
}

func GetDeveloperKey(accKey types.AccountKey) []byte {
	return append(DeveloperSubstore, accKey...)
}

func GetDeveloperListKey() []byte {
	return DeveloperListSubstore
}
