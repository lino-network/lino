package model

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var (
	DelegatorSubstore    = []byte{0x00}
	VoterSubstore        = []byte{0x01}
	ProposalSubstore     = []byte{0x02}
	VoteSubstore         = []byte{0x03}
	ProposalListSubStore = []byte("ProposalList/ProposalListKey")
)

type VoteStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewValidatorManager returns a new ValidatorManager
func NewVoteStorage(key sdk.StoreKey) *VoteStorage {
	cdc := wire.NewCodec()
	vm := &VoteStorage{
		key: key,
		cdc: cdc,
	}

	return vm
}

func (vm VoteStorage) InitGenesis(ctx sdk.Context) error {
	lst := &ProposalList{}

	if err := vm.SetProposalList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm VoteStorage) GetVoter(ctx sdk.Context, accKey types.AccountKey) (*Voter, sdk.Error) {
	store := ctx.KVStore(vm.key)
	voterByte := store.Get(GetVoterKey(accKey))
	if voterByte == nil {
		return nil, ErrGetVoter()
	}
	voter := new(Voter)
	if err := vm.cdc.UnmarshalJSON(voterByte, voter); err != nil {
		return nil, ErrVoterUnmarshalError(err)
	}
	return voter, nil
}

func (vm VoteStorage) SetVoter(ctx sdk.Context, accKey types.AccountKey, voter *Voter) sdk.Error {
	store := ctx.KVStore(vm.key)
	voterByte, err := vm.cdc.MarshalJSON(*voter)
	if err != nil {
		return ErrVoterMarshalError(err)
	}
	store.Set(GetVoterKey(accKey), voterByte)
	return nil
}

func (vm VoteStorage) DeleteVoter(ctx sdk.Context, username types.AccountKey) sdk.Error {
	store := ctx.KVStore(vm.key)
	store.Delete(GetVoterKey(username))
	return nil
}

func (vm VoteStorage) GetVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey) (*Vote, sdk.Error) {
	store := ctx.KVStore(vm.key)
	voteByte := store.Get(GetVoteKey(proposalID, voter))
	if voteByte == nil {
		return nil, ErrGetVote()
	}
	vote := new(Vote)
	if err := vm.cdc.UnmarshalJSON(voteByte, vote); err != nil {
		return nil, ErrVoteUnmarshalError(err)
	}
	return vote, nil
}

func (vm VoteStorage) SetVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey, vote *Vote) sdk.Error {
	store := ctx.KVStore(vm.key)
	voteByte, err := vm.cdc.MarshalJSON(*vote)
	if err != nil {
		return ErrVoteMarshalError(err)
	}
	store.Set(GetVoteKey(proposalID, voter), voteByte)
	return nil
}

func (vm VoteStorage) DeleteVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey) sdk.Error {
	store := ctx.KVStore(vm.key)
	store.Delete(GetVoteKey(proposalID, voter))
	return nil
}

func (vm VoteStorage) GetProposalList(ctx sdk.Context) (*ProposalList, sdk.Error) {
	store := ctx.KVStore(vm.key)
	lstByte := store.Get(GetProposalListKey())
	if lstByte == nil {
		return nil, ErrGetProposal()
	}
	lst := new(ProposalList)
	if err := vm.cdc.UnmarshalJSON(lstByte, lst); err != nil {
		return nil, ErrProposalUnmarshalError(err)
	}
	return lst, nil
}

func (vm VoteStorage) SetProposalList(ctx sdk.Context, lst *ProposalList) sdk.Error {
	store := ctx.KVStore(vm.key)
	lstByte, err := vm.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrProposalMarshalError(err)
	}
	store.Set(GetProposalListKey(), lstByte)
	return nil
}

// onle support change parameter proposal now
func (vm VoteStorage) GetProposal(ctx sdk.Context, proposalID types.ProposalKey) (*ChangeParameterProposal, sdk.Error) {
	store := ctx.KVStore(vm.key)
	proposalByte := store.Get(GetProposalKey(proposalID))
	if proposalByte == nil {
		return nil, ErrGetProposal()
	}
	proposal := new(ChangeParameterProposal)
	if err := vm.cdc.UnmarshalJSON(proposalByte, proposal); err != nil {
		return nil, ErrProposalUnmarshalError(err)
	}
	return proposal, nil
}

// onle support change parameter proposal now
func (vm VoteStorage) SetProposal(ctx sdk.Context, proposalID types.ProposalKey, proposal *ChangeParameterProposal) sdk.Error {
	store := ctx.KVStore(vm.key)
	proposalByte, err := vm.cdc.MarshalJSON(*proposal)
	if err != nil {
		return ErrProposalMarshalError(err)
	}
	store.Set(GetProposalKey(proposalID), proposalByte)
	return nil
}

func (vm VoteStorage) DeleteProposal(ctx sdk.Context, proposalID types.ProposalKey) sdk.Error {
	store := ctx.KVStore(vm.key)
	store.Delete(GetProposalKey(proposalID))
	return nil
}

func (vm VoteStorage) GetDelegation(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) (*Delegation, sdk.Error) {
	store := ctx.KVStore(vm.key)
	delegationByte := store.Get(GetDelegationKey(voter, delegator))
	if delegationByte == nil {
		return nil, ErrGetDelegation()
	}
	delegation := new(Delegation)
	if err := vm.cdc.UnmarshalJSON(delegationByte, delegation); err != nil {
		return nil, ErrDelegationUnmarshalError(err)
	}
	return delegation, nil
}

func (vm VoteStorage) SetDelegation(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey, delegation *Delegation) sdk.Error {
	store := ctx.KVStore(vm.key)
	delegationByte, err := vm.cdc.MarshalJSON(*delegation)
	if err != nil {
		return ErrDelegationMarshalError(err)
	}
	store.Set(GetDelegationKey(voter, delegator), delegationByte)
	return nil
}

func (vm VoteStorage) DeleteDelegation(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) sdk.Error {
	store := ctx.KVStore(vm.key)
	store.Delete(GetDelegationKey(voter, delegator))
	return nil
}

func (vm VoteStorage) GetAllDelegators(ctx sdk.Context, voterName types.AccountKey) ([]types.AccountKey, sdk.Error) {
	store := ctx.KVStore(vm.key)
	iterator := store.Iterator(subspace(GetDelegatorPrefix(voterName)))

	var delegators []types.AccountKey

	for ; iterator.Valid(); iterator.Next() {
		delegationBytes := iterator.Value()
		var delegation Delegation
		err := vm.cdc.UnmarshalJSON(delegationBytes, &delegation)
		if err != nil {
			return nil, ErrDelegationUnmarshalError(err)
		}
		delegators = append(delegators, delegation.Delegator)
	}
	iterator.Close()
	return delegators, nil
}

func (vm VoteStorage) GetAllVotes(ctx sdk.Context, proposalID types.ProposalKey) ([]Vote, sdk.Error) {
	store := ctx.KVStore(vm.key)
	iterator := store.Iterator(subspace(GetVotePrefix(proposalID)))

	var votes []Vote

	for ; iterator.Valid(); iterator.Next() {
		voteBytes := iterator.Value()
		var vote Vote
		err := vm.cdc.UnmarshalJSON(voteBytes, &vote)
		if err != nil {
			return nil, ErrVoteUnmarshalError(err)
		}
		votes = append(votes, vote)
	}
	iterator.Close()
	return votes, nil
}

func (vm VoteStorage) GetNextProposalID() (types.ProposalKey, sdk.Error) {
	nextProposalID += 1
	return types.ProposalKey(strconv.FormatInt(nextProposalID, 10)), nil
}

func GetDelegatorPrefix(me types.AccountKey) []byte {
	return append(append(DelegatorSubstore, me...), types.KeySeparator...)
}

// "delegator substore" + "me(voter)" + "my delegator"
func GetDelegationKey(me types.AccountKey, myDelegator types.AccountKey) []byte {
	return append(GetDelegatorPrefix(me), myDelegator...)
}

func GetVotePrefix(id types.ProposalKey) []byte {
	return append(append(VoteSubstore, id...), types.KeySeparator...)
}

// "vote substore" + "proposalID" + "voter"
func GetVoteKey(proposalID types.ProposalKey, voter types.AccountKey) []byte {
	return append(GetVotePrefix(proposalID), voter...)
}

func GetProposalKey(proposalID types.ProposalKey) []byte {
	return append(ProposalSubstore, proposalID...)
}

func GetProposalListKey() []byte {
	return ProposalListSubStore
}

func GetVoterKey(me types.AccountKey) []byte {
	return append(VoterSubstore, me...)
}

func subspace(prefix []byte) (start, end []byte) {
	end = make([]byte, len(prefix))
	copy(end, prefix)
	end[len(end)-1]++
	return prefix, end
}
