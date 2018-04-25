package validator

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	l10   = types.LNO("10")
	l15   = types.LNO("15")
	l100  = types.LNO("100")
	l200  = types.LNO("200")
	l400  = types.LNO("400")
	l1000 = types.LNO("1000")
	l1100 = types.LNO("1100")
	l1600 = types.LNO("1600")

	c0    = types.Coin{0 * types.Decimals}
	c200  = types.Coin{200 * types.Decimals}
	c400  = types.Coin{400 * types.Decimals}
	c1011 = types.Coin{1011 * types.Decimals}
	c1021 = types.Coin{1021 * types.Decimals}
	c1600 = types.Coin{1600 * types.Decimals}
	c1800 = types.Coin{1800 * types.Decimals}
	c2000 = types.Coin{2000 * types.Decimals}
	c8000 = types.Coin{8000 * types.Decimals}
)

func TestRegisterBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create two test users
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", c8000)

	// let user1 register as validator
	ownerKey, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	assert.Equal(t, acc1Balance, c400.Plus(initCoin))
	assert.Equal(t, true, valManager.IsValidatorExist(ctx, user1))

	// now user1 should be the only validator (WOW, dictator!)
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, verifyList.LowestPower, c1600)
	assert.Equal(t, 1, len(verifyList.OncallValidators))
	assert.Equal(t, 1, len(verifyList.AllValidators))
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// make sure the validator's account info (power&pubKey) is correct
	verifyAccount, _ := valManager.storage.GetValidator(ctx, user1)
	assert.Equal(t, c1600, verifyAccount.Deposit)
	assert.Equal(t, ownerKey.Bytes(), verifyAccount.ABCIValidator.GetPubKey())
}

func TestRegisterFeeNotEnough(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create test user
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	// let user1 register as validator
	ownerKey, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("user1", l400, *ownerKey)
	result := handler(ctx, msg)
	assert.Equal(t, ErrVotingDepositNotEnough().Result(), result)

	// let user register as voter
	voteManager.AddVoter(ctx, "user1", c8000)

	result2 := handler(ctx, msg)
	assert.Equal(t, ErrCommitingDepositNotEnough().Result(), result2)

	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 0, len(verifyList.OncallValidators))
	assert.Equal(t, 0, len(verifyList.AllValidators))
}

func TestRevokeBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create two test users
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", c8000)

	// let user1 register as validator
	ownerKey, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// now user1 should be the only validator
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// let user1 revoke candidancy
	msg2 := NewValidatorRevokeMsg("user1")
	result2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, result2)

	verifyList2, _ := valManager.storage.GetValidatorList(ctx)
	validator, _ := valManager.storage.GetValidator(ctx, "user1")
	assert.Equal(t, 0, len(verifyList2.OncallValidators))
	assert.Equal(t, 0, len(verifyList2.AllValidators))
	assert.Equal(t, c0, validator.Deposit)

}

func TestRevokeNonExistUser(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// let user1(not exists) revoke candidancy
	msg2 := NewValidatorRevokeMsg("user1")
	result2 := handler(ctx, msg2)
	assert.Equal(t, model.ErrGetValidator().Result(), result2)
}

// this is the same situation as we find Byzantine and replace the Byzantine
func TestRevokeOncallValidatorAndSubstitutionExists(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create 21 test users
	users := make([]types.AccountKey, 24)
	for i := 0; i < 24; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1))
		am.AddCoin(ctx, users[i], c2000)

		// let user register as voter first
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), c8000)

		// they will deposit 1000 + 10,20,30...200, 210, 220, 230, 240
		num := (i+1)*10 + 1000
		deposit := types.LNO(strconv.Itoa(num))
		ownerKey, _ := am.GetOwnerKey(ctx, users[i])
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, *ownerKey)
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 21, len(lst.OncallValidators))
	assert.Equal(t, 24, len(lst.AllValidators))
	assert.Equal(t, types.Coin{1040 * types.Decimals}, lst.LowestPower)
	assert.Equal(t, users[3], lst.LowestValidator)

	// lowest validator depoist coins will change the ranks
	ownerKey, _ := am.GetOwnerKey(ctx, users[3])
	deposit := types.LNO(l15)
	msg := NewValidatorDepositMsg("user4", deposit, *ownerKey)
	result := handler(ctx, msg)

	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result)
	assert.Equal(t, types.Coin{1050 * types.Decimals}, lst2.LowestPower)
	assert.Equal(t, users[4], lst2.LowestValidator)

	// now user1, 2, 3 are substitutions
	// user2 can only withdraw 1-20 coins
	withdrawMsg := NewValidatorWithdrawMsg("user2", l100)
	resultWithdraw := handler(ctx, withdrawMsg)
	assert.Equal(t, ErrIllegalWithdraw().Result(), resultWithdraw)

	withdrawMsg2 := NewValidatorWithdrawMsg("user2", l10)
	resultWithdraw2 := handler(ctx, withdrawMsg2)
	assert.Equal(t, sdk.Result{}, resultWithdraw2)
	//revoke a non oncall valodator wont change anything related to oncall list
	revokeMsg := NewValidatorRevokeMsg("user2")
	result2 := handler(ctx, revokeMsg)
	assert.Equal(t, sdk.Result{}, result2)

	lst3, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, types.Coin{1050 * types.Decimals}, lst3.LowestPower)
	assert.Equal(t, users[4], lst3.LowestValidator)
	assert.Equal(t, 23, len(lst3.AllValidators))

	// now only user1(power1010) and user3(power1030) are substitutions
	// the lowest oncall user is user5 with 1050 power
	// revoke user6 (could be byzantine) will make user3 (power 1030) join oncall
	// list become the lowest validator
	revokeMsg2 := NewValidatorRevokeMsg("user6")
	result3 := handler(ctx, revokeMsg2)
	assert.Equal(t, sdk.Result{}, result3)

	lst4, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, types.Coin{1030 * types.Decimals}, lst4.LowestPower)
	assert.Equal(t, users[2], lst4.LowestValidator)
	assert.Equal(t, 22, len(lst4.AllValidators))
}

func TestRevokeAndDepositAgain(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create user
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	// let user register as voter first
	voteManager.AddVoter(ctx, "user1", c8000)

	// let user1 register as validator
	ownerKey, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("user1", l1000, *ownerKey)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 1, len(lst.AllValidators))
	assert.Equal(t, 1, len(lst.OncallValidators))

	// let user1 revoke candidancy
	msg2 := NewValidatorRevokeMsg("user1")
	result2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, result2)

	lstEmpty, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 0, len(lstEmpty.AllValidators))
	assert.Equal(t, 0, len(lstEmpty.OncallValidators))

	// deposit again
	msg3 := NewValidatorDepositMsg("user1", l1000, *ownerKey)
	result3 := handler(ctx, msg3)

	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result3)
	assert.Equal(t, 1, len(lst2.AllValidators))
	assert.Equal(t, 1, len(lst2.OncallValidators))
}

func TestWithdrawBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create test user
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", c8000)

	// let user1 register as validator
	ownerKey, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// now user1 should be the only validator
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// user1 cannot withdraw if is oncall validator
	withdrawMsg := NewValidatorWithdrawMsg("user1", l1600)
	result2 := handler(ctx, withdrawMsg)
	assert.Equal(t, ErrIllegalWithdraw().Result(), result2)

}

func TestDepositBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create test user
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	// let user register as voter first
	voteManager.AddVoter(ctx, "user1", c8000)

	// let user1 register as validator
	ownerKey, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	depositMsg := NewValidatorDepositMsg("user1", l200, *ownerKey)
	result2 := handler(ctx, depositMsg)
	assert.Equal(t, sdk.Result{}, result2)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	assert.Equal(t, acc1Balance, c200.Plus(initCoin))
	assert.Equal(t, true, valManager.IsValidatorExist(ctx, user1))

	// check the lowest power is 1800 now (1600 + 200)
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, c1800, verifyList.LowestPower)
	assert.Equal(t, 1, len(verifyList.OncallValidators))
	assert.Equal(t, 1, len(verifyList.AllValidators))
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// check deposit and power is correct
	validator, _ := valManager.storage.GetValidator(ctx, user1)
	assert.Equal(t, true, validator.Deposit.IsEqual(c1800))
}

func TestDepositWithoutLinoAccount(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// let user1 register as validator
	user1 := createTestAccount(ctx, am, "user1")
	ownerKey, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("qwqwndqwnd", l1600, *ownerKey)
	result := handler(ctx, msg)
	assert.Equal(t, ErrUsernameNotFound().Result(), result)
}

func TestValidatorReplacement(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create 21 test users
	users := make([]types.AccountKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1))
		am.AddCoin(ctx, users[i], c2000)
		// let user register as voter first
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), c8000)
		// they will deposit 10,20,30...200, 210
		num := (i+1)*10 + 1001
		deposit := types.LNO(strconv.Itoa(num))
		ownerKey, _ := am.GetOwnerKey(ctx, users[i])
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, *ownerKey)
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// check validator list, the lowest power is 10
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, true, verifyList.LowestPower.IsEqual(c1011))
	assert.Equal(t, users[0], verifyList.LowestValidator)
	assert.Equal(t, 21, len(verifyList.OncallValidators))
	assert.Equal(t, 21, len(verifyList.AllValidators))

	// create a user failed to join oncall validator list (not enough power)
	user1 := createTestAccount(ctx, am, "noPowerUser")
	am.AddCoin(ctx, user1, c2000)
	// let user register as voter first
	voteManager.AddVoter(ctx, "noPowerUser", c8000)

	//check the user hasn't been added to oncall validators but in the pool
	deposit := types.LNO("1005")
	ownerKey1, _ := am.GetOwnerKey(ctx, user1)
	msg := NewValidatorDepositMsg("noPowerUser", deposit, *ownerKey1)
	result := handler(ctx, msg)

	verifyList2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result)
	assert.Equal(t, true, verifyList.LowestPower.IsEqual(c1011))
	assert.Equal(t, users[0], verifyList.LowestValidator)
	assert.Equal(t, 21, len(verifyList2.OncallValidators))
	assert.Equal(t, 22, len(verifyList2.AllValidators))

	// create a user success to join oncall validator list
	powerfulUser := createTestAccount(ctx, am, "powerfulUser")
	am.AddCoin(ctx, powerfulUser, c2000)
	// let user register as voter first
	voteManager.AddVoter(ctx, "powerfulUser", c8000)

	//check the user has been added to oncall validators and in the pool
	deposit2 := types.LNO("1088")
	ownerKey2, _ := am.GetOwnerKey(ctx, powerfulUser)
	msg2 := NewValidatorDepositMsg("powerfulUser", deposit2, *ownerKey2)
	result2 := handler(ctx, msg2)

	verifyList3, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result2)
	assert.Equal(t, true, verifyList3.LowestPower.IsEqual(c1021))
	assert.Equal(t, users[1], verifyList3.LowestValidator)
	assert.Equal(t, 21, len(verifyList3.OncallValidators))
	assert.Equal(t, 23, len(verifyList3.AllValidators))

	// check user0 has been replaced, and powerful user has been added
	flag := false
	for _, username := range verifyList3.OncallValidators {
		if username == "powerfulUser" {
			flag = true
		}
		if username == "user0" {
			assert.Fail(t, "User0 should have been replaced")
		}
	}
	if !flag {
		assert.Fail(t, "Powerful user should have been added")
	}
}

func TestRemoveBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create two test users
	goodUser := createTestAccount(ctx, am, "goodUser")
	ownerKey1, _ := am.GetOwnerKey(ctx, goodUser)
	badUser := createTestAccount(ctx, am, "badUser")
	ownerKey2, _ := am.GetOwnerKey(ctx, badUser)
	am.AddCoin(ctx, goodUser, c2000)
	am.AddCoin(ctx, badUser, c2000)
	// let user register as voter first
	voteManager.AddVoter(ctx, "goodUser", c8000)
	voteManager.AddVoter(ctx, "badUser", c8000)

	// let both users register as validator
	deposit := types.LNO("1200")
	msg1 := NewValidatorDepositMsg("goodUser", deposit, *ownerKey1)
	msg2 := NewValidatorDepositMsg("badUser", deposit, *ownerKey2)
	handler(ctx, msg1)
	handler(ctx, msg2)

	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 2, len(verifyList.OncallValidators))
	assert.Equal(t, 2, len(verifyList.AllValidators))

	valManager.RemoveValidatorFromAllLists(ctx, "badUser")
	verifyList2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 1, len(verifyList2.OncallValidators))
	assert.Equal(t, 1, len(verifyList2.AllValidators))
	assert.Equal(t, goodUser, verifyList2.OncallValidators[0])
	assert.Equal(t, goodUser, verifyList2.AllValidators[0])
}
