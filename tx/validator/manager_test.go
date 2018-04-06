package validator

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
)

func TestAbsentValidator(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)
	vm.InitGenesis(ctx)

	// create 21 test users
	users := make([]types.AccountKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i))
		am.AddCoin(ctx, users[i], c2000)

		// they will deposit 10,20,30...200, 210
		deposit := types.LNO(sdk.NewRat(int64((i+1)*10) + int64(1001)))
		ownerKey, _ := am.GetOwnerKey(ctx, users[i])
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i), deposit, *ownerKey)
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}
	absentList := []int32{0, 1, 10, 20}
	err := vm.UpdateAbsentValidator(ctx, absentList)
	assert.Nil(t, err)

	validatorList, _ := vm.storage.GetValidatorList(ctx)
	for _, idx := range absentList {
		validator, _ := vm.storage.GetValidator(ctx, validatorList.OncallValidators[idx])
		assert.Equal(t, validator.AbsentVote, 1)
	}

	// absent exceeds limitation
	for i := 0; i < types.AbsentLimitation; i++ {
		err := vm.UpdateAbsentValidator(ctx, absentList)
		assert.Nil(t, err)
	}

	for _, idx := range absentList {
		validator, _ := vm.storage.GetValidator(ctx, validatorList.OncallValidators[idx])
		assert.Equal(t, validator.AbsentVote, 101)
	}
	err = vm.FireIncompetentValidator(ctx, []abci.Evidence{})
	assert.Nil(t, err)
	validatorList2, _ := vm.storage.GetValidatorList(ctx)
	assert.Equal(t, 17, len(validatorList2.OncallValidators))
	assert.Equal(t, 17, len(validatorList2.AllValidators))

	for _, idx := range absentList {
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList2.OncallValidators))
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList2.AllValidators))
	}

	// byzantine
	byzantineList := []int32{3, 8, 14}
	byzantines := []abci.Evidence{}
	for _, idx := range byzantineList {
		ownerKey, _ := am.GetOwnerKey(ctx, users[idx])
		byzantines = append(byzantines, abci.Evidence{PubKey: ownerKey.Bytes()})
	}
	err = vm.FireIncompetentValidator(ctx, byzantines)
	assert.Nil(t, err)

	validatorList3, _ := vm.storage.GetValidatorList(ctx)
	assert.Equal(t, 14, len(validatorList3.OncallValidators))
	assert.Equal(t, 14, len(validatorList3.AllValidators))

	for _, idx := range byzantineList {
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList3.OncallValidators))
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList3.AllValidators))
	}
}

func TestGetOncallList(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)
	vm.InitGenesis(ctx)

	// create 21 test users
	users := make([]types.AccountKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i))
		am.AddCoin(ctx, users[i], c2000)

		// they will deposit 10,20,30...200, 210
		deposit := types.LNO(sdk.NewRat(int64((i+1)*10) + int64(1001)))
		ownerKey, _ := am.GetOwnerKey(ctx, users[i])
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i), deposit, *ownerKey)
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	lst, _ := vm.GetOncallValList(ctx)
	for idx, validator := range lst {
		assert.Equal(t, users[idx], validator.Username)
	}

}
