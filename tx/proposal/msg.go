package proposal

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/proposal/model"
	"github.com/lino-network/lino/types"
)

type ChangeParamMsg interface {
	GetDescription() model.Description
	GetCreator() types.AccountKey
}

type ContentCensorshipMsg interface {
	GetCreator() types.AccountKey
	GetPermLink() types.PermLink
}

type ProtocolUpgradeMsg interface {
	GetCreator() types.AccountKey
}

type ChangeGlobalAllocationMsg struct {
	Creator     types.AccountKey            `json:"creator"`
	Description param.GlobalAllocationParam `json:"description"`
}

//----------------------------------------
// ChangeGlobalAllocationMsg Msg Implementations

func NewChangeGlobalAllocationMsg(creator string, desc param.GlobalAllocationParam) ChangeGlobalAllocationMsg {
	return ChangeGlobalAllocationMsg{
		Creator:     types.AccountKey(creator),
		Description: desc,
	}
}

func (msg ChangeGlobalAllocationMsg) GetDescription() model.Description { return msg.Description }
func (msg ChangeGlobalAllocationMsg) GetCreator() types.AccountKey      { return msg.Creator }
func (msg ChangeGlobalAllocationMsg) Type() string                      { return types.ProposalRouterName }

func (msg ChangeGlobalAllocationMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if !msg.Description.InfraAllocation.
		Add(msg.Description.ContentCreatorAllocation).
		Add(msg.Description.DeveloperAllocation).
		Add(msg.Description.ValidatorAllocation).Equal(sdk.NewRat(1)) {
		return ErrIllegalParameter()
	}

	return nil
}

func (msg ChangeGlobalAllocationMsg) String() string {
	return fmt.Sprintf("ChangeGlobalAllocationMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeGlobalAllocationMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeGlobalAllocationMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeGlobalAllocationMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}
