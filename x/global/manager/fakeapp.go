package manager

//go:generate mockery -name FakeApp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linotypes "github.com/lino-network/lino/types"
)

type FakeApp interface {
	Hourly(ctx sdk.Context) []linotypes.BCEventErr
	Daily(ctx sdk.Context) []linotypes.BCEventErr
	Monthly(ctx sdk.Context) []linotypes.BCEventErr
	Yearly(ctx sdk.Context) []linotypes.BCEventErr
}
