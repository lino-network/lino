package validator

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/stretchr/testify/assert"
)

var (
	c0    = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}}
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

	vm.SetValidatorList(ctx, ValidatorListKey, lst)

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
	verifyList, _ := vm.GetValidatorList(ctx, ValidatorListKey)
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

	vm.SetValidatorList(ctx, ValidatorListKey, lst)

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
