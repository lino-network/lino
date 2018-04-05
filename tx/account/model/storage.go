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
	lam := AccountStorage{
		key: key,
		cdc: cdc,
	}
	return &lam
}

func (lam AccountStorage) AccountExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(lam.key)
	if infoByte := store.Get(GetAccountInfoKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

func (lam AccountStorage) GetInfo(ctx sdk.Context, accKey types.AccountKey) (*AccountInfo, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountInfoDoesntExist()
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrGetAccountInfo().TraceCause(err, "")
	}
	return info, nil
}

func (lam AccountStorage) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *AccountInfo) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte, err := lam.cdc.MarshalBinary(*accInfo)
	if err != nil {
		return ErrSetInfoFailed()
	}
	store.Set(GetAccountInfoKey(accKey), infoByte)
	return nil
}

func (lam AccountStorage) GetBankFromAccountKey(ctx sdk.Context, accKey types.AccountKey) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountBankDoesntExist()
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrGetBankFromAccountKey().TraceCause(err, "")
	}
	return lam.GetBankFromAddress(ctx, info.Address)
}

func (lam AccountStorage) GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	bankByte := store.Get(GetAccountBankKey(address))
	if bankByte == nil {
		return nil, ErrAccountBankDoesntExist()
	}
	bank := new(AccountBank)
	if err := lam.cdc.UnmarshalBinary(bankByte, bank); err != nil {
		return nil, ErrGetBankFromAddress().TraceCause(err, "")
	}
	return bank, nil
}

func (lam AccountStorage) SetBankFromAddress(ctx sdk.Context, address sdk.Address, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	bankByte, err := lam.cdc.MarshalBinary(*accBank)
	if err != nil {
		return ErrSetBankFailed().TraceCause(err, "")
	}
	store.Set(GetAccountBankKey(address), bankByte)
	return nil
}

func (lam AccountStorage) SetBankFromAccountKey(ctx sdk.Context, accKey types.AccountKey, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return ErrGetBankFromAccountKey()
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return ErrGetBankFromAccountKey().TraceCause(err, "")
	}

	return lam.SetBankFromAddress(ctx, info.Address, accBank)
}

func (lam AccountStorage) GetMeta(ctx sdk.Context, accKey types.AccountKey) (*AccountMeta, sdk.Error) {
	store := ctx.KVStore(lam.key)
	metaByte := store.Get(GetAccountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrGetMetaFailed()
	}
	meta := new(AccountMeta)
	if err := lam.cdc.UnmarshalBinary(metaByte, meta); err != nil {
		return nil, ErrGetMetaFailed().TraceCause(err, "")
	}
	return meta, nil
}

func (lam AccountStorage) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *AccountMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalBinary(*accMeta)
	if err != nil {
		return ErrSetMetaFailed().TraceCause(err, "")
	}
	store.Set(GetAccountMetaKey(accKey), metaByte)
	return nil
}

func (lam AccountStorage) IsMyFollower(ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	store := ctx.KVStore(lam.key)
	key := GetFollowerKey(me, follower)
	return store.Has(key)
}

func (lam AccountStorage) SetFollowerMeta(ctx sdk.Context, me types.AccountKey, meta FollowerMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrSetFollowerMeta().TraceCause(err, "")
	}
	store.Set(GetFollowerKey(me, meta.FollowerName), metaByte)
	return nil
}

func (lam AccountStorage) RemoveFollowerMeta(ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	store := ctx.KVStore(lam.key)
	store.Delete(GetFollowerKey(me, follower))
	return nil
}

func (lam AccountStorage) IsMyFollowee(ctx sdk.Context, me types.AccountKey, followee types.AccountKey) bool {
	store := ctx.KVStore(lam.key)
	key := GetFolloweeKey(me, followee)
	return store.Has(key)
}

func (lam AccountStorage) SetFolloweeMeta(ctx sdk.Context, me types.AccountKey, meta FollowingMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrSetFollowingMeta().TraceCause(err, "")
	}
	store.Set(GetFolloweeKey(me, meta.FolloweeName), metaByte)
	return nil
}

func (lam AccountStorage) RemoveFolloweeMeta(ctx sdk.Context, me types.AccountKey, followee types.AccountKey) sdk.Error {
	store := ctx.KVStore(lam.key)
	store.Delete(GetFolloweeKey(me, followee))
	return nil
}

func (lam AccountStorage) GetReward(ctx sdk.Context, accKey types.AccountKey) (*Reward, sdk.Error) {
	store := ctx.KVStore(lam.key)
	rewardByte := store.Get(GetRewardKey(accKey))
	if rewardByte == nil {
		return nil, ErrGetRewardFailed()
	}
	reward := new(Reward)
	if err := lam.cdc.UnmarshalBinary(rewardByte, reward); err != nil {
		return nil, ErrGetRewardFailed().TraceCause(err, "")
	}
	return reward, nil
}

func (lam AccountStorage) SetReward(ctx sdk.Context, accKey types.AccountKey, reward *Reward) sdk.Error {
	store := ctx.KVStore(lam.key)
	rewardByte, err := lam.cdc.MarshalBinary(*reward)
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
func GetFolloweeKey(me types.AccountKey, myFollowee types.AccountKey) []byte {
	return append(GetFollowingPrefix(me), myFollowee...)
}

func GetRewardKey(accKey types.AccountKey) []byte {
	return append(AccountRewardSubstore, accKey...)
}
