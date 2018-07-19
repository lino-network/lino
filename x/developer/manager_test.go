package developer

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestReportConsumption(t *testing.T) {
	ctx, _, dm, _ := setupTest(t, 0)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	dm.RegisterDeveloper(ctx, "developer1", devParam.DeveloperMinDeposit, "", "", "")
	dm.RegisterDeveloper(ctx, "developer2", devParam.DeveloperMinDeposit, "", "", "")

	con1 := types.NewCoinFromInt64(100)
	dm.ReportConsumption(ctx, "developer1", con1)
	p1, _ := dm.GetConsumptionWeight(ctx, "developer1")
	assert.True(t, p1.Cmp(big.NewRat(1, 1)) == 0)

	con2 := types.NewCoinFromInt64(100)
	dm.ReportConsumption(ctx, "developer2", con2)
	p2, _ := dm.GetConsumptionWeight(ctx, "developer1")
	assert.True(t, p2.Cmp(big.NewRat(1, 2)) == 0)

	dm.ClearConsumption(ctx)
	p3, _ := dm.GetConsumptionWeight(ctx, "developer1")
	assert.True(t, p3.Cmp(big.NewRat(1, 2)) == 0)

	testCases := map[string]struct {
		developer1Consumption             types.Coin
		developer2Consumption             types.Coin
		expectDeveloper1ConsumptionWeight sdk.Rat
		expectDeveloper2ConsumptionWeight sdk.Rat
	}{
		"test normal consumption": {
			developer1Consumption:             types.NewCoinFromInt64(2500 * types.Decimals),
			developer2Consumption:             types.NewCoinFromInt64(7500 * types.Decimals),
			expectDeveloper1ConsumptionWeight: sdk.NewRat(1, 4),
			expectDeveloper2ConsumptionWeight: sdk.NewRat(3, 4),
		},
		"test empty consumption": {
			developer1Consumption:             types.NewCoinFromInt64(0),
			developer2Consumption:             types.NewCoinFromInt64(0),
			expectDeveloper1ConsumptionWeight: sdk.NewRat(1, 2),
			expectDeveloper2ConsumptionWeight: sdk.NewRat(1, 2),
		},
		"issue https://github.com/lino-network/lino/issues/150": {
			developer1Consumption:             types.NewCoinFromInt64(3333333),
			developer2Consumption:             types.NewCoinFromInt64(4444444),
			expectDeveloper1ConsumptionWeight: sdk.NewRat(3, 7),
			expectDeveloper2ConsumptionWeight: sdk.NewRat(4, 7),
		},
	}
	for testName, tc := range testCases {
		dm.ReportConsumption(ctx, "developer1", tc.developer1Consumption)
		dm.ReportConsumption(ctx, "developer2", tc.developer2Consumption)

		p1, _ := dm.GetConsumptionWeight(ctx, "developer1")
		if !tc.expectDeveloper1ConsumptionWeight.Equal(p1) {
			t.Errorf("%s: diff developer1 usage weight, got %v, want %v",
				testName, p1, tc.expectDeveloper1ConsumptionWeight)
			return
		}

		p2, _ := dm.GetConsumptionWeight(ctx, "developer2")
		if !tc.expectDeveloper2ConsumptionWeight.Equal(p2) {
			t.Errorf("%s: diff developer2 usage weight, got %v, want %v",
				testName, p2, tc.expectDeveloper2ConsumptionWeight)
			return
		}
		dm.ClearConsumption(ctx)
	}
}
