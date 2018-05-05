package proposal

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/proposal/model"
	"github.com/lino-network/lino/types"
)

type CreateProposalMsg struct {
	Creator types.AccountKey `json:"creator"`
	model.ChangeParameterDescription
}

//----------------------------------------
// CreateProposalMsg Msg Implementations

func NewCreateProposalMsg(voter string, para model.ChangeParameterDescription) CreateProposalMsg {
	return CreateProposalMsg{
		Creator:                    types.AccountKey(voter),
		ChangeParameterDescription: para,
	}
}

func (msg CreateProposalMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

func (msg CreateProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.InfraAllocation.
		Add(msg.ContentCreatorAllocation).
		Add(msg.DeveloperAllocation).
		Add(msg.ValidatorAllocation).
		GT(sdk.NewRat(1)) {
		return ErrIllegalParameter()
	}

	if msg.StorageAllocation.
		Add(msg.CDNAllocation).
		GT(sdk.NewRat(1)) {
		return ErrIllegalParameter()
	}
	return nil
}

func (msg CreateProposalMsg) String() string {
	return fmt.Sprintf("CreateProposalMsg{Creator:%v}", msg.Creator)
}

func (msg CreateProposalMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg CreateProposalMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg CreateProposalMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}
