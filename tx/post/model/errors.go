package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type NotFoundErrFunc func([]byte) sdk.Error

func ErrInvalidLinoAmount() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Invalid Lino amount"))
}

func ErrPostNotFound(key []byte) sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post not found for key: %s", key))
}

func ErrPostMetaNotFound(key []byte) sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post meta not found for key: %s", key))
}

func ErrPostLikeNotFound(key []byte) sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post like not found for key: %s", key))
}

func ErrPostCommentNotFound(key []byte) sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post comment not found for key: %s", key))
}

func ErrPostViewNotFound(key []byte) sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post view not found for key: %s", key))
}

func ErrPostDonationNotFound(key []byte) sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post donation not found for key: %s", key))
}

func ErrPostMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post marshal error: %s", err.Error()))
}

func ErrPostUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodePostUnmarshalError, fmt.Sprintf("Post unmarshal error: %s", err.Error()))
}

func ErrPostCreateNonExistAuthor() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Create with non-exist author"))
}

func ErrPostCreateNoParentPost() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Create with invalid parent post"))
}

func ErrPostAuthorDoesntExist() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post author doesn't exist"))
}

func ErrPostExist() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Post already exists"))
}

func ErrLikePostDoesntExist() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Target post doesn't exists"))
}

func ErrDonatePostDoesntExist() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Target post doesn't exists"))
}

func ErrPostDonateInsufficient() sdk.Error {
	return sdk.NewError(types.CodePostStorageError, fmt.Sprintf("Balance no enough"))
}
