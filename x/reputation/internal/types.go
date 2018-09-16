package internal

import (
	"math/big"
)

// Terminology
// DP: donation power

type Uid = string
type Pid = string
type Time = int64

// use Unix timestamp instead of Time object to avoid time zone change problem.
// i.e, server change its time zone.
type bigInt = *big.Int
type Stake = bigInt
type Dp = bigInt // donation power
type Rep = bigInt
type RoundId = int64

// used in topN.
type PostDpPair struct {
	Pid   Pid
	SumDp Dp
}
