package register

import (
	"encoding/json"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

// RegisterMsg - bind username with address(public key), need to be referred by others (pay for it).
type RegisterMsg struct {
	NewUser              types.AccountKey `json:"new_user"`
	NewMasterPubKey      crypto.PubKey    `json:"new_master_public_key"`
	NewPostPubKey        crypto.PubKey    `json:"new_post_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
}

var _ sdk.Msg = RegisterMsg{}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewRegisterMsg(
	newUser string,
	masterPubkey crypto.PubKey,
	postPubkey crypto.PubKey,
	transactionPubkey crypto.PubKey) RegisterMsg {
	return RegisterMsg{
		NewUser:              types.AccountKey(newUser),
		NewMasterPubKey:      masterPubkey,
		NewPostPubKey:        postPubkey,
		NewTransactionPubKey: transactionPubkey,
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
	return nil
}

func (msg RegisterMsg) String() string {
	return fmt.Sprintf("RegisterMsg{Newuser:%v, Master Key:%v, Post Key:%v, Transaction Key:%v}",
		msg.NewUser, msg.NewMasterPubKey, msg.NewPostPubKey, msg.NewTransactionPubKey)
}

// Implements Msg.
func (msg RegisterMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	// the permission will not be checked at auth
	if keyStr == types.PermissionLevel {
		return types.MasterPermission
	}
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
	return []sdk.Address{msg.NewMasterPubKey.Address()}
}
