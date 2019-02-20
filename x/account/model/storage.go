package model

import (
	"encoding/hex"
	"strconv"

	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/codec"
)

var (
	accountInfoSubstore                = []byte{0x00}
	accountBankSubstore                = []byte{0x01}
	accountMetaSubstore                = []byte{0x02}
	accountFollowerSubstore            = []byte{0x03}
	accountFollowingSubstore           = []byte{0x04}
	accountRewardSubstore              = []byte{0x05}
	accountPendingCoinDayQueueSubstore = []byte{0x06}
	accountRelationshipSubstore        = []byte{0x07}
	accountBalanceHistorySubstore      = []byte{0x08}
	accountGrantPubKeySubstore         = []byte{0x09}
	accountRewardHistorySubstore       = []byte{0x0a}
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

	return AccountStorage{
		key: key,
		cdc: cdc,
	}
}

// DoesAccountExist - returns true when a specific account exist in the KVStore.
func (as AccountStorage) DoesAccountExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	return store.Has(GetAccountInfoKey(accKey))
}

// GetInfo - returns general account info of a specific account, returns error otherwise.
func (as AccountStorage) GetInfo(ctx sdk.Context, accKey types.AccountKey) (*AccountInfo, sdk.Error) {
	store := ctx.KVStore(as.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountInfoNotFound()
	}
	info := new(AccountInfo)
	if err := as.cdc.UnmarshalJSON(infoByte, info); err != nil {
		return nil, ErrFailedToUnmarshalAccountInfo(err)
	}
	return info, nil
}

// SetInfo - sets general account info to a specific account, returns error if any.
func (as AccountStorage) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *AccountInfo) sdk.Error {
	store := ctx.KVStore(as.key)
	infoByte, err := as.cdc.MarshalJSON(*accInfo)
	if err != nil {
		return ErrFailedToMarshalAccountInfo(err)
	}
	store.Set(GetAccountInfoKey(accKey), infoByte)
	return nil
}

// GetBankFromAccountKey - returns bank info of a specific account, returns error
// if any.
func (as AccountStorage) GetBankFromAccountKey(
	ctx sdk.Context, me types.AccountKey) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(as.key)
	bankByte := store.Get(GetAccountBankKey(me))
	if bankByte == nil {
		return nil, ErrAccountBankNotFound()
	}
	bank := new(AccountBank)
	if err := as.cdc.UnmarshalJSON(bankByte, bank); err != nil {
		return nil, ErrFailedToUnmarshalAccountBank(err)
	}
	return bank, nil
}

// SetBankFromAddress - sets bank info for a given address,
// returns error if any.
func (as AccountStorage) SetBankFromAccountKey(ctx sdk.Context, username types.AccountKey, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(as.key)
	bankByte, err := as.cdc.MarshalJSON(*accBank)
	if err != nil {
		return ErrFailedToMarshalAccountBank(err)
	}
	store.Set(GetAccountBankKey(username), bankByte)
	return nil
}

// GetMeta - returns meta of a given account that are tiny and frequently updated fields.
func (as AccountStorage) GetMeta(ctx sdk.Context, accKey types.AccountKey) (*AccountMeta, sdk.Error) {
	store := ctx.KVStore(as.key)
	metaByte := store.Get(GetAccountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrAccountMetaNotFound()
	}
	meta := new(AccountMeta)
	if err := as.cdc.UnmarshalJSON(metaByte, meta); err != nil {
		return nil, ErrFailedToUnmarshalAccountMeta(err)
	}
	return meta, nil
}

// SetMeta - sets meta for a given account, returns error if any.
func (as AccountStorage) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *AccountMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(*accMeta)
	if err != nil {
		return ErrFailedToMarshalAccountMeta(err)
	}
	store.Set(GetAccountMetaKey(accKey), metaByte)
	return nil
}

// IsMyfollower - returns true if `follower` follows `me`.
func (as AccountStorage) IsMyFollower(ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	key := getFollowerKey(me, follower)
	return store.Has(key)
}

// SetFollowerMeta - sets follower meta info for a given account which includes time and follower name.
func (as AccountStorage) SetFollowerMeta(ctx sdk.Context, me types.AccountKey, meta FollowerMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrFailedToMarshalFollowerMeta(err)
	}
	store.Set(getFollowerKey(me, meta.FollowerName), metaByte)
	return nil
}

// RemoveFollowerMeta removes follower meta info of a relationship.
func (as AccountStorage) RemoveFollowerMeta(ctx sdk.Context, me types.AccountKey, follower types.AccountKey) {
	store := ctx.KVStore(as.key)
	store.Delete(getFollowerKey(me, follower))
	return
}

// IsMyFollowing - returns true if `me` follows `following`
func (as AccountStorage) IsMyFollowing(ctx sdk.Context, me types.AccountKey, following types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	key := getFollowingKey(me, following)
	return store.Has(key)
}

// SetFollowerMeta - sets following meta info for a given account which includes
// time and following name.
func (as AccountStorage) SetFollowingMeta(ctx sdk.Context, me types.AccountKey, meta FollowingMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrFailedToMarshalFollowingMeta(err)
	}
	store.Set(getFollowingKey(me, meta.FollowingName), metaByte)
	return nil
}

// RemoveFollowingMeta - removes following meta info of a relationship.
func (as AccountStorage) RemoveFollowingMeta(ctx sdk.Context, me types.AccountKey, following types.AccountKey) {
	store := ctx.KVStore(as.key)
	store.Delete(getFollowingKey(me, following))
	return
}

// GetReward - returns reward info of a given account, returns error if any.
func (as AccountStorage) GetReward(ctx sdk.Context, accKey types.AccountKey) (*Reward, sdk.Error) {
	store := ctx.KVStore(as.key)
	rewardByte := store.Get(getRewardKey(accKey))
	if rewardByte == nil {
		return nil, ErrRewardNotFound()
	}
	reward := new(Reward)
	if err := as.cdc.UnmarshalJSON(rewardByte, reward); err != nil {
		return nil, ErrFailedToUnmarshalReward(err)
	}
	return reward, nil
}

// SetReward - sets the rewards info of a given account, returns error if any.
func (as AccountStorage) SetReward(ctx sdk.Context, accKey types.AccountKey, reward *Reward) sdk.Error {
	store := ctx.KVStore(as.key)
	rewardByte, err := as.cdc.MarshalJSON(*reward)
	if err != nil {
		return ErrFailedToMarshalReward(err)
	}
	store.Set(getRewardKey(accKey), rewardByte)
	return nil
}

// GetPendingCoinDayQueue - returns a pending coin day queue for a given address.
func (as AccountStorage) GetPendingCoinDayQueue(
	ctx sdk.Context, me types.AccountKey) (*PendingCoinDayQueue, sdk.Error) {
	store := ctx.KVStore(as.key)
	pendingCoinDayQueueByte := store.Get(getPendingCoinDayQueueKey(me))
	if pendingCoinDayQueueByte == nil {
		return nil, ErrPendingCoinDayQueueNotFound()
	}
	queue := new(PendingCoinDayQueue)
	if err := as.cdc.UnmarshalJSON(pendingCoinDayQueueByte, queue); err != nil {
		return nil, ErrFailedToUnmarshalPendingCoinDayQueue(err)
	}
	return queue, nil
}

// SetPendingCoinDayQueue - sets a pending coin day queue for a given username.
func (as AccountStorage) SetPendingCoinDayQueue(ctx sdk.Context, me types.AccountKey, pendingCoinDayQueue *PendingCoinDayQueue) sdk.Error {
	store := ctx.KVStore(as.key)
	pendingCoinDayQueueByte, err := as.cdc.MarshalJSON(*pendingCoinDayQueue)
	if err != nil {
		return ErrFailedToMarshalPendingCoinDayQueue(err)
	}
	store.Set(getPendingCoinDayQueueKey(me), pendingCoinDayQueueByte)
	return nil
}

// DeleteGrantPubKey - deletes given pubkey in KV.
func (as AccountStorage) DeleteGrantPubKey(ctx sdk.Context, me types.AccountKey, pubKey crypto.PubKey) {
	store := ctx.KVStore(as.key)
	store.Delete(getGrantPubKeyKey(me, pubKey))
	return
}

// GetGrantPubKey - returns grant user info keyed with pubkey.
func (as AccountStorage) GetGrantPubKey(ctx sdk.Context, me types.AccountKey, pubKey crypto.PubKey) (*GrantPubKey, sdk.Error) {
	store := ctx.KVStore(as.key)
	grantPubKeyByte := store.Get(getGrantPubKeyKey(me, pubKey))
	if grantPubKeyByte == nil {
		return nil, ErrGrantPubKeyNotFound()
	}
	grantPubKey := new(GrantPubKey)
	if err := as.cdc.UnmarshalJSON(grantPubKeyByte, grantPubKey); err != nil {
		return nil, ErrFailedToUnmarshalGrantPubKey(err)
	}
	return grantPubKey, nil
}

// SetGrantPubKey - sets a grant user to KV. Key is pubkey and value is grant user info
func (as AccountStorage) SetGrantPubKey(ctx sdk.Context, me types.AccountKey, pubKey crypto.PubKey, grantPubKey *GrantPubKey) sdk.Error {
	store := ctx.KVStore(as.key)
	grantPubKeyByte, err := as.cdc.MarshalJSON(*grantPubKey)
	if err != nil {
		return ErrFailedToMarshalGrantPubKey(err)
	}
	store.Set(getGrantPubKeyKey(me, pubKey), grantPubKeyByte)
	return nil
}

// GetRelationship - returns the relationship between two accounts
func (as AccountStorage) GetRelationship(ctx sdk.Context, me types.AccountKey, other types.AccountKey) (*Relationship, sdk.Error) {
	store := ctx.KVStore(as.key)
	relationshipByte := store.Get(getRelationshipKey(me, other))
	if relationshipByte == nil {
		return nil, nil
	}
	queue := new(Relationship)
	if err := as.cdc.UnmarshalJSON(relationshipByte, queue); err != nil {
		return nil, ErrFailedToUnmarshalRelationship(err)
	}
	return queue, nil
}

// SetRelationship - sets relationship for two accounts
func (as AccountStorage) SetRelationship(ctx sdk.Context, me types.AccountKey, other types.AccountKey, relationship *Relationship) sdk.Error {
	store := ctx.KVStore(as.key)
	relationshipByte, err := as.cdc.MarshalJSON(*relationship)
	if err != nil {
		return ErrFailedToMarshalRelationship(err)
	}
	store.Set(getRelationshipKey(me, other), relationshipByte)
	return nil
}

// GetRelationship - returns the relationship between two accounts
func (as AccountStorage) GetBalanceHistory(
	ctx sdk.Context, me types.AccountKey, bucketSlot int64) (*BalanceHistory, sdk.Error) {
	store := ctx.KVStore(as.key)
	balanceHistoryBytes := store.Get(getBalanceHistoryKey(me, bucketSlot))
	if balanceHistoryBytes == nil {
		return nil, nil
	}
	history := new(BalanceHistory)
	if err := as.cdc.UnmarshalJSON(balanceHistoryBytes, history); err != nil {
		return nil, ErrFailedToUnmarshalBalanceHistory(err)
	}
	return history, nil
}

// SetBalanceHistory - sets balance history
func (as AccountStorage) SetBalanceHistory(
	ctx sdk.Context, me types.AccountKey, bucketSlot int64, history *BalanceHistory) sdk.Error {
	store := ctx.KVStore(as.key)
	historyBytes, err := as.cdc.MarshalJSON(*history)
	if err != nil {
		return ErrFailedToMarshalBalanceHistory(err)
	}
	store.Set(getBalanceHistoryKey(me, bucketSlot), historyBytes)
	return nil
}

// GetRewardHistory - returns the history of rewards that a user received
func (as AccountStorage) GetRewardHistory(
	ctx sdk.Context, me types.AccountKey, bucketSlot int64) (*RewardHistory, sdk.Error) {
	store := ctx.KVStore(as.key)
	rewardHistoryBytes := store.Get(getRewardHistoryKey(me, bucketSlot))
	if rewardHistoryBytes == nil {
		return nil, nil
	}
	history := new(RewardHistory)
	if err := as.cdc.UnmarshalJSON(rewardHistoryBytes, history); err != nil {
		return nil, ErrFailedToUnmarshalRewardHistory(err)
	}
	return history, nil
}

// SetRewardHistory - sets reward history
func (as AccountStorage) SetRewardHistory(
	ctx sdk.Context, me types.AccountKey, bucketSlot int64, history *RewardHistory) sdk.Error {
	store := ctx.KVStore(as.key)
	historyBytes, err := as.cdc.MarshalJSON(*history)
	if err != nil {
		return ErrFailedToMarshalRewardHistory(err)
	}
	store.Set(getRewardHistoryKey(me, bucketSlot), historyBytes)
	return nil
}

// DeleteRewardHistory - delete reward history from KVStore
func (as AccountStorage) DeleteRewardHistory(ctx sdk.Context, me types.AccountKey, bucketSlot int64) {
	store := ctx.KVStore(as.key)
	store.Delete(getRewardHistoryKey(me, bucketSlot))
	return
}

// GetAccountInfoPrefix - "account info substore"
func GetAccountInfoPrefix() []byte {
	return accountInfoSubstore
}

// GetAccountInfoKey - "account info substore" + "username"
func GetAccountInfoKey(accKey types.AccountKey) []byte {
	return append(GetAccountInfoPrefix(), accKey...)
}

// GetAccountBankKey - "account bank substore" + "username"
func GetAccountBankKey(accKey types.AccountKey) []byte {
	return append(accountBankSubstore, accKey...)
}

// GetAccountMetaKey - "account meta substore" + "username"
func GetAccountMetaKey(accKey types.AccountKey) []byte {
	return append(accountMetaSubstore, accKey...)
}

func getFollowerKey(me types.AccountKey, myFollower types.AccountKey) []byte {
	return append(getFollowerPrefix(me), myFollower...)
}

func getFollowerPrefix(me types.AccountKey) []byte {
	return append(append(accountFollowerSubstore, me...), types.KeySeparator...)
}

// "following substore" + "me" + "my following"
func getFollowingKey(me types.AccountKey, myFollowing types.AccountKey) []byte {
	return append(getFollowingPrefix(me), myFollowing...)
}

func getFollowingPrefix(me types.AccountKey) []byte {
	return append(append(accountFollowingSubstore, me...), types.KeySeparator...)
}

func getRewardKey(accKey types.AccountKey) []byte {
	return append(accountRewardSubstore, accKey...)
}

func getRelationshipKey(me types.AccountKey, other types.AccountKey) []byte {
	return append(getRelationshipPrefix(me), other...)
}

func getRelationshipPrefix(me types.AccountKey) []byte {
	return append(append(accountRelationshipSubstore, me...), types.KeySeparator...)
}

func getPendingCoinDayQueueKey(accKey types.AccountKey) []byte {
	return append(accountPendingCoinDayQueueSubstore, accKey...)
}

func getGrantPubKeyPrefix(me types.AccountKey) []byte {
	return append(append(accountGrantPubKeySubstore, me...), types.KeySeparator...)
}

func getGrantPubKeyKey(me types.AccountKey, pubKey crypto.PubKey) []byte {
	return append(getGrantPubKeyPrefix(me), hex.EncodeToString(pubKey.Bytes())...)
}

func getBalanceHistoryPrefix(me types.AccountKey) []byte {
	return append(append(accountBalanceHistorySubstore, me...), types.KeySeparator...)
}

func getBalanceHistoryKey(me types.AccountKey, bucketSlot int64) []byte {
	return strconv.AppendInt(getBalanceHistoryPrefix(me), bucketSlot, 10)
}

func getRewardHistoryPrefix(me types.AccountKey) []byte {
	return append(append(accountRewardHistorySubstore, me...), types.KeySeparator...)
}

func getRewardHistoryKey(me types.AccountKey, bucketSlot int64) []byte {
	return strconv.AppendInt(getRewardHistoryPrefix(me), bucketSlot, 10)
}

// IterateAccounts - iterate accounts in KVStore
func (as AccountStorage) IterateAccounts(ctx sdk.Context, process func(AccountInfo, AccountBank) (stop bool)) {
	store := ctx.KVStore(as.key)
	iter := sdk.KVStorePrefixIterator(store, accountInfoSubstore)
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		accInfo, err := as.GetInfo(ctx, types.AccountKey(val))
		if err != nil {
			panic(err)
		}
		accBank, err := as.GetBankFromAccountKey(ctx, types.AccountKey(val))
		if err != nil {
			panic(err)
		}
		if process(*accInfo, *accBank) {
			return
		}
		iter.Next()
	}
}
