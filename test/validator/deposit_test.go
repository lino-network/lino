package validator

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	val "github.com/lino-network/lino/x/validator"
	vote "github.com/lino-network/lino/x/vote"
)

// test validator deposit
func TestValidatorDeposit(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	newValidatorPriv := secp256k1.GenPrivKey()

	baseTime := time.Now().Unix() + 100
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, newAccountAppPriv, "500000")

	voteDepositMsg := vote.NewVoterDepositMsg(newAccountName, types.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	// deposit the lowest requirement
	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("100000"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)

	// deposit as the highest validator
	valDepositMsg = val.NewValidatorDepositMsg(
		newAccountName, types.LNO("100"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)
}
