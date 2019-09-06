package manager

import (
	"fmt"
	"reflect"
	"time"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/account/model"
	"github.com/lino-network/lino/x/account/types"
	"github.com/lino-network/lino/x/global"
)

// AccountManager - account manager
type AccountManager struct {
	storage     model.AccountStorage
	paramHolder param.ParamKeeper
	gm          global.GlobalKeeper
}

// NewLinoAccount - new account manager
func NewAccountManager(key sdk.StoreKey, holder param.ParamKeeper, gm global.GlobalKeeper) AccountManager {
	return AccountManager{
		storage:     model.NewAccountStorage(key),
		paramHolder: holder,
		gm:          gm,
	}
}

func (accManager AccountManager) DoesAccountExist(ctx sdk.Context, username linotypes.AccountKey) bool {
	return accManager.storage.DoesAccountExist(ctx, username)
}

// RegisterAccount - register account, deduct fee from referrer address then create a new account
func (accManager AccountManager) RegisterAccount(
	ctx sdk.Context, referrerAddr sdk.AccAddress, registerFee linotypes.Coin,
	username linotypes.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error {
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	if accParams.RegisterFee.IsGT(registerFee) {
		return types.ErrRegisterFeeInsufficient()
	}

	if err := accManager.MinusCoinFromAddress(ctx, referrerAddr, registerFee); err != nil {
		return err
	}
	if err := accManager.CreateAccount(ctx, username, signingKey, transactionKey); err != nil {
		return err
	}
	if err := accManager.gm.AddToValidatorInflationPool(ctx, accParams.RegisterFee); err != nil {
		return err
	}
	return accManager.AddCoinToUsername(ctx, username, registerFee.Minus(accParams.RegisterFee))
}

// CreateAccount - create account, caller should make sure the register fee is valid
func (accManager AccountManager) CreateAccount(
	ctx sdk.Context, username linotypes.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error {
	if accManager.storage.DoesAccountExist(ctx, username) {
		return types.ErrAccountAlreadyExists(username)
	}

	// get address from tx key
	addr := sdk.AccAddress(transactionKey.Address().Bytes())
	bank, err := accManager.storage.GetBank(ctx, addr)
	if err != nil {
		if err.Code() != model.ErrAccountBankNotFound().Code() {
			return err
		}
		bank = &model.AccountBank{}
	}
	if bank.Username != "" {
		return types.ErrAddressAlreadyTaken(addr)
	}

	// set public key to bank
	bank.Username = username
	bank.PubKey = transactionKey

	accManager.storage.SetBank(ctx, addr, bank)

	accountInfo := &model.AccountInfo{
		Username:       username,
		CreatedAt:      ctx.BlockHeader().Time.Unix(),
		TransactionKey: transactionKey,
		SigningKey:     signingKey,
		Address:        addr,
	}
	accManager.storage.SetInfo(ctx, username, accountInfo)
	accountMeta := &model.AccountMeta{}
	accManager.storage.SetMeta(ctx, username, accountMeta)
	return nil
}

// MoveCoinFromUsernameToUsername - move coin from sender to receiver
func (accManager AccountManager) MoveCoinFromUsernameToUsername(ctx sdk.Context, sender, receiver linotypes.AccountKey, coin linotypes.Coin) sdk.Error {
	if err := accManager.MinusCoinFromUsername(ctx, sender, coin); err != nil {
		return err
	}
	if err := accManager.AddCoinToUsername(ctx, receiver, coin); err != nil {
		return err
	}
	return nil
}

// AddCoinToUsername - add coin to address associated username
func (accManager AccountManager) AddCoinToUsername(ctx sdk.Context, username linotypes.AccountKey, coin linotypes.Coin) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}
	return accManager.AddCoinToAddress(ctx, accInfo.Address, coin)
}

// AddCoinToAddress - add coin to address associated username
func (accManager AccountManager) AddCoinToAddress(ctx sdk.Context, addr sdk.Address, coin linotypes.Coin) sdk.Error {
	bank, err := accManager.storage.GetBank(ctx, addr)
	if err != nil {
		if err.Code() != model.ErrAccountBankNotFound().Code() {
			return err
		}
		// if address is not created, created a new one
		bank = &model.AccountBank{
			Saving: linotypes.NewCoinFromInt64(0),
		}
	}
	bank.Saving = bank.Saving.Plus(coin)

	accManager.storage.SetBank(ctx, addr, bank)
	return nil
}

// MinusSavingCoin - minus coin from balance, remove coin day in the tail
func (accManager AccountManager) MinusCoinFromUsername(ctx sdk.Context, username linotypes.AccountKey, coin linotypes.Coin) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		if err.Code() == model.ErrAccountInfoNotFound().Code() {
			return types.ErrAccountNotFound(username)
		}
		return err
	}

	return accManager.MinusCoinFromAddress(ctx, accInfo.Address, coin)
}

// MinusCoinFromAddress - minus coin from address
func (accManager AccountManager) MinusCoinFromAddress(ctx sdk.Context, address sdk.Address, coin linotypes.Coin) sdk.Error {
	if coin.IsZero() {
		return nil
	}
	bank, err := accManager.storage.GetBank(ctx, address)
	if err != nil {
		return err
	}

	accountParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}

	bank.Saving = bank.Saving.Minus(coin)
	if !bank.Saving.IsGTE(accountParams.MinimumBalance) {
		return types.ErrAccountSavingCoinNotEnough()
	}

	accManager.storage.SetBank(ctx, address, bank)
	return nil
}

// UpdateJSONMeta - update user JONS meta data
func (accManager AccountManager) UpdateJSONMeta(
	ctx sdk.Context, username linotypes.AccountKey, JSONMeta string) sdk.Error {
	accountMeta := accManager.storage.GetMeta(ctx, username)
	accountMeta.JSONMeta = JSONMeta
	accManager.storage.SetMeta(ctx, username, accountMeta)
	return nil
}

// GetTransactionKey - get transaction public key
func (accManager AccountManager) GetTransactionKey(
	ctx sdk.Context, username linotypes.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, types.ErrGetTransactionKey(username)
	}
	return accountInfo.TransactionKey, nil
}

// GetAppKey - get app public key
func (accManager AccountManager) GetSigningKey(
	ctx sdk.Context, username linotypes.AccountKey) (crypto.PubKey, sdk.Error) {
	info, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, types.ErrGetSigningKey(username)
	}
	return info.SigningKey, nil
}

// GetSavingFromUsername - get user balance
func (accManager AccountManager) GetSavingFromUsername(ctx sdk.Context, username linotypes.AccountKey) (linotypes.Coin, sdk.Error) {
	info, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return linotypes.Coin{}, types.ErrGetSavingFromBank(err)
	}
	bank, err := accManager.storage.GetBank(ctx, info.Address)
	if err != nil {
		return linotypes.Coin{}, types.ErrGetSavingFromBank(err)
	}
	return bank.Saving, nil
}

// GetSavingFromBank - get user balance
func (accManager AccountManager) GetSavingFromAddress(ctx sdk.Context, address sdk.Address) (linotypes.Coin, sdk.Error) {
	bank, err := accManager.storage.GetBank(ctx, address)
	if err != nil {
		return linotypes.Coin{}, types.ErrGetSavingFromBank(err)
	}
	return bank.Saving, nil
}

// GetSequence - get user sequence number
func (accManager AccountManager) GetSequence(ctx sdk.Context, address sdk.Address) (uint64, sdk.Error) {
	bank, err := accManager.storage.GetBank(ctx, address)
	if err != nil {
		return 0, types.ErrGetSequence(err)
	}
	return bank.Sequence, nil
}

// GetAddress - get user bank address
func (accManager AccountManager) GetAddress(ctx sdk.Context, username linotypes.AccountKey) (sdk.AccAddress, sdk.Error) {
	info, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, types.ErrGetAddress(err)
	}
	return info.Address, nil
}

// GetFrozenMoneyList - get user frozen money list
func (accManager AccountManager) GetFrozenMoneyList(
	ctx sdk.Context, addr sdk.Address) ([]model.FrozenMoney, sdk.Error) {
	bank, err := accManager.storage.GetBank(ctx, addr)
	if err != nil {
		return nil, types.ErrGetFrozenMoneyList(err)
	}
	return bank.FrozenMoneyList, nil
}

// IncreaseSequenceByOne - increase user sequence number by one
func (accManager AccountManager) IncreaseSequenceByOne(ctx sdk.Context, address sdk.Address) sdk.Error {
	bank, err := accManager.storage.GetBank(ctx, address)
	if err != nil {
		return types.ErrIncreaseSequenceByOne(err)
	}
	bank.Sequence++
	accManager.storage.SetBank(ctx, address, bank)
	return nil
}

// AuthorizePermission - userA authorize permission to userB (currently only support auth to a developer)
func (accManager AccountManager) AuthorizePermission(
	ctx sdk.Context, me linotypes.AccountKey, grantTo linotypes.AccountKey,
	validityPeriod int64, grantLevel linotypes.Permission, amount linotypes.Coin) sdk.Error {
	if !accManager.storage.DoesAccountExist(ctx, grantTo) {
		return types.ErrAccountNotFound(grantTo)
	}
	if grantLevel != linotypes.PreAuthorizationPermission && grantLevel != linotypes.AppPermission {
		return types.ErrUnsupportGrantLevel()
	}
	newGrantPubKey := model.GrantPermission{
		GrantTo:    grantTo,
		Permission: grantLevel,
		CreatedAt:  ctx.BlockHeader().Time.Unix(),
		ExpiresAt:  ctx.BlockHeader().Time.Add(time.Duration(validityPeriod) * time.Second).Unix(),
		Amount:     amount,
	}
	pubkeys, err := accManager.storage.GetGrantPermissions(ctx, me, grantTo)
	if err != nil {
		// if grant permission list is empty, create a new one
		if err.Code() == model.ErrGrantPubKeyNotFound().Code() {
			accManager.storage.SetGrantPermissions(
				ctx, me, grantTo, []*model.GrantPermission{&newGrantPubKey})
			return nil
		}
		return err
	}

	// iterate grant public key list
	for i, pubkey := range pubkeys {
		if pubkey.Permission == grantLevel {
			pubkeys[i] = &newGrantPubKey
			accManager.storage.SetGrantPermissions(ctx, me, grantTo, pubkeys)
			return nil
		}
	}
	// If grant permission doesn't have record in store, add to grant public key list
	pubkeys = append(pubkeys, &newGrantPubKey)
	accManager.storage.SetGrantPermissions(ctx, me, grantTo, pubkeys)
	return nil
}

// RevokePermission - revoke permission from a developer
func (accManager AccountManager) RevokePermission(
	ctx sdk.Context, me linotypes.AccountKey, grantTo linotypes.AccountKey, permission linotypes.Permission) sdk.Error {
	pubkeys, err := accManager.storage.GetGrantPermissions(ctx, me, grantTo)
	if err != nil {
		return err
	}

	// iterate grant public key list
	for i, pubkey := range pubkeys {
		if pubkey.Permission == permission {
			if len(pubkeys) == 1 {
				accManager.storage.DeleteAllGrantPermissions(ctx, me, grantTo)
				return nil
			}
			accManager.storage.SetGrantPermissions(ctx, me, grantTo, append(pubkeys[:i], pubkeys[i+1:]...))
			return nil
		}
	}
	return model.ErrGrantPubKeyNotFound()
}

// CheckSigningPubKeyOwner - given a public key, check if it is valid for given permission
func (accManager AccountManager) CheckSigningPubKeyOwner(
	ctx sdk.Context, me linotypes.AccountKey, signKey crypto.PubKey,
	permission linotypes.Permission, amount linotypes.Coin) (linotypes.AccountKey, sdk.Error) {
	accInfo, err := accManager.storage.GetInfo(ctx, me)
	if err != nil {
		return "", err
	}
	//check signing key for all permissions
	if reflect.DeepEqual(accInfo.SigningKey, signKey) {
		return me, nil
	}

	// otherwise check tx key
	if reflect.DeepEqual(accInfo.TransactionKey, signKey) {
		return me, nil
	}

	// if user doesn't use his own key, check his grant user pubkey
	grantPubKeys, err := accManager.storage.GetAllGrantPermissions(ctx, me)
	if err != nil {
		return "", err
	}

	for _, pubKey := range grantPubKeys {
		if pubKey.ExpiresAt < ctx.BlockHeader().Time.Unix() {
			continue
		}
		if permission != pubKey.Permission {
			continue
		}
		signingKey, err := accManager.GetSigningKey(ctx, pubKey.GrantTo)
		if err != nil {
			return "", err
		}
		if !reflect.DeepEqual(signKey, signingKey) {
			// check tx key instead
			txKey, err := accManager.GetTransactionKey(ctx, pubKey.GrantTo)
			if err != nil {
				return "", err
			}
			if !reflect.DeepEqual(signKey, txKey) {
				return "", types.ErrCheckAuthenticatePubKeyOwner(me)
			}
		}
		if permission == linotypes.PreAuthorizationPermission {
			if amount.IsGT(pubKey.Amount) {
				return "", types.ErrPreAuthAmountInsufficient(pubKey.GrantTo, pubKey.Amount, amount)
			}
			// override previous grant public key
			if err := accManager.AuthorizePermission(ctx, me, pubKey.GrantTo, pubKey.ExpiresAt-ctx.BlockHeader().Time.Unix(), pubKey.Permission, pubKey.Amount.Minus(amount)); err != nil {
				return "", nil
			}
		}
		return pubKey.GrantTo, nil
	}
	return "", types.ErrCheckAuthenticatePubKeyOwner(me)
}

// RecoverAccount - reset two public key pairs
func (accManager AccountManager) RecoverAccount(
	ctx sdk.Context, username linotypes.AccountKey, newTransactionPubKey, newSigningKey crypto.PubKey) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}

	newAddr := sdk.AccAddress(newTransactionPubKey.Address().Bytes())
	newBank, err := accManager.storage.GetBank(ctx, newAddr)
	if err != nil {
		if err.Code() != model.ErrAccountBankNotFound().Code() {
			return err
		}
		newBank = &model.AccountBank{}
	}
	if newBank.Username != "" {
		return types.ErrAddressAlreadyTaken(newAddr)
	}

	oldAddr := accInfo.Address
	oldBank, err := accManager.storage.GetBank(ctx, oldAddr)
	if err != nil {
		return err
	}

	newBank.Username = username
	oldBank.Username = ""

	newBank.Sequence = newBank.Sequence + oldBank.Sequence
	newBank.Saving = newBank.Saving.Plus(oldBank.Saving)
	oldBank.Saving = linotypes.NewCoinFromInt64(0)

	accInfo.Address = newAddr
	accInfo.SigningKey = newSigningKey
	accInfo.TransactionKey = newTransactionPubKey

	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	newBank.FrozenMoneyList = append(newBank.FrozenMoneyList, oldBank.FrozenMoneyList...)
	if int64(len(newBank.FrozenMoneyList)) >= accParams.MaxNumFrozenMoney {
		return types.ErrFrozenMoneyListTooLong()
	}

	oldBank.FrozenMoneyList = nil

	accManager.storage.SetInfo(ctx, username, accInfo)
	accManager.storage.SetBank(ctx, newAddr, newBank)
	accManager.storage.SetBank(ctx, oldAddr, oldBank)
	return nil
}

// AddFrozenMoney - add frozen money to user's frozen money list
func (accManager AccountManager) AddFrozenMoney(
	ctx sdk.Context, username linotypes.AccountKey,
	amount linotypes.Coin, start, interval, times int64) sdk.Error {
	info, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}
	accountBank, err := accManager.storage.GetBank(ctx, info.Address)
	if err != nil {
		return err
	}
	accManager.cleanExpiredFrozenMoney(ctx, accountBank)
	frozenMoney := model.FrozenMoney{
		Amount:   amount,
		StartAt:  start,
		Interval: interval,
		Times:    times,
	}

	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}

	if int64(len(accountBank.FrozenMoneyList)) >= accParams.MaxNumFrozenMoney {
		return types.ErrFrozenMoneyListTooLong()
	}

	accountBank.FrozenMoneyList = append(accountBank.FrozenMoneyList, frozenMoney)
	accManager.storage.SetBank(ctx, info.Address, accountBank)
	return nil
}

func (accManager AccountManager) cleanExpiredFrozenMoney(ctx sdk.Context, bank *model.AccountBank) {
	idx := 0
	for idx < len(bank.FrozenMoneyList) {
		frozenMoney := bank.FrozenMoneyList[idx]
		if ctx.BlockHeader().Time.Unix() > frozenMoney.StartAt+frozenMoney.Interval*frozenMoney.Times {
			bank.FrozenMoneyList = append(bank.FrozenMoneyList[:idx], bank.FrozenMoneyList[idx+1:]...)
			continue
		}

		idx++
	}
}

// getter
func (accManager AccountManager) GetInfo(ctx sdk.Context, username linotypes.AccountKey) (*model.AccountInfo, sdk.Error) {
	return accManager.storage.GetInfo(ctx, username)
}

func (accManager AccountManager) GetBank(ctx sdk.Context, username linotypes.AccountKey) (*model.AccountBank, sdk.Error) {
	info, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, err
	}
	return accManager.storage.GetBank(ctx, info.Address)
}

func (accManager AccountManager) GetMeta(ctx sdk.Context, username linotypes.AccountKey) (*model.AccountMeta, sdk.Error) {
	return accManager.storage.GetMeta(ctx, username), nil
}

func (accManager AccountManager) GetGrantPubKeys(ctx sdk.Context, username, grantTo linotypes.AccountKey) ([]*model.GrantPermission, sdk.Error) {
	return accManager.storage.GetGrantPermissions(ctx, username, grantTo)
}

func (accManager AccountManager) GetAllGrantPubKeys(ctx sdk.Context, username linotypes.AccountKey) ([]*model.GrantPermission, sdk.Error) {
	return accManager.storage.GetAllGrantPermissions(ctx, username)
}

// Export -
func (accManager AccountManager) ExportToFile(ctx sdk.Context, filepath string) error {
	panic("export account unimplemented")
}

// ImportFromFile import state from file.
func (accManager AccountManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.AccountTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.AccountTablesIR)
	ctx.Logger().Info(fmt.Sprintf("%s state parsed", filepath))
	// import accounts.
	for _, v := range table.Accounts {
		pubkey := v.Info.ResetKey
		addr := sdk.AccAddress(pubkey.Address())
		if _, err := accManager.storage.GetBank(ctx, addr); err == nil {
			addr = sdk.AccAddress(v.Username)
			// bank already exists, account bank use safe address. User must use
			// their signing key to reset.
			ctx.Logger().Error(
				fmt.Sprintf("bank account already exits, %v; use safe addr: %s", v, addr))
		}
		info := &model.AccountInfo{
			Username:       v.Username,
			CreatedAt:      v.Info.CreatedAt,
			SigningKey:     v.Info.TransactionKey,
			TransactionKey: pubkey,
			Address:        addr,
		}
		totalSaving := v.Bank.Saving
		totalSaving = totalSaving.Plus(v.Reward.UnclaimReward)
		bank := &model.AccountBank{
			Saving:          totalSaving,
			FrozenMoneyList: v.Bank.FrozenMoneyList,
			PubKey:          pubkey,
			Sequence:        v.Meta.Sequence,
			Username:        v.Username,
		}
		accManager.storage.SetInfo(ctx, v.Username, info)
		accManager.storage.SetBank(ctx, addr, bank)
		if v.Meta.JSONMeta != "" {
			meta := &model.AccountMeta{
				JSONMeta: v.Meta.JSONMeta,
			}
			accManager.storage.SetMeta(ctx, v.Username, meta)
		}
	}

	// import grant permissions.
	for _, v := range table.AccountGrantPubKeys {
		grant := v.GrantPubKey
		remainingTime := grant.ExpiresAt - ctx.BlockHeader().Time.Unix()
		if remainingTime > 0 {
			err := accManager.AuthorizePermission(ctx, v.Username, grant.Username,
				remainingTime, grant.Permission, grant.Amount)
			if err != nil {
				panic(err)
			}
		}
	}
	ctx.Logger().Info(fmt.Sprintf("%s state imported", filepath))
	return nil
}

// IterateAccounts - iterate accounts in KVStore
func (accManager AccountManager) IterateAccounts(ctx sdk.Context, process func(model.AccountInfo, model.AccountBank) (stop bool)) {
	accManager.storage.IterateAccounts(ctx, process)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
