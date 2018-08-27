package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
)

var (
	nextProposalIDSubstore  = []byte{0x00}
	ongoingProposalSubStore = []byte{0x01}
	expiredProposalSubStore = []byte{0x02}
)

type ProposalStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewProposalStorage(key sdk.StoreKey) ProposalStorage {
	cdc := wire.NewCodec()

	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&ChangeParamProposal{}, "changeParam", nil)
	cdc.RegisterConcrete(&ProtocolUpgradeProposal{}, "upgrade", nil)
	cdc.RegisterConcrete(&ContentCensorshipProposal{}, "censorship", nil)

	cdc.RegisterInterface((*param.Parameter)(nil), nil)
	cdc.RegisterConcrete(param.GlobalAllocationParam{}, "allocation", nil)
	cdc.RegisterConcrete(param.InfraInternalAllocationParam{}, "infraAllocation", nil)
	cdc.RegisterConcrete(param.EvaluateOfContentValueParam{}, "contentValue", nil)
	cdc.RegisterConcrete(param.VoteParam{}, "voteParam", nil)
	cdc.RegisterConcrete(param.ProposalParam{}, "proposalParam", nil)
	cdc.RegisterConcrete(param.DeveloperParam{}, "developerParam", nil)
	cdc.RegisterConcrete(param.ValidatorParam{}, "validatorParam", nil)
	cdc.RegisterConcrete(param.CoinDayParam{}, "coinDayParam", nil)
	cdc.RegisterConcrete(param.BandwidthParam{}, "bandwidthParam", nil)
	cdc.RegisterConcrete(param.AccountParam{}, "accountParam", nil)
	cdc.RegisterConcrete(param.PostParam{}, "postParam", nil)

	wire.RegisterCrypto(cdc)
	vs := ProposalStorage{
		key: key,
		cdc: cdc,
	}
	return vs
}

func (ps ProposalStorage) InitGenesis(ctx sdk.Context) sdk.Error {
	nextProposalID := &NextProposalID{
		NextProposalID: 1,
	}
	if err := ps.SetNextProposalID(ctx, nextProposalID); err != nil {
		return err
	}
	return nil
}

func (ps ProposalStorage) DoesProposalExist(ctx sdk.Context, proposalID types.ProposalKey) bool {
	store := ctx.KVStore(ps.key)
	return store.Has(GetOngoingProposalKey(proposalID)) || store.Has(GetExpiredProposalKey(proposalID))
}

func (ps ProposalStorage) GetOngoingProposal(ctx sdk.Context, proposalID types.ProposalKey) (Proposal, sdk.Error) {
	store := ctx.KVStore(ps.key)
	proposalByte := store.Get(GetOngoingProposalKey(proposalID))
	if proposalByte == nil {
		return nil, ErrProposalNotFound()
	}
	proposal := new(Proposal)
	if err := ps.cdc.UnmarshalJSON(proposalByte, proposal); err != nil {
		return nil, ErrFailedToUnmarshalProposal(err)
	}
	return *proposal, nil
}

func (ps ProposalStorage) SetOngoingProposal(ctx sdk.Context, proposalID types.ProposalKey, proposal Proposal) sdk.Error {
	store := ctx.KVStore(ps.key)
	proposalByte, err := ps.cdc.MarshalJSON(proposal)
	if err != nil {
		return ErrFailedToMarshalProposal(err)
	}
	store.Set(GetOngoingProposalKey(proposalID), proposalByte)
	return nil
}

func (ps ProposalStorage) DeleteOngoingProposal(ctx sdk.Context, proposalID types.ProposalKey) sdk.Error {
	store := ctx.KVStore(ps.key)
	store.Delete(GetOngoingProposalKey(proposalID))
	return nil
}

func (ps ProposalStorage) GetExpiredProposal(ctx sdk.Context, proposalID types.ProposalKey) (Proposal, sdk.Error) {
	store := ctx.KVStore(ps.key)
	proposalByte := store.Get(GetExpiredProposalKey(proposalID))
	if proposalByte == nil {
		return nil, ErrProposalNotFound()
	}
	proposal := new(Proposal)
	if err := ps.cdc.UnmarshalJSON(proposalByte, proposal); err != nil {
		return nil, ErrFailedToUnmarshalProposal(err)
	}
	return *proposal, nil
}

func (ps ProposalStorage) SetExpiredProposal(ctx sdk.Context, proposalID types.ProposalKey, proposal Proposal) sdk.Error {
	store := ctx.KVStore(ps.key)
	proposalByte, err := ps.cdc.MarshalJSON(proposal)
	if err != nil {
		return ErrFailedToMarshalProposal(err)
	}
	store.Set(GetExpiredProposalKey(proposalID), proposalByte)
	return nil
}

func (ps ProposalStorage) DeleteExpiredProposal(ctx sdk.Context, proposalID types.ProposalKey) sdk.Error {
	store := ctx.KVStore(ps.key)
	store.Delete(GetExpiredProposalKey(proposalID))
	return nil
}

func (ps ProposalStorage) GetOngoingProposalList(ctx sdk.Context) ([]Proposal, sdk.Error) {
	store := ctx.KVStore(ps.key)
	iterator := store.Iterator(subspace(ongoingProposalSubStore))

	var proposalList []Proposal

	for ; iterator.Valid(); iterator.Next() {
		proposalBytes := iterator.Value()
		var p Proposal
		err := ps.cdc.UnmarshalJSON(proposalBytes, &p)
		if err != nil {
			return nil, ErrFailedToUnmarshalProposalList(err)
		}
		proposalList = append(proposalList, p)
	}
	iterator.Close()
	return proposalList, nil
}

func (ps ProposalStorage) GetExpiredProposalList(ctx sdk.Context) ([]Proposal, sdk.Error) {
	store := ctx.KVStore(ps.key)
	iterator := store.Iterator(subspace(expiredProposalSubStore))

	var proposalList []Proposal

	for ; iterator.Valid(); iterator.Next() {
		proposalBytes := iterator.Value()
		var p Proposal
		err := ps.cdc.UnmarshalJSON(proposalBytes, &p)
		if err != nil {
			return nil, ErrFailedToUnmarshalProposalList(err)
		}
		proposalList = append(proposalList, p)
	}
	iterator.Close()
	return proposalList, nil
}

func (ps ProposalStorage) GetNextProposalID(ctx sdk.Context) (*NextProposalID, sdk.Error) {
	store := ctx.KVStore(ps.key)
	nextProposalIDByte := store.Get(getNextProposalIDKey())
	if nextProposalIDByte == nil {
		return nil, ErrNextProposalIDNotFound()
	}
	nextProposalID := new(NextProposalID)
	if err := ps.cdc.UnmarshalJSON(nextProposalIDByte, nextProposalID); err != nil {
		return nil, ErrFailedToUnmarshalNextProposalID(err)
	}
	return nextProposalID, nil
}

func (ps ProposalStorage) SetNextProposalID(ctx sdk.Context, nextProposalID *NextProposalID) sdk.Error {
	store := ctx.KVStore(ps.key)
	nextProposalIDByte, err := ps.cdc.MarshalJSON(*nextProposalID)
	if err != nil {
		return ErrFailedToMarshalNextProposalID(err)
	}
	store.Set(getNextProposalIDKey(), nextProposalIDByte)
	return nil
}

func GetOngoingProposalKey(proposalID types.ProposalKey) []byte {
	return append(ongoingProposalSubStore, proposalID...)
}

func GetExpiredProposalKey(proposalID types.ProposalKey) []byte {
	return append(expiredProposalSubStore, proposalID...)
}

func getNextProposalIDKey() []byte {
	return nextProposalIDSubstore
}

func subspace(prefix []byte) (start, end []byte) {
	end = make([]byte, len(prefix))
	copy(end, prefix)
	end[len(end)-1]++
	return prefix, end
}
