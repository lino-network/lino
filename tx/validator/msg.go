package validator

// nolint
import (
	acc "github.com/lino-network/lino/tx/account"
	"github.com/tendermint/go-crypto"
)

type VoteMsg struct {
	Username    acc.AccountKey `json:"username"`
	ValidatorID ValidatorKey   `json:"validator_id"`
	Weight      int64          `json:"weight"`
}

type RegisterValidatorMsg struct {
	NewValidator ValidatorKey  `json:"new_validator"`
	NewPubKey    crypto.PubKey `json:"new_public_key"`
}
