package model

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
	AccountRewardSubstore    = []byte{0x05}
)

type AccountStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewLinoAccountStorage creates and returns a account manager
func NewAccountStorage(key sdk.StoreKey) *AccountStorage {
	cdc := wire.NewCodec()
	return &AccountStorage{
		key: key,
		cdc: cdc,
	}
}

func (as AccountStorage) AccountExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	if infoByte := store.Get(GetAccountInfoKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

func (as AccountStorage) GetInfo(ctx sdk.Context, accKey types.AccountKey) (*AccountInfo, sdk.Error) {
	store := ctx.KVStore(as.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountInfoDoesntExist()
	}
	info := new(AccountInfo)
	if err := as.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrGetAccountInfo().TraceCause(err, "")
	}
	return info, nil
}

func (as AccountStorage) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *AccountInfo) sdk.Error {
	store := ctx.KVStore(as.key)
	infoByte, err := as.cdc.MarshalBinary(*accInfo)
	if err != nil {
		return ErrSetInfoFailed()
	}
	store.Set(GetAccountInfoKey(accKey), infoByte)
	return nil
}

func (as AccountStorage) GetBankFromAccountKey(ctx sdk.Context, accKey types.AccountKey) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(as.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountBankDoesntExist()
	}
	info := new(AccountInfo)
	if err := as.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrGetBankFromAccountKey().TraceCause(err, "")
	}
	return as.GetBankFromAddress(ctx, info.Address)
}

func (as AccountStorage) GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(as.key)
	bankByte := store.Get(GetAccountBankKey(address))
	if bankByte == nil {
		return nil, ErrAccountBankDoesntExist()
	}
	bank := new(AccountBank)
	if err := as.cdc.UnmarshalBinary(bankByte, bank); err != nil {
		return nil, ErrGetBankFromAddress().TraceCause(err, "")
	}
	return bank, nil
}

func (as AccountStorage) SetBankFromAddress(ctx sdk.Context, address sdk.Address, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(as.key)
	bankByte, err := as.cdc.MarshalBinary(*accBank)
	if err != nil {
		return ErrSetBankFailed().TraceCause(err, "")
	}
	store.Set(GetAccountBankKey(address), bankByte)
	return nil
}

func (as AccountStorage) SetBankFromAccountKey(ctx sdk.Context, accKey types.AccountKey, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(as.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return ErrGetBankFromAccountKey()
	}
	info := new(AccountInfo)
	if err := as.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return ErrGetBankFromAccountKey().TraceCause(err, "")
	}

	return as.SetBankFromAddress(ctx, info.Address, accBank)
}

func (as AccountStorage) GetMeta(ctx sdk.Context, accKey types.AccountKey) (*AccountMeta, sdk.Error) {
	store := ctx.KVStore(as.key)
	metaByte := store.Get(GetAccountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrGetMetaFailed()
	}
	meta := new(AccountMeta)
	if err := as.cdc.UnmarshalBinary(metaByte, meta); err != nil {
		return nil, ErrGetMetaFailed().TraceCause(err, "")
	}
	return meta, nil
}

func (as AccountStorage) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *AccountMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalBinary(*accMeta)
	if err != nil {
		return ErrSetMetaFailed().TraceCause(err, "")
	}
	store.Set(GetAccountMetaKey(accKey), metaByte)
	return nil
}

func (as AccountStorage) IsMyFollower(ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	key := GetFollowerKey(me, follower)
	return store.Has(key)
}

func (as AccountStorage) SetFollowerMeta(ctx sdk.Context, me types.AccountKey, meta FollowerMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrSetFollowerMeta().TraceCause(err, "")
	}
	store.Set(GetFollowerKey(me, meta.FollowerName), metaByte)
	return nil
}

func (as AccountStorage) RemoveFollowerMeta(ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	store := ctx.KVStore(as.key)
	store.Delete(GetFollowerKey(me, follower))
	return nil
}

func (as AccountStorage) IsMyFollowing(ctx sdk.Context, me types.AccountKey, following types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	key := GetFollowingKey(me, following)
	return store.Has(key)
}

func (as AccountStorage) SetFollowingMeta(ctx sdk.Context, me types.AccountKey, meta FollowingMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrSetFollowingMeta().TraceCause(err, "")
	}
	store.Set(GetFollowingKey(me, meta.FollowingName), metaByte)
	return nil
}

func (as AccountStorage) RemoveFollowingMeta(ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	store := ctx.KVStore(as.key)
	store.Delete(GetFollowingKey(me, following))
	return nil
}

func (as AccountStorage) GetReward(ctx sdk.Context, accKey types.AccountKey) (*Reward, sdk.Error) {
	store := ctx.KVStore(as.key)
	rewardByte := store.Get(GetRewardKey(accKey))
	if rewardByte == nil {
		return nil, ErrGetRewardFailed()
	}
	reward := new(Reward)
	if err := as.cdc.UnmarshalBinary(rewardByte, reward); err != nil {
		return nil, ErrGetRewardFailed().TraceCause(err, "")
	}
	return reward, nil
}

func (as AccountStorage) SetReward(ctx sdk.Context, accKey types.AccountKey, reward *Reward) sdk.Error {
	store := ctx.KVStore(as.key)
	rewardByte, err := as.cdc.MarshalBinary(*reward)
	if err != nil {
		return ErrSetRewardFailed().TraceCause(err, "")
	}
	store.Set(GetRewardKey(accKey), rewardByte)
	return nil
}

func GetAccountInfoKey(accKey types.AccountKey) []byte {
	return append(AccountInfoSubstore, accKey...)
}

func GetAccountBankKey(address sdk.Address) []byte {
	return append(AccountBankSubstore, address...)
}

func GetAccountMetaKey(accKey types.AccountKey) []byte {
	return append(AccountMetaSubstore, accKey...)
}

func GetFollowerPrefix(me types.AccountKey) []byte {
	return append(append(AccountFollowerSubstore, me...), types.KeySeparator...)
}

func GetFollowingPrefix(me types.AccountKey) []byte {
	return append(append(AccountFollowingSubstore, me...), types.KeySeparator...)
}

// "follower substore" + "me" + "my follower"
func GetFollowerKey(me types.AccountKey, myFollower types.AccountKey) []byte {
	return append(GetFollowerPrefix(me), myFollower...)
}

// "following substore" + "me" + "my following"
func GetFollowingKey(me types.AccountKey, myFollowing types.AccountKey) []byte {
	return append(GetFollowingPrefix(me), myFollowing...)
}

func GetRewardKey(accKey types.AccountKey) []byte {
	return append(AccountRewardSubstore, accKey...)
}
