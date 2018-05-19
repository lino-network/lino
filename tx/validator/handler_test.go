package validator

import (
	"strconv"
	"testing"

	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestRegisterBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	deposit := coinToString(valParam.ValidatorMinCommitingDeposit)
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetSavingFromBank(ctx, user1)
	assert.Equal(t, acc1Balance, minBalance)
	assert.Equal(t, true, valManager.IsValidatorExist(ctx, user1))

	// now user1 should be the only validator
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, verifyList.LowestPower, valParam.ValidatorMinCommitingDeposit)
	assert.Equal(t, 1, len(verifyList.OncallValidators))
	assert.Equal(t, 1, len(verifyList.AllValidators))
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// make sure the validator's account info (power&pubKey) is correct
	verifyAccount, _ := valManager.storage.GetValidator(ctx, user1)
	assert.Equal(t, valParam.ValidatorMinCommitingDeposit, verifyAccount.Deposit)
	assert.Equal(t, valKey.Bytes(), verifyAccount.ABCIValidator.GetPubKey())
}

func TestRegisterFeeNotEnough(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommitingDeposit).Plus(valParam.ValidatorMinVotingDeposit))

	// let user1 register as validator
	deposit := coinToString(valParam.ValidatorMinCommitingDeposit.Minus(types.NewCoin(1000)))
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")

	result := handler(ctx, msg)
	assert.Equal(t, ErrVotingDepositNotEnough().Result(), result)

	// let user register as voter
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

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

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommitingDeposit)
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
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
	assert.Equal(t, 0, len(verifyList2.OncallValidators))
	assert.Equal(t, 0, len(verifyList2.AllValidators))

}

func TestRevokeNonExistUser(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// let user1(not exists) revoke candidancy
	msg := NewValidatorRevokeMsg("user1")
	result := handler(ctx, msg)
	assert.Equal(t, model.ErrGetValidator().Result(), result)
}

// this is the same situation as we find Byzantine and replace the Byzantine
func TestRevokeOncallValidatorAndSubstitutionExists(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(100000 * types.Decimals)

	// create 21 test users
	users := make([]types.AccountKey, 24)
	valKeys := make([]crypto.PubKey, 24)
	for i := 0; i < 24; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1), minBalance.Plus(valParam.ValidatorMinCommitingDeposit))

		// let user register as voter first
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), valParam.ValidatorMinVotingDeposit)

		// they will deposit min commiting deposit + 10,20,30...200, 210, 220, 230, 240
		num := int64((i+1)*10) + valParam.ValidatorMinCommitingDeposit.ToInt64()/types.Decimals
		deposit := types.LNO(strconv.FormatInt(num, 10))
		valKeys[i] = crypto.GenPrivKeyEd25519().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 21, len(lst.OncallValidators))
	assert.Equal(t, 24, len(lst.AllValidators))
	assert.Equal(t, valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(40*types.Decimals)), lst.LowestPower)
	assert.Equal(t, users[3], lst.LowestValidator)

	// lowest validator depoist coins will change the ranks
	deposit := types.LNO("15")
	msg := NewValidatorDepositMsg("user4", deposit, valKeys[3], "")
	result := handler(ctx, msg)

	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result)
	assert.Equal(t, valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(50*types.Decimals)), lst2.LowestPower)
	assert.Equal(t, users[4], lst2.LowestValidator)

	// now user1, 2, 3 are substitutions
	// user2 can only withdraw 1-20 coins
	withdrawMsg := NewValidatorWithdrawMsg("user2", "100")
	resultWithdraw := handler(ctx, withdrawMsg)
	assert.Equal(t, ErrIllegalWithdraw().Result(), resultWithdraw)

	withdrawMsg2 := NewValidatorWithdrawMsg("user2", coinToString(valParam.ValidatorMinWithdraw))
	resultWithdraw2 := handler(ctx, withdrawMsg2)
	assert.Equal(t, sdk.Result{}, resultWithdraw2)
	//revoke a non oncall valodator wont change anything related to oncall list
	revokeMsg := NewValidatorRevokeMsg("user2")
	result2 := handler(ctx, revokeMsg)
	assert.Equal(t, sdk.Result{}, result2)

	lst3, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(50*types.Decimals)), lst3.LowestPower)
	assert.Equal(t, users[4], lst3.LowestValidator)
	assert.Equal(t, 23, len(lst3.AllValidators))

	// now only user1(min + 10) and user3(min + 30) are substitutions
	// the lowest oncall user is user5 with (min + 50) power
	// revoke user6 (could be byzantine) will make user3 (min + 30) join oncall
	// list become the lowest validator
	revokeMsg2 := NewValidatorRevokeMsg("user6")
	result3 := handler(ctx, revokeMsg2)
	assert.Equal(t, sdk.Result{}, result3)

	lst4, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(30*types.Decimals)), lst4.LowestPower)
	assert.Equal(t, users[2], lst4.LowestValidator)
	assert.Equal(t, 22, len(lst4.AllValidators))
}

func TestRevokeAndDepositAgain(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommitingDeposit).Plus(valParam.ValidatorMinCommitingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommitingDeposit)
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
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
	msg3 := NewValidatorDepositMsg("user1", deposit, valKey, "")
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

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommitingDeposit)
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// now user1 should be the only validator
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// user1 cannot withdraw if is oncall validator
	withdrawMsg := NewValidatorWithdrawMsg("user1", coinToString(valParam.ValidatorMinWithdraw))
	result2 := handler(ctx, withdrawMsg)
	assert.Equal(t, ErrIllegalWithdraw().Result(), result2)

}

func TestDepositBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create test user
	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommitingDeposit)
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetSavingFromBank(ctx, user1)
	assert.Equal(t, acc1Balance, minBalance)
	assert.Equal(t, true, valManager.IsValidatorExist(ctx, user1))

	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, valParam.ValidatorMinCommitingDeposit, verifyList.LowestPower)
	assert.Equal(t, 1, len(verifyList.OncallValidators))
	assert.Equal(t, 1, len(verifyList.AllValidators))
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// check deposit and power is correct
	validator, _ := valManager.storage.GetValidator(ctx, user1)
	assert.Equal(t, true, validator.Deposit.IsEqual(valParam.ValidatorMinCommitingDeposit))
}

func TestCommitingDepositExceedVotingDeposit(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create test user
	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(1000 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinVotingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	deposit := coinToString(valParam.ValidatorMinVotingDeposit.Plus(types.NewCoin(2 * types.Decimals)))
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, ErrCommitingDepositExceedVotingDeposit().Result(), result)
}

func TestDepositWithoutLinoAccount(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)
	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)

	valKey := crypto.GenPrivKeyEd25519().PubKey()
	msg := NewValidatorDepositMsg("qwqwndqwnd", coinToString(valParam.ValidatorMinWithdraw), valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, ErrUsernameNotFound().Result(), result)
}

func TestValidatorReplacement(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoin(100000 * types.Decimals)

	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1), minBalance.Plus(valParam.ValidatorMinCommitingDeposit))
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), valParam.ValidatorMinVotingDeposit)

		// they will deposit min commiting deposit + 10,20,30...200, 210, 220, 230, 240
		num := int64((i+1)*10) + valParam.ValidatorMinCommitingDeposit.ToInt64()/types.Decimals
		deposit := types.LNO(strconv.FormatInt(num, 10))
		valKeys[i] = crypto.GenPrivKeyEd25519().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// check validator list, the lowest power is 10
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, true,
		verifyList.LowestPower.IsEqual(valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(10*types.Decimals))))
	assert.Equal(t, users[0], verifyList.LowestValidator)
	assert.Equal(t, 21, len(verifyList.OncallValidators))
	assert.Equal(t, 21, len(verifyList.AllValidators))

	// create a user failed to join oncall validator list (not enough power)
	createTestAccount(ctx, am, "noPowerUser", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))
	voteManager.AddVoter(ctx, "noPowerUser", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommitingDeposit)
	msg := NewValidatorDepositMsg("noPowerUser", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	//check the user hasn't been added to oncall validators but in the pool
	verifyList2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result)
	assert.Equal(t, true,
		verifyList2.LowestPower.IsEqual(valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(10*types.Decimals))))
	assert.Equal(t, users[0], verifyList2.LowestValidator)
	assert.Equal(t, 21, len(verifyList2.OncallValidators))
	assert.Equal(t, 22, len(verifyList2.AllValidators))

	// create a user success to join oncall validator list
	createTestAccount(ctx, am, "powerfulUser", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))
	// let user register as voter first
	voteManager.AddVoter(ctx, "powerfulUser", valParam.ValidatorMinVotingDeposit)

	//check the user has been added to oncall validators and in the pool
	valKey = crypto.GenPrivKeyEd25519().PubKey()
	deposit = coinToString(valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(88 * types.Decimals)))
	msg = NewValidatorDepositMsg("powerfulUser", deposit, valKey, "")
	result = handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	verifyList3, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, true,
		verifyList3.LowestPower.IsEqual(valParam.ValidatorMinCommitingDeposit.Plus(types.NewCoin(20*types.Decimals))))
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
	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)

	minBalance := types.NewCoin(1 * types.Decimals)
	goodUser := createTestAccount(ctx, am, "goodUser", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))
	createTestAccount(ctx, am, "badUser", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))

	valKey1 := crypto.GenPrivKeyEd25519().PubKey()
	valKey2 := crypto.GenPrivKeyEd25519().PubKey()

	voteManager.AddVoter(ctx, "goodUser", valParam.ValidatorMinVotingDeposit)
	voteManager.AddVoter(ctx, "badUser", valParam.ValidatorMinVotingDeposit)

	// let both users register as validator
	msg1 := NewValidatorDepositMsg("goodUser", coinToString(valParam.ValidatorMinCommitingDeposit), valKey1, "")
	msg2 := NewValidatorDepositMsg("badUser", coinToString(valParam.ValidatorMinCommitingDeposit), valKey2, "")
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

func TestRegisterWithDupKey(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)

	minBalance := types.NewCoin(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))
	createTestAccount(ctx, am, "user2", minBalance.Plus(valParam.ValidatorMinCommitingDeposit))

	valKey1 := crypto.GenPrivKeyEd25519().PubKey()

	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)
	voteManager.AddVoter(ctx, "user2", valParam.ValidatorMinVotingDeposit)

	// let both users register as validator
	msg1 := NewValidatorDepositMsg("user1", coinToString(valParam.ValidatorMinCommitingDeposit), valKey1, "")
	msg2 := NewValidatorDepositMsg("user2", coinToString(valParam.ValidatorMinCommitingDeposit), valKey1, "")
	handler(ctx, msg1)

	result2 := handler(ctx, msg2)
	assert.Equal(t, ErrPubKeyHasBeenRegistered().Result(), result2)

}
