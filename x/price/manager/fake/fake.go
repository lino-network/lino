package fake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

type ValidatorAndVote struct {
	Username linotypes.AccountKey
	Votes    linotypes.Coin
}

//go:generate mockery -name FakeValidator

type FakeValidator interface {
	GetValidatorAndVotes(ctx sdk.Context) []ValidatorAndVote
	DoesValidatorExist(ctx sdk.Context, user linotypes.AccountKey) bool
}

func ToValNames(vals []ValidatorAndVote) (rst []linotypes.AccountKey) {
	for _, val := range vals {
		rst = append(rst, val.Username)
	}
	return
}
