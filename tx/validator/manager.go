package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

var ValidatorAccountPrefix = []byte("ValidatorAccountInfo/")
var ValidatorListPrefix = []byte("ValidatorList/")

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
	listByte := store.Get(validatorListKey(accKey))
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
	store.Set(validatorListKey(accKey), listByte)
	return nil
}

// try to join the validator list.
// the action will success if either
// 1. the validator list is not full
// or 2. someone in the validator list has a lower weight than current validator
func (vm ValidatorManager) TryJoinValidatorList(ctx sdk.Context, accKey acc.AccountKey) bool {
	validator, _ := vm.GetValidatorAccount(ctx, accKey)
	lst, _ := vm.GetValidatorList(ctx, "validatoryKey")
	validatorList := lst.validators
	// add to list directly
	if len(validatorList) < types.ValidatorListSize {
		if validator.totalWeight < lst.minWeight {
			lst.minWeight = validator.totalWeight
		}
		validatorList = append(validatorList, validator.validatorName)
		vm.SetValidatorList(ctx, "validatoryKey", lst)
		return true
	}

	// replace the validator with lowest weight
	if validator.totalWeight > lst.minWeight {

		newMinWeight := validator.totalWeight
		for idx, accKey := range validatorList {
			acc, _ := vm.GetValidatorAccount(ctx, accKey)
			//delete the validator has the lowest weight, add new validator
			if acc.totalWeight == lst.minWeight {
				validatorList = append(validatorList[:idx], validatorList[idx+1:]...)
				validatorList = append(validatorList, validator.validatorName)
			}

			if acc.totalWeight < newMinWeight {
				newMinWeight = acc.totalWeight
			}
		}

		// update the lowest weight in the validator list
		lst.minWeight = newMinWeight
		vm.SetValidatorList(ctx, "validatoryKey", lst)
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
