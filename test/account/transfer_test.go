package account

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"

	crypto "github.com/tendermint/go-crypto"
)

// test normal transfer to account name
func TestTransferToAccount(t *testing.T) {
	newAccountName := "newUser"
	baseTime := time.Now().Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), "100")

	transferMsg := acc.NewTransferMsg(
		test.GenesisUser, newAccountName, types.LNO("200"), "")

	test.SignCheckDeliver(t, lb, transferMsg, 1, true, test.GenesisTransactionPriv, baseTime)

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(300*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(300*types.Decimals))
}
