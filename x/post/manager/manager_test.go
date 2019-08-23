package manager

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/testsuites"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/model"
	types "github.com/lino-network/lino/x/post/types"
	"github.com/stretchr/testify/mock"

	acc "github.com/lino-network/lino/x/account/mocks"
	dev "github.com/lino-network/lino/x/developer/mocks"
	global "github.com/lino-network/lino/x/global/mocks"
	price "github.com/lino-network/lino/x/price/mocks"
	rep "github.com/lino-network/lino/x/reputation/mocks"
)

var dummyErr = linotypes.NewError(linotypes.CodeTestDummyError, "")

type PostManagerTestSuite struct {
	testsuites.CtxTestSuite
	pm PostManager
	// deps
	am     *acc.AccountKeeper
	dev    *dev.DeveloperKeeper
	global *global.GlobalKeeper
	price  *price.PriceKeeper
	rep    *rep.ReputationKeeper
	// mock data
	user1  linotypes.AccountKey
	user2  linotypes.AccountKey
	unreg1 linotypes.AccountKey
	app1   linotypes.AccountKey
	app2   linotypes.AccountKey
	app3   linotypes.AccountKey
	rate   sdk.Dec
}

func TestPostManagerTestSuite(t *testing.T) {
	suite.Run(t, new(PostManagerTestSuite))
}

func (suite *PostManagerTestSuite) SetupTest() {
	testPostKey := sdk.NewKVStoreKey("post")
	suite.SetupCtx(0, time.Unix(0, 0), testPostKey)
	suite.am = &acc.AccountKeeper{}
	suite.dev = &dev.DeveloperKeeper{}
	suite.global = &global.GlobalKeeper{}
	suite.price = &price.PriceKeeper{}
	suite.rep = &rep.ReputationKeeper{}
	suite.pm = NewPostManager(testPostKey, suite.am, suite.global, suite.dev, suite.rep, suite.price)

	// background
	suite.user1 = linotypes.AccountKey("user1")
	suite.user2 = linotypes.AccountKey("user2")
	suite.unreg1 = linotypes.AccountKey("user3")
	suite.app1 = linotypes.AccountKey("app1")
	suite.app2 = linotypes.AccountKey("app2")
	suite.app3 = linotypes.AccountKey("app3")

	// reg accounts
	for _, v := range []linotypes.AccountKey{suite.user1, suite.user2, suite.app1, suite.app2, suite.app3} {
		suite.am.On("DoesAccountExist", mock.Anything, v).Return(true).Maybe()
	}
	// unreg accounts
	for _, v := range []linotypes.AccountKey{suite.unreg1} {
		suite.am.On("DoesAccountExist", mock.Anything, v).Return(false).Maybe()
	}

	// reg dev
	for _, v := range []linotypes.AccountKey{suite.app1, suite.app2, suite.app3} {
		suite.dev.On("DoesDeveloperExist", mock.Anything, v).Return(true).Maybe()
	}
	// unreg devs
	for _, v := range []linotypes.AccountKey{suite.unreg1, suite.user1, suite.user2} {
		suite.dev.On("DoesDeveloperExist", mock.Anything, v).Return(false).Maybe()
	}

	rate, err := sdk.NewDecFromStr("0.099")
	suite.Require().Nil(err)
	suite.global.On("GetConsumptionFrictionRate", mock.Anything).Return(rate, nil).Maybe()
	suite.rate = rate
	// // app1, app2 has issued IDA
	// suite.dev.On("GetIDAPrice", suite.Ctx, suite.app1).Return(linotypes.NewMiniDollar(10),nil)
	// suite.dev.On("GetIDAPrice", suite.Ctx, suite.app2).Return(linotypes.NewMiniDollar(7),nil)
}

func (suite *PostManagerTestSuite) TestCreatePost() {
	user1 := suite.user1
	user2 := suite.user2
	user3 := suite.unreg1
	app1 := suite.app1
	app2 := suite.app2
	testCases := []struct {
		testName     string
		postID       string
		title        string
		content      string
		author       linotypes.AccountKey
		createdby    linotypes.AccountKey
		expectResult sdk.Error
	}{
		{
			testName:     "user does not exists",
			postID:       "postID",
			author:       user3,
			title:        "title1",
			content:      "content1",
			createdby:    user3,
			expectResult: types.ErrAccountNotFound(user3),
		},
		{
			testName:     "createdBy does not exists",
			postID:       "postID",
			title:        "title2",
			content:      "content2",
			author:       user2,
			createdby:    user3,
			expectResult: types.ErrAccountNotFound(user3),
		},
		{
			testName:     "createdBy is not an app",
			postID:       "postID",
			content:      "content3",
			title:        "title3",
			author:       user1,
			createdby:    user2,
			expectResult: types.ErrDeveloperNotFound(user2),
		},
		{
			testName:     "creates (postID, user1) successfully, by author",
			postID:       "postID",
			content:      "content4",
			title:        "title4",
			author:       user1,
			createdby:    user1,
			expectResult: nil,
		},
		{
			testName:     "creates (postID, user2) successfully, by app",
			postID:       "postID",
			content:      "content5",
			title:        "title5",
			author:       user2,
			createdby:    app1,
			expectResult: nil,
		},
		{
			testName:     "(postID, user1) already exists",
			postID:       "postID",
			content:      "content6",
			title:        "title6",
			author:       user1,
			createdby:    app1,
			expectResult: types.ErrPostAlreadyExist(linotypes.GetPermlink(user1, "postID")),
		},
		{
			testName:     "(postID, user2) already exists case 1",
			postID:       "postID",
			content:      "content7",
			title:        "title7",
			author:       user2,
			createdby:    user2,
			expectResult: types.ErrPostAlreadyExist(linotypes.GetPermlink(user2, "postID")),
		},
		{
			testName:     "creates (postID2, user2) successfully",
			postID:       "postID2",
			content:      "content8",
			title:        "title8",
			author:       user2,
			createdby:    app1,
			expectResult: nil,
		},
		{
			testName:     "creates (postID2, user1) successfully",
			postID:       "postID2",
			title:        "title9",
			content:      "content9",
			author:       user1,
			createdby:    app2,
			expectResult: nil,
		},
	}

	for _, tc := range testCases {
		// test valid postInfo
		msg := types.CreatePostMsg{
			PostID:    tc.postID,
			Title:     tc.title,
			Content:   tc.content,
			Author:    tc.author,
			CreatedBy: tc.createdby,
		}
		err := suite.pm.CreatePost(
			suite.Ctx, msg.Author, msg.PostID, msg.CreatedBy, msg.Content, msg.Title)
		suite.Equal(tc.expectResult, err, "%s", tc.testName)
		if tc.expectResult == nil {
			post, err := suite.pm.postStorage.GetPost(
				suite.Ctx, linotypes.GetPermlink(tc.author, tc.postID))
			suite.Nil(err)
			suite.Equal(&model.Post{
				PostID:    tc.postID,
				Title:     tc.title,
				Content:   tc.content,
				Author:    tc.author,
				CreatedBy: tc.createdby,
				CreatedAt: suite.Ctx.BlockHeader().Time.Unix(),
				UpdatedAt: suite.Ctx.BlockHeader().Time.Unix(),
			}, post, "%s", tc.testName)
		}
	}
}

func (suite *PostManagerTestSuite) TestUpdatePost() {
	user1 := suite.user1
	user2 := suite.user2
	app1 := suite.app1
	postID := "post1"
	err := suite.pm.CreatePost(suite.Ctx, user1, postID, app1, "content", "title")
	suite.Require().Nil(err)
	baseTime := suite.Ctx.BlockHeader().Time.Unix()

	testCases := []struct {
		testName   string
		author     linotypes.AccountKey
		postID     string
		title      string
		content    string
		expectErr  sdk.Error
		updateTime int64
	}{
		{
			testName:   "update: author update",
			author:     user1,
			postID:     postID,
			title:      "update to this title",
			content:    "update to this content",
			expectErr:  nil,
			updateTime: baseTime + 10,
		},
		{
			testName:   "update with invalid post id",
			author:     user1,
			postID:     "invalid",
			expectErr:  types.ErrPostNotFound(linotypes.GetPermlink(user1, "invalid")),
			updateTime: baseTime + 100,
		},
		{
			testName:   "update with invalid author",
			author:     user2,
			postID:     postID,
			expectErr:  types.ErrPostNotFound(linotypes.GetPermlink(user2, postID)),
			updateTime: baseTime + 1000,
		},
		{
			testName:   "update with account that does not exist",
			author:     suite.unreg1,
			postID:     postID,
			expectErr:  types.ErrAccountNotFound(suite.unreg1),
			updateTime: baseTime + 10000,
		},
	}

	for _, tc := range testCases {
		suite.NextBlock(time.Unix(tc.updateTime, 0))
		err := suite.pm.UpdatePost(suite.Ctx, tc.author, tc.postID, tc.title, tc.content)
		suite.Equal(tc.expectErr, err)
		if tc.expectErr == nil {
			post, err := suite.pm.postStorage.GetPost(
				suite.Ctx, linotypes.GetPermlink(tc.author, tc.postID))
			suite.Nil(err)
			suite.Equal(&model.Post{
				PostID:    tc.postID,
				Title:     tc.title,
				Content:   tc.content,
				Author:    tc.author,
				CreatedBy: app1,
				CreatedAt: baseTime,
				UpdatedAt: tc.updateTime,
			}, post, "%s", tc.testName)
		}
	}
}

func (suite *PostManagerTestSuite) TestDeletePost() {
	user1 := suite.user1
	app1 := suite.app1
	postID := "post1"
	err := suite.pm.CreatePost(suite.Ctx, user1, postID, app1, "content", "title")
	suite.Require().Nil(err)

	testCases := []struct {
		testName  string
		author    linotypes.AccountKey
		postID    string
		expectErr sdk.Error
	}{
		{
			testName:  "delete successfully",
			author:    user1,
			postID:    postID,
			expectErr: nil,
		},
		{
			testName:  "delete post not exist",
			author:    user1,
			postID:    postID,
			expectErr: types.ErrPostNotFound(linotypes.GetPermlink(user1, postID)),
		},
	}

	for _, tc := range testCases {
		err := suite.pm.DeletePost(suite.Ctx, linotypes.GetPermlink(tc.author, tc.postID))
		suite.Equal(tc.expectErr, err)
		if tc.expectErr == nil {
			suite.False(suite.pm.DoesPostExist(
				suite.Ctx, linotypes.GetPermlink(tc.author, tc.postID)))
		}
	}
}

func (suite *PostManagerTestSuite) TestLinoDonateInvalid() {
	user2 := suite.user2
	user1 := suite.user1
	app1 := suite.app1
	postID := "post1"
	err := suite.pm.CreatePost(suite.Ctx, user1, postID, app1, "content", "title")
	suite.Require().Nil(err)

	testCases := []struct {
		testName  string
		from      linotypes.AccountKey
		amount    linotypes.Coin
		author    linotypes.AccountKey
		postID    string
		app       linotypes.AccountKey
		expectErr sdk.Error
	}{
		{
			testName:  "user does not exists",
			from:      suite.unreg1,
			amount:    linotypes.NewCoinFromInt64(23),
			author:    user1,
			postID:    postID,
			app:       app1,
			expectErr: types.ErrAccountNotFound(suite.unreg1),
		},
		{
			testName:  "author does not exists",
			from:      user1,
			amount:    linotypes.NewCoinFromInt64(23),
			author:    suite.unreg1,
			postID:    postID,
			app:       app1,
			expectErr: types.ErrAccountNotFound(suite.unreg1),
		},
		{
			testName:  "post does not exists",
			from:      user2,
			amount:    linotypes.NewCoinFromInt64(23),
			author:    user1,
			postID:    "fakepost",
			app:       app1,
			expectErr: types.ErrPostNotFound(linotypes.GetPermlink(user1, "fakepost")),
		},
		{
			testName:  "self donation",
			from:      user1,
			amount:    linotypes.NewCoinFromInt64(23),
			author:    user1,
			postID:    postID,
			app:       app1,
			expectErr: types.ErrCannotDonateToSelf(user1),
		},
		{
			testName:  "app not found",
			from:      user2,
			amount:    linotypes.NewCoinFromInt64(23),
			author:    user1,
			postID:    postID,
			app:       user2,
			expectErr: types.ErrDeveloperNotFound(user2),
		},
		{
			testName:  "negative amount",
			from:      user2,
			amount:    linotypes.NewCoinFromInt64(-23),
			author:    user1,
			postID:    postID,
			app:       app1,
			expectErr: types.ErrInvalidDonationAmount(linotypes.NewCoinFromInt64(-23)),
		},
		{
			testName:  "zero amount",
			from:      user2,
			amount:    linotypes.NewCoinFromInt64(0),
			author:    user1,
			postID:    postID,
			app:       app1,
			expectErr: types.ErrInvalidDonationAmount(linotypes.NewCoinFromInt64(0)),
		},
	}

	for _, tc := range testCases {
		err := suite.pm.LinoDonate(suite.Ctx, tc.from, tc.amount, tc.author, tc.postID, tc.app)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		if err != nil {
			continue
		}

	}
}

func (suite *PostManagerTestSuite) TestLinoDonateOK() {
	from := suite.user2
	author := suite.user1
	app := suite.app1
	postID := "post1"
	amount := linotypes.NewCoinFromInt64(100000)
	tax := linotypes.DecToCoin(amount.ToDec().Mul(suite.rate))
	income := amount.Minus(tax)
	dollar := linotypes.NewMiniDollar(1000)
	dp := linotypes.NewMiniDollar(33)
	suite.price.On("CoinToMiniDollar", amount).Return(dollar)
	err := suite.pm.CreatePost(suite.Ctx, author, postID, app, "content", "title")
	suite.Require().Nil(err)

	suite.rep.On("DonateAt",
		mock.Anything, from, linotypes.GetPermlink(author, postID), dollar).Return(
		dp, nil).Once()
	suite.global.On("AddFrictionAndRegisterContentRewardEvent",
		mock.Anything,
		RewardEvent{
			PostAuthor: author,
			PostID:     postID,
			Consumer:   from,
			Evaluate:   dp,
			FromApp:    app,
		},
		tax,
		dp,
	).Return(nil).Once()
	suite.am.On("MinusSavingCoinWithFullCoinDay",
		mock.Anything, from, amount, author, "", linotypes.DonationOut).Return(
		linotypes.NewCoinFromInt64(0), nil).Once()
	suite.am.On("AddSavingCoin",
		mock.Anything, author, income, from, "", linotypes.DonationIn,
	).Return(nil).Once()
	err = suite.pm.LinoDonate(suite.Ctx, from, amount, author, postID, app)
	suite.Nil(err)
	suite.price.AssertExpectations(suite.T())
	suite.rep.AssertExpectations(suite.T())
	suite.global.AssertExpectations(suite.T())
	suite.am.AssertExpectations(suite.T())
}

// TODO(yumin): need to test path that external module returns error for 100% code coverage.
func (suite *PostManagerTestSuite) TestLinoDonateExternalFailure() {}

// TODO(yumin): add tests when it's ready.
func (suite *PostManagerTestSuite) TestIDADonateValidation() {
}

func (suite *PostManagerTestSuite) TestIDADonateOK() {
}
