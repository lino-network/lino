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

type signBytesFactory = func(seq uint64) []byte

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

		// signbyte creator returns the bytes that should be signed.
		signBytesCreator := func(seq uint64) []byte {
			return auth.StdSignBytes(
				ctx.ChainID(), uint64(0), seq, stdTx.Fee, stdTx.GetMsgs(), stdTx.GetMemo())
		}

		// validate each msg.
		for _, msgSigs := range msgAndSigs {
			if err := validateMsg(ctx, am, bm, msgSigs, signBytesCreator, stdTx.Fee); err != nil {
				return ctx, err.Result(), true
			}
		}

		return ctx, sdk.Result{}, false
	}
}

func validateMsg(ctx sdk.Context, am acc.AccountKeeper, bm bandwidth.BandwidthKeeper, msgSigs msgAndSigs, signBytesCreator signBytesFactory, fee auth.StdFee) sdk.Error {
	// validate each signature.
	paid := false
	for i, signer := range msgSigs.signers {
		sig := msgSigs.sigs[i]
		var signerAddr sdk.AccAddress
		if signer.IsAddr {
			err := checkAddrSigner(ctx, am, signer.Addr, sig.PubKey, paid)
			if err != nil {
				return err
			}
			signerAddr = signer.Addr
		} else {
			var err sdk.Error
			signerAddr, err = checkAccountSigner(ctx, am, signer.AccountKey, sig.PubKey)
			if err != nil {
				return err
			}
		}

		// 1. verify seq.
		seq, err := am.GetSequence(ctx, signerAddr)
		if err != nil {
			return err
		}
		// 2. verify signature
		signBytes := signBytesCreator(seq)
		if !sig.PubKey.VerifyBytes(signBytes, sig.Signature) {
			return ErrUnverifiedBytes(fmt.Sprintf(
				"signature verification failed, chain-id:%v, seq:%d",
				ctx.ChainID(), seq))
		}
		// 3. increase seq
		if err := am.IncreaseSequenceByOne(ctx, signerAddr); err != nil {
			return err
		}
		// 4. only pay fee in the end.
		// only the first signer pays the fee
		if !paid {
			if err := bm.CheckBandwidth(ctx, signerAddr, fee); err != nil {
				return err
			}
		}
		paid = true
	}
	return nil
}

func checkAddrSigner(ctx sdk.Context, am acc.AccountKeeper, addr sdk.AccAddress, signKey crypto.PubKey, isPaid bool) sdk.Error {
	// if signer is address
	if err := am.CheckSigningPubKeyOwnerByAddress(ctx, addr, signKey, isPaid); err != nil {
		return err
	}
	return nil
}

// this function return the actual signer of the msg.
func checkAccountSigner(ctx sdk.Context, am acc.AccountKeeper, msgSigner types.AccountKey, signKey crypto.PubKey) (signerAddr sdk.AccAddress, err sdk.Error) {
	// check public key is valid to sign this msg
	// return signer is the actual signer of the msg
	signer, err := am.CheckSigningPubKeyOwner(ctx, msgSigner, signKey)
	if err != nil {
		return nil, err
	}
	// get address of actual signer.
	return am.GetAddress(ctx, signer)
}
