package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
)

type testLinocoinApp struct {
	*LinocoinApp
	*bam.TestApp
}

func newTestLinocoinApp() *testLinocoinApp {
	app := NewLinocoinApp()
	tba := &testLinocoinApp{
		LinocoinApp: app,
	}
	tba.TestApp = bam.NewTestApp(app.BaseApp)
	return tba
}
