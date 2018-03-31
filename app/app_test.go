package app

import (
	"os"
	"testing"

	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
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
