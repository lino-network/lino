package developer

import (
	"math/big"
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestReportConsumption(t *testing.T) {
	ctx, _, dm, _ := setupTest(t, 0)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	dm.RegisterDeveloper(ctx, "developer1", devParam.DeveloperMinDeposit)
	dm.RegisterDeveloper(ctx, "developer2", devParam.DeveloperMinDeposit)

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

	cases := map[string]struct {
		Developer1Consumption             types.Coin
		Developer2Consumption             types.Coin
		ExpectDeveloper1ConsumptionWeight *big.Rat
		ExpectDeveloper2ConsumptionWeight *big.Rat
	}{
		"test normal consumption": {
			types.NewCoinFromInt64(2500 * types.Decimals), types.NewCoinFromInt64(7500 * types.Decimals),
			big.NewRat(1, 4), big.NewRat(3, 4),
		},
		"test empty consumption": {
			types.NewCoinFromInt64(0), types.NewCoinFromInt64(0), big.NewRat(1, 2), big.NewRat(1, 2),
		},
		"issue https://github.com/lino-network/lino/issues/150": {
			types.NewCoinFromInt64(3333333), types.NewCoinFromInt64(4444444),
			big.NewRat(3, 7), big.NewRat(4, 7),
		},
	}
	for testName, cs := range cases {
		dm.ReportConsumption(ctx, "developer1", cs.Developer1Consumption)
		dm.ReportConsumption(ctx, "developer2", cs.Developer2Consumption)

		p1, _ := dm.GetConsumptionWeight(ctx, "developer1")
		if cs.ExpectDeveloper1ConsumptionWeight.Cmp(p1) != 0 {
			t.Errorf(
				"%s: expect developer1 usage weight %v, got %v",
				testName, cs.ExpectDeveloper1ConsumptionWeight, p1)
			return
		}

		p2, _ := dm.GetConsumptionWeight(ctx, "developer2")
		if cs.ExpectDeveloper2ConsumptionWeight.Cmp(p2) != 0 {
			t.Errorf(
				"%s: expect developer2 usage weight %v, got %v",
				testName, cs.ExpectDeveloper2ConsumptionWeight, p2)
			return
		}
		dm.ClearConsumption(ctx)
	}
}
