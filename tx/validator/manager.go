package validator

import (
	"math"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
)

// validator manager is the proxy for all storage structs defined above
type ValidatorManager struct {
	storage            *model.ValidatorStorage `json:"validator_storage"`
	preRoundValidators []model.Validator       `json:"pre_round_validators"`
}

// create NewValidatorManager
func NewValidatorManager(key sdk.StoreKey) *ValidatorManager {
	return &ValidatorManager{
		storage: model.NewValidatorStorage(key),
	}
}

func (vm *ValidatorManager) SetPreRoundValidators(ctx sdk.Context) sdk.Error {
	var err sdk.Error
	vm.preRoundValidators, err = vm.GetOncallValList(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (vm *ValidatorManager) GetUpdateValidatorList(ctx sdk.Context) ([]abci.Validator, sdk.Error) {
	curOncallList, err := vm.GetOncallValList(ctx)
	if err != nil {
		return nil, err
	}
	ABCIValList := []abci.Validator{}
	for _, preValidator := range vm.preRoundValidators {
		if findValidatorInList(preValidator, curOncallList) == -1 {
			preValidator.ABCIValidator.Power = 0
			ABCIValList = append(ABCIValList, preValidator.ABCIValidator)
		}
	}
	for _, validator := range curOncallList {
		ABCIValList = append(ABCIValList, validator.ABCIValidator)
	}
	return ABCIValList, nil
}

func (vm ValidatorManager) InitGenesis(ctx sdk.Context) error {
	lst := &model.ValidatorList{
		LowestPower: types.Coin{0},
	}

	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
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

	// reject if this is an oncall validator
	lst, getErr := vm.storage.GetValidatorList(ctx)
	if getErr != nil {
		return false
	}

	if FindAccountInList(username, lst.OncallValidators) != -1 {
		return false
	}
	//reject if the remaining coins are less than register fee
	res := validator.Deposit.Minus(coin)
	return res.IsGTE(types.ValidatorRegisterFee)
}

func (vm ValidatorManager) GetOncallValList(ctx sdk.Context) ([]model.Validator, sdk.Error) {
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return nil, getListErr
	}

	oncallList := make([]model.Validator, len(lst.OncallValidators))
	for i, validatorName := range lst.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
		if err != nil {
			return nil, err
		}
		oncallList[i] = *validator
	}
	return oncallList, nil
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
		validator.AbsentVote += 1

		if err := vm.storage.SetValidator(ctx, lst.OncallValidators[idx], validator); err != nil {
			return err
		}
	}

	return nil
}

func (vm ValidatorManager) MarkByzantine(ctx sdk.Context, username types.AccountKey) sdk.Error {
	byzantine, err := vm.storage.GetValidator(ctx, username)
	if err != nil {
		return err
	}

	byzantine.IsByzantine = true
	if err := vm.storage.SetValidator(ctx, username, byzantine); err != nil {
		return err
	}
	return nil
}

func (vm ValidatorManager) FireIncompetentValidator(ctx sdk.Context, ByzantineValidators []abci.Evidence) sdk.Error {
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}
	fireList := []types.AccountKey{}

	for _, validatorName := range lst.OncallValidators {
		validator, err := vm.storage.GetValidator(ctx, validatorName)
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

func (vm ValidatorManager) RegisterValidator(ctx sdk.Context, username types.AccountKey, pubKey []byte, coin types.Coin) sdk.Error {
	curValidator := &model.Validator{
		ABCIValidator: abci.Validator{PubKey: pubKey, Power: coin.Amount},
		Username:      username,
		Deposit:       coin,
		IsByzantine:   false,
	}
	// check minimum requirements
	if !coin.IsGTE(types.ValidatorRegisterFee) {
		return ErrRegisterFeeNotEnough()
	}

	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}

	// has alreay in the validator list
	if FindAccountInList(username, lst.AllValidators) != -1 {
		return nil
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
	validator.ABCIValidator.Power = validator.Deposit.Amount
	if setErr := vm.storage.SetValidator(ctx, username, validator); setErr != nil {
		return setErr
	}
	return nil
}

// this method won't check if it is a legal withdraw, caller should check by itself
func (vm ValidatorManager) Withdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin, gm global.GlobalManager) sdk.Error {
	validator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}
	validator.Deposit = validator.Deposit.Minus(coin)

	if err := vm.storage.SetValidator(ctx, username, validator); err != nil {
		return err
	}
	if err := vm.CreateReturnCoinEvent(ctx, username, coin, gm); err != nil {
		return nil
	}
	return nil
}

func (vm ValidatorManager) WithdrawAll(ctx sdk.Context, username types.AccountKey, gm global.GlobalManager) sdk.Error {
	validator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}
	if err := vm.Withdraw(ctx, username, validator.Deposit, gm); err != nil {
		return err
	}
	return nil
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

	valRegisterFee, err := types.LinoToCoin(types.LNO(sdk.NewRat(1000)))
	if err != nil {
		return sdk.ErrInvalidCoins("invalid register fee")
	}

	// check minimum requirements
	if !curValidator.Deposit.IsGTE(valRegisterFee) {
		return ErrRegisterFeeNotEnough()
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
	} else if curValidator.ABCIValidator.Power > lst.LowestPower.Amount {
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
	curValidator, getErr := vm.storage.GetValidator(ctx, username)
	if getErr != nil {
		return getErr
	}

	curValidator.ABCIValidator.Power = 0
	lst, getListErr := vm.storage.GetValidatorList(ctx)
	if getListErr != nil {
		return getListErr
	}

	if FindAccountInList(username, lst.AllValidators) == -1 {
		return ErrNotInTheList()
	}

	lst.AllValidators = remove(username, lst.AllValidators)
	lst.OncallValidators = remove(username, lst.OncallValidators)

	//TODO return deposit to pool?
	// if curValidator.IsByzantine {
	// }
	if err := vm.storage.SetValidatorList(ctx, lst); err != nil {
		return err
	}
	if err := vm.storage.SetValidator(ctx, curValidator.Username, curValidator); err != nil {
		return err
	}

	vm.updateLowestValidator(ctx)
	bestCandidate, findErr := vm.getBestCandidate(ctx, lst)
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
	newLowestPower := int64(math.MaxInt64)
	newLowestValidator := types.AccountKey("")

	for _, validatorKey := range lst.OncallValidators {
		validator, _ := vm.storage.GetValidator(ctx, validatorKey)
		if validator.ABCIValidator.Power < newLowestPower {
			newLowestPower = validator.ABCIValidator.Power
			newLowestValidator = validator.Username
		}
	}
	// set the new lowest power
	lst.LowestPower = types.NewCoin(newLowestPower)
	lst.LowestValidator = newLowestValidator

	vm.storage.SetValidatorList(ctx, lst)
}

// find the person has the biggest power among people in the allValidators lists
// but not in the oncall validator list
func (vm ValidatorManager) getBestCandidate(ctx sdk.Context, lst *model.ValidatorList) (types.AccountKey, sdk.Error) {
	bestCandidate := types.AccountKey("")
	bestCandidatePower := int64(0)

	for i, validatorName := range lst.AllValidators {
		validator, getErr := vm.storage.GetValidator(ctx, lst.AllValidators[i])
		if getErr != nil {
			return bestCandidate, getErr
		}
		// not in the oncall list and has a larger power
		if FindAccountInList(validatorName, lst.OncallValidators) == -1 &&
			validator.ABCIValidator.Power > bestCandidatePower {
			bestCandidate = validator.Username
			bestCandidatePower = validator.ABCIValidator.Power
		}
	}
	return bestCandidate, nil

}

// return coin to an user periodically
func (vm ValidatorManager) CreateReturnCoinEvent(ctx sdk.Context, username types.AccountKey, amount types.Coin, gm global.GlobalManager) sdk.Error {
	pieceRat := amount.ToRat().Quo(sdk.NewRat(types.CoinReturnTimes))
	piece := types.RatToCoin(pieceRat)
	event := ReturnCoinEvent{
		Username: username,
		Amount:   piece,
	}

	// return coin with interval
	for i := int64(1); i <= types.CoinReturnTimes; i++ {
		if err := gm.RegisterEventAtTime(ctx, ctx.BlockHeader().Time+(types.CoinReturnIntervalHr*3600*i), event); err != nil {
			return err
		}
	}
	return nil
}

func FindAccountInList(me types.AccountKey, lst []types.AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}

func findValidatorInList(validator model.Validator, validatorList []model.Validator) int {
	for i, curValidator := range validatorList {
		if validator.Username == curValidator.Username {
			return i
		}
	}
	return -1
}
