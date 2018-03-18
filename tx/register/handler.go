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
	bank, err := am.GetBankFromAddress(ctx, msg.Address)
	if err != nil {
		return ErrAccRegisterFail("Get bank failed").Result()
	}
	if bank.Username != "" {
		return ErrAccRegisterFail("Already registered").Result()
	}
	if RegisterFee.IsGTE(bank.Coins) {
		return ErrAccRegisterFail("Register Fee Doesn't enough").Result()
	}

	accInfo := types.AccountInfo{
		Username: msg.NewUser,
		Created:  types.Height(ctx.BlockHeight()),
		PostKey:  bank.PubKey,
		OwnerKey: bank.PubKey,
		Address:  msg.Address,
	}
	if err := am.SetInfo(ctx, accInfo.Username, &accInfo); err != nil {
		return ErrAccRegisterFail("Set info failed").Result()
	}

	bank.Username = msg.NewUser
	if err = am.SetBank(ctx, accInfo.Address, bank); err != nil {
		return ErrAccRegisterFail("Set bank failed").Result()
	}

	accMeta := types.AccountMeta{
		LastActivity:   types.Height(ctx.BlockHeight()),
		ActivityBurden: types.DefaultActivityBurden,
		LastABBlock:    types.Height(ctx.BlockHeight()),
	}
	if err := am.SetMeta(ctx, accInfo.Username, &accMeta); err != nil {
		return ErrAccRegisterFail("Set meta failed").Result()
	}

	followers := types.Followers{Followers: []types.AccountKey{}}
	if err := am.SetFollowers(ctx, accInfo.Username, &followers); err != nil {
		return ErrAccRegisterFail("Set accInfo failed").Result()
	}
	followings := types.Followings{Followings: []types.AccountKey{}}
	if err := am.SetFollowings(ctx, accInfo.Username, &followings); err != nil {
		return ErrAccRegisterFail("Set following failed").Result()
	}
	return sdk.Result{}
}
