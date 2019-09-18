package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linotypes "github.com/lino-network/lino/types"
)

// ErrAccountNotFound - error when account is not found
func ErrAccountNotFound() sdk.Error {
	return linotypes.NewError(linotypes.CodeAccountNotFound, fmt.Sprintf("username is not found"))
}

// ErrInsufficientDeposit - error when deposit is insufficient
func ErrInsufficientDeposit() sdk.Error {
	return linotypes.NewError(linotypes.CodeInsufficientDeposit, fmt.Sprintf("voting deposit fee not enough"))
}

// ErrInsufficientDeposit - error if required voting deposit less than committing deposit
func ErrUnbalancedAccount() sdk.Error {
	return linotypes.NewError(linotypes.CodeUnbalancedAccount, fmt.Sprintf("committing deposit not enough"))
}

// ErrIllegalWithdraw - error if withdraw less than minimum withdraw requirement
func ErrIllegalWithdraw() sdk.Error {
	return linotypes.NewError(linotypes.CodeIllegalWithdraw, fmt.Sprintf("illegal withdraw"))
}

// ErrInvalidCoin - error if coin in msg is invalid
func ErrInvalidCoin() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidCoin, fmt.Sprintf("no coin to withdraw"))
}

// ErrInvalidUsername - error if username is invalid
func ErrInvalidUsername() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

// ErrInvalidWebsite - error if website is invalid
func ErrInvalidWebsite() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidWebsite, fmt.Sprintf("Invalida website"))
}

// ErrInvalidVotedValidators - error if voted too many validators
func ErrInvalidVotedValidators() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidVotedValidators, fmt.Sprintf("Invalid voted validators"))
}

// ErrValidatorPubKeyAlreadyExist - error if validator public key is already exist
func ErrValidatorPubKeyAlreadyExist() sdk.Error {
	return linotypes.NewError(linotypes.CodeValidatorPubKeyAlreadyExist, fmt.Sprintf("validator public key has been registered"))
}

// ErrQueryFailed - error when query validator store failed
func ErrQueryFailed() sdk.Error {
	return linotypes.NewError(linotypes.CodeValidatorQueryFailed, fmt.Sprintf("query validator store failed"))
}

// not found
func ErrValidatorNotFound() sdk.Error {
	return linotypes.NewError(linotypes.CodeValidatorNotFound, fmt.Sprintf("validator is not found"))
}

func ErrValidatorListNotFound() sdk.Error {
	return linotypes.NewError(linotypes.CodeValidatorListNotFound, fmt.Sprintf("validator list is not found"))
}

func ErrElectionListNotFound() sdk.Error {
	return linotypes.NewError(linotypes.CodeElectionListNotFound, fmt.Sprintf("Election list is not found"))
}

// ErrValidatorAlreadyExist - error if validator is already exist
func ErrValidatorAlreadyExist() sdk.Error {
	return linotypes.NewError(linotypes.CodeValidatorAlreadyExist, fmt.Sprintf("validator has been registered"))
}

// ErrInvalidVoterDuty - error when developer attempting to be regsitered is not a voter.
func ErrInvalidVoterDuty() sdk.Error {
	return linotypes.NewError(
		linotypes.CodeInvalidVoterDuty, fmt.Sprintf("user's duty is not voter"))
}
