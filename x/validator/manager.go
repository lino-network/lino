package validator

import (
	"math"
	"reflect"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/validator/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ValidatorManager - validator manager
type ValidatorManager struct {
	storage     model.ValidatorStorage
	paramHolder param.ParamHolder
}

func NewValidatorManager(key sdk.StoreKey, holder param.ParamHolder) ValidatorManager {
	return ValidatorManager{
		storage:     model.NewValidatorStorage(key),
		paramHolder: holder,
	}
}

// InitGenesis - initialize KVStore
func (vm ValidatorManager) InitGenesis(ctx sdk.Context) error {
	if err := vm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

// DoesValidatorExist - check if validator exists in KVStore or not
func (vm ValidatorManager) DoesValidatorExist(ctx sdk.Context, accKey types.AccountKey) bool {
	return vm.storage.DoesValidatorExist(ctx, accKey)
}

// IsLegalWithdraw - check if withdraw is legal or not
func (vm ValidatorManager) IsLegalWithdraw(
	ctx sdk.Context, username types.AccountKey, coin types.Coin) bool {
	validator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return false
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return false
	}

	// reject if withdraw is less than minimum withdraw
	if !coin.IsGTE(param.ValidatorMinWithdraw) {
		return false
	}

	// reject if this is an oncall validator
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return false
	}

	if types.FindAccountInList(username, lst.OncallValidators) != -1 {
		return false
	}

	// pass if it's not in all validator list
	// reject if the remaining coins are less than min deposit requirement
	res := validator.Deposit.Minus(coin)
	return res.IsGTE(param.ValidatorMinCommittingDeposit)
}

// IsBalancedAccount - make sure voting deposit is much than committing (validator) deposit
func (vm ValidatorManager) IsBalancedAccount(
	ctx sdk.Context, accKey types.AccountKey, votingDeposit types.Coin) bool {
	commitingDeposit, err := vm.GetValidatorDeposit(ctx, accKey)
	if err != nil {
		return false
	}
	return votingDeposit.IsGTE(commitingDeposit)
}

// GetInitValidators return all validators in state.
// XXX(yumin): This is intended to be used only in initChainer
// TODO(yumin): add test coverage.
func (vm ValidatorManager) GetInitValidators(ctx sdk.Context) ([]abci.ValidatorUpdate, sdk.Error) {
	validatorList, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return nil, err
	}
	updates := []abci.ValidatorUpdate{}
	for _, curValidator := range validatorList.OncallValidators {
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
	for _, preValidator := range validatorList.PreBlockValidators {
		// set power to 0 if a previous validator not in oncall list anymore
		if types.FindAccountInList(preValidator, validatorList.OncallValidators) == -1 {
			validator, err := vm.storage.GetValidator(ctx, preValidator)
			if err != nil {
				return nil, err
			}
			if validator.Deposit.IsZero() {
				vm.storage.DeleteValidator(ctx, validator.Username)
			}
			updates = append(updates, abci.ValidatorUpdate{
				PubKey: tmtypes.TM2PB.PubKey(validator.PubKey),
				Power:  0,
			})
		}
	}

	for _, curValidator := range validatorList.OncallValidators {
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

// GetValidatorList - get validator list from KV Store
func (vm ValidatorManager) GetValidatorList(ctx sdk.Context) (*model.ValidatorList, sdk.Error) {
	return vm.storage.GetValidatorList(ctx)
}

// GetValidatorDeposit - get validator deposit
func (vm ValidatorManager) GetValidatorDeposit(ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	validator, err := vm.storage.GetValidator(ctx, accKey)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return validator.Deposit, nil
}

// SetValidatorList - set validator list
func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, lst *model.ValidatorList) sdk.Error {
	return vm.storage.SetValidatorList(ctx, lst)
}

// UpdateSigningStats - based on info in beginBlocker, record last block singing info
func (vm ValidatorManager) UpdateSigningStats(
	ctx sdk.Context, voteInfos []abci.VoteInfo) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	// map address to whether that validator has signed.
	addressSigned := make(map[string]bool)
	for _, voteInfo := range voteInfos {
		addressSigned[string(voteInfo.Validator.Address)] = voteInfo.SignedLastBlock
	}

	// go through oncall validator list to get all address and name mapping
	for _, curValidator := range lst.OncallValidators {
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

// PunishOncallValidator - punish oncall validator if 1) byzantine or 2) missing blocks reach limiation
func (vm ValidatorManager) PunishOncallValidator(
	ctx sdk.Context, username types.AccountKey, penalty types.Coin, punishType types.PunishType) (types.Coin, sdk.Error) {
	actualPenalty := penalty
	validator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return actualPenalty, err
	}

	if penalty.IsGT(validator.Deposit) {
		actualPenalty = validator.Deposit
		validator.Deposit = types.NewCoinFromInt64(0)
	} else {
		validator.Deposit = validator.Deposit.Minus(penalty)
	}

	if punishType == types.PunishAbsentCommit {
		validator.AbsentCommit = 0
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return actualPenalty, err
	}

	// remove this validator if its remaining deposit is not enough
	// OR, we explicitly want to fire this validator
	// all deposit will be added back to inflation pool
	if punishType == types.PunishByzantine || !validator.Deposit.IsGTE(param.ValidatorMinCommittingDeposit) {
		if err := vm.RemoveValidatorFromAllLists(ctx, validator.Username); err != nil {
			return actualPenalty, err
		}
		actualPenalty = actualPenalty.Plus(validator.Deposit)
		validator.Deposit = types.NewCoinFromInt64(0)
	}

	if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
		return actualPenalty, err
	}

	if err := vm.AdjustValidatorList(ctx); err != nil {
		return actualPenalty, err
	}
	return actualPenalty, nil
}

// FireIncompetentValidator - fire oncall validator if 1) deposit insufficient 2) byzantine
func (vm ValidatorManager) FireIncompetentValidator(
	ctx sdk.Context, byzantineValidators []abci.Evidence) (types.Coin, sdk.Error) {
	totalPenalty := types.NewCoinFromInt64(0)
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return totalPenalty, err
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return totalPenalty, err
	}

	for _, validatorName := range lst.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return totalPenalty, err
		}

		for _, evidence := range byzantineValidators {
			if reflect.DeepEqual(validator.ABCIValidator.Address, evidence.Validator.Address) {
				actualPenalty, err := vm.PunishOncallValidator(
					ctx, validator.Username, param.PenaltyByzantine, types.PunishByzantine)
				if err != nil {
					return totalPenalty, err
				}
				totalPenalty = totalPenalty.Plus(actualPenalty)
				break
			}
		}

		if validator.AbsentCommit > param.AbsentCommitLimitation {
			actualPenalty, err := vm.PunishOncallValidator(
				ctx, validator.Username, param.PenaltyMissCommit, types.PunishAbsentCommit)
			if err != nil {
				return totalPenalty, err
			}

			totalPenalty = totalPenalty.Plus(actualPenalty)
		}
	}

	return totalPenalty, nil
}

// PunishValidatorsDidntVote - validators are required to vote Protocol Upgrade and Parameter Change proposal
func (vm ValidatorManager) PunishValidatorsDidntVote(
	ctx sdk.Context, penaltyList []types.AccountKey) (types.Coin, sdk.Error) {
	totalPenalty := types.NewCoinFromInt64(0)
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return totalPenalty, err
	}
	// punish these validators who didn't vote
	for _, validator := range penaltyList {
		actualPenalty, err := vm.PunishOncallValidator(ctx, validator, param.PenaltyMissVote, types.PunishDidntVote)
		if err != nil {
			return totalPenalty, err
		}
		totalPenalty = totalPenalty.Plus(actualPenalty)
	}

	return totalPenalty, nil
}

// RegisterValidator - register validator
func (vm ValidatorManager) RegisterValidator(
	ctx sdk.Context, username types.AccountKey, pubKey crypto.PubKey, coin types.Coin, link string) sdk.Error {
	// check validator minimum committing deposit requirement
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}
	if !coin.IsGTE(param.ValidatorMinCommittingDeposit) {
		return ErrInsufficientDeposit()
	}

	// make sure the pub key has not been registered
	lst, err := vm.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	for _, validatorName := range lst.AllValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}
		// XXX(yumin): ABCIValidator no longer has pubkey, changed to address
		if reflect.DeepEqual(validator.ABCIValidator.Address, pubKey.Address().Bytes()) {
			return ErrValidatorPubKeyAlreadyExist()
		}
	}
	// XXX(yumin): const power?
	curValidator := &model.Validator{
		ABCIValidator: abci.Validator{
			Address: pubKey.Address(),
			Power:   types.TendermintValidatorPower,
		},
		PubKey:   pubKey,
		Username: username,
		Deposit:  coin,
		Link:     link,
	}

	if err := vm.storage.SetValidator(ctx, username, curValidator); err != nil {
		return err
	}
	return nil
}

// Deposit - deposit money to validator
func (vm ValidatorManager) Deposit(
	ctx sdk.Context, username types.AccountKey, coin types.Coin, link string) sdk.Error {
	validator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}
	validator.Deposit = validator.Deposit.Plus(coin)
	if len(link) > 0 {
		validator.Link = link
	}
	if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
		return err
	}
	return nil
}

// ValidatorWithdraw - this method won't check if it is a legal withdraw, caller should check by itself
func (vm ValidatorManager) ValidatorWithdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return ErrInvalidCoin()
	}
	validator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}
	validator.Deposit = validator.Deposit.Minus(coin)
	if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
		return err
	}

	return nil
}

// ValidatorWithdrawAll - revoke validator
func (vm ValidatorManager) ValidatorWithdrawAll(ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	validator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := vm.ValidatorWithdraw(ctx, username, validator.Deposit); err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return validator.Deposit, nil
}

// TryBecomeOncallValidator - try to join the oncall validator list, the action will success if either
// 1. the validator list is not full or 2. someone in the validator list has a lower power than current validator
func (vm ValidatorManager) TryBecomeOncallValidator(ctx sdk.Context, username types.AccountKey) sdk.Error {
	curValidator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}
	// check minimum requirements
	if !curValidator.Deposit.IsGTE(param.ValidatorMinCommittingDeposit) {
		return ErrInsufficientDeposit()
	}

	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	defer vm.updateLowestValidator(ctx)
	// has alreay in the oncall validator list
	if types.FindAccountInList(username, lst.OncallValidators) != -1 {
		return nil
	}

	// add to all validators list if not in the list
	if types.FindAccountInList(username, lst.AllValidators) == -1 {
		lst.AllValidators = append(lst.AllValidators, username)
	}

	// add to list directly if validator list is not full
	if int64(len(lst.OncallValidators)) < param.ValidatorListSize {
		lst.OncallValidators = append(lst.OncallValidators, curValidator.Username)
	} else if curValidator.Deposit.IsGT(lst.LowestPower) {
		// replace the validator with lowest power
		for idx, validatorKey := range lst.OncallValidators {
			validator, err := vm.storage.GetValidator(ctx, validatorKey)
			if err != nil {
				return err
			}
			if validator.Username == lst.LowestValidator {
				lst.OncallValidators[idx] = curValidator.Username
			}
		}
	}

	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	if err := vm.storage.SetValidator(ctx, curValidator.Username, curValidator); err != nil {
		return err
	}

	return nil
}

// RemoveValidatorFromAllLists - remove the user from both oncall and allValidators lists
func (vm ValidatorManager) RemoveValidatorFromAllLists(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	if types.FindAccountInList(username, lst.AllValidators) == -1 {
		return nil
	}

	lst.AllValidators = remove(username, lst.AllValidators)
	lst.OncallValidators = remove(username, lst.OncallValidators)

	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	if err := vm.AdjustValidatorList(ctx); err != nil {
		return err
	}

	return nil
}

// if any change happens in oncall validator(remove, punish),
// we should call this function to adjust validator list
func (vm ValidatorManager) AdjustValidatorList(ctx sdk.Context) sdk.Error {
	if err := vm.updateLowestValidator(ctx); err != nil {
		return err
	}
	bestCandidate, err := vm.getBestCandidate(ctx)
	if err != nil {
		return err
	}

	if bestCandidate != types.AccountKey("") {
		if err := vm.TryBecomeOncallValidator(ctx, bestCandidate); err != nil {
			return err
		}
	}
	return nil
}

func remove(me types.AccountKey, users []types.AccountKey) []types.AccountKey {
	idx := 0
	for idx < len(users) {
		username := users[idx]

		if me == username {
			users = append(users[:idx], users[idx+1:]...)
			continue
		}

		idx++
	}

	return users
}

func (vm ValidatorManager) updateLowestValidator(ctx sdk.Context) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	newLowestPower := types.NewCoinFromInt64(math.MaxInt64)
	newLowestValidator := types.AccountKey("")

	for _, validatorKey := range lst.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorKey)
		if err != nil {
			return err
		}

		if newLowestPower.IsGT(validator.Deposit) {
			newLowestPower = validator.Deposit
			newLowestValidator = validator.Username
		}
	}
	// set the new lowest power
	lst.LowestPower = newLowestPower
	lst.LowestValidator = newLowestValidator

	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

// find the person has the biggest power among people in the allValidators lists
// but not in the oncall validator list
func (vm ValidatorManager) getBestCandidate(ctx sdk.Context) (types.AccountKey, sdk.Error) {
	bestCandidate := types.AccountKey("")
	bestCandidatePower := types.NewCoinFromInt64(0)

	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return bestCandidate, err
	}

	for i, validatorName := range lst.AllValidators {
		validator, err := vm.storage.GetValidator(ctx, lst.AllValidators[i])
		if err != nil {
			return bestCandidate, err
		}
		// not in the oncall list and has a larger power
		if types.FindAccountInList(validatorName, lst.OncallValidators) == -1 &&
			validator.Deposit.IsGT(bestCandidatePower) {
			bestCandidate = validator.Username
			bestCandidatePower = validator.Deposit
		}
	}
	return bestCandidate, nil

}

// Export storage state.
func (vm ValidatorManager) Export(ctx sdk.Context) *model.ValidatorTables {
	return vm.storage.Export(ctx)
}

// Import storage state.
func (vm ValidatorManager) Import(ctx sdk.Context, tb *model.ValidatorTablesIR) {
	vm.storage.Import(ctx, tb)
}
