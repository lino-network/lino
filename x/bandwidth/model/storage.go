package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	bandwidthInfoSubstore = []byte{0x00}
	curBlockInfoSubstore  = []byte{0x01}
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
	info := &BandwidthInfo{
		GeneralMsgEMA: sdk.NewDec(0),
		AppMsgEMA:     sdk.NewDec(0),
	}

	if err := bs.SetBandwidthInfo(ctx, info); err != nil {
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

// GetCurBlockInfo - returns cur block info, returns error otherwise.
func (bs BandwidthStorage) GetCurBlockInfo(ctx sdk.Context) (*CurBlockInfo, sdk.Error) {
	store := ctx.KVStore(bs.key)
	infoByte := store.Get(GetCurBlockInfoKey())
	if infoByte == nil {
		return nil, ErrCurBlockInfoNotFound()
	}
	info := new(CurBlockInfo)
	if err := bs.cdc.UnmarshalBinaryLengthPrefixed(infoByte, info); err != nil {
		return nil, ErrFailedToUnmarshalCurBlockInfo(err)
	}
	return info, nil
}

// SetCurBlockInfo - sets cur block info, returns error if any.
func (bs BandwidthStorage) SetCurBlockInfo(ctx sdk.Context, info *CurBlockInfo) sdk.Error {
	store := ctx.KVStore(bs.key)
	infoByte, err := bs.cdc.MarshalBinaryLengthPrefixed(*info)
	if err != nil {
		return ErrFailedToMarshalCurBlockInfo(err)
	}
	store.Set(GetCurBlockInfoKey(), infoByte)
	return nil
}

func GetBandwidthInfoKey() []byte {
	return bandwidthInfoSubstore
}

func GetCurBlockInfoKey() []byte {
	return curBlockInfoSubstore
}
