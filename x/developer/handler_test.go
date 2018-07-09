package developer

import (
	"strconv"
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	accstore "github.com/lino-network/lino/x/account/model"
)

func TestRegistertBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	deposit := strconv.FormatInt(devParam.DeveloperMinDeposit.ToInt64()/types.Decimals, 10)
	msg := NewDeveloperRegisterMsg("developer1", deposit)
	res := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, res)

	// check acc1's money has been withdrawn
	acc1Saving, _ := am.GetSavingFromBank(ctx, types.AccountKey("developer1"))
	assert.Equal(t, true, acc1Saving.IsEqual(minBalance))
	assert.Equal(t, true, dm.DoesDeveloperExist(ctx, types.AccountKey("developer1")))

	// check acc1 is in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 1, len(lst.AllDevelopers))
	assert.Equal(t, types.AccountKey("developer1"), lst.AllDevelopers[0])

}

func TestRevokeBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	deposit := strconv.FormatInt(devParam.DeveloperMinDeposit.ToInt64()/types.Decimals, 10)
	msg := NewDeveloperRegisterMsg("developer1", deposit)
	handler(ctx, msg)

	msg2 := NewDeveloperRevokeMsg("developer1")
	res2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, res2)
	// check acc1's depoist has not been added back
	acc1Saving, _ := am.GetSavingFromBank(ctx, types.AccountKey("developer1"))
	assert.Equal(t, true, acc1Saving.IsEqual(minBalance))
	assert.Equal(t, false, dm.DoesDeveloperExist(ctx, types.AccountKey("developer1")))

	// check acc1 is not in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 0, len(lst.AllDevelopers))
	assert.Equal(t, false, dm.DoesDeveloperExist(ctx, "developer1"))
}

func TestAddFrozenMoney(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	dm.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user", minBalance)

	testCases := []struct {
		testName               string
		times                  int64
		interval               int64
		returnedCoin           types.Coin
		expectedFrozenListLen  int
		expectedFrozenMoney    types.Coin
		expectedFrozenTimes    int64
		expectedFrozenInterval int64
	}{
		{"return coin to user", 10, 2, types.NewCoinFromInt64(100), 1, types.NewCoinFromInt64(100), 10, 2},
		{"return coin to user multiple times", 100000, 20000, types.NewCoinFromInt64(100000), 2, types.NewCoinFromInt64(100000), 100000, 20000},
	}

	for _, tc := range testCases {
		err := returnCoinTo(
			ctx, "user", gm, am, tc.times, tc.interval, tc.returnedCoin)
		assert.Equal(t, nil, err)
		lst, err := am.GetFrozenMoneyList(ctx, types.AccountKey("user"))
		assert.Equal(t, tc.expectedFrozenListLen, len(lst))
		assert.Equal(t, tc.expectedFrozenMoney, lst[len(lst)-1].Amount)
		assert.Equal(t, tc.expectedFrozenTimes, lst[len(lst)-1].Times)
		assert.Equal(t, tc.expectedFrozenInterval, lst[len(lst)-1].Interval)
	}
}

func TestGrantPermissionMsg(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	param, err := dm.paramHolder.GetDeveloperParam(ctx)
	assert.Nil(t, err)
	accParam, err := dm.paramHolder.GetAccountParam(ctx)
	assert.Nil(t, err)

	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance)
	createTestAccount(ctx, am, "user2", minBalance)
	createTestAccount(ctx, am, "app", minBalance)
	err = dm.RegisterDeveloper(ctx, types.AccountKey("app"), param.DeveloperMinDeposit)
	assert.Nil(t, err)

	testCases := []struct {
		testName     string
		msg          GrantPermissionMsg
		expectResult sdk.Result
	}{
		{"normal grant post permission",
			NewGrantPermissionMsg("user1", "app", 10000, 1, types.PostPermission), sdk.Result{}},
		{"normal grant micropayment permission",
			NewGrantPermissionMsg("user2", "app", 10000, 1, types.MicropaymentPermission), sdk.Result{}},
		{"grant permission to non-exist app",
			NewGrantPermissionMsg("user2", "invalidApp", 10000, 1, types.MicropaymentPermission), ErrDeveloperNotFound().Result()},
		{"grant permission to non-exist user",
			NewGrantPermissionMsg("invalid", "app", 10000, 1, types.MicropaymentPermission), ErrAccountNotFound().Result()},
		{"grant permission exceeds maximum limitation",
			NewGrantPermissionMsg("user1", "app", 10000, 100, types.MicropaymentPermission),
			acc.ErrGrantTimesExceedsLimitation(accParam.MaximumMicropaymentGrantTimes).Result()},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if result.Code != tc.expectResult.Code {
			t.Errorf("%s: test failed, expect %v, got %v", tc.testName, tc.expectResult, result)
		}
	}
}

func TestRevokePermissionMsg(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	param, err := dm.paramHolder.GetDeveloperParam(ctx)
	assert.Nil(t, err)
	accParam, err := dm.paramHolder.GetAccountParam(ctx)
	assert.Nil(t, err)

	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", accParam.RegisterFee)
	createTestAccount(ctx, am, "user2", minBalance)
	appPriv := createTestAccount(ctx, am, "app", minBalance)
	err = dm.RegisterDeveloper(ctx, types.AccountKey("app"), param.DeveloperMinDeposit)
	assert.Nil(t, err)
	err = am.AuthorizePermission(
		ctx, types.AccountKey("user1"), types.AccountKey("app"), 1000, 10, types.PostPermission)
	err = am.AuthorizePermission(
		ctx, types.AccountKey("user1"), types.AccountKey("app"), 1000, 10, types.MicropaymentPermission)
	assert.Nil(t, err)

	testCases := []struct {
		testName     string
		msg          RevokePermissionMsg
		expectResult sdk.Result
	}{
		{"normal revoke post permission",
			NewRevokePermissionMsg("user1", appPriv.Generate(2).PubKey(), types.PostPermission), sdk.Result{}},
		{"revoke non-exist pubkey",
			NewRevokePermissionMsg("user1", appPriv.PubKey(), types.PostPermission), accstore.ErrGrantPubKeyNotFound().Result()},
		{"revoke pubkey permission mismatch",
			NewRevokePermissionMsg("user1", appPriv.Generate(1).PubKey(), types.PostPermission),
			acc.ErrRevokePermissionLevelMismatch(types.MicropaymentPermission, types.PostPermission).Result()},
		{"invalid revoke user",
			NewRevokePermissionMsg("invalid", appPriv.Generate(1).PubKey(), types.MicropaymentPermission), ErrAccountNotFound().Result()},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if result.Code != tc.expectResult.Code {
			t.Errorf("%s: test failed, expect %v, got %v", tc.testName, tc.expectResult, result)
		}
	}
}
