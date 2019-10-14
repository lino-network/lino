package types

// nolint
import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

var _ types.Msg = ValidatorRegisterMsg{}
var _ types.Msg = ValidatorRevokeMsg{}
var _ types.Msg = VoteValidatorMsg{}
var _ types.Msg = ValidatorUpdateMsg{}

// ValidatorRegisterMsg - register to become validator
type ValidatorRegisterMsg struct {
	Username  types.AccountKey `json:"username"`
	ValPubKey crypto.PubKey    `json:"validator_public_key"`
	Link      string           `json:"link"`
}

// ValidatorRegisterMsg Msg Implementations
func NewValidatorRegisterMsg(validator string, pubKey crypto.PubKey, link string) ValidatorRegisterMsg {
	return ValidatorRegisterMsg{
		Username:  types.AccountKey(validator),
		ValPubKey: pubKey,
		Link:      link,
	}
}

// Route - implement sdk.Msg
func (msg ValidatorRegisterMsg) Route() string { return RouterKey }

// Type - implement sdk.Msg
func (msg ValidatorRegisterMsg) Type() string { return "ValidatorRegisterMsg" }

// ValidateBasic - implement sdk.Msg
func (msg ValidatorRegisterMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}

	if len(msg.Link) > types.MaximumLinkURL {
		return ErrInvalidWebsite()
	}

	return nil
}

func (msg ValidatorRegisterMsg) String() string {
	return fmt.Sprintf("ValidatorRegisterMsg{Username:%v, PubKey:%v}", msg.Username, msg.ValPubKey)
}

// GetPermission - implement types.Msg
func (msg ValidatorRegisterMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ValidatorRegisterMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implement sdk.Msg
func (msg ValidatorRegisterMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg ValidatorRegisterMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// ValidatorRevokeMsg - revoke validator
type ValidatorRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

// ValidatorRevokeMsg Msg Implementations
func NewValidatorRevokeMsg(validator string) ValidatorRevokeMsg {
	return ValidatorRevokeMsg{
		Username: types.AccountKey(validator),
	}
}

// Route - implement sdk.Msg
func (msg ValidatorRevokeMsg) Route() string { return RouterKey }

// Type - implement sdk.Msg
func (msg ValidatorRevokeMsg) Type() string { return "ValidatorRevokeMsg" }

// ValidateBasic - implement sdk.Msg
func (msg ValidatorRevokeMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg ValidatorRevokeMsg) String() string {
	return fmt.Sprintf("ValidatorRevokeMsg{Username:%v}", msg.Username)
}

// GetPermission - implement types.Msg
func (msg ValidatorRevokeMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ValidatorRevokeMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implement sdk.Msg
func (msg ValidatorRevokeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg ValidatorRevokeMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// VoteValidatorMsg - vote for validator
type VoteValidatorMsg struct {
	Username        types.AccountKey   `json:"username"`
	VotedValidators []types.AccountKey `json:"voted_validators"`
}

// VoteValidatorMsg Msg Implementations
func NewVoteValidatorMsg(username string, votedValidators []string) VoteValidatorMsg {
	var votedVals []types.AccountKey
	for _, val := range votedValidators {
		votedVals = append(votedVals, types.AccountKey(val))
	}
	return VoteValidatorMsg{
		Username:        types.AccountKey(username),
		VotedValidators: votedVals,
	}
}

// Route - implement sdk.Msg
func (msg VoteValidatorMsg) Route() string { return RouterKey }

// Type - implement sdk.Msg
func (msg VoteValidatorMsg) Type() string { return "VoteValidatorMsg" }

// ValidateBasic - implement sdk.Msg
func (msg VoteValidatorMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}

	if len(msg.VotedValidators) > types.MaxVotedValidators ||
		len(msg.VotedValidators) == 0 {
		return ErrInvalidVotedValidators()
	}

	for _, val := range msg.VotedValidators {
		if !val.IsValid() {
			return ErrInvalidVotedValidators()
		}
	}

	return nil
}

func (msg VoteValidatorMsg) String() string {
	return fmt.Sprintf("VoteValidatorMsg{Username:%v, VotedValidators:%v}", msg.Username, msg.VotedValidators)
}

// GetPermission - implement types.Msg
func (msg VoteValidatorMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg VoteValidatorMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implement sdk.Msg
func (msg VoteValidatorMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg VoteValidatorMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// ValidatorUpdateMsg - register to become validator
type ValidatorUpdateMsg struct {
	Username types.AccountKey `json:"username"`
	Link     string           `json:"link"`
}

// ValidatorUpdateMsg Msg Implementations
func NewValidatorUpdateMsg(validator string, link string) ValidatorUpdateMsg {
	return ValidatorUpdateMsg{
		Username: types.AccountKey(validator),
		Link:     link,
	}
}

// Route - implement sdk.Msg
func (msg ValidatorUpdateMsg) Route() string { return RouterKey }

// Type - implement sdk.Msg
func (msg ValidatorUpdateMsg) Type() string { return "ValidatorUpdateMsg" }

// ValidateBasic - implement sdk.Msg
func (msg ValidatorUpdateMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}

	if len(msg.Link) > types.MaximumLinkURL {
		return ErrInvalidWebsite()
	}

	return nil
}

func (msg ValidatorUpdateMsg) String() string {
	return fmt.Sprintf("ValidatorUpdateMsg{Username:%v, Link:%v}", msg.Username, msg.Link)
}

// GetPermission - implement types.Msg
func (msg ValidatorUpdateMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ValidatorUpdateMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implement sdk.Msg
func (msg ValidatorUpdateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg ValidatorUpdateMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// utils
func getSignBytes(msg sdk.Msg) []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
