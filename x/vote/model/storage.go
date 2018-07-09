package model

import (
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

var (
	delegationSubstore    = []byte{0x00}
	voterSubstore         = []byte{0x01}
	voteSubstore          = []byte{0x02}
	referenceListSubStore = []byte{0x03}
	delegateeSubStore     = []byte{0x04}
)

type VoteStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewVoteStorage(key sdk.StoreKey) VoteStorage {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	vs := VoteStorage{
		key: key,
		cdc: cdc,
	}

	return vs
}

func (vs VoteStorage) InitGenesis(ctx sdk.Context) sdk.Error {
	lst := &ReferenceList{}
	if err := vs.SetReferenceList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vs VoteStorage) DoesVoterExist(ctx sdk.Context, accKey types.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetVoterKey(accKey))
}

func (vs VoteStorage) DoesVoteExist(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetVoteKey(proposalID, voter))
}

func (vs VoteStorage) DoesDelegationExist(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) bool {
	store := ctx.KVStore(vs.key)
	return store.Has(GetDelegationKey(voter, delegator))
}

func (vs VoteStorage) GetVoter(ctx sdk.Context, accKey types.AccountKey) (*Voter, sdk.Error) {
	store := ctx.KVStore(vs.key)
	voterByte := store.Get(GetVoterKey(accKey))
	if voterByte == nil {
		return nil, ErrVoterNotFound()
	}
	voter := new(Voter)
	if err := vs.cdc.UnmarshalJSON(voterByte, voter); err != nil {
		return nil, ErrFailedToUnmarshalVoter(err)
	}
	return voter, nil
}

func (vs VoteStorage) SetVoter(ctx sdk.Context, accKey types.AccountKey, voter *Voter) sdk.Error {
	store := ctx.KVStore(vs.key)
	voterByte, err := vs.cdc.MarshalJSON(*voter)
	if err != nil {
		return ErrFailedToMarshalVoter(err)
	}
	store.Set(GetVoterKey(accKey), voterByte)
	return nil
}

func (vs VoteStorage) DeleteVoter(ctx sdk.Context, username types.AccountKey) sdk.Error {
	store := ctx.KVStore(vs.key)
	store.Delete(GetVoterKey(username))
	return nil
}

func (vs VoteStorage) GetVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey) (*Vote, sdk.Error) {
	store := ctx.KVStore(vs.key)
	voteByte := store.Get(GetVoteKey(proposalID, voter))
	if voteByte == nil {
		return nil, ErrVoteNotFound()
	}
	vote := new(Vote)
	if err := vs.cdc.UnmarshalJSON(voteByte, vote); err != nil {
		return nil, ErrFailedToUnmarshalVote(err)
	}
	return vote, nil
}

func (vs VoteStorage) SetVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey, vote *Vote) sdk.Error {
	store := ctx.KVStore(vs.key)
	voteByte, err := vs.cdc.MarshalJSON(*vote)
	if err != nil {
		return ErrFailedToMarshalVote(err)
	}
	store.Set(GetVoteKey(proposalID, voter), voteByte)
	return nil
}

func (vs VoteStorage) DeleteVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey) sdk.Error {
	store := ctx.KVStore(vs.key)
	store.Delete(GetVoteKey(proposalID, voter))
	return nil
}

func (vs VoteStorage) GetDelegation(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) (*Delegation, sdk.Error) {
	store := ctx.KVStore(vs.key)
	delegationByte := store.Get(GetDelegationKey(voter, delegator))
	if delegationByte == nil {
		return nil, ErrDelegationNotFound()
	}
	delegation := new(Delegation)
	if err := vs.cdc.UnmarshalJSON(delegationByte, delegation); err != nil {
		return nil, ErrFailedToUnmarshalDelegation(err)
	}
	return delegation, nil
}

func (vs VoteStorage) SetDelegation(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey, delegation *Delegation) sdk.Error {
	store := ctx.KVStore(vs.key)
	delegationByte, err := vs.cdc.MarshalJSON(*delegation)
	if err != nil {
		return ErrFailedToMarshalDelegation(err)
	}
	store.Set(GetDelegationKey(voter, delegator), delegationByte)
	store.Set(getDelegateeKey(delegator, voter), delegationByte)
	return nil
}

func (vs VoteStorage) DeleteDelegation(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) sdk.Error {
	store := ctx.KVStore(vs.key)
	store.Delete(GetDelegationKey(voter, delegator))
	store.Delete(getDelegateeKey(delegator, voter))
	return nil
}

func (vs VoteStorage) GetAllDelegators(ctx sdk.Context, voterName types.AccountKey) ([]types.AccountKey, sdk.Error) {
	store := ctx.KVStore(vs.key)
	iterator := store.Iterator(subspace(getDelegationPrefix(voterName)))

	var delegators []types.AccountKey

	for ; iterator.Valid(); iterator.Next() {
		delegationBytes := iterator.Value()
		var delegation Delegation
		err := vs.cdc.UnmarshalJSON(delegationBytes, &delegation)
		if err != nil {
			return nil, ErrFailedToUnmarshalDelegation(err)
		}
		delegators = append(delegators, delegation.Delegator)
	}
	iterator.Close()
	return delegators, nil
}

func (vs VoteStorage) GetAllVotes(ctx sdk.Context, proposalID types.ProposalKey) ([]Vote, sdk.Error) {
	store := ctx.KVStore(vs.key)
	iterator := store.Iterator(subspace(getVotePrefix(proposalID)))

	var votes []Vote

	for ; iterator.Valid(); iterator.Next() {
		voteBytes := iterator.Value()
		var vote Vote
		err := vs.cdc.UnmarshalJSON(voteBytes, &vote)
		if err != nil {
			return nil, ErrFailedToUnmarshalVote(err)
		}
		votes = append(votes, vote)
	}
	iterator.Close()
	return votes, nil
}

func (vs VoteStorage) GetReferenceList(ctx sdk.Context) (*ReferenceList, sdk.Error) {
	store := ctx.KVStore(vs.key)
	lstByte := store.Get(getReferenceListKey())
	if lstByte == nil {
		return nil, ErrReferenceListNotFound()
	}
	lst := new(ReferenceList)
	if err := vs.cdc.UnmarshalJSON(lstByte, lst); err != nil {
		return nil, ErrFailedToUnmarshalReferenceList(err)
	}
	return lst, nil
}

func (vs VoteStorage) SetReferenceList(ctx sdk.Context, lst *ReferenceList) sdk.Error {
	store := ctx.KVStore(vs.key)
	lstByte, err := vs.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrFailedToMarshalReferenceList(err)
	}
	store.Set(getReferenceListKey(), lstByte)
	return nil
}

func getDelegationPrefix(me types.AccountKey) []byte {
	return append(append(delegationSubstore, me...), types.KeySeparator...)
}

// "delegation substore" + "me(voter)" + "my delegator"
func GetDelegationKey(me types.AccountKey, myDelegator types.AccountKey) []byte {
	return append(getDelegationPrefix(me), myDelegator...)
}

func getVotePrefix(id types.ProposalKey) []byte {
	return append(append(voteSubstore, id...), types.KeySeparator...)
}

// "vote substore" + "proposalID" + "voter"
func GetVoteKey(proposalID types.ProposalKey, voter types.AccountKey) []byte {
	return append(getVotePrefix(proposalID), voter...)
}

func GetVoterKey(me types.AccountKey) []byte {
	return append(voterSubstore, me...)
}

func getReferenceListKey() []byte {
	return referenceListSubStore
}

func getDelegateePrefix(me types.AccountKey) []byte {
	return append(append(delegateeSubStore, me...), types.KeySeparator...)
}

func getDelegateeKey(me, delegatee types.AccountKey) []byte {
	return append(getDelegateePrefix(me), delegatee...)
}

func subspace(prefix []byte) (start, end []byte) {
	end = make([]byte, len(prefix))
	copy(end, prefix)
	end[len(end)-1]++
	return prefix, end
}
