package auth

import (
	"fmt"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	acc "github.com/lino-network/lino/x/account"
)

func NewAnteHandler(am acc.AccountManager, gm global.GlobalManager) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (_ sdk.Context, _ sdk.Result, abort bool) {
		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			return ctx, ErrIncorrectStdTxType().Result(), true
		}
		// Assert that there are signatures.
		var sigs = stdTx.GetSignatures()
		if len(sigs) == 0 {
			return ctx,
				ErrNoSignatures().Result(),
				true
		}
		sdkMsg := tx.GetMsg()
		msg, ok := sdkMsg.(types.Msg)
		if !ok {
			return ctx, ErrUnknownMsgType().Result(), true
		}

		// Assert that number of signatures is correct.
		var signers = msg.GetSigners()
		if len(sigs) != len(signers) {
			return ctx,
				ErrWrongNumberOfSigners().Result(),
				true
		}

		sequences := make([]int64, len(sigs))
		for i := 0; i < len(sigs); i++ {
			sequences[i] = sigs[i].Sequence
		}
		// for i := 0; i < len(signers); i++ {
		// 	accNums[i] = sigs[i].AccountNumber
		// }
		fee := stdTx.Fee
		signBytes := auth.StdSignBytes(ctx.ChainID(), []int64{}, sequences, fee, msg)
		// fmt.Println("=========== auth", string(signBytes))

		permission := msg.GetPermission()

		// signers get from msg should be verify first
		for i, signer := range signers {
			_, err := am.CheckSigningPubKeyOwner(ctx, types.AccountKey(signer), sigs[i].PubKey, permission)
			if err != nil {
				return ctx, err.Result(), true
			}

			seq, err := am.GetSequence(ctx, types.AccountKey(signer))
			if err != nil {
				return ctx, err.Result(), true
			}
			if seq != sigs[i].Sequence {
				return ctx, ErrInvalidSequence(
					fmt.Sprintf("Invalid sequence for signer %v. Got %d, expected %d",
						types.AccountKey(signer), sigs[i].Sequence, seq)).Result(), true
			}
			if err := am.IncreaseSequenceByOne(ctx, types.AccountKey(signer)); err != nil {
				return ctx, err.Result(), true
			}

			if !sigs[i].PubKey.VerifyBytes(signBytes, sigs[i].Signature) {
				return ctx, ErrUnverifiedBytes(
					fmt.Sprintf("signature verification failed, chain-id:%v", ctx.ChainID())).Result(), true
			}
			tpsCapacityRatio, err := gm.GetTPSCapacityRatio(ctx)
			if err != nil {
				return ctx, err.Result(), true
			}
			if err = am.CheckUserTPSCapacity(ctx, types.AccountKey(signer), tpsCapacityRatio); err != nil {
				return ctx, err.Result(), true
			}
		}

		// TODO(Lino): verify application signature.
		return ctx, sdk.Result{}, false
	}
}
