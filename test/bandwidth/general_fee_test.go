package bandwidth

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	bandwidthmn "github.com/lino-network/lino/x/bandwidth/manager"
	// bandwidthmodel "github.com/lino-network/lino/x/bandwidth/model"
	vote "github.com/lino-network/lino/x/vote"
)

// test validator deposit
func TestMsgFee(t *testing.T) {

	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	baseTime := time.Now().Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	bandwidthmn.BandwidthManagerTestMode = false
	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, newAccountAppPriv, "5000000000")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("3000000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	voteDepositSmallMsg := vote.NewStakeInMsg(newAccountName, types.LNO("1000"))
	fee := auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(100000000)))}
	test.SignCheckDeliverWithFee(t, lb, voteDepositSmallMsg, 1, true, newAccountTransactionPriv, baseTime, fee)
	// test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{}, lb)
	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 2, true, newAccountTransactionPriv, baseTime+1, 1000)
	test.SignCheckDeliverWithFee(t, lb, voteDepositSmallMsg, 2+1000, true, newAccountTransactionPriv, baseTime+2, fee)
	// 0.50006
	// test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{}, lb)
	// test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 3+1000, true, newAccountTransactionPriv, baseTime+3, 1000)
	// test.SignCheckDeliverWithFee(t, lb, voteDepositSmallMsg, 3+2000, true, newAccountTransactionPriv, baseTime+4, fee)
	// test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{}, lb)
}
