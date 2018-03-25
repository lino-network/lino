package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

var ValidatorAccountPrefix = []byte("ValidatorAccountInfo/")
var ValidatorListPrefix = []byte("ValidatorList/")
var ValidatorListKey = acc.AccountKey("ValidatorListKey")

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

	return ValidatorManager{
		key: key,
		cdc: cdc,
	}
}

// Implements ValidatorManager
func (vm ValidatorManager) IsValidatorExist(ctx sdk.Context, accKey acc.AccountKey) bool {
	store := ctx.KVStore(vm.key)
	if infoByte := store.Get(validatorKey(accKey)); infoByte == nil {
		return false
	}
	return true
}

// Implements ValidatorManager
func (vm ValidatorManager) GetValidator(ctx sdk.Context, accKey acc.AccountKey) (*Validator, sdk.Error) {
	store := ctx.KVStore(vm.key)
	accountByte := store.Get(validatorKey(accKey))
	if accountByte == nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get account failed: account doesn't exist")
	}
	acc := new(Validator)
	if err := vm.cdc.UnmarshalJSON(accountByte, acc); err != nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get account failed")
	}
	return acc, nil
}

// Implements ValidatorManager
func (vm ValidatorManager) SetValidator(ctx sdk.Context, accKey acc.AccountKey, account *Validator) sdk.Error {
	store := ctx.KVStore(vm.key)
	accountByte, err := vm.cdc.MarshalJSON(*account)
	if err != nil {
		return ErrValidatorManagerFail("ValidatorManager set account failed")
	}
	store.Set(validatorKey(accKey), accountByte)
	return nil
}

// Implements ValidatorManager
func (vm ValidatorManager) GetValidatorList(ctx sdk.Context, accKey acc.AccountKey) (*ValidatorList, sdk.Error) {
	store := ctx.KVStore(vm.key)
	listByte := store.Get(validatorListKey(ValidatorListKey))
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
func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, accKey acc.AccountKey, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vm.key)
	listByte, err := vm.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrValidatorManagerFail("ValidatorManager set account list failed")
	}
	store.Set(validatorListKey(ValidatorListKey), listByte)
	return nil
}

// try to join the validator list, the action will success if either
// 1. the validator list is not full
// or 2. someone in the validator list has a lower power than current validator
// return a boolean to indicate if the user has became an oncall validator
func (vm ValidatorManager) TryJoinValidatorList(ctx sdk.Context, username acc.AccountKey, addToPool bool) (bool, sdk.Error) {
	curValidator, getErr := vm.GetValidator(ctx, username)
	if getErr != nil {
		return false, ErrValidatorManagerFail("Get validator failed")
	}
	lst, setErr := vm.GetValidatorList(ctx, ValidatorListKey)
	if setErr != nil {
		return false, ErrValidatorManagerFail("Set validator failed")
	}
	defer vm.SetValidatorList(ctx, ValidatorListKey, lst)
	// add to validator pool if needed
	if addToPool {
		lst.AllValidators = append(lst.AllValidators, username)
	}

	// add to list directly if validator list is not full
	if len(lst.OncallValidators) < types.ValidatorListSize {
		if curValidator.ABCIValidator.Power < lst.LowestPower.AmountOf("lino") || len(lst.OncallValidators) == 0 {
			lst.LowestPower = sdk.Coins{sdk.Coin{Denom: "lino", Amount: curValidator.ABCIValidator.Power}}
			lst.LowestValidator = curValidator.Username
		}
		lst.OncallValidators = append(lst.OncallValidators, curValidator.Username)
		return true, nil
	}

	// replace the validator with lowest power
	if curValidator.ABCIValidator.Power > lst.LowestPower.AmountOf("lino") {
		// 1. iterate through validator list to replace the lowest validator
		for idx, validatorKey := range lst.OncallValidators {
			validator, getErr := vm.GetValidator(ctx, validatorKey)
			if getErr != nil {
				return false, ErrValidatorManagerFail("Get validator failed")
			}
			if validator.Username == lst.LowestValidator {
				lst.OncallValidators[idx] = curValidator.Username
			}
		}

		// 2. iterate through validator list to update lowest power&validator
		newLowestPower := curValidator.ABCIValidator.Power
		newLowestValidator := curValidator.Username

		for _, validatorKey := range lst.OncallValidators {
			validator, getErr := vm.GetValidator(ctx, validatorKey)
			if getErr != nil {
				return false, ErrValidatorManagerFail("Get validator failed")
			}
			if validator.ABCIValidator.Power < newLowestPower {
				newLowestPower = validator.ABCIValidator.Power
				newLowestValidator = validator.Username
			}
		}
		// set the new lowest power
		lst.LowestPower = sdk.Coins{sdk.Coin{Denom: "lino", Amount: newLowestPower}}
		lst.LowestValidator = newLowestValidator
		return true, nil
	}
	return false, nil
}

func validatorKey(accKey acc.AccountKey) []byte {
	return append(ValidatorAccountPrefix, accKey...)
}

func validatorListKey(accKey acc.AccountKey) []byte {
	return append(ValidatorListPrefix, accKey...)
}
