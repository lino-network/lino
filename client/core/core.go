package core

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	txbuilder "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// BroadcastTx - broadcast the transaction bytes to Tendermint
func (ctx CoreContext) BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if res.CheckTx.Code != uint32(0) {
		return res, errors.Errorf("CheckTx failed: (%d) %s",
			res.CheckTx.Code,
			res.CheckTx.Log)
	}
	if res.DeliverTx.Code != uint32(0) {
		return res, errors.Errorf("DeliverTx failed: (%d) %s",
			res.DeliverTx.Code,
			res.DeliverTx.Log)
	}
	return res, err
}

// Query - query from Tendermint with the provided key and storename
func (ctx CoreContext) Query(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return ctx.query(key, storeName, "key")
}

// QuerySubspace - query from Tendermint with the provided storename and subspace
func (ctx CoreContext) QuerySubspace(cdc *wire.Codec, subspace []byte, storeName string) (res []sdk.KVPair, err error) {
	resRaw, err := ctx.query(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}
	cdc.MustUnmarshalBinaryLengthPrefixed(resRaw, &res)
	return
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

// sign and build the transaction from the msg
func (ctx CoreContext) SignAndBuild(msgs []sdk.Msg, cdc *wire.Codec) ([]byte, error) {
	// build the Sign Messsage from the Standard Message
	chainID := ctx.ChainID
	if chainID == "" {
		return nil, errors.Errorf("Chain ID required but not specified")
	}
	sequence := ctx.Sequence
	memo := ctx.Memo
	signMsg := txbuilder.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: 0,
		Sequence:      sequence,
		Msgs:          msgs,
	}

	// sign and build
	bz := signMsg.Bytes()
	if ctx.PrivKey == nil {
		return nil, errors.New("Must provide private key")
	}
	sig, err := ctx.PrivKey.Sign(bz)
	if err != nil {
		return nil, err
	}
	sigs := []auth.StdSignature{{
		PubKey:    ctx.PrivKey.PubKey(),
		Signature: sig,
		// XXX(yumin): client core may be broken now. we need to revisit this part
		// and probably remove all these and use cosmos's build-in support functions.
		// Sequence:  sequence,
	}}

	// marshal bytes
	tx := auth.NewStdTx(signMsg.Msgs, signMsg.Fee, sigs, memo)
	return cdc.MarshalJSON(tx)
}

// sign and build the transaction from the msg
func (ctx CoreContext) SignBuildBroadcast(
	msgs []sdk.Msg, cdc *wire.Codec) (*ctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := ctx.SignAndBuild(msgs, cdc)
	if err != nil {
		return nil, err
	}
	return ctx.BroadcastTx(txBytes)
}

// get passphrase from std input
func (ctx CoreContext) GetPassphraseFromStdin(name string) (pass string, err error) {
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", name)
	return client.GetPassword(prompt, buf)
}

// GetNode prepares a simple rpc.Client
func (ctx CoreContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("Must define node URI")
	}
	return ctx.Client, nil
}
