package account

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
)

func checkBankKVByAddress(t *testing.T, ctx sdk.Context, addr sdk.Address, bank model.AccountBank) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	bankPtr, err := accStorage.GetBankFromAddress(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, bank, *bankPtr, "bank should be equal")
}

func checkPendingStake(t *testing.T, ctx sdk.Context, addr sdk.Address, pendingStakeQueue model.PendingStakeQueue) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	pendingStakeQueuePtr, err := accStorage.GetPendingStakeQueue(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, pendingStakeQueue, *pendingStakeQueuePtr, "pending stake should be equal")
}

func checkAccountInfo(t *testing.T, ctx sdk.Context, accKey types.AccountKey, accInfo model.AccountInfo) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	infoPtr, err := accStorage.GetInfo(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *infoPtr, "accout meta should be equal")
}

func checkAccountMeta(t *testing.T, ctx sdk.Context, accKey types.AccountKey, accMeta model.AccountMeta) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	metaPtr, err := accStorage.GetMeta(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *metaPtr, "accout meta should be equal")
}

func TestIsAccountExist(t *testing.T) {
	ctx, am := setupTest(t, 1)
	createTestAccount(ctx, am, "user1")
	assert.True(t, am.IsAccountExist(ctx, types.AccountKey("user1")))
}

func TestAddCoinToAddress(t *testing.T) {
	ctx, am := setupTest(t, 1)

	// add coin to non-exist account
	err := am.AddCoinToAddress(ctx, sdk.Address("test"), types.NewCoin(0))
	assert.Nil(t, err)

	bank := model.AccountBank{
		Address: sdk.Address("test"),
		Balance: types.NewCoin(0),
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue := model.PendingStakeQueue{
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + CoinDays*24*3600,
			Coin:      types.NewCoin(0),
		}}}
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)

	// add coin to exist bank
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: time.Now().Unix()})
	err = am.AddCoinToAddress(ctx, sdk.Address("test"), types.NewCoin(100))
	assert.Nil(t, err)
	bank = model.AccountBank{
		Address: sdk.Address("test"),
		Balance: types.NewCoin(100),
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList,
		model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + CoinDays*24*3600,
			Coin:      types.NewCoin(100),
		})
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)

	// add coin to exist bank after previous coin day
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 3, Time: (ctx.BlockHeader().Time + 3600*24*CoinDays + 1)})
	err = am.AddCoinToAddress(ctx, sdk.Address("test"), types.NewCoin(100))
	assert.Nil(t, err)
	bank = model.AccountBank{
		Address: sdk.Address("test"),
		Balance: types.NewCoin(200),
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList,
		model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + CoinDays*24*3600,
			Coin:      types.NewCoin(100),
		})
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)
}

func TestCreateAccount(t *testing.T) {
	ctx, am := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")

	// normal test
	assert.False(t, am.IsAccountExist(ctx, accKey))
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), types.NewCoin(100))
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey, priv.PubKey(), types.NewCoin(0))
	assert.Nil(t, err)

	assert.True(t, am.IsAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Address:  priv.PubKey().Address(),
		Balance:  types.NewCoin(100),
		Username: accKey,
	}
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	pendingStakeQueue := model.PendingStakeQueue{
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + CoinDays*24*3600,
			Coin:      types.NewCoin(100),
		}}}
	checkPendingStake(t, ctx, priv.PubKey().Address(), pendingStakeQueue)
	accInfo := model.AccountInfo{
		Username: accKey,
		Created:  1,
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	checkAccountInfo(t, ctx, accKey, accInfo)
	accMeta := model.AccountMeta{
		LastActivity:   ctx.BlockHeight(),
		ActivityBurden: types.DefaultActivityBurden,
	}
	checkAccountMeta(t, ctx, accKey, accMeta)

	// username already took
	err = am.CreateAccount(ctx, accKey, priv.PubKey(), types.NewCoin(0))
	assert.Equal(t, ErrAccountAlreadyExists(accKey), err)

	// bank already registered
	err = am.CreateAccount(ctx, types.AccountKey("newKey"), priv.PubKey(), types.NewCoin(0))
	assert.Equal(t, ErrBankAlreadyRegistered(), err)

	// bank doesn't exist
	priv2 := crypto.GenPrivKeyEd25519()
	err = am.CreateAccount(ctx, types.AccountKey("newKey"), priv2.PubKey(), types.NewCoin(0))
	assert.Equal(t, "Error{311:create account newKey failed,Error{310:account bank doesn't exist,<nil>,0},1}", err.Error())

	// register fee doesn't enough
	err = am.AddCoinToAddress(ctx, priv2.PubKey().Address(), types.NewCoin(100))
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, types.AccountKey("newKey"), priv2.PubKey(), types.NewCoin(101))
	assert.Equal(t, ErrRegisterFeeInsufficient(), err)
}
