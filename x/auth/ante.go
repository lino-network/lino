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

// getAccOrAddrSignersFromMsg allows AddrMsg to override signers
func getAccOrAddrSignersFromMsg(msg sdk.Msg) []types.AccOrAddr {
	switch v := msg.(type) {
	case types.AddrMsg:
		return v.GetAccOrAddrSigners()
	default:
		rst := make([]types.AccOrAddr, 0)
		for _, signer := range msg.GetSigners() {
			rst = append(rst, types.NewAccOrAddrFromAcc(types.AccountKey(signer)))
		}
		return rst
	}
}

type msgAndSigs struct {
	msg     sdk.Msg
	signers []types.AccOrAddr
	sigs    []auth.StdSignature
}

func validateAndExtract(stdTx auth.StdTx) ([]msgAndSigs, sdk.Error) {
	// validate memo
	if len(stdTx.GetMemo()) > maxMemoCharacters {
		return nil, sdk.ErrMemoTooLarge(fmt.Sprintf(
			"maximum number of characters is %d but received %d characters",
			maxMemoCharacters, len(stdTx.GetMemo())))
	}

	// validate sigs
	// 1. that there are signatures.
	// 2. no more than limit.
	var sigs = stdTx.GetSignatures()
	if len(sigs) == 0 {
		return nil, ErrNoSignatures()
	}
	if len(sigs) > types.TxSigLimit {
		return nil, sdk.ErrTooManySignatures(fmt.Sprintf(
			"signatures: %d, limit: %d",
			len(sigs), types.TxSigLimit))
	}

	// extract signers
	msgs := stdTx.GetMsgs()
	rst := make([]msgAndSigs, len(msgs))
	for i, msg := range msgs {
		signers := getAccOrAddrSignersFromMsg(msg)
		nSigRequired := len(signers)
		if len(sigs) < nSigRequired {
			return nil, ErrWrongNumberOfSigners()
		}
		rst[i] = msgAndSigs{
			msg:     msg,
			signers: signers,
			sigs:    sigs[:nSigRequired],
		}
		sigs = sigs[nSigRequired:]
	}
	if len(sigs) != 0 {
		return nil, ErrWrongNumberOfSigners()
	}

	return rst, nil
}

func getMsgPermissionAndConsume(msg sdk.Msg) (types.Permission, types.Coin, sdk.Error) {
	var permission types.Permission
	var consumeAmount types.Coin
	switch v := msg.(type) {
	case types.Msg:
		permission = v.GetPermission()
		consumeAmount = v.GetConsumeAmount()
	case types.AddrMsg:
		permission = types.TransactionPermission
		consumeAmount = types.NewCoinFromInt64(0)
	default:
		return permission, consumeAmount, ErrUnknownMsgType()
	}
	return permission, consumeAmount, nil
}

// NewAnteHandler - return an AnteHandler
func NewAnteHandler(am acc.AccountKeeper, bm bandwidth.BandwidthKeeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, sdk.Result, bool) {
		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			return ctx, ErrIncorrectStdTxType().Result(), true
		}
		msgAndSigs, err := validateAndExtract(stdTx)
		if err != nil {
			return ctx, err.Result(), true
		}

		// validate each msg.
		for _, msgSigs := range msgAndSigs {
			permission, consumeAmount, err := getMsgPermissionAndConsume(msgSigs.msg)
			if err != nil {
				return ctx, err.Result(), true
			}
			// validate each signature.
			paid := false
			for i, signer := range msgSigs.signers {
				sig := msgSigs.sigs[i]
				var signerAddr sdk.AccAddress
				var msgSignerAddr sdk.AccAddress
				if signer.IsAddr {
					err := checkAddrSigner(ctx, am, signer.Addr, sig.PubKey, paid)
					if err != nil {
						return ctx, err.Result(), true
					}
					signerAddr = signer.Addr
					msgSignerAddr = signer.Addr
				} else {
					var err sdk.Error
					signerAddr, msgSignerAddr, err = checkAccountSigner(
						ctx, am, signer.AccountKey, sig.PubKey,
						permission, consumeAmount)
					if err != nil {
						return ctx, err.Result(), true
					}

				}

				// 1. verify seq.
				seq, err := am.GetSequence(ctx, msgSignerAddr)
				if err != nil {
					return ctx, err.Result(), true
				}
				// 2. verify signature
				signBytes := auth.StdSignBytes(
					ctx.ChainID(), uint64(0), seq, stdTx.Fee, stdTx.GetMsgs(), stdTx.GetMemo())
				if !sig.PubKey.VerifyBytes(signBytes, sig.Signature) {
					return ctx, ErrUnverifiedBytes(
						fmt.Sprintf("signature verification failed, chain-id:%v, seq:%d",
							ctx.ChainID(), seq)).Result(), true
				}
				// 3. increase seq
				if err := am.IncreaseSequenceByOne(ctx, msgSignerAddr); err != nil {
					return ctx, err.Result(), true
				}

				// 4. only pay fee in the end.
				// only the first signer pays the fee
				if !paid {
					if err := bm.CheckBandwidth(ctx, signerAddr, stdTx.Fee); err != nil {
						return ctx, err.Result(), true
					}
				}
				paid = true
			}
		}

		return ctx, sdk.Result{}, false
	}
}

func checkAddrSigner(ctx sdk.Context, am acc.AccountKeeper, addr sdk.AccAddress, signKey crypto.PubKey, isPaid bool) sdk.Error {
	// if signer is address
	if err := am.CheckSigningPubKeyOwnerByAddress(ctx, addr, signKey, isPaid); err != nil {
		return err
	}
	return nil
}

// this function return the actual signer of the msg (grant permission) and original signer of the msg
func checkAccountSigner(ctx sdk.Context, am acc.AccountKeeper, msgSigner types.AccountKey, signKey crypto.PubKey, permission types.Permission, amount types.Coin) (signerAddr sdk.AccAddress, msgSignerAddr sdk.AccAddress, err sdk.Error) {
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
	return signerAddr, msgSignerAddr, nil
}
