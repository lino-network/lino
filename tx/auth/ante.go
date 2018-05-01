package auth

import (
	"bytes"
	"fmt"

	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAnteHandler(am acc.AccountManager, gm global.GlobalManager) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (_ sdk.Context, _ sdk.Result, abort bool) {

		// Assert that there are signatures.
		var sigs = tx.GetSignatures()
		if len(sigs) == 0 {
			return ctx,
				sdk.ErrUnauthorized("no signers").Result(),
				true
		}
		msg := tx.GetMsg()

		stdTx, ok := tx.(sdk.StdTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be sdk.StdTx").Result(), true
		}

		sequences := make([]int64, len(sigs))
		for i := 0; i < len(sigs); i++ {
			sequences[i] = sigs[i].Sequence
		}
		fee := stdTx.Fee
		signBytes := sdk.StdSignBytes(ctx.ChainID(), sequences, fee, msg)
		msgType := msg.Type()

		if msgType == types.RegisterRouterName {
			// TODO(Lino): here we get the address :(
			var signerAddrs = msg.GetSigners()

			// Only new user can sign their own register transaction
			if len(sigs) != len(signerAddrs) || len(sigs) != 1 {
				return ctx, sdk.ErrUnauthorized("wrong number of signers").Result(), true
			}
			if !bytes.Equal(sigs[0].PubKey.Address(), signerAddrs[0]) {
				return ctx, sdk.ErrUnauthorized("wrong public key for signer").Result(), true
			}
			if !sigs[0].PubKey.VerifyBytes(signBytes, sigs[0].Signature) {
				return ctx, sdk.ErrUnauthorized("signature verification failed").Result(), true
			}
			return ctx, sdk.Result{}, false
		}

		signers := msg.GetSigners()
		if len(sigs) < len(signers) {
			return ctx, sdk.ErrUnauthorized("wrong number of signers").Result(), true
		}
		// signers get from msg should be verify first
		for i, signer := range signers {
			accKey, err := am.CheckAuthenticatePubKeyOwner(ctx, types.AccountKey(signer), sigs[i].PubKey)
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
