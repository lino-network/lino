package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/types"
)

var (
	memo1       = "test memo"
	invalidMemo = "Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! "

	// len of 101
	tooLongOfUTF8Memo = "12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345"

	// len of 100
	maxLenOfUTF8Title = `12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌1234512345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345`

	// len of 101
	tooLongOfUTF8Title = `12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌123456 12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345`

	// len of 1000
	maxLenOfUTF8Content = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌123`

	// len of 1001
	tooLongOfUTF8Content = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌1234`
)

type PostMsgTestSuite struct {
	suite.Suite
}

func TestPostMsgTestSuite(t *testing.T) {
	suite.Run(t, new(PostMsgTestSuite))
}

func (suite *PostMsgTestSuite) TestCreatePostMsgValidateBasic() {
	author := types.AccountKey("TestAuthor")
	app := types.AccountKey("app")
	testCases := []struct {
		testName       string
		msg            CreatePostMsg
		expectedResult sdk.Error
	}{
		{
			testName: "normal case 1",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: author,
			},
			expectedResult: nil,
		},
		{
			testName: "normal case 2",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: app,
			},
			expectedResult: nil,
		},
		{
			testName: "utf8 title",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     maxLenOfUTF8Title,
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: author,
			},
			expectedResult: nil,
		},
		{
			testName: "utf8 content",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 100)),
				Content:   maxLenOfUTF8Content,
				Author:    author,
				CreatedBy: author,
			},
			expectedResult: nil,
		},
		{
			testName: "empty post id",
			msg: CreatePostMsg{
				PostID:    "",
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: author,
			},
			expectedResult: ErrNoPostID(),
		},
		{
			testName: "no author",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    "",
				CreatedBy: app,
			},
			expectedResult: ErrNoAuthor(),
		},
		{
			testName: "post id is too long",
			msg: CreatePostMsg{
				PostID:    string(make([]byte, types.MaximumLengthOfPostID+1)),
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: author,
			},
			expectedResult: ErrPostIDTooLong(),
		},
		{
			testName: "post title is too long",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 101)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: author,
			},
			expectedResult: ErrPostTitleExceedMaxLength(),
		},
		{
			testName: "post utf8 title is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   tooLongOfUTF8Title,
				Content: string(make([]byte, 1000)),
				Author:  author,

				CreatedBy: author,
			},
			expectedResult: ErrPostTitleExceedMaxLength(),
		},
		{
			testName: "post content is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 100)),
				Content: string(make([]byte, 1001)),
				Author:  author,

				CreatedBy: author,
			},
			expectedResult: ErrPostContentExceedMaxLength(),
		},
		{
			testName: "post utf8 content is too long",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 100)),
				Content: tooLongOfUTF8Content,
				Author:  author,
			},
			expectedResult: ErrPostContentExceedMaxLength(),
		},
		{
			testName: "no createdBy",
			msg: CreatePostMsg{
				PostID:  "TestPostID",
				Title:   string(make([]byte, 100)),
				Content: string(make([]byte, 1000)),
				Author:  author,
			},
			expectedResult: ErrNoCreatedBy(),
		},
	}
	for _, tc := range testCases {
		result := tc.msg.ValidateBasic()
		suite.Equal(tc.expectedResult, result,
			"%s: diff result, got %v, want %v", tc.testName, result, tc.expectedResult)
	}
}

func (suite *PostMsgTestSuite) TestCreatePostMsgPermission() {
	author := types.AccountKey("TestAuthor")
	app := types.AccountKey("app")
	testCases := []struct {
		testName string
		msg      CreatePostMsg
		expected types.Permission
	}{
		{
			testName: "created by author",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: author,
			},
			expected: types.TransactionPermission,
		},
		{
			testName: "created by app",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: app,
			},
			expected: types.AppOrAffiliatedPermission,
		},
		{
			testName: "created by app",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     string(make([]byte, 100)),
				Content:   string(make([]byte, 1000)),
				Author:    author,
				CreatedBy: app,
				Preauth:   true,
			},
			expected: types.AppPermission,
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expected, tc.msg.GetPermission(), "%s: diff result", tc.testName)
	}
}

func (suite *PostMsgTestSuite) TestCreatePostMsgSignBytes() {
	author := types.AccountKey("TestAuthor")
	app := types.AccountKey("app")
	testCases := []struct {
		testName string
		msg      CreatePostMsg
		expected []byte
	}{
		{
			testName: "normal",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     "title",
				Content:   "content",
				Author:    author,
				CreatedBy: app,
				Preauth:   true,
			},
			expected: []byte(`{"type":"lino/createPost","value":{"author":"TestAuthor","content":"content","created_by":"app","post_id":"TestPostID","preauth":true,"title":"title"}}`),
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expected, tc.msg.GetSignBytes(), "%s", tc.testName)
	}
}

func (suite *PostMsgTestSuite) TestCreatePostMsgSigners() {
	author := types.AccountKey("TestAuthor")
	app := types.AccountKey("app")
	testCases := []struct {
		testName string
		msg      CreatePostMsg
		expected []sdk.AccAddress
	}{
		{
			testName: "created by author",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     "title",
				Content:   "content",
				Author:    author,
				CreatedBy: author,
				Preauth:   false,
			},
			expected: []sdk.AccAddress{sdk.AccAddress(author)},
		},
		{
			testName: "created by app",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     "title",
				Content:   "content",
				Author:    author,
				CreatedBy: app,
				Preauth:   false,
			},
			expected: []sdk.AccAddress{sdk.AccAddress(app)},
		},
		{
			testName: "created by preauth",
			msg: CreatePostMsg{
				PostID:    "TestPostID",
				Title:     "title",
				Content:   "content",
				Author:    author,
				CreatedBy: app,
				Preauth:   true,
			},
			expected: []sdk.AccAddress{sdk.AccAddress(author)},
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expected, tc.msg.GetSigners(), "%s", tc.testName)
	}
}

func (suite *PostMsgTestSuite) TestUpdatePostValidateBasic() {
	author := types.AccountKey("TestAuthor")
	testCases := []struct {
		testName string
		msg      UpdatePostMsg
		expected sdk.Error
	}{
		{
			testName: "correct input",
			msg: UpdatePostMsg{
				Author:  author,
				PostID:  "TestPostID",
				Title:   "TestTitle",
				Content: "TestContent",
			},
			expected: nil,
		},
		{
			testName: "should throw error after failing post basic check",
			msg: UpdatePostMsg{
				Author:  types.AccountKey(""),
				PostID:  "TestPostID",
				Title:   "TestTitle",
				Content: "TestContent",
			},
			expected: ErrNoAuthor(),
		},
	}
	for _, c := range testCases {
		suite.Equal(c.expected, c.msg.ValidateBasic())
	}
}

func (suite *PostMsgTestSuite) TestDeletePostValidateBasic() {
	author := types.AccountKey("TestAuthor")
	testCases := []struct {
		testName string
		msg      DeletePostMsg
		expected sdk.Error
	}{
		{
			testName: "correct input",
			msg: DeletePostMsg{
				Author: author,
				PostID: "TestPostID",
			},
			expected: nil,
		},
		{
			testName: "should throw error if author is empty",
			msg: DeletePostMsg{
				Author: types.AccountKey(""),
				PostID: "TestPostID",
			},
			expected: ErrNoAuthor(),
		},
		{
			testName: "should throw error if postID is empty",
			msg: DeletePostMsg{
				Author: author,
				PostID: "",
			},
			expected: ErrNoPostID(),
		},
	}
	for _, c := range testCases {
		suite.Equal(c.expected, c.msg.ValidateBasic())
	}
}

func (suite *PostMsgTestSuite) TestDonateMsgValidateBasic() {
	testCases := []struct {
		testName string
		msg      DonateMsg
		expected sdk.Error
	}{
		{
			testName: "normal case",
			msg:      NewDonateMsg("test", types.LNO("1"), "author", "postID", "", memo1),
			expected: nil,
		},
		{
			testName: "no username",
			msg:      NewDonateMsg("", types.LNO("1"), "author", "postID", "", memo1),
			expected: ErrNoUsername(),
		},
		{
			testName: "zero coin is less than lower bound",
			msg:      NewDonateMsg("test", types.LNO("0"), "author", "postID", "", memo1),
			expected: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName: "negative coin is less than lower bound",
			msg:      NewDonateMsg("test", types.LNO("-1"), "author", "postID", "", memo1),
			expected: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName: "invalid target - no post id",
			msg:      NewDonateMsg("test", types.LNO("1"), "author", "", "", memo1),
			expected: ErrInvalidTarget(),
		},
		{
			testName: "invalid target - no author",
			msg:      NewDonateMsg("test", types.LNO("1"), "", "postID", "", memo1),
			expected: ErrInvalidTarget(),
		},
		{
			testName: "invalid target - no author and post id",
			msg:      NewDonateMsg("test", types.LNO("1"), "", "", "", memo1),
			expected: ErrInvalidTarget(),
		},
		{
			testName: "invalid memo",
			msg:      NewDonateMsg("test", types.LNO("1"), "author", "postID", "", invalidMemo),
			expected: ErrInvalidMemo(),
		},
		{
			testName: "utf8 memo is too long",
			msg:      NewDonateMsg("test", types.LNO("1"), "author", "postID", "", tooLongOfUTF8Memo),
			expected: ErrInvalidMemo(),
		},
	}
	for _, c := range testCases {
		suite.Equal(c.expected, c.msg.ValidateBasic())
	}
}

func (suite *PostMsgTestSuite) TestDonateMsgConsumeAmount() {
	testCases := []struct {
		testName string
		msg      DonateMsg
		expected types.Coin
	}{
		{
			testName: "donateMsg",
			msg: NewDonateMsg(
				"test", types.LNO("1"),
				"author", "postID", "", memo1),
			expected: types.NewCoinFromInt64(1 * types.Decimals),
		},
	}
	for _, c := range testCases {
		suite.Equal(c.expected, c.msg.GetConsumeAmount())
	}
}

// func (suite *PostMsgTestSuite) TestIDADonateMsgValidateBasic() {
// 	testCases := []struct {
// 		testName string
// 		msg IDADonateMsg
// 		expected sdk.Error
// 	}{
// 		{
// 			testName: "

// func TestUpdatePostMsg(t *testing.T) {
// 	testCases := []struct {
// 		testName       string
// 		updatePostMsg  UpdatePostMsg
// 		expectedResult sdk.Error
// 	}{
// 		{
// 			testName: "normal case 1",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", "title", "content", []types.IDToURLMapping{}),
// 			expectedResult: nil,
// 		},
// 		{
// 			testName: "normal case 2",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", "title", "content", []types.IDToURLMapping{}),
// 			expectedResult: nil,
// 		},
// 		{
// 			testName: "utf8 title",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", maxLenOfUTF8Title, "content", []types.IDToURLMapping{}),
// 			expectedResult: nil,
// 		},
// 		{
// 			testName: "utf8 content",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", "title", maxLenOfUTF8Content, []types.IDToURLMapping{}),
// 			expectedResult: nil,
// 		},
// 		{
// 			testName: "no author",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"", "postID", "title", "content", []types.IDToURLMapping{}),
// 			expectedResult: ErrNoAuthor(),
// 		},
// 		{
// 			testName: "no post id",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "", "title", "content", []types.IDToURLMapping{}),
// 			expectedResult: ErrNoPostID(),
// 		},
// 		{
// 			testName: "post tile is too long",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", string(make([]byte, 101)), "content", []types.IDToURLMapping{}),
// 			expectedResult: ErrPostTitleExceedMaxLength(),
// 		},
// 		{
// 			testName: "post utf8 tile is too long",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", tooLongOfUTF8Title, "content", []types.IDToURLMapping{}),
// 			expectedResult: ErrPostTitleExceedMaxLength(),
// 		},
// 		{
// 			testName: "post content is too long",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", string(make([]byte, 100)), string(make([]byte, 1001)),
// 				[]types.IDToURLMapping{}),
// 			expectedResult: ErrPostContentExceedMaxLength(),
// 		},
// 		{
// 			testName: "post utf8 content is too long",
// 			updatePostMsg: NewUpdatePostMsg(
// 				"author", "postID", string(make([]byte, 100)), tooLongOfUTF8Content,
// 				[]types.IDToURLMapping{}),
// 			expectedResult: ErrPostContentExceedMaxLength(),
// 		},
// 	}
// 	for _, tc := range testCases {
// 		result := tc.updatePostMsg.ValidateBasic()
// 		if !assert.Equal(t, result, tc.expectedResult) {
// 			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedResult)
// 		}
// 	}
// }

// func TestDeletePostMsg(t *testing.T) {
// 	testCases := []struct {
// 		testName    string
// 		msg         DeletePostMsg
// 		wantErrCode sdk.CodeType
// 	}{
// 		{
// 			testName: "normal case",
// 			msg: DeletePostMsg{
// 				Author: "author",
// 				PostID: "postID",
// 			},
// 			wantErrCode: sdk.CodeOK,
// 		},
// 		{
// 			testName: "empty author",
// 			msg: DeletePostMsg{
// 				Author: "",
// 				PostID: "postID",
// 			},
// 			wantErrCode: types.CodeNoAuthor,
// 		},
// 		{
// 			testName: "empty postID",
// 			msg: DeletePostMsg{
// 				Author: "author",
// 				PostID: "",
// 			},
// 			wantErrCode: types.CodeNoPostID,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		got := tc.msg.ValidateBasic()
// 		if got == nil && tc.wantErrCode != sdk.CodeOK {
// 			t.Errorf("%s: got non-OK code, got %v, want %v", tc.testName, got, tc.wantErrCode)
// 		}
// 		if got != nil {
// 			if got.Code() != tc.wantErrCode {
// 				t.Errorf("%s: diff err code, got %v, want %v", tc.testName, got, tc.wantErrCode)
// 			}
// 		}
// 	}
// }

// func TestMsgPermission(t *testing.T) {
// 	testCases := []struct {
// 		testName           string
// 		msg                types.Msg
// 		expectedPermission types.Permission
// 	}{
// 		{
// 			testName: "donateMsg",
// 			msg: NewDonateMsg(
// 				"test", types.LNO("1"),
// 				"author", "postID", "", memo1),
// 			expectedPermission: types.PreAuthorizationPermission,
// 		},
// 		{
// 			testName: "create post",
// 			msg: CreatePostMsg{
// 				PostID:       "test",
// 				Title:        "title",
// 				Content:      "content",
// 				Author:       "author",
// 				ParentAuthor: types.AccountKey("parentAuthor"),
// 				ParentPostID: "parentPostID",
// 				SourceAuthor: types.AccountKey("sourceAuthor"),
// 				SourcePostID: "sourcePostID",
// 				Links: []types.IDToURLMapping{
// 					{
// 						Identifier: "#1",
// 						URL:        "https://lino.network",
// 					},
// 				},
// 				RedistributionSplitRate: "0.5",
// 			},
// 			expectedPermission: types.AppPermission,
// 		},
// 		{
// 			testName: "view post",
// 			msg: NewViewMsg(
// 				"test", "author", "postID"),
// 			expectedPermission: types.AppPermission,
// 		},
// 		{
// 			testName: "report post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", true),
// 			expectedPermission: types.AppPermission,
// 		},
// 		{
// 			testName: "upvote post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", false),
// 			expectedPermission: types.AppPermission,
// 		},
// 		{
// 			testName: "update post",
// 			msg: NewUpdatePostMsg(
// 				"author", "postID", "title", "content", []types.IDToURLMapping{}),
// 			expectedPermission: types.AppPermission,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		permission := tc.msg.GetPermission()
// 		if tc.expectedPermission != permission {
// 			t.Errorf("%s: diff permission, got %v, want %v", tc.testName, tc.expectedPermission, permission)
// 		}
// 	}
// }

// func TestGetSignBytes(t *testing.T) {
// 	testCases := []struct {
// 		testName string
// 		msg      types.Msg
// 	}{
// 		{
// 			testName: "donateMsg",
// 			msg: NewDonateMsg(
// 				"test", types.LNO("1"),
// 				"author", "postID", "", memo1),
// 		},
// 		{
// 			testName: "create post",
// 			msg: CreatePostMsg{
// 				PostID:       "test",
// 				Title:        "title",
// 				Content:      "content",
// 				Author:       "author",
// 				ParentAuthor: types.AccountKey("parentAuthor"),
// 				ParentPostID: "parentPostID",
// 				SourceAuthor: types.AccountKey("sourceAuthor"),
// 				SourcePostID: "sourcePostID",
// 				Links: []types.IDToURLMapping{
// 					{
// 						Identifier: "#1",
// 						URL:        "https://lino.network",
// 					},
// 				},
// 				RedistributionSplitRate: "0.5",
// 			},
// 		},
// 		{
// 			testName: "view post",
// 			msg: NewViewMsg(
// 				"test", "author", "postID"),
// 		},
// 		{
// 			testName: "report post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", true),
// 		},
// 		{
// 			testName: "upvote post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", false),
// 		},
// 		{
// 			testName: "update post",
// 			msg: NewUpdatePostMsg(
// 				"author", "postID", "title", "content", []types.IDToURLMapping{}),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, tc.testName)
// 	}
// }

// func TestGetSigners(t *testing.T) {
// 	testCases := []struct {
// 		testName      string
// 		msg           types.Msg
// 		expectSigners []types.AccountKey
// 	}{
// 		{
// 			testName: "donateMsg",
// 			msg: NewDonateMsg(
// 				"test", types.LNO("1"),
// 				"author", "postID", "", memo1),
// 			expectSigners: []types.AccountKey{"test"},
// 		},
// 		{
// 			testName: "create post",
// 			msg: CreatePostMsg{
// 				PostID:       "test",
// 				Title:        "title",
// 				Content:      "content",
// 				Author:       "author",
// 				ParentAuthor: types.AccountKey("parentAuthor"),
// 				ParentPostID: "parentPostID",
// 				SourceAuthor: types.AccountKey("sourceAuthor"),
// 				SourcePostID: "sourcePostID",
// 				Links: []types.IDToURLMapping{
// 					{
// 						Identifier: "#1",
// 						URL:        "https://lino.network",
// 					},
// 				},
// 				RedistributionSplitRate: "0.5",
// 			},
// 			expectSigners: []types.AccountKey{"author"},
// 		},
// 		{
// 			testName: "view post",
// 			msg: NewViewMsg(
// 				"test", "author", "postID"),
// 			expectSigners: []types.AccountKey{"test"},
// 		},
// 		{
// 			testName: "report post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", true),
// 			expectSigners: []types.AccountKey{"test"},
// 		},
// 		{
// 			testName: "upvote post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", false),
// 			expectSigners: []types.AccountKey{"test"},
// 		},
// 		{
// 			testName: "update post",
// 			msg: NewUpdatePostMsg(
// 				"author", "postID", "title", "content", []types.IDToURLMapping{}),
// 			expectSigners: []types.AccountKey{"author"},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		if len(tc.msg.GetSigners()) != len(tc.expectSigners) {
// 			t.Errorf("%s: expect number of signers wrong, got %v, want %v", tc.testName, len(tc.msg.GetSigners()), len(tc.expectSigners))
// 			return
// 		}
// 		for i, signer := range tc.msg.GetSigners() {
// 			if types.AccountKey(signer) != tc.expectSigners[i] {
// 				t.Errorf("%s: expect signer wrong, got %v, want %v", tc.testName, types.AccountKey(signer), tc.expectSigners[i])
// 				return
// 			}
// 		}
// 	}
// }

// func TestGetConsumeAmount(t *testing.T) {
// 	testCases := []struct {
// 		testName     string
// 		msg          types.Msg
// 		expectAmount types.Coin
// 	}{
// 		{
// 			testName: "donateMsg",
// 			msg: NewDonateMsg(
// 				"test", types.LNO("1"),
// 				"author", "postID", "", memo1),
// 			expectAmount: types.NewCoinFromInt64(1 * types.Decimals),
// 		},
// 		{
// 			testName: "create post",
// 			msg: CreatePostMsg{
// 				PostID:       "test",
// 				Title:        "title",
// 				Content:      "content",
// 				Author:       "author",
// 				ParentAuthor: types.AccountKey("parentAuthor"),
// 				ParentPostID: "parentPostID",
// 				SourceAuthor: types.AccountKey("sourceAuthor"),
// 				SourcePostID: "sourcePostID",
// 				Links: []types.IDToURLMapping{
// 					{
// 						Identifier: "#1",
// 						URL:        "https://lino.network",
// 					},
// 				},
// 				RedistributionSplitRate: "0.5",
// 			},
// 			expectAmount: types.NewCoinFromInt64(0),
// 		},
// 		{
// 			testName: "view post",
// 			msg: NewViewMsg(
// 				"test", "author", "postID"),
// 			expectAmount: types.NewCoinFromInt64(0),
// 		},
// 		{
// 			testName: "report post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", true),
// 			expectAmount: types.NewCoinFromInt64(0),
// 		},
// 		{
// 			testName: "upvote post",
// 			msg: NewReportOrUpvoteMsg(
// 				"test", "author", "postID", false),
// 			expectAmount: types.NewCoinFromInt64(0),
// 		},
// 		{
// 			testName: "update post",
// 			msg: NewUpdatePostMsg(
// 				"author", "postID", "title", "content", []types.IDToURLMapping{}),
// 			expectAmount: types.NewCoinFromInt64(0),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		if !tc.expectAmount.IsEqual(tc.msg.GetConsumeAmount()) {
// 			t.Errorf("%s: expect consume amount wrong, got %v, want %v", tc.testName, tc.msg.GetConsumeAmount(), tc.expectAmount)
// 			return
// 		}
// 	}
// }
