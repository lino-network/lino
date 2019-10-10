package model

import (
	"strconv"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/vote/types"
)

var (
	VoterSubstore         = []byte{0x01} // SubStore for voter info.
	LinoStakeStatSubStore = []byte{0x02} // SubStore for lino stake statistic
)

// VoteStorage - vote storage
type VoteStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewVoteStorage - new vote storage
func NewVoteStorage(key sdk.StoreKey) VoteStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	cdc.Seal()

	vs := VoteStorage{
		key: key,
		cdc: cdc,
	}

	return vs
}

// DoesVoterExist - check if voter exist in KVStore or not
func (vs VoteStorage) DoesVoterExist(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetVoterKey(accKey))
}

// GetVoter - get voter from KVStore
func (vs VoteStorage) GetVoter(ctx sdk.Context, accKey linotypes.AccountKey) (*Voter, sdk.Error) {
	store := ctx.KVStore(vs.key)
	voterByte := store.Get(GetVoterKey(accKey))
	if voterByte == nil {
		return nil, types.ErrVoterNotFound()
	}
	voter := new(Voter)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(voterByte, voter)
	return voter, nil
}

// SetVoter - set voter to KVStore
func (vs VoteStorage) SetVoter(ctx sdk.Context, accKey linotypes.AccountKey, voter *Voter) {
	store := ctx.KVStore(vs.key)
	voterByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*voter)
	store.Set(GetVoterKey(accKey), voterByte)
}

// // DeleteVoter - delete voter from KVStore
// // should never be deleted.
// func (vs VoteStorage) DeleteVoter(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
// 	store := ctx.KVStore(vs.key)
// 	store.Delete(GetVoterKey(username))
// 	return nil
// }

// SetLinoStakeStat - set lino power statistic at given day
func (vs VoteStorage) SetLinoStakeStat(ctx sdk.Context, day int64, lps *LinoStakeStat) {
	store := ctx.KVStore(vs.key)
	lpsByte := vs.cdc.MustMarshalBinaryLengthPrefixed(*lps)
	store.Set(GetLinoStakeStatKey(day), lpsByte)
}

// GetLinoStakeStat - get lino power statistic at given day
func (vs VoteStorage) GetLinoStakeStat(ctx sdk.Context, day int64) (*LinoStakeStat, sdk.Error) {
	store := ctx.KVStore(vs.key)
	bz := store.Get(GetLinoStakeStatKey(day))
	if bz == nil {
		return nil, types.ErrStakeStatNotFound(day)
	}
	linoStakeStat := new(LinoStakeStat)
	vs.cdc.MustUnmarshalBinaryLengthPrefixed(bz, linoStakeStat)
	return linoStakeStat, nil
}

// StoreMap - map of all substores
func (vs VoteStorage) StoreMap(ctx sdk.Context) utils.StoreMap {
	store := ctx.KVStore(vs.key)
	substores := []utils.SubStore{
		{
			Store:      store,
			Prefix:     VoterSubstore,
			ValCreator: func() interface{} { return new(Voter) },
			Decoder:    vs.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
		{
			Store:      store,
			Prefix:     LinoStakeStatSubStore,
			ValCreator: func() interface{} { return new(LinoStakeStat) },
			Decoder:    vs.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
	}
	return utils.NewStoreMap(substores)
}

// GetLinoStakeStatKey - get lino power statistic at day from KVStore
func GetLinoStakeStatKey(day int64) []byte {
	return append(LinoStakeStatSubStore, strconv.FormatInt(day, 10)...)
}

// GetVoterKey - "voter substore" + "voter"
func GetVoterKey(me linotypes.AccountKey) []byte {
	return append(VoterSubstore, me...)
}
