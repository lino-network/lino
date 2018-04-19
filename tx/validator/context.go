package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type contextKey string

const (
	contextKeyValidator contextKey = "contextKeyValidator"
)

func WithPreBlockValidators(ctx sdk.Context, validators []types.AccountKey) sdk.Context {
	return ctx.WithValue(contextKeyValidator, validators)
}

func GetPreBlockValidators(ctx sdk.Context) []types.AccountKey {
	v := ctx.Value(contextKeyValidator)
	if v == nil {
		return []types.AccountKey{}
	}
	return v.([]types.AccountKey)
}
