package vote

//
// import (
// 	"strconv"
// 	"testing"
//
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	acc "github.com/lino-network/lino/tx/account"
// 	"github.com/lino-network/lino/types"
// 	"github.com/stretchr/testify/assert"
// )
//
// var (
// 	l0    = types.LNO(sdk.NewRat(0))
// 	l10   = types.LNO(sdk.NewRat(10))
// 	l11   = types.LNO(sdk.NewRat(11))
// 	l20   = types.LNO(sdk.NewRat(20))
// 	l21   = types.LNO(sdk.NewRat(21))
// 	l100  = types.LNO(sdk.NewRat(100))
// 	l200  = types.LNO(sdk.NewRat(200))
// 	l400  = types.LNO(sdk.NewRat(400))
// 	l1000 = types.LNO(sdk.NewRat(1000))
// 	l1011 = types.LNO(sdk.NewRat(1011))
// 	l1021 = types.LNO(sdk.NewRat(1021))
// 	l1022 = types.LNO(sdk.NewRat(1022))
// 	l1600 = types.LNO(sdk.NewRat(1600))
// 	l1800 = types.LNO(sdk.NewRat(1800))
// 	l1900 = types.LNO(sdk.NewRat(1900))
// 	l2000 = types.LNO(sdk.NewRat(2000))
//
// 	c0    = types.Coin{0 * types.Decimals}
// 	c10   = types.Coin{10 * types.Decimals}
// 	c11   = types.Coin{11 * types.Decimals}
// 	c20   = types.Coin{20 * types.Decimals}
// 	c21   = types.Coin{21 * types.Decimals}
// 	c100  = types.Coin{100 * types.Decimals}
// 	c200  = types.Coin{200 * types.Decimals}
// 	c400  = types.Coin{400 * types.Decimals}
// 	c1000 = types.Coin{1000 * types.Decimals}
// 	c1011 = types.Coin{1011 * types.Decimals}
// 	c1021 = types.Coin{1021 * types.Decimals}
// 	c1022 = types.Coin{1022 * types.Decimals}
// 	c1600 = types.Coin{1600 * types.Decimals}
// 	c1800 = types.Coin{1800 * types.Decimals}
// 	c1900 = types.Coin{1900 * types.Decimals}
// 	c2000 = types.Coin{2000 * types.Decimals}
// )
//
// func TestRegisterBasic(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create two test users
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, sdk.Result{}, result)
//
// 	// check acc1's money has been withdrawn
// 	acc1Balance, _ := acc1.GetBankBalance(ctx)
// 	assert.Equal(t, acc1Balance, c400)
// 	assert.Equal(t, true, vm.IsValidatorExist(ctx, acc.AccountKey("user1")))
//
// 	// now user1 should be the only validator (WOW, dictator!)
// 	verifyList, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, verifyList.LowestPower, c1600)
// 	assert.Equal(t, 1, len(verifyList.OncallValidators))
// 	assert.Equal(t, 1, len(verifyList.AllValidators))
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.OncallValidators[0])
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.AllValidators[0])
//
// 	// make sure the validator's account info (power&pubKey) is correct
// 	verifyAccount, _ := vm.GetValidator(ctx, acc.AccountKey("user1"))
// 	assert.Equal(t, c1600.Amount, verifyAccount.ABCIValidator.GetPower())
// 	assert.Equal(t, ownerKey.Bytes(), verifyAccount.ABCIValidator.GetPubKey())
// }
//
// func TestRegisterFeeNotEnough(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create test user
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l400, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, ErrRegisterFeeNotEnough().Result(), result)
//
// 	verifyList, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, 0, len(verifyList.OncallValidators))
// 	assert.Equal(t, 0, len(verifyList.AllValidators))
// }
//
// func TestRevokeBasic(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create two test users
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, sdk.Result{}, result)
//
// 	// now user1 should be the only validator
// 	verifyList, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.OncallValidators[0])
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.AllValidators[0])
//
// 	// let user1 revoke candidancy
// 	msg2 := NewValidatorRevokeMsg("user1")
// 	result2 := handler(ctx, msg2)
// 	assert.Equal(t, sdk.Result{}, result2)
//
// 	verifyList2, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, 0, len(verifyList2.OncallValidators))
// 	assert.Equal(t, 0, len(verifyList2.AllValidators))
//
// }
//
// func TestRevokeNonExistUser(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// let user1(not exists) revoke candidancy
// 	msg2 := NewValidatorRevokeMsg("user1")
// 	result2 := handler(ctx, msg2)
// 	assert.Equal(t, ErrGetValidator().Result(), result2)
// }
//
// // ming zi yao chang <-. <-
// func TestRevokeTwiceWontChangeFreezingPeriod(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create user
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, sdk.Result{}, result)
//
// 	// let user1 revoke candidancy
// 	msg2 := NewValidatorRevokeMsg("user1")
// 	result2 := handler(ctx, msg2)
// 	assert.Equal(t, sdk.Result{}, result2)
//
// 	// check withdraw available time is correct
// 	val, _ := vm.GetValidator(ctx, acc.AccountKey("user1"))
// 	assert.Equal(t, types.ValidatorWithdrawFreezingPeriod, val.WithdrawAvailableAt)
//
// 	// adjust block height and revoke again
// 	ctx.WithBlockHeight(800)
// 	result3 := handler(ctx, msg2)
// 	assert.Equal(t, ErrNotInTheList().Result(), result3)
// 	val2, _ := vm.GetValidator(ctx, acc.AccountKey("user1"))
// 	assert.Equal(t, types.ValidatorWithdrawFreezingPeriod, val2.WithdrawAvailableAt)
//
// }
//
// // this is the same situation as we find Byzantine and replace the Byzantine
// func TestRevokeOncallValidatorAndSubstitutionExists(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create 21 test users
// 	users := make([]*acc.Account, 24)
// 	for i := 0; i < 24; i++ {
// 		users[i] = createTestAccount(ctx, lam, "user"+strconv.Itoa(i+1))
// 		users[i].AddCoin(ctx, c2000)
// 		users[i].Apply(ctx)
// 		// they will deposit 1000 + 10,20,30...200, 210, 220, 230, 240
// 		deposit := types.LNO(sdk.NewRat(int64((i+1)*10) + int64(1000)))
// 		ownerKey, _ := users[i].GetOwnerKey(ctx)
// 		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, *ownerKey)
// 		result := handler(ctx, msg)
// 		assert.Equal(t, sdk.Result{}, result)
// 	}
//
// 	lst, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, 21, len(lst.OncallValidators))
// 	assert.Equal(t, 24, len(lst.AllValidators))
// 	assert.Equal(t, types.Coin{1040 * types.Decimals}, lst.LowestPower)
// 	assert.Equal(t, acc.AccountKey("user4"), lst.LowestValidator)
//
// 	// lowest validator depoist coins will change the ranks
// 	ownerKey, _ := users[4].GetOwnerKey(ctx)
// 	deposit := types.LNO(sdk.NewRat(15))
// 	msg := NewValidatorDepositMsg("user4", deposit, *ownerKey)
// 	result := handler(ctx, msg)
//
// 	lst2, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, sdk.Result{}, result)
// 	assert.Equal(t, types.Coin{1050 * types.Decimals}, lst2.LowestPower)
// 	assert.Equal(t, acc.AccountKey("user5"), lst2.LowestValidator)
//
// 	// now user1, 2, 3 are substitutions
// 	//revoke a non oncall valodator wont change anything related to oncall list
// 	revokeMsg := NewValidatorRevokeMsg("user2")
// 	result2 := handler(ctx, revokeMsg)
// 	assert.Equal(t, sdk.Result{}, result2)
//
// 	lst3, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, types.Coin{1050 * types.Decimals}, lst3.LowestPower)
// 	assert.Equal(t, acc.AccountKey("user5"), lst3.LowestValidator)
// 	assert.Equal(t, 23, len(lst3.AllValidators))
//
// 	// now only user1(power1010) and user3(power1030) are substitutions
// 	// the lowest oncall user is user5 with 1050 power
// 	// revoke user6 (could be byzantine) will make user3 (power 1030) join oncall
// 	// list become the lowest validator
// 	revokeMsg2 := NewValidatorRevokeMsg("user6")
// 	result3 := handler(ctx, revokeMsg2)
// 	assert.Equal(t, sdk.Result{}, result3)
//
// 	lst4, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, types.Coin{1030 * types.Decimals}, lst4.LowestPower)
// 	assert.Equal(t, acc.AccountKey("user3"), lst4.LowestValidator)
// 	assert.Equal(t, 22, len(lst4.AllValidators))
// }
//
// func TestRevokeAndDepositAgain(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create user
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, sdk.Result{}, result)
//
// 	lst, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, 1, len(lst.AllValidators))
// 	assert.Equal(t, 1, len(lst.OncallValidators))
//
// 	// let user1 revoke candidancy
// 	msg2 := NewValidatorRevokeMsg("user1")
// 	result2 := handler(ctx, msg2)
// 	assert.Equal(t, sdk.Result{}, result2)
//
// 	lstEmpty, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, 0, len(lstEmpty.AllValidators))
// 	assert.Equal(t, 0, len(lstEmpty.OncallValidators))
//
// 	// deposit again
// 	msg3 := NewValidatorDepositMsg("user1", l100, *ownerKey)
// 	result3 := handler(ctx, msg3)
//
// 	lst2, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, sdk.Result{}, result3)
// 	assert.Equal(t, 1, len(lst2.AllValidators))
// 	assert.Equal(t, 1, len(lst2.OncallValidators))
// }
//
// func TestWithdrawBasic(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create test user
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, sdk.Result{}, result)
//
// 	// now user1 should be the only validator
// 	verifyList, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.OncallValidators[0])
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.AllValidators[0])
//
// 	// user1 cannot withdraw if is oncall validator
// 	withdrawMsg := NewValidatorWithdrawMsg("user1")
// 	result2 := handler(ctx, withdrawMsg)
// 	assert.Equal(t, ErrDepositNotAvailable().Result(), result2)
//
// 	// user1 cannot withdraw if in the freezing period
// 	revokeMsg := NewValidatorRevokeMsg("user1")
// 	handler(ctx, revokeMsg)
// 	result3 := handler(ctx, withdrawMsg)
// 	assert.Equal(t, ErrDepositNotAvailable().Result(), result3)
// 	acc1Balance, _ := acc1.GetBankBalance(ctx)
// 	assert.Equal(t, true, acc1Balance.IsEqual(c400))
//
// 	// user1 can withdraw if the block height has increased 1000
// 	acc1.Apply(ctx)
// 	ctx = ctx.WithBlockHeight(int64(types.ValidatorWithdrawFreezingPeriod))
// 	result4 := handler(ctx, withdrawMsg)
// 	assert.Equal(t, sdk.Result{}, result4)
//
// 	acc1BalanceNew, _ := acc1.GetBankBalance(ctx)
// 	assert.Equal(t, true, acc1BalanceNew.IsEqual(c2000))
// }
//
// func TestWithdrawTwice(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create two test users
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, sdk.Result{}, result)
//
// 	// withdraw first time
// 	withdrawMsg := NewValidatorWithdrawMsg("user1")
// 	revokeMsg := NewValidatorRevokeMsg("user1")
//
// 	handler(ctx, revokeMsg)
// 	ctx = ctx.WithBlockHeight(int64(types.ValidatorWithdrawFreezingPeriod))
// 	result2 := handler(ctx, withdrawMsg)
// 	assert.Equal(t, sdk.Result{}, result2)
//
// 	// withdraw again
// 	result3 := handler(ctx, withdrawMsg)
// 	assert.Equal(t, ErrNoDeposit().Result(), result3)
//
// 	acc1Balance, _ := acc1.GetBankBalance(ctx)
// 	assert.Equal(t, true, acc1Balance.IsEqual(c2000))
//
// }
//
// func TestDepositBasic(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create test user
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	// let user1 register as validator
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("user1", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, sdk.Result{}, result)
//
// 	depositMsg := NewValidatorDepositMsg("user1", l200, *ownerKey)
// 	result2 := handler(ctx, depositMsg)
// 	assert.Equal(t, sdk.Result{}, result2)
//
// 	// check acc1's money has been withdrawn
// 	acc1Balance, _ := acc1.GetBankBalance(ctx)
// 	assert.Equal(t, acc1Balance, c200)
// 	assert.Equal(t, true, vm.IsValidatorExist(ctx, acc.AccountKey("user1")))
//
// 	// check the lowest power is 1800 now (1600 + 200)
// 	verifyList, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, c1800, verifyList.LowestPower)
// 	assert.Equal(t, 1, len(verifyList.OncallValidators))
// 	assert.Equal(t, 1, len(verifyList.AllValidators))
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.OncallValidators[0])
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList.AllValidators[0])
//
// 	// check deposit and power is correct
// 	validator, _ := vm.GetValidator(ctx, acc.AccountKey("user1"))
// 	assert.Equal(t, true, validator.Deposit.IsEqual(c1800))
// }
//
// func TestDepositWithoutLinoAccount(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// let user1 register as validator
// 	acc1 := createTestAccount(ctx, lam, "user1")
// 	ownerKey, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("qwqwndqwnd", l1600, *ownerKey)
// 	result := handler(ctx, msg)
// 	assert.Equal(t, ErrUsernameNotFound().Result(), result)
// }
//
// func TestValidatorReplacement(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create 21 test users
// 	users := make([]*acc.Account, 21)
// 	for i := 0; i < 21; i++ {
// 		users[i] = createTestAccount(ctx, lam, "user"+strconv.Itoa(i))
// 		users[i].AddCoin(ctx, c2000)
// 		users[i].Apply(ctx)
// 		// they will deposit 10,20,30...200, 210
// 		deposit := types.LNO(sdk.NewRat(int64((i+1)*10) + int64(1001)))
// 		ownerKey, _ := users[i].GetOwnerKey(ctx)
// 		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i), deposit, *ownerKey)
// 		result := handler(ctx, msg)
// 		assert.Equal(t, sdk.Result{}, result)
// 	}
//
// 	// check validator list, the lowest power is 10
// 	verifyList, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, true, verifyList.LowestPower.IsEqual(c1011))
// 	assert.Equal(t, acc.AccountKey("user0"), verifyList.LowestValidator)
// 	assert.Equal(t, 21, len(verifyList.OncallValidators))
// 	assert.Equal(t, 21, len(verifyList.AllValidators))
//
// 	// create a user failed to join oncall validator list (not enough power)
// 	acc1 := createTestAccount(ctx, lam, "noPowerUser")
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
//
// 	//check the user hasn't been added to oncall validators but in the pool
// 	deposit := types.LNO(sdk.NewRat(1005))
// 	ownerKey1, _ := acc1.GetOwnerKey(ctx)
// 	msg := NewValidatorDepositMsg("noPowerUser", deposit, *ownerKey1)
// 	result := handler(ctx, msg)
//
// 	verifyList2, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, sdk.Result{}, result)
// 	assert.Equal(t, true, verifyList.LowestPower.IsEqual(c1011))
// 	assert.Equal(t, acc.AccountKey("user0"), verifyList.LowestValidator)
// 	assert.Equal(t, 21, len(verifyList2.OncallValidators))
// 	assert.Equal(t, 22, len(verifyList2.AllValidators))
//
// 	// create a user success to join oncall validator list
// 	acc2 := createTestAccount(ctx, lam, "powerfulUser")
// 	acc2.AddCoin(ctx, c2000)
// 	acc2.Apply(ctx)
//
// 	//check the user has been added to oncall validators and in the pool
// 	deposit2 := types.LNO(sdk.NewRat(1088))
// 	ownerKey2, _ := acc2.GetOwnerKey(ctx)
// 	msg2 := NewValidatorDepositMsg("powerfulUser", deposit2, *ownerKey2)
// 	result2 := handler(ctx, msg2)
//
// 	verifyList3, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, sdk.Result{}, result2)
// 	assert.Equal(t, true, verifyList3.LowestPower.IsEqual(c1021))
// 	assert.Equal(t, acc.AccountKey("user1"), verifyList3.LowestValidator)
// 	assert.Equal(t, 21, len(verifyList3.OncallValidators))
// 	assert.Equal(t, 23, len(verifyList3.AllValidators))
//
// 	// check user0 has been replaced, and powerful user has been added
// 	flag := false
// 	for _, username := range verifyList3.OncallValidators {
// 		if username == "powerfulUser" {
// 			flag = true
// 		}
// 		if username == "user0" {
// 			assert.Fail(t, "User0 should have been replaced")
// 		}
// 	}
// 	if !flag {
// 		assert.Fail(t, "Powerful user should have been added")
// 	}
// }
//
// func TestRemoveBasic(t *testing.T) {
// 	lam := newLinoAccountManager()
// 	vm := newValidatorManager()
// 	ctx := getContext()
// 	handler := NewHandler(vm, lam)
// 	vm.InitGenesis(ctx)
//
// 	// create two test users
// 	acc1 := createTestAccount(ctx, lam, "goodUser")
// 	ownerKey1, _ := acc1.GetOwnerKey(ctx)
// 	acc2 := createTestAccount(ctx, lam, "badUser")
// 	ownerKey2, _ := acc2.GetOwnerKey(ctx)
// 	acc1.AddCoin(ctx, c2000)
// 	acc1.Apply(ctx)
// 	acc2.AddCoin(ctx, c2000)
// 	acc2.Apply(ctx)
//
// 	// let both users register as validator
// 	deposit := types.LNO(sdk.NewRat(1200))
// 	msg1 := NewValidatorDepositMsg("goodUser", deposit, *ownerKey1)
// 	msg2 := NewValidatorDepositMsg("badUser", deposit, *ownerKey2)
// 	handler(ctx, msg1)
// 	handler(ctx, msg2)
//
// 	verifyList, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, 2, len(verifyList.OncallValidators))
// 	assert.Equal(t, 2, len(verifyList.AllValidators))
//
// 	vm.RemoveValidatorFromAllLists(ctx, "badUser")
// 	verifyList2, _ := vm.GetValidatorList(ctx)
// 	assert.Equal(t, 1, len(verifyList2.OncallValidators))
// 	assert.Equal(t, 1, len(verifyList2.AllValidators))
// 	assert.Equal(t, acc.AccountKey("goodUser"), verifyList2.OncallValidators[0])
// 	assert.Equal(t, acc.AccountKey("goodUser"), verifyList2.AllValidators[0])
// }
