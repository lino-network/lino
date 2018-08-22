package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	post "github.com/lino-network/lino/x/post"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// test donate to a normal post
func TestNormalDonation(t *testing.T) {
	newPostUserTransactionPriv := secp256k1.GenPrivKey()
	newPostUserAppPriv := secp256k1.GenPrivKey()
	newPostUser := "poster"
	postID := "New Post"

	newDonateUserTransactionPriv := secp256k1.GenPrivKey()
	newDonateUser := "donator"
	// recover some stake
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newPostUser, lb, 0,
		secp256k1.GenPrivKey(), newPostUserTransactionPriv, newPostUserAppPriv, "100")
	test.CreateAccount(t, newDonateUser, lb, 1,
		secp256k1.GenPrivKey(), newDonateUserTransactionPriv, secp256k1.GenPrivKey(), "100")

	test.CreateTestPost(
		t, lb, newPostUser, postID, 0, newPostUserAppPriv, "", "", "", "", "0", baseTime)

	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(99*types.Decimals))
	test.CheckBalance(t, newDonateUser, lb, types.NewCoinFromInt64(99*types.Decimals))

	donateMsg := post.NewDonateMsg(
		newDonateUser, types.LNO("50"), newPostUser, postID, "", "")

	test.SignCheckDeliver(t, lb, donateMsg, 0, true, newDonateUserTransactionPriv, baseTime)

	test.CheckBalance(t, newDonateUser, lb, types.NewCoinFromInt64(49*types.Decimals))
	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(9900000+4750000))

	claimMsg := acc.NewClaimMsg(newPostUser)
	test.SignCheckDeliver(t, lb, claimMsg, 1, true, newPostUserTransactionPriv, baseTime)
	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(9900000+4750000))
	test.SignCheckDeliver(
		t, lb, claimMsg, 2, true, newPostUserTransactionPriv, baseTime+test.ConsumptionFreezingPeriodSec+1)
	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(1228089278362))
}
