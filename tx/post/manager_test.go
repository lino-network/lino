package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

// test create post
func TestCreatePost(t *testing.T) {
	ctx, _, pm, _ := setupTest(t, 1)
	author := types.AccountKey("author")
	postID := "TestPostID"
	// test valid postInfo
	postCreateParams := PostCreateParams{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       author,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []types.IDToURLMapping{},
		RedistributionSplitRate: sdk.ZeroRat,
	}
	err := pm.CreatePost(ctx, &postCreateParams)
	assert.Nil(t, err)

	postInfo := model.PostInfo{
		PostID:       postCreateParams.PostID,
		Title:        postCreateParams.Title,
		Content:      postCreateParams.Content,
		Author:       postCreateParams.Author,
		ParentAuthor: postCreateParams.ParentAuthor,
		ParentPostID: postCreateParams.ParentPostID,
		SourceAuthor: postCreateParams.SourceAuthor,
		SourcePostID: postCreateParams.SourcePostID,
		Links:        postCreateParams.Links,
	}

	postMeta := model.PostMeta{
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}
	checkPostKVStore(t, ctx, types.GetPostKey(postCreateParams.Author, postCreateParams.PostID), postInfo, postMeta)
}
