package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type AppCtx struct {
	Broadcasts []*domain.BroadcastSession
}

func NewAppContext() *AppCtx {
	return &AppCtx{}
}

func (appCtx *AppCtx) NewHttpRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	return e
}
