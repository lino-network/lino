package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-crypto"
)

// AccountKey key format in KVStore
type AccountKey []byte

// PostKey key format in KVStore
type PostKey []byte

// Height is identity for each block
type Height uint64
