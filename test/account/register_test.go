package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/test"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
)

// test normal transfer and register
func TestTransferAndRegisterAccount(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"

	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()

	registerMsg := acc.NewRegisterMsg(test.GenesisUser, newAccountName, types.LNO("100"),
		newAccountPriv.PubKey(), newAccountPriv.Generate(1).PubKey(), newAccountPriv.Generate(2).PubKey())
	test.SignCheckDeliver(t, lb, registerMsg, 0, true, test.GenesisTransactionPriv, baseTime)

	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(100*types.Decimals))
	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(100*types.Decimals)))
}

// register failed if register fee is insufficient
func TestRegisterAccountFailed(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"

	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()
	registerMsg := acc.NewRegisterMsg(test.GenesisUser, newAccountName, "0.1",
		newAccountPriv.PubKey(), newAccountPriv.Generate(1).PubKey(), newAccountPriv.Generate(2).PubKey())
	test.SignCheckDeliver(t, lb, registerMsg, 0, false, test.GenesisPriv, baseTime)

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	accManager := acc.NewAccountManager(lb.CapKeyAccountStore, ph)
	assert.False(t, accManager.IsAccountExist(ctx, types.AccountKey(newAccountName)))
	test.CheckBalance(t, test.GenesisUser, lb, test.GetGenesisAccountCoin(test.DefaultNumOfVal))
}
