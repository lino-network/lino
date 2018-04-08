package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
	oldwire "github.com/tendermint/go-wire"
)

var (
	DelegatorSubstore    = []byte{0x00}
	VoterSubstore        = []byte{0x01}
	ProposalSubstore     = []byte{0x02}
	VoteSubstore         = []byte{0x03}
	ProposalListSubStore = []byte("ProposalList/ProposalListKey")
)

const returnCoinEvent = 0x1
const decideProposalEvent = 0x2

var _ = oldwire.RegisterInterface(
	struct{ types.Event }{},
	oldwire.ConcreteType{ReturnCoinEvent{}, returnCoinEvent},
	oldwire.ConcreteType{DecideProposalEvent{}, decideProposalEvent},
)

// vote manager is the proxy for all storage structs defined above
type VoteManager struct {
	storage *model.VoteStorage `json:"vote_storage"`
}

// create NewVoteManager
func NewVoteManager(key sdk.StoreKey) *VoteManager {
	return &VoteManager{
		storage: model.NewVoteStorage(key),
	}
}

func (vm VoteManager) InitGenesis(ctx sdk.Context) error {
	lst := &model.ProposalList{}

	if err := vm.storage.SetProposalList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) IsVoterExist(ctx sdk.Context, accKey types.AccountKey) bool {
	voterByte, _ := vm.storage.GetVoter(ctx, accKey)
	return voterByte != nil
}

func (vm VoteManager) IsProposalExist(ctx sdk.Context, proposalID types.ProposalKey) bool {
	proposalByte, _ := vm.storage.GetProposal(ctx, proposalID)
	return proposalByte != nil
}

func (vm VoteManager) IsDelegationExist(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) bool {
	delegationByte, _ := vm.storage.GetDelegation(ctx, voter, delegator)
	return delegationByte != nil
}

func (vm VoteManager) IsLegalWithdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin) bool {
	voter, getErr := vm.storage.GetVoter(ctx, username)
	if getErr != nil {
		return false
	}
	// reject if withdraw is less than minimum withdraw
	if !coin.IsGTE(types.VoterMinimumWithdraw) {
		return false
	}
	//reject if the remaining coins are less than register fee
	res := voter.Deposit.Minus(coin)
	return res.IsGTE(types.VoterRegisterFee)
}

// onle support change parameter proposal now
func (vm VoteManager) AddProposal(ctx sdk.Context, creator types.AccountKey, des *model.ChangeParameterDescription) sdk.Error {
	newID, getErr := vm.storage.GetNextProposalID()
	if getErr != nil {
		return getErr
	}

	proposal := model.Proposal{
		Creator:      creator,
		ProposalID:   newID,
		AgreeVote:    types.Coin{Amount: 0},
		DisagreeVote: types.Coin{Amount: 0},
	}

	changeParameterProposal := &model.ChangeParameterProposal{
		Proposal:                   proposal,
		ChangeParameterDescription: *des,
	}
	if err := vm.storage.SetProposal(ctx, newID, changeParameterProposal); err != nil {
		return err
	}

	lst, getErr := vm.storage.GetProposalList(ctx)
	if getErr != nil {
		return getErr
	}
	lst.OngoingProposal = append(lst.OngoingProposal, newID)
	if err := vm.storage.SetProposalList(ctx, lst); err != nil {
		return err
	}

	return nil
}

// onle support change parameter proposal now
func (vm VoteManager) AddVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey, res bool) sdk.Error {
	vote := model.Vote{
		Voter:  voter,
		Result: res,
	}
	// will overwrite the old vote
	if err := vm.storage.SetVote(ctx, proposalID, voter, &vote); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) AddDelegation(ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, coin types.Coin) sdk.Error {
	var delegation *model.Delegation
	var getErr sdk.Error

	if !vm.IsDelegationExist(ctx, voterName, delegatorName) {
		delegation = &model.Delegation{
			Delegator: delegatorName,
		}
	} else {
		delegation, getErr = vm.storage.GetDelegation(ctx, voterName, delegatorName)
		if getErr != nil {
			return getErr
		}
	}

	voter, getErr := vm.storage.GetVoter(ctx, voterName)
	if getErr != nil {
		return getErr
	}

	voter.DelegatedPower = voter.DelegatedPower.Plus(coin)
	delegation.Amount = delegation.Amount.Plus(coin)

	if err := vm.storage.SetDelegation(ctx, voterName, delegatorName, delegation); err != nil {
		return err
	}
	if err := vm.storage.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) AddVoter(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	voter := &model.Voter{
		Username: username,
		Deposit:  coin,
	}
	// check minimum requirements for registering as a voter
	if !coin.IsGTE(types.VoterRegisterFee) {
		return ErrRegisterFeeNotEnough()
	}

	if setErr := vm.storage.SetVoter(ctx, username, voter); setErr != nil {
		return setErr
	}
	return nil
}

func (vm VoteManager) Deposit(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	voter.Deposit = voter.Deposit.Plus(coin)
	if setErr := vm.storage.SetVoter(ctx, username, voter); setErr != nil {
		return setErr
	}
	return nil
}

// this method won't check if it is a legal withdraw, caller should check by itself
func (vm VoteManager) Withdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin, gm global.GlobalManager) sdk.Error {
	voter, getErr := vm.storage.GetVoter(ctx, username)
	if getErr != nil {
		return getErr
	}
	voter.Deposit = voter.Deposit.Minus(coin)

	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	if err := vm.CreateReturnCoinEvent(ctx, username, coin, gm); err != nil {
		return nil
	}
	return nil
}

func (vm VoteManager) WithdrawAll(ctx sdk.Context, username types.AccountKey, gm global.GlobalManager) sdk.Error {
	voter, getErr := vm.storage.GetVoter(ctx, username)
	if getErr != nil {
		return getErr
	}
	if err := vm.Withdraw(ctx, username, voter.Deposit, gm); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) ReturnCoinToDelegator(ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, gm global.GlobalManager) sdk.Error {
	voter, getErr := vm.storage.GetVoter(ctx, voterName)
	if getErr != nil {
		return getErr
	}
	delegation, getErr := vm.storage.GetDelegation(ctx, voterName, delegatorName)
	if getErr != nil {
		return getErr
	}

	voter.DelegatedPower = voter.DelegatedPower.Minus(delegation.Amount)
	if err := vm.CreateReturnCoinEvent(ctx, delegatorName, delegation.Amount, gm); err != nil {
		return err
	}
	if err := vm.storage.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}
	if err := vm.storage.DeleteDelegation(ctx, voterName, delegatorName); err != nil {
		return err
	}
	return nil
}

// return coin to an user periodically
func (vm VoteManager) CreateReturnCoinEvent(ctx sdk.Context, username types.AccountKey, amount types.Coin, gm global.GlobalManager) sdk.Error {
	pieceRat := amount.ToRat().Quo(sdk.NewRat(types.CoinReturnTimes))
	piece := types.RatToCoin(pieceRat)
	event := ReturnCoinEvent{
		Username: username,
		Amount:   piece,
	}

	if err := gm.RegisterCoinReturnEvent(ctx, event); err != nil {
		return err
	}

	return nil
}

// decide the proposal in 7 days
func (vm VoteManager) CreateDecideProposalEvent(ctx sdk.Context, gm global.GlobalManager) sdk.Error {
	event := DecideProposalEvent{}
	if err := gm.RegisterEventAtTime(ctx, ctx.BlockHeader().Time+(types.ProposalDecideHr*3600), event); err != nil {
		return err
	}
	return nil
}

func (vm VoteManager) GetVotingPower(ctx sdk.Context, voterName types.AccountKey) (types.Coin, sdk.Error) {
	voter, getErr := vm.storage.GetVoter(ctx, voterName)
	if getErr != nil {
		return types.Coin{}, getErr
	}
	res := voter.Deposit.Plus(voter.DelegatedPower)
	return res, nil
}

func (vm VoteManager) GetAllDelegators(ctx sdk.Context, voterName types.AccountKey) ([]types.AccountKey, sdk.Error) {
	return vm.storage.GetAllDelegators(ctx, voterName)
}

func (vm VoteManager) DeleteVoter(ctx sdk.Context, username types.AccountKey) sdk.Error {
	return vm.storage.DeleteVoter(ctx, username)
}

func (vm VoteManager) DeleteDelegation(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) sdk.Error {
	return vm.storage.DeleteDelegation(ctx, voter, delegator)
}
