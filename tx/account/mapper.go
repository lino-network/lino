package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var AccountInfoPrefix = []byte("AccountInfo/")
var AccountBankPrefix = []byte("AccountBank/")
var AccountMetaPrefix = []byte("AccountMeta/")
var AccountFollowersPrefix = []byte("Followers/")
var AccountFollowingsPrefix = []byte("Followering/")

var _ types.AccountManager = (*linoAccountManager)(nil)

// LinoAccountMapper implements types.AccountMapper
type linoAccountManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// Create and return a sealed account mapper
func NewLinoAccountManager(key sdk.StoreKey) linoAccountManager {
	cdc := wire.NewCodec()
	lam := linoAccountManager{
		key: key,
		cdc: cdc,
	}
	types.RegisterWireLinoAccount(cdc)
	return lam
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) GetInfo(ctx sdk.Context, accKey types.AccountKey) (*types.AccountInfo, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(accountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get info failed: info doesn't exist")
	}
	info := new(types.AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get info failed")
	}
	return info, nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *types.AccountInfo) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte, err := lam.cdc.MarshalBinary(*accInfo)
	if err != nil {
		return ErrCodeAccountManagerFail("linoAccountManager set info failed")
	}
	store.Set(accountInfoKey(accKey), infoByte)
	return nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) GetBankFromAccountKey(ctx sdk.Context, accKey types.AccountKey) (*types.AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(accountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get bank failed: user doesn't exist")
	}
	info := new(types.AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get bank failed: unmarshal failed")
	}
	return lam.GetBankFromAddress(ctx, info.Address)
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*types.AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	bankByte := store.Get(accountBankKey(address))
	if bankByte == nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get bank failed: bank doesn't exist")
	}
	bank := new(types.AccountBank)
	if err := lam.cdc.UnmarshalBinary(bankByte, bank); err != nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get bank failed: unmarshal failed")
	}
	return bank, nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) SetBank(ctx sdk.Context, address sdk.Address, accBank *types.AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	bankByte, err := lam.cdc.MarshalBinary(*accBank)
	if err != nil {
		return ErrCodeAccountManagerFail("linoAccountManager set bank failed")
	}
	store.Set(accountBankKey(address), bankByte)
	return nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) GetMeta(ctx sdk.Context, accKey types.AccountKey) (*types.AccountMeta, sdk.Error) {
	store := ctx.KVStore(lam.key)
	metaByte := store.Get(accountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get meta failed: meta doesn't exist")
	}
	meta := new(types.AccountMeta)
	if err := lam.cdc.UnmarshalBinary(metaByte, meta); err != nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get bank failed: unmarshal failed")
	}
	return meta, nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *types.AccountMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalBinary(*accMeta)
	if err != nil {
		return ErrCodeAccountManagerFail("linoAccountManager set meta failed")
	}
	store.Set(accountMetaKey(accKey), metaByte)
	return nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) GetFollowers(ctx sdk.Context, accKey types.AccountKey) (*types.Followers, sdk.Error) {
	store := ctx.KVStore(lam.key)
	followersByte := store.Get(accountFollowersKey(accKey))
	if followersByte == nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get followers failed: followers doesn't exist")
	}
	followers := new(types.Followers)
	if err := lam.cdc.UnmarshalBinary(followersByte, followers); err != nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get followers failed: unmarshal failed")
	}
	return followers, nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) SetFollowers(ctx sdk.Context, accKey types.AccountKey, followers *types.Followers) sdk.Error {
	store := ctx.KVStore(lam.key)
	followersByte, err := lam.cdc.MarshalBinary(*followers)
	if err != nil {
		return ErrCodeAccountManagerFail("linoAccountManager set meta failed")
	}
	store.Set(accountFollowersKey(accKey), followersByte)
	return nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) GetFollowings(ctx sdk.Context, accKey types.AccountKey) (*types.Followings, sdk.Error) {
	store := ctx.KVStore(lam.key)
	followingsByte := store.Get(accountFollowingsKey(accKey))
	if followingsByte == nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get followings failed: followers doesn't exist")
	}
	followings := new(types.Followings)
	if err := lam.cdc.UnmarshalBinary(followingsByte, followings); err != nil {
		return nil, ErrCodeAccountManagerFail("linoAccountManager get followings failed: unmarshal failed")
	}
	return followings, nil
}

// Implements sdk.AccountMapper.
func (lam linoAccountManager) SetFollowings(ctx sdk.Context, accKey types.AccountKey, followings *types.Followings) sdk.Error {
	store := ctx.KVStore(lam.key)
	followingsByte, err := lam.cdc.MarshalBinary(*followings)
	if err != nil {
		return ErrCodeAccountManagerFail("linoAccountManager set meta failed")
	}
	store.Set(accountFollowingsKey(accKey), followingsByte)
	return nil
}

func accountInfoKey(accKey types.AccountKey) []byte {
	return append(AccountInfoPrefix, accKey...)
}

func accountBankKey(address sdk.Address) []byte {
	return append(AccountBankPrefix, address...)
}

func accountMetaKey(accKey types.AccountKey) []byte {
	return append(AccountMetaPrefix, accKey...)
}

func accountFollowersKey(accKey types.AccountKey) []byte {
	return append(AccountFollowersPrefix, accKey...)
}

func accountFollowingsKey(accKey types.AccountKey) []byte {
	return append(AccountFollowingsPrefix, accKey...)
}
