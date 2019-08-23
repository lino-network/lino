package repv2

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntTestSuite struct {
	suite.Suite
}

func (suite *IntTestSuite) TestNewInt() {
	suite.Equal(big.NewInt(334124), NewInt(334124).Int)
}

func (suite *IntTestSuite) TestClone() {
	a := NewInt(333)
	b := a.Clone()
	b.Add(NewInt(123))
	suite.Equal(NewInt(456), b)
	suite.Equal(NewInt(333), a)
}

func (suite *IntTestSuite) TestMisc() {
	a := NewInt(33)
	a.Add(NewInt(44))
	suite.Equal(NewInt(77), a)

	b := NewInt(44)
	b.Sub(NewInt(33))
	suite.Equal(NewInt(11), b)

	c := NewInt(55)
	c.Mul(NewInt(100))
	suite.Equal(NewInt(5500), c)

	d := NewInt(55)
	d.Div(NewInt(11))
	suite.Equal(NewInt(5), d)

	suite.Equal(NewInt(444), IntMax(NewInt(333), NewInt(444)))
	suite.Equal(NewInt(333), IntMin(NewInt(333), NewInt(444)))
	suite.Equal(NewInt(777), IntAdd(NewInt(333), NewInt(444)))
	suite.Equal(NewInt(-111), IntSub(NewInt(333), NewInt(444)))
	suite.Equal(NewInt(8547), IntMul(NewInt(111), NewInt(77)))
	suite.Equal(NewInt(3), IntDiv(NewInt(333), NewInt(111)))
	suite.Equal(true, IntGreater(NewInt(444), NewInt(111)))
	suite.Equal(false, IntGreater(NewInt(1), NewInt(111)))
	suite.Equal(true, IntGTE(NewInt(333), NewInt(111)))
	suite.Equal(true, IntGTE(NewInt(333), NewInt(333)))
	suite.Equal(false, IntGTE(NewInt(111), NewInt(333)))
	suite.Equal(true, IntLess(NewInt(111), NewInt(333)))
	suite.Equal(false, IntLess(NewInt(444), NewInt(333)))
	suite.Equal(NewInt(3), IntDiv(NewInt(333), NewInt(111)))
	suite.Equal(NewInt(3), IntDiv(NewInt(333), NewInt(111)))
}

func (suite *IntTestSuite) TestIntDivFrac() {
	suite.Panics(func() { IntDivFrac(NewInt(1000), 1, 0) })
	suite.Panics(func() { IntDivFrac(NewInt(1000), 0, 1) })
	suite.Panics(func() { IntDivFrac(NewInt(1000), 0, 0) })
	cases := []struct {
		v        Int
		num      int64
		denum    int64
		expected Int
	}{
		{
			v:        NewInt(80),
			num:      8,
			denum:    10,
			expected: NewInt(100),
		},
		{
			v:        NewInt(100),
			num:      1,
			denum:    3,
			expected: NewInt(300),
		},
		{
			v:        NewInt(77),
			num:      11,
			denum:    7,
			expected: NewInt(49),
		},
	}

	for i, v := range cases {
		suite.Equal(v.expected, IntDivFrac(v.v, v.num, v.denum), "case: %d", i)
	}
}

func (suite *IntTestSuite) TestIntMulFrac() {
	suite.Panics(func() { IntMulFrac(NewInt(1000), 1, 0) })
	cases := []struct {
		v        Int
		num      int64
		denum    int64
		expected Int
	}{
		{
			v:        NewInt(80),
			num:      8,
			denum:    10,
			expected: NewInt(64),
		},
		{
			v:        NewInt(100),
			num:      1,
			denum:    3,
			expected: NewInt(33),
		},
		{
			v:        NewInt(77),
			num:      11,
			denum:    7,
			expected: NewInt(121),
		},
	}

	for i, v := range cases {
		suite.Equal(v.expected, IntMulFrac(v.v, v.num, v.denum), "case: %d", i)
	}
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, new(IntTestSuite))
}
