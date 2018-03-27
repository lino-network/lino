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

type VoteMsg struct {
	Voter         acc.AccountKey `json:"voter"`
	ValidatorName acc.AccountKey `json:"validator_name"`
	Power         sdk.Coins      `json:"power"`
}

type ValidatorRegisterMsg struct {
	Username  acc.AccountKey `json:"username"`
	Deposit   sdk.Coins      `json:"deposit"`
	ValPubKey crypto.PubKey  `json:"validator_public_key"`
}

//----------------------------------------
// Vote Msg Implementations

func NewVoteMsg(voter string, validator string, power sdk.Coins) VoteMsg {
	return VoteMsg{
		Voter:         acc.AccountKey(voter),
		ValidatorName: acc.AccountKey(validator),
		Power:         power,
	}
}

func (msg VoteMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

func (msg VoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Voter) < types.MinimumUsernameLength ||
		len(msg.Voter) > types.MaximumUsernameLength ||
		len(msg.ValidatorName) < types.MinimumUsernameLength ||
		len(msg.ValidatorName) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	// cannot vote a negative amount of votes
	if !msg.Power.IsValid() {
		return sdk.ErrInvalidCoins(msg.Power.String())
	}
	if !msg.Power.IsPositive() {
		return sdk.ErrInvalidCoins("invalid votes")
	}

	return nil
}

func (msg VoteMsg) String() string {
	return fmt.Sprintf("VoteMsg{Voter:%v, ValidatorName:%v, Votes:%v}",
		msg.Voter, msg.ValidatorName, msg.Power)
}

func (msg VoteMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg VoteMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg VoteMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Voter)}
}

//----------------------------------------
// RegisterValidatorMsg Msg Implementations

func NewValidatorRegisterMsg(validator string, deposit sdk.Coins, pubKey crypto.PubKey) ValidatorRegisterMsg {
	return ValidatorRegisterMsg{
		Username:  acc.AccountKey(validator),
		Deposit:   deposit,
		ValPubKey: pubKey,
	}
}

func (msg ValidatorRegisterMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

func (msg ValidatorRegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}
	if !msg.Deposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}
	if !msg.Deposit.IsPositive() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}

	return nil
}

func (msg ValidatorRegisterMsg) String() string {
	return fmt.Sprintf("ValidatorRegisterMsg{Username:%v, Deposit:%v, PubKey:%v}", msg.Username, msg.Deposit, msg.ValPubKey)
}

func (msg ValidatorRegisterMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ValidatorRegisterMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ValidatorRegisterMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
