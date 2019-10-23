package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	// "github.com/lino-network/lino/types"
	// post "github.com/lino-network/lino/x/post"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// test publish a normal post
func TestNormalPublish(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	postID1 := "New Post 1"
	postID2 := "New Post 2"
	// recover some coin day
	baseT := time.Unix(0, 0).Add(3600 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, newAccountAppPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID1, 0, newAccountTransactionPriv, baseTime)
	test.SimulateOneBlock(lb, baseTime+test.PostIntervalSec)
	test.CreateTestPost(
		t, lb, newAccountName, postID2, 1, newAccountTransactionPriv, baseTime+test.PostIntervalSec)
}
