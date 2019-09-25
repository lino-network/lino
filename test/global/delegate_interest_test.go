package global

// import (
// 	"testing"
// 	"time"

// 	"github.com/lino-network/lino/test"
// 	"github.com/lino-network/lino/types"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"

// 	posttypes "github.com/lino-network/lino/x/post/types"
// 	vote "github.com/lino-network/lino/x/vote"
// )

// func TestDelegateInterest(t *testing.T) {
// 	postUserPriv := secp256k1.GenPrivKey()
// 	donatorPriv := secp256k1.GenPrivKey()
// 	u1Priv := secp256k1.GenPrivKey()
// 	u2Priv := secp256k1.GenPrivKey()

// 	postUserName := "poster"
// 	donatorName := "donator"
// 	u1Name := "user1"
// 	u2Name := "user2"

// 	postID := "New Post"

// 	// to recover the coin day
// 	baseT := time.Now().Add(7200 * time.Second)
// 	baseTime := baseT.Unix()
// 	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

// 	test.CreateAccount(t, donatorName, lb, 0,
// 		secp256k1.GenPrivKey(), donatorPriv, secp256k1.GenPrivKey(), "100000")
// 	test.CreateAccount(t, u1Name, lb, 1,
// 		secp256k1.GenPrivKey(), u1Priv, secp256k1.GenPrivKey(), "100000")
// 	test.CreateAccount(t, u2Name, lb, 2,
// 		secp256k1.GenPrivKey(), u2Priv, secp256k1.GenPrivKey(), "100000")
// 	test.CreateAccount(t, postUserName, lb, 3,
// 		secp256k1.GenPrivKey(), postUserPriv, secp256k1.GenPrivKey(), "100000")
// 	test.CreateTestPost(
// 		t, lb, postUserName, postID, 0, postUserPriv, baseTime)

// 	donateMsg := posttypes.NewDonateMsg(
// 		donatorName, types.LNO("2000"), postUserName, postID, "", "")
// 	u1DelegateMsg := vote.NewDelegateMsg(u1Name, u2Name, types.LNO("10000"))
// 	u2DelegateMsg := vote.NewDelegateMsg(u2Name, u1Name, types.LNO("40000"))

// 	test.SignCheckDeliver(t, lb, donateMsg, 0, true, donatorPriv, baseTime)
// 	test.SignCheckDeliver(t, lb, u1DelegateMsg, 0, true, u1Priv, baseTime)
// 	test.SignCheckDeliver(t, lb, u2DelegateMsg, 0, true, u2Priv, baseTime)

// 	test.CheckBalance(t, u1Name, lb, types.NewCoinFromInt64(89999*types.Decimals))
// 	test.CheckBalance(t, u2Name, lb, types.NewCoinFromInt64(59999*types.Decimals))

// 	// 3rd day
// 	baseTime += 3600 * 24 * 3
// 	test.SimulateOneBlock(lb, baseTime)

// 	u1ClaimInterestMsg := vote.NewClaimInterestMsg(u1Name)
// 	u2ClaimInterestMsg := vote.NewClaimInterestMsg(u2Name)

// 	test.SignCheckDeliver(t, lb, u1ClaimInterestMsg, 1, true, u1Priv, baseTime)
// 	test.SignCheckDeliver(t, lb, u2ClaimInterestMsg, 1, true, u2Priv, baseTime)

// 	test.CheckBalance(t, u1Name, lb, types.NewCoinFromInt64((89999+0.15748)*types.Decimals))
// 	test.CheckBalance(t, u2Name, lb, types.NewCoinFromInt64((59999+0.62992)*types.Decimals))

// 	u1DelegateMsg = vote.NewDelegateMsg(u1Name, u2Name, types.LNO("20000"))
// 	u2StakeOutMsg := vote.NewStakeOutMsg(u2Name, types.LNO("10000"))
// 	u2DelegatorWithdrawMsg := vote.NewDelegatorWithdrawMsg(u2Name, u1Name, types.LNO("10000"))

// 	test.SignCheckDeliver(t, lb, u1DelegateMsg, 2, true, u1Priv, baseTime)
// 	test.SignCheckDeliver(t, lb, u2StakeOutMsg, 2, false, u2Priv, baseTime)
// 	test.SignCheckDeliver(t, lb, u2DelegatorWithdrawMsg, 3, true, u2Priv, baseTime)
// 	test.SignCheckDeliver(t, lb, donateMsg, 1, true, donatorPriv, baseTime)

// 	// 4th day
// 	baseTime += 3600 * 24 * 1
// 	test.SimulateOneBlock(lb, baseTime)
// 	test.SignCheckDeliver(t, lb, donateMsg, 2, true, donatorPriv, baseTime)

// 	// 5th day
// 	baseTime += 3600 * 24 * 1
// 	test.SimulateOneBlock(lb, baseTime)
// 	test.SignCheckDeliver(t, lb, donateMsg, 3, true, donatorPriv, baseTime)

// 	u1DelegatorWithdrawMsg := vote.NewDelegatorWithdrawMsg(u1Name, u2Name, types.LNO("30000"))
// 	u2DelegatorWithdrawMsg = vote.NewDelegatorWithdrawMsg(u2Name, u1Name, types.LNO("30000"))

// 	test.SignCheckDeliver(t, lb, u1DelegatorWithdrawMsg, 3, true, u1Priv, baseTime)
// 	test.SignCheckDeliver(t, lb, u2DelegatorWithdrawMsg, 4, true, u2Priv, baseTime)

// 	// 6th day
// 	baseTime += 3600 * 24 * 1
// 	test.SimulateOneBlock(lb, baseTime)

// 	test.SignCheckDeliver(t, lb, u1ClaimInterestMsg, 4, true, u1Priv, baseTime)
// 	test.SignCheckDeliver(t, lb, u2ClaimInterestMsg, 5, true, u2Priv, baseTime)

// 	test.CheckBalance(t, u1Name, lb, types.NewCoinFromInt64((69999+1.10088)*types.Decimals))
// 	test.CheckBalance(t, u2Name, lb, types.NewCoinFromInt64((59999+1.57332)*types.Decimals))

// }
