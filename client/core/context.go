package core

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	linotypes "github.com/lino-network/lino/types"
)

// CoreContext - context used in terminal
type CoreContext struct {
	ChainID         string
	Height          int64
	TrustNode       bool
	NodeURI         string
	FromAddressName string
	Sequence        uint64
	Memo            string
	Client          rpcclient.Client
	PrivKeys        []crypto.PrivKey
	Fees            sdk.Coins
	TxEncoder       sdk.TxEncoder
}

// WithChainID - mount chain id on context
func (c CoreContext) WithChainID(chainID string) CoreContext {
	c.ChainID = chainID
	return c
}

// WithHeight - mount height on context
func (c CoreContext) WithHeight(height int64) CoreContext {
	c.Height = height
	return c
}

// WithTrustNode - mount trust node on context
func (c CoreContext) WithTrustNode(trustNode bool) CoreContext {
	c.TrustNode = trustNode
	return c
}

// WithNodeURI - mount node uri on context
func (c CoreContext) WithNodeURI(nodeURI string) CoreContext {
	c.NodeURI = nodeURI
	return c
}

// WithFromAddressName - mount from address on context
func (c CoreContext) WithFromAddressName(fromAddressName string) CoreContext {
	c.FromAddressName = fromAddressName
	return c
}

// WithSequence - mount sequence number on context
func (c CoreContext) WithSequence(sequence uint64) CoreContext {
	c.Sequence = sequence
	return c
}

// WithClient - mount client on context
func (c CoreContext) WithClient(client rpcclient.Client) CoreContext {
	c.Client = client
	return c
}

// WithPrivKey - mount private key on context
func (c CoreContext) WithPrivKey(privKey crypto.PrivKey) CoreContext {
	c.PrivKeys = append(c.PrivKeys, privKey)
	return c
}

// WithFees - mount fees
func (c CoreContext) WithFees(fees string) CoreContext {
	parsedFees, err := sdk.ParseCoins(fees)
	if err != nil {
		panic(err)
	}
	if parsedFees.Len() != 1 || parsedFees[0].Denom != linotypes.LinoCoinDenom {
		panic("invalid tx fees, unit must be: " + linotypes.LinoCoinDenom)
	}
	c.Fees = parsedFees
	return c
}

// WithCodec - mount cdc on context.
func (c CoreContext) WithTxEncoder(encoder sdk.TxEncoder) CoreContext {
	c.TxEncoder = encoder
	return c
}
