package model

import (
	"github.com/lino-network/lino/types"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	voterSubstore = []byte{0x01}
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

// Export - Export voter state
func (vs VoteStorage) Export(ctx sdk.Context) *VoterTables {
	tables := &VoterTables{}
	store := ctx.KVStore(vs.key)
	// export table.voters
	func() {
		itr := sdk.KVStorePrefixIterator(store, voterSubstore)
		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			k := itr.Key()
			username := types.AccountKey(k[1:])
			val, err := vs.GetVoter(ctx, username)
			if err != nil {
				panic("failed to read voter: " + err.Error())
			}
			row := VoterRow{
				Username: username,
				Voter:    *val,
			}
			tables.Voters = append(tables.Voters, row)
		}
	}()
	// export table.Delegations
	// func() {
	// 	itr := sdk.KVStorePrefixIterator(store, delegationSubstore)
	// 	defer itr.Close()
	// 	for ; itr.Valid(); itr.Next() {
	// 		k := itr.Key()
	// 		meDelegator := string(k[1:])
	// 		strs := strings.Split(meDelegator, types.KeySeparator)
	// 		if len(strs) != 2 {
	// 			panic("failed to split out meDelegator: " + meDelegator)
	// 		}
	// 		voter, delegator := types.AccountKey(strs[0]), types.AccountKey(strs[1])
	// 		val, err := vs.GetDelegation(ctx, voter, delegator)
	// 		if err != nil {
	// 			panic("failed to read delegation: " + err.Error())
	// 		}
	// 		row := DelegationRow{
	// 			Voter:      voter,
	// 			Delegator:  delegator,
	// 			Delegation: *val,
	// 		}
	// 		tables.Delegations = append(tables.Delegations, row)
	// 	}
	// }()

	// list, err := vs.GetReferenceList(ctx)
	// if err != nil {
	// 	panic("failed to get Reference List: " + err.Error())
	// }
	// tables.ReferenceList = ReferenceListTable{
	// 	List: *list,
	// }
	return tables
}

// Import - Import voter state
func (vs VoteStorage) Import(ctx sdk.Context, ir *VoterTablesIR) {
	check := func(e error) {
		if e != nil {
			panic("[vote] Failed to import: " + e.Error())
		}
	}
	// import table.Voters
	for _, v := range ir.Voters {
		err := vs.SetVoter(ctx, v.Username, &v.Voter)
		check(err)
	}
	// import table.Delegations
	// for _, v := range ir.Delegations {
	// 	err := vs.SetDelegation(ctx, v.Voter, v.Delegator, &v.Delegation)
	// 	check(err)
	// }
	// // import table.ReferenceList
	// err := vs.SetReferenceList(ctx, &ir.ReferenceList.List)
	// check(err)
}

// GetVoterKey - "voter substore" + "voter"
func GetVoterKey(me types.AccountKey) []byte {
	return append(voterSubstore, me...)
}
