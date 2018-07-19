package vote

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAddVoter(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	param, _ := vm.paramHolder.GetVoteParam(ctx)

	testCases := []struct {
		testName       string
		username       types.AccountKey
		coin           types.Coin
		expectedResult sdk.Error
	}{
		{
			testName:       "insufficient deposit",
			username:       user1,
			coin:           types.NewCoinFromInt64(100 * types.Decimals),
			expectedResult: ErrInsufficientDeposit(),
		},
		{
			testName:       "normal case",
			username:       user1,
			coin:           param.VoterMinDeposit,
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
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
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
			testName:       "validator has reached min deposit",
			addVoter:       true,
			username:       user1,
			coin:           voteParam.VoterMinDeposit,
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

func TestIsInValidatorList(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
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
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	param, _ := vm.paramHolder.GetVoteParam(ctx)

	vm.AddVoter(ctx, user1, param.VoterMinDeposit.Plus(types.NewCoinFromInt64(100*types.Decimals)))

	testCases := []struct {
		testName       string
		allValidators  []types.AccountKey
		username       types.AccountKey
		withdraw       types.Coin
		expectedResult bool
	}{
		{
			testName:       "illegal withdraw that less than minimum withdraw",
			allValidators:  []types.AccountKey{},
			username:       user1,
			withdraw:       param.VoterMinWithdraw.Minus(types.NewCoinFromInt64(1 * types.Decimals)),
			expectedResult: false,
		},
		{
			testName:       "normal case",
			allValidators:  []types.AccountKey{},
			username:       user1,
			withdraw:       param.VoterMinWithdraw,
			expectedResult: true,
		},
		{
			testName:       "validator can't withdraw",
			allValidators:  []types.AccountKey{user1},
			username:       user1,
			withdraw:       param.VoterMinWithdraw,
			expectedResult: false,
		},
		{
			testName:       "illegal withdraw",
			allValidators:  []types.AccountKey{},
			username:       user1,
			withdraw:       types.NewCoinFromInt64(100),
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
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	param, _ := vm.paramHolder.GetVoteParam(ctx)

	vm.AddVoter(ctx, user1, param.VoterMinDeposit)

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
			withdraw:       param.DelegatorMinWithdraw,
			expectedResult: false,
		},
		{
			testName:       "normal case",
			addDelegation:  true,
			delegatedCoin:  types.NewCoinFromInt64(100 * types.Decimals),
			delegator:      user2,
			voter:          user1,
			withdraw:       param.DelegatorMinWithdraw,
			expectedResult: true,
		},
		{
			testName:       "no delegation exist, can't withdraw 0",
			addDelegation:  false,
			delegatedCoin:  types.NewCoinFromInt64(0),
			delegator:      user2,
			voter:          user1,
			withdraw:       types.NewCoinFromInt64(0),
			expectedResult: false,
		},
		{
			testName:       "no delegation exist, can't withdraw 101",
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
			err := vm.AddDelegation(ctx, tc.voter, tc.delegator, tc.delegatedCoin)
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
