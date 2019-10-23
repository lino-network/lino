package validator

import (
	"strconv"
	"testing"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/test"
	linotypes "github.com/lino-network/lino/types"
	valtypes "github.com/lino-network/lino/x/validator/types"
	types "github.com/lino-network/lino/x/vote/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestRevoke(t *testing.T) {
	// testName := "TestRegisterValidatorOneByOne"

	// start with 1 genesis validator
	baseT := time.Unix(0, 0).Add(100 * time.Second)
	baseTime := baseT.Unix()
	lb := test.NewTestLinoBlockchain(t, 1, baseT)

	// add 1 validator
	newAccountResetPriv := secp256k1.GenPrivKey()
	newAccountTransactionPriv := secp256k1.GenPrivKey()
	newAccountAppPriv := secp256k1.GenPrivKey()

	newValidatorPriv := secp256k1.GenPrivKey()

	newAccountName := "validator"
	newAccountName += strconv.Itoa(1)

	test.CreateAccountWithTime(t, newAccountName, lb, uint64(0),
		newAccountResetPriv, newAccountTransactionPriv, newAccountAppPriv, "500000", baseTime)

	voteDepositMsg := types.NewStakeInMsg(newAccountName, strconv.Itoa(300000))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, newAccountTransactionPriv, baseTime)

	valRegisterMsg := valtypes.NewValidatorRegisterMsg(newAccountName, newValidatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valRegisterMsg, 1, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)

	// create a voter to vote validator1
	newVoterResetPriv := secp256k1.GenPrivKey()
	newVoterTransactionPriv := secp256k1.GenPrivKey()
	newVoterAppPriv := secp256k1.GenPrivKey()
	newVoterName := "voter"

	test.CreateAccountWithTime(t, newVoterName, lb, uint64(1),
		newVoterResetPriv, newVoterTransactionPriv, newVoterAppPriv, "500000", baseTime)
	voterDepositMsg := types.NewStakeInMsg(newVoterName, linotypes.LNO("100000"))
	test.SignCheckDeliver(t, lb, voterDepositMsg, 0, true, newVoterTransactionPriv, baseTime)

	// let voter vote validator1
	voteMsg := valtypes.NewVoteValidatorMsg(newVoterName, []string{"validator1"})
	test.SignCheckDeliver(t, lb, voteMsg, 1, true, newVoterTransactionPriv, baseTime)
	test.CheckReceivedVotes(t, "validator1", linotypes.NewCoinFromInt64(400000*linotypes.Decimals), lb)

	// revoke validator1
	valRevokeMsg := valtypes.NewValidatorRevokeMsg(newAccountName)
	test.SignCheckDeliver(t, lb, valRevokeMsg, 2, true, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)

	// cannot revoke and register in the pending period
	test.SignCheckDeliver(t, lb, valRegisterMsg, 3, false, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)

	test.SignCheckDeliver(t, lb, valRevokeMsg, 4, false, newAccountTransactionPriv, baseTime)
	test.CheckOncallValidatorList(t, newAccountName, false, lb)

	// after pending period, can register again
	ctx := lb.BaseApp.NewContext(true, abci.Header{Time: baseT})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	param := ph.GetValidatorParam(ctx)
	test.SimulateOneBlock(lb, baseTime+param.ValidatorRevokePendingSec+1)
	test.SignCheckDeliver(t, lb, valRegisterMsg, 5, true, newAccountTransactionPriv, baseTime+param.ValidatorRevokePendingSec+2)
	test.CheckOncallValidatorList(t, newAccountName, true, lb)

	// check the validator inherite the previous votes
	test.CheckReceivedVotes(t, "validator1", linotypes.NewCoinFromInt64(400000*linotypes.Decimals), lb)
}
