package client

import "github.com/spf13/cobra"

// nolint
const (
	FlagChainID   = "chain-id"
	FlagNode      = "node"
	FlagHeight    = "height"
	FlagTrustNode = "trust-node"
	FlagName      = "name"
	FlagSequence  = "sequence"
	FlagFee       = "fee"
	FlagPrivKey   = "priv-key"
	FlagPubKey    = "pub-key"

	// Account
	FlagIsFollow = "is-follow"
	FlagFollowee = "followee"
	FlagFollower = "follower"
	FlagSender   = "sender"
	FlagReceiver = "receiver"
	FlagAmount   = "amount"
	FlagMemo     = "memo"

	// Developer
	FlagDeveloper = "developer"
	FlagDeposit   = "deposit"
	FlagUser      = "user"
	FlagReferrer  = "referrer"
	FlagSeconds   = "seconds"

	// Infra
	FlagProvider = "provider"
	FlagUsage    = "usage"

	// Post
	FlagDonator                 = "donator"
	FlagLikeUser                = "likeUser"
	FlagWeight                  = "weight"
	FlagAuthor                  = "author"
	FlagPostID                  = "post-ID"
	FlagTitle                   = "title"
	FlagContent                 = "content"
	FlagParentAuthor            = "parent-author"
	FlagParentPostID            = "parent-post-ID"
	FlagSourceAuthor            = "source-author"
	FlagSourcePostID            = "source-post-ID"
	FlagRedistributionSplitRate = "redistribution-split-rate"
	FlagIsMicropayment          = "is-micropayment"

	// Vote
	FlagVoter      = "voter"
	FlagProposalID = "proposal-id"
	FlagResult     = "result"
	FlagLink       = "link"
)

// LineBreak can be included in a command list to provide a blank line
// to help with readability
var LineBreak = &cobra.Command{Run: func(*cobra.Command, []string) {}}

// GetCommands adds common flags to query commands
func GetCommands(cmds ...*cobra.Command) []*cobra.Command {
	for _, c := range cmds {
		// TODO: make this default false when we support proofs
		c.Flags().Bool(FlagTrustNode, true, "Don't verify proofs for responses")
		c.Flags().String(FlagChainID, "", "Chain ID of tendermint node")
		c.Flags().String(FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
		c.Flags().Int64(FlagHeight, 0, "block height to query, omit to get most recent provable block")
	}
	return cmds
}

// PostCommands adds common flags for commands to post tx
func PostCommands(cmds ...*cobra.Command) []*cobra.Command {
	for _, c := range cmds {
		c.Flags().Int64(FlagSequence, 0, "Sequence number to sign the tx")
		c.Flags().String(FlagChainID, "", "Chain ID of tendermint node")
		c.Flags().String(FlagPrivKey, "", "Private key to sign the transaction")
		c.Flags().String(FlagNode, "tcp://localhost:46657", "<host>:<port> to tendermint rpc interface for this chain")
	}
	return cmds
}
