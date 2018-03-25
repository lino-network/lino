package app

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

// Construct some global addrs and txs for tests.
var (
	chainID = "" // TODO(cosmos)

	priv1 = crypto.GenPrivKeyEd25519()
	addr1 = priv1.PubKey().Address()
	addr2 = crypto.GenPrivKeyEd25519().PubKey().Address()
	coins = sdk.Coins{{"foocoin", 10}}
)

func newLinoBlockchain() *LinoBlockchain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewLinoBlockchain(logger, db)
}

//_______________________________________________________________________

func TestMsgs(t *testing.T) {
	// TODO(Lino)
}

func TestGenesis(t *testing.T) {
	// TODO(Lino)

}
