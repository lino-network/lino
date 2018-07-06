package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
)

var (
	proposalSubstore       = []byte{0x00}
	proposalListSubStore   = []byte{0x01}
	nextProposalIDSubstore = []byte{0x02}
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
	proposalLst := &ProposalList{}
	if err := ps.SetProposalList(ctx, proposalLst); err != nil {
		return err
	}

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
	return store.Has(GetProposalKey(proposalID))
}

func (ps ProposalStorage) GetProposalList(ctx sdk.Context) (*ProposalList, sdk.Error) {
	store := ctx.KVStore(ps.key)
	lstByte := store.Get(GetProposalListKey())
	if lstByte == nil {
		return nil, ErrGetProposal()
	}
	lst := new(ProposalList)
	if err := ps.cdc.UnmarshalJSON(lstByte, lst); err != nil {
		return nil, ErrProposalUnmarshalError(err)
	}
	return lst, nil
}

func (ps ProposalStorage) SetProposalList(ctx sdk.Context, lst *ProposalList) sdk.Error {
	store := ctx.KVStore(ps.key)
	lstByte, err := ps.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrProposalMarshalError(err)
	}
	store.Set(GetProposalListKey(), lstByte)
	return nil
}

// only support change parameter proposal now
func (ps ProposalStorage) GetProposal(ctx sdk.Context, proposalID types.ProposalKey) (Proposal, sdk.Error) {
	store := ctx.KVStore(ps.key)
	proposalByte := store.Get(GetProposalKey(proposalID))
	if proposalByte == nil {
		return nil, ErrGetProposal()
	}
	proposal := new(Proposal)
	if err := ps.cdc.UnmarshalJSON(proposalByte, proposal); err != nil {
		return nil, ErrProposalUnmarshalError(err)
	}
	return *proposal, nil
}

// only support change parameter proposal now
func (ps ProposalStorage) SetProposal(ctx sdk.Context, proposalID types.ProposalKey, proposal Proposal) sdk.Error {
	store := ctx.KVStore(ps.key)
	proposalByte, err := ps.cdc.MarshalJSON(proposal)
	if err != nil {
		return ErrProposalMarshalError(err)
	}
	store.Set(GetProposalKey(proposalID), proposalByte)
	return nil
}

func (ps ProposalStorage) GetNextProposalID(ctx sdk.Context) (*NextProposalID, sdk.Error) {
	store := ctx.KVStore(ps.key)
	nextProposalIDByte := store.Get(GetNextProposalIDKey())
	if nextProposalIDByte == nil {
		return nil, ErrGetNextProposalID()
	}
	nextProposalID := new(NextProposalID)
	if err := ps.cdc.UnmarshalJSON(nextProposalIDByte, nextProposalID); err != nil {
		return nil, ErrNextProposalIDUnmarshalError(err)
	}
	return nextProposalID, nil
}

func (ps ProposalStorage) SetNextProposalID(ctx sdk.Context, nextProposalID *NextProposalID) sdk.Error {
	store := ctx.KVStore(ps.key)
	nextProposalIDByte, err := ps.cdc.MarshalJSON(*nextProposalID)
	if err != nil {
		return ErrNextProposalIDMarshalError(err)
	}
	store.Set(GetNextProposalIDKey(), nextProposalIDByte)
	return nil
}

func (ps ProposalStorage) DeleteProposal(ctx sdk.Context, proposalID types.ProposalKey) sdk.Error {
	store := ctx.KVStore(ps.key)
	store.Delete(GetProposalKey(proposalID))
	return nil
}

func GetProposalKey(proposalID types.ProposalKey) []byte {
	return append(proposalSubstore, proposalID...)
}

func GetProposalListKey() []byte {
	return proposalListSubStore
}

func GetNextProposalIDKey() []byte {
	return nextProposalIDSubstore
}
