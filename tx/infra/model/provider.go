package model

import (
	"github.com/lino-network/lino/types"
)

type InfraProvider struct {
	Username types.AccountKey `json:"username"`
	Usage    int64            `json:"usage"`
}

type InfraProviderList struct {
	AllInfraProviders []types.AccountKey `json:"all_infra_providers"`
}
