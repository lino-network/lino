package auth

import (
	"bytes"
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
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
		// TODO: can tx just implement message?
		msg := tx.GetMsg()

		sequences := make([]int64, len(sigs))
		for i := 0; i < len(sigs); i++ {
			sequences[i] = sigs[i].Sequence
		}
		signBytes := sdk.StdSignBytes(ctx.ChainID(), sequences, sdk.StdFee{}, msg)

		msgType := msg.Type()

		if msgType == types.RegisterRouterName {
			// TODO(Lino): here we get the address. So ugly :(
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
			seq, err := am.GetSequence(ctx, types.AccountKey(signer))
			if err != nil {
				return ctx, err.Result(), true
			}
			if seq != sigs[i].Sequence {
				return ctx, sdk.ErrInvalidSequence(
						fmt.Sprintf("Invalid sequence. Got %d, expected %d", sigs[i].Sequence, seq)).Result(),
					true
			}
			if err := am.IncreaseSequenceByOne(ctx, types.AccountKey(signer)); err != nil {
				return ctx, err.Result(), true
			}

			pubKey, err := am.GetOwnerKey(ctx, types.AccountKey(signer))
			if err != nil {
				return ctx, err.Result(), true
			}
			// TODO(Lino): match postkey and owner key.
			if !reflect.DeepEqual(*pubKey, sigs[i].PubKey) {
				return ctx, sdk.ErrUnauthorized("signer mismatch").Result(), true
			}
			if !sigs[i].PubKey.VerifyBytes(signBytes, sigs[i].Signature) {
				return ctx, sdk.ErrUnauthorized("signature verification failed").Result(), true
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
