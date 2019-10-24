package validator

import (
	"strconv"
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	linotypes "github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	// store "github.com/lino-network/lino/x/validator/model"
	valtypes "github.com/lino-network/lino/x/validator/types"
	types "github.com/lino-network/lino/x/vote/types"
	// abci "github.com/tendermint/tendermint/abci/types"
)

func TestVoteStandby(t *testing.T) {
	// testName := "TestVoteStandby"

	// start with 22 genesis validator
	baseT := time.Unix(0, 0).Add(100 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	// add one more validators but not oncall
	for i := 0; i < 1; i++ {
		newAccountResetPriv := secp256k1.GenPrivKey()
		newAccountTransactionPriv := secp256k1.GenPrivKey()

		newValidatorPriv := secp256k1.GenPrivKey()

		newAccountName := "altval"
		newAccountName += strconv.Itoa(i)

		test.CreateAccount(t, newAccountName, lb, uint64(i),
			newAccountResetPriv, newAccountTransactionPriv, "500000")

		voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("200000"))
		test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

		valDepositMsg := valtypes.NewValidatorRegisterMsg(newAccountName, newValidatorPriv.PubKey(), "")
		test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)

		test.CheckOncallValidatorList(t, newAccountName, false, lb)
		test.CheckStandbyValidatorList(t, newAccountName, true, lb)

	}

	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountName := "voter"

	test.CreateAccount(t, newAccountName, lb, uint64(1),
		newAccountResetPriv, newAccountTransactionPriv, "500000")
	voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("100000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	// let voter vote altval
	voteMsg := valtypes.NewVoteValidatorMsg(newAccountName, []string{"altval0"})
	test.SignCheckDeliver(t, lb, voteMsg, 2, true, newAccountTransactionPriv, baseTime)

	// check altval has become oncall and vote increased
	test.CheckOncallValidatorList(t, "altval0", true, lb)
	test.CheckStandbyValidatorList(t, "altval0", false, lb)
	test.CheckReceivedVotes(t, "altval0", linotypes.NewCoinFromInt64(300000*linotypes.Decimals), lb)

	// check val21 is gone
	test.CheckOncallValidatorList(t, "validator21", false, lb)
	test.CheckStandbyValidatorList(t, "validator21", true, lb)

}
