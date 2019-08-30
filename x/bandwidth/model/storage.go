package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	bandwidthInfoSubstore = []byte{0x00}
	lastBlockInfoSubstore = []byte{0x01}
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

func (bs BandwidthStorage) InitGenesis(ctx sdk.Context) error {
	bandwidthInfo := &BandwidthInfo{
		GeneralMsgEMA: sdk.NewDec(0),
		AppMsgEMA:     sdk.NewDec(0),
		MaxMPS:        sdk.NewDec(0),
	}

	if err := bs.SetBandwidthInfo(ctx, bandwidthInfo); err != nil {
		return err
	}

	lastBlockInfo := &LastBlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
	}

	if err := bs.SetLastBlockInfo(ctx, lastBlockInfo); err != nil {
		return err
	}

	return nil
}

// GetBandwidthInfo - returns bandwidth info, returns error otherwise.
func (bs BandwidthStorage) GetBandwidthInfo(ctx sdk.Context) (*BandwidthInfo, sdk.Error) {
	store := ctx.KVStore(bs.key)
	infoByte := store.Get(GetBandwidthInfoKey())
	if infoByte == nil {
		return nil, ErrBandwidthInfoNotFound()
	}
	info := new(BandwidthInfo)
	if err := bs.cdc.UnmarshalBinaryLengthPrefixed(infoByte, info); err != nil {
		return nil, ErrFailedToUnmarshalBandwidthInfo(err)
	}
	return info, nil
}

// SetBandwidthInfo - sets bandwidth info, returns error if any.
func (bs BandwidthStorage) SetBandwidthInfo(ctx sdk.Context, info *BandwidthInfo) sdk.Error {
	store := ctx.KVStore(bs.key)
	infoByte, err := bs.cdc.MarshalBinaryLengthPrefixed(*info)
	if err != nil {
		return ErrFailedToMarshalBandwidthInfo(err)
	}
	store.Set(GetBandwidthInfoKey(), infoByte)
	return nil
}

// GetLastBlockInfo - returns cur block info, returns error otherwise.
func (bs BandwidthStorage) GetLastBlockInfo(ctx sdk.Context) (*LastBlockInfo, sdk.Error) {
	store := ctx.KVStore(bs.key)
	infoByte := store.Get(GetLastBlockInfoKey())
	if infoByte == nil {
		return nil, ErrLastBlockInfoNotFound()
	}
	info := new(LastBlockInfo)
	if err := bs.cdc.UnmarshalBinaryLengthPrefixed(infoByte, info); err != nil {
		return nil, ErrFailedToUnmarshalLastBlockInfo(err)
	}
	return info, nil
}

// SetLastBlockInfo - sets cur block info, returns error if any.
func (bs BandwidthStorage) SetLastBlockInfo(ctx sdk.Context, info *LastBlockInfo) sdk.Error {
	store := ctx.KVStore(bs.key)
	infoByte, err := bs.cdc.MarshalBinaryLengthPrefixed(*info)
	if err != nil {
		return ErrFailedToMarshalLastBlockInfo(err)
	}
	store.Set(GetLastBlockInfoKey(), infoByte)
	return nil
}

func GetBandwidthInfoKey() []byte {
	return bandwidthInfoSubstore
}

func GetLastBlockInfoKey() []byte {
	return lastBlockInfoSubstore
}
