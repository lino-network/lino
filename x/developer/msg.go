package developer

// nolint
import (
	"fmt"

	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = DeveloperRegisterMsg{}
var _ types.Msg = DeveloperRevokeMsg{}
var _ types.Msg = GrantPermissionMsg{}
var _ types.Msg = RevokePermissionMsg{}

type DeveloperRegisterMsg struct {
	Username    types.AccountKey `json:"username"`
	Deposit     types.LNO        `json:"deposit"`
	Website     string           `json:"website"`
	Description string           `json:"description"`
	AppMetaData string           `json:"app_meta_data"`
}

type DeveloperRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

type GrantPermissionMsg struct {
	Username        types.AccountKey `json:"username"`
	AuthenticateApp types.AccountKey `json:"authenticate_app"`
	ValidityPeriod  int64            `json:"validity_period"`
	GrantLevel      types.Permission `json:"grant_level"`
	Times           int64            `json:"times"`
}

type RevokePermissionMsg struct {
	Username   types.AccountKey `json:"username"`
	PubKey     crypto.PubKey    `json:"public_key"`
	GrantLevel types.Permission `json:"grant_level"`
}

// DeveloperRegisterMsg Msg Implementations
func NewDeveloperRegisterMsg(developer string, deposit types.LNO, website string, description string, appMetaData string) DeveloperRegisterMsg {
	return DeveloperRegisterMsg{
		Username:    types.AccountKey(developer),
		Deposit:     deposit,
		Website:     website,
		Description: description,
		AppMetaData: appMetaData,
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

func (msg DeveloperRegisterMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg DeveloperRegisterMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DeveloperRegisterMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
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

func (msg DeveloperRevokeMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg DeveloperRevokeMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DeveloperRevokeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Grant Msg Implementations
func NewGrantPermissionMsg(
	user, app string, validityPeriod int64, times int64, grantLevel types.Permission) GrantPermissionMsg {
	return GrantPermissionMsg{
		Username:        types.AccountKey(user),
		AuthenticateApp: types.AccountKey(app),
		ValidityPeriod:  validityPeriod,
		GrantLevel:      grantLevel,
		Times:           times,
	}
}

func (msg GrantPermissionMsg) Type() string { return types.DeveloperRouterName }

func (msg GrantPermissionMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if len(msg.AuthenticateApp) < types.MinimumUsernameLength ||
		len(msg.AuthenticateApp) > types.MaximumUsernameLength {
		return ErrInvalidAuthenticateApp()
	}

	if msg.ValidityPeriod <= 0 {
		return ErrInvalidValidityPeriod()
	}

	if msg.Times < 0 {
		return ErrInvalidGrantTimes()
	}

	if msg.GrantLevel == types.ResetPermission ||
		msg.GrantLevel == types.TransactionPermission ||
		msg.GrantLevel == types.GrantMicropaymentPermission ||
		msg.GrantLevel == types.GrantPostPermission {
		return ErrGrantPermissionTooHigh()
	}

	return nil
}

func (msg GrantPermissionMsg) String() string {
	return fmt.Sprintf("GrantPermissionMsg{User:%v, Grant to App:%v, validity period:%v, grant level:%v}",
		msg.Username, msg.AuthenticateApp, msg.ValidityPeriod, msg.GrantLevel)
}

func (msg GrantPermissionMsg) GetPermission() types.Permission {
	if msg.GrantLevel == types.MicropaymentPermission {
		return types.GrantMicropaymentPermission
	}
	return types.GrantPostPermission
}

func (msg GrantPermissionMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg GrantPermissionMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Revoke Msg Implementations
func NewRevokePermissionMsg(user string, pubKey crypto.PubKey, grantLevel types.Permission) RevokePermissionMsg {
	return RevokePermissionMsg{
		Username:   types.AccountKey(user),
		PubKey:     pubKey,
		GrantLevel: grantLevel,
	}
}

func (msg RevokePermissionMsg) Type() string { return types.DeveloperRouterName }

func (msg RevokePermissionMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.GrantLevel == types.ResetPermission ||
		msg.GrantLevel == types.TransactionPermission ||
		msg.GrantLevel == types.GrantMicropaymentPermission ||
		msg.GrantLevel == types.GrantPostPermission {
		return ErrGrantPermissionTooHigh()
	}

	return nil
}

func (msg RevokePermissionMsg) String() string {
	return fmt.Sprintf("RevokePermissionMsg{User:%v, revoke key:%v, grant level:%v}",
		msg.Username, msg.PubKey, msg.GrantLevel)
}

func (msg RevokePermissionMsg) GetPermission() types.Permission {
	if msg.GrantLevel == types.MicropaymentPermission {
		return types.GrantMicropaymentPermission
	}
	return types.GrantPostPermission
}

func (msg RevokePermissionMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg RevokePermissionMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}
