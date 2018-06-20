package developer

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/model"
)

type DeveloperManager struct {
	storage     model.DeveloperStorage `json:"infra_developer_storage"`
	paramHolder param.ParamHolder      `json:"param_holder"`
}

// create NewDeveloperManager
func NewDeveloperManager(key sdk.StoreKey, holder param.ParamHolder) DeveloperManager {
	return DeveloperManager{
		storage:     model.NewDeveloperStorage(key),
		paramHolder: holder,
	}
}

func (dm DeveloperManager) InitGenesis(ctx sdk.Context) error {
	if err := dm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (dm DeveloperManager) DoesDeveloperExist(ctx sdk.Context, username types.AccountKey) bool {
	return dm.storage.DoesDeveloperExist(ctx, username)
}

func (dm DeveloperManager) RegisterDeveloper(
	ctx sdk.Context, username types.AccountKey, deposit types.Coin) sdk.Error {
	param, err := dm.paramHolder.GetDeveloperParam(ctx)
	if err != nil {
		return err
	}
	// check developer mindmum deposit requirement
	if !deposit.IsGTE(param.DeveloperMinDeposit) {
		return ErrDeveloperDepositNotEnough()
	}

	developer := &model.Developer{
		Username: username,
		Deposit:  deposit,
	}
	if err := dm.storage.SetDeveloper(ctx, username, developer); err != nil {
		return err
	}
	if err := dm.AddToDeveloperList(ctx, username); err != nil {
		return err
	}
	return nil
}

func (dm DeveloperManager) AddToDeveloperList(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, err := dm.storage.GetDeveloperList(ctx)
	if err != nil {
		return err
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

func (dm DeveloperManager) RemoveFromDeveloperList(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	lst, err := dm.storage.GetDeveloperList(ctx)
	if err != nil {
		return err
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

func (dm DeveloperManager) ReportConsumption(
	ctx sdk.Context, username types.AccountKey, consumption types.Coin) sdk.Error {
	developer, err := dm.storage.GetDeveloper(ctx, username)
	if err != nil {
		return err
	}
	developer.AppConsumption = developer.AppConsumption.Plus(consumption)
	if err := dm.storage.SetDeveloper(ctx, username, developer); err != nil {
		return err
	}
	return nil
}

func (dm DeveloperManager) GetConsumptionWeight(
	ctx sdk.Context, username types.AccountKey) (*big.Rat, sdk.Error) {
	lst, err := dm.storage.GetDeveloperList(ctx)
	if err != nil {
		return nil, err
	}

	totalConsumption := types.NewCoinFromInt64(0)
	myConsumption := types.NewCoinFromInt64(0)
	for _, developerName := range lst.AllDevelopers {
		curDeveloper, err := dm.storage.GetDeveloper(ctx, developerName)
		if err != nil {
			return nil, err
		}
		totalConsumption = totalConsumption.Plus(curDeveloper.AppConsumption)
		if curDeveloper.Username == username {
			myConsumption = curDeveloper.AppConsumption
		}
	}
	// if not any consumption here, we evenly distribute all inflation
	if totalConsumption.ToRat().Sign() == 0 {
		return big.NewRat(1, int64(len(lst.AllDevelopers))), nil
	}
	return new(big.Rat).Quo(myConsumption.ToRat(), totalConsumption.ToRat()), nil
}

func (dm DeveloperManager) GetDeveloperList(ctx sdk.Context) (*model.DeveloperList, sdk.Error) {
	return dm.storage.GetDeveloperList(ctx)
}

func (dm DeveloperManager) ClearConsumption(ctx sdk.Context) sdk.Error {
	lst, err := dm.storage.GetDeveloperList(ctx)
	if err != nil {
		return err
	}

	for _, developerName := range lst.AllDevelopers {
		curDeveloper, err := dm.storage.GetDeveloper(ctx, developerName)
		if err != nil {
			return err
		}
		curDeveloper.AppConsumption = types.NewCoinFromInt64(0)
		if err := dm.storage.SetDeveloper(ctx, developerName, curDeveloper); err != nil {
			return err
		}
	}
	return nil
}

// this method won't check if it is a legal withdraw, caller should check by itself
func (dm DeveloperManager) Withdraw(
	ctx sdk.Context, username types.AccountKey, coin types.Coin) sdk.Error {
	developer, err := dm.storage.GetDeveloper(ctx, username)
	if err != nil {
		return err
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

func (dm DeveloperManager) WithdrawAll(
	ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	developer, err := dm.storage.GetDeveloper(ctx, username)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	if err := dm.Withdraw(ctx, username, developer.Deposit); err != nil {
		return types.NewCoinFromInt64(0), err
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
