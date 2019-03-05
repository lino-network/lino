package developer

// nolint
import (
	"fmt"
	"unicode/utf8"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = DeveloperRegisterMsg{}
var _ types.Msg = DeveloperUpdateMsg{}
var _ types.Msg = DeveloperRevokeMsg{}
var _ types.Msg = GrantPermissionMsg{}
var _ types.Msg = RevokePermissionMsg{}
var _ types.Msg = PreAuthorizationMsg{}

// DeveloperRegisterMsg - register developer on blockchain
type DeveloperRegisterMsg struct {
	Username    types.AccountKey `json:"username"`
	Deposit     types.LNO        `json:"deposit"`
	Website     string           `json:"website"`
	Description string           `json:"description"`
	AppMetaData string           `json:"app_meta_data"`
}

// DeveloperUpdateMsg - update developer info on blockchain
type DeveloperUpdateMsg struct {
	Username    types.AccountKey `json:"username"`
	Website     string           `json:"website"`
	Description string           `json:"description"`
	AppMetaData string           `json:"app_meta_data"`
}

// DeveloperRevokeMsg - revoke developer on blockchain
type DeveloperRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

// GrantPermissionMsg - user grant permission to app
type GrantPermissionMsg struct {
	Username          types.AccountKey `json:"username"`
	AuthorizedApp     types.AccountKey `json:"authorized_app"`
	ValidityPeriodSec int64            `json:"validity_period_second"`
	GrantLevel        types.Permission `json:"grant_level"`
	Amount            types.LNO        `json:"amount"`
}

// RevokePermissionMsg - user revoke permission from app
type RevokePermissionMsg struct {
	Username   types.AccountKey `json:"username"`
	RevokeFrom types.AccountKey `json:"revoke_from"`
	Permission types.Permission `json:"permission"`
}

// PreAuthorizationMsg - preauth permission to app
type PreAuthorizationMsg struct {
	Username          types.AccountKey `json:"username"`
	AuthorizedApp     types.AccountKey `json:"authorized_app"`
	ValidityPeriodSec int64            `json:"validity_period_second"`
	Amount            types.LNO        `json:"amount"`
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

// Route - implements sdk.Msg
func (msg DeveloperRegisterMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg DeveloperRegisterMsg) Type() string { return "DeveloperRegisterMsg" }

// ValidateBasic - implements sdk.Msg
func (msg DeveloperRegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err
	}

	if len(msg.Website) > types.MaximumLengthOfDeveloperWebsite {
		return ErrInvalidWebsite()
	}

	if utf8.RuneCountInString(msg.Description) > types.MaximumLengthOfDeveloperDesctiption {
		return ErrInvalidDescription()
	}

	if utf8.RuneCountInString(msg.AppMetaData) > types.MaximumLengthOfAppMetadata {
		return ErrInvalidAppMetadata()
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

// GetSigners - implements sdk.Msg
func (msg DeveloperRegisterMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

func (msg DeveloperRegisterMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewDeveloperUpdateMsg - new DeveloperUpdateMsg
func NewDeveloperUpdateMsg(developer string, website string, description string, appMetaData string) DeveloperUpdateMsg {
	return DeveloperUpdateMsg{
		Username:    types.AccountKey(developer),
		Website:     website,
		Description: description,
		AppMetaData: appMetaData,
	}
}

// Route - implements sdk.Msg
func (msg DeveloperUpdateMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg DeveloperUpdateMsg) Type() string { return "DeveloperUpdateMsg" }

// ValidateBasic - implements sdk.Msg
func (msg DeveloperUpdateMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if len(msg.Website) > types.MaximumLengthOfDeveloperWebsite {
		return ErrInvalidWebsite()
	}

	if utf8.RuneCountInString(msg.Description) > types.MaximumLengthOfDeveloperDesctiption {
		return ErrInvalidDescription()
	}

	if utf8.RuneCountInString(msg.AppMetaData) > types.MaximumLengthOfAppMetadata {
		return ErrInvalidAppMetadata()
	}
	return nil
}

func (msg DeveloperUpdateMsg) String() string {
	return fmt.Sprintf(
		"DeveloperUpdateMsg{Username:%v, Website:%v, Description:%v, Metadata:%v}",
		msg.Username, msg.Website, msg.Description, msg.AppMetaData)
}

func (msg DeveloperUpdateMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg DeveloperUpdateMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg DeveloperUpdateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg DeveloperUpdateMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// DeveloperRevokeMsg Msg Implementations
func NewDeveloperRevokeMsg(developer string) DeveloperRevokeMsg {
	return DeveloperRevokeMsg{
		Username: types.AccountKey(developer),
	}
}

// Route - implements sdk.Msg
func (msg DeveloperRevokeMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg DeveloperRevokeMsg) Type() string { return "DeveloperRevokeMsg" }

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

// GetSignBytes - implements sdk.Msg
func (msg DeveloperRevokeMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg DeveloperRevokeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg DeveloperRevokeMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// Grant Msg Implementations
func NewGrantPermissionMsg(
	user, app string, validityPeriodSec int64, grantLevel types.Permission, amount types.LNO) GrantPermissionMsg {
	return GrantPermissionMsg{
		Username:          types.AccountKey(user),
		AuthorizedApp:     types.AccountKey(app),
		ValidityPeriodSec: validityPeriodSec,
		GrantLevel:        grantLevel,
		Amount:            amount,
	}
}

// Route - implements sdk.Msg
func (msg GrantPermissionMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg GrantPermissionMsg) Type() string { return "GrantPermissionMsg" }

// ValidateBasic - implements sdk.Msg
func (msg GrantPermissionMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if len(msg.AuthorizedApp) < types.MinimumUsernameLength ||
		len(msg.AuthorizedApp) > types.MaximumUsernameLength {
		return ErrInvalidAuthorizedApp()
	}

	if msg.ValidityPeriodSec <= 0 ||
		msg.ValidityPeriodSec > types.MaxGranPermValiditySec {
		return ErrInvalidValidityPeriod()
	}

	if msg.GrantLevel == types.ResetPermission ||
		msg.GrantLevel == types.TransactionPermission ||
		msg.GrantLevel == types.GrantAppPermission {
		return ErrGrantPermissionTooHigh()
	}

	if msg.GrantLevel == types.PreAuthorizationPermission ||
		msg.GrantLevel == types.AppAndPreAuthorizationPermission {
		_, err := types.LinoToCoin(msg.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (msg GrantPermissionMsg) String() string {
	return fmt.Sprintf("GrantPermissionMsg{User:%v, Grant to App:%v, validity period:%v, grant level:%v}",
		msg.Username, msg.AuthorizedApp, msg.ValidityPeriodSec, msg.GrantLevel)
}

func (msg GrantPermissionMsg) GetPermission() types.Permission {
	if msg.GrantLevel == types.PreAuthorizationPermission ||
		msg.GrantLevel == types.AppAndPreAuthorizationPermission {
		return types.TransactionPermission
	}
	return types.GrantAppPermission
}

// GetSignBytes - implements sdk.Msg
func (msg GrantPermissionMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg GrantPermissionMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg GrantPermissionMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// Revoke Msg Implementations
func NewRevokePermissionMsg(user string, revokeFrom string, permission int) RevokePermissionMsg {
	return RevokePermissionMsg{
		Username:   types.AccountKey(user),
		RevokeFrom: types.AccountKey(revokeFrom),
		Permission: types.Permission(permission),
	}
}

// Route - implements sdk.Msg
func (msg RevokePermissionMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg RevokePermissionMsg) Type() string { return "RevokePermissionMsg" }

// ValidateBasic - implements sdk.Msg
func (msg RevokePermissionMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	if len(msg.RevokeFrom) < types.MinimumUsernameLength ||
		len(msg.RevokeFrom) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	if msg.Permission == types.ResetPermission ||
		msg.Permission == types.TransactionPermission ||
		msg.Permission == types.GrantAppPermission {
		return ErrInvalidGrantPermission()
	}
	return nil
}

func (msg RevokePermissionMsg) String() string {
	return fmt.Sprintf("RevokePermissionMsg{User:%v, revoke from:%v, permission:%v}",
		msg.Username, msg.RevokeFrom, msg.Permission)
}

func (msg RevokePermissionMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg RevokePermissionMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg RevokePermissionMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg RevokePermissionMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// PreAuthorization Msg Implementations
func NewPreAuthorizationMsg(
	user string, authorizedApp string, validityPeriodSec int64, amount types.LNO) PreAuthorizationMsg {
	return PreAuthorizationMsg{
		Username:          types.AccountKey(user),
		AuthorizedApp:     types.AccountKey(authorizedApp),
		ValidityPeriodSec: validityPeriodSec,
		Amount:            amount,
	}
}

// Route - implements sdk.Msg
func (msg PreAuthorizationMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg PreAuthorizationMsg) Type() string { return "PreAuthorizationMsg" }

// ValidateBasic - implements sdk.Msg
func (msg PreAuthorizationMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	if len(msg.AuthorizedApp) < types.MinimumUsernameLength ||
		len(msg.AuthorizedApp) > types.MaximumUsernameLength {
		return ErrInvalidAuthorizedApp()
	}

	if msg.ValidityPeriodSec <= 0 ||
		msg.ValidityPeriodSec > types.MaxGranPermValiditySec {
		return ErrInvalidValidityPeriod()
	}

	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg PreAuthorizationMsg) String() string {
	return fmt.Sprintf("PreAuthorizationMsg{User:%v, Authorized App:%v, Validate Period:%v, Amount:%v}",
		msg.Username, msg.AuthorizedApp, msg.ValidityPeriodSec, msg.Amount)
}

func (msg PreAuthorizationMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg PreAuthorizationMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg PreAuthorizationMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg PreAuthorizationMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
