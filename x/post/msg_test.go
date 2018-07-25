package post

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	memo1       = "memo1"
	invalidMemo = "Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! "
)

func getCommentAndRepost(
	t *testing.T, parentAuthor, parentPostID, sourceAuthor, sourcePostID string) CreatePostMsg {
	return CreatePostMsg{
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
	testCases := []struct {
		testName       string
		msg            CreatePostMsg
		expectedResult sdk.Error
	}{
		{
			testName: "normal case 1",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "0",
			},
			expectedResult: nil,
		},
		{
			testName: "normal case 2",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "1",
			},
			expectedResult: nil,
		},
		{
			testName: "empty post id",
			msg: CreatePostMsg{
				PostID:  "",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "0",
			},
			expectedResult: ErrNoPostID(),
		},
		{
			testName: "no author",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  "",
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "0",
			},
			expectedResult: ErrNoAuthor(),
		},
		{
			testName: "post title is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 51)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "0",
			},
			expectedResult: ErrPostTitleExceedMaxLength(),
		},
		{
			testName: "post content is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1001)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "0",
			},
			expectedResult: ErrPostContentExceedMaxLength(),
		},
		{
			testName: "negative redistribution split rate is invalid",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "-1",
			},
			expectedResult: ErrInvalidPostRedistributionSplitRate(),
		},
		{
			testName: "redistribution split rate can't be bigger than 1",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "1.01",
			},
			expectedResult: ErrInvalidPostRedistributionSplitRate(),
		},
		{
			testName: "redistribution split rate length is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links:   []types.IDToURLMapping{},
				RedistributionSplitRate: "0.00000000001",
			},
			expectedResult: ErrRedistributionSplitRateLengthTooLong(),
		},
		{
			testName: "identifier length is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links: []types.IDToURLMapping{
					{
						Identifier: string(make([]byte, 21)),
						URL:        string(make([]byte, 50)),
					},
				},
				RedistributionSplitRate: "0",
			},
			expectedResult: ErrIdentifierLengthTooLong()},
		{
			testName: "url length is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 50)),
				Content: string(make([]byte, 1000)),
				Author:  author,
				Links: []types.IDToURLMapping{
					{
						Identifier: string(make([]byte, 20)),
						URL:        string(make([]byte, 51)),
					},
				},
				RedistributionSplitRate: "0",
			},
			expectedResult: ErrURLLengthTooLong(),
		},
	}
	for _, tc := range testCases {
		result := tc.msg.ValidateBasic()
		if !assert.Equal(t, tc.expectedResult, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedResult)
		}
	}
}

func TestUpdatePostMsg(t *testing.T) {
	testCases := []struct {
		testName       string
		updatePostMsg  UpdatePostMsg
		expectedResult sdk.Error
	}{
		{
			testName: "normal case 1",
			updatePostMsg: NewUpdatePostMsg(
				"author", "postID", "title", "content", []types.IDToURLMapping{}, "1.0"),
			expectedResult: nil,
		},
		{
			testName: "normal case 2",
			updatePostMsg: NewUpdatePostMsg(
				"author", "postID", "title", "content", []types.IDToURLMapping{}, "0"),
			expectedResult: nil,
		},
		{
			testName: "no author",
			updatePostMsg: NewUpdatePostMsg(
				"", "postID", "title", "content", []types.IDToURLMapping{}, "0"),
			expectedResult: ErrNoAuthor(),
		},
		{
			testName: "no post id",
			updatePostMsg: NewUpdatePostMsg(
				"author", "", "title", "content", []types.IDToURLMapping{}, "0"),
			expectedResult: ErrNoPostID(),
		},
		{
			testName: "post tile is too long",
			updatePostMsg: NewUpdatePostMsg(
				"author", "postID", string(make([]byte, 51)), "content", []types.IDToURLMapping{}, "0"),
			expectedResult: ErrPostTitleExceedMaxLength(),
		},
		{
			testName: "post content is too long",
			updatePostMsg: NewUpdatePostMsg(
				"author", "postID", string(make([]byte, 50)), string(make([]byte, 1001)),
				[]types.IDToURLMapping{}, "0"),
			expectedResult: ErrPostContentExceedMaxLength(),
		},
		{
			testName: "redistribution split rate can't be bigger than 1",
			updatePostMsg: NewUpdatePostMsg(
				"author", "postID", string(make([]byte, 50)), string(make([]byte, 1000)),
				[]types.IDToURLMapping{}, "1.01"),
			expectedResult: ErrInvalidPostRedistributionSplitRate(),
		},
		{
			testName: "redistribution split rate can't be negative",
			updatePostMsg: NewUpdatePostMsg(
				"author", "postID", string(make([]byte, 50)), string(make([]byte, 1000)),
				[]types.IDToURLMapping{}, "-1"),
			expectedResult: ErrInvalidPostRedistributionSplitRate(),
		},
		{
			testName: "redistruction split rate length is too long",
			updatePostMsg: NewUpdatePostMsg(
				"author", "postID", string(make([]byte, 50)), string(make([]byte, 1000)),
				[]types.IDToURLMapping{}, "0.000000000001"),
			expectedResult: ErrRedistributionSplitRateLengthTooLong(),
		},
	}
	for _, tc := range testCases {
		result := tc.updatePostMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedResult) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedResult)
		}
	}
}

func TestDeletePostMsg(t *testing.T) {
	testCases := []struct {
		testName    string
		msg         DeletePostMsg
		wantErrCode sdk.CodeType
	}{
		{
			testName: "normal case",
			msg: DeletePostMsg{
				Author: "author",
				PostID: "postID",
			},
			wantErrCode: sdk.CodeOK,
		},
		{
			testName: "empty author",
			msg: DeletePostMsg{
				Author: "",
				PostID: "postID",
			},
			wantErrCode: types.CodeNoAuthor,
		},
		{
			testName: "empty postID",
			msg: DeletePostMsg{
				Author: "author",
				PostID: "",
			},
			wantErrCode: types.CodeNoPostID,
		},
	}
	for _, tc := range testCases {
		got := tc.msg.ValidateBasic()
		if got == nil && tc.wantErrCode != sdk.CodeOK {
			t.Errorf("%s: got non-OK code, got %v, want %v", tc.testName, got, tc.wantErrCode)
		}
		if got != nil {
			if got.Code() != tc.wantErrCode {
				t.Errorf("%s: diff err code, got %v, want %v", tc.testName, got, tc.wantErrCode)
			}
		}
	}
}

func TestCommentAndRepost(t *testing.T) {
	parentAuthor := "Parent"
	parentPostID := "ParentPostID"
	sourceAuthor := "Source"
	sourcePostID := "SourcePostID"

	testCases := []struct {
		testName      string
		msg           CreatePostMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case 1",
			msg:           getCommentAndRepost(t, "", "", "", ""),
			expectedError: nil,
		},
		{
			testName:      "normal case 2",
			msg:           getCommentAndRepost(t, parentAuthor, parentPostID, "", ""),
			expectedError: nil,
		},
		{
			testName:      "normal case 3",
			msg:           getCommentAndRepost(t, "", "", sourceAuthor, sourcePostID),
			expectedError: nil,
		},
		{
			testName:      "post can't be comment and re-post at the same time 1",
			msg:           getCommentAndRepost(t, parentAuthor, parentPostID, sourceAuthor, sourcePostID),
			expectedError: ErrCommentAndRepostConflict(),
		},
		{
			testName:      "post can't be comment and re-post at the same time 2",
			msg:           getCommentAndRepost(t, parentAuthor, parentPostID, sourceAuthor, ""),
			expectedError: ErrCommentAndRepostConflict(),
		},
		{
			testName:      "post can't be comment and re-post at the same time 3",
			msg:           getCommentAndRepost(t, parentAuthor, parentPostID, "", sourcePostID),
			expectedError: ErrCommentAndRepostConflict(),
		},
		{
			testName:      "post can't be comment and re-post at the same time 4",
			msg:           getCommentAndRepost(t, parentAuthor, "", sourceAuthor, sourcePostID),
			expectedError: ErrCommentAndRepostConflict(),
		},
		{
			testName:      "post can't be comment and re-post at the same time 5",
			msg:           getCommentAndRepost(t, "", parentPostID, sourceAuthor, sourcePostID),
			expectedError: ErrCommentAndRepostConflict(),
		},
		{
			testName:      "post can't be comment and re-post at the same time 6",
			msg:           getCommentAndRepost(t, parentAuthor, "", sourceAuthor, ""),
			expectedError: ErrCommentAndRepostConflict(),
		},
	}
	for _, tc := range testCases {
		result := tc.msg.ValidateBasic()
		if !assert.Equal(t, tc.expectedError, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestDonationMsg(t *testing.T) {
	testCases := []struct {
		testName      string
		donateMsg     DonateMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case",
			donateMsg:     NewDonateMsg("test", types.LNO("1"), "author", "postID", "", memo1),
			expectedError: nil,
		},
		{
			testName:      "no username",
			donateMsg:     NewDonateMsg("", types.LNO("1"), "author", "postID", "", memo1),
			expectedError: ErrNoUsername(),
		},
		{
			testName:      "zero coin is less than lower bound",
			donateMsg:     NewDonateMsg("test", types.LNO("0"), "author", "postID", "", memo1),
			expectedError: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName:      "negative coin is less than lower bound",
			donateMsg:     NewDonateMsg("test", types.LNO("-1"), "author", "postID", "", memo1),
			expectedError: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName:      "invalid target - no post id",
			donateMsg:     NewDonateMsg("test", types.LNO("1"), "author", "", "", memo1),
			expectedError: ErrInvalidTarget(),
		},
		{
			testName:      "invalid target - no author",
			donateMsg:     NewDonateMsg("test", types.LNO("1"), "", "postID", "", memo1),
			expectedError: ErrInvalidTarget(),
		},
		{
			testName:      "invalid target - no author and post id",
			donateMsg:     NewDonateMsg("test", types.LNO("1"), "", "", "", memo1),
			expectedError: ErrInvalidTarget(),
		},
		{
			testName:      "invalid memo",
			donateMsg:     NewDonateMsg("test", types.LNO("1"), "author", "postID", "", invalidMemo),
			expectedError: ErrInvalidMemo(),
		},
	}

	for _, tc := range testCases {
		result := tc.donateMsg.ValidateBasic()
		if !assert.Equal(t, tc.expectedError, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestReportOrUpvoteMsg(t *testing.T) {
	testCases := []struct {
		testName          string
		reportOrUpvoteMsg ReportOrUpvoteMsg
		expectedError     sdk.Error
	}{
		{
			testName:          "normal case - report",
			reportOrUpvoteMsg: NewReportOrUpvoteMsg("test", "author", "postID", true),
			expectedError:     nil,
		},
		{
			testName:          "normal case - upvote",
			reportOrUpvoteMsg: NewReportOrUpvoteMsg("test", "author", "postID", false),
			expectedError:     nil,
		},
		{
			testName:          "no username",
			reportOrUpvoteMsg: NewReportOrUpvoteMsg("", "author", "postID", true),
			expectedError:     ErrNoUsername(),
		},
		{
			testName:          "invalid target - no post id",
			reportOrUpvoteMsg: NewReportOrUpvoteMsg("test", "author", "", true),
			expectedError:     ErrInvalidTarget(),
		},
		{
			testName:          "invalid target - no author",
			reportOrUpvoteMsg: NewReportOrUpvoteMsg("test", "", "postID", false),
			expectedError:     ErrInvalidTarget(),
		},
		{
			testName:          "invalid target - no author and post id",
			reportOrUpvoteMsg: NewReportOrUpvoteMsg("test", "", "", false),
			expectedError:     ErrInvalidTarget(),
		},
	}

	for _, tc := range testCases {
		result := tc.reportOrUpvoteMsg.ValidateBasic()
		if !assert.Equal(t, tc.expectedError, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestViewMsg(t *testing.T) {
	testCases := []struct {
		testName      string
		viewMsg       ViewMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case",
			viewMsg:       NewViewMsg("test", "author", "postID"),
			expectedError: nil,
		},
		{
			testName:      "no username",
			viewMsg:       NewViewMsg("", "author", "postID"),
			expectedError: ErrNoUsername(),
		},
		{
			testName:      "invalid target - no author",
			viewMsg:       NewViewMsg("test", "", "postID"),
			expectedError: ErrInvalidTarget(),
		},
		{
			testName:      "invalid target - no post id",
			viewMsg:       NewViewMsg("test", "author", ""),
			expectedError: ErrInvalidTarget(),
		},
	}

	for _, tc := range testCases {
		result := tc.viewMsg.ValidateBasic()
		if !assert.Equal(t, tc.expectedError, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestMsgPermission(t *testing.T) {
	testCases := []struct {
		testName           string
		msg                types.Msg
		expectedPermission types.Permission
	}{
		{
			testName: "donateMsg",
			msg: NewDonateMsg(
				"test", types.LNO("1"),
				"author", "postID", "", memo1),
			expectedPermission: types.PreAuthorizationPermission,
		},
		{
			testName: "create post",
			msg: CreatePostMsg{
				PostID:       "test",
				Title:        "title",
				Content:      "content",
				Author:       "author",
				ParentAuthor: types.AccountKey("parentAuthor"),
				ParentPostID: "parentPostID",
				SourceAuthor: types.AccountKey("sourceAuthor"),
				SourcePostID: "sourcePostID",
				Links: []types.IDToURLMapping{
					{
						Identifier: "#1",
						URL:        "https://lino.network",
					},
				},
				RedistributionSplitRate: "0.5",
			},
			expectedPermission: types.AppPermission,
		},
		{
			testName: "view post",
			msg: NewViewMsg(
				"test", "author", "postID"),
			expectedPermission: types.AppPermission,
		},
		{
			testName: "report post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", true),
			expectedPermission: types.AppPermission,
		},
		{
			testName: "upvote post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", false),
			expectedPermission: types.AppPermission,
		},
		{
			testName: "update post",
			msg: NewUpdatePostMsg(
				"author", "postID", "title", "content", []types.IDToURLMapping{}, "0"),
			expectedPermission: types.AppPermission,
		},
	}

	for _, tc := range testCases {
		permission := tc.msg.GetPermission()
		if tc.expectedPermission != permission {
			t.Errorf("%s: diff permission, got %v, want %v", tc.testName, tc.expectedPermission, permission)
		}
	}
}

func TestGetSignBytes(t *testing.T) {
	testCases := []struct {
		testName string
		msg      types.Msg
	}{
		{
			testName: "donateMsg",
			msg: NewDonateMsg(
				"test", types.LNO("1"),
				"author", "postID", "", memo1),
		},
		{
			testName: "create post",
			msg: CreatePostMsg{
				PostID:       "test",
				Title:        "title",
				Content:      "content",
				Author:       "author",
				ParentAuthor: types.AccountKey("parentAuthor"),
				ParentPostID: "parentPostID",
				SourceAuthor: types.AccountKey("sourceAuthor"),
				SourcePostID: "sourcePostID",
				Links: []types.IDToURLMapping{
					{
						Identifier: "#1",
						URL:        "https://lino.network",
					},
				},
				RedistributionSplitRate: "0.5",
			},
		},
		{
			testName: "view post",
			msg: NewViewMsg(
				"test", "author", "postID"),
		},
		{
			testName: "report post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", true),
		},
		{
			testName: "upvote post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", false),
		},
		{
			testName: "update post",
			msg: NewUpdatePostMsg(
				"author", "postID", "title", "content", []types.IDToURLMapping{}, "0"),
		},
	}

	for _, tc := range testCases {
		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, tc.testName)
	}
}

func TestGetSigners(t *testing.T) {
	testCases := []struct {
		testName      string
		msg           types.Msg
		expectSigners []types.AccountKey
	}{
		{
			testName: "donateMsg",
			msg: NewDonateMsg(
				"test", types.LNO("1"),
				"author", "postID", "", memo1),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName: "create post",
			msg: CreatePostMsg{
				PostID:       "test",
				Title:        "title",
				Content:      "content",
				Author:       "author",
				ParentAuthor: types.AccountKey("parentAuthor"),
				ParentPostID: "parentPostID",
				SourceAuthor: types.AccountKey("sourceAuthor"),
				SourcePostID: "sourcePostID",
				Links: []types.IDToURLMapping{
					{
						Identifier: "#1",
						URL:        "https://lino.network",
					},
				},
				RedistributionSplitRate: "0.5",
			},
			expectSigners: []types.AccountKey{"author"},
		},
		{
			testName: "view post",
			msg: NewViewMsg(
				"test", "author", "postID"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName: "report post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", true),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName: "upvote post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", false),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName: "update post",
			msg: NewUpdatePostMsg(
				"author", "postID", "title", "content", []types.IDToURLMapping{}, "0"),
			expectSigners: []types.AccountKey{"author"},
		},
	}

	for _, tc := range testCases {
		if len(tc.msg.GetSigners()) != len(tc.expectSigners) {
			t.Errorf("%s: expect number of signers wrong, got %v, want %v", tc.testName, len(tc.msg.GetSigners()), len(tc.expectSigners))
			return
		}
		for i, signer := range tc.msg.GetSigners() {
			if types.AccountKey(signer) != tc.expectSigners[i] {
				t.Errorf("%s: expect signer wrong, got %v, want %v", tc.testName, types.AccountKey(signer), tc.expectSigners[i])
				return
			}
		}
	}
}

func TestGetConsumeAmount(t *testing.T) {
	testCases := []struct {
		testName     string
		msg          types.Msg
		expectAmount types.Coin
	}{
		{
			testName: "donateMsg",
			msg: NewDonateMsg(
				"test", types.LNO("1"),
				"author", "postID", "", memo1),
			expectAmount: types.NewCoinFromInt64(1 * types.Decimals),
		},
		{
			testName: "create post",
			msg: CreatePostMsg{
				PostID:       "test",
				Title:        "title",
				Content:      "content",
				Author:       "author",
				ParentAuthor: types.AccountKey("parentAuthor"),
				ParentPostID: "parentPostID",
				SourceAuthor: types.AccountKey("sourceAuthor"),
				SourcePostID: "sourcePostID",
				Links: []types.IDToURLMapping{
					{
						Identifier: "#1",
						URL:        "https://lino.network",
					},
				},
				RedistributionSplitRate: "0.5",
			},
			expectAmount: types.NewCoinFromInt64(0),
		},
		{
			testName: "view post",
			msg: NewViewMsg(
				"test", "author", "postID"),
			expectAmount: types.NewCoinFromInt64(0),
		},
		{
			testName: "report post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", true),
			expectAmount: types.NewCoinFromInt64(0),
		},
		{
			testName: "upvote post",
			msg: NewReportOrUpvoteMsg(
				"test", "author", "postID", false),
			expectAmount: types.NewCoinFromInt64(0),
		},
		{
			testName: "update post",
			msg: NewUpdatePostMsg(
				"author", "postID", "title", "content", []types.IDToURLMapping{}, "0"),
			expectAmount: types.NewCoinFromInt64(0),
		},
	}

	for _, tc := range testCases {
		if !tc.expectAmount.IsEqual(tc.msg.GetConsumeAmount()) {
			t.Errorf("%s: expect consume amount wrong, got %v, want %v", tc.testName, tc.msg.GetConsumeAmount(), tc.expectAmount)
			return
		}
	}
}
