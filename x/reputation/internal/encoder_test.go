package internal

import (
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
	rst := decodeUserMeta(bytes)
	assert.Equal(dt, rst)
}
