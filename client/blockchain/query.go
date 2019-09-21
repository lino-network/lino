package blockchain

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	// "github.com/lino-network/lino/utils"

	linotypes "github.com/lino-network/lino/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "blockchain",
		Short:                      "Blockchain-related Queries",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		// txutils.QueryTxCmd(cdc),
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
	)
	cmd.AddCommand(client.GetCommands(
		getCmdTx(cdc),
		getCmdHeight(cdc),
		getCmdMessage(cdc),
	)...)
	return cmd
}

// GetCmdHeight -
func getCmdHeight(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "height",
		Short: "height current block height",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			height, err := rpc.GetChainHeight(cliCtx)
			if err != nil {
				return err
			}
			fmt.Println(height)
			return nil
		},
	}
}

type MsgPrintLayout struct {
	Hash string    `json:"hash"`
	Msgs []sdk.Msg `json:"msgs"`
}

// GetCmdBlock -
func getCmdMessage(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "msg",
		Short: "msg <block-height>, print messages and results of the block height",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			var height *int64
			// optional height
			if len(args) > 0 {
				h, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				if h > 0 {
					tmp := int64(h)
					height = &tmp
				}
			}

			// get the node
			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			// header -> BlockchainInfo
			// header, tx -> Block
			// results -> BlockResults
			res, err := node.Block(height)
			if err != nil {
				return err
			}

			if !cliCtx.TrustNode {
				check, err := cliCtx.Verify(res.Block.Height)
				if err != nil {
					return err
				}

				err = tmliteProxy.ValidateBlockMeta(res.BlockMeta, check)
				if err != nil {
					return err
				}

				err = tmliteProxy.ValidateBlock(res.Block, check)
				if err != nil {
					return err
				}
			}

			decoder := linotypes.TxDecoder(cdc)
			var result []MsgPrintLayout
			for _, txbytes := range res.Block.Data.Txs {
				hexstr := hex.EncodeToString(txbytes.Hash())
				hexstr = strings.ToUpper(hexstr)
				tx, err := decoder(txbytes)
				if err != nil {
					return err
				}
				result = append(result, MsgPrintLayout{
					Hash: hexstr,
					Msgs: tx.GetMsgs(),
				})
			}
			out, err := cdc.MarshalJSONIndent(result, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))

			return nil
		},
	}
}

// GetCmdTx -
func getCmdTx(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tx",
		Short: "tx <hash>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			hashHexStr := args[0]
			hash, err := hex.DecodeString(hashHexStr)
			if err != nil {
				return err
			}

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			resTx, err := node.Tx(hash, !cliCtx.TrustNode)
			if err != nil {
				return err
			}

			if !cliCtx.TrustNode {
				if err = txutils.ValidateTxResult(cliCtx, resTx); err != nil {
					return err
				}
			}

			printTx(cdc, resTx.Tx)

			resBlocks, err := getBlocksForTxResults(cliCtx, []*ctypes.ResultTx{resTx})
			if err != nil {
				return err
			}

			out, err := formatTxResult(cliCtx.Codec, resTx, resBlocks[resTx.Height])
			if err != nil {
				return err
			}

			fmt.Println(out)
			return nil
		},
	}
}

func getBlocksForTxResults(cliCtx context.CLIContext, resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resBlocks := make(map[int64]*ctypes.ResultBlock)

	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}

func formatTxResult(cdc *codec.Codec, resTx *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (sdk.TxResponse, error) {
	tx, err := parseTx(linotypes.TxDecoder(cdc), resTx.Tx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return sdk.NewResponseResultTx(resTx, tx, resBlock.Block.Time.Format(time.RFC3339)), nil
}

func parseTx(decoder sdk.TxDecoder, txBytes []byte) (sdk.Tx, error) {
	tx, err := decoder(txBytes)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func printTx(cdc *codec.Codec, txBytes []byte) error {
	decoder := linotypes.TxDecoder(cdc)
	msg, err := decoder(txBytes)
	if err != nil {
		return err
	}
	out, err2 := cdc.MarshalJSONIndent(msg, "", "  ")
	if err2 != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
