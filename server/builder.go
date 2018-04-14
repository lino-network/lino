package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	cmn "github.com/tendermint/tmlibs/common"
)

func QueryLocalStorage(key cmn.HexBytes, storeName string) (res []byte, err error) {
	node := rpcclient.NewHTTP(viper.GetString(FlagNode), "/websocket")
	path := fmt.Sprintf("/%s/key", types.ValidatorKVStoreKey)
	opts := rpcclient.ABCIQueryOptions{
		Height:  0,
		Trusted: true,
	}
	result, err := node.ABCIQueryWithOptions(path, model.GetValidatorListKey(), opts)
	if err != nil {
		return res, err
	}
	resp := result.Response
	if resp.Code != uint32(0) {
		return res, errors.Errorf("Query failed: (%d) %s", resp.Code, resp.Log)
	}
	return resp.Value, nil
}

// Broadcast the transaction bytes to Tendermint
func BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node := rpcclient.NewHTTP(viper.GetString(FlagNode), "/websocket")

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

// sign and build the transaction from the msg
func SignAndBuild(
	name, passphrase string, sequence int64, msg sdk.Msg, cdc *wire.Codec) ([]byte, error) {

	// build the Sign Messsage from the Standard Message
	chainID := viper.GetString(FlagChainID)
	signMsg := sdk.StdSignMsg{
		ChainID:   chainID,
		Sequences: []int64{sequence},
		Msg:       msg,
	}

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	// sign and build
	bz := signMsg.Bytes()

	sig, pubkey, err := keybase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}
	sigs := []sdk.StdSignature{{
		PubKey:    pubkey,
		Signature: sig,
		Sequence:  sequence,
	}}

	// marshal bytes
	tx := sdk.NewStdTx(signMsg.Msg, signMsg.Fee, sigs)

	return cdc.MarshalBinary(tx)
}

// sign and build the transaction from the msg
func SignBuildBroadcast(name string, passphrase string, sequence int64, msg sdk.Msg, cdc *wire.Codec) (*ctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := SignAndBuild(name, passphrase, sequence, msg, cdc)
	if err != nil {
		return nil, err
	}

	return BroadcastTx(txBytes)
}
