package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrGetVoter() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get voter failed"))
}

func ErrGetVote() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get vote failed"))
}

func ErrGetDelegation() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get delegation failed"))
}

func ErrUsernameNotFound() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Username not found"))
}

func ErrIllegalWithdraw() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Illegal withdraw"))
}

func ErrNoCoinToWithdraw() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("No coin to withdraw"))
}

func ErrRegisterFeeNotEnough() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Register fee not enough"))
}

func ErrInvalidUsername() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Invalid Username"))
}

func ErrValidatorCannotRevoke() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Invalid revoke"))
}

func ErrNotOngoingProposal() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Not ongoing proposal"))
}
