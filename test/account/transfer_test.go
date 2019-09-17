package account

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	acctypes "github.com/lino-network/lino/x/account/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// test normal transfer to account name
func TestTransferToAccount(t *testing.T) {
	newAccountName := "newuser"

	baseT := time.Now()
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), secp256k1.GenPrivKey(), secp256k1.GenPrivKey(), "100")

	transferMsg := acctypes.NewTransferMsg(
		test.GenesisUser, newAccountName, types.LNO("200"), "")

	test.SignCheckDeliver(t, lb, transferMsg, 1, true, test.GenesisTransactionPriv, baseTime)

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(300*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(299*types.Decimals))
}
