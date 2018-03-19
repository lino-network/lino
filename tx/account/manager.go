package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

var AccountInfoPrefix = []byte("AccountInfo/")
var AccountBankPrefix = []byte("AccountBank/")
var AccountMetaPrefix = []byte("AccountMeta/")
var AccountFollowerPrefix = []byte("Follower/")
var AccountFollowingPrefix = []byte("Followering/")

var _ types.AccountManager = (*LinoAccountManager)(nil)

// LinoAccountManager implements types.AccountManager
type LinoAccountManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewLinoAccountManager creates and returns a account manager
func NewLinoAccountManager(key sdk.StoreKey) LinoAccountManager {
	cdc := wire.NewCodec()
	lam := LinoAccountManager{
		key: key,
		cdc: cdc,
	}
	types.RegisterWireLinoAccount(cdc)
	return lam
}

// Implements types.AccountManager.
func (lam LinoAccountManager) CreateAccount(ctx sdk.Context, accKey types.AccountKey, pubkey crypto.PubKey, accBank *types.AccountBank) (*types.AccountInfo, sdk.Error) {
	accInfo := types.AccountInfo{
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
	if err := lam.SetBank(ctx, accInfo.Address, accBank); err != nil {
		return nil, err
	}

	accMeta := types.AccountMeta{
		LastActivity:   types.Height(ctx.BlockHeight()),
		ActivityBurden: types.DefaultActivityBurden,
		LastABBlock:    types.Height(ctx.BlockHeight()),
	}
	if err := lam.SetMeta(ctx, accInfo.Username, &accMeta); err != nil {
		return nil, err
	}

	follower := types.Follower{Follower: []types.AccountKey{}}
	if err := lam.SetFollower(ctx, accInfo.Username, &follower); err != nil {
		return nil, err
	}
	following := types.Following{Following: []types.AccountKey{}}
	if err := lam.SetFollowing(ctx, accInfo.Username, &following); err != nil {
		return nil, err
	}
	return &accInfo, nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) AccountExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(lam.key)
	if infoByte := store.Get(accountInfoKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

// Implements types.AccountManager.
func (lam LinoAccountManager) GetInfo(ctx sdk.Context, accKey types.AccountKey) (*types.AccountInfo, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(accountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get info failed: info doesn't exist")
	}
	info := new(types.AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get info failed")
	}
	return info, nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *types.AccountInfo) sdk.Error {
	store := ctx.KVStore(lam.key)
	infoByte, err := lam.cdc.MarshalBinary(*accInfo)
	if err != nil {
		return ErrAccountManagerFail("LinoAccountManager set info failed")
	}
	store.Set(accountInfoKey(accKey), infoByte)
	return nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) GetBankFromAccountKey(ctx sdk.Context, accKey types.AccountKey) (*types.AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	infoByte := store.Get(accountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get bank failed: user doesn't exist")
	}
	info := new(types.AccountInfo)
	if err := lam.cdc.UnmarshalBinary(infoByte, info); err != nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get bank failed: unmarshal failed")
	}
	return lam.GetBankFromAddress(ctx, info.Address)
}

// Implements types.AccountManager.
func (lam LinoAccountManager) GetBankFromAddress(ctx sdk.Context, address sdk.Address) (*types.AccountBank, sdk.Error) {
	store := ctx.KVStore(lam.key)
	bankByte := store.Get(accountBankKey(address))
	if bankByte == nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get bank failed: bank doesn't exist")
	}
	bank := new(types.AccountBank)
	if err := lam.cdc.UnmarshalBinary(bankByte, bank); err != nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get bank failed: unmarshal failed")
	}
	return bank, nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) SetBank(ctx sdk.Context, address sdk.Address, accBank *types.AccountBank) sdk.Error {
	store := ctx.KVStore(lam.key)
	bankByte, err := lam.cdc.MarshalBinary(*accBank)
	if err != nil {
		return ErrAccountManagerFail("LinoAccountManager set bank failed")
	}
	store.Set(accountBankKey(address), bankByte)
	return nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) GetMeta(ctx sdk.Context, accKey types.AccountKey) (*types.AccountMeta, sdk.Error) {
	store := ctx.KVStore(lam.key)
	metaByte := store.Get(accountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get meta failed: meta doesn't exist")
	}
	meta := new(types.AccountMeta)
	if err := lam.cdc.UnmarshalBinary(metaByte, meta); err != nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get bank failed: unmarshal failed")
	}
	return meta, nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *types.AccountMeta) sdk.Error {
	store := ctx.KVStore(lam.key)
	metaByte, err := lam.cdc.MarshalBinary(*accMeta)
	if err != nil {
		return ErrAccountManagerFail("LinoAccountManager set meta failed")
	}
	store.Set(accountMetaKey(accKey), metaByte)
	return nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) GetFollower(ctx sdk.Context, accKey types.AccountKey) (*types.Follower, sdk.Error) {
	store := ctx.KVStore(lam.key)
	followerByte := store.Get(accountFollowerKey(accKey))
	if followerByte == nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get follower failed: follower doesn't exist")
	}
	follower := new(types.Follower)
	if err := lam.cdc.UnmarshalBinary(followerByte, follower); err != nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get follower failed: unmarshal failed")
	}
	return follower, nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) SetFollower(ctx sdk.Context, accKey types.AccountKey, follower *types.Follower) sdk.Error {
	store := ctx.KVStore(lam.key)
	followerByte, err := lam.cdc.MarshalBinary(*follower)
	if err != nil {
		return ErrAccountManagerFail("LinoAccountManager set meta failed")
	}
	store.Set(accountFollowerKey(accKey), followerByte)
	return nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) GetFollowing(ctx sdk.Context, accKey types.AccountKey) (*types.Following, sdk.Error) {
	store := ctx.KVStore(lam.key)
	followingByte := store.Get(accountFollowingKey(accKey))
	if followingByte == nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get following failed: follower doesn't exist")
	}
	following := new(types.Following)
	if err := lam.cdc.UnmarshalBinary(followingByte, following); err != nil {
		return nil, ErrAccountManagerFail("LinoAccountManager get following failed: unmarshal failed")
	}
	return following, nil
}

// Implements types.AccountManager.
func (lam LinoAccountManager) SetFollowing(ctx sdk.Context, accKey types.AccountKey, following *types.Following) sdk.Error {
	store := ctx.KVStore(lam.key)
	followingByte, err := lam.cdc.MarshalBinary(*following)
	if err != nil {
		return ErrAccountManagerFail("LinoAccountManager set meta failed")
	}
	store.Set(accountFollowingKey(accKey), followingByte)
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

func accountFollowerKey(accKey types.AccountKey) []byte {
	return append(AccountFollowerPrefix, accKey...)
}

func accountFollowingKey(accKey types.AccountKey) []byte {
	return append(AccountFollowingPrefix, accKey...)
}
