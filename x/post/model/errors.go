package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrPostNotFound - error if post is not found in KVStore
func ErrPostNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostNotFound, fmt.Sprintf("post is not found for key: %s", key))
}

// ErrPostMetaNotFound - error if post meta is not found in KVStore
func ErrPostMetaNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostMetaNotFound, fmt.Sprintf("post meta is not found for key: %s", key))
}

// ErrPostReportOrUpvoteNotFound - error if report or upvote is not found in KVStore
func ErrPostReportOrUpvoteNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostReportOrUpvoteNotFound, fmt.Sprintf("post report or upvote not found for key: %s", key))
}

// ErrPostCommentNotFound - error if comment is not found in KVStore
func ErrPostCommentNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostCommentNotFound, fmt.Sprintf("Post comment not found for key: %s", key))
}

// ErrPostViewNotFound - error if view is not found in KVStore
func ErrPostViewNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostViewNotFound, fmt.Sprintf("Post view not found for key: %s", key))
}

// ErrPostDonationNotFound - error if post donation is not found in KVStore
func ErrPostDonationNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostDonationNotFound, fmt.Sprintf("Post donation not found for key: %s", key))
}

// ErrFailedToMarshalPostInfo - error if marshal post info failed
func ErrFailedToMarshalPostInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostInfo, fmt.Sprintf("failed to marshal post info: %s", err.Error()))
}

// ErrFailedToMarshalPostMeta - error if marshal post meta failed
func ErrFailedToMarshalPostMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostMeta, fmt.Sprintf("failed to marshal post meta: %s", err.Error()))
}

// ErrFailedToMarshalPostReportOrUpvote - error if marshal post report or upvote failed
func ErrFailedToMarshalPostReportOrUpvote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostReportOrUpvote, fmt.Sprintf("failed to marshal post report or upvote: %s", err.Error()))
}

// ErrFailedToMarshalPostComment - error if marshal post comment failed
func ErrFailedToMarshalPostComment(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostComment, fmt.Sprintf("failed to marshal post comment: %s", err.Error()))
}

// ErrFailedToMarshalPostView - error if marshal post view failed
func ErrFailedToMarshalPostView(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostView, fmt.Sprintf("failed to marshal post view: %s", err.Error()))
}

// ErrFailedToMarshalPostDonations - error if marshal post donation failed
func ErrFailedToMarshalPostDonations(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostDonations, fmt.Sprintf("failed to marshal post donations: %s", err.Error()))
}

// ErrFailedToUnmarshalPostInfo - error if unmarshal post info failed
func ErrFailedToUnmarshalPostInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostInfo, fmt.Sprintf("failed to unmarshal post info: %s", err.Error()))
}

// ErrFailedToUnmarshalPostMeta - error if unmarshal post meta failed
func ErrFailedToUnmarshalPostMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostMeta, fmt.Sprintf("failed to unmarshal post meta: %s", err.Error()))
}

// ErrFailedToUnmarshalPostReportOrUpvote - error if unmarshal post report or upvote failed
func ErrFailedToUnmarshalPostReportOrUpvote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostReportOrUpvote, fmt.Sprintf("failed to unmarshal post report or upvote: %s", err.Error()))
}

// ErrFailedToUnmarshalPostComment - error if unmarshal post comment failed
func ErrFailedToUnmarshalPostComment(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostComment, fmt.Sprintf("failed to unmarshal post comment: %s", err.Error()))
}

// ErrFailedToUnmarshalPostView - error if unmarshal post view failed
func ErrFailedToUnmarshalPostView(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostView, fmt.Sprintf("failed to unmarshal post view: %s", err.Error()))
}

// ErrFailedToUnmarshalPostDonations - error if unmarshal post donations failed
func ErrFailedToUnmarshalPostDonations(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostDonations, fmt.Sprintf("failed to unmarshal post donations: %s", err.Error()))
}
