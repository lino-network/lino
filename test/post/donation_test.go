package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	post "github.com/lino-network/lino/x/post"

	crypto "github.com/tendermint/go-crypto"
)

// test donate to a normal post
func TestNormalDonation(t *testing.T) {
	newPostUserTransactionPriv := crypto.GenPrivKeyEd25519()
	newPostUserPostPriv := crypto.GenPrivKeyEd25519()
	newPostUser := "poster"
	postID := "New Post"

	newDonateUserTransactionPriv := crypto.GenPrivKeyEd25519()
	newDonateUser := "donator"
	// recover some stake
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newPostUser, lb, 0,
		crypto.GenPrivKeyEd25519(), newPostUserTransactionPriv, crypto.GenPrivKeyEd25519(), newPostUserPostPriv, "100")
	test.CreateAccount(t, newDonateUser, lb, 1,
		crypto.GenPrivKeyEd25519(), newDonateUserTransactionPriv, crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), "100")

	test.CreateTestPost(
		t, lb, newPostUser, postID, 0, newPostUserPostPriv, "", "", "", "", "0", baseTime)

	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(100*types.Decimals))
	test.CheckBalance(t, newDonateUser, lb, types.NewCoinFromInt64(100*types.Decimals))

	donateMsg := post.NewDonateMsg(
		newDonateUser, types.LNO("50"), newPostUser, postID, "", "")

	test.SignCheckDeliver(t, lb, donateMsg, 0, true, newDonateUserTransactionPriv, baseTime)

	test.CheckBalance(t, newDonateUser, lb, types.NewCoinFromInt64(50*types.Decimals))
	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(10000000+4750000))

	claimMsg := acc.NewClaimMsg(newPostUser)
	test.SignCheckDeliver(t, lb, claimMsg, 1, true, newPostUserTransactionPriv, baseTime)
	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(10000000+4750000))
	test.SignCheckDeliver(
		t, lb, claimMsg, 2, true, newPostUserTransactionPriv, baseTime+test.ConsumptionFreezingPeriodHr*3600+1)
	test.CheckBalance(t, newPostUser, lb, types.NewCoinFromInt64(1228089378362))
}
