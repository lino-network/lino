package internal

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	assert := assert.New(t)
	dt := &userMeta{
		CustomerScore:     big.NewInt(10),
		FreeScore:         big.NewInt(10),
		LastSettled:       1,
		LastDonationRound: 2,
	}

	bytes := encodeUserMeta(dt)
	fmt.Println(bytes)
	rst := decodeUserMeta(bytes)
	assert.Equal(dt, rst)
}
