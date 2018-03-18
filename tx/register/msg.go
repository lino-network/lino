package register

import (
	"encoding/json"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/lino-network/lino/types"
)

// RegisterMsg - bind username with address(public key), need to be referred by others (pay for it).
type RegisterMsg struct {
	NewUser types.AccountKey `json:"new_user"`
	Address sdk.Address      `json:"address"`
}

var _ sdk.Msg = RegisterMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewRegisterMsg(newUser string, address sdk.Address) RegisterMsg {
	return RegisterMsg{
		NewUser: types.AccountKey(newUser),
		Address: address,
	}
}

// Implements Msg.
func (msg RegisterMsg) Type() string { return types.RegisterRouterName } // TODO: "account/register"

// Implements Msg.
func (msg RegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.NewUser) < types.MinimumUsernameLength ||
		len(msg.NewUser) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illeagle length")
	}

	match, err := regexp.MatchString(types.UsernameReCheck, string(msg.NewUser))
	if err != nil {
		return ErrInvalidUsername("match error").TraceCause(err, "re error")
	}
	if !match {
		return ErrInvalidUsername("illeagle input")
	}

	if len(msg.Address) == 0 {
		return bank.ErrInvalidAddress(msg.Address.String())
	}
	return nil
}

func (msg RegisterMsg) String() string {
	return fmt.Sprintf("RegisterMsg{Newuser:%v, Address:%v}", msg.NewUser, msg.Address)
}

// Implements Msg.
func (msg RegisterMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg RegisterMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg RegisterMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.NewUser)}
}
