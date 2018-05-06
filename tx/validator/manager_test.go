package validator

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
)

func TestAbsentValidator(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i))
		am.AddCoin(ctx, users[i], c2000)

		// let user register as voter first
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i)), c8000)

		// they will deposit 10,20,30...200, 210
		num := (i+1)*10 + 1001
		deposit := types.LNO(strconv.Itoa(num))
		valKeys[i] = crypto.GenPrivKeyEd25519().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i), deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}
	absentList := []int32{0, 1, 10, 20}
	err := valManager.UpdateAbsentValidator(ctx, absentList)
	assert.Nil(t, err)

	validatorList, _ := valManager.storage.GetValidatorList(ctx)
	for _, idx := range absentList {
		validator, _ := valManager.storage.GetValidator(ctx, validatorList.OncallValidators[idx])
		assert.Equal(t, validator.AbsentCommit, 1)
	}

	// absent exceeds limitation
	for i := 0; i < types.AbsentCommitLimitation; i++ {
		err := valManager.UpdateAbsentValidator(ctx, absentList)
		assert.Nil(t, err)
	}

	for _, idx := range absentList {
		validator, _ := valManager.storage.GetValidator(ctx, validatorList.OncallValidators[idx])
		assert.Equal(t, validator.AbsentCommit, 101)
	}

	_, err = valManager.FireIncompetentValidator(ctx, []abci.Evidence{})
	assert.Nil(t, err)
	validatorList2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 17, len(validatorList2.OncallValidators))
	assert.Equal(t, 17, len(validatorList2.AllValidators))

	for _, idx := range absentList {
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList2.OncallValidators))
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList2.AllValidators))
	}

	// byzantine
	byzantineList := []int32{3, 8, 14}
	byzantines := []abci.Evidence{}
	for _, idx := range byzantineList {
		byzantines = append(byzantines, abci.Evidence{PubKey: valKeys[idx].Bytes()})
	}
	_, err = valManager.FireIncompetentValidator(ctx, byzantines)
	assert.Nil(t, err)

	validatorList3, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 14, len(validatorList3.OncallValidators))
	assert.Equal(t, 14, len(validatorList3.AllValidators))

	for _, idx := range byzantineList {
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList3.OncallValidators))
		assert.Equal(t, -1, FindAccountInList(users[idx], validatorList3.AllValidators))
	}
}

func TestGetOncallList(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create 21 test users
	users := make([]types.AccountKey, 21)
	valKeys := make([]crypto.PubKey, 21)
	for i := 0; i < 21; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i))
		am.AddCoin(ctx, users[i], c2000)
		// let user register as voter first
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i)), c8000)

		// they will deposit 10,20,30...200, 210
		num := (i+1)*10 + 1001
		deposit := types.LNO(strconv.Itoa(num))
		valKeys[i] = crypto.GenPrivKeyEd25519().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i), deposit, valKeys[i], "")
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
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create test users
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)
	user2 := createTestAccount(ctx, am, "user2")
	am.AddCoin(ctx, user2, c2000)
	// let user register as voter first
	voteManager.AddVoter(ctx, types.AccountKey("user1"), c8000)
	voteManager.AddVoter(ctx, types.AccountKey("user2"), c8000)

	// let both users register as validator
	valKey := crypto.GenPrivKeyEd25519().PubKey()
	msg := NewValidatorDepositMsg("user1", l1100, valKey, "")
	handler(ctx, msg)

	valKey2 := crypto.GenPrivKeyEd25519().PubKey()
	msg2 := NewValidatorDepositMsg("user2", l1600, valKey2, "")
	handler(ctx, msg2)

	param, _ := valManager.paramHolder.GetValidatorParam(ctx)
	// punish user2 as byzantine (explicitly remove)
	valManager.PunishOncallValidator(ctx, types.AccountKey("user2"), param.PenaltyByzantine, true)
	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 1, len(lst.OncallValidators))
	assert.Equal(t, 1, len(lst.AllValidators))
	assert.Equal(t, types.AccountKey("user1"), lst.OncallValidators[0])

	validator, _ := valManager.storage.GetValidator(ctx, "user2")
	assert.Equal(t, c0, validator.Deposit)

	// punish user1 as missing vote (wont explicitly remove)
	valManager.PunishOncallValidator(ctx, types.AccountKey("user1"), param.PenaltyMissVote, false)
	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 0, len(lst2.OncallValidators))
	assert.Equal(t, 0, len(lst2.AllValidators))

	validator2, _ := valManager.storage.GetValidator(ctx, "user1")
	assert.Equal(t, c0, validator2.Deposit)
}

func TestPunishmentAndSubstitutionExists(t *testing.T) {
	ctx, am, valManager, voteManager, gm := setupTest(t, 0)
	handler := NewHandler(am, valManager, voteManager, gm)
	valManager.InitGenesis(ctx)

	// create 21 test users
	users := make([]types.AccountKey, 24)
	valKeys := make([]crypto.PubKey, 24)
	for i := 0; i < 24; i++ {
		users[i] = createTestAccount(ctx, am, "user"+strconv.Itoa(i+1))
		am.AddCoin(ctx, users[i], c8000)
		// let user register as voter first
		voteManager.AddVoter(ctx, types.AccountKey("user"+strconv.Itoa(i+1)), c8000)
		num := (i+1)*100 + 1000
		deposit := types.LNO(strconv.Itoa(num))
		valKeys[i] = crypto.GenPrivKeyEd25519().PubKey()
		msg := NewValidatorDepositMsg("user"+strconv.Itoa(i+1), deposit, valKeys[i], "")
		result := handler(ctx, msg)
		assert.Equal(t, sdk.Result{}, result)
	}

	// lowest is user4 with power 1400
	lst, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 21, len(lst.OncallValidators))
	assert.Equal(t, 24, len(lst.AllValidators))
	assert.Equal(t, types.Coin{1400 * types.Decimals}, lst.LowestPower)
	assert.Equal(t, users[3], lst.LowestValidator)

	param, _ := valManager.paramHolder.GetValidatorParam(ctx)
	// punish user4 as missing vote (wont explicitly remove)
	// user3 will become the lowest one with power 1300
	valManager.PunishOncallValidator(ctx, users[3], param.PenaltyMissVote, false)
	lst2, _ := valManager.storage.GetValidatorList(ctx)
	assert.Equal(t, 21, len(lst2.OncallValidators))
	assert.Equal(t, 24, len(lst2.AllValidators))
	assert.Equal(t, types.Coin{1300 * types.Decimals}, lst2.LowestPower)
	assert.Equal(t, users[2], lst2.LowestValidator)

}

func TestGetUpdateValidatorList(t *testing.T) {
	ctx, am, valManager, _, _ := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	user2 := createTestAccount(ctx, am, "user2")

	valKey1 := crypto.GenPrivKeyEd25519().PubKey()
	valKey2 := crypto.GenPrivKeyEd25519().PubKey()

	param, _ := valManager.paramHolder.GetValidatorParam(ctx)

	valManager.RegisterValidator(ctx, user1, valKey1.Bytes(), param.ValidatorMinCommitingDeposit, "")
	valManager.RegisterValidator(ctx, user2, valKey2.Bytes(), param.ValidatorMinCommitingDeposit, "")

	val1, _ := valManager.storage.GetValidator(ctx, user1)
	val2, _ := valManager.storage.GetValidator(ctx, user2)

	val1NoPower := abci.Validator{
		Power:  0,
		PubKey: val1.ABCIValidator.GetPubKey(),
	}

	val2NoPower := abci.Validator{
		Power:  0,
		PubKey: val2.ABCIValidator.GetPubKey(),
	}

	cases := []struct {
		oncallValidators   []types.AccountKey
		preBlockValidators []types.AccountKey
		expectUpdateList   []abci.Validator
	}{
		{[]types.AccountKey{user1}, []types.AccountKey{}, []abci.Validator{val1.ABCIValidator}},
		{[]types.AccountKey{user1, user2}, []types.AccountKey{user1}, []abci.Validator{val1.ABCIValidator, val2.ABCIValidator}},
		{[]types.AccountKey{user1, user2}, []types.AccountKey{user1, user2}, []abci.Validator{val1.ABCIValidator, val2.ABCIValidator}},
		{[]types.AccountKey{user2}, []types.AccountKey{user1, user2}, []abci.Validator{val1NoPower, val2.ABCIValidator}},
		{[]types.AccountKey{}, []types.AccountKey{user2}, []abci.Validator{val2NoPower}},
	}

	for _, cs := range cases {
		lst := &model.ValidatorList{
			OncallValidators:   cs.oncallValidators,
			PreBlockValidators: cs.preBlockValidators,
		}
		valManager.storage.SetValidatorList(ctx, lst)
		actualList, _ := valManager.GetUpdateValidatorList(ctx)
		assert.Equal(t, cs.expectUpdateList, actualList)
	}
}

func TestIsLegalWithdraw(t *testing.T) {
	ctx, am, valManager, _, _ := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	param, _ := valManager.paramHolder.GetValidatorParam(ctx)

	valManager.RegisterValidator(
		ctx, user1, crypto.GenPrivKeyEd25519().PubKey().Bytes(),
		param.ValidatorMinCommitingDeposit.Plus(types.NewCoin(100*types.Decimals)), "")

	cases := []struct {
		oncallValidators []types.AccountKey
		username         types.AccountKey
		withdraw         types.Coin
		expectResult     bool
	}{
		{[]types.AccountKey{}, user1, param.ValidatorMinWithdraw.Minus(types.NewCoin(1)), false},
		{[]types.AccountKey{}, user1, param.ValidatorMinCommitingDeposit, false},
		{[]types.AccountKey{user1}, user1, param.ValidatorMinWithdraw, false},
		{[]types.AccountKey{}, user1, param.ValidatorMinWithdraw, true},
	}

	for _, cs := range cases {
		lst := &model.ValidatorList{
			OncallValidators: cs.oncallValidators,
		}
		valManager.storage.SetValidatorList(ctx, lst)
		res := valManager.IsLegalWithdraw(ctx, cs.username, cs.withdraw)
		assert.Equal(t, cs.expectResult, res)
	}
}
