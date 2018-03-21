package post

import (
	acc "github.com/lino-network/lino/tx/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPost(t *testing.T) {
	pm := newPostManager()
	user := acc.AccountKey("user")
	postID := "post ID"

	post := NewPost(user, postID, &pm)
	assert.Equal(t, user, post.GetAuthor())
	assert.Equal(t, postID, post.GetPostID())
	assert.Equal(t, GetPostKey(user, postID), post.GetPostKey())
	assert.NotNil(t, post.postManager)
	assert.Nil(t, post.postInfo)
	assert.Nil(t, post.postMeta)
	assert.Nil(t, post.postLikes)
	assert.Nil(t, post.postComments)
	assert.Nil(t, post.postViews)
	assert.Nil(t, post.postDonations)
	assert.False(t, post.writePostInfo)
	assert.False(t, post.writePostMeta)
	assert.False(t, post.writePostLikes)
	assert.False(t, post.writePostComments)
	assert.False(t, post.writePostViews)
	assert.False(t, post.writePostDonations)
}
