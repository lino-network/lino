package model

import (
	"github.com/lino-network/lino/types"
)

// InfraProviderRow - infra provider, pk: app
type InfraProviderRow struct {
	App      types.AccountKey `json:"app"`
	Provider InfraProvider    `json:"provider"`
}

// InfraProviderListRow - all providers, pk: none.
type InfraProviderListRow struct {
	List InfraProviderList `json:"list"`
}

// InfraTables infra storage state
type InfraTables struct {
	InfraProviders    []InfraProviderRow
	InfraProviderList InfraProviderListRow
}

// ToIR - same
func (i InfraTables) ToIR() InfraTablesIR {
	return i
}
