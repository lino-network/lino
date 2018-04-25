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
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	baseTime := time.Now().Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0, newAccountPriv, "100")

	transferMsg := acc.NewTransferMsg(
		test.GenesisUser, types.LNO("100"), []byte{}, acc.TransferToUser(newAccountName))

	test.SignCheckDeliver(t, lb, transferMsg, 1, true, test.GenesisPriv, baseTime)

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoin(200*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(200*types.Decimals))
}

// test normal transfer to address
func TestTransferToAddress(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()

	transferMsg := acc.NewTransferMsg(
		test.GenesisUser, types.LNO("100"), []byte{},
		acc.TransferToAddr(newAccountPriv.PubKey().Address()))
	test.SignCheckDeliver(t, lb, transferMsg, 0, true, test.GenesisPriv, baseTime)

	test.CreateAccount(t, newAccountName, lb, 1, newAccountPriv, "100")

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoin(200*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoin(200*types.Decimals))
}
