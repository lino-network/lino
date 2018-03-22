package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

func privAndBank() (crypto.PrivKey, *acc.AccountBank) {
	priv := crypto.GenPrivKeyEd25519()
	accBank := &acc.AccountBank{
		Address: priv.PubKey().Address(),
		Balance: sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	return priv.Wrap(), accBank
}

func TestHandlerCreatePost(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	handler := NewHandler(pm, lam)

	priv, bank := privAndBank()
	user := acc.AccountKey("testuser")
	_, err := lam.CreateAccount(ctx, user, priv.PubKey(), bank)
	assert.Nil(t, err)

	// test valid post
	postInfo := PostInfo{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	msg := NewCreatePostMsg(postInfo)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	// after handler check KVStore
	postMeta := PostMeta{
		Created:      0,
		LastUpdate:   0,
		LastActivity: 0,
		AllowReplies: true,
	}
	postViews := PostViews{Views: []View{}}
	postLikes := PostLikes{Likes: []Like{}}
	postComments := PostComments{Comments: []PostKey{}}
	postDonations := PostDonations{Donations: []Donation{}, Reward: sdk.Coins{}}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "TestPostID"), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invlaid author
	postInfo.Author = acc.AccountKey("invalid")
	msg = NewCreatePostMsg(postInfo)
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostCreateNonExistAuthor().Result())
}

func TestHandlerCreateComment(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	handler := NewHandler(pm, lam)
	priv, bank := privAndBank()
	user := acc.AccountKey("testuser")
	_, err := lam.CreateAccount(ctx, user, priv.PubKey(), bank)
	assert.Nil(t, err)

	postInfo := PostInfo{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	msg := NewCreatePostMsg(postInfo)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// test comment
	postInfo.Author = user
	postInfo.PostID = "comment"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "TestPostID"
	msg = NewCreatePostMsg(postInfo)
	ctx = ctx.WithBlockHeight(1)
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check comment
	postMeta := PostMeta{
		Created:      1,
		LastUpdate:   1,
		LastActivity: 1,
		AllowReplies: true,
	}
	postViews := PostViews{Views: []View{}}
	postLikes := PostLikes{Likes: []Like{}}
	postComments := PostComments{Comments: []PostKey{}}
	postDonations := PostDonations{Donations: []Donation{}, Reward: sdk.Coins{}}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "comment"), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// check parent
	postInfo.PostID = "TestPostID"
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.Created = 0
	postMeta.LastUpdate = 0
	postComments = PostComments{Comments: []PostKey{GetPostKey(user, "comment")}}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "TestPostID"), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invalid parent
	postInfo.PostID = "invalid post"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "invalid parent"
	msg = NewCreatePostMsg(postInfo)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostCommentsNotFound(GetPostKey(user, "invalid parent")).Result())

	// test duplicate comment
	postInfo.Author = user
	postInfo.PostID = "comment"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "TestPostID"
	msg = NewCreatePostMsg(postInfo)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostExist().Result())

	// test cycle comment
	postInfo.Author = user
	postInfo.PostID = "newComment"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "newComment"
	msg = NewCreatePostMsg(postInfo)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostCommentsNotFound(GetPostKey(user, "newComment")).Result())
}
