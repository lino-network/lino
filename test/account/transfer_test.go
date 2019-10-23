package account

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	acctypes "github.com/lino-network/lino/x/account/types"
)

// test normal transfer to account name
func TestTransferToAccount(t *testing.T) {
	newAccountName := "newuser"

	baseT := time.Unix(0, 0)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), secp256k1.GenPrivKey(), "100")

	transferMsg := acctypes.NewTransferMsg(
		test.GenesisUser, newAccountName, types.LNO("200"), "")

	test.SignCheckDeliver(t, lb, transferMsg, 1, true, test.GenesisTransactionPriv, baseTime)

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(300*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(299*types.Decimals))
}

// test transfer between addresses.
func TestTransferToAddress(t *testing.T) {
	newAccountName := "newuser"
	newReceiver := "newreceiver"

	baseT := time.Unix(0, 0)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	senderPriv := secp256k1.GenPrivKey()
	receiverPriv := secp256k1.GenPrivKey()

	test.CreateAccount(t, newAccountName, lb, 0,
		senderPriv, secp256k1.GenPrivKey(), "100")
	test.CreateAccount(t, newReceiver, lb, 1,
		receiverPriv, secp256k1.GenPrivKey(), "100")

	// user -> address
	transferMsg := acctypes.NewTransferV2Msg(
		types.NewAccOrAddrFromAcc(types.AccountKey(test.GenesisUser)),
		types.NewAccOrAddrFromAddr(sdk.AccAddress(senderPriv.PubKey().Address())),
		types.LNO("200"), "")
	test.SignCheckDeliver(t, lb, transferMsg, 2, true, test.GenesisTransactionPriv, baseTime)

	test.CheckBalance(t, test.GenesisUser, lb,
		test.GetGenesisAccountCoin(test.DefaultNumOfVal).Minus(types.NewCoinFromInt64(400*types.Decimals)))
	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(299*types.Decimals))

	// addr -> addr
	transferMsg = acctypes.NewTransferV2Msg(
		types.NewAccOrAddrFromAddr(sdk.AccAddress(senderPriv.PubKey().Address())),
		types.NewAccOrAddrFromAddr(sdk.AccAddress(receiverPriv.PubKey().Address())),
		types.LNO("100"), "")
	test.SignCheckDeliver(t, lb, transferMsg, 1, true, senderPriv, baseTime)

	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(199*types.Decimals))
	test.CheckBalance(t, newReceiver, lb, types.NewCoinFromInt64(199*types.Decimals))

	// addr -> user
	transferMsg = acctypes.NewTransferV2Msg(
		types.NewAccOrAddrFromAddr(sdk.AccAddress(senderPriv.PubKey().Address())),
		types.NewAccOrAddrFromAcc(types.AccountKey(newReceiver)),
		types.LNO("100"), "")
	test.SignCheckDeliver(t, lb, transferMsg, 2, true, senderPriv, baseTime)

	test.CheckBalance(t, newAccountName, lb, types.NewCoinFromInt64(99*types.Decimals))
	test.CheckBalance(t, newReceiver, lb, types.NewCoinFromInt64(299*types.Decimals))
}
