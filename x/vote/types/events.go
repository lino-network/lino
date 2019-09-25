package types

import (
	linotypes "github.com/lino-network/lino/types"
)

// UnassignDutyEvent - unassign duty needs a grace period and after that
// duty and frozen money will be cleared.
type UnassignDutyEvent struct {
	Username linotypes.AccountKey `json:"username"`
}
