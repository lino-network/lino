package post

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ErrGetRedistributionSplitRate(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v redistribution split rate failed", permlink))
}

func ErrGetCreatedTime(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v created time failed", permlink))
}

func ErrGetRootSourcePost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v root source post failed", permlink))
}

func ErrSetRootSourcePost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager set %v root source post failed", permlink))
}

func ErrCreatePost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager created post %v failed", permlink))
}

func ErrCreatePostSourceInvalid(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager created post %v failed, source post is invalid", permlink))
}

func ErrPostExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post %v already exist", permlink))
}

func ErrAddOrUpdateLikeToPost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add or update like to post %v failed", permlink))
}

func ErrReportOrUpvoteToPostExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("report or upvote to post %v already exists", permlink))
}

func ErrAddOrUpdateReportOrUpvoteToPost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add or update report or upvote to post %v failed", permlink))
}

func ErrAddOrUpdateViewToPost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add or update view to post %v failed", permlink))
}

func ErrRevokeReportOrUpvoteToPost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("revoke report or upvote to post %v failed", permlink))
}

func ErrAddDonation(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add donation to post %v failed", permlink))
}

func ErrDeletePost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("delete post %v failed", permlink))
}

func ErrGetPenaltyScore(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("get post %v penalty score failed", permlink))
}

func ErrCreatePostAuthorNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("create post author %v not found", author))
}

func ErrCreateExistPost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("create post failed, post %v already exist", permlink))
}

func ErrUpdatePostNotFound(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("update post failed, post %v not found", permlink))
}

func ErrDeletePostNotFound(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("delete post failed, post %v not found", permlink))
}

func ErrLikePostUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("like post failed, user %v not found", user))
}

func ErrViewPostUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("view post failed, user %v not found", user))
}

func ErrLikeNonExistPost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("like post failed, target post %v not found", permlink))
}

func ErrViewNonExistPost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("view post failed, target post %v not found", permlink))
}

func ErrDonateFailed(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed", permlink))
}

func ErrAccountCheckingCoinNotEnough(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, user checking coin not enough", permlink))
}

func ErrAccountSavingCoinNotEnough(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, user saving coin not enough", permlink))
}

func ErrDonateUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donation failed, user %v not found", user))
}

func ErrDonateAuthorNotFound(permlink types.Permlink, author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donation failed, post %v author %v not found", permlink, author))
}

func ErrDonatePostNotFound(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, post doesn't exist", permlink))
}

func ErrDonatePostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, post is deleted", permlink))
}

func ErrUpdatePostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("update post %v failed, post is deleted", permlink))
}

func ErrReportOrUpvoteFailed(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report or upvote to post %v failed", permlink))
}

func ErrReportOrUpvoteUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report or upvote failed, user %v not found", user))
}

func ErrDonateToSelf(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate failed, user %v donate to self", user))
}

func ErrMicropaymentExceedsLimitation() sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate failed, micropayment exceeds limitation"))
}

func ErrUpdatePostAuthorNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("update post failed, author %v not found", author))
}

func ErrDeletePostAuthorNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("delete post failed, author %v not found", author))
}

func ErrReportAuthorNotFound(permlink types.Permlink, author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report failed, post %v author %v not found", permlink, author))
}

func ErrReportOrUpvotePostDoesntExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report or upvote to post %v failed, post doesn't exist", permlink))
}

func ErrUpvoteUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote failed, user %v not found", user))
}

func ErrUpvoteAuthorNotFound(permlink types.Permlink, author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote failed, post %v author %v not found", permlink, author))
}

func ErrUpvotePostDoesntExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote to post %v failed, post doesn't exist", permlink))
}

func ErrNoPostID() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("No Post ID"))
}

func ErrPostIDTooLong() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Post ID too long"))
}

func ErrNoAuthor() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("No Author"))
}

func ErrCommentAndRepostError() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Post can't be comment and repost at the same time"))
}

func ErrCommentInvalidParent(parentPostKey types.Permlink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("comment post parent %v doesn't exist", parentPostKey))
}

func ErrPostLikeNoUsername() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Like needs username"))
}

func ErrPostLikeWeightOverflow(weight int64) sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Like weight overflow: %v", weight))
}

func ErrPostLikeInvalidTarget() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Like target post invalid"))
}

func ErrPostReportOrUpvoteNoUsername() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("report or upvote needs username"))
}

func ErrPostReportOrUpvoteInvalidTarget() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("report or upvote target post invalid"))
}

func ErrRedistributionSplitRateLengthTooLong() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("redistribution rate string is too long"))
}

func ErrIdentifierLengthTooLong() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("identifier is too long"))
}

func ErrURLLengthTooLong() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("url is too long"))
}

func ErrTooManyURL() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("too many url"))
}

func ErrPostViewNoUsername() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("view msg needs username"))
}

func ErrPostViewTimeInvalid(time int64) sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("view msg time invalid: %v", time))
}

func ErrPostViewInvalidTarget() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("view msg target post invalid"))
}

func ErrPostTitleExceedMaxLength() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Post title exceeds max length limitation"))
}

func ErrPostContentExceedMaxLength() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Post content exceeds max length limitation"))
}

func ErrPostRedistributionSplitRate() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Post redistribution rate invalid"))
}

func ErrPostDonateNoUsername() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Donate needs username"))
}

func ErrPostDonateInvalidTarget() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Donate target post invalid"))
}

func ErrInvalidMemo() sdk.Error {
	return types.NewError(types.CodeInvalidMemo, fmt.Sprintf("invalid memo in Donate"))
}
