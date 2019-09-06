package bandwidth

import (
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	// bandwidthmodel "github.com/lino-network/lino/x/bandwidth/model"
	devtypes "github.com/lino-network/lino/x/developer/types"
	vote "github.com/lino-network/lino/x/vote"
)

// test validator deposit
func TestAppBandwidth(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	baseTime := time.Now().Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, newAccountAppPriv, "5000000000")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("3000000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	registerAppMsg := devtypes.NewDeveloperRegisterMsg(newAccountName, "dummy", "dummy", "dummy")
	test.SignCheckDeliver(t, lb, registerAppMsg, 1, true, newAccountTransactionPriv, baseTime)

	// the tx will fail since app bandwidth info will be updated hourly
	voteDepositSmallMsg := vote.NewStakeInMsg(newAccountName, types.LNO("1000"))
	test.SignCheckDeliver(t, lb, voteDepositSmallMsg, 2, false, newAccountTransactionPriv, baseTime)

	// the tx will success after one hour
	test.SimulateOneBlock(lb, baseTime+3600)
	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 2, true, newAccountTransactionPriv, baseTime+3603, 4800)
	// new bandwidth credit will be -29049
	// test.CheckAppBandwidthInfo(t, bandwidthmodel.AppBandwidthInfo{}, types.AccountKey(newAccountName), lb)
	// can send msg after max 3600 seconds
	test.RepeatSignCheckDeliver(t, lb, voteDepositSmallMsg, 4802, true, newAccountTransactionPriv, baseTime+3603+3600, 1)
	// test.CheckAppBandwidthInfo(t, bandwidthmodel.AppBandwidthInfo{}, types.AccountKey(newAccountName), lb)

}
