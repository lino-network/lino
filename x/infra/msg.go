package infra

// nolint
import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

var _ types.Msg = ProviderReportMsg{}

type ProviderReportMsg struct {
	Username types.AccountKey `json:"username"`
	Usage    int64            `json:"usage"`
}

//----------------------------------------
// ReportMsg Msg Implementations

func NewProviderReportMsg(provider string, usage int64) ProviderReportMsg {
	return ProviderReportMsg{
		Username: types.AccountKey(provider),
		Usage:    usage,
	}
}

func (msg ProviderReportMsg) Type() string { return types.InfraRouterName } // TODO: "account/register"

func (msg ProviderReportMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.Usage <= 0 {
		return ErrInvalidUsage()
	}

	return nil
}

func (msg ProviderReportMsg) String() string {
	return fmt.Sprintf("ProviderReportMsg{Username:%v, Usage:%v}", msg.Username, msg.Usage)
}

func (msg ProviderReportMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ProviderReportMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ProviderReportMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
