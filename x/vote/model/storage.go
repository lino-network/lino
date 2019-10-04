package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
)

var (
	VoterSubstore = []byte{0x01}
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
	vs := VoteStorage{
		key: key,
		cdc: cdc,
	}

	return vs
}

// DoesVoterExist - check if voter exist in KVStore or not
func (vs VoteStorage) DoesVoterExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetVoterKey(accKey))
}

// GetVoter - get voter from KVStore
func (vs VoteStorage) GetVoter(ctx sdk.Context, accKey types.AccountKey) (*Voter, sdk.Error) {
	store := ctx.KVStore(vs.key)
	voterByte := store.Get(GetVoterKey(accKey))
	if voterByte == nil {
		return nil, ErrVoterNotFound()
	}
	voter := new(Voter)
	if err := vs.cdc.UnmarshalBinaryLengthPrefixed(voterByte, voter); err != nil {
		return nil, ErrFailedToUnmarshalVoter(err)
	}
	return voter, nil
}

// SetVoter - set voter to KVStore
func (vs VoteStorage) SetVoter(ctx sdk.Context, accKey types.AccountKey, voter *Voter) sdk.Error {
	store := ctx.KVStore(vs.key)
	voterByte, err := vs.cdc.MarshalBinaryLengthPrefixed(*voter)
	if err != nil {
		return ErrFailedToMarshalVoter(err)
	}
	store.Set(GetVoterKey(accKey), voterByte)
	return nil
}

// DeleteVoter - delete voter from KVStore
func (vs VoteStorage) DeleteVoter(ctx sdk.Context, username types.AccountKey) sdk.Error {
	store := ctx.KVStore(vs.key)
	store.Delete(GetVoterKey(username))
	return nil
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
	}
	return utils.NewStoreMap(substores)
}

// GetVoterKey - "voter substore" + "voter"
func GetVoterKey(me types.AccountKey) []byte {
	return append(VoterSubstore, me...)
}
