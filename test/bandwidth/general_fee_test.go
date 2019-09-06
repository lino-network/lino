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
	bandwidthmodel "github.com/lino-network/lino/x/bandwidth/model"
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
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64((5000000000-3000000-1)*types.Decimals-2523))

	voteDepositSmallMsg := vote.NewStakeInMsg(newAccountName, types.LNO("1000"))
	fee := auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(100000000)))}
	smFee := auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(1)))}
	test.SignCheckDeliverWithFee(t, lb, voteDepositSmallMsg, 1, false, newAccountTransactionPriv, baseTime+1, smFee)
	test.SignCheckDeliverWithFee(t, lb, voteDepositSmallMsg, 1, true, newAccountTransactionPriv, baseTime+1, fee)

	curU := "0.501692631996395802"
	curUDec, _ := sdk.NewDecFromStr(curU)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 1,
		CurMsgFee:            types.NewCoinFromInt64(2523),
		CurU:                 curUDec,
	}, lb)

	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 2, true, newAccountTransactionPriv, baseTime+4, 900)
	test.SimulateOneBlock(lb, baseTime+5)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            types.NewCoinFromInt64(50006),
		CurU:                 curUDec,
	}, lb)

	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 902, true, newAccountTransactionPriv, baseTime+8, 900)
	test.SimulateOneBlock(lb, baseTime+9)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            types.NewCoinFromInt64(565615),
		CurU:                 curUDec,
	}, lb)

	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 1802, true, newAccountTransactionPriv, baseTime+12, 900)
	test.SimulateOneBlock(lb, baseTime+13)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            types.NewCoinFromInt64(4044452),
		CurU:                 curUDec,
	}, lb)

}
