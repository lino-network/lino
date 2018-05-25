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
		test.GenesisUser, types.LNO("100"), "", acc.TransferToUser(newAccountName))

	test.SignCheckDeliver(t, lb, transferMsg, 1, true, test.GenesisTransactionPriv, baseTime)

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(200*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(200*types.Decimals))
}

// test normal transfer to address
func TestTransferToAddress(t *testing.T) {
	newAccountName := "newUser"
	newAccountPriv := crypto.GenPrivKeyEd25519()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()

	transferMsg := acc.NewTransferMsg(
		test.GenesisUser, types.LNO("100"), "",
		acc.TransferToAddr(newAccountPriv.PubKey().Address()))
	test.SignCheckDeliver(t, lb, transferMsg, 0, true, test.GenesisTransactionPriv, baseTime)

	test.CreateAccount(t, newAccountName, lb, 1,
		newAccountPriv, crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), "100")

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(200*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(200*types.Decimals))
}
