package validator

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	val "github.com/lino-network/lino/tx/validator"
	vote "github.com/lino-network/lino/tx/vote"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// test normal revoke
func TestValidatorRevoke(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	newValidatorPriv := crypto.GenPrivKeyEd25519()

	baseTime := time.Now().Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0, newAccountPriv, 5000)

	voteDepositMsg := vote.NewVoterDepositMsg(newAccountName, types.LNO(sdk.NewRat(3000)))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO(sdk.NewRat(1500)), newValidatorPriv.PubKey())
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountPriv, baseTime)
	test.CheckAllValidatorList(t, newAccountName, true, lb)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)

	valRevokeMsg := val.NewValidatorRevokeMsg(newAccountName)
	test.SignCheckDeliver(t, lb, valRevokeMsg, 2, true, newAccountPriv, baseTime)
	test.CheckAllValidatorList(t, newAccountName, false, lb)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(500*types.Decimals))

	test.SignCheckDeliver(t, lb, valRevokeMsg, 3, false, newAccountPriv,
		baseTime+test.CoinReturnIntervalHr*3600+1)
	// check the first coin return
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(71428571))
	for i := int64(1); i < types.CoinReturnTimes; i++ {
		test.SignCheckDeliver(t, lb, valRevokeMsg, 3+i, false, newAccountPriv,
			baseTime+test.CoinReturnIntervalHr*3600*(i+1)+2)
	}
	// will get all coins back after the freezing period
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(2000*types.Decimals))

	// won't get extra coins in the future
	test.SignCheckDeliver(t, lb, valRevokeMsg, 3+test.CoinReturnTimes, false, newAccountPriv,
		baseTime+test.CoinReturnIntervalHr*3600*(test.CoinReturnTimes+1)+3)
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(2000*types.Decimals))

}
