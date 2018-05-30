package developer

// nolint
import (
	"encoding/json"
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = DeveloperRegisterMsg{}
var _ sdk.Msg = DeveloperRevokeMsg{}
var _ sdk.Msg = GrantDeveloperMsg{}

type DeveloperRegisterMsg struct {
	Username types.AccountKey `json:"username"`
	Deposit  types.LNO        `json:"deposit"`
}

type DeveloperRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

type GrantDeveloperMsg struct {
	Username        types.AccountKey `json:"username"`
	AuthenticateApp types.AccountKey `json:"authenticate_app"`
	ValidityPeriod  int64            `json:"validity_period"`
	GrantLevel      types.Permission `json:"grant_level"`
}

// DeveloperRegisterMsg Msg Implementations
func NewDeveloperRegisterMsg(developer string, deposit types.LNO) DeveloperRegisterMsg {
	return DeveloperRegisterMsg{
		Username: types.AccountKey(developer),
		Deposit:  deposit,
	}
}

func (msg DeveloperRegisterMsg) Type() string { return types.DeveloperRouterName }

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
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
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

// DeveloperRevokeMsg Msg Implementations
func NewDeveloperRevokeMsg(developer string) DeveloperRevokeMsg {
	return DeveloperRevokeMsg{
		Username: types.AccountKey(developer),
	}
}

func (msg DeveloperRevokeMsg) Type() string { return types.DeveloperRouterName }

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
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
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

// Grant Msg Implementations
func NewGrantDeveloperMsg(user, app string, validityPeriod int64, grantLevel types.Permission) GrantDeveloperMsg {
	return GrantDeveloperMsg{
		Username:        types.AccountKey(user),
		AuthenticateApp: types.AccountKey(app),
		ValidityPeriod:  validityPeriod,
		GrantLevel:      grantLevel,
	}
}

func (msg GrantDeveloperMsg) Type() string { return types.DeveloperRouterName }

func (msg GrantDeveloperMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if len(msg.AuthenticateApp) < types.MinimumUsernameLength ||
		len(msg.AuthenticateApp) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.ValidityPeriod <= 0 {
		return ErrInvalidValidityPeriod()
	}

	return nil
}

func (msg GrantDeveloperMsg) String() string {
	return fmt.Sprintf("GrantDeveloperMsg{User:%v, Grant to App:%v, validity period:%v, grant level:%v}",
		msg.Username, msg.AuthenticateApp, msg.ValidityPeriod, msg.GrantLevel)
}

func (msg GrantDeveloperMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
	return nil
}

func (msg GrantDeveloperMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg GrantDeveloperMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
