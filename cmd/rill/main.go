package main

import (
	"github.com/yavurb/rill/internal/app"
)

func main() {
	appCtx := app.NewAppContext()
	app := appCtx.NewHttpRouter()

	app.Logger.Fatal(app.Start(":8910"))
}
