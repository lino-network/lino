package account

import (
	"encoding/json"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

// RegisterMsg - bind username with address(public key), need to be referred by others (pay for it).
type RegisterMsg struct {
	Referrer    types.AccountKey `json:"referrer"`
	NewUser     types.AccountKey `json:"new_user"`
	NewPubKey   crypto.PubKey    `json:"new_public_key"`
	RegisterFee sdk.Coins        `json:"register_fee"`
}

var _ sdk.Msg = RegisterMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewRegisterMsg(referrer, newUser string, pubkey crypto.PubKey, coins sdk.Coins) RegisterMsg {
	return RegisterMsg{
		Referrer:    types.AccountKey(referrer),
		NewUser:     types.AccountKey(newUser),
		NewPubKey:   pubkey,
		RegisterFee: coins,
	}
}

// Implements Msg.
func (msg RegisterMsg) Type() string { return "account" } // TODO: "account/register"

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
	if !msg.RegisterFee.IsValid() {
		return bank.ErrInvalidCoins(msg.RegisterFee.String())
	}
	if !msg.RegisterFee.IsPositive() {
		return bank.ErrInvalidCoins(msg.RegisterFee.String())
	}
	return nil
}

func (msg RegisterMsg) String() string {
	return fmt.Sprintf("RegisterMsg{Referrer:%v -> Newuser:%v, New PubKey:%v, RegisterFee: %v}",
		msg.Referrer, msg.NewUser, msg.NewPubKey, msg.RegisterFee)
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
	return []sdk.Address{sdk.Address(msg.Referrer)}
}
