package client

import (
	"encoding/hex"

	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/lino-network/lino/client/core"
)

var ValidateCmd = cosmoscli.ValidateCmd

func parsePrivKey(key string) (crypto.PrivKey, error) {
	var privKey crypto.PrivKey
	privKeyBytes, err := hex.DecodeString(viper.GetString(FlagPrivKey))
	if err != nil {
		return privKey, err
	}
	privKey, _ = cryptoAmino.PrivKeyFromBytes(privKeyBytes)
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
		panic("Missing --" + FlagSequence)
	}

	ctx := core.CoreContext{
		ChainID:         viper.GetString(FlagChainID),
		Height:          viper.GetInt64(FlagHeight),
		TrustNode:       viper.GetBool(FlagTrustNode),
		FromAddressName: viper.GetString(FlagName),
		NodeURI:         nodeURI,
		Sequence:        uint64(viper.GetInt64(FlagSequence)), // XXX(yumin): dangerous, but ok.
		Client:          rpc,
	}
	ctx = ctx.WithFees(viper.GetString(FlagFees))

	hasKey := false
	for _, keyFlag := range []string{FlagPrivKey, FlagPrivKey2} {
		key := viper.GetString(keyFlag)
		if key != "" {
			pk, err := parsePrivKey(key)
			if err != nil {
				panic(err)
			}
			hasKey = true
			ctx = ctx.WithPrivKey(pk)
		}
	}
	if !hasKey {
		panic("Missing --" + FlagPrivKey)
	}
	return ctx
}

type CommandTxCallback func(cmd *cobra.Command, args []string) error
