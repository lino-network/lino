package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *LinocoinApp) initBaseApp() {
	bapp := baseapp.NewBaseApp(appName)
	app.BaseApp = bapp
	app.router = bapp.Router()
	app.initBaseAppTxDecoder()
}

func (app *LinocoinApp) initBaseAppTxDecoder() {
	cdc := makeTxCodec()
	app.BaseApp.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = sdk.StdTx{}
		// StdTx.Msg is an interface whose concrete
		// types are registered in app/msgs.go.
		err := cdc.UnmarshalBinary(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxParse("").TraceCause(err, "")
		}
		return tx, nil
	})
}
