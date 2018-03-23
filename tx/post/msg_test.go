package post

import (
	"testing"

	acc "github.com/lino-network/lino/tx/account"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	author := acc.AccountKey("TestAuthor")
	// test valid post
	post := PostInfo{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       author,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
	}
	createMsg := NewCreatePostMsg(post)
	result := createMsg.ValidateBasic()
	assert.Nil(t, result)

	// test missing post id
	post.PostID = ""

	createMsg = NewCreatePostMsg(post)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostCreateNoPostID())

	post.Author = ""
	post.PostID = "testPost"
	createMsg = NewCreatePostMsg(post)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostCreateNoAuthor())

	// test exceeding max title length
	post.Author = author
	post.Title = string(make([]byte, 51))
	createMsg = NewCreatePostMsg(post)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostTitleExceedMaxLength())

	// test exceeding max content length
	post.Title = string(make([]byte, 50))
	post.Content = string(make([]byte, 1001))
	createMsg = NewCreatePostMsg(post)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostContentExceedMaxLength())
}
