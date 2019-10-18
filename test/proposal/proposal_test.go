package proposal

// import (
// 	"testing"
// 	"time"

// 	"github.com/lino-network/lino/param"
// 	"github.com/lino-network/lino/test"
// 	"github.com/lino-network/lino/types"
// 	"github.com/lino-network/lino/x/proposal"
// 	val "github.com/lino-network/lino/x/validator"
// 	vote "github.com/lino-network/lino/x/vote"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"
// )

// func TestForceValidatorVote(t *testing.T) {
// 	accountTransactionPriv := secp256k1.GenPrivKey()
// 	accountAppPriv := secp256k1.GenPrivKey()
// 	accountName := "newuser"
// 	validatorPriv := secp256k1.GenPrivKey()

// 	accountTransactionPriv2 := secp256k1.GenPrivKey()
// 	accountAppPriv2 := secp256k1.GenPrivKey()
// 	accountName2 := "newuser2"
// 	validatorPriv2 := secp256k1.GenPrivKey()

// 	baseT := time.Unix(0,0).Add(100 * time.Second)
// 	baseTime := baseT.Unix()
// 	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal, baseT)

// 	totalLNO := types.LNO("1000000000")
// 	depositLNO := types.LNO("3000000")

// 	totalCoin, _ := types.LinoToCoin(totalLNO)
// 	depositCoin, _ := types.LinoToCoin(depositLNO)

// 	test.CreateAccount(t, accountName, lb, 0,
// 		secp256k1.GenPrivKey(), accountTransactionPriv, accountAppPriv, totalLNO)

// 	test.CreateAccount(t, accountName2, lb, 1,
// 		secp256k1.GenPrivKey(), accountTransactionPriv2, accountAppPriv2, totalLNO)

// 	voteDepositMsg := vote.NewStakeInMsg(accountName, depositLNO)
// 	test.SignCheckDeliver(t, lb, voteDepositMsg, 0, true, accountTransactionPriv, baseTime)

// 	valDepositMsg := val.NewValidatorDepositMsg(accountName, depositLNO, validatorPriv.PubKey(), "")
// 	test.SignCheckDeliver(t, lb, valDepositMsg, 1, true, accountTransactionPriv, baseTime)

// 	voteDepositMsg2 := vote.NewStakeInMsg(accountName2, depositLNO)
// 	test.SignCheckDeliver(t, lb, voteDepositMsg2, 0, true, accountTransactionPriv2, baseTime)

// 	valDepositMsg2 := val.NewValidatorDepositMsg(accountName2, depositLNO, validatorPriv2.PubKey(), "")
// 	test.SignCheckDeliver(t, lb, valDepositMsg2, 1, true, accountTransactionPriv2, baseTime)

// 	test.CheckOncallValidatorList(t, accountName, true, lb)
// 	test.CheckOncallValidatorList(t, accountName2, true, lb)

// 	desc := param.GlobalAllocationParam{
// 		GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
// 		InfraAllocation:          types.NewDecFromRat(1, 100),
// 		ContentCreatorAllocation: types.NewDecFromRat(1, 100),
// 		DeveloperAllocation:      types.NewDecFromRat(1, 100),
// 		ValidatorAllocation:      types.NewDecFromRat(97, 100),
// 	}

// 	changeAllocationMsg := proposal.NewChangeGlobalAllocationParamMsg(accountName, desc, "")
// 	test.SignCheckDeliver(t, lb, changeAllocationMsg, 2, true, accountTransactionPriv, baseTime)

// 	accBalance := totalCoin.Minus(depositCoin).Minus(depositCoin).Minus(types.NewCoinFromInt64(1 * types.Decimals))
// 	test.CheckBalance(t, accountName, lb, accBalance.Minus(test.ChangeParamMinDeposit))
// 	test.CheckBalance(t, accountName2, lb, accBalance)

// 	test.SimulateOneBlock(lb, baseTime)
// 	// let validator 1 vote and validator 2 not vote.
// 	voteProposalMsg := proposal.NewVoteProposalMsg(accountName, int64(1), true)
// 	test.SignCheckDeliver(t, lb, voteProposalMsg, 3, true, accountTransactionPriv, baseTime)

// 	test.SimulateOneBlock(lb, baseTime+test.ProposalDecideSec+1)
// 	test.SimulateOneBlock(lb, baseTime+(test.ProposalDecideSec+test.ParamChangeExecutionSec)+2)
// 	test.CheckGlobalAllocation(t, lb, desc)

// 	// check validator 2 has been punished for not voting
// 	test.CheckValidatorDeposit(t, accountName, lb, depositCoin)
// 	test.CheckValidatorDeposit(t, accountName2, lb, depositCoin.Minus(test.PenaltyMissVote))
// }
