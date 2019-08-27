package bandwidth

import (
	"testing"
)

func TestAddMsgSignedByUser(t *testing.T) {
	// ctx, bm := setupTest(t, 1)

	// testCases := []struct {
	// 	testName           string
	// 	amount             types.Coin
	// 	from               types.AccountKey
	// 	detailType         types.TransferDetailType
	// 	memo               string
	// 	atWhen             time.Time
	// 	expectCurBlockInfo model.CurBlockInfo
	// }{
	// 	{
	// 		testName:   "add coin to account's saving",
	// 		amount:     c100,
	// 		from:       fromUser1,
	// 		detailType: types.TransferIn,
	// 		memo:       "memo",
	// 		atWhen:     baseTime,
	// 		expectBank: model.AccountBank{
	// 			Saving:  accParam.RegisterFee.Plus(c100),
	// 			CoinDay: accParam.RegisterFee,
	// 		},
	// 		expectPendingCoinDayQueue: model.PendingCoinDayQueue{
	// 			LastUpdatedAt: baseTimeSlot,
	// 			TotalCoinDay:  sdk.ZeroDec(),
	// 			TotalCoin:     c100,
	// 			PendingCoinDays: []model.PendingCoinDay{
	// 				{
	// 					StartTime: baseTimeSlot,
	// 					EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
	// 					Coin:      c100,
	// 				},
	// 			},
	// 		},
	// 	},
	// }
}
