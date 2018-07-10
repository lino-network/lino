package post

import (
	"fmt"
	"testing"
	"time"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/model"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	abci "github.com/tendermint/abci/types"
)

func TestHandlerCreatePost(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	user := createTestAccount(t, ctx, am, "user1")

	// test valid post
	msg := CreatePostMsg{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	assert.True(t, pm.DoesPostExist(ctx, types.GetPermlink(msg.Author, msg.PostID)))

	// test invlaid author
	msg.Author = types.AccountKey("invalid")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccountNotFound(msg.Author).Result())
}

func TestHandlerUpdatePost(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	err := pm.DeletePost(ctx, types.GetPermlink(user1, postID1))
	assert.Nil(t, err)

	testCases := map[string]struct {
		msg        UpdatePostMsg
		wantResult sdk.Result
	}{
		"normal update": {
			msg:        NewUpdatePostMsg(string(user), postID, "update title", "update content", []types.IDToURLMapping(nil), "1"),
			wantResult: sdk.Result{},
		},
		"update author doesn't exist": {
			msg:        NewUpdatePostMsg("invalid", postID, "update title", "update content", []types.IDToURLMapping(nil), "1"),
			wantResult: ErrAccountNotFound("invalid").Result(),
		},
		"update post doesn't exist - invalid post ID": {
			msg:        NewUpdatePostMsg(string(user), "invalid", "update title", "update content", []types.IDToURLMapping(nil), "1"),
			wantResult: ErrPostNotFound(types.GetPermlink(user, "invalid")).Result(),
		},
		"update post doesn't exist - invalid author": {
			msg:        NewUpdatePostMsg(string(user2), postID, "update title", "update content", []types.IDToURLMapping(nil), "1"),
			wantResult: ErrPostNotFound(types.GetPermlink(user2, postID)).Result(),
		},
		"update deleted post": {
			msg:        NewUpdatePostMsg(string(user1), postID1, "update title", "update content", []types.IDToURLMapping(nil), "1"),
			wantResult: ErrUpdatePostIsDeleted(types.GetPermlink(user1, postID1)).Result(),
		},
	}
	for _, tc := range testCases {
		splitRate, err := sdk.NewRatFromDecimal(tc.msg.RedistributionSplitRate)
		assert.Nil(t, err)
		result := handler(ctx, tc.msg)
		assert.Equal(t, tc.wantResult, result)
		if tc.wantResult.Code != sdk.ABCICodeOK {
			continue
		}
		postInfo := model.PostInfo{
			PostID:       tc.msg.PostID,
			Title:        tc.msg.Title,
			Content:      tc.msg.Content,
			Author:       tc.msg.Author,
			SourceAuthor: "",
			SourcePostID: "",
			Links:        tc.msg.Links,
		}

		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time,
			LastUpdatedAt:           ctx.BlockHeader().Time,
			LastActivityAt:          ctx.BlockHeader().Time,
			AllowReplies:            true,
			IsDeleted:               false,
			RedistributionSplitRate: splitRate,
		}
		checkPostKVStore(t, ctx,
			types.GetPermlink(tc.msg.Author, tc.msg.PostID), postInfo, postMeta)
	}
}

func TestHandlerDeletePost(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user1 := createTestAccount(t, ctx, am, "user1")

	testCases := map[string]struct {
		msg        DeletePostMsg
		wantResult sdk.Result
	}{
		"normal delete": {
			msg: DeletePostMsg{
				Author: user,
				PostID: postID,
			},
			wantResult: sdk.Result{},
		},
		"author doesn't exist": {
			msg: DeletePostMsg{
				Author: types.AccountKey("invalid"),
				PostID: postID,
			},
			wantResult: ErrAccountNotFound("invalid").Result(),
		},
		"post doesn't exist - invalid author": {
			msg: DeletePostMsg{
				Author: user1,
				PostID: "postID",
			},
			wantResult: ErrPostNotFound(types.GetPermlink(user1, postID)).Result(),
		},
		"post doesn't exist - invalid postID": {
			msg: DeletePostMsg{
				Author: user,
				PostID: "invalid",
			},
			wantResult: ErrPostNotFound(types.GetPermlink(user, "invalid")).Result(),
		},
	}
	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		assert.Equal(t, tc.wantResult, result)
	}
}

func TestHandlerCreateComment(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	// test comment
	msg := CreatePostMsg{
		PostID:       "comment",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: user,
		ParentPostID: postID,
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       msg.PostID,
		Title:        msg.Title,
		Content:      msg.Content,
		Author:       msg.Author,
		ParentAuthor: msg.ParentAuthor,
		ParentPostID: msg.ParentPostID,
		SourceAuthor: msg.SourceAuthor,
		SourcePostID: msg.SourcePostID,
		Links:        msg.Links,
	}

	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat(),
	}

	checkPostKVStore(t, ctx, types.GetPermlink(user, "comment"), postInfo, postMeta)

	// check parent
	postInfo.PostID = postID
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.CreatedAt = ctx.BlockHeader().Time
	postMeta.LastUpdatedAt = ctx.BlockHeader().Time
	checkPostKVStore(t, ctx, types.GetPermlink(user, postID), postInfo, postMeta)

	// test invalid parent
	msg.PostID = "invalid post"
	msg.ParentAuthor = user
	msg.ParentPostID = "invalid parent"

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostNotFound(types.GetPermlink(user, msg.ParentPostID)).Result())

	// test duplicate comment
	msg.Author = user
	msg.PostID = "comment"
	msg.ParentAuthor = user
	msg.ParentPostID = "TestPostID"

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostAlreadyExist(types.GetPermlink(msg.Author, msg.PostID)).Result())

	// test cycle comment
	msg.Author = user
	msg.PostID = "newComment"
	msg.ParentAuthor = user
	msg.ParentPostID = "newComment"

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostNotFound(types.GetPermlink(user, msg.PostID)).Result())
}

func TestHandlerRepost(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	// test repost
	msg := CreatePostMsg{
		PostID:       "repost",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: user,
		SourcePostID: postID,
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       msg.PostID,
		Title:        msg.Title,
		Content:      msg.Content,
		Author:       msg.Author,
		ParentAuthor: msg.ParentAuthor,
		ParentPostID: msg.ParentPostID,
		SourceAuthor: msg.SourceAuthor,
		SourcePostID: msg.SourcePostID,
		Links:        msg.Links,
	}

	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat(),
	}

	checkPostKVStore(t, ctx, types.GetPermlink(user, "repost"), postInfo, postMeta)

	// test 2 depth repost
	msg.PostID = "repost-repost"
	msg.SourceAuthor = user
	msg.SourcePostID = "repost"
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: time.Now().Unix()})
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check 2 depth repost
	postInfo.PostID = "repost-repost"
	postMeta = model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat(),
	}
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = postID
	checkPostKVStore(t, ctx, types.GetPermlink(user, postInfo.PostID), postInfo, postMeta)
}

func TestHandlerPostLike(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	likeMsg := NewLikeMsg(string(user), 10000, string(user), postID)
	result := handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
	}
	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		TotalLikeCount:          1,
		TotalLikeWeight:         10000,
		RedistributionSplitRate: sdk.ZeroRat(),
	}
	checkPostKVStore(t, ctx, types.GetPermlink(user, postID), postInfo, postMeta)

	// test update like
	likeMsg = NewLikeMsg(string(user), -10000, string(user), postID)
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})
	postMeta.TotalDislikeWeight = 10000
	postMeta.TotalLikeWeight = 0
	checkPostKVStore(t, ctx, types.GetPermlink(user, postID), postInfo, postMeta)

	// test invalid like target post
	likeMsg = NewLikeMsg(string(user), -10000, string(user), "invalid")
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, ErrPostNotFound(types.GetPermlink(user, "invalid")).Result())
	checkPostKVStore(t, ctx, types.GetPermlink(user, postID), postInfo, postMeta)

	// test invalid like username
	likeMsg = NewLikeMsg("invalid", 10000, string(user), postID)
	result = handler(ctx, likeMsg)

	assert.Equal(t, result, ErrAccountNotFound(likeMsg.Username).Result())
	checkPostKVStore(t, ctx, types.GetPermlink(user, postID), postInfo, postMeta)
}

func TestHandlerPostDonate(t *testing.T) {
	ctx, am, ph, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	accParam, err := ph.GetAccountParam(ctx)
	assert.Nil(t, err)

	author, postID := createTestPost(t, ctx, "author", "postID", am, pm, "0")
	author1, deletedPostID := createTestPost(t, ctx, "author1", "delete", am, pm, "0")

	pm.DeletePost(ctx, types.GetPermlink(author1, deletedPostID))

	userWithSufficientSaving := createTestAccount(t, ctx, am, "userWithSufficientSaving")
	err = am.AddSavingCoin(
		ctx, userWithSufficientSaving, types.NewCoinFromInt64(100*types.Decimals),
		referrer, "", types.TransferIn)
	assert.Nil(t, err)

	secondUserWithSufficientSaving := createTestAccount(t, ctx, am, "secondUserWithSufficientSaving")
	err = am.AddSavingCoin(
		ctx, secondUserWithSufficientSaving, types.NewCoinFromInt64(100*types.Decimals),
		referrer, "", types.TransferIn)
	assert.Nil(t, err)

	microPaymentUser := createTestAccount(t, ctx, am, "microPaymentUser")
	err = am.AddSavingCoin(ctx, microPaymentUser, types.NewCoinFromInt64(2), referrer, "", types.TransferIn)
	assert.Nil(t, err)

	cases := []struct {
		TestName            string
		DonateUesr          types.AccountKey
		Amount              types.LNO
		ToAuthor            types.AccountKey
		IsMicropayment      bool
		ToPostID            string
		ExpectErr           sdk.Result
		ExpectPostMeta      model.PostMeta
		ExpectDonatorSaving types.Coin
		ExpectAuthorSaving  types.Coin
		//https://github.com/lino-network/lino/issues/154
		ExpectRegisteredEvent             RewardEvent
		ExpectDonateTimesFromUserToAuthor int64
		ExpectCumulativeConsumption       types.Coin
	}{
		{"donate from sufficient saving",
			userWithSufficientSaving, types.LNO("100"), author, false, postID, sdk.Result{},
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(95 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			accParam.RegisterFee, accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(95 * types.Decimals)),
			RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   userWithSufficientSaving,
				Evaluate:   types.NewCoinFromInt64(2363998),
				Original:   types.NewCoinFromInt64(100 * types.Decimals),
				Friction:   types.NewCoinFromInt64(5 * types.Decimals),
				FromApp:    "",
			}, 1, types.NewCoinFromInt64(100 * types.Decimals),
		},
		{"donate from insufficient saving",
			userWithSufficientSaving, types.LNO("100"), author, false, postID,
			acc.ErrAccountSavingCoinNotEnough().Result(),
			model.PostMeta{},
			accParam.RegisterFee, accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(95 * types.Decimals)),
			RewardEvent{}, 1, types.NewCoinFromInt64(100 * types.Decimals),
		},
		{"donate less money from second user with sufficient saving",
			secondUserWithSufficientSaving, types.LNO("50"), author, false, postID, sdk.Result{},
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalReward:             types.NewCoinFromInt64(14250000),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(50 * types.Decimals)),
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(14250000)),
			RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   secondUserWithSufficientSaving,
				Evaluate:   types.NewCoinFromInt64(1357309),
				Original:   types.NewCoinFromInt64(50 * types.Decimals),
				Friction:   types.NewCoinFromInt64(250000),
				FromApp:    "",
			}, 1, types.NewCoinFromInt64(150 * types.Decimals),
		},
		{"donate second times from second user with sufficient saving",
			secondUserWithSufficientSaving, types.LNO("50"), author, false, postID, sdk.Result{},
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        3,
				TotalReward:             types.NewCoinFromInt64(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			accParam.RegisterFee,
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   secondUserWithSufficientSaving,
				Evaluate:   types.NewCoinFromInt64(1357067),
				Original:   types.NewCoinFromInt64(50 * types.Decimals),
				Friction:   types.NewCoinFromInt64(250000),
				FromApp:    "",
			}, 2, types.NewCoinFromInt64(200 * types.Decimals),
		},
		{"invalid target postID",
			userWithSufficientSaving, types.LNO("1"), author, false, "invalid",
			ErrPostNotFound(types.GetPermlink(author, "invalid")).Result(),
			model.PostMeta{},
			accParam.RegisterFee,
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			RewardEvent{}, 1, types.NewCoinFromInt64(200 * types.Decimals),
		},
		{"invalid target author",
			userWithSufficientSaving, types.LNO("1"), types.AccountKey("invalid"), false, postID,
			ErrPostNotFound(types.GetPermlink(types.AccountKey("invalid"), postID)).Result(),
			model.PostMeta{},
			accParam.RegisterFee,
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			RewardEvent{}, 0, types.NewCoinFromInt64(200 * types.Decimals),
		},
		{"donate to self",
			author, types.LNO("100"), author, false, postID, ErrCannotDonateToSelf(author).Result(),
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalReward:             types.NewCoinFromInt64(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			RewardEvent{}, 0, types.NewCoinFromInt64(20000000),
		},
		{"invalid micropayment",
			microPaymentUser, types.LNO("10000"), author, true, postID,
			ErrMicropaymentExceedsLimitation().Result(),
			model.PostMeta{},
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(2)),
			accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(19000000)),
			RewardEvent{}, 0, types.NewCoinFromInt64(20000000),
		},
		{"micropayment",
			microPaymentUser, types.LNO("0.00001"), author, true, postID, sdk.Result{},
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        4,
				TotalReward:             types.NewCoinFromInt64(19000001),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(1)),
			accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(19000001)),
			RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   microPaymentUser,
				Evaluate:   types.NewCoinFromInt64(5),
				Original:   types.NewCoinFromInt64(1),
				Friction:   types.NewCoinFromInt64(0),
				FromApp:    "",
			}, 1, types.NewCoinFromInt64(20000001),
		},
		{"donate to deleted post",
			microPaymentUser, types.LNO("0.00001"), author1, false, deletedPostID,
			ErrDonatePostIsDeleted(types.GetPermlink(author1, deletedPostID)).Result(),
			model.PostMeta{},
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(1)),
			accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(19000001)),
			RewardEvent{}, 0, types.NewCoinFromInt64(20000001),
		},
	}

	for _, cs := range cases {
		donateMsg := NewDonateMsg(
			string(cs.DonateUesr), cs.Amount, string(cs.ToAuthor), cs.ToPostID, "", memo1, cs.IsMicropayment)
		result := handler(ctx, donateMsg)
		assert.Equal(t, cs.ExpectErr, result)
		if cs.ExpectErr.Code == sdk.ABCICodeOK {
			checkPostMeta(t, ctx, types.GetPermlink(cs.ToAuthor, cs.ToPostID), cs.ExpectPostMeta)
		}
		authorSaving, err := am.GetSavingFromBank(ctx, author)
		assert.Nil(t, err)
		if !authorSaving.IsEqual(cs.ExpectAuthorSaving) {
			t.Errorf(
				"%s: expect author saving %v, got %v",
				cs.TestName, cs.ExpectAuthorSaving, authorSaving)
			return
		}
		donatorSaving, err := am.GetSavingFromBank(ctx, cs.DonateUesr)
		assert.Nil(t, err)
		if !donatorSaving.IsEqual(cs.ExpectDonatorSaving) {
			t.Errorf(
				"%s: expect donator saving %v, got %v",
				cs.TestName, cs.ExpectDonatorSaving, donatorSaving)
			return
		}
		if cs.ExpectErr.Code == sdk.ABCICodeOK {
			eventList :=
				gm.GetTimeEventListAtTime(ctx, ctx.BlockHeader().Time+3600*7*24)
			assert.Equal(t, cs.ExpectRegisteredEvent, eventList.Events[len(eventList.Events)-1])
		}
		times, err := am.GetDonationRelationship(ctx, cs.ToAuthor, cs.DonateUesr)
		assert.Nil(t, err)
		if cs.ExpectDonateTimesFromUserToAuthor != times {
			t.Errorf(
				"%s: expect donate times %v, got %v",
				cs.TestName, cs.ExpectDonateTimesFromUserToAuthor, times)
			return
		}
		cumulativeConsumption, err := gm.GetConsumption(ctx)
		assert.Nil(t, err)
		if !cs.ExpectCumulativeConsumption.IsEqual(cumulativeConsumption) {
			t.Errorf(
				"%s: expect cumulative consumption %v, got %v",
				cs.TestName, cs.ExpectCumulativeConsumption, cumulativeConsumption)
			return
		}
	}
}

func TestHandlerRePostDonate(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0.15")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	err := am.AddSavingCoin(
		ctx, user3, types.NewCoinFromInt64(123*types.Decimals),
		referrer, "", types.TransferIn)
	assert.Nil(t, err)
	// repost
	msg := CreatePostMsg{
		PostID:       "repost",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user2,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: user1,
		SourcePostID: postID,
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	donateMsg := NewDonateMsg(
		string(user3), types.LNO("100"), string(user2), "repost", "", memo1, false)
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, sdk.Result{})
	eventList :=
		gm.GetTimeEventListAtTime(ctx, ctx.BlockHeader().Time+3600*7*24)

	// after handler check KVStore
	// check repost first
	postInfo := model.PostInfo{
		PostID:       msg.PostID,
		Title:        msg.Title,
		Content:      msg.Content,
		Author:       msg.Author,
		ParentAuthor: msg.ParentAuthor,
		ParentPostID: msg.ParentPostID,
		SourceAuthor: msg.SourceAuthor,
		SourcePostID: msg.SourcePostID,
		Links:        msg.Links,
	}
	totalReward, err := types.RatToCoin(sdk.NewRat(15 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	assert.Nil(t, err)
	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		TotalDonateCount:        1,
		TotalReward:             totalReward,
		RedistributionSplitRate: sdk.ZeroRat(),
	}
	checkPostKVStore(t, ctx, types.GetPermlink(user2, "repost"), postInfo, postMeta)
	repostRewardEvent := RewardEvent{
		PostAuthor: user2,
		PostID:     "repost",
		Consumer:   user3,
		Evaluate:   types.NewCoinFromInt64(518227),
		Original:   types.NewCoinFromInt64(15 * types.Decimals),
		Friction:   types.NewCoinFromInt64(75000),
		FromApp:    "",
	}
	assert.Equal(t, repostRewardEvent, eventList.Events[1])

	// check source post
	postMeta.TotalReward, _ = types.RatToCoin(sdk.NewRat(85 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	postInfo.Author = user1
	postInfo.PostID = postID
	postInfo.SourceAuthor = ""
	postInfo.SourcePostID = ""
	postMeta.RedistributionSplitRate = sdk.NewRat(3, 20)

	checkPostKVStore(t, ctx, types.GetPermlink(user1, postID), postInfo, postMeta)

	acc1Saving, _ := am.GetSavingFromBank(ctx, user1)
	acc2Saving, _ := am.GetSavingFromBank(ctx, user2)
	acc3Saving, _ := am.GetSavingFromBank(ctx, user3)
	acc1SavingCoin, _ := types.RatToCoin(sdk.NewRat(85 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	acc2SavingCoin, _ := types.RatToCoin(sdk.NewRat(15 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	assert.Equal(t, acc1Saving, initCoin.Plus(acc1SavingCoin))
	assert.Equal(t, acc2Saving, initCoin.Plus(acc2SavingCoin))
	assert.Equal(t, acc3Saving, initCoin.Plus(types.NewCoinFromInt64(23*types.Decimals)))

	sourceRewardEvent := RewardEvent{
		PostAuthor: user1,
		PostID:     postID,
		Consumer:   user3,
		Evaluate:   types.NewCoinFromInt64(2075784),
		Original:   types.NewCoinFromInt64(85 * types.Decimals),
		Friction:   types.NewCoinFromInt64(425000),
		FromApp:    "",
	}
	assert.Equal(t, sourceRewardEvent, eventList.Events[0])
}

func TestHandlerReportOrUpvote(t *testing.T) {
	ctx, am, ph, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)
	coinDayParam, _ := ph.GetCoinDayParam(ctx)
	accParam, _ := ph.GetAccountParam(ctx)
	postParam, _ := ph.GetPostParam(ctx)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	user4 := createTestAccount(t, ctx, am, "user4")

	baseTime := ctx.BlockHeader().Time + coinDayParam.SecondsToRecoverCoinDayStake
	permlink := types.GetPermlink(user1, postID)
	invalidPermlink := types.GetPermlink("invalid", "invalid")

	testCases := []struct {
		testName               string
		reportOrUpvoteUser     string
		isReport               bool
		targetPostAuthor       string
		targetPostID           string
		lastReportOrUpvoteAt   int64
		expectResult           sdk.Result
		expectTotalReportStake types.Coin
		expectTotalUpvoteStake types.Coin
	}{
		{"user1 report", string(user1), true, string(user1), postID, baseTime - postParam.ReportOrUpvoteInterval,
			sdk.Result{}, accParam.RegisterFee, types.NewCoinFromInt64(0)},
		{"user2 report", string(user2), true, string(user1), postID, baseTime - postParam.ReportOrUpvoteInterval,
			sdk.Result{}, accParam.RegisterFee.Plus(accParam.RegisterFee), types.NewCoinFromInt64(0)},
		{"user3 upvote", string(user3), false, string(user1), postID, baseTime - postParam.ReportOrUpvoteInterval,
			sdk.Result{}, accParam.RegisterFee.Plus(accParam.RegisterFee), accParam.RegisterFee},
		{"user1 wanna change report to upvote", string(user1), false, string(user1), postID, baseTime - postParam.ReportOrUpvoteInterval,
			ErrReportOrUpvoteAlreadyExist(permlink).Result(), accParam.RegisterFee.Plus(accParam.RegisterFee), accParam.RegisterFee},
		{"user4 report to an invalid post", string(user4), true, "invalid", "invalid", baseTime - postParam.ReportOrUpvoteInterval,
			ErrPostNotFound(invalidPermlink).Result(), accParam.RegisterFee.Plus(accParam.RegisterFee), accParam.RegisterFee},
	}

	for _, tc := range testCases {
		lastReportOrUpvoteAtCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: tc.lastReportOrUpvoteAt})
		am.UpdateLastReportOrUpvoteAt(lastReportOrUpvoteAtCtx, types.AccountKey(tc.reportOrUpvoteUser))

		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime})
		msg := NewReportOrUpvoteMsg(tc.reportOrUpvoteUser, tc.targetPostAuthor, tc.targetPostID, tc.isReport)
		result := handler(newCtx, msg)
		assert.Equal(t, tc.expectResult, result, fmt.Sprintf("%s: got %v, want %v", tc.testName, result, sdk.Result{}))
		if tc.expectResult.Code != sdk.ABCICodeOK {
			continue
		}
		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time,
			LastUpdatedAt:           ctx.BlockHeader().Time,
			LastActivityAt:          newCtx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalReportStake:        tc.expectTotalReportStake,
			TotalUpvoteStake:        tc.expectTotalUpvoteStake,
		}
		targetPost := types.GetPermlink(types.AccountKey(tc.targetPostAuthor), tc.targetPostID)
		checkPostMeta(t, ctx, targetPost, postMeta)
		lastReportOrUpvoteAt, _ := am.GetLastReportOrUpvoteAt(ctx, types.AccountKey(tc.reportOrUpvoteUser))
		assert.Equal(t, baseTime, lastReportOrUpvoteAt)
	}
}

func TestHandlerView(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm)

	createTime := ctx.BlockHeader().Time
	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	cases := []struct {
		viewUser             types.AccountKey
		postID               string
		author               types.AccountKey
		viewTime             int64
		expectTotalViewCount int64
		expectUserViewCount  int64
	}{
		{user3, postID, user1, 1, 1, 1},
		{user3, postID, user1, 2, 2, 2},
		{user2, postID, user1, 3, 3, 1},
		{user2, postID, user1, 4, 4, 2},
		{user1, postID, user1, 5, 5, 1},
	}

	for _, cs := range cases {
		postKey := types.GetPermlink(cs.author, cs.postID)
		ctx = ctx.WithBlockHeader(abci.Header{Time: cs.viewTime})
		msg := NewViewMsg(string(cs.viewUser), string(cs.author), cs.postID)
		result := handler(ctx, msg)
		assert.Equal(t, result, sdk.Result{})
		postMeta := model.PostMeta{
			CreatedAt:               createTime,
			LastUpdatedAt:           createTime,
			LastActivityAt:          createTime,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalViewCount:          cs.expectTotalViewCount,
		}
		checkPostMeta(t, ctx, postKey, postMeta)
		view, err := pm.postStorage.GetPostView(ctx, postKey, cs.viewUser)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectUserViewCount, view.Times)
		assert.Equal(t, cs.viewTime, view.LastViewAt)
	}
}
