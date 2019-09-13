package validator

import (
	"strconv"
	"testing"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	val "github.com/lino-network/lino/x/validator"
	store "github.com/lino-network/lino/x/validator/model"
	vote "github.com/lino-network/lino/x/vote"
	abci "github.com/tendermint/tendermint/abci/types"
)

// test validator deposit
func TestValidatorDeposit(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	newValidatorPriv := secp256k1.GenPrivKey()

	baseTime := time.Now().Unix() + 100
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, newAccountAppPriv, "500000")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	// deposit the lowest requirement
	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("100000"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)

	// deposit as the highest validator
	valDepositMsg = val.NewValidatorDepositMsg(
		newAccountName, types.LNO("100"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)
}

func TestRegisterValidatorOneByOne(t *testing.T) {
	testName := "TestRegisterValidatorOneByOne"

	// start with 1 genesis validator
	lb := test.NewTestLinoBlockchain(t, 1)
	baseTime := time.Now().Unix() + 100

	// add 20 validators
	seq := 0
	for ; seq < test.DefaultNumOfVal-1; seq++ {
		newAccountResetPriv := secp256k1.GenPrivKey()
		newAccountTransactionPriv := secp256k1.GenPrivKey()
		newAccountAppPriv := secp256k1.GenPrivKey()

		newValidatorPriv := secp256k1.GenPrivKey()

		newAccountName := "validator"
		newAccountName += strconv.Itoa(seq + 1)

		test.CreateAccount(t, newAccountName, lb, uint64(seq),
			newAccountResetPriv, newAccountTransactionPriv, newAccountAppPriv, "500000")

		voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
		test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

		valDepositMsg := val.NewValidatorDepositMsg(
			newAccountName, strconv.Itoa(120000+100*seq), newValidatorPriv.PubKey(), "")
		test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
		test.CheckOncallValidatorList(t, newAccountName, true, lb)
		test.CheckAllValidatorList(t, newAccountName, true, lb)
	}

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	lst, err := valManager.GetValidatorList(ctx)
	if err != nil {
		t.Errorf("%s: failed to get validator list, got err %v", testName, err)
	}

	if len(lst.AllValidators) != test.DefaultNumOfVal {
		t.Errorf("%s: diff all validators, got %v, want %v", testName, len(lst.AllValidators), test.DefaultNumOfVal)
	}

	if len(lst.OncallValidators) != test.DefaultNumOfVal {
		t.Errorf("%s: diff oncall validators, got %v, want %v", testName, len(lst.OncallValidators), test.DefaultNumOfVal)
	}

	// register more validators, but will not be oncall
	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()

	newValidatorPriv := secp256k1.GenPrivKey()

	newAccountName := "validatorx"
	test.CreateAccount(t, newAccountName, lb, uint64(seq),
		newAccountResetPriv, newAccountTransactionPriv, newAccountAppPriv, "500000")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("100000"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	test.CheckOncallValidatorList(t, newAccountName, false, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)

	// the 22nd validator will be oncall by depositing more money,
	// validator0 will be removed from oncall
	valDepositMsg = val.NewValidatorDepositMsg(
		newAccountName, types.LNO("1"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)

	test.CheckOncallValidatorList(t, "validator0", false, lb)
	test.CheckAllValidatorList(t, "validator0", true, lb)
}

func TestRemoveTheSameLowestDepositValidator(t *testing.T) {
	// start with 21 genesis validator
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix() + 100

	// Add a new validator who has higher deposit
	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()

	newValidatorPriv := secp256k1.GenPrivKey()

	newAccountName := "validatorx"

	test.CreateAccount(t, newAccountName, lb, 0,
		newAccountResetPriv, newAccountTransactionPriv, newAccountAppPriv, "500000")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("110000"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)

	// check removed oncall validator
	test.CheckOncallValidatorList(t, "validator0", false, lb)
	test.CheckAllValidatorList(t, "validator0", true, lb)
}

func TestFireIncompetentValidator(t *testing.T) {
	testName := "TestFireIncompetentValidator"

	// start with 21 genesis validator
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix() + 100

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	valStore := store.NewValidatorStorage(lb.CapKeyValStore)

	lst, err := valManager.GetValidatorList(ctx)
	assert.Nil(t, err)

	// set validator0 fails to commit
	var signingValidators []abci.VoteInfo
	for _, val := range lst.OncallValidators {
		validator, err := valStore.GetValidator(ctx, val)
		if err != nil {
			t.Errorf("%s: failed to get validator, got err %v", testName, err)
		}
		if validator.AbsentCommit != 1 {
			t.Errorf("%s: expect 1 absent commit for %s, got %v", testName, val, validator.AbsentCommit)
		}

		abciVal := abci.VoteInfo{
			Validator:       validator.ABCIValidator,
			SignedLastBlock: true,
		}

		signingValidators = append(signingValidators, abciVal)
	}
	signingValidators[0].SignedLastBlock = false

	// simulate one block
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			Height:  lb.LastBlockHeight() + 1,
			ChainID: "Lino",
			Time:    time.Unix(baseTime+100, 0),
		},
		LastCommitInfo: abci.LastCommitInfo{
			Votes: signingValidators,
		},
	})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()

	// check validator0 absent commit is 2 (each validator is 1 by default)
	ctx = lb.BaseApp.NewContext(true, abci.Header{})
	valStore = store.NewValidatorStorage(lb.CapKeyValStore)
	val0, err := valStore.GetValidator(ctx, "validator0")
	if err != nil {
		t.Errorf("%s: failed to get validator, got err %v", testName, err)
	}
	if val0.AbsentCommit != 2 {
		t.Errorf("%s: expect 2 absent commit for val0, got %v", testName, val0.AbsentCommit)
	}

	val1, err := valStore.GetValidator(ctx, "validator1")
	if err != nil {
		t.Errorf("%s: failed to get validator, got err %v", testName, err)
	}
	if val1.AbsentCommit != 0 {
		t.Errorf("%s: expect 0 absent commit for val1, got %v", testName, val1.AbsentCommit)
	}

	// set val0 to miss 601 times
	for i := 0; i < 599; i++ {
		lb.BeginBlock(abci.RequestBeginBlock{
			Header: abci.Header{
				Height:  lb.LastBlockHeight() + 1,
				ChainID: "Lino",
				Time:    time.Unix(baseTime+200+int64(i), 0),
			},
			LastCommitInfo: abci.LastCommitInfo{
				Votes: signingValidators,
			},
		})
		lb.EndBlock(abci.RequestEndBlock{})
		lb.Commit()
	}

	// check val0 is gone
	test.CheckOncallValidatorList(t, "validator0", false, lb)
	test.CheckAllValidatorList(t, "validator0", false, lb)
}

func TestFireIncompetentValidatorAndThenAddOneWithHighestDepositAsSupplement(t *testing.T) {
	testName := "TestFireIncompetentValidatorAndThenAddOneWithHighestDepositAsSupplement"

	// start with 21 genesis validator
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix() + 100

	// add two more validators but not oncall
	for i := 0; i < 2; i++ {
		newAccountResetPriv := secp256k1.GenPrivKey()
		newAccountTransactionPriv := secp256k1.GenPrivKey()
		newAccountAppPriv := secp256k1.GenPrivKey()

		newValidatorPriv := secp256k1.GenPrivKey()

		newAccountName := "altval"
		newAccountName += strconv.Itoa(i)

		test.CreateAccount(t, newAccountName, lb, uint64(i),
			newAccountResetPriv, newAccountTransactionPriv, newAccountAppPriv, "500000")

		voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
		test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

		valDepositMsg := val.NewValidatorDepositMsg(
			newAccountName, types.LNO("100000"), newValidatorPriv.PubKey(), "")
		test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
		test.CheckOncallValidatorList(t, newAccountName, false, lb)
		test.CheckAllValidatorList(t, newAccountName, true, lb)
	}

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	valStore := store.NewValidatorStorage(lb.CapKeyValStore)

	lst, err := valManager.GetValidatorList(ctx)
	assert.Nil(t, err)

	// set validator0 fails to commit
	var signingValidators []abci.VoteInfo
	for _, val := range lst.OncallValidators {
		validator, err := valStore.GetValidator(ctx, val)
		if err != nil {
			t.Errorf("%s: failed to get validator, got err %v", testName, err)
		}

		abciVal := abci.VoteInfo{
			Validator:       validator.ABCIValidator,
			SignedLastBlock: true,
		}

		signingValidators = append(signingValidators, abciVal)
	}
	signingValidators[0].SignedLastBlock = false

	// set val0 to miss 601 times
	for i := 0; i < 594; i++ {
		lb.BeginBlock(abci.RequestBeginBlock{
			Header: abci.Header{
				Height:  lb.LastBlockHeight() + 1,
				ChainID: "Lino",
				Time:    time.Unix(baseTime+200+int64(i), 0),
			},
			LastCommitInfo: abci.LastCommitInfo{
				Votes: signingValidators,
			},
		})
		lb.EndBlock(abci.RequestEndBlock{})
		lb.Commit()
	}

	// check val0 is gone
	test.CheckOncallValidatorList(t, "validator0", false, lb)
	test.CheckAllValidatorList(t, "validator0", false, lb)

	// check altval0 joins oncall validator, but altval1 not
	test.CheckOncallValidatorList(t, "altval0", true, lb)
	test.CheckAllValidatorList(t, "altval0", true, lb)

	test.CheckOncallValidatorList(t, "altval1", false, lb)
	test.CheckAllValidatorList(t, "altval1", true, lb)

}

func TestFireIncompetentValidatorAndThenAddOneMoreValidator(t *testing.T) {
	testName := "TestFireIncompetentValidatorAndThenAddOneMoreValidator"

	// start with 21 genesis validator
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)
	baseTime := time.Now().Unix() + 100

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	valStore := store.NewValidatorStorage(lb.CapKeyValStore)

	lst, err := valManager.GetValidatorList(ctx)
	assert.Nil(t, err)

	// set validator0 fails to commit
	var signingValidators []abci.VoteInfo
	for _, val := range lst.OncallValidators {
		validator, err := valStore.GetValidator(ctx, val)
		if err != nil {
			t.Errorf("%s: failed to get validator, got err %v", testName, err)
		}

		abciVal := abci.VoteInfo{
			Validator:       validator.ABCIValidator,
			SignedLastBlock: true,
		}

		signingValidators = append(signingValidators, abciVal)
	}
	signingValidators[0].SignedLastBlock = false

	// simulate one block
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			Height:  lb.LastBlockHeight() + 1,
			ChainID: "Lino",
			Time:    time.Unix(baseTime+100, 0),
		},
		LastCommitInfo: abci.LastCommitInfo{
			Votes: signingValidators,
		},
	})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()

	// set val0 to miss 601 times
	for i := 0; i < 599; i++ {
		lb.BeginBlock(abci.RequestBeginBlock{
			Header: abci.Header{
				Height:  lb.LastBlockHeight() + 1,
				ChainID: "Lino",
				Time:    time.Unix(baseTime+200+int64(i), 0),
			},
			LastCommitInfo: abci.LastCommitInfo{
				Votes: signingValidators,
			},
		})
		lb.EndBlock(abci.RequestEndBlock{})
		lb.Commit()
	}

	// check val0 is gone
	test.CheckOncallValidatorList(t, "validator0", false, lb)
	test.CheckAllValidatorList(t, "validator0", false, lb)

	// add one more validator
	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()

	newValidatorPriv := secp256k1.GenPrivKey()

	newAccountName := "altval"

	test.CreateAccount(t, newAccountName, lb, 0,
		newAccountResetPriv, newAccountTransactionPriv, newAccountAppPriv, "500000")

	voteDepositMsg := vote.NewStakeInMsg(newAccountName, types.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(
		newAccountName, types.LNO("100000"), newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)
	test.CheckAllValidatorList(t, newAccountName, true, lb)
}
