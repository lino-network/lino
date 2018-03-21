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

type NotFoundErrFunc func(PostKey) sdk.Error

func ErrPostNotFound(postKey PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post not found for key: %s", postKey))
}

func ErrPostMetaNotFound(postKey PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post meta not found for key: %s", postKey))
}

func ErrPostLikesNotFound(postKey PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post likes not found for key: %s", postKey))
}

func ErrPostCommentsNotFound(postKey PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post comments not found for key: %s", postKey))
}

func ErrPostViewsNotFound(postKey PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post views not found for key: %s", postKey))
}

func ErrPostDonationsNotFound(postKey PostKey) sdk.Error {
	return newError(types.CodePostNotFound, fmt.Sprintf("Post donations not found for key: %s", postKey))
}

func ErrPostMarshalError(err error) sdk.Error {
	return newError(types.CodePostMarshalError, fmt.Sprintf("Post marshal error: %s", err.Error()))
}

func ErrPostUnmarshalError(err error) sdk.Error {
	return newError(types.CodePostUnmarshalError, fmt.Sprintf("Post unmarshal error: %s", err.Error()))
}

func ErrPostCreateNoPostID() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Create with empty post id"))
}

func ErrPostCreateNoAuthor() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Create with empty author"))
}

func ErrPostCreateNonExistAuthor() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Create with non-exist author"))
}

func ErrPostCreateNoParentPost() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Create with invalid parent post"))
}

func ErrPostTitleExceedMaxLength() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Post title exceeds max length limitation"))
}

func ErrPostContentExceedMaxLength() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Post content exceeds max length limitation"))
}

func ErrPostAuthorDoesntExist() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Post author doesn't exist"))
}

func ErrPostExist() sdk.Error {
	return newError(types.CodePostCreateError, fmt.Sprintf("Post already exists"))
}

func ErrLikePostDoesntExist() sdk.Error {
	return newError(types.CodePostLikeError, fmt.Sprintf("Target post doesn't exists"))
}

func ErrDonatePostDoesntExist() sdk.Error {
	return newError(types.CodePostLikeError, fmt.Sprintf("Target post doesn't exists"))
}

func ErrPostLikeNoUsername() sdk.Error {
	return newError(types.CodePostLikeError, fmt.Sprintf("Like needs have username"))
}

func ErrPostLikeWeightOverflow(weight int64) sdk.Error {
	return newError(types.CodePostLikeError, fmt.Sprintf("Like weight overflow: %v", weight))
}

func ErrPostLikeInvalidTarget() sdk.Error {
	return newError(types.CodePostLikeError, fmt.Sprintf("Like target post invalid"))
}

func ErrPostDonateNoUsername() sdk.Error {
	return newError(types.CodePostDonateError, fmt.Sprintf("Donate needs have username"))
}

func ErrPostDonateInvalidTarget() sdk.Error {
	return newError(types.CodePostDonateError, fmt.Sprintf("Donate target post invalid"))
}

func ErrPostDonateInsufficient() sdk.Error {
	return newError(types.CodePostDonateError, fmt.Sprintf("Balance no enough"))
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
