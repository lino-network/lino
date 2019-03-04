package model

import (
	"fmt"
	"strings"

	"github.com/lino-network/lino/types"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	accountInfoSubstore                = []byte{0x00}
	accountBankSubstore                = []byte{0x01}
	accountMetaSubstore                = []byte{0x02}
	accountRewardSubstore              = []byte{0x03}
	accountPendingCoinDayQueueSubstore = []byte{0x04}
	accountGrantPubKeySubstore         = []byte{0x05}
	// XXX(yukai): deprecated.
	// accountFollowerSubstore            = []byte{0x03}
	// accountFollowingSubstore           = []byte{0x04}
	// XXX(yukai): deprecated.
	// accountRelationshipSubstore        = []byte{0x07}
	// XXX(yukai): deprecated.
	// accountBalanceHistorySubstore      = []byte{0x08}
	// XXX(yukai): deprecated.
	// accountRewardHistorySubstore = []byte{0x0a}
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
	if err := as.cdc.UnmarshalBinaryLengthPrefixed(infoByte, info); err != nil {
		return nil, ErrFailedToUnmarshalAccountInfo(err)
	}
	return info, nil
}

// SetInfo - sets general account info to a specific account, returns error if any.
func (as AccountStorage) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *AccountInfo) sdk.Error {
	store := ctx.KVStore(as.key)
	infoByte, err := as.cdc.MarshalBinaryLengthPrefixed(*accInfo)
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
	if err := as.cdc.UnmarshalBinaryLengthPrefixed(bankByte, bank); err != nil {
		return nil, ErrFailedToUnmarshalAccountBank(err)
	}
	return bank, nil
}

// SetBankFromAddress - sets bank info for a given address,
// returns error if any.
func (as AccountStorage) SetBankFromAccountKey(ctx sdk.Context, username types.AccountKey, accBank *AccountBank) sdk.Error {
	store := ctx.KVStore(as.key)
	bankByte, err := as.cdc.MarshalBinaryLengthPrefixed(*accBank)
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
	if err := as.cdc.UnmarshalBinaryLengthPrefixed(metaByte, meta); err != nil {
		return nil, ErrFailedToUnmarshalAccountMeta(err)
	}
	return meta, nil
}

// SetMeta - sets meta for a given account, returns error if any.
func (as AccountStorage) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *AccountMeta) sdk.Error {
	store := ctx.KVStore(as.key)
	metaByte, err := as.cdc.MarshalBinaryLengthPrefixed(*accMeta)
	if err != nil {
		return ErrFailedToMarshalAccountMeta(err)
	}
	store.Set(GetAccountMetaKey(accKey), metaByte)
	return nil
}

// GetReward - returns reward info of a given account, returns error if any.
func (as AccountStorage) GetReward(ctx sdk.Context, accKey types.AccountKey) (*Reward, sdk.Error) {
	store := ctx.KVStore(as.key)
	rewardByte := store.Get(getRewardKey(accKey))
	if rewardByte == nil {
		return nil, ErrRewardNotFound()
	}
	reward := new(Reward)
	if err := as.cdc.UnmarshalBinaryLengthPrefixed(rewardByte, reward); err != nil {
		return nil, ErrFailedToUnmarshalReward(err)
	}
	return reward, nil
}

// SetReward - sets the rewards info of a given account, returns error if any.
func (as AccountStorage) SetReward(ctx sdk.Context, accKey types.AccountKey, reward *Reward) sdk.Error {
	store := ctx.KVStore(as.key)
	rewardByte, err := as.cdc.MarshalBinaryLengthPrefixed(*reward)
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
	if err := as.cdc.UnmarshalBinaryLengthPrefixed(pendingCoinDayQueueByte, queue); err != nil {
		return nil, ErrFailedToUnmarshalPendingCoinDayQueue(err)
	}
	return queue, nil
}

// SetPendingCoinDayQueue - sets a pending coin day queue for a given username.
func (as AccountStorage) SetPendingCoinDayQueue(ctx sdk.Context, me types.AccountKey, pendingCoinDayQueue *PendingCoinDayQueue) sdk.Error {
	store := ctx.KVStore(as.key)
	pendingCoinDayQueueByte, err := as.cdc.MarshalBinaryLengthPrefixed(*pendingCoinDayQueue)
	if err != nil {
		return ErrFailedToMarshalPendingCoinDayQueue(err)
	}
	store.Set(getPendingCoinDayQueueKey(me), pendingCoinDayQueueByte)
	return nil
}

// DeleteAllGrantPubKeys - deletes all grant pubkeys from a granted user in KV.
func (as AccountStorage) DeleteAllGrantPubKeys(ctx sdk.Context, me types.AccountKey, grantTo types.AccountKey) {
	store := ctx.KVStore(as.key)
	store.Delete(getGrantPubKeysKey(me, grantTo))
	return
}

// GetGrantPubKey - returns grant user info keyed with pubkey.
func (as AccountStorage) GetGrantPubKeys(ctx sdk.Context, me types.AccountKey, grantTo types.AccountKey) ([]*GrantPubKey, sdk.Error) {
	store := ctx.KVStore(as.key)
	grantPubKeyByte := store.Get(getGrantPubKeysKey(me, grantTo))
	if grantPubKeyByte == nil {
		return nil, ErrGrantPubKeyNotFound()
	}
	grantPubKeys := new([]*GrantPubKey)
	if err := as.cdc.UnmarshalBinaryLengthPrefixed(grantPubKeyByte, grantPubKeys); err != nil {
		return nil, ErrFailedToUnmarshalGrantPubKey(err)
	}
	return *grantPubKeys, nil
}

// GetAllGrantPubKeys - returns grant user info keyed with pubkey.
func (as AccountStorage) GetAllGrantPubKeys(ctx sdk.Context, me types.AccountKey) ([]*GrantPubKey, sdk.Error) {
	grantPubKeys := make([]*GrantPubKey, 0)
	store := ctx.KVStore(as.key)
	iter := sdk.KVStorePrefixIterator(store, getGrantPubKeyPrefix(me))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return grantPubKeys, nil
		}
		// fmt.Println(string(key[len(getGrantPubKeyPrefix(me)):]))
		val := iter.Value()
		grantPubKeyList := new([]*GrantPubKey)
		if err := as.cdc.UnmarshalBinaryLengthPrefixed(val, grantPubKeyList); err != nil {
			return nil, ErrFailedToUnmarshalGrantPubKey(err)
		}
		grantPubKeys = append(grantPubKeys, *grantPubKeyList...)
		iter.Next()
	}
	return grantPubKeys, nil
}

// SetGrantPubKey - sets a grant user to KV. Key is pubkey and value is grant user info
func (as AccountStorage) SetGrantPubKeys(ctx sdk.Context, me types.AccountKey, grantTo types.AccountKey, grantPubKeys []*GrantPubKey) sdk.Error {
	fmt.Println("===> grant permission", me, grantTo, grantPubKeys)
	store := ctx.KVStore(as.key)
	grantPubKeyByte, err := as.cdc.MarshalBinaryLengthPrefixed(grantPubKeys)
	if err != nil {
		return ErrFailedToMarshalGrantPubKey(err)
	}
	store.Set(getGrantPubKeysKey(me, grantTo), grantPubKeyByte)
	return nil
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

func getRewardKey(accKey types.AccountKey) []byte {
	return append(accountRewardSubstore, accKey...)
}

func getPendingCoinDayQueueKey(accKey types.AccountKey) []byte {
	return append(accountPendingCoinDayQueueSubstore, accKey...)
}

func getGrantPubKeyPrefix(me types.AccountKey) []byte {
	return append(append(accountGrantPubKeySubstore, me...), types.KeySeparator...)
}

func getGrantPubKeysKey(me types.AccountKey, grantTo types.AccountKey) []byte {
	return append(append(getGrantPubKeyPrefix(me), grantTo...), types.KeySeparator...)
}

// Export to table representation.
func (as AccountStorage) Export(ctx sdk.Context) *AccountTables {
	tables := &AccountTables{}
	store := ctx.KVStore(as.key)
	// export tables.account
	func() {
		itr := sdk.KVStorePrefixIterator(store, accountInfoSubstore)
		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			k, _ := itr.Key(), itr.Value()
			username := types.AccountKey(k[1:])

			accInfo, err := as.GetInfo(ctx, username)
			if err != nil {
				panic(err)
			}

			accBank, err := as.GetBankFromAccountKey(ctx, username)
			if err != nil {
				panic(err)
			}

			accMeta, err := as.GetMeta(ctx, username)
			if err != nil {
				panic(err)
			}

			accPending, err := as.GetPendingCoinDayQueue(ctx, username)
			if err != nil {
				panic(err)
			}

			reward, err := as.GetReward(ctx, username)
			if err != nil {
				panic(err)
			}

			// set all states
			accRow := AccountRow{
				Username:            username,
				Info:                *accInfo,
				Bank:                *accBank,
				Meta:                *accMeta,
				Reward:              *reward,
				PendingCoinDayQueue: *accPending,
			}
			tables.Accounts = append(tables.Accounts, accRow)
		}
	}()
	// export tables.GrantPubKeys
	func() {
		itr := sdk.KVStorePrefixIterator(store, accountGrantPubKeySubstore)
		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			usernamePubKeyPermission := string(itr.Key()[1:])
			strs := strings.Split(usernamePubKeyPermission, types.KeySeparator)
			if len(strs) != 3 {
				panic("illegat usernamePubkeyAndPermission: " + usernamePubKeyPermission)
			}
			// username, grantTo, permission := types.AccountKey(strs[0]), types.AccountKey(strs[1]), types.Permission(strs[2])
			// pubKeyBytes, err := hex.DecodeString(pubKeyHex)
			// if err != nil {
			// 	panic("Failed to decode pubkeyHex: " + pubKeyHex + " " + err.Error())
			// }
			// var pubKey crypto.PubKey
			// err = as.cdc.UnmarshalBinaryLengthPrefixed(pubKeyBytes, &pubKey)
			// if err != nil {
			// 	panic("Faield to decode pubkeyBytes to pubkey interface: " + err.Error())
			// }

			// xxx(yukai): need update
			// info, err := as.GetGrantPubKey(ctx, username, grantTo, permission)
			// if err != nil {
			// 	panic("failed GetGrantPubKey: " + err.Error())
			// }
			// row := GrantPubKeyRow{
			// 	Username:    username,
			// 	PubKey:      pubKey,
			// 	GrantPubKey: *info,
			// }
			// tables.AccountGrantPubKeys = append(tables.AccountGrantPubKeys, row)
		}
	}()
	return tables
}

// Import from tablesIR.
func (as AccountStorage) Import(ctx sdk.Context, tb *AccountTablesIR) {
	check := func(err error) {
		if err != nil {
			panic("[as] Failed to import: " + err.Error())
		}
	}
	// import table.accounts
	for _, v := range tb.Accounts {
		err := as.SetInfo(ctx, v.Username, &v.Info)
		check(err)
		err = as.SetBankFromAccountKey(ctx, v.Username, &v.Bank)
		check(err)
		err = as.SetMeta(ctx, v.Username, &v.Meta)
		check(err)
		err = as.SetReward(ctx, v.Username, &v.Reward)
		check(err)
		q := &PendingCoinDayQueue{
			LastUpdatedAt:   v.PendingCoinDayQueue.LastUpdatedAt,
			TotalCoinDay:    sdk.MustNewDecFromStr(v.PendingCoinDayQueue.TotalCoinDay),
			TotalCoin:       v.PendingCoinDayQueue.TotalCoin,
			PendingCoinDays: v.PendingCoinDayQueue.PendingCoinDays,
		}
		err = as.SetPendingCoinDayQueue(ctx, v.Username, q)
		check(err)
	}
	// import AccountGrantPubKeys
	pubKeyMap := make(map[types.AccountKey][]*GrantPubKey)
	for _, v := range tb.AccountGrantPubKeys {
		if _, ok := pubKeyMap[v.GrantPubKey.GrantTo]; ok {
			pubKeyMap[v.GrantPubKey.GrantTo] = append(pubKeyMap[v.GrantPubKey.GrantTo], &v.GrantPubKey)
		} else {
			pubKeyMap[v.GrantPubKey.GrantTo] = []*GrantPubKey{&v.GrantPubKey}
		}
	}
	// for grantTo, grantPubKeyList := range pubKeyMap {
	// 	err := as.SetGrantPubKeys(ctx, v.Username, grantTo, grantPubKeyList)
	// 	check(err)
	// }
}

// IterateAccounts - iterate accounts in KVStore
func (as AccountStorage) IterateAccounts(ctx sdk.Context, process func(AccountInfo, AccountBank) (stop bool)) {
	store := ctx.KVStore(as.key)
	iter := sdk.KVStorePrefixIterator(store, accountInfoSubstore)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value() // TODO(yumin): why value not key?
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
