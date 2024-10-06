package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/app/mods"
	"github.com/yavurb/rill/internal/broadcasts/application"
	"github.com/yavurb/rill/internal/broadcasts/infrastructure/repository"
	"github.com/yavurb/rill/internal/broadcasts/infrastructure/ui"
)

type App struct {
	config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
	}
}

func (app *App) NewHttpRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Validator = mods.NewAppValidator()

	e.Use(middleware.Logger(), middleware.RequestID(), middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: app.config.Cors.AllowOrigins,
		AllowMethods: app.config.Cors.AllowMethods,
	}))

	e.Logger.SetLevel(log.DEBUG)

	broadcastsRepository := repository.NewLocalRepository()
	broadcastsUsecase := application.NewBroadcastUsecase(broadcastsRepository, app.config)
	ui.NewBroadcastsRouter(e, broadcastsUsecase)

	return e
}
