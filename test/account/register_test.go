package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	accmn "github.com/lino-network/lino/x/account/manager"
	acctypes "github.com/lino-network/lino/x/account/types"
	"github.com/lino-network/lino/x/global"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// test normal transfer and register
func TestTransferAndRegisterAccount(t *testing.T) {
	newTransactionPriv := secp256k1.GenPrivKey()
	newSigningPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"

	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()

	registerMsg := acctypes.NewRegisterMsg(test.GenesisUser, newAccountName, types.LNO("100"),
		newTransactionPriv.PubKey(), newSigningPriv.PubKey(), nil)
	test.SignCheckDeliver(t, lb, registerMsg, 0, true, test.GenesisTransactionPriv, baseTime)

	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(99*types.Decimals))
	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(100*types.Decimals)))
}

// register failed if register fee is insufficient
func TestRegisterAccountFailed(t *testing.T) {
	newResetPriv := secp256k1.GenPrivKey()
	newTransactionPriv := secp256k1.GenPrivKey()
	newAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"

	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix()
	registerMsg := acctypes.NewRegisterMsg(test.GenesisUser, newAccountName, "0.1",
		newResetPriv.PubKey(), newTransactionPriv.PubKey(), newAppPriv.PubKey())
	test.SignCheckDeliver(t, lb, registerMsg, 0, false, test.GenesisPriv, baseTime)

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	gm := global.NewGlobalManager(lb.CapKeyGlobalStore, ph)
	accManager := accmn.NewAccountManager(lb.CapKeyAccountStore, ph, &gm)
	assert.False(t, accManager.DoesAccountExist(ctx, types.AccountKey(newAccountName)))
	test.CheckBalance(t, test.GenesisUser, lb, test.GetGenesisAccountCoin(test.DefaultNumOfVal))
}
