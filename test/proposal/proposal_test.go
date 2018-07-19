package proposal

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/test"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal"
	val "github.com/lino-network/lino/x/validator"
	vote "github.com/lino-network/lino/x/vote"
	crypto "github.com/tendermint/tendermint/crypto"
)

func TestForceValidatorVote(t *testing.T) {
	accountTransactionPriv := crypto.GenPrivKeyEd25519()
	accountPostPriv := crypto.GenPrivKeyEd25519()
	accountName := "newuser"
	validatorPriv := crypto.GenPrivKeyEd25519()

	accountTransactionPriv2 := crypto.GenPrivKeyEd25519()
	accountPostPriv2 := crypto.GenPrivKeyEd25519()
	accountName2 := "newuser2"
	validatorPriv2 := crypto.GenPrivKeyEd25519()

	baseTime := time.Now().Unix() + 100
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	totalLNO := types.LNO("1000000000")
	depositLNO := types.LNO("3000000")

	totalCoin, _ := types.LinoToCoin(totalLNO)
	depositCoin, _ := types.LinoToCoin(depositLNO)

	test.CreateAccount(t, accountName, lb, 0,
		crypto.GenPrivKeyEd25519(), accountTransactionPriv, crypto.GenPrivKeyEd25519(), accountPostPriv, totalLNO)

	test.CreateAccount(t, accountName2, lb, 1,
		crypto.GenPrivKeyEd25519(), accountTransactionPriv2, crypto.GenPrivKeyEd25519(), accountPostPriv2, totalLNO)

	voteDepositMsg := vote.NewVoterDepositMsg(accountName, depositLNO)
	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, accountTransactionPriv, baseTime)

	valDepositMsg := val.NewValidatorDepositMsg(accountName, depositLNO, validatorPriv.PubKey(), "")
	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, accountTransactionPriv, baseTime)

	voteDepositMsg2 := vote.NewVoterDepositMsg(accountName2, depositLNO)
	test.SignCheckDeliver(t, lb, voteDepositMsg2, 0, true, accountTransactionPriv2, baseTime)

	valDepositMsg2 := val.NewValidatorDepositMsg(accountName2, depositLNO, validatorPriv2.PubKey(), "")
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

	accBalance := totalCoin.Minus(depositCoin).Minus(depositCoin)
	test.CheckBalance(t, accountName, lb, accBalance.Minus(test.ChangeParamMinDeposit))
	test.CheckBalance(t, accountName2, lb, accBalance)

	test.SimulateOneBlock(lb, baseTime)
	// let validator 1 vote and validator 2 not vote.
	voteProposalMsg := proposal.NewVoteProposalMsg(accountName, int64(1), true)
	test.SignCheckDeliver(t, lb, voteProposalMsg, 3, true, accountTransactionPriv, baseTime)

	test.SimulateOneBlock(lb, baseTime+test.ProposalDecideHr*3600+1)
	test.SimulateOneBlock(lb, baseTime+(test.ProposalDecideHr+test.ParamChangeHr)*3600+2)
	test.CheckGlobalAllocation(t, lb, desc)

	// check validator 2 has been punished for not voting
	test.CheckValidatorDeposit(t, accountName, lb, depositCoin)
	test.CheckValidatorDeposit(t, accountName2, lb, depositCoin.Minus(test.PenaltyMissVote))
}
