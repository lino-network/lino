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
	FlagPrivKey   = "priv_key"
	FlagPubKey    = "pub_key"

	// Account
	FlagIsFollow     = "is_follow"
	FlagFollowee     = "followee"
	FlagFollower     = "follower"
	FlagSender       = "sender"
	FlagReceiverName = "receiver_name"
	FlagReceiverAddr = "receiver_addr"
	FlagAmount       = "amount"
	FlagMemo         = "memo"

	// Developer
	FlagDeveloper = "developer"
	FlagDeposit   = "deposit"
	FlagUser      = "user"
	FlagSeconds   = "seconds"

	// Infra
	FlagProvider = "provider"
	FlagUsage    = "usage"

	// Post
	FlagDonator                 = "donator"
	FlagLikeUser                = "likeUser"
	FlagWeight                  = "weight"
	FlagAuthor                  = "author"
	FlagPostID                  = "post_ID"
	FlagTitle                   = "title"
	FlagContent                 = "content"
	FlagParentAuthor            = "parent_author"
	FlagParentPostID            = "parent_post_ID"
	FlagSourceAuthor            = "source_author"
	FlagSourcePostID            = "source_post_ID"
	FlagRedistributionSplitRate = "redistribution_split_rate"
	FlagFromChecking            = "from_checking"

	// Vote
	FlagVoter      = "voter"
	FlagProposalID = "proposal_id"
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
		c.Flags().String(FlagNode, "tcp://localhost:46657", "<host>:<port> to tendermint rpc interface for this chain")
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
