package vote

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	val "github.com/lino-network/lino/tx/validator"
	vote "github.com/lino-network/lino/tx/vote"
	"github.com/lino-network/lino/types"

	crypto "github.com/tendermint/go-crypto"
)

func TestVoterRevoke(t *testing.T) {
	newAccountTransactionPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	newValidatorPriv := crypto.GenPrivKeyEd25519()

	delegator1TransactionPriv := crypto.GenPrivKeyEd25519()
	delegator2TransactionPriv := crypto.GenPrivKeyEd25519()
	delegator1Name := "delegator1"
	delegator2Name := "delegator2"

	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), newAccountTransactionPriv, crypto.GenPrivKeyEd25519(), "5000")
	test.CreateAccount(t, delegator1Name, lb, 1,
		crypto.GenPrivKeyEd25519(), delegator1TransactionPriv, crypto.GenPrivKeyEd25519(), "2100")
	test.CreateAccount(t, delegator2Name, lb, 2,
		crypto.GenPrivKeyEd25519(), delegator2TransactionPriv, crypto.GenPrivKeyEd25519(), "700")

	voteDepositMsg := vote.NewVoterDepositMsg(newAccountName, types.LNO("3000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("1500"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	// let delegator delegate coins to voter
	delegateMsg := vote.NewDelegateMsg(delegator1Name, newAccountName, types.LNO("2100"))
	delegateMsg2 := vote.NewDelegateMsg(delegator2Name, newAccountName, types.LNO("700"))

	test.SignCheckDeliver(t, lb, delegateMsg, 0, true, delegator1TransactionPriv, baseTime)
	test.SignCheckDeliver(t, lb, delegateMsg2, 0, true, delegator2TransactionPriv, baseTime)

	// delegator can withdraw coins
	delegatorWithdrawMsg := vote.NewDelegatorWithdrawMsg(delegator1Name, newAccountName,
		types.LNO("700"))
	test.SignCheckDeliver(t, lb, delegatorWithdrawMsg, 1, true, delegator1TransactionPriv, baseTime)
	//all validators cannot revoke voter candidancy
	voterRevokeMsg := vote.NewVoterRevokeMsg(newAccountName)
	test.SimulateOneBlock(lb, baseTime)
	test.SignCheckDeliver(t, lb, voterRevokeMsg, 2, false, newAccountTransactionPriv, baseTime)

	//validators can revoke voter candidancy after revoking validator candidancy
	valRevokeMsg := val.NewValidatorRevokeMsg(newAccountName)
	test.SignCheckDeliver(t, lb, valRevokeMsg, 3, true, newAccountTransactionPriv, baseTime)
	test.SimulateOneBlock(lb, baseTime)
	test.SignCheckDeliver(t, lb, voterRevokeMsg, 4, true, newAccountTransactionPriv, baseTime)

	// check delegator withdraw first coin return
	test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalHr*3600+1)
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(114285714))
	test.CheckBalance(t, delegator1Name, lb, types.NewCoin(300*types.Decimals))
	test.CheckBalance(t, delegator2Name, lb, types.NewCoin(100*types.Decimals))

	// check balance after freezing period
	for i := int64(1); i < test.CoinReturnTimes; i++ {
		test.SimulateOneBlock(lb, baseTime+test.CoinReturnIntervalHr*3600*(i+1)+1)
	}
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(5000*types.Decimals))
	test.CheckBalance(t, delegator1Name, lb, types.NewCoin(2100*types.Decimals))
	test.CheckBalance(t, delegator2Name, lb, types.NewCoin(700*types.Decimals))

}
