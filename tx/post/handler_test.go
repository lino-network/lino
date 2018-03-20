package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

func privAndBank() (crypto.PrivKey, *types.AccountBank) {
	priv := crypto.GenPrivKeyEd25519()
	accBank := &types.AccountBank{
		Address: priv.PubKey().Address(),
		Coins:   sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	return priv.Wrap(), accBank
}

func TestHandlerCreatePost(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	handler := NewHandler(pm, lam)

	priv, bank := privAndBank()
	user := types.AccountKey("testuser")
	_, err := lam.CreateAccount(ctx, user, priv.PubKey(), bank)
	assert.Nil(t, err)

	// test valid post
	post := types.Post{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
	}
	msg := NewCreatePostMsg(post)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// test invlaid author
	post.Author = types.AccountKey("invalid")
	msg = NewCreatePostMsg(post)
	result = handler(ctx, msg)
	assert.Equal(t, result, acc.ErrAccountManagerFail("LinoAccountManager get meta failed: meta doesn't exist").Result())

	// test comment
	post.Author = user
	post.ParentAuthor = user
	post.ParentPostID = "TestPostID"
	post.PostID = "newPost"
	msg = NewCreatePostMsg(post)
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
}
