package model

import (
	"github.com/lino-network/lino/types"
)

// InfraProvider - infra provider of blockchain
type InfraProvider struct {
	Username types.AccountKey `json:"username"`
	Usage    int64            `json:"usage"`
}

// InfraProviderList - infra provider list of blockchain
type InfraProviderList struct {
	AllInfraProviders []types.AccountKey `json:"all_infra_providers"`
}
