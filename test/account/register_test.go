package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/test"
	acc "github.com/lino-network/lino/tx/account"
	reg "github.com/lino-network/lino/tx/register"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
)

// test normal transfer and register
func TestTransferAndRegisterAccount(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountAddr := newAccountPriv.PubKey().Address()
	newAccountName := "newUser"

	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()
	transferMsg := acc.NewTransferMsg(
		test.GenesisUser, types.LNO(sdk.NewRat(100)), []byte{}, acc.TransferToAddr(newAccountAddr))

	test.SignCheckDeliver(t, lb, transferMsg, 0, true, test.GenesisPriv, baseTime)

	registerMsg := reg.NewRegisterMsg(newAccountName, newAccountPriv.PubKey())
	test.SignCheckDeliver(t, lb, registerMsg, 0, true, newAccountPriv, baseTime)

	test.CheckBalance(t, newAccountName, lb, types.NewCoin(100*types.Decimals))
	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoin(100*types.Decimals)))
}

// register failed if account balance is zero
func TestRegisterAccountFailed(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"

	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()
	registerMsg := reg.NewRegisterMsg(newAccountName, newAccountPriv.PubKey())
	test.SignCheckDeliver(t, lb, registerMsg, 0, false, newAccountPriv, baseTime)

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	accManager := acc.NewAccountManager(lb.CapKeyAccountStore)
	assert.False(t, accManager.IsAccountExist(ctx, types.AccountKey(newAccountName)))
	test.CheckBalance(t, test.GenesisUser, lb, test.GetGenesisAccountCoin(test.DefaultNumOfVal))
}
