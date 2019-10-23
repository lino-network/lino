package bandwidth

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/lino-network/lino/test"
	linotypes "github.com/lino-network/lino/types"
	bandwidthmn "github.com/lino-network/lino/x/bandwidth/manager"
	bandwidthmodel "github.com/lino-network/lino/x/bandwidth/model"
	types "github.com/lino-network/lino/x/vote/types"
)

// test validator deposit
func TestMsgFee(t *testing.T) {

	newAccountSignPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	baseT := time.Unix(0, 0)
	baseTime := baseT.Unix()

	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)
	bandwidthmn.BandwidthManagerTestMode = false

	test.CreateAccountWithTime(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountSignPriv, "5000000000", baseTime)

	voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("3000000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountSignPriv, baseTime)

	test.CheckBalance(t, newAccountName, lb, linotypes.NewCoinFromInt64((5000000000-3000000-1)*linotypes.Decimals-2523))

	voteDepositSmallMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("1000"))
	fee := auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(linotypes.LinoCoinDenom, sdk.NewInt(100000000)))}
	smFee := auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(linotypes.LinoCoinDenom, sdk.NewInt(1)))}
	test.SignCheckDeliverWithFee(t, lb, voteDepositSmallMsg, 1, false, newAccountSignPriv, baseTime+1, smFee)
	test.SignCheckDeliverWithFee(t, lb, voteDepositSmallMsg, 2, true, newAccountSignPriv, baseTime+1, fee)

	curU := "0.501692631996395802"
	curUDec, _ := sdk.NewDecFromStr(curU)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 1,
		CurMsgFee:            linotypes.NewCoinFromInt64(2523),
		CurU:                 curUDec,
	}, lb)

	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 3, true, newAccountSignPriv, baseTime+4, 900)
	test.SimulateOneBlock(lb, baseTime+5)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            linotypes.NewCoinFromInt64(50006),
		CurU:                 curUDec,
	}, lb)

	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 903, true, newAccountSignPriv, baseTime+8, 900)
	test.SimulateOneBlock(lb, baseTime+9)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            linotypes.NewCoinFromInt64(565615),
		CurU:                 curUDec,
	}, lb)

	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 1803, true, newAccountSignPriv, baseTime+12, 900)
	test.SimulateOneBlock(lb, baseTime+13)
	test.CheckCurBlockInfo(t, bandwidthmodel.BlockInfo{
		TotalMsgSignedByApp:  0,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            linotypes.NewCoinFromInt64(4044452),
		CurU:                 curUDec,
	}, lb)

}
