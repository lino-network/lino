package post

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrAccountNotFound - error when account is not found
func ErrAccountNotFound(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account %v is not found", author))
}

// ErrPostNotFound - error when post is not found
func ErrPostNotFound(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostNotFound, fmt.Sprintf("post %v doesn't exist", permlink))
}

// ErrPostTooOften - error when user posting too often
func ErrPostTooOften(author types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostTooOften, fmt.Sprintf("%v post too often", author))
}

// ErrPostAlreadyExist - error when post is already exist
func ErrPostAlreadyExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodePostAlreadyExist, fmt.Sprintf("post %v already exist", permlink))
}

// ErrInvalidPostRedistributionSplitRate - error when post redistribution split rate is invalid
func ErrInvalidPostRedistributionSplitRate() sdk.Error {
	return types.NewError(types.CodeInvalidPostRedistributionSplitRate, fmt.Sprintf("invalid post redistribution split rate"))
}

// ErrDonatePostIsDeleted - error when donate to a deleted post
func ErrDonatePostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeDonatePostIsDeleted, fmt.Sprintf("donate to post %s failed, post is deleted", permlink))
}

// ErrGetSourcePost - error when get repost's source post failed
func ErrGetSourcePost(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeGetSourcePost, fmt.Sprintf("failed to get source post %s", permlink))
}

// ErrDeveloperNotFound - error when develoepr is not found
func ErrDeveloperNotFound(fromApp types.AccountKey) sdk.Error {
	return types.NewError(types.CodeDeveloperNotFound, fmt.Sprintf("developer %s is not found", fromApp))
}

// ErrCannotDonateToSelf - error when donate to self
func ErrCannotDonateToSelf(user types.AccountKey) sdk.Error {
	return types.NewError(types.CodeCannotDonateToSelf, fmt.Sprintf("donate failed, user %v donate to self", user))
}

// ErrProcessSourceDonation - error when donate to source post failed
func ErrProcessSourceDonation(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeProcessSourceDonation, fmt.Sprintf("failed to process source donation: %s", permlink))
}

// ErrProcessDonation - error when donation failed
func ErrProcessDonation(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeProcessDonation, fmt.Sprintf("failed to process donation: %s", permlink))
}

// ErrUpdatePostIsDeleted - error when update a deleted post
func ErrUpdatePostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeUpdatePostIsDeleted, fmt.Sprintf("update post failed, post %v is deleted", permlink))
}

// ErrReportOrUpvoteAlreadyExist - error when user report or upvote to a post which he already reported or upvoted
func ErrReportOrUpvoteAlreadyExist(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeReportOrUpvoteAlreadyExist, fmt.Sprintf("report or upvote to post %v already exists", permlink))
}

// ErrCreatePostSourceInvalid - error when repost's source post is invalid
func ErrCreatePostSourceInvalid(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeCreatePostSourceInvalid, fmt.Sprintf("create post %v with invalid source", permlink))
}

// ErrReportOrUpvoteTooOften - error when user report too often
func ErrReportOrUpvoteTooOften() sdk.Error {
	return types.NewError(types.CodeReportOrUpvoteTooOften, fmt.Sprintf("report or upvote too often, please wait"))
}

// ErrNoPostID - error when posting without post ID
func ErrNoPostID() sdk.Error {
	return types.NewError(types.CodeNoPostID, fmt.Sprintf("no post ID"))
}

// ErrPostIDTooLong - error when post ID is too long
func ErrPostIDTooLong() sdk.Error {
	return types.NewError(types.CodePostIDTooLong, fmt.Sprintf("post ID is too long"))
}

// ErrNoAuthor - error when posting without user
func ErrNoAuthor() sdk.Error {
	return types.NewError(types.CodeNoAuthor, fmt.Sprintf("no Author"))
}

// ErrCommentAndRepostConflict - error when posting with both source post and parent post
func ErrCommentAndRepostConflict() sdk.Error {
	return types.NewError(types.CodeCommentAndRepostConflict, fmt.Sprintf("post can't be comment and repost at the same time"))
}

// ErrInvalidTarget - error when target post is invalid
func ErrInvalidTarget() sdk.Error {
	return types.NewError(types.CodeInvalidTarget, fmt.Sprintf("target post is invalid"))
}

// ErrRedistributionSplitRateLengthTooLong - error when redistribution split rate's length is too long
func ErrRedistributionSplitRateLengthTooLong() sdk.Error {
	return types.NewError(types.CodeRedistributionSplitRateLengthTooLong, fmt.Sprintf("redistribution rate string is too long"))
}

// ErrIdentifierLengthTooLong - error when post identifier length is too long
func ErrIdentifierLengthTooLong() sdk.Error {
	return types.NewError(types.CodeIdentifierLengthTooLong, fmt.Sprintf("identifier is too long"))
}

// ErrURLLengthTooLong - error when post url length is too long
func ErrURLLengthTooLong() sdk.Error {
	return types.NewError(types.CodeURLLengthTooLong, fmt.Sprintf("url is too long"))
}

// ErrTooManyURL - error when posting with too many url
func ErrTooManyURL() sdk.Error {
	return types.NewError(types.CodeTooManyURL, fmt.Sprintf("too many url"))
}

// ErrPostTitleExceedMaxLength - error when post title is too long
func ErrPostTitleExceedMaxLength() sdk.Error {
	return types.NewError(types.CodePostTitleExceedMaxLength, fmt.Sprintf("post title exceeds max length limitation"))
}

// ErrPostContentExceedMaxLength - error when post content is too long
func ErrPostContentExceedMaxLength() sdk.Error {
	return types.NewError(types.CodePostContentExceedMaxLength, fmt.Sprintf("post content exceeds max length limitation"))
}

// ErrNoUsername - error when posting without username
func ErrNoUsername() sdk.Error {
	return types.NewError(types.CodeNoUsername, fmt.Sprintf("username is missing"))
}

// ErrInvalidMemo - error when donate memo is invalid
func ErrInvalidMemo() sdk.Error {
	return types.NewError(types.CodeInvalidMemo, fmt.Sprintf("invalid memo"))
}
