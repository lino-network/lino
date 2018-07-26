package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	post "github.com/lino-network/lino/x/post"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// test publish a normal post
func TestNormalPublish(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	postID1 := "New Post 1"
	postID2 := "New Post 2"
	// recover some stake
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, newAccountAppPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID1, 0, newAccountAppPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, postID2, 1, newAccountTransactionPriv, "", "", "", "", "0", baseTime)
}

// test publish a repost
func TestNormalRepost(t *testing.T) {
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), secp256k1.GenPrivKey(), newAccountAppPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountAppPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, repostID, 1, newAccountAppPriv,
		newAccountName, postID, "", "", "0", baseTime)

}

// test invalid repost if source post id doesn't exist
func TestInvalidRepost(t *testing.T) {
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), secp256k1.GenPrivKey(), newAccountAppPriv, "100")

	msg := post.CreatePostMsg{
		PostID:                  postID,
		Title:                   string(make([]byte, 50)),
		Content:                 string(make([]byte, 1000)),
		Author:                  types.AccountKey(newAccountName),
		RedistributionSplitRate: "0",
	}
	// reject due to stake
	test.SignCheckDeliver(t, lb, msg, 0, true, newAccountAppPriv, baseTime)
	msg.SourceAuthor = types.AccountKey(newAccountName)
	msg.SourcePostID = "invalid"
	msg.PostID = repostID
	// invalid source post id
	test.SignCheckDeliver(t, lb, msg, 1, false, newAccountAppPriv, baseTime)
}

// test publish a comment
func TestComment(t *testing.T) {
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	postID := "New Post"
	comment := "Comment"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), secp256k1.GenPrivKey(), newAccountAppPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountAppPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, comment, 1, newAccountAppPriv,
		"", "", newAccountName, postID, "0", baseTime)
}
