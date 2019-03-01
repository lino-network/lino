package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func CheckPathContentAndMinLength(path []string, expectMinLength int) sdk.Error {
	if len(path) < expectMinLength {
		return ErrInvalidQueryPath()
	}
	for i := 0; i < expectMinLength; i++ {
		if path[i] == "" {
			return ErrInvalidQueryPath()
		}
	}
	return nil
}
