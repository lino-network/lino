package developer

import (
	"strconv"
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	accstore "github.com/lino-network/lino/x/account/model"
)

func TestRegistertBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	minDeposit, _ := devParam.DeveloperMinDeposit.ToInt64()
	deposit := strconv.FormatInt(minDeposit/types.Decimals, 10)
	testCases := []struct {
		testName     string
		msg          DeveloperRegisterMsg
		expectResult sdk.Result
	}{
		{
			testName: "normal update",
			msg: NewDeveloperRegisterMsg(
				"developer1", deposit, "https://lino.network", "decentralized autonomous video content economy", "app meta data"),
			expectResult: sdk.Result{},
		},
		{
			testName: "invalid username",
			msg: NewDeveloperRegisterMsg(
				"developer2", deposit, "https://lino.network", "decentralized autonomous video content economy", "app meta data"),
			expectResult: ErrAccountNotFound().Result(),
		},
		{
			testName: "developer doesn't exist",
			msg: NewDeveloperRegisterMsg(
				"developer1", deposit, "https://lino.network", "decentralized autonomous video content economy", "app meta data"),
			expectResult: ErrDeveloperAlreadyExist("developer1").Result(),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if result.Code != tc.expectResult.Code {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}
	}

	// check acc1's money has been withdrawn
	acc1Saving, _ := am.GetSavingFromBank(ctx, types.AccountKey("developer1"))
	assert.Equal(t, true, acc1Saving.IsEqual(minBalance))
	assert.Equal(t, true, dm.DoesDeveloperExist(ctx, types.AccountKey("developer1")))

	// check acc1 is in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 1, len(lst.AllDevelopers))
	assert.Equal(t, types.AccountKey("developer1"), lst.AllDevelopers[0])
}

func TestUpdateBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	createTestAccount(ctx, am, "developer2", devParam.DeveloperMinDeposit.Plus(minBalance))
	dm.RegisterDeveloper(ctx, "developer1", devParam.DeveloperMinDeposit, "", "", "")

	testCases := []struct {
		testName     string
		msg          DeveloperUpdateMsg
		expectResult sdk.Result
	}{
		{
			testName: "normal update",
			msg: NewDeveloperUpdateMsg(
				"developer1", "https://lino.network", "decentralized autonomous video content economy", "app meta data"),
			expectResult: sdk.Result{},
		},
		{
			testName: "invalid username",
			msg: NewDeveloperUpdateMsg(
				"invalid", "https://lino.network", "decentralized autonomous video content economy", "app meta data"),
			expectResult: ErrAccountNotFound().Result(),
		},
		{
			testName: "developer doesn't exist",
			msg: NewDeveloperUpdateMsg(
				"developer2", "https://lino.network", "decentralized autonomous video content economy", "app meta data"),
			expectResult: ErrDeveloperNotFound().Result(),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if result.Code != tc.expectResult.Code {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}
	}
}

func TestRevokeBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	devMinDeposit, _ := devParam.DeveloperMinDeposit.ToInt64()
	deposit := strconv.FormatInt(devMinDeposit/types.Decimals, 10)
	msg := NewDeveloperRegisterMsg("developer1", deposit, "", "", "")
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

	// revoke developer doesn't exist
	msg3 := NewDeveloperRevokeMsg("developer2")
	res3 := handler(ctx, msg3)
	assert.Equal(t, ErrDeveloperNotFound().Result(), res3)
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
		{
			testName:               "return coin to user",
			times:                  10,
			interval:               2,
			returnedCoin:           types.NewCoinFromInt64(100),
			expectedFrozenListLen:  1,
			expectedFrozenMoney:    types.NewCoinFromInt64(100),
			expectedFrozenTimes:    10,
			expectedFrozenInterval: 2,
		},
		{
			testName:               "return coin to user again",
			times:                  100000,
			interval:               20000,
			returnedCoin:           types.NewCoinFromInt64(100000),
			expectedFrozenListLen:  2,
			expectedFrozenMoney:    types.NewCoinFromInt64(100000),
			expectedFrozenTimes:    100000,
			expectedFrozenInterval: 20000,
		},
	}

	for _, tc := range testCases {
		err := returnCoinTo(
			ctx, "user", gm, am, tc.times, tc.interval, tc.returnedCoin)
		if err != nil {
			t.Errorf("%s: failed to return coin, got err %v", tc.testName, err)
		}

		lst, err := am.GetFrozenMoneyList(ctx, types.AccountKey("user"))
		if err != nil {
			t.Errorf("%s: failed to return coin, got err %v", tc.testName, err)
		}
		if len(lst) != tc.expectedFrozenListLen {
			t.Errorf("%s: diff list len, got %v, want %v", tc.testName, len(lst), tc.expectedFrozenListLen)
		}
		if !lst[len(lst)-1].Amount.IsEqual(tc.expectedFrozenMoney) {
			t.Errorf("%s: diff amount, got %v, want %v", tc.testName, lst[len(lst)-1].Amount, tc.expectedFrozenMoney)
		}
		if lst[len(lst)-1].Times != tc.expectedFrozenTimes {
			t.Errorf("%s: diff times, got %v, want %v", tc.testName, lst[len(lst)-1].Times, tc.expectedFrozenTimes)
		}
		if lst[len(lst)-1].Interval != tc.expectedFrozenInterval {
			t.Errorf("%s: diff interval, got %v, want %v", tc.testName, lst[len(lst)-1].Interval, tc.expectedFrozenInterval)
		}
	}
}

func TestGrantPermissionMsg(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	param, err := dm.paramHolder.GetDeveloperParam(ctx)
	assert.Nil(t, err)

	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance)
	createTestAccount(ctx, am, "user2", minBalance)
	createTestAccount(ctx, am, "app", minBalance)
	err = dm.RegisterDeveloper(ctx, types.AccountKey("app"), param.DeveloperMinDeposit, "", "", "")
	assert.Nil(t, err)

	testCases := []struct {
		testName     string
		msg          GrantPermissionMsg
		expectResult sdk.Result
	}{
		{
			testName:     "normal grant app permission",
			msg:          NewGrantPermissionMsg("user1", "app", 10000, types.AppPermission, "0"),
			expectResult: sdk.Result{},
		},
		{
			testName:     "grant permission to non-exist app",
			msg:          NewGrantPermissionMsg("user2", "invalidApp", 10000, types.AppPermission, "0"),
			expectResult: ErrDeveloperNotFound().Result(),
		},
		{
			testName:     "grant permission to non-exist user",
			msg:          NewGrantPermissionMsg("invalid", "app", 10000, types.AppPermission, "0"),
			expectResult: ErrAccountNotFound().Result(),
		},
		{
			testName:     "grant permission with invalid amount",
			msg:          NewGrantPermissionMsg("user1", "app", 10000, types.PreAuthorizationPermission, "-1"),
			expectResult: types.ErrInvalidCoins("LNO can't be less than lower bound").Result(),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if result.Code != tc.expectResult.Code {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}
	}
}

func TestHandlePreAuthorizationMsg(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	param, err := dm.paramHolder.GetDeveloperParam(ctx)
	assert.Nil(t, err)

	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance)
	createTestAccount(ctx, am, "user2", minBalance)
	createTestAccount(ctx, am, "app", minBalance)
	err = dm.RegisterDeveloper(ctx, types.AccountKey("app"), param.DeveloperMinDeposit, "", "", "")
	assert.Nil(t, err)

	testCases := []struct {
		testName     string
		msg          PreAuthorizationMsg
		expectResult sdk.Result
	}{
		{
			testName:     "normal preauthorization msg",
			msg:          NewPreAuthorizationMsg("user1", "app", 10000, types.LNO("100")),
			expectResult: sdk.Result{},
		},
		{
			testName:     "grant permission to non-exist app",
			msg:          NewPreAuthorizationMsg("user2", "invalidApp", 10000, types.LNO("100")),
			expectResult: ErrDeveloperNotFound().Result(),
		},
		{
			testName:     "grant permission to non-exist user",
			msg:          NewPreAuthorizationMsg("invalid", "app", 10000, types.LNO("100")),
			expectResult: ErrAccountNotFound().Result(),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if result.Code != tc.expectResult.Code {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
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
	_, _, appAppPriv := createTestAccount(ctx, am, "app", minBalance)

	err = dm.RegisterDeveloper(ctx, types.AccountKey("app"), param.DeveloperMinDeposit, "", "", "")
	assert.Nil(t, err)
	err = am.AuthorizePermission(
		ctx, types.AccountKey("user1"), types.AccountKey("app"), 1000, types.AppPermission, types.NewCoinFromInt64(0))
	assert.Nil(t, err)

	testCases := []struct {
		testName     string
		msg          RevokePermissionMsg
		expectResult sdk.Result
	}{
		{
			testName:     "normal revoke app permission",
			msg:          NewRevokePermissionMsg("user1", appAppPriv.PubKey()),
			expectResult: sdk.Result{},
		},
		{
			testName:     "revoke non-exist pubkey",
			msg:          NewRevokePermissionMsg("user1", appAppPriv.PubKey()),
			expectResult: accstore.ErrGrantPubKeyNotFound().Result(),
		},
		{
			testName:     "invalid revoke user",
			msg:          NewRevokePermissionMsg("invalid", appAppPriv.PubKey()),
			expectResult: ErrAccountNotFound().Result(),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if result.Code != tc.expectResult.Code {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}
	}
}
