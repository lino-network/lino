package validator

import (
	"strconv"
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestRegisterBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	deposit := coinToString(valParam.ValidatorMinCommittingDeposit)
	valKey := secp256k1.GenPrivKey().PubKey()
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetSavingFromBank(ctx, user1)
	assert.Equal(t, acc1Balance, minBalance)
	assert.Equal(t, true, valManager.DoesValidatorExist(ctx, user1))

	// now user1 should be the only validator
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, verifyList.LowestPower, valParam.ValidatorMinCommittingDeposit)
	assert.Equal(t, 1, len(verifyList.OncallValidators))
	assert.Equal(t, 1, len(verifyList.AllValidators))
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// make sure the validator's account info (power&pubKey) is correct
	verifyAccount, _ := valManager.storage.GetValidator(ctx, user1)
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit, verifyAccount.Deposit)
	assert.Equal(t, valKey, verifyAccount.PubKey)
}

func TestRegisterFeeNotEnough(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit).Plus(valParam.ValidatorMinVotingDeposit))

	// let user1 register as validator
	deposit := coinToString(valParam.ValidatorMinCommittingDeposit.Minus(types.NewCoinFromInt64(1000)))
	valKey := secp256k1.GenPrivKey().PubKey()
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")

	result := handler(ctx, msg)
	assert.Equal(t, ErrInsufficientDeposit().Result(), result)

	// let user register as voter
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	result2 := handler(ctx, msg)
	assert.Equal(t, ErrInsufficientDeposit().Result(), result2)

	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 0, len(verifyList.OncallValidators))
	assert.Equal(t, 0, len(verifyList.AllValidators))
}

func TestRevokeBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := secp256k1.GenPrivKey().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommittingDeposit)
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
	assert.Equal(t, model.ErrValidatorNotFound().Result(), result)
}

// this is the same situation as we find Byzantine and replace the Byzantine
func TestRevokeOncallValidatorAndSubstitutionExists(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(100000 * types.Decimals)

	// create 21 test users
	users := make([]types.AccountKey, 24)
	valKeys := make([]crypto.PubKey, 24)
	for i := 0; i < 24; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1), minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

		// let user register as voter first
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), valParam.ValidatorMinVotingDeposit)

		// they will deposit min committing deposit + 10,20,30...200, 210, 220, 230, 240
		valMinCommitDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((i+1)*10) + valMinCommitDeposit/types.Decimals
		deposit := types.LNO(strconv.FormatInt(num, 10))
		valKeys[i] = secp256k1.GenPrivKey().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 21, len(lst.OncallValidators))
	assert.Equal(t, 24, len(lst.AllValidators))
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(40*types.Decimals)), lst.LowestPower)
	assert.Equal(t, users[3], lst.LowestValidator)

	// lowest validator depoist coins will change the ranks
	deposit := types.LNO("15")
	msg := NewValidatorDepositMsg("user4", deposit, valKeys[3], "")
	result := handler(ctx, msg)

	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result)
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(50*types.Decimals)), lst2.LowestPower)
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
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(50*types.Decimals)), lst3.LowestPower)
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
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(30*types.Decimals)), lst4.LowestPower)
	assert.Equal(t, users[2], lst4.LowestValidator)
	assert.Equal(t, 22, len(lst4.AllValidators))
}

func TestRevokeAndDepositAgain(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit).Plus(valParam.ValidatorMinCommittingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := secp256k1.GenPrivKey().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommittingDeposit)
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
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := secp256k1.GenPrivKey().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommittingDeposit)
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
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := secp256k1.GenPrivKey().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommittingDeposit)
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetSavingFromBank(ctx, user1)
	assert.Equal(t, acc1Balance, minBalance)
	assert.Equal(t, true, valManager.DoesValidatorExist(ctx, user1))

	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit, verifyList.LowestPower)
	assert.Equal(t, 1, len(verifyList.OncallValidators))
	assert.Equal(t, 1, len(verifyList.AllValidators))
	assert.Equal(t, user1, verifyList.OncallValidators[0])
	assert.Equal(t, user1, verifyList.AllValidators[0])

	// check deposit and power is correct
	validator, _ := valManager.storage.GetValidator(ctx, user1)
	assert.Equal(t, true, validator.Deposit.IsEqual(valParam.ValidatorMinCommittingDeposit))
}

func TestCommittingDepositExceedVotingDeposit(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create test user
	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(1000 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinVotingDeposit))

	// let user1 register as voter first
	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := secp256k1.GenPrivKey().PubKey()
	deposit := coinToString(valParam.ValidatorMinVotingDeposit.Plus(types.NewCoinFromInt64(2 * types.Decimals)))
	msg := NewValidatorDepositMsg("user1", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, ErrUnbalancedAccount().Result(), result)
}

func TestDepositWithoutLinoAccount(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)
	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)

	valKey := secp256k1.GenPrivKey().PubKey()
	msg := NewValidatorDepositMsg("qwqwndqwnd", coinToString(valParam.ValidatorMinWithdraw), valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountNotFound().Result(), result)
}

func TestValidatorReplacement(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(100000 * types.Decimals)

	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1), minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), valParam.ValidatorMinVotingDeposit)

		// they will deposit min committing deposit + 10,20,30...200, 210, 220, 230, 240
		valMinCommitDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((i+1)*10) + valMinCommitDeposit/types.Decimals
		deposit := types.LNO(strconv.FormatInt(num, 10))
		valKeys[i] = secp256k1.GenPrivKey().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// check validator list, the lowest power is 10
	verifyList, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, true,
		verifyList.LowestPower.IsEqual(valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(10*types.Decimals))))
	assert.Equal(t, users[0], verifyList.LowestValidator)
	assert.Equal(t, 21, len(verifyList.OncallValidators))
	assert.Equal(t, 21, len(verifyList.AllValidators))

	// create a user failed to join oncall validator list (not enough power)
	createTestAccount(ctx, am, "noPowerUser", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
	voteManager.AddVoter(ctx, "noPowerUser", valParam.ValidatorMinVotingDeposit)

	// let user1 register as validator
	valKey := secp256k1.GenPrivKey().PubKey()
	deposit := coinToString(valParam.ValidatorMinCommittingDeposit)
	msg := NewValidatorDepositMsg("noPowerUser", deposit, valKey, "")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	//check the user hasn't been added to oncall validators but in the pool
	verifyList2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result)
	assert.Equal(t, true,
		verifyList2.LowestPower.IsEqual(valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(10*types.Decimals))))
	assert.Equal(t, users[0], verifyList2.LowestValidator)
	assert.Equal(t, 21, len(verifyList2.OncallValidators))
	assert.Equal(t, 22, len(verifyList2.AllValidators))

	// create a user success to join oncall validator list
	createTestAccount(ctx, am, "powerfulUser", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
	// let user register as voter first
	voteManager.AddVoter(ctx, "powerfulUser", valParam.ValidatorMinVotingDeposit)

	//check the user has been added to oncall validators and in the pool
	valKey = secp256k1.GenPrivKey().PubKey()
	deposit = coinToString(valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(88 * types.Decimals)))
	msg = NewValidatorDepositMsg("powerfulUser", deposit, valKey, "")
	result = handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	verifyList3, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, true,
		verifyList3.LowestPower.IsEqual(valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(20*types.Decimals))))
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

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	goodUser := createTestAccount(ctx, am, "goodUser", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
	createTestAccount(ctx, am, "badUser", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

	valKey1 := secp256k1.GenPrivKey().PubKey()
	valKey2 := secp256k1.GenPrivKey().PubKey()

	voteManager.AddVoter(ctx, "goodUser", valParam.ValidatorMinVotingDeposit)
	voteManager.AddVoter(ctx, "badUser", valParam.ValidatorMinVotingDeposit)

	// let both users register as validator
	msg1 := NewValidatorDepositMsg("goodUser", coinToString(valParam.ValidatorMinCommittingDeposit), valKey1, "")
	msg2 := NewValidatorDepositMsg("badUser", coinToString(valParam.ValidatorMinCommittingDeposit), valKey2, "")
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

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
	createTestAccount(ctx, am, "user2", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

	valKey1 := secp256k1.GenPrivKey().PubKey()

	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)
	voteManager.AddVoter(ctx, "user2", valParam.ValidatorMinVotingDeposit)

	// let both users register as validator
	msg1 := NewValidatorDepositMsg("user1", coinToString(valParam.ValidatorMinCommittingDeposit), valKey1, "")
	msg2 := NewValidatorDepositMsg("user2", coinToString(valParam.ValidatorMinCommittingDeposit), valKey1, "")
	handler(ctx, msg1)

	result2 := handler(ctx, msg2)
	assert.Equal(t, ErrValidatorPubKeyAlreadyExist().Result(), result2)

}

func TestAddFrozenMoney(t *testing.T) {
	ctx, am, valManager, _, gm := setupTest(t, 0)
	valManager.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user := createTestAccount(ctx, am, "user", minBalance)

	testCases := []struct {
		testName               string
		times                  int64
		interval               int64
		returnedCoin           types.Coin
		expectedFrozenListLen  int
		expectedFrozenMoney    types.Coin
		expectedFrozenTimes    int64
		expectedFrozenInterval int64
	}{
		{
			testName:               "return coin to user",
			times:                  10,
			interval:               2,
			returnedCoin:           types.NewCoinFromInt64(100),
			expectedFrozenListLen:  1,
			expectedFrozenMoney:    types.NewCoinFromInt64(100),
			expectedFrozenTimes:    10,
			expectedFrozenInterval: 2,
		},
		{
			testName:               "return coin to user again",
			times:                  100000,
			interval:               20000,
			returnedCoin:           types.NewCoinFromInt64(100000),
			expectedFrozenListLen:  2,
			expectedFrozenMoney:    types.NewCoinFromInt64(100000),
			expectedFrozenTimes:    100000,
			expectedFrozenInterval: 20000,
		},
	}

	for _, tc := range testCases {
		err := returnCoinTo(
			ctx, "user", gm, am, tc.times, tc.interval, tc.returnedCoin)
		if err != nil {
			t.Errorf("%s: failed to return coin, got err %v", tc.testName, err)
		}

		lst, err := am.GetFrozenMoneyList(ctx, user)
		if err != nil {
			t.Errorf("%s: failed to get frozen money list, got err %v", tc.testName, err)
		}
		if len(lst) != tc.expectedFrozenListLen {
			t.Errorf("%s: diff list len, got %v, want %v", tc.testName, len(lst), tc.expectedFrozenListLen)
		}
		if !lst[len(lst)-1].Amount.IsEqual(tc.expectedFrozenMoney) {
			t.Errorf("%s: diff amount, got %v, want %v", tc.testName, lst[len(lst)-1].Amount, tc.expectedFrozenMoney)
		}
		if lst[len(lst)-1].Times != tc.expectedFrozenTimes {
			t.Errorf("%s: diff times, got %v, want %v", tc.testName, lst[len(lst)-1].Times, tc.expectedFrozenTimes)
		}
		if lst[len(lst)-1].Interval != tc.expectedFrozenInterval {
			t.Errorf("%s: diff interval, got %v, want %v", tc.testName, lst[len(lst)-1].Interval, tc.expectedFrozenInterval)
		}
	}
}
