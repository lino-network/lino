package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/price/types"
)

var (
	fedPriceSubStore       = []byte{0x00} // validator's latest price.
	priceHistorySubStore   = []byte{0x01} // hourly prices
	currentPriceSubStore   = []byte{0x02} // current price
	lastValidatorsSubStore = []byte{0x03} // validators in last update time.
	feedHistorySubStore    = []byte{0x04} // fed history.
)

// GetFedPriceKey - price key.
func GetFedPriceKey(u linotypes.AccountKey) []byte {
	return append(fedPriceSubStore, u...)
}

// GetPriceHistoryKey - hourly price.
func GetPriceHistoryKey() []byte {
	return priceHistorySubStore
}

// GetCurrentPriceKey - get current price.
func GetCurrentPriceKey() []byte {
	return currentPriceSubStore
}

// GetLastValidatorsKey - get last validators.
func GetLastValidatorsKey() []byte {
	return lastValidatorsSubStore
}

// GetLastValidatorsKey - get last validators.
func GetFeedHistoryKey() []byte {
	return feedHistorySubStore
}

// PriceStorage - price storage
type PriceStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewPriceStorage - returns a new PriceStorage, binary encoded.
func NewPriceStorage(key sdk.StoreKey) PriceStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)

	return PriceStorage{
		key: key,
		cdc: cdc,
	}
}

// GetFedPrice - get fed price of validator from KVStore
func (ps PriceStorage) GetFedPrice(ctx sdk.Context, val linotypes.AccountKey) (*FedPrice, sdk.Error) {
	store := ctx.KVStore(ps.key)
	key := GetFedPriceKey(val)
	infoByte := store.Get(key)
	if infoByte == nil {
		return nil, types.ErrFedPriceNotFound(val)
	}
	price := new(FedPrice)
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(infoByte, price)
	return price, nil
}

// SetFedPrice - set fed price to KVStore
func (ps PriceStorage) SetFedPrice(ctx sdk.Context, price *FedPrice) {
	store := ctx.KVStore(ps.key)
	bytes := ps.cdc.MustMarshalBinaryLengthPrefixed(*price)
	store.Set(GetFedPriceKey(price.Validator), bytes)
}

// GetPriceHistory - return price history.
func (ps PriceStorage) GetPriceHistory(ctx sdk.Context) []TimePrice {
	store := ctx.KVStore(ps.key)
	key := GetPriceHistoryKey()
	bytes := store.Get(key)
	if bytes == nil {
		return nil
	}
	price := make([]TimePrice, 0)
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &price)
	return price
}

// SetPriceHistory - set price history.
func (ps PriceStorage) SetPriceHistory(ctx sdk.Context, prices []TimePrice) {
	store := ctx.KVStore(ps.key)
	bytes := ps.cdc.MustMarshalBinaryLengthPrefixed(prices)
	store.Set(GetPriceHistoryKey(), bytes)
}

// GetCurrentPrice - return current price
func (ps PriceStorage) GetCurrentPrice(ctx sdk.Context) (*TimePrice, sdk.Error) {
	store := ctx.KVStore(ps.key)
	key := GetCurrentPriceKey()
	bytes := store.Get(key)
	if bytes == nil {
		return nil, types.ErrCurrentPriceNotFound()
	}
	price := new(TimePrice)
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, price)
	return price, nil
}

func (ps PriceStorage) SetCurrentPrice(ctx sdk.Context, price *TimePrice) {
	store := ctx.KVStore(ps.key)
	bytes := ps.cdc.MustMarshalBinaryLengthPrefixed(price)
	store.Set(GetCurrentPriceKey(), bytes)
}

func (ps PriceStorage) GetLastValidators(ctx sdk.Context) []linotypes.AccountKey {
	store := ctx.KVStore(ps.key)
	key := GetLastValidatorsKey()
	bytes := store.Get(key)
	if bytes == nil {
		return nil
	}
	vals := make([]linotypes.AccountKey, 0)
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &vals)
	return vals
}

func (ps PriceStorage) SetLastValidators(ctx sdk.Context, last []linotypes.AccountKey) {
	store := ctx.KVStore(ps.key)
	bytes := ps.cdc.MustMarshalBinaryLengthPrefixed(last)
	store.Set(GetLastValidatorsKey(), bytes)
}

func (ps PriceStorage) GetFeedHistory(ctx sdk.Context) []FeedHistory {
	store := ctx.KVStore(ps.key)
	key := GetFeedHistoryKey()
	bytes := store.Get(key)
	if bytes == nil {
		return nil
	}
	history := make([]FeedHistory, 0)
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &history)
	return history
}

func (ps PriceStorage) SetFeedHistory(ctx sdk.Context, history []FeedHistory) {
	store := ctx.KVStore(ps.key)
	bytes := ps.cdc.MustMarshalBinaryLengthPrefixed(history)
	store.Set(GetFeedHistoryKey(), bytes)
}
