package vote

import (
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VoteManager - vote manager
type VoteManager struct {
	storage     model.VoteStorage
	paramHolder param.ParamHolder
}

func NewVoteManager(key sdk.StoreKey, holder param.ParamHolder) VoteManager {
	return VoteManager{
		storage:     model.NewVoteStorage(key),
		paramHolder: holder,
	}
}

// InitGenesis - initialize KV Store
func (vm VoteManager) InitGenesis(ctx sdk.Context) error {
	if err := vm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

// DoesVoterExist - check if voter exist or not
func (vm VoteManager) DoesVoterExist(ctx sdk.Context, accKey types.AccountKey) bool {
	return vm.storage.DoesVoterExist(ctx, accKey)
}

// DoesVoterExist - check if vote exist or not
func (vm VoteManager) DoesVoteExist(ctx sdk.Context, proposalID types.ProposalKey, accKey types.AccountKey) bool {
	return vm.storage.DoesVoteExist(ctx, proposalID, accKey)
}

// DoesVoterExist - check if delegation exist or not
func (vm VoteManager) DoesDelegationExist(ctx sdk.Context, voter types.AccountKey, delegator types.AccountKey) bool {
	return vm.storage.DoesDelegationExist(ctx, voter, delegator)
}

// DoesVoterExist - check if a voter is validator or not
func (vm VoteManager) IsInValidatorList(ctx sdk.Context, username types.AccountKey) bool {
	lst, err := vm.storage.GetReferenceList(ctx)
	if err != nil {
		return false
	}
	for _, validator := range lst.AllValidators {
		if validator == username {
			return true
		}
	}
	return false
}

// IsLegalVoterWithdraw - check withdraw voter is not validator
func (vm VoteManager) IsLegalVoterWithdraw(
	ctx sdk.Context, username types.AccountKey, coin types.Coin) bool {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return false
	}
	// reject if this is a validator
	if vm.IsInValidatorList(ctx, username) {
		return false
	}

	availableStakes := voter.LinoStake.Minus(voter.DelegateToOthers)
	//reject if the remaining coins are not enough
	return availableStakes.IsGTE(coin) && coin.IsPositive()
}

// IsLegalDelegatorWithdraw - check if delegator withdraw is valid or not
func (vm VoteManager) IsLegalDelegatorWithdraw(
	ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, coin types.Coin) bool {
	delegation, err := vm.storage.GetDelegation(ctx, voterName, delegatorName)
	if err != nil {
		return false
	}

	//reject if the remaining delegation are less than zero
	return delegation.Amount.IsGTE(coin) && coin.IsPositive()
}

// CanBecomeValidator - check if vote deposit meet requirement or not
func (vm VoteManager) CanBecomeValidator(ctx sdk.Context, username types.AccountKey) bool {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return false
	}
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return false
	}
	// check minimum voting deposit for validator
	return voter.LinoStake.IsGTE(param.ValidatorMinVotingDeposit)
}

// AddVote - voter vote for a proposal
func (vm VoteManager) AddVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey, res bool) sdk.Error {
	// check if the vote exist
	if vm.DoesVoteExist(ctx, proposalID, voter) {
		return ErrVoteAlreadyExist()
	}

	votingPower, err := vm.GetVotingPower(ctx, voter)
	if err != nil {
		return err
	}

	vote := model.Vote{
		Voter:       voter,
		Result:      res,
		VotingPower: votingPower,
	}

	if err := vm.storage.SetVote(ctx, proposalID, voter, &vote); err != nil {
		return err
	}
	return nil
}

// GetVote - get vote detail based on voter and proposal ID
func (vm VoteManager) GetVote(ctx sdk.Context, proposalID types.ProposalKey, voter types.AccountKey) (*model.Vote, sdk.Error) {
	return vm.storage.GetVote(ctx, proposalID, voter)
}

// AddDelegation - add delegation
func (vm VoteManager) AddDelegation(ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, coin types.Coin) sdk.Error {
	var delegation *model.Delegation
	var err sdk.Error

	if !vm.DoesDelegationExist(ctx, voterName, delegatorName) {
		delegation = &model.Delegation{
			Delegator: delegatorName,
			Amount:    types.NewCoinFromInt64(0),
		}
	} else {
		delegation, err = vm.storage.GetDelegation(ctx, voterName, delegatorName)
		if err != nil {
			return err
		}
	}

	// add voter if not exist
	if !vm.DoesVoterExist(ctx, voterName) {
		if err := vm.AddVoter(ctx, voterName, types.NewCoinFromInt64(0)); err != nil {
			return err
		}
	}

	// add delegatedPower for voter
	voter, err := vm.storage.GetVoter(ctx, voterName)
	if err != nil {
		return err
	}
	voter.DelegatedPower = voter.DelegatedPower.Plus(coin)
	// add delegateToOthers for delegator
	delegator, err := vm.storage.GetVoter(ctx, delegatorName)
	if err != nil {
		return err
	}
	delegator.DelegateToOthers = delegator.DelegateToOthers.Plus(coin)

	delegation.Amount = delegation.Amount.Plus(coin)

	if err := vm.storage.SetDelegation(ctx, voterName, delegatorName, delegation); err != nil {
		return err
	}
	if err := vm.storage.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}
	if err := vm.storage.SetVoter(ctx, delegatorName, delegator); err != nil {
		return err
	}
	return nil
}

// AddVoter - add voter
func (vm VoteManager) AddVoter(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	voter := &model.Voter{
		Username:          username,
		LinoStake:         coin,
		LastPowerChangeAt: ctx.BlockHeader().Time.Unix(),
	}

	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	return nil
}

// AddLinoStake - add lino power
func (vm VoteManager) AddLinoStake(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	voter.LinoStake = voter.LinoStake.Plus(coin)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()
	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	return nil
}

// MinusLinoStake - this method won't check if it is a legal withdraw, caller should check by itself
func (vm VoteManager) MinusLinoStake(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return ErrInvalidCoin()
	}
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	voter.LinoStake = voter.LinoStake.Minus(coin)
	voter.LastPowerChangeAt = ctx.BlockHeader().Time.Unix()
	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	return nil
}

// DelegatorWithdraw - withdraw delegation
func (vm VoteManager) DelegatorWithdraw(
	ctx sdk.Context, voterName types.AccountKey, delegatorName types.AccountKey, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return ErrInvalidCoin()
	}
	// change voter's delegated power
	voter, err := vm.storage.GetVoter(ctx, voterName)
	if err != nil {
		return err
	}
	voter.DelegatedPower = voter.DelegatedPower.Minus(coin)
	if err := vm.storage.SetVoter(ctx, voterName, voter); err != nil {
		return err
	}

	// change delegator's delegateToOthers
	delegator, err := vm.storage.GetVoter(ctx, delegatorName)
	if err != nil {
		return err
	}
	delegator.DelegateToOthers = delegator.DelegateToOthers.Minus(coin)
	if err := vm.storage.SetVoter(ctx, delegatorName, delegator); err != nil {
		return err
	}

	// change this delegation's amount
	delegation, err := vm.storage.GetDelegation(ctx, voterName, delegatorName)
	if err != nil {
		return err
	}
	delegation.Amount = delegation.Amount.Minus(coin)

	if delegation.Amount.IsZero() {
		if err := vm.storage.DeleteDelegation(ctx, voterName, delegatorName); err != nil {
			return err
		}
	} else {
		vm.storage.SetDelegation(ctx, voterName, delegatorName, delegation)
	}

	return nil
}

// ClaimInterest - add lino power interst to user balance
func (vm VoteManager) ClaimInterest(
	ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	claimedInterest := voter.Interest
	voter.Interest = types.NewCoinFromInt64(0)
	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return claimedInterest, nil
}

// AddInterest - add interst
func (vm VoteManager) AddInterest(
	ctx sdk.Context, username types.AccountKey, interest types.Coin) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, username)
	if err != nil {
		return err
	}
	voter.Interest = voter.Interest.Plus(interest)
	if err := vm.storage.SetVoter(ctx, username, voter); err != nil {
		return err
	}
	return nil
}

// GetVotingPower - get voter voting power
func (vm VoteManager) GetVotingPower(ctx sdk.Context, voterName types.AccountKey) (types.Coin, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, voterName)
	if err != nil {
		return types.Coin{}, err
	}
	res := voter.LinoStake.Plus(voter.DelegatedPower).Minus(voter.DelegateToOthers)
	return res, nil
}

// GetPenaltyList - get penalty list if voter is also validator doesn't vote
func (vm VoteManager) GetPenaltyList(
	ctx sdk.Context, proposalID types.ProposalKey, proposalType types.ProposalType,
	oncallValidators []types.AccountKey) (types.PenaltyList, sdk.Error) {
	penaltyList := types.PenaltyList{
		PenaltyList: []types.AccountKey{},
	}

	// get all votes to calculate the voting result
	votes, err := vm.storage.GetAllVotes(ctx, proposalID)
	if err != nil {
		return penaltyList, err
	}

	for _, vote := range votes {
		// remove from list if the validator voted
		for idx, validator := range oncallValidators {
			if validator == vote.Voter {
				oncallValidators = append(oncallValidators[:idx], oncallValidators[idx+1:]...)
				break
			}
		}
		// TODO: Decide if we wanna delete vote.
		// vm.storage.DeleteVote(ctx, proposalID, vote.Voter)
	}

	// put all validators who didn't vote on these two types proposal into penalty list
	if proposalType == types.ChangeParam || proposalType == types.ProtocolUpgrade {
		penaltyList.PenaltyList = oncallValidators
	}
	return penaltyList, nil
}

// GetLinoStake - get lino stake
func (vm VoteManager) GetLinoStake(ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, accKey)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return voter.LinoStake, nil
}

// GetLinoStakeLastChangedAt - get linoStake last changed time
func (vm VoteManager) GetLinoStakeLastChangedAt(ctx sdk.Context, accKey types.AccountKey) (int64, sdk.Error) {
	voter, err := vm.storage.GetVoter(ctx, accKey)
	if err != nil {
		return 0, err
	}
	return voter.LastPowerChangeAt, nil
}

// SetLinoStakeLastChangedAt - set linoStake last changed time
func (vm VoteManager) SetLinoStakeLastChangedAt(ctx sdk.Context, accKey types.AccountKey, lastChangedAt int64) sdk.Error {
	voter, err := vm.storage.GetVoter(ctx, accKey)
	if err != nil {
		return err
	}
	voter.LastPowerChangeAt = lastChangedAt
	if err := vm.storage.SetVoter(ctx, accKey, voter); err != nil {
		return err
	}
	return nil
}

// GetAllDelegators - get all delegators of a voter
func (vm VoteManager) GetAllDelegators(ctx sdk.Context, voterName types.AccountKey) ([]types.AccountKey, sdk.Error) {
	return vm.storage.GetAllDelegators(ctx, voterName)
}

// GetValidatorReferenceList - get all delegatee
func (vm VoteManager) GetValidatorReferenceList(ctx sdk.Context) (*model.ReferenceList, sdk.Error) {
	return vm.storage.GetReferenceList(ctx)
}

// SetValidatorReferenceList - set validator reference list
func (vm VoteManager) SetValidatorReferenceList(ctx sdk.Context, lst *model.ReferenceList) sdk.Error {
	return vm.storage.SetReferenceList(ctx, lst)
}

// Export storage state.
func (vm VoteManager) Export(ctx sdk.Context) *model.VoterTables {
	return vm.storage.Export(ctx)
}

// Import storage state.
func (vm VoteManager) Import(ctx sdk.Context, voter *model.VoterTablesIR) {
	vm.storage.Import(ctx, voter)
}
