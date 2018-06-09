package validator

import (
	"encoding/hex"
	"math"
	"reflect"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/x/validator/model"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
)

type ValidatorManager struct {
	storage     model.ValidatorStorage `json:"validator_storage"`
	paramHolder param.ParamHolder      `json:"param_holder"`
}

func NewValidatorManager(key sdk.StoreKey, holder param.ParamHolder) ValidatorManager {
	return ValidatorManager{
		storage:     model.NewValidatorStorage(key),
		paramHolder: holder,
	}
}

func (vm ValidatorManager) InitGenesis(ctx sdk.Context) error {
	if err := vm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) IsValidatorExist(ctx sdk.Context, accKey types.AccountKey) bool {
	infoByte, _ := vm.storage.GetValidator(ctx, accKey)
	return infoByte != nil
}

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

	if FindAccountInList(username, lst.OncallValidators) != -1 {
		return false
	}

	// pass if it's not in all validator list
	// reject if the remaining coins are less than min deposit requirement
	res := validator.Deposit.Minus(coin)
	return res.IsGTE(param.ValidatorMinCommitingDeposit)
}

func (vm ValidatorManager) IsBalancedAccount(
	ctx sdk.Context, accKey types.AccountKey, votingDeposit types.Coin) bool {
	commitingDeposit, err := vm.GetValidatorDeposit(ctx, accKey)
	if err != nil {
		return false
	}
	return votingDeposit.IsGTE(commitingDeposit)
}

func (vm ValidatorManager) GetUpdateValidatorList(ctx sdk.Context) ([]abci.Validator, sdk.Error) {
	validatorList, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return nil, err
	}
	ABCIValList := []abci.Validator{}
	for _, preValidator := range validatorList.PreBlockValidators {
		// set power to 0 if a previous validator not in oncall list anymore
		if FindAccountInList(preValidator, validatorList.OncallValidators) == -1 {
			validator, err := vm.storage.GetValidator(ctx, preValidator)
			if err != nil {
				return nil, err
			}
			if validator.Deposit.IsZero() {
				vm.storage.DeleteValidator(ctx, validator.Username)
			}

			validator.ABCIValidator.Power = 0
			ABCIValList = append(ABCIValList, validator.ABCIValidator)
		}
	}

	for _, curValidator := range validatorList.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, curValidator)
		if err != nil {
			return nil, err
		}
		ABCIValList = append(ABCIValList, validator.ABCIValidator)
	}
	return ABCIValList, nil
}

func (vm ValidatorManager) GetValidatorList(ctx sdk.Context) (*model.ValidatorList, sdk.Error) {
	return vm.storage.GetValidatorList(ctx)
}

func (vm ValidatorManager) GetValidatorDeposit(ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	validator, err := vm.storage.GetValidator(ctx, accKey)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return validator.Deposit, nil
}

func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, lst *model.ValidatorList) sdk.Error {
	return vm.storage.SetValidatorList(ctx, lst)
}

func (vm ValidatorManager) UpdateAbsentValidator(ctx sdk.Context, absentValidators []int32) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}

	// add number of produced blocks for all validators
	for _, curValidator := range lst.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, curValidator)
		if err != nil {
			return err
		}
		validator.ProducedBlocks += 1
		if err := vm.storage.SetValidator(ctx, curValidator, validator); err != nil {
			return err
		}

	}

	// sort the oncall validator list according to their address
	var addrs []string
	addrToName := make(map[string]types.AccountKey)

	for _, validatorName := range lst.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}

		keyBytes := validator.ABCIValidator.GetPubKey()
		pubKey, cerr := crypto.PubKeyFromBytes(keyBytes)
		if cerr != nil {
			return ErrGetPubKeyFailed()
		}
		addr := hex.EncodeToString(pubKey.Address())
		addrs = append(addrs, addr)
		addrToName[addr] = validatorName
	}

	sort.Strings(addrs)

	for _, idx := range absentValidators {
		if idx >= int32(len(addrs)) {
			return ErrAbsentValidatorNotCorrect()
		}
		validatorName := addrToName[addrs[idx]]
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}
		validator.AbsentCommit += 1
		validator.ProducedBlocks -= 1

		if err := vm.storage.SetValidator(ctx, validatorName, validator); err != nil {
			return err
		}
	}

	return nil
}

func (vm ValidatorManager) PunishOncallValidator(
	ctx sdk.Context, username types.AccountKey, penalty types.Coin, willFire bool) (types.Coin, sdk.Error) {
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

	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return actualPenalty, err
	}

	// remove this validator if its remaining deposit is not enough
	// OR, we explicitly want to fire this validator
	// all deposit will be added back to inflation pool
	if willFire || !validator.Deposit.IsGTE(param.ValidatorMinCommitingDeposit) {
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

func (vm ValidatorManager) FireIncompetentValidator(
	ctx sdk.Context, ByzantineValidators []abci.Evidence) (types.Coin, sdk.Error) {
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
		if validator.AbsentCommit > param.AbsentCommitLimitation {
			actualPenalty, err := vm.PunishOncallValidator(
				ctx, validator.Username, param.PenaltyMissCommit, true)
			if err != nil {
				return totalPenalty, err
			}
			totalPenalty = totalPenalty.Plus(actualPenalty)
		}

		for _, evidence := range ByzantineValidators {
			if reflect.DeepEqual(validator.ABCIValidator.PubKey, evidence.PubKey) {
				actualPenalty, err := vm.PunishOncallValidator(
					ctx, validator.Username, param.PenaltyByzantine, true)
				if err != nil {
					return totalPenalty, err
				}
				totalPenalty = totalPenalty.Plus(actualPenalty)
			}
		}
	}

	return totalPenalty, nil
}

func (vm ValidatorManager) PunishValidatorsDidntVote(
	ctx sdk.Context, penaltyList []types.AccountKey) (types.Coin, sdk.Error) {
	totalPenalty := types.NewCoinFromInt64(0)
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return totalPenalty, err
	}
	// punish these validators who didn't vote
	for _, validator := range penaltyList {
		actualPenalty, err := vm.PunishOncallValidator(ctx, validator, param.PenaltyMissVote, false)
		if err != nil {
			return totalPenalty, err
		}
		totalPenalty = totalPenalty.Plus(actualPenalty)
	}

	return totalPenalty, nil
}

func (vm ValidatorManager) RegisterValidator(
	ctx sdk.Context, username types.AccountKey, pubKey []byte, coin types.Coin, link string) sdk.Error {
	// check validator minimum commiting deposit requirement
	param, err := vm.paramHolder.GetValidatorParam(ctx)
	if err != nil {
		return err
	}
	if !coin.IsGTE(param.ValidatorMinCommitingDeposit) {
		return ErrCommitingDepositNotEnough()
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
		if reflect.DeepEqual(validator.ABCIValidator.PubKey, pubKey) {
			return ErrPubKeyHasBeenRegistered()
		}
	}
	curValidator := &model.Validator{
		ABCIValidator: abci.Validator{PubKey: pubKey, Power: 1000},
		Username:      username,
		Deposit:       coin,
		Link:          link,
	}

	if err := vm.storage.SetValidator(ctx, username, curValidator); err != nil {
		return err
	}
	return nil
}

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

// this method won't check if it is a legal withdraw, caller should check by itself
func (vm ValidatorManager) ValidatorWithdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return ErrNoCoinToWithdraw()
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

// try to join the oncall validator list, the action will success if either
// 1. the validator list is not full
// or 2. someone in the validator list has a lower power than current validator
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
	if !curValidator.Deposit.IsGTE(param.ValidatorMinCommitingDeposit) {
		return ErrCommitingDepositNotEnough()
	}

	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	defer vm.updateLowestValidator(ctx)
	// has alreay in the oncall validator list
	if FindAccountInList(username, lst.OncallValidators) != -1 {
		return nil
	}

	// add to all validators list if not in the list
	if FindAccountInList(username, lst.AllValidators) == -1 {
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

// remove the user from both oncall and allValidators lists
func (vm ValidatorManager) RemoveValidatorFromAllLists(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, err := vm.storage.GetValidatorList(ctx)
	if err != nil {
		return err
	}
	if FindAccountInList(username, lst.AllValidators) == -1 {
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
	for idx, username := range users {
		if me == username {
			users = append(users[:idx], users[idx+1:]...)
		}
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
		if FindAccountInList(validatorName, lst.OncallValidators) == -1 &&
			validator.Deposit.IsGT(bestCandidatePower) {
			bestCandidate = validator.Username
			bestCandidatePower = validator.Deposit
		}
	}
	return bestCandidate, nil

}

func FindAccountInList(me types.AccountKey, lst []types.AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}
