package proposal

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/tx/proposal"
	val "github.com/lino-network/lino/tx/validator"
	vote "github.com/lino-network/lino/tx/vote"
	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/go-crypto"
)

func TestForceValidatorVote(t *testing.T) {
	accountTransactionPriv := crypto.GenPrivKeyEd25519()
	accountPostPriv := crypto.GenPrivKeyEd25519()
	accountName := "newUser"
	validatorPriv := crypto.GenPrivKeyEd25519()

	accountTransactionPriv2 := crypto.GenPrivKeyEd25519()
	accountPostPriv2 := crypto.GenPrivKeyEd25519()
	accountName2 := "newUser2"
	validatorPriv2 := crypto.GenPrivKeyEd25519()

	baseTime := time.Now().Unix() + 100
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, accountName, lb, 0,
		crypto.GenPrivKeyEd25519(), accountTransactionPriv, accountPostPriv, "1000000")

	test.CreateAccount(t, accountName2, lb, 1,
		crypto.GenPrivKeyEd25519(), accountTransactionPriv2, accountPostPriv2, "1000000")

	voteDepositMsg := vote.NewVoterDepositMsg(accountName, types.LNO("3000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, accountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(accountName, types.LNO("3000"), validatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, accountTransactionPriv, baseTime)

	voteDepositMsg2 := vote.NewVoterDepositMsg(accountName2, types.LNO("3000"))
	test.SignCheckDeliver(t, lb, voteDepositMsg2, 0, true, accountTransactionPriv2, baseTime)

	valDepositMsg2 := val.NewValidatorDepositMsg(accountName2, types.LNO("3000"), validatorPriv2.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg2, 1, true, accountTransactionPriv2, baseTime)

	test.CheckOncallValidatorList(t, accountName, true, lb)
	test.CheckOncallValidatorList(t, accountName2, true, lb)

	desc := param.GlobalAllocationParam{
		InfraAllocation:          sdk.NewRat(1, 100),
		ContentCreatorAllocation: sdk.NewRat(1, 100),
		DeveloperAllocation:      sdk.NewRat(1, 100),
		ValidatorAllocation:      sdk.NewRat(97, 100),
	}

	changeAllocationMsg := proposal.NewChangeGlobalAllocationParamMsg(accountName, desc)
	test.SignCheckDeliver(t, lb, changeAllocationMsg, 2, true, accountTransactionPriv, baseTime)

	test.CheckBalance(t, accountName, lb, types.NewCoin(894000*types.Decimals))
	test.CheckBalance(t, accountName2, lb, types.NewCoin(994000*types.Decimals))

	test.SimulateOneBlock(lb, baseTime)
	// let validator 1 vote and validator 2 not vote.
	voteMsg := vote.NewVoteMsg(accountName, int64(1), true)
	test.SignCheckDeliver(t, lb, voteMsg, 3, true, accountTransactionPriv, baseTime)

	test.SimulateOneBlock(lb, baseTime+test.ProposalDecideHr*3600+1)
	test.SimulateOneBlock(lb, baseTime+(test.ProposalDecideHr+test.ParamChangeHr)*3600+2)
	test.CheckGlobalAllocation(t, lb, desc)

	// check validator 2 has been punished for not voting
	test.CheckValidatorDeposit(t, accountName, lb, types.NewCoin(3000*types.Decimals))
	test.CheckValidatorDeposit(t, accountName2, lb, types.NewCoin(3000*types.Decimals).Minus(test.PenaltyMissVote))
}
