package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

// EventExec is a function that can execute events.
type EventExec = func(ctx sdk.Context, event linotypes.Event) sdk.Error

// BCEvent execute blockchain scheduled events.
type BCEventExec = func(ctx sdk.Context) linotypes.BCEventErr
