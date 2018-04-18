package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type contextKey string

const (
	contextKeyOncallValidator contextKey = "contextKeyOncallValidator"
	contextKeyAllValidator    contextKey = "contextKeyAllValidator"
)

func WithOncallValidators(ctx sdk.Context, validators []types.AccountKey) sdk.Context {
	return ctx.WithValue(contextKeyOncallValidator, validators)
}

func GetOncallValidators(ctx sdk.Context) []types.AccountKey {
	v := ctx.Value(contextKeyOncallValidator)
	if v == nil {
		return []types.AccountKey{}
	}
	return v.([]types.AccountKey)
}

func WithAllValidators(ctx sdk.Context, validators []types.AccountKey) sdk.Context {
	return ctx.WithValue(contextKeyAllValidator, validators)
}

func GetAllValidators(ctx sdk.Context) []types.AccountKey {
	v := ctx.Value(contextKeyAllValidator)
	if v == nil {
		return []types.AccountKey{}
	}
	return v.([]types.AccountKey)
}
