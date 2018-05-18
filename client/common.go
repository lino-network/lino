package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/lino-network/lino/client/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-crypto"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func NewCoreContextFromViper() core.CoreContext {
	nodeURI := viper.GetString(FlagNode)
	var rpc rpcclient.Client
	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}
	var privKey crypto.PrivKey
	privKeyStr := viper.GetString(FlagPrivKey)
	if privKeyStr != "" {
		privKeyBytes, _ := hex.DecodeString(viper.GetString(FlagPrivKey))
		privKey, _ = crypto.PrivKeyFromBytes(privKeyBytes)
	}

	return core.CoreContext{
		ChainID:         viper.GetString(FlagChainID),
		Height:          viper.GetInt64(FlagHeight),
		TrustNode:       viper.GetBool(FlagTrustNode),
		FromAddressName: viper.GetString(FlagName),
		NodeURI:         nodeURI,
		Sequence:        viper.GetInt64(FlagSequence),
		Client:          rpc,
		PrivKey:         privKey,
	}
}

type CommandTxCallback func(cmd *cobra.Command, args []string) error

func PrintIndent(inputs ...interface{}) error {
	for _, input := range inputs {
		output, err := json.MarshalIndent(input, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	}
	return nil
}
