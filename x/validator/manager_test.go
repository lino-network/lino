//nolint:errcheck
package validator

import (
	"math/rand"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/validator/model"
)

func TestByzantines(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, &gm)
	err := valManager.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(100000 * types.Decimals)
	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i), minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
		err := voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i)), valParam.ValidatorMinVotingDeposit)
		if err != nil {
			panic(err)
		}

		// they will deposit 10,20,30...200, 210
		validatorMinDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((i+1)*10) + validatorMinDeposit/types.Decimals
		deposit := strconv.FormatInt(num, 10)
		valKeys[i] = secp256k1.GenPrivKey().PubKey()
		name := "user" + strconv.Itoa(i)
		msg := NewValidatorDepositMsg(name, deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// byzantine
	byzantineList := []int32{3, 8, 14}
	byzantines := []abci.Evidence{}
	for _, idx := range byzantineList {
		byzantines = append(byzantines, abci.Evidence{Validator: abci.Validator{
			Address: valKeys[idx].Address(),
			Power:   1000}})
	}
	_, err = valManager.FireIncompetentValidator(ctx, byzantines)
	assert.Nil(t, err)

	validatorList3, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 18, len(validatorList3.OncallValidators))
	assert.Equal(t, 18, len(validatorList3.AllValidators))

	for _, idx := range byzantineList {
		assert.Equal(t, -1, types.FindAccountInList(users[idx], validatorList3.OncallValidators))
		assert.Equal(t, -1, types.FindAccountInList(users[idx], validatorList3.AllValidators))
	}

}

func TestAbsentValidatorWillBeFired(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, &gm)
	err := valManager.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(100000 * types.Decimals)
	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i), minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
		err := voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i)), valParam.ValidatorMinVotingDeposit)
		if err != nil {
			panic(err)
		}

		// they will deposit 10,20,30...200, 210
		validatorMinDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((i+1)*10) + validatorMinDeposit/types.Decimals
		deposit := strconv.FormatInt(num, 10)
		valKeys[i] = secp256k1.GenPrivKey().PubKey()
		name := "user" + strconv.Itoa(i)
		msg := NewValidatorDepositMsg(name, deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// construct signing list
	absentList := []int{0, 1, 10}
	index := 0
	signingList := []abci.VoteInfo{}
	for i := 0; i < 21; i++ {
		signingList = append(signingList, abci.VoteInfo{
			Validator: abci.Validator{
				Address: valKeys[i].Address(),
				Power:   1000},
			SignedLastBlock: true,
		})
		if index < len(absentList) && i == absentList[index] {
			signingList[i].SignedLastBlock = false
			index++
		}
	}
	// shuffle the signing validator array
	destSigningList := make([]abci.VoteInfo, len(signingList))
	perm := rand.Perm(len(signingList))
	for i, v := range perm {
		destSigningList[v] = signingList[i]
	}
	err = valManager.UpdateSigningStats(ctx, signingList)
	assert.Nil(t, err)

	index = 0
	for i := 0; i < 21; i++ {
		validator, _ := valManager.storage.GetValidator(ctx, types.AccountKey("user"+strconv.Itoa(i)))
		if index < len(absentList) && i == absentList[index] {
			assert.Equal(t, int64(1), validator.AbsentCommit)
			assert.Equal(t, int64(0), validator.ProducedBlocks)
			index++
		} else {
			assert.Equal(t, int64(0), validator.AbsentCommit)
			assert.Equal(t, int64(1), validator.ProducedBlocks)
		}
	}

	param, _ := valManager.paramHolder.GetValidatorParam(ctx)
	// absent exceeds limitation
	for i := int64(0); i < param.AbsentCommitLimitation; i++ {
		err := valManager.UpdateSigningStats(ctx, signingList)
		assert.Nil(t, err)
	}

	index = 0
	for i := 0; i < 21; i++ {
		validator, _ := valManager.storage.GetValidator(ctx, types.AccountKey("user"+strconv.Itoa(i)))
		if index < len(absentList) && i == absentList[index] {
			assert.Equal(t, int64(601), validator.AbsentCommit)
			assert.Equal(t, int64(0), validator.ProducedBlocks)
			index++
		} else {
			assert.Equal(t, int64(0), validator.AbsentCommit)
			assert.Equal(t, int64(601), validator.ProducedBlocks)
		}
	}

	_, err = valManager.FireIncompetentValidator(ctx, []abci.Evidence{})
	assert.Nil(t, err)
	validatorList2, _ := valManager.storage.GetValidatorList(ctx)

	assert.Equal(t, 18, len(validatorList2.OncallValidators))
	assert.Equal(t, 18, len(validatorList2.AllValidators))

	for _, idx := range absentList {
		assert.Equal(t, -1, types.FindAccountInList(types.AccountKey("user"+strconv.Itoa(idx)), validatorList2.OncallValidators))
		assert.Equal(t, -1, types.FindAccountInList(types.AccountKey("user"+strconv.Itoa(idx)), validatorList2.AllValidators))
	}
}

func TestAbsentValidatorWontBeFired(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, &gm)
	err := valManager.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(100000 * types.Decimals)
	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i), minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
		err := voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i)), valParam.ValidatorMinVotingDeposit)
		if err != nil {
			panic(err)
		}

		// they will deposit 1000,2000,3000...20000, 21000
		validatorMinDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((i+1)*1000) + validatorMinDeposit/types.Decimals
		deposit := strconv.FormatInt(num, 10)
		valKeys[i] = secp256k1.GenPrivKey().PubKey()
		name := "user" + strconv.Itoa(i)
		msg := NewValidatorDepositMsg(name, deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// construct signing list
	absentList := []int{0, 1, 10}
	index := 0
	signingList := []abci.VoteInfo{}
	for i := 0; i < 21; i++ {
		signingList = append(signingList, abci.VoteInfo{
			Validator: abci.Validator{
				Address: valKeys[i].Address(),
				Power:   1000},
			SignedLastBlock: true,
		})
		if index < len(absentList) && i == absentList[index] {
			signingList[i].SignedLastBlock = false
			index++
		}
	}
	// shuffle the signing validator array
	destSigningList := make([]abci.VoteInfo, len(signingList))
	perm := rand.Perm(len(signingList))
	for i, v := range perm {
		destSigningList[v] = signingList[i]
	}
	err = valManager.UpdateSigningStats(ctx, signingList)
	assert.Nil(t, err)

	index = 0
	for i := 0; i < 21; i++ {
		validator, _ := valManager.storage.GetValidator(ctx, types.AccountKey("user"+strconv.Itoa(i)))
		if index < len(absentList) && i == absentList[index] {
			assert.Equal(t, int64(1), validator.AbsentCommit)
			assert.Equal(t, int64(0), validator.ProducedBlocks)
			index++
		} else {
			assert.Equal(t, int64(0), validator.AbsentCommit)
			assert.Equal(t, int64(1), validator.ProducedBlocks)
		}
	}

	param, _ := valManager.paramHolder.GetValidatorParam(ctx)
	// absent exceeds limitation
	for i := int64(0); i < param.AbsentCommitLimitation; i++ {
		err := valManager.UpdateSigningStats(ctx, signingList)
		assert.Nil(t, err)
	}

	index = 0
	for i := 0; i < 21; i++ {
		validator, _ := valManager.storage.GetValidator(ctx, types.AccountKey("user"+strconv.Itoa(i)))
		if index < len(absentList) && i == absentList[index] {
			assert.Equal(t, int64(601), validator.AbsentCommit)
			assert.Equal(t, int64(0), validator.ProducedBlocks)
			index++
		} else {
			assert.Equal(t, int64(0), validator.AbsentCommit)
			assert.Equal(t, int64(601), validator.ProducedBlocks)
		}
	}

	_, err = valManager.FireIncompetentValidator(ctx, []abci.Evidence{})
	assert.Nil(t, err)
	validatorList2, _ := valManager.storage.GetValidatorList(ctx)

	assert.Equal(t, 21, len(validatorList2.OncallValidators))
	assert.Equal(t, 21, len(validatorList2.AllValidators))

	// check deposit has been deducted by 200
	for _, v := range absentList {
		validator, _ := valManager.storage.GetValidator(ctx, types.AccountKey("user"+strconv.Itoa(v)))

		assert.Equal(t, int64(0), validator.AbsentCommit)

		validatorMinDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((v+1)*1000) + validatorMinDeposit/types.Decimals
		num -= 200
		depositCoin := types.NewCoinFromInt64(num * types.Decimals)
		assert.Equal(t, depositCoin, validator.Deposit)
	}
}

func TestGetOncallList(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, &gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(100000 * types.Decimals)
	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i), minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i)), valParam.ValidatorMinVotingDeposit)

		// they will deposit 10,20,30...200, 210
		validatorMinDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((i+1)*10) + validatorMinDeposit/types.Decimals
		deposit := strconv.FormatInt(num, 10)
		valKeys[i] = secp256k1.GenPrivKey().PubKey()
		name := "user" + strconv.Itoa(i)
		msg := NewValidatorDepositMsg(name, deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	lst, _ := valManager.GetValidatorList(ctx)
	for idx, validator := range lst.OncallValidators {
		assert.Equal(t, users[idx], validator)
	}

}

func TestPunishmentBasic(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, &gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	createTestAccount(ctx, am, "user1", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
	createTestAccount(ctx, am, "user2", minBalance.Plus(valParam.ValidatorMinCommittingDeposit))

	valKey1 := secp256k1.GenPrivKey().PubKey()
	valKey2 := secp256k1.GenPrivKey().PubKey()

	voteManager.AddVoter(ctx, "user1", valParam.ValidatorMinVotingDeposit)
	voteManager.AddVoter(ctx, "user2", valParam.ValidatorMinVotingDeposit)

	// let both users register as validator
	msg1 := NewValidatorDepositMsg("user1", coinToString(valParam.ValidatorMinCommittingDeposit), valKey1, "")
	msg2 := NewValidatorDepositMsg("user2", coinToString(valParam.ValidatorMinCommittingDeposit), valKey2, "")
	handler(ctx, msg1)
	handler(ctx, msg2)

	// punish user2 as byzantine (explicitly remove)
	_, err := valManager.PunishOncallValidator(ctx, types.AccountKey("user2"), valParam.PenaltyByzantine, types.PunishByzantine)
	if err != nil {
		panic(err)
	}
	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 1, len(lst.OncallValidators))
	assert.Equal(t, 1, len(lst.AllValidators))
	assert.Equal(t, types.AccountKey("user1"), lst.OncallValidators[0])

	validator, _ := valManager.storage.GetValidator(ctx, "user2")
	assert.Equal(t, true, validator.Deposit.IsZero())

	// punish user1 as missing vote (wont explicitly remove)
	_, err = valManager.PunishOncallValidator(ctx, types.AccountKey("user1"), valParam.PenaltyMissVote, types.PunishDidntVote)
	if err != nil {
		panic(err)
	}
	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 0, len(lst2.OncallValidators))
	assert.Equal(t, 0, len(lst2.AllValidators))

	validator2, _ := valManager.storage.GetValidator(ctx, "user1")
	assert.Equal(t, true, validator2.Deposit.IsZero())
}

func TestPunishmentAndSubstitutionExists(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, &gm)
	valManager.InitGenesis(ctx)

	valParam, _ := valManager.paramHolder.GetValidatorParam(ctx)
	minBalance := types.NewCoinFromInt64(100000 * types.Decimals)

	// create 24 test users
	users := make([]types.AccountKey, 24)
	valKeys := make([]crypto.PubKey, 24)
	for i := 0; i < 24; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1), minBalance.Plus(valParam.ValidatorMinCommittingDeposit))
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), valParam.ValidatorMinVotingDeposit)

		validatorMinDeposit, _ := valParam.ValidatorMinCommittingDeposit.ToInt64()
		num := int64((i+1)*1000) + validatorMinDeposit/types.Decimals
		deposit := strconv.FormatInt(num, 10)

		valKeys[i] = secp256k1.GenPrivKey().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// lowest is user4 with power (min + 400)
	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 21, len(lst.OncallValidators))
	assert.Equal(t, 24, len(lst.AllValidators))
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(4000*types.Decimals)), lst.LowestPower)
	assert.Equal(t, users[3], lst.LowestValidator)

	// punish user4 as missing vote (wont explicitly remove)
	// user3 will become the lowest one with power (min + 3000)
	_, err := valManager.PunishOncallValidator(ctx, users[3], types.NewCoinFromInt64(2000*types.Decimals), types.PunishDidntVote)
	if err != nil {
		panic(err)
	}
	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 21, len(lst2.OncallValidators))
	assert.Equal(t, 24, len(lst2.AllValidators))
	assert.Equal(t, valParam.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(3000*types.Decimals)), lst2.LowestPower)
	assert.Equal(t, users[2], lst2.LowestValidator)

}

func TestInitValidators(t *testing.T) {
	ctx, am, valManager, _, _ := setupTest(t, 0)
	valManager.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(100 * types.Decimals)

	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)

	valKey1 := secp256k1.GenPrivKey().PubKey()
	valKey2 := secp256k1.GenPrivKey().PubKey()

	param, _ := valManager.paramHolder.GetValidatorParam(ctx)

	err := valManager.RegisterValidator(ctx, user1, valKey1, param.ValidatorMinCommittingDeposit, "")
	if err != nil {
		panic(err)
	}
	err = valManager.RegisterValidator(ctx, user2, valKey2, param.ValidatorMinCommittingDeposit, "")
	if err != nil {
		panic(err)
	}

	val1 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey1),
		Power:  types.TendermintValidatorPower,
	}

	val2 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey2),
		Power:  types.TendermintValidatorPower,
	}

	testCases := []struct {
		testName            string
		oncallValidators    []types.AccountKey
		expectedUpdatedList []abci.ValidatorUpdate
	}{
		{
			testName:            "only one oncall validator",
			oncallValidators:    []types.AccountKey{user1},
			expectedUpdatedList: []abci.ValidatorUpdate{val1},
		},
		{
			testName:            "two oncall validators",
			oncallValidators:    []types.AccountKey{user1, user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val1, val2},
		},
		{
			testName:            "another one",
			oncallValidators:    []types.AccountKey{user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val2},
		},
		{
			testName:            "no validators exists",
			oncallValidators:    []types.AccountKey{},
			expectedUpdatedList: []abci.ValidatorUpdate{},
		},
	}

	for _, tc := range testCases {
		lst := &model.ValidatorList{
			OncallValidators: tc.oncallValidators,
		}
		err := valManager.storage.SetValidatorList(ctx, lst)
		if err != nil {
			t.Errorf("%s: failed to set validator list, got err %v", tc.testName, err)
		}
		actualList, err := valManager.GetInitValidators(ctx)
		if err != nil {
			t.Errorf("%s: failed to get validator list, got err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.expectedUpdatedList, actualList) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, actualList, tc.expectedUpdatedList)
		}
	}
}

func TestGetValidatorUpdates(t *testing.T) {
	ctx, am, valManager, _, _ := setupTest(t, 0)
	valManager.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(100 * types.Decimals)

	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)

	valKey1 := secp256k1.GenPrivKey().PubKey()
	valKey2 := secp256k1.GenPrivKey().PubKey()

	param, _ := valManager.paramHolder.GetValidatorParam(ctx)

	err := valManager.RegisterValidator(ctx, user1, valKey1, param.ValidatorMinCommittingDeposit, "")
	if err != nil {
		panic(err)
	}
	err = valManager.RegisterValidator(ctx, user2, valKey2, param.ValidatorMinCommittingDeposit, "")
	if err != nil {
		panic(err)
	}

	val1 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey1),
		Power:  types.TendermintValidatorPower,
	}

	val2 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey2),
		Power:  types.TendermintValidatorPower,
	}

	val1NoPower := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey1),
		Power:  0,
	}

	val2NoPower := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey2),
		Power:  0,
	}

	testCases := []struct {
		testName            string
		oncallValidators    []types.AccountKey
		preBlockValidators  []types.AccountKey
		expectedUpdatedList []abci.ValidatorUpdate
	}{
		{
			testName:            "only one oncall validator",
			oncallValidators:    []types.AccountKey{user1},
			preBlockValidators:  []types.AccountKey{},
			expectedUpdatedList: []abci.ValidatorUpdate{val1},
		},
		{
			testName:            "two oncall validators and one pre block validator",
			oncallValidators:    []types.AccountKey{user1, user2},
			preBlockValidators:  []types.AccountKey{user1},
			expectedUpdatedList: []abci.ValidatorUpdate{val1, val2},
		},
		{
			testName:            "two oncall validatos and two pre block validators",
			oncallValidators:    []types.AccountKey{user1, user2},
			preBlockValidators:  []types.AccountKey{user1, user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val1, val2},
		},
		{
			testName:            "one oncall validator and two pre block validators",
			oncallValidators:    []types.AccountKey{user2},
			preBlockValidators:  []types.AccountKey{user1, user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val1NoPower, val2},
		},
		{
			testName:            "only one pre block validator",
			oncallValidators:    []types.AccountKey{},
			preBlockValidators:  []types.AccountKey{user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val2NoPower},
		},
	}

	for _, tc := range testCases {
		lst := &model.ValidatorList{
			OncallValidators:   tc.oncallValidators,
			PreBlockValidators: tc.preBlockValidators,
		}
		err := valManager.storage.SetValidatorList(ctx, lst)
		if err != nil {
			t.Errorf("%s: failed to set validator list, got err %v", tc.testName, err)
		}

		actualList, err := valManager.GetValidatorUpdates(ctx)
		if err != nil {
			t.Errorf("%s: failed to get validator list, got err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.expectedUpdatedList, actualList) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, actualList, tc.expectedUpdatedList)
		}
	}
}

func TestIsLegalWithdraw(t *testing.T) {
	ctx, am, valManager, _, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(100 * types.Decimals)

	user1 := createTestAccount(ctx, am, "user1", minBalance)
	param, _ := valManager.paramHolder.GetValidatorParam(ctx)
	valManager.InitGenesis(ctx)
	err := valManager.RegisterValidator(
		ctx, user1, secp256k1.GenPrivKey().PubKey(),
		param.ValidatorMinCommittingDeposit.Plus(types.NewCoinFromInt64(100*types.Decimals)), "")
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		testName         string
		oncallValidators []types.AccountKey
		username         types.AccountKey
		withdraw         types.Coin
		expectedResult   bool
	}{
		{
			testName:         "withdraw amount is a little less than minimum requirement",
			oncallValidators: []types.AccountKey{user1},
			username:         user1,
			withdraw:         param.ValidatorMinWithdraw.Minus(types.NewCoinFromInt64(1)),
			expectedResult:   false,
		},
		{
			testName:         "remaining coin is less than minimum committing deposit",
			oncallValidators: []types.AccountKey{user1},
			username:         user1,
			withdraw:         param.ValidatorMinCommittingDeposit,
			expectedResult:   false,
		},
		{
			testName:         "oncall validator can't withdraw",
			oncallValidators: []types.AccountKey{user1},
			username:         user1,
			withdraw:         param.ValidatorMinWithdraw,
			expectedResult:   false,
		},
		{
			testName:         "withdraw successfully",
			oncallValidators: []types.AccountKey{},
			username:         user1,
			withdraw:         param.ValidatorMinWithdraw,
			expectedResult:   true,
		},
	}

	for _, tc := range testCases {
		lst := &model.ValidatorList{
			OncallValidators: tc.oncallValidators,
		}
		err := valManager.storage.SetValidatorList(ctx, lst)
		if err != nil {
			t.Errorf("%s: failed to set validator list, got err %v", tc.testName, err)
		}

		res := valManager.IsLegalWithdraw(ctx, tc.username, tc.withdraw)
		if res != tc.expectedResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectedResult)
		}
	}
}
