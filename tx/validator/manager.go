package validator

import (
	"encoding/json"
	"math"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
)

var ValidatorAccountPrefix = []byte("ValidatorAccountInfo/")
var ValidatorListPrefixWithKey = []byte("ValidatorList/ValidatorListKey")

// Validator Manager implements types.AccountManager
type ValidatorManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewValidatorManager returns a new ValidatorManager
func NewValidatorMananger(key sdk.StoreKey) ValidatorManager {
	cdc := wire.NewCodec()

	vm := ValidatorManager{
		key: key,
		cdc: cdc,
	}

	return vm
}

func (vm ValidatorManager) InitGenesis(ctx sdk.Context, data json.RawMessage) error {
	lst := &ValidatorList{
		LowestPower: types.Coin{0},
	}

	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) IsValidatorExist(ctx sdk.Context, accKey acc.AccountKey) bool {
	store := ctx.KVStore(vm.key)
	if infoByte := store.Get(GetValidatorKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

func (vm ValidatorManager) GetValidator(ctx sdk.Context, accKey acc.AccountKey) (*Validator, sdk.Error) {
	store := ctx.KVStore(vm.key)
	validatorByte := store.Get(GetValidatorKey(accKey))
	if validatorByte == nil {
		return nil, ErrGetValidator()
	}
	validator := new(Validator)
	if err := vm.cdc.UnmarshalJSON(validatorByte, validator); err != nil {
		return nil, ErrValidatorUnmarshalError(err)
	}
	return validator, nil
}

func (vm ValidatorManager) SetValidator(ctx sdk.Context, accKey acc.AccountKey, validator *Validator) sdk.Error {
	store := ctx.KVStore(vm.key)
	validatorByte, err := vm.cdc.MarshalJSON(*validator)
	if err != nil {
		return ErrValidatorMarshalError(err)
	}
	store.Set(GetValidatorKey(accKey), validatorByte)
	return nil
}

func (vm ValidatorManager) GetValidatorList(ctx sdk.Context) (*ValidatorList, sdk.Error) {
	store := ctx.KVStore(vm.key)
	listByte := store.Get(GetValidatorListKey())
	if listByte == nil {
		return nil, ErrGetValidatorList()
	}
	lst := new(ValidatorList)
	if err := vm.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrValidatorUnmarshalError(err)
	}
	return lst, nil
}

func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vm.key)
	listByte, err := vm.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrSetValidatorList()
	}
	store.Set(GetValidatorListKey(), listByte)
	return nil
}

func (vm ValidatorManager) GetOncallValList(ctx sdk.Context) ([]Validator, sdk.Error) {
	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return nil, getListErr
	}

	oncallList := make([]Validator, len(lst.OncallValidators))
	for i, validatorName := range lst.OncallValidators {
		validator, err := vm.GetValidator(ctx, validatorName)
		if err != nil {
			return nil, err
		}
		oncallList[i] = *validator
	}
	return oncallList, nil
}

func (vm ValidatorManager) UpdateAbsentValidator(ctx sdk.Context, absentValidators []int32) sdk.Error {
	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}

	for _, idx := range absentValidators {
		if idx > int32(len(lst.OncallValidators)) {
			return ErrAbsentValidatorNotCorrect()
		}
		validator, err := vm.GetValidator(ctx, lst.OncallValidators[idx])
		if err != nil {
			return err
		}
		validator.AbsentVote += 1

		if err := vm.SetValidator(ctx, lst.OncallValidators[idx], validator); err != nil {
			return err
		}
	}

	return nil
}

func (vm ValidatorManager) MarkByzantine(ctx sdk.Context, username acc.AccountKey) sdk.Error {
	byzantine, err := vm.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	byzantine.IsByzantine = true
	if err := vm.SetValidator(ctx, username, byzantine); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) FireIncompetentValidator(ctx sdk.Context, ByzantineValidators []abci.Evidence) sdk.Error {
	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}
	fireList := []acc.AccountKey{}

	for _, validatorName := range lst.OncallValidators {
		validator, err := vm.GetValidator(ctx, validatorName)
		if err != nil {
			return err
		}
		if validator.AbsentVote > types.AbsentLimitation {
			fireList = append(fireList, validatorName)
			continue
		}
		for _, evidence := range ByzantineValidators {
			if reflect.DeepEqual(validator.ABCIValidator.PubKey, evidence.PubKey) {
				fireList = append(fireList, validatorName)
				if err := vm.MarkByzantine(ctx, validatorName); err != nil {
					return err
				}
			}
		}
	}
	for _, validatorName := range fireList {
		if err := vm.RemoveValidatorFromAllLists(ctx, validatorName); err != nil {
			return err
		}
	}

	return nil
}

func (vm ValidatorManager) AddToCandidatePool(ctx sdk.Context, username acc.AccountKey) sdk.Error {
	curValidator, getErr := vm.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}

	valRegisterFee, err := types.LinoToCoin(types.LNO(sdk.NewRat(1000)))
	if err != nil {
		return sdk.ErrInvalidCoins("invalid register fee")
	}
	// check minimum requirements
	if !curValidator.Deposit.IsGTE(valRegisterFee) {
		return ErrRegisterFeeNotEnough()
	}

	lst, getListErr := vm.GetValidatorList(ctx)

	if getListErr != nil {
		return getListErr
	}

	// has alreay in the validator list
	if findAccountInList(username, lst.AllValidators) != -1 {
		return nil
	}

	lst.AllValidators = append(lst.AllValidators, username)
	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

// try to join the oncall validator list, the action will success if either
// 1. the validator list is not full
// or 2. someone in the validator list has a lower power than current validator
// return a boolean to indicate if the user has became an oncall validator
// Also, set WithdrawAvailableAt to be infinite if become an oncall validator
func (vm ValidatorManager) TryBecomeOncallValidator(ctx sdk.Context, username acc.AccountKey) sdk.Error {
	curValidator, getErr := vm.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}

	valRegisterFee, err := types.LinoToCoin(types.LNO(sdk.NewRat(1000)))
	if err != nil {
		return sdk.ErrInvalidCoins("invalid register fee")
	}

	// check minimum requirements
	if !curValidator.Deposit.IsGTE(valRegisterFee) {
		return ErrRegisterFeeNotEnough()
	}

	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}
	defer vm.updateLowestValidator(ctx)
	// has alreay in the oncall validator list
	if findAccountInList(username, lst.OncallValidators) != -1 {
		return nil
	}
	// add to list directly if validator list is not full
	if len(lst.OncallValidators) < types.ValidatorListSize {
		lst.OncallValidators = append(lst.OncallValidators, curValidator.Username)
		curValidator.WithdrawAvailableAt = types.InfiniteFreezingPeriod
		//vm.updateLowestValidator(ctx)
	} else if curValidator.ABCIValidator.Power > lst.LowestPower.Amount {
		// replace the validator with lowest power
		for idx, validatorKey := range lst.OncallValidators {
			validator, getErr := vm.GetValidator(ctx, validatorKey)
			if getErr != nil {
				return getErr
			}
			if validator.Username == lst.LowestValidator {
				lst.OncallValidators[idx] = curValidator.Username
			}
		}
		curValidator.WithdrawAvailableAt = types.InfiniteFreezingPeriod
		//vm.updateLowestValidator(ctx)
	}

	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	if err := vm.SetValidator(ctx, curValidator.Username, curValidator); err != nil {
		return err
	}

	return nil
}

// remove the user from both oncall and allValidators lists
// Also, set WithdrawAvailableAt to a freezing period
func (vm ValidatorManager) RemoveValidatorFromAllLists(ctx sdk.Context, username acc.AccountKey) sdk.Error {
	curValidator, getErr := vm.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}

	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}

	if findAccountInList(username, lst.AllValidators) == -1 {
		return ErrNotInTheList()
	}

	lst.AllValidators = remove(username, lst.AllValidators)
	lst.OncallValidators = remove(username, lst.OncallValidators)

	if curValidator.IsByzantine {
		//TODO return deposit to pool?
		curValidator.WithdrawAvailableAt = types.Height(ctx.BlockHeight() + int64(types.ValidatorWithdrawFreezingPeriod))
	} else {
		curValidator.WithdrawAvailableAt = types.Height(ctx.BlockHeight() + int64(types.ValidatorWithdrawFreezingPeriod))
	}

	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	if err := vm.SetValidator(ctx, curValidator.Username, curValidator); err != nil {
		return err
	}

	vm.updateLowestValidator(ctx)
	bestCandidate, findErr := vm.getBestCandidate(ctx, lst)
	if findErr != nil {
		return findErr
	}

	if bestCandidate != acc.AccountKey("") {
		if joinErr := vm.TryBecomeOncallValidator(ctx, bestCandidate); joinErr != nil {
			return joinErr
		}
	}
	return nil
}

func remove(me acc.AccountKey, users []acc.AccountKey) []acc.AccountKey {
	for idx, username := range users {
		if me == username {
			users = append(users[:idx], users[idx+1:]...)
		}
	}
	return users
}

func (vm ValidatorManager) updateLowestValidator(ctx sdk.Context) {
	lst, _ := vm.GetValidatorList(ctx)
	newLowestPower := int64(math.MaxInt64)
	newLowestValidator := acc.AccountKey("")

	for _, validatorKey := range lst.OncallValidators {
		validator, _ := vm.GetValidator(ctx, validatorKey)
		if validator.ABCIValidator.Power < newLowestPower {
			newLowestPower = validator.ABCIValidator.Power
			newLowestValidator = validator.Username
		}
	}
	// set the new lowest power
	lst.LowestPower = types.NewCoin(newLowestPower)
	lst.LowestValidator = newLowestValidator

	vm.SetValidatorList(ctx, lst)
}

// find the person has the biggest power among people in the allValidators lists
// but not in the oncall validator list
func (vm ValidatorManager) getBestCandidate(ctx sdk.Context, lst *ValidatorList) (acc.AccountKey, sdk.Error) {
	bestCandidate := acc.AccountKey("")
	bestCandidatePower := int64(0)

	for i, validatorName := range lst.AllValidators {
		validator, getErr := vm.GetValidator(ctx, lst.AllValidators[i])
		if getErr != nil {
			return bestCandidate, getErr
		}
		// not in the oncall list and has a larger power
		if findAccountInList(validatorName, lst.OncallValidators) == -1 &&
			validator.ABCIValidator.Power > bestCandidatePower {
			bestCandidate = validator.Username
			bestCandidatePower = validator.ABCIValidator.Power
		}
	}
	return bestCandidate, nil

}

func GetValidatorKey(accKey acc.AccountKey) []byte {
	return append(ValidatorAccountPrefix, accKey...)
}

func GetValidatorListKey() []byte {
	return ValidatorListPrefixWithKey
}

func findAccountInList(me acc.AccountKey, lst []acc.AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}
