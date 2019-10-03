package types

import (
	"fmt"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

// DeveloperRegisterMsg - register developer on blockchain
type DeveloperRegisterMsg struct {
	Username    types.AccountKey `json:"username"`
	Website     string           `json:"website"`
	Description string           `json:"description"`
	AppMetaData string           `json:"app_meta_data"`
}

var _ types.Msg = DeveloperRegisterMsg{}

// DeveloperRegisterMsg Msg Implementations
func NewDeveloperRegisterMsg(developer string, website string, description string, appMetaData string) DeveloperRegisterMsg {
	return DeveloperRegisterMsg{
		Username:    types.AccountKey(developer),
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
	if !msg.Username.IsValid() {
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

func (msg DeveloperRegisterMsg) String() string {
	return fmt.Sprintf("DeveloperRegisterMsg{%s}", msg.Username)
}

func (msg DeveloperRegisterMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg DeveloperRegisterMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg DeveloperRegisterMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

func (msg DeveloperRegisterMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// DeveloperUpdateMsg - update developer info on blockchain
type DeveloperUpdateMsg struct {
	Username    types.AccountKey `json:"username"`
	Website     string           `json:"website"`
	Description string           `json:"description"`
	AppMetaData string           `json:"app_meta_data"`
}

var _ types.Msg = DeveloperUpdateMsg{}

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
	if !msg.Username.IsValid() {
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
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg DeveloperUpdateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg DeveloperUpdateMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// DeveloperRevokeMsg - revoke developer on blockchain
type DeveloperRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

var _ types.Msg = DeveloperRevokeMsg{}

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
	if !msg.Username.IsValid() {
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
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg DeveloperRevokeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg DeveloperRevokeMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// GrantPermissionMsg - user grant permission to app
type GrantPermissionMsg struct {
	Username          types.AccountKey `json:"username"`
	AuthorizedApp     types.AccountKey `json:"authorized_app"`
	ValidityPeriodSec int64            `json:"validity_period_second"`
	GrantLevel        types.Permission `json:"grant_level"`
	Amount            types.LNO        `json:"amount"`
}

var _ types.Msg = GrantPermissionMsg{}

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
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}

	if !msg.AuthorizedApp.IsValid() {
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
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg GrantPermissionMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg GrantPermissionMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// RevokePermissionMsg - user revoke permission from app
type RevokePermissionMsg struct {
	Username   types.AccountKey `json:"username"`
	RevokeFrom types.AccountKey `json:"revoke_from"`
	Permission types.Permission `json:"permission"`
}

var _ types.Msg = RevokePermissionMsg{}

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
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}
	if !msg.RevokeFrom.IsValid() {
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
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg RevokePermissionMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg RevokePermissionMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

const (
	// IDA price's unit is 1/10 usd cent. range: [0.001USD, 1USD]
	IDAPriceMin = 1
	IDAPriceMax = 1000
)

// IDAIssueMsg - IDA issue message.
type IDAIssueMsg struct {
	Username types.AccountKey `json:"username"`
	// IDAName  string           `json:"ida_name"`
	IDAPrice int64 `json:"ida_price"`
}

var _ types.Msg = IDAIssueMsg{}

// Route - implements sdk.Msg
func (msg IDAIssueMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg IDAIssueMsg) Type() string { return "IDAIssueMsg" }

// ValidateBasic - implements sdk.Msg
func (msg IDAIssueMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}
	// if len(msg.IDAName) < 3 || len(msg.IDAName) > 10 {
	// 	return ErrInvalidIDAName()
	// }
	// if !allUppercaseLetter(msg.IDAName) {
	// 	return ErrInvalidIDAName()
	// }
	if !(msg.IDAPrice >= IDAPriceMin && msg.IDAPrice <= IDAPriceMax) {
		return ErrInvalidIDAPrice()
	}
	return nil
}

func (msg IDAIssueMsg) String() string {
	return fmt.Sprintf("IDAIssueMsg{%s, %d}", msg.Username, msg.IDAPrice)
}

func (msg IDAIssueMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg IDAIssueMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg IDAIssueMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg IDAIssueMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// IDAMintMsg - Mint more IDA from user pool.
type IDAMintMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

var _ types.Msg = IDAMintMsg{}

// Route - implements sdk.Msg
func (msg IDAMintMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg IDAMintMsg) Type() string { return "IDAMintMsg" }

// ValidateBasic - implements sdk.Msg
func (msg IDAMintMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}
	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg IDAMintMsg) String() string {
	return fmt.Sprintf("IDAMintMsg{username:%v, amount:%v}", msg.Username, msg.Amount)
}

func (msg IDAMintMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg IDAMintMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg IDAMintMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg IDAMintMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// IDATransferMsg - Transfer IDA.
type IDATransferMsg struct {
	App    types.AccountKey `json:"app"`
	Amount types.IDAStr     `json:"amount"`
	From   types.AccountKey `json:"from"`
	To     types.AccountKey `json:"to"`
	Signer types.AccountKey `json:"singer"`
}

var _ types.Msg = IDATransferMsg{}

// Route - implements sdk.Msg
func (msg IDATransferMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg IDATransferMsg) Type() string { return "IDATransferMsg" }

// ValidateBasic - implements sdk.Msg
func (msg IDATransferMsg) ValidateBasic() sdk.Error {
	if !msg.App.IsValid() ||
		!msg.From.IsValid() ||
		!msg.To.IsValid() ||
		!msg.Signer.IsValid() {
		return ErrInvalidUsername()
	}
	if msg.From == msg.To {
		return ErrIDATransferSelf()
	}
	if !(msg.From == msg.App || msg.To == msg.App) {
		return ErrInvalidTransferTarget()
	}
	_, err := msg.Amount.ToMiniIDA()
	if err != nil {
		return err
	}
	return nil
}

func (msg IDATransferMsg) String() string {
	return fmt.Sprintf("IDATransferMsg{%v, %v, %v, %v}", msg.App, msg.Amount, msg.From, msg.To)
}

func (msg IDATransferMsg) GetPermission() types.Permission {
	return types.AppOrAffiliatedPermission
}

// GetSignBytes - implements sdk.Msg
func (msg IDATransferMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg IDATransferMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Signer)}
}

// GetConsumeAmount - implements types.Msg
func (msg IDATransferMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// IDAAuthorizeMsg - update app's permission of IDA of the user.
type IDAAuthorizeMsg struct {
	Username types.AccountKey `json:"username"`
	App      types.AccountKey `json:"app"`
	Activate bool             `json:"activate"`
}

var _ types.Msg = IDAAuthorizeMsg{}

// Route - implements sdk.Msg
func (msg IDAAuthorizeMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg IDAAuthorizeMsg) Type() string { return "IDAAuthorizeMsg" }

// ValidateBasic - implements sdk.Msg
func (msg IDAAuthorizeMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() || !msg.App.IsValid() {
		return ErrInvalidUsername()
	}
	if msg.App == msg.Username {
		return ErrInvalidIDAAuth()
	}
	return nil
}

func (msg IDAAuthorizeMsg) String() string {
	return fmt.Sprintf("IDAAuthorizeMsg{}")
}

func (msg IDAAuthorizeMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg IDAAuthorizeMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg IDAAuthorizeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg IDAAuthorizeMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// UpdateAffiliatedMsg - update affiliate accounts.
type UpdateAffiliatedMsg struct {
	App      types.AccountKey `json:"app"`
	Username types.AccountKey `json:"username"`
	Activate bool             `json:"activate"`
}

var _ types.Msg = UpdateAffiliatedMsg{}

// Route - implements sdk.Msg
func (msg UpdateAffiliatedMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg UpdateAffiliatedMsg) Type() string { return "UpdateAffiliatedMsg" }

// ValidateBasic - implements sdk.Msg
func (msg UpdateAffiliatedMsg) ValidateBasic() sdk.Error {
	if !msg.App.IsValid() || !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg UpdateAffiliatedMsg) String() string {
	return fmt.Sprintf("UpdateAffiliatedMsg{}")
}

func (msg UpdateAffiliatedMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg UpdateAffiliatedMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg UpdateAffiliatedMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.App)}
}

// GetConsumeAmount - implements types.Msg
func (msg UpdateAffiliatedMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// utils
func getSignBytes(msg sdk.Msg) []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// func allUppercaseLetter(s string) bool {
// 	for _, v := range s {
// 		if !(v >= 'A' && v <= 'Z') {
// 			return false
// 		}
// 	}
// 	return true
// }
