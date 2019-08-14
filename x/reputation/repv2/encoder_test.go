package repv2

import (
	// "fmt"
	// "math/big"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EncoderTestSuite struct {
	suite.Suite
}

// func (suite *EncoderTestSuite) TestUserMeta() {
// 	user := &userMeta{
// 		Consumption       : big.NewInt()
// 		Hold
// 		Reputation
// 		LastSettledRound
// 		LastDonationRound
// 		Unsettled

// 	}
// }

// func BenchmarkEncoder(b *testing.B) {
// 	user := &userMeta2{
// 		Consumption:       Int{big.NewInt(10)},
// 		Hold:              Int{big.NewInt(300)},
// 		Reputation:        Int{big.NewInt(5000)},
// 		LastSettledRound:  3,
// 		LastDonationRound: 5,
// 		Unsettled: []Donation2{
// 			{
// 				Pid:    "1344",
// 				Amount: Int{big.NewInt(3333)},
// 				Impact: Int{big.NewInt(55134)},
// 			},
// 			{
// 				Pid:    "1344",
// 				Amount: Int{big.NewInt(3333)},
// 				Impact: Int{big.NewInt(55134)},
// 			},
// 		},
// 	}

// 	sz := 0
// 	for i := 0; i < b.N; i++ {
// 		dt := encodeUserMeta2(user)
// 		sz += len(dt)
// 		_ = decodeUserMeta2(dt)
// 	}

// 	fmt.Println("xcx", sz)
// }

func TestEncoderTestSuite(t *testing.T) {
	suite.Run(t, new(EncoderTestSuite))
}
