package post

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	memo1       = "memo1"
	invalidMemo = "Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! "
)

func testDonationValidate(t *testing.T, donateMsg DonateMsg, expectError sdk.Error) {
	result := donateMsg.ValidateBasic()
	assert.Equal(t, expectError, result)
}

func testReportOrUpvoteValidate(t *testing.T, reportOrUpvoteMsg ReportOrUpvoteMsg, expectError sdk.Error) {
	result := reportOrUpvoteMsg.ValidateBasic()
	assert.Equal(t, expectError, result)
}

func testViewValidate(t *testing.T, viewMsg ViewMsg, expectError sdk.Error) {
	result := viewMsg.ValidateBasic()
	assert.Equal(t, expectError, result)
}

func testLikeValidate(t *testing.T, likeMsg LikeMsg, expectError sdk.Error) {
	result := likeMsg.ValidateBasic()
	assert.Equal(t, expectError, result)
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
			expectResult: ErrNoPostID()},
		{postCreateParams: PostCreateParams{
			PostID: "TestPostID", Title: string(make([]byte, 50)), Content: string(make([]byte, 1000)),
			Author: "", Links: []types.IDToURLMapping{}, RedistributionSplitRate: "0"},
			expectResult: ErrNoAuthor()},
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

func TestUpdatePostMsg(t *testing.T) {
	cases := []struct {
		updatePostMsg UpdatePostMsg
		expectResult  sdk.Error
	}{
		{updatePostMsg: NewUpdatePostMsg(
			"author", "postID", "title", "content", []types.IDToURLMapping{}, "1.0"), expectResult: nil},
		{updatePostMsg: NewUpdatePostMsg(
			"author", "postID", "title", "content", []types.IDToURLMapping{}, "0"), expectResult: nil},
		{updatePostMsg: NewUpdatePostMsg(
			"", "postID", "title", "content", []types.IDToURLMapping{}, "0"), expectResult: ErrNoAuthor()},
		{updatePostMsg: NewUpdatePostMsg(
			"author", "", "title", "content", []types.IDToURLMapping{}, "0"), expectResult: ErrNoPostID()},
		{updatePostMsg: NewUpdatePostMsg(
			"author", "postID", string(make([]byte, 51)), "content", []types.IDToURLMapping{}, "0"),
			expectResult: ErrPostTitleExceedMaxLength()},
		{updatePostMsg: NewUpdatePostMsg(
			"author", "postID", string(make([]byte, 50)), string(make([]byte, 1001)),
			[]types.IDToURLMapping{}, "0"), expectResult: ErrPostContentExceedMaxLength()},
		{updatePostMsg: NewUpdatePostMsg(
			"author", "postID", string(make([]byte, 50)), string(make([]byte, 1000)),
			[]types.IDToURLMapping{}, "1.01"), expectResult: ErrPostRedistributionSplitRate()},
		{updatePostMsg: NewUpdatePostMsg(
			"author", "postID", string(make([]byte, 50)), string(make([]byte, 1000)),
			[]types.IDToURLMapping{}, "-1"), expectResult: ErrPostRedistributionSplitRate()},
	}
	for _, cs := range cases {
		result := cs.updatePostMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectResult)
	}
}

func TestDeletePostMsg(t *testing.T) {
	testCases := map[string]struct {
		msg         DeletePostMsg
		wantErrCode sdk.CodeType
	}{
		"normal case": {
			msg: DeletePostMsg{
				Author: "author",
				PostID: "postID",
			},
			wantErrCode: sdk.CodeOK,
		},
		"empty author": {
			msg: DeletePostMsg{
				Author: "",
				PostID: "postID",
			},
			wantErrCode: types.CodePostMsgError,
		},
		"empty postID": {
			msg: DeletePostMsg{
				Author: "author",
				PostID: "",
			},
			wantErrCode: types.CodePostMsgError,
		},
	}
	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()
		if got == nil && tc.wantErrCode != sdk.CodeOK {
			t.Errorf("%s ValidateBasic: got %v, want %v", testName, got, tc.wantErrCode)
		}
		if got != nil {
			if got.ABCICode() != tc.wantErrCode {
				t.Errorf("%s ValidateBasic: got %v, want %v", testName, got, tc.wantErrCode)
			}
		}
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
		{NewLikeMsg("test", 10000, "author", "postID"), nil},
		{NewLikeMsg("test", -10000, "author", "postID"), nil},
		{NewLikeMsg("test", 10001, "author", "postID"),
			ErrPostLikeWeightOverflow(10001)},
		{NewLikeMsg("test", -10001, "author", "postID"),
			ErrPostLikeWeightOverflow(-10001)},
		{NewLikeMsg("", 10000, "author", "postID"), ErrPostLikeNoUsername()},
		{NewLikeMsg("test", 10000, "", "postID"), ErrPostLikeInvalidTarget()},
		{NewLikeMsg("test", 10000, "author", ""), ErrPostLikeInvalidTarget()},
		{NewLikeMsg("test", 10000, "", ""), ErrPostLikeInvalidTarget()},
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
		{NewDonateMsg("test", types.LNO("1"),
			"author", "postID", "", false, memo1), nil},
		{NewDonateMsg("", types.LNO("1"), "author", "postID", "", false, memo1),
			ErrPostDonateNoUsername()},
		{NewDonateMsg("test", types.LNO("0"), "author", "postID", "", false, memo1),
			sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
		{NewDonateMsg("test", types.LNO("-1"), "author", "postID", "", false, memo1),
			sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
		{NewDonateMsg("test", types.LNO("1"), "author", "", "", false, memo1),
			ErrPostDonateInvalidTarget()},
		{NewDonateMsg("test", types.LNO("1"), "", "postID", "", false, memo1),
			ErrPostDonateInvalidTarget()},
		{NewDonateMsg("test", types.LNO("1"), "", "", "", false, memo1),
			ErrPostDonateInvalidTarget()},
		{NewDonateMsg("test", types.LNO("1"), "author", "postID", "", false, invalidMemo),
			ErrInvalidMemo()},
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
		{NewReportOrUpvoteMsg("test", "author", "postID", true), nil},
		{NewReportOrUpvoteMsg("test", "author", "postID", false), nil},
		{NewReportOrUpvoteMsg("", "author", "postID", true),
			ErrPostReportOrUpvoteNoUsername()},
		{NewReportOrUpvoteMsg("test", "author", "", true),
			ErrPostReportOrUpvoteInvalidTarget()},
		{NewReportOrUpvoteMsg("test", "", "postID", false),
			ErrPostReportOrUpvoteInvalidTarget()},
		{NewReportOrUpvoteMsg("test", "", "", false),
			ErrPostReportOrUpvoteInvalidTarget()},
	}

	for _, cs := range cases {
		testReportOrUpvoteValidate(t, cs.reportOrUpvoteMsg, cs.expectError)
	}
}

func TestViewMsg(t *testing.T) {
	cases := []struct {
		viewMsg     ViewMsg
		expectError sdk.Error
	}{
		{NewViewMsg("test", "author", "postID"), nil},
		{NewViewMsg("", "author", "postID"),
			ErrPostViewNoUsername()},
		{NewViewMsg("test", "", "postID"),
			ErrPostViewInvalidTarget()},
		{NewViewMsg("test", "author", ""),
			ErrPostViewInvalidTarget()},
	}

	for _, cs := range cases {
		testViewValidate(t, cs.viewMsg, cs.expectError)
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"donateMsg from saving": {
			msg: NewDonateMsg(
				"test", types.LNO("1"),
				"author", "postID", "", false, memo1),
			expectPermission: types.TransactionPermission,
		},
		"donateMsg from checking": {
			msg: NewDonateMsg(
				"test", types.LNO("1"),
				"author", "postID", "", true, memo1),
			expectPermission: types.PostPermission,
		},
		"create post": {
			msg: NewCreatePostMsg(PostCreateParams{
				PostID:       "test",
				Title:        "title",
				Content:      "content",
				Author:       "author",
				ParentAuthor: types.AccountKey("parentAuthor"),
				ParentPostID: "parentPostID",
				SourceAuthor: types.AccountKey("sourceAuthor"),
				SourcePostID: "sourcePostID",
				Links: []types.IDToURLMapping{
					types.IDToURLMapping{
						Identifier: "#1",
						URL:        "https://lino.network",
					},
				},
				RedistributionSplitRate: "0.5",
			}),
			expectPermission: types.PostPermission,
		},
		"like post": {
			msg: NewLikeMsg(
				"test", 10000, "author", "postID"),
			expectPermission: types.PostPermission,
		},
		"view post": {
			msg: NewViewMsg(
				"test", "author", "postID"),
			expectPermission: types.PostPermission,
		},
		"report post": {
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", true),
			expectPermission: types.PostPermission,
		},
		"upvote post": {
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", false),
			expectPermission: types.PostPermission,
		},
		"update post": {
			msg: NewUpdatePostMsg(
				"author", "postID", "title", "content", []types.IDToURLMapping{}, "0"),
			expectPermission: types.PostPermission,
		},
	}

	for testName, cs := range cases {
		permissionLevel := cs.msg.Get(types.PermissionLevel)
		if permissionLevel == nil {
			if cs.expectPermission != types.PostPermission {
				t.Errorf(
					"%s: expect permission incorrect, expect %v, got %v",
					testName, cs.expectPermission, types.PostPermission)
				return
			} else {
				continue
			}
		}
		permission, ok := permissionLevel.(types.Permission)
		assert.Equal(t, ok, true)
		if cs.expectPermission != permission {
			t.Errorf(
				"%s: expect permission incorrect, expect %v, got %v",
				testName, cs.expectPermission, permission)
			return
		}
	}
}
