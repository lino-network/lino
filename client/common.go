package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/lino-network/lino/client/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	crypto "github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
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
		privKey, _ = cryptoAmino.PrivKeyFromBytes(privKeyBytes)
	}

	if (viper.GetInt64(FlagSequence) < 0) {
		panic("Error on Sequence < 0, Sequence = " + fmt.Sprintf("%d", viper.GetInt64(FlagSequence)))
	}

	return core.CoreContext{
		ChainID:         viper.GetString(FlagChainID),
		Height:          viper.GetInt64(FlagHeight),
		TrustNode:       viper.GetBool(FlagTrustNode),
		FromAddressName: viper.GetString(FlagName),
		NodeURI:         nodeURI,
		Sequence:        uint64(viper.GetInt64(FlagSequence)), // XXX(yumin): dangerous, but ok.
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
