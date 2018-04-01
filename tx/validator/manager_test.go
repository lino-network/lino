package validator

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
)

func TestAbsentValidator(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newValidatorManager()
	ctx := getContext()
	handler := NewHandler(vm, lam)
	vm.InitGenesis(ctx)

	// create 21 test users
	users := make([]*acc.Account, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, lam, "user"+strconv.Itoa(i))
		users[i].AddCoin(ctx, c2000)
		users[i].Apply(ctx)
		// they will deposit 10,20,30...200, 210
		deposit := types.LNO(sdk.NewRat(int64((i+1)*10) + int64(1001)))
		ownerKey, _ := users[i].GetOwnerKey(ctx)
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i), deposit, *ownerKey)
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}
	absentList := []int32{0, 1, 10, 20}
	err := vm.UpdateAbsentValidator(ctx, absentList)
	assert.Nil(t, err)

	validatorList, _ := vm.GetValidatorList(ctx)
	for _, idx := range absentList {
		validator, _ := vm.GetValidator(ctx, validatorList.OncallValidators[idx])
		assert.Equal(t, validator.AbsentVote, 1)
	}

	// absent exceeds limitation
	for i := 0; i < types.AbsentLimitation; i++ {
		err := vm.UpdateAbsentValidator(ctx, absentList)
		assert.Nil(t, err)
	}

	for _, idx := range absentList {
		validator, _ := vm.GetValidator(ctx, validatorList.OncallValidators[idx])
		assert.Equal(t, validator.AbsentVote, 101)
	}
	err = vm.FireIncompetentValidator(ctx, []abci.Evidence{})
	assert.Nil(t, err)
	validatorList2, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, 17, len(validatorList2.OncallValidators))
	assert.Equal(t, 17, len(validatorList2.AllValidators))

	for _, idx := range absentList {
		assert.Equal(t, -1, FindAccountInList(users[idx].GetUsername(ctx), validatorList2.OncallValidators))
		assert.Equal(t, -1, FindAccountInList(users[idx].GetUsername(ctx), validatorList2.AllValidators))
	}

	// byzantine
	byzantineList := []int32{3, 8, 14}
	byzantines := []abci.Evidence{}
	for _, idx := range byzantineList {
		ownerKey, _ := users[idx].GetOwnerKey(ctx)
		byzantines = append(byzantines, abci.Evidence{PubKey: ownerKey.Bytes()})
	}
	err = vm.FireIncompetentValidator(ctx, byzantines)
	assert.Nil(t, err)

	validatorList3, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, 14, len(validatorList3.OncallValidators))
	assert.Equal(t, 14, len(validatorList3.AllValidators))

	for _, idx := range byzantineList {
		assert.Equal(t, -1, FindAccountInList(users[idx].GetUsername(ctx), validatorList3.OncallValidators))
		assert.Equal(t, -1, FindAccountInList(users[idx].GetUsername(ctx), validatorList3.AllValidators))
	}
}

func TestGetOncallList(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newValidatorManager()
	ctx := getContext()
	handler := NewHandler(vm, lam)
	vm.InitGenesis(ctx)

	// create 21 test users
	users := make([]*acc.Account, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, lam, "user"+strconv.Itoa(i))
		users[i].AddCoin(ctx, c2000)
		users[i].Apply(ctx)
		// they will deposit 10,20,30...200, 210
		deposit := types.LNO(sdk.NewRat(int64((i+1)*10) + int64(1001)))
		ownerKey, _ := users[i].GetOwnerKey(ctx)
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i), deposit, *ownerKey)
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	lst, _ := vm.GetOncallValList(ctx)
	for idx, validator := range lst {
		assert.Equal(t, acc.AccountKey("user"+strconv.Itoa(idx)), validator.Username)
	}

}
