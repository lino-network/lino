package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	accmn "github.com/lino-network/lino/x/account/manager"
	accmodel "github.com/lino-network/lino/x/account/model"
	acctypes "github.com/lino-network/lino/x/account/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestRegisterAccountV2(t *testing.T) {
	newTransactionPriv := secp256k1.GenPrivKey()
	newSigningPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"

	baseT := time.Unix(0, 0)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	registerMsgV2 := acctypes.NewRegisterV2Msg(
		types.NewAccOrAddrFromAcc(types.AccountKey(test.GenesisUser)), newAccountName, types.LNO("100"),
		newTransactionPriv.PubKey(), newSigningPriv.PubKey())
	test.SignCheckDeliverWithMultiSig(
		t, lb, registerMsgV2, []uint64{0, 0}, true,
		[]secp256k1.PrivKeySecp256k1{test.GenesisTransactionPriv, newTransactionPriv}, baseTime)

	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(99*types.Decimals))
	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(100*types.Decimals)))
	test.CheckAccountInfo(t, newAccountName, lb, accmodel.AccountInfo{
		Username:       types.AccountKey(newAccountName),
		TransactionKey: newTransactionPriv.PubKey(),
		SigningKey:     newSigningPriv.PubKey(),
		CreatedAt:      baseTime,
		Address:        sdk.AccAddress(newTransactionPriv.PubKey().Address()),
	})
}

func TestRegisterAccountV2Failed(t *testing.T) {
	newTransactionPriv := secp256k1.GenPrivKey()
	newSigningPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"

	baseT := time.Unix(0, 0)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	registerMsgV2 := acctypes.NewRegisterV2Msg(
		types.NewAccOrAddrFromAcc(types.AccountKey(test.GenesisUser)), newAccountName,
		types.LNO("0.1"),
		newTransactionPriv.PubKey(), newSigningPriv.PubKey())
	test.SignCheckDeliverWithMultiSig(
		t, lb, registerMsgV2, []uint64{0, 0}, false,
		[]secp256k1.PrivKeySecp256k1{test.GenesisTransactionPriv, newTransactionPriv}, baseTime)
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	accManager := accmn.NewAccountManager(lb.CapKeyAccountStore, ph)
	assert.False(t, accManager.DoesAccountExist(ctx, types.AccountKey(newAccountName)))
	test.CheckBalance(t, test.GenesisUser, lb, test.GetGenesisAccountCoin(test.DefaultNumOfVal))
}
