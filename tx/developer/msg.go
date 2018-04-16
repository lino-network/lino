package developer

// nolint
import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type DeveloperRegisterMsg struct {
	Username types.AccountKey `json:"username"`
	Deposit  types.LNO        `json:"deposit"`
}

type DeveloperRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

//----------------------------------------
// DeveloperRegisterMsg Msg Implementations

func NewDeveloperRegisterMsg(developer string, deposit types.LNO) DeveloperRegisterMsg {
	return DeveloperRegisterMsg{
		Username: types.AccountKey(developer),
		Deposit:  deposit,
	}
}

func (msg DeveloperRegisterMsg) Type() string { return types.DeveloperRouterName } // TODO: "account/register"

func (msg DeveloperRegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err
	}
	return nil
}

func (msg DeveloperRegisterMsg) String() string {
	return fmt.Sprintf("DeveloperRegisterMsg{Username:%v, Deposit:%v}", msg.Username, msg.Deposit)
}

func (msg DeveloperRegisterMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg DeveloperRegisterMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DeveloperRegisterMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

//----------------------------------------
// DeveloperRevokeMsg Msg Implementations

func NewDeveloperRevokeMsg(developer string) DeveloperRevokeMsg {
	return DeveloperRevokeMsg{
		Username: types.AccountKey(developer),
	}
}

func (msg DeveloperRevokeMsg) Type() string { return types.DeveloperRouterName } // TODO: "account/register"

func (msg DeveloperRevokeMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg DeveloperRevokeMsg) String() string {
	return fmt.Sprintf("DeveloperRevokeMsg{Username:%v}", msg.Username)
}

func (msg DeveloperRevokeMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg DeveloperRevokeMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DeveloperRevokeMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
