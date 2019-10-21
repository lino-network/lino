package manager

import (
	"fmt"
	"math"
	"reflect"
	"runtime/debug"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	acc "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/validator/model"
	"github.com/lino-network/lino/x/validator/types"
	"github.com/lino-network/lino/x/vote"
	votetypes "github.com/lino-network/lino/x/vote/types"
)

const (
	exportVersion = 1
	importVersion = 1
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
func (vm ValidatorManager) InitGenesis(ctx sdk.Context) {
	lst := &model.ValidatorList{
		LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
		LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
	}
	vm.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) OnBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	// update preblock validators
	validatorList := vm.GetValidatorList(ctx)
	vals := vm.GetCommittingValidators(ctx)
	validatorList.PreBlockValidators = vals
	vm.SetValidatorList(ctx, validatorList)

	// update signing stats.
	updateErr := vm.updateSigningStats(ctx, req.LastCommitInfo.Votes)
	if updateErr != nil {
		panic(updateErr)
	}

	if err := vm.fireIncompetentValidator(ctx, req.ByzantineValidators); err != nil {
		panic(err)
	}
}

func (vm ValidatorManager) UpdateValidator(ctx sdk.Context, username linotypes.AccountKey, link string) sdk.Error {
	val, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}
	val.Link = link
	vm.storage.SetValidator(ctx, username, val)
	return nil
}

// RegisterValidator - register a validator.
func (vm ValidatorManager) RegisterValidator(ctx sdk.Context, username linotypes.AccountKey, valPubKey crypto.PubKey, link string) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
	if linotypes.FindAccountInList(username, lst.Jail) != -1 {
		return vm.rejoinFromJail(ctx, username)
	}

	if vm.IsLegalValidator(ctx, username) {
		return types.ErrValidatorAlreadyExist()
	}

	// must be voter duty
	if duty, err := vm.vote.GetVoterDuty(ctx, username); err != nil || duty != votetypes.DutyVoter {
		return types.ErrInvalidVoterDuty()
	}

	param := vm.paramHolder.GetValidatorParam(ctx)

	// assign validator duty in vote
	if err := vm.vote.AssignDuty(ctx, username, votetypes.DutyValidator, param.ValidatorMinDeposit); err != nil {
		return err
	}

	if err := vm.checkDupPubKey(ctx, valPubKey); err != nil {
		return err
	}

	// recover data if was revoked: inherite the votes.
	prevVotes := vm.getPrevVotes(ctx, username)
	validator := &model.Validator{
		ABCIValidator: abci.Validator{
			Address: valPubKey.Address(),
			Power:   0,
		},
		ReceivedVotes: prevVotes,
		PubKey:        valPubKey,
		Username:      username,
		Link:          link,
	}
	vm.storage.SetValidator(ctx, username, validator)

	// join as candidate validator first and vote itself
	if err := vm.addValidatortToCandidateList(ctx, username); err != nil {
		return err
	}
	if err := vm.onCandidateVotesInc(ctx, username); err != nil {
		return err
	}
	if err := vm.VoteValidator(ctx, username, []linotypes.AccountKey{username}); err != nil {
		return err
	}

	return nil
}

func (vm ValidatorManager) RevokeValidator(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	if !vm.IsLegalValidator(ctx, username) {
		return types.ErrInvalidValidator()
	}

	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	me.HasRevoked = true
	vm.storage.SetValidator(ctx, username, me)

	if err := vm.removeValidatorFromAllLists(ctx, username); err != nil {
		return err
	}
	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}

	param := vm.paramHolder.GetValidatorParam(ctx)
	if err = vm.vote.UnassignDuty(ctx, username, param.ValidatorRevokePendingSec); err != nil {
		return err
	}

	return nil
}

func (vm ValidatorManager) VoteValidator(ctx sdk.Context, username linotypes.AccountKey,
	votedValidators []linotypes.AccountKey) sdk.Error {
	param := vm.paramHolder.GetValidatorParam(ctx)
	if int64(len(votedValidators)) > param.MaxVotedValidators {
		return types.ErrInvalidVotedValidators()
	}
	// check if voted validators exist
	for _, valName := range votedValidators {
		if !vm.IsLegalValidator(ctx, valName) {
			return types.ErrValidatorNotFound(valName)
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
	coin, err := vm.global.GetValidatorHourlyInflation(ctx)
	if err != nil {
		return err
	}

	param := vm.paramHolder.GetValidatorParam(ctx)
	lst := vm.storage.GetValidatorList(ctx)
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
	param := vm.paramHolder.GetValidatorParam(ctx)
	totalStake, err := vm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return err
	}

	if !totalStake.IsGTE(param.ValidatorMinDeposit) {
		return types.ErrInsufficientDeposit()
	}

	vm.removeValidatorFromJailList(ctx, username)
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
	prevList := vm.storage.GetElectionVoteList(ctx, username)
	totalStake, err := vm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return nil, err
	}

	if len(prevList.ElectionVotes) == 0 && len(votedValidators) == 0 {
		return nil, nil
	}
	if len(votedValidators) == 0 {
		return nil, types.ErrInvalidVotedValidators()
	}

	voteStake := linotypes.DecToCoin(
		totalStake.ToDec().Quo(sdk.NewDec(int64(len(votedValidators)))))

	// add all old votes into res set first and default all votes are negative (not in the new list)
	for _, oldVote := range prevList.ElectionVotes {
		changeDec := oldVote.Vote.Neg()
		res = append(res, &model.ElectionVote{
			ValidatorName: oldVote.ValidatorName,
			Vote:          changeDec,
		})
	}

	// a helper function to return the pointer to the matching ElectionVote.
	findInPrev := func(valName linotypes.AccountKey) *model.ElectionVote {
		for _, oldVote := range res {
			if oldVote.ValidatorName == valName {
				return oldVote
			}
		}
		return nil
	}

	for _, validatorName := range votedValidators {
		if prev := findInPrev(validatorName); prev != nil {
			prev.Vote = prev.Vote.Plus(voteStake)
		} else {
			res = append(res, &model.ElectionVote{
				ValidatorName: validatorName,
				Vote:          voteStake,
			})
		}
	}
	return res, nil
}

func (vm ValidatorManager) updateValidatorReceivedVotes(ctx sdk.Context, updates []*model.ElectionVote) sdk.Error {
	fmt.Printf("updates: %+v", updates)
	lst := vm.storage.GetValidatorList(ctx)
	fmt.Printf("validator lst: %+v", lst)
	for _, update := range updates {
		if update.Vote.IsZero() {
			continue
		}

		// revoked validator's record will still be in kv.
		validator, err := vm.storage.GetValidator(ctx, update.ValidatorName)
		if err != nil {
			return err
		}
		validator.ReceivedVotes = validator.ReceivedVotes.Plus(update.Vote)
		vm.storage.SetValidator(ctx, update.ValidatorName, validator)

		// the corresponding validator's received votes increase
		if update.Vote.IsPositive() {
			fmt.Printf("update.ValidatorName: %+v", update.ValidatorName)
			fmt.Printf("lst.Oncall: %+v", lst.Oncall)
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

	vm.storage.SetElectionVoteList(ctx, username, lst)
	return nil
}

// IsLegalValidator - check if the validator is a validator and not revoked.
func (vm ValidatorManager) IsLegalValidator(ctx sdk.Context, accKey linotypes.AccountKey) bool {
	val, err := vm.storage.GetValidator(ctx, accKey)
	if err != nil {
		return false
	}
	return !val.HasRevoked
}

// GetInitValidators return all validators in state.
// XXX(yumin): This is intended to be used only in initChainer
func (vm ValidatorManager) GetInitValidators(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error) {
	committingValidators := vm.GetCommittingValidators(ctx)
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
	validatorList := vm.storage.GetValidatorList(ctx)
	updates := []abci.ValidatorUpdate{}
	committingValidators := vm.GetCommittingValidators(ctx)
	committingSet := linotypes.AccountListToSet(committingValidators)

	for _, preValidator := range validatorList.PreBlockValidators {
		// set power to 0 if a previous validator not in oncall and standby list anymore
		if committingSet[preValidator] == false {
			validator, err := vm.storage.GetValidator(ctx, preValidator)
			if err != nil {
				return nil, err
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
func (vm ValidatorManager) updateSigningStats(ctx sdk.Context, voteInfos []abci.VoteInfo) sdk.Error {
	// map address to whether that validator has signed.
	addressSigned := make(map[string]bool)
	for _, voteInfo := range voteInfos {
		addressSigned[string(voteInfo.Validator.Address)] = voteInfo.SignedLastBlock
	}

	// go through oncall and standby validator list to get all address and name mapping
	committingValidators := vm.GetCommittingValidators(ctx)
	for _, curValidator := range committingValidators {
		validator, err := vm.storage.GetValidator(ctx, curValidator)
		if err != nil {
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
		vm.storage.SetValidator(ctx, curValidator, validator)
	}

	return nil
}

// PunishOncallValidator - punish committing validator
// if 1) byzantine or 2) missing blocks reach limiation
func (vm ValidatorManager) PunishCommittingValidator(ctx sdk.Context, username linotypes.AccountKey,
	penalty linotypes.Coin, punishType linotypes.PunishType) sdk.Error {
	// slash and add slashed coin back into validator inflation pool
	actualPenalty, err := vm.vote.SlashStake(ctx, username, penalty)
	if err != nil {
		return err
	}
	if err := vm.global.AddToValidatorInflationPool(ctx, actualPenalty); err != nil {
		return err
	}
	validator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}
	validator.NumSlash++
	// reset absent commit
	if punishType == linotypes.PunishAbsentCommit {
		validator.AbsentCommit = 0
	}
	vm.storage.SetValidator(ctx, username, validator)

	totalStake, err := vm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return err
	}
	// remove this validator and put into jail if its remaining stake is not enough
	// OR, this is byzantine validator
	// OR, the num of slash exceeds limit
	param := vm.paramHolder.GetValidatorParam(ctx)
	if punishType == linotypes.PunishByzantine ||
		!totalStake.IsGTE(param.ValidatorMinDeposit) ||
		validator.NumSlash > param.SlashLimitation {
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

// FireIncompetentValidator - fire oncall validator if
// 1. absent commit > absent limitation.
// 2. byzantine
func (vm ValidatorManager) fireIncompetentValidator(ctx sdk.Context,
	byzantineValidators []abci.Evidence) sdk.Error {
	param := vm.paramHolder.GetValidatorParam(ctx)
	committingValidators := vm.GetCommittingValidators(ctx)

	for _, validatorName := range committingValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}

		for _, evidence := range byzantineValidators {
			if reflect.DeepEqual(validator.ABCIValidator.Address, evidence.Validator.Address) {
				if err := vm.PunishCommittingValidator(ctx, validator.Username, param.PenaltyByzantine,
					linotypes.PunishByzantine); err != nil {
					return err
				}
				break
			}
		}

		if validator.AbsentCommit > param.AbsentCommitLimitation {
			if err := vm.PunishCommittingValidator(ctx, validator.Username, param.PenaltyMissCommit,
				linotypes.PunishAbsentCommit); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm ValidatorManager) checkDupPubKey(ctx sdk.Context, pubKey crypto.PubKey) sdk.Error {
	// make sure the pub key has not been registered
	allValidators := vm.GetAllValidators(ctx)
	for _, validatorName := range allValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(validator.ABCIValidator.Address, pubKey.Address().Bytes()) {
			return types.ErrValidatorPubKeyAlreadyExist()
		}
	}

	return nil
}

func (vm ValidatorManager) onStakeChange(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst := vm.storage.GetElectionVoteList(ctx, username)
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
	me, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	lst := vm.GetValidatorList(ctx)
	if me.ReceivedVotes.IsGT(lst.LowestOncallVotes) {
		// join the oncall validator list
		vm.removeValidatorFromCandidateList(ctx, username)
		if err := vm.addValidatortToOncallList(ctx, username); err != nil {
			return err
		}
	} else if me.ReceivedVotes.IsGT(lst.LowestStandbyVotes) {
		// join the standby validator list
		vm.removeValidatorFromCandidateList(ctx, username)
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
	me, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// join the oncall validator list
	lst := vm.GetValidatorList(ctx)
	if me.ReceivedVotes.IsGT(lst.LowestOncallVotes) {
		vm.removeValidatorFromStandbyList(ctx, username)
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
	if err := vm.setOncallValidatorPower(ctx, username); err != nil {
		return err
	}
	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) onStandbyVotesDec(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst := vm.GetValidatorList(ctx)
	validator, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	if !validator.ReceivedVotes.IsGTE(lst.LowestStandbyVotes) {
		vm.removeValidatorFromStandbyList(ctx, username)
		if err := vm.addValidatortToCandidateList(ctx, username); err != nil {
			return err
		}
	}

	if err := vm.balanceValidatorList(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) onOncallVotesDec(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	if err := vm.setOncallValidatorPower(ctx, username); err != nil {
		return err
	}

	lst := vm.GetValidatorList(ctx)
	validator, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	if !validator.ReceivedVotes.IsGTE(lst.LowestStandbyVotes) {
		// move to the candidate validator list
		vm.removeValidatorFromOncallList(ctx, username)
		if err := vm.addValidatortToCandidateList(ctx, username); err != nil {
			return err
		}
	} else if !validator.ReceivedVotes.IsGTE(lst.LowestOncallVotes) {
		// move to the standby validator list
		vm.removeValidatorFromOncallList(ctx, username)
		if err := vm.addValidatortToStandbyList(ctx, username); err != nil {
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
	lst := vm.storage.GetValidatorList(ctx)
	curLen := int64(len(lst.Oncall))
	param := vm.paramHolder.GetValidatorParam(ctx)

	for curLen > param.OncallSize {
		lst := vm.storage.GetValidatorList(ctx)
		lowestOncall, _, err := vm.getLowestVotesAndValidator(ctx, lst.Oncall)
		if err != nil {
			return err
		}
		vm.removeValidatorFromOncallList(ctx, lowestOncall)
		if err := vm.addValidatortToStandbyList(ctx, lowestOncall); err != nil {
			return err
		}
		curLen--
	}
	return nil
}

// move lowest votes standby validator to candidate list if current standby size exceeds max
func (vm ValidatorManager) removeExtraStandby(ctx sdk.Context) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
	curLen := int64(len(lst.Standby))
	param := vm.paramHolder.GetValidatorParam(ctx)
	for curLen > param.StandbySize {
		lst := vm.storage.GetValidatorList(ctx)
		lowestStandby, _, err := vm.getLowestVotesAndValidator(ctx, lst.Standby)
		if err != nil {
			return err
		}
		vm.removeValidatorFromStandbyList(ctx, lowestStandby)
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
	lst := vm.storage.GetValidatorList(ctx)
	lenOncall := int64(len(lst.Oncall))
	lenStandby := int64(len(lst.Standby))
	lenCandidate := int64(len(lst.Candidates))
	param := vm.paramHolder.GetValidatorParam(ctx)

	for lenOncall < param.OncallSize && lenStandby > 0 {
		lst := vm.storage.GetValidatorList(ctx)
		highestStandby, _, err := vm.getHighestVotesAndValidator(ctx, lst.Standby)
		if err != nil {
			return err
		}
		vm.removeValidatorFromStandbyList(ctx, highestStandby)
		if err := vm.addValidatortToOncallList(ctx, highestStandby); err != nil {
			return err
		}
		lenOncall++
		lenStandby--
	}

	for lenOncall < param.OncallSize && lenCandidate > 0 {
		lst := vm.storage.GetValidatorList(ctx)
		highestCandidate, _, err := vm.getHighestVotesAndValidator(ctx, lst.Candidates)
		if err != nil {
			return err
		}
		vm.removeValidatorFromCandidateList(ctx, highestCandidate)
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
	lst := vm.storage.GetValidatorList(ctx)
	lenStandby := int64(len(lst.Standby))
	lenCandidate := int64(len(lst.Candidates))
	param := vm.paramHolder.GetValidatorParam(ctx)

	for lenStandby < param.StandbySize && lenCandidate > 0 {
		lst := vm.storage.GetValidatorList(ctx)
		highestCandidate, _, err := vm.getHighestVotesAndValidator(ctx, lst.Candidates)
		if err != nil {
			return err
		}
		vm.removeValidatorFromCandidateList(ctx, highestCandidate)
		if err := vm.addValidatortToStandbyList(ctx, highestCandidate); err != nil {
			return err
		}
		lenStandby++
		lenCandidate--
	}
	return nil
}

func (vm ValidatorManager) removeValidatorFromAllLists(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	vm.removeValidatorFromOncallList(ctx, username)
	vm.removeValidatorFromStandbyList(ctx, username)
	vm.removeValidatorFromCandidateList(ctx, username)
	vm.removeValidatorFromJailList(ctx, username)
	return nil
}

func (vm ValidatorManager) removeValidatorFromOncallList(ctx sdk.Context, username linotypes.AccountKey) {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Oncall = removeFromList(username, lst.Oncall)
	vm.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) removeValidatorFromStandbyList(ctx sdk.Context, username linotypes.AccountKey) {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Standby = removeFromList(username, lst.Standby)
	vm.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) removeValidatorFromCandidateList(ctx sdk.Context, username linotypes.AccountKey) {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Candidates = removeFromList(username, lst.Candidates)
	vm.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) removeValidatorFromJailList(ctx sdk.Context, username linotypes.AccountKey) {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Jail = removeFromList(username, lst.Jail)
	vm.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) addValidatortToOncallList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Oncall = append(lst.Oncall, username)
	if err := vm.setOncallValidatorPower(ctx, username); err != nil {
		return err
	}
	vm.SetValidatorList(ctx, lst)
	return nil
}

func (vm ValidatorManager) addValidatortToStandbyList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Standby = append(lst.Standby, username)
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// set oncall validator committing power equal to 1
	me.ABCIValidator.Power = 1
	vm.storage.SetValidator(ctx, username, me)
	vm.SetValidatorList(ctx, lst)
	return nil
}

func (vm ValidatorManager) addValidatortToCandidateList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Candidates = append(lst.Candidates, username)
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// set oncall validator committing power equal to 0
	me.ABCIValidator.Power = 0
	vm.storage.SetValidator(ctx, username, me)
	vm.SetValidatorList(ctx, lst)
	return nil
}

func (vm ValidatorManager) addValidatortToJailList(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
	lst.Jail = append(lst.Jail, username)
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	// set oncall validator committing power equal to 0 and clear all stats
	me.ABCIValidator.Power = 0
	me.AbsentCommit = 0
	me.NumSlash = 0
	vm.storage.SetValidator(ctx, username, me)
	vm.SetValidatorList(ctx, lst)
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
		if validator.ReceivedVotes.IsGTE(highestValdatorVotes) {
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
		if lowestValdatorVotes.IsGTE(validator.ReceivedVotes) {
			lowestValdator = validator.Username
			lowestValdatorVotes = validator.ReceivedVotes
		}
	}
	return lowestValdator, lowestValdatorVotes, nil
}

func (vm ValidatorManager) updateLowestOncall(ctx sdk.Context) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
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

	vm.SetValidatorList(ctx, lst)
	return nil
}

func (vm ValidatorManager) updateLowestStandby(ctx sdk.Context) sdk.Error {
	lst := vm.storage.GetValidatorList(ctx)
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

	vm.SetValidatorList(ctx, lst)
	return nil
}

func (vm ValidatorManager) setOncallValidatorPower(ctx sdk.Context,
	username linotypes.AccountKey) sdk.Error {
	me, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	votesCoinInt64, err := me.ReceivedVotes.ToInt64()
	if err != nil {
		return err
	}
	// set oncall validator committing power equal to it's votes (lino)
	powerLNO := votesCoinInt64 / linotypes.Decimals
	switch {
	case powerLNO > linotypes.ValidatorMaxPower:
		me.ABCIValidator.Power = linotypes.ValidatorMaxPower
	case powerLNO < 1:
		me.ABCIValidator.Power = 1
	default:
		me.ABCIValidator.Power = powerLNO
	}

	vm.storage.SetValidator(ctx, username, me)
	return nil
}

// getter and setter
func (vm ValidatorManager) GetValidator(ctx sdk.Context, accKey linotypes.AccountKey) (*model.Validator, sdk.Error) {
	return vm.storage.GetValidator(ctx, accKey)
}

func (vm ValidatorManager) GetAllValidators(ctx sdk.Context) []linotypes.AccountKey {
	lst := vm.GetValidatorList(ctx)
	tmp := append(lst.Standby, lst.Candidates...)
	return append(lst.Oncall, tmp...)
}

func (vm ValidatorManager) GetCommittingValidators(ctx sdk.Context) []linotypes.AccountKey {
	lst := vm.GetValidatorList(ctx)
	return append(lst.Oncall, lst.Standby...)
}

func (vm ValidatorManager) GetValidatorList(ctx sdk.Context) *model.ValidatorList {
	return vm.storage.GetValidatorList(ctx)
}

func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, lst *model.ValidatorList) {
	showed := make(map[linotypes.AccountKey]bool)
	for _, oncall := range lst.Oncall {
		if showed[oncall] == true {
			debug.PrintStack()
		}
		showed[oncall] = true
	}
	vm.storage.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) GetElectionVoteList(ctx sdk.Context,
	accKey linotypes.AccountKey) *model.ElectionVoteList {
	return vm.storage.GetElectionVoteList(ctx, accKey)
}

func (vm ValidatorManager) getPrevVotes(ctx sdk.Context, user linotypes.AccountKey) linotypes.Coin {
	val, err := vm.storage.GetValidator(ctx, user)
	if err != nil {
		return linotypes.NewCoinFromInt64(0)
	}
	return val.ReceivedVotes
}

func (vm ValidatorManager) GetCommittingValidatorVoteStatus(ctx sdk.Context) []model.ReceivedVotesStatus {
	lst := vm.GetCommittingValidators(ctx)
	res := []model.ReceivedVotesStatus{}
	for _, name := range lst {
		val, err := vm.storage.GetValidator(ctx, name)
		if err != nil {
			panic(err)
		}
		res = append(res, model.ReceivedVotesStatus{
			ValidatorName: name,
			ReceivedVotes: val.ReceivedVotes,
		})
	}
	return res
}

// ExportToFile -
func (vm ValidatorManager) ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	state := &model.ValidatorTablesIR{
		Version: exportVersion,
	}
	substores := vm.storage.StoreMap(ctx)

	// export validators
	substores[string(model.ValidatorSubstore)].Iterate(func(key []byte, val interface{}) bool {
		validator := val.(*model.Validator)
		state.Validators = append(state.Validators, model.ValidatorIR{
			ABCIValidator: model.ABCIValidatorIR{
				Address: validator.ABCIValidator.Address,
				Power:   validator.ABCIValidator.Power,
			},
			PubKey:         model.NewABCIPubKeyIRFromTM(validator.PubKey),
			Username:       validator.Username,
			ReceivedVotes:  validator.ReceivedVotes,
			HasRevoked:     validator.HasRevoked,
			AbsentCommit:   validator.AbsentCommit,
			ProducedBlocks: validator.ProducedBlocks,
			Link:           validator.Link,
		})
		return false
	})

	// export votes
	substores[string(model.ElectionVoteListSubstore)].Iterate(func(key []byte, val interface{}) bool {
		user := linotypes.AccountKey(key)
		votelist := val.(*model.ElectionVoteList)
		votesIR := make([]model.ElectionVoteIR, 0)
		for _, vote := range votelist.ElectionVotes {
			votesIR = append(votesIR, model.ElectionVoteIR(vote))
		}
		state.Votes = append(state.Votes, model.ElectionVoteListIR{
			Username:      user,
			ElectionVotes: votesIR,
		})
		return false
	})

	// export validator list.
	substores[string(model.ValidatorListSubstore)].Iterate(func(key []byte, val interface{}) bool {
		lst := val.(*model.ValidatorList)
		state.List = model.ValidatorListIR(*lst)
		return false
	})

	return utils.Save(filepath, cdc, state)
}

// ImportFromFile import state from file.
func (vs ValidatorManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.ValidatorTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.ValidatorTablesIR)

	if table.Version != importVersion {
		return fmt.Errorf("unsupported import version: %d", table.Version)
	}

	// import validators.
	for _, val := range table.Validators {
		vs.storage.SetValidator(ctx, val.Username, &model.Validator{
			ABCIValidator: abci.Validator{
				Address: val.ABCIValidator.Address,
				Power:   val.ABCIValidator.Power,
			},
			PubKey:         val.PubKey.ToTM(),
			Username:       val.Username,
			ReceivedVotes:  val.ReceivedVotes,
			HasRevoked:     val.HasRevoked,
			AbsentCommit:   val.AbsentCommit,
			ProducedBlocks: val.ProducedBlocks,
			Link:           val.Link,
		})
	}

	// import votes.
	for _, vote := range table.Votes {
		votes := make([]model.ElectionVote, 0)
		for _, v := range vote.ElectionVotes {
			votes = append(votes, model.ElectionVote(v))
		}
		vs.storage.SetElectionVoteList(ctx, vote.Username, &model.ElectionVoteList{
			ElectionVotes: votes,
		})
	}

	// import validator list
	validatorList := model.ValidatorList(table.List)
	vs.storage.SetValidatorList(ctx, &validatorList)
	return nil
}
