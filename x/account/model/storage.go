package model

import (
	"strings"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/account/types"
)

var (
	AccountInfoSubstore        = []byte{0x00}
	AccountBankSubstore        = []byte{0x01}
	AccountMetaSubstore        = []byte{0x02}
	AccountGrantPubKeySubstore = []byte{0x03}
	AccountPoolSubstore        = []byte{0x04}
	AccountSupplySubstore      = []byte{0x05}
)

// AccountStorage - account storage
type AccountStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts
	cdc *wire.Codec
}

// NewLinoAccountStorage - creates and returns a account manager
func NewAccountStorage(key sdk.StoreKey) AccountStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	cdc.Seal()

	return AccountStorage{
		key: key,
		cdc: cdc,
	}
}

// DoesAccountExist - returns true when a specific account exist in the KVStore.
func (as AccountStorage) DoesAccountExist(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	store := ctx.KVStore(as.key)
	return store.Has(GetAccountInfoKey(accKey))
}

// GetInfo - returns general account info of a specific account, returns error otherwise.
func (as AccountStorage) GetInfo(ctx sdk.Context, accKey linotypes.AccountKey) (*AccountInfo, sdk.Error) {
	store := ctx.KVStore(as.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, types.ErrAccountNotFound(accKey)
	}
	info := new(AccountInfo)
	as.cdc.MustUnmarshalBinaryLengthPrefixed(infoByte, info)
	return info, nil
}

// SetInfo - sets general account info to a specific account, returns error if any.
func (as AccountStorage) SetInfo(ctx sdk.Context, accInfo *AccountInfo) {
	store := ctx.KVStore(as.key)
	infoByte := as.cdc.MustMarshalBinaryLengthPrefixed(*accInfo)
	store.Set(GetAccountInfoKey(accInfo.Username), infoByte)
}

// GetBank - returns bank info of a specific address, returns error if any.
func (as AccountStorage) GetBank(ctx sdk.Context, addr sdk.Address) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(as.key)
	bankByte := store.Get(GetAccountBankKey(addr))
	if bankByte == nil {
		return nil, types.ErrAccountBankNotFound(addr)
	}
	bank := new(AccountBank)
	as.cdc.MustUnmarshalBinaryLengthPrefixed(bankByte, bank)
	return bank, nil
}

// SetBank - sets bank info for a given address, returns error if any.
func (as AccountStorage) SetBank(ctx sdk.Context, addr sdk.Address, accBank *AccountBank) {
	store := ctx.KVStore(as.key)
	bankByte := as.cdc.MustMarshalBinaryLengthPrefixed(*accBank)
	store.Set(GetAccountBankKey(addr), bankByte)
}

// GetMeta - returns meta of a given account that are tiny and frequently updated fields.
func (as AccountStorage) GetMeta(ctx sdk.Context, accKey linotypes.AccountKey) *AccountMeta {
	store := ctx.KVStore(as.key)
	metaByte := store.Get(GetAccountMetaKey(accKey))
	if metaByte == nil {
		return &AccountMeta{
			JSONMeta: "",
		}
	}
	meta := new(AccountMeta)
	as.cdc.MustUnmarshalBinaryLengthPrefixed(metaByte, meta)
	return meta
}

// SetMeta - sets meta for a given account, returns error if any.
func (as AccountStorage) SetMeta(ctx sdk.Context, accKey linotypes.AccountKey, accMeta *AccountMeta) {
	store := ctx.KVStore(as.key)
	metaByte := as.cdc.MustMarshalBinaryLengthPrefixed(*accMeta)
	store.Set(GetAccountMetaKey(accKey), metaByte)
}

// DeleteAllGrantPermissions - deletes all grant pubkeys from a granted user in KV.
func (as AccountStorage) DeleteAllGrantPermissions(ctx sdk.Context, me linotypes.AccountKey, grantTo linotypes.AccountKey) {
	store := ctx.KVStore(as.key)
	store.Delete(getGrantPermKey(me, grantTo))
	return
}

// GetGrantPermissions - returns grant user info keyed with pubkey.
func (as AccountStorage) GetGrantPermissions(ctx sdk.Context, me linotypes.AccountKey, grantTo linotypes.AccountKey) ([]*GrantPermission, sdk.Error) {
	store := ctx.KVStore(as.key)
	grantPubKeyByte := store.Get(getGrantPermKey(me, grantTo))
	if grantPubKeyByte == nil {
		return nil, types.ErrGrantPubKeyNotFound()
	}
	grantPubKeys := new([]*GrantPermission)
	as.cdc.MustUnmarshalBinaryLengthPrefixed(grantPubKeyByte, grantPubKeys)
	return *grantPubKeys, nil
}

// GetAllGrantPermissions - returns grant user info keyed with pubkey.
func (as AccountStorage) GetAllGrantPermissions(ctx sdk.Context, me linotypes.AccountKey) ([]*GrantPermission, sdk.Error) {
	grantPermissions := make([]*GrantPermission, 0)
	store := ctx.KVStore(as.key)
	iter := sdk.KVStorePrefixIterator(store, getGrantPermPrefix(me))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		val := iter.Value()
		grantPermList := new([]*GrantPermission)
		as.cdc.MustUnmarshalBinaryLengthPrefixed(val, grantPermList)
		grantPermissions = append(grantPermissions, *grantPermList...)
	}
	return grantPermissions, nil
}

// SetGrantPermissions - sets a grant user to KV. Key is pubkey and value is grant user info
func (as AccountStorage) SetGrantPermissions(ctx sdk.Context, me linotypes.AccountKey, grantTo linotypes.AccountKey, grantPubKeys []*GrantPermission) {
	store := ctx.KVStore(as.key)
	grantPermByte := as.cdc.MustMarshalBinaryLengthPrefixed(grantPubKeys)
	store.Set(getGrantPermKey(me, grantTo), grantPermByte)
}

func (as AccountStorage) SetPool(ctx sdk.Context, pool *Pool) {
	store := ctx.KVStore(as.key)
	bz := as.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(GetAccountPoolKey(pool.Name), bz)
}

func (as AccountStorage) GetPool(ctx sdk.Context, name linotypes.PoolName) (*Pool, sdk.Error) {
	store := ctx.KVStore(as.key)
	bz := store.Get(GetAccountPoolKey(name))
	if bz == nil {
		return nil, types.ErrPoolNotFound(name)
	}
	pool := new(Pool)
	as.cdc.MustUnmarshalBinaryLengthPrefixed(bz, pool)
	return pool, nil
}

func (as AccountStorage) GetSupply(ctx sdk.Context) *Supply {
	store := ctx.KVStore(as.key)
	bz := store.Get(GetAccountSupplyKey())
	if bz == nil {
		panic("Lino Supply Not Initialized")
	}
	supply := new(Supply)
	as.cdc.MustUnmarshalBinaryLengthPrefixed(bz, supply)
	return supply
}

func (as AccountStorage) SetSupply(ctx sdk.Context, supply *Supply) {
	store := ctx.KVStore(as.key)
	bz := as.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(GetAccountSupplyKey(), bz)
}

func (as AccountStorage) PartialStoreMap(ctx sdk.Context) utils.StoreMap {
	store := ctx.KVStore(as.key)
	stores := []utils.SubStore{
		{
			Store:      store,
			Prefix:     AccountInfoSubstore,
			ValCreator: func() interface{} { return new(AccountInfo) },
			Decoder:    as.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
		{
			Store:      store,
			Prefix:     AccountBankSubstore,
			ValCreator: func() interface{} { return new(AccountBank) },
			Decoder:    as.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
		{
			Store:      store,
			Prefix:     AccountMetaSubstore,
			ValCreator: func() interface{} { return new(AccountMeta) },
			Decoder:    as.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
		{
			Store:      store,
			Prefix:     AccountGrantPubKeySubstore,
			ValCreator: func() interface{} { return new([]*GrantPermission) },
			Decoder:    as.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
		{
			Store:      store,
			Prefix:     AccountPoolSubstore,
			ValCreator: func() interface{} { return new(Pool) },
			Decoder:    as.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
	}
	return utils.NewStoreMap(stores)
}

// GetAccountInfoPrefix - "account info substore"
func GetAccountInfoPrefix() []byte {
	return AccountInfoSubstore
}

// GetAccountInfoKey - "account info substore" + "username"
func GetAccountInfoKey(accKey linotypes.AccountKey) []byte {
	return append(GetAccountInfoPrefix(), accKey...)
}

// GetAccountBankKey - "account bank substore" + "username"
func GetAccountBankKey(addr sdk.Address) []byte {
	return append(AccountBankSubstore, addr.Bytes()...)
}

// GetAccountMetaKey - "account meta substore" + "username"
func GetAccountMetaKey(accKey linotypes.AccountKey) []byte {
	return append(AccountMetaSubstore, accKey...)
}

// GetAccountPoolKey - "AccountPoolSubstore" + "pool name"
func GetAccountPoolKey(poolname linotypes.PoolName) []byte {
	return append(AccountPoolSubstore, poolname...)
}

// GetAccountSupplyKey - AccountSupplySubstore
func GetAccountSupplyKey() []byte {
	return AccountSupplySubstore
}

func getGrantPermPrefix(me linotypes.AccountKey) []byte {
	return append(append(AccountGrantPubKeySubstore, me...), linotypes.KeySeparator...)
}

func getGrantPermKey(me linotypes.AccountKey, grantTo linotypes.AccountKey) []byte {
	return append(append(getGrantPermPrefix(me), grantTo...), linotypes.KeySeparator...)
}

func ParseGrantKey(key []byte) (user linotypes.AccountKey, grantTo linotypes.AccountKey) {
	parsed := strings.Split(string(key), linotypes.KeySeparator)
	// see above for why 3, instead of 2.
	if len(parsed) != 3 {
		panic("illeage grankey:" + string(key))
	}
	return linotypes.AccountKey(parsed[0]), linotypes.AccountKey(parsed[1])
}
