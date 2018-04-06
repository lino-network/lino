package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func testDonationValidate(t *testing.T, donateMsg DonateMsg, expectError sdk.Error) {
	result := donateMsg.ValidateBasic()
	assert.Equal(t, result, expectError)
}

func testLikeValidate(t *testing.T, likeMsg LikeMsg, expectError sdk.Error) {
	result := likeMsg.ValidateBasic()
	assert.Equal(t, result, expectError)
}

func testCommentAndRepostValidate(t *testing.T, postCreateParams PostCreateParams, expectError sdk.Error) {
	createMsg := NewCreatePostMsg(postCreateParams)
	result := createMsg.ValidateBasic()
	assert.Equal(t, expectError, result)
}

func getCommentAndRepost(t *testing.T, parentAuthor, parentPostID, sourceAuthor, sourcePostID string) PostCreateParams {
	return PostCreateParams{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       "author",
		ParentAuthor: types.AccountKey(parentAuthor),
		ParentPostID: parentPostID,
		SourceAuthor: types.AccountKey(sourceAuthor),
		SourcePostID: sourcePostID,
	}
}

func TestCreatePostMsg(t *testing.T) {
	author := types.AccountKey("TestAuthor")
	// test valid post
	postCreateParams := PostCreateParams{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       author,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
	}
	createMsg := NewCreatePostMsg(postCreateParams)
	result := createMsg.ValidateBasic()
	assert.Nil(t, result)

	// test missing post id
	postCreateParams.PostID = ""

	createMsg = NewCreatePostMsg(postCreateParams)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostCreateNoPostID())

	postCreateParams.Author = ""
	postCreateParams.PostID = "testPost"
	createMsg = NewCreatePostMsg(postCreateParams)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostCreateNoAuthor())

	// test exceeding max title length
	postCreateParams.Author = author
	postCreateParams.Title = string(make([]byte, 51))
	createMsg = NewCreatePostMsg(postCreateParams)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostTitleExceedMaxLength())

	// test exceeding max content length
	postCreateParams.Title = string(make([]byte, 50))
	postCreateParams.Content = string(make([]byte, 1001))
	createMsg = NewCreatePostMsg(postCreateParams)
	result = createMsg.ValidateBasic()
	assert.Equal(t, result, ErrPostContentExceedMaxLength())
}

func TestCommentAndRepost(t *testing.T) {
	parentAuthor := "Parent"
	parentPostID := "ParentPostID"
	sourceAuthor := "Source"
	sourcePostID := "SourcePostID"

	cases := []struct {
		postCreateParams PostCreateParams
		expectError      sdk.Error
	}{
		{getCommentAndRepost(t, "", "", "", ""), nil},
		{getCommentAndRepost(t, parentAuthor, parentPostID, "", ""), nil},
		{getCommentAndRepost(t, "", "", sourceAuthor, sourcePostID), nil},
		{getCommentAndRepost(t, parentAuthor, parentPostID, sourceAuthor, sourcePostID), ErrCommentAndRepostError()},
		{getCommentAndRepost(t, parentAuthor, parentPostID, sourceAuthor, ""), ErrCommentAndRepostError()},
		{getCommentAndRepost(t, parentAuthor, parentPostID, "", sourcePostID), ErrCommentAndRepostError()},
		{getCommentAndRepost(t, parentAuthor, "", sourceAuthor, sourcePostID), ErrCommentAndRepostError()},
		{getCommentAndRepost(t, "", parentPostID, sourceAuthor, sourcePostID), ErrCommentAndRepostError()},
		{getCommentAndRepost(t, parentAuthor, "", sourceAuthor, ""), ErrCommentAndRepostError()},
	}
	for _, cs := range cases {
		testCommentAndRepostValidate(t, cs.postCreateParams, cs.expectError)
	}
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
		{NewDonateMsg(types.AccountKey("test"), types.LNO(sdk.NewRat(1)), types.AccountKey("author"), "postID"), nil},
		{NewDonateMsg(types.AccountKey(""), types.LNO(sdk.NewRat(1)), types.AccountKey("author"), "postID"), ErrPostDonateNoUsername()},
		{NewDonateMsg(types.AccountKey("test"), types.LNO(sdk.NewRat(0)), types.AccountKey("author"), "postID"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
		{NewDonateMsg(types.AccountKey("test"), types.LNO(sdk.NewRat(-1)), types.AccountKey("author"), "postID"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
		{NewDonateMsg(types.AccountKey("test"), types.LNO(sdk.NewRat(1)), types.AccountKey("author"), ""), ErrPostDonateInvalidTarget()},
		{NewDonateMsg(types.AccountKey("test"), types.LNO(sdk.NewRat(1)), types.AccountKey(""), "postID"), ErrPostDonateInvalidTarget()},
		{NewDonateMsg(types.AccountKey("test"), types.LNO(sdk.NewRat(1)), types.AccountKey(""), ""), ErrPostDonateInvalidTarget()},
	}

	for _, cs := range cases {
		testDonationValidate(t, cs.donateMsg, cs.expectError)
	}
}
