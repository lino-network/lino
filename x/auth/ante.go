package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/bandwidth"
	"github.com/tendermint/tendermint/crypto"
)

const (
	maxMemoCharacters = 100
)

// NewAnteHandler - return an AnteHandler
func NewAnteHandler(am acc.AccountKeeper, bm bandwidth.BandwidthKeeper) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
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

		memo := stdTx.GetMemo()
		if len(memo) > maxMemoCharacters {
			return ctx,
				sdk.ErrMemoTooLarge(
					fmt.Sprintf("maximum number of characters is %d but received %d characters",
						maxMemoCharacters, len(memo))).Result(),
				true
		}

		fee := stdTx.Fee

		sdkMsgs := tx.GetMsgs()

		var signers []sdk.AccAddress
		for _, msg := range sdkMsgs {
			for _, signer := range msg.GetSigners() {
				signers = append(signers, signer)
			}
		}

		if len(sigs) > types.TxSigLimit {
			return ctx, sdk.ErrTooManySignatures(
					fmt.Sprintf("signatures: %d, limit: %d", len(sigs), types.TxSigLimit),
				).Result(),
				true
		}

		if len(signers) != len(sigs) {
			return ctx,
				ErrWrongNumberOfSigners().Result(),
				true
		}
		// signers get from msg should be verify first
		var idx = 0
		for _, msg := range sdkMsgs {
			msg, ok := msg.(types.Msg)
			if !ok {
				return ctx, ErrUnknownMsgType().Result(), true
			}
			permission := msg.GetPermission()
			msgSigners := msg.GetSigners()
			// Recover msg needs one more signature for new bank address
			consumeAmount := msg.GetConsumeAmount()

			for _, msgSigner := range msgSigners {
				signerAddr, msgSignerAddr, err := getMsgSignerAddrAndSignerAddr(
					ctx, am, types.AccountKey(msgSigner), sigs[idx].PubKey, permission, consumeAmount, idx > 0)
				if err != nil {
					return ctx, err.Result(), true
				}
				seq, err := am.GetSequence(ctx, msgSignerAddr)
				if err != nil {
					return ctx, err.Result(), true
				}
				signBytes := auth.StdSignBytes(ctx.ChainID(), uint64(0), seq, fee, sdkMsgs, stdTx.GetMemo())
				// verify signature
				if !sigs[idx].PubKey.VerifyBytes(signBytes, sigs[idx].Signature) {
					return ctx, ErrUnverifiedBytes(
						fmt.Sprintf("signature verification failed, chain-id:%v, seq:%d",
							ctx.ChainID(), seq)).Result(), true
				}
				// succ
				if err := am.IncreaseSequenceByOne(ctx, msgSignerAddr); err != nil {
					// XXX(yumin): cosmos anth panic here, should we?
					return ctx, err.Result(), true
				}

				// first signer pays the fee
				if idx == 0 {
					if err := bm.CheckBandwidth(ctx, signerAddr, fee); err != nil {
						return ctx, err.Result(), true
					}
				}
				idx++
			}
		}

		// TODO(Lino): verify application signature.
		return ctx, sdk.Result{}, false
	}
}

// this function return the actual signer of the msg (grant permission) and original signer of the msg
func getMsgSignerAddrAndSignerAddr(
	ctx sdk.Context, am acc.AccountKeeper, msgSigner types.AccountKey, signKey crypto.PubKey, permission types.Permission,
	amount types.Coin, isPaid bool) (signerAddr sdk.AccAddress, msgSignerAddr sdk.AccAddress, err sdk.Error) {
	if msgSigner.IsUsername() {
		// if original signer is username
		// check public key is valid to sign this msg
		// return signer is the actual signer of the msg
		signer, err := am.CheckSigningPubKeyOwner(ctx, msgSigner, signKey, permission, amount)
		if err != nil {
			return nil, nil, err
		}
		// get address of actual signer.
		signerAddr, err = am.GetAddress(ctx, signer)
		if err != nil {
			return nil, nil, err
		}
		// get address of original signer.
		msgSignerAddr, err = am.GetAddress(ctx, msgSigner)
		if err != nil {
			return nil, nil, err
		}
	} else {
		// if signer is address
		if err := am.CheckSigningPubKeyOwnerByAddress(ctx, sdk.AccAddress(msgSigner), signKey, isPaid); err != nil {
			return nil, nil, err
		}
		// msg actual signer is the same as original signer
		signerAddr = sdk.AccAddress(msgSigner)
		msgSignerAddr = sdk.AccAddress(msgSigner)
	}
	return
}
