package post

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	author := types.AccountKey("TestAuthor")
	// test valid post
	postInfo := CreateInfo{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       author,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
	}
	createMsg := NewCreateMsg(postInfo)
	result := createMsg.ValidateBasic()
	assert.Nil(t, result)

	// test missing post id
	postInfo.PostID = ""

	createMsg = NewCreateMsg(postInfo)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostCreateNoPostID())

	postInfo.Author = ""
	postInfo.PostID = "testPost"
	createMsg = NewCreateMsg(postInfo)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostCreateNoAuthor())

	// test exceeding max title length

	postInfo.Author = author
	postInfo.Title = string(make([]byte, 51))
	createMsg = NewCreateMsg(postInfo)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostTitleExceedMaxLength())

	// test exceeding max content length
	postInfo.Title = string(make([]byte, 50))
	postInfo.Content = string(make([]byte, 1001))
	createMsg = NewCreateMsg(postInfo)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostContentExceedMaxLength())
}
