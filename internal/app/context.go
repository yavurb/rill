package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yavurb/rill/internal/app/mods"
	"github.com/yavurb/rill/internal/broadcasts/application"
	"github.com/yavurb/rill/internal/broadcasts/infrastructure/repository"
	"github.com/yavurb/rill/internal/broadcasts/infrastructure/ui"
)

type AppCtx struct{}

func NewAppContext() *AppCtx {
	return &AppCtx{}
}

func (appCtx *AppCtx) NewHttpRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Validator = mods.NewAppValidator()

	e.Use(middleware.Logger(), middleware.RequestID(), middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:4321",
			"https://rill.one",
			"http://rill.one",
			"https://rill.lat",
			"http://rill.lat",
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	e.Logger.SetLevel(log.DEBUG)

	broadcastsRepository := repository.NewLocalRepository()
	broadcastsUsecase := application.NewBroadcastUsecase(broadcastsRepository)
	ui.NewBroadcastsRouter(e, broadcastsUsecase)

	return e
}
