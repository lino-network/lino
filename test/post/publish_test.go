package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	post "github.com/lino-network/lino/x/post"

	crypto "github.com/tendermint/tendermint/crypto"
)

// test publish a normal post
func TestNormalPublish(t *testing.T) {
	newAccountTransactionPriv := crypto.GenPrivKeyEd25519()
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newuser"
	postID1 := "New Post 1"
	postID2 := "New Post 2"
	// recover some stake
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), newAccountTransactionPriv, crypto.GenPrivKeyEd25519(), newAccountPostPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID1, 0, newAccountPostPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, postID2, 1, newAccountTransactionPriv, "", "", "", "", "0", baseTime)
}

// test publish a repost
func TestNormalRepost(t *testing.T) {
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newuser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), newAccountPostPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountPostPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, repostID, 1, newAccountPostPriv,
		newAccountName, postID, "", "", "0", baseTime)

}

// test invalid repost if source post id doesn't exist
func TestInvalidRepost(t *testing.T) {
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newuser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), newAccountPostPriv, "100")

	msg := post.CreatePostMsg{
		PostID:                  postID,
		Title:                   string(make([]byte, 50)),
		Content:                 string(make([]byte, 1000)),
		Author:                  types.AccountKey(newAccountName),
		RedistributionSplitRate: "0",
	}
	// reject due to stake
	test.SignCheckDeliver(t, lb, msg, 0, true, newAccountPostPriv, baseTime)
	msg.SourceAuthor = types.AccountKey(newAccountName)
	msg.SourcePostID = "invalid"
	msg.PostID = repostID
	// invalid source post id
	test.SignCheckDeliver(t, lb, msg, 1, false, newAccountPostPriv, baseTime)
}

// test publish a comment
func TestComment(t *testing.T) {
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newuser"
	postID := "New Post"
	comment := "Comment"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), newAccountPostPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountPostPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, comment, 1, newAccountPostPriv,
		"", "", newAccountName, postID, "0", baseTime)
}
