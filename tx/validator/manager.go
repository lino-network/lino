package validator

import (
	"math"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
)

// validator manager is the proxy for all storage structs defined above
type ValidatorManager struct {
	storage *model.ValidatorStorage `json:"validator_storage"`
}

// create NewValidatorManager
func NewValidatorManager(key sdk.StoreKey) *ValidatorManager {
	return &ValidatorManager{
		storage: model.NewValidatorStorage(key),
	}
}

func (vm *ValidatorManager) GetUpdateValidatorList(ctx sdk.Context) ([]abci.Validator, sdk.Error) {
	curOncallList, err := vm.GetOncallValidatorList(ctx)
	if err != nil {
		return nil, err
	}
	ABCIValList := []abci.Validator{}
	preBlockValidators := GetPreBlockValidators(ctx)
	for _, preValidator := range preBlockValidators {
		if FindAccountInList(preValidator, curOncallList) == -1 {
			validator, getErr := vm.storage.GetValidator(ctx, preValidator)
			if getErr != nil {
				return nil, err
			}
			validator.ABCIValidator.Power = 0
			ABCIValList = append(ABCIValList, validator.ABCIValidator)
		}
	}

	for _, curValidator := range curOncallList {
		validator, getErr := vm.storage.GetValidator(ctx, curValidator)
		if getErr != nil {
			return nil, err
		}
		ABCIValList = append(ABCIValList, validator.ABCIValidator)
	}
	return ABCIValList, nil
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

func (vm ValidatorManager) IsLegalWithdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin) bool {
	validator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return false
	}

	// reject if withdraw is less than minimum withdraw
	if !coin.IsGTE(types.ValidatorMinWithdraw) {
		return false
	}

	// reject if this is an oncall validator
	lst, getErr := vm.storage.GetValidatorList(ctx)
	if getErr != nil {
		return false
	}

	if FindAccountInList(username, lst.OncallValidators) != -1 {
		return false
	}
	//reject if the remaining coins are less than min deposit requirement
	res := validator.Deposit.Minus(coin)
	return res.IsGTE(types.ValidatorMinCommitingDeposit)
}

func (vm ValidatorManager) GetOncallValidatorList(ctx sdk.Context) ([]types.AccountKey, sdk.Error) {
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return nil, getListErr
	}
	return lst.OncallValidators, nil
}

func (vm ValidatorManager) GetAllValidatorList(ctx sdk.Context) ([]types.AccountKey, sdk.Error) {
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return nil, getListErr
	}
	return lst.AllValidators, nil
}

func (vm ValidatorManager) UpdateAbsentValidator(ctx sdk.Context, absentValidators []int32) sdk.Error {
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}

	for _, idx := range absentValidators {
		if idx > int32(len(lst.OncallValidators)) {
			return ErrAbsentValidatorNotCorrect()
		}
		validator, err := vm.storage.GetValidator(ctx, lst.OncallValidators[idx])
		if err != nil {
			return err
		}
		validator.AbsentCommit += 1

		if err := vm.storage.SetValidator(ctx, lst.OncallValidators[idx], validator); err != nil {
			return err
		}
	}

	return nil
}

func (vm ValidatorManager) PunishOncallValidator(ctx sdk.Context, username types.AccountKey, penalty types.Coin, gm global.GlobalManager, willFire bool) sdk.Error {
	validator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}
	validator.Deposit = validator.Deposit.Minus(penalty)
	if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
		return err
	}

	// add coins back to inflation pool
	if err := gm.AddToValidatorInflationPool(ctx, penalty); err != nil {
		return err
	}
	// remove this validator if its remaining deposit is not enough
	// OR, we explicitly want to fire this validator
	// it is user's responsibility to do future withdraw/deposit
	if willFire || !validator.Deposit.IsGTE(types.ValidatorMinCommitingDeposit) {
		if err := vm.RemoveValidatorFromAllLists(ctx, validator.Username); err != nil {
			return err
		}
		return nil
	}
	if err := vm.AdjustValidatorList(ctx); err != nil {
		return err
	}
	return nil
}
func (vm ValidatorManager) FireIncompetentValidator(ctx sdk.Context, ByzantineValidators []abci.Evidence, gm global.GlobalManager) sdk.Error {
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}

	for _, validatorName := range lst.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}
		if validator.AbsentCommit > types.AbsentCommitLimitation {
			vm.PunishOncallValidator(ctx, validator.Username, types.PenaltyMissVote, gm, true)
			continue
		}

		for _, evidence := range ByzantineValidators {
			if reflect.DeepEqual(validator.ABCIValidator.PubKey, evidence.PubKey) {
				vm.PunishOncallValidator(ctx, validator.Username, types.PenaltyByzantine, gm, true)
			}
		}
	}

	return nil
}

func (vm ValidatorManager) RegisterValidator(ctx sdk.Context, username types.AccountKey, pubKey []byte, coin types.Coin) sdk.Error {
	// check validator minimum commiting deposit requirement
	if !coin.IsGTE(types.ValidatorMinCommitingDeposit) {
		return ErrCommitingDepositNotEnough()
	}

	curValidator := &model.Validator{
		ABCIValidator: abci.Validator{PubKey: pubKey, Power: 1000},
		Username:      username,
		Deposit:       coin,
	}

	if setErr := vm.storage.SetValidator(ctx, username, curValidator); setErr != nil {
		return setErr
	}
	return nil
}

func (vm ValidatorManager) Deposit(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	validator, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}
	validator.Deposit = validator.Deposit.Plus(coin)
	if setErr := vm.storage.SetValidator(ctx, username, validator); setErr != nil {
		return setErr
	}
	return nil
}

// this method won't check if it is a legal withdraw, caller should check by itself
func (vm ValidatorManager) ValidatorWithdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	if coin.IsZero() {
		return ErrNoCoinToWithdraw()
	}
	validator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}
	validator.Deposit = validator.Deposit.Minus(coin)
	if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) ValidatorWithdrawAll(ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	validator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return types.NewCoin(0), getErr
	}
	if err := vm.ValidatorWithdraw(ctx, username, validator.Deposit); err != nil {
		return types.NewCoin(0), err
	}
	return validator.Deposit, nil
}

// try to join the oncall validator list, the action will success if either
// 1. the validator list is not full
// or 2. someone in the validator list has a lower power than current validator
// return a boolean to indicate if the user has became an oncall validator
// Also, set WithdrawAvailableAt to be infinite if become an oncall validator
func (vm ValidatorManager) TryBecomeOncallValidator(ctx sdk.Context, username types.AccountKey) sdk.Error {
	curValidator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}

	// check minimum requirements
	if !curValidator.Deposit.IsGTE(types.ValidatorMinCommitingDeposit) {
		return ErrCommitingDepositNotEnough()
	}

	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
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
	if len(lst.OncallValidators) < types.ValidatorListSize {
		lst.OncallValidators = append(lst.OncallValidators, curValidator.Username)
		//vm.updateLowestValidator(ctx)
	} else if curValidator.Deposit.Amount > lst.LowestPower.Amount {
		// replace the validator with lowest power
		for idx, validatorKey := range lst.OncallValidators {
			validator, getErr := vm.storage.GetValidator(ctx, validatorKey)
			if getErr != nil {
				return getErr
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
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
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
	vm.updateLowestValidator(ctx)
	bestCandidate, findErr := vm.getBestCandidate(ctx)
	if findErr != nil {
		return findErr
	}

	if bestCandidate != types.AccountKey("") {
		if joinErr := vm.TryBecomeOncallValidator(ctx, bestCandidate); joinErr != nil {
			return joinErr
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

func (vm ValidatorManager) updateLowestValidator(ctx sdk.Context) {
	lst, _ := vm.storage.GetValidatorList(ctx)
	newLowestPower := types.NewCoin(math.MaxInt64)
	newLowestValidator := types.AccountKey("")

	for _, validatorKey := range lst.OncallValidators {
		validator, _ := vm.storage.GetValidator(ctx, validatorKey)
		if validator.Deposit.Amount < newLowestPower.Amount {
			newLowestPower = validator.Deposit
			newLowestValidator = validator.Username
		}
	}
	// set the new lowest power
	lst.LowestPower = newLowestPower
	lst.LowestValidator = newLowestValidator

	vm.storage.SetValidatorList(ctx, lst)
}

// find the person has the biggest power among people in the allValidators lists
// but not in the oncall validator list
func (vm ValidatorManager) getBestCandidate(ctx sdk.Context) (types.AccountKey, sdk.Error) {
	bestCandidate := types.AccountKey("")
	bestCandidatePower := types.NewCoin(0)

	lst, getErr := vm.storage.GetValidatorList(ctx)
	if getErr != nil {
		return bestCandidate, getErr
	}

	for i, validatorName := range lst.AllValidators {
		validator, getErr := vm.storage.GetValidator(ctx, lst.AllValidators[i])
		if getErr != nil {
			return bestCandidate, getErr
		}
		// not in the oncall list and has a larger power
		if FindAccountInList(validatorName, lst.OncallValidators) == -1 &&
			validator.Deposit.Amount > bestCandidatePower.Amount {
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
