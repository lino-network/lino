package model

import (
	"strconv"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

var (
	accountInfoSubstore              = []byte{0x00}
	accountBankSubstore              = []byte{0x01}
	accountMetaSubstore              = []byte{0x02}
	accountFollowerSubstore          = []byte{0x03}
	accountFollowingSubstore         = []byte{0x04}
	accountRewardSubstore            = []byte{0x05}
	accountPendingStakeQueueSubstore = []byte{0x06}
	accountRelationshipSubstore      = []byte{0x07}
	accountGrantListSubstore         = []byte{0x08}
	accountBalanceHistorySubstore    = []byte{0x09}
)

type AccountStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewLinoAccountStorage creates and returns a account manager
func NewAccountStorage(key sdk.StoreKey) AccountStorage {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)

	cdc.RegisterInterface((*Detail)(nil), nil)
	cdc.RegisterConcrete(BalanceIn{}, "transfer/in", nil)
	cdc.RegisterConcrete(BalanceOut{}, "transfer/out", nil)

	cdc.RegisterInterface((*types.TransferObject)(nil), nil)
	cdc.RegisterConcrete(types.AccountKey(""), "transfer/to/acckey", nil)
	cdc.RegisterConcrete(types.PermLink(""), "transfer/to/permlink", nil)
	cdc.RegisterConcrete(types.InternalObject(""), "transfer/to/internal", nil)

	return AccountStorage{
		key: key,
		cdc: cdc,
	}
}

// AccountExist returns true when a specific account exist in the KVStore.
func (as AccountStorage) AccountExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	return store.Has(GetAccountInfoKey(accKey))
}

// GetInfo returns general account info of a specific account, returns error otherwise.
func (as AccountStorage) GetInfo(ctx sdk.Context, accKey types.AccountKey) (*AccountInfo, sdk.Error) {
	store := ctx.KVStore(as.key)
	infoByte := store.Get(GetAccountInfoKey(accKey))
	if infoByte == nil {
		return nil, ErrAccountInfoNotFound()
	}
	info := new(AccountInfo)
	if err := as.cdc.UnmarshalJSON(infoByte, info); err != nil {
		return nil, ErrGetAccountInfo().TraceCause(err, "")
	}
	return info, nil
}

// SetInfo sets general account info to a specific account, returns error if any.
func (as AccountStorage) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *AccountInfo) sdk.Error {
	store := ctx.KVStore(as.key)
	infoByte, err := as.cdc.MarshalJSON(*accInfo)
	if err != nil {
		return ErrSetInfoFailed()
	}
	store.Set(GetAccountInfoKey(accKey), infoByte)
	return nil
}

// GetBankFromAccountKey returns bank info of a specific account, returns error
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
		return nil, ErrGetBankFromAccountKey().TraceCause(err, "")
	}
	return bank, nil
}

// SetBankFromAddress sets bank info for a given address,
// returns error if any.
func (as AccountStorage) SetBankFromAccountKey(ctx sdk.Context, username types.AccountKey, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(as.key)
	bankByte, err := as.cdc.MarshalJSON(*accBank)
	if err != nil {
		return ErrSetBankFailed().TraceCause(err, "")
	}
	store.Set(GetAccountBankKey(username), bankByte)
	return nil
}

// GetMeta returns meta of a given account that are tiny
// and frequently updated fields.
func (as AccountStorage) GetMeta(ctx sdk.Context, accKey types.AccountKey) (*AccountMeta, sdk.Error) {
	store := ctx.KVStore(as.key)
	metaByte := store.Get(GetAccountMetaKey(accKey))
	if metaByte == nil {
		return nil, ErrGetMetaFailed()
	}
	meta := new(AccountMeta)
	if err := as.cdc.UnmarshalJSON(metaByte, meta); err != nil {
		return nil, ErrGetMetaFailed().TraceCause(err, "")
	}
	return meta, nil
}

// SetMeta sets meta for a given account, returns error if any.
func (as AccountStorage) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *AccountMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(*accMeta)
	if err != nil {
		return ErrSetMetaFailed().TraceCause(err, "")
	}
	store.Set(GetAccountMetaKey(accKey), metaByte)
	return nil
}

// IsMyfollower returns true if `follower` follows `me`.
func (as AccountStorage) IsMyFollower(ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	key := getFollowerKey(me, follower)
	return store.Has(key)
}

// SetFollowerMeta sets follower meta info for a given account which includes
// time and follower name.
func (as AccountStorage) SetFollowerMeta(ctx sdk.Context, me types.AccountKey, meta FollowerMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrSetFollowerMeta().TraceCause(err, "")
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

// IsMyFollowing returns true if `me` follows `following`
func (as AccountStorage) IsMyFollowing(ctx sdk.Context, me types.AccountKey, following types.AccountKey) bool {
	store := ctx.KVStore(as.key)
	key := getFollowingKey(me, following)
	return store.Has(key)
}

// SetFollowerMeta sets following meta info for a given account which includes
// time and following name.
func (as AccountStorage) SetFollowingMeta(ctx sdk.Context, me types.AccountKey, meta FollowingMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalJSON(meta)
	if err != nil {
		return ErrSetFollowingMeta().TraceCause(err, "")
	}
	store.Set(getFollowingKey(me, meta.FollowingName), metaByte)
	return nil
}

// RemoveFollowingMeta removes following meta info of a relationship.
func (as AccountStorage) RemoveFollowingMeta(ctx sdk.Context, me types.AccountKey, following types.AccountKey) {
	store := ctx.KVStore(as.key)
	store.Delete(getFollowingKey(me, following))
	return
}

// GetReward returns reward info of a given account, returns error if any.
func (as AccountStorage) GetReward(ctx sdk.Context, accKey types.AccountKey) (*Reward, sdk.Error) {
	store := ctx.KVStore(as.key)
	rewardByte := store.Get(getRewardKey(accKey))
	if rewardByte == nil {
		return nil, ErrGetRewardFailed()
	}
	reward := new(Reward)
	if err := as.cdc.UnmarshalJSON(rewardByte, reward); err != nil {
		return nil, ErrGetRewardFailed().TraceCause(err, "")
	}
	return reward, nil
}

// SetReward sets the rewards info of a given account, returns error if any.
func (as AccountStorage) SetReward(ctx sdk.Context, accKey types.AccountKey, reward *Reward) sdk.Error {
	store := ctx.KVStore(as.key)
	rewardByte, err := as.cdc.MarshalJSON(*reward)
	if err != nil {
		return ErrSetRewardFailed().TraceCause(err, "")
	}
	store.Set(getRewardKey(accKey), rewardByte)
	return nil
}

// GetPendingStakeQueue returns a pending stake queue for a given address.
func (as AccountStorage) GetPendingStakeQueue(
	ctx sdk.Context, me types.AccountKey) (*PendingStakeQueue, sdk.Error) {
	store := ctx.KVStore(as.key)
	pendingStakeQueueByte := store.Get(getPendingStakeQueueKey(me))
	if pendingStakeQueueByte == nil {
		return nil, ErrGetPendingStakeFailed()
	}
	queue := new(PendingStakeQueue)
	if err := as.cdc.UnmarshalJSON(pendingStakeQueueByte, queue); err != nil {
		return nil, ErrGetPendingStakeFailed().TraceCause(err, "")
	}
	return queue, nil
}

// SetPendingStakeQueue sets a pending stake queue for a given username.
func (as AccountStorage) SetPendingStakeQueue(ctx sdk.Context, me types.AccountKey, pendingStakeQueue *PendingStakeQueue) sdk.Error {
	store := ctx.KVStore(as.key)
	pendingStakeQueueByte, err := as.cdc.MarshalJSON(*pendingStakeQueue)
	if err != nil {
		return ErrSetRewardFailed().TraceCause(err, "")
	}
	store.Set(getPendingStakeQueueKey(me), pendingStakeQueueByte)
	return nil
}

// SetGrantKeyList sets a list of grant public keys for a given account.
func (as AccountStorage) SetGrantKeyList(ctx sdk.Context, me types.AccountKey, grantKeyList *GrantKeyList) sdk.Error {
	store := ctx.KVStore(as.key)
	GrantKeyListByte, err := as.cdc.MarshalJSON(*grantKeyList)
	if err != nil {
		return ErrSetGrantListFailed().TraceCause(err, "")
	}
	store.Set(getGrantKeyListKey(me), GrantKeyListByte)
	return nil
}

// GetGrantKeyList returns a list of grant public keys for a given account.
func (as AccountStorage) GetGrantKeyList(ctx sdk.Context, me types.AccountKey) (*GrantKeyList, sdk.Error) {
	store := ctx.KVStore(as.key)
	grantKeyListByte := store.Get(getGrantKeyListKey(me))
	if grantKeyListByte == nil {
		return nil, ErrGetGrantListFailed()
	}
	grantKeyList := new(GrantKeyList)
	if err := as.cdc.UnmarshalJSON(grantKeyListByte, grantKeyList); err != nil {
		return nil, ErrGetGrantListFailed().TraceCause(err, "")
	}
	return grantKeyList, nil
}

// GetRelationship returns the relationship between two accounts.
func (as AccountStorage) GetRelationship(ctx sdk.Context, me types.AccountKey, other types.AccountKey) (*Relationship, sdk.Error) {
	store := ctx.KVStore(as.key)
	relationshipByte := store.Get(getRelationshipKey(me, other))
	if relationshipByte == nil {
		return nil, nil
	}
	queue := new(Relationship)
	if err := as.cdc.UnmarshalJSON(relationshipByte, queue); err != nil {
		return nil, ErrGetRelationshipFailed().TraceCause(err, "")
	}
	return queue, nil
}

// SetRelationship sets relationship for two accounts.
func (as AccountStorage) SetRelationship(ctx sdk.Context, me types.AccountKey, other types.AccountKey, relationship *Relationship) sdk.Error {
	store := ctx.KVStore(as.key)
	relationshipByte, err := as.cdc.MarshalJSON(*relationship)
	if err != nil {
		return ErrSetRelationshipFailed().TraceCause(err, "")
	}
	store.Set(getRelationshipKey(me, other), relationshipByte)
	return nil
}

// GetRelationship returns the relationship between two accounts.
func (as AccountStorage) GetBalanceHistory(
	ctx sdk.Context, me types.AccountKey, transactionSlot int64) (*BalanceHistory, sdk.Error) {
	store := ctx.KVStore(as.key)
	balanceHistoryBytes := store.Get(getBalanceHistoryKey(me, transactionSlot))
	if balanceHistoryBytes == nil {
		return nil, nil
	}
	history := new(BalanceHistory)
	if err := as.cdc.UnmarshalJSON(balanceHistoryBytes, history); err != nil {
		return nil, ErrGetBalanceHistoryFailed().TraceCause(err, "")
	}
	return history, nil
}

// SetBalanceHistory sets balance history.
func (as AccountStorage) SetBalanceHistory(
	ctx sdk.Context, me types.AccountKey, timeSlot int64, history *BalanceHistory) sdk.Error {
	store := ctx.KVStore(as.key)
	historyBytes, err := as.cdc.MarshalJSON(*history)
	if err != nil {
		return ErrSetBalanceHistoryFailed().TraceCause(err, "")
	}
	store.Set(getBalanceHistoryKey(me, timeSlot), historyBytes)
	return nil
}

func GetAccountInfoKey(accKey types.AccountKey) []byte {
	return append(accountInfoSubstore, accKey...)
}

func GetAccountBankKey(accKey types.AccountKey) []byte {
	return append(accountBankSubstore, accKey...)
}

func GetAccountMetaKey(accKey types.AccountKey) []byte {
	return append(accountMetaSubstore, accKey...)
}

// "follower substore" + "me" + "my follower"
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

func getPendingStakeQueueKey(accKey types.AccountKey) []byte {
	return append(accountPendingStakeQueueSubstore, accKey...)
}

func getGrantKeyListKey(me types.AccountKey) []byte {
	return append(accountGrantListSubstore, me...)
}

func getBalanceHistoryPrefix(me types.AccountKey) []byte {
	return append(append(accountBalanceHistorySubstore, me...), types.KeySeparator...)
}

func getBalanceHistoryKey(me types.AccountKey, atWhen int64) []byte {
	return strconv.AppendInt(getBalanceHistoryPrefix(me), atWhen, 10)
}
