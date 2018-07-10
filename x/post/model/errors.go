package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// not found error
func ErrPostNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostNotFound, fmt.Sprintf("post is not found for key: %s", key))
}

func ErrPostMetaNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostMetaNotFound, fmt.Sprintf("post meta is not found for key: %s", key))
}

func ErrPostLikeNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostLikeNotFound, fmt.Sprintf("post like not found for key: %s", key))
}

func ErrPostReportOrUpvoteNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostReportOrUpvoteNotFound, fmt.Sprintf("post report or upvote not found for key: %s", key))
}

func ErrPostCommentNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostCommentNotFound, fmt.Sprintf("Post comment not found for key: %s", key))
}

func ErrPostViewNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostViewNotFound, fmt.Sprintf("Post view not found for key: %s", key))
}

func ErrPostDonationNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostDonationNotFound, fmt.Sprintf("Post donation not found for key: %s", key))
}

// marshal error
func ErrFailedToMarshalPostInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostInfo, fmt.Sprintf("failed to marshal post info: %s", err.Error()))
}

func ErrFailedToMarshalPostMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostMeta, fmt.Sprintf("failed to marshal post meta: %s", err.Error()))
}

func ErrFailedToMarshalPostLike(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostLike, fmt.Sprintf("failed to marshal post like: %s", err.Error()))
}

func ErrFailedToMarshalPostReportOrUpvote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostReportOrUpvote, fmt.Sprintf("failed to marshal post report or upvote: %s", err.Error()))
}

func ErrFailedToMarshalPostComment(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostComment, fmt.Sprintf("failed to marshal post comment: %s", err.Error()))
}

func ErrFailedToMarshalPostView(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostView, fmt.Sprintf("failed to marshal post view: %s", err.Error()))
}

func ErrFailedToMarshalPostDonations(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostDonations, fmt.Sprintf("failed to marshal post donations: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalPostInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostInfo, fmt.Sprintf("failed to unmarshal post info: %s", err.Error()))
}

func ErrFailedToUnmarshalPostMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostMeta, fmt.Sprintf("failed to unmarshal post meta: %s", err.Error()))
}

func ErrFailedToUnmarshalPostLike(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostLike, fmt.Sprintf("failed to unmarshal post like: %s", err.Error()))
}

func ErrFailedToUnmarshalPostReportOrUpvote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostReportOrUpvote, fmt.Sprintf("failed to unmarshal post report or upvote: %s", err.Error()))
}

func ErrFailedToUnmarshalPostComment(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostComment, fmt.Sprintf("failed to unmarshal post comment: %s", err.Error()))
}

func ErrFailedToUnmarshalPostView(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostView, fmt.Sprintf("failed to unmarshal post view: %s", err.Error()))
}

func ErrFailedToUnmarshalPostDonations(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostDonations, fmt.Sprintf("failed to unmarshal post donations: %s", err.Error()))
}
