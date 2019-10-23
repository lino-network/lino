package bandwidth

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/lino-network/lino/test"
	linotypes "github.com/lino-network/lino/types"
	bandwidthmodel "github.com/lino-network/lino/x/bandwidth/model"
	devtypes "github.com/lino-network/lino/x/developer/types"
	types "github.com/lino-network/lino/x/vote/types"
)

// test validator deposit
func TestAppBandwidth(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	baseT := time.Unix(0, 0)
	baseTime := time.Unix(0, 0).Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, "5000000000")

	voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("3000000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	registerAppMsg := devtypes.NewDeveloperRegisterMsg(newAccountName, "dummy", "dummy", "dummy")
	test.SignCheckDeliver(t, lb, registerAppMsg, 2, true, newAccountTransactionPriv, baseTime)

	// the tx will fail since app bandwidth info will be updated hourly
	voteDepositSmallMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("1000"))
	test.SignCheckTxFail(t, lb, voteDepositSmallMsg, 2, newAccountTransactionPriv)

	// the tx will success after one hour
	test.SimulateOneBlock(lb, baseTime+3600)
	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 3, true, newAccountTransactionPriv, baseTime+3603, 1000)
	curCredit := "910.849417274954801802"
	curCreditDec, _ := sdk.NewDecFromStr(curCredit)
	test.CheckAppBandwidthInfo(t, bandwidthmodel.AppBandwidthInfo{
		Username:           "newuser",
		MaxBandwidthCredit: sdk.NewDec(2400),
		CurBandwidthCredit: curCreditDec,
		MessagesInCurBlock: 0,
		ExpectedMPS:        sdk.NewDec(240),
		LastRefilledAt:     baseTime + 3603,
	}, linotypes.AccountKey(newAccountName), lb)

	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 1003, true, newAccountTransactionPriv, baseTime+3603+3, 1000)
	curCredit = "7.099448771136388802"
	curCreditDec, _ = sdk.NewDecFromStr(curCredit)
	test.CheckAppBandwidthInfo(t, bandwidthmodel.AppBandwidthInfo{
		Username:           "newuser",
		MaxBandwidthCredit: sdk.NewDec(2400),
		CurBandwidthCredit: curCreditDec,
		MessagesInCurBlock: 0,
		ExpectedMPS:        sdk.NewDec(240),
		LastRefilledAt:     baseTime + 3603 + 3,
	}, linotypes.AccountKey(newAccountName), lb)

	curCredit = "301.600990083371031122"
	curCreditDec, _ = sdk.NewDecFromStr(curCredit)
	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 2003, true, newAccountTransactionPriv, baseTime+3603+3+3, 720)
	test.CheckAppBandwidthInfo(t, bandwidthmodel.AppBandwidthInfo{
		Username:           "newuser",
		MaxBandwidthCredit: sdk.NewDec(2400),
		CurBandwidthCredit: curCreditDec,
		MessagesInCurBlock: 0,
		ExpectedMPS:        sdk.NewDec(240),
		LastRefilledAt:     baseTime + 3603 + 3 + 3,
	}, linotypes.AccountKey(newAccountName), lb)
}
