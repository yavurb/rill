package main

import (
	"context"
	"fmt"

	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/app"
)

var Environments = []string{"local", "staging", "production"}

func main() {
	cfg := loadConfig()
	app := app.NewApp(cfg)
	httpServer := app.NewHttpRouter()

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

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	httpServer.Logger.Fatal(httpServer.Start(addr))
}

// loadConfig attempts to load the configuration from a list of predefined environments.
// It iterates through the environments and tries to load the config file from each one.
// If a config is successfully loaded, it returns the config. If none of the configs can be loaded,
// it panics with the last encountered error.
func loadConfig() *config.Config {
	var appConfig *config.Config
	var err error

	for _, env := range Environments {
		fmt.Printf("Trying to load config from %s ⚙️\n", env)
		confPath := fmt.Sprintf("config/%s/config.pkl", env)

		appConfig, err = config.LoadFromPath(context.Background(), confPath)
		if err != nil {
			fmt.Printf("Failed to load config from %s\n", confPath)
			fmt.Println(err)
		} else if appConfig != nil {
			fmt.Printf("Loaded config from %s ✔︎\n", env)
			break
		}
	}

	if appConfig == nil {
		panic(err)
	}

	return appConfig
}
