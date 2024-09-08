package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yavurb/rill/internal/app/mods"
	"github.com/yavurb/rill/internal/broadcasts/application"
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/broadcasts/infrastructure/repository"
	"github.com/yavurb/rill/internal/broadcasts/infrastructure/ui"
)

type AppCtx struct {
	Broadcasts []*domain.BroadcastSession
}

func NewAppContext() *AppCtx {
	return &AppCtx{}
}

func (appCtx *AppCtx) NewHttpRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Validator = mods.NewAppValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4321"},
	}))

	e.Logger.SetLevel(log.DEBUG)

	broadcastsRepository := repository.NewLocalRepository(appCtx.Broadcasts)
	broadcastsUsecase := application.NewBroadcastUsecase(broadcastsRepository)
	ui.NewBroadcastsRouter(e, broadcastsUsecase)

	return e
}
