package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	post "github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// test publish a normal post
func TestNormalPublish(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID := "New Post"
	// recover some stake
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0, newAccountPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountPriv, "", "", "", "", sdk.ZeroRat, baseTime)
}

// test publish a repost
func TestNormalRepost(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0, newAccountPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountPriv, "", "", "", "", sdk.ZeroRat, baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, repostID, 1, newAccountPriv,
		newAccountName, postID, "", "", sdk.ZeroRat, baseTime)

}

// test invalid repost if source post id doesn't exist
func TestInvalidRepost(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0, newAccountPriv, "100")

	postCreateParams := post.PostCreateParams{
		PostID:                  postID,
		Title:                   string(make([]byte, 50)),
		Content:                 string(make([]byte, 1000)),
		Author:                  types.AccountKey(newAccountName),
		RedistributionSplitRate: sdk.ZeroRat,
	}
	msg := post.NewCreatePostMsg(postCreateParams)
	// reject due to stake
	test.SignCheckDeliver(t, lb, msg, 0, true, newAccountPriv, baseTime)
	postCreateParams.SourceAuthor = types.AccountKey(newAccountName)
	postCreateParams.SourcePostID = "invalid"
	postCreateParams.PostID = repostID
	msg = post.NewCreatePostMsg(postCreateParams)
	// invalid source post id
	test.SignCheckDeliver(t, lb, msg, 1, false, newAccountPriv, baseTime)
}

// test publish a comment
func TestComment(t *testing.T) {
	newAccountPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID := "New Post"
	comment := "Comment"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0, newAccountPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountPriv, "", "", "", "", sdk.ZeroRat, baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, comment, 1, newAccountPriv,
		"", "", newAccountName, postID, sdk.ZeroRat, baseTime)
}
