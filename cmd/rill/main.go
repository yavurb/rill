package main

import (
	"fmt"

	"github.com/yavurb/rill/internal/app"
)

func main() {
	appCtx := app.NewAppContext()
	app := appCtx.NewHttpRouter()

	fmt.Print(`
 ██▀███   ██▓ ██▓     ██▓    
▓██ ▒ ██▒▓██▒▓██▒    ▓██▒    
▓██ ░▄█ ▒▒██▒▒██░    ▒██░    
▒██▀▀█▄  ░██░▒██░    ▒██░    
░██▓ ▒██▒░██░░██████▒░██████▒
░ ▒▓ ░▒▓░░▓  ░ ▒░▓  ░░ ▒░▓  ░
  ░▒ ░ ▒░ ▒ ░░ ░ ▒  ░░ ░ ▒  ░
  ░░   ░  ▒ ░  ░ ░     ░ ░   
   ░      ░      ░  ░    ░  ░
                             
`)

	app.Logger.Fatal(app.Start(":8910"))
}
