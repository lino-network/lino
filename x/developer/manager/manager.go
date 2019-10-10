package developer

import (
	"fmt"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/developer/model"
	"github.com/lino-network/lino/x/developer/types"
	"github.com/lino-network/lino/x/price"
	"github.com/lino-network/lino/x/vote"
	votetypes "github.com/lino-network/lino/x/vote/types"
)

const (
	maxAffiliatedAccount = 500

	exportVersion = 1
	importVersion = 1
)

type DeveloperManager struct {
	storage model.DeveloperStorage
	// deps
	paramHolder param.ParamKeeper
	vote        vote.VoteKeeper
	acc         account.AccountKeeper
	price       price.PriceKeeper
}

// NewDeveloperManager - create new developer manager
func NewDeveloperManager(key sdk.StoreKey, holder param.ParamKeeper, vote vote.VoteKeeper, acc account.AccountKeeper, price price.PriceKeeper) DeveloperManager {
	return DeveloperManager{
		storage:     model.NewDeveloperStorage(key),
		paramHolder: holder,
		vote:        vote,
		acc:         acc,
		price:       price,
	}
}

// InitGenesis - init developer manager
func (dm DeveloperManager) InitGenesis(ctx sdk.Context, reservePoolAmount linotypes.Coin) sdk.Error {
	if !reservePoolAmount.IsNotNegative() {
		return types.ErrInvalidReserveAmount(reservePoolAmount)
	}
	dm.storage.SetReservePool(ctx, &model.ReservePool{
		Total: reservePoolAmount,
	})
	return nil
}

// DoesDeveloperExist - check if given developer exists and not deleted before.
func (dm DeveloperManager) DoesDeveloperExist(ctx sdk.Context, username linotypes.AccountKey) bool {
	dev, err := dm.storage.GetDeveloper(ctx, username)
	if err != nil {
		return false
	}
	return !dev.IsDeleted
}

func (dm DeveloperManager) GetDeveloper(ctx sdk.Context, username linotypes.AccountKey) (model.Developer, sdk.Error) {
	dev, err := dm.storage.GetDeveloper(ctx, username)
	if err != nil {
		return model.Developer{}, err
	}
	return *dev, nil
}

// GetLiveDevelopers - returns all developers that are live(not deregistered).
func (dm DeveloperManager) GetLiveDevelopers(ctx sdk.Context) []model.Developer {
	rst := make([]model.Developer, 0)
	devs := dm.storage.GetAllDevelopers(ctx)
	for _, dev := range devs {
		if !dev.IsDeleted {
			rst = append(rst, dev)
		}
	}
	return rst
}

// RegisterDeveloper - register a developer.
// Stateful validation:
// 1. account exists.
// 2. has never been a developer.
// 3. VoteDuty is simply a voter, not a validator candidate.
// 4. not an affiliated account.
// 4. has minimum LS.
func (dm DeveloperManager) RegisterDeveloper(ctx sdk.Context, username linotypes.AccountKey, website, description, appMetaData string) sdk.Error {
	if !dm.acc.DoesAccountExist(ctx, username) {
		return types.ErrAccountNotFound()
	}
	if dm.storage.HasDeveloper(ctx, username) {
		return types.ErrDeveloperAlreadyExist(username)
	}
	if duty, err := dm.vote.GetVoterDuty(ctx, username); err != nil || duty != votetypes.DutyVoter {
		return types.ErrInvalidVoterDuty()
	}
	if dm.storage.HasUserRole(ctx, username) {
		return types.ErrInvalidUserRole()
	}
	// check developer minimum LS requirement
	param, err := dm.paramHolder.GetDeveloperParam(ctx)
	if err != nil {
		return err
	}
	staked, err := dm.vote.GetLinoStake(ctx, username)
	if err != nil {
		return err
	}
	if !staked.IsGTE(param.DeveloperMinDeposit) {
		return types.ErrInsufficientDeveloperDeposit()
	}

	// assign duty in vote
	err = dm.vote.AssignDuty(ctx, username, votetypes.DutyApp, param.DeveloperMinDeposit)
	if err != nil {
		return err
	}
	developer := &model.Developer{
		Username:    username,
		Website:     website,
		Description: description,
		AppMetaData: appMetaData,
		IsDeleted:   false,
	}
	dm.storage.SetDeveloper(ctx, *developer)
	return nil
}

// UpdateDeveloper - update developer.
// 1. developer must not be deleted.
func (dm DeveloperManager) UpdateDeveloper(ctx sdk.Context, username linotypes.AccountKey, website, description, appMetadata string) sdk.Error {
	// Use DoesDeveloperExist to make sure we don't forget to check IsDeleted.
	if !dm.DoesDeveloperExist(ctx, username) {
		return types.ErrDeveloperNotFound()
	}
	developer, err := dm.GetDeveloper(ctx, username)
	if err != nil {
		return err
	}
	developer.Website = website
	developer.Description = description
	developer.AppMetaData = appMetadata
	dm.storage.SetDeveloper(ctx, developer)
	return nil
}

// UnregisterDeveloper - unregister a developer
// validation:
// 1. Developer exists.
// 2. No IDA issued or IDA revoked.
// TODO:
// remove all affiliated accounts.
// mark developer as deleted.
func (dm DeveloperManager) UnregisterDeveloper(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	return linotypes.ErrUnimplemented("UnregisterDeveloper")
}

// IssueIDA - Application issue IDA
func (dm DeveloperManager) IssueIDA(ctx sdk.Context, appname linotypes.AccountKey, idaName string, idaPrice int64) sdk.Error {
	if !dm.DoesDeveloperExist(ctx, appname) {
		return types.ErrDeveloperNotFound()
	}
	// not issued before
	if dm.storage.HasIDA(ctx, appname) {
		return types.ErrIDAIssuedBefore()
	}

	// internally we store MiniIDAPrice, which is 10^(-5) * IDAPrice in MiniDollar, then we have
	// 1 IDA = IDAPrice * 10^(-3) Dollar
	// 1 IDA = IDAPrice * 10^(7) MiniDollar
	// 10^(5) MiniIDA = IDAPrice * 10^(7) MiniDollar
	// 1 MiniIDA = IDAPrice * 10^(2) * MiniDollar
	miniIDAPrice := linotypes.NewMiniDollarFromInt(sdk.NewInt(idaPrice).MulRaw(100))
	ida := model.AppIDA{
		App:             appname,
		Name:            idaName,
		MiniIDAPrice:    miniIDAPrice,
		IsRevoked:       false,
		RevokeCoinPrice: linotypes.NewMiniDollar(0),
	}
	dm.storage.SetIDA(ctx, ida)
	dm.storage.SetIDAStats(ctx, appname, model.AppIDAStats{
		Total: linotypes.NewMiniDollar(0),
	})
	return nil
}

// MintIDA - mint some IDA by converting LINO to IDA (internally MiniDollar).
func (dm DeveloperManager) MintIDA(ctx sdk.Context, appname linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	if _, err := dm.validAppIDA(ctx, appname); err != nil {
		return err
	}
	bank := dm.storage.GetIDABank(ctx, appname, appname)
	// ignored
	// if bank.Unauthed {
	// 	return types.ErrIDAUnauthed()
	// }
	miniDollar, err := dm.exchangeMiniDollar(ctx, appname, amount)
	if err != nil {
		return err
	}
	bank.Balance = bank.Balance.Plus(miniDollar)
	dm.storage.SetIDABank(ctx, appname, appname, bank)
	return nil
}

// exchangeMinidollar - exchange Minidollar with @p amount coin, return max exchanged minidollar.
// Return minidollar must be added to somewhere, otherwise it's burned.
// 1. minus account saving.
// 2. sending coin to reserve pool.
func (dm DeveloperManager) exchangeMiniDollar(ctx sdk.Context, appname linotypes.AccountKey, amount linotypes.Coin) (linotypes.MiniDollar, sdk.Error) {
	bought, err := dm.price.CoinToMiniDollar(ctx, amount)
	if err != nil {
		return linotypes.NewMiniDollar(0), err
	}
	if !bought.IsPositive() {
		return linotypes.NewMiniDollar(0), types.ErrExchangeMiniDollarZeroAmount()
	}
	// exchange
	err = dm.acc.MoveToPool(ctx,
		linotypes.DevIDAReservePool, linotypes.NewAccOrAddrFromAcc(appname), amount)
	if err != nil {
		return linotypes.NewMiniDollar(0), err
	}
	pool := dm.storage.GetReservePool(ctx)
	pool.Total = pool.Total.Plus(amount)
	pool.TotalMiniDollar = pool.TotalMiniDollar.Plus(bought)
	dm.storage.SetReservePool(ctx, pool)
	idaStats := dm.storage.GetIDAStats(ctx, appname)
	idaStats.Total = idaStats.Total.Plus(bought)
	dm.storage.SetIDAStats(ctx, appname, *idaStats)
	return bought, nil
}

// AppTransferIDA - transfer IDA back or from app, by app.
func (dm DeveloperManager) AppTransferIDA(ctx sdk.Context, appname, signer linotypes.AccountKey, amount linotypes.MiniIDA, from, to linotypes.AccountKey) sdk.Error {
	if !(from == appname || to == appname) {
		return types.ErrInvalidTransferTarget()
	}
	ida, err := dm.validAppIDA(ctx, appname)
	if err != nil {
		return err
	}
	if !dm.acc.DoesAccountExist(ctx, from) || !dm.acc.DoesAccountExist(ctx, to) {
		return types.ErrAccountNotFound()
	}
	signerApp, err := dm.GetAffiliatingApp(ctx, signer)
	if err != nil || signerApp != appname {
		return types.ErrInvalidSigner()
	}
	minidollar := linotypes.MiniIDAToMiniDollar(amount, ida.MiniIDAPrice)
	return dm.appIDAMove(ctx, appname, from, to, minidollar)
}

// MoveIDA - app move ida, authorization check applied.
// 1. amount must > 0.
// 2. from's bank is not frozen.
func (dm DeveloperManager) MoveIDA(ctx sdk.Context, app, from, to linotypes.AccountKey, amount linotypes.MiniDollar) sdk.Error {
	if _, err := dm.validAppIDA(ctx, app); err != nil {
		return err
	}
	if !dm.acc.DoesAccountExist(ctx, from) || !dm.acc.DoesAccountExist(ctx, to) {
		return types.ErrAccountNotFound()
	}
	return dm.appIDAMove(ctx, app, from, to, amount)
}

func (dm DeveloperManager) validAppIDA(ctx sdk.Context, app linotypes.AccountKey) (*model.AppIDA, sdk.Error) {
	if !dm.DoesDeveloperExist(ctx, app) {
		return nil, types.ErrDeveloperNotFound()
	}
	ida, err := dm.storage.GetIDA(ctx, app)
	if err != nil {
		return nil, err
	}
	if ida.IsRevoked {
		return nil, types.ErrIDARevoked()
	}
	return ida, nil
}

// appIDAMove - authorization check applied. app is not checked.
func (dm DeveloperManager) appIDAMove(ctx sdk.Context, app, from, to linotypes.AccountKey, amount linotypes.MiniDollar) sdk.Error {
	if !amount.IsPositive() {
		return linotypes.ErrInvalidIDAAmount()
	}
	fromBank := dm.storage.GetIDABank(ctx, app, from)
	toBank := dm.storage.GetIDABank(ctx, app, to)
	if fromBank.Unauthed {
		return types.ErrIDAUnauthed()
	}
	if fromBank.Balance.LT(amount) {
		return types.ErrNotEnoughIDA()
	}
	fromBank.Balance = fromBank.Balance.Minus(amount)
	toBank.Balance = toBank.Balance.Plus(amount)
	dm.storage.SetIDABank(ctx, app, from, fromBank)
	dm.storage.SetIDABank(ctx, app, to, toBank)
	return nil
}

func (dm DeveloperManager) GetMiniIDAPrice(ctx sdk.Context, app linotypes.AccountKey) (linotypes.MiniDollar, sdk.Error) {
	ida, err := dm.validAppIDA(ctx, app)
	if err != nil {
		return linotypes.NewMiniDollar(0), err
	}
	return ida.MiniIDAPrice, nil
}

// BurnIDA - Burn some @p amount of IDA on @p user's account and return coins
// poped from reserve pool. NOTE: cannot burn 0 coins.
func (dm DeveloperManager) BurnIDA(ctx sdk.Context, app, user linotypes.AccountKey, amount linotypes.MiniDollar) (linotypes.Coin, sdk.Error) {
	bank, err := dm.GetIDABank(ctx, app, user)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}
	if bank.Unauthed {
		return linotypes.NewCoinFromInt64(0), types.ErrIDAUnauthed()
	}
	if bank.Balance.LT(amount) {
		return linotypes.NewCoinFromInt64(0), types.ErrNotEnoughIDA()
	}
	bought, used, err := dm.price.MiniDollarToCoin(ctx, amount)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}
	if !bought.IsPositive() {
		return linotypes.NewCoinFromInt64(0), types.ErrBurnZeroIDA()
	}
	// after burn, move coins from the reserve pool to the user's account.
	// only called upon donation, so the newly added coins will then be moved
	// to the vote's frictions pool.
	err = dm.acc.MoveFromPool(ctx,
		linotypes.DevIDAReservePool, linotypes.NewAccOrAddrFromAcc(user), bought)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}
	pool := dm.storage.GetReservePool(ctx)
	if !pool.Total.IsGTE(bought) {
		return linotypes.NewCoinFromInt64(0), types.ErrInsuffientReservePool()
	}
	pool.Total = pool.Total.Minus(bought)
	pool.TotalMiniDollar = pool.TotalMiniDollar.Minus(used)
	dm.storage.SetReservePool(ctx, pool)
	idaStats := dm.storage.GetIDAStats(ctx, app)
	idaStats.Total = idaStats.Total.Minus(used)
	dm.storage.SetIDAStats(ctx, app, *idaStats)
	bank.Balance = bank.Balance.Minus(used)
	dm.storage.SetIDABank(ctx, app, user, &bank)
	return bought, nil
}

func (dm DeveloperManager) GetIDA(ctx sdk.Context, app linotypes.AccountKey) (model.AppIDA, sdk.Error) {
	ida, err := dm.validAppIDA(ctx, app)
	if err != nil {
		return model.AppIDA{}, err
	}
	return *ida, err
}

func (dm DeveloperManager) GetIDABank(ctx sdk.Context, app, user linotypes.AccountKey) (model.IDABank, sdk.Error) {
	if !dm.DoesDeveloperExist(ctx, app) {
		return model.IDABank{}, types.ErrDeveloperNotFound()
	}
	if !dm.acc.DoesAccountExist(ctx, user) {
		return model.IDABank{}, types.ErrAccountNotFound()
	}
	return *dm.storage.GetIDABank(ctx, app, user), nil
}

// UpdateAffiliated - add or remove an affiliated account.
func (dm DeveloperManager) UpdateAffiliated(ctx sdk.Context, appname, username linotypes.AccountKey, activate bool) sdk.Error {
	if !dm.DoesDeveloperExist(ctx, appname) {
		return types.ErrDeveloperNotFound()
	}
	if !dm.acc.DoesAccountExist(ctx, username) {
		return types.ErrAccountNotFound()
	}
	app, err := dm.storage.GetDeveloper(ctx, appname)
	if err != nil {
		return err
	}
	if app.NAffiliated >= maxAffiliatedAccount {
		return types.ErrMaxAffiliatedExceeded()
	}
	if activate {
		err := dm.addAffiliated(ctx, appname, username)
		if err != nil {
			return err
		}
		app.NAffiliated += 1
		dm.storage.SetDeveloper(ctx, *app)
	} else {
		err := dm.removeAffiliated(ctx, appname, username)
		if err != nil {
			return err
		}
		app.NAffiliated -= 1
		dm.storage.SetDeveloper(ctx, *app)
	}
	return nil
}

// To activate an affiliated account, check:
// 1. not affiliated to any developer.
// 2. not a developer.
// 3. not on any other duty.
func (dm DeveloperManager) addAffiliated(ctx sdk.Context, app, username linotypes.AccountKey) sdk.Error {
	if dm.storage.HasUserRole(ctx, username) {
		return types.ErrInvalidAffiliatedAccount("is affiliated already")
	}
	// TODO(@yumin): Do we check if username as developer is deleted already?
	if dm.storage.HasDeveloper(ctx, username) {
		return types.ErrInvalidAffiliatedAccount("is/was developer")
	}
	duty, err := dm.vote.GetVoterDuty(ctx, username)
	if err == nil && duty != votetypes.DutyVoter {
		return types.ErrInvalidAffiliatedAccount("on duty of something else")
	}
	dm.storage.SetAffiliatedAcc(ctx, app, username)
	dm.storage.SetUserRole(ctx, username, &model.Role{
		AffiliatedApp: app,
	})
	return nil
}

// To remove an affiliated account from app
// 1. user is the affiliated account of app.
func (dm DeveloperManager) removeAffiliated(ctx sdk.Context, app, username linotypes.AccountKey) sdk.Error {
	role, err := dm.storage.GetUserRole(ctx, username)
	if err != nil {
		return err
	}
	if role.AffiliatedApp != app {
		return types.ErrInvalidAffiliatedAccount("not affiliated account of provided app")
	}
	dm.storage.DelAffiliatedAcc(ctx, app, username)
	dm.storage.DelUserRole(ctx, username)
	return nil
}

// GetAffiliatingApp - get username's affiliating app, or username itself is an app.
func (dm DeveloperManager) GetAffiliatingApp(ctx sdk.Context, username linotypes.AccountKey) (linotypes.AccountKey, sdk.Error) {
	// username is app itself.
	if dm.DoesDeveloperExist(ctx, username) {
		return username, nil
	}
	// user's role.
	role, err := dm.storage.GetUserRole(ctx, username)
	if err != nil {
		return "", err
	}
	return role.AffiliatedApp, nil
}

// GetAffiliated returns all affiliated account of app.
func (dm DeveloperManager) GetAffiliated(ctx sdk.Context, app linotypes.AccountKey) []linotypes.AccountKey {
	if !dm.DoesDeveloperExist(ctx, app) {
		return nil
	}
	return dm.storage.GetAllAffiliatedAcc(ctx, app)
}

// UpdateAuthorization - update app's authorization on user.
func (dm DeveloperManager) UpdateIDAAuth(ctx sdk.Context, app, username linotypes.AccountKey, active bool) sdk.Error {
	// when developer is revoked, no need to update auth
	if !dm.DoesDeveloperExist(ctx, app) {
		return types.ErrDeveloperNotFound()
	}
	if !dm.acc.DoesAccountExist(ctx, username) {
		return types.ErrAccountNotFound()
	}
	if dm.storage.HasAffiliatedAcc(ctx, app, username) {
		return types.ErrInvalidIDAAuth()
	}
	bank := dm.storage.GetIDABank(ctx, app, username)
	if bank.Unauthed == !active {
		return types.ErrInvalidIDAAuth()
	}
	bank.Unauthed = !active
	dm.storage.SetIDABank(ctx, app, username, bank)
	return nil
}

// ReportConsumption - add consumption to a developer.
func (dm DeveloperManager) ReportConsumption(ctx sdk.Context, app linotypes.AccountKey, consumption linotypes.MiniDollar) sdk.Error {
	developer, err := dm.storage.GetDeveloper(ctx, app)
	if err != nil {
		return err
	}
	developer.AppConsumption = developer.AppConsumption.Plus(consumption)
	dm.storage.SetDeveloper(ctx, *developer)
	return nil
}

// DistributeDevInflation - distribute monthly app inflation.
func (dm DeveloperManager) DistributeDevInflation(ctx sdk.Context) sdk.Error {
	// No-op if there is no developer, leave inflations in pool.
	devs := dm.GetLiveDevelopers(ctx)
	if len(devs) == 0 {
		return nil
	}
	inflation, err := dm.acc.GetPool(ctx, linotypes.InflationDeveloperPool)
	if err != nil {
		return err
	}
	distSchema := make([]sdk.Dec, len(devs))
	totalConsumption := linotypes.NewMiniDollar(0)
	for _, dev := range devs {
		totalConsumption = totalConsumption.Plus(dev.AppConsumption)
	}
	if totalConsumption.IsZero() {
		// if not any consumption here, we evenly distribute all inflation
		for i := range devs {
			distSchema[i] = linotypes.NewDecFromRat(1, int64(len(devs)))
		}
	} else {
		for i, dev := range devs {
			distSchema[i] = dev.AppConsumption.ToDec().Quo(totalConsumption.ToDec())
		}
	}

	distributed := linotypes.NewCoinFromInt64(0)
	for i, developer := range devs {
		if i == (len(devs) - 1) {
			if err := dm.acc.MoveFromPool(ctx, linotypes.InflationDeveloperPool,
				linotypes.NewAccOrAddrFromAcc(developer.Username),
				inflation.Minus(distributed)); err != nil {
				return err
			}
			break
		}
		percentage := distSchema[i]
		myShareRat := inflation.ToDec().Mul(percentage)
		myShareCoin := linotypes.DecToCoin(myShareRat)
		distributed = distributed.Plus(myShareCoin)
		if err := dm.acc.MoveFromPool(ctx, linotypes.InflationDeveloperPool,
			linotypes.NewAccOrAddrFromAcc(developer.Username),
			myShareCoin); err != nil {
			return err
		}
	}

	dm.clearConsumption(ctx)
	return nil
}

func (dm DeveloperManager) clearConsumption(ctx sdk.Context) {
	devs := dm.GetLiveDevelopers(ctx)
	for _, dev := range devs {
		dev.AppConsumption = linotypes.NewMiniDollar(0)
		dm.storage.SetDeveloper(ctx, dev)
	}
}

// Permissions: will be removed in upgrade3.
// GrantPermission
func (dm DeveloperManager) GrantPermission(ctx sdk.Context, app, user linotypes.AccountKey, duration int64, level linotypes.Permission, amount linotypes.LNO) sdk.Error {
	if !dm.DoesDeveloperExist(ctx, app) {
		return types.ErrDeveloperNotFound()
	}
	if !dm.acc.DoesAccountExist(ctx, user) {
		return types.ErrAccountNotFound()
	}

	switch level {
	case linotypes.AppPermission:
		if err := dm.acc.AuthorizePermission(
			ctx, user, app, duration, level, linotypes.NewCoinFromInt64(0)); err != nil {
			return err
		}
	case linotypes.PreAuthorizationPermission:
		coin, err := linotypes.LinoToCoin(amount)
		if err != nil {
			return err
		}
		if err := dm.acc.AuthorizePermission(ctx, user, app, duration, level, coin); err != nil {
			return err
		}
	case linotypes.AppAndPreAuthorizationPermission:
		if err := dm.acc.AuthorizePermission(
			ctx, user, app, duration, linotypes.AppPermission, linotypes.NewCoinFromInt64(0)); err != nil {
			return err
		}
		coin, err := linotypes.LinoToCoin(amount)
		if err != nil {
			return err
		}
		if err := dm.acc.AuthorizePermission(
			ctx, user, app, duration, linotypes.PreAuthorizationPermission, coin); err != nil {
			return err
		}
	default:
		return types.ErrInvalidGrantPermission()
	}
	return nil
}

func (dm DeveloperManager) RevokePermission(ctx sdk.Context, user, app linotypes.AccountKey, perm linotypes.Permission) sdk.Error {
	if !dm.acc.DoesAccountExist(ctx, user) {
		return types.ErrAccountNotFound()
	}
	if err := dm.acc.RevokePermission(ctx, user, app, perm); err != nil {
		return err
	}
	return nil
}

func (dm DeveloperManager) GetReservePool(ctx sdk.Context) model.ReservePool {
	return *dm.storage.GetReservePool(ctx)
}

func (dm DeveloperManager) GetIDAStats(ctx sdk.Context, app linotypes.AccountKey) (model.AppIDAStats, sdk.Error) {
	if _, err := dm.validAppIDA(ctx, app); err != nil {
		return model.AppIDAStats{}, err
	}
	stats := *dm.storage.GetIDAStats(ctx, app)
	return stats, nil
}

func (dm DeveloperManager) ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	state := &model.DeveloperTablesIR{
		Version: exportVersion,
	}
	stores := dm.storage.StoreMap(ctx)

	// export developers
	stores[string(model.DeveloperSubstore)].Iterate(func(key []byte, val interface{}) bool {
		dev := val.(*model.Developer)
		state.Developers = append(state.Developers, model.DeveloperIR{
			Username:       dev.Username,
			AppConsumption: dev.AppConsumption,
			Website:        dev.Website,
			Description:    dev.Description,
			AppMetaData:    dev.AppMetaData,
			IsDeleted:      dev.IsDeleted,
			NAffiliated:    dev.NAffiliated,
		})
		return false
	})

	// export IDAs
	stores[string(model.IdaSubstore)].Iterate(func(key []byte, val interface{}) bool {
		ida := val.(*model.AppIDA)
		state.IDAs = append(state.IDAs, model.AppIDAIR(*ida))
		return false
	})

	// export ida balance
	stores[string(model.IdaBalanceSubstore)].Iterate(func(key []byte, val interface{}) bool {
		app, user := model.ParseIDABalanceKey(key)
		bank := val.(*model.IDABank)
		state.IDABanks = append(state.IDABanks, model.IDABankIR{
			App:      app,
			User:     user,
			Balance:  bank.Balance,
			Unauthed: bank.Unauthed,
		})
		return false
	})

	// export reserve pool
	stores[string(model.ReservePoolSubstore)].Iterate(func(key []byte, val interface{}) bool {
		pool := val.(*model.ReservePool)
		state.ReservePool = model.ReservePoolIR(*pool)
		return false
	})

	// export affiliated accounts
	stores[string(model.AffiliatedAccSubstore)].Iterate(func(key []byte, _ interface{}) bool {
		app, user := model.ParseAffiliatedAccKey(key)
		state.AffiliatedAccs = append(state.AffiliatedAccs, model.AffiliatedAccIR{
			App:  app,
			User: user,
		})
		return false
	})

	// export UserRoles
	stores[string(model.UserRoleSubstore)].Iterate(func(key []byte, val interface{}) bool {
		role := val.(*model.Role)
		state.UserRoles = append(state.UserRoles, model.UserRoleIR{
			User:          linotypes.AccountKey(key),
			AffiliatedApp: role.AffiliatedApp,
		})
		return false
	})

	// export IDA stats
	stores[string(model.IdaStatsSubstore)].Iterate(func(key []byte, val interface{}) bool {
		stats := val.(*model.AppIDAStats)
		state.IDAStats = append(state.IDAStats, model.IDAStatsIR{
			App:   linotypes.AccountKey(key),
			Total: stats.Total,
		})
		return false
	})

	return utils.Save(filepath, cdc, state)
}

// Import from file
func (dm DeveloperManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.DeveloperTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.DeveloperTablesIR)

	if table.Version != importVersion {
		return fmt.Errorf("unsupported import version: %d", table.Version)
	}

	// import developers
	for _, dev := range table.Developers {
		dm.storage.SetDeveloper(ctx, model.Developer{
			Username:       dev.Username,
			AppConsumption: dev.AppConsumption,
			Website:        dev.Website,
			Description:    dev.Description,
			AppMetaData:    dev.AppMetaData,
			IsDeleted:      dev.IsDeleted,
			NAffiliated:    dev.NAffiliated,
		})
	}

	// import IDAs
	for _, ida := range table.IDAs {
		dm.storage.SetIDA(ctx, model.AppIDA(ida))
	}

	// import IDABanks
	for _, bank := range table.IDABanks {
		dm.storage.SetIDABank(ctx, bank.App, bank.User, &model.IDABank{
			Balance:  bank.Balance,
			Unauthed: bank.Unauthed,
		})
	}

	// import reserve pool
	pool := model.ReservePool(table.ReservePool)
	dm.storage.SetReservePool(ctx, &pool)

	// import affiliated accounts
	for _, acc := range table.AffiliatedAccs {
		dm.storage.SetAffiliatedAcc(ctx, acc.App, acc.User)
	}

	// import user roles
	for _, role := range table.UserRoles {
		dm.storage.SetUserRole(ctx, role.User, &model.Role{
			AffiliatedApp: role.AffiliatedApp,
		})
	}

	// import ida stats
	for _, stat := range table.IDAStats {
		dm.storage.SetIDAStats(ctx, stat.App, model.AppIDAStats{
			Total: stat.Total,
		})
	}

	return nil
}
