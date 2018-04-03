package app

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/abci/types"
	oldwire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/auth"
	"github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/tx/register"
	val "github.com/lino-network/lino/tx/validator"
	"github.com/lino-network/lino/types"
)

const (
	appName = "LinoBlockchain"
)

// Extended ABCI application
type LinoBlockchain struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	capKeyAccountStore *sdk.KVStoreKey
	capKeyPostStore    *sdk.KVStoreKey
	capKeyValStore     *sdk.KVStoreKey
	capKeyIBCStore     *sdk.KVStoreKey
	capKeyGlobalStore  *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountManager acc.AccountManager
	postManager    post.PostManager
	valManager     val.ValidatorManager
	globalManager  global.GlobalManager

	preValidators []val.Validator
}

func NewLinoBlockchain(logger log.Logger, dbs map[string]dbm.DB) *LinoBlockchain {
	// create your application object
	var lb = &LinoBlockchain{
		BaseApp:            bam.NewBaseApp(appName, logger, dbs["acc"]),
		cdc:                MakeCodec(),
		capKeyAccountStore: sdk.NewKVStoreKey(types.AccountKVStoreKey),
		capKeyPostStore:    sdk.NewKVStoreKey(types.PostKVStoreKey),
		capKeyValStore:     sdk.NewKVStoreKey(types.ValidatorKVStoreKey),
		capKeyGlobalStore:  sdk.NewKVStoreKey(types.GlobalKVStoreKey),
		capKeyIBCStore:     sdk.NewKVStoreKey("ibc"),
	}
	lb.accountManager = acc.NewLinoAccountManager(lb.capKeyAccountStore)
	lb.postManager = post.NewPostMananger(lb.capKeyPostStore)
	lb.valManager = val.NewValidatorMananger(lb.capKeyValStore)
	lb.globalManager = global.NewGlobalManager(lb.capKeyGlobalStore)

	lb.Router().
		AddRoute(types.RegisterRouterName, register.NewHandler(lb.accountManager)).
		AddRoute(types.AccountRouterName, acc.NewHandler(lb.accountManager)).
		AddRoute(types.PostRouterName, post.NewHandler(lb.postManager, lb.accountManager, lb.globalManager)).
		AddRoute(types.ValidatorRouterName, val.NewHandler(lb.valManager, lb.accountManager))

	lb.SetTxDecoder(lb.txDecoder)
	lb.SetInitChainer(lb.initChainer)
	lb.SetEndBlocker(lb.endBlocker)
	// TODO(Cosmos): mounting multiple stores is broken
	// https://github.com/cosmos/cosmos-sdk/issues/532

	lb.MountStoreWithDB(lb.capKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	lb.MountStoreWithDB(lb.capKeyPostStore, sdk.StoreTypeIAVL, dbs["post"])
	lb.MountStoreWithDB(lb.capKeyValStore, sdk.StoreTypeIAVL, dbs["val"])
	lb.MountStoreWithDB(lb.capKeyGlobalStore, sdk.StoreTypeIAVL, dbs["global"])
	lb.SetAnteHandler(auth.NewAnteHandler(lb.accountManager))
	if err := lb.LoadLatestVersion(lb.capKeyAccountStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyPostStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyValStore); err != nil {
		cmn.Exit(err.Error())
	}
	if err := lb.LoadLatestVersion(lb.capKeyGlobalStore); err != nil {
		cmn.Exit(err.Error())
	}
	lb.preValidators = []val.Validator{}
	return lb
}

// custom tx codec
// TODO: use new go-wire
func MakeCodec() *wire.Codec {
	const msgTypeRegister = 0x1
	const msgTypeFollow = 0x2
	const msgTypeUnfollow = 0x3
	const msgTypeTransfer = 0x4
	const msgTypePost = 0x5
	const msgTypeLike = 0x6
	const msgTypeDonate = 0x7
	const msgTypeValidatorDeposit = 0x8
	const msgTypeValidatorWithdraw = 0x9
	const msgTypeValidatorRevoke = 0x10
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{register.RegisterMsg{}, msgTypeRegister},
		oldwire.ConcreteType{acc.FollowMsg{}, msgTypeFollow},
		oldwire.ConcreteType{acc.UnfollowMsg{}, msgTypeUnfollow},
		oldwire.ConcreteType{acc.TransferMsg{}, msgTypeTransfer},
		oldwire.ConcreteType{post.CreatePostMsg{}, msgTypePost},
		oldwire.ConcreteType{post.LikeMsg{}, msgTypeLike},
		oldwire.ConcreteType{post.DonateMsg{}, msgTypeDonate},
		oldwire.ConcreteType{val.ValidatorDepositMsg{}, msgTypeValidatorDeposit},
		oldwire.ConcreteType{val.ValidatorWithdrawMsg{}, msgTypeValidatorWithdraw},
		oldwire.ConcreteType{val.ValidatorRevokeMsg{}, msgTypeValidatorRevoke},
	)

	// TODO(Lino): Register msg type and model.
	cdc := wire.NewCodec()

	return cdc
}

// custom logic for transaction decoding
func (lb *LinoBlockchain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx{}

	// StdTx.Msg is an interface. The concrete types
	// are registered by MakeTxCodec in bank.RegisterWire.
	err := lb.cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}

// custom logic for basecoin initialization
func (lb *LinoBlockchain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	genesisState := new(genesis.GenesisState)
	//err := oldwire.UnmarshalJSON(stateJSON, genesisState)
	if err := json.Unmarshal(stateJSON, genesisState); err != nil {
		panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
	}

	if err := lb.valManager.InitGenesis(ctx); err != nil {
		panic(err)
	}
	if err := lb.globalManager.InitGlobalState(ctx, genesisState.GlobalState); err != nil {
		panic(err)
	}
	for _, gacc := range genesisState.Accounts {
		if err := lb.toAppAccount(ctx, gacc); err != nil {
			panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		// lb.accountMapper.SetAccount(ctx, acc)
	}
	return abci.ResponseInitChain{}
}

// convert GenesisAccount to AppAccount
func (lb *LinoBlockchain) toAppAccount(ctx sdk.Context, ga *genesis.GenesisAccount) sdk.Error {
	// send coins using address (even no account bank associated with this addr)
	bank, err := lb.accountManager.GetBankFromAddress(ctx, ga.PubKey.Address())
	if err == nil {
		// account bank exists
		panic(sdk.ErrGenesisParse("genesis bank already exist"))
	} else {
		fmt.Println(ga)
		coin, err := types.LinoToCoin(types.LNO(sdk.NewRat(ga.Lino)))
		if err != nil {
			panic(err)
		}
		// account bank not found, create a new one for this address
		bank = &acc.AccountBank{
			Address: ga.PubKey.Address(),
			Balance: coin,
		}
		if setErr := lb.accountManager.SetBankFromAddress(ctx, bank.Address, bank); setErr != nil {
			panic(sdk.ErrGenesisParse("set genesis bank failed"))
		}
		account := acc.NewProxyAccount(acc.AccountKey(ga.Name), &lb.accountManager)
		if account.IsAccountExist(ctx) {
			panic(sdk.ErrGenesisParse("genesis account already exist"))
		}
		if err := account.CreateAccount(ctx, acc.AccountKey(ga.Name), ga.PubKey, bank); err != nil {
			panic(err)
		}

		deposit := types.Coin{1000 * types.Decimals}
		// withdraw money from validator's bank
		if err := account.MinusCoin(ctx, deposit); err != nil {
			panic(err)
		}

		if addErr := lb.valManager.RegisterValidator(ctx, account.GetUsername(ctx), ga.ValPubKey.Bytes(), deposit); addErr != nil {
			panic(addErr)
		}
		if joinErr := lb.valManager.TryBecomeOncallValidator(ctx, account.GetUsername(ctx)); joinErr != nil {
			panic(joinErr)
		}
		if err := account.Apply(ctx); err != nil {
			panic(err)
		}
	}
	return nil
}

func (lb *LinoBlockchain) beginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	var err sdk.Error
	lb.preValidators, err = lb.valManager.GetOncallValList(ctx)
	if err != nil {
		panic(err)
	}
	absentValidators := req.GetAbsentValidators()
	if absentValidators != nil {
		// TODO Err handling
		lb.valManager.UpdateAbsentValidator(ctx, absentValidators)
	}
	// TODO Err handling
	lb.valManager.FireIncompetentValidator(ctx, req.GetByzantineValidators())
	return abci.ResponseBeginBlock{}
}

func (lb *LinoBlockchain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	curOncallList, err := lb.valManager.GetOncallValList(ctx)
	if err != nil {
		panic(err)
	}
	ABCIValList := []abci.Validator{}
	for _, preValidator := range lb.preValidators {
		if FindValidatorInList(preValidator, curOncallList) == -1 {
			preValidator.ABCIValidator.Power = 0
			ABCIValList = append(ABCIValList, preValidator.ABCIValidator)
		}
	}
	for _, validator := range curOncallList {
		ABCIValList = append(ABCIValList, validator.ABCIValidator)
	}
	return abci.ResponseEndBlock{ValidatorUpdates: ABCIValList}
}

func FindValidatorInList(validator val.Validator, validatorList []val.Validator) int {
	for i, curValidator := range validatorList {
		if validator.Username == curValidator.Username {
			return i
		}
	}
	return -1
}
