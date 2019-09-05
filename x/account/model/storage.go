package model

import (
	"github.com/lino-network/lino/types"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	accountInfoSubstore        = []byte{0x00}
	accountBankSubstore        = []byte{0x01}
	accountMetaSubstore        = []byte{0x02}
	accountGrantPubKeySubstore = []byte{0x03}
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
	as.cdc.MustUnmarshalBinaryLengthPrefixed(infoByte, info)
	return info, nil
}

// SetInfo - sets general account info to a specific account, returns error if any.
func (as AccountStorage) SetInfo(ctx sdk.Context, accKey types.AccountKey, accInfo *AccountInfo) {
	store := ctx.KVStore(as.key)
	infoByte := as.cdc.MustMarshalBinaryLengthPrefixed(*accInfo)
	store.Set(GetAccountInfoKey(accKey), infoByte)
}

// GetBank - returns bank info of a specific address, returns error if any.
func (as AccountStorage) GetBank(ctx sdk.Context, addr sdk.Address) (*AccountBank, sdk.Error) {
	store := ctx.KVStore(as.key)
	bankByte := store.Get(GetAccountBankKey(addr))
	if bankByte == nil {
		return nil, ErrAccountBankNotFound()
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
func (as AccountStorage) GetMeta(ctx sdk.Context, accKey types.AccountKey) *AccountMeta {
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
func (as AccountStorage) SetMeta(ctx sdk.Context, accKey types.AccountKey, accMeta *AccountMeta) {
	store := ctx.KVStore(as.key)
	metaByte := as.cdc.MustMarshalBinaryLengthPrefixed(*accMeta)
	store.Set(GetAccountMetaKey(accKey), metaByte)
}

// DeleteAllGrantPermissions - deletes all grant pubkeys from a granted user in KV.
func (as AccountStorage) DeleteAllGrantPermissions(ctx sdk.Context, me types.AccountKey, grantTo types.AccountKey) {
	store := ctx.KVStore(as.key)
	store.Delete(getGrantPermKey(me, grantTo))
	return
}

// GetGrantPermissions - returns grant user info keyed with pubkey.
func (as AccountStorage) GetGrantPermissions(ctx sdk.Context, me types.AccountKey, grantTo types.AccountKey) ([]*GrantPermission, sdk.Error) {
	store := ctx.KVStore(as.key)
	grantPubKeyByte := store.Get(getGrantPermKey(me, grantTo))
	if grantPubKeyByte == nil {
		return nil, ErrGrantPubKeyNotFound()
	}
	grantPubKeys := new([]*GrantPermission)
	as.cdc.MustUnmarshalBinaryLengthPrefixed(grantPubKeyByte, grantPubKeys)
	return *grantPubKeys, nil
}

// GetAllGrantPermissions - returns grant user info keyed with pubkey.
func (as AccountStorage) GetAllGrantPermissions(ctx sdk.Context, me types.AccountKey) ([]*GrantPermission, sdk.Error) {
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
func (as AccountStorage) SetGrantPermissions(ctx sdk.Context, me types.AccountKey, grantTo types.AccountKey, grantPubKeys []*GrantPermission) {
	store := ctx.KVStore(as.key)
	grantPermByte := as.cdc.MustMarshalBinaryLengthPrefixed(grantPubKeys)
	store.Set(getGrantPermKey(me, grantTo), grantPermByte)
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
func GetAccountBankKey(addr sdk.Address) []byte {
	return append(accountBankSubstore, addr.Bytes()...)
}

// GetAccountMetaKey - "account meta substore" + "username"
func GetAccountMetaKey(accKey types.AccountKey) []byte {
	return append(accountMetaSubstore, accKey...)
}

func getGrantPermPrefix(me types.AccountKey) []byte {
	return append(append(accountGrantPubKeySubstore, me...), types.KeySeparator...)
}

func getGrantPermKey(me types.AccountKey, grantTo types.AccountKey) []byte {
	return append(append(getGrantPermPrefix(me), grantTo...), types.KeySeparator...)
}

// // Export to table representation.
// func (as AccountStorage) Export(ctx sdk.Context) *AccountTables {
// 	tables := &AccountTables{}
// 	store := ctx.KVStore(as.key)
// 	// export tables.account
// 	func() {
// 		itr := sdk.KVStorePrefixIterator(store, accountInfoSubstore)
// 		defer itr.Close()
// 		for ; itr.Valid(); itr.Next() {
// 			k, _ := itr.Key(), itr.Value()
// 			username := types.AccountKey(k[1:])

// 			info, err := as.GetInfo(ctx, username)
// 			if err != nil {
// 				panic(err)
// 			}

// 			bank, err := as.GetBank(ctx, info.Address)
// 			if err != nil {
// 				panic(err)
// 			}

// 			meta, err := as.GetMeta(ctx, username)
// 			if err != nil {
// 				panic(err)
// 			}

// 			reward, err := as.GetReward(ctx, username)
// 			if err != nil {
// 				panic(err)
// 			}

// 			// set all states
// 			accRow := AccountRow{
// 				Username: username,
// 				Info:     *info,
// 				Bank:     *bank,
// 				Meta:     *meta,
// 				Reward:   *reward,
// 			}
// 			tables.Accounts = append(tables.Accounts, accRow)
// 		}
// 	}()
// 	// export tables.GrantPubKeys
// 	func() {
// 		itr := sdk.KVStorePrefixIterator(store, accountGrantPubKeySubstore)
// 		defer itr.Close()
// 		for ; itr.Valid(); itr.Next() {
// 			usernameApp := string(itr.Key()[1:])
// 			strs := strings.Split(usernameApp, types.KeySeparator)
// 			if len(strs) != 3 {
// 				panic("illegat usernamePubkeyAndPermission: " + usernameApp)
// 			}
// 			username, app := types.AccountKey(strs[0]), types.AccountKey(strs[1])
// 			permissions, err := as.GetGrantPermissions(ctx, username, app)
// 			if err != nil {
// 				panic("failed to fetch permission for " + username + " and " + app)
// 			}
// 			for _, v := range permissions {
// 				row := GrantPubKeyRow{
// 					Username:    username,
// 					PubKey:      nil, // PubKey is deprecated since upgrade1
// 					GrantPubKey: *v,
// 				}
// 				tables.AccountGrantPubKeys = append(tables.AccountGrantPubKeys, row)
// 			}
// 		}
// 	}()
// 	return tables
// }

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
		accBank, err := as.GetBank(ctx, accInfo.Address)
		if err != nil {
			panic(err)
		}
		if process(*accInfo, *accBank) {
			return
		}
		iter.Next()
	}
}
