package post

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	author := types.AccountKey("TestAuthor")
	postInfo := CreateInfo{
		Permlink:       "TestPermlink",
		Title:          "TestTitle",
		Content:        "TestContent",
		Author:         author,
		ParentAuthor:   "",
		ParentPermlink: "",
		SourceAuthor:   "",
		SourcePermlink: "",
		Links:          nil,
	}
	createMsg := NewCreateMsg(postInfo)
	result := createMsg.ValidateBasic()
	assert.Nil(t, result)

	postInfo = CreateInfo{
		Permlink:       "",
		Title:          "TestTitle",
		Content:        "TestContent",
		Author:         author,
		ParentAuthor:   "",
		ParentPermlink: "",
		SourceAuthor:   "",
		SourcePermlink: "",
		Links:          nil,
	}

	createMsg = NewCreateMsg(postInfo)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostCreateNoPermlink())
}
