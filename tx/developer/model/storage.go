package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var DeveloperSubstore = []byte("Developer/")
var DeveloperListSubstore = []byte("Developer/DeveloperListKey")

type DeveloperStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewDeveloperStorage(key sdk.StoreKey) *DeveloperStorage {
	cdc := wire.NewCodec()
	storage := DeveloperStorage{
		key: key,
		cdc: cdc,
	}
	return &storage
}

func (vs DeveloperStorage) GetDeveloper(ctx sdk.Context, accKey types.AccountKey) (*Developer, sdk.Error) {
	store := ctx.KVStore(vs.key)
	providerByte := store.Get(GetDeveloperKey(accKey))
	if providerByte == nil {
		return nil, ErrGetDeveloper()
	}
	provider := new(Developer)
	if err := vs.cdc.UnmarshalJSON(providerByte, provider); err != nil {
		return nil, ErrDeveloperUnmarshalError(err)
	}
	return provider, nil
}

func (vs DeveloperStorage) SetDeveloper(ctx sdk.Context, accKey types.AccountKey, Developer *Developer) sdk.Error {
	store := ctx.KVStore(vs.key)
	DeveloperByte, err := vs.cdc.MarshalJSON(*Developer)
	if err != nil {
		return ErrDeveloperMarshalError(err)
	}
	store.Set(GetDeveloperKey(accKey), DeveloperByte)
	return nil
}

func (vs DeveloperStorage) GetDeveloperList(ctx sdk.Context) (*DeveloperList, sdk.Error) {
	store := ctx.KVStore(vs.key)
	listByte := store.Get(GetDeveloperListKey())
	if listByte == nil {
		return nil, ErrGetDeveloperList()
	}
	lst := new(DeveloperList)
	if err := vs.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrDeveloperUnmarshalError(err)
	}
	return lst, nil
}

func (vs DeveloperStorage) SetDeveloperList(ctx sdk.Context, lst *DeveloperList) sdk.Error {
	store := ctx.KVStore(vs.key)
	listByte, err := vs.cdc.MarshalJSON(*lst)
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
