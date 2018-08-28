package infra

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/infra/model"
)

// InfraManager - infra manager
type InfraManager struct {
	storage     model.InfraProviderStorage
	paramHolder param.ParamHolder
}

// NewInfraManager - create NewInfraManager
func NewInfraManager(key sdk.StoreKey, holder param.ParamHolder) InfraManager {
	return InfraManager{
		storage:     model.NewInfraProviderStorage(key),
		paramHolder: holder,
	}
}

// InitGenesis - initialize infra manager
func (im InfraManager) InitGenesis(ctx sdk.Context) error {
	if err := im.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

// DoesInfraProviderExist - check if infra provide exists in KVStore or not
func (im InfraManager) DoesInfraProviderExist(ctx sdk.Context, username types.AccountKey) bool {
	return im.storage.DoesInfraProviderExist(ctx, username)
}

// RegisterInfraProvider - register infra provider on KVStore
func (im InfraManager) RegisterInfraProvider(ctx sdk.Context, username types.AccountKey) sdk.Error {
	provider := &model.InfraProvider{
		Username: username,
	}
	if err := im.storage.SetInfraProvider(ctx, username, provider); err != nil {
		return err
	}
	if err := im.AddToInfraProviderList(ctx, username); err != nil {
		return err
	}
	return nil
}

// AddToInfraProviderList - add infra provider to list
func (im InfraManager) AddToInfraProviderList(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, err := im.storage.GetInfraProviderList(ctx)
	if err != nil {
		return err
	}

	// already in the list
	if types.FindAccountInList(username, lst.AllInfraProviders) != -1 {
		return nil
	}
	lst.AllInfraProviders = append(lst.AllInfraProviders, username)
	if err := im.storage.SetInfraProviderList(ctx, lst); err != nil {
		return err
	}
	return nil
}

// RemoveFromProviderList - remove infra provider from list
func (im InfraManager) RemoveFromProviderList(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, err := im.storage.GetInfraProviderList(ctx)
	if err != nil {
		return err
	}
	// not in the list
	idx := types.FindAccountInList(username, lst.AllInfraProviders)
	if idx == -1 {
		return nil
	}
	lst.AllInfraProviders = append(lst.AllInfraProviders[:idx], lst.AllInfraProviders[idx+1:]...)
	if err := im.storage.SetInfraProviderList(ctx, lst); err != nil {
		return err
	}
	return nil
}

// ReportUsage - infra provider report usage and get reward
func (im *InfraManager) ReportUsage(ctx sdk.Context, username types.AccountKey, usage int64) sdk.Error {
	provider, err := im.storage.GetInfraProvider(ctx, username)
	if err != nil {
		return err
	}
	provider.Usage += usage
	if err := im.storage.SetInfraProvider(ctx, username, provider); err != nil {
		return err
	}
	return nil
}

// GetUsageWeight - get the usage percentage of given infra provider
func (im *InfraManager) GetUsageWeight(ctx sdk.Context, username types.AccountKey) (sdk.Rat, sdk.Error) {
	lst, err := im.storage.GetInfraProviderList(ctx)
	if err != nil {
		return sdk.NewRat(0), err
	}

	totalUsage := int64(0)
	myUsage := int64(0)
	for _, providerName := range lst.AllInfraProviders {
		curProvider, err := im.storage.GetInfraProvider(ctx, providerName)
		if err != nil {
			return sdk.NewRat(0), err
		}
		totalUsage += curProvider.Usage
		if curProvider.Username == username {
			myUsage = curProvider.Usage
		}
	}
	if totalUsage == int64(0) {
		return sdk.NewRat(1, int64(len(lst.AllInfraProviders))).Round(types.PrecisionFactor), nil
	}
	return sdk.NewRat(myUsage, totalUsage).Round(types.PrecisionFactor), nil
}

// GetInfraProviderList - get the infra provider list
func (im *InfraManager) GetInfraProviderList(ctx sdk.Context) (*model.InfraProviderList, sdk.Error) {
	return im.storage.GetInfraProviderList(ctx)
}

// ClearUsage - clear all infra provider report usage
func (im *InfraManager) ClearUsage(ctx sdk.Context) sdk.Error {
	lst, err := im.storage.GetInfraProviderList(ctx)
	if err != nil {
		return err
	}

	for _, providerName := range lst.AllInfraProviders {
		curProvider, err := im.storage.GetInfraProvider(ctx, providerName)
		if err != nil {
			return err
		}
		curProvider.Usage = 0
		if err := im.storage.SetInfraProvider(ctx, providerName, curProvider); err != nil {
			return err
		}
	}
	return nil
}

// ClearUsage - clear all infra provider report usage
