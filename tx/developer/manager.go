package developer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/developer/model"
	"github.com/lino-network/lino/types"
)

type DeveloperManager struct {
	storage model.DeveloperStorage `json:"infra_developer_storage"`
}

// create NewDeveloperManager
func NewDeveloperManager(key sdk.StoreKey) DeveloperManager {
	return DeveloperManager{
		storage: model.NewDeveloperStorage(key),
	}
}

func (dm DeveloperManager) InitGenesis(ctx sdk.Context) error {
	if err := dm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (dm DeveloperManager) IsDeveloperExist(ctx sdk.Context, username types.AccountKey) bool {
	infoByte, _ := dm.storage.GetDeveloper(ctx, username)
	return infoByte != nil
}

func (dm DeveloperManager) RegisterDeveloper(ctx sdk.Context, username types.AccountKey, deposit types.Coin) sdk.Error {
	// check developer mindmum deposit requirement
	if !deposit.IsGTE(types.DeveloperMinDeposit) {
		return ErrDeveloperDepositNotEnough()
	}

	developer := &model.Developer{
		Username: username,
	}
	if err := dm.storage.SetDeveloper(ctx, username, developer); err != nil {
		return err
	}
	if err := dm.AddToDeveloperList(ctx, username); err != nil {
		return err
	}
	return nil
}

func (dm DeveloperManager) AddToDeveloperList(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, getErr := dm.storage.GetDeveloperList(ctx)
	if getErr != nil {
		return getErr
	}
	// already in the list
	if FindAccountInList(username, lst.AllDevelopers) != -1 {
		return nil
	}
	lst.AllDevelopers = append(lst.AllDevelopers, username)
	if err := dm.storage.SetDeveloperList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (dm DeveloperManager) RemoveFromDeveloperList(ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, getErr := dm.storage.GetDeveloperList(ctx)
	if getErr != nil {
		return getErr
	}
	// not in the list
	idx := FindAccountInList(username, lst.AllDevelopers)
	if idx == -1 {
		return nil
	}
	lst.AllDevelopers = append(lst.AllDevelopers[:idx], lst.AllDevelopers[idx+1:]...)
	if err := dm.storage.SetDeveloperList(ctx, lst); err != nil {
		return err
	}
	return nil
}

func (dm *DeveloperManager) ReportConsumption(ctx sdk.Context, username types.AccountKey, consumption types.Coin) sdk.Error {
	developer, getErr := dm.storage.GetDeveloper(ctx, username)
	if getErr != nil {
		return getErr
	}
	developer.AppConsumption = developer.AppConsumption.Plus(consumption)
	if err := dm.storage.SetDeveloper(ctx, username, developer); err != nil {
		return err
	}
	return nil
}

func (dm *DeveloperManager) GetConsumptionWeight(ctx sdk.Context, username types.AccountKey) (sdk.Rat, sdk.Error) {
	lst, getErr := dm.storage.GetDeveloperList(ctx)
	if getErr != nil {
		return sdk.NewRat(0), getErr
	}

	totalConsumption := types.NewCoin(0)
	myConsumption := types.NewCoin(0)
	for _, developerName := range lst.AllDevelopers {
		curDeveloper, getErr := dm.storage.GetDeveloper(ctx, developerName)
		if getErr != nil {
			return sdk.NewRat(0), getErr
		}
		totalConsumption = totalConsumption.Plus(curDeveloper.AppConsumption)
		if curDeveloper.Username == username {
			myConsumption = curDeveloper.AppConsumption
		}
	}
	if totalConsumption.ToRat().Equal(sdk.ZeroRat) {
		return sdk.ZeroRat, nil
	}
	return myConsumption.ToRat().Quo(totalConsumption.ToRat()), nil
}

func (dm *DeveloperManager) GetDeveloperList(ctx sdk.Context) (*model.DeveloperList, sdk.Error) {
	return dm.storage.GetDeveloperList(ctx)
}

func (dm *DeveloperManager) ClearConsumption(ctx sdk.Context) sdk.Error {
	lst, getErr := dm.storage.GetDeveloperList(ctx)
	if getErr != nil {
		return getErr
	}

	for _, developerName := range lst.AllDevelopers {
		curDeveloper, getErr := dm.storage.GetDeveloper(ctx, developerName)
		if getErr != nil {
			return getErr
		}
		curDeveloper.AppConsumption = types.NewCoin(0)
		if err := dm.storage.SetDeveloper(ctx, developerName, curDeveloper); err != nil {
			return err
		}
	}
	return nil
}

// this method won't check if it is a legal withdraw, caller should check by itself
func (dm DeveloperManager) Withdraw(ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	developer, getErr := dm.storage.GetDeveloper(ctx, username)
	if getErr != nil {
		return getErr
	}
	developer.Deposit = developer.Deposit.Minus(coin)

	if developer.Deposit.IsZero() {
		if err := dm.storage.DeleteDeveloper(ctx, username); err != nil {
			return err
		}
	} else {
		if err := dm.storage.SetDeveloper(ctx, username, developer); err != nil {
			return err
		}
	}

	return nil
}

func (dm DeveloperManager) WithdrawAll(ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	developer, getErr := dm.storage.GetDeveloper(ctx, username)
	if getErr != nil {
		return types.NewCoin(0), getErr
	}
	if err := dm.Withdraw(ctx, username, developer.Deposit); err != nil {
		return types.NewCoin(0), err
	}
	return developer.Deposit, nil
}

func FindAccountInList(me types.AccountKey, lst []types.AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}
