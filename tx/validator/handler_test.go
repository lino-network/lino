package validator

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/stretchr/testify/assert"
)

var (
	c0    = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}}
	c10   = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 10}}
	c11   = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 11}}
	c20   = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 20}}
	c21   = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 21}}
	c100  = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 100}}
	c200  = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 200}}
	c1600 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 1600}}
	c1800 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 1800}}
	c1900 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 1900}}
	c2000 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 2000}}
)

func TestRegisterBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newValidatorManager()
	ctx := getContext()
	handler := NewHandler(vm, lam)

	lst := &ValidatorList{
		LowestPower: sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}},
	}

	vm.SetValidatorList(ctx, lst)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoins(ctx, c2000)
	acc1.Apply(ctx)

	// let user1 register as validator
	deposit := sdk.Coins{sdk.Coin{Denom: "lino", Amount: 200}}
	ownerKey, _ := acc1.GetOwnerKey(ctx)
	msg := NewValidatorRegisterMsg("user1", deposit)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check acc1's money has been withdrawn
	acc1Balance, _ := acc1.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(c1800))
	assert.Equal(t, true, vm.IsValidatorExist(ctx, acc.AccountKey("user1")))

	// now user1 should be the only validator (WOW, dictator!)
	verifyList, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, true, verifyList.LowestPower.IsEqual(c200))
	assert.Equal(t, acc.AccountKey("user1"), verifyList.OncallValidators[0])
	assert.Equal(t, acc.AccountKey("user1"), verifyList.AllValidators[0])
	assert.Equal(t, 1, len(verifyList.OncallValidators))
	assert.Equal(t, 1, len(verifyList.AllValidators))

	// make sure the validator's account info (power&pubKey) is correct
	verifyAccount, _ := vm.GetValidator(ctx, acc.AccountKey("user1"))
	assert.Equal(t, int64(200), verifyAccount.ABCIValidator.GetPower())
	assert.Equal(t, ownerKey.Bytes(), verifyAccount.ABCIValidator.GetPubKey())

}

func TestVoteBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newValidatorManager()
	ctx := getContext()
	handler := NewHandler(vm, lam)

	lst := &ValidatorList{
		LowestPower: sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}},
	}

	vm.SetValidatorList(ctx, lst)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")
	acc1.AddCoins(ctx, c2000)
	acc1.Apply(ctx)

	acc2.AddCoins(ctx, c2000)
	acc2.Apply(ctx)

	// let user1 register as validator
	deposit := sdk.Coins{sdk.Coin{Denom: "lino", Amount: 200}}
	//ownerKey, _ := acc1.GetOwnerKey(ctx)
	msg := NewValidatorRegisterMsg("user1", deposit)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// let user2 vote 2000 to user1
	msgVote := NewVoteMsg("user2", "user1", c200)
	result2 := handler(ctx, msgVote)
	assert.Equal(t, sdk.Result{}, result2)

	// check user1's power has been increased, and user2's money has been withdrawn
	acc1Balance, _ := acc1.GetBankBalance(ctx)
	acc2Balance, _ := acc2.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(c1800))
	assert.Equal(t, true, acc2Balance.IsEqual(c1800))

	verifyAccount, _ := vm.GetValidator(ctx, acc.AccountKey("user1"))
	assert.Equal(t, int64(400), verifyAccount.ABCIValidator.GetPower())

}

func TestValidatorReplacement(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newValidatorManager()
	ctx := getContext()
	handler := NewHandler(vm, lam)

	lst := &ValidatorList{
		LowestPower: sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}},
	}

	vm.SetValidatorList(ctx, lst)

	// create 21 test users
	users := make([]*acc.Account, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, lam, "user"+strconv.Itoa(i))
		users[i].AddCoins(ctx, c2000)
		users[i].Apply(ctx)
		// they will deposit 10,20,30...200, 210
		deposit := sdk.Coins{sdk.Coin{Denom: "lino", Amount: int64((i + 1) * 10)}}
		msg := NewValidatorRegisterMsg("user"+strconv.Itoa(i), deposit)
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// check validator list, the lowest power is 10
	verifyList, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, true, verifyList.LowestPower.IsEqual(c10))
	assert.Equal(t, acc.AccountKey("user0"), verifyList.LowestValidator)
	assert.Equal(t, 21, len(verifyList.OncallValidators))
	assert.Equal(t, 21, len(verifyList.AllValidators))

	// create a user failed to join oncall validator list (not enough power)
	acc1 := createTestAccount(ctx, lam, "noPowerUser")
	acc1.AddCoins(ctx, c2000)
	acc1.Apply(ctx)

	//check the user hasn't been added to oncall validators but in the pool
	deposit := sdk.Coins{sdk.Coin{Denom: "lino", Amount: 5}}
	msg := NewValidatorRegisterMsg("noPowerUser", deposit)
	result := handler(ctx, msg)

	verifyList2, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result)
	assert.Equal(t, true, verifyList.LowestPower.IsEqual(c10))
	assert.Equal(t, acc.AccountKey("user0"), verifyList.LowestValidator)
	assert.Equal(t, 21, len(verifyList2.OncallValidators))
	assert.Equal(t, 22, len(verifyList2.AllValidators))

	// create a user success to join oncall validator list
	acc2 := createTestAccount(ctx, lam, "powerfulUser")
	acc2.AddCoins(ctx, c2000)
	acc2.Apply(ctx)

	//check the user has been added to oncall validators and in the pool
	deposit2 := sdk.Coins{sdk.Coin{Denom: "lino", Amount: 88}}
	msg2 := NewValidatorRegisterMsg("powerfulUser", deposit2)
	result2 := handler(ctx, msg2)

	verifyList3, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result2)
	assert.Equal(t, true, verifyList3.LowestPower.IsEqual(c20))
	assert.Equal(t, acc.AccountKey("user1"), verifyList3.LowestValidator)
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

	// create a user to vote a validator candidate
	acc3 := createTestAccount(ctx, lam, "voter")
	acc3.AddCoins(ctx, c2000)
	acc3.Apply(ctx)

	// let voter vote 11 to user0, and user0 (power21) will replace user1 (power20)
	msgVote := NewVoteMsg("voter", "user0", c11)
	result3 := handler(ctx, msgVote)

	verifyList4, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, sdk.Result{}, result3)
	assert.Equal(t, true, verifyList4.LowestPower.IsEqual(c21))
	assert.Equal(t, acc.AccountKey("user0"), verifyList4.LowestValidator)
	assert.Equal(t, 21, len(verifyList3.OncallValidators))
	assert.Equal(t, 23, len(verifyList3.AllValidators))
}

func TestRemoveBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newValidatorManager()
	ctx := getContext()
	handler := NewHandler(vm, lam)

	lst := &ValidatorList{
		LowestPower: sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}},
	}

	vm.SetValidatorList(ctx, lst)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "goodUser")
	acc2 := createTestAccount(ctx, lam, "badUser")
	acc1.AddCoins(ctx, c2000)
	acc1.Apply(ctx)
	acc2.AddCoins(ctx, c2000)
	acc2.Apply(ctx)

	// let both users register as validator
	deposit := sdk.Coins{sdk.Coin{Denom: "lino", Amount: 200}}
	msg1 := NewValidatorRegisterMsg("goodUser", deposit)
	msg2 := NewValidatorRegisterMsg("badUser", deposit)
	handler(ctx, msg1)
	handler(ctx, msg2)

	verifyList, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, 2, len(verifyList.OncallValidators))
	assert.Equal(t, 2, len(verifyList.AllValidators))

	vm.RemoveValidatorFromAllLists(ctx, "badUser")
	verifyList2, _ := vm.GetValidatorList(ctx)
	assert.Equal(t, 1, len(verifyList2.OncallValidators))
	assert.Equal(t, 1, len(verifyList2.AllValidators))
	assert.Equal(t, acc.AccountKey("goodUser"), verifyList2.OncallValidators[0])
	assert.Equal(t, acc.AccountKey("goodUser"), verifyList2.AllValidators[0])
}
