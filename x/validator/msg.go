package validator

// nolint
import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

var _ sdk.Msg = ValidatorDepositMsg{}
var _ sdk.Msg = ValidatorWithdrawMsg{}
var _ sdk.Msg = ValidatorRevokeMsg{}

type ValidatorDepositMsg struct {
	Username  types.AccountKey `json:"username"`
	Deposit   types.LNO        `json:"deposit"`
	ValPubKey crypto.PubKey    `json:"validator_public_key"`
	Link      string           `json:"link"`
}

type ValidatorWithdrawMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

type ValidatorRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

// ValidatorDepositMsg Msg Implementations
func NewValidatorDepositMsg(validator string, deposit types.LNO, pubKey crypto.PubKey, link string) ValidatorDepositMsg {
	return ValidatorDepositMsg{
		Username:  types.AccountKey(validator),
		Deposit:   deposit,
		ValPubKey: pubKey,
		Link:      link,
	}
}

func (msg ValidatorDepositMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

func (msg ValidatorDepositMsg) ValidateBasic() sdk.Error {
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

func (msg ValidatorDepositMsg) String() string {
	return fmt.Sprintf("ValidatorDepositMsg{Username:%v, Deposit:%v, PubKey:%v}", msg.Username, msg.Deposit, msg.ValPubKey)
}

func (msg ValidatorDepositMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
	return nil
}

func (msg ValidatorDepositMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ValidatorDepositMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

// ValidatorWithdrawMsg Msg Implementations
func NewValidatorWithdrawMsg(validator string, amount types.LNO) ValidatorWithdrawMsg {
	return ValidatorWithdrawMsg{
		Username: types.AccountKey(validator),
		Amount:   amount,
	}
}

func (msg ValidatorWithdrawMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

func (msg ValidatorWithdrawMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg ValidatorWithdrawMsg) String() string {
	return fmt.Sprintf("ValidatorWithdrawMsg{Username:%v}", msg.Username)
}

func (msg ValidatorWithdrawMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
	return nil
}

func (msg ValidatorWithdrawMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ValidatorWithdrawMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

// ValidatorRevokeMsg Msg Implementations
func NewValidatorRevokeMsg(validator string) ValidatorRevokeMsg {
	return ValidatorRevokeMsg{
		Username: types.AccountKey(validator),
	}
}

func (msg ValidatorRevokeMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

func (msg ValidatorRevokeMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg ValidatorRevokeMsg) String() string {
	return fmt.Sprintf("ValidatorRevokeMsg{Username:%v}", msg.Username)
}

func (msg ValidatorRevokeMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
	return nil
}

func (msg ValidatorRevokeMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ValidatorRevokeMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
