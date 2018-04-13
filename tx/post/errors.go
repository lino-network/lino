package post

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrGetRedistributionSplitRate(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v redistribution split rate failed", postKey))
}

func ErrGetRootSourcePost(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v root source post failed", postKey))
}

func ErrSetRootSourcePost(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("post manager set %v root source post failed", postKey))
}

func ErrCreatePost(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("post manager created post %v failed", postKey))
}

func ErrPostExist(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("post %v already exist", postKey))
}

func ErrAddOrUpdateLikeToPost(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("add or update like to post %v failed", postKey))
}

func ErrAddOrUpdateReportOrUpvoteToPost(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("add or update report or upvote to post %v failed", postKey))
}

func ErrAddDonation(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("add donation to post %v failed", postKey))
}

func ErrUpdateLastActivity(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostManagerError, fmt.Sprintf("update post %v last activity failed", postKey))
}

func ErrCreatePostAuthorNotFound(author types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("create post author %v not found", author))
}

func ErrCreateExistPost(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("create post failed, post %v already exist", postKey))
}

func ErrLikePostUserNotFound(user types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("like post failed, user %v not found", user))
}

func ErrLikeNonExistPost(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("like post failed, target post %v not found", postKey))
}

func ErrDonateFailed(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed", postKey))
}

func ErrDonateUserNotFound(user types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("donation failed, user %v not found", user))
}

func ErrDonateAuthorNotFound(postKey types.PostKey, author types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("donation failed, post %v author %v not found", postKey, author))
}

func ErrDonatePostDoesntExist(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, post doesn't exist", postKey))
}

func ErrReportFailed(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("report to post %v failed", postKey))
}

func ErrReportUserNotFound(user types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("report failed, user %v not found", user))
}

func ErrReportAuthorNotFound(postKey types.PostKey, author types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("report failed, post %v author %v not found", postKey, author))
}

func ErrReportPostDoesntExist(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("report to post %v failed, post doesn't exist", postKey))
}

func ErrUpvoteUserNotFound(user types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote failed, user %v not found", user))
}

func ErrUpvoteAuthorNotFound(postKey types.PostKey, author types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote failed, post %v author %v not found", postKey, author))
}

func ErrUpvotePostDoesntExist(postKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote to post %v failed, post doesn't exist", postKey))
}

func ErrPostCreateNoPostID() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Create with empty post id"))
}

func ErrPostCreateNoAuthor() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Create with empty author"))
}

func ErrCommentAndRepostError() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Post can't be comment and repost at the same time"))
}

func ErrCommentInvalidParent(parentPostKey types.PostKey) sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("comment post parent %v doesn't exist", parentPostKey))
}

func ErrPostLikeNoUsername() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Like needs have username"))
}

func ErrPostLikeWeightOverflow(weight int64) sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Like weight overflow: %v", weight))
}

func ErrPostLikeInvalidTarget() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Like target post invalid"))
}

func ErrPostReportOrUpvoteNoUsername() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("report or upvote needs have username"))
}

func ErrPostReportOrUpvoteInvalidTarget() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("report or upvote target post invalid"))
}

func ErrPostTitleExceedMaxLength() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Post title exceeds max length limitation"))
}

func ErrPostContentExceedMaxLength() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Post content exceeds max length limitation"))
}

func ErrPostDonateNoUsername() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Donate needs have username"))
}

func ErrPostDonateInvalidTarget() sdk.Error {
	return sdk.NewError(types.CodePostMsgError, fmt.Sprintf("Donate target post invalid"))
}
