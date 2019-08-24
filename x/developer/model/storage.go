package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/types"
)

const (
	trueStr = "t"
)

var (
	developerSubstore     = []byte{0x00}
	idaSubstore           = []byte{0x01}
	idaBalanceSubstore    = []byte{0x02}
	reservePoolSubstore   = []byte{0x03}
	affiliatedAccSubstore = []byte{0x04} // (app, user)
	userRoleSubstore      = []byte{0x05} // user role.
	idaStatsSubstore      = []byte{0x06}
)

// GetDeveloperKey - "developer substore" + "developer"
func GetDeveloperKey(accKey linotypes.AccountKey) []byte {
	return append(developerSubstore, accKey...)
}

// GetIDAKey - "ida substore" + "developer"
func GetIDAKey(accKey linotypes.AccountKey) []byte {
	return append(idaSubstore, accKey...)
}

// GetIDAStatsKey - "ida stats substore" + "developer"
func GetIDAStatsKey(accKey linotypes.AccountKey) []byte {
	return append(idaStatsSubstore, accKey...)
}

// GetIDABalanceKey - "ida balance substore" + "app" + "/" + "user"
func GetIDABalanceKey(app linotypes.AccountKey, user linotypes.AccountKey) []byte {
	prefix := append(idaBalanceSubstore, app...)
	return append(append(prefix, []byte(linotypes.KeySeparator)...), user...)
}

// GetAffiliatedAccAppPrefix - "affiliated account substore" + "app" + "/"
func GetAffiliatedAccAppPrefix(app linotypes.AccountKey) []byte {
	prefix := append(affiliatedAccSubstore, app...)
	return append(prefix, []byte(linotypes.KeySeparator)...)
}

// GetIDABalanceKey - "affiliated app prefix" + "user"
func GetAffiliatedAccKey(app linotypes.AccountKey, user linotypes.AccountKey) []byte {
	return append(GetAffiliatedAccAppPrefix(app), user...)
}

// GetUserRoleKey - "user" -> role
func GetUserRoleKey(user linotypes.AccountKey) []byte {
	return append(userRoleSubstore, user...)
}

// GetReservePoolKey - reserve pool key
func GetReservePoolKey() []byte {
	return reservePoolSubstore
}

// DeveloperStorage - developer storage
type DeveloperStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

// DeveloperStorage - new developer storage
func NewDeveloperStorage(key sdk.StoreKey) DeveloperStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	return DeveloperStorage{
		key: key,
		cdc: cdc,
	}
}

// HasDeveloper - check if developer in KVStore or not
func (ds DeveloperStorage) HasDeveloper(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	store := ctx.KVStore(ds.key)
	return store.Has(GetDeveloperKey(accKey))
}

// GetDeveloper - get developer from KVStore
func (ds DeveloperStorage) GetDeveloper(ctx sdk.Context, accKey linotypes.AccountKey) (*Developer, sdk.Error) {
	store := ctx.KVStore(ds.key)
	bz := store.Get(GetDeveloperKey(accKey))
	if bz == nil {
		return nil, types.ErrDeveloperNotFound()
	}
	dev := new(Developer)
	ds.cdc.MustUnmarshalBinaryLengthPrefixed(bz, dev)
	return dev, nil
}

// SetDeveloper - set developer to KVStore
func (ds DeveloperStorage) SetDeveloper(ctx sdk.Context, developer Developer) {
	store := ctx.KVStore(ds.key)
	developerByte := ds.cdc.MustMarshalBinaryLengthPrefixed(developer)
	store.Set(GetDeveloperKey(developer.Username), developerByte)
}

// GetAllDevelopers - get developer list from KVStore.
// NOTE, the result includes the all developers even if it's marked in value as deleted.
func (ds DeveloperStorage) GetAllDevelopers(ctx sdk.Context) []Developer {
	store := ctx.KVStore(ds.key)
	rst := make([]Developer, 0)
	itr := sdk.KVStorePrefixIterator(store, developerSubstore)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		dev := new(Developer)
		ds.cdc.MustUnmarshalBinaryLengthPrefixed(itr.Value(), dev)
		rst = append(rst, *dev)
	}
	return rst
}

// HasIDA - check if developer has a IDA in KVStore or not
func (ds DeveloperStorage) HasIDA(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	return ctx.KVStore(ds.key).Has(GetIDAKey(accKey))
}

// GetIDA - get ida of a developer KVStore
func (ds DeveloperStorage) GetIDA(ctx sdk.Context, accKey linotypes.AccountKey) (*AppIDA, sdk.Error) {
	store := ctx.KVStore(ds.key)
	bz := store.Get(GetIDAKey(accKey))
	if bz == nil {
		return nil, types.ErrIDANotFound()
	}
	ida := new(AppIDA)
	ds.cdc.MustUnmarshalBinaryLengthPrefixed(bz, ida)
	return ida, nil
}

// SetIDA - set IDA to KVStore
func (ds DeveloperStorage) SetIDA(ctx sdk.Context, ida AppIDA) {
	store := ctx.KVStore(ds.key)
	bz := ds.cdc.MustMarshalBinaryLengthPrefixed(ida)
	store.Set(GetIDAKey(ida.App), bz)
}

// GetIDABalance - get ida balance of (app, user)
func (ds DeveloperStorage) GetIDABank(ctx sdk.Context, app linotypes.AccountKey, user linotypes.AccountKey) *IDABank {
	store := ctx.KVStore(ds.key)
	bz := store.Get(GetIDABalanceKey(app, user))
	if bz == nil {
		return &IDABank{
			Balance:  linotypes.NewMiniDollar(0),
			Unauthed: false,
		}
	}
	rst := &IDABank{}
	ds.cdc.MustUnmarshalBinaryLengthPrefixed(bz, rst)
	return rst
}

// SetIDABalance - set (app, user)'s ida balance to amount
func (ds DeveloperStorage) SetIDABank(ctx sdk.Context, app linotypes.AccountKey, user linotypes.AccountKey, bank *IDABank) {
	store := ctx.KVStore(ds.key)
	bz := ds.cdc.MustMarshalBinaryLengthPrefixed(bank)
	store.Set(GetIDABalanceKey(app, user), bz)
}

// GetReservePool - get IDA's lino reserve pool.
func (ds DeveloperStorage) GetReservePool(ctx sdk.Context) *ReservePool {
	store := ctx.KVStore(ds.key)
	bz := store.Get(GetReservePoolKey())
	if bz == nil {
		panic("Developer IDA reserve pool MUST be initialized.")
	}
	rst := new(ReservePool)
	ds.cdc.MustUnmarshalBinaryLengthPrefixed(bz, rst)
	return rst
}

// SetReservePool - get IDA's lino reserve pool.
func (ds DeveloperStorage) SetReservePool(ctx sdk.Context, pool *ReservePool) {
	store := ctx.KVStore(ds.key)
	bz := ds.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(GetReservePoolKey(), bz)
}

// SetAffiliatedAcc - set affiliated account.
func (ds DeveloperStorage) SetAffiliatedAcc(ctx sdk.Context, app, user linotypes.AccountKey) {
	store := ctx.KVStore(ds.key)
	store.Set(GetAffiliatedAccKey(app, user), []byte(trueStr))
}

// HasAffiliateAcc - has this affiliated account..
func (ds DeveloperStorage) HasAffiliatedAcc(ctx sdk.Context, app, user linotypes.AccountKey) bool {
	store := ctx.KVStore(ds.key)
	return store.Has(GetAffiliatedAccKey(app, user))
}

// DelAffiliatedAcc - remove this affiliated acc.
func (ds DeveloperStorage) DelAffiliatedAcc(ctx sdk.Context, app, user linotypes.AccountKey) {
	store := ctx.KVStore(ds.key)
	store.Delete(GetAffiliatedAccKey(app, user))
}

func (ds DeveloperStorage) GetAllAffiliatedAcc(ctx sdk.Context, app linotypes.AccountKey) []linotypes.AccountKey {
	store := ctx.KVStore(ds.key)
	rst := make([]linotypes.AccountKey, 0)
	appPrefix := GetAffiliatedAccAppPrefix(app)
	itr := sdk.KVStorePrefixIterator(store, appPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		acc := itr.Key()[len(appPrefix):]
		rst = append(rst, linotypes.AccountKey(acc))
	}
	return rst
}

func (ds DeveloperStorage) SetUserRole(ctx sdk.Context, user linotypes.AccountKey, role *Role) {
	store := ctx.KVStore(ds.key)
	bz := ds.cdc.MustMarshalBinaryLengthPrefixed(role)
	store.Set(GetUserRoleKey(user), bz)
}

func (ds DeveloperStorage) DelUserRole(ctx sdk.Context, user linotypes.AccountKey) {
	ctx.KVStore(ds.key).Delete(GetUserRoleKey(user))
}

func (ds DeveloperStorage) GetUserRole(ctx sdk.Context, user linotypes.AccountKey) (*Role, sdk.Error) {
	bz := ctx.KVStore(ds.key).Get(GetUserRoleKey(user))
	if bz == nil {
		return nil, types.ErrInvalidUserRole()
	}
	rst := new(Role)
	ds.cdc.MustUnmarshalBinaryLengthPrefixed(bz, rst)
	return rst, nil
}

func (ds DeveloperStorage) HasUserRole(ctx sdk.Context, user linotypes.AccountKey) bool {
	return ctx.KVStore(ds.key).Has(GetUserRoleKey(user))
}

// SetIDAStats set ida stats of a app.
func (ds DeveloperStorage) SetIDAStats(ctx sdk.Context, app linotypes.AccountKey, stats AppIDAStats) {
	store := ctx.KVStore(ds.key)
	bz := ds.cdc.MustMarshalBinaryLengthPrefixed(stats)
	store.Set(GetIDAStatsKey(app), bz)
}

// GetIDAStats returns the stats of the IDA.
func (ds DeveloperStorage) GetIDAStats(ctx sdk.Context, app linotypes.AccountKey) *AppIDAStats {
	bz := ctx.KVStore(ds.key).Get(GetIDAStatsKey(app))
	if bz == nil {
		return &AppIDAStats{
			Total: linotypes.NewMiniDollar(0),
		}
	}
	stats := new(AppIDAStats)
	ds.cdc.MustUnmarshalBinaryLengthPrefixed(bz, stats)
	return stats
}
