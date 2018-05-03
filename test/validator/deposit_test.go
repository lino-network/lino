package validator

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	val "github.com/lino-network/lino/tx/validator"
	vote "github.com/lino-network/lino/tx/vote"
	"github.com/lino-network/lino/types"

	crypto "github.com/tendermint/go-crypto"
)

// test validator deposit
func TestValidatorDeposit(t *testing.T) {
	newAccountTransactionPriv := crypto.GenPrivKeyEd25519()
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	newValidatorPriv := crypto.GenPrivKeyEd25519()

	baseTime := time.Now().Unix() + 100
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), newAccountTransactionPriv, newAccountPostPriv, "5000")

	voteDepositMsg := vote.NewVoterDepositMsg(newAccountName, types.LNO("3000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	// deposit the lowest requirement
	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("1000"), newValidatorPriv.PubKey())
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)

	// deposit as the highest validator
	valDepositMsg = val.NewValidatorDepositMsg(
		newAccountName, types.LNO("1"), newValidatorPriv.PubKey())
	test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)
}
