package post

// XXX(yumin): add integration tests

// import (
// 	"testing"
// 	"time"

// 	"github.com/lino-network/lino/param"
// 	"github.com/lino-network/lino/types"
// 	dev "github.com/lino-network/lino/x/developer"
// 	"github.com/lino-network/lino/x/global"
// 	"github.com/lino-network/lino/x/post/model"
// 	rep "github.com/lino-network/lino/x/reputation"
// 	vote "github.com/lino-network/lino/x/vote"
// 	"github.com/stretchr/testify/assert"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	acc "github.com/lino-network/lino/x/account"
// 	accmodel "github.com/lino-network/lino/x/account/model"
// 	abci "github.com/tendermint/tendermint/abci/types"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"
// )

// func  TestHandlerCreatePost(t *testing.T) {
// 	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
// 	handler := NewHandler(pm, am, &gm, dm, rm)
// 	postParam, _ := ph.GetPostParam(ctx)

// 	user := createTestAccount(t, ctx, am, "user1")

// 	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(postParam.PostIntervalSec, 0)})
// 	// test valid post
// 	msg := CreatePostMsg{
// 		PostID:                  "TestPostID",
// 		Title:                   string(make([]byte, 50)),
// 		Content:                 string(make([]byte, 1000)),
// 		Author:                  user,
// 		ParentAuthor:            "",
// 		ParentPostID:            "",
// 		SourceAuthor:            "",
// 		SourcePostID:            "",
// 		Links:                   nil,
// 		RedistributionSplitRate: "0",
// 	}
// 	result := handler(ctx, msg)
// 	assert.Equal(t, result, sdk.Result{})
// 	assert.True(t, pm.DoesPostExist(ctx, types.GetPermlink(msg.Author, msg.PostID)))

// 	// test invlaid author
// 	msg.Author = types.AccountKey("invalid")
// 	result = handler(ctx, msg)
// 	assert.Equal(t, result, ErrAccountNotFound(msg.Author).Result())

// 	// test duplicate post
// 	msg.Author = user
// 	result = handler(ctx, msg)
// 	assert.Equal(t, result, ErrPostAlreadyExist(types.GetPermlink(user, msg.PostID)).Result())

// 	// test post too often
// 	msg.PostID = "Post too often"
// 	result = handler(ctx, msg)
// 	assert.Equal(t, result, ErrPostTooOften(msg.Author).Result())
// }

// func TestHandlerUpdatePost(t *testing.T) {
// 	ctx, am, _, pm, gm, dm, _, rm := setupTest(t, 1)
// 	handler := NewHandler(pm, am, &gm, dm, rm)

// 	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
// 	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
// 	user2 := createTestAccount(t, ctx, am, "user2")
// 	err := pm.DeletePost(ctx, types.GetPermlink(user1, postID1))
// 	assert.Nil(t, err)

// 	testCases := map[string]struct {
// 		msg        UpdatePostMsg
// 		wantResult sdk.Result
// 	}{
// 		"normal update": {
// 			msg:        NewUpdatePostMsg(string(user), postID, "update title", "update content", []types.IDToURLMapping(nil)),
// 			wantResult: sdk.Result{},
// 		},
// 		"update author doesn't exist": {
// 			msg:        NewUpdatePostMsg("invalid", postID, "update title", "update content", []types.IDToURLMapping(nil)),
// 			wantResult: ErrAccountNotFound("invalid").Result(),
// 		},
// 		"update post doesn't exist - invalid post ID": {
// 			msg:        NewUpdatePostMsg(string(user), "invalid", "update title", "update content", []types.IDToURLMapping(nil)),
// 			wantResult: ErrPostNotFound(types.GetPermlink(user, "invalid")).Result(),
// 		},
// 		"update post doesn't exist - invalid author": {
// 			msg:        NewUpdatePostMsg(string(user2), postID, "update title", "update content", []types.IDToURLMapping(nil)),
// 			wantResult: ErrPostNotFound(types.GetPermlink(user2, postID)).Result(),
// 		},
// 		"update deleted post": {
// 			msg:        NewUpdatePostMsg(string(user1), postID1, "update title", "update content", []types.IDToURLMapping(nil)),
// 			wantResult: ErrUpdatePostIsDeleted(types.GetPermlink(user1, postID1)).Result(),
// 		},
// 	}
// 	for testName, tc := range testCases {
// 		result := handler(ctx, tc.msg)
// 		if !assert.Equal(t, tc.wantResult, result) {
// 			t.Errorf("%s: diff result, got %v, want %v", testName, result, tc.wantResult)
// 		}
// 		if !tc.wantResult.IsOK() {
// 			continue
// 		}

// 		postInfo := model.PostInfo{
// 			PostID:       tc.msg.PostID,
// 			Title:        tc.msg.Title,
// 			Content:      tc.msg.Content,
// 			Author:       tc.msg.Author,
// 			SourceAuthor: "",
// 			SourcePostID: "",
// 			Links:        tc.msg.Links,
// 		}

// 		postMeta := model.PostMeta{
// 			CreatedAt:               ctx.BlockHeader().Time.Unix(),
// 			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
// 			LastActivityAt:          ctx.BlockHeader().Time.Unix(),
// 			AllowReplies:            true,
// 			IsDeleted:               false,
// 			TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
// 			TotalReward:             types.NewCoinFromInt64(0),
// 			TotalReportCoinDay:      types.NewCoinFromInt64(0),
// 			RedistributionSplitRate: sdk.ZeroDec(),
// 		}
// 		checkPostKVStore(t, ctx,
// 			types.GetPermlink(tc.msg.Author, tc.msg.PostID), postInfo, postMeta)
// 	}
// }

// func TestHandlerDeletePost(t *testing.T) {
// 	ctx, am, _, pm, gm, dm, _, rm := setupTest(t, 1)
// 	handler := NewHandler(pm, am, &gm, dm, rm)

// 	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
// 	user1 := createTestAccount(t, ctx, am, "user1")

// 	testCases := map[string]struct {
// 		msg        DeletePostMsg
// 		wantResult sdk.Result
// 	}{
// 		"normal delete": {
// 			msg: DeletePostMsg{
// 				Author: user,
// 				PostID: postID,
// 			},
// 			wantResult: sdk.Result{},
// 		},
// 		"author doesn't exist": {
// 			msg: DeletePostMsg{
// 				Author: types.AccountKey("invalid"),
// 				PostID: postID,
// 			},
// 			wantResult: ErrAccountNotFound("invalid").Result(),
// 		},
// 		"post doesn't exist - invalid author": {
// 			msg: DeletePostMsg{
// 				Author: user1,
// 				PostID: "postID",
// 			},
// 			wantResult: ErrPostNotFound(types.GetPermlink(user1, postID)).Result(),
// 		},
// 		"post doesn't exist - invalid postID": {
// 			msg: DeletePostMsg{
// 				Author: user,
// 				PostID: "invalid",
// 			},
// 			wantResult: ErrPostNotFound(types.GetPermlink(user, "invalid")).Result(),
// 		},
// 	}
// 	for testName, tc := range testCases {
// 		result := handler(ctx, tc.msg)
// 		if !assert.Equal(t, tc.wantResult, result) {
// 			t.Errorf("%s: diff result, got %v, want %v", testName, result, tc.wantResult)
// 		}
// 	}
// }

// func TestHandlerCreateComment(t *testing.T) {
// 	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
// 	handler := NewHandler(pm, am, &gm, dm, rm)
// 	postParam, err := ph.GetPostParam(ctx)
// 	assert.Nil(t, err)

// 	baseTime := time.Now()
// 	baseTime1 := baseTime.Add(time.Duration(postParam.PostIntervalSec) * time.Second)
// 	baseTime2 := baseTime1.Add(time.Duration(postParam.PostIntervalSec) * time.Second)
// 	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime})
// 	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

// 	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime1})
// 	// test comment
// 	msg := CreatePostMsg{
// 		PostID:                  "comment",
// 		Title:                   string(make([]byte, 50)),
// 		Content:                 string(make([]byte, 1000)),
// 		Author:                  user,
// 		ParentAuthor:            user,
// 		ParentPostID:            postID,
// 		SourceAuthor:            "",
// 		SourcePostID:            "",
// 		Links:                   nil,
// 		RedistributionSplitRate: "0",
// 	}
// 	result := handler(ctx, msg)
// 	assert.Equal(t, result, sdk.Result{})

// 	// after handler check KVStore
// 	postInfo := model.PostInfo{
// 		PostID:       msg.PostID,
// 		Title:        msg.Title,
// 		Content:      msg.Content,
// 		Author:       msg.Author,
// 		ParentAuthor: msg.ParentAuthor,
// 		ParentPostID: msg.ParentPostID,
// 		SourceAuthor: msg.SourceAuthor,
// 		SourcePostID: msg.SourcePostID,
// 		Links:        msg.Links,
// 	}

// 	postMeta := model.PostMeta{
// 		CreatedAt:               baseTime1.Unix(),
// 		LastUpdatedAt:           baseTime1.Unix(),
// 		LastActivityAt:          baseTime1.Unix(),
// 		AllowReplies:            true,
// 		TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
// 		TotalReward:             types.NewCoinFromInt64(0),
// 		TotalReportCoinDay:      types.NewCoinFromInt64(0),
// 		RedistributionSplitRate: sdk.ZeroDec(),
// 	}

// 	checkPostKVStore(t, ctx, types.GetPermlink(user, "comment"), postInfo, postMeta)

// 	// check parent
// 	postInfo.PostID = postID
// 	postInfo.ParentAuthor = ""
// 	postInfo.ParentPostID = ""
// 	postMeta.CreatedAt = baseTime.Unix()
// 	postMeta.LastUpdatedAt = baseTime.Unix()
// 	checkPostKVStore(t, ctx, types.GetPermlink(user, postID), postInfo, postMeta)

// 	// test post too often
// 	msg.PostID = "post too often"

// 	result = handler(ctx, msg)
// 	assert.Equal(t, result, ErrPostTooOften(user).Result())

// 	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime2})
// 	// test invalid parent
// 	msg.PostID = "invalid post"
// 	msg.ParentAuthor = user
// 	msg.ParentPostID = "invalid parent"

// 	result = handler(ctx, msg)
// 	assert.Equal(t, result, ErrPostNotFound(types.GetPermlink(user, msg.ParentPostID)).Result())

// 	// test duplicate comment
// 	msg.Author = user
// 	msg.PostID = "comment"
// 	msg.ParentAuthor = user
// 	msg.ParentPostID = "TestPostID"

// 	result = handler(ctx, msg)
// 	assert.Equal(t, result, ErrPostAlreadyExist(types.GetPermlink(msg.Author, msg.PostID)).Result())

// 	// test cycle comment
// 	msg.Author = user
// 	msg.PostID = "newComment"
// 	msg.ParentAuthor = user
// 	msg.ParentPostID = "newComment"

// 	result = handler(ctx, msg)
// 	assert.Equal(t, result, ErrPostNotFound(types.GetPermlink(user, msg.PostID)).Result())
// }

// func TestHandlerPostDonate(t *testing.T) {
// 	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
// 	handler := NewHandler(pm, am, &gm, dm, rm)

// 	accParam, err := ph.GetAccountParam(ctx)
// 	assert.Nil(t, err)

// 	author, postID := createTestPost(t, ctx, "author", "postID", am, pm, "0")
// 	author1, deletedPostID := createTestPost(t, ctx, "author1", "delete", am, pm, "0")

// 	pm.DeletePost(ctx, types.GetPermlink(author1, deletedPostID))

// 	userWithSufficientSaving := createTestAccount(t, ctx, am, "userWithSufficientSaving")
// 	err = am.AddSavingCoin(
// 		ctx, userWithSufficientSaving, types.NewCoinFromInt64(100*types.Decimals),
// 		referrer, "", types.TransferIn)
// 	assert.Nil(t, err)

// 	secondUserWithSufficientSaving := createTestAccount(t, ctx, am, "secondUserWithSufficientSaving")
// 	err = am.AddSavingCoin(
// 		ctx, secondUserWithSufficientSaving, types.NewCoinFromInt64(100*types.Decimals),
// 		referrer, "", types.TransferIn)
// 	assert.Nil(t, err)

// 	micropaymentUser := createTestAccount(t, ctx, am, "micropaymentUser")
// 	err = am.AddSavingCoin(
// 		ctx, micropaymentUser, types.NewCoinFromInt64(1*types.Decimals),
// 		referrer, "", types.TransferIn)
// 	assert.Nil(t, err)
// 	testCases := []struct {
// 		testName            string
// 		donateUser          types.AccountKey
// 		amount              types.LNO
// 		toAuthor            types.AccountKey
// 		toPostID            string
// 		expectErr           sdk.Result
// 		expectPostMeta      model.PostMeta
// 		expectDonatorSaving types.Coin
// 		expectAuthorSaving  types.Coin
// 		//https://github.com/lino-network/lino/issues/154
// 		expectRegisteredEvent             RewardEvent
// 		expectDonateTimesFromUserToAuthor int64
// 		expectCumulativeConsumption       types.Coin
// 		expectAuthorReward                accmodel.Reward
// 	}{
// 		{
// 			testName:   "donate from sufficient saving",
// 			donateUser: userWithSufficientSaving,
// 			amount:     types.LNO("100"),
// 			toAuthor:   author,
// 			toPostID:   postID,
// 			expectErr:  sdk.Result{},
// 			expectPostMeta: model.PostMeta{
// 				CreatedAt:               ctx.BlockHeader().Time.Unix(),
// 				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
// 				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
// 				AllowReplies:            true,
// 				TotalDonateCount:        1,
// 				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReportCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReward:             types.NewCoinFromInt64(95 * types.Decimals),
// 				RedistributionSplitRate: sdk.ZeroDec(),
// 			},
// 			expectDonatorSaving: accParam.RegisterFee,
// 			expectAuthorSaving: accParam.RegisterFee.Plus(
// 				types.NewCoinFromInt64(95 * types.Decimals)),
// 			expectRegisteredEvent: RewardEvent{
// 				PostAuthor: author,
// 				PostID:     postID,
// 				Consumer:   userWithSufficientSaving,
// 				Evaluate:   types.NewCoinFromInt64(1 * types.Decimals), // only 1, reputation
// 				Original:   types.NewCoinFromInt64(100 * types.Decimals),
// 				Friction:   types.NewCoinFromInt64(5 * types.Decimals),
// 				FromApp:    "",
// 			},
// 			expectDonateTimesFromUserToAuthor: 1,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(100 * types.Decimals),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(95 * types.Decimals),
// 				OriginalIncome: types.NewCoinFromInt64(95 * types.Decimals),
// 			},
// 		},
// 		{
// 			testName:            "donate from insufficient saving",
// 			donateUser:          userWithSufficientSaving,
// 			amount:              types.LNO("100"),
// 			toAuthor:            author,
// 			toPostID:            postID,
// 			expectErr:           acc.ErrAccountSavingCoinNotEnough().Result(),
// 			expectPostMeta:      model.PostMeta{},
// 			expectDonatorSaving: accParam.RegisterFee,
// 			expectAuthorSaving: accParam.RegisterFee.Plus(
// 				types.NewCoinFromInt64(95 * types.Decimals)),
// 			expectRegisteredEvent:             RewardEvent{},
// 			expectDonateTimesFromUserToAuthor: 1,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(100 * types.Decimals),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(95 * types.Decimals),
// 				OriginalIncome: types.NewCoinFromInt64(95 * types.Decimals),
// 			},
// 		},
// 		{
// 			testName:   "donate less money from second user with sufficient saving",
// 			donateUser: secondUserWithSufficientSaving,
// 			amount:     types.LNO("50"),
// 			toAuthor:   author,
// 			toPostID:   postID,
// 			expectErr:  sdk.Result{},
// 			expectPostMeta: model.PostMeta{
// 				CreatedAt:               ctx.BlockHeader().Time.Unix(),
// 				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
// 				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
// 				AllowReplies:            true,
// 				TotalDonateCount:        2,
// 				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReportCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReward:             types.NewCoinFromInt64(14250000),
// 				RedistributionSplitRate: sdk.ZeroDec(),
// 			},
// 			expectDonatorSaving: accParam.RegisterFee.Plus(
// 				types.NewCoinFromInt64(50 * types.Decimals)),
// 			expectAuthorSaving: accParam.RegisterFee.Plus(types.NewCoinFromInt64(14250000)),
// 			expectRegisteredEvent: RewardEvent{
// 				PostAuthor: author,
// 				PostID:     postID,
// 				Consumer:   secondUserWithSufficientSaving,
// 				Evaluate:   types.NewCoinFromInt64(1 * types.Decimals),
// 				Original:   types.NewCoinFromInt64(50 * types.Decimals),
// 				Friction:   types.NewCoinFromInt64(250000),
// 				FromApp:    "",
// 			},
// 			expectDonateTimesFromUserToAuthor: 1,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(150 * types.Decimals),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(14250000),
// 				OriginalIncome: types.NewCoinFromInt64(14250000),
// 			},
// 		},
// 		{
// 			testName:   "donate second times from second user with sufficient saving (donate stake is zero)",
// 			donateUser: secondUserWithSufficientSaving,
// 			amount:     types.LNO("50"),
// 			toAuthor:   author,
// 			toPostID:   postID,
// 			expectErr:  sdk.Result{},
// 			expectPostMeta: model.PostMeta{
// 				CreatedAt:               ctx.BlockHeader().Time.Unix(),
// 				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
// 				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
// 				AllowReplies:            true,
// 				TotalDonateCount:        3,
// 				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReportCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReward:             types.NewCoinFromInt64(190 * types.Decimals),
// 				RedistributionSplitRate: sdk.ZeroDec(),
// 			},
// 			expectDonatorSaving: accParam.RegisterFee,
// 			expectAuthorSaving:  accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
// 			expectRegisteredEvent: RewardEvent{
// 				PostAuthor: author,
// 				PostID:     postID,
// 				Consumer:   secondUserWithSufficientSaving,
// 				Evaluate:   types.NewCoinFromInt64(0),
// 				Original:   types.NewCoinFromInt64(50 * types.Decimals),
// 				Friction:   types.NewCoinFromInt64(250000),
// 				FromApp:    "",
// 			},
// 			expectDonateTimesFromUserToAuthor: 2,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(200 * types.Decimals),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(190 * types.Decimals),
// 				OriginalIncome: types.NewCoinFromInt64(190 * types.Decimals),
// 			},
// 		},
// 		{
// 			testName:   "micropayment",
// 			donateUser: micropaymentUser,
// 			amount:     types.LNO("0.00001"),
// 			toAuthor:   author,
// 			toPostID:   postID,
// 			expectErr:  sdk.Result{},
// 			expectPostMeta: model.PostMeta{
// 				CreatedAt:               ctx.BlockHeader().Time.Unix(),
// 				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
// 				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
// 				AllowReplies:            true,
// 				TotalDonateCount:        4,
// 				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReportCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReward:             types.NewCoinFromInt64(19000001),
// 				RedistributionSplitRate: sdk.ZeroDec(),
// 			},
// 			expectDonatorSaving: types.NewCoinFromInt64(199999),
// 			expectAuthorSaving:  accParam.RegisterFee.Plus(types.NewCoinFromInt64(19000001)),
// 			expectRegisteredEvent: RewardEvent{
// 				PostAuthor: author,
// 				PostID:     postID,
// 				Consumer:   micropaymentUser,
// 				Evaluate:   types.NewCoinFromInt64(1),
// 				Original:   types.NewCoinFromInt64(1),
// 				Friction:   types.NewCoinFromInt64(0),
// 				FromApp:    "",
// 			},
// 			expectDonateTimesFromUserToAuthor: 1,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(20000001),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(19000001),
// 				OriginalIncome: types.NewCoinFromInt64(19000001),
// 			},
// 		},
// 		{
// 			testName:                          "invalid target postID",
// 			donateUser:                        userWithSufficientSaving,
// 			amount:                            types.LNO("1"),
// 			toAuthor:                          author,
// 			toPostID:                          "invalid",
// 			expectErr:                         ErrPostNotFound(types.GetPermlink(author, "invalid")).Result(),
// 			expectPostMeta:                    model.PostMeta{},
// 			expectDonatorSaving:               accParam.RegisterFee,
// 			expectAuthorSaving:                accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
// 			expectRegisteredEvent:             RewardEvent{},
// 			expectDonateTimesFromUserToAuthor: 1,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(200 * types.Decimals),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(19000001),
// 				OriginalIncome: types.NewCoinFromInt64(19000001),
// 			},
// 		},
// 		{
// 			testName:                          "invalid target author",
// 			donateUser:                        userWithSufficientSaving,
// 			amount:                            types.LNO("1"),
// 			toAuthor:                          types.AccountKey("invalid"),
// 			toPostID:                          postID,
// 			expectErr:                         ErrPostNotFound(types.GetPermlink(types.AccountKey("invalid"), postID)).Result(),
// 			expectPostMeta:                    model.PostMeta{},
// 			expectDonatorSaving:               accParam.RegisterFee,
// 			expectAuthorSaving:                accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
// 			expectRegisteredEvent:             RewardEvent{},
// 			expectDonateTimesFromUserToAuthor: 0,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(200 * types.Decimals),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(19000001),
// 				OriginalIncome: types.NewCoinFromInt64(19000001),
// 			},
// 		},
// 		{
// 			testName:   "donate to self",
// 			donateUser: author,
// 			amount:     types.LNO("100"),
// 			toAuthor:   author,
// 			toPostID:   postID,
// 			expectErr:  ErrCannotDonateToSelf(author).Result(),
// 			expectPostMeta: model.PostMeta{
// 				CreatedAt:               ctx.BlockHeader().Time.Unix(),
// 				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
// 				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
// 				AllowReplies:            true,
// 				TotalDonateCount:        2,
// 				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReportCoinDay:      types.NewCoinFromInt64(0),
// 				TotalReward:             types.NewCoinFromInt64(19000001),
// 				RedistributionSplitRate: sdk.ZeroDec(),
// 			},
// 			expectDonatorSaving:               accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
// 			expectAuthorSaving:                accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
// 			expectRegisteredEvent:             RewardEvent{},
// 			expectDonateTimesFromUserToAuthor: 0,
// 			expectCumulativeConsumption:       types.NewCoinFromInt64(20000000),
// 			expectAuthorReward: accmodel.Reward{
// 				TotalIncome:    types.NewCoinFromInt64(19000001),
// 				OriginalIncome: types.NewCoinFromInt64(19000001),
// 			},
// 		},
// 		{
// 			testName:   "donate to deleted post",
// 			donateUser: userWithSufficientSaving,
// 			amount:     types.LNO("1"),
// 			toAuthor:   author1,
// 			toPostID:   deletedPostID,
// 			expectErr:  ErrDonatePostIsDeleted(types.GetPermlink(author1, deletedPostID)).Result(),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		donateMsg := NewDonateMsg(
// 			string(tc.donateUser), tc.amount, string(tc.toAuthor), tc.toPostID, "", memo1)
// 		result := handler(ctx, donateMsg)
// 		if !assert.Equal(t, tc.expectErr, result) {
// 			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectErr)
// 		}
// 		if tc.expectErr.Code.IsOK() {
// 			checkPostMeta(t, ctx, types.GetPermlink(tc.toAuthor, tc.toPostID), tc.expectPostMeta)
// 		} else {
// 			continue
// 		}

// 		authorSaving, err := am.GetSavingFromBank(ctx, tc.toAuthor)
// 		if err != nil {
// 			t.Errorf("%s: failed to get author saving from bank, got err %v", tc.testName, err)
// 		}
// 		if !authorSaving.IsEqual(tc.expectAuthorSaving) {
// 			t.Errorf("%s: diff author saving, got %v, want %v", tc.testName, authorSaving, tc.expectAuthorSaving)
// 			return
// 		}

// 		donatorSaving, err := am.GetSavingFromBank(ctx, tc.donateUser)
// 		if err != nil {
// 			t.Errorf("%s: failed to get donator saving from bank, got err %v", tc.testName, err)
// 		}
// 		if !donatorSaving.IsEqual(tc.expectDonatorSaving) {
// 			t.Errorf("%s: diff donator saving %v, got %v", tc.testName, donatorSaving, tc.expectDonatorSaving)
// 			return
// 		}

// 		if tc.expectErr.IsOK() {
// 			err := gm.CommitEventCache(ctx)
// 			if err != nil {
// 				t.Errorf("%s: failed to commit event, got err %v", tc.testName, err)
// 			}
// 			eventList := gm.GetTimeEventListAtTime(ctx, ctx.BlockHeader().Time.Unix()+3600*7*24)
// 			if !assert.Equal(t, tc.expectRegisteredEvent, eventList.Events[len(eventList.Events)-1]) {
// 				t.Errorf("%s: diff event, got %v, want %v", tc.testName,
// 					eventList.Events[len(eventList.Events)-1], tc.expectRegisteredEvent)
// 			}

// 			as := accmodel.NewAccountStorage(testAccountKVStoreKey)
// 			reward, err := as.GetReward(ctx, tc.toAuthor)
// 			if err != nil {
// 				t.Errorf("%s: failed to get reward, got err %v", tc.testName, err)
// 			}
// 			tc.expectAuthorReward.FrictionIncome = types.NewCoinFromInt64(0)
// 			tc.expectAuthorReward.InflationIncome = types.NewCoinFromInt64(0)
// 			tc.expectAuthorReward.UnclaimReward = types.NewCoinFromInt64(0)
// 			if !assert.Equal(t, tc.expectAuthorReward, *reward) {
// 				t.Errorf("%s: diff reward, got %v, want %v", tc.testName, *reward, tc.expectAuthorReward)
// 			}
// 		}

// 		cumulativeConsumption, err := gm.GetConsumption(ctx)
// 		if err != nil {
// 			t.Errorf("%s: failed to get consumption, got err %v", tc.testName, err)
// 		}
// 		if !tc.expectCumulativeConsumption.IsEqual(cumulativeConsumption) {
// 			t.Errorf("%s: diff cumulative consumption, got %v, want %v",
// 				tc.testName, cumulativeConsumption, tc.expectCumulativeConsumption)
// 			return
// 		}
// 	}
// }

// func BenchmarkNumDonate(b *testing.B) {
// 	ctx := getContext(0)
// 	ph := param.NewParamHolder(testParamKVStoreKey)
// 	ph.InitParam(ctx)
// 	accManager := acc.NewAccountManager(testAccountKVStoreKey, ph)
// 	postManager := NewPostManager(testPostKVStoreKey, ph)
// 	globalManager := global.NewGlobalManager(testGlobalKVStoreKey, ph)
// 	devManager := dev.NewDeveloperManager(testDeveloperKVStoreKey, ph)
// 	devManager.InitGenesis(ctx)
// 	voteManager := vote.NewVoteManager(testVoteKVStoreKey, ph)
// 	voteManager.InitGenesis(ctx)
// 	repManager := rep.NewReputationManager(testRepKVStoreKey, testRepV2KVStoreKey, ph)

// 	cdc := globalManager.WireCodec()
// 	cdc.RegisterInterface((*types.Event)(nil), nil)
// 	cdc.RegisterConcrete(RewardEvent{}, "event/reward", nil)

// 	InitGlobalManager(ctx, globalManager)
// 	handler := NewHandler(postManager, accManager, &globalManager, devManager, repManager)
// 	splitRate, _ := sdk.NewDecFromStr("0")

// 	resetPriv := secp256k1.GenPrivKey()
// 	txPriv := secp256k1.GenPrivKey()
// 	appPriv := secp256k1.GenPrivKey()

// 	accManager.CreateAccount(ctx, "", types.AccountKey("user1"),
// 		resetPriv.PubKey(), txPriv.PubKey(), appPriv.PubKey(), types.NewCoinFromInt64(100000*int64(b.N)))

// 	accManager.CreateAccount(ctx, "", types.AccountKey("user2"),
// 		resetPriv.PubKey(), txPriv.PubKey(), appPriv.PubKey(), types.NewCoinFromInt64(100000*int64(b.N)))
// 	postManager.CreatePost(
// 		ctx, types.AccountKey("user1"), "postID", "", "", "", "",
// 		string(make([]byte, 1000)), string(make([]byte, 50)),
// 		splitRate, []types.IDToURLMapping{})

// 	b.ResetTimer()
// 	for n := 0; n < b.N; n++ {
// 		handler(ctx, NewDonateMsg("user2", "1", "user1", "postID", "", ""))
// 	}
// }
