package developer

import (
	"strconv"
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestRegistertBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	developer1 := createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	deposit := strconv.FormatInt(devParam.DeveloperMinDeposit.ToInt64()/types.Decimals, 10)
	msg := NewDeveloperRegisterMsg("developer1", deposit)
	res := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, res)

	// check acc1's money has been withdrawn
	acc1Saving, _ := am.GetSavingFromBank(ctx, developer1)
	assert.Equal(t, true, acc1Saving.IsEqual(minBalance))
	assert.Equal(t, true, dm.IsDeveloperExist(ctx, developer1))

	// check acc1 is in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 1, len(lst.AllDevelopers))
	assert.Equal(t, developer1, lst.AllDevelopers[0])

}

func TestRevokeBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	developer1 := createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	deposit := strconv.FormatInt(devParam.DeveloperMinDeposit.ToInt64()/types.Decimals, 10)
	msg := NewDeveloperRegisterMsg("developer1", deposit)
	handler(ctx, msg)

	msg2 := NewDeveloperRevokeMsg("developer1")
	res2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, res2)
	// check acc1's depoist has not been added back
	acc1Saving, _ := am.GetSavingFromBank(ctx, developer1)
	assert.Equal(t, true, acc1Saving.IsEqual(minBalance))
	assert.Equal(t, false, dm.IsDeveloperExist(ctx, developer1))

	// check acc1 is not in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 0, len(lst.AllDevelopers))
	assert.Equal(t, false, dm.IsDeveloperExist(ctx, "developer1"))
}

func TestAddFrozenMoney(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	dm.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user := createTestAccount(ctx, am, "user", minBalance)

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
		lst, err := am.GetFrozenMoneyList(ctx, user)
		assert.Equal(t, tc.expectedFrozenListLen, len(lst))
		assert.Equal(t, tc.expectedFrozenMoney, lst[len(lst)-1].Amount)
		assert.Equal(t, tc.expectedFrozenTimes, lst[len(lst)-1].Times)
		assert.Equal(t, tc.expectedFrozenInterval, lst[len(lst)-1].Interval)

	}
}
