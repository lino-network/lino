package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

// ErrAccountNotFound - error when account is not found
func ErrAccountNotFound(author linotypes.AccountKey) sdk.Error {
	return linotypes.NewError(linotypes.CodeAccountNotFound, fmt.Sprintf("account %v is not found", author))
}

// ErrPostNotFound - error when post is not found
func ErrPostNotFound(permlink linotypes.Permlink) sdk.Error {
	return linotypes.NewError(linotypes.CodePostNotFound, fmt.Sprintf("post %v doesn't exist", permlink))
}

// ErrPostAlreadyExist - error when post is already exist
func ErrPostAlreadyExist(permlink linotypes.Permlink) sdk.Error {
	return linotypes.NewError(linotypes.CodePostAlreadyExist, fmt.Sprintf("post %v already exist", permlink))
}

// ErrPostDeleted - error when post has been deleted.
func ErrPostDeleted(permlink linotypes.Permlink) sdk.Error {
	return linotypes.NewError(linotypes.CodePostDeleted, fmt.Sprintf("permlink %v was deleted", permlink))
}

// ErrDeveloperNotFound - error when develoepr is not found
func ErrDeveloperNotFound(fromApp linotypes.AccountKey) sdk.Error {
	return linotypes.NewError(linotypes.CodeDeveloperNotFound, fmt.Sprintf("developer %s is not found", fromApp))
}

// ErrCannotDonateToSelf - error when donate to self
func ErrCannotDonateToSelf(user linotypes.AccountKey) sdk.Error {
	return linotypes.NewError(linotypes.CodeCannotDonateToSelf, fmt.Sprintf("donate failed, user %v donate to self", user))
}

// ErrInvalidDonationAmount - error when donation amount is invalid.
func ErrInvalidDonationAmount(amount linotypes.Coin) sdk.Error {
	return linotypes.NewError(linotypes.CodeDonationAmountInvalid, fmt.Sprintf("donation amount is invalid: %s", amount))
}

// ErrProcessDonation - error when donation failed
func ErrProcessDonation(permlink linotypes.Permlink) sdk.Error {
	return linotypes.NewError(linotypes.CodeProcessDonation, fmt.Sprintf("failed to process donation: %s", permlink))
}

// ErrNoPostID - error when posting without post ID
func ErrNoPostID() sdk.Error {
	return linotypes.NewError(linotypes.CodeNoPostID, fmt.Sprintf("no post ID"))
}

// ErrPostIDTooLong - error when post ID is too long
func ErrPostIDTooLong() sdk.Error {
	return linotypes.NewError(linotypes.CodePostIDTooLong, fmt.Sprintf("post ID is too long"))
}

// ErrInvalidAuthor - error when posting without user
func ErrInvalidAuthor() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidAuthor, fmt.Sprintf("invalid Author"))
}

// ErrInvalidCreatedBy - error when posting without createdBy
func ErrInvalidCreatedBy() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidCreatedBy, fmt.Sprintf("invalid CreatedBy"))
}

// ErrInvalidTarget - error when target post is invalid
func ErrInvalidTarget() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidTarget, fmt.Sprintf("target post is invalid"))
}

// ErrPostTitleExceedMaxLength - error when post title is too long
func ErrPostTitleExceedMaxLength() sdk.Error {
	return linotypes.NewError(linotypes.CodePostTitleExceedMaxLength, fmt.Sprintf("post title exceeds max length limitation"))
}

// ErrPostContentExceedMaxLength - error when post content is too long
func ErrPostContentExceedMaxLength() sdk.Error {
	return linotypes.NewError(linotypes.CodePostContentExceedMaxLength, fmt.Sprintf("post content exceeds max length limitation"))
}

// ErrInvalidUsername - error when posting without username
func ErrInvalidUsername() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidUsername, fmt.Sprintf("invalid username"))
}

// ErrInvalidMemo - error when donate memo is invalid
func ErrInvalidMemo() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidMemo, fmt.Sprintf("invalid memo"))
}

// ErrQueryFailed - error when query post store failed
func ErrQueryFailed() sdk.Error {
	return linotypes.NewError(linotypes.CodePostQueryFailed, fmt.Sprintf("query post store failed"))
}

// ErrInvalidApp - error when making an IDA donation without specifying app.
func ErrInvalidApp() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidApp, fmt.Sprintf("invalid App"))
}

// ErrNonPositiveIDAAmount - error when ida amount is invalid.
func ErrNonPositiveIDAAmount(v linotypes.MiniIDA) sdk.Error {
	return linotypes.NewError(linotypes.CodeNonPositiveIDAAmount, fmt.Sprintf("nonpositive IDA amount: %v", v))
}

// ErrDonateAmountTooLittle -
func ErrDonateAmountTooLittle() sdk.Error {
	return linotypes.NewError(linotypes.CodeDonateAmountTooLittle, fmt.Sprintf("donation amount is too small"))
}

// ErrInvalidSigner - signes does not match app.
func ErrInvalidSigner() sdk.Error {
	return linotypes.NewError(
		linotypes.CodeInvalidSigner, fmt.Sprintf("signer does not match app, post"))
}
