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

func testReportOrUpvoteValidate(t *testing.T, reportOrUpvoteMsg ReportOrUpvoteMsg, expectError sdk.Error) {
	result := reportOrUpvoteMsg.ValidateBasic()
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

func getCommentAndRepost(
	t *testing.T, parentAuthor, parentPostID, sourceAuthor, sourcePostID string) PostCreateParams {
	return PostCreateParams{
		PostID:                  "TestPostID",
		Title:                   string(make([]byte, 50)),
		Content:                 string(make([]byte, 1000)),
		Author:                  "author",
		ParentAuthor:            types.AccountKey(parentAuthor),
		ParentPostID:            parentPostID,
		SourceAuthor:            types.AccountKey(sourceAuthor),
		SourcePostID:            sourcePostID,
		RedistributionSplitRate: "0",
	}
}

func TestCreatePostMsg(t *testing.T) {
	sdk.ZeroRat.GT(sdk.ZeroRat)
	author := types.AccountKey("TestAuthor")
	cases := []struct {
		postCreateParams PostCreateParams
		expectResult     sdk.Error
	}{
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 50)), Content: string(make([]byte, 1000)),
			Author: author, Links: []types.IDToURLMapping{}, RedistributionSplitRate: "0"}, expectResult: nil},
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 50)), Content: string(make([]byte, 1000)),
			Author: author, Links: []types.IDToURLMapping{}, RedistributionSplitRate: "1"}, expectResult: nil},
		{postCreateParams: PostCreateParams{
			PostID: "", Title: string(make([]byte, 50)), Content: string(make([]byte, 1000)),
			Author: author, Links: []types.IDToURLMapping{}, RedistributionSplitRate: "0"},
			expectResult: ErrPostCreateNoPostID()},
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 50)), Content: string(make([]byte, 1000)),
			Author: "", Links: []types.IDToURLMapping{}, RedistributionSplitRate: "0"},
			expectResult: ErrPostCreateNoAuthor()},
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 51)), Content: string(make([]byte, 1000)),
			Author: author, Links: []types.IDToURLMapping{}, RedistributionSplitRate: "0"},
			expectResult: ErrPostTitleExceedMaxLength()},
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 50)), Content: string(make([]byte, 1001)),
			Author: author, Links: []types.IDToURLMapping{}, RedistributionSplitRate: "0"},
			expectResult: ErrPostContentExceedMaxLength()},
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 50)), Content: string(make([]byte, 1000)),
			Author: author, Links: []types.IDToURLMapping{}, RedistributionSplitRate: "-1"},
			expectResult: ErrPostRedistributionSplitRate()},
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 50)), Content: string(make([]byte, 1000)),
			Author: author, Links: []types.IDToURLMapping{}, RedistributionSplitRate: "1.01"},
			expectResult: ErrPostRedistributionSplitRate()},
	}
	for _, cs := range cases {
		createMsg := NewCreatePostMsg(cs.postCreateParams)
		result := createMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectResult)
	}
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
		{NewLikeMsg(types.AccountKey("test"), 10001, types.AccountKey("author"), "postID"),
			ErrPostLikeWeightOverflow(10001)},
		{NewLikeMsg(types.AccountKey("test"), -10001, types.AccountKey("author"), "postID"),
			ErrPostLikeWeightOverflow(-10001)},
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
		{NewDonateMsg(types.AccountKey("test"), types.LNO("1"),
			types.AccountKey("author"), "postID", ""), nil},
		{NewDonateMsg(types.AccountKey(""), types.LNO("1"), types.AccountKey("author"), "postID", ""),
			ErrPostDonateNoUsername()},
		{NewDonateMsg(types.AccountKey("test"), types.LNO("0"), types.AccountKey("author"), "postID", ""),
			sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
		{NewDonateMsg(types.AccountKey("test"), types.LNO("-1"), types.AccountKey("author"), "postID", ""),
			sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
		{NewDonateMsg(types.AccountKey("test"), types.LNO("1"), types.AccountKey("author"), "", ""),
			ErrPostDonateInvalidTarget()},
		{NewDonateMsg(types.AccountKey("test"), types.LNO("1"), types.AccountKey(""), "postID", ""),
			ErrPostDonateInvalidTarget()},
		{NewDonateMsg(types.AccountKey("test"), types.LNO("1"), types.AccountKey(""), "", ""),
			ErrPostDonateInvalidTarget()},
	}

	for _, cs := range cases {
		testDonationValidate(t, cs.donateMsg, cs.expectError)
	}
}

func TestReportOrUpvoteMsg(t *testing.T) {
	cases := []struct {
		reportOrUpvoteMsg ReportOrUpvoteMsg
		expectError       sdk.Error
	}{
		{NewReportOrUpvoteMsg(types.AccountKey("test"), types.AccountKey("author"), "postID", true, false), nil},
		{NewReportOrUpvoteMsg(types.AccountKey("test"), types.AccountKey("author"), "postID", false, false), nil},
		{NewReportOrUpvoteMsg(types.AccountKey(""), types.AccountKey("author"), "postID", true, false),
			ErrPostReportOrUpvoteNoUsername()},
		{NewReportOrUpvoteMsg(types.AccountKey("test"), types.AccountKey("author"), "", true, true),
			ErrPostReportOrUpvoteInvalidTarget()},
		{NewReportOrUpvoteMsg(types.AccountKey("test"), types.AccountKey(""), "postID", false, true),
			ErrPostReportOrUpvoteInvalidTarget()},
		{NewReportOrUpvoteMsg(types.AccountKey("test"), types.AccountKey(""), "", false, false),
			ErrPostReportOrUpvoteInvalidTarget()},
	}

	for _, cs := range cases {
		testReportOrUpvoteValidate(t, cs.reportOrUpvoteMsg, cs.expectError)
	}
}
