package register

import (
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
	"testing"
)

func TestRegisterUsername(t *testing.T) {
	register := "register"
	priv := crypto.GenPrivKeyEd25519()
	msg := NewRegisterMsg(register, priv.PubKey())
	result := msg.ValidateBasic()
	assert.Nil(t, result)

	// Register Username length invalid
	register = "re"
	msg = NewRegisterMsg(register, priv.PubKey())
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername("illeagle length"))

	register = "registerregisterregis"
	msg = NewRegisterMsg(register, priv.PubKey())
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername("illeagle length"))

	// Minimum Length
	register = "reg"
	msg = NewRegisterMsg(register, priv.PubKey())
	result = msg.ValidateBasic()
	assert.Nil(t, result)

	// Maximum Length
	register = "registerregisterregi"
	msg = NewRegisterMsg(register, priv.PubKey())
	result = msg.ValidateBasic()
	assert.Nil(t, result)

	// Illegel character
	registerList := [...]string{"register#", "_register", "-register", "reg@ister",
		"reg*ister", "register!", "register()", "reg$ister",
		"reg=ister", "register^", "register.", "reg$ister,"}
	for _, register := range registerList {
		msg = NewRegisterMsg(register, priv.PubKey())
		result = msg.ValidateBasic()
		assert.Equal(t, result, ErrInvalidUsername("illeagle input"))
	}
}
