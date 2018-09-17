package global

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	acc "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/post"
	vote "github.com/lino-network/lino/x/vote"
)

func TestStakeInterest(t *testing.T) {
	postUserPriv := secp256k1.GenPrivKey()
	donatorPriv := secp256k1.GenPrivKey()
	u1Priv := secp256k1.GenPrivKey()
	u2Priv := secp256k1.GenPrivKey()

	postUserName := "poster"
	donatorName := "donator"
	u1Name := "user1"
	u2Name := "user2"

	postID := "New Post"

	// to recover the coin day
	baseTime := time.Now().Unix() + 7200
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, donatorName, lb, 0,
		secp256k1.GenPrivKey(), donatorPriv, secp256k1.GenPrivKey(), "100000")
	test.CreateAccount(t, u1Name, lb, 1,
		secp256k1.GenPrivKey(), u1Priv, secp256k1.GenPrivKey(), "100000")
	test.CreateAccount(t, u2Name, lb, 2,
		secp256k1.GenPrivKey(), u2Priv, secp256k1.GenPrivKey(), "100000")
	test.CreateAccount(t, postUserName, lb, 3,
		secp256k1.GenPrivKey(), postUserPriv, secp256k1.GenPrivKey(), "100000")
	test.CreateTestPost(
		t, lb, postUserName, postID, 0, postUserPriv, "", "", "", "", "0", baseTime)

	donateMsg := post.NewDonateMsg(
		donatorName, types.LNO("2000"), postUserName, postID, "", "")
	u1StakeInMsg := vote.NewStakeInMsg(u1Name, types.LNO("10000"))
	u2StakeInMsg := vote.NewStakeInMsg(u2Name, types.LNO("40000"))

	test.SignCheckDeliver(t, lb, donateMsg, 0, true, donatorPriv, baseTime)
	test.SignCheckDeliver(t, lb, u1StakeInMsg, 0, true, u1Priv, baseTime)
	test.SignCheckDeliver(t, lb, u2StakeInMsg, 0, true, u2Priv, baseTime)

	test.CheckBalance(t, u1Name, lb, types.NewCoinFromInt64(89999*types.Decimals))
	test.CheckBalance(t, u2Name, lb, types.NewCoinFromInt64(59999*types.Decimals))

	// 3 days
	baseTime += 3600 * 24 * 3
	test.SimulateOneBlock(lb, baseTime)

	u1ClaimInterestMsg := acc.NewClaimInterestMsg(u1Name)
	u2ClaimInterestMsg := acc.NewClaimInterestMsg(u2Name)

	test.SignCheckDeliver(t, lb, u1ClaimInterestMsg, 1, true, u1Priv, baseTime)
	test.SignCheckDeliver(t, lb, u2ClaimInterestMsg, 1, true, u2Priv, baseTime)

	test.CheckBalance(t, u1Name, lb, types.NewCoinFromInt64((89999+40)*types.Decimals))
	test.CheckBalance(t, u2Name, lb, types.NewCoinFromInt64((59999+160)*types.Decimals))

}
