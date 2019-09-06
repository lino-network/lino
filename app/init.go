package app

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	genutil "github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	// sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagOverwrite = "overwrite"
)

// InitCmd initializes all files for tendermint and application
// XXX(yumin): after upgrade-1, we deprecated previous init function and start to use
// cosmos gaia init.
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(tmcli.HomeFlag))

			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			// gen pubkey
			_, pk, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()
			if !viper.GetBool(flagOverwrite) && common.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			// XXX(yumin): generate genesis file from app state.
			appGenTx, _, validator, err := LinoBlockchainGenTx(cdc, pk)
			if err != nil {
				return err
			}

			appState, err := LinoBlockchainGenState(cdc, []json.RawMessage{appGenTx})
			if err != nil {
				return err
			}

			// TODO(yumin): this is broken and we need update this part.
			if err = genutil.ExportGenesisFile(
				&tmtypes.GenesisDoc{
					GenesisTime:     time.Now(),
					ChainID:         chainID,
					ConsensusParams: nil,
					Validators:      []tmtypes.GenesisValidator{validator},
					AppHash:         []byte(""),
					AppState:        appState,
				},
				genFile); err != nil {
				return err
			}

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			fmt.Printf("Initialized lino configuration and bootstrapping files in %s...\n", viper.GetString(cli.HomeFlag))
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")

	return cmd
}
