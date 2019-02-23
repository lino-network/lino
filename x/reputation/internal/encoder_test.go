package internal

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
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
	rst := decodeUserMeta(bytes)
	assert.Equal(dt, rst)
}
