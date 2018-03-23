package validator

// nolint
import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tx "github.com/cosmos/cosmos-sdk/x/bank"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

type VoteMsg struct {
	Voter         acc.AccountKey `json:"voter"`
	ValidatorName acc.AccountKey `json:"validator_name"`
	Weight        int64          `json:"weight"`
}

type ValidatorRegisterMsg struct {
	ValidatorName acc.AccountKey `json:"validator_name"`
	PubKey        crypto.PubKey  `json:"new_public_key"`
	Deposit       sdk.Coins      `json:"deposit"`
}

//----------------------------------------
// Vote Msg Implementations

func NewVoteMsg(voter string, validator string, weight int64) VoteMsg {
	msg := VoteMsg{
		Voter:         acc.AccountKey(voter),
		ValidatorName: acc.AccountKey(validator),
		Weight:        weight,
	}
	return msg
}

func (msg VoteMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

func (msg VoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Voter) < types.MinimumUsernameLength ||
		len(msg.Voter) > types.MaximumUsernameLength ||
		len(msg.ValidatorName) < types.MinimumUsernameLength ||
		len(msg.ValidatorName) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	// cannot vote a negative amount of votes
	if msg.Weight <= 0 {
		return tx.ErrInvalidCoins("invalid votes")
	}

	return nil
}

func (msg VoteMsg) String() string {
	return fmt.Sprintf("VoteMsg{Voter:%v, ValidatorName:%v, Votes:%v}",
		msg.Voter, msg.ValidatorName, msg.Weight)
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

func NewValidatorRegisterMsg(validator string, pubkey crypto.PubKey, deposit sdk.Coins) ValidatorRegisterMsg {
	msg := ValidatorRegisterMsg{
		ValidatorName: acc.AccountKey(validator),
		PubKey:        pubkey,
		Deposit:       deposit,
	}
	return msg
}

func (msg ValidatorRegisterMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

func (msg ValidatorRegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.ValidatorName) < types.MinimumUsernameLength ||
		len(msg.ValidatorName) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	return nil
}

func (msg ValidatorRegisterMsg) String() string {
	return fmt.Sprintf("VoteMsg{ValidatorName:%v}", msg.ValidatorName)
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
	return []sdk.Address{sdk.Address(msg.ValidatorName)}
}
