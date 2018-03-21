package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func newAmount(amount int64) sdk.Coins {
	return sdk.Coins{
		{"lino", amount},
	}
}

func testDonationValidate(t *testing.T, donateMsg DonateMsg, expectError sdk.Error) {
	result := donateMsg.ValidateBasic()
	assert.Equal(t, result, expectError)
}

func testLikeValidate(t *testing.T, likeMsg LikeMsg, expectError sdk.Error) {
	result := likeMsg.ValidateBasic()
	assert.Equal(t, result, expectError)
}

func TestCreatePostMsg(t *testing.T) {
	author := types.AccountKey("TestAuthor")
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

func TestLikeMsg(t *testing.T) {
	cases := []struct {
		likeMsg     LikeMsg
		expectError sdk.Error
	}{
		{NewLikeMsg(types.AccountKey("test"), 10000, types.AccountKey("author"), "postID"), nil},
		{NewLikeMsg(types.AccountKey("test"), -10000, types.AccountKey("author"), "postID"), nil},
		{NewLikeMsg(types.AccountKey("test"), 10001, types.AccountKey("author"), "postID"), ErrPostLikeWeightOverflow(10001)},
		{NewLikeMsg(types.AccountKey("test"), -10001, types.AccountKey("author"), "postID"), ErrPostLikeWeightOverflow(-10001)},
		{NewLikeMsg(types.AccountKey(""), 10000, types.AccountKey("author"), "postID"), ErrPostLikeNoUsername()},
		{NewLikeMsg(types.AccountKey("test"), 10000, types.AccountKey(""), "postID"), ErrPostLikeInvalidTarget()},
		{NewLikeMsg(types.AccountKey("test"), 10000, types.AccountKey("author"), ""), ErrPostLikeInvalidTarget()},
		{NewLikeMsg(types.AccountKey("test"), 10000, types.AccountKey(""), ""), ErrPostLikeInvalidTarget()},
	}

	for _, cs := range cases {
		testLikeValidate(t, cs.likeMsg, cs.expectError)
	}
}

func TestDonationMsg(t *testing.T) {
	cases := []struct {
		donateMsg   DonateMsg
		expectError sdk.Error
	}{
		{NewDonateMsg(types.AccountKey("test"), newAmount(1), types.AccountKey("author"), "postID"), nil},
		{NewDonateMsg(types.AccountKey(""), newAmount(1), types.AccountKey("author"), "postID"), ErrPostLikeNoUsername()},
		{NewDonateMsg(types.AccountKey("test"), newAmount(0), types.AccountKey("author"), "postID"), bank.ErrInvalidCoins("0lino")},
		{NewDonateMsg(types.AccountKey("test"), newAmount(-1), types.AccountKey("author"), "postID"), bank.ErrInvalidCoins("-1lino")},
		{NewDonateMsg(types.AccountKey("test"), newAmount(1), types.AccountKey("author"), ""), ErrPostLikeInvalidTarget()},
		{NewDonateMsg(types.AccountKey("test"), newAmount(1), types.AccountKey(""), "postID"), ErrPostLikeInvalidTarget()},
		{NewDonateMsg(types.AccountKey("test"), newAmount(1), types.AccountKey(""), ""), ErrPostLikeInvalidTarget()},
	}

	for _, cs := range cases {
		testDonationValidate(t, cs.donateMsg, cs.expectError)
	}
}
