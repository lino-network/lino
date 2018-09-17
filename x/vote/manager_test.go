package vote

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAddVoter(t *testing.T) {
	ctx, am, vm, _, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)

	testCases := []struct {
		testName       string
		username       types.AccountKey
		coin           types.Coin
		expectedResult sdk.Error
	}{
		{
			testName:       "normal case",
			username:       user1,
			coin:           types.NewCoinFromInt64(100 * types.Decimals),
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		res := vm.AddVoter(ctx, tc.username, tc.coin)
		if !assert.Equal(t, tc.expectedResult, res) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectedResult)
		}
	}
}

func TestCanBecomeValidator(t *testing.T) {
	ctx, am, vm, _, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	valParam, _ := vm.paramHolder.GetValidatorParam(ctx)

	testCases := []struct {
		testName       string
		addVoter       bool
		username       types.AccountKey
		coin           types.Coin
		expectedResult bool
	}{
		{
			testName:       "non-voter can't become validator",
			addVoter:       false,
			username:       user1,
			coin:           types.NewCoinFromInt64(0),
			expectedResult: false,
		},
		{
			testName:       "become validator successfully",
			addVoter:       true,
			username:       user1,
			coin:           valParam.ValidatorMinVotingDeposit,
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		if tc.addVoter {
			err := vm.AddVoter(ctx, tc.username, tc.coin)
			if err != nil {
				t.Errorf("%s: failed to add voter, got err %v", tc.testName, err)
			}
		}
		actualRes := vm.CanBecomeValidator(ctx, tc.username)
		if actualRes != tc.expectedResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, actualRes, tc.expectedResult)
		}
	}
}

func TestAddAndClaimInterest(t *testing.T) {
	testName := "TestAddAndClaimInterest"
	ctx, am, vm, _, _ := setupTest(t, 0)

	accKey := types.AccountKey("accKey")
	minBalance := types.NewCoinFromInt64(1000 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance)

	err := vm.AddVoter(ctx, accKey, c100)
	if err != nil {
		t.Errorf("%s: failed to add voter, got err %v", testName, err)
	}

	err = vm.AddInterest(ctx, accKey, c500)
	if err != nil {
		t.Errorf("%s: failed to add interest, got err %v", testName, err)
	}

	voter, err := vm.storage.GetVoter(ctx, accKey)
	if err != nil {
		t.Errorf("%s: failed to get voter, got err %v", testName, err)
	}

	if !assert.Equal(t, c500, voter.Interest) {
		t.Errorf("%s: diff interest", testName)
	}

	_, err = vm.ClaimInterest(ctx, accKey)
	if err != nil {
		t.Errorf("%s: failed to add claim interest, got err %v", testName, err)
	}
	voter, err = vm.storage.GetVoter(ctx, accKey)
	if err != nil {
		t.Errorf("%s: failed to get voter, got err %v", testName, err)
	}

	if !assert.Equal(t, true, voter.Interest.IsZero()) {
		t.Errorf("%s: diff interest", testName)
	}

}

func TestIsInValidatorList(t *testing.T) {
	ctx, am, vm, _, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	user3 := createTestAccount(ctx, am, "user3", minBalance)

	testCases := []struct {
		testName       string
		username       types.AccountKey
		allValidators  []types.AccountKey
		expectedResult bool
	}{
		{
			testName:       "not in empty validator list",
			username:       user1,
			allValidators:  []types.AccountKey{},
			expectedResult: false,
		},
		{
			testName:       "not in validator list",
			username:       user1,
			allValidators:  []types.AccountKey{user2, user3},
			expectedResult: false,
		},
		{
			testName:       "in validator list",
			username:       user1,
			allValidators:  []types.AccountKey{user1},
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		referenceList := &model.ReferenceList{
			AllValidators: tc.allValidators,
		}
		err := vm.storage.SetReferenceList(ctx, referenceList)
		if err != nil {
			t.Errorf("%s: failed to set reference list, got err %v", tc.testName, err)
		}
		res := vm.IsInValidatorList(ctx, tc.username)
		if res != tc.expectedResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectedResult)
		}
	}
}

func TestIsLegalVoterWithdraw(t *testing.T) {
	ctx, am, vm, _, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)

	vm.AddVoter(ctx, user1, types.NewCoinFromInt64(100*types.Decimals))

	testCases := []struct {
		testName       string
		allValidators  []types.AccountKey
		username       types.AccountKey
		withdraw       types.Coin
		expectedResult bool
	}{
		{
			testName:       "normal case",
			allValidators:  []types.AccountKey{},
			username:       user1,
			withdraw:       types.NewCoinFromInt64(1),
			expectedResult: true,
		},
		{
			testName:       "validator can't withdraw",
			allValidators:  []types.AccountKey{user1},
			username:       user1,
			withdraw:       types.NewCoinFromInt64(1),
			expectedResult: false,
		},
		{
			testName:       "illegal withdraw",
			allValidators:  []types.AccountKey{},
			username:       user1,
			withdraw:       types.NewCoinFromInt64(101 * types.Decimals),
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		referenceList := &model.ReferenceList{
			AllValidators: tc.allValidators,
		}

		err := vm.storage.SetReferenceList(ctx, referenceList)
		if err != nil {
			t.Errorf("%s: failed to set reference list, got err %v", tc.testName, err)
		}

		res := vm.IsLegalVoterWithdraw(ctx, tc.username, tc.withdraw)
		if res != tc.expectedResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectedResult)
		}
	}
}

func TestIsLegalDelegatorWithdraw(t *testing.T) {
	ctx, am, vm, _, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	withdraw := types.NewCoinFromInt64(10 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)

	vm.AddVoter(ctx, user1, types.NewCoinFromInt64(101*types.Decimals))

	testCases := []struct {
		testName       string
		addDelegation  bool
		delegatedCoin  types.Coin
		delegator      types.AccountKey
		voter          types.AccountKey
		withdraw       types.Coin
		expectedResult bool
	}{
		{
			testName:       "no delegation exist, can't withdraw",
			addDelegation:  false,
			delegatedCoin:  types.NewCoinFromInt64(0),
			delegator:      user2,
			voter:          user1,
			withdraw:       withdraw,
			expectedResult: false,
		},
		{
			testName:       "normal case",
			addDelegation:  true,
			delegatedCoin:  types.NewCoinFromInt64(100 * types.Decimals),
			delegator:      user2,
			voter:          user1,
			withdraw:       withdraw,
			expectedResult: true,
		},
		{
			testName:       "can't withdraw 0",
			addDelegation:  false,
			delegatedCoin:  types.NewCoinFromInt64(0),
			delegator:      user2,
			voter:          user1,
			withdraw:       types.NewCoinFromInt64(0),
			expectedResult: false,
		},
		{
			testName:       "can't withdraw 101",
			addDelegation:  false,
			delegatedCoin:  types.NewCoinFromInt64(0),
			delegator:      user2,
			voter:          user1,
			withdraw:       types.NewCoinFromInt64(101 * types.Decimals),
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		if tc.addDelegation {
			err := vm.AddVoter(ctx, tc.delegator, types.NewCoinFromInt64(0))
			if err != nil {
				t.Errorf("%s: failed to add voter, got err %v", tc.testName, err)
			}

			err = vm.AddDelegation(ctx, tc.voter, tc.delegator, tc.delegatedCoin)
			if err != nil {
				t.Errorf("%s: failed to add delegation, got err %v", tc.testName, err)
			}
		}
		res := vm.IsLegalDelegatorWithdraw(ctx, tc.voter, tc.delegator, tc.withdraw)
		if res != tc.expectedResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectedResult)
		}
	}
}
