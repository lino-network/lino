package core

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	cmn "github.com/tendermint/tmlibs/common"
)

// Broadcast the transaction bytes to Tendermint
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

// Query from Tendermint with the provided key and storename
func (ctx CoreContext) Query(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return ctx.query(key, storeName, "key")
}

// Query from Tendermint with the provided storename and subspace
func (ctx CoreContext) QuerySubspace(cdc *wire.Codec, subspace []byte, storeName string) (res []sdk.KVPair, err error) {
	resRaw, err := ctx.query(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}
	cdc.MustUnmarshalBinary(resRaw, &res)
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
		Height:  ctx.Height,
		Trusted: ctx.TrustNode,
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

// Get the from address from the name flag
func (ctx CoreContext) GetFromAddress() (from sdk.Address, err error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	name := ctx.FromAddressName
	if name == "" {
		return nil, errors.Errorf("must provide a from address name")
	}

	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.Errorf("No key for: %s", name)
	}

	return info.PubKey.Address(), nil
}

// sign and build the transaction from the msg
func (ctx CoreContext) SignAndBuild(msg sdk.Msg, cdc *wire.Codec) ([]byte, error) {
	// build the Sign Messsage from the Standard Message
	chainID := ctx.ChainID
	if chainID == "" {
		return nil, errors.Errorf("Chain ID required but not specified")
	}
	sequence := ctx.Sequence
	signMsg := auth.StdSignMsg{
		ChainID:   chainID,
		Sequences: []int64{sequence},
		Msg:       msg,
	}

	// sign and build
	bz := signMsg.Bytes()
	if ctx.PrivKey == nil {
		return nil, errors.New("Must provide private key")
	}
	sig := ctx.PrivKey.Sign(bz)
	sigs := []auth.StdSignature{{
		PubKey:    ctx.PrivKey.PubKey(),
		Signature: sig,
		Sequence:  sequence,
	}}

	// marshal bytes
	tx := auth.NewStdTx(signMsg.Msg, signMsg.Fee, sigs)
	return cdc.MarshalJSON(tx)
}

// sign and build the transaction from the msg
func (ctx CoreContext) SignBuildBroadcast(
	msg sdk.Msg, cdc *wire.Codec) (*ctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := ctx.SignAndBuild(msg, cdc)
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
