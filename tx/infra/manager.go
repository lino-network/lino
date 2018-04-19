package infra

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/infra/model"
	"github.com/lino-network/lino/types"
)

type InfraManager struct {
	storage *model.InfraProviderStorage `json:"infra_provider_storage"`
}

// create NewInfraManager
func NewInfraManager(key sdk.StoreKey) *InfraManager {
	return &InfraManager{
		storage: model.NewInfraProviderStorage(key),
	}
}

func (im InfraManager) InitGenesis(ctx sdk.Context) error {
	if err := im.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (im InfraManager) IsInfraProviderExist(ctx sdk.Context, username types.AccountKey) bool {
	infoByte, _ := im.storage.GetInfraProvider(ctx, username)
	return infoByte != nil
}

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

func (im InfraManager) AddToInfraProviderList(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, getErr := im.storage.GetInfraProviderList(ctx)
	if getErr != nil {
		return getErr
	}

	// already in the list
	if FindAccountInList(username, lst.AllInfraProviders) != -1 {
		return nil
	}
	lst.AllInfraProviders = append(lst.AllInfraProviders, username)
	if err := im.storage.SetInfraProviderList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (im InfraManager) RemoveFromProviderList(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, getErr := im.storage.GetInfraProviderList(ctx)
	if getErr != nil {
		return getErr
	}
	// not in the list
	idx := FindAccountInList(username, lst.AllInfraProviders)
	if idx == -1 {
		return nil
	}
	lst.AllInfraProviders = append(lst.AllInfraProviders[:idx], lst.AllInfraProviders[idx+1:]...)
	if err := im.storage.SetInfraProviderList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (im *InfraManager) ReportUsage(ctx sdk.Context, username types.AccountKey, usage int64) sdk.Error {
	provider, getErr := im.storage.GetInfraProvider(ctx, username)
	if getErr != nil {
		return getErr
	}
	provider.Usage += usage
	if err := im.storage.SetInfraProvider(ctx, username, provider); err != nil {
		return err
	}
	return nil
}

func (im *InfraManager) GetUsageWeight(ctx sdk.Context, username types.AccountKey) (sdk.Rat, sdk.Error) {
	lst, getErr := im.storage.GetInfraProviderList(ctx)
	if getErr != nil {
		return sdk.NewRat(0), getErr
	}

	totalUsage := int64(0)
	myUsage := int64(0)
	for _, providerName := range lst.AllInfraProviders {
		curProvider, getErr := im.storage.GetInfraProvider(ctx, providerName)
		if getErr != nil {
			return sdk.NewRat(0), getErr
		}
		totalUsage += curProvider.Usage
		if curProvider.Username == username {
			myUsage = curProvider.Usage
		}
	}
	return sdk.NewRat(myUsage, totalUsage), nil
}

func (im *InfraManager) GetInfraProviderList(ctx sdk.Context) (*model.InfraProviderList, sdk.Error) {
	return im.storage.GetInfraProviderList(ctx)
}

func (im *InfraManager) ClearUsage(ctx sdk.Context) sdk.Error {
	lst, getErr := im.storage.GetInfraProviderList(ctx)
	if getErr != nil {
		return getErr
	}

	for _, providerName := range lst.AllInfraProviders {
		curProvider, getErr := im.storage.GetInfraProvider(ctx, providerName)
		if getErr != nil {
			return getErr
		}
		curProvider.Usage = 0
		if err := im.storage.SetInfraProvider(ctx, providerName, curProvider); err != nil {
			return err
		}
	}
	return nil
}

func FindAccountInList(me types.AccountKey, lst []types.AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}
