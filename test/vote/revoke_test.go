package vote

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	linotypes "github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	// val "github.com/lino-network/lino/x/validator"
	types "github.com/lino-network/lino/x/vote/types"
)

func TestVoterRevoke(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	// newValidatorPriv := secp256k1.GenPrivKey()

	delegator1TransactionPriv := secp256k1.GenPrivKey()
	delegator2TransactionPriv := secp256k1.GenPrivKey()
	delegator1Name := "delegator1"
	delegator2Name := "delegator2"

	// to recover the coin day
	baseT := time.Unix(0,0).Add(7200 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, secp256k1.GenPrivKey(), "500000")
	test.CreateAccount(t, delegator1Name, lb, 1,
		secp256k1.GenPrivKey(), delegator1TransactionPriv, secp256k1.GenPrivKey(), "210100")
	test.CreateAccount(t, delegator2Name, lb, 2,
		secp256k1.GenPrivKey(), delegator2TransactionPriv, secp256k1.GenPrivKey(), "70100")

	voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	// valDepositMsg := val.NewValidatorDepositMsg(
	// 	newAccountName, linotypes.LNO("150000"), newValidatorPriv.PubKey(), "")
	// test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	// // let delegator delegate coins to voter
	// delegateMsg := types.NewDelegateMsg(delegator1Name, newAccountName, linotypes.LNO("210000"))
	// delegateMsg2 := types.NewDelegateMsg(delegator2Name, newAccountName, linotypes.LNO("70000"))

	// test.SignCheckDeliver(t, lb, delegateMsg, 0, true, delegator1TransactionPriv, baseTime)
	// test.SignCheckDeliver(t, lb, delegateMsg2, 0, true, delegator2TransactionPriv, baseTime)

	// // delegator can withdraw coins
	// delegatorWithdrawMsg := types.NewDelegatorWithdrawMsg(delegator1Name, newAccountName,
	// 	linotypes.LNO("70000"))
	// test.SignCheckDeliver(t, lb, delegatorWithdrawMsg, 1, true, delegator1TransactionPriv, baseTime)

	// //all validators cannot revoke voter candidancy
	// stakeOutMsg := types.NewStakeOutMsg(newAccountName, linotypes.LNO("300000"))
	// test.SimulateOneBlock(lb, baseTime)
	// test.SignCheckDeliver(t, lb, stakeOutMsg, 2, false, newAccountTransactionPriv, baseTime)

	// //validators can stake out after revoking validator candidancy
	// valRevokeMsg := val.NewValidatorRevokeMsg(newAccountName)
	// test.SignCheckDeliver(t, lb, valRevokeMsg, 3, true, newAccountTransactionPriv, baseTime)
	// test.SimulateOneBlock(lb, baseTime)
	// test.SignCheckDeliver(t, lb, stakeOutMsg, 4, true, newAccountTransactionPriv, baseTime)

	// // check delegator withdraw first coin return
	// test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalSec+1)
	// test.CheckBalance(t, newAccountName, lb, linotypes.NewCoinFromInt64(11428471429))
	// test.CheckBalance(t, delegator1Name, lb, linotypes.NewCoinFromInt64(10099*linotypes.Decimals))
	// test.CheckBalance(t, delegator2Name, lb, linotypes.NewCoinFromInt64(99*linotypes.Decimals))

	// // check balance after freezing period
	// for i := int64(1); i < test.CoinReturnTimes; i++ {
	// 	test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalSec*(i+1)+1)
	// }
	// test.CheckBalance(t, newAccountName, lb, linotypes.NewCoinFromInt64(499999*linotypes.Decimals))
	// test.CheckBalance(t, delegator1Name, lb, linotypes.NewCoinFromInt64(70099*linotypes.Decimals))
	// test.CheckBalance(t, delegator2Name, lb, linotypes.NewCoinFromInt64(99*linotypes.Decimals))
}
