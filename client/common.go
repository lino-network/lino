package client

import (
	"encoding/hex"
	"fmt"
	"os"

	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/lino-network/lino/client/core"
	"github.com/lino-network/lino/client/encrypt"
)

var ValidateCmd = cosmoscli.ValidateCmd

func ParsePubKey(key string) (crypto.PubKey, error) {
	pubKeyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return cryptoAmino.PubKeyFromBytes(pubKeyBytes)
}

func ParsePrivKey(key string) (crypto.PrivKey, error) {
	// @ tag means that priv-key is encrypted in the file.
	if key[0] == '@' {
		bytes, err := encrypt.DecryptByStdin(key[1:])
		if err != nil {
			exitWith("Failed to decrypt file: %s", err)
		}
		key = string(bytes)
	}

	var privKey crypto.PrivKey
	privKeyBytes, err := hex.DecodeString(key)
	if err != nil {
		return privKey, err
	}
	privKey, err = cryptoAmino.PrivKeyFromBytes(privKeyBytes)
	if err != nil {
		return privKey, err
	}
	return privKey, nil
}

func NewCoreContextFromViper() core.CoreContext {
	nodeURI := viper.GetString(FlagNode)
	var rpc rpcclient.Client
	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}

	seq := viper.GetInt64(FlagSequence)
	if seq < 0 {
		exitWith("Missing --" + FlagSequence)
	}

	ctx := core.CoreContext{
		ChainID:         viper.GetString(FlagChainID),
		Height:          viper.GetInt64(FlagHeight),
		TrustNode:       viper.GetBool(FlagTrustNode),
		FromAddressName: viper.GetString(FlagName),
		Offline:         viper.GetBool(FlagOffline),
		NodeURI:         nodeURI,
		Sequence:        uint64(viper.GetInt64(FlagSequence)), // XXX(yumin): dangerous, but ok.
		Client:          rpc,
	}
	ctx = ctx.WithFees(viper.GetString(FlagFees))

	privKey := viper.GetString(FlagPrivKey)
	if len(privKey) == 0 {
		exitWith("Missing --" + FlagPrivKey)
	}

	pk, err := ParsePrivKey(privKey)
	if err != nil {
		exitWith("Invalid PrivKey: %s", err)
	}
	ctx = ctx.WithPrivKey(pk)
	return ctx
}

func NewCoreBroadcastContextFromViper() core.CoreContext {
	nodeURI := viper.GetString(FlagNode)
	var rpc rpcclient.Client
	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}

	ctx := core.CoreContext{
		ChainID:         viper.GetString(FlagChainID),
		Height:          viper.GetInt64(FlagHeight),
		TrustNode:       viper.GetBool(FlagTrustNode),
		FromAddressName: viper.GetString(FlagName),
		Offline:         viper.GetBool(FlagOffline),
		NodeURI:         nodeURI,
		Client:          rpc,
	}
	return ctx
}

type CommandTxCallback func(cmd *cobra.Command, args []string) error

func exitWith(s string, args ...interface{}) {
	fmt.Printf(s+"\n", args...)
	os.Exit(1)
}
