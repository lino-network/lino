package validator

import (
	"strconv"
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	linotypes "github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	store "github.com/lino-network/lino/x/validator/model"
	valtypes "github.com/lino-network/lino/x/validator/types"
	types "github.com/lino-network/lino/x/vote/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// test validator deposit
func TestValidatorRegister(t *testing.T) {
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountName := "newuser"
	newValidatorPriv := secp256k1.GenPrivKey()

	baseT := time.Unix(0, 0).Add(100 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	test.CreateAccount(t, newAccountName, lb, 0,
		secp256k1.GenPrivKey(), newAccountTransactionPriv, "500000")

	voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("150000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	// deposit the lowest requirement
	valDepositMsg := valtypes.NewValidatorRegisterMsg(
		newAccountName, newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 2, false, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)

	test.SignCheckDeliver(t, lb, voteDepositMsg, 3, true, newAccountTransactionPriv, baseTime)
	test.SignCheckDeliver(t, lb, valDepositMsg, 4, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)

}

func TestRegisterValidatorOneByOne(t *testing.T) {
	testName := "TestRegisterValidatorOneByOne"

	// start with 1 genesis validator
	baseT := time.Unix(0, 0).Add(100 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, 1, baseT)

	// add 21 validators
	seq := 0
	for ; seq < test.DefaultNumOfVal-1; seq++ {
		newAccountResetPriv := secp256k1.GenPrivKey()
		newAccountTransactionPriv := secp256k1.GenPrivKey()

		newValidatorPriv := secp256k1.GenPrivKey()

		newAccountName := "validator"
		newAccountName += strconv.Itoa(seq + 1)

		test.CreateAccount(t, newAccountName, lb, uint64(seq),
			newAccountResetPriv, newAccountTransactionPriv, "500000")

		voteDepositMsg := types.NewStakeInMsg(newAccountName, strconv.Itoa(300000+100*seq))
		test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

		valDepositMsg := valtypes.NewValidatorRegisterMsg(newAccountName, newValidatorPriv.PubKey(), "")
		test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)
		test.CheckOncallValidatorList(t, newAccountName, true, lb)
	}

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	vs := store.NewValidatorStorage(lb.CapKeyValStore)
	lst := vs.GetValidatorList(ctx)

	if len(lst.Oncall) != test.DefaultNumOfVal {
		t.Errorf("%s: diff all validators, got %v, want %v", testName, len(lst.Oncall), test.DefaultNumOfVal)
	}

	// register more validators, but will not be oncall
	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()

	newValidatorPriv := secp256k1.GenPrivKey()

	newAccountName := "validatorx"
	test.CreateAccount(t, newAccountName, lb, uint64(seq),
		newAccountResetPriv, newAccountTransactionPriv, "500000")

	voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("200000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := valtypes.NewValidatorRegisterMsg(
		newAccountName, newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)

	test.CheckOncallValidatorList(t, newAccountName, false, lb)
	test.CheckStandbyValidatorList(t, newAccountName, true, lb)

	// the 22nd validator will be oncall by depositing more money,
	// validator0 will be removed from oncall

	test.SignCheckDeliver(t, lb, voteDepositMsg, 3, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)

	test.CheckOncallValidatorList(t, "validator0", false, lb)
	test.CheckStandbyValidatorList(t, "validator0", true, lb)

}

func TestRemoveTheSameLowestDepositValidator(t *testing.T) {
	// start with 22 genesis validator
	baseT := time.Unix(0, 0).Add(100 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	// Add a new validator who has higher deposit
	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()

	newValidatorPriv := secp256k1.GenPrivKey()

	newAccountName := "validatorx"

	test.CreateAccount(t, newAccountName, lb, 0,
		newAccountResetPriv, newAccountTransactionPriv, "500000")

	voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("300000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

	valDepositMsg := valtypes.NewValidatorRegisterMsg(newAccountName, newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)

	// check removed oncall validator
	test.CheckOncallValidatorList(t, "validator21", false, lb)
}

func TestFireIncompetentValidator(t *testing.T) {
	testName := "TestFireIncompetentValidator"

	// start with 21 genesis validator
	baseT := time.Unix(0, 0).Add(100 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	// ph := param.NewParamHolder(lb.CapKeyParamStore)
	valStore := store.NewValidatorStorage(lb.CapKeyValStore)
	lst := valStore.GetValidatorList(ctx)

	// set validator0 fails to commit
	var signingValidators []abci.VoteInfo
	for _, val := range lst.Oncall {
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

	// check val0 is gone and in jail
	test.CheckOncallValidatorList(t, "validator0", false, lb)
	test.CheckJailValidatorList(t, "validator0", true, lb)
}

func TestFireIncompetentValidator2(t *testing.T) {
	testName := "TestFireIncompetentValidator2"

	// start with 22 genesis validator
	baseT := time.Unix(0, 0).Add(100 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

	// add two more validators but not oncall
	for i := 0; i < 2; i++ {
		newAccountResetPriv := secp256k1.GenPrivKey()
		newAccountTransactionPriv := secp256k1.GenPrivKey()

		newValidatorPriv := secp256k1.GenPrivKey()

		newAccountName := "altval"
		newAccountName += strconv.Itoa(i)

		test.CreateAccount(t, newAccountName, lb, uint64(i),
			newAccountResetPriv, newAccountTransactionPriv, "500000")

		voteDepositMsg := types.NewStakeInMsg(newAccountName, linotypes.LNO("200000"))
		test.SignCheckDeliver(t, lb, voteDepositMsg, 1, true, newAccountTransactionPriv, baseTime)

		valDepositMsg := valtypes.NewValidatorRegisterMsg(newAccountName, newValidatorPriv.PubKey(), "")
		test.SignCheckDeliver(t, lb, valDepositMsg, 2, true, newAccountTransactionPriv, baseTime)

		test.CheckOncallValidatorList(t, newAccountName, false, lb)
		test.CheckStandbyValidatorList(t, newAccountName, true, lb)

	}
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	valStore := store.NewValidatorStorage(lb.CapKeyValStore)

	lst := valStore.GetValidatorList(ctx)
	committingLst := append(lst.Oncall, lst.Standby...)
	// set validator0 fails to commit
	var signingValidators []abci.VoteInfo
	for _, val := range committingLst {
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
	test.CheckJailValidatorList(t, "validator0", true, lb)
	// the votes won't change after slashing
	test.CheckReceivedVotes(t, "validator0", linotypes.NewCoinFromInt64(200000*linotypes.Decimals), lb)

	// check altval1 joins oncall validator, but altval0 not
	test.CheckOncallValidatorList(t, "altval1", true, lb)
	test.CheckStandbyValidatorList(t, "altval1", false, lb)

	test.CheckOncallValidatorList(t, "altval0", false, lb)
	test.CheckStandbyValidatorList(t, "altval0", true, lb)

}
