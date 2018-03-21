package register

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

var RegisterFee = sdk.Coins{sdk.Coin{Denom: "Lino", Amount: 100}}

func NewHandler(am types.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case RegisterMsg:
			return handleRegisterMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleRegisterMsg(ctx sdk.Context, am types.AccountManager, msg RegisterMsg) sdk.Result {
	if am.AccountExist(ctx, msg.NewUser) {
		return ErrAccRegisterFail("Username exist").Result()
	}
	bank, err := am.GetBankFromAddress(ctx, msg.NewPubKey.Address())
	if err != nil {
		return ErrAccRegisterFail("Get bank failed").Result()
	}
	if bank.Username != "" {
		return ErrAccRegisterFail("Already registered").Result()
	}
	if RegisterFee.IsGTE(bank.Balance) {
		return ErrAccRegisterFail("Register Fee Doesn't enough").Result()
	}

	if _, err := am.CreateAccount(ctx, msg.NewUser, msg.NewPubKey, bank); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
