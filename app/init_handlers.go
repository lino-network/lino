package app

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *LinocoinApp) initHandlers() {
	app.initDefaultAnteHandler()
	app.initRouterHandlers()
}

func (app *LinocoinApp) initDefaultAnteHandler() {

	// Deducts fee from payer.
	// Verifies signatures and nonces.
	// Sets Signers to ctx.
	app.BaseApp.SetDefaultAnteHandler(
		auth.NewAnteHandler(app.accountMapper))
}

func (app *LinocoinApp) initRouterHandlers() {

	// All handlers must be added here.
	// The order matters.
	app.router.AddRoute("bank", bank.NewHandler(app.accountMapper))
}
