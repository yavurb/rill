package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/app/mods"
	"github.com/yavurb/rill/internal/broadcasts/application"
	"github.com/yavurb/rill/internal/broadcasts/application/webrtc"
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
	broadcastConnectionUsecase := webrtc.NewBroadcastConnectionUsecase(app.config, e.Logger)
	viewerConnectionUsecase := webrtc.NewViewerConnectionUsecase(app.config, e.Logger)
	broadcastUsecaseParams := application.BroadcastUsecaseParams{
		Repository:       broadcastsRepository,
		Config:           app.config,
		BroadcastUsecase: broadcastConnectionUsecase,
		ViewerUsecase:    viewerConnectionUsecase,
		Logger:           e.Logger,
	}
	broadcastsUsecase := application.NewBroadcastUsecase(broadcastUsecaseParams)
	ui.NewBroadcastsRouter(e, broadcastsUsecase)

	return e
}
