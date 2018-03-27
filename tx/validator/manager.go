package validator

import (
	"fmt"
	"math"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
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

// NewValidatorManager returns a new ValidatorManager that
// uses go-wire to (binary) encode and decode concrete Validator
func NewValidatorMananger(key sdk.StoreKey) ValidatorManager {
	cdc := wire.NewCodec()

	vm := ValidatorManager{
		key: key,
		cdc: cdc,
	}

	return vm
}

// Implements ValidatorManager
func (vm ValidatorManager) Init(ctx sdk.Context) sdk.Error {
	lst := &ValidatorList{
		LowestPower: sdk.Coins{sdk.Coin{Denom: types.Denom, Amount: 0}},
	}

	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

// Implements ValidatorManager
func (vm ValidatorManager) IsValidatorExist(ctx sdk.Context, accKey acc.AccountKey) bool {
	store := ctx.KVStore(vm.key)
	if infoByte := store.Get(ValidatorKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

// Implements ValidatorManager
func (vm ValidatorManager) GetValidator(ctx sdk.Context, accKey acc.AccountKey) (*Validator, sdk.Error) {
	store := ctx.KVStore(vm.key)
	validatorByte := store.Get(ValidatorKey(accKey))
	if validatorByte == nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get validator failed: validator doesn't exist")
	}
	validator := new(Validator)
	if err := vm.cdc.UnmarshalJSON(validatorByte, validator); err != nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get validator failed")
	}
	return validator, nil
}

// Implements ValidatorManager
func (vm ValidatorManager) SetValidator(ctx sdk.Context, accKey acc.AccountKey, validator *Validator) sdk.Error {
	store := ctx.KVStore(vm.key)
	validatorByte, err := vm.cdc.MarshalJSON(*validator)
	if err != nil {
		return ErrValidatorManagerFail("ValidatorManager set validator failed")
	}
	store.Set(ValidatorKey(accKey), validatorByte)
	return nil
}

// Implements ValidatorManager
func (vm ValidatorManager) GetValidatorList(ctx sdk.Context) (*ValidatorList, sdk.Error) {
	store := ctx.KVStore(vm.key)
	listByte := store.Get(ValidatorListKey())
	if listByte == nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get account list failed: account list doesn't exist")
	}
	lst := new(ValidatorList)
	if err := vm.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get account list failed")
	}
	return lst, nil
}

// Implements ValidatorManager
func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vm.key)
	listByte, err := vm.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrValidatorManagerFail("ValidatorManager set account list failed")
	}
	store.Set(ValidatorListKey(), listByte)
	return nil
}

// Implements ValidatorManager
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

// Implements ValidatorManager
func (vm ValidatorManager) UpdateAbsentValidator(ctx sdk.Context, absentValidators []int32) sdk.Error {
	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}
	for _, idx := range absentValidators {
		if idx > int32(len(lst.OncallValidators)) {
			return ErrValidatorManagerFail("absent validator index out of range")
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

// Implements ValidatorManager
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
			}
		}
	}
	fmt.Println(fireList)
	for _, validatorName := range fireList {
		if err := vm.RemoveValidatorFromAllLists(ctx, validatorName); err != nil {
			return err
		}
	}

	return nil
}

// try to join the validator list, the action will success if either
// 1. the validator list is not full
// or 2. someone in the validator list has a lower power than current validator
// return a boolean to indicate if the user has became an oncall validator
func (vm ValidatorManager) TryJoinValidatorList(ctx sdk.Context, username acc.AccountKey, addToPool bool) sdk.Error {
	curValidator, getErr := vm.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}
	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}
	// add to validator pool if needed
	if addToPool {
		lst.AllValidators = append(lst.AllValidators, username)
	}

	// add to list directly if validator list is not full
	if len(lst.OncallValidators) < types.ValidatorListSize {
		if len(lst.OncallValidators) == 0 || curValidator.ABCIValidator.Power < lst.LowestPower.AmountOf(types.Denom) {
			lst.LowestPower = sdk.Coins{sdk.Coin{Denom: types.Denom, Amount: curValidator.ABCIValidator.Power}}
			lst.LowestValidator = curValidator.Username
		}
		lst.OncallValidators = append(lst.OncallValidators, curValidator.Username)
	} else if curValidator.ABCIValidator.Power > lst.LowestPower.AmountOf(types.Denom) {
		// replace the validator with lowest power
		// 1. iterate through validator list to replace the lowest validator
		for idx, validatorKey := range lst.OncallValidators {
			validator, getErr := vm.GetValidator(ctx, validatorKey)
			if getErr != nil {
				return getErr
			}
			if validator.Username == lst.LowestValidator {
				lst.OncallValidators[idx] = curValidator.Username
			}
		}

		// 2. iterate through validator list to update lowest power&validator
		//updateErr := nil
		vm.updateLowestValidator(ctx, lst)
	}
	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	return nil
}

// remove the user from both oncall and allValidators lists
func (vm ValidatorManager) RemoveValidatorFromAllLists(ctx sdk.Context, username acc.AccountKey) sdk.Error {
	lst, getListErr := vm.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}

	lst.AllValidators = remove(username, lst.AllValidators)
	lst.OncallValidators = remove(username, lst.OncallValidators)

	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
	}

	vm.updateLowestValidator(ctx, lst)

	// find the person has the biggest power among people in the allValidators lists
	// but not in the oncall validator list
	bestCandidate := acc.AccountKey("")
	bestCandidatePower := int64(0)

	for i, validatorName := range lst.AllValidators {
		validator, getErr := vm.GetValidator(ctx, lst.AllValidators[i])
		if getErr != nil {
			return getErr
		}

		// not in the oncall list and has a larger power
		if findAccountInList(validatorName, lst.OncallValidators) == -1 &&
			validator.ABCIValidator.Power > bestCandidatePower {
			bestCandidate = validator.Username
			bestCandidatePower = validator.ABCIValidator.Power
		}
	}

	if bestCandidate != acc.AccountKey("") {
		if joinErr := vm.TryJoinValidatorList(ctx, bestCandidate, false); joinErr != nil {
			return joinErr
		}
	}
	if err := vm.SetValidatorList(ctx, lst); err != nil {
		return err
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

func (vm ValidatorManager) updateLowestValidator(ctx sdk.Context, lst *ValidatorList) {
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
	lst.LowestPower = sdk.Coins{sdk.Coin{Denom: types.Denom, Amount: newLowestPower}}
	lst.LowestValidator = newLowestValidator
}

func ValidatorKey(accKey acc.AccountKey) []byte {
	return append(ValidatorAccountPrefix, accKey...)
}

func ValidatorListKey() []byte {
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
