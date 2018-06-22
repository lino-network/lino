package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NotFoundErrFunc func([]byte) sdk.Error

func ErrInvalidLinoAmount() sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Invalid Lino amount"))
}

func ErrPostNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post not found for key: %s", key))
}

func ErrPostMetaNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post meta not found for key: %s", key))
}

func ErrPostLikeNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post like not found for key: %s", key))
}

func ErrPostReportOrUpvoteNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post report or upvote not found for key: %s", key))
}

func ErrPostCommentNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post comment not found for key: %s", key))
}

func ErrPostViewNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post view not found for key: %s", key))
}

func ErrPostDonationNotFound(key []byte) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post donation not found for key: %s", key))
}

func ErrPostMarshalError(err error) sdk.Error {
	return types.NewError(types.CodePostStorageError, fmt.Sprintf("Post marshal error: %s", err.Error()))
}

func ErrPostUnmarshalError(err error) sdk.Error {
	return types.NewError(types.CodePostUnmarshalError, fmt.Sprintf("Post unmarshal error: %s", err.Error()))
}
