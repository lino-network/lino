package client

import (
	"github.com/spf13/cobra"
)

// nolint
const (
	FlagOffline   = "offline"
	FlagChainID   = "chain-id"
	FlagNode      = "node"
	FlagHeight    = "height"
	FlagTrustNode = "trust-node"
	FlagName      = "name"
	FlagSequence  = "sequence"
	FlagPrivKey   = "priv-key"
	FlagPubKey    = "pub-key"
	FlagFees      = "fees"

	// Infra
	FlagProvider = "provider"
	FlagUsage    = "usage"
)

// PostCommands adds common flags for commands to post tx
func PostCommands(cmds ...*cobra.Command) []*cobra.Command {
	for _, c := range cmds {
		c.Flags().Int64(FlagSequence, -1, "Sequence number to sign the tx")
		c.Flags().String(FlagChainID, "", "Chain ID of tendermint node")
		c.Flags().String(FlagPrivKey, "", "Private key to sign the transaction")
		c.Flags().String(FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
		c.Flags().String(FlagFees, "", "Fees to pay along with transaction; eg: 1linocoin")
		c.Flags().Bool(FlagOffline, false, "Print Tx to stdout only")
	}
	return cmds
}
