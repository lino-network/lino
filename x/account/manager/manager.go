package manager

import (
	"bytes"
	"fmt"
	"reflect"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/account/model"
	"github.com/lino-network/lino/x/account/types"
)

const (
	nSecOfOneHour = 3600
	// HoursPerYear - as defined by a julian year of 365.25 days
	nHourOfOneYear = 8766

	exportVersion = 3
	importVersion = 3
)

// AccountManager - account manager
type AccountManager struct {
	storage     model.AccountStorage
	paramHolder param.ParamKeeper
}

// NewLinoAccount - new account manager
func NewAccountManager(key sdk.StoreKey, holder param.ParamKeeper) AccountManager {
	return AccountManager{
		storage:     model.NewAccountStorage(key),
		paramHolder: holder,
	}
}

func (am AccountManager) InitGenesis(ctx sdk.Context, total linotypes.Coin, pools []model.Pool) {
	neverInited := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				neverInited = true
			}
		}()
		am.storage.GetSupply(ctx)
	}()
	if !neverInited {
		panic("Account Module Already Inited")
	}
	am.storage.SetSupply(ctx, &model.Supply{
		LastYearTotal:     total,
		Total:             total,
		ChainStartTime:    ctx.BlockTime().Unix(),
		LastInflationTime: ctx.BlockTime().Unix(),
	})
	for _, pool := range pools {
		am.storage.SetPool(ctx, &pool)
	}
}

func (am AccountManager) GetPool(
	ctx sdk.Context, poolName linotypes.PoolName) (linotypes.Coin, sdk.Error) {
	pool, err := am.storage.GetPool(ctx, poolName)
	if err != nil {
		return linotypes.NewCoinFromInt64(0), err
	}
	return pool.Balance, nil
}

// MoveCoin move coins from the acc/addr to the acc/addr.
func (am AccountManager) MoveCoin(ctx sdk.Context, sender, receiver linotypes.AccOrAddr, coin linotypes.Coin) sdk.Error {
	if coin.IsNegative() {
		return types.ErrNegativeMoveAmount(coin)
	}
	err := am.minusCoin(ctx, sender, coin)
	if err != nil {
		return err
	}
	return am.addCoin(ctx, receiver, coin)
}

// MoveFromPool - move coin from pool to an address or user.
func (am AccountManager) MoveFromPool(ctx sdk.Context, poolName linotypes.PoolName, dest linotypes.AccOrAddr, amount linotypes.Coin) sdk.Error {
	if amount.IsNegative() {
		return types.ErrNegativeMoveAmount(amount)
	}

	pool, err := am.storage.GetPool(ctx, poolName)
	if err != nil {
		return err
	}
	if !pool.Balance.IsGTE(amount) {
		return types.ErrPoolNotEnough(poolName)
	}
	pool.Balance = pool.Balance.Minus(amount)
	am.storage.SetPool(ctx, pool)
	return am.addCoin(ctx, dest, amount)
}

// MoveToPool - move coin from an address or account to pool
func (am AccountManager) MoveToPool(ctx sdk.Context, poolName linotypes.PoolName, from linotypes.AccOrAddr, amount linotypes.Coin) sdk.Error {
	if amount.IsNegative() {
		return types.ErrNegativeMoveAmount(amount)
	}

	err := am.minusCoin(ctx, from, amount)
	if err != nil {
		return err
	}
	pool, err := am.storage.GetPool(ctx, poolName)
	if err != nil {
		return err
	}
	pool.Balance = pool.Balance.Plus(amount)
	am.storage.SetPool(ctx, pool)
	return nil
}

// MoveBetweenPools, move coin between pools.
func (am AccountManager) MoveBetweenPools(ctx sdk.Context, from, to linotypes.PoolName, amount linotypes.Coin) sdk.Error {
	if amount.IsNegative() {
		return types.ErrNegativeMoveAmount(amount)
	}

	fromPool, err := am.storage.GetPool(ctx, from)
	if err != nil {
		return err
	}

	toPool, err := am.storage.GetPool(ctx, to)
	if err != nil {
		return err
	}

	if !fromPool.Balance.IsGTE(amount) {
		return types.ErrPoolNotEnough(from)
	}

	fromPool.Balance = fromPool.Balance.Minus(amount)
	toPool.Balance = toPool.Balance.Plus(amount)
	am.storage.SetPool(ctx, fromPool)
	am.storage.SetPool(ctx, toPool)
	return nil
}

// Mint - distribute the inflation to pools hourly.
func (am AccountManager) Mint(ctx sdk.Context) sdk.Error {
	supply := am.storage.GetSupply(ctx)
	chainStartTime := supply.ChainStartTime
	lastInflation := supply.LastInflationTime
	blockTime := ctx.BlockTime().Unix()

	nLastInflation := (lastInflation - chainStartTime) / nSecOfOneHour
	nCurrent := (blockTime - chainStartTime) / nSecOfOneHour

	if nCurrent <= nLastInflation {
		return nil
	}

	// invarience
	// premise: lastInflation >= chainStartTime
	// nCurrent > nLastInflation ==> blocktime > lastInflation
	// after: lastInflation = blocktime > chainStartTime
	for nth := nLastInflation + 1; nth <= nCurrent; nth++ {
		// mint to pools
		err := am.hourlyMintOn(ctx, supply)
		if err != nil {
			return err
		}
		if nth%nHourOfOneYear == 0 {
			supply.LastYearTotal = supply.Total
		}
	}

	supply.LastInflationTime = blockTime
	am.storage.SetSupply(ctx, supply)
	return nil
}

// will change supply after mint and allocate to pools.
func (am AccountManager) hourlyMintOn(ctx sdk.Context, supply *model.Supply) sdk.Error {
	allocation := am.paramHolder.GetGlobalAllocationParam(ctx)
	growthRate := allocation.GlobalGrowthRate
	minted := linotypes.DecToCoin(
		supply.LastYearTotal.ToDec().Mul(growthRate).
			Mul(linotypes.NewDecFromRat(1, nHourOfOneYear)))
	contentCreator := linotypes.DecToCoin(minted.ToDec().Mul(allocation.ContentCreatorAllocation))
	validator := linotypes.DecToCoin(minted.ToDec().Mul(allocation.ValidatorAllocation))
	developer := minted.Minus(contentCreator).Minus(validator)

	if err := am.mintToPool(ctx, linotypes.InflationConsumptionPool, contentCreator); err != nil {
		return err
	}
	if err := am.mintToPool(ctx, linotypes.InflationValidatorPool, validator); err != nil {
		return err
	}
	if err := am.mintToPool(ctx, linotypes.InflationDeveloperPool, developer); err != nil {
		return err
	}
	supply.Total = supply.Total.Plus(minted)
	return nil
}

func (am AccountManager) mintToPool(ctx sdk.Context, poolName linotypes.PoolName, amount linotypes.Coin) sdk.Error {
	pool, err := am.storage.GetPool(ctx, poolName)
	if err != nil {
		return err
	}
	pool.Balance = pool.Balance.Plus(amount)
	am.storage.SetPool(ctx, pool)
	return nil
}

func (am AccountManager) addCoin(ctx sdk.Context, dest linotypes.AccOrAddr, amount linotypes.Coin) sdk.Error {
	if !dest.IsAddr {
		return am.addCoinToUsername(ctx, dest.AccountKey, amount)
	} else {
		am.addCoinToAddress(ctx, dest.Addr, amount)
	}
	return nil
}

// addCoinToUsername - add coin to address associated username
func (am AccountManager) addCoinToUsername(ctx sdk.Context, username linotypes.AccountKey, coin linotypes.Coin) sdk.Error {
	accInfo, err := am.storage.GetInfo(ctx, username)
	if err != nil {
		return types.ErrAccountNotFound(username)
	}
	am.addCoinToAddress(ctx, accInfo.Address, coin)
	return nil
}

// addCoinToAddress - add coin to address associated username
func (am AccountManager) addCoinToAddress(ctx sdk.Context, addr sdk.AccAddress, coin linotypes.Coin) {
	if coin.IsZero() {
		return
	}
	bank, err := am.storage.GetBank(ctx, addr)
	if err != nil {
		// if address is not created, created a new one
		bank = &model.AccountBank{
			Saving:  linotypes.NewCoinFromInt64(0),
			Pending: linotypes.NewCoinFromInt64(0),
		}
	}
	bank.Saving = bank.Saving.Plus(coin)
	am.storage.SetBank(ctx, addr, bank)
}

func (am AccountManager) minusCoin(ctx sdk.Context, from linotypes.AccOrAddr, amount linotypes.Coin) sdk.Error {
	if !from.IsAddr {
		return am.minusCoinFromUsername(ctx, from.AccountKey, amount)
	} else {
		return am.minusCoinFromAddress(ctx, from.Addr, amount)
	}
}

// minusSavingCoin - minus coin from balance, remove coin day in the tail
func (am AccountManager) minusCoinFromUsername(ctx sdk.Context, username linotypes.AccountKey, coin linotypes.Coin) sdk.Error {
	accInfo, err := am.storage.GetInfo(ctx, username)
	if err != nil {
		return types.ErrAccountNotFound(username)
	}
	return am.minusCoinFromAddress(ctx, accInfo.Address, coin)
}

// minusCoinFromAddress - minus coin from address
func (am AccountManager) minusCoinFromAddress(ctx sdk.Context, address sdk.AccAddress, coin linotypes.Coin) sdk.Error {
	if coin.IsZero() {
		return nil
	}
	bank, err := am.storage.GetBank(ctx, address)
	if err != nil {
		return err
	}

	bank.Saving = bank.Saving.Minus(coin)
	if !bank.Saving.IsGTE(am.paramHolder.GetAccountParam(ctx).MinimumBalance) {
		return types.ErrAccountSavingCoinNotEnough()
	}

	am.storage.SetBank(ctx, address, bank)
	return nil
}

func (am AccountManager) DoesAccountExist(ctx sdk.Context, username linotypes.AccountKey) bool {
	return am.storage.DoesAccountExist(ctx, username)
}

// RegisterAccount - register account, deduct fee from referrer address then create a new account
func (am AccountManager) RegisterAccount(ctx sdk.Context, referrer linotypes.AccOrAddr, registerFee linotypes.Coin, username linotypes.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error {
	minRegFee := am.paramHolder.GetAccountParam(ctx).RegisterFee
	if minRegFee.IsGT(registerFee) {
		return types.ErrRegisterFeeInsufficient()
	}

	if err := am.createAccount(ctx, username, signingKey, transactionKey); err != nil {
		return err
	}

	err := am.MoveToPool(ctx, linotypes.InflationValidatorPool, referrer, minRegFee)
	if err != nil {
		return err
	}

	err = am.MoveCoin(ctx,
		referrer, linotypes.NewAccOrAddrFromAcc(username), registerFee.Minus(minRegFee))
	return err
}

func (am AccountManager) GenesisAccount(ctx sdk.Context, username linotypes.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error {
	return am.createAccount(ctx, username, signingKey, transactionKey)
}

// CreateAccount - create account, caller should make sure the register fee is valid
func (am AccountManager) createAccount(ctx sdk.Context, username linotypes.AccountKey, signingKey, transactionKey crypto.PubKey) sdk.Error {
	if am.storage.DoesAccountExist(ctx, username) {
		return types.ErrAccountAlreadyExists(username)
	}

	// get address from tx key
	addr := sdk.AccAddress(transactionKey.Address())
	bank, err := am.storage.GetBank(ctx, addr)
	if err != nil {
		bank = &model.AccountBank{}
	}
	if bank.Username != "" {
		return types.ErrAddressAlreadyTaken(addr.String())
	}

	// set public key to bank
	bank.Username = username
	bank.PubKey = transactionKey
	am.storage.SetBank(ctx, addr, bank)

	accountInfo := &model.AccountInfo{
		Username:       username,
		CreatedAt:      ctx.BlockHeader().Time.Unix(),
		TransactionKey: transactionKey,
		SigningKey:     signingKey,
		Address:        addr,
	}
	am.storage.SetInfo(ctx, accountInfo)
	return nil
}

// UpdateJSONMeta - update user JONS meta data
func (accManager AccountManager) UpdateJSONMeta(
	ctx sdk.Context, username linotypes.AccountKey, jsonMeta string) sdk.Error {
	accountMeta := accManager.storage.GetMeta(ctx, username)
	accountMeta.JSONMeta = jsonMeta
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

// CheckSigningPubKeyOwner - given a public key, check if it is valid for the user.
func (accManager AccountManager) CheckSigningPubKeyOwner(ctx sdk.Context, me linotypes.AccountKey, signKey crypto.PubKey) (linotypes.AccountKey, sdk.Error) {
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

	return "", types.ErrCheckAuthenticatePubKeyOwner(me)
}

// CheckSigningPubKeyOwnerByAddress - given a public key, check if it is valid for address.
// If tx is already paid then bank can be created.
func (accManager AccountManager) CheckSigningPubKeyOwnerByAddress(
	ctx sdk.Context, address sdk.AccAddress, signKey crypto.PubKey, isPaid bool) sdk.Error {
	bank, err := accManager.storage.GetBank(ctx, address)
	if err != nil {
		if !isPaid || err.Code() != linotypes.CodeAccountBankNotFound {
			return err
		}
		bank = &model.AccountBank{}
	}

	if bank.PubKey == nil {
		if !bytes.Equal(signKey.Address(), address) {
			return sdk.ErrInvalidPubKey(
				fmt.Sprintf("PubKey does not match Signer address %s", address))
		}
		bank.PubKey = signKey
		accManager.storage.SetBank(ctx, address, bank)
	}
	//check signing key for all permissions
	if !reflect.DeepEqual(bank.PubKey, signKey) {
		return types.ErrCheckAuthenticatePubKeyAddress(address)
	}

	return nil
}

// RecoverAccount - reset two public key pairs
func (accManager AccountManager) RecoverAccount(
	ctx sdk.Context, username linotypes.AccountKey, newTransactionPubKey, newSigningKey crypto.PubKey) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}

	newAddr := sdk.AccAddress(newTransactionPubKey.Address())
	newBank, err := accManager.storage.GetBank(ctx, newAddr)
	if err != nil {
		if err.Code() != types.ErrAccountBankNotFound(newAddr).Code() {
			return err
		}
		newBank = &model.AccountBank{
			Saving:  linotypes.NewCoinFromInt64(0),
			Pending: linotypes.NewCoinFromInt64(0),
		}
	}
	if newBank.Username != "" {
		return types.ErrAddressAlreadyTaken(newAddr.String())
	}

	oldAddr := accInfo.Address
	oldBank, err := accManager.storage.GetBank(ctx, oldAddr)
	if err != nil {
		return err
	}

	newBank.Username = username
	oldBank.Username = ""

	newBank.Sequence += oldBank.Sequence
	newBank.Saving = newBank.Saving.Plus(oldBank.Saving)
	newBank.PubKey = newTransactionPubKey
	oldBank.Saving = linotypes.NewCoinFromInt64(0)

	accInfo.Address = newAddr
	accInfo.SigningKey = newSigningKey
	accInfo.TransactionKey = newTransactionPubKey

	newBank.Pending = newBank.Pending.Plus(oldBank.Pending)
	oldBank.Pending = linotypes.NewCoinFromInt64(0)

	accManager.storage.SetInfo(ctx, accInfo)
	accManager.storage.SetBank(ctx, newAddr, newBank)
	accManager.storage.SetBank(ctx, oldAddr, oldBank)
	return nil
}

// AddPending - record pending amount of a user.
func (accManager AccountManager) AddPending(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error {
	info, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}
	bank, err := accManager.storage.GetBank(ctx, info.Address)
	if err != nil {
		return err
	}
	bank.Pending = bank.Pending.Plus(amount)
	accManager.storage.SetBank(ctx, info.Address, bank)
	return nil
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

func (accManager AccountManager) GetBankByAddress(ctx sdk.Context, addr sdk.AccAddress) (*model.AccountBank, sdk.Error) {
	return accManager.storage.GetBank(ctx, addr)
}

func (accManager AccountManager) GetMeta(ctx sdk.Context, username linotypes.AccountKey) (*model.AccountMeta, sdk.Error) {
	return accManager.storage.GetMeta(ctx, username), nil
}

func (accManager AccountManager) GetSupply(ctx sdk.Context) model.Supply {
	return *accManager.storage.GetSupply(ctx)
}

// ExportToFile -
func (am AccountManager) ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	state := &model.AccountTablesIR{
		Version: exportVersion,
	}
	substores := am.storage.PartialStoreMap(ctx)

	// export accounts
	substores[string(model.AccountInfoSubstore)].Iterate(func(key []byte, val interface{}) bool {
		acc := val.(*model.AccountInfo)
		state.Accounts = append(state.Accounts, model.AccountIR(*acc))
		return false
	})

	// export banks
	substores[string(model.AccountBankSubstore)].Iterate(func(key []byte, val interface{}) bool {
		bank := val.(*model.AccountBank)
		addr := key
		state.Banks = append(state.Banks, model.AccountBankIR{
			Address:  addr,
			Saving:   bank.Saving,
			Pending:  bank.Pending,
			PubKey:   bank.PubKey,
			Sequence: bank.Sequence,
			Username: bank.Username,
		})
		return false
	})

	// export metas
	substores[string(model.AccountMetaSubstore)].Iterate(func(key []byte, val interface{}) bool {
		meta := val.(*model.AccountMeta)
		acc := linotypes.AccountKey(key)
		state.Metas = append(state.Metas, model.AccountMetaIR{
			Username: acc,
			JSONMeta: meta.JSONMeta,
		})
		return false
	})

	// pools
	substores[string(model.AccountPoolSubstore)].Iterate(func(key []byte, val interface{}) bool {
		pool := val.(*model.Pool)
		state.Pools = append(state.Pools, model.PoolIR(*pool))
		return false
	})

	// supply
	state.Supply = model.SupplyIR(*am.storage.GetSupply(ctx))

	return utils.Save(filepath, cdc, state)
}

// ImportFromFile import state from file.
func (am AccountManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.AccountTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.AccountTablesIR)

	if table.Version != importVersion {
		return fmt.Errorf("unsupported import version: %d", table.Version)
	}

	banks := make(map[string]int)

	// import accounts.
	for _, v := range table.Accounts {
		info := model.AccountInfo(v)
		if _, err := am.storage.GetInfo(ctx, v.Username); err != nil {
			am.storage.SetInfo(ctx, &info)
			if banks[string(v.Address)] != 0 {
				panic(fmt.Errorf("used address: %s", v.Address))
			}
			banks[string(v.Address)] = 1
		} else {
			panic(fmt.Errorf("duplicated username: %s", v.Username))
		}
	}

	// import banks
	for _, v := range table.Banks {
		bank := model.AccountBank{
			Saving:   v.Saving,
			Pending:  v.Pending,
			PubKey:   v.PubKey,
			Sequence: v.Sequence,
			Username: v.Username,
		}
		if banks[string(v.Address)] > 1 {
			panic(fmt.Errorf("duplicated address: %+v", v))
		}
		banks[string(v.Address)] = 2
		am.storage.SetBank(ctx, sdk.AccAddress(v.Address), &bank)
	}

	// import meta
	for _, meta := range table.Metas {
		am.storage.SetMeta(ctx, meta.Username, &model.AccountMeta{
			JSONMeta: meta.JSONMeta,
		})
	}

	// import pools
	for _, poolir := range table.Pools {
		am.storage.SetPool(ctx, (*model.Pool)(&poolir))
	}

	// import supply
	am.storage.SetSupply(ctx, (*model.Supply)(&table.Supply))

	return nil
}
