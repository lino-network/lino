package post

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ErrAccountNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account %v is not found", author))
}

func ErrPostNotFound(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostNotFound, fmt.Sprintf("post %v doesn't exist", permlink))
}

func ErrPostAlreadyExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostAlreadyExist, fmt.Sprintf("post %v already exist", permlink))
}

func ErrInvalidPostRedistributionSplitRate() sdk.Error {
	return types.NewError(types.CodeInvalidPostRedistributionSplitRate, fmt.Sprintf("invalid post redistribution split rate"))
}

func ErrDonatePostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeDonatePostIsDeleted, fmt.Sprintf("donate to post %s failed, post is deleted", permlink))
}

func ErrGetSourcePost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeGetSourcePost, fmt.Sprintf("failed to get source post %s", permlink))
}

func ErrDeveloperNotFound(fromApp types.AccountKey) sdk.Error {
	return types.NewError(types.CodeDeveloperNotFound, fmt.Sprintf("developer %s is not found", fromApp))
}

func ErrCannotDonateToSelf(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodeCannotDonateToSelf, fmt.Sprintf("donate failed, user %v donate to self", user))
}

func ErrProcessSourceDonation(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeProcessSourceDonation, fmt.Sprintf("failed to process source donation: %s", permlink))
}

func ErrProcessDonation(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeProcessDonation, fmt.Sprintf("failed to process donation: %s", permlink))
}

func ErrUpdatePostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeUpdatePostIsDeleted, fmt.Sprintf("update post failed, post %v is deleted", permlink))
}

func ErrReportOrUpvoteAlreadyExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeReportOrUpvoteAlreadyExist, fmt.Sprintf("report or upvote to post %v already exists", permlink))
}

func ErrCreatePostSourceInvalid(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeCreatePostSourceInvalid, fmt.Sprintf("create post %v with invalid source", permlink))
}

func ErrReportOrUpvoteTooOften() sdk.Error {
	return types.NewError(types.CodeReportOrUpvoteTooOften, fmt.Sprintf("report or upvote too often, please wait"))
}

func ErrNoPostID() sdk.Error {
	return types.NewError(types.CodeNoPostID, fmt.Sprintf("no post ID"))
}

func ErrPostIDTooLong() sdk.Error {
	return types.NewError(types.CodePostIDTooLong, fmt.Sprintf("post ID is too long"))
}

func ErrNoAuthor() sdk.Error {
	return types.NewError(types.CodeNoAuthor, fmt.Sprintf("no Author"))
}

func ErrCommentAndRepostConflict() sdk.Error {
	return types.NewError(types.CodeCommentAndRepostConflict, fmt.Sprintf("post can't be comment and repost at the same time"))
}

func ErrInvalidTarget() sdk.Error {
	return types.NewError(types.CodeInvalidTarget, fmt.Sprintf("target post is invalid"))
}

func ErrRedistributionSplitRateLengthTooLong() sdk.Error {
	return types.NewError(types.CodeRedistributionSplitRateLengthTooLong, fmt.Sprintf("redistribution rate string is too long"))
}

func ErrIdentifierLengthTooLong() sdk.Error {
	return types.NewError(types.CodeIdentifierLengthTooLong, fmt.Sprintf("identifier is too long"))
}

func ErrURLLengthTooLong() sdk.Error {
	return types.NewError(types.CodeURLLengthTooLong, fmt.Sprintf("url is too long"))
}

func ErrTooManyURL() sdk.Error {
	return types.NewError(types.CodeTooManyURL, fmt.Sprintf("too many url"))
}

func ErrPostTitleExceedMaxLength() sdk.Error {
	return types.NewError(types.CodePostTitleExceedMaxLength, fmt.Sprintf("post title exceeds max length limitation"))
}

func ErrPostContentExceedMaxLength() sdk.Error {
	return types.NewError(types.CodePostContentExceedMaxLength, fmt.Sprintf("post content exceeds max length limitation"))
}

func ErrNoUsername() sdk.Error {
	return types.NewError(types.CodeNoUsername, fmt.Sprintf("username is missing"))
}

func ErrInvalidMemo() sdk.Error {
	return types.NewError(types.CodeInvalidMemo, fmt.Sprintf("invalid memo"))
}
