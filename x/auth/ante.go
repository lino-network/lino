package auth

import (
	"fmt"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"

	"github.com/cosmos/cosmos-sdk/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	post "github.com/lino-network/lino/x/post"
)

const (
	maxMemoCharacters = 100
)

// GetMsgDonationAmount - return the amount of donation in of @p msg, if not donation, return 0.
func GetMsgDonationAmount(msg types.Msg) types.Coin {
	donation, ok := msg.(post.DonateMsg)
	if !ok {
		return types.NewCoinFromInt64(0)
	}
	rst, err := types.LinoToCoin(donation.Amount)
	if err != nil {
		return types.NewCoinFromInt64(0)
	}
	return rst
}

// NewAnteHandler - return an AnteHandler
func NewAnteHandler(am acc.AccountManager, gm global.GlobalManager) sdk.AnteHandler {
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
			consumeAmount := msg.GetConsumeAmount()
			for _, msgSigner := range msgSigners {
				// check public key is valid to sign this msg
				_, err := am.CheckSigningPubKeyOwner(ctx, types.AccountKey(msgSigner), sigs[idx].PubKey, permission, consumeAmount)
				if err != nil {
					return ctx, err.Result(), true
				}
				donationAmount := GetMsgDonationAmount(msg)
				// enable no-cost-donation starting BlockchainUpgrade1Update1Height
				if ctx.BlockHeader().Height < types.BlockchainUpgrade1Update1Height ||
					!donationAmount.IsGTE(types.NewCoinFromInt64(types.NoTPSLimitDonationMin)) {
					// get current tps
					tpsCapacityRatio, err := gm.GetTPSCapacityRatio(ctx)
					if err != nil {
						return ctx, err.Result(), true
					}
					// check user tps capacity
					if err = am.CheckUserTPSCapacity(ctx, types.AccountKey(msgSigner), tpsCapacityRatio); err != nil {
						return ctx, err.Result(), true
					}
				}
				// construct sign bytes and verify sequence number.
				seq, err := am.GetSequence(ctx, types.AccountKey(msgSigner))
				if err != nil {
					return ctx, err.Result(), true
				}
				signBytes := auth.StdSignBytes(ctx.ChainID(), uint64(0), uint64(seq), fee, sdkMsgs, stdTx.GetMemo())
				// verify signature
				if !sigs[idx].PubKey.VerifyBytes(signBytes, sigs[idx].Signature) {
					return ctx, ErrUnverifiedBytes(
						fmt.Sprintf("signature verification failed, chain-id:%v, seq:%d",
							ctx.ChainID(), seq)).Result(), true
				}
				// succ
				if err := am.IncreaseSequenceByOne(ctx, types.AccountKey(msgSigner)); err != nil {
					// XXX(yumin): cosmos anth panic here, should we?
					return ctx, err.Result(), true
				}

				idx++
			}
		}

		// TODO(Lino): verify application signature.
		return ctx, sdk.Result{}, false
	}
}
