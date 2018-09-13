package vote

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	val "github.com/lino-network/lino/x/validator"
	vote "github.com/lino-network/lino/x/vote"
)

func TestVoterRevoke(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	newValidatorPriv := secp256k1.GenPrivKey()

	delegator1TransactionPriv := secp256k1.GenPrivKey()
	delegator2TransactionPriv := secp256k1.GenPrivKey()
	delegator1Name := "delegator1"
	delegator2Name := "delegator2"

	// to recover the coin day
	baseTime := time.Now().Unix() + 7200
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, secp256k1.GenPrivKey(), "500000")
	test.CreateAccount(t, delegator1Name, lb, 1,
		secp256k1.GenPrivKey(), delegator1TransactionPriv, secp256k1.GenPrivKey(), "210100")
	test.CreateAccount(t, delegator2Name, lb, 2,
		secp256k1.GenPrivKey(), delegator2TransactionPriv, secp256k1.GenPrivKey(), "70100")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("150000"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	// let delegator delegate coins to voter
	delegateMsg := vote.NewDelegateMsg(delegator1Name, newAccountName, types.LNO("210000"))
	delegateMsg2 := vote.NewDelegateMsg(delegator2Name, newAccountName, types.LNO("70000"))

	test.SignCheckDeliver(t, lb, delegateMsg, 0, true, delegator1TransactionPriv, baseTime)
	test.SignCheckDeliver(t, lb, delegateMsg2, 0, true, delegator2TransactionPriv, baseTime)

	// delegator can withdraw coins
	delegatorWithdrawMsg := vote.NewDelegatorWithdrawMsg(delegator1Name, newAccountName,
		types.LNO("70000"))
	test.SignCheckDeliver(t, lb, delegatorWithdrawMsg, 1, true, delegator1TransactionPriv, baseTime)

	//validators can stake out after revoking validator candidancy
	stakeOutMsg := vote.NewStakeOutMsg(newAccountName, types.LNO("300000"))
	valRevokeMsg := val.NewValidatorRevokeMsg(newAccountName)
	test.SignCheckDeliver(t, lb, valRevokeMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.SimulateOneBlock(lb, baseTime)
	test.SignCheckDeliver(t, lb, stakeOutMsg, 3, true, newAccountTransactionPriv, baseTime)

	// check delegator withdraw first coin return
	test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalSec+1)
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(11428471429))
	test.CheckBalance(t, delegator1Name, lb, types.NewCoinFromInt64(30099*types.Decimals))
	test.CheckBalance(t, delegator2Name, lb, types.NewCoinFromInt64(10099*types.Decimals))

	// check balance after freezing period
	for i := int64(1); i < test.CoinReturnTimes; i++ {
		test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalSec*(i+1)+1)
	}
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(499999*types.Decimals))
	test.CheckBalance(t, delegator1Name, lb, types.NewCoinFromInt64(210099*types.Decimals))
	test.CheckBalance(t, delegator2Name, lb, types.NewCoinFromInt64(70099*types.Decimals))
}
