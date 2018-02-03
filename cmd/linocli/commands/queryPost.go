package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wire "github.com/tendermint/go-wire"
	lcmd "github.com/tendermint/light-client/commands"
	proofcmd "github.com/tendermint/light-client/commands/proofs"
	"github.com/tendermint/light-client/proofs"

	btypes "github.com/lino-network/lino/types"
)

//nolint
const (
	FlagAddress       = "address"
)

var PostQueryCmd = &cobra.Command{
	Use:   "post",
	Short: "Get specific post, with proof",
	RunE:  lcmd.RequireInit(doPostQuery),
}

func init() {
	flags := PostQueryCmd.Flags()
	flags.String(FlagAddress, "", "Destination address for the query")
	flags.Int(FlagSequence, -1, "Sequence number for post")
}

func doPostQuery(cmd *cobra.Command, args []string) error {
	addr, err := proofs.ParseHexKey(viper.GetString(FlagAddress))
	if err != nil {
		return err
	}
	fmt.Println(string(addr)+"#"+viper.GetString(FlagSequence))
	key := wire.BinaryBytes(string(addr)+"#"+string(viper.GetInt(FlagSequence)))

	post := new(btypes.Post)
	proof, err := proofcmd.GetAndParseAppProof(key, &post)
	if err != nil {
		return err
	}

	return proofcmd.OutputProof(post, proof.BlockHeight())
}