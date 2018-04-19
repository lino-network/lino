package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	acc "github.com/lino-network/lino/tx/account"
	post "github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// test donate to a normal post
func TestNormalDonation(t *testing.T) {
	newPostUserPriv := crypto.GenPrivKeyEd25519()
	newPostUser := "poster"
	postID := "New Post"

	newDonateUserPriv := crypto.GenPrivKeyEd25519()
	newDonateUser := "donator"
	// recover some stake
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newPostUser, lb, 0, newPostUserPriv, 100)
	test.CreateAccount(t, newDonateUser, lb, 1, newDonateUserPriv, 100)

	test.CreateTestPost(
		t, lb, newPostUser, postID, 0, newPostUserPriv, "", "", "", "", sdk.ZeroRat, baseTime)

	test.CheckBalance(t, newPostUser, lb, types.NewCoin(100*types.Decimals))
	test.CheckBalance(t, newDonateUser, lb, types.NewCoin(100*types.Decimals))

	donateMsg := post.NewDonateMsg(
		types.AccountKey(newDonateUser), types.LNO(sdk.NewRat(50)),
		types.AccountKey(newPostUser), postID, "")

	test.SignCheckDeliver(t, lb, donateMsg, 0, true, newDonateUserPriv, baseTime)

	test.CheckBalance(t, newDonateUser, lb, types.NewCoin(50*types.Decimals))
	test.CheckBalance(t, newPostUser, lb, types.NewCoin(10000000+4750000))

	claimMsg := acc.NewClaimMsg(newPostUser)
	test.SignCheckDeliver(t, lb, claimMsg, 1, true, newPostUserPriv, baseTime)
	test.CheckBalance(t, newPostUser, lb, types.NewCoin(10000000+4750000))
	test.SignCheckDeliver(
		t, lb, claimMsg, 2, true, newPostUserPriv, baseTime+test.ConsumptionFreezingPeriodHr*3600+1)
	test.CheckBalance(t, newPostUser, lb, types.NewCoin(10000000+4750000))
	test.SignCheckDeliver(
		t, lb, claimMsg, 3, true, newPostUserPriv, baseTime+test.ConsumptionFreezingPeriodHr*3600+2)
	test.CheckBalance(t, newPostUser, lb, types.NewCoin(944687598610))
}
