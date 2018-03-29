package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var (
	AccountInfoSubstore      = []byte{0x00}
	AccountBankSubstore      = []byte{0x01}
	AccountMetaSubstore      = []byte{0x02}
	AccountFollowerSubstore  = []byte{0x03}
	AccountFollowingSubstore = []byte{0x04}
)

// LinoAccountManager implements types.AccountManager
type AccountManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewLinoAccountManager creates and returns a account manager
func NewLinoAccountManager(key sdk.StoreKey) AccountManager {
	cdc := wire.NewCodec()
	lam := AccountManager{
		key: key,
		cdc: cdc,
	}
	RegisterWireLinoAccount(cdc)
	return lam
}

// Implements types.AccountManager.
func (lam AccountManager) AccountExist(ctx sdk.Context, accKey AccountKey) bool {
	store := ctx.KVStore(lam.key)
	if infoByte := store.Get(GetAccountInfoKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

// Implements types.AccountManager.
func (lam AccountManager) GetInfo(ctx sdk.Context, accKey AccountKey) (*AccountInfo, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrGetInfoFailed()
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrAccountUnmarshalError(err)
	}
	return info, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetInfo(ctx sdk.Context, accKey AccountKey, accInfo *AccountInfo) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte, err := lam.cdc.MarshalBinary(*accInfo)
	if err != nil {
		return ErrSetInfoFailed()
	}
	store.Set(GetAccountInfoKey(accKey), infoByte)
	return nil
}

// Implements types.AccountManager.
func (lam AccountManager) GetBankFromAccountKey(ctx sdk.Context, accKey AccountKey) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrGetBankFailed()
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrAccountUnmarshalError(err)
	}
	return lam.GetBankFromAddress(ctx, info.Address)
}

// Implements types.AccountManager.
func (lam AccountManager) GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	bankByte := store.Get(GetAccountBankKey(address))
	if bankByte == nil {
		return nil, ErrSetBankFailed()
	}
	bank := new(AccountBank)
	if err := lam.cdc.UnmarshalBinary(bankByte, bank); err != nil {
		return nil, ErrAccountUnmarshalError(err)
	}
	return bank, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetBankFromAddress(ctx sdk.Context, address sdk.Address, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	bankByte, err := lam.cdc.MarshalBinary(*accBank)
	if err != nil {
		return ErrSetBankFailed()
	}
	store.Set(GetAccountBankKey(address), bankByte)
	return nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetBankFromAccountKey(ctx sdk.Context, accKey AccountKey, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return ErrGetBankFailed()
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return ErrAccountUnmarshalError(err)
	}

	return lam.SetBankFromAddress(ctx, info.Address, accBank)
}

// Implements types.AccountManager.
func (lam AccountManager) GetMeta(ctx sdk.Context, accKey AccountKey) (*AccountMeta, sdk.Error) {
	store := ctx.KVStore(lam.key)
	metaByte := store.Get(GetAccountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrGetMetaFailed()
	}
	meta := new(AccountMeta)
	if err := lam.cdc.UnmarshalBinary(metaByte, meta); err != nil {
		return nil, ErrAccountUnmarshalError(err)
	}
	return meta, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetMeta(ctx sdk.Context, accKey AccountKey, accMeta *AccountMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalBinary(*accMeta)
	if err != nil {
		return ErrSetMetaFailed()
	}
	store.Set(GetAccountMetaKey(accKey), metaByte)
	return nil
}

func (lam AccountManager) IsMyFollower(ctx sdk.Context, me AccountKey, follower AccountKey) bool {
	store := ctx.KVStore(lam.key)
	key := GetFollowerKey(me, follower)
	return store.Has(key)
}

// Implements types.AccountManager.
func (lam AccountManager) SetFollowerMeta(ctx sdk.Context, me AccountKey, follower AccountKey, meta FollowerMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrAccountMarshalError(err)
	}
	store.Set(GetFollowerKey(me, follower), metaByte)
	return nil
}

func (lam AccountManager) RemoveFollowerMeta(ctx sdk.Context, me AccountKey, follower AccountKey) sdk.Error {
	store := ctx.KVStore(lam.key)
	store.Delete(GetFollowerKey(me, follower))
	return nil
}

func (lam AccountManager) IsMyFollowing(ctx sdk.Context, me AccountKey, followee AccountKey) bool {
	store := ctx.KVStore(lam.key)
	key := GetFollowingKey(me, followee)
	return store.Has(key)
}

// Implements types.AccountManager.
func (lam AccountManager) SetFollowingMeta(ctx sdk.Context, me AccountKey, followee AccountKey, meta FollowingMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrAccountMarshalError(err)
	}
	store.Set(GetFollowingKey(me, followee), metaByte)
	return nil
}

func (lam AccountManager) RemoveFollowingMeta(ctx sdk.Context, me AccountKey, followee AccountKey) sdk.Error {
	store := ctx.KVStore(lam.key)
	store.Delete(GetFollowingKey(me, followee))
	return nil
}

func GetAccountInfoKey(accKey AccountKey) []byte {
	return append(AccountInfoSubstore, accKey...)
}

func GetAccountBankKey(address sdk.Address) []byte {
	return append(AccountBankSubstore, address...)
}

func GetAccountMetaKey(accKey AccountKey) []byte {
	return append(AccountMetaSubstore, accKey...)
}

func GetFollowerPrefix(me AccountKey) []byte {
	return append(append(AccountFollowerSubstore, me...), types.KeySeparator...)
}

func GetFollowingPrefix(me AccountKey) []byte {
	return append(append(AccountFollowingSubstore, me...), types.KeySeparator...)
}

// "follower substore" + "me" + "my follower"
func GetFollowerKey(me AccountKey, myFollower AccountKey) []byte {
	return append(GetFollowerPrefix(me), myFollower...)
}

// "following substore" + "me" + "my following"
func GetFollowingKey(me AccountKey, myFollowing AccountKey) []byte {
	return append(GetFollowingPrefix(me), myFollowing...)
}
