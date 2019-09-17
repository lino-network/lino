package validator

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	val "github.com/lino-network/lino/x/validator"
	vote "github.com/lino-network/lino/x/vote"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// test normal revoke
func TestValidatorRevoke(t *testing.T) {
	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	newValidatorPriv := secp256k1.GenPrivKey()

	baseT := time.Now().Add(3600 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	test.CreateAccount(t, newAccountName, lb, 0,
		newAccountResetPriv, newAccountTransactionPriv, newAccountAppPriv, "500000")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
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
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(49999*types.Decimals))
	// check the first coin return
	test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalSec+1)
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(7142757143))

	// will get all coins back after the freezing period
	for i := int64(1); i < test.CoinReturnTimes; i++ {
		test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalSec*(i+1)+1)
	}
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(199999*types.Decimals))

	// won't get extra coins in the future
	test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalSec*(test.CoinReturnTimes+1)+1)
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(199999*types.Decimals))
}
