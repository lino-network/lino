package validator

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"

	val "github.com/lino-network/lino/x/validator"
	vote "github.com/lino-network/lino/x/vote"
	crypto "github.com/tendermint/tendermint/crypto"
)

// test normal revoke
func TestValidatorRevoke(t *testing.T) {
	newAccountMasterPriv := crypto.GenPrivKeyEd25519()
	newAccountTransactionPriv := crypto.GenPrivKeyEd25519()
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newuser"
	newValidatorPriv := crypto.GenPrivKeyEd25519()

	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		newAccountMasterPriv, newAccountTransactionPriv, crypto.GenPrivKeyEd25519(), newAccountPostPriv, "500000")

	voteDepositMsg := vote.NewVoterDepositMsg(newAccountName, types.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("150000"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
	test.CheckAllValidatorList(t, newAccountName, true, lb)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)

	valRevokeMsg := val.NewValidatorRevokeMsg(newAccountName)
	test.SignCheckDeliver(t, lb, valRevokeMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.CheckAllValidatorList(t, newAccountName, false, lb)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(50000*types.Decimals))
	// check the first coin return
	test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalHr*3600+1)
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(7142857143))

	// will get all coins back after the freezing period
	for i := int64(1); i < test.CoinReturnTimes; i++ {
		test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalHr*3600*(i+1)+1)
	}
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(200000*types.Decimals))

	// won't get extra coins in the future
	test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalHr*3600*(test.CoinReturnTimes+1)+1)
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(200000*types.Decimals))

}
