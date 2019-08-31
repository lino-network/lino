package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	bandwidthInfoSubstore = []byte{0x00}
	blockInfoSubstore     = []byte{0x01}
	mpsSubStore           = []byte{0x02}
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

func GetBandwidthInfoKey() []byte {
	return bandwidthInfoSubstore
}

func GetBlockInfoKey() []byte {
	return blockInfoSubstore
}
