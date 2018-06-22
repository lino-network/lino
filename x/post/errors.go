package post

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ErrGetRedistributionSplitRate(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v redistribution split rate failed", permLink))
}

func ErrGetCreatedTime(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v created time failed", permLink))
}

func ErrGetRootSourcePost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager got %v root source post failed", permLink))
}

func ErrSetRootSourcePost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager set %v root source post failed", permLink))
}

func ErrCreatePost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager created post %v failed", permLink))
}

func ErrCreatePostSourceInvalid(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post manager created post %v failed, source post is invalid", permLink))
}

func ErrPostExist(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("post %v already exist", permLink))
}

func ErrAddOrUpdateLikeToPost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add or update like to post %v failed", permLink))
}

func ErrReportOrUpvoteToPostExist(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("report or upvote to post %v already exists", permLink))
}

func ErrAddOrUpdateReportOrUpvoteToPost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add or update report or upvote to post %v failed", permLink))
}

func ErrAddOrUpdateViewToPost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add or update view to post %v failed", permLink))
}

func ErrRevokeReportOrUpvoteToPost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("revoke report or upvote to post %v failed", permLink))
}

func ErrAddDonation(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("add donation to post %v failed", permLink))
}

func ErrDeletePost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("delete post %v failed", permLink))
}

func ErrGetPenaltyScore(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostManagerError, fmt.Sprintf("get post %v penalty score failed", permLink))
}

func ErrCreatePostAuthorNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("create post author %v not found", author))
}

func ErrCreateExistPost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("create post failed, post %v already exist", permLink))
}

func ErrUpdatePostNotFound(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("update post failed, post %v not found", permLink))
}

func ErrDeletePostNotFound(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("delete post failed, post %v not found", permLink))
}

func ErrLikePostUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("like post failed, user %v not found", user))
}

func ErrViewPostUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("view post failed, user %v not found", user))
}

func ErrLikeNonExistPost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("like post failed, target post %v not found", permLink))
}

func ErrViewNonExistPost(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("view post failed, target post %v not found", permLink))
}

func ErrDonateFailed(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed", permLink))
}

func ErrAccountCheckingCoinNotEnough(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, user checking coin not enough", permLink))
}

func ErrAccountSavingCoinNotEnough(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, user saving coin not enough", permLink))
}

func ErrDonateUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donation failed, user %v not found", user))
}

func ErrDonateAuthorNotFound(permLink types.PermLink, author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donation failed, post %v author %v not found", permLink, author))
}

func ErrDonatePostNotFound(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, post doesn't exist", permLink))
}

func ErrDonatePostIsDeleted(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate to post %v failed, post is deleted", permLink))
}

func ErrUpdatePostIsDeleted(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("update post %v failed, post is deleted", permLink))
}

func ErrReportOrUpvoteFailed(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report or upvote to post %v failed", permLink))
}

func ErrReportOrUpvoteUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report or upvote failed, user %v not found", user))
}

func ErrDonateToSelf(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("donate failed, user %v donate to self", user))
}

func ErrUpdatePostAuthorNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("update post failed, author %v not found", author))
}

func ErrDeletePostAuthorNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("delete post failed, author %v not found", author))
}

func ErrReportAuthorNotFound(permLink types.PermLink, author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report failed, post %v author %v not found", permLink, author))
}

func ErrReportOrUpvotePostDoesntExist(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("report or upvote to post %v failed, post doesn't exist", permLink))
}

func ErrUpvoteUserNotFound(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote failed, user %v not found", user))
}

func ErrUpvoteAuthorNotFound(permLink types.PermLink, author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote failed, post %v author %v not found", permLink, author))
}

func ErrUpvotePostDoesntExist(permLink types.PermLink) sdk.Error {
	return types.NewError(types.CodePostHandlerError, fmt.Sprintf("upvote to post %v failed, post doesn't exist", permLink))
}

func ErrNoPostID() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("No Post ID"))
}

func ErrNoAuthor() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("No Author"))
}

func ErrCommentAndRepostError() sdk.Error {
	return types.NewError(types.CodePostMsgError, fmt.Sprintf("Post can't be comment and repost at the same time"))
}

func ErrCommentInvalidParent(parentPostKey types.PermLink) sdk.Error {
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
