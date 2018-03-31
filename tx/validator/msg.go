package validator

// nolint
import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

type ValidatorDepositMsg struct {
	Username  acc.AccountKey `json:"username"`
	Deposit   sdk.Coins      `json:"deposit"`
	ValPubKey crypto.PubKey  `json:"validator_public_key"`
}

type ValidatorWithdrawMsg struct {
	Username acc.AccountKey `json:"username"`
}

type ValidatorRevokeMsg struct {
	Username acc.AccountKey `json:"username"`
}

//----------------------------------------
// ValidatorDepositMsg Msg Implementations

func NewValidatorDepositMsg(validator string, deposit sdk.Coins, pubKey crypto.PubKey) ValidatorDepositMsg {
	return ValidatorDepositMsg{
		Username:  acc.AccountKey(validator),
		Deposit:   deposit,
		ValPubKey: pubKey,
	}
}

func (msg ValidatorDepositMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

func (msg ValidatorDepositMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	if !msg.Deposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}
	if !msg.Deposit.IsPositive() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}

	return nil
}

func (msg ValidatorDepositMsg) String() string {
	return fmt.Sprintf("ValidatorDepositMsg{Username:%v, Deposit:%v, PubKey:%v}", msg.Username, msg.Deposit, msg.ValPubKey)
}

func (msg ValidatorDepositMsg) Get(key interface{}) (value interface{}) {
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

//----------------------------------------
// ValidatorWithdrawMsg Msg Implementations

func NewValidatorWithdrawMsg(validator string) ValidatorWithdrawMsg {
	return ValidatorWithdrawMsg{
		Username: acc.AccountKey(validator),
	}
}

func (msg ValidatorWithdrawMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

func (msg ValidatorWithdrawMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg ValidatorWithdrawMsg) String() string {
	return fmt.Sprintf("ValidatorWithdrawMsg{Username:%v}", msg.Username)
}

func (msg ValidatorWithdrawMsg) Get(key interface{}) (value interface{}) {
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

//----------------------------------------
// ValidatorRevokeMsg Msg Implementations

func NewValidatorRevokeMsg(validator string) ValidatorRevokeMsg {
	return ValidatorRevokeMsg{
		Username: acc.AccountKey(validator),
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
