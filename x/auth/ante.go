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
			return ctx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}
		// Assert that there are signatures.
		var sigs = stdTx.GetSignatures()
		if len(sigs) == 0 {
			return ctx,
				sdk.ErrUnauthorized("no signers").Result(),
				true
		}
		sdkMsg := tx.GetMsg()
		msg, ok := sdkMsg.(types.Msg)
		if !ok {
			return ctx, sdk.ErrInternal("unrecognize msg").Result(), true
		}

		// Assert that number of signatures is correct.
		var signerAddrs = msg.GetSigners()
		if len(sigs) != len(signerAddrs) {
			return ctx,
				sdk.ErrUnauthorized("wrong number of signers").Result(),
				true
		}

		sequences := make([]int64, len(sigs))
		for i := 0; i < len(sigs); i++ {
			sequences[i] = sigs[i].Sequence
		}
		// for i := 0; i < len(signerAddrs); i++ {
		// 	accNums[i] = sigs[i].AccountNumber
		// }
		fee := stdTx.Fee
		signBytes := auth.StdSignBytes(ctx.ChainID(), []int64{}, sequences, fee, msg)
		// fmt.Println("=========== auth", string(signBytes))

		permission := msg.GetPermission()
		signers := msg.GetSigners()
		if len(sigs) < len(signers) {
			return ctx, sdk.ErrUnauthorized("wrong number of signers").Result(), true
		}
		// signers get from msg should be verify first
		for i, signer := range signers {
			accKey, err := am.CheckAuthenticatePubKeyOwner(ctx, types.AccountKey(signer), sigs[i].PubKey, permission)
			if err != nil {
				return ctx, err.Result(), true
			}

			seq, err := am.GetSequence(ctx, accKey)
			if err != nil {
				return ctx, err.Result(), true
			}
			if seq != sigs[i].Sequence {
				return ctx, sdk.ErrInvalidSequence(
					fmt.Sprintf("Invalid sequence for signer %v. Got %d, expected %d",
						accKey, sigs[i].Sequence, seq)).Result(), true
			}
			if err := am.IncreaseSequenceByOne(ctx, accKey); err != nil {
				return ctx, err.Result(), true
			}

			if !sigs[i].PubKey.VerifyBytes(signBytes, sigs[i].Signature) {
				return ctx, sdk.ErrUnauthorized(
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
