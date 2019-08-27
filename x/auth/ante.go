package auth

import (
	"fmt"

	"github.com/lino-network/lino/x/bandwidth"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"

	"github.com/cosmos/cosmos-sdk/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	post "github.com/lino-network/lino/x/post"
	vote "github.com/lino-network/lino/x/vote"
	votetypes "github.com/lino-network/lino/x/vote/types"
)

const (
	maxMemoCharacters = 100
)

// GetMsgDonationAmount - return the amount of donation in of @p msg, if not donation, return 0.
// XXX(yumin): outdated, should be remove after Upgrade2.
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

// GetMsgDonationValidAmount - return the min of (amount of donation in of @p msg, saving)
// if not donation, return 0.
func GetMsgDonationValidAmount(ctx sdk.Context, msg types.Msg, am acc.AccountKeeper, pm post.PostKeeper) types.Coin {
	zero := types.NewCoinFromInt64(0)
	donation, ok := msg.(post.DonateMsg)
	if !ok {
		return zero
	}
	rst, err := types.LinoToCoin(donation.Amount)
	if err != nil {
		return zero
	}

	permlink := types.GetPermlink(donation.Author, donation.PostID)
	if !pm.DoesPostExist(ctx, permlink) {
		return zero
	}

	saving, err := am.GetSavingFromUsername(ctx, donation.Username)
	if err != nil {
		return types.NewCoinFromInt64(0)
	}

	// not valid when saving is less than donation amount.
	if rst.IsGT(saving) {
		return types.NewCoinFromInt64(0)
	}
	return rst
}

// NewAnteHandler - return an AnteHandler
func NewAnteHandler(am acc.AccountKeeper, gm global.GlobalManager,
	pm post.PostKeeper, vm vote.VoteKeeper, bm bandwidth.BandwidthKeeper) sdk.AnteHandler {
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
				signer, err := am.CheckSigningPubKeyOwner(ctx, types.AccountKey(msgSigner), sigs[idx].PubKey, permission, consumeAmount)
				if err != nil {
					return ctx, err.Result(), true
				}

				// donationAmount = GetMsgDonationValidAmount(ctx, msg, am, pm)
				// if !donationAmount.IsGTE(types.NewCoinFromInt64(types.NoTPSLimitDonationMin)) {
				// 	// get current tps
				// 	tpsCapacityRatio, err := gm.GetTPSCapacityRatio(ctx)
				// 	if err != nil {
				// 		return ctx, err.Result(), true
				// 	}
				// 	// check user tps capacity
				// 	// if err = am.CheckUserTPSCapacity(ctx, types.AccountKey(msgSigner), tpsCapacityRatio); err != nil {
				// 	// 	return ctx, err.Result(), true
				// 	// }
				// }

				// construct sign bytes and verify sequence number.
				addr, err := am.GetAddress(ctx, types.AccountKey(msgSigner))
				if err != nil {
					return ctx, err.Result(), true
				}
				seq, err := am.GetSequence(ctx, addr)
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
				if err := am.IncreaseSequenceByOne(ctx, addr); err != nil {
					// XXX(yumin): cosmos anth panic here, should we?
					return ctx, err.Result(), true
				}

				signerDuty := vm.GetVoterDuty(ctx, signer)
				// TODO(zhimao): bandwidth model for app signed message
				if signerDuty == votetypes.DutyApp {
					bm.AddMsgSignedByApp(ctx, 1)
				} else {
					// msg fee for general message
					if !bm.IsUserMsgFeeEnough(ctx, fee) {
						return ctx, ErrIncorrectStdTxType().Result(), true
					}

					// TODO(zhimao): minus message fee
					types.NewCoinFromInt64(fee.Amount.AmountOf("LNO").Int64())
					bm.AddMsgSignedByUser(ctx, 1)
				}
				idx++
			}
		}

		// TODO(Lino): verify application signature.
		return ctx, sdk.Result{}, false
	}
}
