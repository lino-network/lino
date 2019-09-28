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
	author := types.AccountKey("testauthor")
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
			expectedResult: ErrInvalidAuthor(),
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
			expectedResult: ErrInvalidCreatedBy(),
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
			expected: ErrInvalidAuthor(),
		},
	}
	for _, c := range testCases {
		suite.Run(c.testName, func() {
			suite.Equal(c.expected, c.msg.ValidateBasic())
		})
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
			expected: ErrInvalidAuthor(),
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
		suite.Run(c.testName, func() {
			suite.Equal(c.expected, c.msg.ValidateBasic())
		})
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
			expected: ErrInvalidUsername(),
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
			testName: "invalid app - invalid app name",
			msg:      NewDonateMsg("test", types.LNO("1"), "", "", "x", memo1),
			expected: ErrInvalidTarget(),
		},
		{
			testName: "self donate",
			msg:      NewDonateMsg("test", types.LNO("1"), "test", "post1", "app1", memo1),
			expected: ErrCannotDonateToSelf("test"),
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
		suite.Run(c.testName, func() {
			suite.Equal(c.expected, c.msg.ValidateBasic())
		})
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
		suite.Run(c.testName, func() {
			suite.Equal(c.expected, c.msg.GetConsumeAmount())
		})
	}
}

func (suite *PostMsgTestSuite) TestIDADonateMsgValidateBasic() {
	testCases := []struct {
		testName string
		msg      IDADonateMsg
		expected sdk.Error
	}{
		{
			testName: "ok1",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "12345",
				Author:   "user2",
				PostID:   "post1",
				Memo:     memo1,
				Signer:   "singer",
			},
			expected: nil,
		},
		{
			testName: "no username",
			msg: IDADonateMsg{
				Username: "",
				App:      "app1",
				Amount:   "12345",
				Author:   "user2",
				PostID:   "post1",
				Signer:   "singer",
			},
			expected: ErrInvalidUsername(),
		},
		{
			testName: "zero amount is less than lower bound",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "0",
				Author:   "user2",
				PostID:   "post1",
				Signer:   "singer",
			},
			expected: types.ErrInvalidIDAAmount(),
		},
		{
			testName: "negative coin is less than lower bound",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "-1",
				Author:   "user2",
				PostID:   "post1",
				Signer:   "singer",
			},
			expected: types.ErrInvalidIDAAmount(),
		},
		{
			testName: "invalid target - no post id",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "1",
				Author:   "user2",
				PostID:   "",
				Signer:   "singer",
			},
			expected: ErrInvalidTarget(),
		},
		{
			testName: "invalid target - no author",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "1",
				Author:   "",
				PostID:   "post1",
				Signer:   "singer",
			},
			expected: ErrInvalidTarget(),
		},
		{
			testName: "invalid target - no post id",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "1",
				Author:   "user2",
				PostID:   "",
				Signer:   "singer",
			},
			expected: ErrInvalidTarget(),
		},
		{
			testName: "invalid app - invalid app name",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "x",
				Amount:   "1",
				Author:   "user2",
				PostID:   "post1",
				Signer:   "singer",
			},
			expected: ErrInvalidApp(),
		},
		{
			testName: "self donate",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "1",
				Author:   "user1",
				PostID:   "post1",
				Signer:   "singer",
			},
			expected: ErrCannotDonateToSelf("user1"),
		},
		{
			testName: "invalid memo",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "1",
				Author:   "user2",
				PostID:   "post1",
				Memo:     invalidMemo,
				Signer:   "singer",
			},
			expected: ErrInvalidMemo(),
		},
		{
			testName: "utf8 memo is too long",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "1",
				Author:   "user2",
				PostID:   "post1",
				Memo:     tooLongOfUTF8Memo,
				Signer:   "singer",
			},
			expected: ErrInvalidMemo(),
		},
		{
			testName: "invalid signer",
			msg: IDADonateMsg{
				Username: "user1",
				App:      "app1",
				Amount:   "1",
				Author:   "user2",
				PostID:   "post1",
				Memo:     memo1,
				Signer:   "x",
			},
			expected: ErrInvalidUsername(),
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expected, tc.msg.ValidateBasic(), "%s", tc.testName)
	}
}

// func (suite *PostMsgTestSuite) TestMsgPermission() {
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
