package repv2

type Uid string
type Pid string

// Time in this package is an int64, unix timestamp, in seconds.
type Time int64

type bigInt = Int

type LinoCoin = bigInt
type IF = bigInt // Impact factor
type Rep = bigInt

type RoundId int64

// used in topN.
type PostIFPair struct {
	Pid   Pid `json:"p"`
	SumIF IF  `json:"s_if"`
}

// merged donation, used in userUnsettled.
type Donation struct {
	Pid    Pid      `json:"p"`
	Amount LinoCoin `json:"a"`
	Impact IF       `json:"i"`
}
