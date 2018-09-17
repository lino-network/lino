package global

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

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

	// 3rd day
	baseTime += 3600 * 24 * 3
	test.SimulateOneBlock(lb, baseTime)

	u1ClaimInterestMsg := vote.NewClaimInterestMsg(u1Name)
	u2ClaimInterestMsg := vote.NewClaimInterestMsg(u2Name)

	test.SignCheckDeliver(t, lb, u1ClaimInterestMsg, 1, true, u1Priv, baseTime)
	test.SignCheckDeliver(t, lb, u2ClaimInterestMsg, 1, true, u2Priv, baseTime)

	test.CheckBalance(t, u1Name, lb, types.NewCoinFromInt64((89999+0.15748)*types.Decimals))
	test.CheckBalance(t, u2Name, lb, types.NewCoinFromInt64((59999+0.62992)*types.Decimals))

	u1StakeInMsg = vote.NewStakeInMsg(u1Name, types.LNO("20000"))
	u2StakeOutMsg := vote.NewStakeOutMsg(u2Name, types.LNO("10000"))

	test.SignCheckDeliver(t, lb, u1StakeInMsg, 2, true, u1Priv, baseTime)
	test.SignCheckDeliver(t, lb, u2StakeOutMsg, 2, true, u2Priv, baseTime)
	test.SignCheckDeliver(t, lb, donateMsg, 1, true, donatorPriv, baseTime)

	// 4th day
	baseTime += 3600 * 24 * 1
	test.SimulateOneBlock(lb, baseTime)
	test.SignCheckDeliver(t, lb, donateMsg, 2, true, donatorPriv, baseTime)

	// 5th day
	baseTime += 3600 * 24 * 1
	test.SimulateOneBlock(lb, baseTime)
	test.SignCheckDeliver(t, lb, donateMsg, 3, true, donatorPriv, baseTime)

	u1StakeOutMsg := vote.NewStakeOutMsg(u1Name, types.LNO("30000"))
	u2StakeOutMsg = vote.NewStakeOutMsg(u2Name, types.LNO("30000"))

	test.SignCheckDeliver(t, lb, u1StakeOutMsg, 3, true, u1Priv, baseTime)
	test.SignCheckDeliver(t, lb, u2StakeOutMsg, 3, true, u2Priv, baseTime)

	// 6th day
	baseTime += 3600 * 24 * 1
	test.SimulateOneBlock(lb, baseTime)

	test.SignCheckDeliver(t, lb, u1ClaimInterestMsg, 4, true, u1Priv, baseTime)
	test.SignCheckDeliver(t, lb, u2ClaimInterestMsg, 4, true, u2Priv, baseTime)

	test.CheckBalance(t, u1Name, lb, types.NewCoinFromInt64((69999+1.10088)*types.Decimals))
	test.CheckBalance(t, u2Name, lb, types.NewCoinFromInt64((59999+1.57332)*types.Decimals))

}
