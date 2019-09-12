package core

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// BroadcastTx - broadcast the transaction bytes to Tendermint
func (ctx CoreContext) BroadcastTx(txBytes []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	res, err := node.BroadcastTxCommit(txBytes)
	if err != nil {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}

	if !res.CheckTx.IsOK() {
		return sdk.NewResponseFormatBroadcastTxCommit(res), fmt.Errorf(res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return sdk.NewResponseFormatBroadcastTxCommit(res), fmt.Errorf(res.DeliverTx.Log)
	}

	return sdk.NewResponseFormatBroadcastTxCommit(res), nil
}

// Query - query from Tendermint with the provided key and storename
func (ctx CoreContext) Query(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return ctx.query(key, storeName, "key")
}

// Query from Tendermint with the provided storename and path
func (ctx CoreContext) query(key cmn.HexBytes, storeName, endPath string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	node, err := ctx.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
		Prove:  !ctx.TrustNode,
	}
	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}
	resp := result.Response
	if resp.Code != uint32(0) {
		return res, errors.Errorf("Query failed: (%d) %s", resp.Code, resp.Log)
	}
	return resp.Value, nil
}

func (ctx CoreContext) DoTxPrintResponse(msg sdk.Msg) error {
	if err := msg.ValidateBasic(); err != nil {
		panic(err)
	}

	// build and sign the transaction, then broadcast to Tendermint
	res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg})
	if signErr != nil {
		return signErr
	}

	fmt.Println(res.String())
	return nil
}

// sign and build the transaction from the msg
func (ctx CoreContext) SignBuildBroadcast(msgs []sdk.Msg) (sdk.TxResponse, error) {
	txBytes, err := ctx.BuildAndSign(msgs)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return ctx.BroadcastTx(txBytes)
}

func MakeSignature(msg authtypes.StdSignMsg, pk crypto.PrivKey) (sig authtypes.StdSignature, err error) {
	// sign and build
	bz := msg.Bytes()
	if pk == nil {
		return sig, errors.New("Must provide private key")
	}
	sigBytes, err := pk.Sign(bz)
	if err != nil {
		return sig, err
	}
	sig = authtypes.StdSignature{
		PubKey:    pk.PubKey(),
		Signature: sigBytes,
	}
	return sig, err
}

func (ctx CoreContext) Sign(msg authtypes.StdSignMsg) ([]byte, error) {
	sigs := make([]authtypes.StdSignature, 0)
	for _, pk := range ctx.PrivKeys {
		sig, err := MakeSignature(msg, pk)
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, sig)
	}

	return ctx.TxEncoder(
		authtypes.NewStdTx(msg.Msgs, msg.Fee, sigs, msg.Memo))
}

func (ctx CoreContext) BuildSignMsg(msgs []sdk.Msg) (authtypes.StdSignMsg, error) {
	if ctx.ChainID == "" {
		return authtypes.StdSignMsg{}, fmt.Errorf("chain ID required but not specified")
	}
	fees := ctx.Fees
	return authtypes.StdSignMsg{
		ChainID:       ctx.ChainID,
		AccountNumber: 0,
		Sequence:      ctx.Sequence,
		Memo:          ctx.Memo,
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(1, fees),
	}, nil
}

func (ctx CoreContext) BuildAndSign(msgs []sdk.Msg) ([]byte, error) {
	msg, err := ctx.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	return ctx.Sign(msg)
}

// GetNode prepares a simple rpc.Client
func (ctx CoreContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("Must define node URI")
	}
	return ctx.Client, nil
}
