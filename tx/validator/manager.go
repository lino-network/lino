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
func (vm ValidatorManager) GetValidatorAccount(ctx sdk.Context, accKey acc.AccountKey) (*ValidatorAccount, sdk.Error) {
	store := ctx.KVStore(vm.key)
	accountByte := store.Get(validatorKey(accKey))
	if accountByte == nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get account failed: account doesn't exist")
	}
	acc := new(ValidatorAccount)
	if err := vm.cdc.UnmarshalBinary(accountByte, acc); err != nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get account failed")
	}
	return acc, nil
}

// Implements ValidatorManager
func (vm ValidatorManager) SetValidatorAccount(ctx sdk.Context, accKey acc.AccountKey, account *ValidatorAccount) sdk.Error {
	store := ctx.KVStore(vm.key)
	accountByte, err := vm.cdc.MarshalBinary(*account)
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
	if err := vm.cdc.UnmarshalBinary(listByte, lst); err != nil {
		return nil, ErrValidatorManagerFail("ValidatorManager get account list failed")
	}
	return lst, nil
}

// Implements ValidatorManager
func (vm ValidatorManager) SetValidatorList(ctx sdk.Context, accKey acc.AccountKey, lst *ValidatorList) sdk.Error {
	store := ctx.KVStore(vm.key)
	listByte, err := vm.cdc.MarshalBinary(*lst)
	if err != nil {
		return ErrValidatorManagerFail("ValidatorManager set account list failed")
	}
	store.Set(validatorListKey(ValidatorListKey), listByte)
	return nil
}

// try to join the validator list.
// the action will success if either
// 1. the validator list is not full
// or 2. someone in the validator list has a lower power than current validator
func (vm ValidatorManager) TryJoinValidatorList(ctx sdk.Context, validatorName acc.AccountKey, addToPool bool) bool {
	curValidator, _ := vm.GetValidatorAccount(ctx, validatorName)
	lst, _ := vm.GetValidatorList(ctx, ValidatorListKey)
	defer vm.SetValidatorList(ctx, ValidatorListKey, lst)
	// add to validator pool if needed
	if addToPool {
		lst.validatorPool = append(lst.validatorPool, validatorName)
	}

	// add to list directly if validator list is not full
	if len(lst.validators) < types.ValidatorListSize {
		if curValidator.Power < lst.lowestPower.AmountOf("lino") {
			lst.lowestPower = sdk.Coins{sdk.Coin{Denom: "lino", Amount: curValidator.Power}}
			lst.lowestValidator = curValidator.validatorName
		}
		lst.validators = append(lst.validators, curValidator.validatorName)
		return true
	}

	// replace the validator with lowest power
	if curValidator.Power > lst.lowestPower.AmountOf("lino") {
		newLowestPower := curValidator.Power
		newLowestValidator := curValidator.validatorName

		// iterate through validator list to
		// 1. replace the lowest validator, 2.find new lowest validator and power
		for idx, accKey := range lst.validators {
			acc, _ := vm.GetValidatorAccount(ctx, accKey)
			//replacement
			if acc.validatorName == lst.lowestValidator {
				lst.validators[idx] = curValidator.validatorName
			}
			// update lowest power and validator
			if acc.Power < newLowestPower {
				newLowestPower = acc.Power
				newLowestValidator = acc.validatorName
			}
		}

		// set the new lowest power
		lst.lowestPower = sdk.Coins{sdk.Coin{Denom: "lino", Amount: newLowestPower}}
		lst.lowestValidator = newLowestValidator

		return true
	}
	return false
}

func validatorKey(accKey acc.AccountKey) []byte {
	return append(ValidatorAccountPrefix, accKey...)
}

func validatorListKey(accKey acc.AccountKey) []byte {
	return append(ValidatorListPrefix, accKey...)
}
