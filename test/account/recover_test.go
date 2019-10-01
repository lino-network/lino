package account

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	accmodel "github.com/lino-network/lino/x/account/model"
	acctypes "github.com/lino-network/lino/x/account/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestRecoverAccount(t *testing.T) {
	transactionPriv := secp256k1.GenPrivKey()
	signingPriv := secp256k1.GenPrivKey()
	newTransactionPriv := secp256k1.GenPrivKey()
	newSigningPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"

	baseT := time.Now()
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	registerMsgV2 := acctypes.NewRegisterV2Msg(
		types.NewAccOrAddrFromAcc(
			types.AccountKey(test.GenesisUser)), newAccountName, types.LNO("100"),
		transactionPriv.PubKey(), signingPriv.PubKey())
	test.SignCheckDeliverWithMultiSig(
		t, lb, registerMsgV2, []uint64{0, 0}, true,
		[]secp256k1.PrivKeySecp256k1{test.GenesisTransactionPriv, transactionPriv}, baseTime)

	recoverMsg := acctypes.NewRecoverMsg(newAccountName, newTransactionPriv.PubKey(), newSigningPriv.PubKey())
	test.SignCheckDeliverWithMultiSig(
		t, lb, recoverMsg, []uint64{1, 0}, true,
		[]secp256k1.PrivKeySecp256k1{transactionPriv, newTransactionPriv}, baseTime)
	test.CheckAccountInfo(t, newAccountName, lb, accmodel.AccountInfo{
		Username:       types.AccountKey(newAccountName),
		TransactionKey: newTransactionPriv.PubKey(),
		SigningKey:     newSigningPriv.PubKey(),
		CreatedAt:      baseTime,
		Address:        sdk.AccAddress(newTransactionPriv.PubKey().Address()),
	})
}
