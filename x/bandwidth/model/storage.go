package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

var (
	bandwidthInfoSubstore = []byte{0x00}
	blockInfoSubstore     = []byte{0x01}
	appBandwidthSubstore  = []byte{0x02}
)

// BandwidthStorage - bandwidth storage
type BandwidthStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewBandwidthStorage(key sdk.StoreKey) BandwidthStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)

	return BandwidthStorage{
		key: key,
		cdc: cdc,
	}
}

// GetBandwidthInfo - returns bandwidth info, returns error otherwise.
func (bs BandwidthStorage) GetBandwidthInfo(ctx sdk.Context) (*BandwidthInfo, sdk.Error) {
	store := ctx.KVStore(bs.key)
	infoByte := store.Get(GetBandwidthInfoKey())
	if infoByte == nil {
		return nil, ErrBandwidthInfoNotFound()
	}
	info := new(BandwidthInfo)
	bs.cdc.MustUnmarshalBinaryLengthPrefixed(infoByte, info)
	return info, nil
}

// SetBandwidthInfo - sets bandwidth info, returns error if any.
func (bs BandwidthStorage) SetBandwidthInfo(ctx sdk.Context, info *BandwidthInfo) sdk.Error {
	store := ctx.KVStore(bs.key)
	infoByte := bs.cdc.MustMarshalBinaryLengthPrefixed(*info)
	store.Set(GetBandwidthInfoKey(), infoByte)
	return nil
}

// GetBlockInfo - returns cur block info, returns error otherwise.
func (bs BandwidthStorage) GetBlockInfo(ctx sdk.Context) (*BlockInfo, sdk.Error) {
	store := ctx.KVStore(bs.key)
	infoByte := store.Get(GetBlockInfoKey())
	if infoByte == nil {
		return nil, ErrBlockInfoNotFound()
	}
	info := new(BlockInfo)
	bs.cdc.MustUnmarshalBinaryLengthPrefixed(infoByte, info)
	return info, nil
}

// SetBlockInfo - sets cur block info, returns error if any.
func (bs BandwidthStorage) SetBlockInfo(ctx sdk.Context, info *BlockInfo) sdk.Error {
	store := ctx.KVStore(bs.key)
	infoByte := bs.cdc.MustMarshalBinaryLengthPrefixed(*info)
	store.Set(GetBlockInfoKey(), infoByte)
	return nil
}

func (bs BandwidthStorage) GetAppBandwidthInfo(ctx sdk.Context, accKey linotypes.AccountKey) (*AppBandwidthInfo, sdk.Error) {
	store := ctx.KVStore(bs.key)
	infoByte := store.Get(GetAppBandwidthInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAppBandwidthInfoNotFound()
	}
	info := new(AppBandwidthInfo)
	bs.cdc.MustUnmarshalBinaryLengthPrefixed(infoByte, info)
	return info, nil
}

func (bs BandwidthStorage) SetAppBandwidthInfo(ctx sdk.Context, accKey linotypes.AccountKey, info *AppBandwidthInfo) sdk.Error {
	store := ctx.KVStore(bs.key)
	infoByte := bs.cdc.MustMarshalBinaryLengthPrefixed(*info)
	store.Set(GetAppBandwidthInfoKey(accKey), infoByte)
	return nil
}

// GetAllAppBandwidthInfo
func (bs BandwidthStorage) GetAllAppBandwidthInfo(ctx sdk.Context) ([]*AppBandwidthInfo, sdk.Error) {
	allInfo := make([]*AppBandwidthInfo, 0)
	store := ctx.KVStore(bs.key)
	iter := sdk.KVStorePrefixIterator(store, appBandwidthSubstore)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		val := iter.Value()
		info := new(AppBandwidthInfo)
		bs.cdc.MustUnmarshalBinaryLengthPrefixed(val, info)
		allInfo = append(allInfo, info)
	}
	return allInfo, nil
}

func (bs BandwidthStorage) DoesAppBandwidthInfoExist(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	store := ctx.KVStore(bs.key)
	return store.Has(GetAppBandwidthInfoKey(accKey))
}

func GetBandwidthInfoKey() []byte {
	return bandwidthInfoSubstore
}

func GetBlockInfoKey() []byte {
	return blockInfoSubstore
}

// GetAppBandwidthInfoKey - "app bandwidth substore" + "username"
func GetAppBandwidthInfoKey(accKey linotypes.AccountKey) []byte {
	return append(appBandwidthSubstore, accKey...)
}
