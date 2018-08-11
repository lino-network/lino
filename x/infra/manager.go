package infra

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/infra/model"
)

type InfraManager struct {
	storage     model.InfraProviderStorage `json:"infra_provider_storage"`
	paramHolder param.ParamHolder          `json:"param_holder"`
}

// create NewInfraManager
func NewInfraManager(key sdk.StoreKey, holder param.ParamHolder) InfraManager {
	return InfraManager{
		storage:     model.NewInfraProviderStorage(key),
		paramHolder: holder,
	}
}

func (im InfraManager) InitGenesis(ctx sdk.Context) error {
	if err := im.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (im InfraManager) DoesInfraProviderExist(ctx sdk.Context, username types.AccountKey) bool {
	return im.storage.DoesInfraProviderExist(ctx, username)
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
	lst, err := im.storage.GetInfraProviderList(ctx)
	if err != nil {
		return err
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
	lst, err := im.storage.GetInfraProviderList(ctx)
	if err != nil {
		return err
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

func (im *InfraManager) GetInfraProviderList(ctx sdk.Context) (*model.InfraProviderList, sdk.Error) {
	return im.storage.GetInfraProviderList(ctx)
}

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

func FindAccountInList(me types.AccountKey, lst []types.AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}
