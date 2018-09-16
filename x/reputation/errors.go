package reputation

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
