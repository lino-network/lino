package post

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type CodeType = sdk.CodeType

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code CodeType) string {
	switch code {
	case types.CodePostNotFound:
		return "Post Not Found"
	case types.CodePostMarshalError:
		return "Post Marshal Error"
	case types.CodePostUnmarshalError:
		return "Post Unmarshal Error"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

type NotFoundErrFunc func(types.PostKey) sdk.Error

func ErrPostNotFound(postKey types.PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post not found for key: %s", postKey))
}

func ErrPostMetaNotFound(postKey types.PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post meta not found for key: %s", postKey))
}

func ErrPostLikesNotFound(postKey types.PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post likes not found for key: %s", postKey))
}

func ErrPostCommentsNotFound(postKey types.PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post comments not found for key: %s", postKey))
}

func ErrPostViewsNotFound(postKey types.PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post views not found for key: %s", postKey))
}

func ErrPostDonationsNotFound(postKey types.PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post donations not found for key: %s", postKey))
}

func ErrPostMarshalError(err error) sdk.Error {
	return newError(types.CodePostMarshalError, fmt.Sprintf("Post marshal error: %s", err.Error()))
}

func ErrPostUnmarshalError(err error) sdk.Error {
	return newError(types.CodePostUnmarshalError, fmt.Sprintf("Post unmarshal error: %s", err.Error()))
}

func ErrPostCreateNoPermlink() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Create with empty permlink"))
}

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}

func newError(code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
