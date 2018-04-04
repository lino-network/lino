package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
)

type ReturnCoinEvent struct {
	Username acc.AccountKey `json:"username"`
	Amount   types.Coin     `json:"amount"`
}

type DecideProposalEvent struct {
	ProposalID ProposalKey `json:"proposal_id"`
}

func (event ReturnCoinEvent) Execute(ctx sdk.Context, vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	account := acc.NewProxyAccount(event.Username, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound()
	}

	if err := account.AddCoin(ctx, event.Amount); err != nil {
		return err
	}
	if err := account.Apply(ctx); err != nil {
		return err
	}

	return nil
}

// func (event DecideProposalEvent) Execute(ctx sdk.Context, vm VoteManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
//   votes, getErr := vm
// 	return nil
// }
