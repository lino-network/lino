package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/go-crypto"
)

var AccountInfoPrefix = []byte("AccountInfo/")
var AccountBankPrefix = []byte("AccountBank/")
var AccountMetaPrefix = []byte("AccountMeta/")
var AccountFollowerPrefix = []byte("Follower/")
var AccountFollowingPrefix = []byte("Followering/")

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
	if infoByte := store.Get(accountInfoKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

// Implements types.AccountManager.
func (lam AccountManager) GetInfo(ctx sdk.Context, accKey AccountKey) (*AccountInfo, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(accountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountManagerFail("AccountManager get info failed: info doesn't exist")
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrAccountManagerFail("AccountManager get info failed")
	}
	return info, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetInfo(ctx sdk.Context, accKey AccountKey, accInfo *AccountInfo) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte, err := lam.cdc.MarshalBinary(*accInfo)
	if err != nil {
		return ErrAccountManagerFail("LinoAccountManager set info failed")
	}
	store.Set(accountInfoKey(accKey), infoByte)
	return nil
}

// Implements types.AccountManager.
func (lam AccountManager) GetBankFromAccountKey(ctx sdk.Context, accKey AccountKey) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(accountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountManagerFail("AccountManager get bank failed: user doesn't exist")
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrAccountManagerFail("AccountManager get bank failed: unmarshal failed")
	}
	return lam.GetBankFromAddress(ctx, info.Address)
}

// Implements types.AccountManager.
func (lam AccountManager) GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	bankByte := store.Get(accountBankKey(address))
	if bankByte == nil {
		return nil, ErrAccountManagerFail("AccountManager get bank failed: bank doesn't exist")
	}
	bank := new(AccountBank)
	if err := lam.cdc.UnmarshalBinary(bankByte, bank); err != nil {
		return nil, ErrAccountManagerFail("AccountManager get bank failed: unmarshal failed")
	}
	return bank, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetBankFromAddress(ctx sdk.Context, address sdk.Address, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	bankByte, err := lam.cdc.MarshalBinary(*accBank)
	if err != nil {
		return ErrAccountManagerFail("AccountManager set bank failed")
	}
	store.Set(accountBankKey(address), bankByte)
	return nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetBankFromAccountKey(ctx sdk.Context, accKey AccountKey, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(accountInfoKey(accKey))
	if infoByte == nil {
		return ErrAccountManagerFail("AccountManager set bank failed: user doesn't exist")
	}
	info := new(AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return ErrAccountManagerFail("AccountManager set bank failed: unmarshal failed")
	}

	return lam.SetBankFromAddress(ctx, info.Address, accBank)
}

// Implements types.AccountManager.
func (lam AccountManager) GetMeta(ctx sdk.Context, accKey AccountKey) (*AccountMeta, sdk.Error) {
	store := ctx.KVStore(lam.key)
	metaByte := store.Get(accountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrAccountManagerFail("AccountManager get meta failed: meta doesn't exist")
	}
	meta := new(AccountMeta)
	if err := lam.cdc.UnmarshalBinary(metaByte, meta); err != nil {
		return nil, ErrAccountManagerFail("AccountManager get bank failed: unmarshal failed")
	}
	return meta, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetMeta(ctx sdk.Context, accKey AccountKey, accMeta *AccountMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalBinary(*accMeta)
	if err != nil {
		return ErrAccountManagerFail("AccountManager set meta failed")
	}
	store.Set(accountMetaKey(accKey), metaByte)
	return nil
}

// Implements types.AccountManager.
func (lam AccountManager) GetFollower(ctx sdk.Context, accKey AccountKey) (*Follower, sdk.Error) {
	store := ctx.KVStore(lam.key)
	followerByte := store.Get(accountFollowerKey(accKey))
	if followerByte == nil {
		return nil, ErrAccountManagerFail("AccountManager get follower failed: follower doesn't exist")
	}
	follower := new(Follower)
	if err := lam.cdc.UnmarshalBinary(followerByte, follower); err != nil {
		return nil, ErrAccountManagerFail("AccountManager get follower failed: unmarshal failed")
	}
	return follower, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetFollower(ctx sdk.Context, accKey AccountKey, follower *Follower) sdk.Error {
	store := ctx.KVStore(lam.key)
	followerByte, err := lam.cdc.MarshalBinary(*follower)
	if err != nil {
		return ErrAccountManagerFail("AccountManager set meta failed")
	}
	store.Set(accountFollowerKey(accKey), followerByte)
	return nil
}

// Implements types.AccountManager.
func (lam AccountManager) GetFollowing(ctx sdk.Context, accKey AccountKey) (*Following, sdk.Error) {
	store := ctx.KVStore(lam.key)
	followingByte := store.Get(accountFollowingKey(accKey))
	if followingByte == nil {
		return nil, ErrAccountManagerFail("AccountManager get following failed: follower doesn't exist")
	}
	following := new(Following)
	if err := lam.cdc.UnmarshalBinary(followingByte, following); err != nil {
		return nil, ErrAccountManagerFail("AccountManager get following failed: unmarshal failed")
	}
	return following, nil
}

// Implements types.AccountManager.
func (lam AccountManager) SetFollowing(ctx sdk.Context, accKey AccountKey, following *Following) sdk.Error {
	store := ctx.KVStore(lam.key)
	followingByte, err := lam.cdc.MarshalBinary(*following)
	if err != nil {
		return ErrAccountManagerFail("AccountManager set meta failed")
	}
	store.Set(accountFollowingKey(accKey), followingByte)
	return nil
}

func accountInfoKey(accKey AccountKey) []byte {
	return append(AccountInfoPrefix, accKey...)
}

func accountBankKey(address sdk.Address) []byte {
	return append(AccountBankPrefix, address...)
}

func accountMetaKey(accKey AccountKey) []byte {
	return append(AccountMetaPrefix, accKey...)
}

func accountFollowerKey(accKey AccountKey) []byte {
	return append(AccountFollowerPrefix, accKey...)
}

func accountFollowingKey(accKey AccountKey) []byte {
	return append(AccountFollowingPrefix, accKey...)
}

// Implements types.AccountManager.
func (lam AccountManager) CreateAccount(ctx sdk.Context, accKey AccountKey, pubkey crypto.PubKey, accBank *AccountBank) (*AccountInfo, sdk.Error) {
	accInfo := AccountInfo{
		Username: accKey,
		Created:  types.Height(ctx.BlockHeight()),
		PostKey:  pubkey,
		OwnerKey: pubkey,
		Address:  pubkey.Address(),
	}
	if err := lam.SetInfo(ctx, accInfo.Username, &accInfo); err != nil {
		return nil, err
	}

	accBank.Username = accKey
	if err := lam.SetBankFromAddress(ctx, accInfo.Address, accBank); err != nil {
		return nil, err
	}

	accMeta := AccountMeta{
		LastActivity:   types.Height(ctx.BlockHeight()),
		ActivityBurden: types.DefaultActivityBurden,
	}
	if err := lam.SetMeta(ctx, accInfo.Username, &accMeta); err != nil {
		return nil, err
	}

	follower := Follower{Follower: []AccountKey{}}
	if err := lam.SetFollower(ctx, accInfo.Username, &follower); err != nil {
		return nil, err
	}
	following := Following{Following: []AccountKey{}}
	if err := lam.SetFollowing(ctx, accInfo.Username, &following); err != nil {
		return nil, err
	}
	return &accInfo, nil
}
