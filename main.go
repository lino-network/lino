package main

import (
	"github.com/lino-network/lino/app"
)

func main() {
	// Create BaseApp.
	var linoAppPtr = app.NewLinocoinApp()

	linoAppPtr.RunForever()
	return
}