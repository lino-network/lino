package register

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	"github.com/tendermint/go-crypto"
)

func TestRegisterUsername(t *testing.T) {
	register := "register"
	priv := crypto.GenPrivKeyEd25519()
	msg := NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result := msg.ValidateBasic()
	assert.Nil(t, result)

	// Register Username length invalid
	register = "re"
	msg = NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername("illeagle length"))

	register = "registerregisterregis"
	msg = NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername("illeagle length"))

	// Minimum Length
	register = "reg"
	msg = NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result = msg.ValidateBasic()
	assert.Nil(t, result)

	// Maximum Length
	register = "registerregisterregi"
	msg = NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result = msg.ValidateBasic()
	assert.Nil(t, result)

	// Illegel character
	registerList := [...]string{"register#", "_register", "-register", "reg@ister",
		"reg*ister", "register!", "register()", "reg$ister", "reg ister", " register",
		"reg=ister", "register^", "register.", "reg$ister,"}
	for _, register := range registerList {
		msg = NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
		result = msg.ValidateBasic()
		assert.Equal(t, result, ErrInvalidUsername("illeagle input"))
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"provider report msg": {
			NewRegisterMsg("test", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey()),
			types.MasterPermission},
	}

	for testName, cs := range cases {
		permissionLevel := cs.msg.Get(types.PermissionLevel)
		if permissionLevel == nil {
			if cs.expectPermission != types.PostPermission {
				t.Errorf(
					"%s: expect permission incorrect, expect %v, got %v",
					testName, cs.expectPermission, types.PostPermission)
				return
			} else {
				continue
			}
		}
		permission, ok := permissionLevel.(types.Permission)
		assert.Equal(t, ok, true)
		if cs.expectPermission != permission {
			t.Errorf(
				"%s: expect permission incorrect, expect %v, got %v",
				testName, cs.expectPermission, permission)
			return
		}
	}
}
