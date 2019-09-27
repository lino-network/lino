package manager

import (
	"math"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/validator/model"
	"github.com/lino-network/lino/x/validator/types"
	"github.com/lino-network/lino/x/vote"
	votetypes "github.com/lino-network/lino/x/vote/types"
)

// ValidatorManager - validator manager
type ValidatorManager struct {
	storage model.ValidatorStorage

	// deps
	paramHolder param.ParamKeeper
	vote        vote.VoteKeeper
	global      global.GlobalKeeper
	acc         acc.AccountKeeper
}

func NewValidatorManager(key sdk.StoreKey, holder param.ParamKeeper, vote vote.VoteKeeper,
	global global.GlobalKeeper, acc acc.AccountKeeper) ValidatorManager {
	return ValidatorManager{
		storage:     model.NewValidatorStorage(key),
		paramHolder: holder,
		vote:        vote,
		global:      global,
		acc:         acc,
	}
}

// InitGenesis - initialize KVStore
func (vm ValidatorManager) InitGenesis(ctx sdk.Context) error {
	if err := vm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

// RegisterValidator - register a validator.
func (vm ValidatorManager) RegisterValidator(ctx sdk.Context, username linotypes.AccountKey, valPubKey crypto.PubKey, link string) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	if linotypes.FindAccountInList(username, lst.Jail) != -1 {
		return vm.rejoinFromJail(ctx, username)
	}

	if vm.DoesValidatorExist(ctx, username) {
		return types.ErrValidatorAlreadyExist()
	}

	// must be voter duty
	if duty, err := vm.vote.GetVoterDuty(ctx, username); err != nil || duty != votetypes.DutyVoter {
		return types.ErrInvalidVoterDuty()
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	// assign validator duty in vote
	if err = vm.vote.AssignDuty(ctx, username, votetypes.DutyValidator, param.ValidatorMinDeposit); err != nil {
		return err
	}

	if err := vm.CheckDupPubKey(ctx, valPubKey); err != nil {
		return err
	}

	validator := &model.Validator{
		ABCIValidator: abci.Validator{
			Address: valPubKey.Address(),
			Power:   0,
		},
		PubKey:   valPubKey,
		Username: username,
		Link:     link,
	}

	if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
		return err
	}

	// join as candidate validator first and vote itself
	if err := vm.addValidatortToCandidateList(ctx, username); err != nil {
		return err
	}
	if err := vm.VoteValidator(ctx, username, []linotypes.AccountKey{username}); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) RevokeValidator(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	me.HasRevoked = true
	if err := vm.storage.SetValidator(ctx, username, me); err != nil {
		return err
	}

	if err := vm.removeValidatorFromAllLists(ctx, username); err != nil {
		return err
	}
	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	if err = vm.vote.UnassignDuty(ctx, username, param.ValidatorRevokePendingSec); err != nil {
		return err
	}

	return nil
}

func (vm ValidatorManager) VoteValidator(ctx sdk.Context, username linotypes.AccountKey,
	votedValidators []linotypes.AccountKey) sdk.Error {
	// check if voted validators exist
	for _, valName := range votedValidators {
		if !vm.storage.DoesValidatorExist(ctx, valName) {
			return types.ErrValidatorNotFound()
		}
	}

	if !vm.storage.DoesElectionVoteListExist(ctx, username) {
		// init election vote list if not exist
		lst := &model.ElectionVoteList{}
		if err := vm.storage.SetElectionVoteList(ctx, username, lst); err != nil {
			return err
		}
	}

	updates, err := vm.getElectionVoteListUpdates(ctx, username, votedValidators)
	if err != nil {
		return err
	}

	if err := vm.updateValidatorReceivedVotes(ctx, updates); err != nil {
		return err
	}

	if err := vm.setNewElectionVoteList(ctx, username, votedValidators); err != nil {
		return err
	}

	if err := vm.vote.ClaimInterest(ctx, username); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) DistributeInflationToValidator(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	coin, err := vm.global.GetValidatorHourlyInflation(ctx)
	if err != nil {
		return err
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	totalWeight := int64(len(lst.Oncall))*param.OncallInflationWeight +
		int64(len(lst.Standby))*param.StandbyInflationWeight
	index := int64(0)
	// give inflation to each validator according it's weight
	for _, oncall := range lst.Oncall {
		ratPerOncall := coin.ToDec().Mul(sdk.NewDec(param.OncallInflationWeight)).Quo(sdk.NewDec(totalWeight - index))
		err := vm.acc.AddCoinToUsername(ctx, oncall, linotypes.DecToCoin(ratPerOncall))
		if err != nil {
			return err
		}
		coin = coin.Minus(linotypes.DecToCoin(ratPerOncall))
		index += param.OncallInflationWeight
	}

	for _, standby := range lst.Standby {
		ratPerStandby := coin.ToDec().Mul(sdk.NewDec(param.StandbyInflationWeight)).Quo(sdk.NewDec(totalWeight - index))
		err := vm.acc.AddCoinToUsername(ctx, standby, linotypes.DecToCoin(ratPerStandby))
		if err != nil {
			return err
		}
		coin = coin.Minus(linotypes.DecToCoin(ratPerStandby))
		index += param.StandbyInflationWeight
	}
	return nil
}

func (vm ValidatorManager) rejoinFromJail(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	totalStake, err := vm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return err
	}

	if !totalStake.IsGTE(param.ValidatorMinDeposit) {
		return types.ErrInsufficientDeposit()
	}

	if err := vm.removeValidatorFromJailList(ctx, username); err != nil {
		return err
	}
	if err := vm.addValidatortToCandidateList(ctx, username); err != nil {
		return err
	}
	if err := vm.onCandidateVotesInc(ctx, username); err != nil {
		return err
	}
	return nil
}

// calculate the changed votes between current election votes and previous election votes
// negative number means the corresponding validator need to decrease it's received votes
// positive number means the corresponding validator need to increase it's received votes
func (vm ValidatorManager) getElectionVoteListUpdates(ctx sdk.Context, username linotypes.AccountKey,
	votedValidators []linotypes.AccountKey) ([]*model.ElectionVote, sdk.Error) {
	res := []*model.ElectionVote{}
	prevList, err := vm.storage.GetElectionVoteList(ctx, username)
	if err != nil {
		return nil, err
	}
	totalStake, err := vm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return nil, err
	}
	voteStakeDec := totalStake.ToDec().Quo(sdk.NewDec(int64(len(votedValidators))))

	// add all old votes into res set first and default all votes are negative (not in the new list)
	for _, oldVote := range prevList.ElectionVotes {
		changeDec := oldVote.Vote.ToDec().Mul(sdk.NewDec(-1))
		res = append(res, &model.ElectionVote{
			ValidatorName: oldVote.ValidatorName,
			Vote:          linotypes.DecToCoin(changeDec),
		})
	}

	for _, validatorName := range votedValidators {
		found := false
		for _, oldVote := range res {
			if oldVote.ValidatorName == validatorName {
				found = true
				oldVote.Vote = linotypes.DecToCoin(oldVote.Vote.ToDec().Add(voteStakeDec))
			}
		}

		// newly voted validator
		if !found {
			res = append(res, &model.ElectionVote{
				ValidatorName: validatorName,
				Vote:          linotypes.DecToCoin(voteStakeDec),
			})
		}
	}
	return res, nil
}

func (vm ValidatorManager) updateValidatorReceivedVotes(ctx sdk.Context, updates []*model.ElectionVote) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	for _, update := range updates {
		if update.Vote.IsZero() {
			continue
		}

		// the voted validator may have revoked, just continue
		validator, err := vm.storage.GetValidator(ctx, update.ValidatorName)
		if err != nil {
			continue
		}
		validator.ReceivedVotes = validator.ReceivedVotes.Plus(update.Vote)
		if err := vm.storage.SetValidator(ctx, update.ValidatorName, validator); err != nil {
			return err
		}

		// the corresponding validator's received votes increase
		if update.Vote.IsPositive() {
			if linotypes.FindAccountInList(update.ValidatorName, lst.Oncall) != -1 {
				if err := vm.onOncallVotesInc(ctx, update.ValidatorName); err != nil {
					return err
				}
			}
			if linotypes.FindAccountInList(update.ValidatorName, lst.Standby) != -1 {
				if err := vm.onStandbyVotesInc(ctx, update.ValidatorName); err != nil {
					return err
				}
			}
			if linotypes.FindAccountInList(update.ValidatorName, lst.Candidates) != -1 {
				if err := vm.onCandidateVotesInc(ctx, update.ValidatorName); err != nil {
					return err
				}
			}
		} else {
			// the corresponding validator's received votes decrease
			if linotypes.FindAccountInList(update.ValidatorName, lst.Oncall) != -1 {
				if err := vm.onOncallVotesDec(ctx, update.ValidatorName); err != nil {
					return err
				}
			}
			if linotypes.FindAccountInList(update.ValidatorName, lst.Standby) != -1 {
				if err := vm.onStandbyVotesDec(ctx, update.ValidatorName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (vm ValidatorManager) setNewElectionVoteList(ctx sdk.Context, username linotypes.AccountKey,
	votedValidators []linotypes.AccountKey) sdk.Error {
	if len(votedValidators) == 0 {
		return nil
	}
	lst := &model.ElectionVoteList{}
	totalStake, err := vm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return err
	}

	voteStakeDec := totalStake.ToDec().Quo(sdk.NewDec(int64(len(votedValidators))))
	for _, validatorName := range votedValidators {
		electionVote := model.ElectionVote{
			ValidatorName: validatorName,
			Vote:          linotypes.DecToCoin(voteStakeDec),
		}
		lst.ElectionVotes = append(lst.ElectionVotes, electionVote)
	}

	if err := vm.storage.SetElectionVoteList(ctx, username, lst); err != nil {
		return err
	}
	return nil
}

// DoesValidatorExist - check if validator exists in KVStore or not
func (vm ValidatorManager) DoesValidatorExist(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	return vm.storage.DoesValidatorExist(ctx, accKey)
}

// GetInitValidators return all validators in state.
// XXX(yumin): This is intended to be used only in initChainer
func (vm ValidatorManager) GetInitValidators(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error) {
	committingValidators, err := vm.GetCommittingValidators(ctx)
	if err != nil {
		return nil, err
	}
	updates := []abci.ValidatorUpdate{}
	for _, curValidator := range committingValidators {
		validator, err := vm.storage.GetValidator(ctx, curValidator)
		if err != nil {
			return nil, err
		}
		updates = append(updates, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(validator.PubKey),
			Power:  validator.ABCIValidator.Power,
		})
	}
	return updates, nil
}

// GetValidatorUpdates - after a block, compare updated validator set with
// recorded validator set before block execution
func (vm ValidatorManager) GetValidatorUpdates(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error) {
	validatorList, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return nil, err
	}
	updates := []abci.ValidatorUpdate{}
	committingValidators, err := vm.GetCommittingValidators(ctx)
	if err != nil {
		return nil, err
	}

	for _, preValidator := range validatorList.PreBlockValidators {
		// set power to 0 if a previous validator not in oncall and standby list anymore
		if linotypes.FindAccountInList(preValidator, committingValidators) == -1 {
			validator, err := vm.storage.GetValidator(ctx, preValidator)
			if err != nil {
				return nil, err
			}
			// delete revoked validator
			if validator.HasRevoked {
				if err := vm.storage.DeleteValidator(ctx, validator.Username); err != nil {
					return nil, err
				}
			}
			updates = append(updates, abci.ValidatorUpdate{
				PubKey: tmtypes.TM2PB.PubKey(validator.PubKey),
				Power:  0,
			})
		}
	}

	for _, curValidator := range committingValidators {
		validator, err := vm.storage.GetValidator(ctx, curValidator)
		if err != nil {
			return nil, err
		}

		updates = append(updates, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(validator.PubKey),
			Power:  validator.ABCIValidator.Power,
		})

	}
	return updates, nil
}

// UpdateSigningStats - based on info in beginBlocker, record last block singing info
func (vm ValidatorManager) UpdateSigningStats(ctx sdk.Context, voteInfos []abci.VoteInfo) sdk.Error {
	committingValidators, err := vm.GetCommittingValidators(ctx)
	if err != nil {
		return err
	}

	// map address to whether that validator has signed.
	addressSigned := make(map[string]bool)
	for _, voteInfo := range voteInfos {
		addressSigned[string(voteInfo.Validator.Address)] = voteInfo.SignedLastBlock
	}

	// go through oncall and standby validator list to get all address and name mapping
	for _, curValidator := range committingValidators {
		validator, getErr := vm.storage.GetValidator(ctx, curValidator)
		if getErr != nil {
			return err
		}
		signed, exist := addressSigned[string(validator.ABCIValidator.Address)]
		if !exist || !signed {
			validator.AbsentCommit++
		} else {
			validator.ProducedBlocks++
			if validator.AbsentCommit > 0 {
				validator.AbsentCommit--
			}
		}
		if err := vm.storage.SetValidator(ctx, curValidator, validator); err != nil {
			return err
		}
	}

	return nil
}

// PunishOncallValidator - punish committing validator if 1) byzantine or 2) missing blocks reach limiation
func (vm ValidatorManager) punishCommittingValidator(ctx sdk.Context, username linotypes.AccountKey,
	penalty linotypes.Coin, punishType linotypes.PunishType) sdk.Error {
	// slash and add slashed coin back into validator inflation pool
	actualPenalty, err := vm.vote.SlashStake(ctx, username, penalty)
	if err != nil {
		return err
	}

	if err := vm.global.AddToValidatorInflationPool(ctx, actualPenalty); err != nil {
		return err
	}

	if punishType == linotypes.PunishAbsentCommit {
		// reset absent commit
		validator, err := vm.storage.GetValidator(ctx, username)
		if err != nil {
			return err
		}
		validator.AbsentCommit = 0
		if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
			return err
		}

	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	totalStake, err := vm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return err
	}

	// remove this validator and put into jail if its remaining stake is not enough
	// OR, we explicitly want to fire this validator
	if punishType == linotypes.PunishByzantine || !totalStake.IsGTE(param.ValidatorMinDeposit) {
		if err := vm.removeValidatorFromAllLists(ctx, username); err != nil {
			return err
		}
		if err := vm.addValidatortToJailList(ctx, username); err != nil {
			return err
		}

		if err := vm.balanceValidatorList(ctx); err != nil {
			return err
		}
	}

	return nil
}

// FireIncompetentValidator - fire oncall validator if 1) deposit insufficient 2) byzantine
func (vm ValidatorManager) FireIncompetentValidator(ctx sdk.Context,
	byzantineValidators []abci.Evidence) sdk.Error {
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	committingValidators, err := vm.GetCommittingValidators(ctx)
	if err != nil {
		return err
	}

	for _, validatorName := range committingValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}

		for _, evidence := range byzantineValidators {
			if reflect.DeepEqual(validator.ABCIValidator.Address, evidence.Validator.Address) {
				if err := vm.punishCommittingValidator(ctx, validator.Username, param.PenaltyByzantine,
					linotypes.PunishByzantine); err != nil {
					return err
				}
				break
			}
		}

		if validator.AbsentCommit > param.AbsentCommitLimitation {
			if err := vm.punishCommittingValidator(ctx, validator.Username, param.PenaltyMissCommit,
				linotypes.PunishAbsentCommit); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm ValidatorManager) CheckDupPubKey(ctx sdk.Context, pubKey crypto.PubKey) sdk.Error {
	// make sure the pub key has not been registered
	allValidators, err := vm.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	for _, validatorName := range allValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}
		// XXX(yumin): ABCIValidator no longer has pubkey, changed to address
		if reflect.DeepEqual(validator.ABCIValidator.Address, pubKey.Address().Bytes()) {
			return types.ErrValidatorPubKeyAlreadyExist()
		}
	}

	return nil
}

func (vm ValidatorManager) onStakeChange(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	if !vm.storage.DoesElectionVoteListExist(ctx, username) {
		return nil
	}

	lst, err := vm.storage.GetElectionVoteList(ctx, username)
	if err != nil {
		return err
	}

	oldList := []linotypes.AccountKey{}
	for _, electionVote := range lst.ElectionVotes {
		oldList = append(oldList, electionVote.ValidatorName)
	}

	updates, err := vm.getElectionVoteListUpdates(ctx, username, oldList)
	if err != nil {
		return err
	}

	if err := vm.updateValidatorReceivedVotes(ctx, updates); err != nil {
		return err
	}

	if err := vm.setNewElectionVoteList(ctx, username, oldList); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) onCandidateVotesInc(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	me, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	if !me.ReceivedVotes.IsGT(lst.LowestStandbyVotes) {
		return nil
	}

	if err := vm.removeValidatorFromCandidateList(ctx, username); err != nil {
		return err
	}

	if me.ReceivedVotes.IsGT(lst.LowestOncallVotes) {
		// join the oncall validator list
		if err := vm.addValidatortToOncallList(ctx, username); err != nil {
			return err
		}
	} else {
		// join the standby validator list
		if err := vm.addValidatortToStandbyList(ctx, username); err != nil {
			return err
		}
	}

	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) onStandbyVotesInc(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	me, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// join the oncall validator list
	if me.ReceivedVotes.IsGT(lst.LowestOncallVotes) {
		if err := vm.removeValidatorFromStandbyList(ctx, username); err != nil {
			return err
		}
		if err := vm.addValidatortToOncallList(ctx, username); err != nil {
			return err
		}
	}
	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) onOncallVotesInc(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) onStandbyVotesDec(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	validator, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	if validator.ReceivedVotes.IsGTE(lst.LowestStandbyVotes) {
		return nil
	}

	// need to exchange position with the highest votes candidates if it's votes greater than me.
	highestCandidate, highestCandidateVotes, err := vm.getHighestVotesAndValidator(ctx, lst.Candidates)
	if err != nil {
		return err
	}

	if highestCandidateVotes.IsGT(validator.ReceivedVotes) {
		if err := vm.removeValidatorFromCandidateList(ctx, highestCandidate); err != nil {
			return err
		}
		if err := vm.removeValidatorFromStandbyList(ctx, username); err != nil {
			return err
		}
		if err := vm.addValidatortToCandidateList(ctx, username); err != nil {
			return err
		}
		if err := vm.addValidatortToStandbyList(ctx, highestCandidate); err != nil {
			return err
		}

	}
	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) onOncallVotesDec(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	validator, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	if validator.ReceivedVotes.IsGTE(lst.LowestOncallVotes) {
		return nil
	}

	// need to exchange position with the highest votes standby if it's votes greater than me.
	highestStandby, highestStandbyVotes, err := vm.getHighestVotesAndValidator(ctx, lst.Standby)
	if err != nil {
		return err
	}

	if highestStandbyVotes.IsGT(validator.ReceivedVotes) {
		if err := vm.removeValidatorFromStandbyList(ctx, highestStandby); err != nil {
			return err
		}
		if err := vm.removeValidatorFromOncallList(ctx, username); err != nil {
			return err
		}
		if err := vm.addValidatortToOncallList(ctx, highestStandby); err != nil {
			return err
		}
		if err := vm.addValidatortToStandbyList(ctx, username); err != nil {
			return err
		}
	}

	// need to exchange position with the highest votes candidate if it's votes greater than me.
	highestCandidate, highestCandidateVotes, err := vm.getHighestVotesAndValidator(ctx, lst.Candidates)
	if err != nil {
		return err
	}

	if highestCandidateVotes.IsGT(validator.ReceivedVotes) {
		if err := vm.removeValidatorFromCandidateList(ctx, highestCandidate); err != nil {
			return err
		}
		if err := vm.removeValidatorFromStandbyList(ctx, username); err != nil {
			return err
		}
		if err := vm.addValidatortToCandidateList(ctx, username); err != nil {
			return err
		}
		if err := vm.addValidatortToStandbyList(ctx, highestCandidate); err != nil {
			return err
		}
	}

	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) balanceValidatorList(ctx sdk.Context) sdk.Error {
	if err := vm.removeExtraOncall(ctx); err != nil {
		return err
	}
	if err := vm.removeExtraStandby(ctx); err != nil {
		return err
	}
	if err := vm.fillEmptyOncall(ctx); err != nil {
		return err
	}
	if err := vm.fillEmptyStandby(ctx); err != nil {
		return err
	}
	if err := vm.updateLowestOncall(ctx); err != nil {
		return err
	}
	if err := vm.updateLowestStandby(ctx); err != nil {
		return err
	}
	return nil
}

// move lowest votes oncall validator to standby list if current oncall size exceeds max
func (vm ValidatorManager) removeExtraOncall(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	curLen := int64(len(lst.Oncall))
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	for curLen > param.OncallSize {
		lst, err := vm.storage.GetValidatorList(ctx)
		if err != nil {
			return err
		}
		lowestOncall, _, err := vm.getLowestVotesAndValidator(ctx, lst.Oncall)
		if err != nil {
			return err
		}
		if err := vm.removeValidatorFromOncallList(ctx, lowestOncall); err != nil {
			return err
		}
		if err := vm.addValidatortToStandbyList(ctx, lowestOncall); err != nil {
			return err
		}
		curLen--
	}
	return nil
}

// move lowest votes standby validator to candidate list if current standby size exceeds max
func (vm ValidatorManager) removeExtraStandby(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	curLen := int64(len(lst.Standby))
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	for curLen > param.StandbySize {
		lst, err := vm.storage.GetValidatorList(ctx)
		if err != nil {
			return err
		}
		lowestStandby, _, err := vm.getLowestVotesAndValidator(ctx, lst.Standby)
		if err != nil {
			return err
		}
		if err := vm.removeValidatorFromStandbyList(ctx, lowestStandby); err != nil {
			return err
		}
		if err := vm.addValidatortToCandidateList(ctx, lowestStandby); err != nil {
			return err
		}
		curLen--
	}
	return nil
}

// move highest votes standby validator to oncall list if current oncall size not full
// if standby validator not enough, try move highest votes candidate validator
func (vm ValidatorManager) fillEmptyOncall(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	lenOncall := int64(len(lst.Oncall))
	lenStandby := int64(len(lst.Standby))
	lenCandidate := int64(len(lst.Candidates))
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	for lenOncall < param.OncallSize && lenStandby > 0 {
		lst, err := vm.storage.GetValidatorList(ctx)
		if err != nil {
			return err
		}
		highestStandby, _, err := vm.getHighestVotesAndValidator(ctx, lst.Standby)
		if err != nil {
			return err
		}
		if err := vm.removeValidatorFromStandbyList(ctx, highestStandby); err != nil {
			return err
		}
		if err := vm.addValidatortToOncallList(ctx, highestStandby); err != nil {
			return err
		}
		lenOncall++
		lenStandby--
	}

	for lenOncall < param.OncallSize && lenCandidate > 0 {
		lst, err := vm.storage.GetValidatorList(ctx)
		if err != nil {
			return err
		}
		highestCandidate, _, err := vm.getHighestVotesAndValidator(ctx, lst.Candidates)
		if err != nil {
			return err
		}
		if err := vm.removeValidatorFromCandidateList(ctx, highestCandidate); err != nil {
			return err
		}
		if err := vm.addValidatortToOncallList(ctx, highestCandidate); err != nil {
			return err
		}
		lenOncall++
		lenCandidate--
	}

	return nil
}

// move highest votes candidate validator to standby list if current standby size not full
func (vm ValidatorManager) fillEmptyStandby(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	lenStandby := int64(len(lst.Standby))
	lenCandidate := int64(len(lst.Candidates))
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}

	for lenStandby < param.StandbySize && lenCandidate > 0 {
		lst, err := vm.storage.GetValidatorList(ctx)
		if err != nil {
			return err
		}
		highestCandidate, _, err := vm.getHighestVotesAndValidator(ctx, lst.Candidates)
		if err != nil {
			return err
		}
		if err := vm.removeValidatorFromCandidateList(ctx, highestCandidate); err != nil {
			return err
		}
		if err := vm.addValidatortToStandbyList(ctx, highestCandidate); err != nil {
			return err
		}
		lenStandby++
		lenCandidate--
	}
	return nil
}

func (vm ValidatorManager) removeValidatorFromAllLists(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	if err := vm.removeValidatorFromOncallList(ctx, username); err != nil {
		return err
	}
	if err := vm.removeValidatorFromStandbyList(ctx, username); err != nil {
		return err
	}
	if err := vm.removeValidatorFromCandidateList(ctx, username); err != nil {
		return err
	}
	if err := vm.removeValidatorFromJailList(ctx, username); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) removeValidatorFromOncallList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Oncall = removeFromList(username, lst.Oncall)
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) removeValidatorFromStandbyList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Standby = removeFromList(username, lst.Standby)
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) removeValidatorFromCandidateList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Candidates = removeFromList(username, lst.Candidates)
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) removeValidatorFromJailList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Jail = removeFromList(username, lst.Jail)
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) addValidatortToOncallList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Oncall = append(lst.Oncall, username)
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	votesCoinInt64, err := me.ReceivedVotes.ToInt64()
	if err != nil {
		return err
	}
	// set oncall validator committing power equal to it's votes (lino)
	me.ABCIValidator.Power = votesCoinInt64 / linotypes.Decimals
	if err := vm.storage.SetValidator(ctx, username, me); err != nil {
		return err
	}
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) addValidatortToStandbyList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Standby = append(lst.Standby, username)
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// set oncall validator committing power equal to 1
	me.ABCIValidator.Power = 1
	if err := vm.storage.SetValidator(ctx, username, me); err != nil {
		return err
	}
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) addValidatortToCandidateList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Candidates = append(lst.Candidates, username)
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// set oncall validator committing power equal to 0
	me.ABCIValidator.Power = 0
	if err := vm.storage.SetValidator(ctx, username, me); err != nil {
		return err
	}
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) addValidatortToJailList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	lst.Jail = append(lst.Jail, username)
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// set oncall validator committing power equal to 0
	me.ABCIValidator.Power = 0
	if err := vm.storage.SetValidator(ctx, username, me); err != nil {
		return err
	}
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func removeFromList(me linotypes.AccountKey, users []linotypes.AccountKey) []linotypes.AccountKey {
	for i := 0; i < len(users); i++ {
		if me == users[i] {
			return append(users[:i], users[i+1:]...)
		}
	}
	return users
}

func (vm ValidatorManager) getHighestVotesAndValidator(ctx sdk.Context,
	lst []linotypes.AccountKey) (linotypes.AccountKey, linotypes.Coin, sdk.Error) {
	highestValdator := linotypes.AccountKey("")
	highestValdatorVotes := linotypes.NewCoinFromInt64(0)

	for i := range lst {
		validator, err := vm.storage.GetValidator(ctx, lst[i])
		if err != nil {
			return highestValdator, linotypes.NewCoinFromInt64(0), err
		}
		if validator.ReceivedVotes.IsGT(highestValdatorVotes) {
			highestValdator = validator.Username
			highestValdatorVotes = validator.ReceivedVotes
		}
	}
	return highestValdator, highestValdatorVotes, nil
}

func (vm ValidatorManager) getLowestVotesAndValidator(ctx sdk.Context,
	lst []linotypes.AccountKey) (linotypes.AccountKey, linotypes.Coin, sdk.Error) {
	lowestValdator := linotypes.AccountKey("")
	lowestValdatorVotes := linotypes.NewCoinFromInt64(math.MaxInt64)

	for i := range lst {
		validator, err := vm.storage.GetValidator(ctx, lst[i])
		if err != nil {
			return lowestValdator, linotypes.NewCoinFromInt64(0), err
		}
		if lowestValdatorVotes.IsGT(validator.ReceivedVotes) {
			lowestValdator = validator.Username
			lowestValdatorVotes = validator.ReceivedVotes
		}
	}
	return lowestValdator, lowestValdatorVotes, nil
}

func (vm ValidatorManager) updateLowestOncall(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	newLowestVotes := linotypes.NewCoinFromInt64(math.MaxInt64)
	newLowestValidator := linotypes.AccountKey("")

	for _, validatorKey := range lst.Oncall {
		validator, err := vm.storage.GetValidator(ctx, validatorKey)
		if err != nil {
			return err
		}

		if newLowestVotes.IsGT(validator.ReceivedVotes) {
			newLowestVotes = validator.ReceivedVotes
			newLowestValidator = validator.Username
		}
	}

	// set the new lowest power
	if len(lst.Oncall) == 0 {
		lst.LowestOncallVotes = linotypes.NewCoinFromInt64(0)
		lst.LowestOncall = linotypes.AccountKey("")
	} else {
		lst.LowestOncallVotes = newLowestVotes
		lst.LowestOncall = newLowestValidator
	}

	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) updateLowestStandby(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	newLowestVotes := linotypes.NewCoinFromInt64(math.MaxInt64)
	newLowestValidator := linotypes.AccountKey("")

	for _, validatorKey := range lst.Standby {
		validator, err := vm.storage.GetValidator(ctx, validatorKey)
		if err != nil {
			return err
		}

		if newLowestVotes.IsGT(validator.ReceivedVotes) {
			newLowestVotes = validator.ReceivedVotes
			newLowestValidator = validator.Username
		}
	}

	// set the new lowest power
	if len(lst.Standby) == 0 {
		lst.LowestStandbyVotes = linotypes.NewCoinFromInt64(0)
		lst.LowestStandby = linotypes.AccountKey("")
	} else {

		lst.LowestStandbyVotes = newLowestVotes
		lst.LowestStandby = newLowestValidator
	}

	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

// getter and setter
func (vm ValidatorManager) GetValidator(ctx sdk.Context, accKey linotypes.AccountKey) (*model.Validator, sdk.Error) {
	return vm.storage.GetValidator(ctx, accKey)
}

func (vm ValidatorManager) GetAllValidators(ctx sdk.Context) ([]linotypes.AccountKey, sdk.Error) {
	lst, err := vm.GetValidatorList(ctx)
	if err != nil {
		return []linotypes.AccountKey{}, err
	}
	tmp := append(lst.Standby, lst.Candidates...)
	return append(lst.Oncall, tmp...), nil
}

func (vm ValidatorManager) GetCommittingValidators(ctx sdk.Context) ([]linotypes.AccountKey, sdk.Error) {
	lst, err := vm.GetValidatorList(ctx)
	if err != nil {
		return []linotypes.AccountKey{}, err
	}
	return append(lst.Oncall, lst.Standby...), nil
}

func (vm ValidatorManager) GetValidatorList(ctx sdk.Context) (*model.ValidatorList, sdk.Error) {
	return vm.storage.GetValidatorList(ctx)
}

func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, lst *model.ValidatorList) sdk.Error {
	return vm.storage.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) GetElectionVoteList(ctx sdk.Context,
	accKey linotypes.AccountKey) (*model.ElectionVoteList, sdk.Error) {
	return vm.storage.GetElectionVoteList(ctx, accKey)
}

// // Export storage state.
// func (vm ValidatorManager) Export(ctx sdk.Context) *model.ValidatorTables {
// 	return vm.storage.Export(ctx)
// }

// // Import storage state.
// func (vm ValidatorManager) Import(ctx sdk.Context, tb *model.ValidatorTablesIR) {
// 	vm.storage.Import(ctx, tb)
// }
